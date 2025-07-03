package main

import (
	. "github.com/ktye/wg/module"
)

type ftok = func() K

func tok(x K) K {
	s := cat1(src(), Kc(10))
	pp = nn(s)
	s = Cat(s, x)  // src contains all src
	pp += int32(s) // pp is the parser position within src
	pe = pp + nn(x)
	r := mk(Lt, 0)
	for {
		ws()
		if pp == pe {
			break
		}
		for i := int32(193); i < 200; i++ { // tchr, tnms, tvrb, tpct, tvar, tsym, trap
			y := Func[i].(ftok)()
			if y != 0 {
				y |= K(int64(pp-int32(s)) << 32)
				r = cat1(r, y)
				break
			}
		}
	}
	SetI32(16, int32(s)) //SetI64(512, int64(s))
	return r
}
func src() K { return ti(Ct, I32(16)) }
func tchr() K {
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
	r := mk(Ct, 0)
	q := uint32(0)
	for {
		if pp == pe {
			srcp = pe - I32(16)
			trap() //parse
		}
		c := I8(pp)
		pp++
		if c == 34 && q == 0 {
			break
		}
		if c == '\\' && q == 0 {
			q = 1
			continue
		}
		if q != 0 {
			c = cq(c)
			q = 0
		}
		r = cat1(r, Kc(c))
	}
	if nn(r) == 1 {
		return Fst(r)
	}
	return r
}
func cq(c int32) int32 { // \t \n \r \" \\   -> 9 10 13 34 92
	if c == 116 {
		return 9
	}
	if c == 110 {
		return 10
	}
	if c == 114 {
		return 13
	}
	return c
}
func thex() K {
	r := mk(Ct, 0)
	for pp < pe-1 {
		c := I8(pp)
		if is(c, 128) == 0 {
			break
		}
		r = cat1(r, Kc((hx(c)<<4)+hx(I8(1+pp))))
		pp += 2
	}
	if nn(r) == 1 {
		return Fst(r)
	}
	return r
}
func hx(c int32) int32 {
	if is(c, 4) != 0 {
		return c - '0'
	} else {
		return c - 'W'
	}
}
func tnms() K {
	r := tnum()
	for pp < pe-1 && I8(pp) == ' ' {
		pp++
		x := tnum()
		if x == 0 {
			break
		}
		t := tp(r)
		if t < 16 {
			r = Enl(r)
		}
		t = maxtype(r, x)
		r = uptype(r, t)
		r = cat1(r, uptype(x, t))
	}
	return r
}
func tnum() K {
	c := I8(pp)
	if c == '-' || c == '.' {
		if is(I8(pp-1), 64) != 0 {
			return 0 // e.g. x-1 is (x - 1) not (x -1)
		}
	}
	if c == '-' && pp < 1+pe {
		pp++
		r := tunm()
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
			if is(I8(1+p), 4) != 0 {
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
		if c == 'e' || c == 'E' {
			return Kf(pexp(float64(r)))
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
				if pp < pe && I8(pp) == 'a' {
					dx(q)
					return pflz(F64(int32(q)))
				}
				return q
			}
		}
	}
	return Ki(int32(r))
}
func pu() int64 {
	r := int64(0)
	for pp < pe {
		c := I8(pp)
		if is(c, 4) == 0 {
			break
		}
		r = 10*r + int64(c-'0')
		pp++
	}
	return r
}
func pexp(f float64) float64 {
	pp++
	e := int64(1)
	if pp < pe {
		c := I8(pp)
		if c == '-' || c == '+' {
			if c == '-' {
				e = int64(-1)
			}
			pp++
		}
	}
	e *= pu()
	return f * pow(10.0, float64(e))
}
func pflt(i int64) K {
	var c int32
	d := 1.0
	f := float64(i)
	pp++ // .
	for pp < pe {
		c = I8(pp)
		if is(c, 4) == 0 {
			break
		}
		d /= 10.0
		f += d * float64(c-'0')
		pp++
	}
	if pp < pe {
		c = I8(pp)
		if c == 'e' || c == 'E' {
			f = pexp(f)
		}
	}
	if pp < pe {
		c = I8(pp)
		if c == 'a' {
			return pflz(f)
		}
		if c == 'p' {
			return ppi(f)
		}
	}
	return Kf(f)
}
func pflz(f float64) K {
	r := K(0)
	pp++
	if pp < pe {
		r = tunm()
	}
	return Rot(Kf(f), r)
}
func ppi(f float64) K {
	pp++
	return Kf(pi * f)
}

func tvrb() K {
	c := I8(pp)
	if is(c, 1) == 0 {
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
			if is(c, 8) != 0 {
				srcp = pp - I32(16)
				trap() //parse
			}
			o = 97
		}
	}
	return K(o + idx(c, 227, 253))
}
func tpct() K {
	c := I8(pp)
	if is(c, 48) != 0 { // ([{}]); \n
		pp++
		return K(c)
	}
	if c == 10 {
		pp++
		return K(';')
	}
	return 0
}
func tvar() K {
	c := I8(pp)
	if is(c, 2) == 0 {
		return 0
	}
	pp++
	r := Ku(uint64(c))
	for pp < pe {
		c = I8(pp)
		if is(c, 6) == 0 {
			break
		}
		r = cat1(r, ti(ct, c))
		pp++
	}
	return sc(r)
}
func tsym() K {
	r := K(0)
	for I8(pp) == 96 {
		pp++
		if r == 0 {
			r = mk(St, 0)
		}
		s := K(0)
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
	var c int32
	for pp < pe {
		c = I8(pp)
		if c == 10 || c > 32 {
			break
		}
		pp++
	}
	for pp < pe {
		c = I8(pp)
		if c == 47 && I8(pp-1) < 33 {
			pp++
			for pp < pe {
				c = I8(pp)
				if c == 10 {
					break
				}
				pp++
			}
		} else {
			return
		}
	}
}
func is(x, m int32) int32 { return m & I8(100+x) }
