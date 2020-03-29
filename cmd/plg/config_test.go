package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gbrlsnchs/pilgo/cmd/internal/command"
	"github.com/gbrlsnchs/pilgo/fs/fsutil"
	"github.com/google/go-cmp/cmp"
)

var _ command.Interface = new(configCmd)

func TestConfig(t *testing.T) {
	t.Run("Execute", testConfigExecute)
	t.Run("SetFlags", testConfigSetFlags)
}

func testConfigExecute(t *testing.T) {
	testCases := []struct {
		name string
		cmd  configCmd
		want string
		err  error
	}{
		{
			name: "config",
			cmd: configCmd{
				file:    "foo",
				baseDir: "test",
				link:    strptr{addr: newString("f00")},
				targets: commalist{
					"test",
					"testing",
					"testdata",
				},
			},
			err: nil,
		},
	}
	for _, tc := range testCases {
		testdata := filepath.Join("testdata", t.Name())
		t.Run(tc.name, func(t *testing.T) {
			config := filepath.Join(testdata, tc.name+".yml")
			before, err := ioutil.ReadFile(config)
			if err != nil {
				t.Fatal(err)
			}
			// Restore original contents of file.
			defer func(t *testing.T) {
				if err := ioutil.WriteFile(config, before, 0o644); err != nil {
					t.Fatal(err)
				}
			}(t)
			var (
				bd  strings.Builder
				ctx = context.WithValue(context.Background(), command.OptsCtxKey, opts{
					config:   config,
					fsDriver: fsutil.OSDriver{},
				})
			)
			if want, got := tc.err, tc.cmd.Execute(ctx, &bd, nil); !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			golden := readFile(t, filepath.Join(testdata, tc.name+".golden"))
			if want, got := tc.want, bd.String(); got != want {
				t.Errorf(
					`"show" command output mismatch (-want +got):\n%s`,
					cmp.Diff(want, got),
				)
			}
			after := readFile(t, config)
			if want, got := golden, after; string(got) != string(want) {
				t.Errorf("\nwant:\n%s\ngot:\n%s", want, got)
				t.Logf("detailed diff: %s", cmp.Diff(want, got))
			}
		})
	}
}

func testConfigSetFlags(t *testing.T) {
	allowUnexported := cmp.AllowUnexported(configCmd{}, strptr{})
	testCases := []struct {
		flags map[string]string
		want  configCmd
	}{
		{
			flags: nil,
			want: configCmd{
				baseDir: "",
				file:    "",
				link:    strptr{addr: nil},
				targets: nil,
			},
		},
		{
			flags: map[string]string{
				"file":    "test",
				"basedir": "testdata",
				"link":    "7357",
				"targets": "foo,bar,baz",
			},
			want: configCmd{
				file:    "test",
				baseDir: "testdata",
				link: strptr{
					addr: newString("7357"),
				},
				targets: commalist{"foo", "bar", "baz"},
			},
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			var (
				cmd  configCmd
				fset = flag.NewFlagSet("config", flag.PanicOnError)
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
					"(*configCmd).SetFlags mismatch (-want +got):\n%s",
					cmp.Diff(want, got, allowUnexported),
				)
			}
		})
	}
}

// TODO(gbrlsnchs): create reusable helper
func newString(s string) *string { return &s }
