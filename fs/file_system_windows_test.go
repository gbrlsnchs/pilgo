package fs_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	// XXX: Since symlinks in testdata folder aren't recognized
	// on Windows, this creates them only for that platform.
	testdirs := []string{
		"Info",
		"ReadDir",
	}
	for _, dir := range testdirs {
		symlink := filepath.Join("testdata", "TestFileSystem", dir, "symlink")
		if err := os.Remove(symlink); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err := os.Symlink("directory", symlink); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	os.Exit(m.Run())
}
