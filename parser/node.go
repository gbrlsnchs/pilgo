package parser

// Node is a tree node that holds nested file metadata.
type Node struct {
	Target   File
	Link     File
	Children []*Node
	Status   Status
}
