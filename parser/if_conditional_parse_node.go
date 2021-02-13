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
