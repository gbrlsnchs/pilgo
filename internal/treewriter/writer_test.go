package treewriter_test

import (
	"strings"
	"testing"

	"gsr.dev/pilgrim/internal/treewriter"
)

func TestWriter(t *testing.T) {
	// TODO(gbrlsnchs): test output using golden files
	testCases := []struct {
		n     treewriter.Node
		input []byte
		want  string
	}{
		{
			n: &testNode{
				text: "root",
				nodes: []*testNode{
					{text: "foo"},
					{text: "bar"},
				},
			},
			input: nil,
			want: `. root
├── foo
└── bar
`,
		},
		{
			n: &testNode{
				text: "",
				nodes: []*testNode{
					{
						text: "foo",
						nodes: []*testNode{
							{text: "bar"},
						},
					},
					{
						text: "baz",
						nodes: []*testNode{
							{text: "qux"},
						},
					},
				},
			},
			input: nil,
			want: `.
├── foo
│   └── bar
└── baz
    └── qux
`,
		},
		{
			n: &testNode{
				text: "",
				nodes: []*testNode{
					{text: "foo"},
					{
						text: "bar",
						nodes: []*testNode{
							{text: "baz"},
						},
					},
					{
						text: "qux",
						nodes: []*testNode{
							{text: "quux"},
						},
					},
				},
			},
			input: nil,
			want: `.
├── foo
├── bar
│   └── baz
└── qux
    └── quux
`,
		},
	}
	for _, tc := range testCases {
		t.Run("Write", func(t *testing.T) {
			var bd strings.Builder
			tw := treewriter.NewWriter(&bd, tc.n)
			tw.Write(tc.input)
			if want, got := tc.want, bd.String(); got != want {
				t.Errorf("want %q, got %q", want, got)
			}
			t.Logf("\n%s", bd.String())
		})
	}
}
