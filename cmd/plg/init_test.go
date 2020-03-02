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
	"gopkg.in/yaml.v3"
	"gsr.dev/pilgrim"
	"gsr.dev/pilgrim/cmd/internal/command"
)

var _ command.Interface = new(initCmd)

func TestInit(t *testing.T) {
	t.Run("Execute", testInitExecute)
	t.Run("SetFlags", testInitSetFlags)
}

func testInitExecute(t *testing.T) {
	testCases := []struct {
		name   string
		cmd    initCmd
		want   pilgrim.Config
		err    error
		remove bool
	}{
		{
			name: "init",
			cmd:  initCmd{},
			want: pilgrim.Config{
				Targets: []string{
					"bar",
					"foo",
				},
			},
			err:    nil,
			remove: true,
		},
		{
			name: "force",
			cmd:  initCmd{force: true},
			want: pilgrim.Config{
				BaseDir: "/tmp",
				Targets: []string{
					"bar",
					"foo",
				},
			},
			err:    nil,
			remove: false,
		},
		{
			name: "nop",
			cmd:  initCmd{},
			want: pilgrim.Config{
				BaseDir: "/tmp",
				Targets: []string{
					"test",
				},
			},
			err:    errConfigExists,
			remove: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testdata := filepath.Join("testdata", t.Name())
			config := filepath.Join(testdata, defaultConfig)
			if tc.remove {
				defer os.Remove(config)
			} else if tc.cmd.force {
				orig, err := ioutil.ReadFile(config)
				if err != nil {
					t.Fatal(err)
				}
				defer func(t *testing.T) {
					if err := ioutil.WriteFile(config, orig, 0o644); err != nil {
						t.Fatal(err)
					}
				}(t)
			}
			tc.cmd.config = config
			tc.cmd.cwd = testdata
			var bd strings.Builder
			if want, got := tc.err, tc.cmd.Execute(&bd); !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			if want, got := "", bd.String(); got != want {
				t.Errorf("want %q, got %q", want, got)
			}
			var c pilgrim.Config
			b, err := ioutil.ReadFile(config)
			if err != nil {
				t.Fatal(err)
			}
			if err = yaml.Unmarshal(b, &c); err != nil {
				t.Fatal(err)
			}
			if want, got := tc.want, c; !cmp.Equal(got, want) {
				t.Errorf("command \"init\" outcome failed (-want +got):\n%s", cmp.Diff(want, got))
			}
		})
	}
}

func testInitSetFlags(t *testing.T) {
	allowUnexported := cmp.AllowUnexported(initCmd{})
	testCases := []struct {
		flags map[string]string
		want  initCmd
	}{
		{
			flags: nil,
			want:  initCmd{},
		},
		{
			flags: map[string]string{
				"force": "true",
			},
			want: initCmd{force: true},
		},
		{
			flags: map[string]string{
				"force": "false",
			},
			want: initCmd{force: false},
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			var (
				cmd  initCmd
				fset = flag.NewFlagSet("init", flag.PanicOnError)
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
					"(*initCmd).SetFlags mismatch (-want +got):\n%s",
					cmp.Diff(want, got, allowUnexported),
				)
			}
		})
	}
}
