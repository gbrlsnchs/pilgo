package parser_test

import (
	"path/filepath"
	"testing"

	"github.com/gbrlsnchs/pilgo/parser"
)

func TestFile(t *testing.T) {
	t.Run("Dir", testFileDir)
	t.Run("FullPath", testFileFullPath)
}

func testFileDir(t *testing.T) {
	testCases := []struct {
		f    parser.File
		want string
	}{
		{
			f:    parser.File{},
			want: "",
		},
		{
			f: parser.File{
				BaseDir: "test",
				Path:    nil,
			},
			want: "test",
		},
		{
			f: parser.File{
				BaseDir: "test",
				Path:    []string{"foo"},
			},
			want: "test",
		},
		{
			f: parser.File{
				BaseDir: "test",
				Path:    []string{"foo", "bar"},
			},
			want: filepath.Join("test", "foo"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.f.FullPath(), func(t *testing.T) {
			if want, got := tc.want, tc.f.Dir(); got != want {
				t.Errorf("want %q, got %q", want, got)
			}
		})
	}
}

func testFileFullPath(t *testing.T) {
	testCases := []struct {
		f    parser.File
		want string
	}{
		{
			f:    parser.File{},
			want: "",
		},
		{
			f: parser.File{
				BaseDir: "test",
				Path:    nil,
			},
			want: "test",
		},
		{
			f: parser.File{
				BaseDir: "test",
				Path:    []string{"foo"},
			},
			want: filepath.Join("test", "foo"),
		},
		{
			f: parser.File{
				BaseDir: "test",
				Path:    []string{"foo", "bar"},
			},
			want: filepath.Join("test", "foo", "bar"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.want, func(t *testing.T) {
			if want, got := tc.want, tc.f.FullPath(); got != want {
				t.Errorf("want %q, got %q", want, got)
			}
		})
	}
}
