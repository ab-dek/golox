package environment

import (
	"fmt"

	errs "github.com/ab-dek/golox/errors"
	t "github.com/ab-dek/golox/token"
)

type Environment struct {
	enclosing *Environment
	values    map[string]any
}

func NewEnv(enclosing *Environment) *Environment {
	return &Environment{
		enclosing: enclosing,
		values:    make(map[string]any),
	}
}

func (e *Environment) Define(key string, value any) {
	e.values[key] = value
}

func (e *Environment) Assign(name t.Token, value any) {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return
	}

	if e.enclosing != nil {
		e.enclosing.Assign(name, value)
		return
	}

	errMsg := errs.ReportRuntimeError(name, fmt.Sprintf("Undefined variable %s. \n", name.Lexeme))
	panic(errMsg)
}
func (e *Environment) Get(name t.Token) any {
	if value, ok := e.values[name.Lexeme]; ok {
		if value == nil {
			errMsg := errs.ReportRuntimeError(name, fmt.Sprintf("Variable not intialized %s. \n", name.Lexeme))
			panic(errMsg)
		}
		return value
	}

	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}

	errMsg := errs.ReportRuntimeError(name, fmt.Sprintf("Undefined variable %s. \n", name.Lexeme))
	panic(errMsg)
}
