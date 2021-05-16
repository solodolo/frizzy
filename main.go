package main

import (
	"context"
	"flag"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"mettlach.codes/frizzy/config"
	"mettlach.codes/frizzy/file"
	"mettlach.codes/frizzy/lexer"
	"mettlach.codes/frizzy/parser"
	"mettlach.codes/frizzy/processor"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	startDevServer := flag.Bool("d", false, "start a web server to serve files in output directory")
	devServerPort := flag.Int("p", 8080, "the port the web server will listen on")
	clearOutput := flag.Bool("c", false, "clear any existing output")
	flag.Parse()

	configPath := os.Args[len(os.Args)-1]

	if config, err := config.LoadConfig(configPath); err != nil {
		log.Printf("error loading config: %s\n", err)
		return
	} else {
		if *clearOutput {
			log.Printf("removing %s\n", config.OutputPath)
			if err := clearOutputDirectory(config.OutputPath); err != nil {
				log.Print(err)
				return
			}
		}

		templatePathChan, _ := walkFiles(config.GetTemplatePath())
		contentPathChan, _ := walkFiles(config.GetContentPath())
		pagesPathChan, _ := walkFiles(config.GetPagesPath())
		// have to process templates first, then content, then pages
		log.Println("pipelining template files")
		if err := runPipeline(templatePathChan, templatePipeline); err != nil {
			log.Println("exiting")
			return
		} else {
			log.Println("finished template files")
		}

		log.Println("pipelining content files")
		if err := runPipeline(contentPathChan, fullPipeline); err != nil {
			log.Println("exiting")
			return
		} else {
			log.Println("finished content files")
		}

		log.Println("pipelining page files")
		if err := runPipeline(pagesPathChan, fullPipeline); err != nil {
			log.Println("exiting")
			return
		} else {
			log.Println("finished page files")
		}

		if *startDevServer {
			log.Println("starting development server...")
			server := file.DevServer{ServerRoot: config.OutputPath, Port: *devServerPort}
			server.ListenAndServe()
		}

		log.Println("Done")
	}
}

func printUsage() {
	log.Println("usage: frizzy [-c] /path/to/config.json")
}

func clearOutputDirectory(outputDir string) error {
	return os.RemoveAll(outputDir)
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

func runPipeline(pathChan <-chan string, handler func(context.Context, *os.File) []<-chan error) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChans := []<-chan error{}
	files := make([]*os.File, len(pathChan))

	defer func() {
		for _, f := range files {
			f.Close()
		}
	}()

	for inputPath := range pathChan {
		log.Printf("    %s", inputPath)

		f, err := os.Open(inputPath)
		files = append(files, f)

		if err != nil {
			log.Printf("    pipeline error: failed to open %s, %s\n", inputPath, err)
			continue
		}

		errChans = append(errChans, handler(ctx, f)...)
	}

	errChan := mergeErrChans(ctx, errChans)

	for err := range errChan {
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func mergeErrChans(ctx context.Context, errChans []<-chan error) <-chan error {
	wg := sync.WaitGroup{}
	wg.Add(len(errChans))

	errChan := make(chan error)

	merge := func(idx int, ec <-chan error, errChans []<-chan error) {
		defer wg.Done()
		for err := range ec {
			select {
			case errChan <- err:
			case <-ctx.Done():
				return
			}
		}
	}

	for i, ec := range errChans {
		go merge(i, ec, errChans)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	return errChan
}

func walkFiles(inputPath string) (<-chan string, <-chan error) {
	pathChan := make(chan string)
	errChan := make(chan error, 1)

	go func() {
		defer close(pathChan)
		defer close(errChan)

		errChan <- filepath.Walk(inputPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// ignore nested dirs for now
			if info.IsDir() {
				return nil
			}

			pathChan <- path
			return nil
		})
	}()

	return pathChan, errChan
}

func templatePipeline(ctx context.Context, templateFile *os.File) []<-chan error {
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

func fullPipeline(ctx context.Context, contentFile *os.File) []<-chan error {
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
			processorErrChan, _ := processAndRender(ctx, inputPath, fannedNodeChan, curPage, numPages)
			pagedErrChans = append(pagedErrChans, processorErrChan)
		}

		pagedErrChans = append(pagedErrChans, lexErrChan, parserErrChan)
		return pagedErrChans
	} else {
		processorErrChan, _ := processAndRender(ctx, inputPath, nodeChan, 0, 0)
		return []<-chan error{lexErrChan, parserErrChan, processorErrChan}
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
