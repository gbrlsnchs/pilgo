package main

import (
	"strings"

	"github.com/gbrlsnchs/cli"
	"github.com/gbrlsnchs/cli/cliutil"
	"github.com/gbrlsnchs/pilgo/fs"
)

type readMode struct {
	include cliutil.CommaSepOptionSet
	exclude cliutil.CommaSepOptionSet
	hidden  bool
}

func (md *readMode) resolve(files []fs.FileInfo) []string {
	eligible := make([]string, 0, len(files))
	for _, fi := range files {
		fname := fi.Name()
		if fname == "" || !md.hidden && strings.HasPrefix(fname, ".") {
			continue
		}
		if len(md.include) > 0 {
			if _, ok := md.include[fname]; !ok {
				continue
			}
		}
		if _, ok := md.exclude[fname]; ok {
			continue
		}
		eligible = append(eligible, fname)
	}
	return eligible
}

func (md *readMode) option(name string) cli.Option {
	switch name {
	case "include":
		return cli.VarOption{
			OptionDetails: cli.OptionDetails{
				Description: "Comma-separated list of targets to be included.",
				ArgLabel:    "TARGET 1,...,TARGET n",
			},
			Recipient: &md.include,
		}
	case "exclude":
		return cli.VarOption{
			OptionDetails: cli.OptionDetails{
				Description: "Comma-separated list of targets to be excluded.",
				ArgLabel:    "TARGET 1,...,TARGET n",
			},
			Recipient: &md.exclude,
		}
	case "hidden":
		return cli.BoolOption{
			OptionDetails: cli.OptionDetails{
				Description: "Include hidden files.",
				Short:       'H',
			},
			Recipient: &md.hidden,
		}
	default:
		panic("unknown option")
	}
}
