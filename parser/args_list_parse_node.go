package parser

import (
	"fmt"
	"math"
)

// ArgsListParseNode represents a comma separated list of function arguments
type ArgsListParseNode struct {
	ParseNode
	flattenedChildren []*ArgsListParseNode
}

func (receiver *ArgsListParseNode) String() string {
	return fmt.Sprintf("%T", *receiver)
}

func (node *ArgsListParseNode) IsTerminal() bool {
	return false
}

// GetArguments returns an array of nodes representing the arguments
// from left to right
func (receiver *ArgsListParseNode) GetArguments() []TreeNode {
	numParts := float64(len(receiver.children)) / float64(2)
	numArgs := int(math.Ceil(numParts))
	args := make([]TreeNode, 0, numArgs)

	for i := 0; i < len(receiver.children); i += 2 {
		args = append(args, receiver.children[i])
	}

	return args
}
