package config

import (
	"path/filepath"
	"strings"
)

const (
	// DefaultName is the default name of the configuration file for Pilgo.
	DefaultName = "pilgo.yml"
	sep         = string(filepath.Separator)
)

// SetMode is a type for the mode used when setting a configuration.
type SetMode int

const (
	// ModeConfig is the mode for setting every field except targets.
	ModeConfig SetMode = iota
	// ModeScan is the mode for setting only targets.
	ModeScan
)

// Config is a configuration format for Pilgo.
type Config struct {
	BaseDir string             `yaml:"baseDir,omitempty"`
	Link    string             `yaml:"link,omitempty"`
	Targets []string           `yaml:"targets,omitempty"`
	Options map[string]*Config `yaml:"options,omitempty"`
	Flatten bool               `yaml:"flatten,omitempty"`
	UseHome *bool              `yaml:"useHome,omitempty"`
	Tags    []string           `yaml:"tags,omitempty"`
}

// Set sets o to path. The path may be nested, but will be a no-op if the
// parent paths don't exist already. An empty path sets the root configuration.
//
// All fields of the original configuration are overridden. If preserving fields
// is the intention, use Merge instead.
func (c *Config) Set(path string, new *Config, m SetMode) {
	if path == "" {
		*c = *c.resolveNew(new, m)
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
		if c.Options == nil {
			c.Options = make(map[string]*Config, 1)
		}
		if cc := c.Options[tg]; cc != nil {
			*cc = *cc.resolveNew(new, m)
		} else {
			c.Options[tg] = new
		}
		if c.Options[tg].isEmpty() {
			// Don't let empty maps (garbage) in the configuration file.
			delete(c.Options, tg)
			if len(c.Options) == 0 {
				c.Options = nil
			}
			return
		}
	}
}

func (c *Config) isEmpty() bool {
	return c.BaseDir == "" &&
		c.Link == "" &&
		len(c.Targets) == 0 &&
		len(c.Options) == 0 &&
		c.UseHome == nil &&
		!c.Flatten &&
		len(c.Tags) == 0
}

func (c *Config) resolveNew(new *Config, m SetMode) *Config {
	switch m {
	case ModeConfig:
		new.Targets = c.Targets
	case ModeScan:
		tgs := new.Targets
		*new = *c
		new.Targets = tgs
	default:
		// TODO(gbrlsnchs): add a free mode to overwrite everything.
		panic("unknown mode")
	}
	return new
}
