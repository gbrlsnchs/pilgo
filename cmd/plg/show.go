package main

import (
	"flag"
	"fmt"
	"io"
)

type showCmd struct {
	config string
	cwd    string
}

func (cmd showCmd) Execute(stdout io.Writer) error {
	tr, err := buildTree(cmd.config, cmd.cwd)
	if err != nil {
		return err
	}
	// TODO(gbrlsnchs): print errors' details
	fmt.Fprint(stdout, tr)
	return nil
}

func (showCmd) SetFlags(_ *flag.FlagSet) { /* NOP */ }
