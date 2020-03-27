package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
	"gsr.dev/pilgrim/cmd/internal/command"
	"gsr.dev/pilgrim/config"
	"gsr.dev/pilgrim/fs"
	"gsr.dev/pilgrim/linker"
	"gsr.dev/pilgrim/parser"
)

type checkCmd struct {
	fail bool
}

func (cmd checkCmd) Execute(ctx context.Context, stdout, stderr io.Writer) error {
	o := ctx.Value(command.OptsCtxKey).(opts)
	exe := ctx.Value(command.ErrCtxKey).(string)
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
	ln := linker.New(fs)
	if err = ln.Resolve(tr); err != nil {
		var cft *linker.ConflictError
		if errors.As(err, &cft) {
			if !cmd.fail {
				goto printtree
			}
			for _, err := range cft.Errs {
				fmt.Fprintf(stderr, "%s: %v\n", exe, err)
			}
		}
		return err
	}
printtree:
	fmt.Fprint(stdout, tr)
	return nil
}

func (cmd *checkCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&cmd.fail, "fail", false, "exit with fail status if there are conflicts")
}
