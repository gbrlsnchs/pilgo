package parser

import "strings"

// Status is a node's status.
type Status uint8

const (
	// StatusReady means the node is ready to be symlinked.
	StatusReady Status = iota + 1
)

func (s Status) String() string { return strings.ToUpper(s.str()) }

func (s Status) str() string {
	switch s {
	case StatusReady:
		return "ready"
	default:
		return "undefined"
	}
}