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
	Memorycopy2(0, 0, 1<<I32(128))
}
func catch() {
	Memorycopy3(0, 0, int32(65536)*Memorysize2())
}
func try(x K) {
	defer Catch(catch)
	repl(x)
	store()
}
func repl(x K) {
	n := nn(x)
	xp := int32(x)
	s := int32(0)
	if n > 0 {
		s = I8(xp)
		if I8(xp) == 92 && n > 1 { // \
			c := I8(1 + xp)
			if I8(1+xp) == '\\' {
				Exit(0)
			} else if c == 'm' {
				dx(x)
				dx(Out(Ki(I32(128))))
			}
			return
		}
	}
	x = val(x)
	if x != 0 {
		if s == 32 {
			dx(Out(x))
		} else {
			write(cat1(join(Kc(10), Lst(x)), Kc(10)))
		}
	}
}
func doargs() {
	a := ndrop(1, getargv())
	an := nn(a)
	ap := int32(a)
	ee := Ku(25901) // -e
	for i := int32(0); i < an; i++ {
		x := x0(ap)
		if match(x, ee) != 0 { // -e (exit)
			if i < an-1 {
				dx(x)
				x = x1(ap)
				dx(ee)
				dx(a)
				repl(x)
			}
			Exit(0)
		}
		dofile(x, readfile(rx(x)))
		ap += 8
	}
	dx(ee)
	dx(a)
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
	dx(xe)
	dx(x)
	dx(tt)
	dx(kk)
}

func Out(x K) K {
	write(cat1(Kst(rx(x)), Kc(10)))
	return x
}
func Otu(x, y K) K {
	write(cat1(Kst(x), Kc(':')))
	return Out(y)
}
func read() (r K) {
	r = mk(Ct, 504)
	return ntake(ReadIn(int32(r), 504), r)
}
func write(x K) {
	Write(0, 0, int32(x), nn(x))
	dx(x)
}
func getargv() (r K) {
	n := Args()
	r = mk(Lt, n)
	rp := int32(r)
	for i := int32(0); i < n; i++ {
		s := mk(Ct, Arg(i, 0))
		Arg(i, int32(s))
		SetI64(rp, int64(s))
		rp += 8
	}
	return r
}
func readfile(x K) (r K) { // x C
	if nn(x) == 0 {
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
		trap(Io)
	}
	dx(x)
	return y
}
