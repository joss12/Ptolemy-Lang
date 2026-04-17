package src

import (
	"fmt"
	"math"
)

// Eval is the main entry point - evaluates any AST node
func Eval(node Node, env *Environment) Object {
	switch node := node.(type) {

	// Program - evaluate all statements
	case *Program:
		return evalProgram(node, env)

	// Statements
	case *ExpressionStatement:
		result := Eval(node.Expression, env)
		return result

	case *BlockStatement:
		return evalBlockStatement(node, env)

	case *LetStatement:
		val := Eval(node.Value, env)
		if IsError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	case *ConstStatement:
		val := Eval(node.Value, env)
		if IsError(val) {
			return val
		}
		env.SetConstant(node.Name.Value, val)

	case *ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if IsError(val) {
			return val
		}
		return &ReturnValue{Value: val}

	case *IfStatement:
		return evalIfStatement(node, env)

	case *WhileStatement:
		return evalWhileStatement(node, env)

	case *ForStatement:
		return evalForStatement(node, env)

	// Expressions - Literals
	case *NumberLiteral:
		return &NumberObject{Value: node.Value}

	case *StringLiteral:
		return &StringObject{Value: node.Value}

	case *BooleanLiteral:
		return NativeBoolToBooleanObject(node.Value)

	case *NullLiteral:
		return NULL_OBJ_VAL

	case *ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && IsError(elements[0]) {
			return elements[0]
		}
		return &ArrayObject{Elements: elements}

	// Expressions - Operations
	case *PrefixExpression:
		right := Eval(node.Right, env)
		if IsError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right, node.Token.Line)

	case *InfixExpression:
		left := Eval(node.Left, env)
		if IsError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if IsError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right, node.Token.Line)

	case *AssignmentExpression:
		return evalAssignment(node, env)

	// Expressions - Variables
	case *Identifier:
		return evalIdentifier(node, env)

	// Expressions - Index
	case *IndexExpression:
		left := Eval(node.Left, env)
		if IsError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if IsError(index) {
			return index
		}
		return evalIndexExpression(left, index, node.Token.Line)

	// Expressions - Functions
	case *FunctionLiteral:
		fn := &FunctionObject{
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
		}
		if node.Name != "" {
			env.Set(node.Name, fn)
		}
		return fn

	case *CallExpression:
		function := Eval(node.Function, env)
		if IsError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && IsError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args, node.Token.Line)
	}

	return nil
}

// PROGRAM

func evalProgram(program *Program, env *Environment) Object {
	var result Object

	for _, stmt := range program.Statements {
		result = Eval(stmt, env)

		switch result := result.(type) {
		case *ReturnValue:
			return result.Value
		case *ErrorObject:
			return result
		}
	}

	return result
}

// BLOCK STATEMENT

func evalBlockStatement(block *BlockStatement, env *Environment) Object {
	var result Object

	for _, stmt := range block.Statements {
		result = Eval(stmt, env)

		if result != nil {
			rt := result.Type()
			if rt == RETURN_OBJ || rt == ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

// IF STATEMENT

func evalIfStatement(is *IfStatement, env *Environment) Object {
	condition := Eval(is.Condition, env)
	if IsError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(is.Consequence, env)
	} else if is.Alternative != nil {
		return Eval(is.Alternative, env)
	}

	return NULL_OBJ_VAL
}

// WHILE STATEMENT

func evalWhileStatement(ws *WhileStatement, env *Environment) Object {
	var result Object

	for {
		condition := Eval(ws.Condition, env)
		if IsError(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break
		}

		result = Eval(ws.Body, env)

		if result != nil {
			rt := result.Type()
			if rt == RETURN_OBJ || rt == ERROR_OBJ {
				return result
			}
		}
	}

	return NULL_OBJ_VAL
}

// FOR STATEMENT

func evalForStatement(fs *ForStatement, env *Environment) Object {
	forEnv := NewEnclosedEnvironment(env)

	if fs.Init != nil {
		result := Eval(fs.Init, forEnv)
		if IsError(result) {
			return result
		}
	}

	var result Object

	for {
		if fs.Condition != nil {
			condition := Eval(fs.Condition, forEnv)
			if IsError(condition) {
				return condition
			}
			if !isTruthy(condition) {
				break
			}
		}

		result = Eval(fs.Body, forEnv)
		if result != nil {
			rt := result.Type()
			if rt == RETURN_OBJ || rt == ERROR_OBJ {
				return result
			}
		}

		if fs.Increment != nil {
			incResult := Eval(fs.Increment, forEnv)
			if IsError(incResult) {
				return incResult
			}
		}
	}

	return NULL_OBJ_VAL
}

// IDENTIFIERS

func evalIdentifier(node *Identifier, env *Environment) Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return NewError(node.Token.Line, "undefined variable '%s'", node.Value)
}

// ASSIGNMENT

func evalAssignment(node *AssignmentExpression, env *Environment) Object {
	val := Eval(node.Value, env)
	if IsError(val) {
		return val
	}

	switch left := node.Left.(type) {
	case *Identifier:
		if env.IsConstant(left.Value) {
			return NewError(node.Token.Line, "cannot reassign constant '%s'", left.Value)
		}
		if _, ok := env.Update(left.Value, val); ok {
			return val
		}
		return NewError(node.Token.Line, "undefined variable '%s'", left.Value)

	case *IndexExpression:
		obj := Eval(left.Left, env)
		if IsError(obj) {
			return obj
		}

		idx := Eval(left.Index, env)
		if IsError(idx) {
			return idx
		}

		arr, ok := obj.(*ArrayObject)
		if !ok {
			return NewError(node.Token.Line, "index assignment on non-array")
		}

		index := int(idx.(*NumberObject).Value)
		if index < 0 || index >= len(arr.Elements) {
			return NewError(node.Token.Line, "index out of bounds: %d", index)
		}

		arr.Elements[index] = val
		return val
	}

	return NewError(node.Token.Line, "invalid assignment target")
}

// PREFIX EXPRESSIONS

func evalPrefixExpression(operator string, right Object, line int) Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right, line)
	default:
		return NewError(line, "unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right Object) Object {
	switch right {
	case TRUE_OBJ:
		return FALSE_OBJ
	case FALSE_OBJ:
		return TRUE_OBJ
	case NULL_OBJ_VAL:
		return TRUE_OBJ
	default:
		return FALSE_OBJ
	}
}

func evalMinusPrefixOperatorExpression(right Object, line int) Object {
	if right.Type() != NUMBER_OBJ {
		return NewError(line, "unknown operator: -%s", right.Type())
	}
	value := right.(*NumberObject).Value
	return &NumberObject{Value: -value}
}

// INFIX EXPRESSIONS

func evalInfixExpression(operator string, left, right Object, line int) Object {
	switch {
	case left.Type() == NUMBER_OBJ && right.Type() == NUMBER_OBJ:
		return evalNumberInfixExpression(operator, left, right, line)
	case left.Type() == STRING_OBJ && right.Type() == STRING_OBJ:
		return evalStringInfixExpression(operator, left, right, line)

	case left.Type() == STRING_OBJ && right.Type() == NUMBER_OBJ && operator == "+":
		leftVal := left.(*StringObject).Value
		rightVal := right.(*NumberObject).Inspect()
		return &StringObject{Value: leftVal + rightVal}

	case left.Type() == NUMBER_OBJ && right.Type() == STRING_OBJ && operator == "+":
		leftVal := left.(*NumberObject).Inspect()
		rightVal := right.(*StringObject).Value
		return &StringObject{Value: leftVal + rightVal}

	case operator == "==":
		return NativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return NativeBoolToBooleanObject(left != right)
	case operator == "&&":
		return evalLogicalAnd(left, right)
	case operator == "||":
		return evalLogicalOr(left, right)
	case left.Type() != right.Type():
		return NewError(line, "type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	default:
		return NewError(line, "unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalLogicalAnd(left, right Object) Object {

	//if left is falsy, return false
	if !isTruthy(left) {
		return FALSE_OBJ
	}

	//if right is truthy, return right's truthiness
	if isTruthy(right) {
		return TRUE_OBJ
	}

	return FALSE_OBJ
}

func evalLogicalOr(left, right Object) Object {
	//if left is truthy, return true
	if isTruthy(left) {
		return TRUE_OBJ
	}

	//If left is falsy, return right's truthiness
	if isTruthy(right) {
		return TRUE_OBJ
	}

	return FALSE_OBJ
}

func evalNumberInfixExpression(operator string, left, right Object, line int) Object {
	leftVal := left.(*NumberObject).Value
	rightVal := right.(*NumberObject).Value

	switch operator {
	case "+":
		return &NumberObject{Value: leftVal + rightVal}
	case "-":
		return &NumberObject{Value: leftVal - rightVal}
	case "*":
		return &NumberObject{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return NewError(line, "division by zero")
		}
		return &NumberObject{Value: leftVal / rightVal}
	case "%":
		if rightVal == 0 {
			return NewError(line, "division by zero")
		}
		return &NumberObject{Value: math.Mod(leftVal, rightVal)}
	case "<":
		return NativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return NativeBoolToBooleanObject(leftVal > rightVal)
	case "<=":
		return NativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return NativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return NativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return NativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return NewError(line, "unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right Object, line int) Object {
	leftVal := left.(*StringObject).Value
	rightVal := right.(*StringObject).Value

	switch operator {
	case "+":
		return &StringObject{Value: leftVal + rightVal}
	case "==":
		return NativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return NativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return NewError(line, "unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

// INDEX EXPRESSIONS

func evalIndexExpression(left, index Object, line int) Object {
	switch {
	case left.Type() == ARRAY_OBJ && index.Type() == NUMBER_OBJ:
		return evalArrayIndexExpression(left, index, line)
	default:
		return NewError(line, "index operator not supported for %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index Object, line int) Object {
	arrayObject := array.(*ArrayObject)
	idx := int(index.(*NumberObject).Value)
	max := len(arrayObject.Elements) - 1

	if idx < 0 || idx > max {
		return NewError(line, "index out of bounds: %d", idx)
	}

	return arrayObject.Elements[idx]
}

// FUNCTIONS

func evalExpressions(exps []Expression, env *Environment) []Object {
	var result []Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if IsError(evaluated) {
			return []Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn Object, args []Object, line int) Object {
	switch fn := fn.(type) {
	case *FunctionObject:
		extendedEnv, err := extendFunctionEnv(fn, args, line)
		if err != nil {
			return err
		}
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *BuiltinObject:
		return fn.FN(args...)

	default:
		return NewError(line, "not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *FunctionObject, args []Object, line int) (*Environment, *ErrorObject) {
	env := NewEnclosedEnvironment(fn.Env)

	if len(args) != len(fn.Parameters) {
		return nil, NewError(line, "wrong number of arguments: expected %d, got %d",
			len(fn.Parameters), len(args))
	}

	for i, param := range fn.Parameters {
		env.Set(param.Value, args[i])
	}

	return env, nil
}

func unwrapReturnValue(obj Object) Object {
	if returnValue, ok := obj.(*ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

// HELPERS

func isTruthy(obj Object) bool {
	switch obj {
	case NULL_OBJ_VAL:
		return false
	case TRUE_OBJ:
		return true
	case FALSE_OBJ:
		return false
	default:
		return true
	}
}

// BUILT-IN FUNCTIONS

var builtins = map[string]*BuiltinObject{
	"print": {
		FN: func(args ...Object) Object {
			parts := []string{}
			for _, arg := range args {
				parts = append(parts, arg.Inspect())
			}
			fmt.Println(joinStrings(parts, " "))
			return NULL_OBJ_VAL
		},
	},

	"len": {
		FN: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError(0, "wrong number of arguments: len expects 1, got %d", len(args))
			}
			switch arg := args[0].(type) {
			case *ArrayObject:
				return &NumberObject{Value: float64(len(arg.Elements))}
			case *StringObject:
				return &NumberObject{Value: float64(len(arg.Value))}
			default:
				return NewError(0, "argument to len not supported, got %s", args[0].Type())
			}
		},
	},

	"type": {
		FN: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError(0, "wrong number of arguments: type expects 1, got %d", len(args))
			}
			return &StringObject{Value: string(args[0].Type())}
		},
	},

	"push": {
		FN: func(args ...Object) Object {
			if len(args) != 2 {
				return NewError(0, "wrong number of arguments: push expects 2, got %d", len(args))
			}
			arr, ok := args[0].(*ArrayObject)
			if !ok {
				return NewError(0, "first argument to push must be an array, got %s", args[0].Type())
			}
			arr.Elements = append(arr.Elements, args[1])
			return NULL_OBJ_VAL
		},
	},

	"pop": {
		FN: func(args ...Object) Object {
			if len(args) != 1 {
				return NewError(0, "wrong number of arguments: pop expects 1, got %d", len(args))
			}
			arr, ok := args[0].(*ArrayObject)
			if !ok {
				return NewError(0, "argument to pop must be an array, got %s", args[0].Type())
			}
			length := len(arr.Elements)
			if length == 0 {
				return NewError(0, "cannot pop from empty array")
			}
			last := arr.Elements[length-1]
			arr.Elements = arr.Elements[:length-1]
			return last
		},
	},
}

// joinStrings joins string slice with separator
func joinStrings(parts []string, sep string) string {
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += sep
		}
		result += p
	}
	return result
}
