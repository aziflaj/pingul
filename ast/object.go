package ast

import (
	"strings"

	"github.com/aziflaj/pingul/token"
)

// { key: value, key2: value2, ... }
type ObjectLiteral struct {
	Token token.Token      // the '{' token
	Pairs map[string]Expression
}

func (o *ObjectLiteral) expressionNode() {}
func (o *ObjectLiteral) TokenLiteral() []rune {
	return o.Token.Literal
}
func (o *ObjectLiteral) String() string {
	var b strings.Builder

	b.WriteString("{")
	first := true
	for key, value := range o.Pairs {
		if !first {
			b.WriteString(", ")
		}
		first = false
		b.WriteString(key)
		b.WriteString(": ")
		b.WriteString(value.String())
	}
	b.WriteString("}")

	return b.String()
}

// <expression>.<property>
type PropertyAccess struct {
	Token    token.Token // the '.' token
	Object   Expression
	Property string // property name
}

func (p *PropertyAccess) expressionNode() {}
func (p *PropertyAccess) TokenLiteral() []rune {
	return p.Token.Literal
}
func (p *PropertyAccess) String() string {
	var b strings.Builder

	b.WriteString("(")
	b.WriteString(p.Object.String())
	b.WriteString(".")
	b.WriteString(p.Property)
	b.WriteString(")")

	return b.String()
}
