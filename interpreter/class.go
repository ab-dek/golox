package interpreter

import (
	"fmt"
)

type class struct {
	name          string
	methods       map[string]*Function
	staticMethods map[string]*Function
	superclass    *class
}

func newClass(name string, methods map[string]*Function, staticMethods map[string]*Function, superclass *class) *class {
	return &class{
		name:          name,
		methods:       methods,
		staticMethods: staticMethods,
		superclass:    superclass,
	}
}

func (c class) call(i *Interpreter, arguments []any) any {
	instance := newInstance(c)

	initializer := c.FindMethod("init")

	if initializer != nil {
		initializer.bind(instance).call(i, arguments)
	}

	return instance
}

func (c class) arity() int {
	initializer := c.FindMethod("init")
	if initializer == nil {
		return 0
	}

	return initializer.arity()
}

func (c class) FindMethod(name string) *Function {
	if method, ok := c.methods[name]; ok {
		return method
	}

	return nil
}

func (c class) FindStaticMethod(name string) *Function {
	if method, ok := c.staticMethods[name]; ok {
		return method
	}

	return nil
}

func (c class) String() string {
	return fmt.Sprintf("<class %s>", c.name)
}
