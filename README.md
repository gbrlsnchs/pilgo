# Pilgo
[![Linux, macOS and Windows](https://github.com/gbrlsnchs/pilgo/workflows/Linux,%20macOS%20and%20Windows/badge.svg)](https://github.com/gbrlsnchs/pilgo/actions)

## What
Pilgo is a tool for managing dotfiles via symlinks either through a CLI (`plg`) or a Go framework.

### Installation
#### Downloading binary from releases
For now, binaries for Linux, macOS and Windows are available as assets in [releases](https://github.com/gbrlsnchs/pilgo/releases).

#### Installing using Go
You need Go 1.13 or greater in order to build and install Pilgo:

```console
$ go get -u github.com/gbrlsnchs/pilgo/cmd/plg
```

#### Installing via package managers
Currently, Pilgo is not available on any package repositories. If you manage to publish it on the respective repository for the platform you use, please open an issue or a pull request in order for this section to be updated!

## Why
Because using GNU Stow is limiting and you want visual feedback of how things are configured.

### Use case
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

You back up your dotfiles using Git and, when you clone the repository to restore them, you either use other tools or custom scripts.

### The problem
You wish to correctly symlink those files, but you don't want to limit yourself for naming them (like when using GNU Stow).

Also, you want it to be reproducible across anywhere you're installing your configuration, including on Windows (why though).

## How

OK, in the previous example, for all directories except `zsh`, you want to symlink them inside `$XDG_CONFIG_HOME` (or your OS equivalent).

For `zsh`, you can't symlink the whole directory, since you're going to put files in `$HOME`. Therefore, you need to symlink them individually (no, you don't want to symlink `zsh` as your `$HOME`, please no).

### Creating a configuration file
Pilgo uses a configuration file to manage your dotfiles. The configuration file is a simple YAML file called `pilgo.yml`.

Pilgo can initialize a configuration file for you, and by default includes all eligible files in the current directory ():

```console
$ plg init
```
<kbd>**Hint:**</kbd> To list all possible flags, run `plg init -h`.

Now, here's our repository structure after running `plg init`:
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
├── pilgo.yml
└── zsh
    ├── zprofile
    └── zshrc

6 directories, 8 files
```

Notice that we now have a file called `pilgo.yml`, and here's what's within it:
```yaml
targets:
- alacritty
- bspwm
- dunst
- mpd
- mpv
- zsh
```


### Listing files
OK, configuration created, let's **visualize** what is going to happen with the current configuration we have:
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

The output is very self-explanatory: the tree structure shows your dotfiles listed in the `pilgo.yml` file, the arrow shows to where links will point after you create symlinks using Pilgo.

<kbd>**Hint:**</kbd> The default base directory is `$XDG_BASE_DIRECTORY` for Linux distros, `$HOME/Library/Application Support` for macOS and `%AppData%` on Windows.

Well, it seems right except for `zsh`. Let's fix it, shall we?

### Configuring files
So, for `zsh`, we need to change the following:
- Its base directory must be set to home directory
- It shouldn't be symlinked (only its children, individually)
- It should have two targets, `zprofile` and `zshrc`

To do so, we run:
```console
$ plg config -home -flatten zsh
```

<kbd>**Hint:**</kbd> Pilgo substitutes environment variables in order for your `pilgo.yml` to be more portable.

As said before, one advantage of using Pilgo is that you can name files however you want and then configure them in `pilgo.yml` to have a custom symlink name, not needing to name files with an initial dot.

For both `zprofile` and `zshrc`, we'll need to configure them to have a custom name when symlinked:
```console
$ plg config -link=.zprofile zsh/zprofile
$ plg config -link=.zshrc zsh/zshrc
```

And now, if we run the `show` command again:
```console
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

All good, but this is just a projection of what will be done. Can those files really be symlinked without any further issues, like, for example, does a file already exist where we wish to create a link?

Let's run the `check` command to, well... check our tree:
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

OK, some statuses showed up in the right part of the output. Here is what these statuses mean:
  - `READY` means everything is OK for the file to be symlinked
  - `DONE` means the file is already symlinked to the target configured in `pilgo.yml`
  - `EXPAND` means a directory already exists where we want to create our symlink, so it recursively creates symlinks for the target's children if it's a directory and puts them inside the existing directory
  - `SKIP` means the file is just a directory to hold other targets, so it doesn't get symlinked (like in the `zsh` example, where the directory `zsh` was just used to wrap ZSH-related files)
  - `CONFLICT` means it's impossible to expand either the target, the link or both, and the file in place of where our symlink should be created has no relation with it
  - `ERROR` means things broke and gone wrong

### Creating symlinks
Finally, after you visualized what's going to be done, it's time to symlink. Note that, before symlinking, Pilgo checks for conflicts and errors, so you won't have only half of your dotfiles directory symlinked.


It's a two step process, first check, then symlink:
```console
$ plg link
```

If there's an error, the command above should fail and return an exit code greater than zero.

Otherwise, you're done! You only need to configure `pilgo.yml` once. After that, commit it along your dotfiles (in the root directory) and use Pilgo to deal with your dotfiles in other environments.
