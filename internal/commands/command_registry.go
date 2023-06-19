package commands

import "fmt"

type CommandRegistry map[string]Command

// NewCommandRegistry returns a new instance of CommandRegistry.
func NewCommandRegistry(versionCommand *VersionCommand, attuneCommand *AttuneCommand, configCommand *ConfigCommand,
	packageCommand *PackageCommand, helpCommand *HelpCommand) CommandRegistry {
	return CommandRegistry{
		helpCommand.Name():    helpCommand,
		versionCommand.Name(): versionCommand,
		attuneCommand.Name():  attuneCommand,
		configCommand.Name():  configCommand,
		packageCommand.Name(): packageCommand,
	}
}

// GetAllCommands returns a slice containing all commands.
func (commandRegistry CommandRegistry) GetAllCommands() []Command {
	var commandsSlice []Command
	for _, command := range commandRegistry {
		commandsSlice = append(commandsSlice, command)
	}
	return commandsSlice
}

// GetCommand returns the command with the given name, if it exists.
//
// It takes the following parameters:
//   - commandName: The name of the command.
func (commandRegistry CommandRegistry) GetCommand(commandName string) (Command, error) {
	command, isPresent := commandRegistry[commandName]
	if !isPresent {
		return nil, fmt.Errorf("command not valid")
	}

	return command, nil
}
