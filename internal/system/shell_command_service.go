package system

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
)

// ShellCommandService provides functionality for running shell commands.
type ShellCommandService struct {
	createShellCommandFunc CreateShellCommandFunc
	runShellCommandFunc    RunShellCommandFunc
	outputWriter           io.Writer
}

// NewShellCommandService returns a new instance of ShellCommandService.
func NewShellCommandService(createShellCommandFunc CreateShellCommandFunc, runShellCommandFunc RunShellCommandFunc,
	outputWriter io.Writer) *ShellCommandService {
	return &ShellCommandService{
		createShellCommandFunc: createShellCommandFunc,
		runShellCommandFunc:    runShellCommandFunc,
		outputWriter:           outputWriter,
	}
}

// RunShellCommand runs a shell command for the given program and the given arguments. The command's output is printed
// to stdout.
//
// It takes the following parameters:
//   - program: The name of the program to run.
//   - printOutput: Whether to print the output of the program.
//   - resultCaptureRegex: A regular expression that captures the result of the command. If this is nil, the result is
//     an empty string.
//   - args: The arguments to pass to the program, if any.
//
// It returns the result captured by the regular expression, and an error if one occurred. If no result was captured,
// the result is an empty string.
func (shellCommandService *ShellCommandService) RunShellCommand(program string, printOutput bool,
	resultCaptureRegex *regexp.Regexp, args ...string) (string, error) {
	return shellCommandService.runShellCommandFunc(shellCommandService.createShellCommandFunc,
		shellCommandService.outputWriter, program, printOutput, resultCaptureRegex, args...)
}

// CreateShellCommandFunc is a function for creating a shell command to be run.
type CreateShellCommandFunc func(program string, args ...string) ShellCommand

// NewCreateShellCommandFunc returns a new instance of CreateShellCommandFunc.
func NewCreateShellCommandFunc() CreateShellCommandFunc {
	return defaultCreateShellCommandFunc
}

// defaultCreateShellCommandFunc is the default implementation of CreateShellCommandFunc.
func defaultCreateShellCommandFunc(program string, args ...string) ShellCommand {
	return NewRealShellCommand(program, args...)
}

// RunShellCommandFunc is a function for running a shell command.
type RunShellCommandFunc func(createShellCommand CreateShellCommandFunc, outputWriter io.Writer, program string,
	printOutput bool, resultCaptureRegex *regexp.Regexp, args ...string) (string, error)

// NewRunShellCommandFunc returns a new instance of RunShellCommandFunc.
func NewRunShellCommandFunc() RunShellCommandFunc {
	return defaultRunShellCommandFunc
}

// defaultRunShellCommandFunc is the default implementation of RunShellCommandFunc.
func defaultRunShellCommandFunc(createShellCommand CreateShellCommandFunc, outputWriter io.Writer, program string,
	printOutput bool, resultCaptureRegex *regexp.Regexp, args ...string) (string, error) {
	command := createShellCommand(program, args...)

	stdout, err := command.StdoutPipe()
	if err != nil {
		return "", err
	}

	stderr, err := command.StderrPipe()
	if err != nil {
		return "", err
	}

	if err := command.Start(); err != nil {
		return "", err
	}

	var optionalOutputWriter io.Writer
	if printOutput {
		optionalOutputWriter = outputWriter
	}

	errs := make(chan error)
	results := make(chan string)
	go readLines(stdout, optionalOutputWriter, resultCaptureRegex, results, errs)
	go readLines(stderr, optionalOutputWriter, nil, results, errs)

	select {
	case err := <-errs:
		return "", err
	default:
	}

	if err := command.Wait(); err != nil {
		return "", err
	}

	result := ""
	select {
	case result = <-results:
	default:
	}

	if exitCode := command.ExitCode(); exitCode != 0 {
		return "", fmt.Errorf("error running command \"%s %s\", with exit code %d", program, args, exitCode)
	}

	return result, nil
}

// readLines reads all lines of text from the given Reader and prints them to stdout. If the given regular expression
// finds a match, its submatch is written to the given results channel. If any error is encountered, it is written to
// the given error channel.
//
// It takes the following parameters:
//   - reader: The Reader to read from.
//   - outputWriter: The Writer to print the output to. If nil, no output is printed.
//   - resultCaptureRegex: A regular expression that captures any results using a capturing group. If this is nil or has
//     no capturing group, no results are captured.
//   - results: The channel to write any results to.
//   - errs: The channel to write any errors to.
func readLines(reader io.Reader, outputWriter io.Writer, resultCaptureRegex *regexp.Regexp, results chan<- string,
	errs chan<- error) {
	scanner := bufio.NewScanner(reader)

	var cumulativeOutput = ""
	for scanner.Scan() {
		line := scanner.Text()
		cumulativeOutput += line + "\n"
		if outputWriter != nil {
			if _, err := fmt.Fprintln(outputWriter, line); err != nil {
				errs <- err
			}
		}
	}

	if resultCaptureRegex != nil && resultCaptureRegex.MatchString(cumulativeOutput) {
		results <- resultCaptureRegex.FindStringSubmatch(cumulativeOutput)[1]
	}

	if err := scanner.Err(); err != nil {
		errs <- err
	}
}
