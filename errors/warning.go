package errors

import (
	"fmt"
	"os"

	t "github.com/ab-dek/golox/token"
)

func ReportWarning(token t.Token, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Warning %s: %s\n", token.Line, " at '"+token.Lexeme+"'", message)
}
