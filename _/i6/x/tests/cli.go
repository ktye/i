package cli

import . "github.com/ktye/wg/module"

func f(x int32, y int32) int32 {
	return Func[x].(func(int32) int32)(y)
}
