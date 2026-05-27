package ast

import (
	t "github.com/ab-dek/golox/token"
)

type Expr interface {
	Accept(visitor ExprVisitor) any
}

type ExprVisitor interface {
	visitBinary(expr *Binary) any
	visitGrouping(expr *Grouping) any
	visitLiteral(expr *Literal) any
	visitUnary(expr *Unary) any
}

type Binary struct {
	Left     Expr
	Operator t.Token
	Right    Expr
}

func NewBinary(left, right Expr, operator t.Token) *Binary {
	return &Binary{
		Left:     left,
		Right:    right,
		Operator: operator,
	}
}

func (b *Binary) Accept(visitor ExprVisitor) any {
	return visitor.visitBinary(b)
}

type Grouping struct {
	Expression Expr
}

func NewGrouping(expr Expr) *Grouping {
	return &Grouping{
		Expression: expr,
	}
}

func (g *Grouping) Accept(visitor ExprVisitor) any {
	return visitor.visitGrouping(g)
}

type Literal struct {
	Value any
}

func NewLiteral(value any) *Literal {
	return &Literal{
		Value: value,
	}
}

func (l *Literal) Accept(visitor ExprVisitor) any {
	return visitor.visitLiteral(l)
}

type Unary struct {
	Operator t.Token
	Right    Expr
}

func NewUnary(operator t.Token, right Expr) *Unary {
	return &Unary{
		Operator: operator,
		Right:    right,
	}
}

func (u *Unary) Accept(visitor ExprVisitor) any {
	return visitor.visitUnary(u)
}
