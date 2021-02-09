package processor

// Context is the set of methods avaialble for storing parsing
// context data
type Context map[string]Result

// Merge adds the keys and values from other into receiver
// Matching keys in receiver will be overwritten
func (receiver *Context) Merge(other *Context) {
	for k, v := range *other {
		(*receiver)[k] = v
	}
}
