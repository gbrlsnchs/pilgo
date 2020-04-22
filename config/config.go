package config

import (
	"path/filepath"
	"strings"
)

// DefaultName is the default name of the configuration file for Pilgo.
const DefaultName = "pilgo.yml"

const sep = string(filepath.Separator)

// Config is a configuration format for Pilgo.
type Config struct {
	BaseDir string             `yaml:"baseDir,omitempty"`
	Link    *string            `yaml:"link,omitempty"`
	Targets []string           `yaml:"targets,omitempty"`
	Options map[string]*Config `yaml:"options,omitempty"`
	UseHome *bool              `yaml:"useHome,omitempty"`
	Tags    []string           `yaml:"tags,omitempty"`

	opts internalOpts
}

// New creates a new configuration with given targets. Functional options can be
// passed to configure which targets should be included or excluded.
func New(targets []string, opts ...func(*Config)) Config {
	var c Config
	for _, o := range opts {
		o(&c)
	}
	eligible := make([]string, 0, len(targets))
	for _, tg := range targets {
		switch {
		case tg == "":
			fallthrough
		case tg == c.opts.skipName:
			fallthrough
		case !c.opts.addHidden && strings.HasPrefix(tg, "."):
			continue
		}
		if len(c.opts.inc) > 0 {
			_, included := c.opts.inc[tg]
			if !included {
				continue
			}
		}
		_, excluded := c.opts.exc[tg]
		if excluded {
			continue
		}
		eligible = append(eligible, tg)
	}
	c.Targets = eligible
	return c
}

// Set sets o to path. The path may be nested, but will be a no-op if the
// parent paths don't exist already. An empty path sets the root configuration.
func (c *Config) Set(path string, o Config) {
	if path == "" {
		*c = o
		return
	}
	targets := strings.Split(path, sep)
	for i, tg := range targets {
		if i != len(targets)-1 {
			if c.Options == nil {
				return
			}
			next, ok := c.Options[tg]
			if !ok {
				// Don't change what's not set.
				return
			}
			c = next
			continue
		}
		if o.isEmpty() {
			// Don't let empty maps (garbage) in the configuration file.
			delete(c.Options, tg)
			if len(c.Options) == 0 {
				c.Options = nil
			}
			return
		}
		if c.Options == nil {
			c.Options = make(map[string]*Config, 1)
		}
		c.Options[tg] = &o
	}
}

func (c *Config) isEmpty() bool {
	return c.BaseDir == "" &&
		c.Link == nil &&
		len(c.Targets) == 0 &&
		len(c.Options) == 0 &&
		c.UseHome == nil
}

type internalOpts struct {
	inc       map[string]struct{}
	exc       map[string]struct{}
	skipName  string
	addHidden bool
}

// MergeWith is an option to use an existing configuration when
// building a new one, thus inheriting all fields from the existing one.
func MergeWith(c Config) func(*Config) {
	return func(src *Config) {
		opts := src.opts
		*src = c
		src.opts = opts
	}
}

// Include is an option to set which files to include in the configuration.
func Include(set map[string]struct{}) func(*Config) {
	return func(c *Config) {
		c.opts.inc = set
	}
}

// IncludeHidden is an option that allows files prepended by a dot to be
// included in the configuration.
func IncludeHidden(c *Config) {
	c.opts.addHidden = true
}

// Exclude is an option to set which files to exclude in the configuration.
// If an inclusion list is set, this list will have no effect.
func Exclude(set map[string]struct{}) func(*Config) {
	return func(c *Config) {
		c.opts.exc = set
	}
}
