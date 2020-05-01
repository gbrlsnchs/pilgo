package main

import (
	"os"

	"github.com/gbrlsnchs/cli"
	"github.com/gbrlsnchs/pilgo/cmd/internal"
	"github.com/gbrlsnchs/pilgo/config"
	"github.com/gbrlsnchs/pilgo/fs"
	"github.com/gbrlsnchs/pilgo/fs/fsutil"
)

type appConfig struct {
	name          string
	conf          string
	fs            fs.Driver
	getwd         func() (string, error)
	userConfigDir func() (string, error)
	userHomeDir   func() (string, error)
	version       string
}

func (cfg *appConfig) copy() appConfig { return *cfg }

type rootCmd struct {
	// store
	check   checkCmd
	config  configCmd
	init    initCmd
	link    linkCmd
	scan    scanCmd
	show    showCmd
	version versionCmd
}

func main() {
	os.Exit(run())
}

func run() int {
	var (
		root   rootCmd
		appcfg = appConfig{
			name:          "Pilgo",
			fs:            fsutil.OSDriver{},
			getwd:         os.Getwd,
			userConfigDir: os.UserConfigDir,
			userHomeDir:   os.UserHomeDir,
			version:       internal.Version(),
		}
	)
	cli := cli.New(&cli.Command{
		Options: map[string]cli.Option{
			"config": cli.StringOption{
				OptionDetails: cli.OptionDetails{
					Description: "Use a different configuration file.",
					Short:       'c',
				},
				DefValue:  config.DefaultName,
				Recipient: &appcfg.conf,
			},
		},
		Subcommands: map[string]*cli.Command{
			"check": {
				Description: "Check the status of your dotfiles.",
				Options: map[string]cli.Option{
					"fail": cli.BoolOption{
						OptionDetails: cli.OptionDetails{
							Short:       'f',
							Description: "Return an error if there are any conflicts.",
						},
						DefValue:  false,
						Recipient: &root.check.fail,
					},
					"tags": cli.VarOption{
						OptionDetails: cli.OptionDetails{
							Description: "Comma-separated list of tags. Targets with these tags will also be checked.",
							Short:       't',
							ArgLabel:    "TAG 1,...,TAG n",
						},
						Recipient: &root.check.tags,
					},
				},
				Exec: root.check.register(appcfg.copy),
			},
			"config": {
				Description: "Configure a dotfile in the configuration file.",
				Options: map[string]cli.Option{
					"basedir": cli.StringOption{
						OptionDetails: cli.OptionDetails{
							Description: "Set the target's base directory. Works recursively for all nested targets, unless overridden.",
							ArgLabel:    "DIR",
							Short:       'b',
						},
						Recipient: &root.config.baseDir,
					},
					"link": cli.StringOption{
						OptionDetails: cli.OptionDetails{
							Description: "Set the target's link name.",
							ArgLabel:    "NAME",
							Short:       'l',
						},
						Recipient: &root.config.link,
					},
					"usehome": cli.VarOption{
						OptionDetails: cli.OptionDetails{
							Description: "Use home directory as the target's base directory and recursively for all nested targets, unless overridden.",
							Short:       'H',
						},
						Recipient: &root.config.useHome,
					},
					"flatten": cli.BoolOption{
						OptionDetails: cli.OptionDetails{
							Description: "Prevent the target from being included in the link name.",
							Short:       'f',
						},
						Recipient: &root.config.flatten,
					},
					"tags": cli.VarOption{
						OptionDetails: cli.OptionDetails{
							Description: "Comma-separated list of tags to be set for the target.",
							Short:       't',
							ArgLabel:    "TAG 1,...,TAG n",
						},
						Recipient: &root.config.tags,
					},
				},
				Arg: cli.StringArg{
					Label:     "TARGET",
					Required:  false,
					Recipient: &root.config.file,
				},
				Exec: root.config.register(appcfg.copy),
			},
			"init": {
				Description: "Initialize a configuration file.",
				Options: map[string]cli.Option{
					"force": cli.BoolOption{
						OptionDetails: cli.OptionDetails{
							Description: "Overwrite the existing configuration file.",
							Short:       'f',
						},
						Recipient: &root.init.force,
					},
					"include": cli.VarOption{
						OptionDetails: cli.OptionDetails{
							Description: "File to be exclusively included as a target. Repeat option to include more files.",
							ArgLabel:    "FILE",
						},
						Recipient: &root.init.read.include,
					},
					"exclude": cli.VarOption{
						OptionDetails: cli.OptionDetails{
							Description: "File to be excluded from targets. Repeat option to exclude more files.",
							ArgLabel:    "FILE",
						},
						Recipient: &root.init.read.exclude,
					},
					"hidden": cli.BoolOption{
						OptionDetails: cli.OptionDetails{
							Description: "Include hidden files on initialization.",
							Short:       'H',
						},
						Recipient: &root.init.read.hidden,
					},
				},
				Exec: root.init.register(appcfg.copy),
			},
			"link": {
				Description: "Link your dotfiles as set in the configuration file.",
				Exec:        root.link.register(appcfg.copy),
				Options: map[string]cli.Option{
					"tags": cli.VarOption{
						OptionDetails: cli.OptionDetails{
							Description: "Comma-separated list of tags. Targets with these tags will also be linked.",
							Short:       't',
							ArgLabel:    "TAG 1,...,TAG n",
						},
						Recipient: &root.link.tags,
					},
				},
			},
			"scan": {
				Description: "Set targets by scanning a directory.",
				Options: map[string]cli.Option{
					"include": cli.VarOption{
						OptionDetails: cli.OptionDetails{
							Description: "File to be exclusively included in the scan. Repeat option to include more files.",
							ArgLabel:    "FILE",
						},
						Recipient: &root.scan.read.include,
					},
					"exclude": cli.VarOption{
						OptionDetails: cli.OptionDetails{
							Description: "File to be excluded from the scan. Repeat option to exclude more files.",
							ArgLabel:    "FILE",
						},
						Recipient: &root.scan.read.exclude,
					},
					"hidden": cli.BoolOption{
						OptionDetails: cli.OptionDetails{
							Description: "Include hidden files when scanning.",
							Short:       'H',
						},
						Recipient: &root.scan.read.hidden,
					},
				},
				Arg: cli.StringArg{
					Label:     "TARGET",
					Required:  false,
					Recipient: &root.scan.file,
				},
				Exec: root.scan.register(appcfg.copy),
			},
			"show": {
				Description: "Show your dotfiles in a tree view.",
				Exec:        root.show.register(appcfg.copy),
				Options: map[string]cli.Option{
					"tags": cli.VarOption{
						OptionDetails: cli.OptionDetails{
							Description: "Comma-separated list of tags. Targets with these tags will also be shown.",
							Short:       't',
							ArgLabel:    "TAG 1,...,TAG n",
						},
						Recipient: &root.show.tags,
					},
				},
			},
			"version": {
				Description: "Print version.",
				Exec:        root.version.register(appcfg.copy),
			},
		},
	})
	return cli.ParseAndRun(os.Args)
}
