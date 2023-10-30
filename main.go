//go:build !small

package main

import (
	. "github.com/ktye/wg/module"
)

func main() { // _start
	kinit()
	doargs()
	write(Ku(2932601077199979)) // "ktye/k\n"
	store()
	for {
		write(Ku(32))
		x := read()
		try(x)
	}
}

func store() {
	g := (1 << (I32(128) - 16)) - Memorysize2()
	if g > 0 {
		Memorygrow2(g)
	}
	Memorycopy2(0, 0, int32(1)<<I32(128))
}
func catch() {
	Memorycopy3(0, 0, int32(65536)*Memorysize2())
}
func try(x K) {
	defer Catch(catch)
	repl(x)
	store()
}

func doargs() {
	a := ndrop(1, getargv())
	an := nn(a)
	ee := Ku(25901) // -e
	for i := int32(0); i < an; i++ {
		x := x0(a)
		if match(x, ee) != 0 { // -e (exit)
			if i < an-1 {
				dx(x)
				x = x1(a)
				dx(ee)
				repl(x)
			}
			Exit(0)
		}
		dofile(x, readfile(rx(x)))
		a += 8
	}
	dx(ee)
}
func dofile(x K, c K) {
	kk := Ku(27438) // .k
	tt := Ku(29742) // .t
	xe := ntake(-2, rx(x))
	if match(xe, kk) != 0 { // file.k (execute)
		dx(val(c))
	} else if match(xe, tt) != 0 { // file.t (test)
		test(c)
	} else { // file (assign file:bytes..)
		dx(Asn(sc(rx(x)), c))
	}
	dxy(xe, x)
	dxy(tt, kk)
}
func getargv() K {
	n := Args()
	r := mk(Lt, n)
	rp := int32(r)
	for i := int32(0); i < n; i++ {
		s := mk(Ct, Arg(i, 0))
		Arg(i, int32(s))
		SetI64(rp, int64(s))
		rp += 8
	}
	return r
}
