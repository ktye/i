package main

import (
	. "github.com/ktye/wg/module"
)

func Atx(x, y K) (r K) { // x@y
	xt, yt := tp(x), tp(y)
	xp := int32(x)
	if xt < 16 {
		if xt == 0 || xt > tt {
			return cal(x, l1(y))
		}
		if xt == st {
			if xp == 0 && yt == it { // `123 (quoted verb)
				return K(int32(y))
			}
			if xp == 32 { // `k@
				return Kst(y)
			} else if xp == 40 { // `l@
				return Lst(y)
			} else {
				trap(Value)
			}
		}
	}
	if xt > Lt {
		x, r = spl2(x)
		if xt == Tt && yt&15 == it {
			trap(Nyi) // table-row-indexing
		}
		return Atx(r, Fnd(x, y))
	}
	if yt < It {
		y = uptype(y, it)
		yt = tp(y)
	}
	if yt == It || yt == Bt {
		return atv(x, y)
	}
	if yt == it && xt > 16 {
		return ati(x, int32(y))
	}
	if yt == Lt {
		return Ecr(19, l2(x, y))
	}
	trap(Nyi) // f@
	return x
}
func ati(x K, i int32) (r K) { // x BT..LT
	t := tp(x)
	if t < 16 {
		return x
	}
	if i < 0 || i >= nn(x) {
		dx(x)
		return missing(t - 16)
	}
	s := sz(t)
	p := int32(x) + i*s
	switch s >> 2 {
	case 0:
		r = K(U8(p))
	case 1:
		r = K(U32(p))
	case 2:
		r = K(U64(p))
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
func atv(x, y K) (r K) { // x BT..LT
	t := tp(x)
	yn := nn(y)
	if t < 16 {
		dx(y)
		return ntake(yn, x)
	}
	xn := nn(x)
	r = mk(t, yn)
	s := sz(t)
	rp := int32(r)
	xp := int32(x)
	yp := int32(y)
	na := missing(t - 16)
	switch s >> 2 {
	case 0:
		for i := int32(0); i < yn; i++ {
			xi := I32(yp)
			if xi < 0 || xi >= xn {
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
			if xi < 0 || xi >= xn {
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
			if xi < 0 || xi >= xn {
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
			if xi < 0 || xi >= xn {
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
func stv(x, i, y K) (r K) {
	x = use(x)
	xt := tp(x)
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
	xt, yt := tp(x), tp(y)
	if xt < Lt && yt != xt-16 {
		trap(Type)
	}
	xn := nn(x)
	if i < 0 || i >= xn {
		trap(Length)
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
		} else if xt != Ft {
			trap(Nyi)
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

func atdepth(x, y K) (r K) {
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
		if nn(y) != 1 {
			trap(Rank)
		}
		return Ecl(19, l2(x, Fst(y)))
	}
	return atdepth(x, y)
}
