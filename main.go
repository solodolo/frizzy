package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"mettlach.codes/frizzy/config"
	"mettlach.codes/frizzy/lexer"
	"mettlach.codes/frizzy/parser"
	"mettlach.codes/frizzy/processor"
)

func main() {
	if config, err := config.LoadConfig("/home/dmmettlach/workspace/frizzy/config.json"); err != nil {
		log.Fatal(fmt.Errorf("error loading config: %s", err))
	} else {
		processContent(config.GetContentPath())
		processContent(config.GetPagesPath())

		fmt.Println("Done")
	}
}

func processContent(path string) {
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
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
			renderFile(file)
		}
		return nil
	})
}

func renderFile(contentFile *os.File) {
	nodeChan := make(chan parser.TreeNode)
	parserErrChan := make(chan error)

	nodeProcessor := processor.NewNodeProcessor(contentFile.Name(), nil, nil, nil)

	resultChan := make(chan processor.Result)
	lexer := lexer.Lexer{}
	tokChan, lexErrChan := lexer.Lex(contentFile)
	go parser.Parse(tokChan, nodeChan, parserErrChan)
	go nodeProcessor.Process(nodeChan, resultChan)

	for result := range resultChan {
		fmt.Print(result)
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
