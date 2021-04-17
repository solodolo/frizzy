package parser

import "sync"

// TemplateCache caches parsed treenodes for template paths
type TemplateCache struct {
	cache map[string][]TreeNode
}

var once sync.Once
var cache *TemplateCache

func createCache() {
	if cache == nil {
		cache = &TemplateCache{cache: make(map[string][]TreeNode)}
	}
}

// GetTemplateCache returns the TemplateCache singleton
func GetTemplateCache() *TemplateCache {
	once.Do(createCache)
	return cache
}

var mut sync.Mutex

// Insert inserts the key, value pair of the export context
// represented by filename
func (receiver *TemplateCache) Insert(key string, value TreeNode) {
	mut.Lock()
	defer mut.Unlock()
	receiver.cache[key] = append(receiver.cache[key], value)
}

// Get returns the export context of the given filename
func (receiver *TemplateCache) Get(key string) *[]TreeNode {
	nodes := receiver.cache[key]
	return &nodes
}
