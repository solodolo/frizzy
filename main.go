package main

import (
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
	if config, err := config.LoadConfig("/home/dmmettlach/workspace/frizzy/config.json"); err != nil {
		log.Fatal(fmt.Errorf("error loading config: %s", err))
	} else {
		outputPath := config.OutputPath

		const numConsumers = 1
		wg := sync.WaitGroup{}

		// have to process content first, then pages
		paths := []string{config.GetContentPath(), config.GetPagesPath()}
		for _, inputPath := range paths {
			wg.Add(numConsumers)
			pathChan, _ := walkFiles(inputPath)

			for i := 0; i < numConsumers; i++ {
				go func(index int) {
					consumer(pathChan)
					wg.Done()
				}(i)
			}

			wg.Wait()
		}

		fmt.Println("starting development server...")
		server := file.DevServer{ServerRoot: outputPath, Port: 8080}
		server.ListenAndServe()

		fmt.Println("Done")
	}
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

func consumer(pathChan <-chan string) {
	config := config.GetLoadedConfig()

	for path := range pathChan {
		if file, err := os.Open(path); err != nil {
			log.Fatalf("could not open file %s for processing\n", path)
		} else {
			fmt.Printf("processing %s\n", path)
			renderFile(file, config.OutputPath)
		}
	}
}

func createOutputDirs(outputPath string) {
	if err := os.MkdirAll(outputPath, 0750); err != nil {
		log.Fatalf("error creating output dir %q: %s\n", outputPath, err)
	}
}

func renderFile(contentFile *os.File, outputPath string) {

	nodeProcessor := processor.NewNodeProcessor(contentFile.Name(), nil, nil, nil)

	lexer := lexer.Lexer{}
	tokChan, lexErrChan := lexer.Lex(contentFile)
	nodeChan, parserErrChan := parser.Parse(tokChan)
	resultChan := nodeProcessor.Process(nodeChan)

	for result := range resultChan {
		renderHTMLResult(result, contentFile.Name(), outputPath)
	}

	tokErr := <-lexErrChan
	if tokErr != nil {
		log.Fatalf("lexer error: %s", tokErr)
	}
	parseErr := <-parserErrChan
	if parseErr != nil {
		log.Fatalf("parser error: %s", parseErr)
	}

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
