package parser

import "fmt"

// ElseIfStatementParseNode represents a parsed else_if statement
// It should always be the child of an IfStatementParseNode
// The children represent the elements of the else if statement
// including the condition and body
type ElseIfStatementParseNode struct {
	ParseNode
}

func (receiver ElseIfStatementParseNode) String() string {
	return fmt.Sprintf("%T", receiver)
}
