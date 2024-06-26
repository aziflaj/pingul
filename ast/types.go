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

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() []rune {
	return b.Token.Literal
}
func (b *Boolean) String() string {
	return string(b.Token.Literal)
}

type String struct {
	Token token.Token
	Value []rune
}

func (s *String) expressionNode() {}
func (s *String) TokenLiteral() []rune {
	return s.Token.Literal
}
func (s *String) String() string {
	return string(s.Token.Literal)
}

type Nil struct {
	Token token.Token
}

func (n *Nil) expressionNode() {}
func (n *Nil) TokenLiteral() []rune {
	return n.Token.Literal
}
func (n *Nil) String() string {
	return string(n.Token.Literal)
}

type List struct {
	Token token.Token
	Items []Expression
}

func (l *List) expressionNode() {}
func (l *List) TokenLiteral() []rune {
	return l.Token.Literal
}
func (l *List) String() string {
	var b strings.Builder

	b.WriteString("[")
	for i, item := range l.Items {
		b.WriteString(item.String())

		if i < len(l.Items)-1 {
			b.WriteString(", ")
		}
	}
	b.WriteString("]")

	return b.String()
}
