package parser

import (
	"fmt"
	"strconv"

	"github.com/aziflaj/pingul/ast"
	"github.com/aziflaj/pingul/lexer"
	"github.com/aziflaj/pingul/token"
)

// Operator Predecence levels
type OpPrecedence uint8

const (
	LOWEST  = OpPrecedence(iota)
	EQUALS  // ==
	LGT     // less, greater than, lte, gte
	SUM     // +
	PRODUCT // *
	PREFIX  // `-x`, `or x`
	FUNCALL // function call
)

type (
	prefixParseHandler func() ast.Expression
	infixParseHandler  func(ast.Expression) ast.Expression
)

type Parser struct {
	lxr *lexer.Lexer

	currentToken token.Token
	peekToken    token.Token

	errors []string

	prefixParseHandlers map[token.TokenType]prefixParseHandler
	infixParseHandlers  map[token.TokenType]infixParseHandler
}

func New(lxr *lexer.Lexer) *Parser {
	p := &Parser{lxr: lxr}

	// set both currentToken and peekToken
	p.nextToken()
	p.nextToken()

	p.prefixParseHandlers = map[token.TokenType]prefixParseHandler{
		token.IDENTIFIER: p.parseIdentifier,
		token.INT:        p.parseIntegerLiteral,
	}

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

	case token.RETURN:
		return p.parseReturnStatement()

	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseVarStatement() *ast.VarStatement {
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

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currentToken}

	// TODO: parse the expression

	// read until the end of the statement
	for p.currentToken.Type != token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{
		Token:      p.currentToken,
		Expression: p.parseExpression(LOWEST),
	}

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

// Helpers

func (p *Parser) peekError(expected token.Token) {
	msg := fmt.Sprintf("Expected next token to be %s, got %s instead",
		expected, p.peekToken)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(op OpPrecedence) ast.Expression {
	prefixHandler := p.prefixParseHandlers[p.currentToken.Type]

	if prefixHandler == nil {
		return nil
	}

	return prefixHandler()
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	val, err := strconv.ParseInt(string(p.currentToken.Literal), 0, 64)
	if err != nil {
		msg := fmt.Sprintf("Could not parse %q as integer", p.currentToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	return &ast.IntegerLiteral{Token: p.currentToken, Value: val}
}
