package interpreter

import (
	"fmt"
)

type class struct {
	name string
}

func newClass(name string) *class {
	return &class{
		name: name,
	}
}

func (c class) call(i *Interpreter, arguments []any) any {
	instance := newInstance(c)
	return instance
}

func (c class) arity() int {
	return 0
}

func (c class) String() string {
	return fmt.Sprintf("<class %s>", c.name)
}
