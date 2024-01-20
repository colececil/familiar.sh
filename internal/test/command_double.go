package test

import "github.com/colececil/familiar.sh/internal/commands"

// CommandDouble is a test double for the Command interface.
type CommandDouble struct {
	commands.Command
	name  string
	order int
}

// NewCommandDouble returns a new instance of CommandDouble. It takes the name of the command and its order as input.
func NewCommandDouble(name string, order int) *CommandDouble {
	return &CommandDouble{
		name:  name,
		order: order,
	}
}

// Name returns the name of the command.
func (d *CommandDouble) Name() string {
	return d.name
}

// Order returns the order of the command.
func (d *CommandDouble) Order() int {
	return d.order
}
