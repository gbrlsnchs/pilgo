package main

import (
	"github.com/gbrlsnchs/cli"
	"github.com/gbrlsnchs/pilgo/config"
	"gopkg.in/yaml.v3"
)

type scanCmd struct {
	file string
	read readMode
}

func (cmd *scanCmd) register(getcfg func() appConfig) cli.ExecFunc {
	return func(_ cli.Program) error {
		appcfg := getcfg()
		fs := appcfg.fs
		conf := appcfg.conf
		b, err := fs.ReadFile(conf)
		if err != nil {
			return err
		}
		files, err := fs.ReadDir(cmd.file)
		if err != nil {
			return err
		}
		cmd.read.exclude.Set(conf)
		targets := cmd.read.resolve(files)
		var c config.Config
		if err := yaml.Unmarshal(b, &c); err != nil {
			return err
		}
		cc := &config.Config{Targets: targets}
		c.Set(cmd.file, cc, config.ModeScan)
		if b, err = marshalYAML(c); err != nil {
			return err
		}
		fi, err := fs.Stat(conf)
		if err != nil {
			return err
		}
		return fs.WriteFile(conf, b, fi.Perm())
	}
}
