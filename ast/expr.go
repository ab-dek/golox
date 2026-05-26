package ast

import t "github.com/ab-dek/golox/token"

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

func (b *Binary) Accept(visitor ExprVisitor) any {
	return visitor.visitBinary(b)
}

type Grouping struct {
	Expresssion Expr
}

func (g *Grouping) Accept(visitor ExprVisitor) any {
	return visitor.visitGrouping(g)
}

type Literal struct {
	Value any
}

func (l *Literal) Accept(visitor ExprVisitor) any {
	return visitor.visitLiteral(l)
}

type Unary struct {
	Operator t.Token
	Right    Expr
}

func (u *Unary) Accept(visitor ExprVisitor) any {
	return visitor.visitUnary(u)
}
