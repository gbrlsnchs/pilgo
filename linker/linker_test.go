package linker_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gsr.dev/pilgrim/linker"
	"gsr.dev/pilgrim/parser"
)

func TestLinker(t *testing.T) {
	t.Run("Resolve", testResolve)
}

func testResolve(t *testing.T) {
	testCases := []struct {
		fs   linker.FileSystem
		n    *parser.Node
		err  error
		want *parser.Node
	}{
		{
			fs: testFileSystem{
				info: map[string]infoReturn{
					filepath.Join("test", "foo"): {
						returnValue: testFileInfo{
							exists: false,
						},
						err: nil,
					},
				},
			},
			n: &parser.Node{
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
			err: nil,
			want: &parser.Node{
				Target: parser.File{
					BaseDir: "",
					Path:    []string{"foo"},
				},
				Link: parser.File{
					BaseDir: "test",
					Path:    []string{"foo"},
				},
				Children: nil,
				Status:   parser.StatusReady,
			},
		},
		{
			fs: testFileSystem{
				info: map[string]infoReturn{
					filepath.Join("test", "foo"): {
						returnValue: testFileInfo{
							exists: true,
						},
						err: nil,
					},
				},
			},
			n: &parser.Node{
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
					},
				},
			},
			err: nil,
			want: &parser.Node{
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
					},
				},
				Status: parser.StatusSkip,
			},
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			ln := linker.New(tc.fs)
			err := ln.Resolve(tc.n)
			if want, got := tc.err, err; !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			if want, got := tc.want, tc.n; !cmp.Equal(got, want) {
				t.Errorf(
					"(*Linker).Resolve mismatch (-want +got):\n%s",
					cmp.Diff(want, got),
				)
			}
		})
	}
}
