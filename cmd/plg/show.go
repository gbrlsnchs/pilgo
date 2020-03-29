package main

import (
	"context"
	"flag"
	"fmt"
	"io"

	"github.com/gbrlsnchs/pilgo/cmd/internal/command"
	"github.com/gbrlsnchs/pilgo/config"
	"github.com/gbrlsnchs/pilgo/fs"
	"github.com/gbrlsnchs/pilgo/parser"
	"gopkg.in/yaml.v3"
)

type showCmd struct{}

func (showCmd) Execute(ctx context.Context, stdout, _ io.Writer) error {
	o := ctx.Value(command.OptsCtxKey).(opts)
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
