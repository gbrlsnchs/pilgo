package main

import (
	"errors"
	"os"

	"github.com/gbrlsnchs/cli"
	"github.com/gbrlsnchs/cli/cliutil"
	"github.com/gbrlsnchs/pilgo/config"
	"github.com/gbrlsnchs/pilgo/fs"
	"gopkg.in/yaml.v3"
)

var errConfigExists = errors.New("configuration file already exists")

type initCmd struct {
	force   bool
	include cliutil.CommaSepOptionSet
	exclude cliutil.CommaSepOptionSet
	hidden  bool
}

func (cmd *initCmd) register(getcfg func() appConfig) func(cli.Program) error {
	return func(_ cli.Program) error {
		var (
			appcfg = getcfg()
			fs     = fs.New(appcfg.fs)
		)
		fi, err := fs.Stat(appcfg.conf)
		if err != nil {
			return err
		}
		cwd, err := appcfg.getwd()
		if err != nil {
			return err
		}
		files, err := fs.ReadDir(cwd)
		if err != nil {
			return err
		}
		var (
			perm = os.FileMode(0o644)
			c    config.Config
		)
		if fi.Exists() {
			if !cmd.force {
				return errConfigExists
			}
			b, err := fs.ReadFile(appcfg.conf)
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
		cmd.exclude.Set(appcfg.conf)
		opts := []func(*config.Config){
			config.Include(cmd.include),
			config.Exclude(cmd.exclude),
			config.MergeWith(c),
		}
		if cmd.hidden {
			opts = append(opts, config.IncludeHidden)
		}
		c = config.New(targets, opts...)
		b, err := marshalYAML(c)
		if err != nil {
			return err
		}
		return fs.WriteFile(appcfg.conf, b, perm)
	}
}
