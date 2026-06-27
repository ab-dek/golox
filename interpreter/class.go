package interpreter

import "fmt"

type class struct {
	name string
}

func newClass(name string) *class {
	return &class{
		name: name,
	}
}

func (c class) ToString() string {
	return fmt.Sprintf("<class %s>", c.name)
}
