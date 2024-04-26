package eval

import (
	"github.com/aziflaj/pingul/ast"
	"github.com/aziflaj/pingul/object"
)

var scope = object.NewScope()

func Eval(scope *object.Scope, node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(scope, node)

	case *ast.BlockStatement:
		return evalBlock(scope, node)

	case *ast.ReturnStatement:
		return &object.Return{Value: Eval(scope, node.ReturnValue)}

	case *ast.ExpressionStatement:
		return Eval(scope, node.Expression)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return &object.Boolean{Value: node.Value}

	case *ast.String:
		return &object.String{Value: node.Value}

	case *ast.Nil:
		return &object.Nil{}

	case *ast.FuncExpression:
		return &object.Func{Params: node.Params, Body: node.Body}

	case *ast.CallExpression:
		// eval args, left to right
		args := make([]object.Object, len(node.Arguments))
		for i, arg := range node.Arguments {
			args[i] = Eval(scope, arg)
		}

		fun := Eval(scope, node.Function)

		return applyFunction(scope, fun, args)

	case *ast.VarStatement:
		val := Eval(scope, node.Value)
		scope.Set(node.Name.String(), val)
		return val

	case *ast.Identifier:
		// try the intrinsic functions first
		ident, ok := object.IntrinsicFuncs[node.String()]

		if ok {
			return ident
		}

		return scope.Get(node.String())

	case *ast.PrefixExpression:
		right := Eval(scope, node.Right)
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(scope, node.Left)
		right := Eval(scope, node.Right)
		return evalInfixExpression(node.Operator, left, right)

	case *ast.IfExpression:
		cond := Eval(scope, node.Condition)
		return evalIfExpression(scope, cond.IsTruthy(), node.Consequence, node.Alternative)

	default:
		return &object.Nil{}
	}
}

func evalProgram(scope *object.Scope, program *ast.Program) object.Object {
	var result object.Object

	for _, stmt := range program.Statements {
		result = Eval(scope, stmt)

		if val, ok := result.(*object.Return); ok {
			return val.Value
		}
	}

	return result
}

func evalBlock(scope *object.Scope, block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, stmt := range block.Statements {
		result = Eval(scope, stmt)

		if result.Type() == object.RETURN {
			return result
		}
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	if operator == "not" {
		return &object.Boolean{Value: !right.IsTruthy()}
	}

	if operator == "-" && right.Type() == object.INT {
		return &object.Integer{Value: -right.(*object.Integer).Value}
	}

	return &object.Nil{}
}

// if left value is bool, all is bool
// if left value is int, all is int
func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	if left.Type() == object.STRING && right.Type() == object.STRING {
		if operator == "+" {
			return &object.String{Value: append(left.(*object.String).Value, right.(*object.String).Value...)}
		}
	}

	if left.Type() == object.INT {
		// transform right value to integer
		if right.Type() == object.BOOL {
			var btoi int64
			if right.(*object.Boolean).Value {
				btoi = 1
			}

			right = &object.Integer{Value: btoi}
		}

		return evalIntegerInfixExpression(operator, left, right)
	}

	// left value is bool or nil
	leftBool := &object.Boolean{Value: left.IsTruthy()}
	rightBool := &object.Boolean{Value: right.IsTruthy()}

	return evalBooleanInfixExpression(operator, leftBool, rightBool)

	// return &object.Nil{}
}

func evalBooleanInfixExpression(
	operator string,
	left object.Object, right object.Object,
) object.Object {
	leftBool := left.(*object.Boolean).Value
	rightBool := right.(*object.Boolean).Value

	if operator == "==" {
		return &object.Boolean{Value: leftBool == rightBool}
	}

	if operator == "!=" {
		return &object.Boolean{Value: leftBool != rightBool}
	}

	return &object.Nil{}
}

func evalIntegerInfixExpression(
	operator string,
	left object.Object, right object.Object,
) object.Object {
	leftInt := left.(*object.Integer).Value
	rightInt := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftInt + rightInt}
	case "-":
		return &object.Integer{Value: leftInt - rightInt}
	case "*":
		return &object.Integer{Value: leftInt * rightInt}
	case "/":
		return &object.Integer{Value: leftInt / rightInt}
	case "%":
		return &object.Integer{Value: leftInt % rightInt}

	case "==":
		return &object.Boolean{Value: leftInt == rightInt}
	case "!=":
		return &object.Boolean{Value: leftInt != rightInt}

	case ">":
		return &object.Boolean{Value: leftInt > rightInt}
	case "<":
		return &object.Boolean{Value: leftInt < rightInt}
	case ">=":
		return &object.Boolean{Value: leftInt >= rightInt}
	case "<=":
		return &object.Boolean{Value: leftInt <= rightInt}
	}

	return &object.Nil{}
}

func evalIfExpression(scope *object.Scope, cond bool, consequence *ast.BlockStatement, alternative *ast.BlockStatement) object.Object {
	if cond {
		return Eval(scope, consequence)
	}

	if alternative != nil {
		return Eval(scope, alternative)
	}

	return &object.Nil{}
}

func applyFunction(scope *object.Scope, fun object.Object, args []object.Object) object.Object {
	if fun.Type() == object.INTRINSIC_FUNC {
		return fun.(object.IntrinsicFunc)(args...)
	}

	if fun.Type() != object.FUNC {
		return &object.Nil{}
	}

	function := fun.(*object.Func)
	localScope := object.NewLocalScope(scope)

	for i, param := range function.Params {
		p := param.String()
		localScope.Set(p, args[i])
	}

	result := Eval(localScope, function.Body)

	if result.Type() == object.RETURN {
		return result.(*object.Return).Value
	}

	return result
}
