package object

import "fmt"

type FuncTable map[string]IntrinsicFunc

var IntrinsicFuncs = FuncTable{
	"print": func(args ...Object) Object {
		for _, arg := range args {
			fmt.Println(arg.Inspect())
		}
		return &Nil{}
	},
	"len": func(args ...Object) Object {
		if len(args) != 1 {
			return &Nil{}
		}

		switch arg := args[0].(type) {
		case *String:
			return &Integer{Value: int64(len(arg.Value))}
		case *List:
			return &Integer{Value: int64(len(arg.Items))}
		default:
			return &Nil{}
		}
	},
}
