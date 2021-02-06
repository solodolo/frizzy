package parser

import (
	"fmt"
)

func Process(nodeChan <-chan TreeNode) {
	var globalContext Context
	for node := range nodeChan {
		processHeadNode(node, globalContext)
	}
}

func processHeadNode(head TreeNode, context Context) Result {
	// if head is an assignment
	//	add to context
	switch typedNode := head.(type) {
	case NonTerminalParseNode:
		children := head.GetChildren()
		left := processHeadNode(children[0], context)
		operator := processHeadNode(children[1], context).(StringResult)
		right := processHeadNode(children[len(children)-1], context)

		if typedNode.IsAssignment() {
			varName := left.GetResult().(string)

			context[varName] = left
			return nil
		} else if typedNode.IsAddition() {
			addResult, err := processAddition(left, right, string(operator))

			if err != nil {
				panic(fmt.Sprintf("addition error: %s", err))
			}

			return addResult
		} else if typedNode.IsMultiplication() {
			multResult, err := processMultiplication(left, right, string(operator))

			if err != nil {
				panic(fmt.Sprintf("multiplication error: %s", err))
			}
			return multResult
		} else if typedNode.IsLogic() {
			logicResult, err := processLogic(left, right, string(operator))

			if err != nil {
				panic(fmt.Sprintf("logic error: %s", err))
			}
			return logicResult
		}
	case StringParseNode:
		return StringResult(typedNode.Value)
	case NumParseNode:
		return IntResult(typedNode.Value)
	case BoolParseNode:
		return BoolResult(typedNode.Value)
	case VarParseNode:
		contextKey := typedNode.Value
		contextVal, exists := context[contextKey]

		if !exists {
			panic(fmt.Sprintf("context key %s not found", contextKey))
		}

		return contextVal
	}
	//
	// send context down recursively to each child
	//
	// if head doesn't have any children (or if a NumNode, StrNode ?)
	//	if head is an int value
	//		return IntResult
	//	else if head is a string value
	//		return StringResult
	//	else if head is a bool value
	//		return BoolResult
	// else if binary operator node
	//	perform operation with left and right Result
	// else if unary operator node
	//	perform operation with left Result
	// else if func call node
	//	left node is func name and the rest are arguments
	// else if an else if node
	//	if statement is true then return body value
	//	else return nil
	// else if an else node
	//	return body value
	// else if an if node
	//	if statement is true return if body
	//	else return first child body that isn't nil
	// else if a for loop node (want to add to context for children so 1. determine number of loops 2. determine context for each loop 3. traverse children after statement with each context)
	// else print block or non-print block
	//	what do we do here?
	//
	return nil
}

func processAddition(left, right Result, operator string) (Result, error) {
	if operator == "+" {
		leftOp := left.(AddableResult)
		return leftOp.Add(right)
	} else if operator == "-" {
		leftOp := left.(SubtractableResult)
		return leftOp.Subtract(right)
	}

	return nil, fmt.Errorf("Invalid addition operator %q", operator)
}

func processMultiplication(left, right Result, operator string) (Result, error) {
	leftOp := left.(MultipliableResult)
	if operator == "*" {
		return leftOp.Multiply(right)
	} else if operator == "/" {
		return leftOp.Divide(right)
	}

	return nil, fmt.Errorf("Invalid multiplication operator %q", operator)
}

func processLogic(left, right Result, operator string) (Result, error) {
	leftOp := left.(LogicalResult)
	if operator == "<" {
		return leftOp.LessThan(right)
	} else if operator == ">" {
		return leftOp.GreaterThan(right)
	} else if operator == "==" {
		return leftOp.EqualTo(right)
	} else if operator == "<=" {
		return leftOp.LessThanEqual(right)
	} else if operator == ">=" {
		return leftOp.GreaterThanEqual(right)
	}

	return nil, fmt.Errorf("Invalid logic operator %q", operator)
}
