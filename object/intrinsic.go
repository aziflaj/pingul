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
}
