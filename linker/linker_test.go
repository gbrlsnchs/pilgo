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
					"foo": {
						returnValue: testFileInfo{
							exists: true,
						},
						err: nil,
					},
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
					"foo": {
						returnValue: testFileInfo{
							exists: true,
						},
						err: nil,
					},
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
		{
			fs: testFileSystem{
				info: map[string]infoReturn{
					"foo": {
						returnValue: testFileInfo{
							exists: true,
						},
						err: nil,
					},
					filepath.Join("test", "foo"): {
						returnValue: testFileInfo{
							exists:   true,
							linkname: filepath.Join("test", "foo"),
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
				Status:   parser.StatusDone,
			},
		},
		{
			fs: testFileSystem{
				info: map[string]infoReturn{
					"foo": {
						returnValue: testFileInfo{
							exists: true,
						},
						err: nil,
					},
					filepath.Join("test", "foo"): {
						returnValue: testFileInfo{
							exists:   true,
							linkname: filepath.Join("test", "bar"),
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
				Status:   parser.StatusConflict,
			},
		},
		{
			fs: testFileSystem{
				info: map[string]infoReturn{
					"foo": {
						returnValue: testFileInfo{
							exists: true,
						},
						err: nil,
					},
					filepath.Join("test", "foo"): {
						returnValue: testFileInfo{
							exists: true,
							isDir:  false,
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
				Status:   parser.StatusConflict,
			},
		},
		{
			fs: testFileSystem{
				info: map[string]infoReturn{
					"foo": {
						returnValue: testFileInfo{
							exists: false,
						},
						err: nil,
					},
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
				Status:   parser.StatusError,
			},
		},
		{
			fs: testFileSystem{
				info: map[string]infoReturn{
					"foo": {
						returnValue: testFileInfo{
							exists: true,
							isDir:  true,
						},
						err: nil,
					},
					filepath.Join("test", "foo"): {
						returnValue: testFileInfo{
							exists: true,
							isDir:  true,
						},
						err: nil,
					},
				},
				readDir: map[string]readDirReturn{
					"foo": readDirReturn{
						returnValue: []string{
							"bar",
							"baz",
							"qux",
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
				Children: []*parser.Node{
					{
						Target: parser.File{
							BaseDir: "",
							Path: []string{
								"foo",
								"bar",
							},
						},
						Link: parser.File{
							BaseDir: "test",
							Path: []string{
								"foo",
								"bar",
							},
						},
						Children: nil,
					},
					{
						Target: parser.File{
							BaseDir: "",
							Path: []string{
								"foo",
								"baz",
							},
						},
						Link: parser.File{
							BaseDir: "test",
							Path: []string{
								"foo",
								"baz",
							},
						},
						Children: nil,
					},
					{
						Target: parser.File{
							BaseDir: "",
							Path: []string{
								"foo",
								"qux",
							},
						},
						Link: parser.File{
							BaseDir: "test",
							Path: []string{
								"foo",
								"qux",
							},
						},
						Children: nil,
					},
				},
				Status: parser.StatusExpand,
			},
		},
		{
			fs: testFileSystem{
				info: map[string]infoReturn{
					"foo": {
						returnValue: testFileInfo{
							exists: true,
							isDir:  false,
						},
						err: nil,
					},
					filepath.Join("test", "foo"): {
						returnValue: testFileInfo{
							exists: true,
							isDir:  true,
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
				Status:   parser.StatusConflict,
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
