package asn

func f(x int32) int32 {
	var r int32
	r = (1 + x)
	return r
}
func g(x int32) (int32, int32) {
	return (1 + x), (1 - x)
}
func h(x int32) int32 {
	var r int32
	x, r = g(x)
	return r
}
