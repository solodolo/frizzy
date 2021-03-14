package parser

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"mettlach.codes/frizzy/lexer"
)

func testParsesNoErrors(test struct {
	tokens [][]lexer.Token
	nodes  []TreeNode
}, t *testing.T) {
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

	if equal, msg := nodeSlicesEqual(test.nodes, nodes); !equal {
		t.Error(msg)
	}
}

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
		&StringParseNode{},
		&StringParseNode{},
		&StringParseNode{},
	}

	if equal, msg := nodeSlicesEqual(expected, nodes); !equal {
		t.Error(msg)
	}
}

func TestIfBlockParsesNoErrors(t *testing.T) {
	var tests = []struct {
		tokens [][]lexer.Token
		nodes  []TreeNode
	}{
		{
			tokens: [][]lexer.Token{{
				lexer.IfToken{},
				lexer.SymbolToken{Symbol: "("},
				lexer.BoolToken{Value: "true"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>foo</p>"},
				lexer.EndToken{},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.IfToken{},
				lexer.SymbolToken{Symbol: "("},
				lexer.BoolToken{Value: "true"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>foo</p>"},
				lexer.PassthroughToken{Value: "<p>bar</p>"},
				lexer.EndToken{},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.IfToken{},
				lexer.SymbolToken{Symbol: "("},
				lexer.BoolToken{Value: "true"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.BlockToken{Block: "{{"},
				lexer.IdentToken{Identifier: "foo"},
				lexer.BlockToken{Block: "}}"},
				lexer.EndToken{},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.IfToken{},
				lexer.SymbolToken{Symbol: "("},
				lexer.BoolToken{Value: "true"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.BlockToken{Block: "{{"},
				lexer.IdentToken{Identifier: "foo.bar"},
				lexer.BlockToken{Block: "}}"},
				lexer.BlockToken{Block: "{{"},
				lexer.NumToken{Num: "42"},
				lexer.BlockToken{Block: "}}"},
				lexer.EndToken{},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.IfToken{},
				lexer.SymbolToken{Symbol: "("},
				lexer.BoolToken{Value: "true"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.BlockToken{Block: "{{"},
				lexer.IdentToken{Identifier: "foo.bar"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<span>abc</span>"},
				lexer.BlockToken{Block: "{{"},
				lexer.NumToken{Num: "42"},
				lexer.BlockToken{Block: "}}"},
				lexer.EndToken{},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.IfToken{},
				lexer.SymbolToken{Symbol: "("},
				lexer.BoolToken{Value: "true"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.IfToken{},
				lexer.SymbolToken{Symbol: "("},
				lexer.BoolToken{Value: "true"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>foo</p>"},
				lexer.EndToken{},
				lexer.EndToken{},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
	}

	for _, test := range tests {
		testParsesNoErrors(test, t)
	}
}

func TestIfElseIfParsesNoErrors(t *testing.T) {
	var tests = []struct {
		tokens [][]lexer.Token
		nodes  []TreeNode
	}{
		{
			tokens: [][]lexer.Token{{
				lexer.IfToken{},
				lexer.SymbolToken{Symbol: "("},
				lexer.BoolToken{Value: "true"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>foo</p>"},
				lexer.ElseIfToken{},
				lexer.SymbolToken{Symbol: "("},
				lexer.BoolToken{Value: "true"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>bar</p>"},
				lexer.EndToken{},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.IfToken{},
				lexer.SymbolToken{Symbol: "("},
				lexer.BoolToken{Value: "true"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>foo</p>"},
				lexer.ElseIfToken{},
				lexer.SymbolToken{Symbol: "("},
				lexer.BoolToken{Value: "true"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.BlockToken{Block: "{{"},
				lexer.NumToken{Num: "42"},
				lexer.SubOpToken{},
				lexer.NumToken{Num: "41"},
				lexer.BlockToken{Block: "}}"},
				lexer.EndToken{},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.IfToken{},
				lexer.SymbolToken{Symbol: "("},
				lexer.BoolToken{Value: "true"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>foo</p>"},
				lexer.ElseIfToken{},
				lexer.SymbolToken{Symbol: "("},
				lexer.BoolToken{Value: "true"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.BlockToken{Block: "{{"},
				lexer.NumToken{Num: "42"},
				lexer.SubOpToken{},
				lexer.NumToken{Num: "41"},
				lexer.BlockToken{Block: "}}"},
				lexer.ElseIfToken{},
				lexer.SymbolToken{Symbol: "("},
				lexer.BoolToken{Value: "true"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.BlockToken{Block: "{{"},
				lexer.NumToken{Num: "42"},
				lexer.SubOpToken{},
				lexer.NumToken{Num: "41"},
				lexer.BlockToken{Block: "}}"},
				lexer.EndToken{},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
	}

	for _, test := range tests {
		testParsesNoErrors(test, t)
	}
}

func TestIfElseIfElseParsesNoErrors(t *testing.T) {
	var tests = []struct {
		tokens [][]lexer.Token
		nodes  []TreeNode
	}{
		{
			tokens: [][]lexer.Token{{
				lexer.IfToken{},
				lexer.SymbolToken{Symbol: "("},
				lexer.BoolToken{Value: "true"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>foo</p>"},
				lexer.ElseIfToken{},
				lexer.SymbolToken{Symbol: "("},
				lexer.BoolToken{Value: "true"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.BlockToken{Block: "{{"},
				lexer.NumToken{Num: "42"},
				lexer.SubOpToken{},
				lexer.NumToken{Num: "41"},
				lexer.BlockToken{Block: "}}"},
				lexer.ElseIfToken{},
				lexer.SymbolToken{Symbol: "("},
				lexer.BoolToken{Value: "true"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.BlockToken{Block: "{{"},
				lexer.NumToken{Num: "42"},
				lexer.SubOpToken{},
				lexer.NumToken{Num: "41"},
				lexer.BlockToken{Block: "}}"},
				lexer.ElseToken{},
				lexer.BlockToken{Block: "{{"},
				lexer.NumToken{Num: "42"},
				lexer.AddOpToken{},
				lexer.NumToken{Num: "41"},
				lexer.BlockToken{Block: "}}"},
				lexer.EndToken{},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.IfToken{},
				lexer.SymbolToken{Symbol: "("},
				lexer.BoolToken{Value: "true"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>foo</p>"},
				lexer.ElseIfToken{},
				lexer.SymbolToken{Symbol: "("},
				lexer.BoolToken{Value: "true"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},

				lexer.IfToken{},
				lexer.SymbolToken{Symbol: "("},
				lexer.BoolToken{Value: "true"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>foo</p>"},
				lexer.EndToken{},

				lexer.ElseIfToken{},
				lexer.SymbolToken{Symbol: "("},
				lexer.BoolToken{Value: "true"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},

				lexer.IfToken{},
				lexer.SymbolToken{Symbol: "("},
				lexer.BoolToken{Value: "true"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>foo</p>"},
				lexer.EndToken{},

				lexer.ElseToken{},

				lexer.IfToken{},
				lexer.SymbolToken{Symbol: "("},
				lexer.BoolToken{Value: "true"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>foo</p>"},
				lexer.EndToken{},

				lexer.EndToken{},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.IfToken{},
				lexer.SymbolToken{Symbol: "("},
				lexer.BoolToken{Value: "true"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>foo</p>"},
				lexer.ElseIfToken{},
				lexer.SymbolToken{Symbol: "("},
				lexer.BoolToken{Value: "true"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},

				lexer.IfToken{},
				lexer.SymbolToken{Symbol: "("},
				lexer.BoolToken{Value: "true"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>foo</p>"},
				lexer.ElseIfToken{},
				lexer.SymbolToken{Symbol: "("},
				lexer.BoolToken{Value: "true"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>foo</p>"},
				lexer.EndToken{},

				lexer.ElseIfToken{},
				lexer.SymbolToken{Symbol: "("},
				lexer.BoolToken{Value: "true"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},

				lexer.PassthroughToken{Value: "<p>foo</p>"},

				lexer.ElseToken{},

				lexer.IfToken{},
				lexer.SymbolToken{Symbol: "("},
				lexer.BoolToken{Value: "true"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>foo</p>"},
				lexer.ElseToken{},
				lexer.PassthroughToken{Value: "<p>foo</p>"},
				lexer.EndToken{},

				lexer.EndToken{},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
	}

	for _, test := range tests {
		testParsesNoErrors(test, t)
	}
}

func TestIfElseParsesNoErrors(t *testing.T) {
	var tests = []struct {
		tokens [][]lexer.Token
		nodes  []TreeNode
	}{
		{
			tokens: [][]lexer.Token{{
				lexer.IfToken{},
				lexer.SymbolToken{Symbol: "("},
				lexer.BoolToken{Value: "true"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>foo</p>"},
				lexer.ElseToken{},
				lexer.PassthroughToken{Value: "<p>bar</p>"},
				lexer.EndToken{},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
	}

	for _, test := range tests {
		testParsesNoErrors(test, t)
	}
}

func TestForLoopParsesNoErrors(t *testing.T) {
	var tests = []struct {
		tokens [][]lexer.Token
		nodes  []TreeNode
	}{
		{
			tokens: [][]lexer.Token{{
				lexer.ForToken{},
				lexer.IdentToken{Identifier: "foo"},
				lexer.InToken{},
				lexer.StrToken{Str: "bar"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<p>foo</p>"},
				lexer.PassthroughToken{Value: "<p>bar</p>"},
				lexer.EndToken{},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.ForToken{},
				lexer.IdentToken{Identifier: "foo"},
				lexer.InToken{},
				lexer.StrToken{Str: "bar"},
				lexer.BlockToken{Block: "}}"},
				lexer.BlockToken{Block: "{{"},
				lexer.IdentToken{Identifier: "baz"},
				lexer.BlockToken{Block: "}}"},
				lexer.BlockToken{Block: "{{"},
				lexer.IdentToken{Identifier: "fozz"},
				lexer.BlockToken{Block: "}}"},
				lexer.EndToken{},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.ForToken{},
				lexer.IdentToken{Identifier: "foo"},
				lexer.InToken{},
				lexer.StrToken{Str: "bar"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<h1>bar</h1>"},
				lexer.BlockToken{Block: "{{"},
				lexer.IdentToken{Identifier: "baz"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<h1>bar</h1>"},
				lexer.BlockToken{Block: "{{"},
				lexer.IdentToken{Identifier: "fozz"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "<h1>bar</h1>"},
				lexer.EndToken{},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.ForToken{},
				lexer.IdentToken{Identifier: "foo"},
				lexer.InToken{},
				lexer.StrToken{Str: "bar"},
				lexer.BlockToken{Block: "}}"},
				lexer.ForToken{},
				lexer.IdentToken{Identifier: "fizz"},
				lexer.InToken{},
				lexer.StrToken{Str: "buzz"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "foo"},
				lexer.EndToken{},
				lexer.EndToken{},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.ForToken{},
				lexer.IdentToken{Identifier: "foo"},
				lexer.InToken{},
				lexer.IdentToken{Identifier: "print"},
				lexer.SymbolToken{Symbol: "("},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "foo"},
				lexer.EndToken{},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.ForToken{},
				lexer.IdentToken{Identifier: "foo"},
				lexer.InToken{},
				lexer.IdentToken{Identifier: "print"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "foo"},
				lexer.EndToken{},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.ForToken{},
				lexer.IdentToken{Identifier: "foo"},
				lexer.InToken{},
				lexer.IdentToken{Identifier: "post"},
				lexer.SymbolToken{Symbol: "."},
				lexer.IdentToken{Identifier: "title"},
				lexer.BlockToken{Block: "}}"},
				lexer.PassthroughToken{Value: "foo"},
				lexer.EndToken{},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
	}

	for _, test := range tests {
		testParsesNoErrors(test, t)
	}
}

func TestBlockParsesNoErrors(t *testing.T) {
	var tests = []struct {
		tokens [][]lexer.Token
		nodes  []TreeNode
	}{
		{
			tokens: [][]lexer.Token{{
				lexer.BlockToken{Block: "{{"},
				lexer.StrToken{Str: "foo"},
				lexer.BlockToken{Block: "}}"},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.BlockToken{Block: "{{"},
				lexer.IdentToken{Identifier: "foo"},
				lexer.BlockToken{Block: "}}"},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.BlockToken{Block: "{{"},
				lexer.StrToken{Str: "foo"},
				lexer.AddOpToken{},
				lexer.StrToken{Str: "bar"},
				lexer.BlockToken{Block: "}}"},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.BlockToken{Block: "{{"},
				lexer.NegationOpToken{},
				lexer.BoolToken{Value: "true"},
				lexer.BlockToken{Block: "}}"},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.BlockToken{Block: "{{"},
				lexer.SubOpToken{},
				lexer.NumToken{Num: "10"},
				lexer.BlockToken{Block: "}}"},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.BlockToken{Block: "{{"},
				lexer.IdentToken{Identifier: "print"},
				lexer.SymbolToken{Symbol: "("},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.BlockToken{Block: "{{"},
				lexer.IdentToken{Identifier: "print"},
				lexer.SymbolToken{Symbol: "("},
				lexer.StrToken{Str: "foo"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.BlockToken{Block: "{{"},
				lexer.IdentToken{Identifier: "print"},
				lexer.SymbolToken{Symbol: "("},
				lexer.StrToken{Str: "foo"},
				lexer.SymbolToken{Symbol: ","},
				lexer.NumToken{Num: "5"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},

		{
			tokens: [][]lexer.Token{{
				lexer.BlockToken{Block: "{{"},
				lexer.IdentToken{Identifier: "print"},
				lexer.SymbolToken{Symbol: "("},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
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

		if equal, msg := nodeSlicesEqual(test.nodes, nodes); !equal {
			t.Error(msg)
		}
	}
}

func TestBlocksAndPassthroughsParsesNoErrors(t *testing.T) {
	var tests = []struct {
		tokens [][]lexer.Token
		nodes  []TreeNode
	}{
		{
			tokens: [][]lexer.Token{{
				lexer.PassthroughToken{},
				lexer.BlockToken{Block: "{{"},
				lexer.StrToken{Str: "foo"},
				lexer.BlockToken{Block: "}}"},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.BlockToken{Block: "{{"},
				lexer.IdentToken{Identifier: "foo"},
				lexer.BlockToken{Block: "}}"},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.BlockToken{Block: "{{"},
				lexer.StrToken{Str: "foo"},
				lexer.AddOpToken{},
				lexer.StrToken{Str: "bar"},
				lexer.BlockToken{Block: "}}"},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.BlockToken{Block: "{{"},
				lexer.NegationOpToken{},
				lexer.BoolToken{Value: "true"},
				lexer.BlockToken{Block: "}}"},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.BlockToken{Block: "{{"},
				lexer.SubOpToken{},
				lexer.NumToken{Num: "10"},
				lexer.BlockToken{Block: "}}"},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.BlockToken{Block: "{{"},
				lexer.IdentToken{Identifier: "print"},
				lexer.SymbolToken{Symbol: "("},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.BlockToken{Block: "{{"},
				lexer.IdentToken{Identifier: "print"},
				lexer.SymbolToken{Symbol: "("},
				lexer.StrToken{Str: "foo"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.BlockToken{Block: "{{"},
				lexer.IdentToken{Identifier: "print"},
				lexer.SymbolToken{Symbol: "("},
				lexer.StrToken{Str: "foo"},
				lexer.SymbolToken{Symbol: ","},
				lexer.NumToken{Num: "5"},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
			},
		},

		{
			tokens: [][]lexer.Token{{
				lexer.BlockToken{Block: "{{"},
				lexer.IdentToken{Identifier: "print"},
				lexer.SymbolToken{Symbol: "("},
				lexer.SymbolToken{Symbol: ")"},
				lexer.BlockToken{Block: "}}"},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
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

		if equal, msg := nodeSlicesEqual(test.nodes, nodes); !equal {
			t.Error(msg)
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
					lexer.AddOpToken{},
				},
			},
			errMsg: `unexpected symbol "ADD_OP"`,
		},
		{
			tokens: [][]lexer.Token{
				{
					lexer.BlockToken{Block: "{{"},
					lexer.ForToken{},
					lexer.SubOpToken{},
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
		{
			tokens: [][]lexer.Token{
				{
					lexer.BlockToken{Block: "{{"},
				},
				{
					lexer.SymbolToken{Symbol: "?"},
					lexer.BlockToken{Block: "}}"},
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
		<-nodeChan

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
				&NonTerminalParseNode{},
				&NonTerminalParseNode{Value: "B"},
				&NonTerminalParseNode{Value: "C"},
				&NonTerminalParseNode{Value: "D"},
				&BlockParseNode{Value: "{{"},
				&NonTerminalParseNode{Value: "V"},
				&BlockParseNode{Value: "}}"},
			},
		},
		{
			tokens: [][]lexer.Token{
				{
					lexer.BlockToken{Block: "{{:"},
					lexer.IdentToken{Identifier: "post"},
					lexer.SymbolToken{Symbol: "."},
					lexer.IdentToken{Identifier: "title"},
					lexer.BlockToken{Block: "}}"},
					lexer.EOLToken{},
				},
			},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
				&NonTerminalParseNode{Value: "B"},
				&NonTerminalParseNode{Value: "C"},
				&NonTerminalParseNode{Value: "E"},
				&BlockParseNode{Value: "{{:"},
				&NonTerminalParseNode{Value: "F"},
				&BlockParseNode{Value: "}}"},
				&NonTerminalParseNode{Value: "K"},
				&NonTerminalParseNode{Value: "L"},
				&NonTerminalParseNode{Value: "M"},
				&NonTerminalParseNode{Value: "U"},
				&NonTerminalParseNode{Value: "N"},
				&NonTerminalParseNode{Value: "O"},
				&NonTerminalParseNode{Value: "P"},
				&IdentParseNode{Value: "post.title"},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.BlockToken{Block: "{{:"},
				lexer.StrToken{Str: "foobar"},
				lexer.BlockToken{Block: "}}"},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
				&NonTerminalParseNode{Value: "B"},
				&NonTerminalParseNode{Value: "C"},
				&NonTerminalParseNode{Value: "E"},
				&BlockParseNode{Value: "{{:"},
				&NonTerminalParseNode{Value: "F"},
				&BlockParseNode{Value: "}}"},
				&NonTerminalParseNode{Value: "K"},
				&NonTerminalParseNode{Value: "L"},
				&NonTerminalParseNode{Value: "M"},
				&NonTerminalParseNode{Value: "U"},
				&NonTerminalParseNode{Value: "N"},
				&NonTerminalParseNode{Value: "O"},
				&NonTerminalParseNode{Value: "P"},
				&StringParseNode{Value: "foobar"},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.BlockToken{Block: "{{"},
				lexer.NumToken{Num: "2"},
				lexer.AddOpToken{},
				lexer.NumToken{Num: "5"},
				lexer.BlockToken{Block: "}}"},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
				&NonTerminalParseNode{Value: "B"},
				&NonTerminalParseNode{Value: "C"},
				&NonTerminalParseNode{Value: "D"},
				&BlockParseNode{Value: "{{"},
				&NonTerminalParseNode{Value: "F"},
				&BlockParseNode{Value: "}}"},
				&NonTerminalParseNode{Value: "K"},
				&NonTerminalParseNode{Value: "L"},
				&NonTerminalParseNode{Value: "M"},
				&NonTerminalParseNode{Value: "U"},
				&NonTerminalParseNode{Value: "N"},
				&NonTerminalParseNode{Value: "N"},
				&StringParseNode{Value: "+"},
				&NonTerminalParseNode{Value: "O"},
				&NonTerminalParseNode{Value: "O"},
				&NonTerminalParseNode{Value: "P"},
				&NonTerminalParseNode{Value: "P"},
				&NumParseNode{Value: 5},
				&NumParseNode{Value: 2},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.BlockToken{Block: "{{"},
				lexer.NumToken{Num: "2"},
				lexer.AddOpToken{},
				lexer.NumToken{Num: "5"},
				lexer.BlockToken{Block: "}}"},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
				&NonTerminalParseNode{Value: "B"},
				&NonTerminalParseNode{Value: "C"},
				&NonTerminalParseNode{Value: "D"},
				&BlockParseNode{Value: "{{"},
				&NonTerminalParseNode{Value: "F"},
				&BlockParseNode{Value: "}}"},
				&NonTerminalParseNode{Value: "K"},
				&NonTerminalParseNode{Value: "L"},
				&NonTerminalParseNode{Value: "M"},
				&NonTerminalParseNode{Value: "U"},
				&NonTerminalParseNode{Value: "N"},
				&NonTerminalParseNode{Value: "N"},
				&StringParseNode{Value: "+"},
				&NonTerminalParseNode{Value: "O"},
				&NonTerminalParseNode{Value: "O"},
				&NonTerminalParseNode{Value: "P"},
				&NonTerminalParseNode{Value: "P"},
				&NumParseNode{Value: 5},
				&NumParseNode{Value: 2},
			},
		},
		{
			tokens: [][]lexer.Token{{
				lexer.BlockToken{Block: "{{"},
				lexer.NumToken{Num: "2"},
				lexer.SubOpToken{},
				lexer.NumToken{Num: "5"},
				lexer.BlockToken{Block: "}}"},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
				&NonTerminalParseNode{Value: "B"},
				&NonTerminalParseNode{Value: "C"},
				&NonTerminalParseNode{Value: "D"},
				&BlockParseNode{Value: "{{"},
				&NonTerminalParseNode{Value: "F"},
				&BlockParseNode{Value: "}}"},
				&NonTerminalParseNode{Value: "K"},
				&NonTerminalParseNode{Value: "L"},
				&NonTerminalParseNode{Value: "M"},
				&NonTerminalParseNode{Value: "U"},
				&NonTerminalParseNode{Value: "N"},
				&NonTerminalParseNode{Value: "N"},
				&StringParseNode{Value: "-"},
				&NonTerminalParseNode{Value: "O"},
				&NonTerminalParseNode{Value: "O"},
				&NonTerminalParseNode{Value: "P"},
				&NonTerminalParseNode{Value: "P"},
				&NumParseNode{Value: 5},
				&NumParseNode{Value: 2},
			},
		},
		{
			tokens: [][]lexer.Token{
				{
					lexer.BlockToken{Block: "{{"},
					lexer.NumToken{Num: "2"},
				},
				{
					lexer.AddOpToken{},
					lexer.NumToken{Num: "5"},
				},
				{
					lexer.BlockToken{Block: "}}"},
					lexer.EOLToken{},
				},
			},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
				&NonTerminalParseNode{Value: "B"},
				&NonTerminalParseNode{Value: "C"},
				&NonTerminalParseNode{Value: "D"},
				&BlockParseNode{Value: "{{"},
				&NonTerminalParseNode{Value: "F"},
				&BlockParseNode{Value: "}}"},
				&NonTerminalParseNode{Value: "K"},
				&NonTerminalParseNode{Value: "L"},
				&NonTerminalParseNode{Value: "M"},
				&NonTerminalParseNode{Value: "U"},
				&NonTerminalParseNode{Value: "N"},
				&NonTerminalParseNode{Value: "N"},
				&StringParseNode{Value: "+"},
				&NonTerminalParseNode{Value: "O"},
				&NonTerminalParseNode{Value: "O"},
				&NonTerminalParseNode{Value: "P"},
				&NonTerminalParseNode{Value: "P"},
				&NumParseNode{Value: 5},
				&NumParseNode{Value: 2},
			},
		},
		{
			tokens: [][]lexer.Token{
				{
					lexer.BlockToken{Block: "{{"},
				},
				{
					lexer.BlockToken{Block: "}}"},
					lexer.EOLToken{},
				},
			},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
				&NonTerminalParseNode{Value: "B"},
				&NonTerminalParseNode{Value: "C"},
				&NonTerminalParseNode{Value: "D"},
				&BlockParseNode{Value: "{{"},
				&NonTerminalParseNode{Value: "V"},
				&BlockParseNode{Value: "}}"},
			},
		},
		{
			tokens: [][]lexer.Token{
				{
					lexer.BlockToken{Block: "{{"},
				},
				{
					lexer.IfToken{},
					lexer.SymbolToken{Symbol: "("},
					lexer.IdentToken{Identifier: "foo"},
					lexer.RelOpToken{Operator: "<"},
					lexer.IdentToken{Identifier: "bar"},
					lexer.SymbolToken{Symbol: ")"},
				},
				{
					lexer.IdentToken{Identifier: "print"},
					lexer.SymbolToken{Symbol: "("},
					lexer.StrToken{Str: "foo"},
					lexer.SymbolToken{Symbol: ")"},
					lexer.SymbolToken{Symbol: ";"},
					lexer.IdentToken{Identifier: "print"},
					lexer.SymbolToken{Symbol: "("},
					lexer.StrToken{Str: "not bar"},
					lexer.SymbolToken{Symbol: ")"},
					lexer.SymbolToken{Symbol: ";"},
				},
				{
					lexer.ElseIfToken{},
					lexer.SymbolToken{Symbol: "("},
					lexer.IdentToken{Identifier: "foo"},
					lexer.RelOpToken{Operator: ">"},
					lexer.IdentToken{Identifier: "bar"},
					lexer.SymbolToken{Symbol: ")"},
				},
				{
					lexer.IdentToken{Identifier: "print"},
					lexer.SymbolToken{Symbol: "("},
					lexer.StrToken{Str: "bar"},
					lexer.SymbolToken{Symbol: ")"},
					lexer.SymbolToken{Symbol: ";"},
					lexer.IdentToken{Identifier: "print"},
					lexer.SymbolToken{Symbol: "("},
					lexer.StrToken{Str: "not foo"},
					lexer.SymbolToken{Symbol: ")"},
					lexer.SymbolToken{Symbol: ";"},
				},
				{
					lexer.ElseIfToken{},
					lexer.SymbolToken{Symbol: "("},
					lexer.IdentToken{Identifier: "foo"},
					lexer.RelOpToken{Operator: ">"},
					lexer.IdentToken{Identifier: "bar"},
					lexer.SymbolToken{Symbol: ")"},
				},
				{
					lexer.IdentToken{Identifier: "print"},
					lexer.SymbolToken{Symbol: "("},
					lexer.StrToken{Str: "bar"},
					lexer.SymbolToken{Symbol: ")"},
					lexer.SymbolToken{Symbol: ";"},
					lexer.IdentToken{Identifier: "print"},
					lexer.SymbolToken{Symbol: "("},
					lexer.StrToken{Str: "not foo"},
					lexer.SymbolToken{Symbol: ")"},
					lexer.SymbolToken{Symbol: ";"},
				},
				{
					lexer.ElseIfToken{},
					lexer.SymbolToken{Symbol: "("},
					lexer.IdentToken{Identifier: "foo"},
					lexer.RelOpToken{Operator: ">"},
					lexer.IdentToken{Identifier: "bar"},
					lexer.SymbolToken{Symbol: ")"},
				},
				{
					lexer.IdentToken{Identifier: "print"},
					lexer.SymbolToken{Symbol: "("},
					lexer.StrToken{Str: "bar"},
					lexer.SymbolToken{Symbol: ")"},
					lexer.SymbolToken{Symbol: ";"},
					lexer.IdentToken{Identifier: "print"},
					lexer.SymbolToken{Symbol: "("},
					lexer.StrToken{Str: "not foo"},
					lexer.SymbolToken{Symbol: ")"},
					lexer.SymbolToken{Symbol: ";"},
				},
				{
					lexer.ElseToken{},
				},
				{
					lexer.IdentToken{Identifier: "print"},
					lexer.SymbolToken{Symbol: "("},
					lexer.StrToken{Str: "not bar not foo"},
					lexer.SymbolToken{Symbol: ","},
					lexer.NumToken{Num: "5"},
					lexer.SymbolToken{Symbol: ","},
					lexer.StrToken{Str: "bazzz"},
					lexer.SymbolToken{Symbol: ")"},
					lexer.SymbolToken{Symbol: ";"},
					lexer.EndToken{},
				},
				{
					lexer.BlockToken{Block: "}}"},
					lexer.EOLToken{},
				},
			},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
				&NonTerminalParseNode{Value: "B"},
				&NonTerminalParseNode{Value: "C"},
				&NonTerminalParseNode{Value: "D"},
				&BlockParseNode{Value: "{{"},
				&NonTerminalParseNode{Value: "F"},
				&BlockParseNode{Value: "}}"},
				&IfStatementParseNode{},
				&IfConditionalParseNode{},
				&ElseIfConditionalParseNode{},
				&ElseParseNode{},
				&IdentParseNode{Value: "end"},
				&IdentParseNode{Value: "if"},
				&StringParseNode{Value: "("},
				&NonTerminalParseNode{Value: "K"},
				&StringParseNode{Value: ")"},
				&NonTerminalParseNode{Value: "V"},
				&ElseIfListParseNode{},
				&IdentParseNode{Value: "else"},
				&NonTerminalParseNode{Value: "V"},
				&NonTerminalParseNode{Value: "L"},
				&NonTerminalParseNode{Value: "G"},
				&ElseIfListParseNode{},
				&IdentParseNode{Value: "else_if"},
				&StringParseNode{Value: "("},
				&NonTerminalParseNode{Value: "K"},
				&StringParseNode{Value: ")"},
				&NonTerminalParseNode{Value: "V"},
				&NonTerminalParseNode{Value: "G"},
				&NonTerminalParseNode{Value: "M"},
				&NonTerminalParseNode{Value: "G"},
				&NonTerminalParseNode{Value: "F"},
				&StringParseNode{Value: ";"},
				&ElseIfListParseNode{},
				&IdentParseNode{Value: "else_if"},
				&StringParseNode{Value: "("},
				&NonTerminalParseNode{Value: "K"},
				&StringParseNode{Value: ")"},
				&NonTerminalParseNode{Value: "V"},
				&NonTerminalParseNode{Value: "L"},
				&NonTerminalParseNode{Value: "G"},
				&NonTerminalParseNode{Value: "F"},
				&StringParseNode{Value: ";"},
				&NonTerminalParseNode{Value: "U"},
				&NonTerminalParseNode{Value: "F"},
				&StringParseNode{Value: ";"},
				&FuncCallParseNode{},
				&IdentParseNode{Value: "else_if"},
				&StringParseNode{Value: "("},
				&NonTerminalParseNode{Value: "K"},
				&StringParseNode{Value: ")"},
				&NonTerminalParseNode{Value: "V"},
				&NonTerminalParseNode{Value: "L"},
				&NonTerminalParseNode{Value: "G"},
				&NonTerminalParseNode{Value: "M"},
				&NonTerminalParseNode{Value: "G"},
				&NonTerminalParseNode{Value: "F"},
				&StringParseNode{Value: ";"},
				&FuncCallParseNode{},
				&NonTerminalParseNode{Value: "U"},
				&StringParseNode{Value: "<"},
				&NonTerminalParseNode{Value: "N"},
				&FuncCallParseNode{},
				&IdentParseNode{Value: "print"},
				&StringParseNode{Value: "("},
				&ArgsParseNode{},
				&StringParseNode{Value: ")"},
				&NonTerminalParseNode{Value: "L"},
				&NonTerminalParseNode{Value: "G"},
				&NonTerminalParseNode{Value: "M"},
				&NonTerminalParseNode{Value: "G"},
				&NonTerminalParseNode{Value: "F"},
				&StringParseNode{Value: ";"},
				&NonTerminalParseNode{Value: "U"},
				&NonTerminalParseNode{Value: "F"},
				&StringParseNode{Value: ";"},
				&FuncCallParseNode{},
				&IdentParseNode{Value: "print"},
				&StringParseNode{Value: "("},
				&ArgsParseNode{},
				&StringParseNode{Value: ")"},
				&NonTerminalParseNode{Value: "N"},
				&NonTerminalParseNode{Value: "O"},
				&IdentParseNode{Value: "print"},
				&StringParseNode{Value: "("},
				&ArgsParseNode{},
				&StringParseNode{Value: ")"},
				&ArgsListParseNode{},
				&NonTerminalParseNode{Value: "M"},
				&NonTerminalParseNode{Value: "G"},
				&NonTerminalParseNode{Value: "F"},
				&StringParseNode{Value: ";"},
				&NonTerminalParseNode{Value: "U"},
				&NonTerminalParseNode{Value: "F"},
				&StringParseNode{Value: ";"},
				&FuncCallParseNode{},
				&NonTerminalParseNode{Value: "U"},
				&StringParseNode{Value: ">"},
				&NonTerminalParseNode{Value: "N"},
				&FuncCallParseNode{},
				&IdentParseNode{Value: "print"},
				&StringParseNode{Value: "("},
				&ArgsParseNode{},
				&StringParseNode{Value: ")"},
				&ArgsListParseNode{},
				&NonTerminalParseNode{Value: "O"},
				&NonTerminalParseNode{Value: "P"},
				&ArgsListParseNode{},
				&NonTerminalParseNode{Value: "K"},
				&NonTerminalParseNode{Value: "U"},
				&NonTerminalParseNode{Value: "F"},
				&StringParseNode{Value: ";"},
				&FuncCallParseNode{},
				&NonTerminalParseNode{Value: "U"},
				&StringParseNode{Value: ">"},
				&NonTerminalParseNode{Value: "N"},
				&FuncCallParseNode{},
				&IdentParseNode{Value: "print"},
				&StringParseNode{Value: "("},
				&ArgsParseNode{},
				&StringParseNode{Value: ")"},
				&NonTerminalParseNode{Value: "N"},
				&NonTerminalParseNode{Value: "O"},
				&IdentParseNode{Value: "print"},
				&StringParseNode{Value: "("},
				&ArgsParseNode{},
				&StringParseNode{Value: ")"},
				&ArgsListParseNode{},
				&NonTerminalParseNode{Value: "K"},
				&NonTerminalParseNode{Value: "P"},
				&IdentParseNode{Value: "bar"},
				&NonTerminalParseNode{Value: "K"},
				&NonTerminalParseNode{Value: "L"},
				&NonTerminalParseNode{Value: "U"},
				&StringParseNode{Value: ">"},
				&NonTerminalParseNode{Value: "N"},
				&FuncCallParseNode{},
				&IdentParseNode{Value: "print"},
				&StringParseNode{Value: "("},
				&ArgsParseNode{},
				&StringParseNode{Value: ")"},
				&NonTerminalParseNode{Value: "N"},
				&NonTerminalParseNode{Value: "O"},
				&IdentParseNode{Value: "print"},
				&StringParseNode{Value: "("},
				&ArgsParseNode{},
				&StringParseNode{Value: ")"},
				&ArgsListParseNode{},
				&NonTerminalParseNode{Value: "O"},
				&NonTerminalParseNode{Value: "P"},
				&ArgsListParseNode{},
				&NonTerminalParseNode{Value: "K"},
				&NonTerminalParseNode{Value: "L"},
				&IdentParseNode{Value: "foo"},
				&NonTerminalParseNode{Value: "L"},
				&NonTerminalParseNode{Value: "M"},
				&NonTerminalParseNode{Value: "N"},
				&NonTerminalParseNode{Value: "O"},
				&IdentParseNode{Value: "print"},
				&StringParseNode{Value: "("},
				&ArgsParseNode{},
				&StringParseNode{Value: ")"},
				&ArgsListParseNode{},
				&NonTerminalParseNode{Value: "O"},
				&NonTerminalParseNode{Value: "P"},
				&ArgsListParseNode{},
				&NonTerminalParseNode{Value: "K"},
				&NonTerminalParseNode{Value: "P"},
				&IdentParseNode{Value: "bar"},
				&NonTerminalParseNode{Value: "K"},
				&NonTerminalParseNode{Value: "L"},
				&NonTerminalParseNode{Value: "M"},
				&NonTerminalParseNode{Value: "M"},
				&NonTerminalParseNode{Value: "U"},
				&NonTerminalParseNode{Value: "O"},
				&NonTerminalParseNode{Value: "P"},
				&ArgsListParseNode{},
				&NonTerminalParseNode{Value: "K"},
				&NonTerminalParseNode{Value: "P"},
				&IdentParseNode{Value: "bar"},
				&NonTerminalParseNode{Value: "K"},
				&NonTerminalParseNode{Value: "L"},
				&IdentParseNode{Value: "foo"},
				&NonTerminalParseNode{Value: "L"},
				&NonTerminalParseNode{Value: "M"},
				&NonTerminalParseNode{Value: "U"},
				&NonTerminalParseNode{Value: "U"},
				&NonTerminalParseNode{Value: "N"},
				&NonTerminalParseNode{Value: "P"},
				&IdentParseNode{Value: "bar"},
				&NonTerminalParseNode{Value: "K"},
				&NonTerminalParseNode{Value: "L"},
				&IdentParseNode{Value: "foo"},
				&NonTerminalParseNode{Value: "L"},
				&NonTerminalParseNode{Value: "M"},
				&NonTerminalParseNode{Value: "M"},
				&NonTerminalParseNode{Value: "U"},
				&NonTerminalParseNode{Value: "N"},
				&NonTerminalParseNode{Value: "N"},
				&NonTerminalParseNode{Value: "O"},
				&IdentParseNode{Value: "foo"},
				&NonTerminalParseNode{Value: "L"},
				&NonTerminalParseNode{Value: "M"},
				&NonTerminalParseNode{Value: "M"},
				&NonTerminalParseNode{Value: "U"},
				&NonTerminalParseNode{Value: "U"},
				&NonTerminalParseNode{Value: "N"},
				&NonTerminalParseNode{Value: "O"},
				&NonTerminalParseNode{Value: "O"},
				&NonTerminalParseNode{Value: "P"},
				&NonTerminalParseNode{Value: "M"},
				&NonTerminalParseNode{Value: "U"},
				&NonTerminalParseNode{Value: "U"},
				&NonTerminalParseNode{Value: "N"},
				&NonTerminalParseNode{Value: "N"},
				&NonTerminalParseNode{Value: "O"},
				&NonTerminalParseNode{Value: "P"},
				&NonTerminalParseNode{Value: "P"},
				&StringParseNode{Value: "not bar"},
				&NonTerminalParseNode{Value: "U"},
				&NonTerminalParseNode{Value: "N"},
				&NonTerminalParseNode{Value: "N"},
				&NonTerminalParseNode{Value: "O"},
				&NonTerminalParseNode{Value: "O"},
				&NonTerminalParseNode{Value: "P"},
				&StringParseNode{Value: "not bar not foo"},
				&StringParseNode{Value: "foo"},
				&NonTerminalParseNode{Value: "N"},
				&NonTerminalParseNode{Value: "O"},
				&NonTerminalParseNode{Value: "O"},
				&NonTerminalParseNode{Value: "P"},
				&NonTerminalParseNode{Value: "P"},
				&StringParseNode{Value: "not foo"},
				&NonTerminalParseNode{Value: "O"},
				&NonTerminalParseNode{Value: "P"},
				&NonTerminalParseNode{Value: "P"},
				&StringParseNode{Value: "not foo"},
				&StringParseNode{Value: "bar"},
				&NonTerminalParseNode{Value: "P"},
				&StringParseNode{Value: "not foo"},
				&StringParseNode{Value: "bar"},
				&StringParseNode{Value: "bar"},
			},
		},
		{
			tokens: [][]lexer.Token{
				{
					lexer.BlockToken{Block: "{{"},
				},
				{
					lexer.ForToken{},
					lexer.SymbolToken{Symbol: "("},
					lexer.IdentToken{Identifier: "foo"},
					lexer.InToken{},
					lexer.StrToken{Str: "foo/bar"},
					lexer.SymbolToken{Symbol: ")"},
				},
				{
					lexer.IdentToken{Identifier: "print"},
					lexer.SymbolToken{Symbol: "("},
					lexer.StrToken{Str: "foo"},
					lexer.SymbolToken{Symbol: ")"},
					lexer.SymbolToken{Symbol: ";"},
					lexer.EndToken{},
				},
				{
					lexer.BlockToken{Block: "}}"},
					lexer.EOLToken{},
				},
			},
			nodes: []TreeNode{
				&NonTerminalParseNode{Value: ""},
				&NonTerminalParseNode{Value: "B"},
				&NonTerminalParseNode{Value: "C"},
				&NonTerminalParseNode{Value: "D"},
				&BlockParseNode{Value: "{{"},
				&NonTerminalParseNode{Value: "F"},
				&BlockParseNode{Value: "}}"},
				&ForLoopParseNode{},
				&IdentParseNode{Value: "for"},
				&StringParseNode{Value: "("},
				&IdentParseNode{Value: "foo"},
				&IdentParseNode{Value: "in"},
				&StringParseNode{Value: "foo/bar"},
				&StringParseNode{Value: ")"},
				&NonTerminalParseNode{Value: "V"},
				&IdentParseNode{Value: "end"},
				&NonTerminalParseNode{Value: "G"},
				&NonTerminalParseNode{Value: "F"},
				&StringParseNode{Value: ";"},
				&FuncCallParseNode{},
				&IdentParseNode{Value: "print"},
				&StringParseNode{Value: "("},
				&ArgsParseNode{},
				&StringParseNode{Value: ")"},
				&ArgsListParseNode{},
				&NonTerminalParseNode{Value: "K"},
				&NonTerminalParseNode{Value: "L"},
				&NonTerminalParseNode{Value: "M"},
				&NonTerminalParseNode{Value: "U"},
				&NonTerminalParseNode{Value: "N"},
				&NonTerminalParseNode{Value: "O"},
				&NonTerminalParseNode{Value: "P"},
				&StringParseNode{Value: "foo"},
			},
		},
		{
			tokens: [][]lexer.Token{
				{
					lexer.BlockToken{Block: "{{"},
				},
				{
					lexer.ForToken{},
					lexer.SymbolToken{Symbol: "("},
					lexer.IdentToken{Identifier: "foo"},
					lexer.InToken{},
					lexer.IdentToken{Identifier: "post"},
					lexer.SymbolToken{Symbol: "."},
					lexer.IdentToken{Identifier: "title"},
					lexer.SymbolToken{Symbol: ")"},
				},
				{
					lexer.IdentToken{Identifier: "print"},
					lexer.SymbolToken{Symbol: "("},
					lexer.StrToken{Str: "foo"},
					lexer.SymbolToken{Symbol: ")"},
					lexer.SymbolToken{Symbol: ";"},
					lexer.EndToken{},
				},
				{
					lexer.BlockToken{Block: "}}"},
					lexer.EOLToken{},
				},
			},
			nodes: []TreeNode{
				&NonTerminalParseNode{Value: ""},
				&NonTerminalParseNode{Value: "B"},
				&NonTerminalParseNode{Value: "C"},
				&NonTerminalParseNode{Value: "D"},
				&BlockParseNode{Value: "{{"},
				&NonTerminalParseNode{Value: "F"},
				&BlockParseNode{Value: "}}"},
				&ForLoopParseNode{},
				&IdentParseNode{Value: "for"},
				&StringParseNode{Value: "("},
				&IdentParseNode{Value: "foo"},
				&IdentParseNode{Value: "in"},
				&IdentParseNode{Value: "post.children"},
				&StringParseNode{Value: ")"},
				&NonTerminalParseNode{Value: "V"},
				&IdentParseNode{Value: "end"},
				&NonTerminalParseNode{Value: "G"},
				&NonTerminalParseNode{Value: "F"},
				&StringParseNode{Value: ";"},
				&FuncCallParseNode{},
				&IdentParseNode{Value: "print"},
				&StringParseNode{Value: "("},
				&ArgsParseNode{},
				&StringParseNode{Value: ")"},
				&ArgsListParseNode{},
				&NonTerminalParseNode{Value: "K"},
				&NonTerminalParseNode{Value: "L"},
				&NonTerminalParseNode{Value: "M"},
				&NonTerminalParseNode{Value: "U"},
				&NonTerminalParseNode{Value: "N"},
				&NonTerminalParseNode{Value: "O"},
				&NonTerminalParseNode{Value: "P"},
				&StringParseNode{Value: "foo"},
			},
		},
		{
			tokens: [][]lexer.Token{
				{
					lexer.BlockToken{Block: "{{"},
				},
				{
					lexer.ForToken{},
					lexer.SymbolToken{Symbol: "("},
					lexer.IdentToken{Identifier: "foo"},
					lexer.InToken{},
					lexer.IdentToken{Identifier: "bar"},
					lexer.SymbolToken{Symbol: "("},
					lexer.SymbolToken{Symbol: ")"},
					lexer.SymbolToken{Symbol: ")"},
				},
				{
					lexer.IdentToken{Identifier: "print"},
					lexer.SymbolToken{Symbol: "("},
					lexer.StrToken{Str: "foo"},
					lexer.SymbolToken{Symbol: ")"},
					lexer.SymbolToken{Symbol: ";"},
					lexer.EndToken{},
				},
				{
					lexer.BlockToken{Block: "}}"},
					lexer.EOLToken{},
				},
			},
			nodes: []TreeNode{
				&NonTerminalParseNode{Value: ""},
				&NonTerminalParseNode{Value: "B"},
				&NonTerminalParseNode{Value: "C"},
				&NonTerminalParseNode{Value: "D"},
				&BlockParseNode{Value: "{{"},
				&NonTerminalParseNode{Value: "F"},
				&BlockParseNode{Value: "}}"},
				&ForLoopParseNode{},
				&IdentParseNode{Value: "for"},
				&StringParseNode{Value: "("},
				&IdentParseNode{Value: "foo"},
				&IdentParseNode{Value: "in"},
				&FuncCallParseNode{},
				&StringParseNode{Value: ")"},
				&NonTerminalParseNode{Value: "V"},
				&IdentParseNode{Value: "end"},
				&IdentParseNode{Value: "bar"},
				&StringParseNode{Value: "("},
				&ArgsParseNode{},
				&StringParseNode{Value: ")"},
				&NonTerminalParseNode{Value: "G"},
				&NonTerminalParseNode{Value: "F"},
				&StringParseNode{Value: ";"},
				&FuncCallParseNode{},
				&IdentParseNode{Value: "print"},
				&StringParseNode{Value: "("},
				&ArgsParseNode{},
				&StringParseNode{Value: ")"},
				&ArgsListParseNode{},
				&NonTerminalParseNode{Value: "K"},
				&NonTerminalParseNode{Value: "L"},
				&NonTerminalParseNode{Value: "M"},
				&NonTerminalParseNode{Value: "U"},
				&NonTerminalParseNode{Value: "N"},
				&NonTerminalParseNode{Value: "O"},
				&NonTerminalParseNode{Value: "P"},
				&StringParseNode{Value: "foo"},
			},
		},
		{
			tokens: [][]lexer.Token{
				{
					lexer.BlockToken{Block: "{{"},
					lexer.IdentToken{Identifier: "post"},
					lexer.SymbolToken{Symbol: "."},
					lexer.IdentToken{Identifier: "title"},
					lexer.AssignOpToken{Operator: "="},
					lexer.NumToken{Num: "5"},
					lexer.SymbolToken{Symbol: ";"},
					lexer.IdentToken{Identifier: "bar"},
					lexer.SymbolToken{Symbol: "."},
					lexer.IdentToken{Identifier: "baz"},
					lexer.AssignOpToken{Operator: "="},
					lexer.NumToken{Num: "1"},
					lexer.SymbolToken{Symbol: ";"},
					lexer.IdentToken{Identifier: "bar"},
					lexer.SymbolToken{Symbol: "."},
					lexer.IdentToken{Identifier: "baz"},
					lexer.AssignOpToken{Operator: "="},
					lexer.NumToken{Num: "1"},
					lexer.SymbolToken{Symbol: ";"},
					lexer.BlockToken{Block: "}}"},
					lexer.EOLToken{},
				},
			},
			nodes: []TreeNode{
				&NonTerminalParseNode{},
				&NonTerminalParseNode{Value: "B"},
				&NonTerminalParseNode{Value: "C"},
				&NonTerminalParseNode{Value: "D"},
				&BlockParseNode{Value: "{{"},
				&NonTerminalParseNode{Value: "V"},
				&BlockParseNode{Value: "}}"},
				&NonTerminalParseNode{Value: "G"},
				&NonTerminalParseNode{Value: "G"},
				&NonTerminalParseNode{Value: "F"},
				&StringParseNode{Value: ";"},
				&NonTerminalParseNode{Value: "G"},
				&NonTerminalParseNode{Value: "F"},
				&StringParseNode{Value: ";"},
				&NonTerminalParseNode{Value: "K"},
				&NonTerminalParseNode{Value: "F"},
				&StringParseNode{Value: ";"},
				&NonTerminalParseNode{Value: "K"},
				&IdentParseNode{Value: "bar.baz"},
				&StringParseNode{Value: "="},
				&NonTerminalParseNode{Value: "K"},
				&NonTerminalParseNode{Value: "K"},
				&IdentParseNode{Value: "bar.baz"},
				&StringParseNode{Value: "="},
				&NonTerminalParseNode{Value: "K"},
				&NonTerminalParseNode{Value: "L"},
				&IdentParseNode{Value: "foo.title"},
				&StringParseNode{Value: "="},
				&NonTerminalParseNode{Value: "K"},
				&NonTerminalParseNode{Value: "L"},
				&NonTerminalParseNode{Value: "M"},
				&NonTerminalParseNode{Value: "L"},
				&NonTerminalParseNode{Value: "M"},
				&NonTerminalParseNode{Value: "U"},
				&NonTerminalParseNode{Value: "M"},
				&NonTerminalParseNode{Value: "U"},
				&NonTerminalParseNode{Value: "N"},
				&NonTerminalParseNode{Value: "U"},
				&NonTerminalParseNode{Value: "N"},
				&NonTerminalParseNode{Value: "O"},
				&NonTerminalParseNode{Value: "N"},
				&NonTerminalParseNode{Value: "O"},
				&NonTerminalParseNode{Value: "P"},
				&NonTerminalParseNode{Value: "O"},
				&NonTerminalParseNode{Value: "P"},
				&NumParseNode{Value: 1},
				&NonTerminalParseNode{Value: "P"},
				&NumParseNode{Value: 1},
				&NumParseNode{Value: 5},
			},
		},
	}

	for i, test := range tests {
		tokChan := getTokChan(test.tokens)
		nodeChan := make(chan TreeNode)
		errChan := make(chan error)

		go Parse(tokChan, nodeChan, errChan)

		nodes := []TreeNode{}

		for node := range nodeChan {
			nodes = append(nodes, node)
		}

		if err := <-errChan; err != nil {
			t.Errorf("Test %d failed.\nExpected no errors. Got %q.", i, err.Error())
		}

		flatChildren := flattenChildren(nodes[0])

		if equal, msg := nodeSlicesEqual(test.nodes, flatChildren); !equal {
			t.Errorf("Test %d failed.\n%s", i, msg)
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
	allChildren := []TreeNode{node}
	queue := node.GetChildren()

	for len(queue) > 0 {
		child := queue[0]
		queue = queue[1:]
		queue = append(queue, child.GetChildren()...)

		allChildren = append(allChildren, child)
	}

	return allChildren
}

func nodeSlicesEqual(expected, got []TreeNode) (bool, string) {
	if len(expected) != len(got) {
		return false, fmt.Sprintf("Expected %d nodes. Got %d.", len(expected), len(got))
	}

	for i := 0; i < len(expected); i++ {
		expectedNode := expected[i]
		gotNode := got[i]

		if reflect.TypeOf(expectedNode) != reflect.TypeOf(gotNode) {
			return false, fmt.Sprintf("Expected type %q. Got type %q.", reflect.TypeOf(expected[i]), reflect.TypeOf(got[i]))
		}

		switch node := expectedNode.(type) {
		case *NonTerminalParseNode:
			gn := gotNode.(*NonTerminalParseNode)
			if gn.Value != node.Value {
				return false, fmt.Sprintf("Expected value %q. Got value %q.", node.Value, gn.Value)
			}
		case *StringParseNode:
			gn := gotNode.(*StringParseNode)
			if gn.Value != node.Value {
				return false, fmt.Sprintf("Expected value %q. Got value %q.", node.Value, gn.Value)
			}
		case *NumParseNode:
			gn := gotNode.(*NumParseNode)
			if gn.Value != node.Value {
				return false, fmt.Sprintf("Expected number %d. Got number %d.", node.Value, gn.Value)
			}
		case *IdentParseNode:
			gn := gotNode.(*IdentParseNode)
			if gn.Value != node.Value {
				return false, fmt.Sprintf("Expected ident %q. Got ident %q.", node.Value, gn.Value)
			}
		}
	}

	return true, ""
}
