package parser

import "fmt"

// IfStatementParseNode represents the parent node of a parsed
// if statement
// The children are the elemnts of the if statment including
// each if/else_if/else block
//
// tree structure
//|-- parser.IfStatementParseNode
//		|-- parser.IdentParseNode: if
//		|-- parser.StringParseNode: "("
//		|-- parser.NonTerminalParseNode: K
//		|-- parser.StringParseNode: ")"
//		|-- parser.NonTerminalParseNode: V
//		|-- parser.ElseIfStatementParseNode
//		|-- ...
//		|-- parser.ElseParseNode
//		|-- parser.IdentParseNode: end
type IfStatementParseNode struct {
	ParseNode
}

func (receiver IfStatementParseNode) String() string {
	return fmt.Sprintf("%T", receiver)
}

// GetSelectedBody examines each block of this if statement in order
// and returns the TreeNodes in the body of the first true block
func (receiver IfStatementParseNode) GetSelectedBody() []TreeNode {
	return nil
}
