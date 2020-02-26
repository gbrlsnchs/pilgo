package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/subcommands"
	"gsr.dev/pilgrim"
	"gsr.dev/pilgrim/cmd/internal/command"
	"gsr.dev/pilgrim/parser"
)

const defaultConfig = "pilgrim.json"

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
		command.Synopsis("Show tree view of files to be symlinked."),
		command.Usage(`show:
	Show tree view of files to be symlinked.`),
		command.Stdout(os.Stdout),
		command.Stderr(os.Stderr),
	), "")
	cmd.Register(command.New(
		checkCmd{config: defaultConfig, cwd: cwd},
		command.Name("check"),
		command.Synopsis("Check symlinks and show them in a tree view."),
		command.Usage(`check:
	Check symlinks and show them in a tree view.`),
		command.Stdout(os.Stdout),
		command.Stderr(os.Stderr),
	), "")
	flag.Parse()
	return int(cmd.Execute(context.TODO()))
}

func buildTree(config, cwd string) (*parser.Tree, error) {
	b, err := ioutil.ReadFile(config)
	if err != nil {
		return nil, err
	}
	var c pilgrim.Config
	if err = json.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	var p parser.Parser
	tr, err := p.Parse(c, parser.Cwd(cwd))
	if err != nil {
		return nil, err
	}
	return tr, err
}
