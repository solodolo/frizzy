package parser

import (
	"fmt"
	"strconv"

	"mettlach.codes/frizzy/lexer"
)

func Parse(tokChan chan []lexer.Token, nodeChan chan TreeNode) {
	stateStack := &[]int{}
	nodeStack := &[]TreeNode{}

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
				err := parseTokens(tokSlice, stateStack, nodeStack, nodeChan)

				if err != nil {
					// handle error
				}
			}
		}
	}
}

func parseTokens(tokens []lexer.Token, stateStack *[]int, nodeStack *[]TreeNode, nodeChan chan TreeNode) (err error) {
	if len(*stateStack) == 0 {
		*stateStack = append(*stateStack, 0)
	}

	currentState := (*stateStack)[len(*stateStack)-1]

	for i := 0; i < len(tokens); {
		token := tokens[i]
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

			// Add tree node to stack
			node := getNodeForToken(token)
			*nodeStack = append(*nodeStack, node)

			// go to next token
			i++
		} else if IsReduceAction(action) {
			ruleNum, e := strconv.Atoi(action[1:])

			if e != nil {
				err = fmt.Errorf("could not convert %q to a valid grammar rule", action[1:])
				break
			}

			left, right := GetProductionParts(ruleNum)
			// Number of symbols to pop
			numToPop := len(right)
			*stateStack = (*stateStack)[:len(*stateStack)-numToPop]
			poppedNodes := (*nodeStack)[len(*nodeStack)-numToPop:]
			*nodeStack = (*nodeStack)[:len(*nodeStack)-numToPop]

			// Get the row and column to lookup in goto table
			lookupState := (*stateStack)[len(*stateStack)-1]
			lookupCol, exists := SymbolColMapping[left]

			if !exists {
				err = fmt.Errorf("parse table column not found for symbol %q", left)
			}

			gotoState, e := strconv.Atoi(LR1ParseTable[lookupState][lookupCol])

			if e != nil {
				err = fmt.Errorf("could not convert goto state %q to valid state", gotoState)
				break
			}

			// Push goto state on top of stack
			*stateStack = append(*stateStack, gotoState)

			// Push nonterminal onto stack
			node := &StringParseNode{Value: left}
			node.children = poppedNodes
			*nodeStack = append(*nodeStack, node)
		}
	}

	return
}

// Creates the appropriate tree node for a given token
func getNodeForToken(token lexer.Token) TreeNode {
	var node TreeNode

	switch tok := token.(type) {
	case lexer.NumToken:
		num, _ := strconv.Atoi(tok.Num)
		node = &NumParseNode{Value: num}
	case lexer.BoolToken:
		truthy := tok.Value == "true"
		node = &BoolParseNode{Value: truthy}
	case lexer.IdentToken:
		ident := tok.Identifier
		node = &IdentParseNode{Value: ident}
	case lexer.VarToken:
		varName := tok.Variable
		node = &VarParseNode{Value: varName}
	default:
		str := tok.GetValue()
		node = &StringParseNode{Value: str}
	}

	return node
}
