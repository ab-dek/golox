package interpreter

type callable interface {
	call(i *Interpreter, arguments []any) any
	arity() int
}
