package renderer

import (
	"context"

	"mettlach.codes/frizzy/parser"
	"mettlach.codes/frizzy/processor"
)

// RenderNullResults reads the results and drops them
func RenderNullResults(ctx context.Context, inputPath string, nodeChan <-chan parser.TreeNode, curPage, numPages int) (<-chan error, <-chan error) {
	nodeProcessor := processor.NewNodeProcessor(inputPath, nil, nil, nil, nil, curPage, numPages)

	outputPath := processor.GetMarkdownOutputPath(inputPath, curPage)
	nodeProcessor.ExportStore.Insert([]string{"_href"}, processor.StringResult(outputPath))
	processorChan, processErrorChan := nodeProcessor.Process(nodeChan, ctx)
	resultChan := processor.PostProcessMarkdown(inputPath, processorChan)

	doneChan := make(chan error)
	go func() {
		defer close(doneChan)
		for range resultChan {
		}
	}()

	return processErrorChan, doneChan
}
