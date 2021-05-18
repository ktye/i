package k

import (
	"bytes"
	"math"
	"math/cmplx"
)

type fcmp func(T, T) T

//= 1 2=2 /01b
func equal(x, y T) T {
	return nc(x, y, equal, func(x, y int) bool { return x == y }, feq, zeq, func(x, y string) bool { return x == y })
}

//> 2<!4 /0001b
func more(x, y T) T {
	return nc(x, y, more, func(x, y int) bool { return x > y }, func(x, y float64) bool { return x > y }, zgt, func(x, y string) bool { return x > y })
}

//< 2>!4 /1100b
//< `a`b`c<`alpha /100b
func less(x, y T) T {
	return nc(x, y, less, func(x, y int) bool { return x < y }, func(x, y float64) bool { return x < y }, zlt, func(x, y string) bool { return x < y })
}

func nc(x, y T, f fcmp, fi func(int, int) bool, ff func(float64, float64) bool, fz func(complex128, complex128) bool, fs func(string, string) bool) T {
	if r := nclist(x, y, f); r != nil {
		return r
	}
	xt, yt := numtype(x), numtype(y)
	if xt < yt {
		x, xt = uptype(x, yt), yt
	}
	if yt < xt {
		y, yt = uptype(y, xt), xt
	}
	xv, ix := x.(vector)
	yv, iy := y.(vector)
	if ix == false && iy == false {
		switch xt {
		case booltype:
			return fi(ib(x.(bool)), ib(y.(bool)))
		case bytetype:
			return fi(int(x.(byte)), int(y.(byte)))
		case inttype:
			return fi(x.(int), y.(int))
		case floattype:
			return ff(x.(float64), y.(float64))
		case complextype:
			return fz(x.(complex128), y.(complex128))
		case stringtype:
			return fs(x.(string), y.(string))
			panic("type")
		}
	} else if ix == false && iy {
		defer dx(y)
		r := make([]bool, yv.ln())
		switch xt {
		case booltype:
			ii, v := ib(x.(bool)), yv.(B).v
			for i := range r {
				r[i] = fi(ii, ib(v[i]))
			}
		case bytetype:
			ii, v := int(x.(int)), yv.(C).v
			for i := range r {
				r[i] = fi(ii, int(v[i]))
			}
		case inttype:
			ii, v := x.(int), yv.(I).v
			for i := range r {
				r[i] = fi(ii, v[i])
			}
		case floattype:
			ii, v := x.(float64), yv.(F).v
			for i := range v {
				r[i] = ff(ii, v[i])
			}
		case complextype:
			ii, v := x.(complex128), yv.(Z).v
			for i := range v {
				r[i] = fz(ii, v[i])
			}
		case stringtype:
			ii, v := x.(string), yv.(S).v
			for i := range v {
				r[i] = fs(ii, v[i])
			}
		default:
			panic("type")
		}
		return KB(r)
	} else if ix && iy == false {
		defer dx(x)
		r := make([]bool, xv.ln())
		switch xt {
		case booltype:
			ii, v := ib(y.(bool)), xv.(B).v
			for i := range r {
				r[i] = fi(ib(v[i]), ii)
			}
		case bytetype:
			ii, v := int(y.(int)), xv.(C).v
			for i := range r {
				r[i] = fi(int(v[i]), ii)
			}
		case inttype:
			ii, v := y.(int), xv.(I).v
			for i := range r {
				r[i] = fi(v[i], ii)
			}
		case floattype:
			ii, v := y.(float64), xv.(F).v
			for i := range v {
				r[i] = ff(v[i], ii)
			}
		case complextype:
			ii, v := y.(complex128), xv.(Z).v
			for i := range v {
				r[i] = fz(v[i], ii)
			}
		case stringtype:
			ii, v := y.(string), xv.(S).v
			for i := range v {
				r[i] = fs(v[i], ii)
			}
		default:
			panic("type")
		}
		return KB(r)
	} else {
		if yv.ln() != xv.ln() {
			panic("length")
		}
		defer dx(x)
		defer dx(y)
		r := make([]bool, xv.ln())
		switch xt {
		case booltype:
			xx, yy := xv.(B), yv.(B)
			for i, ii := range xx.v {
				r[i] = fi(ib(ii), ib(yy.v[i]))
			}
		case bytetype:
			xx, yy := xv.(C), yv.(C)
			for i, ii := range xx.v {
				r[i] = fi(int(ii), int(yy.v[i]))
			}
		case inttype:
			xx, yy := xv.(I), yv.(I)
			for i, ii := range xx.v {
				r[i] = fi(ii, yy.v[i])
			}
		case floattype:
			xx, yy := xv.(F), yv.(F)
			for i, ii := range xx.v {
				r[i] = ff(ii, yy.v[i])
			}
		case complextype:
			xx, yy := xv.(Z), yv.(Z)
			for i, ii := range xx.v {
				r[i] = fz(ii, yy.v[i])
			}
		case stringtype:
			xx, yy := xv.(S), yv.(S)
			for i, ii := range xx.v {
				r[i] = fs(ii, yy.v[i])
			}
		default:
			panic("type")
		}
		return KB(r)
	}
	panic("type")
}
func nclist(x, y T, f fcmp) T {
	_, xl := x.(L)
	_, yl := y.(L)
	if xl == false && yl == false {
		return nil
	} else if xl == true && yl == false {
		return eachleft(f, x, y)
	} else if xl == false && yl == true {
		return eachright(f, x, y)
	} else {
		return each2(f, x, y)
	}
}
func feq(x, y float64) bool    { return (x == y || math.IsNaN(x) && math.IsNaN(y)) }
func zeq(x, y complex128) bool { return (x == y || cmplx.IsNaN(x) && cmplx.IsNaN(y)) }
func blt(a, b bool) bool       { return a == false && b }
func zlt(a, b complex128) bool {
	if real(a) == real(b) {
		return imag(a) < imag(b)
	}
	return real(a) < real(b)
}
func zgt(a, b complex128) bool {
	if real(a) == real(b) {
		return imag(a) > imag(b)
	}
	return real(a) > real(b)
}
func llt(x, y T) bool {
	switch v := x.(type) {
	case C:
		if w, o := y.(C); o {
			return bytes.Compare(v.v, w.v) < 0
		}
	}
	panic("type")
}
