package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Johnnydeeps/pokedexcli/internal/pokecache"
)

func cleanInput(text string) []string {
	lowered := strings.ToLower(text)
	words := strings.Fields(lowered)
	return words
}

type config struct {
	nextLocationsURL *string
	prevLocationsURL *string
	cache            pokecache.Cache
}

func startRepl(cfg *config) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		ok := scanner.Scan()
		if !ok {
			if err := scanner.Err(); err != nil {
				return
			}
			return
		}

		text := scanner.Text()

		words := cleanInput(text)
		if len(words) == 0 {
			continue
		}

		cliInput := words[0]
		commands := getCommands()
		command, exists := commands[cliInput]
		if exists {
			err := command.callback(cfg)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message with valid commands for the pokedex",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Get the next page of locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Get the previous page of locations",
			callback:    commandMapb,
		},
	}
}
