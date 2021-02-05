package parser

type Result interface {
	GetResult() interface{}
}

type AddableResult interface {
	Add(right Result) (Result, error)
}

type MultipliableResult interface {
	Multiply(right Result) (Result, error)
}
