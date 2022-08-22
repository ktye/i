package main

import (
	. "github.com/ktye/wg/module"
)

func Fnd(x, y K) (r K) { // x?y
	xt, yt := tp(x), tp(y)
	if xt < 16 {
		if yt == Tt {
			return grp(x, y)
		} else {
			return deal(x, y)
		}
	}
	if xt > Lt {
		if xt == Tt {
			trap(Nyi) // t?..
		}
		r, x = spl2(x)
		return Atx(r, Fnd(x, y))
	} else if xt == yt {
		yn := nn(y)
		if SMALL == false {
			if xt < Ft && yn > 2 {
				if yn > 4 && xt == Ct || yn > 8 {
					r = fndXs(x, y, xt, yn)
					if r != 0 {
						return r
					}
				}
			}
		}
		r = mk(It, yn)
		rp := int32(r)
		if xt == Lt {
			yp := int32(y)
			for i := int32(0); i < yn; i++ {
				SetI32(rp, fndl(x, x0(int32(yp))))
				rp += 4
				yp += 8
			}
		} else {
			for i := int32(0); i < yn; i++ {
				yi := ati(rx(y), i)
				SetI32(rp, fnd(x, yi, xt-16))
				dx(yi)
				rp += 4
			}
		}
	} else if xt == yt+16 {
		r = Ki(fnd(x, y, yt))
	} else if xt == Lt {
		r = Ki(fndl(x, rx(y)))
	} else if yt == Lt {
		return Ecr(18, l2(x, y))
	} else {
		trap(Type)
	}
	dx(x)
	dx(y)
	return r
}
func fnd(x, y K, t T) (r int32) {
	xn := nn(x)
	if xn == 0 {
		return nai
	}
	xp, yp := int32(x), int32(y)
	xe := ep(x)
	switch t - 2 {
	case 0: // ct
		r = idxc(yp, xp, xe)
	case 1: // it
		r = idxi(yp, xp, xe)
	case 2: // st
		r = idxi(yp, xp, xe)
	case 3: // ft
		r = idxf(F64(yp), xp, xe)
	default: // zt
		r = idxz(F64(yp), F64(yp+8), xp, xe)
	}
	if r < 0 {
		return nai
	}
	return r
}
func idxc(x, p, e int32) (r int32) {
	r = inC(x, p, e)
	if r == 0 {
		return -1
	}
	return r - p
}
func idxi(x, p, e int32) (r int32) {
	r = inI(x, p, e)
	if r == 0 {
		return -1
	}
	return (r - p) >> 2
}
func idxf(x float64, p, e int32) (r int32) {
	r = inF(x, p, e)
	if r == 0 {
		return -1
	}
	return (r - p) >> 3
}
func idxz(re, im float64, p, e int32) (r int32) {
	r = inZ(re, im, p, e)
	if r == 0 {
		return -1
	}
	return (r - p) >> 4
}
func fndl(x, y K) (r int32) {
	xn := nn(x)
	xp := int32(x)
	dx(y)
	for r < xn {
		if match(K(I64(xp)), y) != 0 {
			return r
		}
		r++
		xp += 8
	}
	return nai
}
func idx(x, a, b int32) int32 {
	for i := a; i < b; i++ {
		if x == I8(i) {
			return i - a
		}
	}
	return -1
}

func Find(x, y K) (r K) { // find[pattern;string] returns all matches (It)
	xt, yt := tp(x), tp(y)
	if xt != yt || xt != Ct {
		trap(Type)
	}
	xn, yn := nn(x), nn(y)
	if xn == 0 || yn == 0 {
		dx(x)
		dx(y)
		return mk(It, 0)
	}
	r = mk(It, 0)
	xp, yp := int32(x), int32(y)
	y0 := yp
	e := yp + yn
	for yp < e { // todo rabin-karp / knuth-morris / boyes-moore..
		if findat(xp, yp, xn) != 0 {
			r = cat1(r, Ki(yp-y0))
			yp += xn
		} else {
			yp++
		}
		continue
	}
	dx(x)
	dx(y)
	return r
}
func findat(xp, yp, n int32) int32 {
	for i := int32(0); i < n; i++ {
		if I8(xp+i) != I8(yp+i) {
			return 0
		}
		continue
	}
	return 1
}
