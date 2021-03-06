package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/gbrlsnchs/cli/clitest"
	"github.com/gbrlsnchs/cli/cliutil"
	"github.com/gbrlsnchs/pilgo/config"
	"github.com/gbrlsnchs/pilgo/fs/fstest"
	"github.com/gbrlsnchs/pilgo/internal"
	"github.com/google/go-cmp/cmp"
)

func TestConfig(t *testing.T) {
	testCases := []struct {
		name string
		drv  fstest.InMemoryDriver
		cmd  configCmd
		want fstest.InMemoryDriver
		err  error
	}{
		{
			name: "default",
			drv: fstest.InMemoryDriver{
				CurrentDir: "home/dotfiles",
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
									"foo": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     []byte("foo"),
										Children: nil,
									},
									"bar": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     []byte("bar"),
										Children: nil,
									},
									"default.yml": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data: yamlData(config.Config{
											BaseDir: "test",
											Targets: []string{
												"foo",
												"bar",
											},
										}),
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
			cmd: configCmd{
				file:    "foo",
				baseDir: "test_foo",
				link:    "f00",
				tags:    cliutil.CommaSepOptionList{"test", "tag"},
			},
			want: fstest.InMemoryDriver{
				CurrentDir: "home/dotfiles",
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
									"foo": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     []byte("foo"),
										Children: nil,
									},
									"bar": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     []byte("bar"),
										Children: nil,
									},
									"default.yml": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data: yamlData(config.Config{
											BaseDir: "test",
											Targets: []string{
												"foo",
												"bar",
											},
											Options: map[string]*config.Config{
												"foo": {
													BaseDir: "test_foo",
													Link:    "f00",
													Targets: nil,
													Tags:    []string{"test", "tag"},
												},
											},
										}),
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
			err: nil,
		},
		{
			name: "home",
			drv: fstest.InMemoryDriver{
				CurrentDir: "home/dotfiles",
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
									"foo": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     []byte("foo"),
										Children: nil,
									},
									"bar": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     []byte("bar"),
										Children: nil,
									},
									"home.yml": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data: yamlData(config.Config{
											BaseDir: "test",
											Targets: []string{
												"foo",
												"bar",
											},
										}),
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
			cmd: configCmd{
				file:    "foo",
				baseDir: "test_foo",
				link:    "f00",
				useHome: boolptr{addr: internal.NewBool(true)},
			},
			want: fstest.InMemoryDriver{
				CurrentDir: "home/dotfiles",
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
									"foo": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     []byte("foo"),
										Children: nil,
									},
									"bar": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     []byte("bar"),
										Children: nil,
									},
									"home.yml": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data: yamlData(config.Config{
											BaseDir: "test",
											Targets: []string{
												"foo",
												"bar",
											},
											Options: map[string]*config.Config{
												"foo": {
													BaseDir: "test_foo",
													Link:    "f00",
													Targets: nil,
													UseHome: internal.NewBool(true),
												},
											},
										}),
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
			err: nil,
		},
		{
			name: "flatten",
			drv: fstest.InMemoryDriver{
				CurrentDir: "home/dotfiles",
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
										Children: map[string]fstest.File{
											"foo": {
												Linkname: "",
												Perm:     os.ModePerm,
												Data:     []byte("bar"),
												Children: nil,
											},
											"bar": {
												Linkname: "",
												Perm:     os.ModePerm,
												Data:     []byte("bar"),
												Children: nil,
											},
										},
									},
									"flatten.yml": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data: yamlData(config.Config{
											BaseDir: "test",
											Targets: []string{
												"test",
											},
										}),
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
			cmd: configCmd{
				file:    "test",
				flatten: true,
			},
			want: fstest.InMemoryDriver{
				CurrentDir: "home/dotfiles",
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
										Children: map[string]fstest.File{
											"foo": {
												Linkname: "",
												Perm:     os.ModePerm,
												Data:     []byte("bar"),
												Children: nil,
											},
											"bar": {
												Linkname: "",
												Perm:     os.ModePerm,
												Data:     []byte("bar"),
												Children: nil,
											},
										},
									},
									"flatten.yml": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data: yamlData(config.Config{
											BaseDir: "test",
											Targets: []string{
												"test",
											},
											Options: map[string]*config.Config{
												"test": {
													Flatten: true,
												},
											},
										}),
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
			err: nil,
		},
		{
			name: "flatten to home",
			drv: fstest.InMemoryDriver{
				CurrentDir: "home/dotfiles",
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
										Children: map[string]fstest.File{
											"foo": {
												Linkname: "",
												Perm:     os.ModePerm,
												Data:     []byte("bar"),
												Children: nil,
											},
											"bar": {
												Linkname: "",
												Perm:     os.ModePerm,
												Data:     []byte("bar"),
												Children: nil,
											},
										},
									},
									"flatten_to_home.yml": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data: yamlData(config.Config{
											BaseDir: "test",
											Targets: []string{
												"test",
											},
										}),
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
			cmd: configCmd{
				file:    "test",
				useHome: boolptr{addr: internal.NewBool(true)},
				flatten: true,
			},
			want: fstest.InMemoryDriver{
				CurrentDir: "home/dotfiles",
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
										Children: map[string]fstest.File{
											"foo": {
												Linkname: "",
												Perm:     os.ModePerm,
												Data:     []byte("bar"),
												Children: nil,
											},
											"bar": {
												Linkname: "",
												Perm:     os.ModePerm,
												Data:     []byte("bar"),
												Children: nil,
											},
										},
									},
									"flatten_to_home.yml": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data: yamlData(config.Config{
											BaseDir: "test",
											Targets: []string{
												"test",
											},
											Options: map[string]*config.Config{
												"test": {
													Flatten: true,
													UseHome: internal.NewBool(true),
												},
											},
										}),
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
			err: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				appcfg = appConfig{
					conf: filepath.Base(t.Name()) + ".yml",
					fs:   &tc.drv,
				}
				exec = tc.cmd.register(appcfg.copy)
				prg  = clitest.NewProgram("config")
				err  = exec(prg)
			)
			if want, got := tc.err, err; !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			if want, got := "", prg.Output(); got != want {
				t.Fatalf("\"config\" command output mismatch (-want +got):\n%s",
					cmp.Diff(want, got))
			}
			if want, got := tc.want, tc.drv; !cmp.Equal(got, want) {
				t.Fatalf("\"config\" command has unintended effects in the file system: (-want +got):\n%s",
					cmp.Diff(want, got))
			}
		})
	}
}
