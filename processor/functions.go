package processor

import (
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"mettlach.codes/frizzy/config"
	"mettlach.codes/frizzy/file"
	"mettlach.codes/frizzy/parser"
)

func PrintRaw(args ...Result) Result {
	return Print(args[0])
}

func PaginateRaw(args ...Result) Result {
	// TODO: replace nils with errors
	if len(args) < 4 {
		return nil
	}

	var filePathString, templatePathString string
	var curPageInt, numPerPageInt int

	if curPage, ok := args[0].(IntResult); ok {
		curPageInt = int(curPage)
	} else {
		return nil
	}

	// Path to content to be paginated
	if filePath, ok := args[1].(StringResult); ok {
		filePathString = string(filePath)
	} else {
		return nil
	}

	// Path to the template to use for each content file on the page
	if templatePath, ok := args[2].(StringResult); ok {
		templatePathString = string(templatePath)
	} else {
		return nil
	}

	// Number of content items per page
	if numPerPage, ok := args[3].(IntResult); ok {
		numPerPageInt = int(numPerPage)
	} else {
		return nil
	}

	contentPaths := file.GetContentPaths(filePathString)
	return Paginate(contentPaths, templatePathString, curPageInt, numPerPageInt)
}

func TemplateRaw(args ...Result) Result {
	var ret Result

	if len(args) != 1 {
		log.Printf("`template` expects one argument, got %d", len(args))
	} else if templatePath, ok := args[0].(StringResult); !ok {
		log.Printf("invalid template argument %s\n", args[0])
	} else {
		ret = Template(string(templatePath))
	}

	return ret
}

// Print takes a result and returns the string version of it
// TODO: This probably isn't needed unless the standard block
// should not return anything
func Print(result Result) StringResult {
	return StringResult(result.String())
}

func Paginate(contentPaths []string, templatePath string, curPage int, numPerPage int) Result {
	paginationContext := BuildPaginationContext(contentPaths, templatePath, curPage, numPerPage)

	templateCache := parser.GetTemplateCache()
	templateNodes := templateCache.Get(templatePath)

	output := ""
	processor := NewNodeProcessor(templatePath, paginationContext, nil, nil, nil)
	for _, node := range *templateNodes {
		result, _ := processor.processHeadNode(node)
		output += result.String()
	}

	return StringResult(output)
}

func BuildPaginationContext(contentPaths []string, templatePath string, curPage int, numPerPage int) *Context {
	if numPerPage == 0 {
		return nil
	}

	numPages := int(math.Ceil(float64(len(contentPaths)) / float64(numPerPage)))
	exportStore := GetExportStore()
	// create a page context
	pageContext := &Context{
		"curPage":      &ContextNode{result: IntResult(curPage)},
		"templatePath": &ContextNode{result: StringResult(templatePath)},
		"numPages":     &ContextNode{result: IntResult(numPages)},
		"prevPage":     &ContextNode{result: StringResult("")},
		"nextPage":     &ContextNode{result: StringResult("")},
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

	return pageContext
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
