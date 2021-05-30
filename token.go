package main

import (
	. "github.com/ktye/wg/module"
)

type ftok = func() K

func tok(x K) (r K) {
	var y K
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
		for i := int32(192); i < 199; i++ { // tbln, tnum, tvrb, tpct, tvar, tsym, tchr
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
func tnum() K {
	p := pp
	c := I8(p)
	if p > int32(src) {
		if c == '-' || c == '.' {
			if is(I8(p-1), 64) {
				return 0 // e.g. x-1 is (x - 1) not (x -1)
			}
		}
	}
	r := pi()
	if r == 0 && p == pp {
		if c == '.' && is(I8(1+pp), 4) {
			return pflt(r)
		}
		return 0
	}
	if pp < pe && I8(pp) == '.' {
		return pflt(r)
	}
	return Ki(int32(r))
}
func pi() (r int64) {
	if I8(pp) == '-' {
		pp++
		p := pp
		r = -pu()
		if r == 0 {
			pp = p - 1
		}
	} else {
		r = pu()
	}
	return r
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
	return Kf(f)
}

func tvrb() (r K) {
	c := I8(pp)
	if !is(c, 1) {
		return 0
	}
	pp++
	if c == 92 && I8(pp-2) == 32 { // \out
		return K(29)
	}
	o := int32(1)
	if pp < pe && I8(pp) == 58 { // :
		pp++
		if is(c, 8) {
			o = 2 // ':
		} else {
			o = 97 // +:
		}
	}
	return K(o + index(c, 228, 253))
}
func tpct() (r K) {
	c := I8(pp)
	if is(c, 48) { // ([{}]);
		pp++
		return K(c)
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
			s = tvar()
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
func tchr() (r K) {
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
func is(x, m int32) bool { return m&I8(100+x) != 0 }
