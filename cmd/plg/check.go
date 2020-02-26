package main

import (
	"flag"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
	"gsr.dev/pilgrim"
	"gsr.dev/pilgrim/fs/osfs"
	"gsr.dev/pilgrim/linker"
	"gsr.dev/pilgrim/parser"
)

type checkCmd struct {
	config string
	cwd    string
}

func (cmd checkCmd) Execute(stdout io.Writer) error {
	var fs osfs.FileSystem
	b, err := fs.ReadFile(cmd.config)
	if err != nil {
		return err
	}
	var c pilgrim.Config
	if yaml.Unmarshal(b, &c); err != nil {
		return err
	}
	var p parser.Parser
	tr, err := p.Parse(c, parser.Cwd(cmd.cwd))
	if err != nil {
		return err
	}
	ln := linker.New(fs)
	if err := tr.Walk(ln.Resolve); err != nil {
		return err
	}
	// TODO(gbrlsnchs): print errors' details
	fmt.Fprint(stdout, tr)
	return nil
}

func (checkCmd) SetFlags(_ *flag.FlagSet) { /* NOP */ }
