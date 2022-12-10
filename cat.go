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
func Enl(x K) (r K) {
	t := tp(x)
	if t < 7 {
		t += 16
		r = mk(t, 1)
		rp := int32(r)
		xp := int32(x)
		s := sz(t)
		switch s >> 2 {
		case 0:
			SetI8(rp, xp)
		case 1:
			SetI32(rp, xp)
		case 2:
			SetI64(rp, I64(xp))
		case 3:
		case 4:
			SetI64(rp, I64(xp))
			SetI64(rp+8, I64(xp+8))
		}
		dx(x)
		return r
	}
	if t == Dt {
		if tp(K(I64(int32(x)))) == St {
			return Flp(Ech(13, l1(x))) // +,'x
		}
	}
	return l1(x)
}
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
		r = mk(Lt, xn)
		if xn == 0 {
			dx(x)
			return r
		}
		xp, rp := int32(x), int32(r)
		e := ep(x)
		switch xt - 18 {
		case 0: //Ct
			for xp < e {
				SetI64(rp, int64(Kc(I8(xp))))
				rp += 8
				xp++
				continue
			}
		case 1: //It
			for xp < e {
				SetI64(rp, int64(Ki(I32(xp))))
				rp += 8
				xp += 4
				continue
			}
		case 2: //St
			for xp < e {
				SetI64(rp, int64(Ks(I32(xp))))
				rp += 8
				xp += 4
				continue
			}
		case 3: //Ft
			for xp < e {
				SetI64(rp, int64(Kf(F64(xp))))
				rp += 8
				xp += 8
				continue
			}
		default: //Zt
			for xp < e {
				SetI64(rp, int64(Kz(F64(xp), F64(xp+8))))
				rp += 8
				xp += 16
				continue
			}
		}
		dx(x)
	} else if xt == Lt {
		r = x
	} else if xt == Tt {
		var k K
		xn := nn(x)
		k = x0(x)
		x = r1(x)
		r = mk(Lt, xn)
		rp := int32(r)
		x = Flp(x)
		xp := int32(x)
		for i := int32(0); i < xn; i++ {
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
func flat(x K) (r K) { // ((..);(..)) -> (...)
	r = mk(Lt, 0)
	xn := nn(x)
	xp := int32(x)
	for i := int32(0); i < xn; i++ {
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
func dcat(x, y K) (r K) { // d,d  t,t
	var q K
	t := tp(x)
	if t == Tt {
		if match(K(I64(int32(x))), K(I64(int32(y)))) == 0 {
			return ucat(explode(x), explode(y))
		}
	}
	r = x0(x)
	x = r1(x)
	q = x0(y)
	y = r1(y)
	if t == Dt {
		return Key(Cat(r, q), Cat(x, y))
	} else {
		dx(q)
		x = Ech(13, l2(x, y))
		return key(r, x, t)
	}
}
func ucats(x K) (r K) { // ,/ unitype-lists
	xn := nn(x)
	if xn == 0 {
		return x
	}
	xp := int32(x)
	var rt T
	rn := int32(0)
	for i := int32(0); i < xn; i++ {
		xi := K(I64(xp))
		t := tp(xi)
		if i == 0 {
			rt = t
		}
		if rt != t || rt < 16 || t > Zt {
			return 0
		}
		rn += nn(xi)
		xp += 8
	}
	r = mk(rt, rn)
	s := sz(rt)
	rp := int32(r)
	xp = int32(x)
	for i := int32(0); i < xn; i++ {
		xi := K(I64(xp))
		rn = s * nn(xi)
		Memorycopy(rp, int32(xi), rn)
		rp += rn
		xp += 8
	}
	dx(x)
	return r
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
		//F64x2store(rp, F64x2load(yp))
		Memorycopy(rp, yp, 16)
		dx(y)
	}
	return r
}
func uspc(x K, xt T, ny int32) K {
	var r K
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
