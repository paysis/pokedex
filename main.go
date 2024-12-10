package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/paysis/pokedex/internal/pokedex"
)

type config struct {
	Previous string
	Next     string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

var commandRegistry = map[string]cliCommand{}

func main() {
	registerCommands()

	cfg := config{}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")

		// ignore scanner error, panic
		scanner.Scan()
		input := scanner.Text()
		input = strings.TrimSpace(input)
		input = strings.ToLower(input)

		if len(input) == 0 {
			continue
		}

		fields := strings.Fields(input)
		commandName := fields[0]

		command, ok := commandRegistry[commandName]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		err := command.callback(&cfg)
		if err != nil {
			_ = fmt.Errorf("something went wrong: %v", err)
		}
	}
}

func registerCommands() {
	// register commands
	commandRegistry["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}

	commandRegistry["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	}

	commandRegistry["map"] = cliCommand{
		name:        "map",
		description: "Display the next 20 location areas",
		callback:    commandMap,
	}

	commandRegistry["mapb"] = cliCommand{
		name:        "mapb",
		description: "Display the previous 20 location areas, if any",
		callback:    commandMapBack,
	}
}

func commandHelp(_ *config) error {
	_, err := fmt.Println("Welcome to the Pokedex!")
	if err != nil {
		return err
	}

	_, err = fmt.Print("Usage:\n\n")
	if err != nil {
		return err
	}

	for k, v := range commandRegistry {
		_, err := fmt.Printf("%s: %s\n", k, v.description)
		if err != nil {
			return err
		}
	}
	return nil
}

func commandExit(_ *config) error {
	if _, err := fmt.Println("Closing the Pokedex... Goodbye!"); err != nil {
		return err
	}
	os.Exit(0)
	return nil
}

func commandMap(cfg *config) error {
	var url string
	if cfg.Next != "" {
		url = cfg.Next
	}
	locationArea, err := pokedex.GetLocation(url)
	if err != nil {
		return err
	}

	cfg.Next = locationArea.Next
	cfg.Previous = locationArea.Previous

	for _, area := range locationArea.Results {
		fmt.Println(area.Name)
	}

	return nil
}

func commandMapBack(cfg *config) error {
	if cfg.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	url := cfg.Previous
	locationArea, err := pokedex.GetLocation(url)
	if err != nil {
		return err
	}

	cfg.Next = locationArea.Next
	cfg.Previous = locationArea.Previous

	for _, area := range locationArea.Results {
		fmt.Println(area.Name)
	}

	return nil
}
