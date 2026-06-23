package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	errs "github.com/ab-dek/golox/errors"
	i "github.com/ab-dek/golox/interpreter"
	p "github.com/ab-dek/golox/parser"
	r "github.com/ab-dek/golox/resolver"
	sc "github.com/ab-dek/golox/scanner"
	t "github.com/ab-dek/golox/token"
)

func main() {
	args := os.Args
	if len(args) > 2 {
		fmt.Println("Usage: golox [script]")
		os.Exit(64)
	} else if len(args) == 2 {
		runFile(args[1])
	} else {
		runPrompt()
	}
}

func runFile(scriptPath string) {
	source, err := os.ReadFile(scriptPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s : %v\n", scriptPath, err)
		os.Exit(66)
	}

	run(string(source))

	if errs.HadError {
		os.Exit(65)
	}
	if errs.HadRuntimeError {
		os.Exit(70)
	}
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	interpreter := i.NewInterpreter()

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}
		line := scanner.Text()

		scanner := sc.NewScanner(line)

		tokens := scanner.ScanTokens()
		lastToken := tokens[len(tokens)-2]
		parser := p.NewParser(tokens)

		if lastToken.TokenType == t.SEMICOLON || lastToken.TokenType == t.RIGHT_BRACE {
			interpreter.Interpret(parser.Parse()) // parsing a statement
		} else {
			value := interpreter.EvalExpr(parser.ParseExpr()) // parsing expression
			fmt.Printf("%v \n", value)
		}

		errs.HadError = false
	}

	if err := scanner.Err(); err != nil {
		log.Printf("error reading input: %v", err)
	}
}

func run(source string) {
	scanner := sc.NewScanner(source)
	tokens := scanner.ScanTokens()

	// for _, token := range tokens {
	// 	fmt.Println("-----------------------------")
	// 	fmt.Println(token.ToString())
	// }

	parser := p.NewParser(tokens)
	stmts := parser.Parse()

	if errs.HadError {
		return
	}

	// printer := e.NewPrinter()
	// fmt.Println(printer.Print(expr))

	interpreter := i.NewInterpreter()

	resolver := r.NewResolver(interpreter)
	resolver.ResolveStmts(stmts)

	if errs.HadError {
		return
	}

	interpreter.Interpret(stmts)
}
