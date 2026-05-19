package main

import (
	"bufio"
	"fmt"
	"os"

	sc "github.com/ab-dek/golox/scanner"
)

var hadError bool

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

	if hadError {
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
		hadError = false
	}
}

func run(source string) {
	scanner := sc.NewScanner(source)
	tokens := scanner.ScanTokens()

	for token := range tokens {
		fmt.Println(token)
	}
}

func Error(line int, message string) {
	Report(line, "", message)
}

func Report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error %s: %s\n", line, where, message)
}
