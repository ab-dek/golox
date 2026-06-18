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
	VisitWhile(stmt WhileStmt) any
	VisitBreak(stmt BreakStmt) any
	VisitContinue(stmt ContinueStmt) any
	VisitFunc(stmt FuncStmt) any
	VisitReturn(stmt ReturnStmt) any
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

type WhileStmt struct {
	Condition Expr
	Body      Stmt
	Increment Expr // value is only set when desugaring a for loop
}

func NewWhileStmt(condition, increment Expr, body Stmt) *WhileStmt {
	return &WhileStmt{
		Condition: condition,
		Body:      body,
		Increment: increment,
	}
}

func (w WhileStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitWhile(w)
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

type FuncStmt struct {
	Name   t.Token
	Params []t.Token
	Body   []Stmt
}

func NewFunc(name t.Token, params []t.Token, body []Stmt) *FuncStmt {
	return &FuncStmt{
		Name:   name,
		Params: params,
		Body:   body,
	}
}

func (f FuncStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitFunc(f)
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

type ReturnStmt struct {
	Keyword t.Token
	Value   Expr
}

func NewReturn(keyword t.Token, value Expr) *ReturnStmt {
	return &ReturnStmt{
		Keyword: keyword,
		Value:   value,
	}
}

func (r ReturnStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitReturn(r)
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

type BreakStmt struct{}

func NewBreakStmt() *BreakStmt {
	return &BreakStmt{}
}

func (l BreakStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitBreak(l)
}

type ContinueStmt struct{}

func NewContinueStmt() *ContinueStmt {
	return &ContinueStmt{}
}

func (c ContinueStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitContinue(c)
}
