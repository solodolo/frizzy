package parser

import (
	"fmt"
	"testing"

	"mettlach.codes/frizzy/lexer"
)

func TestIfStatementReturnsCorrectIfConditional(t *testing.T) {
	var tests = []struct {
		tokens []lexer.Token
		val    bool
	}{
		{
			[]lexer.Token{
				lexer.IfToken{},
				lexer.BoolToken{Value: "true"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>foo</p>"},
				lexer.EndToken{},
				lexer.EOLToken{},
			},
			true,
		},
		{
			[]lexer.Token{
				lexer.IfToken{},
				lexer.BoolToken{Value: "false"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>foo</p>"},
				lexer.ElseIfToken{},
				lexer.BoolToken{Value: "true"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>bar</p>"},
				lexer.EndToken{},
				lexer.EOLToken{},
			},
			false,
		},
		{
			[]lexer.Token{
				lexer.IfToken{},
				lexer.BoolToken{Value: "true"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>foo</p>"},
				lexer.ElseIfToken{},
				lexer.BoolToken{Value: "true"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>bar</p>"},
				lexer.ElseToken{},
				lexer.PassthroughToken{Value: "foo"},
				lexer.EndToken{},
				lexer.EOLToken{},
			},
			false,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("if tokens %d", i), func(t *testing.T) {
			ifTokens := test.tokens
			stateStack := []int{}
			nodeStack := []TreeNode{}

			_, head, _ := parseTokens(ifTokens, &stateStack, &nodeStack)
			ifStatement := extractToken(head, []int{0}).(*IfStatementParseNode)
			conditional := ifStatement.GetIfConditional()
			conditionalChildren := conditional.GetChildren()

			if len(conditionalChildren) != 1 {
				t.Errorf(
					"expected conditional to have 1 child, got %d",
					len(conditionalChildren),
				)
			} else {
				conditional = conditionalChildren[0]
				if _, ok := conditional.(*BoolParseNode); !ok {
					t.Errorf("expected BoolParseNode, got %T", conditional)
				}
			}
		})
	}
}

func TestIfStatementReturnsCorrectElseIfConditionals(t *testing.T) {
	var tests = []struct {
		tokens   []lexer.Token
		expected []TreeNode
	}{
		{
			[]lexer.Token{
				lexer.IfToken{},
				lexer.BoolToken{Value: "true"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>foo</p>"},
				lexer.EndToken{},
				lexer.EOLToken{},
			},
			[]TreeNode{},
		},
		{
			[]lexer.Token{
				lexer.IfToken{},
				lexer.BoolToken{Value: "false"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>foo</p>"},
				lexer.ElseIfToken{},
				lexer.BoolToken{Value: "true"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>bar</p>"},
				lexer.EndToken{},
				lexer.EOLToken{},
			},
			[]TreeNode{&BoolParseNode{Value: true}},
		},
		{
			[]lexer.Token{
				lexer.IfToken{},
				lexer.BoolToken{Value: "true"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>foo</p>"},
				lexer.ElseIfToken{},
				lexer.BoolToken{Value: "true"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>bar</p>"},
				lexer.ElseToken{},
				lexer.PassthroughToken{Value: "foo"},
				lexer.EndToken{},
				lexer.EOLToken{},
			},
			[]TreeNode{&BoolParseNode{Value: true}},
		},
		{
			[]lexer.Token{
				lexer.IfToken{},
				lexer.BoolToken{Value: "true"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>foo</p>"},
				lexer.ElseIfToken{},
				lexer.BoolToken{Value: "true"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>bar</p>"},
				lexer.ElseIfToken{},
				lexer.BoolToken{Value: "false"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>bar</p>"},
				lexer.EndToken{},
				lexer.EOLToken{},
			},
			[]TreeNode{&BoolParseNode{Value: true}, &BoolParseNode{Value: false}},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("else if conditionals %d", i), func(t *testing.T) {
			ifTokens := test.tokens
			stateStack := []int{}
			nodeStack := []TreeNode{}

			_, head, _ := parseTokens(ifTokens, &stateStack, &nodeStack)
			ifStatement := extractToken(head, []int{0}).(*IfStatementParseNode)
			conditionals := ifStatement.GetElseIfConditionals()
			conditionalChildren := make([]TreeNode, 0, len(conditionals))

			for _, conditional := range conditionals {
				children := conditional.GetChildren()
				if len(children) > 0 {
					conditionalChildren = append(conditionalChildren, children[0])
				}
			}

			if !nodeSlicesEqual(conditionalChildren, test.expected) {
				t.Errorf("expected %v to equal %v", conditionalChildren, test.expected)
			}
		})
	}
}

func TestIfStatementReturnsCorrectElseIfBody(t *testing.T) {
	var tests = []struct {
		tokens   []lexer.Token
		expected []TreeNode
	}{
		{
			[]lexer.Token{
				lexer.IfToken{},
				lexer.BoolToken{Value: "true"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>foo</p>"},
				lexer.EndToken{},
				lexer.EOLToken{},
			},
			[]TreeNode{},
		},
		{
			[]lexer.Token{
				lexer.IfToken{},
				lexer.BoolToken{Value: "false"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>foo</p>"},
				lexer.ElseIfToken{},
				lexer.BoolToken{Value: "true"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>bar</p>"},
				lexer.EndToken{},
				lexer.EOLToken{},
			},
			[]TreeNode{&StringParseNode{Value: "<p>bar</p>"}},
		},
		{
			[]lexer.Token{
				lexer.IfToken{},
				lexer.BoolToken{Value: "true"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>foo</p>"},
				lexer.ElseIfToken{},
				lexer.BoolToken{Value: "true"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "whatever"},
				lexer.ElseToken{},
				lexer.PassthroughToken{Value: "foo"},
				lexer.EndToken{},
				lexer.EOLToken{},
			},
			[]TreeNode{&StringParseNode{Value: "whatever"}},
		},
		{
			[]lexer.Token{
				lexer.IfToken{},
				lexer.BoolToken{Value: "true"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>foo</p>"},
				lexer.ElseIfToken{},
				lexer.BoolToken{Value: "true"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "first"},
				lexer.ElseIfToken{},
				lexer.BoolToken{Value: "false"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "second"},
				lexer.EndToken{},
				lexer.EOLToken{},
			},
			[]TreeNode{&StringParseNode{Value: "first"}, &StringParseNode{Value: "second"}},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("else if bodies %d", i), func(t *testing.T) {
			ifTokens := test.tokens
			stateStack := []int{}
			nodeStack := []TreeNode{}

			_, head, _ := parseTokens(ifTokens, &stateStack, &nodeStack)
			ifStatement := extractToken(head, []int{0}).(*IfStatementParseNode)

			for i := range test.expected {
				body, ok := ifStatement.GetElseIfBody(i)

				if !ok {
					t.Errorf("expected body at index %d, found none", i)
					continue
				}

				bodyChildren := body.GetChildren()
				if len(bodyChildren) != 1 {
					t.Errorf("expected 1 child, got %d", len(bodyChildren))
				} else if bodyChildren[0].String() != test.expected[i].String() {
					t.Errorf("expected %v, got %v", test.expected[i], body)
				}
			}
		})
	}
}
