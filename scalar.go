package main

import (
	. "github.com/ktye/wg/module"
)

type f1i = func(int32) int32
type f1f = func(float64) float64
type f1z = func(float64, float64) K
type f2i = func(int32, int32) int32
type fi3 = func(int32, int32, int32)

func Neg(x K) K              { return nm(220, x) } //220
func negi(x int32) int32     { return -x }
func negf(x float64) float64 { return -x }
func negz(x, y float64) K    { return Kz(-x, -y) }

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
	return nm(223, x) //227
}
func absi(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}
func absf(x float64) float64 { return F64abs(x) }
func absZ(x K) K {
	n := nn(x)
	r := mk(Ft, n)
	rp := int32(r)
	xp := int32(x)
	for n > 0 {
		n--
		SetF64(rp, hypot(F64(xp), F64(xp+8)))
		xp += 16
		rp += 8
		continue
	}
	dx(x)
	return r
}

func Sqr(x K) K {
	if tp(x)&15 != ft {
		x = Add(Kf(0), x)
	}
	return nm(244, x) //300
}
func sqrf(x float64) float64 { return F64sqrt(x) }

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
		dxy(x, y)
		if yt == ft {
			return Kf(hypot(F64(xp), F64(yp)))
		} else if yt == zt {
			return Kf(hypot(F64(xp), hypot(F64(yp), F64(yp+8))))
		}
	}
	trap() //nyi
	return 0
}
func Img(x K) K { // imag x
	xt := tp(x)
	if xt > Zt {
		return Ech(33, l1(x))
	}
	if xt == Zt {
		xp := 8 + int32(x)
		n := nn(x)
		r := mk(Ft, n)
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
	xp := int32(x)
	if tp(x) == zt {
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

func Add(x, y K) K          { return nd(226, 2, x, y) } //234
func addi(x, y int32) int32 { return x + y }
func addf(xp, yp, rp int32) { SetF64(rp, F64(xp)+F64(yp)) }
func addz(xp, yp, rp int32) { SetF64(rp, F64(xp)+F64(yp)); SetF64(rp+8, F64(xp+8)+F64(yp+8)) }
func Sub(x, y K) K          { return nd(229, 3, x, y) } //245
func subi(x, y int32) int32 { return x - y }
func subf(xp, yp, rp int32) { SetF64(rp, F64(xp)-F64(yp)) }
func subz(xp, yp, rp int32) { SetF64(rp, F64(xp)-F64(yp)); SetF64(rp+8, F64(xp+8)-F64(yp+8)) }

func Mul(x, y K) K          { return nd(232, 4, x, y) } //256
func muli(x, y int32) int32 { return x * y }
func mulf(xp, yp, rp int32) { SetF64(rp, F64(xp)*F64(yp)) }
func mulz(xp, yp, rp int32) {
	xr, xi := F64(xp), F64(xp+8)
	yr, yi := F64(yp), F64(yp+8)
	SetF64(rp, xr*yr-xi*yi)
	SetF64(rp+8, xr*yi+xi*yr)
}

func Mod(x, y K) K          { return nd(244, 41, x, y) } //300
func modi(x, y int32) int32 { return x % y }
func idiv(x, y K, mod int32) K {
	if mod != 0 {
		return Mod(x, y)
	}
	return Div(x, y)
}
func Div(x, y K) K          { return nd(235, 5, x, y) } //267
func divi(x, y int32) int32 { return x / y }
func divf(xp, yp, rp int32) { SetF64(rp, F64(xp)/F64(yp)) }
func divz(xp, yp, rp int32) {
	xr, xi := F64(xp), F64(xp+8)
	yr, yi := F64(yp), F64(yp+8)
	r, d, e, f := 0.0, 0.0, 0.0, 0.0
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
	SetF64(rp, e)
	SetF64(rp+8, f)
}

func Min(x, y K) K { return nd(238, 6, x, y) } //278
func mini(x, y int32) int32 {
	if x < y {
		return x
	}
	return y
}
func minf(xp, yp, rp int32) { SetF64(rp, F64min(F64(xp), F64(yp))) }
func minz(xp, yp, rp int32) {
	if cmZ(xp, yp) > 0 {
		xp = yp
	}
	SetI64(rp, I64(xp))
	SetI64(rp+8, I64(xp+8))
}

func Max(x, y K) K { return nd(241, 7, x, y) } //289
func maxi(x, y int32) int32 {
	if x > y {
		return x
	} else {
		return y
	}
}
func maxf(xp, yp, rp int32) { SetF64(rp, F64max(F64(xp), F64(yp))) }
func maxz(xp, yp, rp int32) {
	if cmZ(xp, yp) < 0 {
		xp = yp
	}
	SetI64(rp, I64(xp))
	SetI64(rp+8, I64(xp+8))
}

// compare: 0(match) -1(x<y) 1(x>y)
func cmi(x, y int32) int32 { return I32B(x > y) - I32B(x < y) }
func cmC(x, y int32) int32 { x, y = I8(x), I8(y); return I32B(x > y) - I32B(x < y) }
func cmI(x, y int32) int32 { x, y = I32(x), I32(y); return I32B(x > y) - I32B(x < y) }
func cmF(x, y int32) int32 {
	a, b := I64(x), I64(y)
	if 2 == I32B(a < 0)+I32B(b < 0) {
		a, b = -a, -b
	}
	return I32B(a > b) - I32B(a < b)
}
func cmZ(x, y int32) int32 {
	r := cmF(x, y)
	if r == 0 {
		return cmF(x+8, y+8)
	} else {
		return r
	}
}

func Eql(x, y K) K { return nc(10, 0, x, y) } //308
func Les(x, y K) K { // x<y   `file<c
	if tp(x) == st && tp(y) == Ct {
		if int32(x) == 0 {
			write(rx(y))
			return y
		}
		return writefile(cs(x), y)
	}
	return nc(8, -1, x, y) //323
}
func Mor(x, y K) K { return nc(9, 1, x, y) } //338

func Ang(x K) K { // angle x
	var r K
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
func Rot(x, y K) K { // r@deg
	var r K
	if tp(x) > Zt {
		return Ech(35, l2(x, y))
	}
	x = uptype(x, zt)
	if y == 0 {
		return x
	}
	if tp(y)&15 > ft {
		trap() //type
	}
	y = uptype(y, ft)
	yp := int32(y)
	if tp(y) == ft {
		r = Kz(0, 0)
		cosin(F64(yp), int32(r))
	} else {
		yn := nn(y)
		r = mk(Zt, yn)
		rp := int32(r)
		for yn > 0 {
			yn--
			cosin(F64(yp), rp)
			yp += 8
			rp += 16
		}
	}
	dx(y)
	return Mul(r, x)
}
func Sin(x K) K { return nf(44, x, 0) } // sin x
func Cos(x K) K { return nf(45, x, 0) } // cos x
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
		trap() //type
	}
	dx(x)
	return F64(int32(x))
}
func nf(f int32, x, y K) K {
	var r K
	xt := tp(x)
	if xt >= Lt {
		if y == 0 {
			return Ech(K(f), l1(x))
		} else {
			return Ech(K(f-64), l2(y, x))
		}
	}
	if xt&15 < ft {
		x = uptype(x, ft)
		xt = tp(x)
	}
	xp := int32(x)
	xn := int32(1)
	if xt == ft {
		r = Kf(0)
	} else {
		xn = nn(x)
		r = mk(Ft, xn)
	}
	if xn > 0 {
		f += 233 - 60*I32B(f == 106) //exp,log,sin,cos,pow only pow uses y
		dr := int32(r) - xp
		e := xp + 8*xn
		for xp < e {
			Func[f].(fi3)(xp, int32(y), xp+dr)
			xp += 8
			continue
		}
	}
	dxy(x, y)
	return r
}
func nm(f int32, x K) K { //monadic
	var r K
	xt := tp(x)
	if xt > Lt {
		r = x0(x)
		return key(r, nm(f, r1(x)), xt)
	}
	xp := int32(x)
	if xt == Lt {
		n := nn(x)
		r = mk(Lt, n)
		rp := int32(r)
		for n > 0 {
			n--
			SetI64(rp, int64(nm(f, x0(K(xp)))))
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
			trap() //type
			return 0
		case 3:
			r = Kf(Func[1+f].(f1f)(F64(xp)))
			dx(x)
			return r
		case 4:
			r = Func[2+f].(f1z)(F64(xp), F64(xp+8))
			dx(x)
			return r
		default:
			trap() //type
			return 0
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
		for rp < e {
			SetI8(rp, Func[f].(f1i)(I8(xp)))
			xp++
			rp++
			continue
		}
	case 1:
		for rp < e {
			SetI32(rp, Func[f].(f1i)(I32(xp)))
			xp += 4
			rp += 4
			continue
		}
	case 2:
		trap() //type
	default: //F/Z (only called for neg)
		for rp < e {
			SetF64(rp, Func[1+f].(f1f)(F64(xp)))
			xp += 8
			rp += 8
			continue
		}
	}
	dx(x)
	return r
}
func nd(f, ff int32, x, y K) K { //dyadic
	var r K
	var n int32
	t := dtypes(x, y)
	if t > Lt {
		r = dkeys(x, y)
		return key(r, Func[64+ff].(f2)(dvals(x), dvals(y)), t)
	}
	if t == Lt {
		return Ech(K(ff), l2(x, y))
	}
	t = maxtype(x, y)
	x = uptype(x, t)
	y = uptype(y, t)
	av := conform(x, y)
	xp, yp := int32(x), int32(y)

	if av == 0 { //atom-atom
		switch t - 2 {
		case 0: // ct
			return Kc(Func[f].(f2i)(xp, yp))
		case 1: // it
			return Ki(Func[f].(f2i)(xp, yp))
		case 2: // st
			trap() //type
			return 0
		default: // ft zt
			r = mk(16+t, 1) //Kf(0.0)
			dxy(x, y)
			Func[f-4+int32(t)].(fi3)(xp, yp, int32(r))
			return Fst(r)
		}
	}

	ix := sz(t + 16)
	iy := ix
	if av == 1 { //av
		x = Enl(x)
		xp = int32(x)
		ix = 0
		n = nn(y)
		r = use1(y)
	} else if av == 2 { //va
		n = nn(x)
		y = Enl(y)
		yp = int32(y)
		iy = 0
		r = use1(x)
	} else {
		n = nn(x)
		r = use2(x, y)
	}
	if n == 0 {
		dxy(x, y)
		return r
	}

	rp := int32(r)
	e := ep(r)
	dz := int32(8) << I32B(t > ft)
	switch t - 2 {
	case 0: // ct
		for rp < e {
			SetI8(rp, Func[f].(f2i)(I8(xp), I8(yp)))
			xp += ix
			yp += iy
			rp++
			continue
		}
	case 1: // it
		for rp < e {
			SetI32(rp, Func[f].(f2i)(I32(xp), I32(yp)))
			xp += ix
			yp += iy
			rp += 4
			continue
		}
	case 2: // st
		trap() //type
	default: // ft zt
		for rp < e {
			Func[f-4+int32(t)].(fi3)(xp, yp, rp)
			xp += ix
			yp += iy
			rp += dz
			continue
		}
	}
	dxy(x, y)
	return r
}
func nc(ff, q int32, x, y K) K { //compare
	var r K
	var n int32
	t := dtypes(x, y)
	if t > Lt {
		r = dkeys(x, y)
		return key(r, nc(ff, q, dvals(x), dvals(y)), t)
	}
	if t == Lt {
		return Ech(K(ff), l2(x, y))
	}
	t = maxtype(x, y)
	x = uptype(x, t)
	y = uptype(y, t)
	av := conform(x, y)
	xp, yp := int32(x), int32(y)
	if av == 0 { // atom-atom
		dxy(x, y)
		// 11(derived), 12(proj), 13(lambda), 14(native)?
		return Ki(I32B(q == Func[245+t].(f2i)(xp, yp)))
	}
	ix := sz(t + 16)
	iy := ix
	if av == 1 { //av
		x = Enl(x)
		xp = int32(x)
		ix = 0
		n = nn(y)
	} else if av == 2 { //va
		n = nn(x)
		y = Enl(y)
		yp = int32(y)
		iy = 0
	} else {
		n = nn(x)
	}
	r = mk(It, n)
	if n == 0 {
		dxy(x, y)
		return r
	}
	rp := int32(r)
	e := ep(r)
	for rp < e {
		SetI32(rp, I32B(q == Func[250+t].(f2i)(xp, yp)))
		xp += ix
		yp += iy
		rp += 4
		continue
	}
	dxy(x, y)
	return r
}
func conform(x, y K) int32 { // 0:atom-atom 1:atom-vector, 2:vector-atom, 3:vector-vector
	r := 2*I32B(tp(x) > 16) + I32B(tp(y) > 16)
	if r == 3 {
		if nn(x) != nn(y) {
			trap() //length
		}
	}
	return r
}
func dtypes(x, y K) T {
	xt, yt := tp(x), tp(y)
	return T(maxi(int32(xt), int32(yt)))
}
func dkeys(x, y K) K {
	if tp(x) > Lt {
		return x0(x)
	}
	return x0(y)
}
func dvals(x K) K {
	if tp(x) > Lt {
		return r1(x)
	}
	return x
}
func maxtype(x, y K) T {
	xt, yt := tp(x)&15, tp(y)&15
	t := T(maxi(int32(xt), int32(yt)))
	if t == 0 {
		t = it
	}
	return t
}
func uptype(x K, dst T) K {
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
			f := float64(xp)
			if xt == ft {
				f = F64(xp)
				dx(x)
			}
			return Kz(f, 0)
		} else {
			trap() //type
			return 0
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
	r := mk(dst+16, xn)
	rp := int32(r)
	e := ep(r)
	if dst == it {
		for rp < e {
			SetI32(rp, I8(xp))
			xp++
			rp += 4
		}
	} else if dst == ft {
		for rp < e {
			SetF64(rp, float64(I32(xp)))
			xp += 4
			rp += 8
		}
	} else if dst == zt {
		for rp < e {
			SetF64(rp, F64(xp))
			SetF64(rp+8, 0.0)
			xp += 8
			rp += 16
		}
	} else {
		trap() //type
	}
	dx(x)
	return r
}
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
func use(x K) K {
	xt := tp(x)
	if xt < 16 || xt > Lt {
		trap() //type
	}
	if I32(int32(x)-4) == 1 {
		return x
	}
	nx := nn(x)
	r := mk(xt, nx)
	Memorycopy(int32(r), int32(x), sz(xt)*nx)
	if xt == Lt {
		rl(r)
	}
	dx(x)
	return r
}
