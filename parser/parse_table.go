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
		"A -> AA",
		"AA -> B",
		"AA -> ε",
		"B -> B PASSTHROUGH",
		"B -> B C",
		"B -> PASSTHROUGH",
		"B -> C",
		"C -> D",
		"C -> E",
		"C -> H",
		"C -> J",
		"D -> {{ K }}",
		"D -> {{ K -}",
		"BB -> E F G END",
		"E -> {{ W }} I",
		"F -> {{ X }} I F",
		"F -> ε",
		"G -> {{ Y }} I",
		"G -> ε",
		"H -> {{ Z }} I END",
		"I -> CC",
		// "I -> ε",
		// "CC -> CC PASSTHROUGH",
		"CC -> PASSTHROUGH",
		// "CC -> CC D",
		// "CC -> D",
		"J -> {{: K }}",
		"J -> {{: K -}",
		"K -> P",
		"K -> M",
		"L -> ID . L",
		"L -> ID",
		"M -> ID ( N )",
		"N -> O",
		"N -> ε",
		"O -> O , P",
		"O -> P",
		"P -> L = P",
		"P -> Q",
		"Q -> Q LOGIC_OP R",
		"Q -> R",
		"R -> R REL_OP S",
		"R -> S",
		"S -> S + T",
		"S -> S - T",
		"S -> T",
		"T -> T MULT_OP U",
		"T -> U",
		"U -> ! U",
		"U -> - U",
		"U -> V",
		"V -> STRING",
		"V -> NUM",
		"V -> BOOL",
		"V -> L",
		"V -> ( P )",
		"W -> if ( P )",
		"X -> else_if ( P )",
		"Y -> else ( P )",
		"Z -> for ( ID in STRING )",
		"Z -> for ( ID in L )",
		"Z -> for ( ID in M )",
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
