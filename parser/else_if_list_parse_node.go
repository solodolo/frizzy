package parser

import "fmt"

// ElseIfListParseNode represents one or more else_ifs with
// their conditionals and bodies
// Will be a tree like
//	|-- ElseIfListParseNode
//		|-- ElseIfListParseNode
//			|-- IdentParseNode: else_if
//			|-- Expression
//			|-- BlockParseNode
//			|-- ContentParseNode
type ElseIfListParseNode struct {
	ParseNode
}

func (receiver *ElseIfListParseNode) String() string {
	return fmt.Sprintf("%T", *receiver)
}

func (node *ElseIfListParseNode) IsTerminal() bool {
	return false
}

// GetConditionals returns the true/false statements inside of this
// and each nested else_if block
func (receiver *ElseIfListParseNode) GetConditionals() []TreeNode {
	numConditionals := len(receiver.children) / 4
	conditionals := make([]TreeNode, 0, numConditionals)

	for i := 0; i < numConditionals; i++ {
		offset := i*4 + 1
		conditionals = append(conditionals, receiver.children[offset])
	}

	return conditionals
}

func (receiver *ElseIfListParseNode) GetBodies() []TreeNode {
	numBodies := len(receiver.children) / 4
	bodies := make([]TreeNode, 0, numBodies)

	for i := 0; i < numBodies; i++ {
		offset := i*4 + 3
		bodies = append(bodies, receiver.children[offset])
	}

	return bodies
}
