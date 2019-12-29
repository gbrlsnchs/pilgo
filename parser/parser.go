package parser

import (
	"sort"

	"gsr.dev/pilgrim"
)

// Parser is a configuration parser.
type Parser struct{}

// Parse parses a configuration file and returns its tree representation.
func (p *Parser) Parse(c pilgrim.Config) (*Tree, error) {
	root := &Node{Children: parseChildren(c.BaseDir, c)}
	return &Tree{root}, nil
}

func parseTarget(baseDir, target string, c pilgrim.Config) *Node {
	n := &Node{Target: File{"", target}}
	if c.BaseDir != "" {
		baseDir = c.BaseDir
	}
	link := target
	if c.Link != nil {
		link = *c.Link
	}
	n.Link = File{baseDir, link}
	return n
}

func parseChildren(baseDir string, c pilgrim.Config) []*Node {
	var children []*Node
	tglen := len(c.Targets)
	if tglen > 0 {
		sort.Strings(c.Targets)
		children = make([]*Node, tglen)
		for i, tg := range c.Targets {
			children[i] = parseTarget(baseDir, tg, c.Options[tg])
		}
	}
	return children
}
