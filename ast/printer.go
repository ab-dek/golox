package ast

import (
	"fmt"
	"strings"
)

type printer struct{}

func NewPrinter() *printer {
	return &printer{}
}

func (p *printer) Print(expr Expr) string {
	return expr.Accept(p).(string)
}

func (p *printer) VisitBinary(expr *Binary) any {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *printer) VisitGrouping(expr *Grouping) any {
	return p.parenthesize("group", expr.Expression)
}

func (p *printer) VisitLiteral(expr *Literal) any {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.Value)
}

func (p *printer) VisitUnary(expr *Unary) any {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (p *printer) VisitTernary(expr *Ternary) any {
	return p.parenthesize("? :", expr.Condition, expr.Then, expr.Else)
}

func (p *printer) parenthesize(name string, exprs ...Expr) string {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(name)

	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(p.Print(expr))
	}

	builder.WriteString(")")

	return builder.String()
}
