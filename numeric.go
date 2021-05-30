package main

import (
	. "github.com/ktye/wg/module"
)

func Neg(x K) (r K) {
	if tp(x) == It {
		nx := nn(x)
		r = mk(It, nx)
		rp := int32(r)
		xp := int32(x)
		for i := int32(0); i < nx; i++ {
			SetI32(rp, -I32(xp))
			xp += 4
			rp += 4
		}
		dx(x)
		return r
	}
	return Ki(-int32(x))
}

func Sqr(x K) K { trap(Nyi); return x }

func nd(x, y K, f int32) (r K) {
	x, y = uptypes(x, y)
	n, av, t := conform(x, y)
	switch t - 2 {
	case 0: // ct
		r = ndc(x, y, f, n, av)
	case 1: // it
		r = ndi(x, y, f, n, av)
	case 2: // st
		trap(Type)
	case 3: // ft
		r = ndf(x, y, 1+f, n, av)
	case 4: // zt
		r = ndz(x, y, 2+f, n, av)
	default:
		trap(Type)
	}
	return r
}
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

type f2b = func(uint64, uint64) uint64
type f2i = func(int32, int32) int32
type f2f = func(float64, float64) float64
type f2z = func(float64, float64, float64, float64) (float64, float64)
type f2c = func(float64, float64) int32
type f2d = func(float64, float64, float64, float64) int32

func use2(x, y K) K {
	if I64(int32(y)-8) == 1 {
		return rx(y)
	}
	return use1(x)
}
func use1(x K) K {
	if I64(int32(x)-8) == 1 {
		return rx(x)
	}
	return mk(tp(x), nn(x))
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
func isnan(x float64) bool             { return x != x }

func Rot(x, y K) (r K) { return nd(uptype(x, zt), y, 250) } // r angle deg
func rot(xr, xi, yr, yi float64) (zr float64, zi float64) {
	cs, sn := 0.0, 0.0
	if yr == 0.0 {
		cs = 1.0
	} else if yr == 90.0 {
		sn = 1.0
	} else if yr == 180.0 {
		cs = -1.0
	} else if yr == 270.0 {
		sn = -1.0
	} else {
		trap(Nyi)
	}
	return xr*cs - xi*sn, xr*sn + xi*cs
}

func abz(p, q float64) float64 {
	p, q = F64abs(p), F64abs(q)
	if p < q {
		p, q = q, p
	}
	if p == 0 {
		return 0
	}
	q = q / p
	return p * F64sqrt(1+q*q)
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
