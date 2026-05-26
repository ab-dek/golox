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
	return fmt.Sprintf("%v", expr.Accept(p))
}

func (p *printer) visitBinary(expr *Binary) any {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *printer) visitGrouping(expr *Grouping) any {
	return p.parenthesize("group", expr.Expresssion)
}

func (p *printer) visitLiteral(expr *Literal) any {
	if expr.Value == nil {
		return "nil"
	}
	return expr.Value
}

func (p *printer) visitUnary(expr *Unary) any {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
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
