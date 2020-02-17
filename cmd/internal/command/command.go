package command

import (
	"context"
	"flag"

	"github.com/google/subcommands"
)

// Command is a CLI command.
type Command struct {
	cmd      Interface
	name     string
	synopsis string
	usage    string
}

// New creates a new command.
func New(cmd Interface, opts ...func(*Command)) *Command {
	c := &Command{cmd: cmd}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Name returns the command's name.
func (c *Command) Name() string { return c.name }

// Synopsis returns the command's synopsis.
func (c *Command) Synopsis() string { return c.synopsis }

// Usage returns the command's usage instructions.
func (c *Command) Usage() string { return c.usage }

// SetFlags sets all necessary flags.
func (c *Command) SetFlags(f *flag.FlagSet) {
	c.cmd.SetFlags(f)
}

// Execute executes the command.
func (c *Command) Execute(ctx context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	var err error
	select {
	case <-ctx.Done():
		err = ctx.Err()
	default:
		err = c.cmd.Execute()
	}
	if err != nil {
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

// Name sets the command's name.
func Name(s string) func(*Command) {
	return func(c *Command) {
		c.name = s
	}
}

// Synopsis sets the command's synopsis.
func Synopsis(s string) func(*Command) {
	return func(c *Command) {
		c.synopsis = s
	}
}

// Usage sets the command's usage instructions.
func Usage(s string) func(*Command) {
	return func(c *Command) {
		c.usage = s
	}
}

// Interface is a command to be wrapped by Command.
type Interface interface {
	Execute() error
	SetFlags(*flag.FlagSet)
}
