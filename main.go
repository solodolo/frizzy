package main

import (
	"bufio"
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
		// if err := unix.Chroot(config.RootPath); err != nil {
		// 	log.Fatalf("chroot error: %s", err)
		// }
		contentPath := config.ContentPath

		filepath.Walk(contentPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Fatalf("error walking content dir: %s", err)
			}
			// ignore nested dirs for now
			if info.IsDir() {
				return nil
			}

			if contentFile, err := os.Open(path); err != nil {
				log.Fatalf("could not open content file %s for processing", path)
			} else {
				log.Default().Printf("processing %s\n", path)
				tokChan := make(chan []lexer.Token)
				errChan := make(chan error)

				nodeChan := make(chan parser.TreeNode)
				parserErrChan := make(chan error)

				nodeProcessor := processor.NodeProcessor{
					Context: &processor.Context{},
				}

				resultChan := make(chan processor.Result)
				go lexer.Lex(bufio.NewScanner(contentFile), tokChan, errChan)
				go parser.Parse(tokChan, nodeChan, parserErrChan)
				go nodeProcessor.Process(nodeChan, resultChan)

				for result := range resultChan {
					fmt.Println(result)
				}

				select {
				case tokErr := <-errChan:
					if tokErr != nil {
						log.Fatalf("lexer error: %s", tokErr)
					}
				case parseErr := <-parserErrChan:
					if parseErr != nil {
						log.Fatalf("parser error: %s", parseErr)
					}
				}
			}
			return nil
		})

		fmt.Println("Done")
	}
}
