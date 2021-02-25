package processor

// FunctionModule represents a type that can take a
// function name and result arguments and returns a result
type FunctionModule interface {
	CallFunction(string, ...Result) (Result, bool)
	registerFunc(string, func(...Result) Result)
}

type BuiltinFunctionModule map[string]func(...Result) Result

func (receiver *BuiltinFunctionModule) registerFunc(funcName string, function func(...Result) Result) {
	(*receiver)[funcName] = function
}

func (receiver *BuiltinFunctionModule) CallFunction(funcName string, funcArgs ...Result) (Result, bool) {
	if function, ok := (*receiver)[funcName]; ok {
		return function(funcArgs...), true
	}

	return nil, false
}

func NewBuiltinFunctionModule() BuiltinFunctionModule {
	module := BuiltinFunctionModule{}
	module.registerFunc("print", PrintRaw)
	module.registerFunc("paginate", PaginateRaw)

	return module
}
