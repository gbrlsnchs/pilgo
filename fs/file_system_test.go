package fs_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/andybalholm/crlf"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/text/transform"
	"gsr.dev/pilgrim/fs"
)

var (
	filePerms = map[string]os.FileMode{
		"darwin":  0o644,
		"linux":   0o644,
		"windows": 0o666,
	}
	directoryPerms = map[string]os.FileMode{
		"darwin":  0o755,
		"linux":   0o755,
		"windows": 0o777,
	}
	symlinkPerms = map[string]os.FileMode{
		"darwin":  0o755,
		"linux":   0o777,
		"windows": 0o666,
	}
)

func TestFileSystem(t *testing.T) {
	t.Run("Info", testFileSystemInfo)
	t.Run("ReadDir", testFileSystemReadDir)
	t.Run("ReadFile", testFileSystemReadFile)
	t.Run("WriteFile", testFileSystemWriteFile)
}

func testFileSystemInfo(t *testing.T) {
	testCases := []struct {
		filename string
		exists   bool
		isDir    bool
		linkname string
		perm     os.FileMode
		err      error
	}{
		{
			filename: "exists",
			exists:   true,
			isDir:    false,
			linkname: "",
			perm:     filePerms[runtime.GOOS],
			err:      nil,
		},
		{
			filename: "not_exists",
			exists:   false,
			isDir:    false,
			linkname: "",
			perm:     0,
			err:      nil,
		},
		{
			filename: "directory",
			exists:   true,
			isDir:    true,
			linkname: "",
			perm:     directoryPerms[runtime.GOOS],
			err:      nil,
		},
		{
			filename: "symlink",
			exists:   true,
			isDir:    false,
			linkname: "directory",
			perm:     symlinkPerms[runtime.GOOS],
			err:      nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			var (
				fs       fs.FileSystem
				filename = filepath.Join("testdata", t.Name())
			)
			fi, err := fs.Info(filename)
			if want, got := tc.err, err; !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			t.Run("Exists", func(t *testing.T) {
				if want, got := tc.exists, fi.Exists(); want != got {
					t.Fatalf("want %t, got %t", want, got)
				}
			})
			t.Run("IsDir", func(t *testing.T) {
				if want, got := tc.isDir, fi.IsDir(); want != got {
					t.Fatalf("want %t, got %t", want, got)
				}
			})
			t.Run("Linkname", func(t *testing.T) {
				if want, got := tc.linkname, fi.Linkname(); want != got {
					t.Fatalf("want %q, got %q", want, got)
				}
			})
			t.Run("Perm", func(t *testing.T) {
				if want, got := tc.perm, fi.Perm(); got != want {
					t.Fatalf("want %#o, got %#o", want, got)
				}
			})
		})
	}
}

func testFileSystemReadDir(t *testing.T) {
	testCases := []struct {
		filename string
		want     []string
		err      error
	}{
		{
			filename: "directory",
			want:     []string{"bar", "foo"},
			err:      nil,
		},
		{
			filename: "symlink",
			want:     []string{"bar", "foo"},
			err:      nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			var (
				fs       fs.FileSystem
				filename = filepath.Join("testdata", t.Name())
			)
			files, err := fs.ReadDir(filename)
			if want, got := tc.err, err; !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			if want, got := tc.want, files; !cmp.Equal(got, want) {
				t.Errorf(
					"(*FileSystem).ReadDir mismatch (-want +got):\n%s",
					cmp.Diff(want, got),
				)
			}
		})
	}
}

func testFileSystemReadFile(t *testing.T) {
	testCases := []struct {
		filename string
		want     string
		err      error
	}{
		{
			filename: "test.txt",
			want:     "read file test\n",
			err:      nil,
		},
		{
			filename: "nonexistent.txt",
			want:     "",
			err:      os.ErrNotExist,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			var (
				fs       fs.FileSystem
				filename = filepath.Join("testdata", t.Name())
			)
			b, err := fs.ReadFile(filename)
			if want, got := tc.err, err; !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			if want, got := tc.want, b; string(got) != want {
				t.Errorf("want %q, got %q", want, got)
			}
		})
	}
}

func testFileSystemWriteFile(t *testing.T) {
	testCases := []struct {
		filename string
		wantData string
		wantPerm os.FileMode
		err      error
	}{
		{
			filename: "test.txt",
			wantData: "write file test\n",
			wantPerm: filePerms[runtime.GOOS],
			err:      nil,
		},
	}
	normalize := new(crlf.Normalize)
	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			var (
				fs       fs.FileSystem
				filename = filepath.Join("testdata", t.Name())
			)
			err := fs.WriteFile(filename, []byte(tc.wantData), tc.wantPerm)
			if want, got := tc.err, err; !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			f, err := os.Open(filename)
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(filename)
			defer f.Close()
			r := transform.NewReader(f, normalize)
			b, err := ioutil.ReadAll(r)
			if err != nil {
				t.Fatal(err)
			}
			if want, got := tc.wantData, b; string(got) != want {
				t.Errorf("want %q, got %q", want, got)
			}
			fi, err := os.Stat(filename)
			if err != nil {
				t.Fatal(err)
			}
			if want, got := tc.wantPerm, fi.Mode().Perm(); got != want {
				t.Errorf("want %#o, got %#o", want, got)
			}
		})
	}

}
