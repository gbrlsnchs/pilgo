package main

import (
	"context"
	"flag"
	"fmt"
	"io"

	"github.com/gbrlsnchs/pilgo/cmd/internal/command"
)

type versionCmd struct {
	v string
}

func (cmd versionCmd) Execute(ctx context.Context, stdout, _ io.Writer) error {
	exe := ctx.Value(command.ErrCtxKey).(string)
	version := cmd.v
	if version == "" {
		version = "unknown version"
	}
	fmt.Fprintf(stdout, "%s %s\n", exe, version)
	return nil
}

func (versionCmd) SetFlags(_ *flag.FlagSet) { /* NOP */ }
