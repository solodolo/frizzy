package parser

import (
	"fmt"
	"testing"
)

func getNodeWithBlockChildren(count int) ElseIfListParseNode {
	node := ElseIfListParseNode{}
	current := &node

	for i := 0; i < count; i++ {
		next := &ElseIfListParseNode{}
		children := []TreeNode{
			next,
			&IdentParseNode{Value: "else_if"},
			&StringParseNode{Value: "("},
			&NonTerminalParseNode{Value: "K"},
			&StringParseNode{Value: ")"},
		}
		current.SetChildren(children)
		current = next
	}

	children := []TreeNode{
		&IdentParseNode{Value: "else_if"},
		&StringParseNode{Value: "("},
		&NonTerminalParseNode{Value: "K"},
		&StringParseNode{Value: ")"},
	}

	current.SetChildren(children)

	return node
}

func getNodeWithoutBlockChildren(count int) ElseIfListParseNode {
	node := ElseIfListParseNode{}
	var current TreeNode = &node

	for i := 0; i < count; i++ {
		next := &IdentParseNode{Value: "else_if"}
		children := []TreeNode{
			next,
			&StringParseNode{Value: "("},
			&NonTerminalParseNode{Value: "K"},
			&StringParseNode{Value: ")"},
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

func TestGetNumBlockChildrenWithBlockChildrenReturnsCorrectNumber(t *testing.T) {
	expected := 5
	node := getNodeWithBlockChildren(expected)
	got := node.getNumBlockChildren()

	if got != expected {
		t.Errorf("expected node to count %d children, counted %d", expected, got)
	}
}

func TestGetConditionalWithBlockChildrenReturnsCorrectNode(t *testing.T) {
	node := getNodeWithBlockChildren(3)
	conditional := node.getConditional()
	typedConditional, ok := conditional.(*NonTerminalParseNode)

	if !ok {
		t.Errorf("expected conditional to be NonTerminalParseNode, got %T", conditional)
	} else if typedConditional.Value != "K" {
		t.Errorf("expected conditional to have value 'K', found %q", typedConditional.Value)
	}
}

func TestGetConditionalWithoutBlockChildrenReturnsCorrectNode(t *testing.T) {
	node := getNodeWithoutBlockChildren(3)
	conditional := node.getConditional()
	typedConditional, ok := conditional.(*NonTerminalParseNode)

	if !ok {
		t.Errorf("expected conditional to be NonTerminalParseNode, got %T", conditional)
	} else if typedConditional.Value != "K" {
		t.Errorf("expected conditional to have value 'K', found %q", typedConditional.Value)
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
			fmt.Println(conditional)
			if typedConditional, ok := conditional.(*NonTerminalParseNode); !ok || typedConditional.Value != "K" {
				t.Errorf("expected NonTerminalParseNode with value K, got %T", conditional)
			}
		}
	}
}
