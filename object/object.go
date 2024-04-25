package object

import (
	"fmt"
	"strings"

	"github.com/aziflaj/pingul/ast"
)

type ObjectType string

const (
	INT    = ObjectType("INT")
	BOOL   = ObjectType("BOOL")
	NIL    = ObjectType("NIL")
	RETURN = ObjectType("RETURN")
	FUNC   = ObjectType("FUNC")
)

type Object interface {
	Type() ObjectType
	Inspect() string

	IsTruthy() bool
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INT }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%s(%d)", i.Type(), i.Value) }
func (i *Integer) IsTruthy() bool   { return i.Value != 0 }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOL }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%s(%t)", b.Type(), b.Value) }
func (b *Boolean) IsTruthy() bool   { return b.Value }

type Nil struct{}

func (n *Nil) Type() ObjectType { return NIL }
func (n *Nil) Inspect() string  { return string(n.Type()) }
func (n *Nil) IsTruthy() bool   { return false }

type Return struct {
	Value Object
}

func (r *Return) Type() ObjectType { return RETURN }
func (r *Return) Inspect() string  { return r.Value.Inspect() }
func (r *Return) IsTruthy() bool   { return r.Value.IsTruthy() }

type Func struct {
	Params []*ast.Identifier
	Body   *ast.BlockStatement
}

func (f *Func) Type() ObjectType { return FUNC }
func (f *Func) IsTruthy() bool   { return true }
func (f *Func) Inspect() string {
	var b strings.Builder

	b.WriteString("func(")
	for i, p := range f.Params {
		b.WriteString(p.String())
		if i < len(f.Params)-1 {
			b.WriteString(", ")
		}
	}
	b.WriteString(") {\n")
	b.WriteString(f.Body.String())
	b.WriteString("\n}")

	return b.String()
}
