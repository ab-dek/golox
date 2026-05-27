package main

import (
	"fmt"

	e "github.com/ab-dek/golox/ast"
	t "github.com/ab-dek/golox/token"
)

func main() {
	expr := &e.Binary{
		Left: &e.Unary{
			Operator: t.Token{TokenType: t.MINUS, Lexeme: "-", Literal: nil, Line: 1},
			Right:    &e.Literal{Value: 123},
		},
		Operator: t.Token{TokenType: t.STAR, Lexeme: "*", Literal: nil, Line: 1},
		Right: &e.Grouping{
			Expression: &e.Literal{Value: 45.67},
		},
	}

	printer := e.NewPrinter()
	fmt.Println(printer.Print(expr))

	rpnPrinter := e.NewRPN()
	fmt.Println(rpnPrinter.PrintRPN(expr))
}
