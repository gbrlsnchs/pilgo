package main

import (
	"context"
	"flag"
	"os"
	"path/filepath"

	"github.com/google/subcommands"
	"gsr.dev/pilgrim/cmd/internal/command"
	"gsr.dev/pilgrim/config"
	"gsr.dev/pilgrim/fs"
	"gsr.dev/pilgrim/fs/fsutil"
)

const defaultConfig = config.DefaultName

func main() {
	os.Exit(run())
}

func run() int {
	exe := filepath.Base(os.Args[0])
	// NOTE: subcommands.NewCommander must be used in order to avoid bugs
	// when this function is tested using "go-cmdtest".
	//
	// The bug consists of having subcommands.DefaultCommand register homonymous
	// commands more than once for a single ".ct" file.
	//
	// Since subcommands.DefaultCommand is global, those commands are appended to the same command
	// slice and, when subcommands.Execute is called, each command is checked by name within that
	// command is checked by name within that slice and, since they're homonyms, only the first one
	// is picked. Because each individual command in a ".ct" file gets new stdout and stderr
	// redirection in order to be checked, each run writes to writers from the first run.
	//
	// This note can be removed when https://github.com/google/subcommands/issues/32 gets resolved.
	cmd := subcommands.NewCommander(flag.CommandLine, exe)
	commands := []struct {
		command.Interface
		name     string
		synopsis string
	}{
		{showCmd{}, "show", "show tree view of files to be symlinked"},
		{&checkCmd{}, "check", "check symlinks and show them in a tree view"},
		{&initCmd{}, "init", "initialize a configuration file"},
		{&configCmd{}, "config", "configure a file's options"},
		{linkCmd{}, "link", "create symlinks"},
	}
	for _, c := range commands {
		cmd.Register(command.New(
			c.Interface,
			command.Name(c.name),
			command.Synopsis(c.synopsis),
			command.Stdout(os.Stdout),
			command.Stderr(os.Stderr),
		), "")
	}
	flag.Parse()
	ctx := context.TODO()
	ctx = context.WithValue(ctx, command.ErrCtxKey, exe)
	ctx = context.WithValue(ctx, command.OptsCtxKey, opts{
		config:        defaultConfig,
		fsDriver:      fsutil.OSDriver{},
		getwd:         os.Getwd,
		userConfigDir: os.UserConfigDir,
	})
	status := cmd.Execute(ctx)
	return int(status)
}

type opts struct {
	config        string
	fsDriver      fs.Driver
	getwd         func() (string, error)
	userConfigDir func() (string, error)
}
