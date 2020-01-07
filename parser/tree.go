package parser

// Tree is a simple tree data structure
// that holds parsed configuration data.
type Tree struct {
	Root *Node
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
