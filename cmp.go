package k

import . "github.com/ktye/wg/module"

func eqc(x, y K) int32 {
	n := nn(x)
	if n != nn(y) {
		return 0
	}
	xp := int32(x)
	yp := int32(y)
	for i := int32(0); i < n; i++ {
		if I8(xp) != I8(yp) {
			return 0
		}
		xp++
		yp++
	}
	return 1
}
