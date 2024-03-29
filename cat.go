package main

import (
	. "github.com/ktye/wg/module"
)

func Cat(x, y K) K {
	xt, yt := tp(x), tp(y)
	if xt == Tt && yt == Dt {
		return dcat(x, y)
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
	var r K
	xt := tp(x)
	if xt < 16 || xt == Dt {
		return l1(x)
	} else if xt < Lt {
		xn := nn(x)
		r = mk(Lt, nn(x))
		rp := int32(r)
		for i := int32(0); i < xn; i++ {
			SetI64(rp+8*i, int64(ati(rx(x), i)))
		}
		dx(x)
		return r
	} else if xt == Tt { // Tt
		xn := nn(x)
		k := x0(x)
		x = Flp(r1(x))
		r = mk(Lt, 0)
		for i := int32(0); i < xn; i++ {
			r = cat1(r, Key(rx(k), ati(rx(x), i)))
		}
		dxy(x, k)
		return r
	}
	return x
}
func ucat(x, y K) K { // Bt,Bt .. Tt,Tt
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
func ucat1(x, y, z K) K { return cat1(ucat(x, y), z) }
func cat1(x, y K) K {
	t := tp(x)
	x = uspc(x, t, 1)
	if t == Lt {
		y = l1(rx(y))
		x = ti(Ft, int32(x))
	}
	return ti(t, int32(sti(x, nn(x)-1, y)))
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
