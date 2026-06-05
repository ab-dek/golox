package interpreter

import (
	"fmt"
	"math"

	"github.com/ab-dek/golox/ast"
	e "github.com/ab-dek/golox/environment"
	errs "github.com/ab-dek/golox/errors"
	t "github.com/ab-dek/golox/token"
)

type interpreter struct {
	env *e.Environment
}

func NewInterpreter() *interpreter {
	return &interpreter{
		env: e.NewEnv(),
	}
}

func (i *interpreter) Interpret(stmts []ast.Stmt) {
	for _, stmt := range stmts {
		i.execute(stmt)
	}
}

func (i *interpreter) execute(stmt ast.Stmt) {
	stmt.Accept(i)
}

func (i *interpreter) evaluate(expr ast.Expr) any {
	return expr.Accept(i)
}

func (i *interpreter) VisitExpr(stmt *ast.ExprStmt) any {
	i.evaluate(stmt.Expr)
	return nil
}

func (i *interpreter) VisitPrint(stmt *ast.PrintStmt) any {
	value := i.evaluate(stmt.Expr)
	fmt.Printf("%v \n", value)
	return nil
}

func (i *interpreter) VisitVar(stmt *ast.VarStmt) any {
	var value any
	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}
	i.env.Define(stmt.Name.Lexeme, value)
	return nil
}

func (i *interpreter) VisitBinary(expr *ast.Binary) any {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)
	fmt.Printf("evaluating binary epxr: %v %v %v \n", left, expr.Operator.Lexeme, right)

	switch expr.Operator.TokenType {
	case t.BANG_EQUAL:
		return !(left == right)
	case t.EQUAL_EQUAL:
		return left == right
	case t.GREATER:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) > right.(float64)
	case t.GREATER_EQUAL:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) >= right.(float64)
	case t.LESS:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) < right.(float64)
	case t.LESS_EQUAL:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) <= right.(float64)
	case t.MINUS:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) - right.(float64)
	case t.PLUS:
		if leftFloat, ok := left.(float64); ok {
			if rightFloat, ok := right.(float64); ok {
				return leftFloat + rightFloat
			}
		}

		if leftString, ok := left.(string); ok {
			if rightString, ok := right.(string); ok {
				return leftString + rightString
			}
		}
		errs.ReportRuntimeError(expr.Operator, "Operand must be a number.")
		panic(errs.RuntimeError{})
	case t.MODULO:
		i.checkNumberOperands(expr.Operator, left, right)
		if right.(float64) == 0 {
			errs.ReportRuntimeError(expr.Operator, "Cannot modulo by 0.")
			panic(errs.RuntimeError{})
		}
		return math.Mod(left.(float64), right.(float64))
	case t.SLASH:
		i.checkNumberOperands(expr.Operator, left, right)
		if right.(float64) == 0 {
			errs.ReportRuntimeError(expr.Operator, "Cannot divide by 0.")
			panic(errs.RuntimeError{})
		}
		return left.(float64) / right.(float64)
	case t.STAR:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) * right.(float64)
	}

	return nil
}

func (i *interpreter) VisitGrouping(expr *ast.Grouping) any {
	fmt.Printf("unwrapping parenthesis: %v \n", expr.Expression)
	return i.evaluate(expr.Expression)
}

func (i *interpreter) VisitLiteral(expr *ast.Literal) any {
	fmt.Printf("extracting literal value: %v \n", expr.Value)
	return expr.Value
}

func (i *interpreter) VisitUnary(expr *ast.Unary) any {
	right := i.evaluate(expr.Right)
	fmt.Printf("evaluating unary epxr: %v %v \n", expr.Operator.Lexeme, right)

	switch expr.Operator.TokenType {
	case t.MINUS:
		i.checkNumberOperand(expr.Operator, right)
		return -right.(float64)
	case t.BANG:
		return !i.isTruthy(right)
	}

	return nil
}

func (i *interpreter) VisitVariable(expr *ast.Variable) any {
	return i.env.Get(expr.Name)
}

func (i *interpreter) VisitTernary(expr *ast.Ternary) any {
	// TODO: implement me
	panic("ternary interpreter unimplemented")
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

func (i *interpreter) checkNumberOperand(operator t.Token, operand any) {
	if _, isFloat := operand.(float64); isFloat {
		return
	}
	errs.ReportRuntimeError(operator, "Operand must be a number.")
	panic(errs.RuntimeError{})
}

func (i *interpreter) checkNumberOperands(operator t.Token, right, left any) {
	if _, isFloat := right.(float64); isFloat {
		if _, isFloat := left.(float64); isFloat {
			return
		}
	}
	errs.ReportRuntimeError(operator, "Operand must be a number.")
	panic(errs.RuntimeError{})
}
