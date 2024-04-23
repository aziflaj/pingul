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

// - <expression>
// not <expression>
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
	b.WriteString("(")
	b.WriteString(p.Right.String())
	b.WriteString(")")
	b.WriteString(")")

	return b.String()
}

// <expression> <infix operator> <expression>
type InfixExpression struct {
	Token    token.Token // the operator token
	Left     Expression
	Operator string
	Right    Expression
}

func (i *InfixExpression) expressionNode() {}
func (i *InfixExpression) TokenLiteral() []rune {
	return i.Token.Literal
}
func (i *InfixExpression) String() string {
	var b strings.Builder

	b.WriteString("(")
	b.WriteString(i.Left.String())
	b.WriteString(" ")
	b.WriteString(i.Operator)
	b.WriteString(" ")
	b.WriteString(i.Right.String())
	b.WriteString(")")

	return b.String()
}
