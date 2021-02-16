package processor

import (
	"sync"
)

// Files can export variables for use by other files
// All assignment variables are exported

// ExportStore is a singleton to read and write export vars
type ExportStore struct {
	exports map[string]Context
}

var once sync.Once
var store *ExportStore

func createStore() {
	if store == nil {
		store = &ExportStore{}
	}
}

// GetExportStore returns the exportStore singleton
func GetExportStore() *ExportStore {
	once.Do(createStore)
	return store
}

// Insert inserts the key, value pair of the export context
// represented by filename
func (receiver *ExportStore) Insert(filename, key string, value Result) {
	if _, ok := receiver.exports[filename]; ok {
		receiver.exports[filename][key] = ContextNode{result: value}
	} else {
		val := ContextNode{result: value}
		receiver.exports[filename] = Context{key: val}
	}
}

// Get returns the export context of the given filename
func (receiver *ExportStore) Get(filename string) Context {
	return receiver.exports[filename]
}
