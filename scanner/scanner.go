package scanner

import (
	"fmt"
	"strconv"

	errs "github.com/ab-dek/golox/errors"
	t "github.com/ab-dek/golox/token"
)

var keywords = map[string]t.TokenType{
	"and":    t.AND,
	"class":  t.CLASS,
	"else":   t.ELSE,
	"false":  t.FALSE,
	"for":    t.FOR,
	"fun":    t.FUN,
	"if":     t.IF,
	"nil":    t.NIL,
	"or":     t.OR,
	"print":  t.PRINT,
	"return": t.RETURN,
	"super":  t.SUPER,
	"this":   t.THIS,
	"true":   t.TRUE,
	"var":    t.VAR,
	"while":  t.WHILE,
}

type Scanner struct {
	source               string
	tokens               []t.Token
	start, current, line int
}

func NewScanner(source string) *Scanner {
	return &Scanner{source: source, tokens: make([]t.Token, 0), start: 0, current: 0, line: 1}
}

func (s *Scanner) ScanTokens() []t.Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, *t.NewToken(t.EOF, "", nil, s.line))
	return s.tokens
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(t.LEFT_PAREN)
	case ')':
		s.addToken(t.RIGHT_PAREN)
	case '{':
		s.addToken(t.LEFT_BRACE)
	case '}':
		s.addToken(t.RIGHT_BRACE)
	case ',':
		s.addToken(t.COMMA)
	case '.':
		s.addToken(t.DOT)
	case '-':
		s.addToken(t.MINUS)
	case '+':
		s.addToken(t.PLUS)
	case ';':
		s.addToken(t.SEMICOLON)
	case '*':
		s.addToken(t.STAR)
	case '!':
		if s.match('=') {
			s.addToken(t.BANG_EQUAL)
		} else {
			s.addToken(t.BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(t.EQUAL_EQUAL)
		} else {
			s.addToken(t.EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(t.LESS_EQUAL)
		} else {
			s.addToken(t.LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(t.GREATER_EQUAL)
		} else {
			s.addToken(t.GREATER)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else if s.match('*') {
			s.multiLineComment()
		} else {
			s.addToken(t.SLASH)
		}
	case ' ':
	case '\r':
	case '\t':
		// Ignore whitespace.
	case '\n':
		s.line++
	case '"':
		s.string()
	default:
		if s.isDigit(c) {
			s.number()
		} else if s.isAlpha(c) {
			s.identifier()
		} else {
			errs.Error(s.line, fmt.Sprintf("unexpected character: \"%s\"", string(c)))
		}
	}
}

func (s *Scanner) advance() byte {
	s.current++
	return s.source[s.current-1]
}

func (s *Scanner) addToken(tokenType t.TokenType) {
	s.addTokenWithLiteral(tokenType, nil)
}

func (s *Scanner) addTokenWithLiteral(tokenType t.TokenType, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, *t.NewToken(tokenType, text, literal, s.line))
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}

	if s.source[s.current] != expected {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		errs.Error(s.line, "unterminated string.")
	}

	// the closing "
	s.advance()

	// trim the surrouding quotes
	value := s.source[s.start+1 : s.current-1]
	s.addTokenWithLiteral(t.STRING, value)
}

func (s *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	// look for a fractional part
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		// consume the .
		s.advance()
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	value := s.source[s.start:s.current]
	floatValue, _ := strconv.ParseFloat(value, 64)
	s.addTokenWithLiteral(t.NUMBER, floatValue)
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return 0
	}

	return s.source[s.current+1]
}

func (s *Scanner) isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	tokenType := t.IDENTIFIER
	if t, ok := keywords[text]; ok {
		tokenType = t
	}
	s.addToken(tokenType)
}

func (s *Scanner) isAlphaNumeric(c byte) bool {
	return s.isDigit(c) || s.isAlpha(c)
}

func (s *Scanner) multiLineComment() {
	for !s.isAtEnd() {
		if s.peek() == '*' && s.peekNext() == '/' {
			break
		}
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		errs.Error(s.line, "unterminated multi-line comment.")
	}

	// consuming closing */ tags
	s.advance()
	s.advance()
}
