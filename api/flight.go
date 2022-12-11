package api

import (
	"database/sql"
	"flightpath/service/localcache"
	"log"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var (
	ErrInvalidPath        = errors.New("invalid path")
	ErrSourceAlrdyDefined = errors.New("source already defined")
	ErrCycleTour          = errors.New("cycle tour")
)

type flight struct {
	db *sql.DB
	lc localcache.Service
}

func NewFlightStore(db *sql.DB, lc localcache.Service) *flight {
	return &flight{
		db: db,
		lc: lc,
	}
}

func (f *flight) InsertFlights(personID string, path [][]string, departure int64) error {
	route, err := resolvePath(path)
	if err != nil {
		return errors.Wrapf(err, "resolvePath failed")
	}
	// upsert by (person, departure time)
	insertQuery := "INSERT INTO `paths` (`person_id`, `path`, `departure_time`, `insert_time`) VALUES (?,?,?,?)" +
		" ON DUPLICATE KEY UPDATE `person_id`=VALUES(person_id), `departure_time`=VALUES(departure_time)"
	if _, err := f.db.Exec(insertQuery,
		personID,
		strings.Join(route, ","),
		departure,
		time.Now().Unix()); err != nil {
		return errors.Wrapf(err, "db.Exec failed")
	}
	// Try delete cache data
	f.lc.Delete(personID)
	return nil
}

func (f *flight) GetFlights(personID string) ([]FlightPath, error) {
	// Fetch from cache first
	if val, ok := f.lc.Get(personID); ok {
		log.Printf("cache hit %v", personID)
		return val.([]FlightPath), nil
	}

	// fetch from DB
	query := "SELECT `path`, `departure_time` FROM `paths` WHERE `person_id`=?"
	rows, err := f.db.Query(query, personID)
	if err != nil {
		return nil, errors.Wrapf(err, "db.Query failed")
	}
	defer rows.Close()

	resp := []FlightPath{}
	for rows.Next() {
		var path string
		var departure int64
		if err := rows.Scan(&path, &departure); err != nil {
			return nil, errors.Wrapf(err, "rows.Scan failed")
		}
		p := strings.Split(path, ",")
		resp = append(resp, FlightPath{
			Path:          p,
			Src:           p[0],
			Dest:          p[len(p)-1],
			DepartureTime: departure,
		})
	}

	f.lc.SetDefault(personID, resp)
	return resp, nil
}

func resolvePath(p [][]string) ([]string, error) {
	// build directed graph
	g := map[string]string{}
	edgeDegree := map[string]int{}
	for _, edge := range p {
		f, t := edge[0], edge[1]
		if _, exist := g[f]; exist {
			return nil, errors.Wrapf(ErrSourceAlrdyDefined, "%v", edge)
		}
		g[f] = t
		// -1 for node with outgoing edge, and 1 for incoming
		edgeDegree[f] -= 1
		edgeDegree[t] += 1
	}

	// validate by checking the degree. Must have exact one node with degree = 1, which is starting city.
	// on the other side, the only node with final degree = 1 is end city. All the intermideate cities
	// should be 0 degree.
	var start, end string
	for city, degree := range edgeDegree {
		if degree < -1 || degree > 1 {
			// Impossible case
			return nil, ErrInvalidPath
		}
		switch degree {
		case -1:
			if start == "" {
				start = city
			} else {
				// multiple starting city
				return nil, ErrInvalidPath
			}
		case 1:
			if end == "" {
				end = city
			} else {
				// multiple ending city
				return nil, ErrInvalidPath
			}
		}
	}
	if start == "" || end == "" {
		// no starting/ending city. Cycle tour
		return nil, ErrInvalidPath
	}

	// Finally rearrange them by going through the path from starting city
	route := []string{}
	cur := start
	visit := map[string]struct{}{}
	for {
		if _, exist := visit[cur]; exist {
			return nil, ErrCycleTour
		}
		route = append(route, cur)
		visit[cur] = struct{}{}
		nxt, exist := g[cur]
		if !exist {
			break
		}
		cur = nxt
	}
	return route, nil
}
