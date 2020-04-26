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

func TestShow(t *testing.T) {
	os.Setenv("MY_ENV_VAR", "show.txt")
	defer os.Unsetenv("MY_ENV_VAR")
	testCases := []struct {
		name string
		drv  fstest.InMemoryDriver
		want fstest.InMemoryDriver
		cmd  showCmd
		tags cliutil.CommaSepOptionSet
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
									"default.yml": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data: yamlData(config.Config{
											Targets: []string{
												"$MY_ENV_VAR",
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
									"default.yml": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data: yamlData(config.Config{
											Targets: []string{
												"$MY_ENV_VAR",
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
			cmd: showCmd{},
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
									"home.yml": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data: yamlData(config.Config{
											Targets: []string{
												"$MY_ENV_VAR",
												"test",
											},
											UseHome: internal.NewBool(true),
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
									"home.yml": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data: yamlData(config.Config{
											Targets: []string{
												"$MY_ENV_VAR",
												"test",
											},
											UseHome: internal.NewBool(true),
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
			cmd: showCmd{},
			err: nil,
		},
		{
			name: "tags exclude",
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
									"tags_exclude.yml": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data: yamlData(config.Config{
											Targets: []string{
												"$MY_ENV_VAR",
												"test",
												"foo",
											},
											Options: map[string]*config.Config{
												"test": {
													Tags: []string{"test"},
												},
												"foo": {
													Tags: []string{"test", "foo"},
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
									"tags_exclude.yml": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data: yamlData(config.Config{
											Targets: []string{
												"$MY_ENV_VAR",
												"test",
												"foo",
											},
											Options: map[string]*config.Config{
												"test": {
													Tags: []string{"test"},
												},
												"foo": {
													Tags: []string{"test", "foo"},
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
			cmd:  showCmd{},
			tags: nil,
			err:  nil,
		},
		{
			name: "tags include multi",
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
									"tags_include_multi.yml": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data: yamlData(config.Config{
											Targets: []string{
												"$MY_ENV_VAR",
												"test",
												"foo",
											},
											Options: map[string]*config.Config{
												"test": {
													Tags: []string{"test"},
												},
												"foo": {
													Tags: []string{"test", "foo"},
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
									"tags_include_multi.yml": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data: yamlData(config.Config{
											Targets: []string{
												"$MY_ENV_VAR",
												"test",
												"foo",
											},
											Options: map[string]*config.Config{
												"test": {
													Tags: []string{"test"},
												},
												"foo": {
													Tags: []string{"test", "foo"},
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
			cmd: showCmd{},
			tags: cliutil.CommaSepOptionSet{
				"test": struct{}{},
			},
			err: nil,
		},
		{
			name: "tags include single",
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
									"tags_include_single.yml": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data: yamlData(config.Config{
											Targets: []string{
												"$MY_ENV_VAR",
												"test",
												"foo",
											},
											Options: map[string]*config.Config{
												"test": {
													Tags: []string{"test"},
												},
												"foo": {
													Tags: []string{"test", "foo"},
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
									"tags_include_single.yml": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data: yamlData(config.Config{
											Targets: []string{
												"$MY_ENV_VAR",
												"test",
												"foo",
											},
											Options: map[string]*config.Config{
												"test": {
													Tags: []string{"test"},
												},
												"foo": {
													Tags: []string{"test", "foo"},
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
			cmd: showCmd{},
			tags: cliutil.CommaSepOptionSet{
				"foo": struct{}{},
			},
			err: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				appcfg = appConfig{
					conf:          filepath.Base(t.Name()) + ".yml",
					fs:            &tc.drv,
					getwd:         func() (string, error) { return fstest.AbsPath("home", "dotfiles"), nil },
					userConfigDir: func() (string, error) { return fstest.AbsPath("home", "config"), nil },
					userHomeDir:   func() (string, error) { return fstest.AbsPath("home"), nil },
					tags:          tc.tags,
				}
				exec = tc.cmd.register(appcfg.copy)
				prg  = clitest.NewProgram("show")
				err  = exec(prg)
			)
			if want, got := tc.err, err; !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			golden := filepath.Join("testdata", t.Name()) + ".golden"
			b, err := readFile(golden)
			if err != nil {
				t.Fatal(err)
			}
			if want, got := string(b), prg.Output(); got != want {
				t.Errorf("\"show\" command stdout mismatch (-want +got):\n%s",
					cmp.Diff(want, got))
			}
			if want, got := string(b), prg.CombinedOutput(); got != want {
				t.Errorf("\"show\" command combined output mismatch (-want +got):\n%s",
					cmp.Diff(want, got))
			}
			if want, got := "", prg.ErrOutput(); got != want {
				t.Errorf("\"show\" command stderr mismatch (-want +got):\n%s",
					cmp.Diff(want, got))
			}
			if want, got := tc.want, tc.drv; !cmp.Equal(got, want) {
				t.Fatalf("\"show\" command has unintended effects in the file system: (-want +got):\n%s",
					cmp.Diff(want, got))
			}
		})
	}
}
