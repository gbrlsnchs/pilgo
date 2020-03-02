package main

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"gopkg.in/yaml.v3"
	"gsr.dev/pilgrim"
	"gsr.dev/pilgrim/fs/osfs"
)

type configCmd struct {
	config string
	cwd    string

	file    string
	baseDir string
	link    strptr
	targets targetList
}

func (cmd configCmd) Execute(_ io.Writer) error {
	var fs osfs.FileSystem
	b, err := fs.ReadFile(cmd.config)
	if err != nil {
		return err
	}
	var c pilgrim.Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return err
	}
	c.Set(cmd.file, pilgrim.Config{
		BaseDir: cmd.baseDir,
		Link:    cmd.link.addr,
		Targets: cmd.targets,
	})
	if b, err = marshalYAML(c); err != nil {
		return err
	}
	fi, err := fs.Info(cmd.config)
	if err != nil {
		return err
	}
	return fs.WriteFile(cmd.config, b, fi.Perm())
}

func (cmd *configCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&cmd.file, "file", "", "file to be configured")
	f.StringVar(&cmd.baseDir, "basedir", "", "set field \"baseDir\"")
	f.Var(&cmd.link, "link", "set field \"link\"")
	f.Var(&cmd.targets, "targets", "comma-separated list of targets")
}

type targetList []string

func (tgl *targetList) Set(value string) error {
	*tgl = strings.Split(value, ",")
	return nil
}

func (tgl targetList) String() string { return strings.Join(tgl, ",") }

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
