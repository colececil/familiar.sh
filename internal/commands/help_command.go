package commands

import (
	"fmt"
)

// HelpCommand represents the "help" command.
type HelpCommand struct {
	// Commands is a slice containing all available commands.
	Commands []Command
}

// NewHelpCommand creates a new instance of HelpCommand.
func NewHelpCommand(versionCommand *VersionCommand, attuneCommand *AttuneCommand, configCommand *ConfigCommand,
	packageCommand *PackageCommand) *HelpCommand {
	return &HelpCommand{
		Commands: []Command{
			versionCommand,
			attuneCommand,
			configCommand,
			packageCommand,
		},
	}
}

// Name returns the name of the command, as it appears on the command line while being used.
func (helpCommand *HelpCommand) Name() string {
	return "help"
}

// Description returns a short description of the command.
func (helpCommand *HelpCommand) Description() string {
	return "List help information."
}

// Documentation returns detailed documentation for the command.
func (helpCommand *HelpCommand) Documentation() string {
	return "The `help` command lists information about all available Familiar CLI commands. If you provide a command name as an argument, it will display detailed documentation for that command.\n\nUsage:\n  familiar help\n  familiar help <command>"
}

// Execute runs the command with the given arguments.
//
// It takes the following parameters:
//   - args: A slice containing the arguments to pass in to the command.
//
// If there is an error executing the command, Execute will return an error that can be displayed to the user.
func (helpCommand *HelpCommand) Execute(args []string) error {
	if len(args) == 0 {
		// No command name was provided - display a list of all available commands.
		fmt.Println(`Usage: familiar <command> [<args>]

Available commands are listed below. Run "familiar help <command>" to get detailed documentation for a specific command.`)
		fmt.Printf("  %-15s %s\n", helpCommand.Name(), helpCommand.Description())
		for _, command := range helpCommand.Commands {
			fmt.Printf("  %-15s %s\n", command.Name(), command.Description())
		}
		fmt.Println()
		return nil
	} else {
		// A command name was provided - display detailed documentation for the command.
		name := args[0]
		isPresent := false
		if name == "help" {
			fmt.Printf("%s - %s\n\n%s\n", name, helpCommand.Description(), helpCommand.Documentation())
		} else {
			for _, command := range helpCommand.Commands {
				if command.Name() == name {
					isPresent = true
					fmt.Printf("%s - %s\n\n%s\n", name, command.Description(), command.Documentation())
					break
				}
			}
			if !isPresent {
				return fmt.Errorf("unknown command: %s", name)
			}
		}
		return nil
	}
}
