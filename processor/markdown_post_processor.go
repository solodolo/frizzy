package processor

import (
	"path/filepath"

	"github.com/gomarkdown/markdown"
)

// MarkdownPostProcessor handles turning markdown into html
type MarkdownPostProcessor struct {
	Filepath string
}

// Call turns a processed markdown result into a processed html result
// If the input is not markdown, it is passed through
func (receiver *MarkdownPostProcessor) Call(result Result) Result {
	if filepath.Ext(receiver.Filepath) == ".md" {
		inputBytes := []byte(result.String())
		mdBytes := markdown.ToHTML(inputBytes, nil, nil)
		return StringResult(mdBytes)
	}

	return result
}
