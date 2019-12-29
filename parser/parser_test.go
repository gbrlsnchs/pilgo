package parser_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gsr.dev/pilgrim"
	"gsr.dev/pilgrim/parser"
)

func TestParser(t *testing.T) {
	t.Run("Parse", testParserParse)
}

func testParserParse(t *testing.T) {
	testCases := []struct {
		c   pilgrim.Config
		tr  *parser.Tree
		err error
	}{
		{
			c: pilgrim.Config{
				BaseDir: "test",
				Link:    nil,
				Targets: []string{
					"foo",
				},
			},
			tr: &parser.Tree{
				Root: &parser.Node{Children: []*parser.Node{
					{
						Target: parser.File{"", "foo"},
						Link:   parser.File{"test", "foo"},
					},
				}},
			},
			err: nil,
		},
		{
			c: pilgrim.Config{
				BaseDir: "test",
				Link:    nil,
				Targets: []string{
					"foo",
					"bar",
				},
			},
			tr: &parser.Tree{
				Root: &parser.Node{Children: []*parser.Node{
					{
						Target: parser.File{"", "bar"},
						Link:   parser.File{"test", "bar"},
					},
					{
						Target: parser.File{"", "foo"},
						Link:   parser.File{"test", "foo"},
					},
				}},
			},
			err: nil,
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			var p parser.Parser
			tr, err := p.Parse(tc.c)
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
