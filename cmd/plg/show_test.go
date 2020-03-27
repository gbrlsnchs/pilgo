package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gsr.dev/pilgrim/cmd/internal/command"
	"gsr.dev/pilgrim/fs/fsutil"
)

var _ command.Interface = showCmd{}

func TestShow(t *testing.T) {
	t.Run("Execute", testShowExecute)
	t.Run("SetFlags", testShowSetFlags)
}

func testShowExecute(t *testing.T) {
	os.Setenv("MY_ENV_VAR", "home")
	defer os.Unsetenv("MY_ENV_VAR")
	testCases := []struct {
		name string
		cmd  showCmd
		err  error
	}{
		{
			name: "show",
			cmd:  showCmd{},
			err:  nil,
		},
	}
	for _, tc := range testCases {
		testdata := filepath.Join("testdata", t.Name())
		t.Run(tc.name, func(t *testing.T) {
			config := filepath.Join(testdata, defaultConfig)
			before, err := ioutil.ReadFile(config)
			if err != nil {
				t.Fatal(err)
			}
			var (
				bd  strings.Builder
				ctx = context.WithValue(context.Background(), command.OptsCtxKey, opts{
					config:   config,
					fsDriver: fsutil.OSDriver{},
					getwd: func() (string, error) {
						return filepath.Join(testdata, "targets"), nil
					},
					userConfigDir: func() (string, error) { return "user_config_dir", nil },
				})
			)
			if want, got := tc.err, tc.cmd.Execute(ctx, &bd, nil); !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			golden := readFile(t, filepath.Join("testdata", t.Name())+".golden")
			if want, got := golden, bd.String(); got != want {
				t.Errorf(
					"\"show\" command output mismatch (-want +got):\n%s",
					cmp.Diff(want, got),
				)
			}
			after, err := ioutil.ReadFile(config)
			if err != nil {
				t.Fatal(err)
			}
			// This guarantees the config file has only been read, not written.
			if want, got := before, after; bytes.Compare(got, want) != 0 {
				t.Errorf("%s has been modified after command", config)
			}
		})
	}
}

func testShowSetFlags(t *testing.T) {
	allowUnexported := cmp.AllowUnexported(showCmd{})
	testCases := []struct {
		flags map[string]string
		want  showCmd
	}{
		{
			flags: nil,
			want:  showCmd{},
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			var (
				cmd  showCmd
				fset = flag.NewFlagSet("show", flag.PanicOnError)
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
					"(*showCmd).SetFlags mismatch (-want +got):\n%s",
					cmp.Diff(want, got, allowUnexported),
				)
			}
		})
	}
}
