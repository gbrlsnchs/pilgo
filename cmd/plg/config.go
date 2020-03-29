package main

import (
	"context"
	"flag"
	"fmt"
	"io"

	"github.com/gbrlsnchs/pilgo/cmd/internal/command"
	"github.com/gbrlsnchs/pilgo/config"
	"github.com/gbrlsnchs/pilgo/fs"
	"gopkg.in/yaml.v3"
)

type configCmd struct {
	file    string
	baseDir string
	link    strptr
	targets commalist
}

func (cmd configCmd) Execute(ctx context.Context, stdout, _ io.Writer) error {
	o := ctx.Value(command.OptsCtxKey).(opts)
	fs := fs.New(o.fsDriver)
	b, err := fs.ReadFile(o.config)
	if err != nil {
		return err
	}
	var c config.Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return err
	}
	c.Set(cmd.file, config.Config{
		BaseDir: cmd.baseDir,
		Link:    cmd.link.addr,
		Targets: cmd.targets,
	})
	if b, err = marshalYAML(c); err != nil {
		return err
	}
	fi, err := fs.Stat(o.config)
	if err != nil {
		return err
	}
	return fs.WriteFile(o.config, b, fi.Perm())
}

func (cmd *configCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&cmd.file, "file", "", "file to be configured")
	f.StringVar(&cmd.baseDir, "basedir", "", "set field \"baseDir\"")
	f.Var(&cmd.link, "link", "set field \"link\"")
	f.Var(&cmd.targets, "targets", "comma-separated list of targets")
}

type strptr struct {
	addr *string
}

func (sp *strptr) Set(value string) error {
	sp.addr = &value
	return nil
}

func (sp strptr) String() string {
	if sp.addr == nil {
		return fmt.Sprint(sp.addr)
	}
	return *sp.addr
}
