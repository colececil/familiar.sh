package commands

import (
	"fmt"
	"github.com/colececil/familiar.sh/internal/config"
)

// ConfigCommand represents the "config" command.
type ConfigCommand struct {
}

// Name returns the name of the command, as it appears on the command line while being used.
func (configCommand *ConfigCommand) Name() string {
	return "config"
}

// Description returns a short description of the command.
func (configCommand *ConfigCommand) Description() string {
	return "Print the contents of the shared configuration file or set the config file location"
}

// Documentation returns detailed documentation for the command.
func (configCommand *ConfigCommand) Documentation() string {
	return `The "config" command has the following subcommands:

location: Print the config file location or set the config file location to the given path.

Run "familiar help config location" for more information about the "location" subcommand.`
}

// Execute runs the command with the given arguments.
//
// It takes the following parameters:
//   - args: A slice containing the arguments to pass in to the command.
//
// If there is an error executing the command, Execute will return an error that can be displayed to the user.
func (configCommand *ConfigCommand) Execute(args []string) error {
	if len(args) == 0 {
		// Print the contents of the configuration file.
		return nil
	}

	switch args[0] {
	case "location":
		if len(args) == 1 {
			location, err := config.GetConfigLocation()
			if len(location) > 0 {
				fmt.Println(location)
			}
			return err
		} else {
			return config.SetConfigLocation(args[1])
		}
	default:
		return fmt.Errorf("unknown subcommand %q", args[0])
	}
}
