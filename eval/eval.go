package eval

import (
	"github.com/aziflaj/pingul/ast"
	"github.com/aziflaj/pingul/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return &object.Boolean{Value: node.Value}

	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)

	default:
		return &object.Nil{}
	}
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	if operator == "not" {
		if right.Type() == object.BOOL {
			return &object.Boolean{Value: !right.(*object.Boolean).Value}
		} else if right.Type() == object.INT {
			return &object.Boolean{Value: right.(*object.Integer).Value == 0}
		}
	}

	if operator == "-" && right.Type() == object.INT {
		return &object.Integer{Value: -right.(*object.Integer).Value}
	}

	return &object.Nil{}
}

// if left value is bool, all is bool
// if left value is int, all is int
func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	if left.Type() == object.BOOL {
		// transform right value to boolean
		if right.Type() == object.INT {
			right = &object.Boolean{Value: right.(*object.Integer).Value != 0}
		}

		return evalBooleanInfixExpression(operator, left, right)
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

	return &object.Nil{}
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