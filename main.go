package main

import (
	"context"
	"flag"
	"fmt"
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
		if err := runPipeline(templatePathChan, templateCacher); err != nil {
			log.Println("exiting")
			return
		} else {
			log.Println("finished template files")
		}

		log.Println("pipelining content files")
		if err := runPipeline(contentPathChan, fileRenderer); err != nil {
			log.Println("exiting")
			return
		} else {
			log.Println("finished content files")
		}

		log.Println("pipelining page files")
		if err := runPipeline(pagesPathChan, fileRenderer); err != nil {
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

func getPaginationInfo(file *os.File) (string, int, bool) {
	var (
		contentPath string
		numPerPage  int
		ok          bool = true
	)

	paginationRegex := regexp.MustCompile(`{{\s*paginate\(("[^"]+"),\s*("[^"]+"),\s*(\d+)\)\s*}}`)
	if bytes, err := io.ReadAll(file); err != nil {
		ok = false
	} else if found := paginationRegex.FindSubmatch(bytes); found != nil {
		contentPath = string(found[1])
		if numPerPage, err = strconv.Atoi(string(found[3])); err != nil {
			ok = false
		}
	} else {
		ok = false
	}

	return contentPath, numPerPage, ok
}

func runPipeline(pathChan <-chan string, handler func(context.Context, *os.File, int) []<-chan error) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChans := []<-chan error{}
	for inputPath := range pathChan {
		log.Printf("    %s -> ", inputPath)
		f, err := os.Open(inputPath)

		if err != nil {
			log.Printf("pipeline error: %s\n", err)
			continue
		}

		if contentPath, numPerPage, ok := getPaginationInfo(f); ok {
			contentPaths := file.GetContentPaths(contentPath)
			numPages := int(math.Ceil(float64(len(contentPaths)) / float64(numPerPage)))

			for curPage := 1; curPage <= numPages; curPage++ {
				paginationCtx := &processor.Context{}
				paginationCtx.Insert([]string{"curPage"}, processor.IntResult(curPage))
				errChans = append(errChans, handler(ctx, f, curPage)...)
			}
		} else {
			errChans = append(errChans, handler(ctx, f, -1)...)
		}
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

func templateCacher(ctx context.Context, templateFile *os.File, curPage int) []<-chan error {
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

	doneChan := make(chan error)
	go func(templateCache *parser.TemplateCache, cacheKey string) {
		defer close(doneChan)
		for node := range nodeChan {
			templateCache.Insert(cacheKey, node)
		}
	}(templateCache, cacheKey)

	return []<-chan error{lexErrChan, parserErrChan, doneChan}
}

func fileRenderer(ctx context.Context, contentFile *os.File, curPage int) []<-chan error {
	lexer := lexer.Lexer{}
	tokChan, lexErrChan := lexer.Lex(contentFile, ctx)
	nodeChan, parserErrChan := parser.Parse(tokChan, ctx)

	var processorCtx *processor.Context

	if curPage > 0 {
		processorCtx.Insert([]string{"curPage"}, processor.IntResult(curPage))
	}

	nodeProcessor := processor.NewNodeProcessor(contentFile.Name(), processorCtx, nil, nil, nil)
	processorChan, processorErrChan := nodeProcessor.Process(nodeChan, ctx)

	postProcessor := processor.MarkdownPostProcessor{Filepath: contentFile.Name()}
	resultChan := postProcessor.Call(processorChan)

	doneChan := make(chan error)
	go func(resultChan <-chan processor.Result, contentFile *os.File) {
		defer close(doneChan)
		for result := range resultChan {
			outputPath := getOutputPath(contentFile.Name(), curPage)
			log.Println(outputPath)
			renderHTMLResult(result, outputPath)
		}
	}(resultChan, contentFile)

	return []<-chan error{lexErrChan, parserErrChan, processorErrChan, doneChan}
}

func getOutputPath(inputPath string, curPage int) string {
	config := config.GetLoadedConfig()
	outputPath := config.OutputPath

	relativeInputPath := file.TrimRootPrefix(inputPath)
	fullPath := filepath.Join(outputPath, relativeInputPath)

	if curPage < 1 {
		fullPath = strings.TrimSuffix(fullPath, filepath.Ext(fullPath)) + ".html"
	} else {
		trimmed := strings.TrimSuffix(fullPath, filepath.Ext(fullPath))
		fullPath = filepath.Join(trimmed, fmt.Sprint(curPage)) + ".html"
	}

	return fullPath
}

func renderHTMLResult(result processor.Result, outputPath string) {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0750); err != nil {
		log.Fatalf("could not create output dir %q: %s\n", filepath.Dir(outputPath), err)
	}

	f, err := os.Create(outputPath)

	if err != nil {
		fmt.Printf("could not create output file %q: %s\n", outputPath, err)
	} else {
		defer f.Close()
		f.WriteString(result.String())
	}
}
