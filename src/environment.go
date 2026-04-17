package src

type Environment struct {
	store     map[string]Object
	constants map[string]bool
	outer     *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	c := make(map[string]bool)
	return &Environment{store: s, constants: c, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func (e *Environment) SetConstant(name string, val Object) Object {
	e.store[name] = val
	e.constants[name] = true
	return val
}

func (e *Environment) IsConstant(name string) bool {
	if _, ok := e.constants[name]; ok {
		return true
	}
	if e.outer != nil {
		return e.outer.IsConstant(name)
	}
	return false
}

func (e *Environment) Update(name string, val Object) (Object, bool) {
	_, ok := e.store[name]
	if ok {
		e.store[name] = val
		return val, true
	}
	if e.outer != nil {
		return e.outer.Update(name, val)
	}
	return nil, false
}
