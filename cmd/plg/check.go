package main

import (
	"errors"
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
	if err := tr.Walk(resolveFunc(ln)); err != nil {
		return err
	}
	// TODO(gbrlsnchs): print errors' details
	fmt.Fprint(stdout, tr)
	return nil
}

func (checkCmd) SetFlags(_ *flag.FlagSet) { /* NOP */ }

func resolveFunc(ln *linker.Linker) func(n *parser.Node) error {
	return func(n *parser.Node) error {
		if err := ln.Resolve(n); err != nil && !isConflict(err) {
			return err
		}
		return nil
	}
}

func isConflict(err error) bool {
	return errors.Is(err, linker.ErrLinkExists) ||
		errors.Is(err, linker.ErrLinkNotExpands) ||
		errors.Is(err, linker.ErrTargetNotExists) ||
		errors.Is(err, linker.ErrTargetNotExpands)
}
