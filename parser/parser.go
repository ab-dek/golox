package parser

import (
	"github.com/ab-dek/golox/ast"
	errs "github.com/ab-dek/golox/errors"
	t "github.com/ab-dek/golox/token"
)

/*
grammmer:

program→ statement* EOF ;
statement→ exprStmt
		   | printStmt ;
exprStmt→ expression ";" ;
printStmt→ "print" expression ";" ;
expression→ conditional ;
conditional→ equality ( "?" expression ":" conditional )? ;
equality→ comparison ( ( "!=" | "==" ) comparison )* ;
comparison→ term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term→ factor ( ( "-" | "+" ) factor )* ;
factor→ unary ( ( "/" | "*" ) unary )* ;
unary→ ( "!" | "-" ) unary
	   | primary ;
primary→ NUMBER | STRING | "true" | "false" | "nil"
| "(" expression ")" ;
*/

type Parser struct {
	tokens  []t.Token
	current int
}

func NewParser(tokens []t.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() []ast.Stmt {
	var statements []ast.Stmt
	for !p.IsAtEnd() {
		statements = append(statements, p.statement())
	}
	return statements
}

func (p *Parser) statement() ast.Stmt {
	if p.match(t.PRINT) {
		return p.printStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) printStatement() ast.Stmt {
	expr := p.expression()
	p.consume(t.SEMICOLON, "Expect ';' after value.")
	return ast.NewPrintStmt(expr)
}

func (p *Parser) expressionStatement() ast.Stmt {
	expr := p.expression()
	p.consume(t.SEMICOLON, "Expect ';' after value.")
	return ast.NewExprStmt(expr)
}

func (p *Parser) expression() ast.Expr {
	return p.conditional()
}

func (p *Parser) conditional() ast.Expr {
	expr := p.equality()
	if p.match(t.QUESTION_MARK) {
		questionMark := p.previous()
		thenBranch := p.expression()

		p.consume(t.COLON, "Expect ':' after then-branch of conditional operator.")

		elseBranch := p.conditional()
		expr = ast.NewTernary(expr, thenBranch, elseBranch, questionMark)
	}
	return expr
}

func (p *Parser) equality() ast.Expr {
	expr := p.comparision()

	for p.match(t.BANG_EQUAL, t.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparision()
		expr = ast.NewBinary(expr, right, operator)
	}

	return expr
}

func (p *Parser) comparision() ast.Expr {
	expr := p.term()

	for p.match(t.GREATER, t.GREATER_EQUAL, t.LESS, t.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = ast.NewBinary(expr, right, operator)
	}

	return expr
}

func (p *Parser) term() ast.Expr {
	expr := p.factor()

	for p.match(t.MINUS, t.PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = ast.NewBinary(expr, right, operator)
	}

	return expr
}

func (p *Parser) factor() ast.Expr {
	expr := p.unary()

	for p.match(t.STAR, t.SLASH, t.MODULO) {
		operator := p.previous()
		right := p.unary()
		expr = ast.NewBinary(expr, right, operator)
	}

	return expr
}

func (p *Parser) unary() ast.Expr {
	if p.match(t.BANG, t.MINUS) {
		operator := p.previous()
		right := p.unary()
		return ast.NewUnary(operator, right)
	}

	return p.primary()
}

func (p *Parser) primary() ast.Expr {
	switch {
	case p.match(t.FALSE):
		return ast.NewLiteral(false)
	case p.match(t.TRUE):
		return ast.NewLiteral(true)
	case p.match(t.NIL):
		return ast.NewLiteral(nil)
	case p.match(t.NUMBER, t.STRING):
		return ast.NewLiteral(p.previous().Literal)
	case p.match(t.LEFT_PAREN):
		expr := p.expression()
		p.consume(t.RIGHT_PAREN, "Expect ')' after expression.")
		return ast.NewGrouping(expr)
	}
	p.error(p.peek(), "Expect expression.")
	return nil
}

func (p *Parser) match(types ...t.TokenType) bool {
	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokenType t.TokenType) bool {
	if p.IsAtEnd() {
		return false
	}
	return p.peek().TokenType == tokenType
}

func (p *Parser) advance() t.Token {
	if !p.IsAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) IsAtEnd() bool {
	return p.peek().TokenType == t.EOF
}

func (p *Parser) peek() t.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() t.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(tokenType t.TokenType, message string) t.Token {
	if p.check(tokenType) {
		return p.advance()
	}
	p.error(p.peek(), message)
	return t.Token{}
}

type parseError struct{}

func (p *Parser) error(token t.Token, message string) {
	errs.ParseError(token, message)
	panic(parseError{})
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.IsAtEnd() {
		if p.previous().TokenType == t.SEMICOLON {
			return
		}

		switch p.peek().TokenType {
		case t.CLASS, t.FUN, t.VAR, t.FOR, t.IF, t.WHILE, t.PRINT, t.RETURN:
			return
		}

		p.advance()
	}
}
