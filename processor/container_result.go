package processor

import "fmt"

// ContainerResult is a result that holds a context
// useful for results that contain more results like
// in a for loop over a context range
type ContainerResult struct {
	context *Context
}

// GetResult returns the context of this container
func (receiver ContainerResult) GetResult() interface{} {
	return receiver.context
}

func (receiver ContainerResult) String() string {
	return fmt.Sprintf("%T", receiver)
}
