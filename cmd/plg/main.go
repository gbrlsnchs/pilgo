package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/subcommands"
	"gopkg.in/yaml.v3"
	"gsr.dev/pilgrim/cmd/internal/command"
)

// TODO(gbrlsnchs): use `pilgrim.DefaultConfig`
const defaultConfig = "pilgrim.yml"

func main() {
	os.Exit(run())
}

func run() int {
	exe := os.Args[0]
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", exe, err)
		return int(subcommands.ExitFailure)
	}
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
	cmd := subcommands.NewCommander(flag.CommandLine, filepath.Base(exe))
	cmd.Register(command.New(
		showCmd{config: defaultConfig, cwd: cwd},
		command.Name("show"),
		command.Synopsis("show tree view of files to be symlinked"),
		command.Stdout(os.Stdout),
		command.Stderr(os.Stderr),
	), "")
	cmd.Register(command.New(
		checkCmd{config: defaultConfig, cwd: cwd},
		command.Name("check"),
		command.Synopsis("check symlinks and show them in a tree view"),
		command.Stdout(os.Stdout),
		command.Stderr(os.Stderr),
	), "")
	cmd.Register(command.New(
		&initCmd{config: defaultConfig, cwd: cwd},
		command.Name("init"),
		command.Synopsis("initialize a configuration file"),
		command.Stdout(os.Stdout),
		command.Stderr(os.Stderr),
	), "")
	cmd.Register(command.New(
		&configCmd{config: defaultConfig, cwd: cwd},
		command.Name("config"),
		command.Synopsis("configure a file's options"),
		command.Stdout(os.Stdout),
		command.Stderr(os.Stderr),
	), "")
	flag.Parse()
	ctx := context.TODO()
	return int(cmd.Execute(
		context.WithValue(ctx, command.ErrCtxKey, exe),
	))
}

func marshalYAML(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
