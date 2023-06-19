package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: No command specified.")
		os.Exit(1)
	}

	commandName := os.Args[1]

	commandRegistry := InitializeCommandRegistry()
	command, err := commandRegistry.GetCommand(commandName)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	if err := command.Execute(os.Args[2:]); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
}
