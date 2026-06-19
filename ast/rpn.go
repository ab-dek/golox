package ast

import (
	"fmt"
)

// reverse polish notation
type RPN struct{}

func NewRPN() *RPN {
	return &RPN{}
}

func (r *RPN) PrintRPN(expr Expr) string {
	return expr.Accept(r).(string)
}

func (r *RPN) VisitAssignment(expr Assignment) any {
	//TODO: implement me
	panic("unimplemented")
}

func (r *RPN) VisitBinary(expr Binary) any {
	left := expr.Left.Accept(r).(string)
	right := expr.Right.Accept(r).(string)
	operator := expr.Operator.Lexeme

	return fmt.Sprintf("%s %s %s", left, right, operator)
}

func (r *RPN) VisitGrouping(expr Grouping) any {
	return expr.Expression.Accept(r)
}

func (r *RPN) VisitLiteral(expr Literal) any {
	return fmt.Sprintf("%v", expr.Value)
}

func (r *RPN) VisitLogical(expr Logical) any {
	left := expr.Left.Accept(r).(string)
	right := expr.Right.Accept(r).(string)
	operator := expr.Operator.Lexeme

	return fmt.Sprintf("%s %s %s", left, right, operator)
}

func (r *RPN) VisitUnary(expr Unary) any {
	right := expr.Right.Accept(r).(string)
	operator := expr.Operator.Lexeme

	return fmt.Sprintf("%s %s", right, operator)
}

func (r *RPN) VisitCall(expr Call) any {
	panic("unimplemented")
}

func (r *RPN) VisitVariable(expr Variable) any {
	panic("unimplemented")
}

func (r *RPN) VisitTernary(expr Ternary) any {
	condition := expr.Condition.Accept(r).(string)
	then := expr.Then.Accept(r).(string)
	Else := expr.Else.Accept(r).(string)
	return fmt.Sprintf("%s %s %s", condition, then, Else)
}

func (r *RPN) VisitFuncExpr(expr FuncExpr) any {
	panic("unimplemented")
}
