package parser

import "fmt"

// ElseIfStatementParseNode represents the parent of one or more else_if
// statements represented by an ElseIfListParseNode
// This should always be the child of an IfStatementParseNode
// There should be a single ElseIfListParseNode child node
type ElseIfStatementParseNode struct {
	ParseNode
}

func (receiver *ElseIfStatementParseNode) String() string {
	return fmt.Sprintf("%T", receiver)
}

// GetConditionals simply wraps ElseIfListParseNode GetConditionals
func (receiver *ElseIfStatementParseNode) GetConditionals() []TreeNode {
	child := receiver.children[0].(*ElseIfListParseNode)
	return child.GetConditionals()
}
