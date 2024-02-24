package system

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

// ShellCommandRunner is an interface that provides a method to run a shell command and capture its output.
type ShellCommandRunner interface {
	// Run runs the shell command that was specified when creating the ShellCommandRunner. It optionally takes in a
	// regular expression that is matched against the complete  stdout output of the command. If there is a match, the
	// match of its first capturing group will be returned. If there is no match or no regular expression is provided,
	// an empty string will be returned.
	//
	// Additionally, if an outputWriter was provided when creating the ShellCommandRunner, both the stdout and stderr
	// output of the command will be written to it.
	//
	// If there are any issues running the command or processing its output, an error will be returned.
	Run(resultCaptureRegex *regexp.Regexp) (string, error)
}

// NewShellCommandRunner creates a new instance of a type that implements the ShellCommandRunner interface, using the
// default implementation. If an outputWriter is provided, the output of the command will be written to it when it is
// run. If the full output does not need to be captured, the outputWriter can be set to nil. The command to run is
// specified by the program and args parameters.
func NewShellCommandRunner(outputWriter io.Writer, program string, args ...string) ShellCommandRunner {
	return ShellCommandRunner(
		&shellCommandRunner{
			cmd:          exec.Command(program, args...),
			outputWriter: outputWriter,
		},
	)
}

// shellCommandRunner is the default implementation of the ShellCommandRunner interface. It runs the specified command
// using exec.Command.
type shellCommandRunner struct {
	cmd          *exec.Cmd
	outputWriter io.Writer
}

// Run implements ShellCommandRunner.Run by running the command and processing its output.
func (r *shellCommandRunner) Run(resultCaptureRegex *regexp.Regexp) (string, error) {
	stdout, err := r.cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	stderr, err := r.cmd.StderrPipe()
	if err != nil {
		return "", err
	}

	if err := r.cmd.Start(); err != nil {
		return "", err
	}

	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(2)
	errorChannel := make(chan error, 2)
	resultChannel := make(chan string, 1)

	go readLines(stdout, r.outputWriter, resultCaptureRegex, waitGroup, resultChannel, errorChannel)
	go readLines(stderr, r.outputWriter, nil, waitGroup, resultChannel, errorChannel)

	if err := r.cmd.Wait(); err != nil {
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

	if exitCode := r.cmd.ProcessState.ExitCode(); exitCode != 0 {
		return "",
			fmt.Errorf("error running command \"%s\", with exit code %d", strings.Join(r.cmd.Args, " "), exitCode)
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
