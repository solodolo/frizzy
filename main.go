package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
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

	configPath := os.Args[1]

	if config, err := config.LoadConfig(configPath); err != nil {
		log.Printf("error loading config: %s\n", err)
		return
	} else {
		templatePathChan, _ := walkFiles(config.GetTemplatePath())
		contentPathChan, _ := walkFiles(config.GetContentPath())
		pagesPathChan, _ := walkFiles(config.GetPagesPath())
		// have to process templates first, then content, then pages
		log.Println("pipelining template files")
		if err := runPipeline(templatePathChan, templateCacher); err != nil {
			log.Println("exiting")
			return
		}

		log.Println("pipelining content files")
		if err := runPipeline(contentPathChan, fileRenderer); err != nil {
			log.Println("exiting")
			return
		}

		log.Println("pipelining page files")
		if err := runPipeline(pagesPathChan, fileRenderer); err != nil {
			log.Println("exiting")
			return
		}

		log.Println("starting development server...")
		server := file.DevServer{ServerRoot: config.OutputPath, Port: 8080}
		server.ListenAndServe()

		log.Println("Done")
	}
}

func printUsage() {
	log.Println("usage: frizzy /path/to/config.json")
}

func runPipeline(pathChan <-chan string, handler func(context.Context, *os.File) []<-chan error) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChans := []<-chan error{}
	for inputPath := range pathChan {
		log.Printf("    %s\n", inputPath)
		f, err := os.Open(inputPath)

		if err != nil {
			log.Printf("pipeline error: %s\n", err)
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

	merge := func(ec <-chan error) {
		defer wg.Done()
		for err := range ec {
			select {
			case errChan <- err:
			case <-ctx.Done():
				return
			}
		}
	}

	for _, ec := range errChans {
		go merge(ec)
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

func createOutputDirs(outputPath string) {
	if err := os.MkdirAll(outputPath, 0750); err != nil {
		log.Fatalf("error creating output dir %q: %s\n", outputPath, err)
	}
}

func templateCacher(ctx context.Context, templateFile *os.File) []<-chan error {
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

	go func(templateCache *parser.TemplateCache, cacheKey string) {
		for node := range nodeChan {
			templateCache.Insert(cacheKey, node)
		}
	}(templateCache, cacheKey)

	return []<-chan error{lexErrChan, parserErrChan}
}

func fileRenderer(ctx context.Context, contentFile *os.File) []<-chan error {
	config := config.GetLoadedConfig()
	outputPath := config.OutputPath

	lexer := lexer.Lexer{}
	tokChan, lexErrChan := lexer.Lex(contentFile, ctx)
	nodeChan, parserErrChan := parser.Parse(tokChan, ctx)

	nodeProcessor := processor.NewNodeProcessor(contentFile.Name(), nil, nil, nil, nil)
	resultChan := nodeProcessor.Process(nodeChan, ctx)

	go func(contentFile *os.File, outputPath string) {
		for result := range resultChan {
			renderHTMLResult(result, contentFile.Name(), outputPath)
		}
	}(contentFile, outputPath)

	return []<-chan error{lexErrChan, parserErrChan}
}

func renderHTMLResult(result processor.Result, inputPath, outputPath string) {
	relativeInputPath := file.GetRelativePathTo(inputPath)
	fullPath := filepath.Join(outputPath, relativeInputPath)

	if !strings.HasSuffix(fullPath, ".html") {
		fullPath = strings.TrimSuffix(fullPath, filepath.Ext(fullPath)) + ".html"
	}

	if err := os.MkdirAll(filepath.Dir(fullPath), 0750); err != nil {
		log.Fatalf("could not create output dir %q: %s\n", filepath.Dir(fullPath), err)
	}

	f, err := os.Create(fullPath)
	defer f.Close()

	if err != nil {
		fmt.Printf("could not create output file %q: %s\n", fullPath, err)
	} else {
		f.WriteString(result.String())
	}
}
