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
		drv            fstest.SpyDriver
		tr             *parser.Tree
		mkdirAllCalled bool
		mkdirAllArgs   fstest.CallStack
		symlinkCalled  bool
		symlinkArgs    fstest.CallStack
		err            error
	}{
		{
			drv: fstest.SpyDriver{
				StatReturn: map[string]fs.FileInfo{
					"foo": fstest.StubFile{
						ExistsReturn: true,
					},
					filepath.Join("test", "foo"): fstest.StubFile{
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
			drv: fstest.SpyDriver{
				StatReturn: map[string]fs.FileInfo{
					"foo": fstest.StubFile{
						ExistsReturn: true,
					},
					filepath.Join("test", "foo"): fstest.StubFile{
						ExistsReturn: false,
					},
					"conf": fstest.StubFile{
						ExistsReturn: true,
					},
					"dirs": fstest.StubFile{
						ExistsReturn: true,
					},
					filepath.Join("dirs", "conf"): fstest.StubFile{
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
			drv: fstest.SpyDriver{
				StatReturn: map[string]fs.FileInfo{
					"foo": fstest.StubFile{
						ExistsReturn: true,
						IsDirReturn:  true,
					},
					filepath.Join("test", "foo"): fstest.StubFile{
						ExistsReturn: true,
					},
					"conf": fstest.StubFile{
						ExistsReturn: true,
					},
					filepath.Join("dirs", "conf"): fstest.StubFile{
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
			drv: fstest.SpyDriver{
				StatReturn: map[string]fs.FileInfo{
					"foo": fstest.StubFile{
						ExistsReturn: true,
						IsDirReturn:  true,
					},
					filepath.Join("test", "foo"): fstest.StubFile{
						ExistsReturn: true,
					},
					"conf": fstest.StubFile{
						ExistsReturn: true,
					},
					filepath.Join("dirs", "conf"): fstest.StubFile{
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
			drv: fstest.SpyDriver{
				ReadDirReturn: map[string][]fs.FileInfo{
					"expand": {
						fstest.StubFile{NameReturn: "expand_child", ExistsReturn: true},
					},
				},
				ReadDirErr: map[string]error{
					"expand": nil,
				},
				StatReturn: map[string]fs.FileInfo{
					// done
					"done": fstest.StubFile{
						ExistsReturn: true,
					},
					filepath.Join("test", "done"): fstest.StubFile{
						ExistsReturn:   true,
						LinknameReturn: "done",
					},
					// ready
					"ready": fstest.StubFile{
						ExistsReturn: true,
					},
					filepath.Join("test", "ready"): fstest.StubFile{
						ExistsReturn: false,
					},
					// expand and children
					"expand": fstest.StubFile{
						ExistsReturn: true,
						IsDirReturn:  true,
					},
					filepath.Join("test", "expand"): fstest.StubFile{
						ExistsReturn: true,
						IsDirReturn:  true,
					},
					filepath.Join("expand", "expand_child"): fstest.StubFile{
						ExistsReturn: true,
					},
					filepath.Join("test", "expand", "expand_child"): fstest.StubFile{
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
		drv       fstest.SpyDriver
		n         *parser.Node
		err       error
		conflicts []error
		want      *parser.Node
	}{
		{
			drv: fstest.SpyDriver{
				StatReturn: map[string]fs.FileInfo{
					"foo": fstest.StubFile{
						ExistsReturn: true,
					},
					filepath.Join("test", "foo"): fstest.StubFile{
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
			drv: fstest.SpyDriver{
				StatReturn: map[string]fs.FileInfo{
					// targets
					"foo": fstest.StubFile{
						ExistsReturn: true,
					},
					filepath.Join("foo", "bar"): fstest.StubFile{
						ExistsReturn: true,
					},
					// links
					filepath.Join("test", "foo"): fstest.StubFile{
						ExistsReturn: false,
					},
					filepath.Join("test", "foo", "bar"): fstest.StubFile{
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
						Status: parser.StatusReady,
					},
				},
				Status: parser.StatusSkip,
			},
		},
		{
			drv: fstest.SpyDriver{
				StatReturn: map[string]fs.FileInfo{
					"foo": fstest.StubFile{
						ExistsReturn: true,
					},
					filepath.Join("test", "foo"): fstest.StubFile{
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
			drv: fstest.SpyDriver{
				StatReturn: map[string]fs.FileInfo{
					"foo": fstest.StubFile{
						ExistsReturn: true,
					},
					filepath.Join("test", "foo"): fstest.StubFile{
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
			conflicts: []error{linker.ErrLinkExists},
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
			drv: fstest.SpyDriver{
				StatReturn: map[string]fs.FileInfo{
					"foo": fstest.StubFile{
						ExistsReturn: true,
						IsDirReturn:  true,
					},
					filepath.Join("test", "foo"): fstest.StubFile{
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
			conflicts: []error{linker.ErrLinkNotExpands},
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
			drv: fstest.SpyDriver{
				StatReturn: map[string]fs.FileInfo{
					"foo": fstest.StubFile{
						ExistsReturn: false,
					},
					filepath.Join("test", "foo"): fstest.StubFile{
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
			conflicts: []error{linker.ErrTargetNotExists},
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
			drv: fstest.SpyDriver{
				StatReturn: map[string]fs.FileInfo{
					// targets
					"foo": fstest.StubFile{
						ExistsReturn: true,
						IsDirReturn:  true,
					},
					filepath.Join("foo", "bar"): fstest.StubFile{
						ExistsReturn: true,
					},
					filepath.Join("foo", "baz"): fstest.StubFile{
						ExistsReturn: true,
					},
					filepath.Join("foo", "qux"): fstest.StubFile{
						ExistsReturn: true,
					},
					// links
					filepath.Join("test", "foo"): fstest.StubFile{
						ExistsReturn: true,
						IsDirReturn:  true,
					},
					filepath.Join("test", "foo", "bar"): fstest.StubFile{
						ExistsReturn: false,
					},
					filepath.Join("test", "foo", "bar"): fstest.StubFile{
						ExistsReturn: false,
					},
					filepath.Join("test", "foo", "baz"): fstest.StubFile{
						ExistsReturn: false,
					},
					filepath.Join("test", "foo", "qux"): fstest.StubFile{
						ExistsReturn: false,
					},
				},
				StatErr: map[string]error{
					"foo":                        nil,
					filepath.Join("test", "foo"): nil,
				},
				ReadDirReturn: map[string][]fs.FileInfo{
					"foo": {
						fstest.StubFile{NameReturn: "bar"},
						fstest.StubFile{NameReturn: "baz"},
						fstest.StubFile{NameReturn: "qux"},
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
						Status:   parser.StatusReady,
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
						Status:   parser.StatusReady,
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
						Status:   parser.StatusReady,
					},
				},
				Status: parser.StatusExpand,
			},
		},
		{
			drv: fstest.SpyDriver{
				StatReturn: map[string]fs.FileInfo{
					"foo": fstest.StubFile{
						ExistsReturn: true,
						IsDirReturn:  false,
					},
					filepath.Join("test", "foo"): fstest.StubFile{
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
			conflicts: []error{linker.ErrTargetNotExpands},
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
			drv: fstest.SpyDriver{
				StatReturn: map[string]fs.FileInfo{
					"foo": fstest.StubFile{
						ExistsReturn: true,
					},
					"test": fstest.StubFile{
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
			tr := &parser.Tree{
				Root: &parser.Node{Children: []*parser.Node{tc.n}},
			}
			err := ln.Resolve(tr)
			// TODO(gbrlsnchs): check error message has correct file path
			var cft *linker.ConflictError
			if errors.As(err, &cft) {
				if want, got := len(tc.conflicts), len(cft.Errs); got != want {
					t.Fatalf("want %d, got %d", want, got)
				}
				for i, err := range cft.Errs {
					t.Run("conflict", func(t *testing.T) {
						if want, got := tc.conflicts[i], err; !errors.Is(got, want) {
							t.Fatalf("want %v, got %v", want, got)
						}
					})
				}
			} else {
				if want, got := tc.err, err; !errors.Is(got, want) {
					t.Fatalf("want %v, got %v", want, got)
				}
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
