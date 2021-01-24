package parser

import (
	"fmt"
	"strconv"

	"mettlach.codes/frizzy/lexer"
)

// Parse reads tokens from tokChan and parses them into tree nodes
// which are then fed to nodeChan
func Parse(tokChan chan []lexer.Token, nodeChan chan TreeNode, errChan chan error) {
	defer close(nodeChan)
	defer close(errChan)

	stateStack := &[]int{}
	nodeStack := &[]TreeNode{}

	for tokens := range tokChan {
		// When open block is found
		// 	start in table state 0
		// 	parse until an accept is reached
		// Go until open block is found
		i := 0
		for i < len(tokens) {
			token := tokens[i]
			// Passthrough tokens can just be sent on
			if ptToken, ok := token.(lexer.PassthroughToken); ok {
				node := StringParseNode{Value: ptToken.GetValue()}
				nodeChan <- node
				i++
			} else {
				j, err := parseTokens(tokens[i:], stateStack, nodeStack, nodeChan)

				if err != nil {
					errChan <- err
				}

				i += j
			}
		}
	}

	errChan <- nil
}

func parseTokens(tokens []lexer.Token, stateStack *[]int, nodeStack *[]TreeNode, nodeChan chan TreeNode) (i int, err error) {
	if len(*stateStack) == 0 {
		*stateStack = append(*stateStack, 0)
	}

	for i < len(tokens) {
		token := tokens[i]
		currentState := (*stateStack)[len(*stateStack)-1]
		currentSymbol := token.GetGrammarSymbol()
		col, exists := SymbolColMapping[currentSymbol]

		if !exists {
			err = fmt.Errorf("parse table column not found for symbol %q", currentSymbol)
			break
		}

		action := LR1ParseTable[currentState][col]

		if action == "" {
			err = fmt.Errorf("invalid action taken for state %d, symbol %q", currentState, currentSymbol)
			break
		} else if IsShiftAction(action) {
			err = handleShiftAction(action, token, stateStack, nodeStack)

			if err != nil {
				break
			}

			// go to next token
			i++
		} else if IsReduceAction(action) {
			err = handleReduceAction(action, stateStack, nodeStack)

			if err != nil {
				break
			}
		} else {
			// Accept is the only remaining option
			head := NonTerminalParseNode{Value: "A"}
			head.children = *nodeStack

			// Send head to channel
			go func() {
				nodeChan <- head
			}()

			// Clear stacks
			*stateStack = []int{}
			*nodeStack = []TreeNode{}
			i++
			break
		}
	}

	return
}

// Shifts the next state onto the state stack and creates a tree node
// for the symbol
func handleShiftAction(action string, token lexer.Token, stateStack *[]int, nodeStack *[]TreeNode) (err error) {
	nextState, e := strconv.Atoi(action[1:])

	if e != nil {
		err = fmt.Errorf("could not convert %q to a valid state", action[1:])
		return
	}

	*stateStack = append(*stateStack, nextState)

	// Add tree node to stack
	node := getNodeForToken(token)
	*nodeStack = append(*nodeStack, node)

	return
}

// Looks up reduction rule A -> B and shifts |B| symbols off each stack.
// Then pushes GOTO[stateStack.back(), A] onto state stack and A onto node stack
func handleReduceAction(action string, stateStack *[]int, nodeStack *[]TreeNode) (err error) {
	ruleNum, e := strconv.Atoi(action[1:])

	if e != nil {
		err = fmt.Errorf("could not convert %q to a valid grammar rule", action[1:])
		return
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
		return
	}

	// Push goto state on top of stack
	*stateStack = append(*stateStack, gotoState)

	// Push nonterminal onto stack
	node := &NonTerminalParseNode{Value: left}
	node.children = poppedNodes
	*nodeStack = append(*nodeStack, node)

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
