package treewriter_test

import "gsr.dev/pilgrim/parser/internal/treewriter"

type testNode struct {
	text  string
	nodes []*testNode
}

func (tr *testNode) At(i int) treewriter.Node { return tr.nodes[i] }
func (tr *testNode) Len() int                 { return len(tr.nodes) }
func (tr *testNode) String() string           { return tr.text }
