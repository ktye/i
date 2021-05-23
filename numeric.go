package k

import (
	. "github.com/ktye/wg/module"
)

func Neg(x K) K { return Ki(-int32(x)) }

func Sqr(x K) K { trap(Nyi); return x }

func Add(x, y K) K { return Ki(int32(x) + int32(y)) }
func Sub(x, y K) K { return Ki(int32(x) - int32(y)) }
func Mul(x, y K) K { return Ki(int32(x) * int32(y)) }
func Div(x, y K) K { return Ki(int32(x) / int32(y)) }

func Min(x, y K) K { trap(Nyi); return x }
func Max(x, y K) K { trap(Nyi); return x }

func uptypes(x, y K) (K, K) {
	xt, yt := tp(x)&15, tp(y)&15
	if xt < yt {
		x, xt = uptype(x, xt), yt
	}
	if yt < xt {
		y = uptype(y, yt)
	}
	return x, y
}
func uptype(x K, dst T) (r K) {
	xt := tp(x)
	xp := int32(x)
	if xt < 16 {
		if dst == ct {
			return Kc(xp)
		} else if dst == it {
			return Ki(xp)
		} else if dst == ft {
			return Kf(float64(xp))
		} else if dst == zt {
			if xt == ft {
				return Kz(x, 0)
			} else {
				return Kz(Kf(float64(xp)), 0)
			}
		} else {
			trap(Type)
		}
	}
	if xt < It && dst == Ft {
		x, xt = uptype(x, it), It
	}
	if xt < Ft && dst == Zt {
		x, xt = uptype(x, ft), Ft
	}
	dst += 16
	xn := nn(x)
	r = mk(dst, xn)
	rp := int32(r)
	if dst == ct {
		Memorycopy(rp, xp, xn)
	} else if dst == it {
		for i := int32(0); i < xn; i++ {
			SetI32(rp, I8(xp))
			xp++
			rp += 4
		}
	} else if dst == ft {
		for i := int32(0); i < xn; i++ {
			SetF64(rp, float64(I32(xp)))
			xp += 4
			rp += 8
		}
	} else if dst == zt {
		return KZ(rx(x), 0)
	} else {
		trap(Type)
	}
	dx(x)
	return r
}
