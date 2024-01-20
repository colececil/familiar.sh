package commands

import (
	"fmt"
	"slices"
)

type CommandRegistry map[string]Command

// NewCommandRegistry returns a new instance of CommandRegistry. The registry is populated with the given Command
// instances.
func NewCommandRegistry(commands ...Command) CommandRegistry {
	commandRegistry := make(CommandRegistry)
	for _, command := range commands {
		commandRegistry[command.Name()] = command
	}
	commandRegistry.validate()
	return commandRegistry
}

// GetAllCommands returns a slice containing all commands.
func (r CommandRegistry) GetAllCommands() []Command {
	commands := make([]Command, 0, len(r))
	for _, command := range r {
		commands = append(commands, command)
	}
	slices.SortFunc(commands, func(command1, command2 Command) int {
		return command1.Order() - command2.Order()
	})
	return commands
}

// GetCommand returns the command with the given name, if it exists.
//
// It takes the following parameters:
//   - commandName: The name of the command.
func (r CommandRegistry) GetCommand(commandName string) (Command, error) {
	command, isPresent := r[commandName]
	if !isPresent {
		return nil, fmt.Errorf("command not valid")
	}

	return command, nil
}

// Validate checks if the CommandRegistry contains the expected commands, with no extras. It also checks that the
// commands' orders are valid. If the validation fails, it panics.
func (r CommandRegistry) validate() {
	panicMessage := "command registry does not contain the expected commands"
	expectedCommands := []string{"help", "version", "attune", "config", "package"}

	if len(r) != len(expectedCommands) {
		panic(panicMessage)
	}

	for _, expectedCommand := range expectedCommands {
		if r[expectedCommand] == nil {
			panic(panicMessage)
		}
	}

	orderNumbersEncountered := make(map[int]bool)
	for _, command := range r {
		order := command.Order()
		if order < 1 || order > len(r) || orderNumbersEncountered[order] {
			panic(panicMessage)
		}
		orderNumbersEncountered[order] = true
	}
}
