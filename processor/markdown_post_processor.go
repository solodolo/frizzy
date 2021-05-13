package processor

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"mettlach.codes/frizzy/config"
	"mettlach.codes/frizzy/file"
)

// Call turns a processed markdown result into a processed html result
// If the input is not markdown, it is passed through
func PostProcessMarkdown(inputPath string, resultChan <-chan Result) <-chan Result {
	if filepath.Ext(inputPath) == ".md" {
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

func getFullOutputPath(inputPath string) string {
	config := config.GetLoadedConfig()
	outputPath := config.OutputPath

	relativeInputPath := file.TrimRootPrefix(inputPath)
	return filepath.Join(outputPath, relativeInputPath)
}

// GetMarkdownOutputPath returns the html output path given
// the input path of a file
func GetMarkdownOutputPath(inputPath string) string {
	fullPath := getFullOutputPath(inputPath)
	fullPath = strings.TrimSuffix(fullPath, filepath.Ext(fullPath)) + ".html"

	return fullPath
}

// GetPagedMarkdownOutputPath returns the paginated output
// path given the input path and current page of a file
func GetPagedMarkdownOutputPath(inputPath string, curPage int) string {
	fullPath := getFullOutputPath(inputPath)
	trimmed := strings.TrimSuffix(fullPath, filepath.Ext(fullPath))
	fullPath = fmt.Sprintf("%s_%03d.html", trimmed, curPage)

	return fullPath
}
