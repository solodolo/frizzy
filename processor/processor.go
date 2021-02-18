package processor

import (
	"fmt"
	"log"
	"strings"

	"mettlach.codes/frizzy/file"
	"mettlach.codes/frizzy/parser"
)

// Process reads each node from nodeChan and walks through its tree
// turning parse nodes into output
func Process(nodeChan <-chan parser.TreeNode, resultChan chan<- Result, context *Context) {
	defer close(resultChan)

	for node := range nodeChan {
		resultChan <- processHeadNode(node, context)
	}
}

func processHeadNode(head parser.TreeNode, context *Context) Result {
	// if head is an assignment
	//	add to context
	switch typedNode := head.(type) {
	case *parser.NonTerminalParseNode:
		children := head.GetChildren()

		if typedNode.IsAssignment() {
			left, right, _ := getBinaryOperatorAndOperands(children, context)
			varName := left.GetResult().(string)

			(*context)[varName] = ContextNode{result: right}
			return nil
		} else if typedNode.IsAddition() {
			addResult, err := processAddition(getBinaryOperatorAndOperands(children, context))

			if err != nil {
				panic(fmt.Sprintf("addition error: %s", err))
			}

			return addResult
		} else if typedNode.IsMultiplication() {
			multResult, err := processMultiplication(getBinaryOperatorAndOperands(children, context))

			if err != nil {
				panic(fmt.Sprintf("multiplication error: %s", err))
			}
			return multResult
		} else if typedNode.IsRelation() {
			relResult, err := processRel(getBinaryOperatorAndOperands(children, context))

			if err != nil {
				panic(fmt.Sprintf("rel error: %s", err))
			}
			return relResult
		} else if typedNode.IsLogic() {
			logicResult, err := processLogic(getBinaryOperatorAndOperands(children, context))

			if err != nil {
				panic(fmt.Sprintf("logic error: %s", err))
			}
			return logicResult
		} else if typedNode.IsUnary() {
			unaryResult, err := processUnary(getUnaryOperatorAndOperand(children, context))

			if err != nil {
				panic(fmt.Sprintf("unary error: %s", err))
			}

			return unaryResult
		} else {
			// TODO: This smells bad
			if len(children) == 1 {
				return processHeadNode(children[0], context)
			}
			ret := ""
			for _, child := range children {
				ret += processHeadNode(child, context).String()
			}
			return StringResult(ret)
		}
	case *parser.ForLoopParseNode:
		inputResult := processHeadNode(typedNode.GetLoopInput(), context)
		inputContexts := getLoopInputContexts(&inputResult)
		loopBody := typedNode.GetLoopBody()

		bodyText := ""

		// TODO: is this loop needed or is the body a single node?
		for _, inputContext := range inputContexts {
			inputContext.Merge(context)
			for i := range loopBody {
				bodyText += processHeadNode(loopBody[i], inputContext).String() + "\n"
			}
		}

		return StringResult(bodyText)
	case *parser.IfStatementParseNode:
		ifCondition := processHeadNode(typedNode.GetIfConditional(), context).(BoolResult)
		// check if first
		if bool(ifCondition) {
			ifBody := processHeadNode(typedNode.GetIfBody(), context)
			return StringResult(ifBody.String())
		}

		// check any else_ifs
		elseIfConditions := typedNode.GetElseIfConditionals()
		for i, elseCondition := range elseIfConditions {
			condition := processHeadNode(elseCondition, context).(BoolResult)
			if bool(condition) {
				if elseIfBody, ok := typedNode.GetElseIfBody(i); ok {
					body := processHeadNode(elseIfBody, context)
					return StringResult(body.String())
				}
			}
		}

		// finally try for the else
		if elseBody, ok := typedNode.GetElseBody(); ok {
			body := processHeadNode(elseBody, context)
			return StringResult(body.String())
		}

		// nothing is true
		return StringResult("")
	case *parser.FuncCallParseNode:
		funcName := typedNode.GetFuncName()
		args := typedNode.GetArgs()
		processedArgs := []Result{}

		for _, arg := range args {
			processedArgs = append(processedArgs, processHeadNode(arg, context))
		}

		switch funcName {
		case "print":
			return Print(processedArgs[0])
		case "Paginate":
			return Paginate(processedArgs[0], processedArgs[1], processedArgs[2])
		default:
			log.Fatalf("call to undefined function %q", funcName)
		}

	case *parser.StringParseNode:
		return StringResult(typedNode.Value)
	case *parser.NumParseNode:
		return IntResult(typedNode.Value)
	case *parser.BoolParseNode:
		return BoolResult(typedNode.Value)
	case *parser.VarParseNode:
		contextKey := typedNode.Value
		keys := strings.Split(contextKey, ".")

		// loop through each key and recursively look up the next level
		current := ContextNode{child: context}
		for _, key := range keys {
			contextNode, exists := current.At(key)

			if !exists {
				log.Fatalf("key %q not found in context", contextKey)
			}

			current = contextNode
		}

		// current is the last context node so we can
		// return its result or its further nested context
		if current.HasResult() {
			return current.result
		}

		return ContainerResult{current.child}
	default:
		return StringResult("")
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

// Returns the operator and operands of the binary operation represented in ops
// e.g. given 5 + 4, ops = []parser.TreeNode{5, '+', 4}
func getBinaryOperatorAndOperands(ops []parser.TreeNode, context *Context) (Result, Result, string) {
	left := processHeadNode(ops[0], context)
	operator := processHeadNode(ops[1], context).(StringResult)
	right := processHeadNode(ops[len(ops)-1], context)

	return left, right, string(operator)
}

// Returns the operator and operand of the unary operation in ops
// e.g. given !false, ops = []parser.TreeNode{"!", false}
func getUnaryOperatorAndOperand(ops []parser.TreeNode, context *Context) (Result, string) {
	operator := processHeadNode(ops[0], context).(StringResult)
	right := processHeadNode(ops[len(ops)-1], context)

	return right, string(operator)
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
	if operator == "&&" {
		return leftOp.LogicalAnd(right)
	} else if operator == "||" {
		return leftOp.LogicalOr(right)
	}

	return nil, fmt.Errorf("Invalid logcal operator %q", operator)
}

// TODO: Figure out a better way to look up the appropriate functions
func processRel(left, right Result, operator string) (Result, error) {
	if operator == "==" {
		leftOp := left.(EqualityResult)
		return leftOp.EqualTo(right)
	} else if operator == "!=" {
		leftOp := left.(EqualityResult)
		return leftOp.NotEqualTo(right)
	} else {
		leftOp := left.(RelResult)
		if operator == "<" {
			return leftOp.LessThan(right)
		} else if operator == ">" {
			return leftOp.GreaterThan(right)
		} else if operator == "<=" {
			return leftOp.LessThanEqual(right)
		} else if operator == ">=" {
			return leftOp.GreaterThanEqual(right)
		}
	}

	return nil, fmt.Errorf("Invalid relation operator %q", operator)
}

func processUnary(right Result, operator string) (Result, error) {
	if operator == "!" {
		rightOp := right.(NotResult)
		return rightOp.Not()
	} else if operator == "-" {
		rightOp := right.(NegativeResult)
		return rightOp.Negative()
	}

	return nil, fmt.Errorf("Invalid unary operator %q", operator)
}

// returns the contexts of each file in contextPath dir
func getLoopContentContexts(contextPath string) []*Context {
	store := GetExportStore()
	// typedInput is a path to content
	contentPaths := file.GetContentPaths(contextPath)
	ret := make([]*Context, len(contentPaths))

	// iterate through dir in content dir
	for i, path := range contentPaths {
		ret[i] = store.Get(path)
	}
	return ret
	// return array of export store context for each file
}

// getLoopInputContexts returns an array of contexts that should
// be sent on each iteration of a for loop
func getLoopInputContexts(input *Result) []*Context {
	switch typedInput := (*input).(type) {
	case *StringResult:
		// return a context for each file in the path
		return getLoopContentContexts(string(*typedInput))
	case *ContainerResult:
		// return a context
		return []*Context{typedInput.context}
	default:
		return nil
	}
}
