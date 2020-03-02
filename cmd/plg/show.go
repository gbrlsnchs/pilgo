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

type showCmd struct{}

func (showCmd) Execute(stdout io.Writer, v interface{}) error {
	o := v.(opts)
	var fs osfs.FileSystem
	b, err := fs.ReadFile(o.config)
	if err != nil {
		return err
	}
	var c pilgrim.Config
	if yaml.Unmarshal(b, &c); err != nil {
		return err
	}
	cwd, err := o.getwd()
	if err != nil {
		return err
	}
	var p parser.Parser
	tr, err := p.Parse(c, parser.Cwd(cwd))
	if err != nil {
		return err
	}
	fmt.Fprint(stdout, tr)
	return nil
}

func (showCmd) SetFlags(_ *flag.FlagSet) { /* NOP */ }
