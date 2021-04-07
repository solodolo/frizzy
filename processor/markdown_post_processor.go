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
func (receiver *MarkdownPostProcessor) Call(result Result) Result {
	if filepath.Ext(receiver.Filepath) == ".md" {
		parser := parser.NewWithExtensions(parser.FencedCode)
		inputBytes := []byte(result.String())
		mdBytes := string(markdown.ToHTML(inputBytes, parser, nil))
		return StringResult(mdBytes)
	}

	return result
}
