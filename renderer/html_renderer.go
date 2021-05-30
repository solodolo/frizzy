package renderer

import (
	"fmt"
	"os"
	"path/filepath"

	"mettlach.codes/frizzy/processor"
)

func RenderHtmlResults(resultChan <-chan processor.Result, outputPath string) <-chan error {
	errChan := make(chan error, 1)
	go func() {
		defer close(errChan)

		for result := range resultChan {
			outputErr := renderHTMLResult(result, outputPath)

			if outputErr != nil {
				errChan <- fmt.Errorf("failed to write to %s, %s", outputPath, outputErr)
				return
			}
		}
	}()

	return errChan
}

func renderHTMLResult(result processor.Result, outputPath string) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0750); err != nil {
		return err
	}

	if f, err := os.Create(outputPath); err != nil {
		return err
	} else {
		defer f.Close()
		f.WriteString(result.String())
	}

	return nil
}
