package lexer

import (
	"fmt"
	"testing"
)

func TestGetLineTokensReturnsCorrectTokenType(t *testing.T) {
	var tests = []struct {
		symbol  string
		tokType string
	}{
		{"*", "lexer.OpToken"},
		{"/", "lexer.OpToken"},
		{"%", "lexer.OpToken"},
		{"+", "lexer.OpToken"},
		{"-", "lexer.OpToken"},
		{"!", "lexer.OpToken"},
		{"==", "lexer.OpToken"},
		{"!=", "lexer.OpToken"},
		{"<=", "lexer.OpToken"},
		{">=", "lexer.OpToken"},
		{"=", "lexer.OpToken"},
		{"<", "lexer.OpToken"},
		{">", "lexer.OpToken"},
		{"123", "lexer.NumToken"},
		{`"foobar"`, "lexer.StrToken"},
		{"for", "lexer.IdentToken"},
		{"post.title", "lexer.VarToken"},
		{"true", "lexer.BoolToken"},
		{"false", "lexer.BoolToken"},
		{";", "lexer.SymbolToken"},
		{"(", "lexer.SymbolToken"},
		{")", "lexer.SymbolToken"},
		{" ", "lexer.PassthroughToken"},
		{"{{", "lexer.BlockToken"},
	}

	for _, test := range tests {
		tokens := getLineTokens(test.symbol)

		if tokType := fmt.Sprintf("%T", tokens[0]); tokType != test.tokType {
			t.Errorf("Expected '%s' to return type %s. Got type %s.", test.symbol, test.tokType, tokType)
		}
	}
}
