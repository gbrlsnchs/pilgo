package parser

import "strings"

// Status is a node's status.
type Status uint8

const (
	// StatusReady means the node is ready to be symlinked.
	StatusReady Status = iota + 1
	// StatusSkip means the node has children and thus might be skipped.
	StatusSkip
	// StatusDone means the symlink already exists and is pointing exactly
	// to the specified node.
	StatusDone
)

func (s Status) String() string { return strings.ToUpper(s.str()) }

func (s Status) str() string {
	switch s {
	case StatusReady:
		return "ready"
	case StatusSkip:
		return "skip"
	case StatusDone:
		return "done"
	default:
		return "undefined"
	}
}
