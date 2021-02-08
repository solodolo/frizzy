package parser

import (
	"fmt"
)

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
