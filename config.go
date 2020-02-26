package pilgrim

// Config is a configuration format for Pilgrim.
type Config struct {
	BaseDir string            `yaml:"baseDir,omitempty"`
	Link    *string           `yaml:"link,omitempty"`
	Targets []string          `yaml:"targets,omitempty"`
	Options map[string]Config `yaml:"options,omitempty"`
}
