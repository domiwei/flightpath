package api

import (
	"database/sql"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var (
	ErrInvalidPath = errors.New("invalid path")
)

type flight struct {
	db *sql.DB
}

func NewFlightStore(db *sql.DB) *flight {
	return &flight{
		db: db,
	}
}

func (f *flight) InsertFlights(personID string, path [][]string, departure int64) error {
	route, err := resolvePath(path)
	if err != nil {
		return errors.Wrapf(err, "resolvePath failed")
	}
	insertQuery := "INSERT INTO `paths` (`person_id`, `path`, `departure_time`, `insert_time`) VALUES (?,?,?,?)"
	if _, err := f.db.Exec(insertQuery, personID, strings.Join(route, ","), departure, time.Now().Unix()); err != nil {
		return errors.Wrapf(err, "db.Exec failed")
	}
	return nil
}

func (f *flight) GetFlights(personID string) (*FlightpathsResponse, error) {
	// fetch from DB
	query := "SELECT `path`, `departure_time` FROM `paths` WHERE `person_id`=?"
	rows, err := f.db.Query(query, personID)
	if err != nil {
		return nil, errors.Wrapf(err, "db.Query failed")
	}
	defer rows.Close()

	resp := FlightpathsResponse{}
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
			Dest:          p[1],
			DepartureTime: departure,
		})
	}
	return &resp, nil
}

func resolvePath(p [][]string) ([]string, error) {
	// build directed graph
	g := map[string]string{}
	edgeDegree := map[string]int{}
	for _, edge := range p {
		f, t := edge[0], edge[1]
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
			return nil, ErrInvalidPath
		}
		switch degree {
		case -1:
			if start == "" {
				start = city
			} else {
				return nil, ErrInvalidPath
			}
		case 1:
			if end == "" {
				end = city
			} else {
				return nil, ErrInvalidPath
			}
		}
	}
	if start == "" || end == "" {
		return nil, ErrInvalidPath
	}

	// Finally rearrange them by going through the path from starting city
	route := []string{}
	cur := start
	for {
		route = append(route, cur)
		nxt, exist := g[cur]
		if !exist {
			break
		}
		cur = nxt
	}
	return route, nil
}
