package main

import (
	"flag"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
	"gsr.dev/pilgrim"
	"gsr.dev/pilgrim/fs/osfs"
	"gsr.dev/pilgrim/parser"
)

type showCmd struct {
	config string
	cwd    string
}

func (cmd showCmd) Execute(stdout io.Writer) error {
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
	fmt.Fprint(stdout, tr)
	return nil
}

func (showCmd) SetFlags(_ *flag.FlagSet) { /* NOP */ }
