package parser

import "fmt"

// ElseParseNode represents a parsed else statement
// It should always be the child of an IfStatementParseNode
// The children represent the elements of the else statement
// specifically the body
type ElseParseNode struct {
	ParseNode
}

func (receiver ElseParseNode) String() string {
	return fmt.Sprintf("%T", receiver)
}

// GetBody returns the body of this else
// Body node should always exist but might not have any children
// if the else section is empty
func (receiver *ElseParseNode) GetBody() TreeNode {
	return receiver.children[1]
}
