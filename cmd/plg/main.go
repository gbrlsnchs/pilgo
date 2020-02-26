package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/subcommands"
	"gsr.dev/pilgrim/cmd/internal/command"
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
		&showCmd{config: defaultConfig, cwd: cwd},
		command.Name("show"),
		command.Synopsis("Show tree view of files to be symlinked."),
		command.Usage(`show [-check]:
	Show tree view of files to be symlinked.`),
		command.Stdout(os.Stdout),
		command.Stderr(os.Stderr),
	), "")
	flag.Parse()
	return int(cmd.Execute(context.TODO()))
}
