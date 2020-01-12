package parser

import (
	"fmt"
	"strings"

	"gsr.dev/pilgrim/internal/treewriter"
)

const dullStatus = StatusSkip | StatusExpand

// Node is a tree node that holds nested file metadata.
type Node struct {
	Target   File
	Link     File
	Children []*Node
	Status   Status
}

// At returns a child node at index i.
func (n *Node) At(i int) treewriter.Node { return n.Children[i] }

// Len returns the number of children of n.
func (n *Node) Len() int { return len(n.Children) }

func (n *Node) String() string {
	if len(n.Target.Path) == 0 {
		return ""
	}
	var (
		bd     strings.Builder
		symbol = "<-"
	)
	if n.Status&dullStatus != 0 {
		symbol = ""
	}
	fmt.Fprintf(&bd, "%s\t%s", n.Target.base(), symbol)
	if n.Status&dullStatus == 0 {
		if fullPath := n.Link.FullPath(); fullPath != "" {
			fmt.Fprintf(&bd, " %s", fullPath)
		}
	}
	if n.Status > 0 {
		fmt.Fprintf(&bd, "\t(%s)", n.Status)
	}
	return bd.String()
}
