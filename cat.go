package k

import (
	. "github.com/ktye/wg/module"
)

func enl(x K) K {
	trap(Nyi)
	return x
}
func ucat(x, y K) (r K) { // Bt,Bt .. Lt,Lt
	t, nx, ny := tp(x), nn(x), nn(y)
	if t == Zt {
		trap(Nyi)
	}
	s := sz(t)
	var rp int32
	if I32(int32(x)-8) == 1 && bucket(s*nx) == bucket(s*(nx+ny)) {
		r = rx(x)
		rp = int32(r)
	} else {
		r = mk(t, nx+ny)
		rp = int32(r)
		Memorycopy(rp, int32(x), s*nx)
		if t == Lt {
			rl(x)
		}
	}
	Memorycopy(rp+s*nx, int32(y), s*ny)
	if t == Lt {
		rl(y)
	}
	dx(x)
	dx(y)
	SetI32(rp-4, nx+ny)
	return r
}

func lcat(x, y K) (r K) {
	nx := nn(x)
	n8 := nx * 8
	if I32(int32(x)-8) == 1 && bucket(n8) == bucket(8+n8) {
		r = x
	} else {
		r = mk(Lt, 1+nx)
		Memorycopy(int32(r), int32(x), n8)
		rl(x)
		dx(x)
	}
	SetI64(int32(r)+n8, int64(y))
	SetI32(int32(r)-4, 1+nx)
	return r
}
