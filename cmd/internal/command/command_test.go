package command_test

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/google/subcommands"
	"gsr.dev/pilgrim/cmd/internal/command"
)

var _ subcommands.Command = new(command.Command)

func TestCommand(t *testing.T) {
	t.Run("Name", testCommandName)
	t.Run("Synopsis", testCommandSynopsis)
	t.Run("Usage", testCommandUsage)
	t.Run("SetFlags", testCommandSetFlags)
	t.Run("Execute", testCommandExecute)
}

func testCommandName(t *testing.T) {
	name := "foo"
	c := command.New(new(ifaceSpy), command.Name("foo"))
	if want, got := name, c.Name(); got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func testCommandSynopsis(t *testing.T) {
	synopsis := "bar"
	c := command.New(new(ifaceSpy), command.Synopsis("bar"))
	if want, got := synopsis, c.Synopsis(); got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func testCommandUsage(t *testing.T) {
	usage := "baz"
	c := command.New(new(ifaceSpy), command.Usage("baz"))
	if want, got := usage, c.Usage(); got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func testCommandSetFlags(t *testing.T) {
	f := new(flag.FlagSet)
	spy := new(ifaceSpy)
	c := command.New(spy)
	c.SetFlags(f)
	if want, got := f, spy.fset; got != want {
		t.Errorf("flag set mismatch")
	}
}

type ifaceSpy struct {
	w    io.Writer
	txt  string
	err  error
	fset *flag.FlagSet
}

func (spy *ifaceSpy) Execute() error {
	fmt.Fprintf(spy.w, "%s", spy.txt)
	return spy.err
}

func (spy *ifaceSpy) SetFlags(f *flag.FlagSet) {
	spy.fset = f
}

func testCommandExecute(t *testing.T) {
	var (
		bd          strings.Builder
		stderr strings.Builder
		spy         = &ifaceSpy{&bd, "foobar", nil, nil}
		c           = command.New(spy, command.Stderr(&stderr))
		ctx, cancel = context.WithCancel(context.Background())
		status      = c.Execute(ctx, nil)
	)
	if want, got := subcommands.ExitSuccess, status; got != want {
		t.Errorf("want %d, got %d", want, got)
	}
	if want, got := "foobar", bd.String(); got != want {
		t.Errorf("want %q, got %q", want, got)
	}
	if want, got := "", stderr.String(); got != want {
		t.Errorf("want %q, got %q", want, got)
	}
	bd.Reset()
	stderr.Reset()
	spy.err = errors.New("oops")
	status = c.Execute(ctx, nil)
	if want, got := subcommands.ExitFailure, status; got != want {
		t.Errorf("want %d, got %d", want, got)
	}
	if want, got := "command: oops\n", stderr.String(); got != want {
		t.Errorf("want %q, got %q", want, got)
	}
	bd.Reset()
	stderr.Reset()
	spy.err = nil
	cancel()
	status = c.Execute(ctx, nil)
	if want, got := subcommands.ExitFailure, status; got != want {
		t.Errorf("want %d, got %d", want, got)
	}
	if want, got := fmt.Sprintf("command: %v\n", context.Canceled), stderr.String(); got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}
