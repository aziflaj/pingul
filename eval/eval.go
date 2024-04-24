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
