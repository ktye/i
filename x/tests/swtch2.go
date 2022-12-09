package swtch2

func f(x int32) int32 {
	return x
}
func g(x int32) int32 {
	var r int32
	switch x {
	case 0:
		r = f(x)
	default:
		if x > 5 {
			r = (x - 3)
		} else {
			r = (x - 2)
		}
		r = r
	}
	return r
}
