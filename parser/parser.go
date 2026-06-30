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
declaration→ varDecl | funcDecl | classDecl | statement ;

classDecl→ "class" IDENTIFIER "{" function* "}" ;
varDecl→ "var" IDENTIFIER ( "=" expression )? ";" ;

funDecl→ "fun" function;
function→ IDENTIFIER funExpr ;
funcExpr→ "(" parameters? ")" block ;
parameters→ IDENTIFIER ( "," IDENTIFIER )* ;

statement→ exprStmt | printStmt | ifStmt | whileStmt | forStmt | breakStmt | returnStmt | continueStmt | block ;
forStmt→ "for" "(" ( varDecl | exprStmt | ";" )
		  expression? ";"
		  expression? ")" statement ;
whileStmt→ "while" "(" expression ")" statement ;
ifStmt→ "if" "(" expression ")" statement ( "else" statement )? ;
block→ "{" declaration* "}" ;
exprStmt→ expression ";" ;
printStmt→ "print" expression ";" ;
returnStmt→ "return" expression? ";" ;

expression→ assignment ;
assignment→ ( call "." )? IDENTIFIER "=" assignment | logic_or ;
logic_or→ logic_and ( "or" logic_and )* ;
logic_and→ ternary ( "and" ternary )* ;
ternary→ equality ( "?" expression ":" conditional )? ;
equality→ comparison ( ( "!=" | "==" ) comparison )* ;
comparison→ term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term→ factor ( ( "-" | "+" ) factor )* ;
factor→ unary ( ( "/" | "*" ) unary )* ;
unary→ ( "!" | "-" ) unary | call ;
call→ primary ( "(" arguments? ")" | "." IDENTIFIER )* ;
argument→ expression ( "," expression )* ;
primary→ NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" | IDENTIFIER | "fun" funExpr ;
*/

type Parser struct {
	tokens      []t.Token
	current     int
	loopNesting int
	returned    bool // return called, break or continue
}

func NewParser(tokens []t.Token) *Parser {
	return &Parser{
		tokens:      tokens,
		current:     0,
		loopNesting: 0,
		returned:    false,
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
			fmt.Printf("%v \n", err)
			errs.HadError = true
			p.synchronize()
		}
	}()
	if p.match(t.VAR) {
		return p.varDeclaration()
	}
	if p.check(t.FUN) && p.checkNext(t.IDENTIFIER) {
		p.consume(t.FUN, "")
		return p.funDeclaration("function")
	}
	if p.match(t.CLASS) {
		return p.class()
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

func (p *Parser) funDeclaration(kind string) ast.Stmt {
	name := p.consume(t.IDENTIFIER, fmt.Sprintf("Expect a %v name.", kind))

	return ast.NewFunc(name, p.funExpr(kind))
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
	if p.match(t.WHILE) {
		p.loopNesting++
		defer func() {
			p.loopNesting--
		}()
		return p.whileStatement()
	}
	if p.match(t.FOR) {
		p.loopNesting++
		defer func() {
			p.loopNesting--
		}()
		return p.forStatement()
	}
	if p.match(t.BREAK) {
		return p.breakStatement()
	}
	if p.match(t.CONTINUE) {
		return p.continueStatement()
	}
	if p.match(t.RETURN) {
		return p.returnStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) forStatement() ast.Stmt {
	p.consume(t.LEFT_PAREN, "Expect '(' after 'for'.")

	var initializer ast.Stmt
	if p.match(t.VAR) {
		initializer = p.varDeclaration()
	} else if p.match(t.SEMICOLON) {
		initializer = nil
	} else {
		initializer = p.expressionStatement()
	}

	var condition ast.Expr
	if !p.check(t.SEMICOLON) {
		condition = p.expression()
	}
	p.consume(t.SEMICOLON, "Expect ';' after loop condition.")

	var increment ast.Expr
	if !p.check(t.RIGHT_PAREN) {
		increment = p.expression()
	}
	p.consume(t.RIGHT_PAREN, "Expect ')' after for clause.")

	body := p.statement()

	if condition == nil {
		condition = ast.NewLiteral(true)
	}
	body = ast.NewWhileStmt(condition, increment, body)

	if initializer != nil {
		body = ast.NewBlock([]ast.Stmt{
			initializer,
			body,
		})
	}

	return body
}

func (p *Parser) whileStatement() ast.Stmt {
	p.consume(t.LEFT_PAREN, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(t.RIGHT_PAREN, "Expect ')' after while conditional.")

	body := p.statement()
	return ast.NewWhileStmt(condition, nil, body)
}

func (p *Parser) printStatement() ast.Stmt {
	expr := p.expression()
	p.consume(t.SEMICOLON, "Expect ';' after value.")
	return ast.NewPrintStmt(expr)
}

func (p *Parser) block() []ast.Stmt {
	var statements []ast.Stmt
	for !p.check(t.RIGHT_BRACE) && !p.isAtEnd() {
		if p.returned {
			errs.ReportWarning(p.peek(), "Un-reachable code.")
			p.returned = false
		}

		statements = append(statements, p.declaration())
	}

	p.consume(t.RIGHT_BRACE, "Expect '}' after block.")
	return statements
}

func (p *Parser) class() ast.Stmt {
	name := p.consume(t.IDENTIFIER, "Expect a class name.")
	p.consume(t.LEFT_BRACE, "Expect '{' before class body.")

	var methods []ast.FuncStmt
	for !p.check(t.RIGHT_BRACE) && !p.isAtEnd() {
		method := p.funDeclaration("method").(*ast.FuncStmt)
		methods = append(methods, *method)
	}

	p.consume(t.RIGHT_BRACE, "Expect '}' after class body.")

	return ast.NewClassStmt(name, methods)
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

	p.consume(t.LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(t.RIGHT_PAREN, "Expect ')' after if condition.")

	thenStmt := p.statement()
	if p.match(t.ELSE) {
		elseStmt = p.statement()
	}

	return ast.NewIfStmt(condition, thenStmt, elseStmt)
}

func (p *Parser) breakStatement() ast.Stmt {
	breakKeyword := p.previous()
	p.consume(t.SEMICOLON, "Expect ';' after statement.")
	if p.loopNesting == 0 {
		p.error(breakKeyword, "Cannot use 'break' statement outside of a loop.")
	}

	p.returned = true

	return ast.NewBreakStmt()
}

func (p *Parser) continueStatement() ast.Stmt {
	continueKeyword := p.previous()
	p.consume(t.SEMICOLON, "Expect ';' after statement.")
	if p.loopNesting == 0 {
		p.error(continueKeyword, "Cannot use 'continue' statement outside of a loop.")
	}

	p.returned = true

	return ast.NewContinueStmt()
}

func (p *Parser) returnStatement() ast.Stmt {
	returnKeyword := p.previous()
	var value ast.Expr
	if !p.check(t.SEMICOLON) {
		value = p.expression()
	}

	p.consume(t.SEMICOLON, "Expect ';' after statement.")

	p.returned = true

	return ast.NewReturn(returnKeyword, value)
}

func (p *Parser) assignment() ast.Expr {
	expr := p.or()
	switch {
	case p.match(t.EQUAL):
		equals := p.previous()
		value := p.assignment()
		if variable, ok := expr.(*ast.VarExpr); ok {
			return ast.NewAssignment(variable.Name, value)
		} else if get, ok := expr.(*ast.Get); ok {
			return ast.NewSet(get.Object, value, get.Name)
		}

		p.error(equals, "Invalid assignment target.")
	case p.match(t.PLUS_EQUAL):
		equals := p.previous()
		value := p.assignment()
		var addition *ast.Binary
		if variable, ok := expr.(*ast.VarExpr); ok {
			addition = ast.NewBinary(variable, value, *t.NewToken(t.PLUS, "+", "", equals.Line))
			return ast.NewAssignment(variable.Name, addition)
		} else if get, ok := expr.(*ast.Get); ok {
			return ast.NewSet(get.Object, addition, get.Name)
		}

		p.error(equals, "Invalid assignment target.")
	case p.match(t.MINUS_EQUAL):
		equals := p.previous()
		value := p.assignment()
		var subtraction *ast.Binary
		if variable, ok := expr.(*ast.VarExpr); ok {
			subtraction = ast.NewBinary(variable, value, *t.NewToken(t.MINUS, "-", "", equals.Line))
			return ast.NewAssignment(variable.Name, subtraction)
		} else if get, ok := expr.(*ast.Get); ok {
			return ast.NewSet(get.Object, subtraction, get.Name)
		}

		p.error(equals, "Invalid assignment target.")
	case p.match(t.STAR_EQUAL):
		equals := p.previous()
		value := p.assignment()
		var multiplication *ast.Binary
		if variable, ok := expr.(*ast.VarExpr); ok {
			multiplication = ast.NewBinary(variable, value, *t.NewToken(t.STAR, "*", "", equals.Line))
			return ast.NewAssignment(variable.Name, multiplication)
		} else if get, ok := expr.(*ast.Get); ok {
			return ast.NewSet(get.Object, multiplication, get.Name)
		}

		p.error(equals, "Invalid assignment target.")
	case p.match(t.SLASH_EQUAL):
		equals := p.previous()
		value := p.assignment()
		var division *ast.Binary
		if variable, ok := expr.(*ast.VarExpr); ok {
			division = ast.NewBinary(variable, value, *t.NewToken(t.SLASH, "/", "", equals.Line))
			return ast.NewAssignment(variable.Name, division)
		} else if get, ok := expr.(*ast.Get); ok {
			return ast.NewSet(get.Object, division, get.Name)
		}

		p.error(equals, "Invalid assignment target.")
	}
	return expr
}

func (p *Parser) or() ast.Expr {
	expr := p.and()
	for p.match(t.OR) {
		operator := p.previous()
		right := p.and()
		expr = ast.NewLogical(expr, right, operator)
	}
	return expr
}

func (p *Parser) and() ast.Expr {
	expr := p.ternary()
	for p.match(t.AND) {
		operator := p.previous()
		right := p.ternary()
		expr = ast.NewLogical(expr, right, operator)
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

	return p.call()
}

func (p *Parser) call() ast.Expr {
	expr := p.primary()

	for {
		if p.match(t.LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else if p.match(t.DOT) {
			name := p.consume(t.IDENTIFIER, "Expect property name after '.'.")
			expr = ast.NewGet(expr, name)
		} else {
			break
		}
	}
	return expr
}

func (p *Parser) finishCall(callee ast.Expr) ast.Expr {
	var arguments []ast.Expr
	if !p.check(t.RIGHT_PAREN) {
		arguments = append(arguments, p.expression())
		for p.match(t.COMMA) {
			if len(arguments) >= 255 {
				errs.ParseError(p.peek(), "Can't have more than 255 arguments")
			}
			arguments = append(arguments, p.expression())
		}
	}

	rightParen := p.consume(t.RIGHT_PAREN, "Expect ')' after arguments.")
	return ast.NewCall(callee, rightParen, arguments)
}

func (p *Parser) funExpr(kind string) ast.FuncExpr {
	p.consume(t.LEFT_PAREN, "Expect a '(' after function name.")

	var parameters []t.Token
	if !p.check(t.RIGHT_PAREN) {
		parameters = append(parameters, p.consume(t.IDENTIFIER, "Expect parameter name."))
		for p.match(t.COMMA) {
			if len(parameters) >= 255 {
				errs.ParseError(p.peek(), "Can't have more than 255 parameters")
			}
			parameters = append(parameters, p.consume(t.IDENTIFIER, "Expect parameter name."))
		}
	}
	p.consume(t.RIGHT_PAREN, "Expect ')' after parameters.")

	p.consume(t.LEFT_BRACE, fmt.Sprintf("Expect '{' before %v body.", kind))
	block := p.block()

	return *ast.NewFuncExpr(parameters, block)
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
	case p.match(t.THIS):
		return ast.NewThis(p.previous())
	case p.match(t.IDENTIFIER):
		return ast.NewVariable(p.previous())
	case p.match(t.LEFT_PAREN):
		expr := p.expression()
		p.consume(t.RIGHT_PAREN, "Expect ')' after expression.")
		return ast.NewGrouping(expr)
	case p.match(t.FUN):
		return p.funExpr("function")
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

func (p *Parser) checkNext(tokenType t.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.tokens[p.current+1].TokenType == tokenType
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
	errMsg := errs.ParseError(token, message)
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
