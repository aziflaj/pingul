package parser

import (
	"github.com/aziflaj/pingul/ast"
	"github.com/aziflaj/pingul/lexer"
)

// Parser wraps the generated yacc parser and provides a compatible API
type Parser struct {
	lexer  *lexer.LexerImpl
	errors []string
}

// New creates a new Parser
func New(lxr *lexer.Lexer) *Parser {
	return &Parser{
		lexer:  lxr.Impl(),
		errors: []string{},
	}
}

// NewFromLexerImpl creates a parser from a LexerImpl (for yacc use)
func NewFromLexerImpl(impl *lexer.LexerImpl) *Parser {
	return &Parser{
		lexer:  impl,
		errors: []string{},
	}
}

// Errors returns the parse errors
func (p *Parser) Errors() []string {
	return p.errors
}

// ParseProgram parses the input and returns an AST program
func (p *Parser) ParseProgram() *ast.Program {
	parseErrors = []string{} // Reset global errors

	yaccLexer := &YaccLexer{
		impl:    p.lexer,
		program: nil,
	}

	yyParse(yaccLexer)

	// Deduplicate consecutive identical errors
	uniqueErrors := []string{}
	var lastError string
	for _, err := range parseErrors {
		if err != lastError {
			uniqueErrors = append(uniqueErrors, err)
			lastError = err
		}
	}

	p.errors = uniqueErrors

	if yaccLexer.program != nil {
		return yaccLexer.program
	}

	return &ast.Program{Statements: []ast.Statement{}}
}

// ParseFromString is a helper that creates a parser from a string input
func ParseFromString(input string) (*ast.Program, []string) {
	impl := lexer.NewLexerImpl(input)
	parser := NewFromLexerImpl(impl)
	program := parser.ParseProgram()
	return program, parser.Errors()
}
