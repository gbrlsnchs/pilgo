package main

import (
	"fmt"
	"regexp"

	"github.com/gbrlsnchs/cli"
)

var rgx = regexp.MustCompile(`^.*?(\d)`)

type versionCmd struct{}

func (*versionCmd) register(getcfg func() appConfig) func(cli.Program) error {
	return func(prg cli.Program) error {
		appcfg := getcfg()
		v := appcfg.version
		if v == "" {
			v = "unknown version"
		} else {
			v = rgx.ReplaceAllString(v, "$1")
		}
		fmt.Fprintf(prg.Stdout(), "%s %s\n", appcfg.name, v)
		return nil
	}
}
