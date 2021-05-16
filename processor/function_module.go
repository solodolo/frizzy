package processor

import "fmt"

// FunctionModule represents a type that can take a
// function name and result arguments and returns a result
type FunctionModule interface {
	CallFunction(string, ...Result) (Result, error)
	registerFunc(string, func(...Result) (Result, error))
}

// BuiltinFunctionModule acts as an interface between callers
// and builtin Frizzy functions
type BuiltinFunctionModule map[string]func(...Result) (Result, error)

func (receiver *BuiltinFunctionModule) registerFunc(funcName string, function func(...Result) (Result, error)) {
	(*receiver)[funcName] = function
}

// CallFunction calls the function with the given name and passes in
// the provided arguments
func (receiver *BuiltinFunctionModule) CallFunction(funcName string, funcArgs ...Result) (Result, error) {
	if function, ok := (*receiver)[funcName]; ok {
		return function(funcArgs...)
	}

	return nil, fmt.Errorf("function %s is not registered", funcName)
}

// NewBuiltinFunctionModule creates a BuiltinFunctionModule object
// populated with a mapping for each of the builtin functions
func NewBuiltinFunctionModule() *BuiltinFunctionModule {
	module := &BuiltinFunctionModule{}
	module.registerFunc("paginate", PaginateRaw)
	module.registerFunc("pagesBefore", PagesBeforeRaw)
	module.registerFunc("template", TemplateRaw)

	return module
}
