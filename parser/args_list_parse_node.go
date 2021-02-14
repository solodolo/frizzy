package parser

import "fmt"

type ArgsListParseNode struct {
	ParseNode
	flattenedChildren []*ArgsListParseNode
}

func (receiver *ArgsListParseNode) String() string {
	return fmt.Sprintf("%T", *receiver)
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
