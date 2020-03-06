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
	usage := "foo (bar):\n"
	c := command.New(new(ifaceSpy), command.Name("foo"), command.Synopsis("bar"))
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
	txt  string
	err  error
	fset *flag.FlagSet
}

func (spy *ifaceSpy) Execute(w io.Writer, opts interface{}) error {
	fmt.Fprintf(w, "%s", spy.txt)
	return spy.err
}

func (spy *ifaceSpy) SetFlags(f *flag.FlagSet) {
	spy.fset = f
}

func testCommandExecute(t *testing.T) {
	testCases := []struct {
		wantStatus subcommands.ExitStatus
		out        string
		wantStdout string
		err        error
		cancel     bool
		wantStderr string
		wantOutput string
	}{
		{
			wantStatus: subcommands.ExitSuccess,
			out:        "test\n",
			wantStdout: "test\n",
			err:        nil,
			cancel:     false,
			wantStderr: "",
			wantOutput: "test\n",
		},
		{
			wantStatus: subcommands.ExitFailure,
			out:        "test\n",
			wantStdout: "test\n",
			err:        errors.New("oops"),
			cancel:     false,
			wantStderr: "command: oops\n",
			wantOutput: "command: oops\ntest\n",
		},
		{
			wantStatus: subcommands.ExitFailure,
			out:        "test\n",
			wantStdout: "", // context is checked before execution
			err:        nil,
			cancel:     true,
			wantStderr: fmt.Sprintf("command: %v\n", context.Canceled),
			wantOutput: fmt.Sprintf("command: %v\n", context.Canceled),
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
				output         strings.Builder
				stdout, stderr strings.Builder
				spy            = &ifaceSpy{tc.out, tc.err, nil}
				c              = command.New(
					spy,
					command.Stdout(io.MultiWriter(&output, &stdout)),
					command.Stderr(io.MultiWriter(&output, &stderr)),
				)
				status = c.Execute(context.WithValue(ctx, command.ErrCtxKey, "command"), nil)
			)
			if want, got := tc.wantStatus, status; got != want {
				t.Errorf("want %d, got %d", want, got)
			}
			if want, got := tc.wantStdout, stdout.String(); got != want {
				t.Errorf("want %q, got %q", want, got)
			}
			if want, got := tc.wantStderr, stderr.String(); got != want {
				t.Errorf("want %q, got %q", want, got)
			}
			if want, got := tc.wantOutput, output.String(); got != want {
				t.Errorf("want %q, got %q", want, got)
			}
		})
	}
}
