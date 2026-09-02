package environment

import (
	"fmt"

	errs "github.com/ab-dek/golox/errors"
	t "github.com/ab-dek/golox/token"
)

type Environment struct {
	Enclosing *Environment
	values    map[string]any
}

func NewEnv(enclosing *Environment) *Environment {
	return &Environment{
		Enclosing: enclosing,
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

	if e.Enclosing != nil {
		e.Enclosing.Assign(name, value)
		return
	}

	errMsg := errs.RuntimeError(name, fmt.Sprintf("Undefined variable %s. \n", name.Lexeme))
	panic(errMsg)
}

func (e *Environment) Get(name t.Token) any {
	if value, ok := e.values[name.Lexeme]; ok {
		if value == nil {
			errMsg := errs.RuntimeError(name, fmt.Sprintf("Variable not intialized %s. \n", name.Lexeme))
			panic(errMsg)
		}
		return value
	}

	if e.Enclosing != nil {
		return e.Enclosing.Get(name)
	}

	errMsg := errs.RuntimeError(name, fmt.Sprintf("Undefined variable %s. \n", name.Lexeme))
	panic(errMsg)
}

func (e *Environment) GetAt(distance int, name string) any {
	return e.ancestor(distance).values[name]
}

func (e *Environment) ancestor(distance int) *Environment {
	env := e
	for range distance {
		env = env.Enclosing
	}

	return env
}

func (e *Environment) AssignAt(distance int, name t.Token, value any) {
	e.ancestor(distance).values[name.Lexeme] = value
}
