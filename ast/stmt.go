package ast

type Stmt interface {
	Accept(visitor StmtVisitor) any
}

type StmtVisitor interface {
	VisitExpr(stmt *ExprStmt) any
	VisitPrint(stmt *PrintStmt) any
}

type ExprStmt struct {
	Expr Expr
}

func NewExprStmt(expr Expr) *ExprStmt {
	return &ExprStmt{
		Expr: expr,
	}
}

func (e *ExprStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitExpr(e)
}

type PrintStmt struct {
	Expr Expr
}

func NewPrintStmt(expr Expr) *PrintStmt {
	return &PrintStmt{
		Expr: expr,
	}
}

func (p *PrintStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitPrint(p)
}
