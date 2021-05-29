package main

import (
	. "github.com/ktye/wg/module"
)

func parse(x K) (r K) {
	x = tok(x)
	pp = int32(x)
	pe = pp + 8*nn(x)
	return es()
}
func es() (r K) {
	r = mk(Lt, 0)
	for {
		n, _ := next()
		if n == 0 {
			break
		}
		if n == 59 {
			continue
		}
		pp -= 8
		x := e(t())
		if x == 0 {
			break
		}
		if nn(r) != 0 {
			r = cat1(r, 5)
		}
		r = Cat(r, x)
	}
	return r
}
func e(x K, xv int32) (r K) { // Lt
	if x == 0 {
		return 0
	}
	y, yv := t()
	if y == 0 {
		return x
	}
	if yv != 0 && xv == 0 {
		r = e(t())
		if r == 0 { // 1+ (projection)
			r = l2(K(st)<<59, 0)
			x = cat3(x, Ki(2), 27, 1)
			y = cat1(y, 20)
		}
		x, y = pasn(x, y)
		r = ucat(r, x)
		// r = ucat(r, y)
		return dyadic(r, y) // dyadic
	}
	r = e(y, yv)
	if xv == 0 {
		return cat1(cat1(ucat(r, x), 19), 2) // juxtaposition
	} else if r == y && xv+yv == 2 {
		return cat1(cat1(ucat(r, x), 27), 2) // composition
	}
	r = ucat(r, x)
	return monadic(r) // monadic
}
func t() (r K, verb int32) { // Lt
	var ln, s int32
	r, s = next()
	if r == 0 {
		return 0, 0
	}
	if r < 127 && is(int32(r), 32) {
		pp -= 8
		return 0, 0
	}
	if r == K('(') {
		r = rlist(plist(41))
	} else if r == K('{') {
		r = plam(s)
	} else if tp(r) == st {
		r = l2(r, 0)
	} else {
		rt := tp(r)
		if rt == 0 {
			r, verb = r|K(s)<<32, 1
		} else if rt == St {
			if nn(r) == 1 {
				r = Fst(r)
			}
		}
		r = l1(r)
	}
	for {
		n, _ := next()
		if n == 0 {
			break
		}
		a := int32(n)
		tn := tp(n)
		if tn == 0 && a > 20 && a < 27 { // +/
			r, verb = cat1(cat1(r, n), 1), 1
		} else if n == 91 { // [
			n, ln = plist(93)
			n, ln = pspec(r, n, ln)
			if ln < 0 {
				return n, 0
			}
			n = rlist(n, ln)
			r = cat1(cat1(Cat(n, r), K(19+ib(ln != 1))), 2)
		} else {
			pp -= 8
			break
		}
	}
	return r, verb
}
func pasn(x, y K) (K, K) {
	l := K(I64(int32(y)))
	v := int32(l)
	if nn(y) == 1 && tp(l) == 0 && v == 1 || v > 96 {
		dx(y)
		xn := nn(x)
		if xn > 2 { // indexed
			if v > 96 { // indexed-modified
				l -= 96
			}
			s := ati(rx(x), xn-4)
			x = ucat(l1(l), x) // @[`x;i;+;(rhs)]
			lp := int32(x) + 8*(nn(x)-1)
			SetI64(lp, 2+I64(lp))
			y = l2(s, l-K(int32(l)))
			//fmt.Println("=>", sK(x), sK(y))
		} else if v == 1 {
			s := Fst(x) // (`x;0) 0:Lup
			if loc != 0 {
				loc = cat1(loc, s)
			}
			x = l1(s)
			y = l1(l - 1) // l-1 is 0|srcmark
		} else { // modified
			y = cat1(cat1(l2(l-96, 2), Fst(rx(x))), l-K(int32(l)))
		}
	}
	return x, y
}
func plam(s0 int32) (r K) {
	loc = mk(St, 0)
	c := es() // todo: translate srcp
	n, s1 := next()
	if n != 125 {
		trap(Parse)
	}
	cn := nn(c)
	cp := int32(c)
	ar := int32(0)
	for i := int32(0); i < cn; i++ {
		if I64(cp) == 0 {
			if y := I32(cp-8) >> 3; y > 1 && y < 4 {
				ar = maxi(ar, y)
			}
		}
		cp += 8
	}
	i := Add(seq(1+s1-s0), Ki(s0-1))
	s := atv(rx(src), i)
	loc = Cat(ntake(ar, xyz), Unq(loc))
	cn = nn(loc)
	r = cat1(l3(c, loc, mk(Lt, cn)), s)
	loc = 0
	rp := int32(r)
	SetI32(rp-4, ar)
	return l1(K(rp) | K(lf)<<59)
}
func pspec(r, n K, ln int32) (K, int32) {
	if nn(r) == 1 && ln > 2 {
		v := K(I64(int32(r)))
		if tp(v) == 0 && int32(v) == 17 {
			dx(r)
			return cond(n, ln), -1
		}
	}
	return n, ln
}
func cond(x K, xn int32) (r K) {
	xp := int32(x) + 8*xn
	var next, sum int32
	state := int32(1)
	for xp != int32(x) {
		xp -= 8
		r = K(I64(xp))
		if sum > 0 {
			state = 1 - state
			if state != 0 {
				r = cat1(cat1(r, Ki(next)), 7) // jif
			} else {
				r = cat1(cat1(r, Ki(sum)), 6) // j
			}
			SetI64(xp, int64(r))
		}
		next = 8 * nn(r)
		sum += next
	}
	return flat(x)
}
func plist(c K) (r K, n int32) {
	r = mk(Lt, 0)
	for {
		b, _ := next()
		if b == 0 || b == c {
			break
		}
		if n == 0 {
			pp -= 8
		}
		if n != 0 && b != 59 {
			trap(Parse)
		}
		n++
		x := e(t())
		r = cat1(r, x)
		if x == 0 {
			r = cat1(r, K(st)<<59) // <null> is ` 0(lup) (Rev)
		}
	}
	return r, n
	/*
		if n == 1 {
			return Fst(r), 1
		}
		return r, n // cat3(flat(Rev(r)), Ki(n), 27, 1), 1
	*/
}
func rlist(x K, n int32) K {
	if n == 1 {
		return Fst(x)
	}
	return cat3(flat(Rev(x)), Ki(n), 27, 1)
}

func next() (r K, s int32) {
	if pp == pe {
		return 0, 0
	}
	r = K(I64(pp))
	s = 0xffffff & int32(r>>32)
	r = r &^ (K(0xffffff) << 32)
	pp += 8
	return r, s
}
func lastp(x K) K { return K(I64(int32(x) + 8*(nn(x)-1))) }
func dyadic(x, y K) K {
	l := lastp(y)
	if l < 2 {
		x = cat3(x, Ki(2), 27, 1)
		y = cat1(y, 20) // .
	}
	return cat1(ucat(x, y), 2)
}
func monadic(x K) K {
	l := lastp(x)
	if l < 2 {
		return cat1(cat1(x, 19), 2) // @
	}
	return cat1(x, 1)
}
