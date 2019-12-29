package parser

import (
	"sort"

	"gsr.dev/pilgrim"
)

// Parser is a configuration parser.
type Parser struct{}

// Parse parses a configuration file and returns its tree representation.
func (p *Parser) Parse(c pilgrim.Config) (*Tree, error) {
	root := &Node{Children: parseChildren(c.BaseDir, nil, c)}
	return &Tree{root}, nil
}

func parseTarget(baseDir string, target []string, c pilgrim.Config) *Node {
	n := &Node{Target: File{"", target}}
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
	n.Children = parseChildren(baseDir, target, c)
	return n
}

func parseChildren(baseDir string, parentTarget []string, c pilgrim.Config) []*Node {
	var children []*Node
	tglen := len(c.Targets)
	if tglen > 0 {
		sort.Strings(c.Targets)
		children = make([]*Node, tglen)
		for i, tg := range c.Targets {
			children[i] = parseTarget(baseDir, append(parentTarget, tg), c.Options[tg])
		}
	}
	return children
}
