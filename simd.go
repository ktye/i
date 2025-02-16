//go:build simd4||simd5

package main

import (
	. "github.com/ktye/wg/module"
)

func init() {
	Functions(280, negI, negF, negF)
	Functions(283, absI, absF, nyi)

	Functions(286, ltC, eqC, gtC, ltI, eqI, gtI, ltI, eqI, gtI)
	Functions(295, ltcC, eqcC, gtcC, ltiI, eqiI, gtiI, ltiI, eqiI, gtiI)
	Functions(305, sqrF)

	Functions(306, addI, subI, mulI, nyi, minI, maxI)
	Functions(312, addiI, subiI, muliI, nyi, miniI, maxiI)
	Functions(318, addF, subF, mulF, nyi, minF, maxF)
	Functions(324, addfF, subfF, nyi, nyi, minfF, maxfF)
}

type fii = func(int32, int32)
type fi4 = func(int32, int32, int32, int32)

func seq(n int32) K {
	i := Iota()
	a := int32(vl >> 2)
	n = maxi(n, 0)
	r := mk(It, n)
	rp := int32(r)
	e := rp + ev(4*n)
	for rp < e {
		VIstore(rp, i)
		i = i.Add(VIsplat(a))
		rp += vl
	}
	return r
}
func minis(x, y, e int32) int32 {
       	if e - y > 256 {
		a := VIsplat(x)
		u := eu(e)
		for y < u {
			a = a.Min_s(VIload(y))
			y += vl
			continue
		}
		x = a.Hmin_s()
	}
	for y < e {
                x = mini(x, I32(y))
                y += 4
        }
        return x
}
func maxis(x, y, e int32) int32 {
       	if e - y > 256 {
		a := VIsplat(x)
		u := eu(e)
		for y < u {
			a = a.Max_s(VIload(y))
			y += vl
			continue
		}
		x = a.Hmax_s()
	}
	for y < e {
                x = maxi(x, I32(y))
                y += 4
        }
        return x
}
func minfs(y, e int32) float64 {
        f := inf
       	if e - y > 256 {
		a := VFsplat(f)
		u := eu(e)
		for y < u {
			a = a.Pmin(VFload(y))
			y += vl
			continue
		}
		f = a.Hmin()
	}
        for y < e {
                f = F64min(f, F64(y))
                y += 8
        }
        return f
}
func maxfs(y, e int32) float64 {
        f := -inf
       	if e - y > 256 {
		a := VFsplat(f)
		u := eu(e)
		for y < u {
			a = a.Pmax(VFload(y))
			y += vl
			continue
		}
		f = a.Hmax()
	}
        for y < e {
                f = F64max(f, F64(y))
                y += 8
        }
        return f
}
func sumi(xp, e int32) int32 {
        r := int32(0)
	if e - xp > 256 {
		a := VIsplat(0)
		u := eu(e)
		for xp < u {
			a = a.Add(VIload(xp))
			xp += vl
			continue
		}
		r = a.Hsum()
	}
        for xp < e {
                r += I32(xp)
                xp += 4
        }
        return r
}
func sumf(xp, e int32) float64 {
        r := 0.0
	if e - xp > 256 {
		a := VFsplat(0.0)
		u := eu(e)
		for xp < u {
			a = a.Add(VFload(xp))
			xp += vl
			continue
		}
		r = a.Hsum()
	}
        for xp < e {
                r += F64(xp)
                xp += 8
        }
        return r
}
func negI(xp, e int32) {
	for xp < e {
		VIstore(xp, VIload(xp).Neg())
		xp += vl
		continue
	}
}
func negF(xp, e int32) {
	for xp < e {
		VFstore(xp, VFload(xp).Neg())
		xp += vl
		continue
	}
}
func absI(xp, e int32) {
	for xp < e {
		VIstore(xp, VIload(xp).Abs())
		xp += vl
		continue
	}
}
func absF(xp, e int32) {
	for xp < e {
		VFstore(xp, VFload(xp).Abs())
		xp += vl
		continue
	}
}
func sqrF(xp, e int32) {
	for xp < e {
		VFstore(xp, VFload(xp).Sqrt())
		xp += vl
		continue
	}
}
func fscale(x, y, e int32) {
	f := VFsplat(F64(x))
	for y < e {
		VFstore(y, f.Mul(VFload(y)))
		y += vl
	}
}
func scale(x, y K) K {
	xp := int32(x)
	r := use(y)
	rp := int32(r)
	e := ev(ep(r))
	if tp(y) == Zt && F64(xp + 8) != 0 {
		for rp < e {
			mulz(xp, rp, rp)
			rp += 16
		}
	} else {
		fscale(xp, rp, e)
	}
	dx(x)
	return r
}
func addI(x, y, r, e int32) {
	for r < e {
		VIstore(r, VIload(x).Add(VIload(y)))
		x += vl
		y += vl
		r += vl
		continue
	}
}
func addF(x, y, r, e int32) {
	for r < e {
		VFstore(r, VFload(x).Add(VFload(y)))
		x += vl
		y += vl
		r += vl
		continue
	}
}
func subI(x, y, r, e int32) {
	for r < e {
		VIstore(r, VIload(x).Sub(VIload(y)))
		x += vl
		y += vl
		r += vl
		continue
	}
}
func subF(x, y, r, e int32) {
	for r < e {
		VFstore(r, VFload(x).Sub(VFload(y)))
		x += vl
		y += vl
		r += vl
		continue
	}
}
func mulI(x, y, r, e int32) {
	for r < e {
		VIstore(r, VIload(x).Mul(VIload(y)))
		x += vl
		y += vl
		r += vl
		continue
	}
}
func mulF(x, y, r, e int32) {
	for r < e {
		VFstore(r, VFload(x).Mul(VFload(y)))
		x += vl
		y += vl
		r += vl
		continue
	}
}
func minI(x, y, r, e int32) {
	for r < e {
		VIstore(r, VIload(x).Min_s(VIload(y)))
		x += vl
		y += vl
		r += vl
		continue
	}
}
func minF(x, y, r, e int32) {
	for r < e {
		VFstore(r, VFload(x).Pmin(VFload(y)))
		x += vl
		y += vl
		r += vl
		continue
	}
}
func maxI(x, y, r, e int32) {
	for r < e {
		VIstore(r, VIload(x).Max_s(VIload(y)))
		x += vl
		y += vl
		r += vl
		continue
	}
}
func maxF(x, y, r, e int32) {
	for r < e {
		VFstore(r, VFload(x).Pmax(VFload(y)))
		x += vl
		y += vl
		r += vl
		continue
	}
}
func addiI(x, y, r, e int32) {
	i := VIsplat(I32(x))
	for r < e {
		VIstore(r, VIload(y).Add(i))
		y += vl
		r += vl
		continue
	}
}
func addfF(x, y, r, e int32) {
	i := VFsplat(F64(x))
	for r < e {
		VFstore(r, VFload(y).Add(i))
		y += vl
		r += vl
		continue
	}
}
func subiI(x, y, r, e int32) {
	i := VIsplat(I32(x))
	for r < e {
		VIstore(r, i.Sub(VIload(y)))
		y += vl
		r += vl
		continue
	}
}
func subfF(x, y, r, e int32) {
	i := VFsplat(F64(x))
	for r < e {
		VFstore(r, i.Sub(VFload(y)))
		y += vl
		r += vl
		continue
	}
}
func muliI(x, y, r, e int32) {
	i := VIsplat(I32(x))
	for r < e {
		VIstore(r, VIload(y).Mul(i))
		y += vl
		r += vl
		continue
	}
}
func miniI(x, y, r, e int32) {
	i := VIsplat(I32(x))
	for r < e {
		VIstore(r, VIload(y).Min_s(i))
		y += vl
		r += vl
		continue
	}
}
func minfF(x, y, r, e int32) {
	i := VFsplat(F64(x))
	for r < e {
		VFstore(r, VFload(y).Pmin(i))
		y += vl
		r += vl
		continue
	}
}
func maxiI(x, y, r, e int32) {
	i := VIsplat(I32(x))
	for r < e {
		VIstore(r, VIload(y).Max_s(i))
		y += vl
		r += vl
		continue
	}
}
func maxfF(x, y, r, e int32) {
	i := VFsplat(F64(x))
	for r < e {
		VFstore(r, VFload(y).Pmax(i))
		y += vl
		r += vl
		continue
	}
}
func ltC(x, y, r, e int32) {
	i := VI1()
	for r < e {
		VIstore(r, i.And(VIloadB(x).Lt_s(VIloadB(y))))
		x += vl >> 2
		y += vl >> 2
		r += vl
		continue
	}
}
func ltcC(x, y, r, e int32) {
	c := VIsplat(I8(x))
	i := VI1()
	for r < e {
		VIstore(r, i.And(c.Lt_s(VIloadB(y))))
		y += vl >> 2
		r += vl
		continue
	}
}
func eqC(x, y, r, e int32) {
	i := VI1()
	for r < e {
		VIstore(r, i.And(VIloadB(x).Eq(VIloadB(y))))
		x += vl >> 2
		y += vl >> 2
		r += vl
		continue
	}
}
func eqcC(x, y, r, e int32) {
	c := VIsplat(I8(x))
	i := VI1()
	for r < e {
		VIstore(r, i.And(c.Eq(VIloadB(y))))
		y += vl >> 2
		r += vl
		continue
	}
}
func gtC(x, y, r, e int32) {
	i := VI1()
	for r < e {
		VIstore(r, i.And(VIloadB(x).Gt_s(VIloadB(y))))
		x += vl >> 2
		y += vl >> 2
		r += vl
		continue
	}
}
func gtcC(x, y, r, e int32) {
	c := VIsplat(I8(x))
	i := VI1()
	for r < e {
		VIstore(r, i.And(c.Gt_s(VIloadB(y))))
		y += vl >> 2
		r += vl
		continue
	}
}
func ltI(x, y, r, e int32) {
	i := VI1()
	for r < e {
		VIstore(r, i.And(VIload(x).Lt_s(VIload(y))))
		x += vl
		y += vl
		r += vl
		continue
	}
}
func ltiI(x, y, r, e int32) {
	c := VIsplat(I32(x))
	i := VI1()
	for r < e {
		VIstore(r, i.And(c.Lt_s(VIload(y))))
		y += vl
		r += vl
		continue
	}
}
func eqI(x, y, r, e int32) {
	i := VI1()
	for r < e {
		VIstore(r, i.And(VIload(x).Eq(VIload(y))))
		x += vl
		y += vl
		r += vl
		continue
	}
}
func eqiI(x, y, r, e int32) {
	c := VIsplat(I32(x))
	i := VI1()
	for r < e {
		VIstore(r, i.And(c.Eq(VIload(y))))
		y += vl
		r += vl
		continue
	}
}
func gtI(x, y, r, e int32) {
	i := VI1()
	for r < e {
		VIstore(r, i.And(VIload(x).Gt_s(VIload(y))))
		x += vl
		y += vl
		r += vl
		continue
	}
}
func gtiI(x, y, r, e int32) {
	c := VIsplat(I32(x))
	i := VI1()
	for r < e {
		VIstore(r, i.And(c.Gt_s(VIload(y))))
		y += vl
		r += vl
		continue
	}
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
	x = use(x)
	xp = int32(x)
	e := ep(x)
	if e == xp {
		return x
	}
	switch xt - 18 {
	case 0:
		for xp < e {
			SetI8(xp, Func[f].(f1i)(I8(xp)))
			xp++
			continue
		}
	case 1:
		Func[f+60].(fii)(xp, ev(e))
	case 2:
		trap() //type
	default: //F/Z (only called for neg)
		Func[f+61].(fii)(xp, ev(e))
	}
	return x
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
		if t > st && f == 232 {
			return scale(x, y)
		}
		xp = int32(x)
		ix = 0
		n = nn(y)
		r = use1(y)
	} else if av == 2 { //va
		n = nn(x)
		y = Enl(y)
		yp = int32(y)
		if f < 238 { // +*&|
			xp = yp
			yp = int32(x)
			ix = 0
			av = 1
		} else if f == 241 && t > st {
			return scale(Div(Kf(1.0), y), x)
		} else {
			iy = 0
		}
		r = use1(x)
	} else {
		n = nn(x)
		if I32(int32(y)-4) == 1 {
			r = rx(y)
		} else {
			r = use1(x)
		}
		if t == 6 && ff < 4 { // Z+Z Z-Z
			t--
		}
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
		if f < 241 {
			if av == 2 && ff == 3 { // v-a
				SetI32(yp, -I32(yp))
				addiI(yp, xp, rp, ev(e))
			} else {
				if av != 3 {
					ff += 6
				}
				Func[304+ff].(fi4)(xp, yp, rp, ev(e))
			}
		} else {
			for rp < e {
				SetI32(rp, Func[f].(f2i)(I32(xp), I32(yp)))
				xp += ix
				yp += iy
				rp += 4
				continue
			}
		}
	case 2: // st
		trap() //type
	default: // ft zt
		if f < 241 && t == 5 {
			if av == 2 && ff == 3 { // v-a
				SetF64(yp, -F64(yp))
				addfF(yp, xp, rp, ev(e))
			} else {
				if av != 3 {
					ff += 6
				}
				Func[316+ff].(fi4)(xp, yp, rp, ev(e))
			}
		} else {
			for rp < e {
				Func[f-4+int32(t)].(fi3)(xp, yp, rp)
				xp += ix
				yp += iy
				rp += dz
				continue
			}
		}
	}
	dxy(x, y)
	return r
}
//todo: see v/v.go
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
		r = Enl(y)
		y = x
		yp = int32(y)
		x = r
		xp = int32(x)
		q = -q
		ix = 0
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
	if t < 5 { //cis
		if av < 3 {
			q += 9
		}
		Func[T(q)+281+3*t].(fi4)(xp, yp, rp, ev(e))
	} else {
		for rp < e {
			SetI32(rp, I32B(q == Func[250+t].(f2i)(xp, yp)))
			xp += ix
			yp += iy
			rp += 4
			continue
		}
	}
	dxy(x, y)
	return r
}
