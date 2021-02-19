package parser

import "fmt"

// StatementListParseNode represents a list of ';' separated
// statements
// Each statement after the first will be nested under the previous
// like:
//
// |-- parser.StatementListParseNode:
//    |-- parser.StatementListParseNode
//       |-- parser.StatementListParseNode
type StatementListParseNode struct {
	ParseNode
	flattenedChildren []*StatementListParseNode
}

func (receiver *StatementListParseNode) String() string {
	return fmt.Sprintf("%T", *receiver)
}

// GetStatements returns an array of nodes representing the statements
// from top to bottom
func (receiver *StatementListParseNode) GetStatements() []TreeNode {
	flat := receiver.GetFlattenedChildren()
	children := make([]TreeNode, len(flat))

	for i, child := range flat {
		children[i] = child.getStatement()
	}

	return children
}

// GetFlattenedChildren returns the StatementListParseNode children of
// this node in a flattened array
func (receiver *StatementListParseNode) GetFlattenedChildren() []*StatementListParseNode {
	if receiver.flattenedChildren == nil {
		receiver.cacheFlattenedChildren()
	}

	return receiver.flattenedChildren
}

// cacheFlattenedChildren stores this receiver and any StatementListParseNode
// children in a flat array
func (receiver *StatementListParseNode) cacheFlattenedChildren() {
	receiver.flattenedChildren = []*StatementListParseNode{receiver}

	current := receiver
	for current.hasNestedChildren() {
		next := current.children[0].(*StatementListParseNode)
		// bottom child is first else_if so prepend
		receiver.flattenedChildren = append([]*StatementListParseNode{next}, receiver.flattenedChildren...)
		current = next
	}
}

func (receiver *StatementListParseNode) hasNestedChildren() bool {
	if len(receiver.children) == 0 {
		return false
	}
	_, ok := receiver.children[0].(*StatementListParseNode)
	return ok
}

func (receiver *StatementListParseNode) getStatement() TreeNode {
	// the first child is either another StatementListParseNode or
	// a statement node
	if receiver.hasNestedChildren() {
		return receiver.children[1]
	}

	return receiver.children[0]
}
