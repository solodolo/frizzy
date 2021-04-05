package processor

import (
	"math"
	"strconv"

	"mettlach.codes/frizzy/file"
	"mettlach.codes/frizzy/parser"
)

func PrintRaw(args ...Result) Result {
	return Print(args[0])
}

func PaginateRaw(args ...Result) Result {
	// TODO: replace nils with errors
	if len(args) < 3 {
		return nil
	}

	var filePathString, templatePathString string
	var numPerPageInt int

	if filePath, ok := args[0].(StringResult); ok {
		filePathString = string(filePath)
	} else {
		return nil
	}

	if templatePath, ok := args[1].(StringResult); ok {
		templatePathString = string(templatePath)
	} else {
		return nil
	}

	if numPerPage, ok := args[2].(IntResult); ok {
		numPerPageInt = int(numPerPage)
	} else {
		return nil
	}

	contentPaths := file.GetContentPaths(filePathString)
	return Paginate(contentPaths, templatePathString, numPerPageInt)
}

// Print takes a result and returns the string version of it
// TODO: This probably isn't needed unless the standard block
// should not return anything
func Print(result Result) StringResult {
	return StringResult(result.String())
}

func Paginate(contentPaths []string, templatePath string, numPerPage int) Result {
	paginationContexts := buildPaginationContexts(contentPaths, numPerPage)

	templateCache := parser.GetTemplateCache()
	templateNodes := templateCache.Get(templatePath)
	resultChan := make(chan Result)
	for _, paginationContext := range paginationContexts {
		nodeChan := make(chan parser.TreeNode)
		go func(nodeChan chan parser.TreeNode, templateNodes *[]parser.TreeNode) {
			for _, node := range *templateNodes {
				nodeChan <- node
			}
		}(nodeChan, templateNodes)

		processor := NewNodeProcessor(templatePath, paginationContext, nil, nil)
		go processor.Process(nodeChan, resultChan)
	}
	return <-resultChan
}

func buildPaginationContexts(contentPaths []string, numPerPage int) []*Context {
	if numPerPage == 0 {
		return []*Context{}
	}

	numPages := int(math.Ceil(float64(len(contentPaths)) / float64(numPerPage)))
	ret := make([]*Context, numPages)
	exportStore := GetExportStore()
	// for each page
	for curPage := 1; curPage <= numPages; curPage++ {
		// create a page context
		pageContext := &Context{
			"curPage":  &ContextNode{result: IntResult(curPage)},
			"numPages": &ContextNode{result: IntResult(numPages)},
			"prevPage": &ContextNode{result: StringResult("")},
			"nextPage": &ContextNode{result: StringResult("")},
		}

		// get the paths of the content files that will be on this page
		offset := (curPage - 1) * numPerPage
		last := minInt(len(contentPaths), offset+numPerPage)
		contentPathsOnPage := contentPaths[offset:last]

		// get the context for each content file on this page
		contextsOnPage := &Context{}
		for i, contentPath := range contentPathsOnPage {
			// key content like an array
			key := strconv.Itoa(i)
			(*contextsOnPage)[key] = &ContextNode{child: exportStore.Get(contentPath)}
		}

		// add content file contexts to pageContext
		(*pageContext)["content"] = &ContextNode{child: contextsOnPage}
		ret[curPage-1] = pageContext
	}

	return ret
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
