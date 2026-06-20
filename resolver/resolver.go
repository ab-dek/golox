package resolver

import (
	"github.com/ab-dek/golox/ast"
	i "github.com/ab-dek/golox/interpreter"
	s "github.com/ab-dek/golox/stack"
)

type Resolver struct {
	interpreter *i.Interpreter
	scopes      s.Stack[scope]
}

func NewResolver(interpreter *i.Interpreter) *Resolver {
	return &Resolver{
		interpreter: interpreter,
	}
}

func (r *Resolver) Resolve(stmts []ast.Stmt) {
	for _, stmt := range stmts {
		r.resolveStmt(stmt)
	}
}

func (r *Resolver) resolveStmt(stmt ast.Stmt) {
	stmt.Accept(r)
}

func (r *Resolver) resolveExpr(expr ast.Expr) {
	expr.Accept(r)
}

func (r *Resolver) VisitBlock(stmt ast.Block) any {
	r.beginScope()
	r.Resolve(stmt.Statements)
	r.endScope()

	return nil
}

// VisitBreak implements [ast.StmtVisitor].
func (r *Resolver) VisitBreak(stmt ast.BreakStmt) any {
	panic("unimplemented")
}

// VisitContinue implements [ast.StmtVisitor].
func (r *Resolver) VisitContinue(stmt ast.ContinueStmt) any {
	panic("unimplemented")
}

// VisitExpr implements [ast.StmtVisitor].
func (r *Resolver) VisitExpr(stmt ast.ExprStmt) any {
	panic("unimplemented")
}

// VisitFunc implements [ast.StmtVisitor].
func (r *Resolver) VisitFunc(stmt ast.FuncStmt) any {
	panic("unimplemented")
}

// VisitIf implements [ast.StmtVisitor].
func (r *Resolver) VisitIf(stmt ast.IfStmt) any {
	panic("unimplemented")
}

// VisitPrint implements [ast.StmtVisitor].
func (r *Resolver) VisitPrint(stmt ast.PrintStmt) any {
	panic("unimplemented")
}

// VisitReturn implements [ast.StmtVisitor].
func (r *Resolver) VisitReturn(stmt ast.ReturnStmt) any {
	panic("unimplemented")
}

// VisitVar implements [ast.StmtVisitor].
func (r *Resolver) VisitVar(stmt ast.VarStmt) any {
	panic("unimplemented")
}

// VisitWhile implements [ast.StmtVisitor].
func (r *Resolver) VisitWhile(stmt ast.WhileStmt) any {
	panic("unimplemented")
}

// VisitAssignment implements [ast.ExprVisitor].
func (r *Resolver) VisitAssignment(expr ast.Assignment) any {
	panic("unimplemented")
}

// VisitBinary implements [ast.ExprVisitor].
func (r *Resolver) VisitBinary(expr ast.Binary) any {
	panic("unimplemented")
}

// VisitCall implements [ast.ExprVisitor].
func (r *Resolver) VisitCall(expr ast.Call) any {
	panic("unimplemented")
}

// VisitFuncExpr implements [ast.ExprVisitor].
func (r *Resolver) VisitFuncExpr(expr ast.FuncExpr) any {
	panic("unimplemented")
}

// VisitGrouping implements [ast.ExprVisitor].
func (r *Resolver) VisitGrouping(expr ast.Grouping) any {
	panic("unimplemented")
}

// VisitLiteral implements [ast.ExprVisitor].
func (r *Resolver) VisitLiteral(expr ast.Literal) any {
	panic("unimplemented")
}

// VisitLogical implements [ast.ExprVisitor].
func (r *Resolver) VisitLogical(expr ast.Logical) any {
	panic("unimplemented")
}

// VisitTernary implements [ast.ExprVisitor].
func (r *Resolver) VisitTernary(expr ast.Ternary) any {
	panic("unimplemented")
}

// VisitUnary implements [ast.ExprVisitor].
func (r *Resolver) VisitUnary(expr ast.Unary) any {
	panic("unimplemented")
}

// VisitVariable implements [ast.ExprVisitor].
func (r *Resolver) VisitVariable(expr ast.Variable) any {
	panic("unimplemented")
}
func (r *Resolver) beginScope() {
	r.scopes.Push(make(scope))
}

func (r *Resolver) endScope() {
	r.scopes.Pop()
}
