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
	link    strptr
	targets cliutil.CommaSepOptionList
	useHome boolptr
}

func (cmd *configCmd) register(getcfg func() appConfig) func(cli.Program) error {
	return func(_ cli.Program) error {
		appcfg := getcfg()
		fs := fs.New(appcfg.fs)
		b, err := fs.ReadFile(appcfg.conf)
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
			UseHome: cmd.useHome.addr,
		})
		if b, err = marshalYAML(c); err != nil {
			return err
		}
		fi, err := fs.Stat(appcfg.conf)
		if err != nil {
			return err
		}
		return fs.WriteFile(appcfg.conf, b, fi.Perm())
	}
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
		return ""
	}
	return *sp.addr
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
