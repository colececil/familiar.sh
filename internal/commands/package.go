package commands

import (
	"fmt"
	"github.com/colececil/familiar.sh/internal/config"
)

// PackageCommand represents the "package" command.
type PackageCommand struct {
}

// Name returns the name of the command, as it appears on the command line while being used.
func (packageCommand *PackageCommand) Name() string {
	return "package"
}

// Description returns a short description of the command.
func (packageCommand *PackageCommand) Description() string {
	return "Manage packages for a given package manager."
}

// Documentation returns detailed documentation for the command.
func (packageCommand *PackageCommand) Documentation() string {
	return `The "package" command provides subcommands for adding, removing, and listing packages for a given package manager. It also allows you to specify the version of a package to install. It has the following subcommands:

  status
  search
  info
  add
  remove
  update`
}

// Execute runs the command with the given arguments.
//
// It takes the following parameters:
//   - args: A slice containing the arguments to pass in to the command.
//
// If there is an error executing the command, Execute will return an error that can be displayed to the user.
func (packageCommand *PackageCommand) Execute(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("subcommand must be included")
	}

	switch args[0] {
	case "add":
		if len(args) != 2 {
			return fmt.Errorf("wrong number of arguments")
		}

		err := config.AddPackageManager(args[1])
		if err != nil {
			return err
		}

		fmt.Println("Package manager added.")
		return nil
	case "remove":
		if len(args) != 2 {
			return fmt.Errorf("wrong number of arguments")
		}

		err := config.RemovePackageManager(args[1])
		if err != nil {
			return err
		}

		fmt.Println("Package manager removed.")
		return nil
	default:
		return fmt.Errorf("unknown subcommand %q", args[0])
	}
}
