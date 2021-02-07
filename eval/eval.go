package eval

import (
	"github.com/rumpl/monkey-lang/ast"
	"github.com/rumpl/monkey-lang/object"
)

var (
	Null  = &object.Null{}
	True  = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.BlockStatement:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node)
	}

	return nil
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return True
	}

	return False
}

func evalStatements(statements []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range statements {
		result = Eval(stmt)
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return Null
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case True:
		return False
	case False:
		return True
	case Null:
		return True
	default:
		return False
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.IntegerObj {
		return Null
	}
	integer := right.(*object.Integer)
	return &object.Integer{Value: -1 * integer.Value}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.IntegerObj && right.Type() == object.IntegerObj:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return Null
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	li := left.(*object.Integer)
	ri := right.(*object.Integer)

	switch operator {
	case "+":
		return &object.Integer{Value: li.Value + ri.Value}
	case "-":
		return &object.Integer{Value: li.Value - ri.Value}
	case "*":
		return &object.Integer{Value: li.Value * ri.Value}
	case "/":
		return &object.Integer{Value: li.Value / ri.Value}
	case "<":
		return nativeBoolToBooleanObject(li.Value < ri.Value)
	case ">":
		return nativeBoolToBooleanObject(li.Value > ri.Value)
	case "!=":
		return nativeBoolToBooleanObject(li.Value != ri.Value)
	case "==":
		return nativeBoolToBooleanObject(li.Value == ri.Value)
	default:
		return Null
	}
}

func evalIfExpression(obj *ast.IfExpression) object.Object {
	condition := Eval(obj.Condition)

	if isTruthy(condition) {
		return Eval(obj.Consequence)
	} else if obj.Alternative != nil {
		return Eval(obj.Alternative)
	}
	return Null
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case Null:
		return false
	case True:
		return true
	case False:
		return false
	default:
		return true
	}
}
