package ast

import (
	"fmt"

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

// var <identifier> = <expression>;
type VarStatement struct {
	Token token.Token // the token.VAR token
	Name  *Identifier
	Value Expression
}

func (s *VarStatement) statementNode() {} // because Types and stuff
func (s *VarStatement) TokenLiteral() []rune {
	fmt.Println(s.Token)
	return s.Token.Literal
}
