package system

import (
	"io"
	"os/exec"
)

type ShellCommand interface {
	StdoutPipe() (io.ReadCloser, error)
	StderrPipe() (io.ReadCloser, error)
	Start() error
	Wait() error
	ExitCode() int
}

type RealShellCommand struct {
	cmd *exec.Cmd
}

func NewRealShellCommand(program string, args ...string) *RealShellCommand {
	return &RealShellCommand{
		cmd: exec.Command(program, args...),
	}
}

func (realShellCommand *RealShellCommand) StdoutPipe() (io.ReadCloser, error) {
	return realShellCommand.cmd.StdoutPipe()
}

func (realShellCommand *RealShellCommand) StderrPipe() (io.ReadCloser, error) {
	return realShellCommand.cmd.StderrPipe()
}

func (realShellCommand *RealShellCommand) Start() error {
	return realShellCommand.cmd.Start()
}

func (realShellCommand *RealShellCommand) Wait() error {
	return realShellCommand.cmd.Wait()
}

func (realShellCommand *RealShellCommand) ExitCode() int {
	return realShellCommand.cmd.ProcessState.ExitCode()
}
