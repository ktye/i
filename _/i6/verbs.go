package main

import (
	. "github.com/ktye/wg/module"
)

func nyi(x K) K { trap(); return 0 }
func Idy(x K) K { return x } // :x
func Dex(x, y K) K { // x:y
	dx(x)
	return y
}
func Flp(x K) K { // +x
	xt := tp(x)
	if xt == Lt {
		n := nn(x)
		xp := int32(x)
		m := Ki(maxcount(xp, n))
		x = Atx(Rdc(13, l1(Ecr(15, l2(m, x)))), Ecl(2, l2(Til(m), Mul(m, Til(Ki(n))))))
	} else if xt > Lt {
		r := x0(x)
		x = r1(x)
		if xt == Tt {
			x = Key(r, x)
		} else {
			if tp(r) != St || tp(x) != Lt {
				trap() //type
			}
			m := maxcount(int32(x), nn(x))
			x = Ech(15, l2(Ki(m), x)) // (|/#'x)#'x
			r = l2(r, x)
			SetI32(int32(r)-12, m)
			x = ti(Tt, int32(r))
		}
	}
	return x
}
func maxcount(xp int32, n int32) int32 { // |/#l
	r := int32(0)
	for n > 0 {
		n--
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
	if t == Dt {
		return Fst(Val(x))
	}
	return ati(x, 0)
}
func Las(x K) K { // *|x
	t := tp(x)
	if t < 16 {
		return x
	}
	if t == Dt {
		x = Val(x)
	}
	n := nn(x)
	if n == 0 {
		return Fst(x)
	}
	return ati(x, n-1)
}

func Cnt(x K) K { // #x
	t := tp(x)
	dx(x)
	if t < 16 {
		return Ki(1)
	}
	return Ki(nn(x))
}
func Not(x K) K { // ~x
	if tp(x)&15 == st {
		x = Eql(Ks(0), x)
	} else {
		x = Eql(Ki(0), x)
	}
	return x
}
func Til(x K) K {
	xt := tp(x)
	if xt > Lt {
		t := x0(x)
		dx(x)
		return t
	}
	if xt == it {
		return seq(int32(x))
	}
	if xt == It {
		return odo(x) // {x\!*/x}
	}
	trap() //type
	return 0
}
func odo(x K) K {
	y := Scn(4, Enl(rx(x)))
	n := maxi(0, int32(Las(rx(y))))
	m := nn(x)
	r := mk(Lt, m)
	i := int32(0)
	m *= 4
	for i < m {
		SetI64(int32(r)+i<<1, int64(odo1(I32(int32(x)+i), divi(n, I32(int32(y)+i)), n)))
		i += 4
	}
	dx(x)
	dx(y)
	return r
}
func odo1(a, b, n int32) K {
	r := mk(It, n)
	k := int32(0)
	i := int32(0)
	p := int32(r)
	e := ep(r)
	for p < e {
		SetI32(p, k)
		i++
		if i == b {
			i = 0
			k++
			if k == a {
				k = 0
			}
		}
		p += 4
	}
	return r
}
func Unq(x K) K { // ?x
	var r K
	xt := tp(x)
	if xt < 16 {
		return roll(x)
	}
	if xt >= Lt {
		if xt == Dt {
			trap() //type
		}
		if xt == Tt {
			r = x0(x)
			x = r1(x)
			return key(r, Flp(Unq(Flp(x))), xt)
		}
	}
	rx(rx(x))
	return atv(x, Wer(Eql(seq(nn(x)), Fnd(x, x)))) // x@&(!#x)==x?x
}
func Grp(x K) K {
	if tp(x) == it {
		return mat(10, x)
	}
	return kx(96, x) //= grp.
}
func mat(f, x K) K {
	x = rx(seq(int32(x)))
	return Ecr(f, l2(x, x))
}
func Key(x, y K) K { return key(x, y, Dt) } // x!y
func key(x, y K, t T) K { // Dt or Tt
	xt := tp(x)
	yt := tp(y)
	if xt < 16 {
		if xt == it {
			return Mod(y, x)
		}
		if xt == st {
			if yt == Tt { // s!t (key table)
				x = rx(x)
				y = rx(y)
				return Key(Tak(x, y), Drp(x, y))
			}
		}
		x = Enl(x) //allow `a!,1 2 3 short for (`a)!,1 2 3
	}
	xn := nn(x)
	if t == Tt {
		if xn > 0 {
			xn = nn(K(I64(int32(y))))
		}
	} else if yt < 16 {
		trap() //type
	} else if xn != nn(y) {
		trap() //length
	}
	x = l2(x, y)
	SetI32(int32(x)-12, xn)
	return ti(t, int32(x))
}
func Tak(x, y K) K { // x#y
	xt := tp(x)
	yt := tp(y)
	if yt == Dt {
		x = rx(x)
		if xt == it {
			r := x0(y)
			y = r1(y)
			r = Tak(x, r)
			y = Tak(x, y)
			return Key(r, y)
		} else {
			return Key(x, Atx(y, x))
		}
	} else if yt == Tt {
		if xt&15 == st {
			if xt == st {
				x = Enl(x)
			}
			x = rx(x)
			return key(x, Atx(y, x), yt)
		} else {
			return Ecr(15, l2(x, y))
		}
	}
	if xt == it {
		return ntake(int32(x), y)
	}
	y = rx(y)
	if xt > 16 && xt == yt {
		return atv(y, Wer(In(y, x))) // set take
	}
	return Atx(y, Wer(Cal(x, l1(y)))) // f#
}
func ntake(n int32, y K) K {
	var r K
	t := tp(y)
	if n < 0 {
		if tp(y) < 16 {
			return ntake(-n, y)
		}
		n += nn(y)
		if n < 0 {
			return ucat(ntake(-n, missing(t-16)), y)
		}
		return ndrop(n, y)
	}
	if t < 16 {
		return atv(Enl(y), Wer(Ki(n)))
	}
	yn := nn(y)
	s := sz(t)
	yp := int32(y)
	if I32(yp-4) == 1 && bucket(s*yn) == bucket(s*n) && n <= yn && t < Lt {
		SetI32(yp-12, n)
		return y
	}
	if n > yn {
		r = seqm(n, yn)
	} else {
		r = seq(n)
	}
	return atv(y, r)
}
func seqm(n, m int32) K {
	r := mk(It, n)
	k := int32(0)
	p := int32(r)
	e := p + n<<2
	for p < e {
		SetI32(p, k)
		k++
		if k == m {
			k = 0
		}
		p += 4
	}
	return r
}
func Drp(x, y K) K { // x_y
	xt := tp(x)
	yt := tp(y)
	if yt > Lt {
		if yt == Dt || (yt == Tt && xt&15 == st) {
			r := x0(y)
			y = r1(y)
			if xt < 16 {
				x = Enl(x)
			}
			x = rx(Wer(Not(In(rx(r), x))))
			return key(Atx(r, x), Atx(y, x), yt)
		} else {
			return Ecr(16, l2(x, y))
		}
	}
	if xt == it {
		return ndrop(int32(x), y)
	}
	if xt > 16 && xt == yt {
		return atv(y, Wer(Not(In(rx(y), x)))) // set drop
	}
	if yt == it {
		return atv(x, Wer(Not(Eql(y, seq(nn(x))))))
	}
	return Atx(y, Wer(Not(Cal(x, l1(rx(y)))))) // f#
}
func ndrop(n int32, y K) K {
	var r K
	yt := tp(y)
	if yt < 16 || yt > Lt {
		trap() //type
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
		r = uf(r)
	}
	dx(y)
	return r
}

func Cut(x, y K) K { // x^y
	yt := tp(y)
	if yt == it || yt == ft {
		return Pow(y, x)
	}
	xt := tp(x)
	if xt == It {
		return cuts(x, y)
	}
	if xt == Ct && yt == Ct { // "set"^"abc"
		x = rx(Wer(In(rx(y), x)))
		return rcut(y, Cat(Ki(0), Add(Ki(1), x)), Cat(x, Ki(nn(y))))
	}
	if xt != it || yt < 16 {
		trap() //type
	}
	xp := int32(x)
	if xp <= 0 {
		xp = nn(y) / -xp
	}
	r := mk(Lt, xp)
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
	dxy(x, y)
	return r
}
func cuts(x, y K) K { return rcut(y, x, cat1(ndrop(1, rx(x)), Ki(nn(y)))) }
func rcut(x, a, b K) K { // a, b start-stop ranges
	n := nn(a)
	ap, bp := int32(a), int32(b)
	r := mk(Lt, n)
	rp := int32(r)
	for n > 0 {
		n--
		o := I32(ap)
		m := I32(bp) - o
		if m < 0 {
			trap() //value
		}
		SetI64(rp, int64(atv(rx(x), Add(Ki(o), seq(m)))))
		rp += 8
		ap += 4
		bp += 4
	}
	dxy(a, b)
	dx(x)
	return r
}
func split(x, y K) K {
	xt, yt := tp(x), tp(y)
	xn := int32(1)
	if yt == xt+16 {
		x = Wer(Eql(x, rx(y)))
	} else {
		if xt == yt && xt == Ct {
			xn = nn(x)
			x = Find(x, rx(y))
		} else {
			trap() //type
		}
	}
	x = rx(x)
	return rcut(y, Cat(Ki(0), Add(Ki(xn), x)), cat1(x, Ki(nn(y))))
}
func join(x, y K) K { // {(-#x)_,/y,\x}
	n := -int32(Cnt(rx(x)))
	return ndrop(n, Rdc(13, l1(Ecl(13, l2(y, x)))))
}
func Bin(x, y K) K { // x'y
	var r K
	xt := tp(x)
	yt := tp(y)
	if xt == yt || yt == Lt {
		return Ecr(40, l2(x, y))
	} else if xt == yt+16 {
		r = Ki(ibin(x, y, xt))
	} else {
		trap() //type
	}
	dxy(x, y)
	return r
}
func ibin(x, y K, t T) int32 {
	var h int32
	k := int32(0)
	n := nn(x)
	xp := int32(x)
	yp := int32(y)
	j := n - 1
	s := sz(t)
	switch s >> 2 {
	case 0:
		for {
			if k > j {
				return k - 1
			}
			h = (k + j) >> 1
			if I8(xp+h) > yp {
				j = h - 1
			} else {
				k = h + 1
			}
		}
	case 1:
		for {
			if k > j {
				return k - 1
			}
			h = (k + j) >> 1
			if I32(xp+4*h) > yp {
				j = h - 1
			} else {
				k = h + 1
			}
		}
	default:
		f := F64(yp)
		for {
			if k > j {
				return k - 1
			}
			h = (k + j) >> 1
			if F64(xp+8*h) > f {
				j = h - 1
			} else {
				k = h + 1
			}
		}
	}
	return 0 // not reached
}
func Flr(x K) K { // _x
	var r K
	rp := int32(0)
	xt := tp(x)
	xp := int32(x)

	if xt < 16 {
		switch xt - 2 {
		case 0: // c
			return Kc(lc(xp))
		case 1: // i
			return Kc(xp)
		case 2: // s
			return Ki(int32(xp))
		case 3: // f
			dx(x)
			return Ki(int32(F64floor(F64(xp))))
		case 4: // z
			dx(x)
			return Kf(F64(xp))
		default:
			return x
		}
	}
	xn := nn(x)
	switch xt - 18 {
	case 0: //C
		return lower(x)
	case 1: //I
		r = mk(Ct, xn)
		rp = int32(r)
		e := rp + xn
		for rp < e {
			SetI8(rp, I32(xp))
			xp += 4
			rp++
		}
	case 2: //S
		x = use(x)
		return ti(It, int32(x))
	case 3: //F
		r = mk(It, xn)
		rp = int32(r)
		for xn > 0 {
			xn--
			SetI32(rp, int32(F64floor(F64(xp))))
			xp += 8
			rp += 4
		}
	case 4: // Z
		return reim(x, 0)
	default: // L/D/T
		return Ech(16, l1(x))
	}
	dx(x)
	return r
}
func reim(x K, o int32) K {
	var r K
	if tp(x) < 16 {
		r = Kf(F64(int32(x) + o))
	} else {
		r = mk(Ft, nn(x))
		p := int32(r)
		o += int32(x)
		e := ep(r)
		for p < e {
			SetF64(p, F64(o))
			o += 16
			p += 8
		}
	}
	dx(x)
	return r
}
func lower(x K) K {
	x = use(x)
	p := int32(x)
	e := p + nn(x)
	for p < e {
		SetI8(p, lc(I8(p)))
		p++
	}
	return x
}
func lc(x int32) int32 { return x + 32*I32B(uint32(x-65) < 26) }

func Rev(x K) K { // |x
	var r K
	t := tp(x)
	if t < 16 {
		if t == it {
			return Rev(Grp(x)) //antidiag
		}
		trap()
	}
	if t == Dt {
		r = x0(x)
		return Key(Rev(r), Rev(r1(x)))
	}
	xn := nn(x)
	if xn < 2 {
		return x
	}
	r = mk(It, xn)
	rp := int32(r)
	for xn > 0 {
		xn--
		SetI32(rp, xn)
		rp += 4
	}
	return atv(x, r)
}

func Wer(x K) K { // &x
	r := K(0)
	t := tp(x)
	if t < 16 {
		x = Enl(x)
		t = tp(x)
	}
	if t == Dt {
		r = x0(x)
		return Atx(r, Wer(r1(x)))
	}
	xn := nn(x)
	xp := int32(x)
	if t == It {
		n := sumi(xp, ep(x))
		r = mk(It, n)
		rp := int32(r)
		for i := int32(0); i < xn; i++ {
			j := I32(xp)
			for j > 0 {
				j--
				SetI32(rp, i)
				rp += 4
			}
			xp += 4
		}
	} else if xn == 0 {
		r = mk(It, 0)
	} else {
		trap() //type
	}
	dx(x)
	return r
}
func Fwh(x K) K { // *&x
	t := tp(x)
	if t == It {
		dx(x)
		p := int32(x)
		e := ep(x)
		for p < e {
			if I32(p) != 0 {
				return Ki((p - int32(x)) >> 2)
			}
			p += 4
		}
		return Ki(0)
	}
	return Fst(Wer(x))
}
func Typ(x K) K { // @x
	dx(x)
	return sc(Ku(uint64(I8(253 + int32(tp(x))))))
}
func Tok(x K) K { // `t@"src"
	if tp(x) == Ct {
		return tok(x)
	} else {
		return x
	}
}
func Val(x K) K {
	xt := tp(x)
	if xt == st {
		return lup(x)
	}
	if xt == Ct {
		return val(x)
	}
	if xt&15 == zt {
		rx(x)
		return ucat(Enl(reim(x, 0)), Enl(reim(x, 8)))
	}
	if xt == lf || xt == xf { // lambda: (code;string;locals;arity)
		//xp := int32(x)  // native: (ptr;string;arity)
		r := l2(x0(x), x1(x))
		if xt == lf {
			r = cat1(r, x2(x))
		}
		dx(x)
		return cat1(r, Ki(nn(x)))
	}
	if xt == Lt {
		return exec(x) // .L e.g. 1+2 is (1;2;`66)
	}
	if xt > Lt {
		return r1(x)
	} else {
		trap() //type
		return 0
	}
}
func val(x K) K {
	x = parse(tok(x))
	xn := nn(x)
	xp := int32(x) + 8*(xn-1)
	a := int32(0)
	if xn > 2 && I64(xp) == 64 {
		a = 1
	}
	x = exec(x)
	if a != 0 {
		dx(x)
		return 0
	}
	return x
}
func Enc(x, y K) K {
	yt := tp(y)
	if yt == It {
		return cal(lup(Ks(104)), l2(x, y))
	}
	if yt != it {
		trap()
	}
	yi := int32(y)
	n := int32(0)
	if tp(x) == It {
		n = nn(x)
	}
	r := mk(It, 0)
	for {
		n--
		xi := int32(ati(rx(x), n))
		r = cat1(r, Ki(modi(yi, xi)))
		yi = divi(yi, xi)
		if n == 0 {
			break
		}
		if n < 0 && uint32(yi+1) < 2 {
			if yi == -1 {
				r = cat1(r, Ki(-1))
			}
			break
		}
	}
	dx(x)
	return Rev(r)
}
func Dec(x, y K) K { // x//y   {z+x*y}/[0;x;y]
	if tp(y) < 16 {
		trap() //type
	}
	r := Fst(rx(y))
	n := nn(y)
	for i := int32(1); i < n; i++ {
		r = Add(ati(rx(y), i), Mul(ati(rx(x), i), r))
	}
	dxy(x, y)
	return r
}
