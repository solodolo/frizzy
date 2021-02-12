package parser

import "fmt"

type ElseIfListParseNode struct {
	ParseNode
}

func (receiver *ElseIfListParseNode) String() string {
	return fmt.Sprintf("%T", *receiver)
}

// GetConditionals returns the true/false statements inside of this
// and each nested else_if block
func (receiver *ElseIfListParseNode) GetConditionals() []TreeNode {
	numBlockChildren := receiver.getNumBlockChildren()
	conditionals := make([]TreeNode, numBlockChildren+1)

	current := receiver
	for i := 0; i < numBlockChildren; i++ {
		conditionals[i] = current.getConditional()
		current = current.children[0].(*ElseIfListParseNode)
	}

	return conditionals
}

func (receiver *ElseIfListParseNode) getNumBlockChildren() int {
	count := 0
	current := receiver

	for current.hasBlockChildren() {
		count++
		current = current.children[0].(*ElseIfListParseNode)
	}

	return count
}

func (receiver *ElseIfListParseNode) hasBlockChildren() bool {
	if len(receiver.children) == 0 {
		return false
	}
	_, ok := receiver.children[0].(*ElseIfListParseNode)
	return ok
}

func (receiver *ElseIfListParseNode) getConditional() TreeNode {
	// if there are multiple else_ifs then the first child will be
	// a type of ElseIfListParseNode
	// Otherwise it will be else_if(statement)
	var offset int
	if receiver.hasBlockChildren() {
		offset = 1
	}

	return receiver.children[2+offset]
}
