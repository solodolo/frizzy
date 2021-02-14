package parser

import "fmt"

// IfConditionalParseNode represents a single if block including
// the conditional and body
type IfConditionalParseNode struct {
	ParseNode
}

func (receiver *IfConditionalParseNode) String() string {
	return fmt.Sprintf("%T", *receiver)
}

// GetConditional returns the true/false statement for this if
// children are ["if", "(", conditional, ...]
func (receiver *IfConditionalParseNode) GetConditional() TreeNode {
	return receiver.children[2]
}

// GetBody returns the nodes that make up this if conditions body
func (receiver *IfConditionalParseNode) GetBody() TreeNode {
	return receiver.children[4]
}
