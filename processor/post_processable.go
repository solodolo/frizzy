package processor

// PostProcessable is an object that can take a Result chan and return
// another Result chan
type PostProcessable interface {
	Call(<-chan Result) <-chan Result
}
