package token

import "fmt"

type Token struct {
	TokenType TokenType
	Lexeme    string
	Literal   any
	Line      int
}

func NewToken(tokenType TokenType, lexeme string, literal any, line int) *Token {
	return &Token{TokenType: tokenType, Lexeme: lexeme, Literal: literal, Line: line}
}

func (t *Token) ToString() string {
	return fmt.Sprintf("%s %s %v", t.TokenType.String(), t.Lexeme, t.Literal)
}
