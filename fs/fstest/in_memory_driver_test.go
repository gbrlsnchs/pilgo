package fstest_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gbrlsnchs/pilgo/fs"
	"github.com/gbrlsnchs/pilgo/fs/fstest"
	"github.com/google/go-cmp/cmp"
)

var _ fs.Driver = new(fstest.InMemoryDriver)

func TestAbsPath(t *testing.T) {
	want := fmt.Sprintf("~%s", filepath.Join("foo", "bar"))
	got := fstest.AbsPath("foo", "bar")
	if got != want {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func TestInMemoryDriver(t *testing.T) {
	t.Run("MkdirAll", testInMemoryDriverMkdirAll)
	t.Run("ReadDir", testInMemoryDriverReadDir)
	t.Run("ReadFile", testInMemoryDriverReadFile)
	t.Run("Stat", testInMemoryDriverStat)
	t.Run("Symlink", testInMemoryDriverSymlink)
	t.Run("WriteFile", testInMemoryDriverWriteFile)
}

func testInMemoryDriverMkdirAll(t *testing.T) {
	testCases := []struct {
		drv     fstest.InMemoryDriver
		dirname string
		want    fstest.InMemoryDriver
		err     error
	}{
		{
			drv:     fstest.InMemoryDriver{},
			dirname: "foo",
			want: fstest.InMemoryDriver{
				Files: map[string]fstest.File{
					"foo": fstest.File{
						Linkname: "",
						Perm:     os.ModePerm,
						Children: make(map[string]fstest.File, 0),
					},
				},
			},
			err: nil,
		},
		{
			drv: fstest.InMemoryDriver{
				Files: map[string]fstest.File{
					"foo": fstest.File{
						Linkname: "",
						Perm:     os.ModePerm,
						Children: make(map[string]fstest.File, 0),
					},
				},
			},
			dirname: filepath.Join("foo", "bar"),
			want: fstest.InMemoryDriver{
				Files: map[string]fstest.File{
					"foo": fstest.File{
						Linkname: "",
						Perm:     os.ModePerm,
						Children: map[string]fstest.File{
							"bar": {
								Linkname: "",
								Perm:     os.ModePerm,
								Children: make(map[string]fstest.File, 0),
							},
						},
					},
				},
			},
			err: nil,
		},
		{
			drv: fstest.InMemoryDriver{
				CurrentDir: "foo",
				Files: map[string]fstest.File{
					"foo": fstest.File{
						Linkname: "",
						Perm:     os.ModePerm,
						Children: make(map[string]fstest.File, 0),
					},
				},
			},
			dirname: "bar",
			want: fstest.InMemoryDriver{
				CurrentDir: "foo",
				Files: map[string]fstest.File{
					"foo": fstest.File{
						Linkname: "",
						Perm:     os.ModePerm,
						Children: map[string]fstest.File{
							"bar": {
								Linkname: "",
								Perm:     os.ModePerm,
								Children: make(map[string]fstest.File, 0),
							},
						},
					},
				},
			},
			err: nil,
		},
		{
			drv: fstest.InMemoryDriver{
				CurrentDir: "foo",
				Files: map[string]fstest.File{
					"foo": fstest.File{
						Linkname: "",
						Perm:     os.ModePerm,
						Children: make(map[string]fstest.File, 0),
					},
				},
			},
			dirname: fstest.AbsPath("foo", "bar"),
			want: fstest.InMemoryDriver{
				CurrentDir: "foo",
				Files: map[string]fstest.File{
					"foo": fstest.File{
						Linkname: "",
						Perm:     os.ModePerm,
						Children: map[string]fstest.File{
							"bar": {
								Linkname: "",
								Perm:     os.ModePerm,
								Children: make(map[string]fstest.File, 0),
							},
						},
					},
				},
			},
			err: nil,
		},
		{
			drv:     fstest.InMemoryDriver{},
			dirname: filepath.Join("foo", "bar"),
			want: fstest.InMemoryDriver{
				Files: map[string]fstest.File{
					"foo": fstest.File{
						Linkname: "",
						Perm:     os.ModePerm,
						Children: map[string]fstest.File{
							"bar": {
								Linkname: "",
								Perm:     os.ModePerm,
								Children: make(map[string]fstest.File, 0),
							},
						},
					},
				},
			},
			err: nil,
		},
		{
			drv: fstest.InMemoryDriver{
				Files: map[string]fstest.File{
					"foo": fstest.File{
						Linkname: "",
						Perm:     os.ModePerm,
						Children: nil,
					},
				},
			},
			dirname: filepath.Join("foo", "bar"),
			want: fstest.InMemoryDriver{
				Files: map[string]fstest.File{
					"foo": fstest.File{
						Linkname: "",
						Perm:     os.ModePerm,
						Children: nil,
					},
				},
			},
			err: fstest.ErrExist,
		},
		{
			drv: fstest.InMemoryDriver{
				Files: map[string]fstest.File{
					"foo": fstest.File{
						Linkname: "",
						Perm:     os.ModePerm,
						Children: make(map[string]fstest.File, 0),
					},
				},
			},
			dirname: filepath.Join("foo", "bar") + string(filepath.Separator),
			want: fstest.InMemoryDriver{
				Files: map[string]fstest.File{
					"foo": fstest.File{
						Linkname: "",
						Perm:     os.ModePerm,
						Children: map[string]fstest.File{
							"bar": {
								Linkname: "",
								Perm:     os.ModePerm,
								Children: make(map[string]fstest.File, 0),
							},
						},
					},
				},
			},
			err: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.dirname, func(t *testing.T) {
			err := tc.drv.MkdirAll(tc.dirname)
			if want, got := tc.err, err; !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			if want, got := tc.want, tc.drv; !cmp.Equal(got, want) {
				t.Fatalf("(-want +got):\n%s", cmp.Diff(want, got))
			}
		})
	}
}

func testInMemoryDriverReadDir(t *testing.T) {
	testCases := []struct {
		drv     fstest.InMemoryDriver
		dirname string
		want    []fs.FileInfo
		err     error
	}{
		{
			drv: fstest.InMemoryDriver{
				Files: map[string]fstest.File{
					"test": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     nil,
						Children: map[string]fstest.File{
							"bar": {
								Linkname: "foo",
								Perm:     os.ModePerm,
								Data:     nil,
								Children: nil,
							},
							"foo": {
								Linkname: "",
								Perm:     os.ModePerm,
								Data:     []byte("test_foo"),
								Children: nil,
							},
						},
					},
				},
			},
			dirname: "test",
			want: []fs.FileInfo{
				fstest.FileStat{
					Label: "bar",
					File: fstest.File{
						Linkname: "foo",
						Perm:     os.ModePerm,
						Data:     nil,
						Children: nil,
					},
				},
				fstest.FileStat{
					Label: "foo",
					File: fstest.File{
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     []byte("test_foo"),
						Children: nil,
					},
				},
			},
			err: nil,
		},
		{
			drv: fstest.InMemoryDriver{
				CurrentDir: "test",
				Files: map[string]fstest.File{
					"test": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     nil,
						Children: map[string]fstest.File{
							"bar": {
								Linkname: "foo",
								Perm:     os.ModePerm,
								Data:     nil,
								Children: nil,
							},
							"foo": {
								Linkname: "",
								Perm:     os.ModePerm,
								Data:     []byte("test_foo"),
								Children: nil,
							},
						},
					},
				},
			},
			dirname: "",
			want: []fs.FileInfo{
				fstest.FileStat{
					Label: "bar",
					File: fstest.File{
						Linkname: "foo",
						Perm:     os.ModePerm,
						Data:     nil,
						Children: nil,
					},
				},
				fstest.FileStat{
					Label: "foo",
					File: fstest.File{
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     []byte("test_foo"),
						Children: nil,
					},
				},
			},
			err: nil,
		},
		{
			drv: fstest.InMemoryDriver{
				CurrentDir: "test",
				Files: map[string]fstest.File{
					"test": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     nil,
						Children: map[string]fstest.File{
							"bar": {
								Linkname: "foo",
								Perm:     os.ModePerm,
								Data:     nil,
								Children: nil,
							},
							"foo": {
								Linkname: "",
								Perm:     os.ModePerm,
								Data:     []byte("test_foo"),
								Children: nil,
							},
						},
					},
				},
			},
			dirname: fstest.AbsPath("test"),
			want: []fs.FileInfo{
				fstest.FileStat{
					Label: "bar",
					File: fstest.File{
						Linkname: "foo",
						Perm:     os.ModePerm,
						Data:     nil,
						Children: nil,
					},
				},
				fstest.FileStat{
					Label: "foo",
					File: fstest.File{
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     []byte("test_foo"),
						Children: nil,
					},
				},
			},
			err: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.dirname, func(t *testing.T) {
			orig := tc.drv // copy for further comparison
			files, err := tc.drv.ReadDir(tc.dirname)
			if want, got := tc.err, err; !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			if want, got := tc.want, files; !cmp.Equal(got, want) {
				t.Fatalf("(-want +got):\n%s", cmp.Diff(want, got))
			}
			if want, got := orig, tc.drv; !cmp.Equal(got, want) {
				t.Fatalf("(-want +got):\n%s", cmp.Diff(want, got))
			}
		})
	}
}

func testInMemoryDriverReadFile(t *testing.T) {
	testCases := []struct {
		drv      fstest.InMemoryDriver
		filename string
		want     []byte
		err      error
	}{
		{
			drv: fstest.InMemoryDriver{
				Files: map[string]fstest.File{
					"foo": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     []byte("test"),
						Children: nil,
					},
				},
			},
			filename: "foo",
			want:     []byte("test"),
			err:      nil,
		},
		{
			drv: fstest.InMemoryDriver{
				Files: map[string]fstest.File{
					"foo": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     []byte("test_foo"),
						Children: map[string]fstest.File{
							"bar": {
								Linkname: "",
								Perm:     os.ModePerm,
								Data:     []byte("test_bar"),
							},
						},
					},
				},
			},
			filename: filepath.Join("foo", "bar"),
			want:     []byte("test_bar"),
			err:      nil,
		},
		{
			drv: fstest.InMemoryDriver{
				CurrentDir: "foo",
				Files: map[string]fstest.File{
					"foo": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     []byte("test_foo"),
						Children: map[string]fstest.File{
							"bar": {
								Linkname: "",
								Perm:     os.ModePerm,
								Data:     []byte("test_bar"),
							},
						},
					},
				},
			},
			filename: "bar",
			want:     []byte("test_bar"),
			err:      nil,
		},
		{
			drv: fstest.InMemoryDriver{
				CurrentDir: "foo",
				Files: map[string]fstest.File{
					"foo": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     []byte("test_foo"),
						Children: map[string]fstest.File{
							"bar": {
								Linkname: "",
								Perm:     os.ModePerm,
								Data:     []byte("test_bar"),
							},
						},
					},
				},
			},
			filename: fstest.AbsPath("foo", "bar"),
			want:     []byte("test_bar"),
			err:      nil,
		},
		{
			drv: fstest.InMemoryDriver{
				Files: map[string]fstest.File{
					"foo": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     []byte("test"),
						Children: nil,
					},
				},
			},
			filename: "bar",
			want:     nil,
			err:      fstest.ErrNotExist,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			orig := tc.drv // copy for further comparison
			b, err := tc.drv.ReadFile(tc.filename)
			if want, got := tc.err, err; !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			if want, got := tc.want, b; string(got) != string(want) {
				t.Fatalf("want %q, got %q", want, got)
			}
			if want, got := orig, tc.drv; !cmp.Equal(got, want) {
				t.Fatalf("(-want +got):\n%s", cmp.Diff(want, got))
			}
		})
	}
}

func testInMemoryDriverStat(t *testing.T) {
	testCases := []struct {
		drv      fstest.InMemoryDriver
		filename string
		want     fstest.FileStat
		err      error
	}{
		{
			drv: fstest.InMemoryDriver{
				Files: map[string]fstest.File{
					"foo": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     []byte("test"),
						Children: nil,
					},
				},
			},
			filename: "foo",
			want: fstest.FileStat{
				Label: "foo",
				File: fstest.File{
					Linkname: "",
					Perm:     os.ModePerm,
					Data:     []byte("test"),
					Children: nil,
				},
			},
			err: nil,
		},
		{
			drv: fstest.InMemoryDriver{
				Files: map[string]fstest.File{
					"foo": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     []byte("test_foo"),
						Children: map[string]fstest.File{
							"bar": {
								Linkname: "",
								Perm:     os.ModePerm,
								Data:     []byte("test_bar"),
							},
						},
					},
				},
			},
			filename: filepath.Join("foo", "bar"),
			want: fstest.FileStat{
				Label: "bar",
				File: fstest.File{
					Linkname: "",
					Perm:     os.ModePerm,
					Data:     []byte("test_bar"),
				},
			},
			err: nil,
		},
		{
			drv: fstest.InMemoryDriver{
				CurrentDir: "foo",
				Files: map[string]fstest.File{
					"foo": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     []byte("test_foo"),
						Children: map[string]fstest.File{
							"bar": {
								Linkname: "",
								Perm:     os.ModePerm,
								Data:     []byte("test_bar"),
							},
						},
					},
				},
			},
			filename: "bar",
			want: fstest.FileStat{
				Label: "bar",
				File: fstest.File{
					Linkname: "",
					Perm:     os.ModePerm,
					Data:     []byte("test_bar"),
				},
			},
			err: nil,
		},
		{
			drv: fstest.InMemoryDriver{
				CurrentDir: "foo",
				Files: map[string]fstest.File{
					"foo": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     []byte("test_foo"),
						Children: map[string]fstest.File{
							"bar": {
								Linkname: "",
								Perm:     os.ModePerm,
								Data:     []byte("test_bar"),
							},
						},
					},
				},
			},
			filename: fstest.AbsPath("foo", "bar"),
			want: fstest.FileStat{
				Label: "bar",
				File: fstest.File{
					Linkname: "",
					Perm:     os.ModePerm,
					Data:     []byte("test_bar"),
				},
			},
			err: nil,
		},
		{
			drv: fstest.InMemoryDriver{
				Files: map[string]fstest.File{
					"foo": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     []byte("test"),
						Children: nil,
					},
				},
			},
			filename: "bar",
			want:     fstest.FileStat{},
			err:      nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			t.Run(tc.filename, func(t *testing.T) {
				orig := tc.drv // copy for further comparison
				b, err := tc.drv.Stat(tc.filename)
				if want, got := tc.err, err; !errors.Is(got, want) {
					t.Fatalf("want %v, got %v", want, got)
				}
				if want, got := tc.want, b; !cmp.Equal(got, want) {
					t.Fatalf("(-want +got):\n%s", cmp.Diff(want, got))
				}
				if want, got := orig, tc.drv; !cmp.Equal(got, want) {
					t.Fatalf("(-want +got):\n%s", cmp.Diff(want, got))
				}
			})
		})
	}
}

func testInMemoryDriverSymlink(t *testing.T) {
	testCases := []struct {
		drv     fstest.InMemoryDriver
		oldname string
		newname string
		want    fstest.InMemoryDriver
		err     error
	}{
		{
			drv: fstest.InMemoryDriver{
				Files: map[string]fstest.File{
					"foo": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     []byte("test"),
						Children: nil,
					},
				},
			},
			oldname: "foo",
			newname: "bar",
			want: fstest.InMemoryDriver{
				Files: map[string]fstest.File{
					"bar": {
						Linkname: "foo",
						Perm:     os.ModePerm,
						Data:     nil,
						Children: nil,
					},
					"foo": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     []byte("test"),
						Children: nil,
					},
				},
			},
			err: nil,
		},
		{
			drv: fstest.InMemoryDriver{
				Files: map[string]fstest.File{
					"bar": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     []byte("test_bar"),
						Children: nil,
					},
					"foo": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     []byte("test_foo"),
						Children: nil,
					},
				},
			},
			oldname: "foo",
			newname: "bar",
			want: fstest.InMemoryDriver{
				Files: map[string]fstest.File{
					"bar": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     []byte("test_bar"),
						Children: nil,
					},
					"foo": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     []byte("test_foo"),
						Children: nil,
					},
				},
			},
			err: fstest.ErrExist,
		},
		{
			drv: fstest.InMemoryDriver{
				Files: map[string]fstest.File{
					"foo": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     []byte("test"),
						Children: map[string]fstest.File{
							"bar": {
								Linkname: "foo",
								Perm:     os.ModePerm,
								Data:     nil,
								Children: nil,
							},
						},
					},
				},
			},
			oldname: filepath.Join("foo", "bar"),
			newname: "bar",
			want: fstest.InMemoryDriver{
				Files: map[string]fstest.File{
					"foo": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     []byte("test"),
						Children: map[string]fstest.File{
							"bar": {
								Linkname: "foo",
								Perm:     os.ModePerm,
								Data:     nil,
								Children: nil,
							},
						},
					},
					"bar": {
						Linkname: filepath.Join("foo", "bar"),
						Perm:     os.ModePerm,
						Data:     nil,
						Children: nil,
					},
				},
			},
			err: nil,
		},
		{
			drv: fstest.InMemoryDriver{
				CurrentDir: "foo",
				Files: map[string]fstest.File{
					"foo": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     []byte("test"),
						Children: map[string]fstest.File{
							"bar": {
								Linkname: "foo",
								Perm:     os.ModePerm,
								Data:     nil,
								Children: nil,
							},
						},
					},
				},
			},
			oldname: "bar",
			newname: fstest.AbsPath("bar"), // XXX: empty, non-absolute path also works
			want: fstest.InMemoryDriver{
				CurrentDir: "foo",
				Files: map[string]fstest.File{
					"foo": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     []byte("test"),
						Children: map[string]fstest.File{
							"bar": {
								Linkname: "foo",
								Perm:     os.ModePerm,
								Data:     nil,
								Children: nil,
							},
						},
					},
					"bar": {
						Linkname: filepath.Join("foo", "bar"),
						Perm:     os.ModePerm,
						Data:     nil,
						Children: nil,
					},
				},
			},
			err: nil,
		},
		{
			drv: fstest.InMemoryDriver{
				CurrentDir: "foo",
				Files: map[string]fstest.File{
					"foo": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     []byte("test"),
						Children: map[string]fstest.File{
							"bar": {
								Linkname: fstest.AbsPath("foo", "bar"),
								Perm:     os.ModePerm,
								Data:     nil,
								Children: nil,
							},
						},
					},
				},
			},
			oldname: fstest.AbsPath("foo", "bar"),
			newname: fstest.AbsPath("bar"),
			want: fstest.InMemoryDriver{
				CurrentDir: "foo",
				Files: map[string]fstest.File{
					"foo": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     []byte("test"),
						Children: map[string]fstest.File{
							"bar": {
								Linkname: fstest.AbsPath("foo", "bar"),
								Perm:     os.ModePerm,
								Data:     nil,
								Children: nil,
							},
						},
					},
					"bar": {
						Linkname: filepath.Join("foo", "bar"),
						Perm:     os.ModePerm,
						Data:     nil,
						Children: nil,
					},
				},
			},
			err: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.oldname+" "+tc.newname, func(t *testing.T) {
			err := tc.drv.Symlink(tc.oldname, tc.newname)
			if want, got := tc.err, err; !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			if want, got := tc.want, tc.drv; !cmp.Equal(got, want) {
				t.Fatalf("(-want +got):\n%s", cmp.Diff(want, got))
			}
		})
	}
}

func testInMemoryDriverWriteFile(t *testing.T) {
	testCases := []struct {
		drv      fstest.InMemoryDriver
		filename string
		data     []byte
		perm     os.FileMode
		want     fstest.InMemoryDriver
		err      error
	}{
		{
			drv:      fstest.InMemoryDriver{},
			filename: "foo",
			data:     []byte("bar"),
			perm:     0o644,
			want: fstest.InMemoryDriver{
				Files: map[string]fstest.File{
					"foo": {
						Linkname: "",
						Perm:     0o644,
						Data:     []byte("bar"),
						Children: nil,
					},
				},
			},
			err: nil,
		},
		{
			drv: fstest.InMemoryDriver{
				Files: map[string]fstest.File{
					"foo": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     nil,
						Children: make(map[string]fstest.File),
					},
				},
			},
			filename: "foo",
			data:     []byte("bar"),
			perm:     0o644,
			want: fstest.InMemoryDriver{
				Files: map[string]fstest.File{
					"foo": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     nil,
						Children: make(map[string]fstest.File),
					},
				},
			},
			err: fstest.ErrExist,
		},
		{
			drv: fstest.InMemoryDriver{
				Files: map[string]fstest.File{
					"foo": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     nil,
						Children: nil,
					},
				},
			},
			filename: "foo",
			data:     []byte("bar"),
			perm:     0o644,
			want: fstest.InMemoryDriver{
				Files: map[string]fstest.File{
					"foo": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     []byte("bar"),
						Children: nil,
					},
				},
			},
			err: nil,
		},
		{
			drv: fstest.InMemoryDriver{
				Files: map[string]fstest.File{
					"foo": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     nil,
						Children: make(map[string]fstest.File),
					},
				},
			},
			filename: filepath.Join("foo", "bar"),
			data:     []byte("bar"),
			perm:     0o644,
			want: fstest.InMemoryDriver{
				Files: map[string]fstest.File{
					"foo": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     nil,
						Children: map[string]fstest.File{
							"bar": {
								Linkname: "",
								Perm:     0o644,
								Data:     []byte("bar"),
								Children: nil,
							},
						},
					},
				},
			},
			err: nil,
		},
		{
			drv: fstest.InMemoryDriver{
				CurrentDir: "foo",
				Files: map[string]fstest.File{
					"foo": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     nil,
						Children: make(map[string]fstest.File),
					},
				},
			},
			filename: "bar",
			data:     []byte("bar"),
			perm:     0o644,
			want: fstest.InMemoryDriver{
				CurrentDir: "foo",
				Files: map[string]fstest.File{
					"foo": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     nil,
						Children: map[string]fstest.File{
							"bar": {
								Linkname: "",
								Perm:     0o644,
								Data:     []byte("bar"),
								Children: nil,
							},
						},
					},
				},
			},
			err: nil,
		},
		{
			drv: fstest.InMemoryDriver{
				CurrentDir: "foo",
				Files: map[string]fstest.File{
					"foo": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     nil,
						Children: make(map[string]fstest.File),
					},
				},
			},
			filename: fstest.AbsPath("foo", "bar"),
			data:     []byte("bar"),
			perm:     0o644,
			want: fstest.InMemoryDriver{
				CurrentDir: "foo",
				Files: map[string]fstest.File{
					"foo": {
						Linkname: "",
						Perm:     os.ModePerm,
						Data:     nil,
						Children: map[string]fstest.File{
							"bar": {
								Linkname: "",
								Perm:     0o644,
								Data:     []byte("bar"),
								Children: nil,
							},
						},
					},
				},
			},
			err: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			err := tc.drv.WriteFile(tc.filename, tc.data, tc.perm)
			if want, got := tc.err, err; !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			if want, got := tc.want, tc.drv; !cmp.Equal(got, want) {
				t.Fatalf("(-want +got):\n%s", cmp.Diff(want, got))
			}
		})
	}
}

func TestFileStat(t *testing.T) {
	t.Run("Name", func(t *testing.T) {
		testCases := []struct {
			f    fstest.FileStat
			want string
		}{
			{fstest.FileStat{Label: "foo", File: fstest.File{}}, "foo"},
		}
		for _, tc := range testCases {
			t.Run("", func(t *testing.T) {
				if want, got := tc.want, tc.f.Name(); got != want {
					t.Fatalf("want %q, got %q", want, got)
				}
			})
		}
	})
	t.Run("Exists", func(t *testing.T) {
		testCases := []struct {
			f    fstest.FileStat
			want bool
		}{
			{fstest.FileStat{Label: "foo", File: fstest.File{}}, true},
			{fstest.FileStat{}, false},
		}
		for _, tc := range testCases {
			t.Run("", func(t *testing.T) {
				if want, got := tc.want, tc.f.Exists(); got != want {
					t.Fatalf("want %t, got %t", want, got)
				}
			})
		}
	})
	t.Run("IsDir", func(t *testing.T) {
		testCases := []struct {
			f    fstest.FileStat
			want bool
		}{
			{fstest.FileStat{File: fstest.File{Children: nil}}, false},
			{fstest.FileStat{File: fstest.File{Children: make(map[string]fstest.File, 0)}}, true},
			{fstest.FileStat{File: fstest.File{Children: make(map[string]fstest.File, 5)}}, true},
		}
		for _, tc := range testCases {
			t.Run("", func(t *testing.T) {
				if want, got := tc.want, tc.f.IsDir(); got != want {
					t.Fatalf("want %t, got %t", want, got)
				}
			})
		}
	})
	t.Run("Linkname", func(t *testing.T) {
		testCases := []struct {
			f    fstest.FileStat
			want string
		}{
			{
				fstest.FileStat{File: fstest.File{Linkname: "foo"}}, "foo",
			},
		}
		for _, tc := range testCases {
			t.Run("", func(t *testing.T) {
				if want, got := tc.want, tc.f.Linkname(); got != want {
					t.Fatalf("want %q, got %q", want, got)
				}
			})
		}
	})
	t.Run("Perm", func(t *testing.T) {
		testCases := []struct {
			f    fstest.FileStat
			want os.FileMode
		}{
			{fstest.FileStat{File: fstest.File{Perm: os.ModePerm}}, os.ModePerm},
		}
		for _, tc := range testCases {
			t.Run("", func(t *testing.T) {
				if want, got := tc.want, tc.f.Perm(); got != want {
					t.Fatalf("want %#o, got %#o", got, want)
				}
			})
		}
	})
}
