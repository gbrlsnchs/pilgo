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
}

// Set sets o to path. The path may be nested, but will be a no-op if the
// parent paths don't exist already. An empty path sets the root configuration.
//
// All fields of the original configuration are overridden. If preserving fields
// is the intention, use Merge instead.
func (c *Config) Set(path string, new *Config) {
	if path == "" {
		*c = *new
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
			*cc = *new
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
		c.Link == nil &&
		len(c.Targets) == 0 &&
		len(c.Options) == 0 &&
		c.UseHome == nil
}
