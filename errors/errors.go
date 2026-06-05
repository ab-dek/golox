package errors

import (
	"fmt"
	"os"

	t "github.com/ab-dek/golox/token"
)

var HadError bool
var HadRuntimeError bool

func LexError(line int, message string) {
	Report(line, "", message)
}

type ParseError struct{}

func ReportParseError(token t.Token, message string) {
	if token.TokenType == t.EOF {
		Report(token.Line, " at end", message)
	} else {
		Report(token.Line, " at '"+token.Lexeme+"'", message)
	}
}

func Report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error %s: %s\n", line, where, message)
}

type RuntimeError struct{}

func ReportRuntimeError(token t.Token, message string) {
	Report(token.Line, "", message)
	HadRuntimeError = true
}
