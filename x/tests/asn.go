package asn

func f(x int32) int32 {
	var r int32
	r = (1 + x)
	return r
}
func h(x int32) int32 {
	var r int32
	f(x)
	r = f(x)
	return r
}
