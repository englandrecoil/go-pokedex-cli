package pokesave

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/englandrecoil/go-pokedex-cli/internal/pokeapi"
)

const savePath = "./saves/"

func SaveProgress(cfg *pokeapi.Config) error {
	if err := os.MkdirAll(savePath, 0755); err != nil {
		return fmt.Errorf("save progress error: %w", err)
	}

	data, err := json.MarshalIndent(cfg.PokemonCaught, "", " ")
	if err != nil {
		return fmt.Errorf("save progress error: %w", err)
	}

	if err = os.WriteFile("./saves/pokedex.json", data, 0644); err != nil {
		return fmt.Errorf("save progress error: %w", err)
	}
	return nil
}

func LoadProgress(cfg *pokeapi.Config) error {
	pokedexData, err := os.ReadFile("./saves/pokedex.json")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
	}

	if err := json.Unmarshal(pokedexData, &cfg.PokemonCaught); err != nil {
		return fmt.Errorf("load progress error: %w", err)
	}
	return nil
}
