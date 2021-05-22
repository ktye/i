package k

import . "github.com/ktye/wg/module"

func index(x, a, b int32) int32 {
	for i := a; i < b; i++ {
		if x == I8(i) {
			return i - a
		}
	}
	return -1
}
func fndc(x K, c int32) int32 {
	xp := int32(x)
	xn := nn(x)
	for i := int32(0); i < xn; i++ {
		if I8(xp) == c {
			return i
		}
		xp++
	}
	return -1
}
