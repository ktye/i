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
func Ecp(f, x K) K { trap(Nyi); return x }
func Rdc(f, x K) (r K) { // f/x
	xt := tp(x)
	if xt < 16 {
		dx(f)
		return x
	}
	xn := nn(x)
	if xn == 0 {
		return zero(xt - 16)
	}
	if xt > Lt {
		trap(Nyi)
	}
	if tp(f) == 0 {
		fp := int32(f)
		if fp == 1 {
			return lst(x)
		}
		if fp == 13 {
			return cats(x)
		}
		if xt != Lt && fp < 9 {
			return rdx(fp, x, xn) // +-*% &|
		}
	}
	r = ati(rx(x), 0)
	for i := int32(1); i < xn; i++ {
		r = cal(f, l2(r, ati(rx(x), i)))
	}
	dx(x)
	dx(f)
	return r
}
func rdx(fp int32, x K, n int32) (r K) { // (+-*% &|)/
	xt := tp(x)
	s := sz(xt)
	xp := int32(x)
	xn := nn(x)
	if s == 1 {
		xi := I8(xp)
		fp = 214 + 3*fp
		for i := int32(1); i < xn; i++ {
			xp++
			xi = Func[fp].(f2i)(xi, I8(xp))
		}
		if xt == bt {
			r = Ki(xi)
		} else {
			r = Kc(xi)
		}
	} else if s == 4 {
		xi := I32(xp)
		fp = 214 + 3*fp
		for i := int32(1); i < xn; i++ {
			xp += 4
			xi = Func[fp].(f2i)(xi, I32(xp))
		}
		r = Ki(xi)
	} else {
		if xt == Zt {
			trap(Nyi)
		}
		xi := F64(xp)
		fp = 214 + 3*fp
		for i := int32(1); i < xn; i++ {
			xp += 8
			xi = Func[fp].(f2f)(xi, F64(xp))
		}
		r = Kf(xi)
	}
	dx(x)
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
