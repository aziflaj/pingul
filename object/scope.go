package object

type Scope struct {
	table map[string]Object

	// All scopes are local, except the global scope
	// which is the outermost scope
	outter *Scope
}

func NewScope() *Scope {
	return &Scope{
		table:  make(map[string]Object),
		outter: nil,
	}
}

func NewLocalScope(outter *Scope) *Scope {
	s := NewScope()
	s.outter = outter

	return s
}

func (s *Scope) Get(name string) Object {
	obj, ok := s.table[name]

	if !ok {
		if s.outter != nil {
			return s.outter.Get(name)
		}

		return s.Set(name, &Nil{})
	}

	return obj
}

// Always set on the local scope
func (s *Scope) Set(name string, obj Object) Object {
	s.table[name] = obj
	return obj
}
