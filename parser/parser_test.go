package parser

import (
	"reflect"
	"strings"
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
	for node := range nodeChan {
		nodes = append(nodes, node)
	}

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
	var tests = []struct {
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
		{
			tokens: [][]lexer.Token{
				{
					lexer.BlockToken{Block: "{{"},
					lexer.BlockToken{Block: "}}"},
					lexer.EOLToken{},
				},
				{
					lexer.BlockToken{Block: "{{:"},
					lexer.StrToken{Str: "bar"},
					lexer.BlockToken{Block: "}}"},
					lexer.EOLToken{},
				},
			},
			nodes: []TreeNode{
				NonTerminalParseNode{},
				NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{
				{
					lexer.BlockToken{Block: "{{"},
					lexer.BlockToken{Block: "}}"},
					lexer.EOLToken{},
					lexer.PassthroughToken{},
				},
				{
					lexer.BlockToken{Block: "{{:"},
					lexer.StrToken{Str: "bar"},
					lexer.BlockToken{Block: "}}"},
					lexer.EOLToken{},
				},
			},
			nodes: []TreeNode{
				NonTerminalParseNode{},
				StringParseNode{},
				NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{
				{
					lexer.BlockToken{Block: "{{"},
					lexer.NumToken{Num: "5"},
					lexer.MultOpToken{Operator: "*"},
					lexer.NumToken{Num: "4"},
					lexer.BlockToken{Block: "}}"},
					lexer.EOLToken{},
				},
			},
			nodes: []TreeNode{
				NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{
				{
					lexer.BlockToken{Block: "{{"},
					lexer.UnaryOpToken{Operator: "!"},
					lexer.BoolToken{},
					lexer.BlockToken{Block: "}}"},
					lexer.EOLToken{},
				},
			},
			nodes: []TreeNode{
				NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{
				{
					lexer.BlockToken{Block: "{{:"},
					lexer.VarToken{Variable: "post.title"},
					lexer.BlockToken{Block: "}}"},
					lexer.EOLToken{},
				},
			},
			nodes: []TreeNode{
				NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{
				{
					lexer.BlockToken{Block: "{{:"},
					lexer.IdentToken{Identifier: "bar"},
					lexer.BlockToken{Block: "}}"},
					lexer.EOLToken{},
				},
			},
			nodes: []TreeNode{
				NonTerminalParseNode{},
			},
		},
	}

	for _, test := range tests {
		tokChan := getTokChan(test.tokens)
		nodeChan := make(chan TreeNode)
		errChan := make(chan error)

		go Parse(tokChan, nodeChan, errChan)

		nodes := []TreeNode{}
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

func TestParserSendsErrorWithIncorrectToken(t *testing.T) {
	var tests = []struct {
		tokens [][]lexer.Token
		errMsg string
	}{
		{
			tokens: [][]lexer.Token{
				{
					lexer.BlockToken{Block: "{{"},
					lexer.BlockToken{Block: "{{"},
				},
			},
			errMsg: `unexpected symbol "{{"`,
		},
		{
			tokens: [][]lexer.Token{
				{
					lexer.BlockToken{Block: "{{:"},
					lexer.BlockToken{Block: "}}"},
				},
			},
			errMsg: `unexpected symbol "}}"`,
		},
		{
			tokens: [][]lexer.Token{
				{
					lexer.BlockToken{Block: "{{"},
					lexer.ForToken{},
					lexer.MultOpToken{Operator: "*"},
				},
			},
			errMsg: `unexpected symbol "MULT_OP"`,
		},
		{
			tokens: [][]lexer.Token{
				{
					lexer.BlockToken{Block: "{{"},
					lexer.ForToken{},
					lexer.MultOpToken{Operator: "/"},
				},
			},
			errMsg: `unexpected symbol "MULT_OP"`,
		},
		{
			tokens: [][]lexer.Token{
				{
					lexer.BlockToken{Block: "{{"},
					lexer.ForToken{},
					lexer.MultOpToken{Operator: "%"},
				},
			},
			errMsg: `unexpected symbol "MULT_OP"`,
		},
		{
			tokens: [][]lexer.Token{
				{
					lexer.BlockToken{Block: "{{"},
					lexer.ForToken{},
					lexer.AddOpToken{Operator: "+"},
				},
			},
			errMsg: `unexpected symbol "ADD_OP"`,
		},
		{
			tokens: [][]lexer.Token{
				{
					lexer.BlockToken{Block: "{{"},
					lexer.ForToken{},
					lexer.AddOpToken{Operator: "-"},
				},
			},
			errMsg: `unexpected symbol "ADD_OP"`,
		},
		{
			tokens: [][]lexer.Token{
				{
					lexer.BlockToken{Block: "{{"},
					lexer.SymbolToken{Symbol: "^"},
				},
			},
			errMsg: `unrecognized symbol "^"`,
		},
		{
			tokens: [][]lexer.Token{
				{
					lexer.BlockToken{Block: "{{"},
					lexer.BlockToken{Block: "}}"},
				},
				{
					lexer.SymbolToken{Symbol: "?"},
				},
			},
			errMsg: `unrecognized symbol "?"`,
		},
	}

	for _, test := range tests {
		tokChan := getTokChan(test.tokens)
		nodeChan := make(chan TreeNode)
		errChan := make(chan error)
		go Parse(tokChan, nodeChan, errChan)
		go func() { <-nodeChan }()

		err := <-errChan

		if err == nil {
			t.Errorf("Expected error %q. Got nil.", test.errMsg)
		}

		gotMsg := err.Error()
		if strings.Index(gotMsg, test.errMsg) == -1 {
			t.Errorf("Expected error %q to include %q.", gotMsg, test.errMsg)
		}
	}
}

func TestParserBuildsCorrectTree(t *testing.T) {
	var tests = []struct {
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
	}

	for _, test := range tests {
		tokChan := getTokChan(test.tokens)
		nodeChan := make(chan TreeNode)
		errChan := make(chan error)

		go Parse(tokChan, nodeChan, errChan)

		nodes := []TreeNode{}

		go func() {
			for node := range nodeChan {
				nodes = append(nodes, node)
			}
		}()

		if err := <-errChan; err != nil {
			t.Errorf("Expected no errors. Got %q.", err.Error())
		}

		flatChildren := flattenChildren(nodes[0])

		if !nodeSlicesEqual(flatChildren, test.nodes) {
			t.Errorf("Expected %#v. Got %#v.", test.nodes, nodes)
		}
	}
}

func getTokChan(tokens [][]lexer.Token) chan []lexer.Token {
	tokChan := make(chan []lexer.Token)

	go func(tokens [][]lexer.Token) {
		defer close(tokChan)

		for _, tokenLine := range tokens {
			tokChan <- tokenLine
		}
	}(tokens)

	return tokChan
}

func flattenChildren(node TreeNode) []TreeNode {
	allChildren := node.GetChildren()
	for _, child := range allChildren {
		allChildren = append(allChildren, flattenChildren(child)...)
	}

	return allChildren
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
