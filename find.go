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
	if xt == Dt {
		r, x = spl2(x)
		return Atx(r, Fnd(x, y))
	} else if xt == yt {
		yn := nn(y)
		if xt < Ft && yn > 2 {
			if xt == Bt {
				return fndBs(x, y)
			}
			if yn > 4 && xt == Ct || yn > 8 {
				r = fndXs(x, y, xt, yn)
				if r != 0 {
					return r
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
		return 0
	}
	xp, yp := int32(x), int32(y)
	xe := ep(x)
	ve := xe &^ 15
	switch t - 1 {
	case 0: // bt
		r = idxc(yp, xp, ve, xe)
	case 1: // ct
		r = idxc(yp, xp, ve, xe)
	case 2: // it
		r = idxi(yp, xp, ve, xe)
	case 3: // st
		r = idxi(yp, xp, ve, xe)
	case 4: // ft
		r = idxf(F64(xp), yp, xe)
	case 5: // zt
		r = idxz(F64(xp), F64(xp+8), yp, xe)
	default:
		r = int32(trap(Type))
	}
	if r < 0 {
		return xn
	}
	return r
}
func fndBs(x, y K) K {
	a := fnd(x, Kb(0), 1)
	b := fnd(x, Kb(1), 1) - a
	dx(x)
	return Add(Ki(a), Mul(Ki(b), y))
}
func fndXs(x, y K, t T, yn int32) (r K) {
	xn := nn(x)
	a := int32(min(0, int32(x), t, xn))
	b := 1 + (int32(max(0, int32(x), t, xn))-a)>>(3*I32B(t == St))
	if b > 256 && b > yn {
		return 0
	}
	if t == St {
		x, y = Div(Flr(x), Ki(8)), Div(Flr(y), Ki(8))
		a >>= 3
	}
	r = ntake(b, Ki(xn))
	rp := int32(r) - 4*a
	x0 := int32(x)
	xp := ep(x)
	if t == Ct {
		for xp > x0 {
			xp--
			SetI8(rp+4*I8(xp), xp-x0)
		}
	} else {
		for xp > x0 {
			xp -= 4
			SetI32(rp+4*I32(xp), (xp-x0)>>2)
		}
	}
	dx(x)
	r = Atx(r, Add(Ki(-a), y))
	rp = int32(r)
	xp = rp + 4*yn
	for rp < xp {
		if I32(rp) == nai {
			SetI32(rp, xn)
		}
		rp += 4
		continue
	}
	return r
}
func idxc(x, p, ve, e int32) (r int32) {
	r = inC(x, p, ve, e)
	if r == 0 {
		return -1
	}
	e = r + 16
	for i := r; i < e; i++ {
		if I8(i) == x {
			return i - p
		}
		continue
	}
	trap(Err)
	return r //not reached
}
func idxi(x, p, ve, e int32) (r int32) {
	r = inI(x, p, ve, e)
	if r == 0 {
		return -1
	}
	e = r + 16
	for i := r; i < e; i += 4 {
		if I32(i) == x {
			return (i - p) >> 2
		}
		continue
	}
	trap(Err)
	return r //not reached
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
	return xn
}

func index(x, a, b int32) int32 {
	for i := a; i < b; i++ {
		if x == I8(i) {
			return i - a
		}
	}
	return -1
}
func fndc(x K, c int32) int32 {
	xp := int32(x)
	xn := nn(x)
	for i := int32(0); i < xn; i++ {
		if I8(xp) == c {
			return i
		}
		xp++
	}
	return -1
}

func Find(x, y K) (r K) { // find[pattern;string] returns all matches (It)
	xt, yt := tp(x), tp(y)
	if xt != yt && xt != Ct {
		trap(Type)
	}
	xn, yn := nn(x), nn(y)
	if xn == 0 {
		trap(Length)
	}
	if yn == 0 {
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
