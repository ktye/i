package main

import (
	. "github.com/ktye/wg/module"
)

func ech(x K) K { return l2t(x, 0, df) } // '
func ecp(x K) K { return l2t(x, 1, df) } // ':
func rdc(x K) K { return l2t(x, 2, df) } // /
func ecr(x K) K { return l2t(x, 3, df) } // /:
func scn(x K) K { return l2t(x, 4, df) } // \
func ecl(x K) K { return l2t(x, 5, df) } // \:

func Ech(f, x K) (r K) {
	if nn(x) == 1 {
		x = Fst(x)
	} else {
		return ecn(f, x)
	}
	if tp(x) < 16 {
		x = enl(x)
	}
	xt := tp(x)
	if xt > Lt {
		panic(Nyi)
	}
	xn := nn(x)
	r = mk(Lt, xn)
	rp := int32(r)
	for i := int32(0); i < xn; i++ {
		SetI64(rp, int64(Atx(rx(f), ati(rx(x), i))))
		rp += 8
	}
	dx(f)
	dx(x)
	return uf(r)
}
func ecn(f, x K) K { trap(Nyi); return x }
func Ecp(f, x K) K { trap(Nyi); return x }
func Rdc(f, x K) (r K) { // x f/y   (x=0):f/y
	var y K
	if xn := nn(x); xn == 1 {
		y = Fst(x)
		x = 0
	} else if xn == 2 {
		x, y = spl2(x)
	} else {
		trap(Rank)
	}
	yt := tp(y)
	if yt < 16 {
		if x == 0 {
			dx(f)
			return x
		} else {
			return cal(f, l2(x, y))
		}
	}
	yn := nn(y)
	if yn == 0 {
		if x == 0 {
			return zero(yt - 16)
		} else {
			dx(f)
			dx(y)
			return x
		}
	}
	if yt > Lt {
		trap(Nyi)
	}
	xt := tp(x)
	if tp(f) == 0 {
		fp := int32(f)
		if fp == 1 {
			dx(x)
			return lst(y)
		}
		if fp == 13 {
			return cats(x, y)
		}
		if yt != Lt && fp < 9 && (xt == yt-16 || xt == 0) {
			return rdx(fp, x, y, yn) // +-*% &|
		}
	}
	i := int32(0)
	if x == 0 {
		x, i = ati(rx(y), 0), 1
	}
	for i < yn {
		x = cal(rx(f), l2(x, ati(rx(y), i)))
		i++
	}
	dx(y)
	dx(f)
	return x
}

func rdx(fp int32, x, y K, n int32) (r K) { // (+-*% &|)/
	yt := tp(y)
	s := sz(yt)
	yp := int32(y)
	i := int32(0)
	if x == 0 {
		x, i = Fst(rx(y)), 1
	}
	xp := int32(x)
	if s == 1 {
		fp = 214 + 3*fp
		for i < n {
			xp = Func[fp].(f2i)(xp, I8(yp+i))
			i++
		}
		if yt == Bt {
			r = Ki(xp)
		} else {
			r = Kc(xp)
		}
	} else if s == 4 {
		fp = 214 + 3*fp
		for i < n {
			xp = Func[fp].(f2i)(xp, I32(yp+4*i))
			i++
		}
		r = Ki(xp)
	} else {
		if yt == Zt {
			trap(Nyi)
		}
		xf := F64(xp)
		fp = 214 + 3*fp
		for i < n {
			xf = Func[fp].(f2f)(xf, F64(yp+8*i))
			i++
		}
		r = Kf(xf)
	}
	dx(x)
	dx(y)
	return r
}
func Ecr(f, x K) K { trap(Nyi); return x }
func Scn(f, x K) K { trap(Nyi); return x }
func Ecl(f, x K) K { trap(Nyi); return x }

func uf(x K) (r K) {
	xn := nn(x)
	xp := int32(x)
	var rt T
	for i := int32(0); i < xn; i++ {
		t := tp(K(I64(xp)))
		if i == 0 {
			rt = t
		} else if t != rt {
			return x
		}
		xp += 8
	}
	if rt == 0 || rt > zt {
		return x
	}
	if rt > st {
		trap(Nyi)
	}
	rt += 16
	r = mk(rt, xn)
	s := sz(rt)
	rp := int32(r)
	xp = int32(x)
	if s == 1 {
		for i := int32(0); i < xn; i++ {
			SetI8(rp, I32(xp))
			xp += 8
			rp++
		}
	} else if s == 4 {
		for i := int32(0); i < xn; i++ {
			SetI32(rp, I32(xp))
			xp += 8
			rp += 4
		}
	} else {
		trap(Nyi)
	}
	dx(x)
	return r
}
