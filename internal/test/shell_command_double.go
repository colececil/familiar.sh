package test

import (
	"github.com/colececil/familiar.sh/internal/system"
	"os/exec"
	"strings"
)

// ShellCommandDouble is a test double for exec.Command.
type ShellCommandDouble struct {
	CreateShellCommandFunc system.CreateShellCommandFunc
}

var createShellCommandExpectedInputs map[createShellCommandFuncInputs]*exec.Cmd

// NewShellCommandDouble returns a new instance of ShellCommandDouble.
func NewShellCommandDouble() *ShellCommandDouble {
	createShellCommandExpectedInputs = make(map[createShellCommandFuncInputs]*exec.Cmd)
	return &ShellCommandDouble{
		CreateShellCommandFunc: createShellCommandFuncDouble,
	}
}

// createShellCommandFuncInputs represents the input parameters used by the test double to determine the instance of
// *exec.Cmd to return.
type createShellCommandFuncInputs struct {
	program string
	args    string
}

// SetOutputForExpectedInputs sets the output to return when the test double's CreateShellCommandFunc function is called
// with the given inputs.
func (shellCommandDouble *ShellCommandDouble) SetOutputForExpectedInputs(output *exec.Cmd, expectedProgram string,
	expectedArgs ...string) {
	inputs := createShellCommandFuncInputs{
		program: expectedProgram,
		args:    strings.Join(expectedArgs, " "),
	}
	createShellCommandExpectedInputs[inputs] = output
}

// createShellCommandFuncDouble is the implementation of the test double's CreateShellCommandFunc function. If an output
// has been set for the given inputs, it will be returned. Otherwise, exec.Command("") is returned.
func createShellCommandFuncDouble(program string, args ...string) *exec.Cmd {
	inputs := createShellCommandFuncInputs{
		program: program,
		args:    strings.Join(args, " "),
	}

	output, isPresent := createShellCommandExpectedInputs[inputs]
	if !isPresent {
		return exec.Command("")
	}

	return output
}
