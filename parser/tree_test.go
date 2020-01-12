package parser_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gsr.dev/pilgrim/parser"
)

func TestTree(t *testing.T) {
	t.Run("String", testTreeString)
	t.Run("Walk", testTreeWalk)
}

func testTreeString(t *testing.T) {
	// TODO(gbrlsnchs): test output using golden files
	testCases := []struct {
		tr   *parser.Tree
		want string
	}{
		{
			tr: &parser.Tree{&parser.Node{
				Children: []*parser.Node{
					{
						Target: parser.File{"", []string{"foo"}},
						Link:   parser.File{"test", []string{"foo"}},
					},
					{
						Target: parser.File{"", []string{"bar"}},
						Link:   parser.File{"test", []string{"bar"}},
					},
				},
			}},
			want: `.
├── foo <- test/foo
└── bar <- test/bar
`,
		},
		{
			tr: &parser.Tree{&parser.Node{
				Children: []*parser.Node{
					{
						Target: parser.File{"", []string{"foo"}},
						Link:   parser.File{"test", []string{"foo"}},
						Children: []*parser.Node{
							{
								Target: parser.File{"", []string{"foo", "bar"}},
								Link:   parser.File{"test", []string{"foo", "bar"}},
							},
						},
					},
					{
						Target: parser.File{"", []string{"baz"}},
						Link:   parser.File{"test", []string{"baz"}},
						Children: []*parser.Node{
							{
								Target: parser.File{"", []string{"baz", "qux"}},
								Link:   parser.File{"test", []string{"baz", "qux"}},
							},
						},
					},
				},
			}},
			want: `.
├── foo     <- test/foo
│   └── bar <- test/foo/bar
└── baz     <- test/baz
    └── qux <- test/baz/qux
`,
		},
		{
			tr: &parser.Tree{&parser.Node{
				Children: []*parser.Node{
					{
						Target: parser.File{"", []string{"foo"}},
						Link:   parser.File{"test", []string{"foo"}},
						Children: []*parser.Node{
							{
								Target:   parser.File{"", []string{"foo", "bar"}},
								Link:     parser.File{"test", []string{"foo", "bar"}},
								Children: nil,
								Status:   parser.StatusDone,
							},
						},
						Status: parser.StatusSkip,
					},
					{
						Target: parser.File{"", []string{"baz"}},
						Link:   parser.File{"test", []string{"baz"}},
						Children: []*parser.Node{
							{
								Target:   parser.File{"", []string{"baz", "qux"}},
								Link:     parser.File{"test", []string{"baz", "qux"}},
								Children: nil,
								Status:   parser.StatusReady,
							},
						},
						Status: parser.StatusExpand,
					},
					{
						Target:   parser.File{"", []string{"quux"}},
						Link:     parser.File{"test", []string{"quux"}},
						Children: nil,
						Status:   parser.StatusConflict,
					},
					{
						Target:   parser.File{"", []string{"corge"}},
						Link:     parser.File{"test", []string{"corge"}},
						Children: nil,
						Status:   parser.StatusError,
					},
				},
			}},
			want: `.
├── foo                     (SKIP)
│   └── bar <- test/foo/bar (DONE)
├── baz                     (EXPAND)
│   └── qux <- test/baz/qux (READY)
├── quux    <- test/quux    (CONFLICT)
└── corge   <- test/corge   (ERROR)
`,
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			if want, got := filepath.FromSlash(tc.want), tc.tr.String(); got != want {
				t.Errorf("want\n%s\ngot\n%s", want, got)
			}
		})
	}
}

func testTreeWalk(t *testing.T) {
	testCases := []struct {
		tr   *parser.Tree
		want []*parser.Node
		err  error
	}{
		{
			tr: &parser.Tree{
				Root: &parser.Node{Children: []*parser.Node{
					{
						Target: parser.File{
							BaseDir: "",
							Path:    []string{"foo"},
						},
						Link: parser.File{
							BaseDir: "test",
							Path:    []string{"foo"},
						},
						Children: nil,
					},
				}},
			},
			want: []*parser.Node{
				{
					Target: parser.File{
						BaseDir: "",
						Path:    []string{"foo"},
					},
					Link: parser.File{
						BaseDir: "test",
						Path:    []string{"foo"},
					},
					Children: nil,
				},
			},
			err: nil,
		},
		{
			tr: &parser.Tree{
				Root: &parser.Node{Children: []*parser.Node{
					{
						Target: parser.File{
							BaseDir: "",
							Path:    []string{"foo"},
						},
						Link: parser.File{
							BaseDir: "test",
							Path:    []string{"foo"},
						},
						Children: []*parser.Node{
							{
								Target: parser.File{
									BaseDir: "",
									Path:    []string{"foo", "bar"},
								},
								Link: parser.File{
									BaseDir: "test",
									Path:    []string{"foo", "bar"},
								},
								Children: nil,
							},
						},
					},
					{
						Target: parser.File{
							BaseDir: "",
							Path:    []string{"baz"},
						},
						Link: parser.File{
							BaseDir: "test",
							Path:    []string{"baz"},
						},
						Children: nil,
					},
				}},
			},
			want: []*parser.Node{
				{
					Target: parser.File{
						BaseDir: "",
						Path:    []string{"foo"},
					},
					Link: parser.File{
						BaseDir: "test",
						Path:    []string{"foo"},
					},
					Children: []*parser.Node{
						{
							Target: parser.File{
								BaseDir: "",
								Path:    []string{"foo", "bar"},
							},
							Link: parser.File{
								BaseDir: "test",
								Path:    []string{"foo", "bar"},
							},
							Children: nil,
						},
					},
				},
				{
					Target: parser.File{
						BaseDir: "",
						Path:    []string{"foo", "bar"},
					},
					Link: parser.File{
						BaseDir: "test",
						Path:    []string{"foo", "bar"},
					},
					Children: nil,
				},
				{
					Target: parser.File{
						BaseDir: "",
						Path:    []string{"baz"},
					},
					Link: parser.File{
						BaseDir: "test",
						Path:    []string{"baz"},
					},
					Children: nil,
				},
			},
			err: nil,
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			var it iterator
			err := tc.tr.Walk(it.iterate)
			if want, got := err, tc.err; !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			if want, got := tc.want, it.argstack; !cmp.Equal(got, want) {
				t.Errorf(
					"(*Tree).Walk mismatch (-want +got):\n%s",
					cmp.Diff(want, got),
				)
			}
		})
	}
}

type iterator struct {
	argstack []*parser.Node
}

func (it *iterator) iterate(n *parser.Node) error {
	it.argstack = append(it.argstack, n)
	return nil
}
