package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/englandrecoil/go-pokedex-cli/internal/pokeapi"
	"github.com/englandrecoil/go-pokedex-cli/internal/pokecache"
	"github.com/fatih/color"
)

var cliName string = "pokedex "
var errUndefinedCommand error = errors.New("command not found")
var interval int = 1

func printPrompt() {
	fmt.Print(cliName, "> ")
}

func printWelcomeMessage() {
	color.Set(color.FgGreen)
	defer color.Unset()

	fmt.Printf("Welcome to the Pokedex!\n\n")
	fmt.Println("Please note that the Pokedex CLI is using a cache to quickly")
	fmt.Println("access data and reduce the load on the PokeAPI servers")
	fmt.Printf("\n")
	fmt.Println("Use 'help' command to find out about Pokémon world exploration commands.")
	fmt.Println()
}

func defineCommand(input string, cfg *pokeapi.Config) error {
	cleanedInput := strings.Fields(strings.ToLower(input))

	if command, exists := commands[cleanedInput[0]]; exists {
		if err := command.callback(cfg, cleanedInput...); err != nil {
			return err
		}
		return nil
	}
	err := fmt.Errorf("%s: %w", input, errUndefinedCommand)

	return err
}

func init() {
	if noColorVariable := os.Getenv("NO_COLOR"); noColorVariable == "" {
		color.NoColor = false
	} else {
		color.NoColor = true
	}
}

func main() {
	cfg := &pokeapi.Config{
		NextURL:       nil,
		PreviousURL:   nil,
		Cache:         pokecache.NewCache(time.Duration(interval) * time.Hour),
		PokemonCaught: make(map[string]pokeapi.Pokemon),
	}

	reader := bufio.NewScanner(os.Stdin)
	red := color.New(color.FgRed).PrintlnFunc()

	printWelcomeMessage()
	printPrompt()

	for reader.Scan() {
		if err := defineCommand(reader.Text(), cfg); err != nil {
			if errors.Is(err, errUndefinedCommand) {
				red(err)
				fmt.Println("Use 'help' to view the available commands")
			} else {
				red(err)
			}
		}

		printPrompt()
	}
}
