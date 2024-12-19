//go:build !(simd4||simd5)

package main

import (
	. "github.com/ktye/wg/module"
)

const b0 = 5 //smallest bucket
const bs = 15 //lower bits set
const vl = 16 //vector length

func seq(n int32) K {
        n = maxi(n, 0)
        r := mk(It, n)
        for n > 0 {
                n--
                SetI32(int32(r)+4*n, n)
        }
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
		for xp < e {
			SetI32(xp, Func[f].(f1i)(I32(xp)))
			xp += 4
			continue
		}
	case 2:
		trap() //type
	default: //F/Z (only called for neg)
		for xp < e {
			SetF64(xp, Func[1+f].(f1f)(F64(xp)))
			xp += 8
			continue
		}
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
		if I32(int32(y)-4) == 1 {
			r = rx(y)
		} else {
			r = use1(x)
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
