package command

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/google/subcommands"
)

type ContextKey int

const (
	ErrCtxKey ContextKey = iota
	OptsCtxKey
)

// Command is a CLI command.
type Command struct {
	cmd      Interface
	name     string
	synopsis string

	stdout, stderr io.Writer
}

// New creates a new command.
func New(cmd Interface, opts ...func(*Command)) *Command {
	c := &Command{
		cmd:    cmd,
		stdout: ioutil.Discard,
		stderr: ioutil.Discard,
	}
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
func (c *Command) Usage() string { return fmt.Sprintf("%s (%s):\n", c.name, c.synopsis) }

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
		err = c.cmd.Execute(c.stdout, ctx.Value(OptsCtxKey))
	}
	if err != nil {
		fmt.Fprintf(c.stderr, "%v: %v\n", ctx.Value(ErrCtxKey), err)
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

// Stdout sets a standard output the command when it is executed.
// The default value is ioutil.Discard.
func Stdout(w io.Writer) func(c *Command) {
	return func(c *Command) {
		c.stdout = w
	}
}

// Stderr sets an output for errors returned when the command is executed.
// The default value is ioutil.Discard.
func Stderr(w io.Writer) func(c *Command) {
	return func(c *Command) {
		c.stderr = w
	}
}

// Interface is a command to be wrapped by Command.
type Interface interface {
	Execute(io.Writer, interface{}) error
	SetFlags(*flag.FlagSet)
}
