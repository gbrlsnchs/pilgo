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

func TestScan(t *testing.T) {
	testCases := []struct {
		name string
		drv  fstest.InMemoryDriver
		cmd  scanCmd
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
											Flatten: true,
											Tags:    []string{"test"},
											Options: map[string]*config.Config{
												"foo": {
													Link: "f00",
												},
											},
											Targets: []string{
												"foo",
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
			cmd: scanCmd{
				file: "",
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
											Flatten: true,
											Tags:    []string{"test"},
											Options: map[string]*config.Config{
												"foo": {
													Link: "f00",
												},
											},
											Targets: []string{
												"bar",
												"foo",
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
			name: "include",
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
									"include.yml": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data: yamlData(config.Config{
											Targets: []string{
												"foo",
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
			cmd: scanCmd{
				file: "",
				read: readMode{
					include: cliutil.MultiValueOptionSet{
						"bar": struct{}{},
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
									"include.yml": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data: yamlData(config.Config{
											Targets: []string{
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
			err: nil,
		},
		{
			name: "exclude",
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
									"exclude.yml": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data: yamlData(config.Config{
											Targets: []string{
												"foo",
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
			cmd: scanCmd{
				file: "",
				read: readMode{
					exclude: cliutil.MultiValueOptionSet{
						"foo": struct{}{},
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
									"exclude.yml": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data: yamlData(config.Config{
											Targets: []string{
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
			err: nil,
		},
		{
			name: "hidden",
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
									".git": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     []byte("foo"),
										Children: nil,
									},
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
									"hidden.yml": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data: yamlData(config.Config{
											Targets: []string{
												"foo",
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
			cmd: scanCmd{
				file: "",
				read: readMode{
					hidden: true,
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
									".git": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     []byte("foo"),
										Children: nil,
									},
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
									"hidden.yml": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data: yamlData(config.Config{
											Targets: []string{
												".git",
												"bar",
												"foo",
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
			name: "hidden exclude",
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
									".git": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     []byte("foo"),
										Children: nil,
									},
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
									"hidden_exclude.yml": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data: yamlData(config.Config{
											Targets: []string{
												"foo",
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
			cmd: scanCmd{
				file: "",
				read: readMode{
					hidden: true,
					exclude: cliutil.MultiValueOptionSet{
						".git": struct{}{},
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
									".git": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     []byte("foo"),
										Children: nil,
									},
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
									"hidden_exclude.yml": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data: yamlData(config.Config{
											Targets: []string{
												"bar",
												"foo",
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
