package main

import (
	"context"
	"errors"
	"flag"
	"io"
	"os"

	"github.com/gbrlsnchs/pilgo/cmd/internal/command"
	"github.com/gbrlsnchs/pilgo/config"
	"github.com/gbrlsnchs/pilgo/fs"
	"gopkg.in/yaml.v3"
)

var errConfigExists = errors.New("configuration file already exists")

type initCmd struct {
	force   bool
	include commaset
	exclude commaset
}

func (cmd initCmd) Execute(ctx context.Context, stdout, _ io.Writer) error {
	o := ctx.Value(command.OptsCtxKey).(opts)
	var (
		fs = fs.New(o.fsDriver)
		c  config.Config
	)
	fi, err := fs.Stat(o.config)
	if err != nil {
		return err
	}
	cwd, err := o.getwd()
	if err != nil {
		return err
	}
	files, err := fs.ReadDir(cwd)
	if err != nil {
		return err
	}
	perm := os.FileMode(0o644)
	if fi.Exists() {
		if !cmd.force {
			return errConfigExists
		}
		b, err := fs.ReadFile(o.config)
		if err != nil {
			return err
		}
		if err = yaml.Unmarshal(b, &c); err != nil {
			return err
		}
		perm = fi.Perm()
	}
	targets := make([]string, len(files))
	for i, fi := range files {
		targets[i] = fi.Name()
	}
	b, err := marshalYAML(c.Init(targets, cmd.include, cmd.exclude))
	if err != nil {
		return err
	}
	return fs.WriteFile(o.config, b, perm)
}

func (cmd *initCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&cmd.force, "force", false, "override targets from existing configuration file")
	f.Var(&cmd.include, "include", "comma-separated list of targets to be included")
	f.Var(&cmd.exclude, "exclude", "comma-separated list of targets to be excluded")
}
