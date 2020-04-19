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
}

// Parse parses a configuration file and returns its tree representation.
func (p *Parser) Parse(c config.Config, opts ...ParseOption) (*Tree, error) {
	for _, opt := range opts {
		if err := opt(p); err != nil {
			return nil, err
		}
	}
	root := &Node{Children: p.parseChildren(c, nil, nil)}
	return &Tree{root}, nil
}

func (p *Parser) parseChildren(c config.Config, ptargets, plinks []string) []*Node {
	var children []*Node
	tglen := len(c.Targets)
	if tglen > 0 {
		sort.Strings(c.Targets)
		children = make([]*Node, tglen)
		for i, tg := range c.Targets {
			tg = p.expandVar(tg)
			cc := c.Options[tg]
			if cc.UseHome == nil {
				cc.UseHome = c.UseHome
			}
			if cc.BaseDir == "" {
				cc.BaseDir = c.BaseDir
			}
			cc.BaseDir = p.expandVar(cc.BaseDir)
			children[i] = p.parseTarget(cc,
				append(ptargets, tg),
				append(plinks, tg))
		}
	}
	return children
}

func (p *Parser) parseTarget(c config.Config, targets, links []string) *Node {
	n := &Node{Target: File{p.cwd, targets}}
	if c.Link != nil {
		// Replace last element from links. This is a link rename.
		lnlen := len(links)
		linkname := *c.Link
		links[lnlen-1] = linkname
		if linkname == "" {
			s := links[:lnlen-1]
			// We need to create a new slice to avoid reusing the
			// same underlying array between children.
			links = append(make([]string, 0, len(s)), s...)
		}
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
