package parser

import "fmt"

// ElseIfConditionalParseNode represents the parent of one or more else_if
// statements represented by an ElseIfListParseNode
// This should always be the child of an IfStatementParseNode
// There should be a single ElseIfListParseNode child node
type ElseIfConditionalParseNode struct {
	ParseNode
}

func (receiver *ElseIfConditionalParseNode) String() string {
	return fmt.Sprintf("%T", *receiver)
}

func (receiver *ElseIfConditionalParseNode) getElseIfList() (*ElseIfListParseNode, bool) {
	if len(receiver.children) > 0 {
		return receiver.children[0].(*ElseIfListParseNode), true
	}
	return nil, false
}

// GetConditionals simply wraps ElseIfListParseNode GetConditionals
func (receiver *ElseIfConditionalParseNode) GetConditionals() []TreeNode {
	conditionals := []TreeNode{}
	if elseIf, ok := receiver.getElseIfList(); ok {
		conditionals = elseIf.GetConditionals()
	}
	return conditionals
}

// GetElseIfAt returns the nth child else_if where n = index
// If index is out of bounds or else_if isn't found
// the second return value will be false
func (receiver *ElseIfConditionalParseNode) GetElseIfAt(index int) (TreeNode, bool) {
	if elseIf, ok := receiver.getElseIfList(); ok {
		flat := elseIf.GetFlattenedBlockChildren()
		if len(flat) > index {
			return flat[index], ok
		}
	}

	return nil, false
}
