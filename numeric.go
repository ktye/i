package main

import (
	. "github.com/ktye/wg/module"
)

type f1i = func(int32) int32
type f1f = func(float64) float64
type f1z = func(float64, float64) (float64, float64)
type f2b = func(uint64, uint64) uint64
type f2i = func(int32, int32) int32
type f2iK = func(int32, K) K
type f2fK = func(float64, K) K
type f2zK = func(float64, float64, K) K
type f2f = func(float64, float64) float64
type f2z = func(float64, float64, float64, float64) (float64, float64)
type f2c = func(float64, float64) int32
type f2d = func(float64, float64, float64, float64) int32

func use2(x, y K) K {
	if I64(int32(y)-16) == 1 {
		return rx(y)
	}
	return use1(x)
}
func use1(x K) K {
	if I64(int32(x)-16) == 1 {
		return rx(x)
	}
	return mk(tp(x), nn(x))
}

func nm(f int32, x K) (r K) {
	xt := tp(x)
	xp := int32(x)
	if xt < 16 {
		switch xt - 1 {
		case 0:
			r = Ki(Func[f].(f1i)(xp))
		case 1:
			r = Kc(Func[f].(f1i)(xp))
		case 2:
			r = Ki(Func[1+f].(f1i)(xp))
		case 3:
			trap(Type)
		case 4:
			r = Kf(Func[2+f].(f1f)(F64(xp)))
			dx(x)
		case 5:
			r = Kz(Func[3+f].(f1z)(F64(xp), F64(xp)))
			dx(x)
		default:
			trap(Type)
		}
		return r
	}
	switch xt - 17 {
	case 0:
		r = Func[5+f].(f1)(uptype(x, it))
	case 1:
		r = Func[4+f].(f1)(x)
	case 2:
		r = Func[5+f].(f1)(x)
	case 3:
		trap(Type)
	case 4:
		r = Func[6+f].(f1)(x)
	case 5:
		r = Func[7+f].(f1)(x)
	default:
		trap(Type)
	}
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
func negC(x K) (r K) {
	r = use1(x)
	xp, rp := int32(x), int32(r)
	n := nn(r)
	for i := int32(0); i < n; i++ {
		SetI8(rp+i, absc(I8(xp+i)))
	}
	dx(x)
	return r
}
func negI(x K) (r K) {
	r = use1(x)
	n := vn(r)
	xp, rp := int32(x), int32(r)
	for i := int32(0); i < n; i++ {
		I32x4store(rp, I32x4neg(I32x4load(xp)))
		rp += 16
	}
	dx(x)
	return r
}
func negF(x K) (r K) {
	r = use1(x)
	n := vn(r)
	xp, rp := int32(x), int32(r)
	for i := int32(0); i < n; i++ {
		F64x2store(rp, F64x2neg(F64x2load(xp)))
		rp += 16
	}
	dx(x)
	return r
}
func negZ(x K) (r K) { return ZF(negF(x)) }

func Abs(x K) K {
	if tp(x) == zt {
		xp := int32(x)
		dx(x)
		return Kf(hypot(F64(xp), F64(xp+4)))
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
func hypot(p, q float64) float64 {
	//todo
	//switch {
	//case IsInf(p, 0) || IsInf(q, 0):
	//	return Inf(1)
	//case IsNaN(p) || IsNaN(q):
	//	return NaN()
	//}
	p, q = F64abs(p), F64abs(q)
	if p < q {
		t := p
		p = q
		q = t
	}
	if p == 0.0 {
		return 0.0
	}
	q = q / p
	return p * F64sqrt(1+q*q)
}
func absC(x K) (r K) {
	r = use1(x)
	xp, rp := int32(x), int32(r)
	n := nn(r)
	for i := int32(0); i < n; i++ {
		SetI8(rp+i, absc(I8(xp+i)))
	}
	dx(x)
	return r
}
func absI(x K) (r K) {
	r = use1(x)
	n := vn(r)
	xp, rp := int32(x), int32(r)
	for i := int32(0); i < n; i++ {
		I32x4store(rp, I32x4abs(I32x4load(xp)))
		xp += 16
		rp += 16
	}
	dx(x)
	return r
}
func absF(x K) (r K) {
	r = use1(x)
	n := vn(r)
	xp, rp := int32(x), int32(r)
	for i := int32(0); i < n; i++ {
		F64x2store(rp, F64x2abs(I32x4load(xp)))
		xp += 16
		rp += 16
	}
	dx(x)
	return r
}
func absZ(x K) (r K) {
	n := nn(r)
	r = mk(Ft, n)
	rp := int32(r)
	xp := int32(x)
	for i := int32(0); i < n; i++ {
		SetF64(rp, hypot(F64(xp), F64(xp+8)))
		xp += 16
		rp += 8
	}
	dx(x)
	return r
}

func Sqr(x K) (r K) {
	xt := tp(x)
	xp := int32(x)
	if xt == ft {
		r = Kf(F64sqrt(F64(xp)))
	} else if xt == Ft {
		r = sqrF(x)
	} else {
		trap(Type)
	}
	return r
}
func sqrF(x K) (r K) {
	r = use1(x)
	n := vn(r)
	xp, rp := int32(x), int32(r)
	for i := int32(0); i < n; i++ {
		F64x2store(rp, F64x2sqrt(F64x2load(xp)))
		rp += 16
	}
	dx(x)
	return r
}

func nd(f int32, x, y K) (r K) {
	x, y = uptypes(x, y)
	xp, yp := int32(x), int32(y)
	av, t := conform(x, y)
	if av == 0 { // atom-atom
		switch t - 2 {
		case 0: // ct
			r = Kc(Func[f].(f2i)(xp, yp))
		case 1: // it
			r = Ki(Func[1+f].(f2i)(xp, yp))
		case 2: // st
			trap(Type)
		case 3:
			dx(x)
			r = Kf(Func[2+f].(f2f)(F64(xp), F64(yp)))
		case 4:
			dx(x)
			r = Kz(Func[3+f].(f2z)(F64(xp), F64(xp+8), F64(yp), F64(yp+8)))
		default:
			trap(Type)
		}
		return r
	} else if av == 1 { // atom-vector
		switch t - 2 {
		case 0: // ct
			r = Func[4+f].(f2iK)(xp, y)
		case 1: // it
			r = Func[5+f].(f2iK)(xp, y)
		case 2: // st
			trap(Type)
		case 3: // ft
			dx(x)
			r = Func[6+f].(f2fK)(F64(xp), y)
		case 4: // zt
			dx(x)
			r = Func[7+f].(f2zK)(F64(xp), F64(xp+8), y)
		default:
			trap(Type)
		}
		return r
	} else { // vector-vector
		switch t - 2 {
		case 0: // ct
			r = Func[8+f].(f2)(x, y)
		case 1: // it
			r = Func[9+f].(f2)(x, y)
		case 2: // st
			trap(Type)
		case 3: // ft
			r = Func[10+f].(f2)(x, y)
		case 4: // zt
			r = Func[11+f].(f2)(x, y)
		default:
			trap(Type)
		}
	}
	return r
}
func conform(x, y K) (av int32, r T) { // 0:atom-atom 1:atom-vector, 2:vector-vector
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
		trap(Length)
	}
	if nn(y) != xn {
		trap(Length)
	}
	return 2, xt - 16
}
func Add(x, y K) (r K) {
	if tp(y) < 16 {
		return nd(236, y, x)
	}
	return nd(236, x, y)
}
func addi(x, y int32) int32                          { return x + y }
func addf(x, y float64) float64                      { return x + y }
func addz(xr, xi, yr, yi float64) (float64, float64) { return xr + yr, xi + yi }
func addcC(x int32, y K) (r K) {
	n := vn(y)
	v := I8x16splat(x)
	r = use1(y)
	yp, rp := int32(y), int32(r)
	for i := int32(0); i < n; i++ {
		I8x16store(rp, I8x16add(v, I8x16load(yp)))
		yp += 16
		rp += 16
	}
	dx(y)
	return r
}
func addiI(x int32, y K) (r K) {
	n := vn(y)
	v := I32x4splat(x)
	r = use1(y)
	yp, rp := int32(y), int32(r)
	for i := int32(0); i < n; i++ {
		I32x4store(rp, I32x4add(v, I32x4load(yp)))
		yp += 16
		rp += 16
	}
	dx(y)
	return r
}
func addfF(x float64, y K) (r K) {
	n := vn(y)
	v := F64x2splat(x)
	r = use1(y)
	yp, rp := int32(y), int32(r)
	for i := int32(0); i < n; i++ {
		F64x2store(rp, F64x2add(v, F64x2load(yp)))
		yp += 16
		rp += 16
	}
	dx(y)
	return r
}
func addzZ(re, im float64, y K) (r K) {
	n := vn(y)
	v := F64x2replace_lane1(F64x2splat(re), im)
	r = use1(y)
	yp, rp := int32(y), int32(r)
	for i := int32(0); i < n; i++ {
		F64x2store(rp, F64x2add(v, F64x2load(yp)))
		yp += 16
		rp += 16
	}
	dx(y)
	return r
}
func addC(x, y K) (r K) {
	r = use2(x, y)
	rp, xp, yp := int32(r), int32(x), int32(y)
	n := vn(x)
	for i := int32(0); i < n; i++ {
		I8x16store(rp, I8x16add(I8x16load(xp), I8x16load(yp)))
		rp += 16
		xp += 16
		yp += 16
	}
	dx(x)
	dx(y)
	return r
}
func addI(x, y K) (r K) {
	r = use2(x, y)
	rp, xp, yp := int32(r), int32(x), int32(y)
	n := vn(x)
	for i := int32(0); i < n; i++ {
		I32x4store(rp, I32x4add(I32x4load(xp), I32x4load(yp)))
		rp += 16
		xp += 16
		yp += 16
	}
	dx(x)
	dx(y)
	return r
}
func addF(x, y K) (r K) {
	r = use2(x, y)
	rp, xp, yp := int32(r), int32(x), int32(y)
	n := vn(x)
	for i := int32(0); i < n; i++ {
		F64x2store(rp, F64x2add(F64x2load(xp), F64x2load(yp)))
		rp += 16
		xp += 16
		yp += 16
	}
	dx(x)
	dx(y)
	return r
}
func addZ(x, y K) (r K) { return ZF(addF(x, y)) }

func Sub(x, y K) (r K) {
	if tp(y) < 16 {
		return nd(236, Neg(y), x)
	}
	return nd(239, x, y)
}

func Mul(x, y K) K                                   { trap(Nyi); return 0 }
func mulz(xr, xi, yr, yi float64) (float64, float64) { return xr*yr - xi*yi, xr*yi + xi*yr }
func mulzZ(re, im float64, y K) (r K) {
	n := nn(y)
	r = use1(y)
	yp, rp := int32(y), int32(r)
	for i := int32(0); i < n; i++ {
		xx, yy := mulz(re, im, F64(yp), F64(yp+8))
		SetF64(rp, xx)
		SetF64(rp, yy)
		yp += 16
		rp += 16
	}
	dx(y)
	return r
}

func Div(x, y K) K { trap(Nyi); return 0 }
func Min(x, y K) K { trap(Nyi); return 0 }
func Max(x, y K) K { trap(Nyi); return 0 }
func Les(x, y K) K { trap(Nyi); return 0 }
func Mor(x, y K) K { trap(Nyi); return 0 }
func Eql(x, y K) K { trap(Nyi); return 0 }

/*
func nc(x, y K, f int32) (r K) {
	x, y = uptypes(x, y)
	n, av, t := conform(x, y)
	switch t - 2 {
	case 0: // ct
		r = ndc(x, y, f, n, av)
		if av == 3 {
			r = Kb(int32(r))
		} else {
			r = K(int32(r)) | K(Bt)<<59
		}
	case 1: // it
		r = nci(x, y, f, n, av)
	case 2: // st
		trap(Type)
	case 3: // ft
		r = ncf(x, y, 1+f, n, av)
	case 4: // zt
		r = ncf(x, y, 2+f, n, av)
	default:
		trap(Type)
	}
	return r
}
func ndc(x, y K, f, n, av int32) (r K) {
	xp := int32(x)
	yp := int32(y)
	if av == 3 {
		return Kc(Func[f].(f2i)(xp, yp))
	}
	switch av {
	case 0:
		r = use2(x, y)
		rp := int32(r)
		for i := int32(0); i < n; i++ {
			SetI8(rp, Func[f].(f2i)(I8(xp), I8(yp)))
			xp++
			yp++
			rp++
		}
	case 1:
		r = use1(y)
		rp := int32(r)
		for i := int32(0); i < n; i++ {
			SetI8(rp, Func[f].(f2i)(xp, I8(yp)))
			yp++
			rp++
		}
	case 2:
		r = use1(x)
		rp := int32(r)
		for i := int32(0); i < n; i++ {
			SetI8(rp, Func[f].(f2i)(I8(xp), yp))
			xp++
			rp++
		}
	}
	dx(x)
	dx(y)
	return r
}
func ndi(x, y K, f, n, av int32) (r K) {
	xp := int32(x)
	yp := int32(y)
	if av == 3 {
		return Ki(Func[f].(f2i)(xp, yp))
	}
	switch av {
	case 0:
		r = use2(x, y)
		rp := int32(r)
		for i := int32(0); i < n; i++ {
			SetI32(rp, Func[f].(f2i)(I32(xp), I32(yp)))
			xp += 4
			yp += 4
			rp += 4
		}
	case 1:
		r = use1(y)
		rp := int32(r)
		for i := int32(0); i < n; i++ {
			SetI32(rp, Func[f].(f2i)(xp, I32(yp)))
			yp += 4
			rp += 4
		}
	case 2:
		r = use1(x)
		rp := int32(r)
		for i := int32(0); i < n; i++ {
			SetI32(rp, Func[f].(f2i)(I32(xp), yp))
			xp += 4
			rp += 4
		}
	}
	dx(x)
	dx(y)
	return r
}
func nci(x, y K, f, n, av int32) (r K) {
	xp := int32(x)
	yp := int32(y)
	if av == 3 {
		return Kb(Func[f].(f2i)(xp, yp))
	}
	r = mk(Bt, n)
	rp := int32(r)
	switch av {
	case 0:
		for i := int32(0); i < n; i++ {
			SetI8(rp, Func[f].(f2i)(I32(xp), I32(yp)))
			xp += 4
			yp += 4
			rp++
		}
	case 1:
		for i := int32(0); i < n; i++ {
			SetI8(rp, Func[f].(f2i)(xp, I32(yp)))
			yp += 4
			rp++
		}
	case 2:
		for i := int32(0); i < n; i++ {
			SetI8(rp, Func[f].(f2i)(I32(xp), yp))
			xp += 4
			rp++
		}
	}
	dx(x)
	dx(y)
	return r
}
func ndf(x, y K, f, n, av int32) (r K) {
	xp, yp := int32(x), int32(y)
	xd, yd := int32(0), int32(0)
	switch av {
	case 0:
		xd, yd = 8, 8
		r = use2(x, y)
	case 1:
		yd = 8
		r = use1(y)
	case 2:
		xd = 8
		r = use1(x)
	case 3:
		r = Kf(0)
	}
	rp := int32(r)
	for i := int32(0); i < n; i++ {
		SetF64(rp, Func[f].(f2f)(F64(xp), F64(yp)))
		xp += xd
		yp += yd
		rp += 8
	}
	dx(x)
	dx(y)
	return r
}
func ndz(x, y K, f, n, av int32) (r K) {
	xp, yp := int32(x), int32(y)
	xd, yd := int32(0), int32(0)
	switch av {
	case 0:
		xd, yd = 16, 16
		r = use2(x, y)
	case 1:
		yd = 16
		r = use1(y)
	case 2:
		xd = 16
		r = use1(x)
	case 3:
		r = Kz(0, 0)
	}
	rp := int32(r)
	for i := int32(0); i < n; i++ {
		re, im := Func[f].(f2z)(F64(xp), F64(xp+8), F64(yp), F64(yp+8))
		SetF64(rp, re)
		SetF64(rp+8, im)
		xp += xd
		yp += yd
		rp += 16
	}
	dx(x)
	dx(y)
	return r
}
func ncf(x, y K, f, n, av int32) (r K) {
	t := tp(x) & 15
	xp, yp := int32(x), int32(y)
	xd, yd := int32(0), int32(0)
	switch av {
	case 0:
		xd, yd = 8, 8
	case 1:
		yd = 8
	case 2:
		xd = 8
	}
	r = mk(Bt, n)
	rp := int32(r)
	if t == zt {
		xd *= 2
		yd *= 2
		for i := int32(0); i < n; i++ {
			SetI8(rp+i, Func[f].(f2d)(F64(xp), F64(xp+8), F64(yp), F64(yp+8)))
			xp += xd
			yp += yd
		}
	} else {
		for i := int32(0); i < n; i++ {
			SetI8(rp+i, Func[f].(f2c)(F64(xp), F64(yp)))
			xp += xd
			yp += yd
		}
	}
	dx(x)
	dx(y)
	if av == 3 {
		dx(r)
		return Kb(I8(rp))
	}
	return r
}



// 0:vector-vector, 1:atom-vector, 2:vector-atom, 3: atom-atom
func conform(x, y K) (n int32, av int32, r T) {
	xt, yt := tp(x), tp(y)
	if xt < 16 {
		if yt < 16 {
			return 1, 3, xt
		} else {
			return nn(y), 1, xt
		}
	}
	xn := nn(x)
	if yt < 16 {
		return xn, 2, yt
	}
	if nn(y) != xn {
		trap(Length)
	}
	return xn, 0, xt - 16
}

func Add(x, y K) K                                   { return nd(x, y, 220) }
func addi(x, y int32) int32                          { return x + y }
func addf(x, y float64) float64                      { return x + y }
func addz(xr, xi, yr, yi float64) (float64, float64) { return xr + yr, xi + yi }

func Sub(x, y K) K                                   { return nd(x, y, 223) }
func subi(x, y int32) int32                          { return x - y }
func subf(x, y float64) float64                      { return x - y }
func subz(xr, xi, yr, yi float64) (float64, float64) { return xr - yr, xi + yi }

func Mul(x, y K) K                                   { return nd(x, y, 226) }
func muli(x, y int32) int32                          { return x * y }
func mulf(x, y float64) float64                      { return x * y }
func mulz(xr, xi, yr, yi float64) (float64, float64) { return xr*yr - xi*yi, xr*yi + xi*yr }

func Div(x, y K) K              { return nd(x, y, 229) }
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

func Min(x, y K) K { return nd(x, y, 235) }
func mini(x, y int32) int32 {
	if x < y {
		return x
	}
	return y
}
func minf(x, y float64) float64 { return F64min(x, y) }
func minz(xr, xi, yr, yi float64) (float64, float64) {
	if gtz(xr, xi, yr, yi) != 0 {
		return yr, yi
	}
	return xr, xi
}

func Max(x, y K) K { return nd(x, y, 238) }
func maxi(x, y int32) int32 {
	if x > y {
		return x
	}
	return y
}
func maxf(x, y float64) float64 { return F64max(x, y) }
func maxz(xr, xi, yr, yi float64) (float64, float64) {
	if ltz(xr, xi, yr, yi) != 0 {
		return yr, yi
	}
	return xr, xi
}

func Les(x, y K) K           { return nc(x, y, 241) }
func lti(x, y int32) int32   { return ib(x < y) }
func ltf(x, y float64) int32 { return ib(x < y) }
func ltz(xr, xi, yr, yi float64) int32 {
	if eqf(xr, yr) != 0 {
		return ltf(xi, yi)
	}
	return ltf(xr, yr)
}

func Mor(x, y K) K           { return nc(x, y, 244) }
func gti(x, y int32) int32   { return ib(x > y) }
func gtf(x, y float64) int32 { return ib(x > y) }
func gtz(xr, xi, yr, yi float64) int32 {
	if eqf(xr, yr) != 0 {
		return gtf(xi, yi)
	}
	return gtf(xr, yr)
}

func Eql(x, y K) K                     { return nc(x, y, 247) }
func eqi(x, y int32) int32             { return ib(x == y) }
func eqf(x, y float64) int32           { return ib(isnan(x) && isnan(y) || x == y) }
func eqz(xr, xi, yr, yi float64) int32 { return eqf(xr, yr) & eqf(xi, yi) }

*/

func mini(x, y int32) int32 {
	if x < y {
		return x
	}
	return y
}
func maxi(x, y int32) int32 {
	if x > y {
		return x
	}
	return y
}

func isnan(x float64) bool { return x != x }

func Rot(x, y K) (r K) { // r angle deg
	x = uptype(x, zt)
	y = uptype(y, ft)
	if tp(y) != ft {
		panic(Type) // only atom
	}
	deg := F64(int32(y))
	dx(y)
	im, re := sincos(deg)
	if tp(x) == zt {
		xp := int32(x)
		dx(x)
		return Kz(mulz(re, im, F64(xp), F64(xp+8)))
	}
	return mulzZ(re, im, x)
}
func sincos(deg float64) (s float64, c float64) {
	if deg == 0 {
		c = 1.0
	} else if deg == 90 {
		s = 1.0
	} else if deg == 180 {
		c = -1.0
	} else if deg == 270 {
		s = -1.0
	} else {
		trap(Nyi)
	}
	return s, c
}

func uptypes(x, y K) (K, K) {
	xt, yt := tp(x)&15, tp(y)&15
	rt := T(maxi(int32(xt), int32(yt)))
	if rt == bt {
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
