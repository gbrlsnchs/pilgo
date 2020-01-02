package parser

// Status is a node's status.
type Status uint8

const (
	// StatusReady means the node is ready to be symlinked.
	StatusReady Status = iota + 1
)
