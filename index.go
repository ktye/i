package k

import (
	. "github.com/ktye/wg/module"
)

func Atx(x, y K) K { // x@y
	xt, yt := tp(x), tp(y)
	if xt < 16 {
		if xt == 0 || xt > tt {
			return cal(x, l1(y))
		}
	}
	if xt > Lt {
		trap(Nyi) // d@ t@
	}
	if yt == It && xt > 16 {
		return atv(x, y)
	}
	if yt == it && xt > 16 {
		return ati(x, int32(y))
	}
	trap(Nyi) // f@
	return x
}
func ati(x K, i int32) (r K) { // x BT..LT
	t := tp(x)
	s := sz(t)
	p := int32(x) + i*s
	if t == Zt {
		var im K
		if ip := I64(8 + int32(x)); ip != 0 {
			im = ati(rx(K(ip)), i)
		}
		r = ati(rx(K(int32(x))), i)
		dx(x)
		return Kz(r, im)
	}
	if s == 1 {
		r = K(U8(p))
	} else if s == 4 {
		r = K(U32(p))
	} else {
		r = K(U64(p))
	}
	if t == Ft {
		r = Kf(F64reinterpret_i64(uint64(r)))
	} else if t == Lt {
		r = rx(r)
		dx(x)
		return r
	}
	dx(x)
	return r | K(t-16)<<59
}
func atv(x, y K) (r K) { // x BT..LT
	t := tp(x)
	if t == Zt {
		var im K
		if ip := I64(8 + int32(x)); ip != 0 {
			im = atv(rx(K(ip)), rx(y))
		}
		r = atv(rx(K(int32(x))), y)
		dx(x)
		return KZ(r, im)
	}
	xn, yn := nn(x), nn(y)
	r = mk(t, yn)
	s := sz(t)
	rp := int32(r)
	xp := int32(x)
	yp := int32(y)
	if s == 1 {
		for i := int32(0); i < yn; i++ {
			xi := I32(yp)
			if xi < 0 || xi >= xn {
				SetI8(rp, 0)
			} else {
				SetI8(rp, I8(xp+xi))
			}
			rp++
			yp += 4
		}
	} else if s == 4 {
		for i := int32(0); i < yn; i++ {
			xi := I32(yp)
			if xi < 0 || xi >= xn {
				SetI32(rp, 0)
			} else {
				SetI32(rp, I32(xp+4*xi))
			}
			rp += 4
			yp += 4
		}
	} else {
		for i := int32(0); i < yn; i++ {
			xi := I32(yp)
			if xi < 0 || xi >= xn {
				SetI64(rp, 0)
			} else {
				SetI64(rp, I64(xp+8*xi))
			}
			rp += 8
			yp += 4
		}
	}
	if t == Lt {
		rl(r)
	}
	dx(x)
	dx(y)
	return r
}
func use(x K) (K, T) {
	xt := tp(x)
	if xt < 16 {
		trap(Type)
	}
	x = ucat(x, mk(xt, 0))
	return x, xt
}
func stv(x, i, y K) (r K) {
	var xt T
	x, xt = use(x)
	if It != tp(i) {
		trap(Type)
	}
	if xt != tp(y) {
		trap(Type)
	}
	xn := nn(x)
	n := nn(i)
	if n != nn(y) {
		trap(Length)
	}
	s := sz(xt)
	xp := int32(x)
	yp := int32(y)
	ip := int32(i)
	for j := int32(0); j < n; j++ {
		xi := I32(ip + 4*j)
		if xi < 0 || xi >= xn {
			trap(Value)
		}
	}
	if s == 1 {
		for j := int32(0); j < n; j++ {
			SetI8(xp+I32(ip), I8(yp))
			ip += 4
			yp++
		}
	} else if s == 4 {
		for j := int32(0); j < n; j++ {
			SetI32(xp+4*I32(ip), I32(yp))
			ip += 4
			yp += 4
		}
	} else if s == 8 {
		if xt == Ft || xt == Zt {
			trap(Nyi)
		}
		if xt == Lt {
			rl(y)
			for j := int32(0); j < n; j++ {
				dx(K(I64(xp + 8*I32(ip))))
				ip += 4
			}
			ip = int32(i)
		}
		for j := int32(0); j < n; j++ {
			SetI64(xp+8*I32(ip), I64(yp))
			ip += 4
			yp += 8
		}
	}
	dx(i)
	dx(y)
	return x
}
func sti(x K, i int32, y K) K {
	var xt, yt T
	x, xt = use(x)
	yt = tp(y)
	if xt < Lt && yt != xt-16 {
		trap(Type)
	}
	xn := nn(x)
	if i < 0 || i >= xn {
		trap(Length)
	}
	s := sz(xt)
	xp := int32(x)
	if xt == Ft || xt == Zt {
		trap(Nyi)
	}
	if s == 1 {
		SetI8(xp+i, int32(y))
	} else if s == 4 {
		SetI32(xp+4*i, int32(y))
	} else {
		xp += 8 * i
		if xt == Lt {
			dx(K(I64(xp)))
		}
		SetI64(xp, int64(y))
	}
	return x
}
