package main

import (
	. "github.com/ktye/wg/module"
)

func nyi(x K) K { trap(Nyi); return x }
func Dex(x, y K) K { // x:y
	dx(x)
	return y
}
func Flp(x K) (r K) { // +x
	xt := tp(x)
	switch xt - Lt {
	case 0: // Lt   n:#x;  m:|/#x (,/m#/:x)[(!m)+\:m*!n]
		n := nn(x)
		xp := int32(x)
		m := Ki(maxcount(xp, n))
		x = Atx(Rdc(13, l1(Ecr(15, l2(m, x)))), Ecl(2, l2(Til(m), Mul(m, Til(Ki(n))))))
		return x
	case 1: // Dt
		return td(x)
	case 2: // Tt
		return Key(spl2(x))
	default:
		return Abs(x)
	}
}
func maxcount(xp int32, n int32) (r int32) { // |/#l
	for i := int32(0); i < n; i++ {
		x := K(I64(xp))
		xp += 8
		if tp(x) < 16 {
			r = maxi(1, r)
		} else {
			r = maxi(nn(x), r)
		}
	}
	return r
}
func Fst(x K) K { // *x
	t := tp(x)
	if t < 16 {
		return x
	}
	if t > Lt {
		x = Val(x)
		t = tp(x)
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
	xt := tp(x)
	if xt > Lt {
		r = x0(int32(x))
		dx(x)
		return r
	}
	if xt != it {
		trap(Type)
	}
	return seq(int32(x))
}
func seq(n int32) (r K) {
	n = maxi(n, 0)
	r = mk(It, n)
	if n == 0 {
		return r
	}
	p := int32(r)
	SetI32(p, 0)
	SetI32(p+4, 1)
	SetI32(p+8, 2)
	SetI32(p+12, 3)
	v := I32x4load(p)
	w := I32x4splat(4)
	e := ep(r)
	for p < e {
		I32x4store(p, v)
		v = v.Add(w)
		p += 16
		continue
	}
	return r
}
func Unq(x K) (r K) { // ?x
	xt := tp(x)
	if xt < 16 || xt >= Lt {
		trap(Type)
	}
	xn := nn(x)
	r = mk(xt, 0)
	for i := int32(0); i < xn; i++ {
		xi := ati(rx(x), i)
		if int32(In(xi, rx(r))) == 0 {
			r = cat1(r, xi)
		} else {
			dx(xi)
		}
	}
	dx(x)
	return r
}
func unqs(x K) (r K) { // ?^x
	xt := tp(x)
	if xt < 16 {
		trap(Type)
	}
	return Atx(x, Wer(Ecp(12, l1(rx(x))))) // x(&~':x)
}
func Grp(x K) (r K) { // =x
	u := Unq(rx(x))
	n := nn(u)
	r = mk(Lt, n)
	rp := int32(r)
	for i := int32(0); i < n; i++ {
		SetI64(rp, int64(Wer(Eql(ati(rx(u), i), rx(x)))))
		rp += 8
	}
	dx(x)
	return Key(u, r)
}
func Key(x, y K) (r K) { // x!y
	xt, yt := tp(x), tp(y)
	if xt < 16 {
		x = Enl(x)
		xt = tp(x)
	}
	if yt < 16 {
		y = Enl(y)
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
	SetI32(int32(r)-12, nn(x))
	return K(int32(r)) | K(Dt)<<59
}
func key(x, y K, t T) K { return K(int32(Key(x, y))) | K(t)<<59 } // Dt or Tt

func Tak(x, y K) (r K) { // x#y
	xt := tp(x)
	if xt == it {
		return ntake(int32(x), y)
	}
	yt := tp(y)
	if xt > 16 && xt == yt {
		return atv(y, Wer(In(rx(y), x))) // set take
	}
	trap(Nyi) // f take
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
	yt := tp(y)
	if xt > 16 && xt == yt {
		return atv(y, Wer(Not(In(rx(y), x)))) // set drop
	}
	trap(Nyi) // f drop
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
	if I32(yp-4) == 1 && bucket(s*yn) == bucket(s*rn) && yt < Lt {
		r = rx(y)
		SetI32(yp-12, rn)
	} else {
		r = mk(yt, rn)
	}
	rp := int32(r)
	Memorycopy(rp, yp+s*n, s*rn)
	if yt == Lt {
		rl(r)
	}
	dx(y)
	return r
}

func Cut(x, y K) (r K) { // x^y
	xt := tp(x)
	if xt == It {
		return cuts(x, y)
	}
	yt := tp(y)
	if xt == Ct && yt == Ct { // "set"^"abc"
		x = Wer(In(rx(y), x))
		return rcut(y, Cat(Ki(0), Add(Ki(1), rx(x))), Cat(x, Ki(nn(y))))
	}
	if xt != it || yt < 16 {
		trap(Type)
	}
	xp := int32(x)
	if xp <= 0 {
		trap(Value)
	}
	r = mk(Lt, xp)
	rp := int32(r)
	e := ep(r)
	n := nn(y) / xp
	x = seq(n)
	for rp < e {
		SetI64(rp, int64(atv(rx(y), rx(x))))
		x = Add(Ki(n), x)
		rp += 8
		continue
	}
	dx(x)
	dx(y)
	return r
}
func cuts(x, y K) K { return rcut(y, rx(x), cat1(ndrop(1, x), Ki(nn(y)))) }
func rcut(x, a, b K) (r K) { // a, b start-stop ranges
	n := nn(a)
	ap, bp := int32(a), int32(b)
	r = mk(Lt, n)
	rp := int32(r)
	for i := int32(0); i < n; i++ {
		o := I32(ap)
		n := I32(bp) - o
		if n < 0 {
			trap(Value)
		}
		SetI64(rp, int64(atv(rx(x), Add(Ki(o), seq(n)))))
		rp += 8
		ap += 4
		bp += 4
	}
	dx(a)
	dx(b)
	dx(x)
	return r
}
func split(x, y K) (r K) {
	xt, yt := tp(x), tp(y)
	xn := int32(1)
	if 16+xt != yt {
		if xt == Ct && yt == Ct {
			xn = nn(x)
			x = Find(x, rx(y))
		} else {
			trap(Nyi)
		}
	} else {
		x = Wer(Eql(x, rx(y)))
	}
	return rcut(y, Cat(Ki(0), Add(Ki(xn), rx(x))), cat1(x, Ki(nn(y))))
}
func join(x, y K) (r K) {
	xt := tp(x)
	if xt < 16 {
		x = Enl(x)
		xt = tp(x)
	}
	yt := tp(y)
	if yt != Lt {
		trap(Type)
	}
	yp := int32(y)
	yn := nn(y)
	r = mk(xt, 0)
	for i := int32(0); i < yn; i++ {
		v := x0(yp)
		if tp(v) != xt {
			trap(Type)
		}
		if i > 0 {
			r = ucat(r, rx(x))
		}
		r = ucat(r, v)
		yp += 8
	}
	dx(x)
	dx(y)
	return r
}

func Flr(x K) (r K) { // _x
	xt := tp(x)
	xp := int32(x)
	if xt < 16 {
		switch xt - 3 {
		case 0: // i
			return Kc(xp)
		case 1: // s
			return trap(Type)
		case 2: // f
			dx(x)
			return Ki(int32(F64floor(F64(xp))))
		case 3: // z
			dx(x)
			return Kf(F64(xp))
		default:
			return trap(Type)
		}
	}
	xn := nn(x)
	var rp int32
	switch xt - 19 {
	case 0: //I
		r = mk(Ct, xn)
		rp = int32(r)
		for i := int32(0); i < xn; i++ {
			SetI8(rp, I32(xp))
			xp += 4
		}
	case 1: //S
		trap(Type)
	case 2: //F
		r = mk(It, xn)
		rp = int32(r)
		for i := int32(0); i < xn; i++ {
			SetI32(rp, int32(F64floor(F64(xp))))
			xp += 8
			rp += 4
		}
	case 3: // Z
		r = mk(Ft, xn)
		rp = int32(r)
		for i := int32(0); i < xn; i++ {
			SetI64(rp, I64(xp))
			xp += 16
			rp += 8
		}
	case 4: // L
		return Ech(16, l1(x))
	default: // todo D/T
		trap(Type)
	}
	dx(x)
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
		x = Enl(x)
		t = tp(x)
	}
	var n, rp int32
	xn := nn(x)
	xp := int32(x)
	if t == Bt {
		n = sumb(xp, xn)
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
		n = sumi(xp, xn)
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
func firstWhere(x K) K { // *&x (todo idiom)
	t := tp(x)
	if t == Bt {
		dx(x)
		return Ki(maxi(0, firstWhereB(int32(x), nn(x))))
	}
	return Fst(Wer(x))
}
func firstWhereB(xp, n int32) int32 { // *&B
	e := xp + n
	ve := e &^ 15
	p := xp
f:
	for p < ve {
		if I8x16load(p).Any_true() != 0 {
			break f
		}
		p += 16
	}
	for p < e {
		if I8(p) != 0 {
			return p - xp
		}
		p++
	}
	return -1
}

func Typ(x K) (r K) { // @x
	dx(x)
	return sc(Enl(Kc(I8(520 + int32(tp(x))))))
}
func Val(x K) (r K) {
	xt := tp(x)
	if xt == st {
		return lup(x)
	}
	if xt == Ct {
		return val(x)
	}
	if xt == lf {
		xp := int32(x)
		r = cat1(l3(x0(xp), x1(xp), x0(xp+24)), Ki(nn(x))) // (code;locals;string;arity)
		dx(x)
		return r
	}
	if xt == Lt {
		return exec(x) // .L e.g. 1+2 is (1;2;`66)
	}
	if xt > Lt {
		r = x1(int32(x))
		dx(x)
		return r
	}
	trap(Nyi)
	return x
}
func val(x K) (r K) {
	s := src
	x = parse(x)
	xn := nn(x)
	xp := int32(x) + 8*(xn-1)
	a := int32(0)
	if xn > 2 && I64(xp) == 64 {
		a = 1
	}
	x = exec(x)
	dx(src)
	src = s
	if a != 0 {
		dx(x)
		return 0
	}
	return x
}
