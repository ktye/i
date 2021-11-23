package main

import (
	"github.com/ktye/wg/wasi_unstable"

	. "github.com/ktye/wg/module"
)

//func LK(x K) K         { return K(wasi_unstable.L64(uint64(x))) }
//func LI(x int32) int32 { return wasi_unstable.L32(x) }
func main() { // _start
	kinit()
	doargs()
	write(Ku(2932601077199979)) // "ktye/k\n"
	for {
		write(Ku(32))
		x := read()
		try(x)
	}
}
func catch() {
	Memorycopy3(0, 0, int32(65536)*Memorysize2())
}
func try(x K) {
	defer Catch(catch)
	repl(x)
	g := (1 << (I32(128) - 16)) - Memorysize2()
	if g > 0 {
		Memorygrow2(g)
	}
	Memorycopy2(0, 0, 1<<I32(128))
}
func repl(x K) {
	n := nn(x)
	xp := int32(x)
	s := int32(0)
	if n > 0 {
		s = I8(xp)
		if I8(xp) == 92 { // \
			if n == 1 {
				help()
			} else {
				c := I8(1 + xp)
				if I8(1+xp) == '\\' {
					wasi_unstable.Proc_exit(0)
				} else if c == 't' {
					bench(ndrop(2, x))
				} else if c == 'm' {
					dx(x)
					dx(Out(Ki(I32(128))))
				} else if c == 'c' {
					dx(x)
					reset()
				}
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
	a := ndrop(1, getargs())
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
			wasi_unstable.Proc_exit(0)
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
func bench(x K) {
	i := fnd(x, Kc(32), ct)
	if i == nai {
		trap(Parse)
	}
	n := maxi(1, int32(prs(it, ntake(i, rx(x)))))
	x = parse(tok(ndrop(i, x)))
	t := time_()
	for n > 0 {
		dx(exec(rx(x)))
		n--
		continue
	}
	t = time_() - t
	dx(Out(Kf(1e-6 * float64(t))))
}

func Out(x K) K {
	write(cat1(Kst(rx(x)), Kc(10)))
	return x
}
func Otu(x, y K) K {
	write(cat1(Kst(x), Kc(':')))
	return Out(y)
}
func time_() int64 {
	wasi_unstable.Clock_time_get(1, wasi_unstable.Timestamp(0), 512)
	return I64(512)
}
func read() (r K) {
	r = mk(Ct, 504)
	rp := int32(r)
	SetI32(512, rp)
	SetI32(516, 504)
	if wasi_unstable.Fd_read(0, 512, 1, 512) != 0 {
		trap(Io)
	}
	return ntake(maxi(0, I32(512)-1), r)
}
func write(x K) {
	if writepn(int32(x), nn(x)) != 0 {
		trap(Io)
	}
	dx(x)
}
func writepn(p, n int32) int32 {
	SetI32(512, p)
	SetI32(516, n)
	return wasi_unstable.Fd_write(1, 512, 1, 512)
}
func getargs() K {
	wasi_unstable.Args_sizes_get(512, 516)
	n := I32(516)
	a := mk(It, I32(512))
	r := mk(Ct, n)
	//SetI32(512, int32(a))
	wasi_unstable.Args_get(int32(a), int32(r))
	dx(a)
	return split(Kc(0), ndrop(-1, r))
}
func readfile(x K) (r K) { // x C
	// fd=3 is root directory, e.g. wavm run --mount-root . k.wat
	if wasi_unstable.Path_open(3, 0, int32(x), nn(x), 0, 31, 31, 0, 512) != 0 {
		trap(Io)
	}
	fd := I32(512)

	if wasi_unstable.Fd_seek(fd, 0, 2, 512) != 0 {
		trap(Io)
	}
	n := I32(512)
	if wasi_unstable.Fd_seek(fd, 0, 0, 512) != 0 {
		trap(Io)
	}

	r = mk(Ct, n)
	rp := int32(r)
	SetI32(512, rp)
	SetI32(516, n)
	if wasi_unstable.Fd_read(fd, 512, 1, 512) != 0 {
		trap(Io)
	}
	if I32(512) != n {
		trap(Io)
	}
	wasi_unstable.Fd_close(fd)
	dx(x)
	return r
}
func writefile(x, y K) K { // x, y C
	if wasi_unstable.Path_open(3, 0, int32(x), nn(x), 9, 2047, 2047, 0, 512) != 0 {
		trap(Io)
	}
	fd := I32(512)

	yp := int32(y)
	SetI32(512, yp)
	SetI32(516, nn(y))
	if wasi_unstable.Fd_write(fd, 512, 1, 512) != 0 {
		trap(Io)
	}
	wasi_unstable.Fd_close(fd)
	dx(x)
	return y
}
func iwrite(x int32) { write(cat1(Kst(Ki(x)), Kc(10))) }
func help() {
	trap(Nyi)
}
