package resolver

import (
	"github.com/ab-dek/golox/ast"
	errs "github.com/ab-dek/golox/errors"
	i "github.com/ab-dek/golox/interpreter"
	s "github.com/ab-dek/golox/stack"
	t "github.com/ab-dek/golox/token"
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

func (r *Resolver) ResolveStmts(stmts []ast.Stmt) {
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
	r.ResolveStmts(stmt.Statements)
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
	r.resolveExpr(stmt.Expr)
	return nil
}

// VisitFuncStmt implements [ast.StmtVisitor].
func (r *Resolver) VisitFuncStmt(stmt ast.FuncStmt) any {
	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.resolveFunction(stmt.Function)

	return nil
}

// VisitIf implements [ast.StmtVisitor].
func (r *Resolver) VisitIf(stmt ast.IfStmt) any {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.Then)

	if stmt.Else != nil {
		r.resolveStmt(stmt.Else)
	}

	return nil
}

// VisitPrint implements [ast.StmtVisitor].
func (r *Resolver) VisitPrint(stmt ast.PrintStmt) any {
	r.resolveExpr(stmt.Expr)
	return nil
}

// VisitReturn implements [ast.StmtVisitor].
func (r *Resolver) VisitReturn(stmt ast.ReturnStmt) any {
	if stmt.Value != nil {
		r.resolveExpr(stmt.Value)
	}

	return nil
}

// VisitVar implements [ast.StmtVisitor].
func (r *Resolver) VisitVar(stmt ast.VarStmt) any {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolveExpr(stmt.Initializer)
	}
	r.define(stmt.Name)
	return nil
}

// VisitWhile implements [ast.StmtVisitor].
func (r *Resolver) VisitWhile(stmt ast.WhileStmt) any {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.Body)

	return nil
}

// VisitAssignment implements [ast.ExprVisitor].
func (r *Resolver) VisitAssignment(expr ast.Assignment) any {
	r.resolveExpr(expr.Value)
	r.resolveLocal(expr, expr.Name)

	return nil
}

// VisitBinary implements [ast.ExprVisitor].
func (r *Resolver) VisitBinary(expr ast.Binary) any {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)

	return nil
}

// VisitCall implements [ast.ExprVisitor].
func (r *Resolver) VisitCall(expr ast.Call) any {
	r.resolveExpr(expr.Callee)

	for _, arg := range expr.Arguments {
		r.resolveExpr(arg)
	}

	return nil
}

// VisitFuncExpr implements [ast.ExprVisitor].
func (r *Resolver) VisitFuncExpr(expr ast.FuncExpr) any {
	r.resolveFunction(expr)

	return nil
}

// VisitGrouping implements [ast.ExprVisitor].
func (r *Resolver) VisitGrouping(expr ast.Grouping) any {
	r.resolveExpr(expr.Expression)
	return nil
}

// VisitLiteral implements [ast.ExprVisitor].
func (r *Resolver) VisitLiteral(expr ast.Literal) any {
	return nil
}

// VisitLogical implements [ast.ExprVisitor].
func (r *Resolver) VisitLogical(expr ast.Logical) any {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)

	return nil
}

// VisitTernary implements [ast.ExprVisitor].
func (r *Resolver) VisitTernary(expr ast.Ternary) any {
	r.resolveExpr(expr.Condition)
	r.resolveExpr(expr.Else)
	r.resolveExpr(expr.Then)

	return nil
}

// VisitUnary implements [ast.ExprVisitor].
func (r *Resolver) VisitUnary(expr ast.Unary) any {
	r.resolveExpr(expr.Right)
	return nil
}

// VisitVariable implements [ast.ExprVisitor].
func (r *Resolver) VisitVariable(expr ast.Variable) any {
	scope, err := r.scopes.Peek()
	if err == nil && scope[expr.Name.Lexeme] == false {
		errMsg := errs.ParseError(expr.Name, "Can't read local variable in its own initializer.")
		errs.ReportError(errMsg)
	}

	r.resolveLocal(expr, expr.Name)

	return nil
}

func (r *Resolver) beginScope() {
	r.scopes.Push(make(scope))
}

func (r *Resolver) endScope() {
	r.scopes.Pop()
}

func (r *Resolver) declare(name t.Token) {
	if r.scopes.IsEmpty() {
		return
	}

	scope, _ := r.scopes.Peek()
	scope[name.Lexeme] = false
}

func (r *Resolver) define(name t.Token) {
	if r.scopes.IsEmpty() {
		return
	}

	scope, _ := r.scopes.Peek()
	scope[name.Lexeme] = true
}

func (r *Resolver) resolveLocal(expr ast.Expr, name t.Token) {
	lenStack := r.scopes.Size()
	for i := lenStack - 1; i >= 0; i-- {
		scope, _ := r.scopes.Get(i)
		if _, ok := scope[name.Lexeme]; ok {
			r.interpreter.Resolve(expr, lenStack)
		}
	}
}

func (r *Resolver) resolveFunction(function ast.FuncExpr) {
	r.beginScope()

	for _, param := range function.Params {
		r.declare(param)
		r.define(param)
	}

	r.ResolveStmts(function.Body)
	r.endScope()
}
