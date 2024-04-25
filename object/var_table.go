package object

type VarTable struct {
	table map[string]Object
}

func NewVarTable() *VarTable {
	return &VarTable{table: make(map[string]Object)}
}

func (vt *VarTable) Get(name string) Object {
	obj, ok := vt.table[name]

	if !ok {
		return vt.Set(name, &Nil{})
	}

	return obj
}

func (vt *VarTable) Set(name string, obj Object) Object {
	vt.table[name] = obj
	return obj
}
