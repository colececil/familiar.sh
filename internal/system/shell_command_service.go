package system

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"
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

	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(2)
	errorChannel := make(chan error, 2)
	resultChannel := make(chan string, 1)

	go readLines(stdout, optionalOutputWriter, resultCaptureRegex, waitGroup, resultChannel, errorChannel)
	go readLines(stderr, optionalOutputWriter, nil, waitGroup, resultChannel, errorChannel)

	if err := command.Wait(); err != nil {
		return "", err
	}

	waitGroup.Wait()

	result := ""
	select {
	case err := <-errorChannel:
		return "", err
	case result = <-resultChannel:
	default:
	}

	if exitCode := command.ExitCode(); exitCode != 0 {
		return "", fmt.Errorf("error running command \"%s %s\", with exit code %d", program, strings.Join(args, " "),
			exitCode)
	}

	return result, nil
}

// readLines reads all lines of text from the given Reader and prints them to the given output writer. If the given
// regular expression finds a match, its submatch is written to the given results channel. If any error is encountered,
// it is written to the given error channel.
//
// It takes the following parameters:
//   - reader: The Reader to read from.
//   - outputWriter: The Writer to print the output to. If nil, no output is printed.
//   - resultCaptureRegex: A regular expression that captures any results using a capturing group. If this is nil or has
//     no capturing group, no results are captured.
//   - resultChannel: The channel to write the result to, if there is one.
//   - errorChannel: The channel to write an error to, if there is one.
//   - waitGroup: The WaitGroup to notify when the function has finished.
func readLines(reader io.Reader, outputWriter io.Writer, resultCaptureRegex *regexp.Regexp, waitGroup *sync.WaitGroup,
	resultChannel chan<- string, errorChannel chan<- error) {

	defer waitGroup.Done()

	scanner := bufio.NewScanner(reader)
	var cumulativeOutput = ""
	for scanner.Scan() {
		line := scanner.Text()
		cumulativeOutput += line + "\n"
		if outputWriter != nil {
			if _, err := fmt.Fprintln(outputWriter, line); err != nil {
				errorChannel <- err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		errorChannel <- err
		return
	}

	if resultCaptureRegex != nil {
		matchAndSubmatches := resultCaptureRegex.FindStringSubmatch(cumulativeOutput)
		if len(matchAndSubmatches) > 1 {
			resultChannel <- matchAndSubmatches[1]
		}
	}
}
