package main

import (
	. "github.com/ktye/wg/module"
)

func nyi(x K) K { return trap(Nyi) }
func Idy(x K) K { return x } // :x
func Dex(x, y K) K { // x:y
	dx(x)
	return y
}
func Flp(x K) K { // +x
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
		return Key(x0(x), r1(x))
	default:
		return x
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
func Til(x K) (r K) {
	xt := tp(x)
	if xt > Lt {
		r = x0(x)
		dx(x)
		return r
	}
	if xt == it {
		return seq(int32(x))
	}
	if xt == It {
		return kx(120, x) // odo
	}
	return trap(Type)
}
func seq(n int32) (r K) {
	n = maxi(n, 0)
	r = mk(It, n)
	if n == 0 {
		return r
	}
	seqi(int32(r), ep(r))
	return r
}
func seqi(p, e int32) {
	i := int32(0)
	for p < e {
		SetI32(p, i)
		i++
		p += 4
		continue
	}
}
func Unq(x K) (r K) { // ?x
	xt := tp(x)
	if xt < 16 {
		return roll(x)
	}
	if xt >= Lt {
		if xt == Dt {
			trap(Type)
		}
		if xt == Tt {
			r = x0(x)
			x = r1(x)
			return key(r, Flp(Unq(Flp(x))), xt)
		}
		return kx(96, x) // .uqf
	}
	xn := nn(x)
	r = mk(xt, 0)
	for i := int32(0); i < xn; i++ {
		xi := ati(rx(x), i)
		if int32(In(rx(xi), rx(r))) == 0 {
			r = cat1(r, xi)
		} else {
			dx(xi)
		}
	}
	dx(x)
	return r
}
func Uqs(x K) K { // ?^x
	xt := tp(x)
	if xt < 16 {
		trap(Type)
	}
	return kx(88, x) // .uqs
}
func Grp(x K) K    { return kx(128, x) }                             // =x grp.
func grp(x, y K) K { return Atx(Drp(rx(x), rx(y)), Grp(Atx(y, x))) } // s?T
func Key(x, y K) K { return key(x, y, Dt) }                          // x!y
func key(x, y K, t T) (r K) { // Dt or Tt
	xt, yt := tp(x), tp(y)
	if xt < 16 || xt == Dt {
		if yt < 16 || yt == Dt {
			return Key(Enl(x), Enl(y))
		} else {
			if xt == st && yt == Tt {
				return keyt(x, y)
			}
			x = ntake(nn(y), x)
		}
	}
	xn := nn(x)
	if yt < 16 || yt == Dt {
		y = ntake(nn(x), y)
	}
	if xn != nn(y) {
		trap(Length)
	}
	if t == Tt {
		if xn > 0 {
			xn = nn(K(I64(int32(y))))
		}
	}
	r = l2(x, y)
	SetI32(int32(r)-12, xn)
	return K(int32(r)) | K(t)<<59
}
func keyt(x, y K) K { return Key(Tak(rx(x), rx(y)), Drp(x, y)) } // `s!t (key table: (`s#t)!`s_t)

func Tak(x, y K) (r K) { // x#y
	xt := tp(x)
	yt := tp(y)
	if yt == Dt {
		if xt == it {
			r = x0(y)
			y = r1(y)
			r = Tak(rx(x), r)
			y = Tak(x, y)
			return Key(r, y)
		} else {
			return Key(rx(x), Atx(y, x))
		}
	} else if yt == Tt {
		if xt&15 == st {
			if xt == st {
				x = Enl(x)
			}
			return key(rx(x), Atx(y, x), yt)
		} else {
			return Ecr(15, l2(x, y))
		}
	}
	if xt == it {
		return ntake(int32(x), y)
	}
	if xt > 16 && xt == yt {
		return atv(y, Wer(In(rx(y), x))) // set take
	}
	return Atx(rx(y), Wer(Cal(x, l1(y)))) // f#
}
func ntake(n int32, y K) (r K) {
	t := tp(y)
	if n == nai {
		if t < 16 {
			n = 1
		} else {
			n = nn(y)
		}
	}
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
	yt := tp(y)
	if yt > Lt {
		if yt == Dt || (yt == Tt && xt&15 == st) {
			r = x0(y)
			y = r1(y)
			if xt < 16 {
				x = Enl(x)
			}
			x = Wer(Not(In(rx(r), x)))
			return key(Atx(r, rx(x)), Atx(y, x), yt)
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
	return Atx(rx(y), Wer(Not(Cal(x, l1(y))))) // f#
}
func ndrop(n int32, y K) (r K) {
	yt := tp(y)
	if yt < 16 || yt > Lt {
		trap(Type)
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

func Cut(x, y K) (r K) { // x^y
	yt := tp(y)
	if yt == it || yt == ft {
		return Pow(y, x)
	}
	xt := tp(x)
	if xt == It {
		return cuts(x, y)
	}
	if xt == Ct && yt == Ct { // "set"^"abc"
		x = Wer(In(rx(y), x))
		return rcut(y, Cat(Ki(0), Add(Ki(1), rx(x))), Cat(x, Ki(nn(y))))
	}
	if xt != it || yt < 16 {
		trap(Type)
	}
	xp := int32(x)
	if xp <= 0 {
		xp = nn(y) / -xp
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
			trap(Type)
		}
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
		v := x0(K(yp))
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
func lin(x, y, z K) K { return cal(Val(Ks(112)), l3(x, y, z)) } // x y'z  (z.k: `".lin")
func Bin(x, y K) (r K) { // x'y
	xt := tp(x)
	yt := tp(y)
	if xt < 16 || xt > Ft { // n' win?
		if xt == it && yt > 16 {
			return win(int32(x), y)
		} else {
			return trap(Type)
		}
	}
	if xt == yt || yt == Lt {
		return Ecr(40, l2(x, y))
	} else if xt == yt+16 {
		r = Ki(ibin(x, y, xt))
	} else {
		trap(Type)
	}
	dx(x)
	dx(y)
	return r
}
func ibin(x, y K, t T) int32 {
	var k, j, h int32
	n := nn(x)
	xp := int32(x)
	yp := int32(y)
	j = n - 1
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
func win(n int32, x K) (r K) {
	y := seq(n)
	r = mk(Lt, 0)
	m := 1 + nn(x) - n
	for i := int32(0); i < m; i++ {
		r = ucat(r, l1(atv(rx(x), rx(y))))
		y = Add(Ki(1), y)
	}
	dx(x)
	dx(y)
	return r
}

func Flr(x K) (r K) { // _x
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
	var rp int32
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
		return K(int32(x)) | K(It)<<59
		//return Ech(16, l1(x))
	case 3: //F
		r = mk(It, xn)
		rp = int32(r)
		for i := int32(0); i < xn; i++ {
			SetI32(rp, int32(F64floor(F64(xp))))
			xp += 8
			rp += 4
		}
	case 4: // Z
		r = mk(Ft, xn)
		rp = int32(r)
		for i := int32(0); i < xn; i++ {
			SetI64(rp, I64(xp))
			xp += 16
			rp += 8
		}
	default: // L/D/T
		return Ech(16, l1(x))
	}
	dx(x)
	return r
}
func lower(x K) (r K) {
	r = use(x)
	p := int32(r)
	e := p + nn(r)
	for p < e {
		SetI8(p, lc(I8(p)))
		p++
	}
	return r
}
func lc(x int32) int32 {
	if x >= 'A' && x <= 'Z' {
		return x + 32
	} else {
		return x
	}
}

func Rev(x K) (r K) { // |x
	t := tp(x)
	if t < 16 {
		return x
	}
	if t == Dt {
		return Key(Rev(x0(x)), Rev(r1(x)))
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
	if t == Dt {
		return Atx(x0(x), Wer(r1(x)))
	}
	var n, rp int32
	xn := nn(x)
	xp := int32(x)
	if t == It {
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
	} else if xn == 0 {
		r = mk(It, 0)
	} else {
		trap(Type)
	}
	dx(x)
	return r
}
func Fwh(x K) K { // *&x
	t := tp(x)
	if t == It {
		dx(x)
		return Ki(fwh(int32(x), nn(x)))
	}
	return Fst(Wer(x))
}
func fwh(xp, n int32) int32 { // *&I
	p := xp
	e := xp + 4*n
	for p < e {
		if I8(p) != 0 {
			return (p - xp) >> 2
		}
		p += 4
	}
	return nai
}

func Typ(x K) K { // @x
	dx(x)
	return sc(Enl(Kc(I8(520 + int32(tp(x))))))
}
func Tok(x K) K { // `t@"src"
	if tp(x) == Ct {
		return tok(x)
	} else {
		return x
	}
}
func Val(x K) (r K) {
	xt := tp(x)
	if xt == st {
		return lup(x)
	}
	if xt == Ct {
		return val(x)
	}
	if xt == lf || xt == xf { // lambda: (code;locals;string;arity)
		//xp := int32(x)  // native: (ptr;string;arity)
		r = l2(x0(x), x1(x))
		if xt == lf {
			r = cat1(r, x2(x))
		}
		r = cat1(r, Ki(nn(x)))
		dx(x)
		return r
	}
	if xt == Lt {
		return exec(x) // .L e.g. 1+2 is (1;2;`66)
	}
	if xt > Lt {
		r = x1(x)
		dx(x)
		return r
	} else {
		return trap(Type)
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
func Fin(x K) K { return Wer(Not(Xpt(x))) } // find x (index of non-missing)
func Xpt(x K) K { // fill x
	xt := tp(x)
	if xt >= Lt {
		return Ech(38, l1(x))
	} else {
		return Eql(missing(xt-T(16*I32B(xt > 16))), x)
	}
}
func Fil(x, y K) (r K) { // x fill y
	xt, yt := tp(x), tp(y)
	if yt >= Lt {
		return Ecr(38, l2(x, y))
	}
	if xt == yt && xt < 16 {
		if int32(Xpt(y)) != 0 {
			return x
		} else {
			return y
		}
	}
	r = Wer(Xpt(rx(y)))
	if xt == yt-16 {
		return Amd(y, r, 1, x)
	}
	if xt == yt {
		return stv(y, r, atv(x, rx(r)))
	} else {
		return trap(Type)
	}
}
func Enc(x, y K) (r K) { // x\\y
	xt := tp(x)
	n := int32(0)
	if xt == It {
		n = nn(x)
	}
	r = mk(It, 0)
	yn := int32(Cnt(rx(y)))
l:
	for {
		n--
		xi := ati(rx(x), n)
		r = Cat(r, Enl(idiv(rx(y), xi, 1)))
		y = idiv(y, xi, 0)
		if n == 0 || (n < 0 && int32(y) == 0) {
			break
		}
		if tp(y) > 16 && n < 0 {
			if sumi(int32(y), yn) == 0 {
				break l
			}
		}
	}
	dx(x)
	dx(y)
	return Rev(r)
}
func Dec(x, y K) (r K) { // x//y   {z+x*y}/[0;x;y]
	if tp(y) < 16 {
		trap(Type)
	}
	r = Fst(rx(y))
	n := nn(y)
	for i := int32(1); i < n; i++ {
		r = Add(ati(rx(y), i), Mul(ati(rx(x), i), r))
	}
	dx(x)
	dx(y)
	return r
}
func sumi(xp, xn int32) (r int32) {
	e := xp + 4*xn
	for xp < e {
		r += I32(xp)
		xp += 4
	}
	return r
}
