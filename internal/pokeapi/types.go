package pokeapi

type Config struct {
	NextURL     *string
	PreviousURL *string
}

type LocationAreasResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type Direction int

const (
	Next Direction = iota
	Previous
)
