package parser

import (
	"reflect"
	"testing"

	"mettlach.codes/frizzy/lexer"
)

func TestPassthroughTokensReturnCorrectNodeTypes(t *testing.T) {
	tokChan := make(chan []lexer.Token)

	go func() {
		defer close(tokChan)
		tokens := []lexer.Token{
			lexer.PassthroughToken{},
			lexer.PassthroughToken{},
			lexer.PassthroughToken{},
		}

		tokChan <- tokens
	}()

	nodes := []TreeNode{}
	nodeChan := make(chan TreeNode)
	errChan := make(chan error)

	go Parse(tokChan, nodeChan, errChan)
	go func() {
		for node := range nodeChan {
			nodes = append(nodes, node)
		}
	}()

	err := <-errChan

	if err != nil {
		t.Errorf("Expected no errors. Got %q.", err.Error())
	}

	expected := []TreeNode{
		StringParseNode{},
		StringParseNode{},
		StringParseNode{},
	}

	equal := nodeSlicesEqual(expected, nodes)

	if !equal {
		t.Errorf("Expected %v. Got %v.", expected, nodes)
	}
}

func TestNonPassthroughTokensReturnCorrectNodeTypes(t *testing.T) {
	var testing = []struct {
		tokens [][]lexer.Token
		nodes  []TreeNode
	}{
		{
			tokens: [][]lexer.Token{{
				lexer.BlockToken{Block: "{{"},
				lexer.BlockToken{Block: "}}"},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.BlockToken{Block: "{{"},
				lexer.StrToken{Str: "foo"},
				lexer.BlockToken{Block: "}}"},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.BlockToken{Block: "{{"},
				lexer.NumToken{Num: "5"},
				lexer.BlockToken{Block: "}}"},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.BlockToken{Block: "{{"},
				lexer.IfToken{},
				lexer.SymbolToken{Symbol: "("},
				lexer.BoolToken{Value: "false"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.EndToken{},
				lexer.BlockToken{Block: "}}"},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.BlockToken{Block: "{{"},
				lexer.IfToken{},
				lexer.SymbolToken{Symbol: "("},
				lexer.BoolToken{Value: "false"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.ElseIfToken{},
				lexer.SymbolToken{Symbol: "("},
				lexer.BoolToken{Value: "true"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.EndToken{},
				lexer.BlockToken{Block: "}}"},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{
				{
					lexer.BlockToken{Block: "{{"},
					lexer.IfToken{},
					lexer.SymbolToken{Symbol: "("},
					lexer.BoolToken{Value: "false"},
					lexer.SymbolToken{Symbol: ")"},
				},
				{
					lexer.ElseIfToken{},
					lexer.SymbolToken{Symbol: "("},
					lexer.BoolToken{Value: "true"},
					lexer.SymbolToken{Symbol: ")"},
					lexer.EndToken{},
					lexer.BlockToken{Block: "}}"},
					lexer.EOLToken{},
				},
			},
			nodes: []TreeNode{
				NonTerminalParseNode{},
			},
		},
	}

	for _, test := range testing {
		tokChan := make(chan []lexer.Token)

		go func(tokens [][]lexer.Token) {
			defer close(tokChan)

			for _, tokenLine := range tokens {
				tokChan <- tokenLine
			}
		}(test.tokens)

		nodes := []TreeNode{}
		nodeChan := make(chan TreeNode)
		errChan := make(chan error)

		go Parse(tokChan, nodeChan, errChan)
		for node := range nodeChan {
			nodes = append(nodes, node)
		}

		err := <-errChan

		if err != nil {
			t.Errorf("Expected no errors. Got %q.", err.Error())
		}

		if !nodeSlicesEqual(test.nodes, nodes) {
			t.Errorf("Expected %v. Got %v.", test.nodes, nodes)
		}
	}
}

func nodeSlicesEqual(expected, got []TreeNode) bool {
	if len(expected) != len(got) {
		return false
	}

	for i := 0; i < len(expected); i++ {
		if reflect.TypeOf(expected[i]) != reflect.TypeOf(got[i]) {
			return false
		}
	}

	return true
}
