package processor

import (
	"sync"

	"mettlach.codes/frizzy/file"
)

// Files can export variables for use by other files
// All assignment variables are exported

// ExportStorage represents a type that can store results
// into a context and return that context
type ExportStorage interface {
	Insert([]string, Result)
	GetContext() *Context
	GetFileContext(string) *Context
	GetNamespace() string
}

// ExportFileStore is a wrapper around ExportStore for a
// specific file
type ExportFileStore struct {
	Filepath string
}

func NewExportFileStore(filepath string) *ExportFileStore {
	store := &ExportFileStore{Filepath: filepath}
	store.InsertSpecialValues()

	return store
}

// InsertSpecialValues stores meta context values like a link to the file
func (receiver *ExportFileStore) InsertSpecialValues() {
	receiver.Insert([]string{"_href"},
		StringResult(file.GetRelativePathTo(receiver.Filepath)))
}

// Insert inserts the key, value pair of the export context
// represented by Filepath
func (receiver *ExportFileStore) Insert(contextKeys []string, value Result) {
	exportStore := GetExportStore()
	exportStore.Insert(receiver.Filepath, contextKeys, value)
}

// GetContext returns the export context associated with this
// ExportFileStore Filepath
func (receiver *ExportFileStore) GetContext() *Context {
	exportStore := GetExportStore()
	return exportStore.Get(receiver.Filepath)
}

// GetFileContext returns the context associated with the given filePath
func (receiver *ExportFileStore) GetFileContext(filePath string) *Context {
	exportStore := GetExportStore()
	return exportStore.Get(filePath)
}

// GetNamespace returns the filepath for this store
func (receiver *ExportFileStore) GetNamespace() string {
	return receiver.Filepath
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
func (receiver *ExportStore) Insert(filename string, contextKeys []string, value Result) {
	mut.Lock()
	defer mut.Unlock()

	if _, ok := receiver.exports[filename]; !ok {
		receiver.exports[filename] = &Context{}
	}

	receiver.exports[filename].Insert(contextKeys, value)
}

// Get returns the export context of the given filename
func (receiver *ExportStore) Get(filename string) *Context {
	return receiver.exports[filename]
}
