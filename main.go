package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/englandrecoil/go-pokedex-cli/internal/api"
)

var cliName string = "pokedex "

type command struct {
	name        string
	description string
	callback    func() error
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
}

func printPrompt() {
	fmt.Print(cliName, "> ")
}

func commandHelp() error {
	fmt.Printf("Welcome to the Pokedex!\n")
	fmt.Println("\nUsage:")
	fmt.Println("  help\t\tDisplays a help message")
	fmt.Println("  exit\t\tExit the Pokedex")
	fmt.Println("  clear\t\tClear the terminal screen")
	fmt.Println("  map\t\tDisplays the names of the next 20 location areas")
	fmt.Println("")
	return nil
}

func commandExit() error {
	defer os.Exit(0)
	return nil
}

func commandClear() error {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("clear command error: %w", err)
	}
	return nil
}

func commandMap() error {
	if err := api.PrintAreas(); err != nil {
		return fmt.Errorf("map command error: %s", err)
	}
	return nil
}

func defineCommand(input string) error {
	if input, exists := commands[input]; exists {
		if err := input.callback(); err != nil {
			return err
		}
		return nil
	}
	// err := fmt.Errorf("%s: command not found", input)
	err := fmt.Errorf("%s: %w", input, errUndefinedCommand)
	return err
}

var errUndefinedCommand error = errors.New("command not found")

func main() {
	reader := bufio.NewScanner(os.Stdin)
	printPrompt()
	for reader.Scan() {
		// define command. If exists - do something later
		if err := defineCommand(reader.Text()); err != nil {
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
