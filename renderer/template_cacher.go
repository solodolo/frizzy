package renderer

import "mettlach.codes/frizzy/parser"

func CacheTemplateResults(nodeChan <-chan parser.TreeNode, templateCache *parser.TemplateCache, cacheKey string) <-chan error {
	errChan := make(chan error)
	go func(templateCache *parser.TemplateCache, cacheKey string) {
		defer close(errChan)

		for node := range nodeChan {
			templateCache.Insert(cacheKey, node)
		}
	}(templateCache, cacheKey)

	return errChan
}
