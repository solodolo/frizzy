package processor

import (
	"math"
	"strconv"

	"mettlach.codes/frizzy/file"
	"mettlach.codes/frizzy/parser"
)

func Print(result Result) StringResult {
	return StringResult(result.String())
}

func Paginate(filePath, templatePath StringResult, numPerPage int) StringResult {
	contentPaths := file.GetContentPaths(string(filePath))
	numPages := math.Ceil(float64(len(contentPaths)) / float64(numPerPage))
	paginationContexts := buildPaginationContexts(contentPaths, int(numPages), numPerPage)

	templateCache := parser.GetTemplateCache()
	templateNodes := templateCache.Get(string(templatePath))
	for _, paginationContext := range paginationContexts {
		nodeChan := make(chan parser.TreeNode)
		go func(nodeChan chan parser.TreeNode, templateNodes *[]parser.TreeNode) {
			for _, node := range *templateNodes {
				nodeChan <- node
			}
		}(nodeChan, templateNodes)

		Process(nodeChan, paginationContext)
	}
	return StringResult("")
}

func buildPaginationContexts(contentPaths []string, numPages, numPerPage int) []*Context {
	ret := make([]*Context, numPages)
	exportStore := GetExportStore()
	// for each page
	for curPage := 1; curPage <= numPages; curPage++ {
		// create a page context
		pageContext := &Context{
			"curPage":  ContextNode{result: IntResult(curPage)},
			"numPages": ContextNode{result: IntResult(numPages)},
			"prevPage": ContextNode{result: StringResult("")},
			"nextPage": ContextNode{result: StringResult("")},
		}

		// get the paths of the content files that will be on this page
		offset := (curPage - 1) * numPerPage
		contentPathsOnPage := contentPaths[offset : offset+numPerPage]

		// get the context for each content file on this page
		contextsOnPage := &Context{}
		for i, contentPath := range contentPathsOnPage {
			// key content like an array
			key := strconv.Itoa(i)
			(*contextsOnPage)[key] = ContextNode{child: exportStore.Get(contentPath)}
		}

		// add content file contexts to pageContext
		(*pageContext)["content"] = ContextNode{child: contextsOnPage}
		ret[curPage-1] = pageContext
	}

	return ret
}
