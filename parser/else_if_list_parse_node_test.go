package parser

import (
	"testing"
)

func getNodeWithBlockChildren(count int) ElseIfListParseNode {
	node := ElseIfListParseNode{}
	current := &node

	for i := 0; i < count; i++ {
		next := &ElseIfListParseNode{}
		children := []TreeNode{next}
		current.SetChildren(children)
		current = next
	}

	return node
}

func getNodeWithoutBlockChildren(count int) ElseIfListParseNode {
	node := ElseIfListParseNode{}
	var current TreeNode = &node

	for i := 0; i < count; i++ {
		next := &IdentParseNode{}
		children := []TreeNode{next}
		current.SetChildren(children)
		current = next
	}

	return node
}

func TestHasBlockChildrenReturnsTrueWithBlockChildren(t *testing.T) {
	node := getNodeWithBlockChildren(3)
	if !node.hasBlockChildren() {
		t.Errorf("expected node with block children to return true")
	}
}

func TestHasBlockChildrenReturnsFalseWithNoBlockChildren(t *testing.T) {
	node := getNodeWithoutBlockChildren(3)
	if node.hasBlockChildren() {
		t.Errorf("expected node with no block children to return false")
	}
}

func TestGetNumBlockChildrenWithBlockChildrenReturnsCorrectNumber(t *testing.T) {
	expected := 3
	node := getNodeWithBlockChildren(expected)
	got := node.getNumBlockChildren()

	if got != expected {
		t.Errorf("expected node to count %d children, counted %d", expected, got)
	}
}
