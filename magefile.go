// +build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/sh"
)

const (
	module         = "github.com/gbrlsnchs/pilgo"
	versionEnvName = "PILGO_VERSION"
	dstdir         = ".bin"
)

var (
	ldflags = fmt.Sprintf("-X %s/cmd/internal.version=$%s", module, versionEnvName)
	env     = map[string]string{
		versionEnvName: pilgoVersion(),
		"CGO":          "0",
		"GOARCH":       "amd64",
	}
	exes      = []string{"plg"}
	platforms = []string{
		"darwin",
		"linux",
		"windows",
	}
)

var Default = Build

func Build() error {
	for _, ptf := range platforms {
		env["GOOS"] = ptf
		for _, x := range exes {
			fmt.Printf("Building %q for %q...\n", x, ptf)
			src := filepath.Join(module, "cmd", x)
			dst := filepath.Join(dstdir, ptf, executable(x, ptf))
			err := sh.RunWith(env, "go", "build", "-o", dst, "-ldflags", ldflags, src)
			if err != nil {
				return err
			}

		}
	}
	return nil
}

func Test() error {
	return sh.RunWith(env, "go", "test", "-race", "-ldflags", ldflags, "./...")
}

func pilgoVersion() string {
	v := os.Getenv(versionEnvName)
	if v != "" {
		return strings.TrimPrefix(v, "refs/tags/")
	}
	return "test_version"
}

func executable(x string, platform string) string {
	if platform == "windows" {
		return x + ".exe"
	}
	return x
}
