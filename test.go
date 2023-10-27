package main

func test(x K) {
	if tp(x) != Ct {
		trap() //type
	}
	l := ndrop(-1, split(Kc(10), rx(x)))
	n := nn(l)
	dx(l)
	for i := int32(0); i < n; i++ {
		testi(rx(x), i)
	}
	dx(x)
}
func testi(l K, i int32) {
	x := split(Ku(12064), ati(split(Kc(10), l), i))
	if nn(x) != 2 {
		trap() //length
	}
	y := x1(x)
	x = r0(x)
	dx(Out(ucat(ucat(rx(x), Ku(12064)), rx(y))))
	x = Kst(val(x))
	if match(x, y) == 0 {
		x = Out(x)
		trap() //test fails
	}
	dx(x)
	dx(y)
}
