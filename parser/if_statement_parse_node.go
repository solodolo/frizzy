package parser

import "fmt"

// IfStatementParseNode represents the parent node of a parsed
// if statement
// The children are the elemnts of the if statment including
// each if/else_if/else block
// |-- parser.IfStatementParseNode
// 	|-- parser.IdentParseNode: if
// 	|-- parser.NonTerminalParseNode: expression
// 	|
// 	|-- parser.BlockParseNode
// 	|-- parser.ContentParseNode
// 	|
// 	|-- parser.ElseIfListParseNode
// 	|
// 	|-- parser.IdentParseNode: end
type IfStatementParseNode struct {
	ParseNode
}

func (receiver *IfStatementParseNode) String() string {
	return fmt.Sprintf("%T", *receiver)
}

// GetIfConditional returns the child if statement conditional
func (receiver *IfStatementParseNode) GetIfConditional() TreeNode {
	return receiver.children[1]
}

// GetIfBody returns the body of this node's if statement
func (receiver *IfStatementParseNode) GetIfBody() TreeNode {
	return receiver.children[3]
}

func (receiver *IfStatementParseNode) getElseIfList() (*ElseIfListParseNode, bool) {
	elseIf, ok := receiver.children[4].(*ElseIfListParseNode)
	return elseIf, ok
}

// GetElseIfConditionals returns an array of conditional nodes from each
// of this node's else_if children
// Return value may be empty
func (receiver *IfStatementParseNode) GetElseIfConditionals() []TreeNode {
	conditionals := []TreeNode{}

	if elseIfList, ok := receiver.getElseIfList(); ok {
		for _, conditional := range elseIfList.GetConditionals() {
			conditionals = append(conditionals, conditional)
		}
	}
	return conditionals
}

// GetElseIfBody returns the body of the index'th else_if or false if not found
func (receiver *IfStatementParseNode) GetElseIfBody(index int) (TreeNode, bool) {
	if elseIfList, ok := receiver.getElseIfList(); ok {
		if elseIfAt, ok := elseIfList.GetElseIfAt(index); ok {
			return elseIfAt.GetBody(), true
		}
	}

	return nil, false
}

func (receiver *IfStatementParseNode) hasElse() bool {
	numChildren := len(receiver.children)
	elseIdent, ok := receiver.children[numChildren-3].(*IdentParseNode)

	return ok && elseIdent.Value == "else"
}

// GetElseBody returns the boyd of the else section or false if there isn't one
func (receiver *IfStatementParseNode) GetElseBody() (TreeNode, bool) {
	if receiver.hasElse() {
		numChildren := len(receiver.children)
		return receiver.children[numChildren-2], true
	}

	return nil, false
}
