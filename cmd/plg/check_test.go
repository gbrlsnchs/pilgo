package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/gbrlsnchs/cli/clitest"
	"github.com/gbrlsnchs/pilgo/config"
	"github.com/gbrlsnchs/pilgo/fs/fstest"
	"github.com/gbrlsnchs/pilgo/linker"
	"github.com/google/go-cmp/cmp"
)

func TestCheck(t *testing.T) {
	os.Setenv("MY_ENV_VAR", "check.txt")
	defer os.Unsetenv("MY_ENV_VAR")
	testCases := []struct {
		name      string
		drv       fstest.InMemoryDriver
		cmd       checkCmd
		want      fstest.InMemoryDriver
		conflicts bool
		err       error
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
									"test": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     []byte("foo"),
										Children: nil,
									},
									"check.txt": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     []byte("bar"),
										Children: nil,
									},
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
			cmd: checkCmd{},
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
										Children: nil,
									},
									"check.txt": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     []byte("bar"),
										Children: nil,
									},
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
			conflicts: false,
			err:       nil,
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
									"test": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     []byte("foo"),
										Children: nil,
									},
									"check.txt": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     []byte("bar"),
										Children: nil,
									},
									"home.yml": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data: yamlData(config.Config{
											Targets: []string{
												"$MY_ENV_VAR",
												"test",
											},
											UseHome: newBool(true),
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
			cmd: checkCmd{},
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
										Children: nil,
									},
									"check.txt": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     []byte("bar"),
										Children: nil,
									},
									"home.yml": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data: yamlData(config.Config{
											Targets: []string{
												"$MY_ENV_VAR",
												"test",
											},
											UseHome: newBool(true),
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
			conflicts: false,
			err:       nil,
		},
		{
			name: "fail",
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
										Children: nil,
									},
									"check.txt": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     []byte("bar"),
										Children: nil,
									},
									"fail.yml": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data: yamlData(config.Config{
											Targets: []string{
												"$MY_ENV_VAR",
												"nonexistent",
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
			cmd: checkCmd{fail: true},
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
										Children: nil,
									},
									"check.txt": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data:     []byte("bar"),
										Children: nil,
									},
									"fail.yml": {
										Linkname: "",
										Perm:     os.ModePerm,
										Data: yamlData(config.Config{
											Targets: []string{
												"$MY_ENV_VAR",
												"nonexistent",
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
			conflicts: true,
			err:       nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				appcfg = appConfig{
					conf:          tc.name + ".yml",
					fs:            &tc.drv,
					getwd:         func() (string, error) { return fstest.AbsPath("home", "dotfiles"), nil },
					userConfigDir: func() (string, error) { return fstest.AbsPath("home", "config"), nil },
					userHomeDir:   func() (string, error) { return fstest.AbsPath("home"), nil },
				}
				exec = tc.cmd.register(appcfg.copy)
				prg  = clitest.NewProgram("check")
				err  = exec(prg)
			)
			var rcv *linker.ConflictError
			conflicts := errors.As(err, &rcv)
			if !conflicts {
				if want, got := tc.err, err; !errors.Is(got, want) {
					t.Fatalf("want %v, got %v", want, got)
				}
			} else {
				if want, got := tc.conflicts, conflicts; got != want {
					t.Fatalf("want %t, got %t", want, got)
				}
			}
			outputs := []string{"stdout", "stderr", "combined"}
			gots := map[string]string{
				"stdout":   prg.Output(),
				"stderr":   prg.ErrOutput(),
				"combined": prg.CombinedOutput(),
			}
			for _, out := range outputs {
				t.Run(out, func(t *testing.T) {
					golden := filepath.Join(testdir, t.Name()+".golden")
					b, err := readFile(golden)
					if err != nil {
						t.Log(err) // err means output should be empty
					}
					if want, got := string(b), gots[out]; got != want {
						t.Fatalf("\"check\" command %s output mismatch (-want +got):\n%s",
							out,
							cmp.Diff(want, got))
					}
				})
			}
			if want, got := tc.want, tc.drv; !cmp.Equal(got, want) {
				t.Fatalf("\"check\" command has unintended effects in the file system: (-want +got):\n%s",
					cmp.Diff(want, got))
			}
		})
	}
}
