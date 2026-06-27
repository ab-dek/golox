package interpreter

import "fmt"

type instance struct {
	class class
}

func newInstance(class class) *instance {
	return &instance{
		class: class,
	}
}

func (i instance) String() string {
	return fmt.Sprintf("<%s instance>", i.class.name)
}
