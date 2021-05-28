package k

import (
	"github.com/ktye/wg/wasi_unstable"

	. "github.com/ktye/wg/module"
)

func Out(x K) K {
	write(cat1(Kst(rx(x)), Kc(10)))
	return x
}
func Otu(x, y K) K {
	write(cat1(Kst(x), Kc(':')))
	return Out(y)
}
func write(x K) {
	SetI32(512, int32(x))
	SetI32(516, nn(x))
	wasi_unstable.Fd_write(1, 512, 1, 512)
	dx(x)
}
