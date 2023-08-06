package test

import (
	"bytes"
	"io"
)

// ShellCommandDouble is a test double for exec.Command.
type ShellCommandDouble struct {
	stdoutOutput string
	stderrOutput string
	exitCode     int
}

// NewShellCommandDouble returns a new instance of ShellCommandDouble.
func NewShellCommandDouble(stdoutOutput string, stderrOutput string, exitCode int) *ShellCommandDouble {
	return &ShellCommandDouble{
		stdoutOutput: stdoutOutput,
		stderrOutput: stderrOutput,
		exitCode:     exitCode,
	}
}

func (shellCommandDouble *ShellCommandDouble) StdoutPipe() (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewBufferString(shellCommandDouble.stdoutOutput)), nil
}

func (shellCommandDouble *ShellCommandDouble) StderrPipe() (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewBufferString(shellCommandDouble.stderrOutput)), nil
}

func (shellCommandDouble *ShellCommandDouble) Start() error {
	// Todo: Make `Wait` throw an error if this method is not called first.
	return nil
}

func (shellCommandDouble *ShellCommandDouble) Wait() error {
	// Todo: Make this wait until copying from stdout and stderr is complete.
	return nil
}

func (shellCommandDouble *ShellCommandDouble) ExitCode() int {
	return shellCommandDouble.exitCode
}
