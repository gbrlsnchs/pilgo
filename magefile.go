// +build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/sh"
)

const (
	module         = "github.com/gbrlsnchs/pilgo"
	versionEnvName = "PILGO_VERSION"
)

var (
	ldflags = fmt.Sprintf("-X %s/cmd/internal.version=$%s", module, versionEnvName)
	env     = map[string]string{versionEnvName: pilgoVersion()}
)

func Test() (err error) {
	return sh.RunWith(env, "go", "test", "-race", "-ldflags", ldflags, "./...")
}

func pilgoVersion() string {
	v := os.Getenv(versionEnvName)
	if v != "" {
		return v
	}
	return "test_version"
}
