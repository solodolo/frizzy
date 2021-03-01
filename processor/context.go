package processor

import (
	"sort"
	"strings"
)

// ContextNode hosts either a Result or another
// Context level
type ContextNode struct {
	result Result
	child  *Context
}

// HasResult returns true if this node contains a result
func (receiver *ContextNode) HasResult() bool {
	return receiver.result != nil
}

// HasContext returns true if this node contains a
// nested context
func (receiver *ContextNode) HasContext() bool {
	return receiver.child != nil
}

// At returns the ContextNode stored under key or false
// if key does not exist
func (receiver *ContextNode) At(key string) (*ContextNode, bool) {
	if receiver.HasContext() {
		val, ok := (*receiver.child)[key]
		return val, ok
	}
	return nil, false
}

// Context is a recursive key-value store
// for storing Result types
type Context map[string]*ContextNode

// Merge adds the keys and values from other into receiver
// Matching keys in receiver will be overwritten
func (receiver *Context) Merge(other *Context) *Context {
	merged := &Context{}
	for k, v := range *receiver {
		(*merged)[k] = v
	}

	for k, v := range *other {
		(*merged)[k] = v
	}
	return merged
}

// Keys returns the sorted keys of receiver
func (receiver *Context) Keys() []string {
	keys := make([]string, 0, len(*receiver))
	for key := range *receiver {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	return keys
}

// Values converts a nested context into an array of contexts
// reducing its level by one
func (receiver *Context) Values() []*Context {
	keys := receiver.Keys()
	values := make([]*Context, 0, len(*receiver))
	for _, key := range keys {
		values = append(values, (*receiver)[key].child)
	}

	return values
}

// At returns the ContextNode stored under key or false
// if key does not exist
func (receiver *Context) At(key string) (*ContextNode, bool) {
	return receiver.AtNested(strings.Split(key, "."))
}

// AtNested iterates through keys, looking up nested context
// levels, and returns the last node found
func (receiver *Context) AtNested(keys []string) (*ContextNode, bool) {
	current := &ContextNode{child: receiver}
	for _, key := range keys {
		contextNode, exists := current.At(key)

		if !exists {
			return nil, false
		}

		current = contextNode
	}
	return current, true
}

// Insert iterates through keys, adding or looking up nested
// context levels, then inserting the result at the last key
func (receiver *Context) Insert(keys []string, value Result) {
	current := &ContextNode{child: receiver}
	for _, key := range keys {
		if at, ok := current.At(key); ok {
			current = at
		} else {
			next := &ContextNode{}
			if current.child == nil {
				current.child = &Context{key: next}
			} else {
				(*current.child)[key] = next
			}
			current = next
		}
	}
	current.result = value
}
