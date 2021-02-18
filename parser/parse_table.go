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
		"A -> B",
		"B -> C",
		"C -> D",
		"C -> E",
		"D -> {{ F }}",
		"D -> {{ V }}",
		"E -> {{: F }}",
		"H -> ID ( I )",
		"I -> J",
		"I -> ε",
		"J -> J , K",
		"J -> K",
		"K -> VAR_NAME = K",
		"K -> M",
		// "K -> L",
		// "L -> M",
		"M -> M LOGIC_OP U",
		"M -> U",
		"U -> U REL_OP N",
		"U -> N",
		"N -> N ADD_OP O",
		"N -> O",
		"O -> O MULT_OP L",
		"O -> L",
		"L -> UNARY_OP L",
		"L -> P",
		// "O -> O MULT_OP P",
		// "O -> P",
		"P -> VAR_NAME",
		"P -> STRING",
		"P -> NUM",
		"P -> BOOL",
		"P -> ID",
		"P -> ( K )",
		"Q -> W S T end",
		"W -> IF ( K ) V",
		"S -> X",
		"S -> ε",
		"X -> X ELSE_IF ( K ) V",
		"X -> ELSE_IF ( K ) V",
		"T -> ELSE V",
		"T -> ε",
		"R -> FOR ( ID IN STRING ) V END",
		"R -> FOR ( ID IN VAR_NAME ) V END",
		"R -> FOR ( ID IN H ) V END",
		"F -> K",
		"F -> H",
		"F -> Q",
		"F -> R",
		"V -> G",
		"V -> ε",
		"G -> G F ;",
		"G -> F ;",
	}

	rawRows := strings.Split(rawTable, "\n")
	rawRows = rawRows[1:] // discard header

	LR1ParseTable = make([][]string, len(rawRows))

	for i, rawRow := range rawRows {
		LR1ParseTable[i] = strings.Split(rawRow, ", ")
	}

	SymbolColMapping = map[string]int{
		"{{":       0,
		"}}":       1,
		"{{:":      2,
		"ID":       3,
		"(":        4,
		")":        5,
		"VAR_NAME": 6,
		"UNARY_OP": 7,
		"LOGIC_OP": 8,
		"ADD_OP":   9,
		"MULT_OP":  10,
		"STRING":   11,
		"NUM":      12,
		"IF":       13,
		"ELSE_IF":  14,
		"ELSE":     15,
		"END":      16,
		"FOR":      17,
		"IN":       18,
		";":        19,
		",":        20,
		"=":        21,
		"REL_OP":   22,
		"BOOL":     23,
		"$":        24,
		"B":        25,
		"C":        26,
		"D":        27,
		"E":        28,
		"F":        29,
		"G":        30,
		"H":        31,
		"I":        32,
		"J":        33,
		"K":        34,
		"L":        35,
		"M":        36,
		"N":        37,
		"O":        38,
		"P":        39,
		"Q":        40,
		"R":        41,
		"S":        42,
		"T":        43,
		"U":        44,
		"V":        45,
		"W":        46,
		"X":        47,
	}
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
