package errors

import (
	"fmt"
	"os"

	t "github.com/ab-dek/golox/token"
)

var HadError bool
var HadRuntimeError bool

func LexError(line int, message string) {
	report(line, "", message)
}

func ParseError(token t.Token, message string) string {
	if token.TokenType == t.EOF {
		return fmt.Sprintf("[line %d] Error %s: %s\n", token.Line, "at end", message)
	} else {
		return fmt.Sprintf("[line %d] Error %s: %s\n", token.Line, " at '"+token.Lexeme+"'", message)
	}
}

func report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error %s: %s\n", line, where, message)
}

func ReportError(message string) {
	fmt.Fprint(os.Stderr, message)
}

func RuntimeError(token t.Token, message string) string {
	HadRuntimeError = true
	return fmt.Sprintf("[line %d] %s\n", token.Line, message)
}
