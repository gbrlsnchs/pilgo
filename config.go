package pilgrim

// Config is a configuration format for Pilgrim.
type Config struct {
	BaseDir string
	Link    *string
	Targets []string
	Options map[string]Config
}
