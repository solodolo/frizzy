package parser

import (
	_ "embed"
	"strings"
)

var (
	// GrammarProductions is a list of the grammar productions
	GrammarProductions []string
	// LR1ParseTable is the action and goto LR1 table for the grammar
	LR1ParseTable [][]string
	// SymbolColMapping is a mapping between grammar symbols and their columns in the parse table
	SymbolColMapping map[string]int
)

//go:embed lr1_table.txt
var rawTable string

func init() {
	GrammarProductions = []string{
		"A -> program",
		"program -> content",

		"content -> content PASSTHROUGH",
		"content -> PASSTHROUGH",
		"content -> content blocks",
		"content -> blocks",

		"blocks -> block",
		"blocks -> print_block",
		"blocks -> if_statement_block",
		"blocks -> for_block",

		"block -> {{ statement }}",
		"block -> {{ statement -}",

		"print_block -> {{: statement }}",
		"print_block -> {{: statement -}",

		"if_statement_block -> {{if expression }} content END",
		"if_statement_block -> {{if expression }} content else_if_list END",
		"if_statement_block -> {{if expression }} content {{else}} content END",
		"if_statement_block -> {{if expression }} content else_if_list {{else}} content END",

		"else_if_list -> else_if_list {{else_if expression }} content",
		"else_if_list -> {{else_if expression }} content",

		"for_block -> {{for ID in STRING }} content END",
		"for_block -> {{for ID in var_name }} content END",
		"for_block -> {{for ID in func_call }} content END",

		"statement -> expression",
		"statement -> func_call",

		"var_name -> ID . var_name",
		"var_name -> ID",

		"func_call -> ID ( args )",

		"args -> arg_list",
		"args -> ε",

		"arg_list -> arg_list , expression",
		"arg_list -> expression",

		"expression -> var_name = expression",
		"expression -> logic_expression",

		"logic_expression -> logic_expression LOGIC_OP rel_expression",
		"logic_expression -> rel_expression",

		"rel_expression -> rel_expression REL_OP add_expression",
		"rel_expression -> add_expression",

		"add_expression -> add_expression + mult_expression",
		"add_expression -> add_expression - mult_expression",
		"add_expression -> mult_expression",

		"mult_expression -> mult_expression MULT_OP unary_expression",
		"mult_expression -> unary_expression",

		"unary_expression -> ! unary_expression",
		"unary_expression -> - unary_expression",
		"unary_expression -> term_expression",

		"term_expression -> STRING",
		"term_expression -> NUM",
		"term_expression -> BOOL",
		"term_expression -> var_name",
		"term_expression -> ( expression )",
	}

	rawRows := strings.Split(rawTable, "\n")
	SymbolColMapping = parseSymbolColMapping(rawRows[0])
	rawRows = rawRows[1:] // discard header

	LR1ParseTable = make([][]string, len(rawRows))

	for i, rawRow := range rawRows {
		LR1ParseTable[i] = strings.Split(rawRow, ", ")
	}
}

func parseSymbolColMapping(header string) map[string]int {
	symbolMap := map[string]int{}

	pieces := strings.Split(header, "','")
	for i, piece := range pieces {
		piece = strings.ReplaceAll(piece, "'", "")
		symbolMap[piece] = i
	}

	return symbolMap
}

// IsShiftAction determines if the given parse table action is a
// shift action
func IsShiftAction(action string) bool {
	return action[0] == 's'
}

// IsReduceAction determines if the given parse table action is a
// reduce action
func IsReduceAction(action string) bool {
	return action[0] == 'r'
}

// IsAcceptAction determines if the given parse table action is an
// accept action
func IsAcceptAction(action string) bool {
	return action == "acct"
}

// GetProductionParts returns the left and right pieces
// of grammar rule n.
// If rule n is S -> n b C, left = ["S"] and right = ["n", "b", "C"]
func GetProductionParts(n int) (left string, right []string) {
	rule := GrammarProductions[n]

	pieces := strings.Split(rule, " -> ")

	left = pieces[0]
	right = strings.Split(pieces[1], " ")

	if right[0] == "ε" {
		right = []string{}
	}

	return
}
