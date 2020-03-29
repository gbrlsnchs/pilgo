package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/andybalholm/crlf"
	cmdtest "github.com/google/go-cmdtest"
	"golang.org/x/text/transform"
)

const failureStatus = 0xDEADC0DE // 3735929054

var update = flag.Bool("update", false, "update test files with results")

func TestCLI(t *testing.T) {
	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	testdata := filepath.Join("testdata", t.Name())
	ts, err := cmdtest.Read(filepath.Join(testdata, runtime.GOOS))
	if err != nil {
		t.Fatal(err)
	}
	// Utility commands.
	ts.Commands["cp"] = cmdtest.InProcessProgram("cp", cpCmd(filepath.Join(pwd, testdata)))

	// Pilgo commands.
	ts.Commands["plg"] = cmdtest.InProcessProgram("plg", run)
	ts.Run(t, *update)
}

func cpCmd(pwd string) func() int {
	return func() int {
		var (
			argv            = os.Args[1:]
			argc            = len(argv)
			original, clone string
		)
		if argc > 0 {
			original = argv[0]
			clone = original
			if argc > 1 {
				clone = argv[1]
			}
		}
		if original == "" || clone == "" {
			return 0 // NOP
		}
		b, err := ioutil.ReadFile(filepath.Join(pwd, original))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return failureStatus
		}
		dir, _ := filepath.Split(clone)
		if dir != "" {
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return failureStatus
			}
		}
		if err := ioutil.WriteFile(clone, b, 0o644); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return failureStatus
		}
		return 0
	}

}

// misc
func readFile(t *testing.T, name string) string {
	golden, err := ioutil.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	goldenlf, _, err := transform.String(
		new(crlf.Normalize),
		filepath.FromSlash(string(golden)),
	)
	if err != nil {
		t.Fatal(err)
	}
	return goldenlf
}
