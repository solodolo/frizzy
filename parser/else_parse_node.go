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
