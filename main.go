package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var commandRegistry = map[string]cliCommand{}

func main() {
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

		err := command.callback()
		if err != nil {
			_ = fmt.Errorf("something went wrong: %v", err)
		}
	}
}

func commandHelp() error {
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

func commandExit() error {
	if _, err := fmt.Println("Closing the Pokedex... Goodbye!"); err != nil {
		return err
	}
	os.Exit(0)
	return nil
}
