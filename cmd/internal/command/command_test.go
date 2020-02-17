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
	testCases := []struct {
		wantStatus subcommands.ExitStatus
		out        string
		wantOut    string
		err        error
		cancel     bool
		wantErrMsg string
	}{
		{
			wantStatus: subcommands.ExitSuccess,
			out:        "test",
			wantOut:    "test",
			err:        nil,
			cancel:     false,
			wantErrMsg: "",
		},
		{
			wantStatus: subcommands.ExitFailure,
			out:        "test",
			wantOut:    "test",
			err:        errors.New("oops"),
			cancel:     false,
			wantErrMsg: "command: oops\n",
		},
		{
			wantStatus: subcommands.ExitFailure,
			out:        "test",
			wantOut:    "", // context is checked before execution
			err:        nil,
			cancel:     true,
			wantErrMsg: fmt.Sprintf("command: %v\n", context.Canceled),
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if tc.cancel {
				cancel()
			}
			var (
				bd     strings.Builder
				stderr strings.Builder
				spy    = &ifaceSpy{&bd, tc.out, tc.err, nil}
				c      = command.New(spy, command.Stderr(&stderr))
				status = c.Execute(ctx, nil)
			)
			if want, got := tc.wantStatus, status; got != want {
				t.Errorf("want %d, got %d", want, got)
			}
			if want, got := tc.wantOut, bd.String(); got != want {
				t.Errorf("want %q, got %q", want, got)
			}
			if want, got := tc.wantErrMsg, stderr.String(); got != want {
				t.Errorf("want %q, got %q", want, got)
			}
		})
	}
}
