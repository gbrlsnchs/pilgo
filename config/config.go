package config

import (
	"path/filepath"
	"strings"
)

// DefaultName is the default name of the configuration file for Pilgo.
const DefaultName = "pilgo.yml"

// Config is a configuration format for Pilgo.
type Config struct {
	BaseDir string            `yaml:"baseDir,omitempty"`
	Link    *string           `yaml:"link,omitempty"`
	Targets []string          `yaml:"targets,omitempty"`
	Options map[string]Config `yaml:"options,omitempty"`
	UseHome *bool             `yaml:"useHome,omitempty"`

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

// Set sets o to the resolved option in path.
func (c *Config) Set(path string, o Config) {
	var last string
	path, last = filepath.Split(path)
	if path == "" {
		if last == "" {
			*c = o
			return
		}
		if c.Options == nil {
			c.Options = make(map[string]Config, 1)
		}
		c.Options[last] = o
		return
	}
	path = path[:len(path)-1] // trim suffix separator
	var (
		keys = strings.Split(path, string(filepath.Separator))
		oo   Config
		ok   bool
	)
	for _, k := range keys {
		if oo, ok = c.Options[k]; !ok {
			return
		}
		if oo.Options == nil {
			oo = Config{
				BaseDir: oo.BaseDir,
				Link:    oo.Link,
				Targets: oo.Targets,
				Options: make(map[string]Config, 1),
			}
			c.Options[k] = oo
		}
	}
	oo.Options[last] = o
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
