package parser_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gsr.dev/pilgrim/parser"
)

func TestTree(t *testing.T) {
	t.Run("Walk", testTreeWalk)
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
