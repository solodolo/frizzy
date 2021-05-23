package pipeline

import (
	"testing"

	"mettlach.codes/frizzy/config"
)

func BenchmarkTestRun(b *testing.B) {
	config, _ := config.LoadConfig("../test_files/config.json")

	for i := 0; i < b.N; i++ {
		templatePathChan, _ := WalkFiles(config.GetTemplatePath())
		contentPathChan, _ := WalkFiles(config.GetContentPath())
		pagesPathChan, _ := WalkFiles(config.GetPagesPath())

		RunPipeline(templatePathChan, TemplateCacheHandler)

		RunPipeline(contentPathChan, FullPipelineHandler)

		RunPipeline(pagesPathChan, FullPipelineHandler)
	}
}
