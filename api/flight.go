package api

import "database/sql"

type flight struct {
	db *sql.DB
}

func NewFlightStore(db *sql.DB) *flight {
	return &flight{
		db: db,
	}
}

func (f *flight) InsertFlights(personID string, path [][]string, departure int64) error {
	return nil
}

func (f *flight) GetFlights(personID string) (*FlightpathsResponse, error) {
	return nil, nil
}
