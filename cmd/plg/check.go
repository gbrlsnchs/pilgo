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

type checkCmd struct{}

func (checkCmd) Execute(stdout io.Writer, v interface{}) error {
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
	baseDir, err := o.userConfigDir()
	if err != nil {
		return err
	}
	cwd, err := o.getwd()
	if err != nil {
		return err
	}
	var p parser.Parser
	tr, err := p.Parse(c, parser.BaseDir(baseDir), parser.Cwd(cwd), parser.Envsubst)
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
