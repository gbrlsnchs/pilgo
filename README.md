# Pilgo
[![Linux, macOS and Windows](https://github.com/gbrlsnchs/pilgo/workflows/Linux,%20macOS%20and%20Windows/badge.svg)](https://github.com/gbrlsnchs/pilgo/actions)
[![Latest release](https://img.shields.io/github/v/release/gbrlsnchs/pilgo?include_prereleases)](https://github.com/gbrlsnchs/pilgo/releases/latest)
[![License](https://img.shields.io/github/license/gbrlsnchs/pilgo)](https://github.com/gbrlsnchs/pilgo/blob/master/LICENSE)

1. [Overview](#overview)
    1. [Introduction](#introduction)
    2. [Main features](#main-features)
    3. [Roadmap](#roadmap)
    4. [FAQ](#faq)
1. [Installation](#installation)
    1. [Download from releases](#download-from-releases)
    2. [Install using Go](#install-using-go)
    3. [Linux](#linux)
        1. [Arch Linux](#arch-linux)
2. [Instructions](#instructions)
    1. [Problem](#problem)
    2. [Solution](#solution)
        1. [`init`](#init)
        2. [`show`](#show)
        3. [`config` and `scan`](#config-and-scan)
        4. [`check`](#check)
        5. [`link`](#link)

## Overview
### Introduction
Pilgo is a configuration-based dotfiles manager. That means it'll manage your dotfiles based on what's written in a configuration file, more specifically, a YAML file.

It contains a set of commands to configure the YAML file and, of course, manage your dotfiles.

### Main features
- Everything is expressed via a configuration file
    - Flexibility to choose any directory layout for your dotfiles
- Visualization of the configuration in a nice tree view
- Atomic symlinking

### Roadmap
Progress is being tracked using [GitHub Projects](https://help.github.com/en/github/managing-your-work-on-github/about-project-boards). You can check the roadmap [here](https://github.com/gbrlsnchs/pilgo/projects/1).

### FAQ
The FAQ is part of [Pilgo's Wiki](https://github.com/gbrlsnchs/pilgo/wiki). You can check FAQ [here](https://github.com/gbrlsnchs/pilgo/wiki/FAQ).

## Installation
Pilgo is available to download via a few ways. If you feel Pilgo is missing from a specific repository from some OS, please, feel free to open an issue or a PR.

If you manage to publish it on an OS's repository, please, create an issue or a PR in order for the respective section to be updated.

### Download from releases
There are binaries for Linux, macOS and Windows available as assets in [releases](https://github.com/gbrlsnchs/pilgo/releases).

### Install using Go
You can install using Go 1.13 or greater:
```console
$ go get -u github.com/gbrlsnchs/pilgo/cmd/plg
```

### Linux
#### Arch Linux
[![pilgo on AUR](https://img.shields.io/aur/version/pilgo?label=pilgo)](https://aur.archlinux.org/packages/pilgo/)
[![pilgo-bin on AUR](https://img.shields.io/aur/version/pilgo-bin?label=pilgo-bin)](https://aur.archlinux.org/packages/pilgo-bin/)

Pilgo is available on the [AUR](https://wiki.archlinux.org/index.php/Arch_User_Repository):
- [pilgo](https://aur.archlinux.org/packages/pilgo/)
- [pilgo-bin](https://aur.archlinux.org/packages/pilgo-bin/) (binary package)

You can install it using your [AUR helper](https://wiki.archlinux.org/index.php/AUR_helpers) of choice.

Example:
```console
$ yay -Sy pilgo
```

## Instructions
### Problem
Imagine you organize your dotfiles with the following structure:
```console
$ tree
.
├── alacritty
│   └── alacritty.yml
├── bspwm
│   └── bspwmrc
├── dunst
│   └── dunstrc
├── mpd
│   └── mpd.conf
├── mpv
│   └── mpv.conf
└── zsh
    ├── zprofile
    └── zshrc

6 directories, 7 files
```

You back up your dotfiles using Git and, when you clone the repository to restore them, you want to create symlinks of your dotfiles in the appropriate places. Also, you want it to be reproducible across anywhere you clone your dotfiles.

### Solution
#### `init`
You can use Pilgo to create symlinks for you. If you were to use GNU Stow, you'd have to restructure your directory layout. That's not the case for Pilgo, because it is a configuration-based tool.

For that reason, you'll need to initialize a configuration file in your dotfiles directory:
```console
$ plg init
```

<kbd>**Hint:**</kbd> <small>Run `plg init -h` to check all options for `init`.</small>

After running the `init` command, Pilgo creates a YAML file called `pilgo.yml`. It has read the dotfiles directory and listed all files inside it as targets to be symlinked:
```console
$ cat pilgo.yml
targets:
- alacritty
- bspwm
- dunst
- mpd
- mpv
- zsh
```

#### `show`
After the configuration has been created, you can visualize it in a tree view:
```console
$ plg show
.
├── alacritty <- /home/me/.config/alacritty
├── bspwm     <- /home/me/.config/bspwm
├── dunst     <- /home/me/.config/dunst
├── mpd       <- /home/me/.config/mpd
├── mpv       <- /home/me/.config/mpv
└── zsh       <- /home/me/.config/zsh
```

#### `config` and `scan`
If you have ever used Zsh, you'll notice that the configuration is not quite right. That happens because Pilgo creates the configuration file following sane defaults, that is:
- It uses `~/.config` (or the equivalent for other OSes) as the base directory for symlinks
- It assumes symlinks have the same name as their targets

The nice part is you don't need to change your directory layout because of that. You simply fine-tune your Pilgo configuration using the proper command:
```console
$ plg scan zsh
$ plg config -usehome -flatten zsh
```

With the `scan` command, Pilgo scans `zsh` to add its files as its targets.
With the `config` command, we set two properties for the `zsh` directory:
- `-usehome` sets it to use the home directory as the base directory (instead of `~/.config` or equivalent)
- `-flatten` skips adding the `zsh` to the symlink path, skipping directly to its children

Note that `config` modifies `pilgo.yml` for you. Here's how it is after the modification:
```console
$ cat pilgo.yml
targets:
- alacritty
- bspwm
- dunst
- mpd
- mpv
- zsh
options:
  zsh:
    flatten: true
    useHome: true
    targets:
    - zprofile
    - zshrc
```

You can also visualize the new configuration:
```console
$ plg show
.
├── alacritty    <- /home/me/.config/alacritty
├── bspwm        <- /home/me/.config/bspwm
├── dunst        <- /home/me/.config/dunst
├── mpd          <- /home/me/.config/mpd
├── mpv          <- /home/me/.config/mpv
└── zsh       
    ├── zprofile <- /home/me/zprofile
    └── zshrc    <- /home/me/zshrc
```

<kbd>**Hint**:</kbd> <small>If you're not sure what has been configured, you can always run `plg show` or even open `pilgo.yml` and check its content.</small>

However, we need to set one last thing. Both `zprofile` and `zshrc` need to have their symlinks prepended with a dot:
```console
$ plg config -link=.zprofile zsh/zprofile
$ plg config -link=.zshrc zsh/zshrc
$ plg show
.
├── alacritty    <- /home/me/.config/alacritty
├── bspwm        <- /home/me/.config/bspwm
├── dunst        <- /home/me/.config/dunst
├── mpd          <- /home/me/.config/mpd
├── mpv          <- /home/me/.config/mpv
└── zsh       
    ├── zprofile <- /home/me/.zprofile
    └── zshrc    <- /home/me/.zshrc
```

You're done with fine-tuning the configuration. You'll probably not have to change it again for some time. You'll only need to fine-tune it again if you add files with restrictions similar to Zsh's.

#### `check`
Finally, you can link your dotfiles. But before that, you can also check whether your dotfiles are ready to be symlinked, which means there are no conflicts. You can check your files by using the `check` command:
```console
$ plg check
.
├── alacritty    <- /home/me/.config/alacritty     (READY)
├── bspwm        <- /home/me/.config/bspwm         (DONE)
├── dunst                                          (EXPAND)
│   └── dunstrc  <- /home/me/.config/dunst/dunstrc (READY)
├── mpd          <- /home/me/.config/mpd           (ERROR)
├── mpv          <- /home/me/.config/mpv           (DONE)
└── zsh                                            (SKIP)
    ├── zprofile <- /home/me/.zprofile             (READY)
    └── zshrc    <- /home/me/.zshrc                (CONFLICT)
```

`check` is just a preview of how Pilgo will handle your dotfiles. We can note some changes in the output, specially after each symlink name.

Those are states. Each state means something different:
- `READY` means the target can be symlinked without further issues
- `DONE` means the targets is already correctly symlinked, that is, the symlink already points to the same target configured in your Pilgo configuration
- `EXPAND` means a directory exists where Pilgo would create the symlink, but since the target is also a directory, Pilgo can expand it and then symlink files inside it
- `ERROR` means there's something wrong with your target or with your symlink
- `CONFLICT` means one of the following occured:
    - A regular file already exists where your target would be symlinked and it can't be expanded
    - A regular file already exists where your target would be symlinked and your target can't be expanded

Note that Pilgo doesn't solve conflicts automatically, since it could be a destructive action prone to user error. You have to manually resolve conflicts, which consists of removing files from where symlinks would be created.

<kbd>**Hint:**</kbd> <small>You can have more details for errors by running `plg check -fail`.</small>

#### `link`
Lastly, if there are no conflicts or errors, you can simply run:
```console
$ plg link
```

<kbd>**Hint:**</kbd> <small>The `link` command always checks all dotfiles before linking, so you don't end up with only half of them symlinked. If there are conflicts or errors, it will return an error status and abort.</small>

And if you check again, you'll see:
```console
.
├── alacritty    <- /home/me/.config/alacritty     (DONE)
├── bspwm        <- /home/me/.config/bspwm         (DONE)
├── dunst                                          (EXPAND)
│   └── dunstrc  <- /home/me/.config/dunst/dunstrc (DONE)
├── mpd          <- /home/me/.config/mpd           (DONE)
├── mpv          <- /home/me/.config/mpv           (DONE)
└── zsh                                            (SKIP)
    ├── zprofile <- /home/me/.zprofile             (DONE)
    └── zshrc    <- /home/me/.zshrc                (DONE)
```

You can commit `pilgo.yml` to your dotfiles repository and, since you have already configured everything, next time you have to symlink things, you just have to run `plg link`.
