package main

import (
	"errors"
	"os"

	"github.com/gbrlsnchs/cli"
	"github.com/gbrlsnchs/pilgo/config"
	"github.com/gbrlsnchs/pilgo/fs"
)

var errConfigExists = errors.New("configuration file already exists")

type initCmd struct {
	force bool
	read  readMode
}

func (cmd *initCmd) register(getcfg func() appConfig) func(cli.Program) error {
	return func(_ cli.Program) error {
		var (
			appcfg = getcfg()
			fs     = fs.New(appcfg.fs)
		)
		conf := appcfg.conf
		fi, err := fs.Stat(conf)
		if err != nil {
			return err
		}
		fexists := fi.Exists()
		if fexists && !cmd.force {
			return errConfigExists
		}
		cwd, err := appcfg.getwd()
		if err != nil {
			return err
		}
		files, err := fs.ReadDir(cwd)
		if err != nil {
			return err
		}
		cmd.read.exclude.Set(conf)
		targets := cmd.read.resolve(files)
		perm := os.FileMode(0o644)
		if fexists {
			perm = fi.Perm()
		}
		c := &config.Config{Targets: targets}
		b, err := marshalYAML(c)
		if err != nil {
			return err
		}
		return fs.WriteFile(conf, b, perm)
	}
}
