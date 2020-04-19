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
				link:    strptr{addr: newString("f00")},
				targets: cliutil.CommaSepOptionList{
					"test",
					"testing",
					"testdata",
				},
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
											Options: map[string]config.Config{
												"foo": {
													BaseDir: "test_foo",
													Link:    newString("f00"),
													Targets: []string{
														"test",
														"testing",
														"testdata",
													},
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
				link:    strptr{addr: newString("f00")},
				targets: cliutil.CommaSepOptionList{
					"test",
					"testing",
					"testdata",
				},
				useHome: boolptr{addr: newBool(true)},
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
											Options: map[string]config.Config{
												"foo": {
													BaseDir: "test_foo",
													Link:    newString("f00"),
													Targets: []string{
														"test",
														"testing",
														"testdata",
													},
													UseHome: newBool(true),
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
											Options: map[string]config.Config{
												"test": {
													Link: newString(""),
													Targets: []string{
														"bar",
														"foo",
													},
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
				useHome: boolptr{addr: newBool(true)},
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
											Options: map[string]config.Config{
												"test": {
													Link: newString(""),
													Targets: []string{
														"bar",
														"foo",
													},
													UseHome: newBool(true),
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

// TODO(gbrlsnchs): create reusable helper
func newString(s string) *string { return &s }

func newBool(b bool) *bool { return &b }
