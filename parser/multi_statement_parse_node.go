package parser

import "fmt"

// MultiStatementParseNode represents the parent node
// of zero or more nested statements
type MultiStatementParseNode struct {
	ParseNode
}

func (receiver *MultiStatementParseNode) String() string {
	return fmt.Sprintf("%T", *receiver)
}

// GetStatements returns an array of each child statment of this
// multi statement node
func (receiver *MultiStatementParseNode) GetStatements() []TreeNode {
	if len(receiver.children) > 0 {
		statementList := receiver.children[0].(*StatementListParseNode)
		return statementList.GetStatements()
	}

	return []TreeNode{}
}
