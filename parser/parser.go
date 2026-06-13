package parser

import (
	"fmt"

	"github.com/ab-dek/golox/ast"
	errs "github.com/ab-dek/golox/errors"
	t "github.com/ab-dek/golox/token"
)

/*
grammmer:

program→ declaration* EOF ;
declaration→ varDecl
			 | statement ;
varDecl→ "var" IDENTIFIER ( "=" expression )? ";" ;
statement→ exprStmt
		   | printStmt
		   | ifStmt
		   | block ;
ifStmt→ "if" "(" expression ")" statement ( "else" statement )? ;
block→ "{" declaration* "}" ;
exprStmt→ expression ";" ;
printStmt→ "print" expression ";" ;
expression→ assignment ;
assignment→ IDENTIFIER "=" assignment
			| ternary ;
ternary→ equality ( "?" expression ":" conditional )? ;
equality→ comparison ( ( "!=" | "==" ) comparison )* ;
comparison→ term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term→ factor ( ( "-" | "+" ) factor )* ;
factor→ unary ( ( "/" | "*" ) unary )* ;
unary→ ( "!" | "-" ) unary
	   | primary ;
primary→ NUMBER | STRING | "true" | "false" | "nil"
		 | "(" expression ")"
| IDENTIFIER ;
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
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	return statements
}

func (p *Parser) ParseExpr() ast.Expr {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("%v \n", err)
		}
	}()
	return p.expression()
}

func (p *Parser) declaration() ast.Stmt {
	defer func() {
		if err := recover(); err != nil {
			errs.HadError = true
			p.synchronize()
		}
	}()
	if p.match(t.VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) varDeclaration() ast.Stmt {
	name := p.consume(t.IDENTIFIER, "Expect a variable name.")
	var initializer ast.Expr
	if p.match(t.EQUAL) {
		initializer = p.expression()
	}
	p.consume(t.SEMICOLON, "Expect ';' after variable declaration.")
	return ast.NewVarStmt(name, initializer)
}

func (p *Parser) statement() ast.Stmt {
	if p.match(t.PRINT) {
		return p.printStatement()
	}
	if p.match(t.LEFT_BRACE) {
		return ast.NewBlock(p.block())
	}
	if p.match(t.IF) {
		return p.ifStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) printStatement() ast.Stmt {
	expr := p.expression()
	p.consume(t.SEMICOLON, "Expect ';' after value.")
	return ast.NewPrintStmt(expr)
}

func (p *Parser) block() []ast.Stmt {
	var statements []ast.Stmt
	for !p.check(t.RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	p.consume(t.RIGHT_BRACE, "Expect '}' after block.")
	return statements
}

func (p *Parser) expressionStatement() ast.Stmt {
	expr := p.expression()
	p.consume(t.SEMICOLON, "Expect ';' after value.")
	return ast.NewExprStmt(expr)
}

func (p *Parser) expression() ast.Expr {
	return p.assignment()
}

func (p *Parser) ifStatement() ast.Stmt {
	var elseStmt ast.Stmt

	p.consume(t.LEFT_PAREN, "Expect '(' after if conditional expression")
	condition := p.expression()
	p.consume(t.RIGHT_PAREN, "Expect ')' after if conditional expression")

	thenStmt := p.statement()
	if p.match(t.ELSE) {
		elseStmt = p.statement()
	}

	return ast.NewIfStmt(condition, thenStmt, elseStmt)
}

func (p *Parser) assignment() ast.Expr {
	expr := p.ternary()
	switch {
	case p.match(t.EQUAL):
		equals := p.previous()
		value := p.assignment()
		if variable, ok := expr.(*ast.Variable); ok {
			return ast.NewAssignment(variable.Name, value)
		}

		p.error(equals, "Invalid assignment target.")
	case p.match(t.PLUS_EQUAL):
		equals := p.previous()
		value := p.assignment()
		if variable, ok := expr.(*ast.Variable); ok {
			addition := ast.NewBinary(variable, value, *t.NewToken(t.PLUS, "+", "", equals.Line))
			return ast.NewAssignment(variable.Name, addition)
		}

		p.error(equals, "Invalid assignment target.")
	case p.match(t.MINUS_EQUAL):
		equals := p.previous()
		value := p.assignment()
		if variable, ok := expr.(*ast.Variable); ok {
			subtraction := ast.NewBinary(variable, value, *t.NewToken(t.MINUS, "-", "", equals.Line))
			return ast.NewAssignment(variable.Name, subtraction)
		}

		p.error(equals, "Invalid assignment target.")
	case p.match(t.STAR_EQUAL):
		equals := p.previous()
		value := p.assignment()
		if variable, ok := expr.(*ast.Variable); ok {
			multiplication := ast.NewBinary(variable, value, *t.NewToken(t.STAR, "*", "", equals.Line))
			return ast.NewAssignment(variable.Name, multiplication)
		}

		p.error(equals, "Invalid assignment target.")
	case p.match(t.SLASH_EQUAL):
		equals := p.previous()
		value := p.assignment()
		if variable, ok := expr.(*ast.Variable); ok {
			addition := ast.NewBinary(variable, value, *t.NewToken(t.SLASH, "/", "", equals.Line))
			return ast.NewAssignment(variable.Name, addition)
		}

		p.error(equals, "Invalid assignment target.")
	}
	return expr
}

func (p *Parser) ternary() ast.Expr {
	expr := p.equality()
	if p.match(t.QUESTION_MARK) {
		questionMark := p.previous()
		thenBranch := p.expression()

		p.consume(t.COLON, "Expect ':' after then-branch of conditional operator.")

		elseBranch := p.ternary()
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

	for p.match(t.STAR, t.SLASH, t.PERCENT) {
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
	case p.match(t.IDENTIFIER):
		return ast.NewVariable(p.previous())
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
	if p.isAtEnd() {
		return false
	}
	return p.peek().TokenType == tokenType
}

func (p *Parser) advance() t.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
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

func (p *Parser) error(token t.Token, message string) {
	errMsg := errs.ReportParseError(token, message)
	panic(errMsg)
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
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
