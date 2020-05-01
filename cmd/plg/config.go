package main

import (
	"strconv"

	"github.com/gbrlsnchs/cli"
	"github.com/gbrlsnchs/cli/cliutil"
	"github.com/gbrlsnchs/pilgo/config"
	"github.com/gbrlsnchs/pilgo/fs"
	"gopkg.in/yaml.v3"
)

type configCmd struct {
	file    string
	baseDir string
	link    string
	useHome boolptr
	flatten bool
	tags    cliutil.CommaSepOptionList
}

func (cmd *configCmd) register(getcfg func() appConfig) func(cli.Program) error {
	return func(_ cli.Program) error {
		appcfg := getcfg()
		fs := fs.New(appcfg.fs)
		conf := appcfg.conf
		b, err := fs.ReadFile(conf)
		if err != nil {
			return err
		}
		var c config.Config
		if err := yaml.Unmarshal(b, &c); err != nil {
			return err
		}
		cc := &config.Config{
			BaseDir: cmd.baseDir,
			Link:    cmd.link,
			Flatten: cmd.flatten,
			UseHome: cmd.useHome.addr,
			Tags:    cmd.tags,
		}
		c.Set(cmd.file, cc, config.ModeConfig)
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

type boolptr struct {
	addr *bool
}

func (bp *boolptr) Set(value string) error {
	b, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}
	bp.addr = &b
	return nil
}

func (bp *boolptr) String() string {
	if bp.addr == nil {
		return ""
	}
	return strconv.FormatBool(*bp.addr)
}

func (bp *boolptr) IsBoolFlag() bool { return true }
