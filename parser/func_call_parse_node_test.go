package parser

import (
	"strings"
	"testing"

	"mettlach.codes/frizzy/lexer"
)

func TestArgsReturnsCorrectNameParts(t *testing.T) {
	var tests = []struct {
		toks      []lexer.Token
		nameParts []string
	}{
		{[]lexer.Token{
			lexer.BlockToken{Block: "{{"},
			lexer.IdentToken{Identifier: "foo"},
			lexer.SymbolToken{Symbol: "("},
			lexer.StrToken{Str: "a"},
			lexer.SymbolToken{Symbol: ")"},
			lexer.BlockToken{Block: "}}"},
			lexer.EOLToken{},
		}, []string{"a"}},
		{[]lexer.Token{
			lexer.BlockToken{Block: "{{"},
			lexer.IdentToken{Identifier: "foo"},
			lexer.SymbolToken{Symbol: "("},
			lexer.StrToken{Str: "a"},
			lexer.SymbolToken{Symbol: ","},
			lexer.StrToken{Str: "b"},
			lexer.SymbolToken{Symbol: ")"},
			lexer.BlockToken{Block: "}}"},
			lexer.EOLToken{},
		}, []string{"a", "b"}},
		{[]lexer.Token{
			lexer.BlockToken{Block: "{{"},
			lexer.IdentToken{Identifier: "foo"},
			lexer.SymbolToken{Symbol: "("},
			lexer.StrToken{Str: "a"},
			lexer.SymbolToken{Symbol: ","},
			lexer.StrToken{Str: "b"},
			lexer.SymbolToken{Symbol: ","},
			lexer.StrToken{Str: "c"},
			lexer.SymbolToken{Symbol: ","},
			lexer.StrToken{Str: "d"},
			lexer.SymbolToken{Symbol: ","},
			lexer.StrToken{Str: "e"},
			lexer.SymbolToken{Symbol: ")"},
			lexer.BlockToken{Block: "}}"},
			lexer.EOLToken{},
		}, []string{"a", "b", "c", "d", "e"}},
	}

	for _, test := range tests {
		t.Run(strings.Join(test.nameParts, ","), func(t *testing.T) {
			stateStack := []int{}
			nodeStack := []TreeNode{}

			_, head, _ := parseTokens(test.toks, &stateStack, &nodeStack)
			funcCall := extractToken(head, []int{1}).(*FuncCallParseNode)
			got := funcCall.GetArguments()
			gotStrs := make([]string, 0, len(got))

			for _, g := range got {
				node := g.GetChildren()[0].(*StringParseNode)
				gotStrs = append(gotStrs, node.Value)
			}

			if strings.Join(gotStrs, "") != strings.Join(test.nameParts, "") {
				t.Errorf("expected %q, got %q", test.nameParts, gotStrs)
			}
		})
	}
}
