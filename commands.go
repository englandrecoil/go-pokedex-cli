package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"os/exec"
	"path"
	"strconv"

	"github.com/englandrecoil/go-pokedex-cli/internal/pokeapi"
	"github.com/englandrecoil/go-pokedex-cli/internal/pokedraw"
	"github.com/fatih/color"
)

type command struct {
	name        string
	description string
	callback    func(*pokeapi.Config, ...string) error
}

var commands = map[string]command{
	"help": {
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	},
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
	"clear": {
		name:        "clear",
		description: "Clear the terminal screen",
		callback:    commandClear,
	},
	"map": {
		name:        "map",
		description: "Displays the names of the next 20 location areas",
		callback:    commandMap,
	},
	"mapb": {
		name:        "mapb",
		description: "Displays the names of the previous 20 location areas",
		callback:    commandBackMap,
	},
	"explore": {
		name:        "explore {location}",
		description: "List all the Pokémon in a given area",
		callback:    commandExplore,
	},
	"cache": {
		name:        "cache",
		description: "Set the caching interval after which cleaning will occur",
		callback:    commandCache,
	},
	"catch": {
		name:        "catch",
		description: "Catch Pokemon with a certain chance",
		callback:    commandCatch,
	},
	"inspect": {
		name:        "inspect {pokemon_name}",
		description: "Inspect the caught pokemon",
		callback:    commandInspect,
	},
	"pokedex": {
		name:        "pokedex",
		description: "Displays all caught Pokemon",
		callback:    commandPokedex,
	},
}

func commandHelp(cfg *pokeapi.Config, param ...string) error {
	color.RGB(255, 179, 26).Set()
	defer color.Unset()

	fmt.Println("Usage:")
	fmt.Println("  help\t\t\t\tDisplays a help message")
	fmt.Println("  exit\t\t\t\tExit the Pokedex")
	fmt.Println("  clear\t\t\t\tClear the terminal screen")
	fmt.Println("  pokedex\t\t\tDisplays all caught Pokemon")
	fmt.Println("  map\t\t\t\tDisplays the names of the next 20 location areas")
	fmt.Println("  mapb\t\t\t\tDisplays the names of the previous 20 location areas")
	fmt.Println("  explore {location_area}\tDisplays all the Pokémon in a given area")
	fmt.Println("  inspect {pokemon_name}\tInspect the caught pokemon\t")
	fmt.Println("  catch {pokemon_name}\t\tCatch Pokemon with a certain chance")
	fmt.Println("  cache {integer_number}\tSet the caching interval(in hours) after which cleaning will occur")
	fmt.Println("  \t\t\t\t(default value is 1 hour)")
	fmt.Println()
	return nil
}

func commandExit(cfg *pokeapi.Config, params ...string) error {
	dir, err := os.ReadDir("./images")
	if err != nil {
		return fmt.Errorf("command exit error: can't clean image folder: %s", err)
	}

	for _, value := range dir {
		os.RemoveAll(path.Join([]string{"./images", value.Name()}...))
	}

	defer os.Exit(0)
	return nil
}

func commandClear(cfg *pokeapi.Config, params ...string) error {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("clear command error: %w", err)
	}
	return nil
}

func commandMap(cfg *pokeapi.Config, params ...string) error {
	locations, err := pokeapi.GetLocationAreas(cfg, pokeapi.Next)
	if err != nil {
		return fmt.Errorf("map command error: %s", err)
	}

	for _, value := range locations.Results {
		fmt.Println(" - " + value.Name)
	}

	return nil
}

func commandBackMap(cfg *pokeapi.Config, params ...string) error {
	locations, err := pokeapi.GetLocationAreas(cfg, pokeapi.Previous)
	if err != nil {
		return fmt.Errorf("map command error: %s", err)
	}

	for _, value := range locations.Results {
		fmt.Println(" - " + value.Name)
	}

	return nil
}

func commandCache(cfg *pokeapi.Config, params ...string) error {
	if len(params) == 1 {
		return errors.New("cache command error: no value provided")
	}
	inputInterval := params[1]

	var err error
	interval, err = strconv.Atoi(inputInterval)
	if err != nil {
		return errors.New("cache command error: interval must be an integer number")
	}
	if interval <= 0 {
		return errors.New("cache command error: the number must be greater than 0")
	}

	fmt.Printf("%d hour interval was set\n", interval)
	return nil
}

func commandExplore(cfg *pokeapi.Config, params ...string) error {
	if len(params) == 1 {
		return errors.New("explore command error: no location provided")
	}

	location, err := pokeapi.GetLocationArea(cfg, params[1])
	if err != nil {
		return fmt.Errorf("explore command error: %s", err)
	}

	color.RGB(102, 153, 255).Set()
	fmt.Printf("Exploring %s...\n", params[1])
	fmt.Println("Found Pokemon:")
	color.Unset()

	for _, value := range location.PokemonEncounters {
		fmt.Printf(" - " + value.Pokemon.Name + "\n")
	}
	return nil
}

func commandCatch(cfg *pokeapi.Config, params ...string) error {
	color.RGB(102, 153, 255).Set()
	defer color.Unset()

	if len(params) == 1 {
		return errors.New("catch command error: no Pokemon name provided")
	}

	if _, exists := cfg.PokemonCaught[params[1]]; exists {
		fmt.Printf("You already caught %s!\n", params[1])
		return nil
	}

	pokemon, err := pokeapi.GetPokemon(cfg, params[1])
	if err != nil {
		return fmt.Errorf("catch command error: %s", err)
	}

	const treshold = 50
	chance := rand.IntN(pokemon.BaseExperience) + treshold

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
	if pokemon.BaseExperience > chance {
		fmt.Printf("%s escaped!\n", pokemon.Name)
		delete(cfg.PokemonCaught, pokemon.Name)
		return nil
	}
	fmt.Printf("%s was caught!\n", pokemon.Name)
	fmt.Println("You may now inspect it with the 'inspect' command.")

	return nil
}

func commandInspect(cfg *pokeapi.Config, params ...string) error {
	color.RGB(102, 153, 255).Set()
	defer color.Unset()

	if len(params) == 1 {
		return errors.New("inspect command error: no Pokemon name provided")
	}

	if _, exists := cfg.PokemonCaught[params[1]]; !exists {
		fmt.Println("You have not caught that pokemon!")
		return nil
	}

	pokemon, err := pokeapi.GetPokemon(cfg, params[1])
	if err != nil {
		return fmt.Errorf("inspect command error: %s", err)
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")

	for _, value := range pokemon.Stats {
		fmt.Printf(" - "+"%s: "+"%d\n", value.Stat.Name, value.BaseStat)
	}
	fmt.Println("Types:")
	for _, value := range pokemon.Types {
		fmt.Printf(" - "+"%s\n", value.Type.Name)
	}

	// display image
	image, err := pokeapi.GetImage(cfg, pokemon.Sprites.Other.OfficialArtwork.FrontDefault)
	if err != nil {
		return err
	}

	if err = pokedraw.DisplayImage(pokemon.Name, image); err != nil {
		return fmt.Errorf("display image error: %s", err)
	}

	return nil
}

func commandPokedex(cfg *pokeapi.Config, params ...string) error {
	color.RGB(102, 153, 255).Set()
	defer color.Unset()

	if len(cfg.PokemonCaught) == 0 {
		fmt.Println("Your pokedex is empty! Try to catch Pokemon with 'catch' command")
		return nil
	}

	fmt.Println("Your pokedex:")
	for _, value := range cfg.PokemonCaught {
		fmt.Println(" - " + value.Name + "ID:" + strconv.Itoa(value.ID))
	}

	return nil
}
