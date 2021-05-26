package parser

import (
	"fmt"
)

// FuncCallParseNode represents a function call in our grammar
type FuncCallParseNode struct {
	ParseNode
}

func (receiver *FuncCallParseNode) String() string {
	return fmt.Sprintf("%T", *receiver)
}

func (node *FuncCallParseNode) IsTerminal() bool {
	return false
}

// GetFuncName returns the name of the function that
// this node represents
func (receiver *FuncCallParseNode) GetFuncName() string {
	nameNode := receiver.children[0].(*IdentParseNode)
	return nameNode.Value
}

// GetArgs returns a slice of nodes representing the arguments
// to this function
func (receiver *FuncCallParseNode) GetArgs() []TreeNode {
	if argsList, ok := receiver.children[2].(*ArgsListParseNode); ok {
		return argsList.GetArguments()
	}
	return []TreeNode{receiver.children[2]}
}
