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
	c := command.New(new(ifaceMock), command.Name("foo"))
	if want, got := name, c.Name(); got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func testCommandSynopsis(t *testing.T) {
	synopsis := "bar"
	c := command.New(new(ifaceMock), command.Synopsis("bar"))
	if want, got := synopsis, c.Synopsis(); got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func testCommandUsage(t *testing.T) {
	usage := "foo (bar):\n"
	c := command.New(new(ifaceMock), command.Name("foo"), command.Synopsis("bar"))
	if want, got := usage, c.Usage(); got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func testCommandSetFlags(t *testing.T) {
	f := new(flag.FlagSet)
	mo := new(ifaceMock)
	c := command.New(mo)
	c.SetFlags(f)
	if want, got := f, mo.fset; got != want {
		t.Errorf("flag set mismatch")
	}
}

var panicErr = errors.New("panic!")

type ifaceMock struct {
	exec func(ctx context.Context, stdout, stderr io.Writer) error
	fset *flag.FlagSet
}

func (mo *ifaceMock) Execute(ctx context.Context, stdout, stderr io.Writer) error {
	return mo.exec(ctx, stdout, stderr)
}

func (mo *ifaceMock) SetFlags(f *flag.FlagSet) {
	mo.fset = f
}

func testCommandExecute(t *testing.T) {
	testCases := []struct {
		desc             string
		exec             func(ctx context.Context, stdout, stderr io.Writer) error
		cancelBeforeExec bool
		wantStatus       subcommands.ExitStatus
		wantStdout       string
		wantStderr       string
		wantOutput       string
		err              error
	}{
		{
			desc: "stdout only",
			exec: func(_ context.Context, stdout, _ io.Writer) error {
				fmt.Fprintln(stdout, "stdout")
				return nil
			},
			cancelBeforeExec: false,
			wantStatus:       subcommands.ExitSuccess,
			wantStdout:       "stdout\n",
			wantStderr:       "",
			wantOutput:       "stdout\n",
			err:              nil,
		},
		{
			desc: "stderr only",
			exec: func(_ context.Context, _, stderr io.Writer) error {
				fmt.Fprintln(stderr, "stderr")
				return nil
			},
			cancelBeforeExec: false,
			wantStatus:       subcommands.ExitSuccess,
			wantStdout:       "",
			wantStderr:       "stderr\n",
			wantOutput:       "stderr\n",
			err:              nil,
		},
		{
			desc: "stdout first and stderr later",
			exec: func(_ context.Context, stdout, stderr io.Writer) error {
				fmt.Fprintln(stdout, "stdout")
				fmt.Fprintln(stderr, "stderr")
				return nil
			},
			cancelBeforeExec: false,
			wantStatus:       subcommands.ExitSuccess,
			wantStdout:       "stdout\n",
			wantStderr:       "stderr\n",
			wantOutput:       "stdout\nstderr\n",
			err:              nil,
		},
		{
			desc: "stderr first and stdout later",
			exec: func(_ context.Context, stdout, stderr io.Writer) error {
				fmt.Fprintln(stderr, "stderr")
				fmt.Fprintln(stdout, "stdout")
				return nil
			},
			cancelBeforeExec: false,
			wantStatus:       subcommands.ExitSuccess,
			wantStdout:       "stdout\n",
			wantStderr:       "stderr\n",
			wantOutput:       "stderr\nstdout\n",
			err:              nil,
		},
		{
			desc: "error before stdout",
			exec: func(_ context.Context, stdout, _ io.Writer) error {
				fmt.Fprintln(stdout, "stdout")
				return errors.New("error!")
			},
			cancelBeforeExec: false,
			wantStatus:       subcommands.ExitFailure,
			wantStdout:       "stdout\n",
			wantStderr:       "command: error!\n",
			wantOutput:       "command: error!\nstdout\n",
			err:              nil,
		},
		{
			desc: "error before stderr",
			exec: func(_ context.Context, _, stderr io.Writer) error {
				fmt.Fprintln(stderr, "stderr")
				return errors.New("error!")
			},
			cancelBeforeExec: false,
			wantStatus:       subcommands.ExitFailure,
			wantStdout:       "",
			wantStderr:       "command: error!\nstderr\n",
			wantOutput:       "command: error!\nstderr\n",
			err:              nil,
		},
		{
			desc: "context canceled",
			exec: func(_ context.Context, stdout, stderr io.Writer) error {
				// Nothing in this function should be printed.
				fmt.Fprintln(stdout, "stdout")
				fmt.Fprintln(stderr, "stderr")
				return nil
			},
			cancelBeforeExec: true,
			wantStatus:       subcommands.ExitFailure,
			wantStdout:       "",
			wantStderr:       fmt.Sprintf("command: %s\n", context.Canceled),
			wantOutput:       fmt.Sprintf("command: %s\n", context.Canceled),
			err:              nil,
		},
		{
			desc: "panic",
			exec: func(_ context.Context, stdout, stderr io.Writer) error {
				fmt.Fprintln(stdout, "stdout")
				fmt.Fprintln(stderr, "stderr")
				panic(errors.New("panic!"))
			},
			cancelBeforeExec: false,
			wantStatus:       subcommands.ExitFailure,
			wantStdout:       "stdout\n",
			wantStderr:       "command: panic!\nstderr\n",
			wantOutput:       "command: panic!\nstdout\nstderr\n",
			err:              nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if tc.cancelBeforeExec {
				cancel()
			}
			var (
				stdout, stderr strings.Builder
				output         strings.Builder
			)
			c := command.New(&ifaceMock{tc.exec, nil},
				command.Stdout(io.MultiWriter(&stdout, &output)),
				command.Stderr(io.MultiWriter(&stderr, &output)),
			)
			status := c.Execute(context.WithValue(ctx, command.ErrCtxKey, "command"), nil)

			if want, got := tc.wantStatus, status; got != want {
				t.Errorf("status: want %d, got %d", want, got)
			}
			if want, got := tc.wantStdout, stdout.String(); got != want {
				t.Errorf("stdout: want %q, got %q", want, got)
			}
			if want, got := tc.wantStderr, stderr.String(); got != want {
				t.Errorf("stderr: want %q, got %q", want, got)
			}
			if want, got := tc.wantOutput, output.String(); got != want {
				t.Errorf("output: want %q, got %q", want, got)
			}
		})
	}
}
