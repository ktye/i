package main

import (
	. "github.com/ktye/wg/module"
)

func Kst(x K) (r K) {
	xt := tp(x)
	if xt < 16 {
		r = Str(x)
		if xt == ct {
			r = emb(34, 34, r)
		} else if xt == st {
			r = ucat(Ku(96), r)
		}
	} else {
		xn := nn(x)
		if xn == 0 {
			dx(x)
			return kst0(xt - 17)
		}
		if xt == Lt {
			x = Ech(28, l1(x)) // Kst
		} else if xt < Lt && xt != Ct {
			x = Str(x)
		}
		switch xt - 17 {
		case 0:
			r = cat1(Ech(4, l1(x)), Kc('b'))
		case 1:
			r = emb(34, 34, x)
		case 2:
			r = join(Kc(' '), x)
		case 3:
			r = ucat(Ku(96), join(Kc('`'), x))
		case 4:
			r = join(Kc(' '), x)
		case 5:
			r = join(Kc(' '), x)
		case 6:
			if xn == 1 {
				r = Fst(x)
			} else {
				r = emb(40, 41, join(Kc(';'), x))
			}
			r = r
		case 7: // Dt
			x, r = spl2(x)
			x = Kst(x)
			if xn == 1 {
				x = emb(40, 41, x)
				xn = 0
			}
			r = ucat(cat1(x, Kc('!')), Kst(r)) // todo ()!..
		default:
			xn = 0
			r = ucat(Ku(43), Kst(Flp(x)))
		}
		if xn == 1 {
			r = ucat(Ku(44), r)
		}
	}
	return r
}
func kst0(t T) (r K) {
	switch t {
	case 0:
		r = 1647321904 // 0#0b
	case 1:
		r = 8738 // ""
	case 2:
		r = 12321 // !0
	case 3:
		r = 6300464 // 0#`
	case 4:
		r = 774906672 // 0#0.
	case 5:
		r = 1630544688 // 0#0a
	case 6:
		r = 10536 // ()
	default:
		r = trap(Nyi)
	}
	return Ku(uint64(r))
}
func Lst(x K) (r K) { // `l@  matrix-output (list-of-chars)
	xt := tp(x)
	if xt < Lt {
		return Str(x)
	}
	n := nn(x)
	switch xt - Lt {
	case 0: // Lt
		r = mk(Lt, n)
		for i := int32(0); i < n; i++ {
			xi := x0(int32(x) + 8*i)
			ti := tp(xi)
			if ti == ct {
				xi = Enl(xi)
			} else if ti != Ct {
				xi = Kst(xi)
			}
			SetI64(int32(r)+8*i, int64(xi))
		}
		return r
	case 1: // Dt
		return trap(Nyi)
	default: // Tt
		return trap(Nyi)
	}
}
func Str(x K) (r K) {
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
			r = Str(x0(xp))
			p := x1(xp)
			if p%2 != 0 {
				p = cat1(Str(20+p), Kc(':'))
			} else {
				p = Str(21 + p)
			}
			r = ucat(r, p)
		case 2: //pf
			f, l, i := spl3(rx(x))
			f = Str(f)
			dx(i)
			if nn(i) == 1 && I32(int32(i)) == 1 {
				r = ucat(Str(Fst(l)), f)
			} else {
				r = ucat(f, emb('[', ']', join(Kc(';'), Str(l))))
			}
		case 3: // lf
			r = x3(xp)
		}
		dx(x)
		return r
	} else {
		switch xt {
		case 0:
			if xp > 448 {
				return Str(K(xp) - 448)
			}
			switch xp >> 6 {
			case 0: //  0..63  monadic
			case 1: // 64..127 dyadic
				xp -= 64
			case 2: // 128     dyadic indirect
				xp -= 128
			case 3: // 192     tetradic
				xp -= 192
			default:
				return ucat(Ku('`'), Ki(xp))
			}
			r = Ku(uint64(I8(227 + xp)))
		case bt:
			r = Ku(uint64(25136 + xp)) // 0b 1b
		case ct:
			r = Ku(uint64(xp))
		case it:
			r = si(xp)
		case st:
			r = cs(x)
		case ft:
			r = sf(F64(xp))
		case zt:
			r = sfz(F64(xp), F64(xp+8))
		default:
			r = trap(Err)
		}
	}
	dx(x)
	return r
}
func emb(a, b int32, x K) (r K) { return cat1(Cat(Kc(a), x), Kc(b)) }
func si(x int32) (r K) {
	if x == 0 {
		return Ku(uint64('0'))
	} else if x < 0 {
		return ucat(Ku(uint64('-')), si(-x))
	}
	r = mk(Ct, 0)
	for x != 0 {
		r = cat1(r, Kc('0'+x%10))
		x /= 10
	}
	return Rev(r)
}
func sf(x float64) (r K) {
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
	r = mk(Ct, 0)
	u = uint64(x)
	if u == 0 {
		r = cat1(r, Kc('0'))
	}
	for u > 0 {
		r = cat1(r, Kc(int32('0'+u%10)))
		u /= 10
	}
	r = Rev(r)
	r = cat1(r, '.')
	x -= F64floor(x)
	for i := int32(0); i < 6; i++ {
		x *= 10
		r = cat1(r, Kc('0'+(int32(x)%10)))
		continue
	}
	n := nn(r)
	rp := int32(r)
	var c int32
	for i := int32(0); i < n; i++ {
		if I8(rp) == '0' {
			c++
		} else {
			c = 0
		}
		rp++
	}
	return ndrop(-c, r)
}
func sfz(re, im float64) (r K) {
	z := hypot(re, im)
	a := ang2(im, re)
	r = cat1(trdot(sf(z)), Kc('a'))
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

func Cst(x, y K) (r K) { // x$y
	yt := tp(y)
	if yt == ct {
		y, yt = Enl(y), Ct
	}
	if yt == Ct {
		if tp(x) != st {
			trap(Type)
		}
		if int32(x) == 0 { // `$"sym"
			return sc(y)
		}
		return prs(ts(x), y)
	}
	return trap(Nyi) // todo conversions
}
func prs(t T, y K) (r K) { // s$C
	yp, yn := int32(y), nn(y)
	p, e := pp, pe
	pp = yp
	pe = yp + yn
	tt := t & 15
	if tt == 1 {
		r = tbln()
	}
	if tt == 2 {
		if t == Ct {
			return y // `C$
		} else {
			return Fst(y) // `c$"x"
		}
	}
	if t == 4 {
		r = Fst(tsym()) // `s$"`a"
	} else if t > 2 && t < 6 {
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
	for i := int32(520); i < 546; i++ {
		if I8(i) == c {
			return T(i - 520)
		}
		continue
	}
	trap(Value)
	return 0
}
