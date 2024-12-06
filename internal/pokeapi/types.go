package pokeapi

import "github.com/englandrecoil/go-pokedex-cli/internal/pokecache"

type Config struct {
	NextURL     *string
	PreviousURL *string
	Cache       *pokecache.Cache
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
