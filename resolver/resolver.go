package resolver

import (
	"fmt"

	"github.com/ab-dek/golox/ast"
	errs "github.com/ab-dek/golox/errors"
	i "github.com/ab-dek/golox/interpreter"
	s "github.com/ab-dek/golox/stack"
	t "github.com/ab-dek/golox/token"
)

type functionType int

const (
	NONE functionType = iota
	FUNCTION
	METHOD
)

type Resolver struct {
	interpreter     *i.Interpreter
	scopes          s.Stack[scope]
	currentFunction functionType
}

func NewResolver(interpreter *i.Interpreter) *Resolver {
	return &Resolver{
		interpreter:     interpreter,
		currentFunction: NONE,
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
	r.checkUnusedVariable("variable")
	r.endScope()

	return nil
}

func (r *Resolver) VisitClassStmt(stmt ast.ClassStmt) any {
	r.declare(stmt.Name)
	r.define(stmt.Name)

	r.beginScope()
	scope, _ := r.scopes.Peek()
	scope["this"] = &varInfo{
		token:    *t.NewToken(t.IDENTIFIER, "this", nil, 0),
		resolved: true,
	}

	for _, method := range stmt.Methods {
		r.resolveFunction(method.Function, METHOD)
	}
	r.endScope()

	return nil
}

// VisitBreak implements [ast.StmtVisitor].
func (r *Resolver) VisitBreak(stmt ast.BreakStmt) any {
	return nil
}

// VisitContinue implements [ast.StmtVisitor].
func (r *Resolver) VisitContinue(stmt ast.ContinueStmt) any {
	return nil
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

	r.resolveFunction(stmt.Function, FUNCTION)

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
	if r.currentFunction == NONE {
		errMsg := errs.ParseError(stmt.Keyword, "Can't return from a top-level code.")
		errs.ReportError(errMsg)
	}

	if stmt.Value != nil {
		r.resolveExpr(stmt.Value)
	}

	return nil
}

// VisitVarStmt implements [ast.StmtVisitor].
func (r *Resolver) VisitVarStmt(stmt ast.VarStmt) any {
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

// VisitGet implements [ast.ExprVisitor].
func (r *Resolver) VisitGet(expr ast.Get) any {
	r.resolveExpr(expr.Object)
	return nil
}

// VisitFuncExpr implements [ast.ExprVisitor].
func (r *Resolver) VisitFuncExpr(expr ast.FuncExpr) any {
	r.resolveFunction(expr, FUNCTION)

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

// VisitSet implements [ast.ExprVisitor].
func (r *Resolver) VisitSet(expr ast.Set) any {
	r.resolveExpr(expr.Value)
	r.resolveExpr(expr.Object)

	return nil
}

// VisitThis implements [ast.ExprVisitor].
func (r *Resolver) VisitThis(expr ast.This) any {
	r.resolveLocal(expr, expr.Keyword)
	return nil
}

// VisitUnary implements [ast.ExprVisitor].
func (r *Resolver) VisitUnary(expr ast.Unary) any {
	r.resolveExpr(expr.Right)
	return nil
}

// VisitVarExpr implements [ast.ExprVisitor].
func (r *Resolver) VisitVarExpr(expr ast.VarExpr) any {
	scope, err := r.scopes.Peek()
	if err == nil { // check if scope stack is not empty
		if v, ok := scope[expr.Name.Lexeme]; ok && !v.resolved {
			errMsg := errs.ParseError(expr.Name, "Can't read local variable in its own initializer.")
			errs.ReportError(errMsg)
		}
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
	if _, ok := scope[name.Lexeme]; ok {
		errMsg := errs.ParseError(name, "Already variable with this name in this scope.")
		errs.ReportError(errMsg)
	}
	scope[name.Lexeme] = &varInfo{
		token:    name,
		resolved: false,
	}
}

func (r *Resolver) define(name t.Token) {
	if r.scopes.IsEmpty() {
		return
	}

	scope, _ := r.scopes.Peek()
	scope[name.Lexeme].resolved = true
}

func (r *Resolver) resolveLocal(expr ast.Expr, name t.Token) {
	lenStack := r.scopes.Size()
	for i := lenStack - 1; i >= 0; i-- {
		scope, _ := r.scopes.Get(i)
		if v, ok := scope[name.Lexeme]; ok {
			r.interpreter.Resolve(expr, lenStack-1-i)
			v.used = true
			return
		}
	}
}

func (r *Resolver) resolveFunction(function ast.FuncExpr, funcType functionType) {
	enclosingFunc := r.currentFunction
	r.currentFunction = funcType

	r.beginScope()

	for _, param := range function.Params {
		r.declare(param)
		r.define(param)
	}

	r.ResolveStmts(function.Body)
	r.checkUnusedVariable("parameter")
	r.endScope()

	r.currentFunction = enclosingFunc
}

func (r *Resolver) checkUnusedVariable(kind string) {
	scope, _ := r.scopes.Peek()
	for _, value := range scope {
		if !value.used {
			errs.ReportWarning(value.token, fmt.Sprintf("%s never used.", kind))
		}
	}
}
