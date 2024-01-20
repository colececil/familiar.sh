package commands

import "fmt"

// VersionCommand represents the "version" command.
type VersionCommand struct {
	// Version is the current version of Familiar.
	Version FamiliarVersionString
}

// FamiliarVersionString represents a version of Familiar.sh.
type FamiliarVersionString string

// NewVersionCommand returns a new instance of VersionCommand.
func NewVersionCommand(version FamiliarVersionString) *VersionCommand {
	return &VersionCommand{
		Version: version,
	}
}

// Name returns the name of the command, as it appears on the command line while being used.
func (c *VersionCommand) Name() string {
	return "version"
}

// Order returns the order in which the command should be listed in the help command.
func (c *VersionCommand) Order() int {
	return 2
}

// Description returns a short description of the command.
func (c *VersionCommand) Description() string {
	return "Print the installed version of Familiar."
}

// Documentation returns detailed documentation for the command.
func (c *VersionCommand) Documentation() string {
	return `Print the version of Familiar that is currently installed.

The version is displayed in the form "vX.Y.Z", where X is the major version number, Y is the minor version number, and Z is the patch number.`
}

// Execute runs the command with the given arguments.
//
// It takes the following parameters:
//   - args: A slice containing the arguments to pass in to the command.
//
// If there is an error executing the command, Execute will return an error that can be displayed to the user.
func (c *VersionCommand) Execute(args []string) error {
	fmt.Println("Familiar.sh v" + c.Version)
	return nil
}
