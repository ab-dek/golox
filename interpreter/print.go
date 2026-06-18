package interpreter

import "fmt"

// currently not used
type print struct{}

func (p print) call(i *interpreter, arguments []any) any {
	fmt.Printf("%v", arguments[0])
	return nil
}

func (p print) arity() int {
	return 1
}
