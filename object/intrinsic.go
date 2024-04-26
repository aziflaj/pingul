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
	"head": func(args ...Object) Object {
		if len(args) != 1 {
			return &Nil{}
		}

		if list, ok := args[0].(*List); ok {
			if len(list.Items) > 0 {
				return list.Items[0]
			}
		}

		return &Nil{}
	},
	"tail": func(args ...Object) Object {
		if len(args) != 1 {
			return &Nil{}
		}

		if list, ok := args[0].(*List); ok {
			if len(list.Items) > 0 {
				return &List{Items: list.Items[1:]}
			}
		}

		return &Nil{}
	},

	"append": func(args ...Object) Object {
		if len(args) != 2 {
			return &Nil{}
		}

		if list, ok := args[0].(*List); ok {
			return &List{Items: append(list.Items, args[1])}
		}

		return &Nil{}
	},

	"prepend": func(args ...Object) Object {
		if len(args) != 2 {
			return &Nil{}
		}

		if list, ok := args[0].(*List); ok {
			return &List{Items: append([]Object{args[1]}, list.Items...)}
		}

		return &Nil{}
	},

	"pop": func(args ...Object) Object {
		if len(args) != 1 {
			return &Nil{}
		}

		if list, ok := args[0].(*List); ok {
			length := len(list.Items)
			if length > 0 {
				popped := list.Items[length-1]
				list.Items = list.Items[:length-1]
				return popped
			}
		}

		return &Nil{}
	},

	"shift": func(args ...Object) Object {
		if len(args) != 1 {
			return &Nil{}
		}

		if list, ok := args[0].(*List); ok {
			length := len(list.Items)
			if length > 0 {
				shifted := list.Items[0]
				list.Items = list.Items[1:]
				return shifted
			}
		}

		return &Nil{}
	},
}
