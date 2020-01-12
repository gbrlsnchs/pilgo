package treewriter

// Node is a tree node.
type Node interface {
	At(int) Node
	Len() int
}
