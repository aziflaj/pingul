package parser

import (
	"fmt"

	"github.com/aziflaj/pingul/ast"
	"github.com/aziflaj/pingul/lexer"
	"github.com/aziflaj/pingul/token"
)

type Parser struct {
	lxr *lexer.Lexer

	currentToken token.Token
	peekToken    token.Token

	errors []string
}

func New(lxr *lexer.Lexer) *Parser {
	p := &Parser{lxr: lxr}

	// set both currentToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) ParseProgram() *ast.Program {
	prg := &ast.Program{
		Statements: []ast.Statement{},
	}

	for p.currentToken.Type != token.EOF {
		stmt := p.parseStatement()

		if stmt != nil {
			prg.Statements = append(prg.Statements, stmt)
		}

		p.nextToken()
	}

	return prg
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lxr.NextToken()
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.VAR:
		return p.parseVarStatement()
	default:
		return nil
	}
}

func (p *Parser) parseVarStatement() *ast.VarStatement {
	fmt.Println("Parsing var statement")
	fmt.Println("Current token:", p.currentToken)
	stmt := &ast.VarStatement{Token: p.currentToken}

	if p.peekToken.Type != token.IDENTIFIER {
		return nil
	}

	p.nextToken()
	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if p.peekToken.Type != token.ASSIGNMENT {
		p.peekError(token.Token{Type: token.ASSIGNMENT})
		return nil
	}

	// TODO: parse the expression

	// read until the end of the statement
	for p.currentToken.Type != token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) peekError(expected token.Token) {
	msg := fmt.Sprintf("Expected next token to be %s, got %s instead",
		expected, p.peekToken)
	p.errors = append(p.errors, msg)
}
