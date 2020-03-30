package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"strings"
	"testing"

	"github.com/gbrlsnchs/pilgo/cmd/internal/command"
	"github.com/google/go-cmp/cmp"
)

var _ command.Interface = versionCmd{}

func TestVersion(t *testing.T) {
	t.Run("Execute", testVersionExecute)
	t.Run("SetFlags", testVersionSetFlags)
}

func testVersionExecute(t *testing.T) {
	testCases := []struct {
		desc string
		cmd  versionCmd
		want string
		err  error
	}{
		{
			desc: "unknown version",
			cmd:  versionCmd{},
			want: "plg unknown version\n",
			err:  nil,
		},
		{
			desc: "show version",
			cmd:  versionCmd{"v0.0.0"},
			want: "plg v0.0.0\n",
			err:  nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			var (
				bd  strings.Builder
				ctx = context.WithValue(context.Background(), command.ErrCtxKey, "plg")
			)
			if want, got := tc.err, tc.cmd.Execute(ctx, &bd, nil); !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			if want, got := tc.want, bd.String(); got != want {
				t.Fatalf("want %q, got %q", want, got)
			}
		})
	}
}

func testVersionSetFlags(t *testing.T) {
	allowUnexported := cmp.AllowUnexported(versionCmd{})
	testCases := []struct {
		flags map[string]string
		want  versionCmd
	}{
		{
			flags: nil,
			want:  versionCmd{},
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			var (
				cmd  versionCmd
				fset = flag.NewFlagSet("version", flag.PanicOnError)
				args = make([]string, 0, len(tc.flags))
			)
			for name, value := range tc.flags {
				args = append(args, fmt.Sprintf("-%s=%s", name, value))
			}
			cmd.SetFlags(fset)
			t.Logf("args: %v", args)
			if err := fset.Parse(args); err != nil {
				t.Fatal(err)
			}
			if want, got := tc.want, cmd; !cmp.Equal(got, want, allowUnexported) {
				t.Errorf(
					"(*versionCmd).SetFlags mismatch (-want +got):\n%s",
					cmp.Diff(want, got, allowUnexported),
				)
			}
		})
	}
}
