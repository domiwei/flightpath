package api

import (
	"database/sql"
	"flightpath/service/localcache"
	"log"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

func NewHandler(r *gin.RouterGroup, db *sql.DB, lc localcache.Service) {
	h := handler{
		flightStore: NewFlightStore(db, lc),
	}
	r.POST("/:personID/flights", h.postFlights)
	r.GET("/:personID/flights", h.getFlights)
}

type handler struct {
	flightStore *flight
}

type postFlightParams struct {
	Path          [][]string `json:"path"`
	DepartureTime int64      `json:"departureTime"`
}

func (h handler) postFlights(c *gin.Context) {
	pid := c.Param("personID")
	var p postFlightParams
	if err := c.BindJSON(&p); err != nil {
		log.Printf("c.BindJSON failed. %v", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if err := h.flightStore.InsertFlights(pid, p.Path, p.DepartureTime); err != nil {
		log.Printf("insertFlight failed. %v", err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

type getFlightParams struct {
	DepartureGreaterThan int64 `json:"departureGreaterThan"`
}

func (h handler) getFlights(c *gin.Context) {
	pid := c.Param("personID")
	p := getFlightParams{}
	if err := c.BindQuery(&p); err != nil {
		log.Printf("c.BindQuery failed. %v", err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	paths, err := h.flightStore.GetFlights(pid)
	if err != nil {
		log.Printf("GetFlight failed. %v", err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// filter out those satisfying departureTime constraint
	filteredPaths := []FlightPath{}
	for _, path := range paths {
		if path.DepartureTime > p.DepartureGreaterThan {
			filteredPaths = append(filteredPaths, path)
		}
	}

	// Sort by departureTime in asceding order
	sort.Slice(filteredPaths, func(i, j int) bool {
		return filteredPaths[i].DepartureTime < filteredPaths[j].DepartureTime
	})

	c.JSON(http.StatusOK, filteredPaths)
}
