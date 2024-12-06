package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetLocationAreas(cfg *Config, direction Direction) (locationAreas LocationAreasResponse, err error) {
	url := "https://pokeapi.co/api/v2/location-area/"
	locationAreas = LocationAreasResponse{}

	// pagination(direction) rules for the very first ever request
	if cfg.NextURL == nil && cfg.PreviousURL == nil {
		if direction == Previous {
			return locationAreas, fmt.Errorf("no more locations")
		}
		if err := makeAPICall(url, &locationAreas, cfg); err != nil {
			return locationAreas, err
		}
		cfg.NextURL, cfg.PreviousURL = locationAreas.Next, locationAreas.Previous
		return locationAreas, nil
	}

	// pagination(direction) rules for the next requests
	switch direction {
	case Next:
		if cfg.NextURL == nil {
			return locationAreas, fmt.Errorf("no more locations")
		}
		url = *cfg.NextURL
	case Previous:
		if cfg.PreviousURL == nil {
			return locationAreas, fmt.Errorf("no more locations")
		}
		url = *cfg.PreviousURL
	}

	// check if the data is already in the cache
	if data, exists := cfg.Cache.Get(url); exists {
		if err = json.Unmarshal(data, &locationAreas); err != nil {
			return locationAreas, fmt.Errorf("error decoding cached data: %s", err)
		}
		return locationAreas, nil
	}

	if err := makeAPICall(url, &locationAreas, cfg); err != nil {
		return locationAreas, err
	}

	cfg.NextURL, cfg.PreviousURL = locationAreas.Next, locationAreas.Previous

	return locationAreas, nil
}

func makeAPICall[T any](url string, target *T, cfg *Config) error {
	// HTTP processing
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("can't initialize request for server: %s", err)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("can't send request to server: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("non-OK HTTP status: %s", res.Status)
	}

	// cache processing
	bodyData, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %s", err)
	}
	cfg.Cache.Add(url, bodyData)

	if err = json.Unmarshal(bodyData, target); err != nil {
		return fmt.Errorf("error decoding response body: %s", err)
	}

	return nil
}
