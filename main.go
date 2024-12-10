package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")

		// ignore scanner error, panic
		scanner.Scan()
		input := scanner.Text()
		input = strings.TrimSpace(input)
		input = strings.ToLower(input)

		fields := strings.Fields(input)

		firstWord := fields[0]
		fmt.Printf("Your command was: %s\n", firstWord)
	}
}
