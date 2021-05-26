package k

import (
	. "github.com/ktye/wg/module"
)

func parse(x K) (r K) {
	x = tok(x)
	pp = int32(x)
	pe = pp + 8*nn(x)

	r = e(t())
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
		r = ucat(r, y)
		return dyadic(r) // dyadic
	}
	r = e(y, yv)
	if xv == 0 || (r == y && xv+yv == 2) {
		return cat1(cat1(ucat(r, x), 19), 2) // juxtaposition or train
	}
	r = ucat(r, x)
	return monadic(r) // monadic
}
func t() (r K, verb int32) { // Lt
	var ln int32
	r = next()
	if r == 0 {
		return 0, 0
	}
	if r < 127 && is(int32(r), 32) {
		pp -= 8
		return 0, 0
	}
	if r == K('(') {
		r, ln = plist(41)
		if ln == 1 {
			r = Fst(r)
		} else {
			r = cat3(flat(Rev(r)), Ki(ln), 27, 1)
		}
	} else if r == K('{') {
		r = plam()
	} else if tp(r) == st {
		r = l2(r, 0)
	} else {
		r, verb = l1(r), ib(tp(r) == 0)
	}
	for {
		n := next()
		if n == 0 {
			break
		}
		a := int32(n)
		tn := tp(n)
		if tn == 0 && a > 20 && a < 27 { // +/
			r, verb = cat1(cat1(r, n), 1), 1
		} else if n == 91 { // [
			n, ln = plist(93)
			verb = 0
			if ln == 1 {
				r = cat1(cat1(Cat(Fst(n), r), 19), 2)
			} else {
				n = cat3(flat(Rev(n)), Ki(ln), 27, 1)
				r = cat1(cat1(Cat(n, r), 20), 2)
			}
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
			x = ntake(nn(x)-1, x) // (`x;0) 0:Lup
			y = l1(l - 1)         // l-1 is 0|srcmark
		} else { // modified
			y = cat1(cat1(l2(l-96, 2), Fst(rx(x))), l-K(int32(l)))
		}
	}
	return x, y
}
func plam() K { trap(Nyi); return 0 }
func plist(c K) (r K, n int32) {
	r = mk(Lt, 0)
	b := next()
	if b == 0 || b == c {
		return r, 0
	}
	pp -= 8
	for {
		n++
		x := e(t())
		r = cat1(r, x)
		if x == 0 {
			r = cat1(r, K(st)<<59) // <null> is ` 0(lup) (Rev)
		}
		b = next()
		if b == c {
			break
		}
		if b != 59 { // ;
			trap(Parse)
		}
	}
	return r, n
}

func next() (r K) {
	if pp == pe {
		return 0
	}
	r = K(I64(pp))
	pp += 8
	return r
}
func lastp(x K) K { return K(I64(int32(x) + 8*(nn(x)-1))) }
func dyadic(x K) K {
	l := lastp(x)
	if l < 2 {
		x = cat1(x, 20) // .
	}
	return cat1(x, 2)
}
func monadic(x K) K {
	l := lastp(x)
	if l < 2 {
		return cat1(cat1(x, 19), 2) // @
	}
	return cat1(x, 1)
}
