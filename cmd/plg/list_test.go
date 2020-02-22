package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/andybalholm/crlf"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/text/transform"
)

func TestList(t *testing.T) {
	t.Run("Execute", testListExecute)
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
