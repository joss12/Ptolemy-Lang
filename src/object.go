package src

import (
	"fmt"
	"strings"
)

type ObjectType string

const (
	NUMBER_OBJ    ObjectType = "NUMBER"
	STRING_OBJ    ObjectType = "STRING"
	BOOLEAN_OBJ   ObjectType = "BOOLEAN"
	NULL_OBJ      ObjectType = "NULL"
	UNDEFINED_OBJ ObjectType = "UNDEFINED"
	ARRAY_OBJ     ObjectType = "ARRAY"
	FUNCTION_OBJ  ObjectType = "FUNCTION"
	BUILTIN_OBJ   ObjectType = "BUILTIN"
	RETURN_OBJ    ObjectType = "RETURN_VALUE"
	ERROR_OBJ     ObjectType = "ERROR"
)

// Object is the base interface for all runtime values

type Object interface {
	Type() ObjectType
	Inspect() string
}

// Number
type NumberObject struct {
	Value float64
}

func (n *NumberObject) Type() ObjectType { return NUMBER_OBJ }
func (n *NumberObject) Inspect() string {
	//Print  intergers without decimal point
	if n.Value == float64(int64(n.Value)) {
		return fmt.Sprintf("%d", int64(n.Value))
	}
	return fmt.Sprintf("%g", n.Value)
}

// String
type StringObject struct {
	Value string
}

func (s *StringObject) Type() ObjectType { return STRING_OBJ }
func (s *StringObject) Inspect() string  { return s.Value }

// BOOLEAN
type BooleanObject struct {
	Value bool
}

func (b *BooleanObject) Type() ObjectType { return BOOLEAN_OBJ }
func (b *BooleanObject) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

// NULL
type NullObject struct{}

func (n *NullObject) Type() ObjectType { return NULL_OBJ }
func (n *NullObject) Inspect() string  { return "null" }

// ARRAY
type ArrayObject struct {
	Elements []Object
}

func (a *ArrayObject) Type() ObjectType { return ARRAY_OBJ }
func (a *ArrayObject) Inspect() string {
	elements := []string{}
	for _, el := range a.Elements {
		elements = append(elements, el.Inspect())
	}
	return "[" + strings.Join(elements, ", ") + "]"
}

// FUNCTION
type FunctionObject struct {
	Parameters []*Identifier
	Body       *BlockStatement
	Env        *Environment //the captured closer environment
}

func (f *FunctionObject) Type() ObjectType { return FUNCTION_OBJ }
func (f *FunctionObject) Inspect() string {
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	return fmt.Sprintf("fn(%s){...}", strings.Join(params, ", "))
}

// BUILTIN
type BuiltinFunction func(args ...Object) Object

type BuiltinObject struct {
	FN BuiltinFunction
}

func (b *BuiltinObject) Type() ObjectType { return BUILTIN_OBJ }
func (b *BuiltinObject) Inspect() string  { return "builtin function" }

// RETURN VALUE (Wrapper)
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

// ERROR
type ErrorObject struct {
	Message string
	Line    int
}

func (e *ErrorObject) Type() ObjectType { return ERROR_OBJ }
func (e *ErrorObject) Inspect() string {
	if e.Line > 0 {
		return fmt.Sprintf("Returntime Error at line %d: %s", e.Line, e.Message)
	}
	return fmt.Sprintf("Runtime Error: %s", e.Message)
}

// SINGLETON VALUES
var (
	TRUE_OBJ     = &BooleanObject{Value: true}
	FALSE_OBJ    = &BooleanObject{Value: false}
	NULL_OBJ_VAL = &NullObject{}
)

// NativeBoolToBooleanObject converts a Go bool to our Boolean object
// Uses singletons to avoid creating new objects every time
func NativeBoolToBooleanObject(input bool) *BooleanObject {
	if input {
		return TRUE_OBJ
	}
	return FALSE_OBJ
}

// NewError creates a new error object
func NewError(line int, format string, a ...interface{}) *ErrorObject {
	return &ErrorObject{
		Message: fmt.Sprintf(format, a...),
		Line:    line,
	}
}

// IsError checks  if an object is an error
func IsError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR_OBJ
	}
	return false
}
