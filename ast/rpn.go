package ast

import "fmt"

// reverse polish notation
type RPN struct{}

func NewRPN() *RPN {
	return &RPN{}
}

func (r *RPN) PrintRPN(expr Expr) string {
	return expr.Accept(r).(string)
}

func (r *RPN) visitBinary(expr *Binary) any {
	left := expr.Left.Accept(r).(string)
	right := expr.Right.Accept(r).(string)
	operator := expr.Operator.Lexeme

	return fmt.Sprintf("%s %s %s", left, right, operator)
}

func (r *RPN) visitGrouping(expr *Grouping) any {
	return expr.Expresssion.Accept(r)
}

func (r *RPN) visitLiteral(expr *Literal) any {
	return fmt.Sprintf("%v", expr.Value)
}

func (r *RPN) visitUnary(expr *Unary) any {
	right := expr.Right.Accept(r).(string)
	operator := expr.Operator.Lexeme

	return fmt.Sprintf("%s %s", right, operator)
}
