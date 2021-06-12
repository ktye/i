package main

import . "github.com/ktye/wg/module"

func test(x K, clr int32) { // clr=1(reset after each line)
	if tp(x) != Ct {
		trap(Type)
	}
	xn := nn(x)
	y := x
	l := split(Kc(10), rx(x))
	n := nn(l)
	dx(l)
	for i := int32(0); i < n; i++ {
		testi(rx(x), i)
		if clr != 0 {
			dx(x)
			reset()
			x = mk(Ct, xn)
			Memorycopy(int32(x), int32(y), xn)
			y = x
		}
	}
	dx(x)
}
func testi(l K, i int32) {
	var x, y K
	x = split(Ku(12064), ati(split(Kc(10), l), i))
	if nn(x) != 2 {
		trap(Length)
	}
	x, y = spl2(x)
	dx(Out(ucat(ucat(rx(x), Ku(12064)), rx(y))))
	x = Kst(val(x))
	if match(x, y) == 0 {
		x = Out(x)
		trap(Err)
	}
	dx(x)
	dx(y)
}
