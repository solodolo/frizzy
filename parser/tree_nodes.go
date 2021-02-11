package parser

import (
	"fmt"
)

type TreeNode interface {
	GetChildren() []TreeNode
	PrintTree()
	fmt.Stringer
}

type ParseNode struct {
	children []TreeNode
}

func (node *ParseNode) String() string {
	return ""
}

func (node *ParseNode) GetChildren() []TreeNode {
	return node.children
}

func (node *ParseNode) PrintTree() {
	treeStrs := genTreeStr(node, 0, true)

	for _, str := range treeStrs {
		fmt.Println(str)
	}
}

func genTreeStr(node TreeNode, level int, lastChild bool) []string {
	children := node.GetChildren()
	numChildren := len(children)
	treeStrs := []string{fmt.Sprintf("%s", node)}

	if numChildren == 0 {
		return treeStrs
	}

	for i := 0; i < numChildren; i++ {
		child := children[i]
		childStrs := genTreeStr(child, level+1, i >= numChildren-1)

		childStrs[0] = "|-- " + childStrs[0]

		connect := len(childStrs) > 1 && i < numChildren-1

		for j := 1; j < len(childStrs); j++ {
			if connect {
				childStrs[j] = "|  " + childStrs[j]
			} else {
				childStrs[j] = "   " + childStrs[j]
			}
		}

		treeStrs = append(treeStrs, childStrs...)
	}

	return treeStrs
}

type NonTerminalParseNode struct {
	Value string
	ParseNode
}

func (node NonTerminalParseNode) String() string {
	return fmt.Sprintf("%T: %s", node, node.Value)
}

func (node NonTerminalParseNode) IsAssignment() bool {
	return node.Value == "K" && len(node.children) > 1
}

func (node NonTerminalParseNode) IsAddition() bool {
	return node.Value == "N" && len(node.children) > 1
}

func (node NonTerminalParseNode) IsMultiplication() bool {
	return node.Value == "O" && len(node.children) > 1
}

func (node NonTerminalParseNode) IsLogic() bool {
	return node.Value == "M" && len(node.children) > 1
}

func (node NonTerminalParseNode) IsRelation() bool {
	return node.Value == "U" && len(node.children) > 1
}

func (node NonTerminalParseNode) IsUnary() bool {
	return node.Value == "L" && len(node.children) > 1
}

type StringParseNode struct {
	Value string
	ParseNode
}

func (node StringParseNode) String() string {
	return fmt.Sprintf("%T: %q", node, node.Value)
}

type NumParseNode struct {
	Value int
	ParseNode
}

func (node NumParseNode) String() string {
	return fmt.Sprintf("%T: %d", node, node.Value)
}

type BoolParseNode struct {
	Value bool
	ParseNode
}

func (node BoolParseNode) String() string {
	return fmt.Sprintf("%T: %t", node, node.Value)
}

type IdentParseNode struct {
	Value string
	ParseNode
}

func (node IdentParseNode) String() string {
	return fmt.Sprintf("%T: %s", node, node.Value)
}

type VarParseNode struct {
	Value string
	ParseNode
}

func (node VarParseNode) String() string {
	return fmt.Sprintf("%T: %s", node, node.Value)
}
