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
	var (
		status = subcommands.ExitSuccess
		err    error
		lw     = &lineWriter{stdout: c.stdout, stderr: c.stderr}
	)
	select {
	case <-ctx.Done():
		err = ctx.Err()
	default:
		// NOTE: By using this pattern, we're able to share the underlying line slice from lw
		// while implementing different writing methods.
		stdout := (*stdoutWriter)(lw)
		stderr := (*stderrWriter)(lw)
		err = c.execute(ctx, stdout, stderr)
	}
	if err != nil {
		fmt.Fprintf(c.stderr, "%v: %v\n", ctx.Value(ErrCtxKey), err)
		status = subcommands.ExitFailure
	}
	for _, line := range lw.lines {
		fmt.Fprint(line.w, line.txt)
	}
	return status
}

func (c *Command) execute(ctx context.Context, stdout, stderr io.Writer) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	return c.cmd.Execute(ctx, stdout, stderr)
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
	Execute(ctx context.Context, stdout, stderr io.Writer) error
	SetFlags(*flag.FlagSet)
}

type line struct {
	w   io.Writer
	txt string
}

type lineWriter struct {
	stdout, stderr io.Writer
	lines          []line
}

type stdoutWriter lineWriter

func (w *stdoutWriter) Write(b []byte) (int, error) {
	l := line{w.stdout, string(b)}
	w.lines = append(w.lines, l)
	return len(b), nil
}

type stderrWriter lineWriter

func (w *stderrWriter) Write(b []byte) (int, error) {
	l := line{w.stderr, string(b)}
	w.lines = append(w.lines, l)
	return len(b), nil
}
