// +build mage

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/renameio"
	"github.com/magefile/mage/sh"
)

const (
	module         = "github.com/gbrlsnchs/pilgo"
	versionEnvName = "PILGO_VERSION"
	dstdir         = ".bin"
	sep            = string(filepath.Separator)
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

type linuxTxt string

const (
	slash        linuxTxt = "/"
	fileNotFound linuxTxt = "no such file or directory"
	readDirErr   linuxTxt = "readdirent: not a directory"
)

type replacement struct {
	old linuxTxt
	new string
}

// GenCLITests copies Linux E2E CLI tests and properly adapts it to other platforms.
func GenCLITests() error {
	const linuxdir = "linux"
	oses := map[string][]replacement{
		"darwin": {
			{old: readDirErr, new: "fdopendir: not a directory"},
		},
		"windows": {
			{old: slash, new: "\\"},
			{old: fileNotFound, new: "The system cannot find the file specified."},
		},
	}
	for sys, rep := range oses {
		testdir := filepath.Join("cmd", "plg", "testdata", "TestCLI")
		windir := filepath.Join(testdir, sys)
		if err := os.RemoveAll(windir); err != nil {
			return err
		}
		if err := os.Mkdir(windir, 0o755); err != nil {
			return err
		}
		srcdir := filepath.Join(testdir, linuxdir)
		files, err := ioutil.ReadDir(filepath.Join(testdir, linuxdir))
		if err != nil {
			return err
		}
		for _, fi := range files {
			fname := fi.Name()
			b, err := ioutil.ReadFile(filepath.Join(srcdir, fname))
			if err != nil {
				return err
			}
			for _, r := range rep {
				b = bytes.ReplaceAll(b, []byte(r.old), []byte(r.new))
			}
			if err := renameio.WriteFile(filepath.Join(windir, fname),
				b, fi.Mode().Perm()); err != nil {
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
