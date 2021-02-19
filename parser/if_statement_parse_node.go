package parser

import "fmt"

// IfStatementParseNode represents the parent node of a parsed
// if statement
// The children are the elemnts of the if statment including
// each if/else_if/else block
// |  |-- parser.IfStatementParseNode
// |     |-- parser.IfConditionalParseNode
// |     |-- parser.ElseIfConditionalParseNode
// |     |  |-- parser.ElseIfListParseNode
// |     |     |-- parser.ElseIfListParseNode...
// |     |-- parser.ElseParseNode
// |     |-- parser.IdentParseNode: end
type IfStatementParseNode struct {
	ParseNode
}

func (receiver *IfStatementParseNode) String() string {
	return fmt.Sprintf("%T", *receiver)
}

func (receiver *IfStatementParseNode) getIfConditional() *IfConditionalParseNode {
	return receiver.children[0].(*IfConditionalParseNode)
}

// GetIfConditional returns the child if statement conditional
func (receiver *IfStatementParseNode) GetIfConditional() TreeNode {
	ifConditional := receiver.getIfConditional()
	return ifConditional.GetConditional()
}

// GetIfBody returns the body of this node's if statement
func (receiver *IfStatementParseNode) GetIfBody() TreeNode {
	ifConditional := receiver.getIfConditional()
	return ifConditional.GetBody()
}

func (receiver *IfStatementParseNode) getElseIfConditional() (*ElseIfConditionalParseNode, bool) {
	elseIf, ok := receiver.children[1].(*ElseIfConditionalParseNode)
	return elseIf, ok
}

// GetElseIfConditionals returns an array of conditional nodes from each
// of this node's else_if children
// Return value may be empty
func (receiver *IfStatementParseNode) GetElseIfConditionals() []TreeNode {
	conditionals := []TreeNode{}

	if elseIfConditional, ok := receiver.getElseIfConditional(); ok {
		for _, conditional := range elseIfConditional.GetConditionals() {
			conditionals = append(conditionals, conditional)
		}
	}
	return conditionals
}

// GetElseIfBody returns the body of the index'th else_if or false if not found
func (receiver *IfStatementParseNode) GetElseIfBody(index int) (TreeNode, bool) {
	if elseIfConditional, ok := receiver.getElseIfConditional(); ok {
		if elseIfAt, ok := elseIfConditional.GetElseIfAt(index); ok {
			return elseIfAt.(*ElseIfListParseNode).GetBody(), true
		}
	}

	return nil, false
}

func (receiver *IfStatementParseNode) getElseNode() (*ElseParseNode, bool) {
	// else should be second to last child (last child is end)
	elseNode, ok := receiver.children[len(receiver.children)-2].(*ElseParseNode)
	return elseNode, ok
}

// GetElseBody returns the boyd of the else section or false if there isn't one
func (receiver *IfStatementParseNode) GetElseBody() (TreeNode, bool) {
	if elseNode, ok := receiver.getElseNode(); ok {
		return elseNode.GetBody()
	}

	return nil, false
}
