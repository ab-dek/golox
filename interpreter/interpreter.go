package interpreter

import (
	e "github.com/ab-dek/golox/ast"
	t "github.com/ab-dek/golox/token"
)

type interpreter struct{}

func (i *interpreter) evaluate(expr e.Expr) any {
	return expr.Accept(i)
}

func (i *interpreter) visitBinary(expr *e.Binary) any {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.TokenType {
	case t.BANG_EQUAL:
		return !(left == right)
	case t.EQUAL_EQUAL:
		return left == right
	case t.GREATER:
		return left.(float64) > right.(float64)
	case t.GREATER_EQUAL:
		return left.(float64) >= right.(float64)
	case t.LESS:
		return left.(float64) < right.(float64)
	case t.LESS_EQUAL:
		return left.(float64) <= right.(float64)
	case t.MINUS:
		return left.(float64) - right.(float64)
	case t.PLUS:
		if leftFloat, ok := left.(float64); ok {
			if rightFloat, ok := left.(float64); ok {
				return leftFloat + rightFloat
			}
		}

		if leftString, ok := left.(string); ok {
			if rightString, ok := left.(string); ok {
				return leftString + rightString
			}
		}
	case t.MODULO:
		// TODO: implement me
		return nil
	case t.SLASH:
		return left.(float64) / right.(float64)
	case t.STAR:
		return left.(float64) * right.(float64)
	}

	return nil
}

func (i *interpreter) visitGrouping(expr *e.Grouping) any {
	return i.evaluate(expr)
}

func (i *interpreter) visitLiteral(expr *e.Literal) any {
	return expr.Value
}

func (i *interpreter) visitUnary(expr *e.Unary) any {
	right := i.evaluate(expr.Right)

	switch expr.Operator.TokenType {
	case t.MINUS:
		return -right.(float64)
	case t.BANG:
		return !i.isTruthy(right)
	}

	return nil
}

func (i *interpreter) visitTernary(expr *e.Ternary) any {
	// TODO: implement me
	return nil
}

func (i *interpreter) isTruthy(value any) bool {
	if value == nil {
		return false
	}

	if v, ok := value.(bool); ok {
		return v
	}

	return true
}
