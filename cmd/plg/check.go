package main

import (
	"errors"
	"flag"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
	"gsr.dev/pilgrim/config"
	"gsr.dev/pilgrim/fs"
	"gsr.dev/pilgrim/fs/fsutil"
	"gsr.dev/pilgrim/linker"
	"gsr.dev/pilgrim/parser"
)

var errGotConflicts = errors.New("there are one or more conflicts")

type checkCmd struct {
	fail bool
}

func (cmd checkCmd) Execute(stdout io.Writer, v interface{}) error {
	o := v.(opts)
	fs := fs.New(fsutil.OSDriver{})
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
	var (
		ln      = linker.New(fs)
		errlist []error
	)
	if err := tr.Walk(func(n *parser.Node) error {
		if err := ln.Resolve(n); err != nil {
			if !isConflict(err) {
				return err
			}
			if cmd.fail {
				errlist = append(errlist, err)
			}
		}
		return nil
	}); err != nil {
		return err
	}
	if len(errlist) > 0 {
		for _, err := range errlist {
			fmt.Fprintf(stdout, "\t- %v\n", err)
		}
		return errGotConflicts
	}
	fmt.Fprint(stdout, tr)
	return nil
}

func (cmd *checkCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&cmd.fail, "fail", false, "exit with fail status if there are conflicts")
}

func isConflict(err error) bool {
	return errors.Is(err, linker.ErrLinkExists) ||
		errors.Is(err, linker.ErrLinkNotExpands) ||
		errors.Is(err, linker.ErrTargetNotExists) ||
		errors.Is(err, linker.ErrTargetNotExpands)
}
