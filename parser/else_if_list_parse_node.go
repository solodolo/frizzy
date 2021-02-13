package parser

import "fmt"

// ElseIfListParseNode represents one or more else_ifs with
// their conditionals and bodies
type ElseIfListParseNode struct {
	ParseNode
	flattenedBlockChildren []*ElseIfListParseNode
}

func (receiver *ElseIfListParseNode) String() string {
	return fmt.Sprintf("%T", *receiver)
}

// GetConditionals returns the true/false statements inside of this
// and each nested else_if block
func (receiver *ElseIfListParseNode) GetConditionals() []TreeNode {
	blocks := receiver.GetFlattenedBlockChildren()
	conditionals := make([]TreeNode, len(blocks))

	for i, block := range blocks {
		conditionals[i] = block.GetConditional()
	}

	return conditionals
}

// GetFlattenedBlockChildren returns the ElseIfListParseNode children of
// this node in a flattened array
func (receiver *ElseIfListParseNode) GetFlattenedBlockChildren() []*ElseIfListParseNode {
	if receiver.flattenedBlockChildren == nil {
		receiver.cacheFlattenedBlockChildren()
	}

	return receiver.flattenedBlockChildren
}

// cacheFlattenedBlockChildren stores this receiver and any ElseIfListParseNode
// children in a flat array
func (receiver *ElseIfListParseNode) cacheFlattenedBlockChildren() {
	receiver.flattenedBlockChildren = []*ElseIfListParseNode{receiver}

	current := receiver
	for current.hasBlockChildren() {
		next := current.children[0].(*ElseIfListParseNode)
		receiver.flattenedBlockChildren = append(receiver.flattenedBlockChildren, next)
		current = next
	}
}

func (receiver *ElseIfListParseNode) hasBlockChildren() bool {
	if len(receiver.children) == 0 {
		return false
	}
	_, ok := receiver.children[0].(*ElseIfListParseNode)
	return ok
}

// GetConditional returns the true/false statement for this else_if
func (receiver *ElseIfListParseNode) GetConditional() TreeNode {
	// if there are multiple else_ifs then the first child will be
	// a type of ElseIfListParseNode
	// Otherwise it will be else_if(statement)
	var offset int
	if receiver.hasBlockChildren() {
		offset = 1
	}

	return receiver.children[2+offset]
}
