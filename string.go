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
		} else if xt != Ct {
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
			r = ucat(Kc(96), join(Kc('`'), x))
		case 4:
			r = join(Kc(' '), x)
		case 5:
			r = join(Kc(' '), x)
		case 6:
			r = emb(40, 41, join(Kc(';'), x))
		default:
			trap(Nyi)
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
		trap(Nyi)
	}
	return Ku(uint64(r))
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
			r = flat(Str(K(xp) | K(Lt)<<59))
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
			r = ucat(Str(x0(xp)), emb('[', ']', join(Kc(';'), Str(x1(xp)))))
		case 3: // lf
			r = x3(xp)
		}
		dx(x)
		return r
	} else {
		switch xt {
		case 0:
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
			trap(Nyi)
		default:
			trap(Err)
		}
	}
	dx(x)
	return r
}
func emb(a, b int32, x K) (r K) { return cat1(Cat(Kc(a), x), Kc(b)) }
func join(x K, y K) (r K) {
	yn := nn(y)
	yp := int32(y)
	r = mk(Ct, 0)
	for i := int32(0); i < yn; i++ {
		if i > 0 {
			r = cat1(r, x)
		}
		r = ucat(r, x0(yp))
		yp += 8
	}
	return r
}
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
	if x < 0 {
		return ucat(Ku(uint64('-')), sf(-x))
	}
	r = mk(Ct, 0)
	u := uint64(x)
	if u == 0 {
		r = cat1(r, Kc('0'))
	}
	for u > 0 {
		r = cat1(r, Kc(int32('0'+u%10)))
		u /= 10
	}
	r = Rev(r)
	r = cat1(r, '.')
	x -= float64(u)
	for i := int32(0); i < 6; i++ {
		x *= 10
		r = cat1(r, Kc('0'+(int32(x)%10)))
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
