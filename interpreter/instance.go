package interpreter

import (
	"fmt"
)

type instance struct {
	class  class
	fields map[string]*any
}

func newInstance(class class) *instance {
	return &instance{
		class:  class,
		fields: make(map[string]*any),
	}
}

func (i instance) String() string {
	return fmt.Sprintf("<%s instance>", i.class.name)
}
