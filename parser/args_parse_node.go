package parser

import "fmt"

type ArgsParseNode struct {
	ParseNode
}

func (receiver *ArgsParseNode) String() string {
	return fmt.Sprintf("%T", *receiver)
}

// GetArguments returns an array of nodes representing the arguments
// from left to right
// Return will be empty if no arguments are present
func (receiver *ArgsParseNode) GetArguments() []TreeNode {
	if len(receiver.children) > 0 {
		argsList := receiver.children[0].(*ArgsListParseNode)
		return argsList.GetArguments()
	}

	return []TreeNode{}
}
