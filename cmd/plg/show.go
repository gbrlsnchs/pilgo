package main

import (
	"fmt"

	"github.com/gbrlsnchs/cli"
	"github.com/gbrlsnchs/pilgo/config"
	"github.com/gbrlsnchs/pilgo/fs"
	"github.com/gbrlsnchs/pilgo/parser"
	"gopkg.in/yaml.v3"
)

type showCmd struct{}

func (*showCmd) register(getcfg func() appConfig) func(cli.Program) error {
	return func(prg cli.Program) error {
		appcfg := getcfg()
		fs := fs.New(appcfg.fs)
		b, err := fs.ReadFile(appcfg.conf)
		if err != nil {
			return err
		}
		var c config.Config
		if yaml.Unmarshal(b, &c); err != nil {
			return err
		}
		userConfigDir, err := appcfg.userConfigDir()
		if err != nil {
			return err
		}
		homeConfigDir, err := appcfg.userHomeDir()
		if err != nil {
			return err
		}
		cwd, err := appcfg.getwd()
		if err != nil {
			return err
		}
		var p parser.Parser
		tr, err := p.Parse(c,
			parser.BaseDirs(map[parser.Mode]string{
				parser.UserMode: userConfigDir,
				parser.HomeMode: homeConfigDir,
			}),
			parser.Cwd(cwd), parser.Envsubst)
		if err != nil {
			return err
		}
		fmt.Fprint(prg.Stdout(), tr)
		return nil
	}
}
