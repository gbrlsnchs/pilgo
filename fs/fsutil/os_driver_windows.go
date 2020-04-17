// +build windows

package fsutil

import (
	"io/ioutil"
	"os"
)

// WriteFile writes data to filename with permission perm.
func (OSDriver) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return ioutil.WriteFile(filename, data, perm)
}
