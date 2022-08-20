package main

import (
	. "github.com/ktye/wg/module"
)

type f1_ = func(int32, int32, int32)
type f1i = func(int32) int32
type f1f = func(float64) float64
type f1z = func(float64, float64) (float64, float64)
type f2b = func(uint64, uint64) uint64
type f2i = func(int32, int32) int32

//type c2cC = func(I8x16, I8x16, int32, int32, int32)
//type f2cC = func(I8x16, int32, int32, int32)
//type f2iI = func(I32x4, int32, int32, int32)
//type f2vF = func(F64x2, int32, int32, int32)
type ff3i = func(float64, int32, int32, int32)
type fF3i = func(float64, float64, int32, int32, int32)

//type c2Cc = func(int32, I8x16, I8x16, int32, int32)
//type f2Cc = func(int32, I8x16, int32, int32)
type f4i = func(int32, int32, int32, int32)
type f2Ff = func(int32, float64, int32, int32)
type f2Zz = func(int32, float64, float64, int32, int32)

//type f2v = func(int32, int32, int32, int32)

//type f2vc = func(I8x16, int32, int32, int32, int32)
type f2f = func(float64, float64) float64
type f2z = func(float64, float64, float64, float64) (float64, float64)
type f2c = func(float64, float64) int32
type f2d = func(float64, float64, float64, float64) int32

func use2(x, y K) K {
	if I32(int32(y)-4) == 1 {
		return rx(y)
	}
	return use1(x)
}
func use1(x K) K {
	if I32(int32(x)-4) == 1 {
		return rx(x)
	}
	return mk(tp(x), nn(x))
}
func use(x K) (r K) {
	xt := tp(x)
	if xt < 16 || xt > Lt {
		trap(Type)
	}
	if I32(int32(x)-4) == 1 {
		return x
	}
	nx := nn(x)
	r = mk(xt, nx)
	Memorycopy(int32(r), int32(x), sz(xt)*nx)
	if xt == Lt {
		rl(r)
	}
	dx(x)
	return r
}

func nm(f int32, x K) (r K) { //monadic
	xt := tp(x)
	if xt > Lt {
		r, x = spl2(x)
		return key(r, nm(f, x), xt)
	}
	xp := int32(x)
	if xt == Lt {
		n := nn(x)
		r = mk(Lt, n)
		rp := int32(r)
		for i := int32(0); i < n; i++ {
			SetI64(rp, int64(nm(f, x0(xp))))
			xp += 8
			rp += 8
		}
		dx(x)
		return uf(r)
	}
	if xt < 16 {
		switch xt - 2 {
		case 0:
			return Kc(Func[f].(f1i)(xp))
		case 1:
			return Ki(Func[f].(f1i)(xp))
		case 2:
			return trap(Type)
		case 3:
			r = Kf(Func[1+f].(f1f)(F64(xp)))
			dx(x)
			return r
		case 4:
			r = Kz(Func[2+f].(f1z)(F64(xp), F64(xp+8)))
			dx(x)
			return r
		default:
			return trap(Type)
		}
	}
	r = use1(x)
	rp := int32(r)
	e := ep(r)
	if e == rp {
		dx(x)
		return r
	}
	switch xt - 18 {
	case 0:
		Func[3+f].(f1_)(xp, rp, e)
	case 1:
		Func[4+f].(f1_)(xp, rp, e)
	case 2:
		trap(Type)
	case 3:
		Func[5+f].(f1_)(xp, rp, e)
	default:
		Func[6+f].(f1_)(xp, rp, e)
	}
	dx(x)
	return r
}

func Neg(x K) K                            { return nm(220, x) }
func negi(x int32) int32                   { return -x }
func negf(x float64) float64               { return -x }
func negz(x, y float64) (float64, float64) { return -x, -y }
func negC(xp, rp, e int32) {
	for rp < e {
		SetI8(rp, -I8(xp))
		xp++
		rp++
		continue
	}
}
func negI(xp, rp, e int32) {
	for rp < e {
		SetI32(rp, -I32(xp))
		xp += 4
		rp += 4
		continue
	}
}
func negF(xp, rp, e int32) {
	for rp < e {
		SetF64(rp, -F64(xp))
		xp += 8
		rp += 8
		continue
	}
}
func negZ(xp, rp, e int32) { negF(xp, rp, e) }

func Abs(x K) K {
	xt := tp(x)
	if xt > Zt {
		return Ech(32, l1(x))
	}
	if xt == zt {
		xp := int32(x)
		dx(x)
		return Kf(hypot(F64(xp), F64(xp+8)))
	} else if xt == Zt {
		return absZ(x)
	}
	return nm(227, x)
}
func absi(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}
func absf(x float64) float64 { return F64abs(x) }
func absC(xp, rp, e int32) {
	for rp < e {
		SetI8(rp, absi(I8(xp)))
		xp++
		rp++
		continue
	}
}
func absI(xp, rp, e int32) {
	for rp < e {
		SetI32(rp, absi(I32(xp)))
		xp += 4
		rp += 4
		continue
	}
}
func absF(xp, rp, e int32) {
	for rp < e {
		SetF64(rp, F64abs(F64(xp)))
		xp += 8
		rp += 8
		continue
	}
}
func absZ(x K) (r K) {
	n := nn(x)
	r = mk(Ft, n)
	rp := int32(r)
	xp := int32(x)
	for i := int32(0); i < n; i++ {
		SetF64(rp, hypot(F64(xp), F64(xp+8)))
		xp += 16
		rp += 8
		continue
	}
	dx(x)
	return r
}
func Hyp(x, y K) K { // e.g.  norm:0. abs/x
	xt := tp(x)
	yt := tp(y)
	if xt > Zt || yt > Zt {
		return Ech(32, l2(x, y))
	}
	if xt == zt {
		x, xt = Abs(x), ft
	}
	if xt == ft {
		xp := int32(x)
		yp := int32(y)
		dx(x)
		dx(y)
		if yt == ft {
			return Kf(hypot(F64(xp), F64(yp)))
		} else if yt == zt {
			return Kf(hypot(F64(xp), hypot(F64(yp), F64(yp+8))))
		}
	}
	return trap(Nyi)
}

func Sqr(x K) K {
	if tp(x)&15 != ft {
		x = Add(Kf(0), x)
	}
	return nm(300, x)
}
func sqrf(x float64) float64 { return F64sqrt(x) }
func sqrF(xp, rp, e int32) {
	for rp < e {
		SetF64(rp, F64sqrt(F64(xp)))
		xp += 8
		rp += 8
		continue
	}
}

func Img(x K) (r K) { // imag x
	xt := tp(x)
	if xt > Zt {
		return Ech(33, l1(x))
	}
	if xt == Zt {
		xp := 8 + int32(x)
		n := nn(x)
		r = mk(Ft, n)
		rp := int32(r)
		e := rp + 8*n
		for rp < e {
			SetI64(rp, I64(xp))
			xp += 16
			rp += 8
		}
		dx(x)
		return r
	}
	dx(x)
	if xt == zt {
		return Kf(F64(int32(x) + 8))
	}
	if xt < zt {
		return Kf(0.0)
	} else {
		return ntake(nn(x), Kf(0.0))
	}
}
func Cpx(x, y K) K { return Add(x, Mul(Kz(0.0, 1.0), y)) } // x imag y
func Cnj(x K) K { // conj x
	xt := tp(x)
	if xt > Zt {
		return Ech(34, l1(x))
	}
	if xt&15 < zt {
		return x
	}
	xt = tp(x)
	xp := int32(x)
	if xt == zt {
		dx(x)
		return Kz(F64(xp), -F64(xp+8))
	}
	x = use(x)
	xp = 8 + int32(x)
	e := xp + 16*nn(x)
	for xp < e {
		SetF64(xp, -F64(xp))
		xp += 16
	}
	return x
}

func nd(f, ff int32, x, y K) (r K) { //dyadic
	var av int32
	var t T
	r, t, x, y = dctypes(x, y)
	if r != 0 {
		return key(r, Func[64+ff].(f2)(x, y), t)
	}
	if t == Lt {
		return Ech(K(ff), l2(x, y))
	}
	x, y = uptypes(x, y, 1)
	xp, yp := int32(x), int32(y)
	av, t = conform(x, y)
	if av == 0 { // atom-atom
		switch t - 2 {
		case 0: // ct
			return Kc(Func[f].(f2i)(xp, yp))
		case 1: // it
			return Ki(Func[f].(f2i)(xp, yp))
		case 2: // st
			return trap(Type)
		case 3:
			dx(x)
			dx(y)
			return Kf(Func[1+f].(f2f)(F64(xp), F64(yp)))
		default:
			dx(x)
			dx(y)
			return Kz(Func[2+f].(f2z)(F64(xp), F64(xp+8), F64(yp), F64(yp+8)))
		}
	}
	if av == 1 { // atom-vector
		r = use1(y)
		if nn(r) == 0 {
			dx(x)
			dx(y)
			return r
		}
		e := ep(r)
		yp, rp := int32(y), int32(r)
		switch t - 2 {
		case 0: // ct
			Func[3+f].(f4i)(xp, yp, rp, e)
		case 1: // it
			Func[4+f].(f4i)(xp, yp, rp, e)
		case 2: // st
			trap(Type)
		case 3: // ft
			Func[5+f].(ff3i)(F64(xp), yp, rp, e)
		default: // zt
			Func[6+f].(fF3i)(F64(xp), F64(xp+8), yp, rp, e)
		}
		dx(x)
		dx(y)
		return r
	} else { // vector-vector
		r = use2(x, y)
		if nn(r) == 0 {
			dx(x)
			dx(y)
			return r
		}
		rp, xp, yp := int32(r), int32(x), int32(y)
		e := ep(r)
		switch t - 2 {
		case 0: // ct
			Func[7+f].(f4i)(xp, yp, rp, e)
		case 1: // it
			Func[8+f].(f4i)(xp, yp, rp, e)
		case 2: // st
			trap(Type)
		case 3: // ft
			Func[9+f].(f4i)(xp, yp, rp, e)
		default: // zt
			Func[10+f].(f4i)(xp, yp, rp, e)
		}
		dx(x)
		dx(y)
		return r
	}
}
func nc(f, ff int32, x, y K) (r K) { //compare
	var av int32
	var t T
	r, t, x, y = dctypes(x, y)
	if r != 0 {
		return key(r, nc(f, ff, x, y), t)
	}
	if t == Lt {
		return Ech(K(ff), l2(x, y))
	}
	x, y = uptypes(x, y, 0)
	xp, yp := int32(x), int32(y)
	av, t = conform(x, y)
	if av == 0 { // atom-atom
		switch t - 2 {
		case 0: // ct
			return Ki(Func[f].(f2i)(xp, yp))
		case 1: // it
			return Ki(Func[f].(f2i)(xp, yp))
		case 2: // st
			return Ki(Func[f].(f2i)(xp, yp))
		case 3:
			dx(x)
			dx(y)
			return Ki(Func[1+f].(f2c)(F64(xp), F64(yp)))
		default:
			dx(x)
			dx(y)
			return Ki(Func[2+f].(f2d)(F64(xp), F64(xp+8), F64(yp), F64(yp+8)))
		}
	} else if av == 1 { // atom-vector
		yn := nn(y)
		r = mk(It, yn)
		if yn == 0 {
			dx(x)
			dx(y)
			return r
		}
		rp := int32(r)
		e := ep(r)
		switch t - 2 {
		case 0: // ct
			Func[3+f].(f4i)(xp, yp, rp, e)
		case 1: // it
			Func[4+f].(f4i)(xp, yp, rp, e)
		case 2: // st
			Func[4+f].(f4i)(xp, yp, rp, e)
		case 3: // ft
			dx(x)
			Func[5+f].(ff3i)(F64(xp), yp, rp, e)
		default: // zt
			dx(x)
			Func[6+f].(fF3i)(F64(xp), F64(xp+8), yp, rp, e)
		}
		dx(y)
		return r
	}
	if av == 3 {
		xn := nn(x)
		r = mk(It, xn)
		if xn == 0 {
			dx(x)
			dx(y)
			return r
		}
		rp := int32(r)
		e := ep(r)
		switch t - 2 {
		case 0: // ct
			Func[7+f].(f4i)(xp, yp, rp, e)
		case 1: // it
			Func[8+f].(f4i)(xp, yp, rp, e)
		case 2: // st
			Func[8+f].(f4i)(xp, yp, rp, e)
		case 3: // ft
			dx(y)
			Func[9+f].(f2Ff)(xp, F64(yp), rp, e)
		default: // zt
			dx(y)
			Func[10+f].(f2Zz)(xp, F64(yp), F64(yp+8), rp, e)
		}
		dx(x)
		return r
	} else { // vector-vector
		n := nn(x)
		if t == ct {
			r = use2(x, y)
			r = K(It)<<59 | K(uint32(r))
		} else {
			r = mk(It, nn(x))
		}
		if n == 0 {
			dx(x)
			dx(y)
			return r
		}
		// t   1  2  3  4  5  6
		// f+ 11 11 12 12 13 14
		Func[f+11+I32B(t == zt)+int32(t-1)/2].(f4i)(xp, yp, int32(r), ep(r))
		dx(x)
		dx(y)
		return r
	}
}
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
func Add(x, y K) K {
	if tp(y) < 16 {
		return nd(234, 2, y, x)
	}
	return nd(234, 2, x, y)
}
func addi(x, y int32) int32                          { return x + y }
func addf(x, y float64) float64                      { return x + y }
func addz(xr, xi, yr, yi float64) (float64, float64) { return xr + yr, xi + yi }
func addcC(x, yp, rp, e int32) {
	for rp < e {
		SetI8(rp, x+I8(yp))
		yp++
		rp++
		continue
	}
}
func addiI(x, yp, rp, e int32) {
	for rp < e {
		SetI32(rp, x+I32(yp))
		yp += 4
		rp += 4
		continue
	}
}
func addfF(x float64, yp, rp, e int32) {
	for rp < e {
		SetF64(rp, x+F64(yp))
		yp += 8
		rp += 8
		continue
	}
}
func addzZ(re, im float64, yp, rp, e int32) {
	for rp < e {
		SetF64(rp, re+F64(yp))
		SetF64(rp+8, im+F64(yp+8))
		yp += 16
		rp += 16
		continue
	}
}
func addC(xp, yp, rp, e int32) {
	for rp < e {
		SetI8(rp, I8(xp)+I8(yp))
		xp++
		yp++
		rp++
		continue
	}
}
func addI(xp, yp, rp, e int32) {
	for rp < e {
		SetI32(rp, I32(xp)+I32(yp))
		xp += 4
		yp += 4
		rp += 4
		continue
	}
}
func addF(xp, yp, rp, e int32) {
	for rp < e {
		SetF64(rp, F64(xp)+F64(yp))
		xp += 8
		yp += 8
		rp += 8
		continue
	}
}

func Sub(x, y K) K {
	if tp(y) < 16 {
		return nd(234, 2, Neg(y), x)
	}
	return nd(245, 3, x, y)
}
func subi(x, y int32) int32     { return x - y }
func subf(x, y float64) float64 { return x - y }

//func subz(xr, xi, yr, yi float64) (float64, float64) { return xr - yr, xi - yi }
func subcC(x, yp, rp, e int32) {
	for rp < e {
		SetI8(rp, x-I8(yp))
		yp++
		rp++
		continue
	}
}
func subiI(x, yp, rp, e int32) {
	for rp < e {
		SetI32(rp, x-I8(yp))
		yp += 4
		rp += 4
		continue
	}
}
func subfF(x float64, yp, rp, e int32) {
	for rp < e {
		SetF64(rp, x-F64(yp))
		yp += 8
		rp += 8
		continue
	}
}
func subzZ(re, im float64, yp, rp, e int32) { //addzZ
	for rp < e {
		SetF64(rp, re-F64(yp))
		SetF64(rp+8, im-F64(yp+8))
		yp += 16
		rp += 16
		continue
	}
}

func subC(xp, yp, rp, e int32) {
	for rp < e {
		SetI8(rp, I8(xp)-I8(yp))
		xp++
		yp++
		rp++
		continue
	}
}
func subI(xp, yp, rp, e int32) {
	for rp < e {
		SetI32(rp, I32(xp)-I32(yp))
		xp += 4
		yp += 4
		rp += 4
		continue
	}
}
func subF(xp, yp, rp, e int32) {
	for rp < e {
		SetF64(rp, F64(xp)-F64(yp))
		xp += 8
		yp += 8
		rp += 8
		continue
	}
}

func Mul(x, y K) K {
	xt, yt := tp(x), tp(y)
	if xt < zt && yt == Zt {
		return scalez(x, y)
	}
	if yt < zt && xt == Zt {
		return scalez(y, x)
	}
	if yt < 16 {
		return nd(256, 4, y, x)
	}
	return nd(256, 4, x, y)
}
func muli(x, y int32) int32                          { return x * y }
func mulf(x, y float64) float64                      { return x * y }
func mulz(xr, xi, yr, yi float64) (float64, float64) { return xr*yr - xi*yi, xr*yi + xi*yr }
func mulcC(x, yp, rp, e int32) {
	for rp < e {
		SetI8(rp, x*I8(yp))
		yp++
		rp++
		continue
	}
}
func muliI(x, yp, rp, e int32) {
	for rp < e {
		SetI32(rp, x*I32(yp))
		yp += 4
		rp += 4
		continue
	}
}
func mulfF(x float64, yp, rp, e int32) {
	for rp < e {
		SetF64(rp, x*F64(yp))
		yp += 8
		rp += 8
		continue
	}
}
func mulzZ(re, im float64, yp, rp, e int32) {
	for rp < e {
		xx, yy := mulz(re, im, F64(yp), F64(yp+8))
		SetF64(rp, xx)
		SetF64(rp+8, yy)
		yp += 16
		rp += 16
		continue
	}
}
func mulC(xp, yp, rp, e int32) {
	for rp < e {
		SetI8(rp, I8(xp)*I8(yp))
		xp++
		yp++
		rp++
		continue
	}
}
func mulI(xp, yp, rp, e int32) {
	for rp < e {
		SetI32(rp, I32(xp)*I32(yp))
		xp += 4
		yp += 4
		rp += 4
		continue
	}
}
func mulF(xp, yp, rp, e int32) {
	for rp < e {
		SetF64(rp, F64(xp)*F64(yp))
		xp += 8
		yp += 8
		rp += 8
		continue
	}
}
func mulZ(xp, yp, rp, e int32) {
	for rp < e {
		xx, yy := mulz(F64(xp), F64(xp+8), F64(yp), F64(yp+8))
		SetF64(rp, xx)
		SetF64(rp+8, yy)
		xp += 16
		yp += 16
		rp += 16
		continue
	}
}
func scale(s float64, x K) (r K) {
	if tp(x) < Ft {
		x = uptype(x, ft)
	}
	r = use1(x)
	if nn(r) == 0 {
		dx(x)
		return r
	}
	e := ep(r)
	xp, rp := int32(x), int32(r)
	mulfF(s, xp, rp, e)
	dx(x)
	return r
}
func scalez(x, z K) K { // xt<=ft, z:Zt
	if tp(x) < ft {
		x = uptype(x, ft)
	}
	s := F64(int32(x))
	dx(x)
	return scale(s, z)
}

func Div(x, y K) K {
	xt, yt := tp(x), tp(y)
	if xt&15 < ft && yt&15 < ft {
		return idiv(x, y, 0)
	}
	if yt < 16 && xt > 16 && xt < Lt {
		if yt < ft {
			y = uptype(y, ft)
		} else if yt == zt {
			return Mul(Div(Kz(1.0, 0.0), y), x)
		}
		s := 1.0 / F64(int32(y))
		dx(y)
		return scale(s, x)
	}
	return nd(267, 5, x, y)
}
func divi(x, y int32) int32     { return x / y }
func divf(x, y float64) float64 { return x / y }
func divz(xr, xi, yr, yi float64) (e float64, f float64) {
	var r, d float64
	if F64abs(yr) >= F64abs(yi) {
		r = yi / yr
		d = yr + r*yi
		e = (xr + xi*r) / d
		f = (xi - xr*r) / d
	} else {
		r = yr / yi
		d = yi + r*yr
		e = (xr*r + xi) / d
		f = (xi*r - xr) / d
	}
	return e, f
}
func divfF(x float64, yp, rp, e int32) {
	for rp < e {
		SetF64(rp, x/F64(yp))
		yp += 8
		rp += 8
		continue
	}
}
func divzZ(re, im float64, yp, rp, e int32) {
	for rp < e {
		xx, yy := divz(re, im, F64(yp), F64(yp+8))
		SetF64(rp, xx)
		SetF64(rp+8, yy)
		yp += 16
		rp += 16
		continue
	}
}
func divF(xp, yp, rp, e int32) {
	for rp < e {
		SetF64(rp, F64(xp)/F64(yp))
		rp += 8
		xp += 8
		yp += 8
		continue
	}
}
func divZ(xp, yp, rp, e int32) {
	for rp < e {
		xx, yy := divz(F64(xp), F64(xp+8), F64(yp), F64(yp+8))
		SetF64(rp, xx)
		SetF64(rp+8, yy)
		xp += 16
		yp += 16
		rp += 16
		continue
	}
}
func idiv(x, y K, mod int32) (r K) {
	x, y = uptypes(x, y, 1)
	av, t := conform(x, y)
	if t != it {
		trap(Type)
	}
	xp := int32(x)
	yp := int32(y)
	switch av {
	case 0: //a%a
		if mod != 0 {
			r = Ki(xp % yp)
		} else {
			r = Ki(xp / yp)
		}
		return r
	case 1: //a%v
		r = use(y)
		rp := int32(r)
		e := rp + 4*nn(r)
		if mod != 0 {
			for rp < e {
				SetI32(rp, xp%I32(rp))
				rp += 4
			}
		} else {
			for rp < e {
				SetI32(rp, xp/I32(rp))
				rp += 4
			}
		}
		return r
	case 2: //v%v
		r = use2(x, y)
		rp := int32(r)
		e := rp + 4*nn(r)
		if mod != 0 {
			for rp < e {
				SetI32(rp, I32(xp)%I32(yp))
				xp += 4
				yp += 4
				rp += 4
			}
		} else {
			for rp < e {
				SetI32(rp, I32(xp)/I32(yp))
				xp += 4
				yp += 4
				rp += 4
			}
		}
		dx(x)
		dx(y)
		return r
	default: // v%a
		x = use(x)
		xp = int32(x)
		xn := nn(x)
		e := xp + 4*xn
		if yp > 0 && xn > 0 && mod == 0 {
			divIi(xp, yp, e)
		}
		if mod != 0 {
			for xp < e {
				SetI32(xp, I32(xp)%yp)
				xp += 4
			}
		}
		return x
	}
}
func divIi(xp, yp, e int32) {
	s := int32(31) - I32clz(uint32(yp))
	if yp == int32(1)<<s { // x % powers of 2
		for xp < e {
			SetI32(xp, I32(xp)>>s)
			xp += 4
			continue
		}
		return
	}
	for xp < e {
		SetI32(xp, I32(xp)/yp)
		xp += 4
	}
}
func Mod(x, y K) K {
	xt, yt := tp(x), tp(y)
	if xt&15 < ft && yt&15 < ft {
		return idiv(x, y, 1)
	}
	if xt >= Lt || yt >= Lt {
		return nd(0, 41, x, y)
	} else {
		return trap(Type)
	}
}

func Min(x, y K) (r K) {
	if tp(y) < 16 {
		return nd(278, 7, y, x)
	}
	return nd(278, 7, x, y)
}
func mini(x, y int32) int32 {
	if x < y {
		return x
	}
	return y
}
func minf(x, y float64) float64 { return F64min(x, y) }
func minz(xr, xi, yr, yi float64) (float64, float64) {
	if ltz(xr, xi, yr, yi) != 0 {
		return xr, xi
	}
	return yr, yi
}
func mincC(x, yp, rp, e int32) {
	for rp < e {
		SetI8(rp, mini(x, I8(yp)))
		yp++
		rp++
		continue
	}
}
func miniI(x, yp, rp, e int32) {
	for rp < e {
		SetI32(rp, mini(x, I32(yp)))
		yp += 4
		rp += 4
		continue
	}
}
func minfF(x float64, yp, rp, e int32) {
	for rp < e {
		SetF64(rp, minf(x, F64(yp)))
		yp += 8
		rp += 8
		continue
	}
}
func minzZ(re, im float64, yp, rp, e int32) {
	for rp < e {
		xx, yy := minz(re, im, F64(yp), F64(yp+8))
		SetF64(rp, xx)
		SetF64(rp+8, yy)
		yp += 16
		rp += 16
		continue
	}
}
func minC(xp, yp, rp, e int32) {
	for rp < e {
		SetI8(rp, mini(I8(xp), I8(yp)))
		xp++
		yp++
		rp++
		continue
	}
}
func minI(xp, yp, rp, e int32) {
	for rp < e {
		SetI32(rp, mini(I32(xp), I32(yp)))
		xp += 4
		yp += 4
		rp += 4
		continue
	}
}
func minF(xp, yp, rp, e int32) {
	for rp < e {
		SetF64(rp, minf(F64(xp), F64(yp)))
		xp += 8
		yp += 8
		rp += 8
		continue
	}
}
func minZ(xp, yp, rp, e int32) {
	for rp < e {
		xx, yy := minz(F64(xp), F64(xp+8), F64(yp), F64(yp+8))
		SetF64(rp, xx)
		SetF64(rp+8, yy)
		xp += 16
		yp += 16
		rp += 16
		continue
	}
}

func Max(x, y K) (r K) {
	if tp(y) < 16 {
		return nd(289, 8, y, x)
	}
	return nd(289, 8, x, y)
}
func maxi(x, y int32) int32 {
	if x > y {
		return x
	} else {
		return y
	}
}
func maxf(x, y float64) float64 { return F64max(x, y) }
func maxz(xr, xi, yr, yi float64) (float64, float64) {
	if gtz(xr, xi, yr, yi) != 0 {
		return xr, xi
	}
	return yr, yi
}
func maxcC(x, yp, rp, e int32) {
	for rp < e {
		SetI8(rp, maxi(x, I8(yp)))
		yp++
		rp++
		continue
	}
}
func maxiI(x, yp, rp, e int32) {
	for rp < e {
		SetI32(rp, maxi(x, I32(yp)))
		yp += 4
		rp += 4
		continue
	}
}
func maxfF(x float64, yp, rp, e int32) {
	for rp < e {
		SetF64(rp, F64max(x, F64(yp)))
		yp += 8
		rp += 8
		continue
	}
}
func maxzZ(re, im float64, yp, rp, e int32) {
	for rp < e {
		xx, yy := maxz(re, im, F64(yp), F64(yp+8))
		SetF64(rp, xx)
		SetF64(rp+8, yy)
		yp += 16
		rp += 16
		continue
	}
}
func maxC(xp, yp, rp, e int32) {
	for rp < e {
		SetI8(rp, maxi(I8(xp), I8(yp)))
		xp++
		yp++
		rp++
		continue
	}
}
func maxI(xp, yp, rp, e int32) {
	for rp < e {
		SetI32(rp, maxi(I32(xp), I32(yp)))
		xp += 4
		yp += 4
		rp += 4
		continue
	}
}
func maxF(xp, yp, rp, e int32) {
	for rp < e {
		SetF64(rp, F64max(F64(xp), F64(yp)))
		xp += 8
		yp += 8
		rp += 8
		continue
	}
}
func maxZ(xp, yp, rp, e int32) {
	for rp < e {
		xx, yy := maxz(F64(xp), F64(xp+8), F64(yp), F64(yp+8))
		SetF64(rp, xx)
		SetF64(rp+8, yy)
		xp += 16
		yp += 16
		rp += 16
		continue
	}
}

func Eql(x, y K) K                     { return nc(338, 11, x, y) }
func eqi(x, y int32) int32             { return I32B(x == y) }
func eqf(x, y float64) int32           { return I32B((x != x) && (y != y) || x == y) }
func eqz(xr, xi, yr, yi float64) int32 { return eqf(xr, yr) & eqf(xi, yi) }
func eqcC(v, yp, rp, e int32) {
	for rp < e {
		SetI32(rp, I32B(v == I8(yp)))
		yp++
		rp += 4
		continue
	}
}
func eqiI(x int32, yp, rp, e int32) {
	for rp < e {
		SetI32(rp, I32B(x == I32(yp)))
		yp += 4
		rp += 4
		continue
	}
}
func eqfF(x float64, yp, rp, e int32) {
	for rp < e {
		SetI32(rp, eqf(x, F64(yp)))
		yp += 8
		rp += 4
		continue
	}
}
func eqzZ(re, im float64, yp, rp, e int32) {
	for rp < e {
		SetI32(rp, eqz(re, im, F64(yp), F64(yp+8)))
		yp += 16
		rp += 4
		continue
	}
}
func eqCc(xp, y, rp, e int32)                    { eqcC(y, xp, rp, e) }
func eqIi(xp, y, rp, e int32)                    { eqiI(y, xp, rp, e) }
func eqFf(xp int32, y float64, rp, e int32)      { eqfF(y, xp, rp, e) }
func eqZz(xp int32, re, im float64, rp, e int32) { eqzZ(re, im, xp, rp, e) }
func eqC(xp, yp, rp, e int32) {
	for rp < e {
		SetI32(rp, I32B(I8(xp) == I8(yp)))
		xp++
		yp++
		rp += 4
		continue
	}
}
func eqI(xp, yp, rp, e int32) {
	for rp < e {
		SetI32(rp, I32B(I32(xp) == I32(yp)))
		xp += 4
		yp += 4
		rp += 4
		continue
	}
}
func eqF(xp, yp, rp, e int32) {
	for rp < e {
		SetI32(rp, eqf(F64(xp), F64(yp)))
		xp += 8
		yp += 8
		rp += 4
		continue
	}
}
func eqZ(xp, yp, rp, e int32) {
	for rp < e {
		SetI32(rp, eqz(F64(xp), F64(xp+8), F64(yp), F64(yp+8)))
		xp += 16
		yp += 16
		rp += 4
		continue
	}
}

func Les(x, y K) K { // x<y   `file<c
	if tp(x) == st && tp(y) == Ct {
		if int32(x) == 0 {
			write(rx(y))
			return y
		}
		return writefile(cs(x), y)
	}
	return nc(308, 9, x, y)
}
func lti(x, y int32) int32   { return I32B(x < y) }
func ltf(x, y float64) int32 { return I32B(x < y || x != x) }
func ltz(xr, xi, yr, yi float64) int32 {
	if eqf(xr, yr) != 0 {
		return ltf(xi, yi)
	}
	return ltf(xr, yr)
}
func ltcC(v, yp, rp, e int32) {
	for rp < e {
		SetI32(rp, I32B(v < I8(yp)))
		yp++
		rp += 4
		continue
	}
}
func ltiI(x int32, yp, rp, e int32) {
	for rp < e {
		SetI32(rp, I32B(x < I32(yp)))
		yp += 4
		rp += 4
		continue
	}
}
func ltfF(x float64, yp, rp, e int32) {
	for rp < e {
		SetI32(rp, ltf(x, F64(yp)))
		yp += 8
		rp += 4
		continue
	}
}
func ltzZ(re, im float64, yp, rp, e int32) {
	for rp < e {
		if re == F64(yp) {
			SetI32(rp, ltf(im, F64(yp+8)))
		} else {
			SetI32(rp, ltf(re, F64(yp)))
		}
		yp += 16
		rp += 4
		continue
	}
}
func ltCc(xp int32, v, rp, e int32) {
	for rp < e {
		SetI32(rp, I32B(I8(xp) < v))
		xp++
		rp += 4
		continue
	}
}
func ltIi(xp, y int32, rp, e int32) {
	for rp < e {
		SetI32(rp, I32B(I32(xp) < y))
		xp += 4
		rp += 4
		continue
	}
}
func ltFf(xp int32, y float64, rp, e int32) {
	for rp < e {
		SetI32(rp, ltf(F64(xp), y))
		xp += 8
		rp += 4
		continue
	}
}
func ltZz(xp int32, re, im float64, rp, e int32) {
	for rp < e {
		if F64(xp) == re {
			SetI32(rp, ltf(F64(xp+8), im))
		} else {
			SetI32(rp, ltf(F64(xp), re))
		}
		xp += 16
		rp += 4
		continue
	}
}
func ltC(xp, yp, rp, e int32) {
	for rp < e {
		SetI32(rp, I32B(I8(xp) < I8(yp)))
		xp++
		yp++
		rp += 4
		continue
	}
}
func ltI(xp, yp, rp, e int32) {
	for rp < e {
		SetI32(rp, I32B(I32(xp) < I32(yp)))
		xp += 4
		yp += 4
		rp += 4
		continue
	}
}
func ltF(xp, yp, rp, e int32) {
	for rp < e {
		SetI32(rp, ltf(F64(xp), F64(yp)))
		xp += 8
		yp += 8
		rp += 4
		continue
	}
}
func ltZ(xp, yp, rp, e int32) {
	for rp < e {
		SetI32(rp, ltz(F64(xp), F64(xp+8), F64(yp), F64(yp+8)))
		xp += 16
		yp += 16
		rp += 4
		continue
	}
}

func Mor(x, y K) K           { return nc(323, 10, x, y) }
func gti(x, y int32) int32   { return I32B(x > y) }
func gtf(x, y float64) int32 { return I32B(x > y || y != y) }
func gtz(xr, xi, yr, yi float64) int32 {
	if eqf(xr, yr) != 0 {
		return gtf(xi, yi)
	}
	return gtf(xr, yr)
}
func gtcC(v, yp, rp, e int32) {
	for rp < e {
		SetI32(rp, I32B(v > I8(yp)))
		yp++
		rp += 4
		continue
	}
}
func gtiI(x, yp, rp, e int32) {
	for rp < e {
		SetI32(rp, I32B(x > I32(yp)))
		yp += 4
		rp += 4
		continue
	}
}
func gtfF(x float64, yp, rp, e int32) {
	for rp < e {
		SetI32(rp, gtf(x, F64(yp)))
		yp += 8
		rp += 4
		continue
	}
}
func gtzZ(re, im float64, yp, rp, e int32) {
	for rp < e {
		if re == F64(yp) {
			SetI32(rp, gtf(im, F64(yp+8)))
		} else {
			SetI32(rp, gtf(re, F64(yp)))
		}
		yp += 16
		rp += 4
		continue
	}
}
func gtCc(xp, v, rp, e int32) {
	for rp < e {
		SetI32(rp, I32B(I8(xp) > v))
		xp++
		rp += 4
		continue
	}
}
func gtIi(xp, y, rp, e int32) {
	for rp < e {
		SetI32(rp, I32B(I32(xp) > y))
		xp += 4
		rp += 4
		continue
	}
}
func gtFf(xp int32, y float64, rp, e int32) {
	for rp < e {
		SetI32(rp, gtf(F64(xp), y))
		xp += 8
		rp += 4
		continue
	}
}
func gtZz(xp int32, re, im float64, rp, e int32) {
	for rp < e {
		if F64(xp) == re {
			SetI32(rp, gtf(F64(xp+8), im))
		} else {
			SetI32(rp, gtf(F64(xp), re))
		}
		xp += 16
		rp += 4
		continue
	}
}
func gtC(xp, yp, rp, e int32) {
	for rp < e {
		SetI32(rp, I32B(I8(xp) > I8(yp)))
		xp++
		yp++
		rp += 4
		continue
	}
}
func gtI(xp, yp, rp, e int32) {
	for rp < e {
		SetI32(rp, I32B(I32(xp) > I32(yp)))
		xp += 4
		yp += 4
		rp += 4
		continue
	}
}
func gtF(xp, yp, rp, e int32) {
	for rp < e {
		SetI32(rp, gtf(F64(xp), F64(yp)))
		xp += 8
		yp += 8
		rp += 4
		continue
	}
}
func gtZ(xp, yp, rp, e int32) {
	for rp < e {
		SetI32(rp, gtz(F64(xp), F64(xp+8), F64(yp), F64(yp+8)))
		xp += 16
		yp += 16
		rp += 4
		continue
	}
}

//func isnan(x float64) bool { return x != x }

func Ang(x K) (r K) { // angle x
	xt := tp(x)
	if xt > Zt {
		return Ech(35, l1(x))
	}
	if xt < zt {
		dx(x)
		return Kf(0)
	}
	xp := int32(x)
	if xt == zt {
		dx(x)
		return Kf(ang2(F64(xp+8), F64(xp)))
	}
	n := nn(x)
	if xt == Zt {
		r = mk(Ft, n)
		rp := int32(r)
		e := rp + 8*n
		for rp < e {
			SetF64(rp, ang2(F64(xp+8), F64(xp)))
			xp += 16
			rp += 8
		}
	} else {
		r = ntake(n, Kf(0))
	}
	dx(x)
	return r
}
func Rot(x, y K) (r K) { // r@deg
	if tp(x) > Zt {
		return Ech(35, l2(x, y))
	}
	x = uptype(x, zt)
	if y == 0 {
		return x
	}
	if tp(y)&15 > ft {
		trap(Type)
	}
	y = uptype(y, ft)
	yt := tp(y)
	yp := int32(y)
	if yt == ft {
		r = Kz(cosin(F64(yp)))
	} else {
		yn := nn(y)
		r = mk(Zt, yn)
		rp := int32(r)
		for i := int32(0); i < yn; i++ {
			c, s := cosin(F64(yp))
			SetF64(rp, c)
			SetF64(rp+8, s)
			yp += 8
			rp += 16
		}
	}
	dx(y)
	return Mul(r, x)
}
func Sin(x K) K { return nf(44, x, 0) } // sin x
func Cos(x K) K { return nf(39, x, 0) } // cos x
func Exp(x K) K { return nf(42, x, 0) } // exp x
func Log(x K) K { return nf(43, x, 0) } // log x
func Pow(y, x K) K { // x^y
	if tp(x)&15 == it {
		if tp(y) == it {
			if int32(y) >= 0 {
				return ipow(x, int32(y))
			}
		}
	}
	return nf(106, x, y)
}
func Lgn(x, y K) K { // n log y
	xf := fk(x)
	if xf == 10.0 {
		xf = 0.4342944819032518
	} else if xf == 2.0 {
		xf = 1.4426950408889634
	} else {
		xf = 1.0 / log(xf)
	}
	return Mul(Kf(xf), Log(y))
}
func fk(x K) float64 {
	t := tp(x)
	if t == it {
		return float64(int32(x))
	}
	if t != ft {
		trap(Type)
	}
	dx(x)
	return F64(int32(x))
}
func nf(f int32, x, y K) (r K) {
	xt := tp(x)
	if xt >= Lt {
		if y == 0 {
			return Ech(K(f), l1(x))
		} else {
			return Ech(K(f-64), l2(y, x))
		}
	}
	var yf float64
	if y != 0 {
		yf = fk(y)
	}
	if xt&15 < ft {
		x = uptype(x, ft)
		xt = tp(x)
	}
	xp := int32(x)
	if xt == ft {
		r = Kf(0)
		ff(f, int32(r), xp, 1, yf)
	} else {
		xn := nn(x)
		r = mk(Ft, xn)
		if xn > 0 {
			ff(f, int32(r), xp, xn, yf)
		}
	}
	dx(x)
	return r
}
func ff(f, rp, xp, n int32, yf float64) {
	e := xp + 8*n
	switch f - 42 {
	case 0:
		for xp < e {
			SetF64(rp, exp(F64(xp)))
			rp += 8
			xp += 8
			continue
		}
	case 1:
		for xp < e {
			SetF64(rp, log(F64(xp)))
			rp += 8
			xp += 8
			continue
		}
	default:
		if f == 106 { // pow 42+64
			for xp < e {
				SetF64(rp, pow(F64(xp), yf))
				rp += 8
				xp += 8
				continue
			}
		} else { // sin cos
			for xp < e {
				c, s := cosin_(F64(xp))
				SetF64(rp, s)
				if f == 39 {
					SetF64(rp, c)
				}
				rp += 8
				xp += 8
				continue
			}
		}
	}
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
