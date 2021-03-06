package parser

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"mettlach.codes/frizzy/lexer"
)

// Parse reads tokens from tokChan and parses them into tree nodes
// which are then fed to nodeChan
func Parse(tokChan <-chan []lexer.Token, ctx context.Context) (<-chan TreeNode, <-chan error) {
	nodeChan := make(chan TreeNode)
	errChan := make(chan error, 1)

	go readAndParseTokens(ctx, tokChan, nodeChan, errChan)
	return nodeChan, errChan
}

// Read lines of tokens from tokChan and turn them into TreeNodes sent to nodeChan
func readAndParseTokens(ctx context.Context, tokChan <-chan []lexer.Token, nodeChan chan TreeNode, parseErrChan chan error) {
	defer close(nodeChan)
	defer close(parseErrChan)

	stateStack := &[]int{}
	nodeStack := &[]TreeNode{}

	for tokens := range tokChan {
		// Track how many tokens have been read
		i := 0
		for i < len(tokens) {
			j, head, err := parseTokens(tokens[i:], stateStack, nodeStack)

			if err != nil {
				parseErrChan <- err
				return
			}

			if head != nil {
				select {
				case nodeChan <- head:
				case <-ctx.Done():
					return
				}
			}

			i += j
		}
	}
}

func parseTokens(tokens []lexer.Token, stateStack *[]int, nodeStack *[]TreeNode) (i int, head TreeNode, err error) {
	if len(*stateStack) == 0 {
		*stateStack = append(*stateStack, 0)
	}

	for i < len(tokens) {
		token := tokens[i]
		currentState := (*stateStack)[len(*stateStack)-1]
		currentSymbol := token.GetGrammarSymbol()
		col, exists := SymbolColMapping[currentSymbol]

		if !exists {
			err = getParseError(token, "unrecognized symbol %q", currentSymbol)
			break
		}

		action := LR1ParseTable[currentState][col]

		if action == "" {
			err = getParseError(token, "unexpected symbol %q", currentSymbol)
			break
		} else if IsShiftAction(action) {
			err = handleShiftAction(action, token, stateStack, nodeStack)

			if err != nil {
				break
			}

			// go to next token
			i++
		} else if IsReduceAction(action) {
			err = handleReduceAction(action, token, stateStack, nodeStack)

			if err != nil {
				break
			}
		} else {
			// Accept is the only remaining option
			head = &NonTerminalParseNode{}
			head.SetChildren(*nodeStack)

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
		err = getParseError(token, "could not convert %q to a valid state", action[1:])
		return
	}

	*stateStack = append(*stateStack, nextState)

	// Add tree node to stack
	node := getTerminalNodeForToken(token)
	*nodeStack = append(*nodeStack, node)

	return
}

// Looks up reduction rule A -> B and shifts |B| symbols off each stack.
// Then pushes GOTO[stateStack.back(), A] onto state stack and A onto node stack
func handleReduceAction(action string, token lexer.Token, stateStack *[]int, nodeStack *[]TreeNode) (err error) {
	ruleNum, e := strconv.Atoi(action[1:])

	if e != nil {
		err = getParseError(token, "could not convert %q to a valid grammar rule", action[1:])
		return
	}

	left, right := GetProductionParts(ruleNum)
	// Number of symbols to pop
	numToPop := len(right)
	*stateStack = (*stateStack)[:len(*stateStack)-numToPop]

	// Get the row and column to lookup in goto table
	lookupState := (*stateStack)[len(*stateStack)-1]
	lookupCol, exists := SymbolColMapping[left]

	if !exists {
		err = getParseError(token, "unrecognized symbol %q", left)
	}

	gotoState, e := strconv.Atoi(LR1ParseTable[lookupState][lookupCol])

	if e != nil {
		err = getParseError(token, "could not convert goto state %q to valid state", gotoState)
		return
	}

	// Push goto state on top of stack
	*stateStack = append(*stateStack, gotoState)

	isTerminal := false
	if numToPop == 1 {
		isTerminal = (*nodeStack)[len(*nodeStack)-numToPop].IsTerminal()
	}

	node := getNonTerminalNodeForReduction(left)

	// Collapse single non-terminal children
	if numToPop != 1 || isTerminal {
		// Stack symbols that will be popped become children of new node
		children := make([]TreeNode, numToPop)
		copy(children, (*nodeStack)[len(*nodeStack)-numToPop:])

		// If we have nested types e.g. nested ArgLists, then flatten the tree
		if len(children) > 0 && reflect.TypeOf(children[0]) == reflect.TypeOf(node) {
			grandChildren := children[0].GetChildren()
			// Drop the nested child and append remaining children to dropped child's children
			children = append(grandChildren, children[1:]...)
		}

		// Create non-terminal
		node.SetChildren(children)

		// Actually pop symbols
		*nodeStack = (*nodeStack)[:len(*nodeStack)-numToPop]

		// Append new symbol
		*nodeStack = append(*nodeStack, node)
	}
	return
}

// Creates the appropriate tree node for a given token
func getTerminalNodeForToken(token lexer.Token) TreeNode {
	var node TreeNode

	switch tok := token.(type) {
	case lexer.NumToken:
		num, _ := strconv.Atoi(tok.Num)
		node = &NumParseNode{Value: num}
	case lexer.BoolToken:
		truthy := tok.Value == "true"
		node = &BoolParseNode{Value: truthy}
	case lexer.IdentToken, lexer.ForToken, lexer.InToken, lexer.IfToken, lexer.ElseIfToken, lexer.ElseToken, lexer.EndToken:
		ident := tok.GetValue()
		node = &IdentParseNode{Value: ident}
	case lexer.SymbolToken:
		symbol := tok.Symbol
		node = &SymbolParseNode{Value: symbol}
	default:
		str := tok.GetValue()
		node = &StringParseNode{Value: str}
	}

	return node
}

func getNonTerminalNodeForReduction(reduction string) TreeNode {
	switch reduction {
	case "content":
		return &ContentParseNode{}
	case "func_call":
		return &FuncCallParseNode{}
	case "args":
		return &ArgsParseNode{}
	case "arg_list":
		return &ArgsListParseNode{}
	case "if_statement_block":
		return &IfStatementParseNode{}
	case "for_block":
		return &ForLoopParseNode{}
	case "else_if_list":
		return &ElseIfListParseNode{}
	case "var_name":
		return &VarNameParseNode{}
	case "block", "print_block":
		return &BlockParseNode{}
	default:
		return &NonTerminalParseNode{Value: reduction}
	}
}

// Returns an parse error formatted for the current token
func getParseError(token lexer.Token, msg string, msgFmt ...interface{}) error {
	lineNum, lineCol := token.GetLineNum(), token.GetLineCol()

	fmtedMsg := fmt.Sprintf(msg, msgFmt...)
	return fmt.Errorf("parse error line %d col %d: %s", lineNum, lineCol, fmtedMsg)
}
