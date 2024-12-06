package pokeapi

import (
	"encoding/json"
	"fmt"
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
		if err := makeAPICall(url, &locationAreas); err != nil {
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

	if err := makeAPICall(url, &locationAreas); err != nil {
		return locationAreas, err
	}

	cfg.NextURL, cfg.PreviousURL = locationAreas.Next, locationAreas.Previous

	return locationAreas, nil
}

func makeAPICall[T any](url string, target *T) error {
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

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("error decoding response body: %s", err)
	}

	return nil
}
