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

type c2cC = func(I8x16, I8x16, int32, int32, int32)
type f2cC = func(I8x16, int32, int32, int32)
type f2iI = func(I32x4, int32, int32, int32)
type f2vF = func(F64x2, int32, int32, int32)
type f2fF = func(float64, int32, int32, int32)
type f2zZ = func(float64, float64, int32, int32, int32)

type c2Cc = func(int32, I8x16, I8x16, int32, int32)
type f2Cc = func(int32, I8x16, int32, int32)
type f2ii = func(int32, int32, int32, int32)
type f2Ff = func(int32, float64, int32, int32)
type f2Zz = func(int32, float64, float64, int32, int32)

type f2v = func(int32, int32, int32, int32)
type f2vc = func(I8x16, int32, int32, int32, int32)
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

func nm(f int32, x K) (r K) {
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
		switch xt - 1 {
		case 0:
			return Ki(Func[f].(f1i)(xp))
		case 1:
			return Kc(Func[f].(f1i)(xp))
		case 2:
			return Ki(Func[1+f].(f1i)(xp))
		case 3:
			return trap(Type)
		case 4:
			r = Kf(Func[2+f].(f1f)(F64(xp)))
			dx(x)
			return r
		case 5:
			r = Kz(Func[3+f].(f1z)(F64(xp), F64(xp+8)))
			dx(x)
			return r
		default:
			return trap(Type)
		}
	}
	if xt == Bt {
		x, xt = uptype(x, it), It
		xp = int32(x)
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
		Func[4+f].(f1_)(xp, rp, e)
	case 1:
		Func[5+f].(f1_)(xp, rp, e)
	case 2:
		trap(Type)
	case 3:
		Func[6+f].(f1_)(xp, rp, e)
	case 4:
		Func[7+f].(f1_)(xp, rp, e)
	default:
		trap(Type)
	}
	dx(x)
	return r
}

func Neg(x K) K { return nm(220, x) }
func negc(x int32) int32 { // lower
	if x > 'A' && x < 'Z' {
		x += 32
	}
	return x
}
func negi(x int32) int32                   { return -x }
func negf(x float64) float64               { return -x }
func negz(x, y float64) (float64, float64) { return -x, -y }
func negC(xp, rp, e int32) {
	for rp < e {
		SetI8(rp, absc(I8(xp)))
		xp++
		rp++
		continue
	}
}
func negI(xp, rp, e int32) {
	for rp < e {
		I32x4store(rp, I32x4load(xp).Neg())
		xp += 16
		rp += 16
		continue
	}
}
func negF(xp, rp, e int32) {
	for rp < e {
		F64x2store(rp, F64x2load(xp).Neg())
		xp += 16
		rp += 16
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
	return nm(228, x)
}
func absc(x int32) int32 { // upper
	if x > 'a' && x < 'z' {
		x -= 32
	}
	return x
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
		SetI8(rp, absc(I8(xp)))
		xp++
		rp++
		continue
	}
}
func absI(xp, rp, e int32) {
	for rp < e {
		I32x4store(rp, I32x4load(xp).Abs())
		xp += 16
		rp += 16
		continue
	}
}
func absF(xp, rp, e int32) {
	for rp < e {
		F64x2store(rp, F64x2load(xp).Abs())
		xp += 16
		rp += 16
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
func Hypot(x, y K) (r K) { // e.g.  norm:0. abs/x
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

func Sqr(x K) (r K) {
	if tp(x)&15 != ft {
		x = Add(Kf(0), x)
	}
	return nm(383, x)
}
func sqrf(x float64) float64 { return F64sqrt(x) }
func sqrF(xp, rp, e int32) {
	for rp < e {
		F64x2store(rp, F64x2load(xp).Sqrt())
		xp += 16
		rp += 16
		continue
	}
}

func Imag(x K) (r K) { // imag x
	xt := tp(x)
	if xt > Zt {
		return Ech(33, l1(x))
	}
	if xt < zt {
		return zero(xt)
	}
	if xt == zt {
		dx(x)
		return Kf(F64(int32(x) + 8))
	}
	if xt < Zt {
		x = uptype(x, zt)
	}
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
func Cmpl(x, y K) (r K) { return Add(x, Mul(Kz(0.0, 1.0), y)) } // x imag y
func Conj(x K) (r K) { // conj x
	xt := tp(x)
	if xt > Zt {
		return Ech(34, l1(x))
	}
	if xt&15 < zt {
		x = uptype(x, zt)
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

func nd(f, ff int32, x, y K) (r K) {
	var av int32
	var t T
	r, t, x, y = dctypes(x, y)
	if r != 0 {
		return key(r, nd(f, ff, x, y), t)
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
			return Ki(Func[1+f].(f2i)(xp, yp))
		case 2: // st
			return trap(Type)
		case 3:
			dx(x)
			dx(y)
			return Kf(Func[2+f].(f2f)(F64(xp), F64(yp)))
		case 4:
			dx(x)
			dx(y)
			return Kz(Func[3+f].(f2z)(F64(xp), F64(xp+8), F64(yp), F64(yp+8)))
		default:
			return trap(Type)
		}
	} else if av == 1 { // atom-vector
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
			v := I8x16splat(xp)
			Func[4+f].(f2cC)(v, yp, rp, e)
		case 1: // it
			v := I32x4splat(xp)
			Func[5+f].(f2iI)(v, yp, rp, e)
		case 2: // st
			trap(Type)
		case 3: // ft
			v := F64x2splat(F64(xp))
			Func[6+f].(f2vF)(v, yp, rp, e)
		case 4: // zt
			Func[7+f].(f2zZ)(F64(xp), F64(xp+8), yp, rp, e)
		default:
			trap(Type)
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
			Func[8+f].(f2v)(xp, yp, rp, e)
		case 1: // it
			Func[9+f].(f2v)(xp, yp, rp, e)
		case 2: // st
			trap(Type)
		case 3: // ft
			Func[10+f].(f2v)(xp, yp, rp, e)
		case 4: // zt
			Func[11+f].(f2v)(xp, yp, rp, e)
		default:
			trap(Type)
		}
		dx(x)
		dx(y)
		return r
	}
	return r
}
func nc(f, ff int32, x, y K) (r K) {
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
		switch t - 1 {
		case 0: // bt
			return Kb(Func[f].(f2i)(xp, yp))
		case 1: // ct
			return Kb(Func[f].(f2i)(xp, yp))
		case 2: // it
			return Kb(Func[f].(f2i)(xp, yp))
		case 3: // st
			return Kb(Func[f].(f2i)(xp, yp))
		case 4:
			dx(x)
			dx(y)
			return Kb(Func[1+f].(f2c)(F64(xp), F64(yp)))
		case 5:
			dx(x)
			dx(y)
			return Kb(Func[2+f].(f2d)(F64(xp), F64(xp+8), F64(yp), F64(yp+8)))
		default:
			return trap(Type)
		}
	} else if av == 1 { // atom-vector
		yn := nn(y)
		r = mk(Bt, yn)
		if yn == 0 {
			dx(x)
			dx(y)
			return r
		}
		rp := int32(r)
		e := ep(r)
		switch t - 1 {
		case 0: // bt
			Func[3+f].(c2cC)(I8x16splat(xp), I8x16splat(1), yp, rp, e)
		case 1: // ct
			Func[3+f].(c2cC)(I8x16splat(xp), I8x16splat(1), yp, rp, e)
		case 2: // it
			Func[4+f].(f2ii)(xp, yp, rp, e)
		case 3: // st
			Func[4+f].(f2ii)(xp, yp, rp, e)
		case 4: // ft
			dx(x)
			Func[5+f].(f2fF)(F64(xp), yp, rp, e)
		case 5: // zt
			dx(x)
			Func[6+f].(f2zZ)(F64(xp), F64(xp+8), yp, rp, e)
		default:
			trap(Type)
		}
		dx(y)
		return r
	} else if av == 3 {
		xn := nn(x)
		r = mk(Bt, xn)
		if xn == 0 {
			dx(x)
			dx(y)
			return r
		}
		rp := int32(r)
		e := ep(r)
		switch t - 1 {
		case 0: // bt
			Func[7+f].(c2Cc)(xp, I8x16splat(yp), I8x16splat(1), rp, e)
		case 1: // ct
			Func[7+f].(c2Cc)(xp, I8x16splat(yp), I8x16splat(1), rp, e)
		case 2: // it
			Func[8+f].(f2ii)(xp, yp, rp, e)
		case 3: // st
			Func[8+f].(f2ii)(xp, yp, rp, e)
		case 4: // ft
			dx(y)
			Func[9+f].(f2Ff)(xp, F64(yp), rp, e)
		case 5: // zt
			dx(y)
			Func[10+f].(f2Zz)(xp, F64(yp), F64(yp+8), rp, e)
		default:
			trap(Type)
		}
		dx(x)
		return r
	} else { // vector-vector
		n := nn(x)
		if t == Bt || t == Ct {
			r = use2(x, y)
			r = K(Bt)<<59 | K(uint32(r))
		} else {
			r = mk(Bt, nn(x))
		}
		if n == 0 {
			dx(x)
			dx(y)
			return r
		}
		rp := int32(r)
		e := ep(r)
		switch t - 1 {
		case 0: // bt
			Func[11+f].(f2vc)(I8x16splat(1), xp, yp, rp, e)
		case 1: // ct
			Func[11+f].(f2vc)(I8x16splat(1), xp, yp, rp, e)
		case 2: // it
			Func[12+f].(f2v)(xp, yp, rp, e)
		case 3: // st
			Func[12+f].(f2v)(xp, yp, rp, e)
		case 4: // ft
			Func[13+f].(f2v)(xp, yp, rp, e)
		case 5: // zt
			Func[14+f].(f2v)(xp, yp, rp, e)
		default:
			trap(Type)
		}
		dx(x)
		dx(y)
		return r
	}
	return r
}
func conform(x, y K) (av int32, r T) { // 0:atom-atom 1:atom-vector, 2:vector-vector, 3:vector-atom
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
func Add(x, y K) (r K) {
	if tp(y) < 16 {
		return nd(236, 2, y, x)
	}
	return nd(236, 2, x, y)
}
func addi(x, y int32) int32                          { return x + y }
func addf(x, y float64) float64                      { return x + y }
func addz(xr, xi, yr, yi float64) (float64, float64) { return xr + yr, xi + yi }
func addcC(x I8x16, yp, rp, e int32) {
	for rp < e {
		I8x16store(rp, x.Add(I8x16load(yp)))
		yp += 16
		rp += 16
		continue
	}
}
func addiI(x I32x4, yp, rp, e int32) {
	for rp < e {
		I32x4store(rp, x.Add(I32x4load(yp)))
		yp += 16
		rp += 16
		continue
	}
}
func addfF(x F64x2, yp, rp, e int32) {
	for rp < e {
		F64x2store(rp, x.Add(F64x2load(yp)))
		yp += 16
		rp += 16
		continue
	}
}
func addzZ(re, im float64, yp, rp, e int32) {
	x := F64x2splat(re).Replace_lane1(im)
	addfF(x, yp, rp, e)
}
func addC(xp, yp, rp, e int32) {
	for rp < e {
		I8x16store(rp, I8x16load(xp).Add(I8x16load(yp)))
		xp += 16
		yp += 16
		rp += 16
		continue
	}
}
func addI(xp, yp, rp, e int32) {
	for rp < e {
		I32x4store(rp, I32x4load(xp).Add(I32x4load(yp)))
		xp += 16
		yp += 16
		rp += 16
		continue
	}
}
func addF(xp, yp, rp, e int32) {
	for rp < e {
		F64x2store(rp, F64x2load(xp).Add(F64x2load(yp)))
		xp += 16
		yp += 16
		rp += 16
		continue
	}
}
func addZ(xp, yp, rp, e int32) { addF(xp, yp, rp, e) }

func Sub(x, y K) (r K) {
	if tp(y) < 16 {
		return nd(236, 2, Neg(y), x)
	}
	return nd(248, 3, x, y)
}
func subi(x, y int32) int32                          { return x - y }
func subf(x, y float64) float64                      { return x - y }
func subz(xr, xi, yr, yi float64) (float64, float64) { return xr - yr, xi - yi }
func subcC(x I8x16, yp, rp, e int32) {
	for rp < e {
		I8x16store(rp, x.Sub(I8x16load(yp)))
		yp += 16
		rp += 16
		continue
	}
}
func subiI(x I32x4, yp, rp, e int32) {
	for rp < e {
		I32x4store(rp, x.Sub(I32x4load(yp)))
		yp += 16
		rp += 16
		continue
	}
}
func subfF(x F64x2, yp, rp, e int32) {
	for rp < e {
		F64x2store(rp, x.Sub(F64x2load(yp)))
		yp += 16
		rp += 16
		continue
	}
}
func subzZ(re, im float64, yp, rp, e int32) {
	x := F64x2splat(re).Replace_lane1(im)
	subfF(x, yp, rp, e)
}
func subC(xp, yp, rp, e int32) {
	for rp < e {
		I8x16store(rp, I8x16load(xp).Sub(I8x16load(yp)))
		xp += 16
		yp += 16
		rp += 16
		continue
	}
}
func subI(xp, yp, rp, e int32) {
	for rp < e {
		I32x4store(rp, I32x4load(xp).Sub(I32x4load(yp)))
		xp += 16
		yp += 16
		rp += 16
		continue
	}
}
func subF(xp, yp, rp, e int32) {
	for rp < e {
		F64x2store(rp, F64x2load(xp).Sub(F64x2load(yp)))
		xp += 16
		yp += 16
		rp += 16
		continue
	}
}
func subZ(xp, yp, rp, n int32) { subF(xp, yp, rp, n) }

func Mul(x, y K) (r K) {
	xt, yt := tp(x), tp(y)
	if xt < zt && yt == Zt {
		return scalez(x, y)
	}
	if yt < zt && xt == Zt {
		return scalez(y, x)
	}
	if yt < 16 {
		return nd(260, 4, y, x)
	}
	return nd(260, 4, x, y)
}
func muli(x, y int32) int32                          { return x * y }
func mulf(x, y float64) float64                      { return x * y }
func mulz(xr, xi, yr, yi float64) (float64, float64) { return xr*yr - xi*yi, xr*yi + xi*yr }
func mulcC(x I8x16, yp, rp, e int32) {
	c := x.Extract_lane_s0() // no I8x16.mul
	for rp < e {
		SetI8(rp, c*I8(yp))
		yp++
		rp++
		continue
	}
}
func muliI(x I32x4, yp, rp, e int32) {
	for rp < e {
		I32x4store(rp, x.Mul(I32x4load(yp)))
		yp += 16
		rp += 16
		continue
	}
}
func mulfF(x F64x2, yp, rp, e int32) {
	for rp < e {
		F64x2store(rp, x.Mul(F64x2load(yp)))
		yp += 16
		rp += 16
		continue
	}
}
func mulzZ(re, im float64, yp, rp, e int32) { // todo simd
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
		I32x4store(rp, I32x4load(xp).Mul(I32x4load(yp)))
		xp += 16
		yp += 16
		rp += 16
		continue
	}
}
func mulF(xp, yp, rp, e int32) {
	for rp < e {
		F64x2store(rp, F64x2load(xp).Mul(F64x2load(yp)))
		xp += 16
		yp += 16
		rp += 16
		continue
	}
}
func mulZ(xp, yp, rp, e int32) { // todo simd
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
	v := F64x2splat(s)
	xp, rp := int32(x), int32(r)
	for rp < e {
		F64x2store(rp, v.Mul(F64x2load(xp)))
		xp += 16
		rp += 16
		continue
	}
	dx(x)
	return r
}
func scalez(x, z K) (r K) { // xt<=ft, z:Zt
	if tp(x) < ft {
		x = uptype(x, ft)
	}
	s := F64(int32(x))
	dx(x)
	return scale(s, z)
}

func Div(x, y K) (r K) {
	xt, yt := tp(x), tp(y)
	if yt < 16 {
		if yt == ft && xt > 16 {
			s := 1.0 / F64(int32(y))
			dx(y)
			return scale(s, x)
		}
	}
	if xt&15 < ft && yt&15 < ft {
		return idiv(x, y) // no simd for ints
	}
	return nd(272, 5, x, y)
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
func divfF(x F64x2, yp, rp, e int32) {
	for rp < e {
		F64x2store(rp, x.Div(F64x2load(yp)))
		yp += 16
		rp += 16
		continue
	}
}
func divzZ(re, im float64, yp, rp, e int32) { // todo simd
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
		F64x2store(rp, F64x2load(xp).Div(F64x2load(yp)))
		rp += 16
		xp += 16
		yp += 16
		continue
	}
}
func divZ(xp, yp, rp, e int32) { // todo simd
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
func idiv(x, y K) (r K) {
	av, t := conform(x, y)
	if t != it {
		trap(Type)
	}
	xp := int32(x)
	yp := int32(y)
	switch av {
	case 0: //a%a
		return Ki(xp / yp)
	case 1: //a%v
		r = use1(y)
		rp := int32(r)
		e := rp + 4*nn(r)
		for rp < e {
			SetI32(rp, xp/I32(yp))
			xp += 4
			rp += 4
		}
		dx(y)
		return r
	case 2: //v%v
		r = use2(x, y)
		rp := int32(r)
		e := rp + 4*nn(r)
		for rp < e {
			SetI32(rp, I32(xp)/I32(yp))
			xp += 4
			yp += 4
			rp += 4
		}
		dx(x)
		dx(y)
		return r
	default: // v%a
		r = use1(x)
		rp := int32(r)
		e := rp + 4*nn(r)
		for rp < e {
			SetI32(rp, I32(xp)/yp)
			xp += 4
			rp += 4
		}
		dx(x)
		return r
	}
}

func Min(x, y K) (r K) {
	xt, yt := tp(x), tp(y)
	if xt&15 == bt && yt&15 == bt { // keep bool, no uptype to int
		x = uptype(x, ct)
		y = uptype(y, ct)
		r = Min(x, y)
		return K(tp(r)-1)<<59 | K(uint32(r))
	}
	if tp(y) < 16 {
		return nd(284, 7, y, x)
	}
	return nd(284, 7, x, y)
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
func mincC(x I8x16, yp, rp, e int32) {
	for rp < e {
		I8x16store(rp, x.Min_s(I8x16load(yp)))
		yp += 16
		rp += 16
		continue
	}
}
func miniI(x I32x4, yp, rp, e int32) {
	for rp < e {
		I32x4store(rp, x.Min_s(I32x4load(yp)))
		yp += 16
		rp += 16
		continue
	}
}
func minfF(x F64x2, yp, rp, e int32) {
	for rp < e {
		F64x2store(rp, x.Pmin(F64x2load(yp)))
		yp += 16
		rp += 16
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
		I8x16store(rp, I8x16load(xp).Min_s(I8x16load(yp)))
		xp += 16
		yp += 16
		rp += 16
		continue
	}
}
func minI(xp, yp, rp, e int32) {
	for rp < e {
		I32x4store(rp, I32x4load(xp).Min_s(I32x4load(yp)))
		xp += 16
		yp += 16
		rp += 16
		continue
	}
}
func minF(xp, yp, rp, e int32) {
	for rp < e {
		F64x2store(rp, F64x2load(xp).Pmin(F64x2load(yp)))
		xp += 16
		yp += 16
		rp += 16
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
	xt, yt := tp(x), tp(y)
	if xt&15 == bt && yt&15 == bt { // keep bool, no uptype to int
		x = uptype(x, ct)
		y = uptype(y, ct)
		r = Max(x, y)
		return K(tp(r)-1)<<59 | K(uint32(r))
	}
	if tp(y) < 16 {
		return nd(296, 8, y, x)
	}
	return nd(296, 8, x, y)
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
func maxcC(x I8x16, yp, rp, e int32) {
	for rp < e {
		I8x16store(rp, x.Max_s(I8x16load(yp)))
		yp += 16
		rp += 16
		continue
	}
}
func maxiI(x I32x4, yp, rp, e int32) {
	for rp < e {
		I32x4store(rp, x.Max_s(I32x4load(yp)))
		yp += 16
		rp += 16
		continue
	}
}
func maxfF(x F64x2, yp, rp, e int32) {
	for rp < e {
		F64x2store(rp, x.Pmax(F64x2load(yp)))
		yp += 16
		rp += 16
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
		I8x16store(rp, I8x16load(xp).Max_s(I8x16load(yp)))
		xp += 16
		yp += 16
		rp += 16
		continue
	}
}
func maxI(xp, yp, rp, e int32) {
	for rp < e {
		I32x4store(rp, I32x4load(xp).Max_s(I32x4load(yp)))
		xp += 16
		yp += 16
		rp += 16
		continue
	}
}
func maxF(xp, yp, rp, e int32) {
	for rp < e {
		F64x2store(rp, F64x2load(xp).Pmax(F64x2load(yp)))
		xp += 16
		yp += 16
		rp += 16
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

func Eql(x, y K) K                     { return nc(308, 11, x, y) }
func eqi(x, y int32) int32             { return I32B(x == y) }
func eqf(x, y float64) int32           { return I32B(isnan(x) && isnan(y) || x == y) }
func eqz(xr, xi, yr, yi float64) int32 { return eqf(xr, yr) & eqf(xi, yi) }
func eqcC(v, w I8x16, yp, rp, e int32) {
	for rp < e {
		I8x16store(rp, v.Eq(I8x16load(yp)).And(w))
		yp += 16
		rp += 16
		continue
	}
}
func eqiI(x int32, yp, rp, e int32) {
	for rp < e {
		SetI8(rp, eqi(x, I32(yp)))
		yp += 4
		rp++
		continue
	}
}
func eqfF(x float64, yp, rp, e int32) {
	for rp < e {
		SetI8(rp, eqf(x, F64(yp)))
		yp += 8
		rp++
		continue
	}
}
func eqzZ(re, im float64, yp, rp, e int32) {
	for rp < e {
		SetI8(rp, eqz(re, im, F64(yp), F64(yp+8)))
		yp += 16
		rp++
		continue
	}
}
func eqCc(xp int32, v, w I8x16, rp, e int32) {
	for rp < e {
		I8x16store(rp, I8x16load(xp).Eq(v).And(w))
		xp += 16
		rp += 16
		continue
	}
}
func eqIi(xp, y int32, rp, e int32) {
	for rp < e {
		SetI8(rp, eqi(I32(xp), y))
		xp += 4
		rp++
		continue
	}
}
func eqFf(xp int32, y float64, rp, e int32) {
	for rp < e {
		SetI8(rp, eqf(F64(xp), y))
		xp += 8
		rp++
		continue
	}
}
func eqZz(xp int32, re, im float64, rp, e int32) {
	for rp < e {
		SetI8(rp, eqz(F64(xp), F64(xp+8), re, im))
		xp += 16
		rp++
		continue
	}
}
func eqC(w I8x16, xp, yp, rp, e int32) {
	for rp < e {
		I8x16store(rp, I8x16load(xp).Eq(I8x16load(yp)).And(w))
		xp += 16
		yp += 16
		rp += 16
		continue
	}
}
func eqI(xp, yp, rp, e int32) {
	for rp < e {
		SetI8(rp, eqi(I32(xp), I32(yp)))
		xp += 4
		yp += 4
		rp++
		continue
	}
}
func eqF(xp, yp, rp, e int32) {
	for rp < e {
		SetI8(rp, eqf(F64(xp), F64(yp)))
		xp += 8
		yp += 8
		rp++
		continue
	}
}
func eqZ(xp, yp, rp, e int32) {
	for rp < e {
		SetI8(rp, eqz(F64(xp), F64(xp+8), F64(yp), F64(yp+8)))
		xp += 16
		yp += 16
		rp++
		continue
	}
}

func Les(x, y K) (r K) { // x<y   `file<c
	if tp(x) == st && tp(y) == Ct {
		return writefile(cs(x), y)
	}
	return nc(323, 9, x, y)
}
func lti(x, y int32) int32   { return I32B(x < y) }
func ltf(x, y float64) int32 { return I32B(x < y || x != x) }
func ltz(xr, xi, yr, yi float64) int32 {
	if eqf(xr, yr) != 0 {
		return ltf(xi, yi)
	}
	return ltf(xr, yr)
}
func ltcC(v, w I8x16, yp, rp, e int32) {
	for rp < e {
		I8x16store(rp, v.Lt_s(I8x16load(yp)).And(w))
		yp += 16
		rp += 16
		continue
	}
}
func ltiI(x int32, yp, rp, e int32) {
	for rp < e {
		SetI8(rp, lti(x, I32(yp)))
		yp += 4
		rp++
		continue
	}
}
func ltfF(x float64, yp, rp, e int32) {
	for rp < e {
		SetI8(rp, ltf(x, F64(yp)))
		yp += 8
		rp++
		continue
	}
}
func ltzZ(re, im float64, yp, rp, e int32) {
	for rp < e {
		if re == F64(yp) {
			SetI8(rp, ltf(im, F64(yp+8)))
		} else {
			SetI8(rp, ltf(re, F64(yp+8)))
		}
		yp += 16
		rp++
		continue
	}
}
func ltCc(xp int32, v, w I8x16, rp, e int32) {
	for rp < e {
		I8x16store(rp, I8x16load(xp).Lt_s(v).And(w))
		xp += 16
		rp++
		continue
	}
}
func ltIi(xp, y int32, rp, e int32) {
	for rp < e {
		SetI8(rp, lti(I32(xp), y))
		xp += 4
		rp++
		continue
	}
}
func ltFf(xp int32, y float64, rp, e int32) {
	for rp < e {
		SetI8(rp, ltf(F64(xp), y))
		xp += 8
		rp++
		continue
	}
}
func ltZz(xp int32, re, im float64, rp, e int32) {
	for rp < e {
		if F64(xp) == re {
			SetI8(rp, ltf(F64(xp+8), im))
		} else {
			SetI8(rp, ltf(F64(xp), re))
		}
		xp += 16
		rp++
		continue
	}
}
func ltC(w I8x16, xp, yp, rp, e int32) {
	for rp < e {
		I8x16store(rp, I8x16load(xp).Lt_s(I8x16load(yp)).And(w))
		xp += 16
		yp += 16
		rp += 16
		continue
	}
}
func ltI(xp, yp, rp, e int32) {
	for rp < e {
		SetI8(rp, lti(I32(xp), I32(yp)))
		xp += 4
		yp += 4
		rp++
		continue
	}
}
func ltF(xp, yp, rp, e int32) {
	for rp < e {
		SetI8(rp, ltf(F64(xp), F64(yp)))
		xp += 8
		yp += 8
		rp++
		continue
	}
}
func ltZ(xp, yp, rp, e int32) {
	for rp < e {
		SetI8(rp, ltz(F64(xp), F64(xp+8), F64(yp), F64(yp+8)))
		xp += 16
		yp += 16
		rp++
		continue
	}
}

func Mor(x, y K) (r K)       { return nc(338, 10, x, y) }
func gti(x, y int32) int32   { return I32B(x > y) }
func gtf(x, y float64) int32 { return I32B(x > y || y != y) }
func gtz(xr, xi, yr, yi float64) int32 {
	if eqf(xr, yr) != 0 {
		return gtf(xi, yi)
	}
	return gtf(xr, yr)
}
func gtcC(v, w I8x16, yp, rp, e int32) {
	for rp < e {
		I8x16store(rp, v.Gt_s(I8x16load(yp)).And(w))
		yp += 16
		rp += 16
		continue
	}
}
func gtiI(x int32, yp, rp, e int32) {
	for rp < e {
		SetI8(rp, gti(x, I32(yp)))
		yp += 4
		rp++
		continue
	}
}
func gtfF(x float64, yp, rp, e int32) {
	for rp < e {
		SetI8(rp, gtf(x, F64(yp)))
		yp += 8
		rp++
		continue
	}
}
func gtzZ(re, im float64, yp, rp, e int32) {
	for rp < e {
		if re == F64(yp) {
			SetI8(rp, gtf(im, F64(yp+8)))
		} else {
			SetI8(rp, gtf(re, F64(yp+8)))
		}
		yp += 16
		rp++
		continue
	}
}
func gtCc(xp int32, v, w I8x16, rp, e int32) {
	for rp < e {
		I8x16store(rp, I8x16load(xp).Gt_s(v).And(w))
		xp += 16
		rp++
		continue
	}
}
func gtIi(xp, y int32, rp, e int32) {
	for rp < e {
		SetI8(rp, gti(I32(xp), y))
		xp += 4
		rp++
		continue
	}
}
func gtFf(xp int32, y float64, rp, e int32) {
	for rp < e {
		SetI8(rp, gtf(F64(xp), y))
		xp += 8
		rp++
		continue
	}
}
func gtZz(xp int32, re, im float64, rp, e int32) {
	for rp < e {
		if F64(xp) == re {
			SetI8(rp, gtf(F64(xp+8), im))
		} else {
			SetI8(rp, gtf(F64(xp), re))
		}
		xp += 16
		rp++
		continue
	}
}
func gtC(w I8x16, xp, yp, rp, e int32) {
	for rp < e {
		I8x16store(rp, I8x16load(xp).Gt_s(I8x16load(yp)).And(w))
		xp += 16
		yp += 16
		rp += 16
		continue
	}
}
func gtI(xp, yp, rp, e int32) {
	for rp < e {
		SetI8(rp, gti(I32(xp), I32(yp)))
		xp += 4
		yp += 4
		rp++
		continue
	}
}
func gtF(xp, yp, rp, e int32) {
	for rp < e {
		SetI8(rp, gtf(F64(xp), F64(yp)))
		xp += 8
		yp += 8
		rp++
		continue
	}
}
func gtZ(xp, yp, rp, e int32) {
	for rp < e {
		SetI8(rp, gtz(F64(xp), F64(xp+8), F64(yp), F64(yp+8)))
		xp += 16
		yp += 16
		rp++
		continue
	}
}

func isnan(x float64) bool { return x != x }

func Angle(x K) (r K) { // angle x
	xt := tp(x)
	if xt > Zt {
		return Ech(35, l1(x))
	}
	if xt < zt {
		return Ki(0)
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
		r = ntake(n, Ki(0))
	}
	dx(x)
	return r
}
func Rot(x, y K) (r K) { // r angle deg
	if tp(x) > Zt {
		return Ech(35, l2(x, y))
	}
	x = uptype(x, zt)
	if y == 0 {
		return x
	}
	y = uptype(y, ft)
	if tp(y) != ft {
		panic(Type) // only atom
	}
	deg := F64(int32(y))
	dx(y)
	return Mul(Kz(cosin(deg)), x)
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
	if rt == bt && b2i != 0 {
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
			trap(Type)
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
