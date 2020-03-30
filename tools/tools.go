// +build tools

package tools

import (
	// The following imports are tools intended
	// to be used while developing or running CI.
	_ "github.com/magefile/mage"
	_ "golang.org/x/lint/golint"
	_ "golang.org/x/tools/gopls"
)
