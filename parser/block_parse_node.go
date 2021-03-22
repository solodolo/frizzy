package parser

import "fmt"

// BlockParseNode contains a statement that does something but
// should not be rendered to the output
type BlockParseNode struct {
	ParseNode
}

func (node *BlockParseNode) String() string {
	return fmt.Sprintf("%T", *node)
}

// GetContent returns the TreeNode representing the
// content that should be output in the final template
func (receiver *BlockParseNode) GetContent() TreeNode {
	return receiver.children[1]
}

func (receiver *BlockParseNode) IsPrintable() bool {
	openBlock := receiver.children[0].(*StringParseNode)
	return openBlock.Value == "{{:"
}
