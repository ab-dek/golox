package interpreter

import (
	"fmt"
)

type class struct {
	name    string
	methods map[string]*Function
}

func newClass(name string, methods map[string]*Function) *class {
	return &class{
		name:    name,
		methods: methods,
	}
}

func (c class) call(i *Interpreter, arguments []any) any {
	instance := newInstance(c)
	return instance
}

func (c class) arity() int {
	return 0
}

func (c class) FindMethod(name string) *Function {
	if method, ok := c.methods[name]; ok {
		return method
	}

	return nil
}

func (c class) String() string {
	return fmt.Sprintf("<class %s>", c.name)
}
