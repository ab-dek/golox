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

func ParseError(token t.Token, message string) {
	if token.TokenType == t.EOF {
		Report(token.Line, " at end", message)
	} else {
		Report(token.Line, " at '"+token.Lexeme+"'", message)
	}
}

func Report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error %s: %s\n", line, where, message)
}

func RuntimeError(token t.Token, message string) {
	Report(token.Line, "", message)
	HadRuntimeError = true
}
