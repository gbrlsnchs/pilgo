package main

import (
	"fmt"

	"github.com/gbrlsnchs/cli"
)

type versionCmd struct{}

func (*versionCmd) register(getcfg func() appConfig) func(cli.Program) error {
	return func(prg cli.Program) error {
		appcfg := getcfg()
		exe := prg.Name()
		v := appcfg.version
		if v == "" {
			v = "unknown version"
		}
		fmt.Fprintf(prg.Stdout(), "%s %s\n", exe, v)
		return nil
	}
}
