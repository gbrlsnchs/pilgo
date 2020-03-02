package pilgrim

import (
	"path/filepath"
	"strings"
)

// DefaultConfig is the default name of the configuration file for Pilgrim.
const DefaultConfig = "pilgrim.yml"

// Config is a configuration format for Pilgrim.
type Config struct {
	BaseDir string            `yaml:"baseDir,omitempty"`
	Link    *string           `yaml:"link,omitempty"`
	Targets []string          `yaml:"targets,omitempty"`
	Options map[string]Config `yaml:"options,omitempty"`
}

// Init returns a copy of c with only eligible files from targets, which are any
// files that don't start with a dot nor is named equal DefaultConfig's value.
func (c Config) Init(targets []string) Config {
	eligible := make([]string, 0, len(targets))
	for _, tg := range targets {
		if strings.HasPrefix(tg, ".") || tg == DefaultConfig {
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
