package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	cmdtest "github.com/google/go-cmdtest"
)

const failureStatus = 0xDEADC0DE // 3735929054

var update = flag.Bool("update", false, "update test files with results")

func TestCLI(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	testdata := filepath.Join("testdata", t.Name())
	ts, err := cmdtest.Read(filepath.Join(testdata, runtime.GOOS))
	if err != nil {
		t.Fatal(err)
	}
	// Utility commands.
	ts.Commands["cp"] = cpCmd(filepath.Join(cwd, testdata))

	// Pilgo commands.
	ts.Commands["plg"] = cmdtest.InProcessProgram("plg", run)
	ts.Run(t, *update)
}

func cpCmd(pwd string) func([]string, string) ([]byte, error) {
	return func(args []string, inputFile string) ([]byte, error) {
		switch len(args) {
		case 0:
			return nil, errors.New("missing file operand")
		case 1:
			return nil, fmt.Errorf("missing file operand after %s", args[0])
		}
		orig, clone := args[0], args[1]
		if clone == "." {
			clone = orig
		}
		b, err := ioutil.ReadFile(filepath.Join(pwd, orig))
		if err != nil {
			return nil, err
		}
		return nil, ioutil.WriteFile(clone, b, 0o600)
	}

}
