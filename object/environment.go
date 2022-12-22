package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()

	env.outer = outer

	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)

	return &Environment{store: s}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]

	// 「私の手元にないんで、店長にきいてみますね」というノリ
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
		return obj, ok
	}

	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val

	return val
}

func (e *Environment) GetStoredEnv(name string) *Environment {
	_, ok := e.store[name]

	if ok {
		return e
	}

	if e.outer != nil {
		return e.outer.GetStoredEnv(name)
	}

	return nil
}
