package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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
		processContent(config.GetContentPath(), outputPath)
		processContent(config.GetPagesPath(), outputPath)

		fmt.Println("Done")
	}
}

func createOutputDirs(outputPath string) {
	if err := os.MkdirAll(outputPath, 0750); err != nil {
		log.Fatalf("error creating output dir %q: %s\n", outputPath, err)
	}

}

func processContent(inputPath, outputPath string) {
	filepath.Walk(inputPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("error walking dir: %s", err)
		}
		// ignore nested dirs for now
		if info.IsDir() {
			return nil
		}

		if file, err := os.Open(path); err != nil {
			log.Fatalf("could not open file %s for processing", path)
		} else {
			log.Default().Printf("processing %s\n", path)
			renderFile(file, outputPath)
		}
		return nil
	})
}

func renderFile(contentFile *os.File, outputPath string) {
	nodeChan := make(chan parser.TreeNode)
	parserErrChan := make(chan error)

	nodeProcessor := processor.NewNodeProcessor(contentFile.Name(), nil, nil, nil)

	resultChan := make(chan processor.Result)
	lexer := lexer.Lexer{}
	tokChan, lexErrChan := lexer.Lex(contentFile)
	go parser.Parse(tokChan, nodeChan, parserErrChan)
	go nodeProcessor.Process(nodeChan, resultChan)

	for result := range resultChan {
		renderHTMLResult(result, contentFile.Name(), outputPath)
	}

	select {
	case tokErr := <-lexErrChan:
		if tokErr != nil {
			log.Fatalf("lexer error: %s", tokErr)
		}
	case parseErr := <-parserErrChan:
		if parseErr != nil {
			log.Fatalf("parser error: %s", parseErr)
		}
	}
}

func renderHTMLResult(result processor.Result, inputPath, outputPath string) {
	relativeInputPath := file.GetRelativePathTo(inputPath)
	fullPath := filepath.Join(outputPath, relativeInputPath)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0750); err != nil {
		log.Fatalf("could not create output dir %q: %s\n", filepath.Dir(fullPath), err)
	}

	if f, err := os.Create(fullPath); err != nil {
		fmt.Printf("could not create output file %q: %s\n", fullPath, err)
	} else {
		f.WriteString(result.String())
		f.Close()
	}
}
