package system

import (
	"io"
	"os/exec"
)

// ShellCommand is an interface that provides methods for running a shell command.
type ShellCommand interface {
	// Start starts the command.
	Start() error

	// Wait waits for the command to finish, after it has been started by Start.
	Wait() error

	// StdoutPipe returns a pipe that will be connected to the command's standard output when the command starts.
	StdoutPipe() (io.ReadCloser, error)

	// StderrPipe returns a pipe that will be connected to the command's standard error when the command starts.
	StderrPipe() (io.ReadCloser, error)

	// ExitCode returns the exit code of the command, or -1 if the command hasn't exited.
	ExitCode() int

	// String returns the string representation of the command.
	String() string
}

// ShellCommandFactoryFunc is a function that creates a new instance of ShellCommand.
type ShellCommandFactoryFunc func(program string, args ...string) ShellCommand

// NewShellCommand creates a new instance of ShellCommand, using the default implementation. Its parameters are the
// program and arguments of the command to run.
func NewShellCommand(program string, args ...string) ShellCommand {
	return &shellCommand{exec.Command(program, args...)}
}

// shellCommand is the default implementation of the ShellCommand interface. It uses exec.Cmd.
type shellCommand struct {
	*exec.Cmd
}

// ExitCode implements ShellCommand.ExitCode by retrieving it from the underlying exec.Cmd.
func (c *shellCommand) ExitCode() int {
	if c.ProcessState == nil {
		return -1
	}
	return c.ProcessState.ExitCode()
}
