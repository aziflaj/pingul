package ast

import "strings"

type Node interface {
	TokenLiteral() []rune
	String() string
}

// Something that can be executed
type Statement interface {
	Node
	statementNode()
}

// Something that can be evaluated to a value
type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() []rune {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return []rune("")
}

func (p *Program) String() string {
	var b strings.Builder

	for _, s := range p.Statements {
		b.WriteString(s.String())
	}

	return b.String()
}
