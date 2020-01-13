package parser

import (
	"strings"
	"text/tabwriter"

	"gsr.dev/pilgrim/internal/treewriter"
)

// Tree is a simple tree data structure
// that holds parsed configuration data.
type Tree struct {
	Root *Node
}

func (tr *Tree) String() string {
	var (
		bd strings.Builder
		w  = tabwriter.NewWriter(&bd, 0, 0, 1, ' ', 0)
		tw = treewriter.NewWriter(w, (*printableNode)(tr.Root))
	)
	tw.Write(nil)
	w.Flush()
	return bd.String()
}

// Walk traverses the tree using depth-first search and runs fn for each node found.
func (tr *Tree) Walk(fn func(*Node) error) error {
	for _, n := range tr.Root.Children {
		if err := walk(n, fn); err != nil {
			return err
		}
	}
	return nil
}

func walk(n *Node, fn func(*Node) error) error {
	if err := fn(n); err != nil {
		return err
	}
	for _, c := range n.Children {
		if err := walk(c, fn); err != nil {
			return err
		}
	}
	return nil
}
