package parser

// Context is the set of methods avaialble for storing parsing
// context data
type Context interface {
	GetStringValue(key string) string
	GetIntValue(key string) int

	SetStringValue(key, value string)
	SetIntValue(key string, value int)
}

// ParseContext is a key-value store to store parsing context data
type ParseContext struct {
	stringMap map[string]string
	intMap    map[string]int
}

// GetStringValue returns the string value of the given key
// and a bool indicating if found
func (context *ParseContext) GetStringValue(key string) (val string, found bool) {
	val, found = context.stringMap[key]
	return val, found
}

// GetIntValue returns the int value of the given key
// and a bool indicating if found
func (context *ParseContext) GetIntValue(key string) (val int, found bool) {
	val, found = context.intMap[key]
	return val, found
}
