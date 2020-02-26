package main

import (
	"flag"
	"fmt"
	"io"

	"gsr.dev/pilgrim/linker"
	"gsr.dev/pilgrim/osfs"
)

type checkCmd struct {
	config string
	cwd    string
}

func (cmd checkCmd) Execute(stdout io.Writer) error {
	tr, err := buildTree(cmd.config, cmd.cwd)
	if err != nil {
		return err
	}
	var (
		fs osfs.FileSystem
		ln = linker.New(fs)
	)
	if err := tr.Walk(ln.Resolve); err != nil {
		return err
	}
	// TODO(gbrlsnchs): print errors' details
	fmt.Fprint(stdout, tr)
	return nil
}

func (checkCmd) SetFlags(_ *flag.FlagSet) { /* NOP */ }
