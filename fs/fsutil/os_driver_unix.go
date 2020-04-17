// +build !windows

package fsutil

import (
	"os"

	"github.com/google/renameio"
)

// WriteFile writes data to filename atomically with permission perm.
func (OSDriver) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return renameio.WriteFile(filename, data, perm)
}
