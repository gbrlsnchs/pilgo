package main

import (
	"errors"
	"flag"
	"io"
	"os"

	"gopkg.in/yaml.v3"
	"gsr.dev/pilgrim"
	"gsr.dev/pilgrim/fs/osfs"
)

var errConfigExists = errors.New("configuration file already exists")

type initCmd struct {
	force bool
}

func (cmd initCmd) Execute(_ io.Writer, v interface{}) error {
	o := v.(opts)
	var (
		fs osfs.FileSystem
		c  pilgrim.Config
	)
	fi, err := fs.Info(o.config)
	if err != nil {
		return err
	}
	cwd, err := o.getwd()
	if err != nil {
		return err
	}
	targets, err := fs.ReadDir(cwd)
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
	b, err := marshalYAML(c.Init(targets))
	if err != nil {
		return err
	}
	return fs.WriteFile(o.config, b, perm)
}

func (cmd *initCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&cmd.force, "force", false, "override targets from existing configuration file")
}
