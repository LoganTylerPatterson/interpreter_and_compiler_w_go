package evaluator

import (
	"interpreter/ast"
	"interpreter/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.BooleanExpression:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Value)
		return evaluatePrefixExpression(node.Op, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evaluateInfixExpression(left, right, node.Op)
	case *ast.Program:
		return evaluateStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	}
	return nil
}

func evaluateInfixExpression(left object.Object, right object.Object, op string) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evaluateIntegerInfixExpression(left, right, op)
	default:
		return NULL
	}
}

func evaluateIntegerInfixExpression(left object.Object, right object.Object, op string) object.Object {
	lVal := left.(*object.Integer).Value
	rVal := right.(*object.Integer).Value
	switch op {
	case "-":
		return &object.Integer{Value: lVal - rVal}
	case "+":
		return &object.Integer{Value: lVal + rVal}
	case "/":
		return &object.Integer{Value: lVal / rVal}
	case "*":
		return &object.Integer{Value: lVal * rVal}
	default:
		return NULL
	}
}

func evaluatePrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evaluateNotExpression(right)
	case "-":
		return evaluateMinusPrefixExpression(right)
	default:
		return NULL
	}
}

func evaluateMinusPrefixExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return NULL
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evaluateNotExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func nativeBoolToBooleanObject(input bool) object.Object {
	if input {
		return TRUE
	}
	return FALSE
}

func evaluateStatements(stmts []ast.Statement) object.Object {
	var result object.Object
	for _, statement := range stmts {
		result = Eval(statement)
	}
	return result
}
