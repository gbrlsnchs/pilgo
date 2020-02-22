package main

import (
	"flag"
	"path/filepath"
	"testing"

	cmdtest "github.com/google/go-cmdtest"
)

var update = flag.Bool("update", false, "update test files with results")

func TestCLI(t *testing.T) {
	ts, err := cmdtest.Read(filepath.Join("testdata", t.Name()))
	if err != nil {
		t.Fatal(err)
	}
	ts.Run(t, *update)
}
