package main

import (
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
)

var _ command.Interface = new(checkCmd)

func TestCheck(t *testing.T) {
	t.Run("Execute", testCheckExecute)
	t.Run("SetFlags", testCheckSetFlags)
}

func testCheckExecute(t *testing.T) {
	os.Setenv("MY_ENV_VAR", "home")
	defer os.Unsetenv("MY_ENV_VAR")
	testCases := []struct {
		name string
		cmd  checkCmd
		err  error
	}{
		{
			name: "default",
			cmd:  checkCmd{},
			err:  nil,
		},
		{
			name: "fail",
			cmd:  checkCmd{fail: true},
			err:  errGotConflicts,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := filepath.Join("testdata", t.Name(), defaultConfig)
			cwd := filepath.Join("testdata", t.Name(), "targets")
			before, err := ioutil.ReadFile(config)
			if err != nil {
				t.Fatal(err)
			}
			var bd strings.Builder
			err = tc.cmd.Execute(&bd, opts{
				config:        config,
				getwd:         func() (string, error) { return cwd, nil },
				userConfigDir: func() (string, error) { return "user_config_dir", nil },
			})
			if want, got := tc.err, err; !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			golden := readFile(t, filepath.Join("testdata", t.Name())+".golden")
			if want, got := golden, bd.String(); got != want {
				t.Errorf("\nwant:\n%s\ngot:\n%s\n", want, got)
				t.Logf(
					"\"check\" command output mismatch (-want +got):\n%s",
					cmp.Diff(want, got),
				)
			}
			after, err := ioutil.ReadFile(config)
			if err != nil {
				t.Fatal(err)
			}
			// This guarantees the config file has only been read, not written.
			if want, got := before, after; string(got) != string(want) {
				t.Errorf("%s has been modified after command", config)
			}
		})
	}
}

func testCheckSetFlags(t *testing.T) {
	allowUnexported := cmp.AllowUnexported(checkCmd{})
	testCases := []struct {
		flags map[string]string
		want  checkCmd
	}{
		{
			flags: nil,
			want:  checkCmd{},
		},
		{
			flags: map[string]string{
				"fail": "true",
			},
			want: checkCmd{fail: true},
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			var (
				cmd  checkCmd
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
