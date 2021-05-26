package parser

import "fmt"

// ArgsListParseNode represents a comma separated list of function arguments
type ArgsListParseNode struct {
	ParseNode
	flattenedChildren []*ArgsListParseNode
}

func (receiver *ArgsListParseNode) String() string {
	return fmt.Sprintf("%T", *receiver)
}

func (node *ArgsListParseNode) IsTerminal() bool {
	return false
}

// GetArguments returns an array of nodes representing the arguments
// from left to right
func (receiver *ArgsListParseNode) GetArguments() []TreeNode {
	if nestedArgsList, ok := receiver.children[0].(*ArgsListParseNode); ok {
		return append(nestedArgsList.GetArguments(), receiver.children[2])
	}
	return []TreeNode{receiver.children[0], receiver.children[2]}
}

// GetFlattenedChildren returns the ArgsListParseNode children of
// this node in a flattened array
func (receiver *ArgsListParseNode) GetFlattenedChildren() []*ArgsListParseNode {
	if receiver.flattenedChildren == nil {
		receiver.cacheFlattenedChildren()
	}

	return receiver.flattenedChildren
}

// cacheFlattenedChildren stores this receiver and any ArgsListParseNode
// children in a flat array
func (receiver *ArgsListParseNode) cacheFlattenedChildren() {
	receiver.flattenedChildren = []*ArgsListParseNode{receiver}

	current := receiver
	for current.hasNestedChildren() {
		next := current.children[0].(*ArgsListParseNode)
		// bottom child is first else_if so prepend
		receiver.flattenedChildren = append([]*ArgsListParseNode{next}, receiver.flattenedChildren...)
		current = next
	}
}

func (receiver *ArgsListParseNode) hasNestedChildren() bool {
	if len(receiver.children) == 0 {
		return false
	}
	_, ok := receiver.children[0].(*ArgsListParseNode)
	return ok
}

func (receiver *ArgsListParseNode) getArg() TreeNode {
	if receiver.hasNestedChildren() {
		return receiver.children[2]
	}
	return receiver.children[0]
}
