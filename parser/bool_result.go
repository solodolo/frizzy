package parser

type BoolResult bool

func (result BoolResult) GetResult() interface{} {
	return result
}
