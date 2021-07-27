package main

import (
	. "github.com/ktye/wg/module"
)

type ftok = func() K

func tok(x K) (r K) {
	var y K
	x = cmt(x)
	src = x
	pp = int32(x)
	pe = pp + nn(x)
	p := pp // srcp > 0
	r = mk(Lt, 0)
	for {
		ws()
		if pp == pe {
			break
		}
		for i := int32(192); i < 199; i++ { // tbln, tchr, tnms, tvrb, tpct, tvar, tsym
			y = Func[i].(ftok)()
			if y != 0 {
				y |= K(int64(pp-p) << 32)
				r = cat1(r, y)
				break
			}
			if i == 198 { // todo last-1
				trap(Parse)
			}
		}
	}
	return r
}
func tbln() (r K) {
	n := pe - pp
	for i := int32(0); i < n; i++ {
		c := I8(pp + i)
		if c != '0' && c != '1' {
			if i < 1 || c != 'b' {
				return 0
			}
			return pbln(i)
		}
	}
	return 0
}
func pbln(n int32) (r K) {
	r = mk(Bt, n)
	rp := int32(r)
	for i := int32(0); i < n; i++ {
		SetI8(rp, I8(pp+i)-'0')
		rp++
	}
	pp += 1 + n
	if n == 1 {
		return Fst(r)
	}
	return r
}
func tchr() (r K) {
	if I8(pp) == '0' && pp < pe { // 0x01ab (lower case only)
		if I8(1+pp) == 'x' {
			pp += 2
			return thex()
		}
	}
	if I8(pp) != 34 {
		return 0
	}
	pp++
	r = mk(Ct, 0)
	for {
		if pp == pe {
			trap(Parse)
		}
		c := I8(pp)
		pp++
		if c == 34 {
			break
		}
		r = cat1(r, Kc(c))
	}
	if nn(r) == 1 {
		return Fst(r)
	}
	return r
}
func thex() (r K) {
	r = mk(Ct, 0)
	for pp < pe-1 {
		c := I8(pp)
		if !is(c, 128) {
			break
		}
		r = cat1(r, Kc((hx(c)<<4)+hx(I8(1+pp))))
		pp += 2
	}
	return r
}
func hx(c int32) int32 {
	if is(c, 4) {
		return c - '0'
	} else {
		return c - 'a'
	}
}
func tnms() (r K) {
	r = tnum()
	for pp < pe-1 && I8(pp) == ' ' {
		pp++
		x := tnum()
		if x == 0 {
			break
		}
		r = ncat(r, x)
	}
	return r
}
func tnum() (r K) {
	c := I8(pp)
	if pp > int32(src) {
		if c == '-' || c == '.' {
			if is(I8(pp-1), 64) {
				return 0 // e.g. x-1 is (x - 1) not (x -1)
			}
		}
	}
	if c == '-' && pp < 1+pe {
		pp++
		r = tunm()
		if r == 0 {
			pp--
			return 0
		}
		return Neg(r)
	}
	return tunm()
}
func tunm() K {
	p := pp
	r := pu()
	if r == 0 && p == pp {
		if I8(p) == '.' {
			if is(I8(1+p), 4) {
				return pflt(r)
			}
		}
		return 0
	}
	if pp < pe {
		c := I8(pp)
		if c == '.' {
			return pflt(r)
		}
		if c == 'p' {
			return ppi(float64(r))
		}
		if c == 'a' {
			return pflz(float64(r))
		}
		if r == 0 {
			if c == 'N' {
				pp++
				return missing(it)
			}
			if c == 'n' || c == 'w' {
				q := Kf(0)
				SetI64(int32(q), int64(0x7FF8000000000001)) // 0n
				if c == 'w' {
					SetF64(int32(q), inf) // 0w
				}
				pp++
				return q
			}
		}
	}
	return Ki(int32(r))
}
func pu() (r int64) {
	for pp < pe {
		c := I8(pp)
		if is(c, 4) == false {
			break
		}
		r = 10*r + int64(c-'0')
		pp++
	}
	return r
}
func pflt(i int64) K {
	f := float64(i)
	d := 1.0
	pp++ // .
	for pp < pe {
		c := I8(pp)
		if is(c, 4) == false {
			break
		}
		d /= 10.0
		f += d * float64(c-'0')
		pp++
	}
	if pp < pe {
		if I8(pp) == 'a' {
			return pflz(f)
		}
		if I8(pp) == 'p' {
			return ppi(f)
		}
	}
	return Kf(f)
}
func pflz(f float64) (r K) {
	pp++
	if pp < pe {
		r = tunm()
	}
	return Rot(Kf(f), r)
}
func ppi(f float64) (r K) {
	pp++
	return Kf(pi * f)
}

func tvrb() (r K) {
	c := I8(pp)
	if !is(c, 1) {
		p := pp
		r = tvar()
		if r != 0 { // builtins
			rp := int32(r)
			if rp > 48 && rp < 161 { // `in`find`abs`imag`conj`angle`solve`dot`sin`cos`".("`exp`log
				return K(23 + rp>>3)
			} else {
				pp = p
			}
		}
		return 0
	}
	pp++
	if c == 92 && I8(pp-2) == 32 { // \out
		return K(29)
	}
	o := int32(1)
	if pp < pe {
		if I8(pp) == 58 { // :
			pp++
			if is(c, 8) {
				o = 2 // ':
			} else {
				o = 97 // +:
			}
		}
	}
	return K(o + index(c, 228, 253))
}
func tpct() (r K) {
	c := I8(pp)
	if is(c, 48) { // ([{}]); \n
		pp++
		return K(c)
	}
	if c == 10 {
		pp++
		return K(';')
	}
	return 0
}
func tvar() (r K) {
	c := I8(pp)
	if !is(c, 2) {
		return 0
	}
	pp++
	r = Ku(uint64(c))
	for pp < pe {
		c = I8(pp)
		if !is(c, 6) {
			break
		}
		r = cat1(r, K(c)|K(ct)<<59)
		pp++
	}
	return sc(r)
}
func tsym() (r K) {
	var s K
	for I8(pp) == 96 {
		pp++
		if r == 0 {
			r = mk(St, 0)
		}
		s = 0
		if pp < pe {
			s = tchr()
			if tp(s) == ct {
				s = sc(Enl(s))
			} else if s != 0 {
				s = sc(s)
			} else {
				s = tvar()
			}
		}
		if s == 0 {
			s = K(st) << 59
		}
		r = cat1(r, s)
		if pp == pe {
			break
		}
	}
	return r
}
func ws() {
	for pp < pe {
		c := I8(pp)
		if c == 10 || c > 32 {
			break
		}
		pp++
	}
	return
}
func cmt(x K) (r K) {
	n := nn(x)
	r = mk(Ct, n)
	b := int32(1)
	h := int32(0)
	for i := int32(0); i < n; i++ {
		c := I8(int32(x) + i)
		if h != 0 {
			if c == 10 {
				h = 0
			}
		} else {
			if c == 47 && b != 0 {
				h = 1
			} else if c < 33 {
				b = 1
			} else {
				b = 0
			}
		}
		if h != 0 {
			c = 32
		}
		SetI8(int32(r)+i, c)
	}
	dx(x)
	return r
}
func is(x, m int32) bool {
	return m&I8(100+x) != 0
}
