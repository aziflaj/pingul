package ast

type Node interface {
	TokenLiteral() []rune
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
