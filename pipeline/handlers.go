package pipeline

import (
	"context"
	"io"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"mettlach.codes/frizzy/config"
	"mettlach.codes/frizzy/file"
	"mettlach.codes/frizzy/lexer"
	"mettlach.codes/frizzy/parser"
	"mettlach.codes/frizzy/processor"
	"mettlach.codes/frizzy/renderer"
)

func TemplateCacheHandler(ctx context.Context, templateFile *os.File) <-chan error {
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

	cacherErrs := renderer.CacheTemplateResults(nodeChan, templateCache, cacheKey)
	return mergeIntoStandardErrs(ctx, templateFile.Name(), lexErrChan, parserErrChan, cacherErrs)
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

func FullPipelineHtmlRenderer(ctx context.Context, contentFile *os.File) <-chan error {
	renderer := processAndRender
	return FullPipelineHandler(ctx, contentFile, renderer)
}

func FullPipelineNullRenderer(ctx context.Context, contentFile *os.File) <-chan error {
	renderer := renderer.RenderNullResults
	return FullPipelineHandler(ctx, contentFile, renderer)
}

func FullPipelineHandler(ctx context.Context, contentFile *os.File, renderer func(context.Context, string, <-chan parser.TreeNode, int, int) (<-chan error, <-chan error)) <-chan error {
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
			processorErrChan, rendererErrChan := renderer(ctx, inputPath, fannedNodeChan, curPage, numPages)
			pagedErrChans = append(pagedErrChans, processorErrChan, rendererErrChan)
		}

		pagedErrChans = append(pagedErrChans, lexErrChan, parserErrChan)
		return mergeIntoStandardErrs(ctx, contentFile.Name(), pagedErrChans...)
	} else {
		processorErrChan, rendererErrChan := renderer(ctx, inputPath, nodeChan, 0, 0)
		return mergeIntoStandardErrs(
			ctx,
			contentFile.Name(),
			lexErrChan,
			parserErrChan,
			processorErrChan,
			rendererErrChan,
		)
	}
}

func mergeIntoStandardErrs(ctx context.Context, filename string, errChans ...<-chan error) <-chan error {
	wg := sync.WaitGroup{}
	wg.Add(len(errChans))

	errChan := make(chan error)

	merge := func(idx int, ec <-chan error) {
		defer wg.Done()
		for err := range ec {
			if err == nil {
				log.Println(ec)
			}
			stdErr := &StandardError{Filename: filename, Message: err.Error()}
			select {
			case errChan <- stdErr:
			case <-ctx.Done():
				return
			}
		}
	}

	for i, ec := range errChans {
		go merge(i, ec)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	return errChan
}

func fanOutNodes(nodeChan <-chan parser.TreeNode, numPages int) []chan parser.TreeNode {
	nodeChanFan := make([]chan parser.TreeNode, numPages)
	for i := range nodeChanFan {
		nodeChanFan[i] = make(chan parser.TreeNode, 10)
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
	rendererErrChan := renderer.RenderHtmlResults(resultChan, outputPath)

	return processorErrChan, rendererErrChan
}
