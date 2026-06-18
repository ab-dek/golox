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

func (f Function) call(i *interpreter, arguments []any) (result any) {
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
