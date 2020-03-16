package linker_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gsr.dev/pilgrim/fs"
	"gsr.dev/pilgrim/fs/fstest"
	"gsr.dev/pilgrim/linker"
	"gsr.dev/pilgrim/parser"
)

func TestLinker(t *testing.T) {
	t.Run("Resolve", testResolve)
}

func testResolve(t *testing.T) {
	testCases := []struct {
		drv  fstest.Driver
		n    *parser.Node
		err  error
		want *parser.Node
	}{
		{
			drv: fstest.Driver{
				StatReturn: map[string]fs.FileInfo{
					"foo": fstest.FileInfo{
						ExistsReturn: true,
					},
					filepath.Join("test", "foo"): fstest.FileInfo{
						ExistsReturn: false,
					},
				},
				StatErr: map[string]error{
					"foo":                        nil,
					filepath.Join("test", "foo"): nil,
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
			drv: fstest.Driver{
				StatReturn: map[string]fs.FileInfo{
					"foo": fstest.FileInfo{
						ExistsReturn: true,
					},
					filepath.Join("test", "foo"): fstest.FileInfo{
						ExistsReturn: true,
					},
				},
				StatErr: map[string]error{
					"foo":                        nil,
					filepath.Join("test", "foo"): nil,
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
			drv: fstest.Driver{
				StatReturn: map[string]fs.FileInfo{
					"foo": fstest.FileInfo{
						ExistsReturn: true,
					},
					filepath.Join("test", "foo"): fstest.FileInfo{
						ExistsReturn:   true,
						LinknameReturn: filepath.Join("", "foo"),
					},
				},
				StatErr: map[string]error{
					"foo":                        nil,
					filepath.Join("test", "foo"): nil,
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
			drv: fstest.Driver{
				StatReturn: map[string]fs.FileInfo{
					"foo": fstest.FileInfo{
						ExistsReturn: true,
					},
					filepath.Join("test", "foo"): fstest.FileInfo{
						ExistsReturn:   true,
						LinknameReturn: filepath.Join("test", "bar"),
					},
				},
				StatErr: map[string]error{
					"foo":                        nil,
					filepath.Join("test", "bar"): nil,
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
			err: linker.ErrLinkExists,
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
			drv: fstest.Driver{
				StatReturn: map[string]fs.FileInfo{
					"foo": fstest.FileInfo{
						ExistsReturn: true,
						IsDirReturn:  true,
					},
					filepath.Join("test", "foo"): fstest.FileInfo{
						ExistsReturn: true,
						IsDirReturn:  false,
					},
				},
				StatErr: map[string]error{
					"foo":                        nil,
					filepath.Join("test", "foo"): nil,
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
			err: linker.ErrLinkNotExpands,
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
			drv: fstest.Driver{
				StatReturn: map[string]fs.FileInfo{
					"foo": fstest.FileInfo{
						ExistsReturn: false,
					},
					filepath.Join("test", "foo"): fstest.FileInfo{
						ExistsReturn: false,
					},
				},
				StatErr: map[string]error{
					"foo":                        nil,
					filepath.Join("test", "foo"): nil,
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
			err: linker.ErrTargetNotExists,
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
			drv: fstest.Driver{
				StatReturn: map[string]fs.FileInfo{
					"foo": fstest.FileInfo{
						ExistsReturn: true,
						IsDirReturn:  true,
					},
					filepath.Join("test", "foo"): fstest.FileInfo{
						ExistsReturn: true,
						IsDirReturn:  true,
					},
				},
				StatErr: map[string]error{
					"foo":                        nil,
					filepath.Join("test", "foo"): nil,
				},
				ReadDirReturn: map[string][]fs.FileInfo{
					"foo": {
						fstest.FileInfo{NameReturn: "bar"},
						fstest.FileInfo{NameReturn: "baz"},
						fstest.FileInfo{NameReturn: "qux"},
					},
					filepath.Join("test", "foo"): nil,
				},
				ReadDirErr: map[string]error{
					"foo":                        nil,
					filepath.Join("test", "foo"): nil,
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
			drv: fstest.Driver{
				StatReturn: map[string]fs.FileInfo{
					"foo": fstest.FileInfo{
						ExistsReturn: true,
						IsDirReturn:  false,
					},
					filepath.Join("test", "foo"): fstest.FileInfo{
						ExistsReturn: true,
						IsDirReturn:  true,
					},
				},
				StatErr: map[string]error{
					"foo":                        nil,
					filepath.Join("test", "foo"): nil,
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
			err: linker.ErrTargetNotExpands,
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
			drv: fstest.Driver{
				StatReturn: map[string]fs.FileInfo{
					"foo": fstest.FileInfo{
						ExistsReturn: true,
					},
					"test": fstest.FileInfo{
						ExistsReturn: true,
					},
				},
				StatErr: map[string]error{
					"foo":  nil,
					"test": nil,
				},
			},
			n: &parser.Node{
				Target: parser.File{
					BaseDir: "",
					Path:    []string{"foo"},
				},
				Link: parser.File{
					BaseDir: "test",
					Path:    []string{},
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
					Path:    []string{},
				},
				Children: nil,
				Status:   parser.StatusSkip,
			},
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			ln := linker.New(fs.New(&tc.drv))
			err := ln.Resolve(tc.n)
			// TODO(gbrlsnchs): check error message has correct file path
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
