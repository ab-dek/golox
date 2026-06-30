package interpreter

import (
	"fmt"

	"github.com/ab-dek/golox/ast"
	e "github.com/ab-dek/golox/environment"
)

type Function struct {
	ast.FuncStmt
	Closure *e.Environment
}

func NewFunction(declaration ast.FuncStmt, closure *e.Environment) *Function {
	return &Function{
		FuncStmt: declaration,
		Closure:  closure,
	}
}

func (f Function) call(i *Interpreter, arguments []any) (result any) {
	defer func() {
		if err := recover(); err != nil {
			if ReturnValue, ok := err.(Return); ok {
				result = ReturnValue.value
				return
			}
			panic(err)
		}
	}()

	env := e.NewEnv(f.Closure)
	for i, param := range f.Function.Params {
		env.Define(param.Lexeme, arguments[i])
	}

	i.executeBlock(f.Function.Body, env)

	return nil
}

func (f Function) arity() int {
	return len(f.Function.Params)
}

func (f Function) String() string {
	return fmt.Sprintf("<fn %v>", f.Name.Lexeme)
}

func (f *Function) bind(instance *instance) *Function {
	env := e.NewEnv(f.Closure)
	env.Define("this", instance)
	return NewFunction(f.FuncStmt, env)
}
