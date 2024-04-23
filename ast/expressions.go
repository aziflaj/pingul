package ast

import (
	"strings"

	"github.com/aziflaj/pingul/token"
)

type IntegerLiteral struct {
	Token token.Token // the token.INT token
	Value int64
}

func (i *IntegerLiteral) expressionNode() {}
func (i *IntegerLiteral) TokenLiteral() []rune {
	return i.Token.Literal
}

func (i *IntegerLiteral) String() string {
	return string(i.Token.Literal)
}

type PrefixExpression struct {
	Token    token.Token // the prefix token, e.g. `-` (negative sign) or `not`
	Operator string
	Right    Expression
}

func (p *PrefixExpression) expressionNode() {}
func (p *PrefixExpression) TokenLiteral() []rune {
	return p.Token.Literal
}
func (p *PrefixExpression) String() string {
	var b strings.Builder

	b.WriteString("(")
	b.WriteString(p.Operator)
	b.WriteString(p.Right.String())
	b.WriteString(")")

	return b.String()
}
