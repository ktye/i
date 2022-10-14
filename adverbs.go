package main

import (
	. "github.com/ktye/wg/module"
)

type rdf = func(K, int32, T, int32) K

func ech(x K) K { return l2t(x, 0, df) } // '
func ecp(x K) K { return l2t(x, 1, df) } // ':
func rdc(x K) K { return l2t(x, 2, df) } // /
func ecr(x K) K { return l2t(x, 3, df) } // /:
func scn(x K) K { return l2t(x, 4, df) } // \
func ecl(x K) K { return l2t(x, 5, df) } // \:

// maybe?
//  '     f1'x   f2' [x;y]   f3'[x;y;z]...   a' x bin  a'[x;y] lin
//  ':    f2':x  f2':[x;y]   err             a':x  in
//  /     f2/x   f2/ [x;y]   f3/[x;y;z]...   a/x  dec
//  \     f2\x   f2\ [x;y]   f3/[x;y;z]...   a\x  enc
//  /:    fix/:x f2/:[x;y]   err             a/:x join
//  \:    fix/:x f2/:[x;y]   err             a\:x split

func Ech(f, x K) (r K) {
	t := tp(f)
	if isfunc(t) == 0 {
		if nn(x) == 2 {
			x, r = spl2(x)
			return lin(x, f, r)
		}
		return Bin(f, Fst(x))
	}
	if nn(x) == 1 {
		x = Fst(x)
	} else {
		return ecn(f, x)
	}
	if tp(x) < 16 {
		trap(Type)
	}
	xt := tp(x)
	if xt == Dt {
		x, r = spl2(x)
		return Key(x, Ech(f, l1(r)))
	}
	if xt == Tt {
		x = explode(x)
	}
	xn := nn(x)
	r = mk(Lt, xn)
	rp := int32(r)
	for i := int32(0); i < xn; i++ {
		SetI64(rp, int64(Atx(rx(f), ati(rx(x), i))))
		rp += 8
	}
	dx(f)
	dx(x)
	return uf(r)
}
func ecn(f, x K) (r K) {
	if nn(x) == 2 {
		r, x = spl2(x)
		if r == 0 {
			return Ech(f, l1(x))
		}
		if tp(f) == 0 && int32(f) == 13 {
			if tp(r) == Tt && tp(x) == Tt { // T,'T (horcat)
				if nn(r) != nn(x) {
					trap(Length)
				}
				p, q := spl2(r)
				r, x = spl2(x)
				return key(Cat(p, r), Cat(q, x), Tt)
			}
		}
		return ec2(f, r, x)
	}
	return Ech(20, l2(f, Flp(x)))
}
func ec2(f, x, y K) (r K) {
	var t T
	r, t, x, y = dctypes(x, y)
	if r != 0 {
		return key(r, ec2(f, x, y), t)
	}
	n, _ := conform(x, y)
	switch n {
	case 0: // a-a
		return Cal(f, l2(x, y))
	case 1: // a-v
		n = nn(y)
	case 2: // v-v
		n = nn(x)
	default: // v-a
		n = nn(x)
	}
	r = mk(Lt, n)
	rp := int32(r)
	for i := int32(0); i < n; i++ {
		SetI64(rp, int64(Cal(rx(f), l2(ati(rx(x), i), ati(rx(y), i)))))
		rp += 8
	}
	dx(f)
	dx(x)
	dx(y)
	return uf(r)
}
func Ecp(f, x K) (r K) {
	t := tp(f)
	if isfunc(t) == 0 {
		if nn(x) != 1 {
			trap(Rank)
		}
		return In(f, Fst(x))
	}
	xn := nn(x)
	var y K
	var m int32
	if xn == 1 {
		x = Fst(x)
		y = Fst(rx(x)) // could be missing(xt)
		m = 1
	} else if xn == 2 {
		y, x = spl2(x)
	} else {
		trap(Rank)
	}
	xt := tp(x)
	if xt < 16 {
		trap(Type)
	}
	if xt > Lt {
		trap(Nyi)
	}
	xn = nn(x)
	if 1 > xn-m {
		dx(f)
		return x
	}

	yt := tp(y)
	if tp(f) == 0 && xt < Zt && yt == xt-16 {
		fp := int32(f)
		if fp > 1 && fp < 8 && (xt == It || xt == Ft) {
			return epx(fp, x, y, xn) // +-*%&| 234567
		}
		if fp == 11 {
			fp = 10 // ~ =
		}
		if fp > 9 && fp < 11 {
			return epc(fp, x, y, xn) // <>= (~)
		}
	}

	r = mk(Lt, xn)
	rp := int32(r)
	SetI64(rp, int64(cal(rx(f), l2(ati(rx(x), 0), y))))
	for i := int32(1); i < xn; i++ {
		rp += 8
		SetI64(rp, int64(cal(rx(f), l2(ati(rx(x), i), ati(rx(x), i-1)))))
	}
	dx(f)
	dx(x)
	return uf(r)
}

func Rdc(f, x K) (r K) { // x f/y   (x=0):f/y
	t := tp(f)
	if isfunc(t) == 0 {
		if nn(x) == 2 {
			r, x = spl2(x)
			if t == it && isfunc(tp(r)) != 0 {
				return ndo(int32(f), r, x)
			} else {
				trap(Type)
			}
		}
		return Mod(Fst(x), f)
	}
	if arity(f) > 2 {
		return rdn(f, x, 0)
	}
	if t == df { // x//y
		r = x0(int32(f))
		dx(f)
		return Dec(r, Fst(x))
	}
	var y K
	if xn := nn(x); xn == 1 {
		y = Fst(x)
		x = 0
	} else if xn == 2 {
		x, y = spl2(x)
	} else {
		trap(Rank)
	}
	yt := tp(y)
	if yt == Dt {
		y = Val(y)
		yt = tp(y)
	}
	if yt < 16 {
		if x == 0 {
			dx(f)
			return y
		} else {
			return cal(f, l2(x, y))
		}
	}

	yn := nn(y)
	xt := tp(x)
	if tp(f) == 0 {
		fp := int32(f)
		if fp > 1 && fp < 8 && (xt == 0 || yt == xt+16) { // sum,prd,min,max (reduce.go)
			if yt == Tt {
				return Ech(rdc(f), l2(x, Flp(y)))
			}
			r = Func[365+fp].(rdf)(x, int32(y), yt, yn)
			if r != 0 {
				dx(x)
				dx(y)
				return r
			}
		}
		if x == 0 && fp == 13 { // ,/
			if yt < Lt {
				return y
			}
			r = ucats(y)
			if r != 0 {
				return r
			}
		}
	}

	if yn == 0 {
		if x == 0 {
			return ov0(f, y)
		} else {
			dx(f)
			dx(y)
			return x
		}
	}

	i := int32(0)
	if x == 0 {
		x, i = ati(rx(y), 0), 1
	}
	for i < yn {
		x = cal(rx(f), l2(x, ati(rx(y), i)))
		i++
	}
	dx(y)
	dx(f)
	return x
}
func rdn(f, x, l K) (r K) { // {x+y*z}/x  {x+y*z}\x
	r = Fst(rx(x))
	x = Flp(ndrop(1, x))
	n := nn(x)
	for i := int32(0); i < n; i++ {
		r = Cal(rx(f), Cat(l1(r), ati(rx(x), i)))
		if l != 0 {
			l = cat1(l, rx(r))
		}
	}
	dx(f)
	dx(x)
	if l != 0 {
		dx(r)
		return uf(l)
	}
	return r
}

func Ecr(f, x K) (r K) { // f/:x   x f/:y   x/:y(join)
	t := tp(f)
	if isfunc(t) == 0 {
		if nn(x) != 1 {
			trap(Rank)
		}
		return join(f, Fst(x))
	}
	xn := nn(x)
	switch xn - 1 {
	case 0: // fixed-point
		return fix(f, Fst(x), 0)
	case 1:
		var y K
		x, y = spl2(x)
		yt := tp(y)
		if yt < 16 {
			return cal(f, l2(x, y))
		}
		if yt > Lt {
			var t T
			r, t, x, y = dctypes(x, y)
			return key(r, Ecr(f, l2(x, y)), t)
		}
		yn := nn(y)
		r = mk(Lt, yn)
		rp := int32(r)
		for i := int32(0); i < yn; i++ {
			SetI64(rp, int64(cal(rx(f), l2(rx(x), ati(rx(y), i)))))
			rp += 8
		}
		dx(f)
		dx(x)
		dx(y)
		return uf(r)
	default:
		return trap(Rank)
	}
}
func fix(f, x, l K) (r K) {
	y := rx(x)
	for {
		r = Atx(rx(f), rx(x))
		if match(r, x) != 0 {
			break
		}
		if match(r, y) != 0 {
			break
		}
		dx(x)
		x = r
		if l != 0 {
			l = cat1(l, rx(x))
		}
	}
	dx(f)
	dx(r)
	dx(y)
	if l != 0 {
		dx(x)
		return l
	}
	return x
}
func ndo(n int32, f, x K) K {
	for n > 0 {
		x = cal(rx(f), l1(x))
		n--
	}
	dx(f)
	return x
}

func Scn(f, x K) (r K) {
	t := tp(f)
	if isfunc(t) == 0 {
		if nn(x) != 1 {
			trap(Rank)
		}
		return Div(Fst(x), f)
	}
	if arity(f) > 2 {
		return rdn(f, x, mk(Lt, 0))
	}
	if t == df { // x\\y
		r = x0(int32(f))
		dx(f)
		return Enc(r, Fst(x))
	}
	//kdb:if int32(f)==29{trap(Err);}
	var y K
	if xn := nn(x); xn == 1 {
		y = Fst(x)
		x = 0
	} else if xn == 2 {
		x, y = spl2(x)
	} else {
		trap(Rank)
	}
	yt := tp(y)
	if yt < 16 {
		if x == 0 {
			dx(f)
			return y
		} else {
			return cal(f, l2(x, y))
		}
	}
	yn := nn(y)
	if yn == 0 {
		dx(f)
		dx(x)
		return y
	}
	if yt == Dt {
		r, y = spl2(y)
		return Key(r, Scn(f, l2(x, y)))
	}

	xt := tp(x)
	if tp(f) == 0 {
		fp := int32(f)
		if (fp == 2 || fp == 4) && (xt == 0 || yt == xt+16) { // sums,prds (reduce.go)
			if yt == Tt {
				return Flp(Ech(scn(f), l2(x, Flp(y)))) // +f\'[x;+y]
			}
			r = Func[372+fp].(rdf)(x, int32(y), yt, yn)
			if r != 0 {
				dx(x)
				dx(y)
				return r
			}
		}
	}

	r = mk(Lt, yn)
	rp := int32(r)
	i := int32(0)
	if x == 0 {
		x, i = ati(rx(y), 0), 1
		SetI64(rp, int64(rx(x)))
		rp += 8
	}
	for i < yn {
		x = cal(rx(f), l2(x, ati(rx(y), i)))
		SetI64(rp, int64(rx(x)))
		rp += 8
		i++
	}
	dx(y)
	dx(x)
	dx(f)
	return uf(r)
}
func Ecl(f, x K) (r K) { // f\:x   x f\:y   x\:y(split)
	t := tp(f)
	if isfunc(t) == 0 {
		if nn(x) != 1 {
			trap(Rank)
		}
		return split(f, Fst(x))
	}
	xn := nn(x)
	switch xn - 1 {
	case 0: // fixed-point-scan
		x = Fst(x)
		return fix(f, rx(x), Enl(x))
	case 1:
		var y K
		x, y = spl2(x)
		if tp(x) < 16 {
			return cal(f, l2(x, y))
		}
		xn := nn(x)
		r = mk(Lt, xn)
		rp := int32(r)
		for i := int32(0); i < xn; i++ {
			SetI64(rp, int64(cal(rx(f), l2(ati(rx(x), i), rx(y)))))
			rp += 8
		}
		dx(f)
		dx(x)
		dx(y)
		return uf(r)
	default:
		return trap(Rank)
	}
}

func uf(x K) (r K) {
	xn := nn(x)
	xp := int32(x)
	var rt T
	for i := int32(0); i < xn; i++ {
		t := tp(K(I64(xp)))
		if i == 0 {
			rt = t
		} else if t != rt {
			return x
		}
		xp += 8
	}
	if rt == Dt {
		return ufd(x)
	}
	if rt == 0 || rt > zt {
		return x
	}
	rt += 16
	r = mk(rt, xn)
	s := sz(rt)
	rp := int32(r)
	xp = int32(x)
	switch s >> 2 {
	case 0:
		for i := int32(0); i < xn; i++ {
			SetI8(rp, I32(xp))
			xp += 8
			rp++
		}
	case 1:
		for i := int32(0); i < xn; i++ {
			SetI32(rp, I32(xp))
			xp += 8
			rp += 4
		}
	case 2:
		for i := int32(0); i < xn; i++ {
			SetI64(rp, I64(I32(xp)))
			xp += 8
			rp += 8
		}
	default:
		for i := int32(0); i < xn; i++ {
			s := I32(xp)
			SetI64(rp, I64(s))
			SetI64(rp+8, I64(s+8))
			xp += 8
			rp += 16
		}
	}
	dx(x)
	return r
}
func ufd(x K) (r K) {
	xp := int32(x)
	r = Til(x0(xp))
	if tp(r) != St {
		dx(r)
		return x
	}
	n := nn(x)
	for i := int32(0); i < n; i++ {
		if match(r, K(I64(int32(I64(xp))))) == 0 {
			dx(r)
			return x
		}
		xp += 8
	}
	return key(r, Flp(Ech(20, l1(x))), Tt)
}

func ov0(f, x K) K {
	dx(f)
	dx(x)
	return missing(tp(x))
}
func epx(f int32, x, y K, n int32) (r K) { // ( +-*%&| )':
	xt := tp(x)
	xp := int32(x)
	r = mk(xt, n)
	rp := int32(r)
	f = 212 + 11*f
	yp := int32(y)
	if xt == It {
		SetI32(rp, Func[f].(f2i)(I32(xp), yp))
		for i := int32(1); i < n; i++ {
			xp += 4
			rp += 4
			SetI32(rp, Func[f].(f2i)(I32(xp), I32(xp-4)))
		}
	} else {
		f++
		SetF64(rp, Func[f].(f2f)(F64(xp), F64(yp)))
		for i := int32(1); i < n; i++ {
			xp += 8
			rp += 8
			SetF64(rp, Func[f].(f2f)(F64(xp), F64(xp-8)))
		}
	}
	dx(x)
	dx(y)
	return r
}
func epc(f int32, x, y K, n int32) (r K) { // ( <>= )':
	xt := tp(x)
	xp := int32(x)
	s := sz(xt)
	r = mk(It, n)
	rp := int32(r)
	f = 188 + 15*f
	switch s >> 2 {
	case 0:
		SetI32(rp, Func[f].(f2i)(I8(xp), int32(y)))
		for i := int32(1); i < n; i++ {
			xp++
			rp += 4
			SetI32(rp, Func[f].(f2i)(I8(xp), I8(xp-1)))
		}
	case 1:
		SetI32(rp, Func[f].(f2i)(I32(xp), int32(y)))
		for i := int32(1); i < n; i++ {
			xp += 4
			rp += 4
			SetI32(rp, Func[f].(f2i)(I32(xp), I32(xp-4)))
		}
	default:
		f++
		SetI32(rp, Func[f].(f2c)(F64(xp), F64(int32(y))))
		for i := int32(1); i < n; i++ {
			xp += 8
			rp += 4
			SetI32(rp, Func[f].(f2c)(F64(xp), F64(xp-8)))
		}
	}
	dx(x)
	dx(y)
	return r
}
