package fsutil_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/andybalholm/crlf"
	"golang.org/x/text/transform"
	"gsr.dev/pilgrim/fs/fsutil"
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

func TestOSDriver(t *testing.T) {
	t.Run("ReadDir", testOSDriverReadDir)
	t.Run("ReadFile", testOSDriverReadFile)
	t.Run("Stat", testOSDriverStat)
	t.Run("WriteFile", testOSDriverWriteFile)
}

func testOSDriverReadDir(t *testing.T) {
	type wantInfo struct {
		name     string
		exists   bool
		isDir    bool
		linkname string
		perm     os.FileMode
	}
	testCases := []struct {
		filename string
		want     []wantInfo
		err      error
	}{
		{
			filename: "directory",
			want: []wantInfo{
				{"bar", true, false, "", filePerms[runtime.GOOS]},
				{"foo", true, false, "", filePerms[runtime.GOOS]},
			},
			err: nil,
		},
		{
			filename: "symlink",
			want: []wantInfo{
				{"bar", true, false, "", filePerms[runtime.GOOS]},
				{"foo", true, false, "", filePerms[runtime.GOOS]},
			},
			err: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			var (
				drv      fsutil.OSDriver
				filename = filepath.Join("testdata", t.Name())
			)
			files, err := drv.ReadDir(filename)
			if want, got := tc.err, err; !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			if want, got := len(tc.want), len(files); got != want {
				t.Fatalf("want %d files, got %d", want, got)
			}
			for i, fi := range files {
				tc := tc.want[i]
				t.Run("Name", func(t *testing.T) {
					if want, got := tc.name, fi.Name(); want != got {
						t.Fatalf("want %q, got %q", want, got)
					}
				})
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
			}
		})
	}
}

func testOSDriverReadFile(t *testing.T) {
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
				drv      fsutil.OSDriver
				filename = filepath.Join("testdata", t.Name())
			)
			b, err := drv.ReadFile(filename)
			if want, got := tc.err, err; !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			if want, got := tc.want, b; string(got) != want {
				t.Errorf("want %q, got %q", want, got)
			}
		})
	}
}

func testOSDriverStat(t *testing.T) {
	testCases := []struct {
		filename string
		name     string
		exists   bool
		isDir    bool
		linkname string
		perm     os.FileMode
		err      error
	}{
		{
			filename: "exists",
			name:     "exists",
			exists:   true,
			isDir:    false,
			linkname: "",
			perm:     filePerms[runtime.GOOS],
			err:      nil,
		},
		{
			filename: "not_exists",
			name:     "",
			exists:   false,
			isDir:    false,
			linkname: "",
			perm:     0,
			err:      nil,
		},
		{
			filename: "directory",
			name:     "directory",
			exists:   true,
			isDir:    true,
			linkname: "",
			perm:     directoryPerms[runtime.GOOS],
			err:      nil,
		},
		{
			filename: "symlink",
			name:     "symlink",
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
				drv      fsutil.OSDriver
				filename = filepath.Join("testdata", t.Name())
			)
			fi, err := drv.Stat(filename)
			if want, got := tc.err, err; !errors.Is(got, want) {
				t.Fatalf("want %v, got %v", want, got)
			}
			t.Run("Name", func(t *testing.T) {
				if want, got := tc.name, fi.Name(); want != got {
					t.Fatalf("want %q, got %q", want, got)
				}
			})
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

func testOSDriverWriteFile(t *testing.T) {
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
				drv      fsutil.OSDriver
				filename = filepath.Join("testdata", t.Name())
			)
			err := drv.WriteFile(filename, []byte(tc.wantData), tc.wantPerm)
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
