package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"

	"gsr.dev/pilgrim"
	"gsr.dev/pilgrim/linker"
	"gsr.dev/pilgrim/osfs"
	"gsr.dev/pilgrim/parser"
)

type showCmd struct {
	config string

	// flags
	check bool
	cwd   string
}

// Execute builds a tree of symlinks based on a configuration file
// and writes it to stdout.
func (cmd showCmd) Execute(stdout io.Writer) error {
	b, err := ioutil.ReadFile(cmd.config)
	if err != nil {
		return err
	}
	var c pilgrim.Config
	if err = json.Unmarshal(b, &c); err != nil {
		return err
	}
	var p parser.Parser
	tr, err := p.Parse(c, parser.Cwd(cmd.cwd))
	if err != nil {
		return err
	}
	if cmd.check {
		ln := linker.New(osfs.FileSystem{})
		if err := tr.Walk(ln.Resolve); err != nil {
			return err
		}
	}
	// TODO(gbrlsnchs): print errors' details
	fmt.Fprint(stdout, tr)
	return nil
}

// SetFlags resolves flags.
func (cmd *showCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&cmd.check, "check", false, "check all files")
}
