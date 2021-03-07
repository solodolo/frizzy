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
		"Y -> ID . Y",
		"Y -> ID",
		"H -> ID ( I )",
		"I -> J",
		"I -> ε",
		"J -> J , K",
		"J -> K",
		// "K -> VAR_NAME = K",
		"K -> Y = K",
		"K -> M",
		"M -> M LOGIC_OP U",
		"M -> U",
		"U -> U REL_OP N",
		"U -> N",
		"N -> N + O",
		"N -> N - O",
		"N -> O",
		"O -> O MULT_OP L",
		"O -> L",
		"L -> ! L",
		"L -> - L",
		"L -> P",
		// "P -> VAR_NAME",
		"P -> STRING",
		"P -> NUM",
		"P -> BOOL",
		// "P -> ID",
		"P -> Y",
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
		// "R -> FOR ( ID IN VAR_NAME ) V END",
		"R -> FOR ( ID IN Y ) V END",
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
		".":        6,
		"!":        7,
		"LOGIC_OP": 8,
		"+":        9,
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
		"-":        24,
		"$":        25,
		"B":        26,
		"C":        27,
		"D":        28,
		"E":        29,
		"F":        30,
		"G":        31,
		"H":        32,
		"I":        33,
		"J":        34,
		"K":        35,
		"L":        36,
		"M":        37,
		"N":        38,
		"O":        39,
		"P":        40,
		"Q":        41,
		"R":        42,
		"S":        43,
		"T":        44,
		"U":        45,
		"V":        46,
		"W":        47,
		"X":        48,
		"Y":        49,
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
