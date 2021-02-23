package processor

import (
	"sync"
)

// Files can export variables for use by other files
// All assignment variables are exported

// ExportStorage represents a type that can store results
// into a context and return that context
type ExportStorage interface {
	Insert(string, Result)
	GetContext() *Context
	GetFileContext(string) *Context
}

// ExportFileStore is a wrapper around ExportStore for a
// specific file
type ExportFileStore struct {
	filePath string
}

// Insert inserts the key, value pair of the export context
// represented by filePath
func (receiver *ExportFileStore) Insert(contextKey string, value Result) {
	exportStore := GetExportStore()
	exportStore.Insert(receiver.filePath, contextKey, value)
}

// GetContext returns the export context associated with this
// ExportFileStore filePath
func (receiver *ExportFileStore) GetContext() *Context {
	exportStore := GetExportStore()
	return exportStore.Get(receiver.filePath)
}

// GetFileContext returns the context associated with the given filePath
func (receiver *ExportFileStore) GetFileContext(filePath string) *Context {
	exportStore := GetExportStore()
	return exportStore.Get(filePath)
}

// ExportStore is a singleton to read and write export vars
type ExportStore struct {
	exports map[string]*Context
}

var once sync.Once
var store *ExportStore

func createStore() {
	if store == nil {
		store = &ExportStore{exports: make(map[string]*Context)}
	}
}

// GetExportStore returns the exportStore singleton
func GetExportStore() *ExportStore {
	once.Do(createStore)
	return store
}

var mut sync.Mutex

// Insert inserts the key, value pair of the export context
// represented by filename
func (receiver *ExportStore) Insert(filename, contextKey string, value Result) {
	mut.Lock()
	defer mut.Unlock()

	if _, ok := receiver.exports[filename]; ok {
		context := receiver.exports[filename]
		(*context)[contextKey] = ContextNode{result: value}
	} else {
		val := ContextNode{result: value}
		receiver.exports[filename] = &Context{contextKey: val}
	}
}

// Get returns the export context of the given filename
func (receiver *ExportStore) Get(filename string) *Context {
	return receiver.exports[filename]
}
