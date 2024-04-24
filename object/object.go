package object

import "fmt"

type ObjectType string

const (
	INT  = ObjectType("INT")
	BOOL = ObjectType("BOOL")
	NIL  = ObjectType("NIL")
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INT }
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%s(%d)", i.Type(), i.Value)
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOL }
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%s(%t)", b.Type(), b.Value)
}

type Nil struct{}

func (n *Nil) Type() ObjectType { return NIL }
func (n *Nil) Inspect() string  { return string(n.Type()) }
