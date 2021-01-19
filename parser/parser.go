package parser

import (
	"fmt"

	"mettlach.codes/frizzy/lexer"
)

func Parse(tokChan chan []lexer.Token) {
	// Start with state zero
	stateStack := []int{0}

	for tokens := range tokChan {
		// Go until open block is found
		// When open block is found
		// 	start in table state 0
		// 	parse until an accept is reached
		parseTokens(tokens)
	}
}

func parseTokens(tokens []lexer.Token) {
	fmt.Println(tokens)
}
