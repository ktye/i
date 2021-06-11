package main

import (
	. "github.com/ktye/wg/module"
)

func Fnd(x, y K) (r K) { // x?y
	xt, yt := tp(x), tp(y)
	if xt < 16 {
		trap(Type)
	}
	if xt == Lt {
		if yt == Lt {
			return Ecr(18, l2(x, y))
		} else {
			return fndl(x, y)
		}
	}
	if xt == yt+16 {
		r = Ki(fnd(x, y, yt))
	} else if xt == yt {
		yn := nn(y)
		r = mk(It, yn)
		rp := int32(r)
		for i := int32(0); i < yn; i++ {
			yi := ati(rx(y), i)
			SetI32(rp, fnd(x, yi, xt-16))
			dx(yi)
			rp += 4
		}
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
		trap(Type)
	}
	if r < 0 {
		return xn
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
func fndl(x, y K) (r K) {
	xn := nn(x)
	xp := int32(x)
	for i := int32(0); i < xn; i++ {
		if match(K(I64(xp)), y) != 0 {
			r = Ki(i)
			break
		}
		xp += 8
	}
	if r == 0 {
		r = Ki(xn)
	}
	dx(x)
	dx(y)
	return r
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
	e := yp + yn - xn
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
