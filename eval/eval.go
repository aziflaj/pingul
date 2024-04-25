package eval

import (
	"github.com/aziflaj/pingul/ast"
	"github.com/aziflaj/pingul/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node)
	case *ast.BlockStatement:
		return evalBlock(node)

	case *ast.ReturnStatement:
		return &object.Return{Value: Eval(node.ReturnValue)}

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return &object.Boolean{Value: node.Value}

	case *ast.Nil:
		return &object.Nil{}

	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)

	case *ast.IfExpression:
		cond := Eval(node.Condition)
		return evalIfExpression(cond.IsTruthy(), node.Consequence, node.Alternative)

	default:
		return &object.Nil{}
	}
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, stmt := range program.Statements {
		result = Eval(stmt)

		if val, ok := result.(*object.Return); ok {
			return val.Value
		}
	}

	return result
}

func evalBlock(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, stmt := range block.Statements {
		result = Eval(stmt)

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

func evalIfExpression(cond bool, consequence *ast.BlockStatement, alternative *ast.BlockStatement) object.Object {
	if cond {
		return Eval(consequence)
	}

	if alternative != nil {
		return Eval(alternative)
	}

	return &object.Nil{}
}
