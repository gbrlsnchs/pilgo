package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/andybalholm/crlf"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/text/transform"
	"gsr.dev/pilgrim/cmd/internal/command"
)

var _ command.Interface = new(listCmd)
var allowUnexported = cmp.AllowUnexported(listCmd{})

func TestList(t *testing.T) {
	t.Run("Execute", testListExecute)
	t.Run("SetFlags", testListSetFlags)
}

func testListExecute(t *testing.T) {
	testCases := []struct {
		name string
		cmd  listCmd
		want string
		err  error
	}{
		{
			name: "default",
			cmd:  listCmd{},
			err:  nil,
		},
		{
			name: "check",
			cmd:  listCmd{check: true},
			err:  nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testname := filepath.Join("testdata", t.Name())
			tc.cmd.config = testname + ".json"
			before, err := ioutil.ReadFile(tc.cmd.config)
			if err != nil {
				t.Fatal(err)
			}
			tc.cmd.cwd = filepath.Join(testname, "targets")
			var bd strings.Builder
			if want, got := tc.err, tc.cmd.Execute(&bd); !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			golden, err := ioutil.ReadFile(testname + ".golden")
			if err != nil {
				t.Fatal(err)
			}
			goldenlf, _, err := transform.String(
				new(crlf.Normalize),
				filepath.FromSlash(string(golden)),
			)
			if err != nil {
				t.Fatal(err)
			}
			if want, got := goldenlf, bd.String(); got != want {
				t.Errorf(
					"Command list output mismatch (-want +got):\n%s",
					cmp.Diff(want, got),
				)
			}
			after, err := ioutil.ReadFile(tc.cmd.config)
			if err != nil {
				t.Fatal(err)
			}
			// This guarantees the config file has only been read, not written.
			if want, got := before, after; bytes.Compare(got, want) != 0 {
				t.Errorf("%s has been modified after listing", tc.cmd.config)
			}
		})
	}
}

func testListSetFlags(t *testing.T) {
	testCases := []struct {
		flags map[string]string
		want  listCmd
	}{
		{
			flags: nil,
			want: listCmd{
				check: false,
			},
		},
		{
			flags: map[string]string{
				"check": "false",
			},
			want: listCmd{
				check: false,
			},
		},
		{
			flags: map[string]string{
				"check": "true",
			},
			want: listCmd{
				check: true,
			},
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			var (
				cmd  listCmd
				fset = flag.NewFlagSet("list", flag.PanicOnError)
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
					"(*listCmd).SetFlags mismatch (-want +got):\n%s",
					cmp.Diff(want, got, allowUnexported),
				)
			}
		})
	}
}
