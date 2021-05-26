package parser

import (
	"fmt"
)

// ForLoopParseNode represents the parent node in a tree
// containing a parsed for loop
// The children will be the elements of this loop including
// the loop condition and body
type ForLoopParseNode struct {
	ParseNode
}

func (receiver *ForLoopParseNode) String() string {
	return fmt.Sprintf("%T", receiver)
}

func (node *ForLoopParseNode) IsTerminal() bool {
	return false
}

// GetLoopIdent returns the loop identifier
// Given {{for foo in bar}}, returns TreeNode{foo}
func (receiver *ForLoopParseNode) GetLoopIdent() TreeNode {
	// expects children to be {"for", "foo", "in", "bar", "}}", content}
	// where "foo" is the loop ident
	return receiver.children[1]
}

// GetLoopInput returns the loop input
// Given {{for foo in bar}}, returns TreeNode{bar}
func (receiver *ForLoopParseNode) GetLoopInput() TreeNode {
	// expects children to be {"for", "foo", "in", "bar", "}}", content}
	// where "bar" is the loop input
	return receiver.children[3]
}

// GetLoopBody returns the ContentParseNode body of this loop
// Given {{for foo in bar}} content, returns TreeNode{content}
func (receiver *ForLoopParseNode) GetLoopBody() TreeNode {
	return receiver.children[5]
}
