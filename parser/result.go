package parser

type Result interface {
	GetResult() interface{}
}

type IntResult int

func (result IntResult) GetResult() interface{} {
	return result
}

type StringResult string

func (result StringResult) GetResult() interface{} {
	return result
}
