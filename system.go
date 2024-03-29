package main

import (
	. "github.com/ktye/wg/module"
)

func repl(x K) {
	c := I8(int32(x))
	x = val(x)
	if x != 0 {
		if c == 32 {
			dx(Out(x))
		} else {
			write(cat1(join(Kc(10), Lst(x)), Kc(10)))
		}
	}
}

func Out(x K) K {
	write(cat1(Kst(rx(x)), Kc(10)))
	return x
}
func Otu(x, y K) K {
	write(cat1(Kst(x), Kc(':')))
	return Out(y)
}
func write(x K) {
	Write(0, 0, int32(x), nn(x))
	dx(x)
}
func readfile(x K) K { // x C
	var r K
	if nn(x) == 0 {
		dx(x)
		r = mk(Ct, 496)
		r = ntake(ReadIn(int32(r), 496), r)
		return r
	}
	n := Read(int32(x), nn(x), 0)
	if n < 0 {
		dx(x)
		return mk(Ct, 0)
	}
	r = mk(Ct, n)
	Read(int32(x), nn(x), int32(r))
	dx(x)
	return r
}
func writefile(x, y K) K { // x, y C
	r := Write(int32(x), nn(x), int32(y), nn(y))
	if r != 0 {
		trap() //io
	}
	dx(x)
	return y
}
