package pipeline

import "fmt"

type StandardError struct {
	Filename string
	Message  string
}

func (receiver *StandardError) Error() string {
	return fmt.Sprintf("%s: %s", receiver.Filename, receiver.Message)
}
