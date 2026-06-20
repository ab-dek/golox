package ast

import (
	t "github.com/ab-dek/golox/token"
)

type Expr interface {
	Accept(visitor ExprVisitor) any
}

type ExprVisitor interface {
	VisitAssignment(expr Assignment) any
	VisitBinary(expr Binary) any
	VisitGrouping(expr Grouping) any
	VisitLiteral(expr Literal) any
	VisitLogical(expr Logical) any
	VisitUnary(expr Unary) any
	VisitCall(expr Call) any
	VisitVarExpr(expr VarExpr) any
	VisitTernary(expr Ternary) any
	VisitFuncExpr(expr FuncExpr) any
}

type Assignment struct {
	Name  t.Token
	Value Expr
}

func NewAssignment(name t.Token, value Expr) *Assignment {
	return &Assignment{
		Name:  name,
		Value: value,
	}
}

func (a Assignment) Accept(visitor ExprVisitor) any {
	return visitor.VisitAssignment(a)
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

func (b Binary) Accept(visitor ExprVisitor) any {
	return visitor.VisitBinary(b)
}

type Call struct {
	Callee    Expr
	Paren     t.Token
	Arguments []Expr
}

func NewCall(callee Expr, paren t.Token, arguments []Expr) *Call {
	return &Call{
		Callee:    callee,
		Paren:     paren,
		Arguments: arguments,
	}
}

func (c Call) Accept(visitor ExprVisitor) any {
	return visitor.VisitCall(c)
}

type Grouping struct {
	Expression Expr
}

func NewGrouping(expr Expr) *Grouping {
	return &Grouping{
		Expression: expr,
	}
}

func (g Grouping) Accept(visitor ExprVisitor) any {
	return visitor.VisitGrouping(g)
}

type Literal struct {
	Value any
}

func NewLiteral(value any) *Literal {
	return &Literal{
		Value: value,
	}
}

func (l Literal) Accept(visitor ExprVisitor) any {
	return visitor.VisitLiteral(l)
}

type Logical struct {
	Left     Expr
	Operator t.Token
	Right    Expr
}

func NewLogical(left, right Expr, operator t.Token) *Logical {
	return &Logical{
		Left:     left,
		Right:    right,
		Operator: operator,
	}
}

func (l Logical) Accept(visitor ExprVisitor) any {
	return visitor.VisitLogical(l)
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

func (u Unary) Accept(visitor ExprVisitor) any {
	return visitor.VisitUnary(u)
}

type VarExpr struct {
	Name t.Token
}

func NewVariable(name t.Token) *VarExpr {
	return &VarExpr{
		Name: name,
	}
}

func (v VarExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitVarExpr(v)
}

type Ternary struct {
	Condition Expr
	Question  t.Token
	Then      Expr
	Else      Expr
}

func NewTernary(condition, thenBranch, elseBranch Expr, question t.Token) *Ternary {
	return &Ternary{
		Condition: condition,
		Then:      thenBranch,
		Else:      elseBranch,
		Question:  question,
	}
}

func (t Ternary) Accept(visitor ExprVisitor) any {
	return visitor.VisitTernary(t)
}

type FuncExpr struct {
	Params []t.Token
	Body   []Stmt
}

func NewFuncExpr(params []t.Token, body []Stmt) *FuncExpr {
	return &FuncExpr{
		Params: params,
		Body:   body,
	}
}

func (f FuncExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitFuncExpr(f)
}
