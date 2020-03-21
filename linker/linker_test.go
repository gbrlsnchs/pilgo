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
	t.Run("Link", testLink)
	t.Run("Resolve", testResolve)
}

func testLink(t *testing.T) {
	errLink := errors.New("Link")
	testCases := []struct {
		drv            fstest.Driver
		tr             *parser.Tree
		mkdirAllCalled bool
		mkdirAllArgs   fstest.CallStack
		symlinkCalled  bool
		symlinkArgs    fstest.CallStack
		err            error
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
			tr: &parser.Tree{Root: &parser.Node{Children: []*parser.Node{
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
			}}},
			mkdirAllCalled: true,
			mkdirAllArgs: fstest.CallStack{
				fstest.Args{"test"},
			},
			symlinkCalled: true,
			symlinkArgs: fstest.CallStack{
				fstest.Args{"foo", filepath.Join("test", "foo")},
			},
			err: nil,
		},
		{
			drv: fstest.Driver{
				StatReturn: map[string]fs.FileInfo{
					"foo": fstest.FileInfo{
						ExistsReturn: true,
					},
					filepath.Join("test", "foo"): fstest.FileInfo{
						ExistsReturn: false,
					},
					"conf": fstest.FileInfo{
						ExistsReturn: true,
					},
					"dirs": fstest.FileInfo{
						ExistsReturn: true,
					},
					filepath.Join("dirs", "conf"): fstest.FileInfo{
						ExistsReturn: false,
					},
				},
				StatErr: map[string]error{
					"foo":                        nil,
					filepath.Join("test", "foo"): nil,
				},
			},
			tr: &parser.Tree{Root: &parser.Node{Children: []*parser.Node{
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
				{
					Target: parser.File{
						BaseDir: "",
						Path:    []string{"conf"},
					},
					Link: parser.File{
						BaseDir: "dirs",
						Path:    []string{"conf"},
					},
					Children: nil,
				},
			}}},
			mkdirAllCalled: true,
			mkdirAllArgs: fstest.CallStack{
				fstest.Args{"test"},
				fstest.Args{"dirs"},
			},
			symlinkCalled: true,
			symlinkArgs: fstest.CallStack{
				fstest.Args{"foo", filepath.Join("test", "foo")},
				fstest.Args{"conf", filepath.Join("dirs", "conf")},
			},
			err: nil,
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
					},
					"conf": fstest.FileInfo{
						ExistsReturn: true,
					},
					filepath.Join("dirs", "conf"): fstest.FileInfo{
						ExistsReturn: false,
					},
				},
				StatErr: map[string]error{
					"foo":                        nil,
					filepath.Join("test", "foo"): nil,
				},
			},
			tr: &parser.Tree{Root: &parser.Node{Children: []*parser.Node{
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
				{
					Target: parser.File{
						BaseDir: "",
						Path:    []string{"conf"},
					},
					Link: parser.File{
						BaseDir: "dirs",
						Path:    []string{"conf"},
					},
					Children: nil,
				},
			}}},
			mkdirAllCalled: false,
			mkdirAllArgs:   nil,
			symlinkCalled:  false,
			symlinkArgs:    nil,
			err:            linker.ErrLinkNotExpands,
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
					},
					"conf": fstest.FileInfo{
						ExistsReturn: true,
					},
					filepath.Join("dirs", "conf"): fstest.FileInfo{
						ExistsReturn: false,
					},
				},
				StatErr: map[string]error{
					"foo":                        errLink,
					filepath.Join("test", "foo"): nil,
				},
			},
			tr: &parser.Tree{Root: &parser.Node{Children: []*parser.Node{
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
				{
					Target: parser.File{
						BaseDir: "",
						Path:    []string{"conf"},
					},
					Link: parser.File{
						BaseDir: "dirs",
						Path:    []string{"conf"},
					},
					Children: nil,
				},
			}}},
			mkdirAllCalled: false,
			mkdirAllArgs:   nil,
			symlinkCalled:  false,
			symlinkArgs:    nil,
			err:            errLink,
		},
		{
			drv: fstest.Driver{
				ReadDirReturn: map[string][]fs.FileInfo{
					"expand": {
						fstest.FileInfo{NameReturn: "expand_child", ExistsReturn: true},
					},
				},
				ReadDirErr: map[string]error{
					"expand": nil,
				},
				StatReturn: map[string]fs.FileInfo{
					// done
					"done": fstest.FileInfo{
						ExistsReturn: true,
					},
					filepath.Join("test", "done"): fstest.FileInfo{
						ExistsReturn:   true,
						LinknameReturn: "done",
					},
					// ready
					"ready": fstest.FileInfo{
						ExistsReturn: true,
					},
					filepath.Join("test", "ready"): fstest.FileInfo{
						ExistsReturn: false,
					},
					// expand and children
					"expand": fstest.FileInfo{
						ExistsReturn: true,
						IsDirReturn:  true,
					},
					filepath.Join("test", "expand"): fstest.FileInfo{
						ExistsReturn: true,
						IsDirReturn:  true,
					},
					filepath.Join("expand", "expand_child"): fstest.FileInfo{
						ExistsReturn: true,
					},
					filepath.Join("test", "expand", "expand_child"): fstest.FileInfo{
						ExistsReturn: false,
					},
				},
				StatErr: map[string]error{
					// done
					"done":                        nil,
					filepath.Join("test", "done"): nil,
					// ready
					"ready":                        nil,
					filepath.Join("test", "ready"): nil,
					// expand
					"expand":                        nil,
					filepath.Join("test", "expand"): nil,
				},
			},
			tr: &parser.Tree{Root: &parser.Node{Children: []*parser.Node{
				{
					Target: parser.File{
						BaseDir: "",
						Path:    []string{"done"},
					},
					Link: parser.File{
						BaseDir: "test",
						Path:    []string{"done"},
					},
					Children: nil,
				},
				{
					Target: parser.File{
						BaseDir: "",
						Path:    []string{"ready"},
					},
					Link: parser.File{
						BaseDir: "test",
						Path:    []string{"ready"},
					},
					Children: nil,
				},
				{
					Target: parser.File{
						BaseDir: "",
						Path:    []string{"expand"},
					},
					Link: parser.File{
						BaseDir: "test",
						Path:    []string{"expand"},
					},
					Children: nil,
				},
			}}},
			mkdirAllCalled: true,
			mkdirAllArgs: fstest.CallStack{
				fstest.Args{"test"},
				fstest.Args{filepath.Join("test", "expand")},
			},
			symlinkCalled: true,
			symlinkArgs: fstest.CallStack{
				fstest.Args{"ready", filepath.Join("test", "ready")},
				fstest.Args{
					filepath.Join("expand", "expand_child"),
					filepath.Join("test", "expand", "expand_child"),
				},
			},
			err: nil,
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			fs := fs.New(&tc.drv)
			ln := linker.New(fs)
			err := ln.Link(tc.tr)
			if want, got := tc.err, err; !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			t.Run("MkdirAll", func(t *testing.T) {
				hasBeenCalled, args := tc.drv.HasBeenCalled(tc.drv.MkdirAll)
				if want, got := tc.mkdirAllCalled, hasBeenCalled; got != want {
					t.Fatalf("want %t, got %t", want, got)
				}
				if want, got := tc.mkdirAllArgs, args; !cmp.Equal(got, want) {
					t.Fatalf("(-want +got):\n%s", cmp.Diff(want, got))
				}
			})
			t.Run("Symlink", func(t *testing.T) {
				hasBeenCalled, args := tc.drv.HasBeenCalled(tc.drv.Symlink)
				if want, got := tc.symlinkCalled, hasBeenCalled; got != want {
					t.Fatalf("want %t, got %t", want, got)
				}
				if want, got := tc.symlinkArgs, args; !cmp.Equal(got, want) {
					t.Fatalf("(-want +got):\n%s", cmp.Diff(want, got))
				}
			})
		})
	}
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
