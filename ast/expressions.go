package ast

import (
	"strings"

	"github.com/aziflaj/pingul/token"
)

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

// <expression>[<expression>]
type IndexExpression struct {
	Token token.Token // the '[' token
	List  Expression
	Index Expression
}

func (i *IndexExpression) expressionNode() {}
func (i *IndexExpression) TokenLiteral() []rune {
	return i.Token.Literal
}
func (i *IndexExpression) String() string {
	var b strings.Builder

	b.WriteString("(")
	b.WriteString(i.List.String())
	b.WriteString("[")
	b.WriteString(i.Index.String())
	b.WriteString("])")

	return b.String()
}

// if (<expression>) <block> else <block>
type IfExpression struct {
	Token       token.Token // the 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (i *IfExpression) expressionNode() {}
func (i *IfExpression) TokenLiteral() []rune {
	return i.Token.Literal
}
func (i *IfExpression) String() string {
	var b strings.Builder

	b.WriteString("if ")
	b.WriteString(i.Condition.String())
	b.WriteString(" ")
	b.WriteString(i.Consequence.String())

	if i.Alternative != nil {
		b.WriteString(" else ")
		b.WriteString(i.Alternative.String())
	}

	return b.String()
}

// func(<params>) { <block> }
type FuncExpression struct {
	Token  token.Token // the 'func' token
	Params []*Identifier
	Body   *BlockStatement
}

func (f *FuncExpression) expressionNode() {}
func (f *FuncExpression) TokenLiteral() []rune {
	return f.Token.Literal
}
func (f *FuncExpression) String() string {
	var b strings.Builder

	b.WriteString("func(")

	for i, param := range f.Params {
		b.WriteString(param.String())

		if i < len(f.Params)-1 {
			b.WriteString(", ")
		}
	}

	b.WriteString(") {")
	b.WriteString(f.Body.String())
	b.WriteString("}")

	return b.String()
}

// <expression>(<arguments>)
type CallExpression struct {
	Token     token.Token // the '(' token
	Function  Expression  // Identifier or FuncExpression
	Arguments []Expression
}

func (c *CallExpression) expressionNode() {}
func (c *CallExpression) TokenLiteral() []rune {
	return c.Token.Literal
}
func (c *CallExpression) String() string {
	var b strings.Builder

	b.WriteString(c.Function.String())
	b.WriteString("(")

	for i, arg := range c.Arguments {
		b.WriteString(arg.String())

		if i < len(c.Arguments)-1 {
			b.WriteString(", ")
		}
	}

	b.WriteString(")")

	return b.String()
}
