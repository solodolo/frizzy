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

func (receiver ForLoopParseNode) String() string {
	return fmt.Sprintf("%T", receiver)
}

// GetLoopIdent returns the loop identifier
// Given for(foo in "bar"), returns IdentParseNode{foo}
func (receiver ForLoopParseNode) GetLoopIdent() IdentParseNode {
	// expects children to be {"for", "(", "foo", ...}
	// where "foo" is the loop ident
	return receiver.children[2].(IdentParseNode)
}

// GetLoopInput returns the loop input
// Given for(foo in "bar") returns TreeNode{"bar"}
func (receiver ForLoopParseNode) GetLoopInput() TreeNode {
	// expects children to be {"for", "(", "foo", "in", "bar",...}
	// where "bar" is the loop input
	return receiver.children[4]
}

// GetLoopBody returns the children of receiver that are part of the body
// The body is any child between the 'for(foo in "bar")' and the 'end'
// i.e.
// for(foo in 'bar')
//		some body business
//		more body business
// end
func (receiver ForLoopParseNode) GetLoopBody() []TreeNode {
	children := receiver.children
	body := make([]TreeNode, len(children)-7)
	copy(body, children[6:len(receiver.children)-1])

	return body
}
