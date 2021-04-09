package processor

// FunctionModule represents a type that can take a
// function name and result arguments and returns a result
type FunctionModule interface {
	CallFunction(string, ...Result) (Result, bool)
	registerFunc(string, func(...Result) Result)
}

// BuiltinFunctionModule acts as an interface between callers
// and builtin Frizzy functions
type BuiltinFunctionModule map[string]func(...Result) Result

func (receiver *BuiltinFunctionModule) registerFunc(funcName string, function func(...Result) Result) {
	(*receiver)[funcName] = function
}

// CallFunction calls the function with the given name and passes in
// the provided arguments
func (receiver *BuiltinFunctionModule) CallFunction(funcName string, funcArgs ...Result) (Result, bool) {
	if function, ok := (*receiver)[funcName]; ok {
		return function(funcArgs...), true
	}

	return nil, false
}

// NewBuiltinFunctionModule creates a BuiltinFunctionModule object
// populated with a mapping for each of the builtin functions
func NewBuiltinFunctionModule() *BuiltinFunctionModule {
	module := &BuiltinFunctionModule{}
	module.registerFunc("print", PrintRaw)
	module.registerFunc("paginate", PaginateRaw)
	module.registerFunc("template", TemplateRaw)

	return module
}
