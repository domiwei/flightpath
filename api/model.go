package api

type FlightpathsResponse []FlightPath

type FlightPath struct {
	Path          []string `json:"path"`
	Src           string   `json:"source"`
	Dest          string   `json:"destination"`
	DepartureTime int64    `json:"departureTime"`
}

type dbPath struct {
	Path          string `json:"path"`
	Person        string `json:"person_id"`
	DepartureTime int64  `json:"departureTime"`
	InsertTime    int64  `json:"insertTime"`
}
