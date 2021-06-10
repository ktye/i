package main

import (
	"github.com/ktye/wg/wasi_unstable"

	. "github.com/ktye/wg/module"
)

//func LK(x K) K         { return K(wasi_unstable.L64(uint64(x))) }
//func LI(x int32) int32 { return wasi_unstable.L32(x) }
func main() { // _start
	kinit()
	// dx(Asn(Ks(8), readfile(Ku(500237086328)))) // x.txt
	getargs()
	write(Ku(2932601077199979)) // "ktye/k\n"
	for {
		write(Ku(32))
		x := read()
		repl(x)
	}
}
func repl(x K) {
	n := nn(x)
	xp := int32(x)
	if n > 0 && I8(xp) == 92 { // \
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
	x = val(x)
	if x != 0 {
		dx(Out(x))
	}
}
func bench(x K) {
	i := fndc(x, 32)
	if i < 0 {
		trap(Parse)
	}
	n := maxi(1, int32(prs(it, ntake(i, rx(x)))))
	x = parse(ndrop(i, x))
	t := time()
	for n > 0 {
		dx(exec(rx(x)))
		n--
		continue
	}
	t = time() - t
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
func time() int64 {
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
	SetI32(512, int32(x))
	SetI32(516, nn(x))
	if wasi_unstable.Fd_write(1, 512, 1, 512) != 0 {
		trap(Io)
	}
	dx(x)
}
func getargs() int32 {
	wasi_unstable.Args_sizes_get(512, 516)
	n := I32(516)
	r := mk(Ct, n)
	SetI32(512, 516)
	wasi_unstable.Args_get(512, int32(r))
	dx(Asn(sc(Ku(24878)), split(Kc(0), ndrop(-1, r)))) //.a
	return n
}
func readfile(x K) (r K) { // single-read max 504
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
	return r
}
func iwrite(x int32) { write(cat1(Kst(Ki(x)), Kc(10))) }
func help() {
	trap(Nyi)
}
