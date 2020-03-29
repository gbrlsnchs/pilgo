package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"path/filepath"

	"gopkg.in/yaml.v3"
	"gsr.dev/pilgrim/cmd/internal/command"
	"gsr.dev/pilgrim/config"
	"gsr.dev/pilgrim/fs"
	"gsr.dev/pilgrim/linker"
	"gsr.dev/pilgrim/parser"
)

type linkCmd struct{}

func (linkCmd) Execute(ctx context.Context, stdout, stderr io.Writer) error {
	opts := ctx.Value(command.OptsCtxKey).(opts)
	exe := ctx.Value(command.ErrCtxKey).(string)
	fs := fs.New(opts.fsDriver)
	cwd, err := opts.getwd()
	if err != nil {
		return err
	}
	b, err := fs.ReadFile(filepath.Join(cwd, opts.config))
	if err != nil {
		return err
	}
	var c config.Config
	if yaml.Unmarshal(b, &c); err != nil {
		return err
	}
	baseDir, err := opts.userConfigDir()
	if err != nil {
		return err
	}
	var p parser.Parser
	tr, err := p.Parse(c, parser.BaseDir(baseDir), parser.Cwd(cwd), parser.Envsubst)
	if err != nil {
		return err
	}
	ln := linker.New(fs)
	if err := ln.Link(tr); err != nil {
		var cft *linker.ConflictError
		if errors.As(err, &cft) {
			for _, err := range cft.Errs {
				fmt.Fprintf(stderr, "%s: %v\n", exe, err)
			}
		}
		return err
	}
	return nil
}

func (linkCmd) SetFlags(_ *flag.FlagSet) { /* NOP */ }
