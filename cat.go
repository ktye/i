package k

import (
	. "github.com/ktye/wg/module"
)

func Cat(x, y K) (r K) {
	xt, yt := tp(x), tp(y)
	if xt&15 == yt&15 {
		if xt < 16 {
			x = enl(x)
		}
		if yt < 16 {
			return cat1(x, y)
		} else {
			return ucat(x, y)
		}
	} else if xt == Lt && yt < 16 {
		return cat1(x, y)
	}
	trap(Nyi)
	return x
}
func enl(x K) (r K) {
	t := tp(x)
	if t < 16 {
		t += 16
		r = mk(t, 1)
		rp := int32(r)
		xp := int32(x)
		s := sz(t)
		if s == 1 {
			SetI8(rp, xp)
		} else if s == 4 {
			SetI32(rp, xp)
		} else {
			SetI64(rp, I64(xp))
			if t == Zt {
				SetI64(4+rp, I64(4+xp))
			}
		}
	} else {
		r = l1(x)
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
	ny := nn(y)
	r, rp, s := uspc(x, xt, ny)
	if xt == Lt {
		rl(y)
	}
	Memorycopy(rp, int32(y), s*ny)
	dx(y)
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
		if xt == Zt {
			trap(Nyi)
		} else if xt == Ft {
			SetI64(rp, I64(yp))
			dx(y)
		} else {
			SetI64(rp, int64(y))
		}
	}
	return r
}
func cat3(x, a, b, c K) K { // Lt
	r, rp, _ := uspc(x, Lt, 3)
	SetI64(rp, int64(a))
	SetI64(rp+8, int64(b))
	SetI64(rp+16, int64(c))
	return r
}
func uspc(x K, xt T, ny int32) (K, int32, int32) {
	var r K
	nx := nn(x)
	s := sz(xt)
	if I32(int32(x)-8) == 1 && bucket(s*nx) == bucket(s*(nx+ny)) {
		r = x
	} else {
		r = mk(xt, nx+ny)
		Memorycopy(int32(r), int32(x), s*nx)
		if xt == Lt {
			rl(x)
		}
		dx(x)
	}
	SetI32(int32(r)-4, nx+ny)
	return r, int32(r) + s*nx, s
}
func ncat(x, y K) (r K) {
	xt := tp(x)
	if xt < 16 {
		x = enl(x)
	}
	return cat1(uptypes(x, y))
}
func spl2(l K) (x, y K) {
	var lp int32
	x, lp = nxl(int32(l))
	y, lp = nxl(lp)
	dx(l)
	return x, y
}
func spl3(l K) (x, y, z K) {
	var lp int32
	x, lp = nxl(int32(l))
	y, lp = nxl(lp)
	z, lp = nxl(lp)
	dx(l)
	return x, y, z
}
func nxl(lp int32) (K, int32) { return K(I64(lp)), lp + 8 }
