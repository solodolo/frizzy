package parser

import (
	"testing"
)

func getNodeWithBlockChildren(count int) *ElseIfListParseNode {
	node := &ElseIfListParseNode{}
	current := node

	for i := 0; i < count; i++ {
		next := &ElseIfListParseNode{}
		children := []TreeNode{
			next,
			&IdentParseNode{Value: "else_if"},
			&NonTerminalParseNode{Value: "expression"},
			&StringParseNode{Value: "}}"},
			&StringParseNode{Value: ""},
		}
		current.SetChildren(children)
		current = next
	}

	children := []TreeNode{
		&IdentParseNode{Value: "else_if"},
		&NonTerminalParseNode{Value: "expression"},
		&StringParseNode{Value: "}}"},
		&StringParseNode{Value: ""},
	}

	current.SetChildren(children)

	return node
}

func getNodeWithoutBlockChildren(count int) *ElseIfListParseNode {
	node := &ElseIfListParseNode{}
	var current TreeNode = node

	for i := 0; i < count; i++ {
		next := &IdentParseNode{Value: "else_if"}
		children := []TreeNode{
			next,
			&NonTerminalParseNode{Value: "expression"},
			&StringParseNode{Value: "}}"},
			&StringParseNode{Value: ""},
		}
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

func TestGetConditionalWithBlockChildrenReturnsCorrectNode(t *testing.T) {
	node := getNodeWithBlockChildren(3)
	conditional := node.GetConditional()
	typedConditional, ok := conditional.(*NonTerminalParseNode)

	if !ok {
		t.Errorf("expected conditional to be NonTerminalParseNode, got %T", conditional)
	} else if typedConditional.Value != "expression" {
		t.Errorf("expected conditional to have value 'expression', found %q", typedConditional.Value)
	}
}

func TestGetConditionalWithoutBlockChildrenReturnsCorrectNode(t *testing.T) {
	node := getNodeWithoutBlockChildren(3)
	conditional := node.GetConditional()
	typedConditional, ok := conditional.(*NonTerminalParseNode)

	if !ok {
		t.Errorf("expected conditional to be NonTerminalParseNode, got %T", conditional)
	} else if typedConditional.Value != "expression" {
		t.Errorf("expected conditional to have value 'expression', found %q", typedConditional.Value)
	}
}

func TestGetConditionalsWithBlockChildrenReturnsAllConditionals(t *testing.T) {
	numChildren := 4
	node := getNodeWithBlockChildren(numChildren)
	conditionals := node.GetConditionals()

	if len(conditionals) != numChildren+1 {
		t.Errorf("expected %d conditionals, got %d", numChildren+1, len(conditionals))
	} else {
		for _, conditional := range conditionals {
			if typedConditional, ok := conditional.(*NonTerminalParseNode); !ok || typedConditional.Value != "expression" {
				t.Errorf("expected NonTerminalParseNode with value expression, got %T", conditional)
			}
		}
	}
}

func TestGetConditionalsWithoutBlockChildrenReturnsAllConditionals(t *testing.T) {
	expected := 1
	node := getNodeWithoutBlockChildren(4)
	conditionals := node.GetConditionals()

	if len(conditionals) != expected {
		t.Errorf("expected %d conditionals, got %d", expected, len(conditionals))
	} else {
		typedConditional := conditionals[0].(*NonTerminalParseNode)
		if typedConditional.Value != "expression" {
			t.Errorf("expected NonTerminalParseNode with value expression, got %T", conditionals[0])
		}
	}
}

// bottom child should be first in flattened
// and head node should be last in flattened
func TestFlattenedChildrenAreCorrectlyOrdered(t *testing.T) {
	numChildren := 3
	node := getNodeWithBlockChildren(numChildren)
	flattened := node.GetFlattenedBlockChildren()

	if flattened[len(flattened)-1] != node {
		t.Errorf("expected last flattened else_if to be %p, got %p", node, flattened[len(flattened)-1])
	} else {
		// get last child
		current := node
		for i := 0; i < numChildren; i++ {
			current = current.children[0].(*ElseIfListParseNode)
		}

		if flattened[0] != current {
			t.Errorf("expected first flattened else_if to be %p, got %p ", current, flattened[0])
		}
	}
}
