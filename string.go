package main

import (
	. "github.com/ktye/wg/module"
)

func Kst(x K) K { return Atx(Ks(32), x) } // `k@
func Lst(x K) K { return Atx(Ks(40), x) } // `l@
func Str(x K) K {
	r := K(0)
	xt := tp(x)
	if xt > 16 {
		return Ech(17, l1(x))
	}
	xp := int32(x)
	if xt > dt {
		switch xt - cf {
		case 0: // cf
			rx(x)
			r = ucats(Rev(Str(K(xp) | K(Lt)<<59)))
		case 1: // df
			r = Str(x0(x))
			p := x1(x)
			if int32(p)%2 != 0 {
				p = cat1(Str(20+p), Kc(':'))
			} else {
				p = Str(21 + p)
			}
			r = ucat(r, p)
		case 2: //pf
			f := x0(x)
			l := x1(x)
			i := x2(x)
			ft := tp(f)
			f = Str(f)
			dx(i)
			if nn(i) == 1 && I32(int32(i)) == 1 && (ft == 0 || ft == df) {
				r = ucat(Kst(Fst(l)), f)
			} else {
				r = ucat(f, emb('[', ']', ndrop(-1, ndrop(1, Kst(l)))))
			}
		case 3: //lf
			r = x2(x)
		default: // native
			r = x1(x)
		}
		dx(x)
		return r
	} else {
		switch xt {
		case 0:
			if xp > 448 {
				return Str(K(xp) - 448)
			}
			ip := xp
			switch xp >> 6 {
			case 0: //  0..63  monadic
				if xp == 0 {
					return mk(Ct, 0)
				}
			case 1: // 64..127 dyadic
				ip -= 64
			case 2: // 128     dyadic indirect
				ip -= 128
			case 3: // 192     tetradic
				ip -= 192
				//default:
				//	return ucat(Ku('`'), si(xp))
			}
			if ip > 25 || ip == 0 {
				return ucat(Ku('`'), si(xp))
			}
			r = Ku(uint64(I8(227 + ip)))
		case 1: //not reached
			r = 0
		case ct:
			r = Ku(uint64(xp))
		case it:
			r = si(xp)
		case st:
			r = cs(x)
		case ft:
			r = sf(F64(xp))
		default:
			r = sfz(F64(xp), F64(xp+8))
		}
	}
	dx(x)
	return r
}
func emb(a, b int32, x K) K { return cat1(Cat(Kc(a), x), Kc(b)) }
func si(x int32) K {
	if x == 0 {
		return Ku(uint64('0'))
	} else if x == nai {
		return Ku(20016) // 0N
	} else if x < 0 {
		return ucat(Ku(uint64('-')), si(-x))
	}
	r := mk(Ct, 0)
	for x != 0 {
		r = cat1(r, Kc('0'+x%10))
		x /= 10
	}
	return Rev(r)
}
func sf(x float64) K {
	c := int32(0)
	if x != x {
		return Ku(28208) // 0n
	}
	u := uint64(I64reinterpret_f64(x))
	if u == uint64(0x7FF0000000000000) {
		return Ku(30512) // 0w
	} else if u == uint64(0xFFF0000000000000) {
		return Ku(7811117) // -0w
	}
	if x < 0 {
		return ucat(Ku(uint64('-')), sf(-x))
	}
	if x > 0 && (x >= 1e6 || x <= 1e-6) {
		return se(x)
	}
	r := mk(Ct, 0)
	i := int64(x)
	if i == 0 {
		r = cat1(r, Kc('0'))
	}
	for i != 0 {
		r = cat1(r, Kc(int32('0'+i%10)))
		i /= 10
	}

	r = Rev(r)
	r = cat1(r, Kc('.'))
	x -= F64floor(x)
	for i := int32(0); i < 6; i++ {
		x *= 10
		r = cat1(r, Kc('0'+(int32(x)%10)))
		continue
	}
	n := nn(r)
	rp := int32(r)
	for n > 0 {
		n--
		if I8(rp) == '0' {
			c++
		} else {
			c = 0
		}
		rp++
	}
	return ndrop(-c, r)
}
func se(x float64) K {
	f := x
	e := int64(0)
	if frexp1(x) != 0 {
		f = frexp2(x)
		e = frexp3(x)
	}
	x = 0.3010299956639812 * float64(e) // log10(2)*
	ei := int32(F64floor(x))
	x = x - float64(ei)
	return ucat(cat1(sf(f*pow(10.0, x)), Kc('e')), si(ei))
}
func sfz(re, im float64) K {
	if (re != re) || (im != im) {
		return Ku(6385200) // 0na
	}
	z := hypot(re, im)
	a := ang2(im, re)
	r := cat1(trdot(sf(z)), Kc('a'))
	if a != 0.0 {
		r = ucat(r, trdot(sf(a)))
	}
	return r
}
func trdot(x K) K {
	n := nn(x)
	if I8(int32(x)+n-1) == '.' {
		return ndrop(-1, x)
	}
	return x
}

func Cst(x, y K) K { // x$y
	yt := tp(y)
	if yt > Zt {
		return Ecr(17, l2(x, y))
	}
	if yt == ct {
		y, yt = Enl(y), Ct
	}
	if tp(x) != st || yt != Ct {
		trap() //type
	}
	if int32(x) == 0 { // `$"sym"
		return sc(y)
	}
	return prs(ts(x), y)
}
func prs(t T, y K) K { // s$C
	r := K(0)
	yp, yn := int32(y), nn(y)
	p, e := pp, pe
	pp = yp
	pe = yp + yn
	tt := t & 15
	if tt == 2 {
		if t == Ct {
			return y // `C$
		} else {
			return Fst(y) // `c$"x"
		}
	}
	if t == 4 {
		r = Fst(tsym()) // `s$"`a"
	} else if t > 2 && t <= 6 {
		r = tnum()
		if tp(r) < t && r != 0 {
			r = uptype(r, t) // `f$"1"
		}
	}
	if t > Ct && t < Lt {
		if pp == pe {
			r = mk(t, 0) // `I$"" -> !0
		} else {
			if t == 20 {
				r = tsym() // `S$"`a`b"
			} else {
				r = tnms()
				if tp(r)&15 < t&15 && r != 0 {
					r = uptype(r, t&15) // `F$"1 2"
				}
			}
			if tp(r) == t-16 {
				r = Enl(r) // `I$"1" -> ,1
			}
		}
	}
	if tp(r) != t || pp < pe {
		dx(r)
		r = 0
	}
	pp, pe = p, e
	dx(y)
	return r //0(parse error)
}
func ts(x K) T {
	c := int32(Fst(cs(x)))
	for i := int32(521); i < 546; i++ {
		if I8(i) == c {
			return T(i - 520)
		}
		continue
	}
	trap() //value
	return 0
}

func Rtp(y K, x K) K { // `c@ `i@ `s@ `f@ `z@ (reinterpret data)
	t := ts(y)
	xt := tp(x)
	t += T(16 * I32B(t < 16))
	if xt < 16 || t < 17 || t > Zt {
		trap() //type
	}
	n := nn(x) * sz(xt)
	s := sz(t)
	if n%s != 0 {
		trap() //length
	}
	x = use(x)
	SetI32(int32(x)-12, n/s)
	x = K(t)<<59 | K(int32(x))
	return x
}
