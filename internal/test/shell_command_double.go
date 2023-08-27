package test

import (
	"errors"
	"io"
	"os"
)

// ShellCommandDouble is a test double for exec.Command.
type ShellCommandDouble struct {
	stdoutWriter    io.WriteCloser
	stderrWriter    io.WriteCloser
	started         bool
	exited          bool
	waiting         bool
	exitCode        int
	executionFunc   ShellCommandDoubleExecutionFunc
	errorChannel    chan error
	exitCodeChannel chan int
}

type ShellCommandDoubleExecutionFunc func(stdoutWriter io.Writer, stderrWriter io.Writer, errorChannel chan<- error,
	completionChannel chan<- int)

// NewShellCommandDouble returns a new instance of ShellCommandDouble.
func NewShellCommandDouble(executionFunc ShellCommandDoubleExecutionFunc) *ShellCommandDouble {
	return &ShellCommandDouble{
		exitCode:      -1,
		executionFunc: executionFunc,
	}
}

func (shellCommandDouble *ShellCommandDouble) StdoutPipe() (io.ReadCloser, error) {
	if shellCommandDouble.stdoutWriter != nil {
		return nil, errors.New("stdout already set")
	}
	if shellCommandDouble.started {
		return nil, errors.New("shell command already started")
	}
	pipeReader, pipeWriter, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	shellCommandDouble.stdoutWriter = pipeWriter
	return pipeReader, nil
}

func (shellCommandDouble *ShellCommandDouble) StderrPipe() (io.ReadCloser, error) {
	if shellCommandDouble.stderrWriter != nil {
		return nil, errors.New("stderr already set")
	}
	if shellCommandDouble.started {
		return nil, errors.New("shell command already started")
	}
	pipeReader, pipeWriter, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	shellCommandDouble.stderrWriter = pipeWriter
	return pipeReader, nil
}

func (shellCommandDouble *ShellCommandDouble) Start() error {
	shellCommandDouble.started = true
	shellCommandDouble.errorChannel = make(chan error)
	shellCommandDouble.exitCodeChannel = make(chan int)
	go shellCommandDouble.executionFunc(shellCommandDouble.stdoutWriter, shellCommandDouble.stderrWriter,
		shellCommandDouble.errorChannel, shellCommandDouble.exitCodeChannel)
	return nil
}

func (shellCommandDouble *ShellCommandDouble) Wait() error {
	if !shellCommandDouble.started {
		return errors.New("shell command not yet started")
	}
	if shellCommandDouble.waiting {
		return errors.New("already called")
	}

	shellCommandDouble.waiting = true
	select {
	case err := <-shellCommandDouble.errorChannel:
		return err
	case exitCode := <-shellCommandDouble.exitCodeChannel:
		shellCommandDouble.exitCode = exitCode
		shellCommandDouble.exited = true
		if err := shellCommandDouble.stdoutWriter.Close(); err != nil {
			return err
		}
		if err := shellCommandDouble.stderrWriter.Close(); err != nil {
			return err
		}
		return nil
	}
}

func (shellCommandDouble *ShellCommandDouble) ExitCode() int {
	return shellCommandDouble.exitCode
}
