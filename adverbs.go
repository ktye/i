package main

import (
	. "github.com/ktye/wg/module"
)

func ech(x K) K { return l2t(x, 0, df) } // '
func ecp(x K) K { return l2t(x, 1, df) } // ':
func rdc(x K) K { return l2t(x, 2, df) } // /
func ecr(x K) K { return l2t(x, 3, df) } // /:
func scn(x K) K { return l2t(x, 4, df) } // \
func ecl(x K) K { return l2t(x, 5, df) } // \:

func Ech(f, x K) (r K) {
	//fmt.Printf("each f %s rc=%d\n", sK(f), rc(f))

	if nn(x) == 1 {
		x = Fst(x)
	} else {
		return ecn(f, x)
	}
	if tp(x) < 16 {
		x = Enl(x)
	}
	xt := tp(x)
	if xt == Dt {
		x, r = spl2(x)
		x = Out(x)
		r = Out(r)
		return Key(x, Ech(f, l1(r)))
	}
	if xt > Lt {
		panic(Nyi)
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
	//fmt.Printf("<ech f %s rc=%d\n", sK(f), rc(f))
	return uf(r)
}
func ecn(f, x K) (r K) {
	if nn(x) == 2 {
		r, x = spl2(x)
		return ec2(f, r, x)
	}
	return Ech(20, Flp(x))
}
func ec2(f, x, y K) (r K) {
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
	return r
}
func Ecp(f, x K) (r K) {
	xn := nn(x)
	var y K
	if xn == 1 {
		x = Fst(x)
		y = Fst(rx(x))
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
	if 2 > xn+I32B(y != 0) {
		dx(f)
		return x
	}

	yt := tp(y)
	if tp(f) == 0 && xt != Lt && yt == xt-16 {
		fp := int32(f)
		if fp > 2 && fp < 6 || fp == 7 || fp == 8 {
			return epx(fp, x, y, xn) // +-*% &|
		}
		if fp == 12 {
			fp = 11 // ~ =
		}
		if fp > 8 && fp < 12 {
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
func epx(f int32, x, y K, n int32) (r K) { // ( +-*% &| )':
	xt := tp(x)
	xp := int32(x)
	s := sz(xt)
	r = mk(xt, n)
	rp := int32(r)
	f = 212 + 12*f
	switch s >> 2 {
	case 0:
		SetI8(rp, Func[f].(f2i)(I32(xp), int32(y)))
		for i := int32(1); i < n; i++ {
			xp++
			rp++
			SetI8(rp, Func[f].(f2i)(I32(xp), I32(xp-1)))
		}
	case 1:
		SetI32(rp, Func[f].(f2i)(I32(xp), int32(y)))
		for i := int32(1); i < n; i++ {
			xp += 4
			rp += 4
			SetI32(rp, Func[f].(f2i)(I32(xp), I32(xp-4)))
		}
	case 2:
		f++
		SetF64(rp, Func[f].(f2f)(F64(xp), F64(int32(y))))
		for i := int32(1); i < n; i++ {
			xp += 8
			rp += 8
			SetF64(rp, Func[f].(f2f)(F64(xp), F64(xp-8)))
		}
	default:
		trap(Nyi)
	}
	dx(x)
	dx(y)
	return r
}
func epc(f int32, x, y K, n int32) (r K) { // ( <>= )':
	xt := tp(x)
	xp := int32(x)
	s := sz(xt)
	r = mk(Bt, n)
	rp := int32(r)
	f = 143 + 15*f
	switch s >> 2 {
	case 0:
		SetI8(rp, Func[f].(f2i)(I32(xp), int32(y)))
		for i := int32(1); i < n; i++ {
			xp++
			rp++
			SetI8(rp, Func[f].(f2i)(I32(xp), I32(xp-1)))
		}
	case 1:
		SetI8(rp, Func[f].(f2i)(I32(xp), int32(y)))
		for i := int32(1); i < n; i++ {
			xp += 4
			rp++
			SetI8(rp, Func[f].(f2i)(I32(xp), I32(xp-4)))
		}
	case 2:
		f++
		SetI8(rp, Func[f].(f2c)(F64(xp), F64(int32(y))))
		for i := int32(1); i < n; i++ {
			xp += 8
			rp++
			SetI8(rp, Func[f].(f2c)(F64(xp), F64(xp-8)))
		}
	default:
		trap(Nyi)
	}
	dx(x)
	dx(y)
	return r
}
func Rdc(f, x K) (r K) { // x f/y   (x=0):f/y
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
			return x
		} else {
			return cal(f, l2(x, y))
		}
	}

	xt := tp(x)
	yn := nn(y)
	if tp(f) == 0 {
		fp := int32(f)
		if fp > 1 && fp < 9 && (xt == 0 || yt == xt+16) { // sum,prd,min,max (reduce.go)
			r = Func[367+fp].(rdf)(x, int32(y), yt, yn)
			if r != 0 {
				dx(x)
				dx(y)
				return r
			}
		}
		if yt == Lt && x == 0 && fp == 13 {
			r = ucats(y)
			if r != 0 {
				return r
			}
		}
	}

	if yn == 0 {
		if x == 0 {
			return zero(yt - 16)
		} else {
			dx(f)
			dx(y)
			return x
		}
	}
	if yt > Lt {
		trap(Nyi)
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

func Ecr(f, x K) (r K) { // f/:x   x f/:y   x/:y(join)
	t := tp(f)
	if isfunc(t) == 0 {
		if nn(x) != 1 {
			trap(Rank)
		}
		return join(f, Fst(x))
	}
	xt := tp(x)
	if xt != Lt {
		trap(Type)
	}
	xn := nn(x)
	switch xn - 1 {
	case 0: // fixed-point
		return trap(Nyi)
	case 1:
		var y K
		x, y = spl2(x)
		if tp(y) < 16 {
			y = Enl(y)
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

func Scn(f, x K) (r K) {
	var y K
	if xn := nn(x); xn == 1 {
		y = Fst(x)
		x = 0
	} else if xn == 2 {
		x, y = spl2(x)
	} else {
		trap(Rank)
	}
	yn := nn(y)
	if yn == 0 {
		dx(f)
		dx(x)
		return y
	}
	yt := tp(y)
	if yt < 16 {
		if x == 0 {
			dx(f)
			return x
		} else {
			return cal(f, l2(x, y))
		}
	}

	xt := tp(x)
	if tp(f) == 0 {
		fp := int32(f)
		if fp > 1 && fp < 9 && (xt == 0 || yt == xt+16) { // sums,prds,mins,maxs (reduce.go)
			r = Func[374+fp].(rdf)(x, int32(y), yt, yn)
			if r != 0 {
				dx(x)
				dx(y)
				return r
			}
		}
	}

	if yt > Lt {
		trap(Nyi)
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
		return trap(Nyi)
	case 1:
		var y K
		x, y = spl2(x)
		if tp(x) < 16 {
			x = Enl(x)
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
