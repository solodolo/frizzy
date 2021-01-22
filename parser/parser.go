package parser

import (
	"fmt"
	"strconv"

	"mettlach.codes/frizzy/lexer"
)

func Parse(tokChan chan []lexer.Token, nodeChan chan ParentNode) {
	stateStack := &[]int{}

	for tokens := range tokChan {
		// When open block is found
		// 	start in table state 0
		// 	parse until an accept is reached
		// Go until open block is found
		for i, token := range tokens {
			// Passthrough tokens can just be sent on
			if ptToken, ok := token.(lexer.PassthroughToken); ok {
				node := &StringParseNode{Value: ptToken.GetValue()}
				nodeChan <- node
			} else {
				tokSlice := tokens[i:]
				err := parseTokens(tokSlice, stateStack, nodeChan)

				if err != nil {
					// handle error
				}
			}
		}
	}
}

func parseTokens(tokens []lexer.Token, stateStack *[]int, nodeChan chan ParentNode) (err error) {
	if len(*stateStack) == 0 {
		*stateStack = append(*stateStack, 0)
	}

	currentState := (*stateStack)[len(*stateStack)-1]

	for _, token := range tokens {
		currentSymbol := token.GetGrammarSymbol()
		col, exists := SymbolColMapping[currentSymbol]

		if !exists {
			err = fmt.Errorf("parse table column not found for symbol %q", currentSymbol)
			break
		}

		action := LR1ParseTable[currentState][col]

		if action == "" {
			err = fmt.Errorf("invalid action taken for %d, %q", currentState, currentSymbol)
			break
		} else if IsShiftAction(action) {
			nextState, e := strconv.Atoi(action[1:])

			if e != nil {
				err = fmt.Errorf("could not convert %q to a valid state", action[1:])
				break
			}
			*stateStack = append(*stateStack, nextState)
		}
	}

	return
}
