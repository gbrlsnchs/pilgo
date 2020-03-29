package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gbrlsnchs/pilgo/cmd/internal/command"
	"github.com/gbrlsnchs/pilgo/fs/fstest"
	"github.com/gbrlsnchs/pilgo/linker"
	"github.com/google/go-cmp/cmp"
)

var _ command.Interface = linkCmd{}

func TestLink(t *testing.T) {
	t.Run("Execute", testLinkExecute)
	t.Run("SetFlags", testLinkSetFlags)
}

func testLinkExecute(t *testing.T) {
	os.Setenv("MY_ENV_VAR", "my_file")
	defer os.Unsetenv("MY_ENV_VAR")
	testCases := []struct {
		name string
		drv  fstest.InMemoryDriver
		cmd  linkCmd
		want fstest.InMemoryDriver
		err  error
	}{
		{
			name: "default",
			drv: fstest.InMemoryDriver{
				Files: map[string]fstest.File{
					"home": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     nil,
						Children: map[string]fstest.File{
							"dotfiles": {
								Linkname: "",
								Perm:     os.ModePerm,
								Data:     nil,
								Children: map[string]fstest.File{
									"test": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     []byte("foo"),
										Children: nil,
									},
									"my_file": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     []byte("bar"),
										Children: nil,
									},
									defaultConfig: {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     testDefaultConfig,
										Children: nil,
									},
								},
							},
							"config": {
								Linkname: "",
								Perm:     os.ModePerm,
								Data:     nil,
								Children: make(map[string]fstest.File, 0),
							},
						},
					},
				},
			},
			cmd: linkCmd{},
			want: fstest.InMemoryDriver{
				Files: map[string]fstest.File{
					"home": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     nil,
						Children: map[string]fstest.File{
							"dotfiles": {
								Linkname: "",
								Perm:     os.ModePerm,
								Data:     nil,
								Children: map[string]fstest.File{
									"test": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     []byte("foo"),
										Children: nil,
									},
									"my_file": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     []byte("bar"),
										Children: nil,
									},
									defaultConfig: {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     testDefaultConfig,
										Children: nil,
									},
								},
							},
							"config": {
								Linkname: "",
								Perm:     os.ModePerm,
								Data:     nil,
								Children: map[string]fstest.File{
									"test": {
										Linkname: filepath.Join("home", "dotfiles", "test"),
										Data:     nil,
										Perm:     os.ModePerm,
										Children: nil,
									},
									"my_file": {
										Linkname: filepath.Join("home", "dotfiles", "my_file"),
										Data:     nil,
										Perm:     os.ModePerm,
										Children: nil,
									},
								},
							},
						},
					},
				},
			},
			err: nil,
		},
		{
			name: "conflict",
			drv: fstest.InMemoryDriver{
				Files: map[string]fstest.File{
					"home": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     nil,
						Children: map[string]fstest.File{
							"dotfiles": {
								Linkname: "",
								Perm:     os.ModePerm,
								Data:     nil,
								Children: map[string]fstest.File{
									"test": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     []byte("foo"),
										Children: make(map[string]fstest.File, 0),
									},
									"my_file": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     []byte("bar"),
										Children: nil,
									},
									defaultConfig: {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     testDefaultConfig,
										Children: nil,
									},
								},
							},
							"config": {
								Linkname: "",
								Perm:     os.ModePerm,
								Data:     nil,
								Children: map[string]fstest.File{
									"test": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     []byte("bar"),
										Children: nil,
									},
								},
							},
						},
					},
				},
			},
			cmd: linkCmd{},
			want: fstest.InMemoryDriver{
				Files: map[string]fstest.File{
					"home": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     nil,
						Children: map[string]fstest.File{
							"dotfiles": {
								Linkname: "",
								Perm:     os.ModePerm,
								Data:     nil,
								Children: map[string]fstest.File{
									"test": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     []byte("foo"),
										Children: make(map[string]fstest.File, 0),
									},
									"my_file": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     []byte("bar"),
										Children: nil,
									},
									defaultConfig: {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     testDefaultConfig,
										Children: nil,
									},
								},
							},
							"config": {
								Linkname: "",
								Perm:     os.ModePerm,
								Data:     nil,
								Children: map[string]fstest.File{
									"test": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     []byte("bar"),
										Children: nil,
									},
								},
							},
						},
					},
				},
			},
			err: (*linker.ConflictError)(nil),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				bd   strings.Builder
				opts = opts{
					config:        defaultConfig,
					fsDriver:      &tc.drv,
					getwd:         func() (string, error) { return filepath.Join("home", "dotfiles"), nil },
					userConfigDir: func() (string, error) { return filepath.Join("home", "config"), nil },
				}
				ctx = context.WithValue(context.Background(), command.OptsCtxKey, opts)
				err = tc.cmd.Execute(context.WithValue(ctx, command.ErrCtxKey, "plg"), nil, &bd)
			)
			if !errors.As(err, &tc.err) {
				if want, got := tc.err, err; !errors.Is(got, want) {
					t.Fatalf("want %v, got %v", want, got)
				}
			}
			if want, got := tc.want, tc.drv; !cmp.Equal(got, want) {
				t.Fatalf("(-want +got):\n%s", cmp.Diff(want, got))
			}
		})
	}
}

func testLinkSetFlags(t *testing.T) {
	allowUnexported := cmp.AllowUnexported(linkCmd{})
	testCases := []struct {
		flags map[string]string
		want  linkCmd
	}{
		{
			flags: nil,
			want:  linkCmd{},
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			var (
				cmd  linkCmd
				fset = flag.NewFlagSet("link", flag.PanicOnError)
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
					"linkCmd.SetFlags mismatch (-want +got):\n%s",
					cmp.Diff(want, got, allowUnexported),
				)
			}
		})
	}
}
