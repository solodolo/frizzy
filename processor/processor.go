package processor

import (
	"context"
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
	CurPage        int
	NumPages       int
}

func NewNodeProcessor(
	filepath string,
	context *Context,
	pathReader file.GetPathFunc,
	exportStore ExportStorage,
	funcModule FunctionModule,
	curPage int,
	numPages int,
) *NodeProcessor {
	processor := &NodeProcessor{
		Context:        context,
		PathReader:     pathReader,
		ExportStore:    exportStore,
		FunctionModule: funcModule,
		CurPage:        curPage,
		NumPages:       numPages,
	}

	if context == nil {
		processor.Context = &Context{}
	}

	processor.Context.Insert([]string{"curPage"}, IntResult(curPage))
	processor.Context.Insert([]string{"numPages"}, IntResult(numPages))

	if pathReader == nil {
		processor.PathReader = file.GetContentPaths
	}

	if exportStore == nil {
		processor.ExportStore = NewExportFileStore(filepath)
	}

	if funcModule == nil {
		processor.FunctionModule = NewBuiltinFunctionModule()
	}

	return processor
}

// Process reads each node from nodeChan and walks through its tree
// turning parse nodes into output
func (receiver *NodeProcessor) Process(nodeChan <-chan parser.TreeNode, ctx context.Context) (<-chan Result, <-chan error) {
	resultChan := make(chan Result)
	errChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errChan)

		for node := range nodeChan {
			result, err := receiver.processHeadNode(node)
			if err != nil {
				errChan <- err
				return
			}

			select {
			case resultChan <- result:
			case <-ctx.Done():
				return
			}
		}
	}()

	return resultChan, errChan
}

func (receiver *NodeProcessor) processHeadNode(head parser.TreeNode) (Result, error) {
	var processResult Result
	var processError error

	switch typedNode := head.(type) {
	case *parser.NonTerminalParseNode:
		children := head.GetChildren()

		if typedNode.IsAssignment() {
			if keys, right, err := receiver.getAssignmentKeysAndValue(children); err == nil {
				receiver.insertInContext(keys, right)
				receiver.doExport(keys, right)
				processResult = StringResult("")
			} else {
				processError = fmt.Errorf("invalid assignment: %s", err)
			}
		} else if typedNode.IsAddition() {
			processResult, processError = processAddition(receiver.getBinaryOperatorAndOperands(children))
		} else if typedNode.IsMultiplication() {
			processResult, processError = processMultiplication(receiver.getBinaryOperatorAndOperands(children))
		} else if typedNode.IsRelation() {
			processResult, processError = processRel(receiver.getBinaryOperatorAndOperands(children))
		} else if typedNode.IsLogic() {
			processResult, processError = processLogic(receiver.getBinaryOperatorAndOperands(children))
		} else if typedNode.IsUnary() {
			processResult, processError = processUnary(receiver.getUnaryOperatorAndOperand(children))
		} else {
			// TODO: This smells bad
			if len(children) == 1 {
				processResult, processError = receiver.processHeadNode(children[0])
			} else {
				ret := ""
				for _, child := range children {
					tmpRet, err := receiver.processHeadNode(child)
					if err != nil {
						processError = err
						break
					}
					ret += tmpRet.String()
				}

				processResult = StringResult(ret)
			}
		}
	case *parser.ContentParseNode:
		children := typedNode.GetChildren()
		resultText := make([]string, 0, len(children))

		for _, child := range children {
			childResult, err := receiver.processHeadNode(child)
			if err != nil {
				processError = err
				break
			}
			resultText = append(resultText, childResult.String())
		}

		processResult = StringResult(strings.Join(resultText, ""))
	case *parser.ForLoopParseNode:
		inputResult, err := receiver.processHeadNode(typedNode.GetLoopInput())

		if err != nil {
			processError = err
		} else {
			inputContexts := receiver.getLoopInputContexts(&inputResult)
			loopBody := typedNode.GetLoopBody()
			loopIdent := typedNode.GetLoopIdent().(*parser.IdentParseNode)

			processResult = receiver.generateLoopBody(loopBody, loopIdent, inputContexts)
		}
	case *parser.IfStatementParseNode:
		ifResult, err := receiver.processHeadNode(typedNode.GetIfConditional())

		if err != nil {
			processError = err
			break
		}

		ifCondition := ifResult.(BoolResult)
		// check if first
		if bool(ifCondition) {
			ifBody, err := receiver.processHeadNode(typedNode.GetIfBody())
			if err != nil {
				processError = err
			} else {
				processResult = StringResult(ifBody.String())
			}
			break
		}

		// check any else_ifs
		elseIfConditions := typedNode.GetElseIfConditionals()
		for i, elseCondition := range elseIfConditions {
			elseIfResult, err := receiver.processHeadNode(elseCondition)
			if err != nil {
				processError = err
				break
			}

			condition := elseIfResult.(BoolResult)
			if bool(condition) {
				if elseIfBody, ok := typedNode.GetElseIfBody(i); ok {
					body, err := receiver.processHeadNode(elseIfBody)

					if err != nil {
						processError = err
					} else {
						processResult = StringResult(body.String())
					}
					break
				}
			}
		}

		if processError != nil || processResult != nil {
			break
		}

		// finally try for the else
		if elseBody, ok := typedNode.GetElseBody(); ok {
			body, err := receiver.processHeadNode(elseBody)

			if err != nil {
				processError = err
			} else {
				processResult = StringResult(body.String())
			}
			break
		}

		// nothing is true
		processResult = StringResult("")
	case *parser.FuncCallParseNode:
		funcName := typedNode.GetFuncName()

		processedArgs := []Result{}
		args := typedNode.GetArgs()

		if funcName == "paginate" || funcName == "pagesBefore" || funcName == "pagesAfter" {
			processedArgs = append(processedArgs, IntResult(receiver.CurPage))
		}

		if funcName == "pagesAfter" {
			processedArgs = append(processedArgs, IntResult(receiver.NumPages))
		}

		if funcName == "pagesBefore" || funcName == "pagesAfter" {
			inputPath := receiver.ExportStore.GetNamespace()
			processedArgs = append(processedArgs, StringResult(inputPath))
		}

		for _, arg := range args {
			argResult, err := receiver.processHeadNode(arg)

			if err != nil {
				processError = err
				break
			}

			processedArgs = append(processedArgs, argResult)
		}

		if processError == nil {
			if result, err := receiver.callFunction(funcName, processedArgs); err == nil {
				processResult = result
			} else {
				processError = err
			}
		}

	case *parser.StringParseNode:
		processResult = StringResult(typedNode.Value)
	case *parser.NumParseNode:
		processResult = IntResult(typedNode.Value)
	case *parser.BoolParseNode:
		processResult = BoolResult(typedNode.Value)
	case *parser.VarNameParseNode:
		keys := typedNode.GetVarNameParts()
		if node, ok := receiver.lookupInContext(keys); ok {

			// current is the last context node so we can
			// return its result or its further nested context
			if node.HasResult() {
				processResult = node.result
			} else {
				processResult = ContainerResult{node.child}
			}
		} else {
			log.Printf("context lookup failed: %q", strings.Join(keys, "."))
			processResult = StringResult("")
		}
	case *parser.BlockParseNode:
		content := typedNode.GetContent()
		parsed, err := receiver.processHeadNode(content)

		if err != nil {
			processError = err
		} else {
			result := ""
			if typedNode.IsPrintable() {
				result = parsed.String()
			}

			processResult = StringResult(result)
		}
	default:
		processResult = StringResult("")
	}

	return processResult, processError
}

// Returns the left side of the assignment as a string and the right as a processed Result
func (receiver *NodeProcessor) getAssignmentKeysAndValue(ops []parser.TreeNode) ([]string, Result, error) {
	if left, ok := ops[0].(*parser.VarNameParseNode); ok {
		nameParts := left.GetVarNameParts()
		right, err := receiver.processHeadNode(ops[len(ops)-1])

		return nameParts, right, err
	}

	return nil, nil, fmt.Errorf("invalid assignment to %T", ops[0])
}

// Returns the operator and operands of the binary operation represented in ops
// e.g. given 5 + 4, ops = []parser.TreeNode{5, '+', 4}
func (receiver *NodeProcessor) getBinaryOperatorAndOperands(ops []parser.TreeNode) (Result, Result, string) {
	left, _ := receiver.processHeadNode(ops[0])
	operatorResult, _ := receiver.processHeadNode(ops[1])
	operator := operatorResult.(StringResult)
	right, _ := receiver.processHeadNode(ops[len(ops)-1])

	return left, right, string(operator)
}

// Returns the operator and operand of the unary operation in ops
// e.g. given !false, ops = []parser.TreeNode{"!", false}
func (receiver *NodeProcessor) getUnaryOperatorAndOperand(ops []parser.TreeNode) (Result, string) {
	operatorResult, _ := receiver.processHeadNode(ops[0])
	operator := operatorResult.(StringResult)
	right, _ := receiver.processHeadNode(ops[len(ops)-1])

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

func (receiver *NodeProcessor) callFunction(funcName string, args []Result) (Result, error) {
	return receiver.FunctionModule.CallFunction(funcName, args...)
}

func (receiver *NodeProcessor) generateLoopBody(body parser.TreeNode, loopIdent *parser.IdentParseNode, contexts []*Context) StringResult {
	bodyText := ""

	context := receiver.Context
	merged := &Context{}
	namespace := receiver.ExportStore.GetNamespace()

	for _, inputContext := range contexts {
		(*merged)[loopIdent.Value] = &ContextNode{child: context.Merge(inputContext)}
		loopProcessor := NewNodeProcessor(namespace, merged, nil, nil, nil, 0, 0)
		bodyResult, _ := loopProcessor.processHeadNode(body)

		bodyText += bodyResult.String()
	}

	return StringResult(bodyText)
}

// getPaths is responsible for looking up the path of each file in a loop
// or any other situation when multiple files are referenced from the current
// file being processed
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
