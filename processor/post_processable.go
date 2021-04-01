package processor

// PostProcessable is an object that can take a result and return
// another result
type PostProcessable interface {
	Call(Result) Result
}
