package fs_test

import (
	"errors"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gsr.dev/pilgrim/fs"
	"gsr.dev/pilgrim/fs/fstest"
)

func TestFileSystem(t *testing.T) {
	t.Run("MkdirAll", testFileSystemMkdirAll)
	t.Run("ReadDir", testFileSystemReadDir)
	t.Run("ReadFile", testFileSystemReadFile)
	t.Run("Stat", testFileSystemStat)
	t.Run("WriteFile", testFileSystemWriteFile)
}

func testFileSystemMkdirAll(t *testing.T) {
	testCases := []struct {
		drv fs.Driver
		err error
	}{
		{nil, fs.ErrNoDriver},
		{new(fstest.Driver), nil},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			defer checkPanic(t, tc.err)
			fs := fs.New(tc.drv)
			_ = fs.MkdirAll("test")
			drv := tc.drv.(*fstest.Driver)
			hasBeenCalled, args := drv.HasBeenCalled(drv.MkdirAll)
			if want, got := true, hasBeenCalled; got != want {
				t.Fatalf("want %t, got %t", want, got)
			}
			if want, got := (fstest.Args{"test"}), args; !cmp.Equal(got, want) {
				t.Fatalf("FileSystem.MkdirAll mismatch (-want +got):\n%s", cmp.Diff(want, got))
			}
		})
	}
}

func testFileSystemReadDir(t *testing.T) {
	testCases := []struct {
		drv fs.Driver
		err error
	}{
		{nil, fs.ErrNoDriver},
		{new(fstest.Driver), nil},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			defer checkPanic(t, tc.err)
			fs := fs.New(tc.drv)
			_, _ = fs.ReadDir("test")
			drv := tc.drv.(*fstest.Driver)
			hasBeenCalled, args := drv.HasBeenCalled(drv.ReadDir)
			if want, got := true, hasBeenCalled; got != want {
				t.Fatalf("want %t, got %t", want, got)
			}
			if want, got := (fstest.Args{"test"}), args; !cmp.Equal(got, want) {
				t.Fatalf("FileSystem.ReadDir mismatch (-want +got):\n%s", cmp.Diff(want, got))
			}
		})
	}
}

func testFileSystemReadFile(t *testing.T) {
	testCases := []struct {
		drv fs.Driver
		err error
	}{
		{nil, fs.ErrNoDriver},
		{new(fstest.Driver), nil},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			defer checkPanic(t, tc.err)
			fs := fs.New(tc.drv)
			_, _ = fs.ReadFile("test")
			drv := tc.drv.(*fstest.Driver)
			hasBeenCalled, args := drv.HasBeenCalled(drv.ReadFile)
			if want, got := true, hasBeenCalled; got != want {
				t.Fatalf("want %t, got %t", want, got)
			}
			if want, got := (fstest.Args{"test"}), args; !cmp.Equal(got, want) {
				t.Fatalf("FileSystem.ReadFile mismatch (-want +got):\n%s", cmp.Diff(want, got))
			}
		})
	}
}

func testFileSystemStat(t *testing.T) {
	testCases := []struct {
		drv fs.Driver
		err error
	}{
		{nil, fs.ErrNoDriver},
		{new(fstest.Driver), nil},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			defer checkPanic(t, tc.err)
			fs := fs.New(tc.drv)
			_, _ = fs.Stat("test")
			drv := tc.drv.(*fstest.Driver)
			hasBeenCalled, args := drv.HasBeenCalled(drv.Stat)
			if want, got := true, hasBeenCalled; got != want {
				t.Fatalf("want %t, got %t", want, got)
			}
			if want, got := (fstest.Args{"test"}), args; !cmp.Equal(got, want) {
				t.Fatalf("FileSystem.Stat mismatch (-want +got):\n%s", cmp.Diff(want, got))
			}
		})
	}
}

func testFileSystemWriteFile(t *testing.T) {
	testCases := []struct {
		drv fs.Driver
		err error
	}{
		{nil, fs.ErrNoDriver},
		{new(fstest.Driver), nil},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			defer checkPanic(t, tc.err)
			fs := fs.New(tc.drv)
			_ = fs.WriteFile("test", []byte("testing"), 0o777)
			drv := tc.drv.(*fstest.Driver)
			hasBeenCalled, args := drv.HasBeenCalled(drv.WriteFile)
			if want, got := true, hasBeenCalled; got != want {
				t.Fatalf("want %t, got %t", want, got)
			}
			if want, got := (fstest.Args{"test", []byte("testing"), os.FileMode(0o777)}), args; !cmp.Equal(got, want) {
				t.Fatalf("FileSystem.WriteFile mismatch (-want +got):\n%s", cmp.Diff(want, got))
			}
		})
	}
}

func checkPanic(t *testing.T, want error) {
	if r := recover(); r != nil {
		err := r.(error)
		if got := err; !errors.Is(got, want) {
			t.Fatalf("want %v, got %v", want, got)
		}
	}
}
