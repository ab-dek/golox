package environment

import (
	"fmt"

	errs "github.com/ab-dek/golox/errors"
	t "github.com/ab-dek/golox/token"
)

type Environment struct {
	values map[string]any
}

func NewEnv() *Environment {
	return &Environment{
		values: make(map[string]any),
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
	errs.ReportRuntimeError(name, fmt.Sprintf("Undefined variable %s. \n", name.Lexeme))
	panic(errs.RuntimeError{})
}
func (e *Environment) Get(name t.Token) any {
	if value, ok := e.values[name.Lexeme]; ok {
		return value
	}
	errs.ReportRuntimeError(name, fmt.Sprintf("Undefined variable %s. \n", name.Lexeme))
	panic(errs.RuntimeError{})
}
