package api

type FlightpathsResponse []FlightPath

type FlightPath struct {
	Path          []string `json:"path"`
	Src           string   `json:"source"`
	Dest          string   `json:"destination"`
	DepartureTime int64    `json:"departureTime"`
}
