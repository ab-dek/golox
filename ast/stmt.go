package ast

import (
	t "github.com/ab-dek/golox/token"
)

type Stmt interface {
	Accept(visitor StmtVisitor) any
}

type StmtVisitor interface {
	VisitBlock(stmt Block) any
	VisitExpr(stmt ExprStmt) any
	VisitPrint(stmt PrintStmt) any
	VisitVar(stmt VarStmt) any
	VisitIf(stmt IfStmt) any
}

type Block struct {
	Statements []Stmt
}

func NewBlock(statements []Stmt) *Block {
	return &Block{
		Statements: statements,
	}
}

func (b Block) Accept(visitor StmtVisitor) any {
	return visitor.VisitBlock(b)
}

type ExprStmt struct {
	Expr Expr
}

func NewExprStmt(expr Expr) *ExprStmt {
	return &ExprStmt{
		Expr: expr,
	}
}

func (e ExprStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitExpr(e)
}

type IfStmt struct {
	Condition Expr
	Then      Stmt
	Else      Stmt
}

func NewIfStmt(condition Expr, thenBlock, elseBlock Stmt) *IfStmt {
	return &IfStmt{
		Condition: condition,
		Then:      thenBlock,
		Else:      elseBlock,
	}
}

func (i IfStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitIf(i)
}

type PrintStmt struct {
	Expr Expr
}

func NewPrintStmt(expr Expr) *PrintStmt {
	return &PrintStmt{
		Expr: expr,
	}
}

func (p PrintStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitPrint(p)
}

type VarStmt struct {
	Name        t.Token
	Initializer Expr
}

func NewVarStmt(token t.Token, initializer Expr) *VarStmt {
	return &VarStmt{
		Name:        token,
		Initializer: initializer,
	}
}

func (v VarStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitVar(v)
}
