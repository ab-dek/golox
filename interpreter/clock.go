package interpreter

import "time"

type clock struct{}

func (c clock) call(i *Interpreter, arguments []any) any {
	return float64(time.Now().UnixMilli()) / 1000.0
}
func (c clock) arity() int {
	return 0
}
