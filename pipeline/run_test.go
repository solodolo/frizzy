package pipeline

import (
	"testing"

	"mettlach.codes/frizzy/config"
)

var testConfig *config.Config

func TestMain(m *testing.M) {
	testConfig, _ = config.LoadConfig("../test_files/test_config.json")
	m.Run()
}

func BenchmarkTestRun(b *testing.B) {
	for i := 0; i < b.N; i++ {
		templatePathChan, _ := WalkFiles(testConfig.GetTemplatePath())
		contentPathChan, _ := WalkFiles(testConfig.GetContentPath())
		pagesPathChan, _ := WalkFiles(testConfig.GetPagesPath())

		RunPipeline(templatePathChan, TemplateCacheHandler)

		RunPipeline(contentPathChan, FullPipelineNullRenderer)

		RunPipeline(pagesPathChan, FullPipelineNullRenderer)
	}
}
