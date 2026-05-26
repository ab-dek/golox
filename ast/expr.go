package ast

import t "github.com/ab-dek/golox/token"

type Expr interface {
	exprNode()
}

type BinaryExpr struct {
	leftExpr  Expr
	operator  t.Token
	rightExpr Expr
}

func (b *BinaryExpr) exprNode() {}

type GroupingExpr struct {
	expresssion Expr
}

func (g *GroupingExpr) exprNode() {}

type LiteralExpr struct {
	value any
}

func (l *LiteralExpr) exprNode() {}

type UnaryExpr struct {
	operator  t.Token
	rigthExpr Expr
}

func (u *UnaryExpr) exprNode() {}
