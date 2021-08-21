package main

import (
	. "github.com/ktye/wg/module"
)

func Cat(x, y K) (r K) {
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
		return cat1(x, y)
	}
	return uf(Cat(explode(x), explode(y)))
}
func Enl(x K) (r K) {
	t := tp(x)
	if t < 16 {
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
	} else if t == Dt {
		r = Flp(Ech(13, l1(x)))
	} else {
		r = l1(x)
	}
	return r
}
func explode(x K) (r K) {
	xt := tp(x)
	if xt < 16 {
		r = l1(x)
	} else if xt < Lt {
		xn := nn(x)
		r = mk(Lt, xn)
		if xn == 0 {
			dx(x)
			return r
		}
		xp, rp := int32(x), int32(r)
		e := ep(x)
		switch xt - 17 {
		case 0: //Bt
			for xp < e {
				SetI64(rp, int64(Kb(I8(xp))))
				rp += 8
				xp++
				continue
			}
		case 1: //Ct
			for xp < e {
				SetI64(rp, int64(Kc(I8(xp))))
				rp += 8
				xp++
				continue
			}
		case 2: //It
			for xp < e {
				SetI64(rp, int64(Ki(I32(xp))))
				rp += 8
				xp += 4
				continue
			}
		case 3: //St
			for xp < e {
				SetI64(rp, int64(Ks(I32(xp))))
				rp += 8
				xp += 4
				continue
			}
		case 4: //Ft
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
		k, x = spl2(x)
		r = mk(Lt, xn)
		rp := int32(r)
		x = Flp(x)
		xp := int32(x)
		for i := int32(0); i < xn; i++ {
			SetI64(rp, int64(Key(rx(k), x0(xp))))
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
		r = Cat(r, x0(xp))
		xp += 8
	}
	dx(x)
	return r
}
func ucat(x, y K) K { // Bt,Bt .. Lt,Lt
	xt := tp(x)
	//if xt != tp(y) {
	//	panic("ucat")
	//}
	if xt > Lt {
		return dcat(x, y)
	}
	ny := nn(y)
	r, rp, s := uspc(x, xt, ny)
	if xt == Lt {
		rl(y)
	}
	Memorycopy(rp, int32(y), s*ny)
	dx(y)
	return r
}
func dcat(x, y K) (r K) { // d,d  t,t
	var q K
	t := tp(x)
	r, x = spl2(x)
	q, y = spl2(y)
	if t == Dt {
		return Key(Cat(r, q), Cat(x, y))
	} else {
		if match(r, q) == 0 {
			trap(Value)
		}
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
	var rn int32
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
	r, rp, s := uspc(x, xt, 1)
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
		F64x2store(rp, F64x2load(yp))
		dx(y)
	}
	return r
}
func uspc(x K, xt T, ny int32) (K, int32, int32) {
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
	return r, int32(r) + s*nx, s
}
func ncat(x, y K) (r K) {
	xt := tp(x)
	if xt < 16 {
		x = Enl(x)
	}
	return cat1(uptypes(x, y, 0))
}
func spl2(l K) (x, y K) {
	lp := int32(l)
	x, y = x0(lp), x1(lp)
	dx(l)
	return x, y
}
func spl3(l K) (x, y, z K) {
	lp := int32(l)
	x, y, z = x0(lp), x1(lp), x2(lp)
	dx(l)
	return x, y, z
}
