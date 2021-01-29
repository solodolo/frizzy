package parser

import (
	"fmt"
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

	if equal, msg := nodeSlicesEqual(expected, nodes); !equal {
		t.Error(msg)
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
					lexer.ElseToken{},
				},
				{
					lexer.IdentToken{Identifier: "print"},
					lexer.SymbolToken{Symbol: "("},
					lexer.StrToken{Str: "not bar not foo"},
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
		// {
		// 	tokens: [][]lexer.Token{{
		// 		lexer.BlockToken{Block: "{{"},
		// 		lexer.BlockToken{Block: "}}"},
		// 		lexer.EOLToken{},
		// 	}},
		// 	nodes: []TreeNode{
		// 		NonTerminalParseNode{},
		// 		NonTerminalParseNode{Value: "B"},
		// 		NonTerminalParseNode{Value: "C"},
		// 		NonTerminalParseNode{Value: "D"},
		// 		StringParseNode{Value: "{{"},
		// 		NonTerminalParseNode{Value: "G"},
		// 		StringParseNode{Value: "}}"},
		// 	},
		// },
		{
			tokens: [][]lexer.Token{{
				lexer.BlockToken{Block: "{{:"},
				lexer.StrToken{Str: "foobar"},
				lexer.BlockToken{Block: "}}"},
				lexer.EOLToken{},
			}},
			nodes: []TreeNode{
				NonTerminalParseNode{},
				NonTerminalParseNode{Value: "B"},
				NonTerminalParseNode{Value: "C"},
				NonTerminalParseNode{Value: "E"},
				StringParseNode{Value: "{{:"},
				NonTerminalParseNode{Value: "F"},
				StringParseNode{Value: "}}"},
				NonTerminalParseNode{Value: "K"},
				NonTerminalParseNode{Value: "L"},
				NonTerminalParseNode{Value: "M"},
				NonTerminalParseNode{Value: "U"},
				NonTerminalParseNode{Value: "N"},
				NonTerminalParseNode{Value: "O"},
				NonTerminalParseNode{Value: "P"},
				StringParseNode{Value: "foobar"},
			},
		},
		// {
		// 	tokens: [][]lexer.Token{{
		// 		lexer.BlockToken{Block: "{{"},
		// 		lexer.NumToken{Num: "5"},
		// 		lexer.AddOpToken{Operator: "+"},
		// 		lexer.NumToken{Num: "2"},
		// 		lexer.BlockToken{Block: "}}"},
		// 		lexer.EOLToken{},
		// 	}},
		// 	nodes: []TreeNode{
		// 		NonTerminalParseNode{},
		// 		NonTerminalParseNode{Value: "B"},
		// 		NonTerminalParseNode{Value: "C"},
		// 		NonTerminalParseNode{Value: "D"},
		// 		StringParseNode{Value: "{{"},
		// 		NonTerminalParseNode{Value: "F"},
		// 		StringParseNode{Value: "}}"},
		// 		NonTerminalParseNode{Value: "K"},
		// 		NonTerminalParseNode{Value: "L"},
		// 		NonTerminalParseNode{Value: "M"},
		// 		NonTerminalParseNode{Value: "U"},
		// 		NonTerminalParseNode{Value: "N"},
		// 		NonTerminalParseNode{Value: "N"},
		// 		StringParseNode{Value: "+"},
		// 		NonTerminalParseNode{Value: "O"},
		// 		NonTerminalParseNode{Value: "O"},
		// 		NonTerminalParseNode{Value: "P"},
		// 		NonTerminalParseNode{Value: "P"},
		// 		NumParseNode{Value: 5},
		// 		NumParseNode{Value: 2},
		// 	},
		// },
		// {
		// 	tokens: [][]lexer.Token{{
		// 		lexer.BlockToken{Block: "{{"},
		// 		lexer.NumToken{Num: "5"},
		// 		lexer.AddOpToken{Operator: "+"},
		// 		lexer.NumToken{Num: "2"},
		// 		lexer.BlockToken{Block: "}}"},
		// 		lexer.EOLToken{},
		// 	}},
		// 	nodes: []TreeNode{
		// 		NonTerminalParseNode{},
		// 		NonTerminalParseNode{Value: "B"},
		// 		NonTerminalParseNode{Value: "C"},
		// 		NonTerminalParseNode{Value: "D"},
		// 		StringParseNode{Value: "{{"},
		// 		NonTerminalParseNode{Value: "F"},
		// 		StringParseNode{Value: "}}"},
		// 		NonTerminalParseNode{Value: "K"},
		// 		NonTerminalParseNode{Value: "L"},
		// 		NonTerminalParseNode{Value: "M"},
		// 		NonTerminalParseNode{Value: "U"},
		// 		NonTerminalParseNode{Value: "N"},
		// 		NonTerminalParseNode{Value: "N"},
		// 		StringParseNode{Value: "+"},
		// 		NonTerminalParseNode{Value: "O"},
		// 		NonTerminalParseNode{Value: "O"},
		// 		NonTerminalParseNode{Value: "P"},
		// 		NonTerminalParseNode{Value: "P"},
		// 		NumParseNode{Value: 5},
		// 		NumParseNode{Value: 2},
		// 	},
		// },
		// {
		// 	tokens: [][]lexer.Token{{
		// 		lexer.BlockToken{Block: "{{"},
		// 		lexer.NumToken{Num: "5"},
		// 		lexer.AddOpToken{Operator: "-"},
		// 		lexer.NumToken{Num: "2"},
		// 		lexer.BlockToken{Block: "}}"},
		// 		lexer.EOLToken{},
		// 	}},
		// 	nodes: []TreeNode{
		// 		NonTerminalParseNode{},
		// 		NonTerminalParseNode{Value: "B"},
		// 		NonTerminalParseNode{Value: "C"},
		// 		NonTerminalParseNode{Value: "D"},
		// 		StringParseNode{Value: "{{"},
		// 		NonTerminalParseNode{Value: "F"},
		// 		StringParseNode{Value: "}}"},
		// 		NonTerminalParseNode{Value: "K"},
		// 		NonTerminalParseNode{Value: "L"},
		// 		NonTerminalParseNode{Value: "M"},
		// 		NonTerminalParseNode{Value: "U"},
		// 		NonTerminalParseNode{Value: "N"},
		// 		NonTerminalParseNode{Value: "N"},
		// 		StringParseNode{Value: "-"},
		// 		NonTerminalParseNode{Value: "O"},
		// 		NonTerminalParseNode{Value: "O"},
		// 		NonTerminalParseNode{Value: "P"},
		// 		NonTerminalParseNode{Value: "P"},
		// 		NumParseNode{Value: 5},
		// 		NumParseNode{Value: 2},
		// 	},
		// },
		// {
		// 	tokens: [][]lexer.Token{
		// 		{
		// 			lexer.BlockToken{Block: "{{"},
		// 			lexer.NumToken{Num: "5"},
		// 		},
		// 		{
		// 			lexer.AddOpToken{Operator: "+"},
		// 			lexer.NumToken{Num: "2"},
		// 		},
		// 		{
		// 			lexer.BlockToken{Block: "}}"},
		// 			lexer.EOLToken{},
		// 		},
		// 	},
		// 	nodes: []TreeNode{
		// 		NonTerminalParseNode{},
		// 		NonTerminalParseNode{Value: "B"},
		// 		NonTerminalParseNode{Value: "C"},
		// 		NonTerminalParseNode{Value: "D"},
		// 		StringParseNode{Value: "{{"},
		// 		NonTerminalParseNode{Value: "F"},
		// 		StringParseNode{Value: "}}"},
		// 		NonTerminalParseNode{Value: "K"},
		// 		NonTerminalParseNode{Value: "L"},
		// 		NonTerminalParseNode{Value: "M"},
		// 		NonTerminalParseNode{Value: "U"},
		// 		NonTerminalParseNode{Value: "N"},
		// 		NonTerminalParseNode{Value: "N"},
		// 		StringParseNode{Value: "+"},
		// 		NonTerminalParseNode{Value: "O"},
		// 		NonTerminalParseNode{Value: "O"},
		// 		NonTerminalParseNode{Value: "P"},
		// 		NonTerminalParseNode{Value: "P"},
		// 		NumParseNode{Value: 5},
		// 		NumParseNode{Value: 2},
		// 	},
		// },
		// {
		// 	tokens: [][]lexer.Token{
		// 		{
		// 			lexer.BlockToken{Block: "{{"},
		// 		},
		// 		{
		// 			lexer.BlockToken{Block: "}}"},
		// 			lexer.EOLToken{},
		// 		},
		// 	},
		// 	nodes: []TreeNode{
		// 		NonTerminalParseNode{},
		// 		NonTerminalParseNode{Value: "B"},
		// 		NonTerminalParseNode{Value: "C"},
		// 		NonTerminalParseNode{Value: "D"},
		// 		StringParseNode{Value: "{{"},
		// 		NonTerminalParseNode{Value: "G"},
		// 		StringParseNode{Value: "}}"},
		// 	},
		// },
		// {
		// 	tokens: [][]lexer.Token{
		// 		{
		// 			lexer.BlockToken{Block: "{{"},
		// 		},
		// 		{
		// 			lexer.IfToken{},
		// 			lexer.SymbolToken{Symbol: "("},
		// 			lexer.IdentToken{Identifier: "foo"},
		// 			lexer.RelOpToken{Operator: "<"},
		// 			lexer.IdentToken{Identifier: "bar"},
		// 			lexer.SymbolToken{Symbol: ")"},
		// 		},
		// 		{
		// 			lexer.IdentToken{Identifier: "print"},
		// 			lexer.SymbolToken{Symbol: "("},
		// 			lexer.StrToken{Str: "foo"},
		// 			lexer.SymbolToken{Symbol: ")"},
		// 			lexer.SymbolToken{Symbol: ";"},
		// 			lexer.IdentToken{Identifier: "print"},
		// 			lexer.SymbolToken{Symbol: "("},
		// 			lexer.StrToken{Str: "not bar"},
		// 			lexer.SymbolToken{Symbol: ")"},
		// 			lexer.SymbolToken{Symbol: ";"},
		// 		},
		// 		{
		// 			lexer.ElseIfToken{},
		// 			lexer.SymbolToken{Symbol: "("},
		// 			lexer.IdentToken{Identifier: "foo"},
		// 			lexer.RelOpToken{Operator: ">"},
		// 			lexer.IdentToken{Identifier: "bar"},
		// 			lexer.SymbolToken{Symbol: ")"},
		// 		},
		// 		{
		// 			lexer.IdentToken{Identifier: "print"},
		// 			lexer.SymbolToken{Symbol: "("},
		// 			lexer.StrToken{Str: "bar"},
		// 			lexer.SymbolToken{Symbol: ")"},
		// 			lexer.SymbolToken{Symbol: ";"},
		// 			lexer.IdentToken{Identifier: "print"},
		// 			lexer.SymbolToken{Symbol: "("},
		// 			lexer.StrToken{Str: "not foo"},
		// 			lexer.SymbolToken{Symbol: ")"},
		// 			lexer.SymbolToken{Symbol: ";"},
		// 		},
		// 		{
		// 			lexer.ElseToken{},
		// 		},
		// 		{
		// 			lexer.IdentToken{Identifier: "print"},
		// 			lexer.SymbolToken{Symbol: "("},
		// 			lexer.StrToken{Str: "not bar not foo"},
		// 			lexer.SymbolToken{Symbol: ")"},
		// 			lexer.SymbolToken{Symbol: ";"},
		// 			lexer.EndToken{},
		// 		},
		// 		{
		// 			lexer.BlockToken{Block: "}}"},
		// 			lexer.EOLToken{},
		// 		},
		// 	},
		// 	nodes: []TreeNode{
		// 		NonTerminalParseNode{},
		// 		NonTerminalParseNode{Value: "B"},
		// 		NonTerminalParseNode{Value: "C"},
		// 		NonTerminalParseNode{Value: "D"},
		// 		StringParseNode{Value: "{{"},
		// 		NonTerminalParseNode{Value: "G"},
		// 		StringParseNode{Value: "}}"},
		// 	},
		// },
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

		printTreeStr(nodes[0])

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

func printTreeStr(node TreeNode) {
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
		case NonTerminalParseNode:
			gn := gotNode.(NonTerminalParseNode)
			if gn.Value != node.Value {
				return false, fmt.Sprintf("Expected value %q. Got value %q.", node.Value, gn.Value)
			}
		case StringParseNode:
			gn := gotNode.(StringParseNode)
			if gn.Value != node.Value {
				return false, fmt.Sprintf("Expected value %q. Got value %q.", node.Value, gn.Value)
			}
		case NumParseNode:
			gn := gotNode.(NumParseNode)
			if gn.Value != node.Value {
				return false, fmt.Sprintf("Expected number %d. Got number %d.", node.Value, gn.Value)
			}
		case IdentParseNode:
			gn := gotNode.(IdentParseNode)
			if gn.Value != node.Value {
				return false, fmt.Sprintf("Expected ident %q. Got ident %q.", node.Value, gn.Value)
			}
		}
	}

	return true, ""
}
