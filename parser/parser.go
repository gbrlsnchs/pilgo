package parser

import (
	"os"
	"sort"

	"github.com/gbrlsnchs/pilgo/config"
)

// Mode is the type of configuration.
// Each configuration has a distinct base directory.
type Mode int

const (
	// UserMode refers to the default user configuration directory.
	UserMode Mode = iota
	// HomeMode refers to the home directory.
	HomeMode
)

// Parser is a configuration parser.
type Parser struct {
	cwd      string
	baseDirs map[Mode]string
	envsubst bool
	tags     map[string]struct{}
}

// Parse parses a configuration file and returns its tree representation.
func (p *Parser) Parse(c *config.Config, opts ...ParseOption) (*Tree, error) {
	for _, opt := range opts {
		if err := opt(p); err != nil {
			return nil, err
		}
	}
	root := &Node{Children: p.parseChildren(c, nil, nil)}
	return &Tree{root}, nil
}

func (p *Parser) parseChildren(c *config.Config, ptargets, plinks []string) []*Node {
	var children []*Node
	tglen := len(c.Targets)
	if tglen > 0 {
		for i, tg := range c.Targets {
			c.Targets[i] = p.expandVar(tg)
		}
		sort.Strings(c.Targets)
		children = make([]*Node, 0, tglen)
		for _, tg := range c.Targets {
			cc := c.Options[tg]
			if cc == nil {
				cc = new(config.Config) // use default config
			}
			if len(cc.Tags) > 0 {
				if len(p.tags) == 0 {
					continue
				}
				shouldInclude := false
				for _, t := range cc.Tags {
					if _, ok := p.tags[t]; ok {
						shouldInclude = true
						break
					}
				}
				if !shouldInclude {
					continue
				}
			}
			if cc.UseHome == nil {
				cc.UseHome = c.UseHome
			}
			if cc.BaseDir == "" {
				cc.BaseDir = c.BaseDir
			}
			cc.BaseDir = p.expandVar(cc.BaseDir)
			children = append(children, p.parseTarget(cc,
				append(ptargets, tg),
				append(plinks, tg)))
		}
	}
	return children
}

func (p *Parser) parseTarget(c *config.Config, targets, links []string) *Node {
	n := &Node{Target: File{p.cwd, targets}}
	lnlen := len(links)
	if c.Link != "" {
		// Replace last element from links. This is a link rename.
		linkname := c.Link
		links[lnlen-1] = linkname
	}
	if c.Flatten {
		s := links[:lnlen-1]
		// We need to create a new slice to avoid reusing the
		// same underlying array between children.
		links = append(make([]string, 0, len(s)), s...)
	}
	if c.BaseDir == "" {
		mode := UserMode
		if c.UseHome != nil && *c.UseHome {
			mode = HomeMode
		}
		c.BaseDir = p.baseDirs[mode]
	}
	n.Link = File{c.BaseDir, links}
	n.Children = p.parseChildren(c, targets, links)
	return n
}

func (p *Parser) expandVar(s string) string {
	if p.envsubst {
		return os.ExpandEnv(s)
	}
	return s
}

// ParseOption is a funcional option that intend to modify a Parser.
type ParseOption func(*Parser) error

func BaseDirs(dirs map[Mode]string) ParseOption {
	return func(p *Parser) error {
		p.baseDirs = dirs
		return nil
	}
}

// Cwd sets a common cwd for targets.
func Cwd(dirname string) ParseOption {
	return func(p *Parser) error {
		p.cwd = dirname
		return nil
	}
}

// Envsubst enables replacing ${var} or $var with their environment values.
func Envsubst(p *Parser) error {
	p.envsubst = true
	return nil
}

// Tags filters targets by their tags.
func Tags(tags map[string]struct{}) ParseOption {
	return func(p *Parser) error {
		p.tags = tags
		return nil
	}
}
