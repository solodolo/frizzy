package processor

// ContextNode hosts either a Result or another
// Context level
type ContextNode struct {
	result Result
	child  *Context
}

// HasResult returns true if this node contains a result
func (receiver ContextNode) HasResult() bool {
	return receiver.result != nil
}

// HasContext returns true if this node contains a
// nested context
func (receiver ContextNode) HasContext() bool {
	return receiver.child != nil
}

// At returns the ContextNode stored under key or false
// if key does not exist
func (receiver ContextNode) At(key string) (ContextNode, bool) {
	if receiver.HasContext() {
		val, ok := (*receiver.child)[key]
		return val, ok
	}
	return ContextNode{}, false
}

// Context is a recursive key-value store
// for storing Result types
type Context map[string]ContextNode

// Merge adds the keys and values from other into receiver
// Matching keys in receiver will be overwritten
func (receiver *Context) Merge(other *Context) {
	for k, v := range *other {
		(*receiver)[k] = v
	}
}

// At returns the ContextNode stored under key or false
// if key does not exist
func (receiver *Context) At(key string) (ContextNode, bool) {
	val, ok := (*receiver)[key]
	return val, ok
}
