package main

import (
	"github.com/ktye/wg/wasi_unstable"

	. "github.com/ktye/wg/module"
)

//func LK(x K) K         { return K(wasi_unstable.L64(uint64(x))) }
//func LI(x int32) int32 { return wasi_unstable.L32(x) }
func main() { // _start
	kinit()
	xx := Ku(23644)            // \\
	write(Ku(117851310093419)) // "ktye/k"
	for {
		write(Ku(8202)) // "\n "
		x := read()
		if match(x, xx) != 0 {
			wasi_unstable.Proc_exit(0)
		}
		x = val(x)
		if x != 0 {
			write(Kst(x))
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
