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

// GetFuncName returns the name of the function that
// this node represents
func (receiver *FuncCallParseNode) GetFuncName() string {
	nameNode := receiver.children[0].(*IdentParseNode)
	return nameNode.Value
}

func (receiver *FuncCallParseNode) getArgNode() *ArgsParseNode {
	return receiver.children[2].(*ArgsParseNode)
}

// GetArgs returns a slice of nodes representing the arguments
// to this function
func (receiver *FuncCallParseNode) GetArgs() []TreeNode {
	argNode := receiver.getArgNode()
	return argNode.GetArguments()
}
