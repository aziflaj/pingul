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

	token.LPAREN: FUNCALL,
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
		token.TRUE:       p.parseBoolean,
		token.FALSE:      p.parseBoolean,
		token.LPAREN:     p.parseGroupExpressions,
		token.IF:         p.parseIfExpression,
		token.FUNC:       p.parseFuncExpression,
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

		token.LPAREN: p.parseFuncCall,
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
	p.nextToken()

	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currentToken}

	p.nextToken()
	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekToken.Type == token.SEMICOLON {
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

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.currentToken,
		Value: p.currentToken.Type == token.TRUE,
	}
}

func (p *Parser) parseGroupExpressions() ast.Expression {
	p.nextToken()

	parsedExpr := p.parseExpression(LOWEST)

	if p.peekToken.Type != token.RPAREN {
		return nil
	}

	return parsedExpr
}

func (p *Parser) parseIfExpression() ast.Expression {
	expr := &ast.IfExpression{Token: p.currentToken}

	if p.peekToken.Type != token.LPAREN {
		p.errors = append(p.errors, "Expected '(' after 'if'")
		return nil
	}

	p.nextToken()
	expr.Condition = p.parseExpression(LOWEST)

	if p.peekToken.Type != token.RPAREN {
		p.errors = append(p.errors, "Expected ')' after condition")
		return nil
	}
	p.nextToken()

	if p.peekToken.Type != token.LBRACE {
		p.errors = append(p.errors, "Expected '{' after condition")
		return nil
	}

	p.nextToken()
	expr.Consequence = p.parseBlockStatement()

	if p.peekToken.Type == token.ELSE {
		p.nextToken()

		if p.peekToken.Type != token.LBRACE {
			p.errors = append(p.errors, "Expected '{' after 'else'")
			return nil
		}

		p.nextToken()
		expr.Alternative = p.parseBlockStatement()
	}

	return expr
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currentToken}

	p.nextToken()
	// read until the end of the block
	for p.currentToken.Type != token.RBRACE {
		stmt := p.parseStatement()

		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		p.nextToken()
	}

	return block
}

func (p *Parser) parseFuncExpression() ast.Expression {
	expr := &ast.FuncExpression{Token: p.currentToken}

	if p.peekToken.Type != token.LPAREN {
		p.errors = append(p.errors, "Expected '(' after 'func'")
		return nil
	}

	p.nextToken()
	expr.Params = p.parseFuncParams()

	if p.peekToken.Type != token.LBRACE {
		p.errors = append(p.errors, "Expected '{' after parameters")
		return nil
	}

	p.nextToken()
	expr.Body = p.parseBlockStatement()

	return expr
}

func (p *Parser) parseFuncParams() []*ast.Identifier {
	params := []*ast.Identifier{}

	if p.peekToken.Type == token.RPAREN {
		p.nextToken()
		return params
	}

	p.nextToken()
	params = append(params, &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	})

	for p.peekToken.Type == token.COMMA {
		p.nextToken() // skip the comma
		p.nextToken() // move to the next param

		params = append(params, &ast.Identifier{
			Token: p.currentToken,
			Value: p.currentToken.Literal,
		})
	}

	if p.peekToken.Type != token.RPAREN {
		p.errors = append(p.errors, "Expected ')' after parameters")
		return nil
	}

	p.nextToken()

	return params
}

func (p *Parser) parseFuncCall(left ast.Expression) ast.Expression {
	expr := &ast.CallExpression{Token: p.currentToken, Function: left}
	expr.Arguments = p.parseFuncCallArgs()

	return expr
}

func (p *Parser) parseFuncCallArgs() []ast.Expression {
	args := []ast.Expression{}

	if p.peekToken.Type == token.RPAREN {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekToken.Type == token.COMMA {
		p.nextToken() // skip the comma
		p.nextToken() // move to the next arg
		args = append(args, p.parseExpression(LOWEST))
	}

	if p.peekToken.Type != token.RPAREN {
		p.errors = append(p.errors, "Expected ')' after arguments")
		return nil
	}

	p.nextToken()

	return args
}
