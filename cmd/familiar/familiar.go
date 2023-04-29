package main

import (
	"fmt"
	"github.com/colececil/familiar.sh/internal/commands"
	"os"
)

// This will be overwritten during the build process.
var projectVersion = "0.0.0"

var versionCommand = &commands.VersionCommand{Version: projectVersion}
var attuneCommand = &commands.AttuneCommand{}
var configCommand = &commands.ConfigCommand{}
var packageCommand = &commands.PackageCommand{}
var helpCommand = &commands.HelpCommand{Commands: []commands.Command{
	versionCommand,
	attuneCommand,
	configCommand,
	packageCommand,
}}
var commandMap = map[string]commands.Command{
	helpCommand.Name():    helpCommand,
	versionCommand.Name(): versionCommand,
	attuneCommand.Name():  attuneCommand,
	configCommand.Name():  configCommand,
	packageCommand.Name(): packageCommand,
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
