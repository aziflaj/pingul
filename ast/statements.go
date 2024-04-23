package ast

import (
	"strings"

	"github.com/aziflaj/pingul/token"
)

type Identifier struct {
	Token token.Token // the token.IDENTIFIER token
	Value []rune
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() []rune {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return string(i.Value)
}

// var <identifier> = <expression>;
type VarStatement struct {
	Token token.Token // the token.VAR token
	Name  *Identifier
	Value Expression
}

func (s *VarStatement) statementNode() {} // because Types and stuff
func (s *VarStatement) TokenLiteral() []rune {
	return s.Token.Literal
}

func (s *VarStatement) String() string {
	var b strings.Builder

	b.WriteString(string(s.TokenLiteral()))
	b.WriteString(" ")
	b.WriteString(s.Name.String())
	b.WriteString(" = ")

	if s.Value != nil {
		b.WriteString(s.Value.String())
	}

	b.WriteString(";")

	return b.String()
}

// return <expression>;
type ReturnStatement struct {
	Token       token.Token // the token.RETURN token
	ReturnValue Expression
}

func (s *ReturnStatement) statementNode() {}
func (s *ReturnStatement) TokenLiteral() []rune {
	return s.Token.Literal
}

func (s *ReturnStatement) String() string {
	var b strings.Builder

	b.WriteString(string(s.TokenLiteral()))
	b.WriteString(" ")

	if s.ReturnValue != nil {
		b.WriteString(s.ReturnValue.String())
	}

	b.WriteString(";")

	return b.String()
}

// <expression>;
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (s *ExpressionStatement) statementNode() {}
func (s *ExpressionStatement) TokenLiteral() []rune {
	return s.Token.Literal
}

func (s *ExpressionStatement) String() string {
	if s.Expression != nil {
		return s.Expression.String()
	}

	return ""
}
