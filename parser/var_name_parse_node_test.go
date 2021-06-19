package parser

import (
	"strings"
	"testing"

	"mettlach.codes/frizzy/lexer"
)

func TestNestedVarNameReturnsCorrectNameParts(t *testing.T) {
	var tests = []struct {
		toks      []lexer.Token
		nameParts []string
	}{
		{[]lexer.Token{
			lexer.BlockToken{Block: "{{"},
			lexer.IdentToken{Identifier: "foo"},
			lexer.BlockToken{Block: "}}"},
			lexer.EOLToken{},
		}, []string{"foo"}},
		{[]lexer.Token{
			lexer.BlockToken{Block: "{{"},
			lexer.IdentToken{Identifier: "foo"},
			lexer.SymbolToken{Symbol: "."},
			lexer.IdentToken{Identifier: "bar"},
			lexer.BlockToken{Block: "}}"},
			lexer.EOLToken{},
		}, []string{"foo", "bar"}},
		{[]lexer.Token{
			lexer.BlockToken{Block: "{{"},
			lexer.IdentToken{Identifier: "milky_way"},
			lexer.SymbolToken{Symbol: "."},
			lexer.IdentToken{Identifier: "sol_system"},
			lexer.SymbolToken{Symbol: "."},
			lexer.IdentToken{Identifier: "earth"},
			lexer.SymbolToken{Symbol: "."},
			lexer.IdentToken{Identifier: "africa"},
			lexer.SymbolToken{Symbol: "."},
			lexer.IdentToken{Identifier: "malawi"},
			lexer.BlockToken{Block: "}}"},
			lexer.EOLToken{},
		}, []string{"milky_way", "sol_system", "earth", "africa", "malawi"}},
	}

	for _, test := range tests {
		t.Run(strings.Join(test.nameParts, "."), func(t *testing.T) {
			stateStack := []int{}
			nodeStack := []TreeNode{}

			_, head, _ := parseTokens(test.toks, &stateStack, &nodeStack)
			varName := extractToken(head, []int{1}).(*VarNameParseNode)
			got := varName.GetVarNameParts()

			if strings.Join(got, "") != strings.Join(test.nameParts, "") {
				t.Errorf("expected %q, got %q", got, test.nameParts)
			}
		})
	}
}
