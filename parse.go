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
		r = ucat(r, x)
		r = ucat(r, dyadic(y))
		return r // todo dyadic
	}
	r = e(y, yv)
	r = ucat(r, monadic(x))
	return r
}
func t() (r K, verb int32) { // Lt
	r = next()
	if r == 0 {
		return 0, 0
	}
	// ...
	return l1(r), ib(tp(r) == 0)
}

func next() (r K) {
	if pp == pe {
		return 0
	}
	r = K(I64(pp))
	pp += 8
	return r
}
func monadic(x K) K {
	p := int32(x)
	SetI64(p, 40+I64(p))
	return x
}
func dyadic(x K) K {
	p := int32(x)
	SetI64(p, 80+I64(p))
	return x
}
