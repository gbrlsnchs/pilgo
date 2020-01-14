package parser

import (
	"sort"

	"gsr.dev/pilgrim"
)

// Parser is a configuration parser.
type Parser struct {
	cwd string
}

// Parse parses a configuration file and returns its tree representation.
func (p *Parser) Parse(c pilgrim.Config, opts ...ParseOption) (*Tree, error) {
	for _, opt := range opts {
		if err := opt(p); err != nil {
			return nil, err
		}
	}
	root := &Node{Children: p.parseChildren(c.BaseDir, nil, c)}
	return &Tree{root}, nil
}

func (p *Parser) parseTarget(baseDir string, target []string, c pilgrim.Config) *Node {
	n := &Node{Target: File{p.cwd, target}}
	if c.BaseDir != "" {
		baseDir = c.BaseDir
	}
	tglen := len(target)
	link := make([]string, tglen)
	copy(link, target)
	if c.Link != nil {
		link[tglen-1] = *c.Link
	}
	n.Link = File{baseDir, link}
	n.Children = p.parseChildren(baseDir, target, c)
	return n
}

func (p *Parser) parseChildren(baseDir string, parentTarget []string, c pilgrim.Config) []*Node {
	var children []*Node
	tglen := len(c.Targets)
	if tglen > 0 {
		sort.Strings(c.Targets)
		children = make([]*Node, tglen)
		for i, tg := range c.Targets {
			children[i] = p.parseTarget(baseDir, append(parentTarget, tg), c.Options[tg])
		}
	}
	return children
}

// ParseOption is a funcional option that intend to modify a Parser.
type ParseOption func(*Parser) error

// Cwd sets a common cwd for targets.
func Cwd(dirname string) ParseOption {
	return func(p *Parser) error {
		p.cwd = dirname
		return nil
	}
}
