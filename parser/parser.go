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
	PREFIX  // `-x`, `not x`
	FUNCALL // function call
)

var precedences = map[token.TokenType]OpPrecedence{
	token.EQUAL:     EQUALS,
	token.NOT_EQUAL: EQUALS,

	token.GREATER_THAN:          LGT,
	token.LESS_THAN:             LGT,
	token.GREATER_THAN_OR_EQUAL: LGT,
	token.LESS_THAN_OR_EQUAL:    LGT,

	token.PLUS:  SUM,
	token.MINUS: SUM,

	token.MULTIPLY: PRODUCT,
	token.DIVIDE:   PRODUCT,
	token.MODULUS:  PRODUCT,
}

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
		token.NOT:        p.parsePrefixExpression,
		token.MINUS:      p.parsePrefixExpression,
	}

	p.infixParseHandlers = map[token.TokenType]infixParseHandler{
		token.EQUAL:     p.parseInfixExpression,
		token.NOT_EQUAL: p.parseInfixExpression,

		token.GREATER_THAN:          p.parseInfixExpression,
		token.LESS_THAN:             p.parseInfixExpression,
		token.GREATER_THAN_OR_EQUAL: p.parseInfixExpression,
		token.LESS_THAN_OR_EQUAL:    p.parseInfixExpression,

		token.PLUS:  p.parseInfixExpression,
		token.MINUS: p.parseInfixExpression,

		token.MULTIPLY: p.parseInfixExpression,
		token.DIVIDE:   p.parseInfixExpression,
		token.MODULUS:  p.parseInfixExpression,
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

/////////// Helpers ///////////

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lxr.NextToken()
}

func (p *Parser) currentPrecedence() OpPrecedence {
	if prec, ok := precedences[p.currentToken.Type]; ok {
		return prec
	}

	return LOWEST
}

func (p *Parser) peekPrecedence() OpPrecedence {
	if prec, ok := precedences[p.peekToken.Type]; ok {
		return prec
	}

	return LOWEST
}

func (p *Parser) peekError(expected token.Token) {
	msg := fmt.Sprintf("Expected next token to be %s, got %s instead",
		expected, p.peekToken)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(op OpPrecedence) ast.Expression {
	prefixHandler := p.prefixParseHandlers[p.currentToken.Type]

	if prefixHandler == nil {
		msg := fmt.Sprintf("No prefix parse function for %s found", p.currentToken)
		p.errors = append(p.errors, msg)
		return nil
	}
	leftExpr := prefixHandler()

	for p.peekToken.Type != token.SEMICOLON && op < p.peekPrecedence() {
		infixHandler := p.infixParseHandlers[p.peekToken.Type]
		if infixHandler == nil {
			return leftExpr
		}

		p.nextToken()
		leftExpr = infixHandler(leftExpr)
	}

	return leftExpr
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

func (p *Parser) parsePrefixExpression() ast.Expression {
	expr := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: string(p.currentToken.Literal),
	}

	p.nextToken()
	expr.Right = p.parseExpression(PREFIX)

	return expr
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{
		Token:    p.currentToken,
		Left:     left,
		Operator: string(p.currentToken.Literal),
	}

	precedence := p.currentPrecedence()
	p.nextToken()
	expr.Right = p.parseExpression(precedence)

	return expr
}
