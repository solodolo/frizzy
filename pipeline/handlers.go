package pipeline

import (
	"context"
	"io"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	"mettlach.codes/frizzy/config"
	"mettlach.codes/frizzy/file"
	"mettlach.codes/frizzy/lexer"
	"mettlach.codes/frizzy/parser"
	"mettlach.codes/frizzy/processor"
)

func TemplateCacheHandler(ctx context.Context, templateFile *os.File) []<-chan error {
	lexer := lexer.Lexer{}
	tokChan, lexErrChan := lexer.Lex(templateFile, ctx)
	nodeChan, parserErrChan := parser.Parse(tokChan, ctx)

	config := config.GetLoadedConfig()
	templatePath := config.GetTemplatePath()
	templateCache := parser.GetTemplateCache()
	cacheKey := strings.TrimPrefix(templateFile.Name(), templatePath)
	if cacheKey[0] == '/' {
		cacheKey = cacheKey[1:]
	}

	cacherErrs := processor.CacheTemplateResults(nodeChan, templateCache, cacheKey)
	return []<-chan error{lexErrChan, parserErrChan, cacherErrs}
}

func getNumPages(inputFile *os.File) (int, bool) {
	var (
		numPages int
		ok       bool = true
	)

	paginationRegex := regexp.MustCompile(`[{{|{{:]\s*paginate\(("[^"]+"),\s*("[^"]+"),\s*(\d+)\)\s*}}`)
	if bytes, err := io.ReadAll(inputFile); err != nil {
		ok = false
	} else if found := paginationRegex.FindSubmatch(bytes); found != nil {
		contentPath := strings.ReplaceAll(string(found[1]), "\"", "")
		if numPerPage, err := strconv.Atoi(string(found[3])); err != nil {
			ok = false
		} else {
			contentPaths := file.GetContentPaths(contentPath)
			numPages = int(math.Ceil(float64(len(contentPaths)) / float64(numPerPage)))
		}
	} else {
		ok = false
	}

	return numPages, ok
}

func FullPipelineHandler(ctx context.Context, contentFile *os.File) []<-chan error {
	numPages, paginated := getNumPages(contentFile)
	// Rewind the file so it can be lexed from the start
	contentFile.Seek(0, 0)

	inputPath := contentFile.Name()
	lexer := lexer.Lexer{}
	tokChan, lexErrChan := lexer.Lex(contentFile, ctx)
	nodeChan, parserErrChan := parser.Parse(tokChan, ctx)

	if paginated {
		// fan out node chan to processors for each page
		nodeChans := fanOutNodes(nodeChan, numPages)
		// processor and renderer chans for each page plus
		// lexer and parser err chans
		pagedErrChans := make([]<-chan error, 0, numPages*2+2)

		for i, fannedNodeChan := range nodeChans {
			curPage := i + 1
			processorErrChan, rendererErrChan := processAndRender(ctx, inputPath, fannedNodeChan, curPage, numPages)
			pagedErrChans = append(pagedErrChans, processorErrChan, rendererErrChan)
		}

		pagedErrChans = append(pagedErrChans, lexErrChan, parserErrChan)
		return pagedErrChans
	} else {
		processorErrChan, rendererErrChan := processAndRender(ctx, inputPath, nodeChan, 0, 0)
		return []<-chan error{lexErrChan, parserErrChan, processorErrChan, rendererErrChan}
	}
}

func fanOutNodes(nodeChan <-chan parser.TreeNode, numPages int) []chan parser.TreeNode {
	nodeChanFan := make([]chan parser.TreeNode, numPages)
	for i := range nodeChanFan {
		nodeChanFan[i] = make(chan parser.TreeNode, 1)
	}

	go func() {
		defer func() {
			for _, fanNodeChan := range nodeChanFan {
				close(fanNodeChan)
			}
		}()

		for node := range nodeChan {
			for _, fanChan := range nodeChanFan {
				fanChan <- node
			}
		}
	}()

	return nodeChanFan
}

func processAndRender(ctx context.Context, inputPath string, nodeChan <-chan parser.TreeNode, curPage, numPages int) (<-chan error, <-chan error) {
	nodeProcessor := processor.NewNodeProcessor(inputPath, nil, nil, nil, nil, curPage, numPages)

	outputPath := processor.GetMarkdownOutputPath(inputPath, curPage)
	nodeProcessor.ExportStore.Insert([]string{"_href"}, processor.StringResult(outputPath))
	processorChan, processorErrChan := nodeProcessor.Process(nodeChan, ctx)
	resultChan := processor.PostProcessMarkdown(inputPath, processorChan)
	rendererErrChan := processor.RenderHtmlResults(resultChan, outputPath)

	return processorErrChan, rendererErrChan
}
