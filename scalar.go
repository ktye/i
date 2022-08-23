package main

import (
	. "github.com/ktye/wg/module"
)

func conform(x, y K) (int32, T) { // 0:atom-atom 1:atom-vector, 2:vector-vector, 3:vector-atom
	xt, yt := tp(x), tp(y)
	if xt < 16 {
		if yt < 16 {
			return 0, xt
		} else {
			return 1, xt
		}
	}
	xn := nn(x)
	if yt < 16 {
		return 3, yt
	}
	if nn(y) != xn {
		trap(Length)
	}
	return 2, xt - 16
}
func dctypes(x, y K) (K, T, K, K) {
	xt, yt := tp(x), tp(y)
	t := T(maxi(int32(xt), int32(yt)))
	if xt < Dt && yt < Dt {
		return 0, t, x, y
	}
	var k K
	if xt > Lt {
		k, x = spl2(x)
		if yt > Lt {
			var yk K
			yk, y = spl2(y)
			if match(k, yk) == 0 {
				trap(Value)
			}
			dx(yk)
		}
	} else if yt > Lt {
		k, y = spl2(y)
	}
	return k, t, x, y
}
func uptypes(x, y K, b2i int32) (K, K) {
	xt, yt := tp(x)&15, tp(y)&15
	rt := T(maxi(int32(xt), int32(yt)))
	if rt == 0 {
		rt = it
	}
	if xt < rt {
		x = uptype(x, rt)
	}
	if yt < rt {
		y = uptype(y, rt)
	}
	return x, y
}
func uptype(x K, dst T) (r K) {
	xt := tp(x)
	xp := int32(x)
	if xt&15 == dst {
		return x
	}
	if xt < 16 {
		if dst == ct {
			return Kc(xp)
		} else if dst == it {
			return Ki(xp)
		} else if dst == ft {
			return Kf(float64(xp))
		} else if dst == zt {
			var f float64
			if xt == ft {
				f = F64(xp)
				dx(x)
			} else {
				f = float64(xp)
			}
			return Kz(f, 0)
		} else {
			return trap(Type)
		}
	}
	if xt < It && dst == ft {
		x, xt = uptype(x, it), It
	}
	if xt < Ft && dst == zt {
		x, xt = uptype(x, ft), Ft
	}
	xn := nn(x)
	xp = int32(x)
	r = mk(dst+16, xn)
	rp := int32(r)
	if dst == it {
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
		for i := int32(0); i < xn; i++ {
			SetF64(rp, F64(xp))
			SetF64(rp+8, 0.0)
			xp += 8
			rp += 16
		}
	} else {
		trap(Type)
	}
	dx(x)
	return r
}
