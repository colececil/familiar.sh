package test

import (
	"fmt"
	"github.com/colececil/familiar.sh/internal/system"
	"regexp"
	"strings"
)

// ShellCommandServiceDouble contains a test double for system.ShellCommandService. The actual test double can be
// accessed via its ShellCommandService field.
type ShellCommandServiceDouble struct {
	ShellCommandService *system.ShellCommandService
}

var outputsGivenExpectedInputs map[shellCommandFuncInputs]string
var inputsUsed map[shellCommandFuncInputs]bool

// NewShellCommandServiceDouble returns a new instance of ShellCommandServiceDouble.
func NewShellCommandServiceDouble() *ShellCommandServiceDouble {
	outputsGivenExpectedInputs = make(map[shellCommandFuncInputs]string)
	inputsUsed = make(map[shellCommandFuncInputs]bool)
	return &ShellCommandServiceDouble{
		ShellCommandService: system.NewShellCommandService(runShellCommandFuncDouble),
	}
}

// shellCommandFuncInputs represents the input parameters used by the test double to determine the output of a shell
// command.
type shellCommandFuncInputs struct {
	program     string
	printOutput bool
	args        string
}

// SetOutputForExpectedInputs sets the output to return when the test double's RunShellCommand function is called with
// the given inputs.
func (shellCommandServiceDouble *ShellCommandServiceDouble) SetOutputForExpectedInputs(output string,
	expectedProgram string, expectedPrintOutput bool, expectedArgs ...string) {
	inputs := shellCommandFuncInputs{
		program:     expectedProgram,
		printOutput: expectedPrintOutput,
		args:        strings.Join(expectedArgs, " "),
	}
	outputsGivenExpectedInputs[inputs] = output
}

// WasCalledWith returns whether the test double's RunShellCommand function was called with the given inputs.
func (shellCommandServiceDouble *ShellCommandServiceDouble) WasCalledWith(program string, printOutput bool,
	args ...string) bool {
	inputs := shellCommandFuncInputs{
		program:     program,
		printOutput: printOutput,
		args:        strings.Join(args, " "),
	}
	return inputsUsed[inputs]
}

// runShellCommandFuncDouble is the implementation for the test double's RunShellCommand function. If an output has been
// set for the given inputs, resultCaptureRegex is run on the output and the result is returned.
func runShellCommandFuncDouble(program string, printOutput bool, resultCaptureRegex *regexp.Regexp,
	args ...string) (string, error) {
	inputs := shellCommandFuncInputs{
		program:     program,
		printOutput: printOutput,
		args:        strings.Join(args, " "),
	}
	inputsUsed[inputs] = true

	output, isPresent := outputsGivenExpectedInputs[inputs]
	if !isPresent {
		return "", fmt.Errorf("unexpected input")
	}

	var result string
	if resultCaptureRegex != nil && resultCaptureRegex.MatchString(output) {
		result = resultCaptureRegex.FindStringSubmatch(output)[1]
	}

	return result, nil
}
