package interpreter

type callable interface {
	call(i *interpreter, arguments []any) any
	arity() int
}
