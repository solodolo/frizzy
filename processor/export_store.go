// Files can export variables for use by other files
// All assignment variables are exported

// exportStore is a singleton to read and write export vars
type exportStore struct {
	exports map[string]Context
}

var store *exportStore

// GetExportStore returns the exportStore singleton
func GetExportStore() *exportStore {
	if store == nil {
		store = &exportStore{}
	}

	return store
}

// Inserts inserts the key, value pair of the export context
// represented by filename
func (receiver *exportStore) Insert(filename, key string, value Result) {
	if ctx, ok := receiver.exports[filename]; ok {
		receiver.exports[filename][key] = value
	} else {
		receiver.exports[filename] = Context{key: value}
	}
}

// ContextFor returns the export context of the given filename
func (receiver *exportStore) ContextFor(filename string) Context {
	return receiver.exports[filename]
}