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
	// StatusConflict means a symlink already exists but points to a different target.
	StatusConflict
	// StatusError means the target doesn't exist.
	StatusError
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
	case StatusConflict:
		return "conflict"
	case StatusError:
		return "error"
	default:
		return "undefined"
	}
}
