package processor

import (
	"path/filepath"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

// MarkdownPostProcessor handles turning markdown into html
type MarkdownPostProcessor struct {
	Filepath string
}

// Call turns a processed markdown result into a processed html result
// If the input is not markdown, it is passed through
func (receiver *MarkdownPostProcessor) Call(resultChan <-chan Result) <-chan Result {
	if filepath.Ext(receiver.Filepath) == ".md" {
		postProcessChan := make(chan Result)
		go func() {
			defer close(postProcessChan)
			for result := range resultChan {
				parser := parser.NewWithExtensions(parser.FencedCode)
				inputBytes := []byte(result.String())
				mdBytes := string(markdown.ToHTML(inputBytes, parser, nil))
				postProcessChan <- StringResult(mdBytes)
			}
		}()

		return postProcessChan
	}

	return resultChan
}
