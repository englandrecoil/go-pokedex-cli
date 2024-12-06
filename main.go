package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/englandrecoil/go-pokedex-cli/internal/pokeapi"
	"github.com/englandrecoil/go-pokedex-cli/internal/pokecache"
)

var cliName string = "pokedex "
var errUndefinedCommand error = errors.New("command not found")
var interval int = 10
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
	"bmap": {
		name:        "bmap",
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
}

type command struct {
	name        string
	description string
	callback    func(*pokeapi.Config, ...string) error
}

func printPrompt() {
	fmt.Print(cliName, "> ")
}

func commandHelp(cfg *pokeapi.Config, param ...string) error {
	if param == nil {
		fmt.Println("PARAM IS NIL")
	}
	fmt.Println("Usage:")
	fmt.Println("  help\t\t\t\tDisplays a help message")
	fmt.Println("  exit\t\t\t\tExit the Pokedex")
	fmt.Println("  clear\t\t\t\tClear the terminal screen")
	fmt.Println("  map\t\t\t\tDisplays the names of the next 20 location areas")
	fmt.Println("  bmap\t\t\t\tDisplays the names of the previous 20 location areas")
	fmt.Println("  explore {location_area}\tDisplays all the Pokémon in a given area")
	fmt.Println("  cache\t\t\t\tSet the caching interval after which cleaning will occur")
	fmt.Println("  \t\t\t\t(default value is 10 seconds)")
	fmt.Println("")

	return nil
}

func commandExit(cfg *pokeapi.Config, params ...string) error {
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
	locationAreas, err := pokeapi.GetLocationAreas(cfg, pokeapi.Next)
	if err != nil {
		return fmt.Errorf("map command error: %s", err)
	}

	for _, value := range locationAreas.Results {
		fmt.Println(value.Name)
	}

	return nil
}

func commandBackMap(cfg *pokeapi.Config, params ...string) error {
	locationAreas, err := pokeapi.GetLocationAreas(cfg, pokeapi.Previous)
	if err != nil {
		return fmt.Errorf("map command error: %s", err)
	}

	for _, value := range locationAreas.Results {
		fmt.Println(value.Name)
	}

	return nil
}

func commandCache(cfg *pokeapi.Config, params ...string) error {
	fmt.Println("Enter the caching interval")
	reader := bufio.NewScanner(os.Stdin)
	reader.Scan()

	var err error
	interval, err = strconv.Atoi(reader.Text())
	if err != nil {
		return errors.New("cache command error: interval must be a number")
	}

	fmt.Printf("%d seconds interval was set\n", interval)
	return nil
}

func commandExplore(cfg *pokeapi.Config, params ...string) error {
	if len(params) == 1 {
		return fmt.Errorf("explore command error: no location provided")
	}

	return nil
}

func defineCommand(input string, cfg *pokeapi.Config) error {
	cleanedInput := strings.Fields(input)

	if command, exists := commands[cleanedInput[0]]; exists {
		if err := command.callback(cfg, cleanedInput...); err != nil {
			return err
		}
		return nil
	}
	err := fmt.Errorf("%s: %w", input, errUndefinedCommand)
	return err
}

func printWelcomeMessage() {
	fmt.Printf("Welcome to the Pokedex!\n\n")
	fmt.Println("Please note that the Pokedex CLI is using a cache to quickly")
	fmt.Println("access data and reduce the load on the PokeAPI servers")
	fmt.Printf("\n")
	fmt.Println("Use 'help' command to find out about Pokémon world exploration commands.")
	fmt.Println()
}

func main() {
	cfg := &pokeapi.Config{}
	cfg.Cache = pokecache.NewCache(time.Duration(interval))
	reader := bufio.NewScanner(os.Stdin)

	printWelcomeMessage()
	printPrompt()

	for reader.Scan() {
		if err := defineCommand(reader.Text(), cfg); err != nil {
			if errors.Is(err, errUndefinedCommand) {
				fmt.Println(err)
				fmt.Println("Use 'help' to view the available commands")
			} else {
				fmt.Println(err)
			}
		}

		printPrompt()
	}
}
