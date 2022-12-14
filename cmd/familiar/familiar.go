package main

import (
	"fmt"
	"github.com/colececil/familiar.sh/internal/commands"
	"os"
)

// This will be overwritten during the build process.
var projectVersion = "0.0.0"

var versionCommand = &commands.VersionCommand{Version: projectVersion}
var configCommand = &commands.ConfigCommand{}
var helpCommand = &commands.HelpCommand{Commands: []commands.Command{
	versionCommand,
	configCommand,
}}
var commandMap = map[string]commands.Command{
	"help":    helpCommand,
	"version": versionCommand,
	"config":  configCommand,
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: No command specified.")
		os.Exit(1)
	}

	commandName := os.Args[1]
	command, isPresent := commandMap[commandName]
	if !isPresent {
		fmt.Printf("Error: \"%s\" is not a valid command.\n", commandName)
		os.Exit(1)
	}

	if err := command.Execute(os.Args[2:]); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
}
