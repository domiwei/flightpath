package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewHandler(r *gin.RouterGroup, db *sql.DB) {
	h := handler{
		flightStore: NewFlightStore(db),
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
		c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := h.flightStore.InsertFlights(pid, p.Path, p.DepartureTime); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h handler) getFlights(c *gin.Context) {

}
