package pilgrim

import "strings"

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
