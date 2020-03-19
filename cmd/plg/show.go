package main

import (
	"flag"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
	"gsr.dev/pilgrim/config"
	"gsr.dev/pilgrim/fs"
	"gsr.dev/pilgrim/parser"
)

type showCmd struct{}

func (showCmd) Execute(stdout io.Writer, v interface{}) error {
	o := v.(opts)
	fs := fs.New(o.fsDriver)
	b, err := fs.ReadFile(o.config)
	if err != nil {
		return err
	}
	var c config.Config
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
	fmt.Fprint(stdout, tr)
	return nil
}

func (showCmd) SetFlags(_ *flag.FlagSet) { /* NOP */ }
