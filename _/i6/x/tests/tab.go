package tab

import . "github.com/ktye/wg/module"

func init() {

	Functions(0, f)
	Functions(1, g)
	Functions(2, h)

}

func f() {
}
func g(x int32) int32 {
	return x
}
func h(x int32, y int32) int32 {
	return (x + y)
}
