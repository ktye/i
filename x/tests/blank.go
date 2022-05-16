package blank

func f() (int32, int32) {
	return 1, 2
}
func g() int32 {
	var x int32
	x, _ = f()
	return x
}
