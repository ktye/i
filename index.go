package main

import (
	. "github.com/ktye/wg/module"
)

func Atx(x, y K) K { // x@y
	r := K(0)
	xt, yt := tp(x), tp(y)
	xp := int32(x)
	if xt < 16 {
		if xt == 0 || xt > tt {
			return cal(x, l1(y))
		}
		if xt == st {
			if xp == 0 {
				if yt == it { // `123 (quoted verb)
					return K(int32(y))
				}
			}
			return cal(Val(sc(cat1(cs(x), Kc('.')))), l1(y))
		}
	}
	if xt > Lt && yt < Lt {
		r = x0(x)
		x = r1(x)
		if xt == Tt {
			if yt&15 == it {
				return key(r, Ecl(19, l2(x, y)), Dt+T(I32B(yt == It)))
			}
		}
		return Atx(x, Fnd(r, y))
	}
	if yt&15 == ft {
		return Rot(x, y)
	}
	if yt < It {
		y = uptype(y, it)
		yt = tp(y)
	}
	if yt == It {
		return atv(x, y)
	}
	if yt == it {
		return ati(x, int32(y))
	}
	if yt == Lt {
		return Ecr(19, l2(x, y))
	}
	if yt == Dt {
		r = x0(y)
		return Key(r, Atx(x, r1(y)))
	}
	return trap(Type) // f@
}
func ati(x K, i int32) K { // x CT..LT
	r := K(0)
	t := tp(x)
	if t < 16 {
		return x
	}
	if t > Lt {
		return Atx(x, Ki(i))
	}
	if i < 0 || i >= nn(x) {
		dx(x)
		return missing(t - 16)
	}
	s := sz(t)
	p := int32(x) + i*s
	switch s >> 2 {
	case 0:
		r = K(uint32(I8(p)))
	case 1:
		r = K(uint32(I32(p)))
	case 2:
		r = K(uint64(I64(p)))
	default:
		dx(x)
		return Kz(F64(p), F64(p+8))
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
func atv(x, y K) K { // x CT..LT
	t := tp(x)
	if t == Tt {
		return Atx(x, y)
	}
	yn := nn(y)
	if t < 16 {
		dx(y)
		return ntake(yn, x)
	}
	xn := nn(x)
	r := mk(t, yn)
	s := sz(t)
	rp := int32(r)
	xp := int32(x)
	yp := int32(y)

	na := missing(t - 16)
	switch s >> 2 {
	case 0:
		for i := int32(0); i < yn; i++ {
			xi := I32(yp)
			if uint32(xi) >= uint32(xn) {
				SetI8(rp, int32(na))
			} else {
				SetI8(rp, I8(xp+xi))
			}
			rp++
			yp += 4
		}
	case 1:
		for i := int32(0); i < yn; i++ {
			xi := I32(yp)
			if uint32(xi) >= uint32(xn) {
				SetI32(rp, int32(na))
			} else {
				SetI32(rp, I32(xp+4*xi))
			}
			rp += 4
			yp += 4
		}
	case 2:
		for i := int32(0); i < yn; i++ {
			xi := I32(yp)
			if uint32(xi) >= uint32(xn) {
				if t == Lt {
					SetI64(rp, int64(na))
				} else {
					SetI64(rp, I64(int32(na)))
				}
			} else {
				SetI64(rp, I64(xp+8*xi))
			}
			rp += 8
			yp += 4
		}
	default:
		for i := int32(0); i < yn; i++ {
			xi := I32(yp)
			if uint32(xi) >= uint32(xn) {
				SetI64(rp, I64(int32(na)))
				SetI64(rp+8, I64(int32(na)))
			} else {
				xi *= 16
				SetI64(rp, I64(xp+xi))
				SetI64(rp+8, I64(8+xp+xi))
			}
			rp += 16
			yp += 4
		}
	}
	if t == Lt {
		rl(r)
		r = uf(r)
	}
	dx(na)
	dx(x)
	dx(y)
	return r
}
func stv(x, i, y K) K {
	if It != tp(i) {
		trap(Type)
	}
	n := nn(i)
	if n == 0 {
		dx(y)
		dx(i)
		return x
	}
	if n != nn(y) {
		trap(Length)
	}
	x = use(x)
	xt := tp(x)
	xn := nn(x)
	s := sz(xt)
	xp := int32(x)
	yp := int32(y)
	ip := int32(i)
	for j := int32(0); j < n; j++ {
		xi := uint32(I32(ip + 4*j))
		if xi >= uint32(xn) {
			trap(Index)
		}
	}
	switch s >> 2 {
	case 0:
		for j := int32(0); j < n; j++ {
			SetI8(xp+I32(ip), I8(yp))
			ip += 4
			yp++
		}
	case 1:
		for j := int32(0); j < n; j++ {
			SetI32(xp+4*I32(ip), I32(yp))
			ip += 4
			yp += 4
		}
	case 2:
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
		if xt == Lt {
			x = uf(x)
		}
	default:
		for j := int32(0); j < n; j++ {
			xp = int32(x) + 16*I32(ip)
			SetI64(xp, I64(yp))
			SetI64(xp+8, I64(yp+8))
			ip += 4
			yp += 16
		}
	}
	dx(i)
	dx(y)
	return x
}
func sti(x K, i int32, y K) K {
	x = use(x)
	xt := tp(x)
	//if xt < Lt && yt != xt-16 {
	//	trap(Type)
	//}
	xn := nn(x)
	if i < 0 || i >= xn {
		trap(Index)
	}
	s := sz(xt)
	xp := int32(x)
	yp := int32(y)
	switch s >> 2 {
	case 0:
		SetI8(xp+i, yp)
	case 1:
		SetI32(xp+4*i, yp)
	case 2:
		xp += 8 * i
		if xt == Lt {
			dx(K(I64(xp)))
			SetI64(xp, int64(rx(y)))
			x = uf(x)
		} else {
			SetI64(xp, I64(yp))
		}
	default:
		xp += 16 * i
		SetI64(xp, I64(yp))
		SetI64(xp+8, I64(yp+8))
	}
	dx(y)
	return x
}

func atdepth(x, y K) K {
	xt := tp(x)
	if xt < 16 {
		trap(Type)
	}
	f := Fst(rx(y))
	if f == 0 {
		f = seq(nn(x))
	}
	x = Atx(x, f)
	if nn(y) == 1 {
		dx(y)
		return x
	}
	y = ndrop(1, y)
	if tp(f) > 16 {
		if nn(y) == 1 && xt == Tt {
			return Atx(x, Fst(y))
		}
		return Ecl(20, l2(x, y))
	}
	return atdepth(x, y)
}
