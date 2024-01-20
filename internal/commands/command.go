package commands

// A Command represents a command that can be executed via the Familiar CLI.
//
// It encapsulates the information necessary to execute and provide documentation for the command.
type Command interface {
	// Name returns the name of the command, as it appears on the command line while being used.
	Name() string

	// Order returns the order in which the command should be listed in the help command.
	Order() int

	// Description returns a short description of the command.
	Description() string

	// Documentation returns detailed documentation for the command.
	Documentation() string

	// Execute runs the command with the given arguments.
	//
	// It takes the following parameters:
	//   - args: A slice containing the arguments to pass in to the command.
	//
	// If there is an error executing the command, Execute will return an error that can be displayed to the user.
	Execute(args []string) error
}
