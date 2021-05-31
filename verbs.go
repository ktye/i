package main

import (
	. "github.com/ktye/wg/module"
)

func nyi(x K) K { trap(Nyi); return x }
func Dex(x, y K) K { // x:y
	dx(x)
	return y
}
func Flp(x K) K { // +x
	trap(Nyi)
	return x
}
func Fst(x K) K { // *x
	t := tp(x)
	if t < 16 {
		return x
	}
	n := nn(x)
	if n == 0 {
		dx(x)
		return zero(t - 16)
	}
	return ati(x, 0)
}
func Lst(x K) K { // *|x
	t := tp(x)
	if t < 16 {
		return x
	}
	n := nn(x)
	if n == 0 {
		return Fst(x)
	}
	return ati(x, n-1)
}

func Cnt(x K) (r K) { // #x
	t := tp(x)
	dx(x)
	if t < 16 {
		return Ki(1)
	}
	return Ki(nn(x))
}
func Til(x K) (r K) {
	if tp(x) != it {
		trap(Type)
	}
	return seq(int32(x))
}
func seq(n int32) (r K) {
	n = maxi(n, 0)
	r = mk(It, n)
	p := int32(r)
	for i := int32(0); i < n; i++ {
		SetI32(p, i)
		p += 4
	}
	return r
}
func Unq(x K) (r K) { // ?x
	return x // todo
}
func Key(x, y K) (r K) { // x!y
	xt, yt := tp(x), tp(y)
	if xt < 16 {
		x = enl(x)
		xt = tp(x)
	}
	if yt < 16 {
		y = enl(y)
		yt = tp(y)
	}
	nx, ny := nn(x), nn(y)
	if nx != ny {
		if nx == 1 && ny > 1 {
			x = ntake(ny, x)
		} else if ny == 1 && nx > 1 {
			y = ntake(ny, x)
		} else {
			trap(Length)
		}
	}
	r = l2(x, y)
	return K(int32(r)) | K(dt)<<59
}

func Tak(x, y K) (r K) { // x#y
	xt := tp(x)
	if xt == it {
		return ntake(int32(x), y)
	}
	trap(Nyi) // set take
	return 0
}
func ntake(n int32, y K) (r K) {
	t := tp(y)
	if n < 0 {
		if tp(y) < 16 {
			trap(Type)
		}
		n += nn(y)
		if n < 0 {
			return ucat(ntake(-n, zero(t-16)), y)
		}
		return ndrop(n, y)
	}
	yp := int32(y)
	if t < 5 {
		t += 16
		r = mk(t, n)
		s := sz(t)
		rp := int32(r)
		if s == 1 {
			Memoryfill(rp, yp, n)
		} else {
			for i := int32(0); i < n; i++ {
				SetI32(rp, yp)
				rp += 4
			}
		}
		return r
	} else if t == ft {
		r = mk(Ft, n)
		rp := int32(r)
		f := F64(yp)
		for i := int32(0); i < n; i++ {
			SetF64(rp, f)
			rp += 8
		}
		dx(y)
		return r
	} else if t == zt {
		r = mk(Zt, n)
		rp := int32(r)
		re, im := F64(yp), F64(yp+8)
		for i := int32(0); i < n; i++ {
			SetF64(rp, re)
			SetF64(rp+8, im)
			rp += 16
		}
		dx(y)
		return r
	} else if t < 16 {
		r = mk(Lt, n)
		rp := int32(r)
		for i := int32(0); i < n; i++ {
			SetI64(rp, int64(rx(y)))
			rp += 8
		}
		dx(y)
		return r
	}
	return Atx(y, seq(n))
}
func Drp(x, y K) (r K) { // x_y
	xt := tp(x)
	if xt == it {
		return ndrop(int32(x), y)
	}
	trap(Nyi) // set drop
	return 0
}
func ndrop(n int32, y K) (r K) {
	yt := tp(y)
	if yt < 16 {
		trap(Type)
	}
	if yt > Lt {
		trap(Nyi)
	}
	yn := nn(y)
	if n < 0 {
		return ntake(maxi(0, yn+n), y)
	}
	rn := yn - n
	if rn < 0 {
		dx(y)
		return mk(yt, 0)
	}
	s := sz(yt)
	yp := int32(y)
	if I32(yp-16) == 1 && bucket(s*yn) == bucket(s*rn) {
		r = rx(y)
		SetI32(yp-12, rn)
	} else {
		r = mk(yt, rn)
	}
	rp := int32(r)
	Memorycopy(rp, yp+s*n, s*rn)
	dx(y)
	return r
}

func Rev(x K) (r K) { // |x
	t := tp(x)
	if t < 16 {
		return x
	}
	if t >= Ft && t != Lt {
		panic(Nyi)
	}
	xn := nn(x)
	if xn < 2 {
		return x
	}
	r = mk(It, xn)
	rp := int32(r) + 4*xn
	for i := int32(0); i < xn; i++ {
		rp -= 4
		SetI32(rp, i)
	}
	return atv(x, r)
}

func Wer(x K) (r K) { // &x
	t := tp(x)
	if t < 16 {
		x = enl(x)
		t = tp(x)
	}
	var n, rp int32
	xn := nn(x)
	xp := int32(x)
	if t == Bt {
		n = sumb(x)
		r = mk(It, n)
		rp = int32(r)
		for i := int32(0); i < xn; i++ {
			if I8(xp) != 0 {
				SetI32(rp, i)
				rp += 4
			}
			xp++
		}
	} else if t == It {
		n = sumi(x)
		r = mk(It, n)
		rp = int32(r)
		for i := int32(0); i < xn; i++ {
			j := I32(xp)
			for k := int32(0); k < j; k++ {
				SetI32(rp, i)
				rp += 4
			}
			xp += 4
		}
	} else {
		trap(Type)
	}
	dx(x)
	return r
}

func Typ(x K) (r K) { // @x
	r = Ki(int32(tp(x)))
	dx(x)
	return r
}
func Val(x K) (r K) {
	xt := tp(x)
	if xt == Ct {
		return val(x)
	}
	trap(Nyi)
	return x
}
func val(x K) (r K) {
	x = parse(x)
	xn := nn(x)
	xp := int32(x) + 8*(xn-1)
	a := int32(0)
	if xn > 2 && I64(xp) == 2 && I64(xp-8) == 0 {
		a = 1
	}
	x = exec(x)
	if a != 0 {
		dx(x)
		return 0
	}
	return x
}

func sumb(x K) (r int32) {
	p := int32(x)
	e := p + nn(x)
	for p < e {
		r += I8(p)
		p++
	}
	return r
}
func sumi(x K) (r int32) {
	p := int32(x)
	e := p + 4*nn(x)
	for p < e {
		r += I8(p)
		p += 4
	}
	return r
}
