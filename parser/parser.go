package parser

import (
	"os"
	"sort"

	"gsr.dev/pilgrim/config"
)

// Parser is a configuration parser.
type Parser struct {
	cwd      string
	baseDir  string
	envsubst bool
}

// Parse parses a configuration file and returns its tree representation.
func (p *Parser) Parse(c config.Config, opts ...ParseOption) (*Tree, error) {
	for _, opt := range opts {
		if err := opt(p); err != nil {
			return nil, err
		}
	}
	if c.BaseDir != "" {
		p.baseDir = c.BaseDir
	}
	root := &Node{Children: p.parseChildren(p.baseDir, nil, nil, c)}
	return &Tree{root}, nil
}

func (p *Parser) parseChildren(baseDir string, parentTarget, parentLink []string, c config.Config) []*Node {
	var children []*Node
	tglen := len(c.Targets)
	if tglen > 0 {
		sort.Strings(c.Targets)
		children = make([]*Node, tglen)
		for i, tg := range c.Targets {
			tg = p.expandVar(tg)
			children[i] = p.parseTarget(
				p.expandVar(baseDir),
				append(parentTarget, tg),
				append(parentLink, tg),
				c.Options[tg],
			)
		}
	}
	return children
}

func (p *Parser) parseTarget(baseDir string, target, link []string, c config.Config) *Node {
	n := &Node{Target: File{p.cwd, target}}
	if c.BaseDir != "" {
		baseDir = c.BaseDir
	}
	if c.Link != nil {
		lnlen := len(link)
		linkname := *c.Link
		link[lnlen-1] = linkname
		if linkname == "" {
			s := link[:lnlen-1]
			// We need to create a new slice to avoid reusing the
			// same underlying array between children.
			link = append(make([]string, 0, len(s)), s...)
		}
	}
	n.Link = File{baseDir, link}
	n.Children = p.parseChildren(baseDir, target, link, c)
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

// BaseDir sets the default base directory when none is provided.
func BaseDir(dirname string) ParseOption {
	return func(p *Parser) error {
		p.baseDir = dirname
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
