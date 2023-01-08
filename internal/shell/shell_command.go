package shell

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"regexp"
)

// RunShellCommand runs a shell command for the given program and the given arguments. The command's output is printed
// to stdout.
//
// It takes the following parameters:
//   - program: The name of the program to run.
//   - resultCaptureRegex: A regular expression that captures the result of the command. If this is nil, the result is
//     an empty string.
//   - args: The arguments to pass to the program, if any.
//
// It returns the result captured by the regular expression, and an error if one occurred. If no result was captured,
// the result is an empty string.
func RunShellCommand(program string, resultCaptureRegex *regexp.Regexp, args ...string) (string, error) {
	command := exec.Command(program, args...)

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

	errs := make(chan error)
	results := make(chan string)
	go readLines(stdout, resultCaptureRegex, results, errs)
	go readLines(stderr, resultCaptureRegex, results, errs)

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

	return result, nil
}

// readLines reads all lines of text from the given Reader and prints them to stdout. If the given regular expression
// finds a match, its submatch is written to the given results channel. If any error is encountered, it is written to
// the given error channel.
//
// It takes the following parameters:
//   - reader: The Reader to read from.
//   - resultCaptureRegex: A regular expression that captures any results. If this is nil, no results are captured.
//   - results: The channel to write any results to.
//   - errs: The channel to write any errors to.
func readLines(reader io.Reader, resultCaptureRegex *regexp.Regexp, results chan<- string, errs chan<- error) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()

		if resultCaptureRegex != nil && resultCaptureRegex.MatchString(line) {
			results <- resultCaptureRegex.FindStringSubmatch(line)[1]
		}

		fmt.Println(line)
	}
	if err := scanner.Err(); err != nil {
		errs <- err
	}
}
