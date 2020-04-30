// +build mage

package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"text/template"

	"github.com/google/renameio"
	"github.com/magefile/mage/sh"
)

const (
	module = "github.com/gbrlsnchs/pilgo"
	dstdir = ".bin"
	sep    = string(filepath.Separator)
)

var (
	env = map[string]string{
		"CGO":    "0",
		"GOARCH": "amd64",
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
	tag, err := exec.Command("git", "tag", "--points-at=HEAD").Output() // git tag --points-at=HEAD
	if err != nil {
		return err
	}
	if string(tag) == "" {
		return errors.New("aborting build: no tag found")
	}
	for _, ptf := range platforms {
		env["GOOS"] = ptf
		for _, x := range exes {
			fmt.Printf("Building %q for %q...\n", x, ptf)
			src := filepath.Join(module, "cmd", x)
			dst := filepath.Join(dstdir, ptf, executable(x, ptf))
			ldflags := ldflagsVar("internal.version", string(tag))
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

type testPatch struct {
	reps      replacement
	buildData func([]byte) map[string]interface{}
}

// GenCLITests copies Linux E2E CLI tests and properly adapts it to other platforms.
func GenCLITests() error {
	const linuxdir = "linux"
	oses := map[string]map[string][]testPatch{
		"darwin": {
			".*": {
				{
					reps: replacement{
						old: readDirErr,
						new: "fdopendir: not a directory",
					},
				},
			},
		},
		"windows": {
			".*": {
				{
					reps: replacement{
						old: slash,
						new: "\\",
					},
				},
				{
					reps: replacement{
						old: fileNotFound,
						new: "The system cannot find the file specified.",
					},
				},
			},
			"config.ct": {
				{
					reps: replacement{
						old: readDirErr,
						new: "Readdir {{.scandir}}: The system cannot find the path specified.",
					},
					buildData: func(b []byte) map[string]interface{} {
						r := regexp.MustCompile(`plg config -scandir (\w+) --> FAIL`)
						subs := r.FindSubmatch(b)
						t := make(map[string]interface{})
						if len(subs) > 1 {
							t["scandir"] = string(subs[1])
						}
						return t
					},
				},
			},
		},
	}
	for sys, regexes := range oses {
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
			for rgx, repls := range regexes {
				match, err := regexp.MatchString(rgx, fname)
				if err != nil {
					return err
				}
				if !match {
					continue
				}
				for _, p := range repls {
					r := p.reps
					b = bytes.ReplaceAll(b, []byte(r.old), []byte(r.new))
					if p.buildData != nil {
						data := p.buildData(b)
						tmpl := template.Must(
							template.New(fname).Parse(string(b)))
						buf := bytes.NewBuffer(make([]byte, 0, len(b)))
						if err = tmpl.Execute(buf, data); err != nil {
							return err
						}
						b = buf.Bytes()
					}
				}
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
	ldflags := ldflagsVar("internal.version", "test_version")
	return sh.RunWith(env, "go", "test", "-race", "-ldflags", ldflags, "./...")
}

func executable(x string, platform string) string {
	if platform == "windows" {
		return x + ".exe"
	}
	return x
}

func ldflagsVar(name, value string) string {
	return fmt.Sprintf("-X %s/cmd/%s=%s", module, name, value)
}
