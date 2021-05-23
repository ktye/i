package k

import (
	. "github.com/ktye/wg/module"
)

func ech(x K) K { return l2t(x, 0, df) } // '
func ecp(x K) K { return l2t(x, 1, df) } // ':
func rdc(x K) K { return l2t(x, 2, df) } // /
func ecr(x K) K { return l2t(x, 3, df) } // /:
func scn(x K) K { return l2t(x, 4, df) } // \
func ecl(x K) K { return l2t(x, 5, df) } // \:

func Ech(f, x K) K { trap(Nyi); return x }
func Ecp(f, x K) K { trap(Nyi); return x }
func Rdc(f, x K) (r K) { // f/x
	if tp(x) != It {
		trap(Nyi)
	}
	xp := int32(x)
	xn := nn(x)
	r = Ki(I32(xp))
	g := 64 + int32(f)
	for i := int32(1); i < xn; i++ {
		xp += 4
		r = Func[g].(f2)(r, Ki(I32(xp)))
	}
	dx(f)
	dx(x)
	return r
}
func Ecr(f, x K) K { trap(Nyi); return x }
func Scn(f, x K) K { trap(Nyi); return x }
func Ecl(f, x K) K { trap(Nyi); return x }
