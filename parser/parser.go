package parser

import (
	e "github.com/ab-dek/golox/ast"
	t "github.com/ab-dek/golox/token"
)

type Parser struct {
	tokens  []t.Token
	current int
}

func (p *Parser) expression() e.Expr {
	return p.equality()
}

func (p *Parser) equality() e.Expr {
	expr := p.comparision()
	// TODO: implement me
	return expr
}

func (p *Parser) comparision() e.Expr {
	expr := p.term()
	// TODO: implement me
	return expr
}

func (p *Parser) term() e.Expr {
	expr := p.factor()
	// TODO: implement me
	return expr
}

func (p *Parser) factor() e.Expr {
	expr := p.unary()
	// TODO: implement me
	return expr
}

func (p *Parser) unary() e.Expr {
	expr := p.primary()
	// TODO: implement me
	return expr
}

func (p *Parser) primary() e.Expr {
	expr := &e.Literal{Value: nil}
	// TODO: implement me
	return expr
}

func (p *Parser) match(types ...t.TokenType) bool {
	// TODO: implement me
	return false
}
