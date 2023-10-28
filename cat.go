package main

import (
	. "github.com/ktye/wg/module"
)

func Cat(x, y K) K {
	xt, yt := tp(x), tp(y)
	if xt == Tt && yt == Dt {
		y, yt = Enl(y), Tt
	}
	if xt&15 == yt&15 {
		if xt < 16 {
			x = Enl(x)
		}
		if yt < 16 {
			return cat1(x, y)
		} else {
			return ucat(x, y)
		}
	} else if xt == Lt && yt < 16 {
		if nn(x) > 0 {
			return cat1(x, y)
		}
	}
	x = uf(Cat(explode(x), explode(y)))
	if nn(x) == 0 {
		dx(x)
		return mk(xt|16, 0)
	}
	return x
}
func Enl(x K) K { return uf(l1(x)) }
func explode(x K) K {
	r := K(0)
	xt := tp(x)
	if xt < 16 {
		r = l1(x)
	} else if xt == Dt {
		r = mk(Lt, 1)
		SetI64(int32(r), int64(x))
	} else if xt < Lt {
		xn := nn(x)
		r = mk(Lt, nn(x))
		rp := int32(r)
		for i := int32(0); i < xn; i++ {
			SetI64(rp+8*i, int64(ati(rx(x), i)))
		}
		dx(x)
	} else if xt == Lt {
		r = x
	} else if xt == Tt {
		xn := nn(x)
		k := x0(x)
		x = r1(x)
		r = mk(Lt, xn)
		rp := int32(r)
		x = Flp(x)
		xp := int32(x)
		e := ep(r)
		for rp < e {
			SetI64(rp, int64(Key(rx(k), x0(K(xp)))))
			xp += 8
			rp += 8
		}
		dx(x)
		dx(k)
		return r
	}
	return r
}
func flat(x K) K { // ((..);(..)) -> (...)
	r := mk(Lt, 0)
	xp := int32(x)
	e := ep(x)
	for xp < e {
		r = Cat(r, x0(K(xp)))
		xp += 8
	}
	dx(x)
	return r
}
func ucat(x, y K) K { // Bt,Bt .. Lt,Lt
	xt := tp(x)
	if xt > Lt {
		return dcat(x, y)
	}
	xn := nn(x)
	yn := nn(y)
	r := uspc(x, xt, yn)
	s := sz(xt)
	if xt == Lt {
		rl(y)
	}
	Memorycopy(int32(r)+s*xn, int32(y), s*yn)
	dx(y)
	return r
}
func dcat(x, y K) K { // d,d  t,t
	t := tp(x)
	if t == Tt {
		if match(K(I64(int32(x))), K(I64(int32(y)))) == 0 {
			return ucat(explode(x), explode(y))
		}
	}
	r := x0(x)
	x = r1(x)
	q := x0(y)
	y = r1(y)
	if t == Dt {
		r = Cat(r, q)
		return Key(r, Cat(x, y))
	} else {
		dx(q)
		x = Ech(13, l2(x, y))
		return key(r, x, t)
	}
}
func cat1(x, y K) K {
	xt := tp(x)
	xn := nn(x)
	r := uspc(x, xt, 1)
	s := sz(xt)
	rp := int32(r) + s*xn
	yp := int32(y)
	if s == 1 {
		SetI8(rp, yp)
	} else if s == 4 {
		SetI32(rp, yp)
	} else if s == 8 {
		if xt == Ft {
			SetI64(rp, I64(yp))
			dx(y)
		} else {
			SetI64(rp, int64(y))
		}
	} else if s == 16 {
		Memorycopy(rp, yp, 16)
		dx(y)
	}
	return r
}
func uspc(x K, xt T, ny int32) K {
	r := K(0)
	nx := nn(x)
	s := sz(xt)
	if I32(int32(x)-4) == 1 && bucket(s*nx) == bucket(s*(nx+ny)) {
		r = x
	} else {
		r = mk(xt, nx+ny)
		Memorycopy(int32(r), int32(x), s*nx)
		if xt == Lt {
			rl(x)
		}
		dx(x)
	}
	SetI32(int32(r)-12, nx+ny)
	return r
}
func ncat(x, y K) K {
	xt := tp(x)
	if xt < 16 {
		x = Enl(x)
	}
	xt = maxtype(x, y)
	x = uptype(x, xt)
	y = uptype(y, xt)
	return cat1(x, y)
}
