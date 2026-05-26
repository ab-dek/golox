package ast

import t "github.com/ab-dek/golox/token"

type Expr interface {
	accept(visitor ExprVisitor)
}

type ExprVisitor interface {
	visitBinary(expr *Binary)
	visitGrouping(expr *Grouping)
	visitLiteral(expr *Literal)
	visitUnary(expr *Unary)
}

type Binary struct {
	left     Expr
	operator t.Token
	right    Expr
}

func (b *Binary) accept(visitor ExprVisitor) {
	visitor.visitBinary(b)
}

type Grouping struct {
	expresssion Expr
}

func (g *Grouping) accept(visitor ExprVisitor) {
	visitor.visitGrouping(g)
}

type Literal struct {
	value any
}

func (l *Literal) accept(visitor ExprVisitor) {
	visitor.visitLiteral(l)
}

type Unary struct {
	operator t.Token
	rigth    Expr
}

func (u *Unary) accept(visitor ExprVisitor) {
	visitor.visitUnary(u)
}
