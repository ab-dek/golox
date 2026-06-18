package interpreter

import (
	"fmt"

	"github.com/ab-dek/golox/ast"
	e "github.com/ab-dek/golox/environment"
)

type Function ast.FuncStmt

func NewFunction(funcStmt ast.FuncStmt) *Function {
	function := Function(funcStmt)
	return &function
}

func (f Function) call(i *interpreter, arguments []any) any {
	env := e.NewEnv(i.global)
	for i, param := range f.Params {
		env.Define(param.Lexeme, arguments[i])
	}

	i.executeBlock(f.Body, env)

	return nil
}

func (f Function) arity() int {
	return len(f.Params)
}

func (f Function) ToString() string {
	return fmt.Sprintf("<fn %v>", f.Name.Lexeme)
}
