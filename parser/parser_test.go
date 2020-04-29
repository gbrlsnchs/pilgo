package parser_test

import (
	"errors"
	"os"
	"testing"

	"github.com/gbrlsnchs/pilgo/config"
	"github.com/gbrlsnchs/pilgo/internal"
	"github.com/gbrlsnchs/pilgo/parser"
	"github.com/google/go-cmp/cmp"
)

func TestParser(t *testing.T) {
	t.Run("Parse", testParserParse)
}

func testParserParse(t *testing.T) {
	os.Setenv("MY_ENV_VAR", "home")
	defer os.Unsetenv("MY_ENV_VAR")
	testCases := []struct {
		c    config.Config
		opts []parser.ParseOption
		tr   *parser.Tree
		err  error
	}{
		{
			c: config.Config{
				BaseDir: "test",
				Link:    "",
				Targets: []string{
					"foo",
				},
			},
			tr: &parser.Tree{
				Root: &parser.Node{Children: []*parser.Node{
					{
						Target: parser.File{"", []string{"foo"}},
						Link:   parser.File{"test", []string{"foo"}},
					},
				}},
			},
			err: nil,
		},
		{
			c: config.Config{
				BaseDir: "test",
				Link:    "",
				Targets: []string{
					"foo",
					"bar",
				},
			},
			tr: &parser.Tree{
				Root: &parser.Node{Children: []*parser.Node{
					{
						Target: parser.File{"", []string{"bar"}},
						Link:   parser.File{"test", []string{"bar"}},
					},
					{
						Target: parser.File{"", []string{"foo"}},
						Link:   parser.File{"test", []string{"foo"}},
					},
				}},
			},
			err: nil,
		},
		{
			c: config.Config{
				BaseDir: "test",
				Link:    "",
				Targets: []string{
					"foo",
				},
				Options: map[string]*config.Config{
					"foo": {Link: "bar"},
				},
			},
			tr: &parser.Tree{
				Root: &parser.Node{Children: []*parser.Node{
					{
						Target: parser.File{"", []string{"foo"}},
						Link:   parser.File{"test", []string{"bar"}},
					},
				}},
			},
			err: nil,
		},
		{
			c: config.Config{
				BaseDir: "test",
				Link:    "",
				Targets: []string{
					"foo",
				},
				Options: map[string]*config.Config{
					"foo": {
						Targets: []string{
							"bar",
						},
						Options: map[string]*config.Config{
							"bar": {
								Targets: []string{
									"baz",
								},
							},
						},
					},
				},
			},
			tr: &parser.Tree{
				Root: &parser.Node{Children: []*parser.Node{
					{
						Target: parser.File{"", []string{"foo"}},
						Link:   parser.File{"test", []string{"foo"}},
						Children: []*parser.Node{
							{
								Target: parser.File{
									"",
									[]string{"foo", "bar"},
								},
								Link: parser.File{
									"test",
									[]string{"foo", "bar"},
								},
								Children: []*parser.Node{
									{
										Target: parser.File{
											"",
											[]string{"foo", "bar", "baz"},
										},
										Link: parser.File{
											"test",
											[]string{"foo", "bar", "baz"},
										},
									},
								},
							},
						},
					},
				}},
			},
			err: nil,
		},
		{
			c: config.Config{
				BaseDir: "test",
				Link:    "",
				Targets: []string{
					"foo",
					"bar",
				},
			},
			opts: []parser.ParseOption{parser.Cwd("targets")},
			tr: &parser.Tree{
				Root: &parser.Node{Children: []*parser.Node{
					{
						Target: parser.File{"targets", []string{"bar"}},
						Link:   parser.File{"test", []string{"bar"}},
					},
					{
						Target: parser.File{"targets", []string{"foo"}},
						Link:   parser.File{"test", []string{"foo"}},
					},
				}},
			},
			err: nil,
		},
		{
			c: config.Config{
				Link: "",
				Targets: []string{
					"foo",
					"bar",
				},
			},
			opts: []parser.ParseOption{
				parser.BaseDirs(map[parser.Mode]string{
					parser.UserMode: "user",
					parser.HomeMode: "home",
				}),
			},
			tr: &parser.Tree{
				Root: &parser.Node{Children: []*parser.Node{
					{
						Target: parser.File{"", []string{"bar"}},
						Link:   parser.File{"user", []string{"bar"}},
					},
					{
						Target: parser.File{"", []string{"foo"}},
						Link:   parser.File{"user", []string{"foo"}},
					},
				}},
			},
			err: nil,
		},
		{
			c: config.Config{
				Link: "",
				Targets: []string{
					"foo",
					"bar",
				},
				UseHome: internal.NewBool(true),
			},
			opts: []parser.ParseOption{
				parser.BaseDirs(map[parser.Mode]string{
					parser.UserMode: "user",
					parser.HomeMode: "home",
				}),
			},
			tr: &parser.Tree{
				Root: &parser.Node{Children: []*parser.Node{
					{
						Target: parser.File{"", []string{"bar"}},
						Link:   parser.File{"home", []string{"bar"}},
					},
					{
						Target: parser.File{"", []string{"foo"}},
						Link:   parser.File{"home", []string{"foo"}},
					},
				}},
			},
			err: nil,
		},
		{
			c: config.Config{
				Link: "",
				Targets: []string{
					"test",
				},
				Options: map[string]*config.Config{
					"test": {
						Targets: []string{
							"foo",
							"bar",
						},
						UseHome: internal.NewBool(true),
					},
				},
			},
			opts: []parser.ParseOption{
				parser.BaseDirs(map[parser.Mode]string{
					parser.UserMode: "user",
					parser.HomeMode: "home",
				}),
			},
			tr: &parser.Tree{
				Root: &parser.Node{Children: []*parser.Node{
					{
						Target: parser.File{"", []string{"test"}},
						Link:   parser.File{"home", []string{"test"}},
						Children: []*parser.Node{
							{
								Target: parser.File{"", []string{"test", "bar"}},
								Link:   parser.File{"home", []string{"test", "bar"}},
							},
							{
								Target: parser.File{"", []string{"test", "foo"}},
								Link:   parser.File{"home", []string{"test", "foo"}},
							},
						},
					},
				}},
			},
			err: nil,
		},
		{
			c: config.Config{
				Link: "",
				Targets: []string{
					"foo",
					"bar",
				},
				BaseDir: "testing",
				UseHome: internal.NewBool(true), // this is overridden
			},
			opts: []parser.ParseOption{
				parser.BaseDirs(map[parser.Mode]string{
					parser.UserMode: "user",
					parser.HomeMode: "home",
				}),
			},
			tr: &parser.Tree{
				Root: &parser.Node{Children: []*parser.Node{
					{
						Target: parser.File{"", []string{"bar"}},
						Link:   parser.File{"testing", []string{"bar"}},
					},
					{
						Target: parser.File{"", []string{"foo"}},
						Link:   parser.File{"testing", []string{"foo"}},
					},
				}},
			},
			err: nil,
		},
		{
			c: config.Config{
				BaseDir: "test",
				Link:    "",
				Targets: []string{
					"foo",
				},
				Options: map[string]*config.Config{
					"foo": {
						BaseDir: "home",
						Link:    "",
						Targets: []string{
							"bar",
						},
					},
				},
			},
			tr: &parser.Tree{
				Root: &parser.Node{Children: []*parser.Node{
					{
						Target: parser.File{"", []string{"foo"}},
						Link:   parser.File{"home", []string{"foo"}},
						Children: []*parser.Node{
							{
								Target:   parser.File{"", []string{"foo", "bar"}},
								Link:     parser.File{"home", []string{"foo", "bar"}},
								Children: nil,
							},
						},
					},
				}},
			},
			err: nil,
		},
		{
			c: config.Config{
				BaseDir: "test",
				Link:    "",
				Targets: []string{
					"foo",
				},
				Options: map[string]*config.Config{
					"foo": {
						BaseDir: "home",
						Flatten: true,
						Targets: []string{
							"bar",
						},
					},
				},
			},
			tr: &parser.Tree{
				Root: &parser.Node{Children: []*parser.Node{
					{
						Target: parser.File{"", []string{"foo"}},
						Link:   parser.File{"home", []string{}},
						Children: []*parser.Node{
							{
								Target:   parser.File{"", []string{"foo", "bar"}},
								Link:     parser.File{"home", []string{"bar"}},
								Children: nil,
							},
						},
					},
				}},
			},
			err: nil,
		},
		{
			c: config.Config{
				BaseDir: "test",
				Link:    "",
				Targets: []string{
					"foo",
				},
				Options: map[string]*config.Config{
					"foo": {
						BaseDir: "home",
						Link:    "golang",
						Targets: []string{
							"bar",
						},
					},
				},
			},
			tr: &parser.Tree{
				Root: &parser.Node{Children: []*parser.Node{
					{
						Target: parser.File{"", []string{"foo"}},
						Link:   parser.File{"home", []string{"golang"}},
						Children: []*parser.Node{
							{
								Target:   parser.File{"", []string{"foo", "bar"}},
								Link:     parser.File{"home", []string{"golang", "bar"}},
								Children: nil,
							},
						},
					},
				}},
			},
			err: nil,
		},
		{
			c: config.Config{
				BaseDir: "/tmp",
				Link:    "",
				Targets: []string{
					"foo",
					"bar",
				},
			},
			opts: []parser.ParseOption{
				parser.BaseDirs(map[parser.Mode]string{
					parser.UserMode: "user",
					parser.HomeMode: "home",
				}),
			},
			tr: &parser.Tree{
				Root: &parser.Node{Children: []*parser.Node{
					{
						Target: parser.File{"", []string{"bar"}},
						Link:   parser.File{"/tmp", []string{"bar"}},
					},
					{
						Target: parser.File{"", []string{"foo"}},
						Link:   parser.File{"/tmp", []string{"foo"}},
					},
				}},
			},
			err: nil,
		},
		{
			c: config.Config{
				BaseDir: "test",
				Link:    "",
				Targets: []string{
					"foo",
				},
				Options: map[string]*config.Config{
					"foo": {
						Flatten: true,
						Targets: []string{
							"foobar",
							"footest",
						},
					},
				},
			},
			tr: &parser.Tree{
				Root: &parser.Node{Children: []*parser.Node{
					{
						Target: parser.File{"", []string{"foo"}},
						Link:   parser.File{"test", []string{}},
						Children: []*parser.Node{
							{
								Target: parser.File{"", []string{"foo", "foobar"}},
								Link:   parser.File{"test", []string{"foobar"}},
							},
							{
								Target: parser.File{"", []string{"foo", "footest"}},
								Link:   parser.File{"test", []string{"footest"}},
							},
						},
					},
				}},
			},
			err: nil,
		},
		{
			c: config.Config{
				BaseDir: "$MY_ENV_VAR",
				Link:    "",
				Targets: []string{
					"foo",
				},
			},
			opts: []parser.ParseOption{parser.Envsubst},
			tr: &parser.Tree{
				Root: &parser.Node{Children: []*parser.Node{
					{
						Target: parser.File{"", []string{"foo"}},
						Link:   parser.File{"home", []string{"foo"}},
					},
				}},
			},
			err: nil,
		},
		{
			c: config.Config{
				BaseDir: "test",
				Link:    "",
				Targets: []string{
					"$MY_ENV_VAR",
				},
			},
			opts: []parser.ParseOption{parser.Envsubst},
			tr: &parser.Tree{
				Root: &parser.Node{Children: []*parser.Node{
					{
						Target: parser.File{"", []string{"home"}},
						Link:   parser.File{"test", []string{"home"}},
					},
				}},
			},
			err: nil,
		},
		{
			c: config.Config{
				BaseDir: "test",
				Link:    "",
				Targets: []string{
					"foo",
					"bar",
				},
				Options: map[string]*config.Config{
					"bar": {
						Tags: []string{"test"},
					},
				},
			},
			tr: &parser.Tree{
				Root: &parser.Node{Children: []*parser.Node{
					{
						Target: parser.File{"", []string{"foo"}},
						Link:   parser.File{"test", []string{"foo"}},
					},
				}},
			},
			err: nil,
		},
		{
			c: config.Config{
				BaseDir: "test",
				Link:    "",
				Targets: []string{
					"foo",
					"bar",
				},
				Options: map[string]*config.Config{
					"bar": {
						Tags: []string{"test"},
					},
				},
			},
			tr: &parser.Tree{
				Root: &parser.Node{Children: []*parser.Node{
					{
						Target: parser.File{"", []string{"bar"}},
						Link:   parser.File{"test", []string{"bar"}},
					},
					{
						Target: parser.File{"", []string{"foo"}},
						Link:   parser.File{"test", []string{"foo"}},
					},
				}},
			},
			opts: []parser.ParseOption{parser.Tags(map[string]struct{}{
				"test": {},
			})},
			err: nil,
		},
		{
			c: config.Config{
				BaseDir: "test",
				Link:    "",
				Targets: []string{
					"foo",
					"bar",
				},
				Options: map[string]*config.Config{
					"bar": {
						Tags: []string{"test", "test_tag"},
					},
				},
			},
			tr: &parser.Tree{
				Root: &parser.Node{Children: []*parser.Node{
					{
						Target: parser.File{"", []string{"bar"}},
						Link:   parser.File{"test", []string{"bar"}},
					},
					{
						Target: parser.File{"", []string{"foo"}},
						Link:   parser.File{"test", []string{"foo"}},
					},
				}},
			},
			opts: []parser.ParseOption{parser.Tags(map[string]struct{}{
				"test_tag": {},
			})},
			err: nil,
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			var p parser.Parser
			tr, err := p.Parse(&tc.c, tc.opts...)
			if want, got := tc.err, err; !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			if want, got := tc.tr, tr; !cmp.Equal(got, want) {
				t.Fatalf(
					"(*Parser).Parse mismatch: (-want +got):\n%s",
					cmp.Diff(want, got),
				)
			}
		})
	}
}
