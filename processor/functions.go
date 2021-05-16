package processor

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"

	"mettlach.codes/frizzy/config"
	"mettlach.codes/frizzy/file"
	"mettlach.codes/frizzy/parser"
)

func PaginateRaw(args ...Result) (Result, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("paginate expects 4 args, got %d", len(args))
	}

	var filePathString, templatePathString string
	var curPageInt, numPerPageInt int

	if curPage, ok := args[0].(IntResult); ok {
		curPageInt = int(curPage)
	} else {
		return nil, fmt.Errorf("expected current page to be an int, got %T", args[0])
	}

	// Path to content to be paginated
	if filePath, ok := args[1].(StringResult); ok {
		filePathString = string(filePath)
	} else {
		return nil, fmt.Errorf("expected file path to be an string, got %T", args[1])
	}

	// Path to the template to use for each content file on the page
	if templatePath, ok := args[2].(StringResult); ok {
		templatePathString = string(templatePath)
	} else {
		return nil, fmt.Errorf("expected template path to be an string, got %T", args[2])
	}

	// Number of content items per page
	if numPerPage, ok := args[3].(IntResult); ok {
		numPerPageInt = int(numPerPage)
	} else {
		return nil, fmt.Errorf("expected number per page to be an int, got %T", args[3])
	}

	contentPaths := file.GetContentPaths(filePathString)
	return Paginate(contentPaths, templatePathString, curPageInt, numPerPageInt)
}

// PagesBeforeRaw converts the Result type arguments to PagesRaw
// into the actual types that PagesRaw expects
func PagesBeforeRaw(args ...Result) (Result, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("pages before expects 3 args, got %d", len(args))
	}

	var curPageInt, numBeforeInt int
	var inputPathStr string

	if curPage, ok := args[0].(IntResult); ok {
		curPageInt = int(curPage)
	} else {
		return nil, fmt.Errorf("expected current page to be an int, got %T", args[0])
	}

	if inputPath, ok := args[1].(StringResult); ok {
		inputPathStr = string(inputPath)
	} else {
		return nil, fmt.Errorf("expected input path to be a string, got %T", args[1])
	}

	if numBefore, ok := args[2].(IntResult); ok {
		numBeforeInt = int(numBefore)
	} else {
		return nil, fmt.Errorf("expected num before to be an int, got %T", args[2])
	}

	return PagesBefore(curPageInt, numBeforeInt, inputPathStr)
}

func PagesAfterRaw(args ...Result) (Result, error) {
	if len(args) < 4 {
		return nil, fmt.Errorf("pages after exepects 4 args, got %d", len(args))
	}

	var curPageInt, numPagesInt, numAfterInt int
	var inputPathStr string

	if curPage, ok := args[0].(IntResult); ok {
		curPageInt = int(curPage)
	} else {
		return nil, fmt.Errorf("expected current page to be an int, got %T", args[0])
	}

	if numPages, ok := args[1].(IntResult); ok {
		numPagesInt = int(numPages)
	} else {
		return nil, fmt.Errorf("expected num pages to be an int, got %T", args[1])
	}

	if inputPath, ok := args[2].(StringResult); ok {
		inputPathStr = string(inputPath)
	} else {
		return nil, fmt.Errorf("expected input path to be a string, got %T", args[2])
	}

	if numAfter, ok := args[3].(IntResult); ok {
		numAfterInt = int(numAfter)
	} else {
		return nil, fmt.Errorf("expected num after to be an int, got %T", args[3])
	}

	return PagesAfter(curPageInt, numPagesInt, numAfterInt, inputPathStr)
}

func TemplateRaw(args ...Result) (Result, error) {
	var (
		ret Result
		err error
	)

	if len(args) != 1 {
		err = fmt.Errorf("`template` expects one argument, got %d", len(args))
	} else if templatePath, ok := args[0].(StringResult); !ok {
		err = fmt.Errorf("invalid template argument %s", args[0])
	} else {
		ret = Template(string(templatePath))
	}

	return ret, err
}

func Paginate(contentPaths []string, templatePath string, curPage int, numPerPage int) (Result, error) {
	paginationContext, err := buildPaginationContext(contentPaths, curPage, numPerPage)

	if err != nil {
		return nil, err
	}

	templateCache := parser.GetTemplateCache()
	templateNodes := templateCache.Get(templatePath)

	output := ""
	processor := NewNodeProcessor(templatePath, paginationContext, nil, nil, nil, 0, 0)
	for _, node := range *templateNodes {
		result, _ := processor.processHeadNode(node)
		output += result.String()
	}

	return StringResult(output), nil
}

func buildPaginationContext(contentPaths []string, curPage int, numPerPage int) (*Context, error) {
	var (
		pageContext *Context
		err         error
	)

	if numPerPage < 1 {
		err = fmt.Errorf("expected number of items per page to be > 0, got %d", numPerPage)
	} else if curPage < 1 {
		err = fmt.Errorf("expected current page to be > 0, got %d", curPage)
	} else {
		numPages := int(math.Ceil(float64(len(contentPaths)) / float64(numPerPage)))
		exportStore := GetExportStore()
		// create a page context
		pageContext = &Context{
			"curPage":  &ContextNode{result: IntResult(curPage)},
			"numPages": &ContextNode{result: IntResult(numPages)},
		}

		// get the paths of the content files that will be on this page
		offset := minInt(len(contentPaths), (curPage-1)*numPerPage)
		last := minInt(len(contentPaths), offset+numPerPage)
		contentPathsOnPage := contentPaths[offset:last]

		// get the context for each content file on this page
		contextsOnPage := &Context{}
		for i, contentPath := range contentPathsOnPage {
			// key content like an array
			key := fmt.Sprint(i)
			(*contextsOnPage)[key] = &ContextNode{child: exportStore.Get(contentPath)}
		}

		// add content file contexts to pageContext
		(*pageContext)["content"] = &ContextNode{child: contextsOnPage}
	}
	return pageContext, err
}

// PagesBefore builds a collection of contexts for the numBefore pages
// prior to the current page
// These can be iterated through to create pagination links
func PagesBefore(curPage, numBefore int, inputPath string) (ContainerResult, error) {
	ctx := &Context{}
	prevPage := maxInt((curPage - numBefore), 1)

	for ; prevPage < curPage; prevPage++ {
		prevCtx := &Context{}
		(*prevCtx)["_pageNum"] = &ContextNode{result: IntResult(prevPage)}
		(*prevCtx)["_pageHref"] = &ContextNode{
			result: StringResult(GetMarkdownOutputPath(inputPath, prevPage)),
		}

		key := fmt.Sprint(prevPage)
		(*ctx)[key] = &ContextNode{child: prevCtx}
	}

	return ContainerResult{context: ctx}, nil
}

func PagesAfter(curPage, numPages, numAfter int, inputPath string) (ContainerResult, error) {
	ctx := &Context{}
	endPage := minInt(numPages, curPage+numAfter)

	for nextPage := curPage + 1; nextPage <= endPage; nextPage++ {
		nextCtx := &Context{}
		(*nextCtx)["_pageNum"] = &ContextNode{result: IntResult(nextPage)}
		(*nextCtx)["_pageHref"] = &ContextNode{result: StringResult(GetMarkdownOutputPath(inputPath, nextPage))}

		key := fmt.Sprint(nextPage)
		(*ctx)[key] = &ContextNode{child: nextCtx}
	}

	return ContainerResult{context: ctx}, nil
}

func Template(templatePath string) Result {
	config := config.GetLoadedConfig()
	fullPath := filepath.Join(config.GetTemplatePath(), templatePath)

	if f, err := os.Open(fullPath); err != nil {
		log.Printf("could not open template file %s\n", fullPath)
	} else {
		if bytes, err := io.ReadAll(f); err != nil {
			log.Printf("could not read template file %s\n", fullPath)
		} else {
			return StringResult(bytes)
		}
	}

	return StringResult("")
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
