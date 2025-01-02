package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"os/exec"
	"runtime"
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
		description: "Displays all caught Pokémon",
		callback:    commandPokedex,
	},
	"color": {
		name:        "color",
		description: "Configures the display of color output",
		callback:    commandColor,
	},
	"battle": {
		name:        "battle {pokemon_name} {pokemon_name}",
		description: "Simulate battles between captured Pokémon",
		callback:    commandBattle,
	},
}

func commandHelp(cfg *pokeapi.Config, param ...string) error {
	color.Set(color.FgYellow)
	defer color.Unset()

	fmt.Println("Usage:")
	fmt.Println("  pokedex\t\t\tDisplays all caught Pokémon")
	fmt.Println()
	fmt.Println("  map\t\t\t\tDisplays the names of the next 20 location areas")
	fmt.Println()
	fmt.Println("  mapb\t\t\t\tDisplays the names of the previous 20 location areas")
	fmt.Println()
	fmt.Println("  explore {location_area}\tDisplays all the Pokémon in a given area")
	fmt.Println()
	fmt.Println("  inspect {pokemon_name}\tInspect the caught Pokémon\t")
	fmt.Println()
	fmt.Println("  catch {pokemon_name}\t\tCatch Pokémon with a certain chance")
	fmt.Println()
	fmt.Println("  battle {pokemon_name1}\t\tSimulate battles between two captured Pokémon")
	fmt.Println("  {pokemon_name2}")
	fmt.Println()
	fmt.Println("  help\t\t\t\tDisplays a help message")
	fmt.Println()
	fmt.Println("  exit\t\t\t\tExit the Pokedex")
	fmt.Println()
	fmt.Println("  clear\t\t\t\tClear the terminal screen")
	fmt.Println()
	fmt.Println("  cache {integer_number}\tSet the caching interval(in hours) after which")
	fmt.Println("  \t\t\t\tcleaning will occur (default value is 1 hour)")
	fmt.Println()
	fmt.Println("  color {on/off}\t\tConfigures the display of color output. Only works")
	fmt.Println("  \t\t\t\tif the environment variable 'NO_COLOR' is empty")
	fmt.Println("  \t\t\t\t(default option is set to the NO_COLORS value)")
	fmt.Println()
	return nil
}

func commandExit(cfg *pokeapi.Config, params ...string) error {
	defer os.Exit(0)
	return nil
}

func commandClear(cfg *pokeapi.Config, params ...string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	case "darwin", "linux":
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("clear command error: %w", err)
	}

	return nil
}

func commandColor(cfg *pokeapi.Config, params ...string) error {
	if len(params) == 1 {
		return fmt.Errorf("color command error: no argument provided")
	}

	switch params[1] {
	case "off":
		color.NoColor = true
		return nil
	case "on":
		color.NoColor = false
		return nil
	default:
		return fmt.Errorf("color command error: unknown argument")
	}
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

	color.Set(color.FgBlue)
	defer color.Unset()

	fmt.Printf("Exploring %s...\n", params[1])
	fmt.Println("Found Pokemon:")
	color.Unset()

	for _, value := range location.PokemonEncounters {
		fmt.Printf(" - "+"%s"+"\n", value.Pokemon.Name)
	}
	return nil
}

func commandCatch(cfg *pokeapi.Config, params ...string) error {
	color.Set(color.FgBlue)
	defer color.Unset()

	if len(params) == 1 {
		return errors.New("catch command error: no Pokemon name provided")
	}

	if _, exists := cfg.PokemonCaught[params[1]]; exists {
		color.Set(color.FgRed)
		fmt.Printf("You already caught %s!\n", params[1])
		return nil
	}

	pokemon, err := pokeapi.GetPokemon(cfg, params[1])
	if err != nil {
		return fmt.Errorf("catch command error: %s", err)
	}

	_, err = pokeapi.GetImage(cfg, pokemon.Sprites.Other.OfficialArtwork.FrontDefault)
	if err != nil {
		return err
	}

	const treshold = 40
	chance := rand.IntN(pokemon.BaseExperience) + treshold

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)

	if pokemon.BaseExperience > chance {
		color.Set(color.FgRed)
		fmt.Printf("%s escaped!\n", pokemon.Name)
		delete(cfg.PokemonCaught, pokemon.Name)
		return nil
	}
	color.Set(color.FgGreen)
	fmt.Printf("%s was caught!\n", pokemon.Name)
	color.Set(color.FgBlue)
	fmt.Println("You may now inspect it with the 'inspect' command.")

	return nil
}

func commandInspect(cfg *pokeapi.Config, params ...string) error {
	color.Set(color.FgBlue)
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

	fmt.Println(color.BlueString("Name: ") + pokemon.Name)
	fmt.Println(color.BlueString("Height: ") + strconv.Itoa(pokemon.Height))
	fmt.Println(color.BlueString("Weight: ") + strconv.Itoa(pokemon.Weight))
	fmt.Println(color.BlueString("Stats: "))

	for _, value := range pokemon.Stats {
		fmt.Printf(" - "+color.BlueString("%s: ")+"%d\n", value.Stat.Name, value.BaseStat)
	}
	fmt.Println(color.BlueString("Types: "))
	for _, value := range pokemon.Types {
		fmt.Printf(" - "+"%s\n", value.Type.Name)
	}

	image, err := pokeapi.GetImage(cfg, pokemon.Sprites.Other.OfficialArtwork.FrontDefault)
	if err != nil {
		return err
	}

	fmt.Println(color.BlueString("Image: "))

	if err = pokedraw.DisplayImage(image); err != nil {
		return fmt.Errorf("display image error: %s", err)
	}
	fmt.Println()

	return nil
}

func commandPokedex(cfg *pokeapi.Config, params ...string) error {
	if len(cfg.PokemonCaught) == 0 {
		fmt.Println(color.BlueString("Your pokedex is empty! Try to catch Pokemon with 'catch' command"))
		return nil
	}

	fmt.Println(color.BlueString("Your pokedex:"))
	for _, value := range cfg.PokemonCaught {
		fmt.Println(" - " + value.Name)
	}

	return nil
}

func commandBattle(cfg *pokeapi.Config, params ...string) error {
	color.Set(color.FgBlue)
	defer color.Unset()

	if len(params) < 3 {
		return errors.New("battle command error: no Pokemon names provided")
	}
	if len(params) > 3 {
		return errors.New("battle command error: wrong number of arguments. Type `help` to to see available commands")
	}

	_, existsFirst := cfg.PokemonCaught[params[1]]
	_, existsSecond := cfg.PokemonCaught[params[2]]
	if !existsFirst {
		fmt.Printf("%s is not in your Pokedex!\n", params[1])
	}
	if !existsSecond {
		fmt.Printf("%s is not in your Pokedex!\n", params[2])
	}
	if !(existsFirst && existsSecond) {
		return nil
	}

	firstPokemon, err := pokeapi.GetPokemon(cfg, params[1])
	if err != nil {
		return fmt.Errorf("battle command error: %s", err)
	}
	secondPokemon, err := pokeapi.GetPokemon(cfg, params[2])
	if err != nil {
		return fmt.Errorf("battle command error: %s", err)
	}

	color.Set(color.FgBlue)
	defer color.Unset()

	fmt.Printf("The %s vs %s battle has begun\n", params[1], params[2])

	firstContestant := pokeapi.Battler{}
	firstContestant.Experience = firstPokemon.BaseExperience
	firstContestant.Name = firstPokemon.Name
	for _, value := range firstPokemon.Stats {
		switch value.Stat.Name {
		case "hp":
			firstContestant.Health = value.BaseStat
		case "attack":
			firstContestant.Attack = value.BaseStat
		case "defense":
			firstContestant.Defense = value.BaseStat
		case "special-defense":
			firstContestant.Parry = value.BaseStat
		default:
			continue
		}
	}

	secondContestant := pokeapi.Battler{}
	secondContestant.Experience = secondPokemon.BaseExperience
	secondContestant.Name = secondPokemon.Name
	for _, value := range secondPokemon.Stats {
		switch value.Stat.Name {
		case "hp":
			secondContestant.Health = value.BaseStat
		case "attack":
			secondContestant.Attack = value.BaseStat
		case "defense":
			secondContestant.Defense = value.BaseStat
		case "special-defense":
			secondContestant.Parry = value.BaseStat
		default:
			continue
		}
	}

	if err = startBattle(firstContestant, secondContestant); err != nil {
		return fmt.Errorf("battle command error: failed to start battle: %s", err)
	}

	return nil
}
