package cnd2

func f(x int32) int32 {
	if x > 3 {
		return x
	} else {
		x = (2 * x)
		return -x
	}
}
