package processor

import (
	"fmt"
	"log"
	"strings"

	"mettlach.codes/frizzy/file"
	"mettlach.codes/frizzy/parser"
)

// NodeProcessor holds a context and an optional
// channel to receive export assignments during
// this processing run
type NodeProcessor struct {
	Context        *Context
	PathReader     file.GetPathFunc
	ExportStore    ExportStorage
	FunctionModule FunctionModule
}

// Process reads each node from nodeChan and walks through its tree
// turning parse nodes into output
func (receiver *NodeProcessor) Process(nodeChan <-chan parser.TreeNode, resultChan chan<- Result) {
	defer close(resultChan)

	for node := range nodeChan {
		resultChan <- receiver.processHeadNode(node)
	}
}

func (receiver *NodeProcessor) processHeadNode(head parser.TreeNode) Result {
	// if head is an assignment
	//	add to context
	switch typedNode := head.(type) {
	case *parser.NonTerminalParseNode:
		children := head.GetChildren()

		if typedNode.IsAssignment() {
			if keys, right, err := receiver.getAssignmentKeysAndValue(children); err == nil {
				receiver.insertInContext(keys, right)
				receiver.doExport(keys, right)
			} else {
				log.Fatal(fmt.Sprintf("invalid assignment: %s", err))
			}

			return StringResult("")
		} else if typedNode.IsAddition() {
			addResult, err := processAddition(receiver.getBinaryOperatorAndOperands(children))

			if err != nil {
				panic(fmt.Sprintf("addition error: %s", err))
			}

			return addResult
		} else if typedNode.IsMultiplication() {
			multResult, err := processMultiplication(receiver.getBinaryOperatorAndOperands(children))

			if err != nil {
				panic(fmt.Sprintf("multiplication error: %s", err))
			}
			return multResult
		} else if typedNode.IsRelation() {
			relResult, err := processRel(receiver.getBinaryOperatorAndOperands(children))

			if err != nil {
				panic(fmt.Sprintf("rel error: %s", err))
			}
			return relResult
		} else if typedNode.IsLogic() {
			logicResult, err := processLogic(receiver.getBinaryOperatorAndOperands(children))

			if err != nil {
				panic(fmt.Sprintf("logic error: %s", err))
			}
			return logicResult
		} else if typedNode.IsUnary() {
			unaryResult, err := processUnary(receiver.getUnaryOperatorAndOperand(children))

			if err != nil {
				panic(fmt.Sprintf("unary error: %s", err))
			}

			return unaryResult
		} else {
			// TODO: This smells bad
			if len(children) == 1 {
				return receiver.processHeadNode(children[0])
			}
			ret := ""
			for _, child := range children {
				ret += receiver.processHeadNode(child).String()
			}
			return StringResult(ret)
		}
	case *parser.ContentParseNode:
		children := typedNode.GetChildren()
		resultText := make([]string, 0, len(children))

		for _, child := range children {
			resultText = append(resultText, receiver.processHeadNode(child).String())
		}

		return StringResult(strings.Join(resultText, ""))
	case *parser.ForLoopParseNode:
		inputResult := receiver.processHeadNode(typedNode.GetLoopInput())
		inputContexts := receiver.getLoopInputContexts(&inputResult)
		loopBody := typedNode.GetLoopBody()
		loopIdent := typedNode.GetLoopIdent().(*parser.IdentParseNode)

		return receiver.generateLoopBody(loopBody, loopIdent, inputContexts)
	case *parser.IfStatementParseNode:
		ifCondition := receiver.processHeadNode(typedNode.GetIfConditional()).(BoolResult)
		// check if first
		if bool(ifCondition) {
			ifBody := receiver.processHeadNode(typedNode.GetIfBody())
			return StringResult(ifBody.String())
		}

		// check any else_ifs
		elseIfConditions := typedNode.GetElseIfConditionals()
		for i, elseCondition := range elseIfConditions {
			condition := receiver.processHeadNode(elseCondition).(BoolResult)
			if bool(condition) {
				if elseIfBody, ok := typedNode.GetElseIfBody(i); ok {
					body := receiver.processHeadNode(elseIfBody)
					return StringResult(body.String())
				}
			}
		}

		// finally try for the else
		if elseBody, ok := typedNode.GetElseBody(); ok {
			body := receiver.processHeadNode(elseBody)
			return StringResult(body.String())
		}

		// nothing is true
		return StringResult("")
	case *parser.FuncCallParseNode:
		funcName := typedNode.GetFuncName()
		args := typedNode.GetArgs()
		processedArgs := []Result{}

		for _, arg := range args {
			processedArgs = append(processedArgs, receiver.processHeadNode(arg))
		}

		result, ok := receiver.callFunction(funcName, processedArgs)

		if ok {
			return result
		}

		// TODO: Replace with error
		return nil

	case *parser.StringParseNode:
		return StringResult(typedNode.Value)
	case *parser.NumParseNode:
		return IntResult(typedNode.Value)
	case *parser.BoolParseNode:
		return BoolResult(typedNode.Value)
	case *parser.VarNameParseNode:
		keys := typedNode.GetVarNameParts()
		if node, ok := receiver.lookupInContext(keys); ok {

			// current is the last context node so we can
			// return its result or its further nested context
			if node.HasResult() {
				return node.result
			}

			return ContainerResult{node.child}
		}

		log.Fatalf("context lookup failed: %q", strings.Join(keys, "."))
		return nil
	case *parser.BlockParseNode:
		content := typedNode.GetContent()
		parsed := receiver.processHeadNode(content)

		result := ""
		if typedNode.IsPrintable() {
			result = parsed.String()
		}

		return StringResult(result)
	default:
		return StringResult("")
	}
}

// Returns the left side of the assignment as a string and the right as a processed Result
func (receiver *NodeProcessor) getAssignmentKeysAndValue(ops []parser.TreeNode) ([]string, Result, error) {
	if left, ok := ops[0].(*parser.VarNameParseNode); ok {
		nameParts := left.GetVarNameParts()
		right := receiver.processHeadNode(ops[len(ops)-1])

		return nameParts, right, nil
	}

	return nil, nil, fmt.Errorf("invalid assignment to %T", ops[0])
}

// Returns the operator and operands of the binary operation represented in ops
// e.g. given 5 + 4, ops = []parser.TreeNode{5, '+', 4}
func (receiver *NodeProcessor) getBinaryOperatorAndOperands(ops []parser.TreeNode) (Result, Result, string) {
	left := receiver.processHeadNode(ops[0])
	operator := receiver.processHeadNode(ops[1]).(StringResult)
	right := receiver.processHeadNode(ops[len(ops)-1])

	return left, right, string(operator)
}

// Returns the operator and operand of the unary operation in ops
// e.g. given !false, ops = []parser.TreeNode{"!", false}
func (receiver *NodeProcessor) getUnaryOperatorAndOperand(ops []parser.TreeNode) (Result, string) {
	operator := receiver.processHeadNode(ops[0]).(StringResult)
	right := receiver.processHeadNode(ops[len(ops)-1])

	return right, string(operator)
}

func (receiver *NodeProcessor) doExport(keys []string, value Result) {
	if receiver.ExportStore != nil {
		receiver.ExportStore.Insert(keys, value)
	}
}

func (receiver *NodeProcessor) doGetContext(filePath string) *Context {
	if receiver.ExportStore != nil {
		return receiver.ExportStore.GetFileContext(filePath)
	}

	return &Context{}
}

func (receiver *NodeProcessor) callFunction(funcName string, args []Result) (Result, bool) {
	if receiver.FunctionModule == nil {
		module := NewBuiltinFunctionModule()
		receiver.FunctionModule = &module
	}

	return receiver.FunctionModule.CallFunction(funcName, args...)
}

func (receiver *NodeProcessor) generateLoopBody(body parser.TreeNode, loopIdent *parser.IdentParseNode, contexts []*Context) StringResult {
	bodyText := ""

	context := receiver.Context
	merged := &Context{}
	for _, inputContext := range contexts {
		(*merged)[loopIdent.Value] = &ContextNode{child: context.Merge(inputContext)}
		loopProcessor := &NodeProcessor{Context: merged}
		bodyText += loopProcessor.processHeadNode(body).String()
	}

	return StringResult(bodyText)
}

// getPaths reads paths using the provided PathReader or defaulting to
// file.GetContentPaths if PathReader is nil
func (receiver *NodeProcessor) getPaths(subpath string) []string {
	pathReader := receiver.PathReader
	if pathReader == nil {
		pathReader = file.GetContentPaths
	}
	return pathReader(subpath)
}

func (receiver *NodeProcessor) lookupInContext(keys []string) (*ContextNode, bool) {
	return receiver.Context.AtNested(keys)
}

func (receiver *NodeProcessor) insertInContext(keys []string, value Result) {
	receiver.Context.Insert(keys, value)
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
func (receiver *NodeProcessor) getLoopContentContexts(contextPath string) []*Context {
	// typedInput is a path to content
	contentPaths := receiver.getPaths(contextPath)
	ret := make([]*Context, len(contentPaths))

	// iterate through dir in content dir
	for i, path := range contentPaths {
		ret[i] = receiver.doGetContext(path)
	}

	// return array of export store context for each file
	return ret
}

// getLoopInputContexts returns an array of contexts that should
// be sent on each iteration of a for loop
func (receiver *NodeProcessor) getLoopInputContexts(input *Result) []*Context {
	switch typedInput := (*input).(type) {
	case StringResult:
		// return a context for each file in the path
		return receiver.getLoopContentContexts(string(typedInput))
	case ContainerResult:
		// return a context
		return typedInput.context.Values()
	default:
		return nil
	}
}
