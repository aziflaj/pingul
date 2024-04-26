package object

import (
	"fmt"
	"strings"

	"github.com/aziflaj/pingul/ast"
)

type ObjectType string

const (
	INT            = ObjectType("INT")
	BOOL           = ObjectType("BOOL")
	STRING         = ObjectType("STRING")
	LIST           = ObjectType("LIST")
	NIL            = ObjectType("NIL")
	RETURN         = ObjectType("RETURN")
	FUNC           = ObjectType("FUNC")
	INTRINSIC_FUNC = ObjectType("INTRINSIC_FUNC")
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

type String struct {
	Value []rune
}

func (s *String) Type() ObjectType { return STRING }
func (s *String) Inspect() string  { return fmt.Sprintf("%s(%s)", s.Type(), string(s.Value)) }
func (s *String) IsTruthy() bool   { return string(s.Value) != "" }

type List struct {
	Items []Object
}

func (l *List) Type() ObjectType { return LIST }
func (l *List) Inspect() string {
	var b strings.Builder

	b.WriteString("[")
	for i, e := range l.Items {
		b.WriteString(e.Inspect())
		if i < len(l.Items)-1 {
			b.WriteString(", ")
		}
	}
	b.WriteString("]")

	return b.String()
}
func (l *List) IsTruthy() bool { return len(l.Items) > 0 }

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

type IntrinsicFunc func(args ...Object) Object

func (i IntrinsicFunc) Type() ObjectType { return INTRINSIC_FUNC }
func (i IntrinsicFunc) IsTruthy() bool   { return true }
func (i IntrinsicFunc) Inspect() string  { return "func(...) { intrinsic }" }
