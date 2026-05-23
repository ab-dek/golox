package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	errs "github.com/ab-dek/golox/errors"
	sc "github.com/ab-dek/golox/scanner"
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
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		run(line)
		errs.HadError = false
	}
	if err := scanner.Err(); err != nil {
		log.Printf("error reading input: %v", err)
	}
}

func run(source string) {
	fmt.Println("lexing source code: ", source)
	scanner := sc.NewScanner(source)
	tokens := scanner.ScanTokens()

	for _, token := range tokens {
		fmt.Println(token.ToString())
	}
}
