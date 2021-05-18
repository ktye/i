package k

import (
	"math"
	"math/cmplx"
)

type numbase byte

const (
	booltype numbase = iota + 1
	bytetype
	inttype
	floattype
	complextype
	stringtype
)

func numtype(x T) numbase {
	switch x.(type) {
	case bool, B:
		return booltype
	case byte, C:
		return bytetype
	case int, I:
		return inttype
	case float64, F:
		return floattype
	case complex128, Z:
		return complextype
	case string, S:
		return stringtype
	default:
		return 0
	}
}
func uptype(x T, dst numbase) T {
	dx(x)
	switch v := x.(type) {
	case bool:
		i := ib(v)
		if dst == inttype {
			return i
		} else if dst == floattype {
			return float64(i)
		} else if dst == complextype {
			return complex(float64(i), 0)
		}
	case byte:
		if dst == inttype {
			return int(v)
		}
		panic("type")
	case int:
		if dst == floattype {
			return float64(v)
		} else if dst == complextype {
			return complex(float64(v), 0)
		}
	case float64:
		if dst == complextype {
			return complex(v, 0)
		}
	case B:
		if dst == inttype {
			r := mk(0, len(v.v)).(I)
			for i, b := range v.v {
				r.v[i] = ib(b)
			}
			return r
		} else if dst == floattype {
			r := mk(0.0, len(v.v)).(F)
			for i, b := range v.v {
				r.v[i] = float64(ib(b))
			}
			return r
		} else if dst == complextype {
			r := mk(complex(0, 0), len(v.v)).(Z)
			for i, b := range v.v {
				r.v[i] = complex(float64(ib(b)), 0)
			}
			return r
		}
	case C:
		if dst == inttype {
			r := mk(0, len(v.v)).(I)
			for i, b := range v.v {
				r.v[i] = int(b)
			}
			return r
		}
		panic("type")
	case I:
		if dst == floattype {
			r := mk(0.0, len(v.v)).(F)
			for i, b := range v.v {
				r.v[i] = float64(b)
			}
			return r
		} else if dst == complextype {
			r := mk(complex(0, 0), len(v.v)).(Z)
			for i, b := range v.v {
				r.v[i] = complex(float64(b), 0)
			}
			return r
		}
	case F:
		if dst == complextype {
			r := mk(complex(0, 0), len(v.v)).(Z)
			for i, b := range v.v {
				r.v[i] = complex(b, 0)
			}
			return r
		}
	}
	panic("uptype")
}
func ib(v bool) int {
	if v {
		return 1
	}
	return 0
}

func nm(x T, fi func(int) int, ff func(float64) float64, fz func(complex128) complex128) T {
	xt := numtype(x)
	if xt < inttype {
		x, xt = uptype(x, inttype), inttype
	}
	if xt == inttype && fi == nil {
		x, xt = uptype(x, floattype), floattype
	}
	xv, o := x.(vector)
	if o == false {
		switch v := x.(type) {
		case int:
			return fi(v)
		case float64:
			return ff(v)
		case complex128:
			return fz(v)
		}
	}
	xv = use(xv)
	switch v := x.(type) {
	case I:
		for i := range v.v {
			v.v[i] = fi(v.v[i])
		}
	case F:
		for i := range v.v {
			v.v[i] = ff(v.v[i])
		}
	case Z:
		for i := range v.v {
			v.v[i] = fz(v.v[i])
		}
	default:
		panic("type(nm)")
	}
	return xv
}
func nd(x, y T, f2 nf2, fi func(int, int) int, ff func(float64, float64) float64, fz func(complex128, complex128) complex128) T {
	if r := ndlist(x, y, f2); r != nil {
		return r
	}
	xt, yt := numtype(x), numtype(y)
	if xt < inttype {
		x, xt = uptype(x, inttype), inttype
	}
	if yt < inttype {
		y, yt = uptype(y, inttype), inttype
	}
	if xt < yt {
		x, xt = uptype(x, yt), yt
	}
	if yt < xt {
		y, yt = uptype(y, xt), xt
	}
	xv, isxv := x.(vector)
	yv, isyv := y.(vector)
	if isxv == false && isyv == false {
		switch xt {
		case inttype:
			return fi(x.(int), y.(int))
		case floattype:
			return ff(x.(float64), y.(float64))
		case complextype:
			return fz(x.(complex128), y.(complex128))
		}
	} else if isxv == false && isyv {
		yv = use(yv)
		switch xt {
		case inttype:
			ii, v := x.(int), yv.(I).v
			for i := range v {
				v[i] = fi(ii, v[i])
			}
			return yv
		case floattype:
			ii, v := x.(float64), yv.(F).v
			for i := range v {
				v[i] = ff(ii, v[i])
			}
			return yv
		case complextype:
			ii, v := x.(complex128), yv.(Z).v
			for i := range v {
				v[i] = fz(ii, v[i])
			}
			return yv
		}
	} else if isxv == true && isyv == false {
		xv = use(xv)
		switch xt {
		case inttype:
			ii, v := y.(int), xv.(I).v
			for i := range v {
				v[i] = fi(v[i], ii)
			}
			return xv
		case floattype:
			ii, v := y.(float64), xv.(F).v
			for i := range v {
				v[i] = ff(v[i], ii)
			}
			return xv
		case complextype:
			ii, v := y.(complex128), xv.(Z).v
			for i := range v {
				v[i] = fz(v[i], ii)
			}
			return xv
		}
	} else {
		if yv.ln() != xv.ln() {
			panic("length")
		}
		rv := use2(yv, xv)
		switch xt {
		case inttype:
			ri, xi, yi := rv.(I), x.(I), y.(I)
			for i := range ri.v {
				ri.v[i] = fi(xi.v[i], yi.v[i])
			}
			return rv
		case floattype:
			ri, xi, yi := rv.(F), x.(F), y.(F)
			for i := range ri.v {
				ri.v[i] = ff(xi.v[i], yi.v[i])
			}
			return rv
		case complextype:
			ri, xi, yi := rv.(Z), x.(Z), y.(Z)
			for i := range ri.v {
				ri.v[i] = fz(xi.v[i], yi.v[i])
			}
			return rv
		}
	}
	panic("type")
}
func ndlist(x, y T, f nf2) T {
	_, xl := x.(L)
	_, yl := y.(L)
	if xl == false && yl == false {
		return nil
	} else if xl == true && yl == false {
		return eachleft(f2(f), x, y)
	} else if xl == false && yl == true {
		return eachright(f2(f), x, y)
	} else {
		return each2(f2(f), x, y)
	}
}

type nf1 func(T) T
type nf2 func(T, T) T

//+ (1;2) /1 2
//+ (1;2 3) /(1;2 3)
//+ (1)+2 /3
//+ (1+2;4)+3 /6 7
//+ 1+1 /2
//+ 2 3+6 /8 9
//+ (1 2;3 4)+1 /(2 3;4 5)
//+ 1+(2 3;4) /(3 4;5)
//+ (1;2 3)+(2.;4 5) /(3.;6 8)
func add(x, y T) T                    { return nd(x, y, add, addi, addf, addz) }
func addi(x, y int) int               { return x + y }
func addf(x, y float64) float64       { return x + y }
func addz(x, y complex128) complex128 { return x + y }

//- 2 3-6 /-4 -3
//- (1-2)-3 /-4
//- 1-2-3 /2
//- 1-(2-3) /2
func neg(x T) T {
	return nm(x, func(x int) int { return -x }, func(x float64) float64 { return -x }, func(x complex128) complex128 { return -x })
}

//- 7 8-6 5 /1 3
func sub(x, y T) T                    { return nd(x, y, sub, subi, subf, subz) }
func subi(x, y int) int               { return x - y }
func subf(x, y float64) float64       { return x - y }
func subz(x, y complex128) complex128 { return x - y }

//* 7*8 /56
func mul(x, y T) T                    { return nd(x, y, mul, muli, mulf, mulz) }
func muli(x, y int) int               { return x * y }
func mulf(x, y float64) float64       { return x * y }
func mulz(x, y complex128) complex128 { return x * y }

//% %9 /3.
func sqrt(x T) T {
	return nm(x, nil, func(x float64) float64 { return math.Sqrt(x) }, nil)
}

//% 9%2 /4
//% 9%2. /4.5
func div(x, y T) T                    { return nd(x, y, div, divi, divf, divz) }
func divi(x, y int) int               { return x / y }
func divf(x, y float64) float64       { return x / y }
func divz(x, y complex128) complex128 { return x / y }

//& 3&4 /3
//& 3 4.&3.5 /3 3.5
func min(x, y T) T { return nd(x, y, min, mini, math.Min, minz) }
func mini(x, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}

func minz(x, y complex128) complex128 {
	if real(x) < real(y) {
		return x
	} else if real(x) == real(y) && imag(x) < imag(y) {
		return x
	}
	return y
}

//| 3|4 /4
//| 3 4.|3.5 /3.5 4
func max(x, y T) T { return nd(x, y, max, maxi, math.Max, maxz) }
func maxi(x, y int) int {
	if x > y {
		return x
	} else {
		return y
	}
}
func maxz(x, y complex128) complex128 {
	if real(x) > real(y) {
		return x
	} else if real(x) == real(y) && imag(x) > imag(y) {
		return x
	}
	return y
}

//abs abs -1 2 /1 2
//abs abs 2a30 /2.
func abs(x T) T {
	switch v := x.(type) {
	case int:
		if v < 0 {
			return -v
		}
		return v
	case I:
		x = use(v).(I)
		for i, u := range v.v {
			if u < 0 {
				v.v[i] = -u
			}
		}
		return x
	case float64:
		return math.Abs(v)
	case F:
		x = use(v).(F)
		for i, u := range v.v {
			v.v[i] = math.Abs(u)
		}
		return x
	case complex128:
		return cmplx.Abs(v)
	case Z:
		r := mk(0.0, len(v.v)).(F)
		for i, u := range v.v {
			r.v[i] = cmplx.Abs(u)
		}
		dx(x)
		return r
	default:
		panic("type")
	}
}

//angle angle 1a20 2a45 /20 45.
//angle angle 1.2 /0.
func angle(x T) T {
	switch v := x.(type) {
	case int, float64:
		return 0.0
	case F, I:
		return mk(0.0, count(x).(int))
	case complex128:
		return zang(v)
	case Z:
		r := mk(0.0, len(v.v)).(F)
		for i, u := range v.v {
			r.v[i] = zang(u)
		}
		dx(x)
		return r
	default:
		panic("type")
	}
}
func zang(z complex128) (a float64) {
	a = 180.0 / math.Pi * cmplx.Phase(z)
	if a < 0 {
		a += 360
	}
	return a
}

//angle 1a20 angle 25 /1a45
func rotate(x, y T) T {
	var z T
	switch v := y.(type) {
	case int, I:
		return rotate(x, uptype(y, floattype))
	case float64:
		z = expi(v)
	case F:
		r := mk(complex(0, 0), len(v.v)).(Z)
		for i := range v.v {
			r.v[i] = expi(v.v[i])
		}
		dx(y)
		z = r
	default:
		panic("type")
	}
	return mul(x, z)
}
func expi(a float64) complex128 {
	s, c := math.Sincos(a * math.Pi / 180.0)
	return complex(c, s)
}
func rrot(r, a float64) complex128 {
	z := expi(a)
	return complex(r*real(z), r*imag(z))
}
func zrot(x complex128, a float64) complex128 {
	z := expi(a)
	return x * z
}

//real real 1a300 /0.5
func zreal(x T) T {
	switch v := x.(type) {
	case complex128:
		return real(v)
	case Z:
		r := mk(0.0, len(v.v)).(F)
		for i, u := range v.v {
			r.v[i] = real(u)
		}
		return r
	default:
		panic("type")
	}
}

//imag imag 1a60 /0.8660254037844386
func zimag(x T) T {
	switch v := x.(type) {
	case complex128:
		return imag(v)
	case Z:
		r := mk(0.0, len(v.v)).(F)
		for i, u := range v.v {
			r.v[i] = imag(u)
		}
		return r
	default:
		panic("type")
	}
}

//imag 1 imag 1 /1.4142135623730951a45
func complx(x, y T) T { return add(x, mul(complex(0, 1), y)) }

//conj conj 1a60 /1a300
func conj(x T) T {
	switch v := x.(type) {
	case complex128:
		return cmplx.Conj(v)
	case Z:
		v = use(v).(Z)
		for i, u := range v.v {
			v.v[i] = cmplx.Conj(u)
		}
		return v
	default:
		panic("type")
	}
}

//_ _-0.1 3.3 /-1 3
func (k *K) floor(x T) T {
	switch v := x.(type) {
	case float64:
		return ffloor(v)
	case F:
		r := mk(0, len(v.v)).(I)
		for i := range v.v {
			r.v[i] = ffloor(v.v[i])
		}
		r.init()
		dx(x)
		return r
	case int, I:
		return x
	case L:
		return each(f1(k.floor), x)
	default:
		panic("type")
	}
}
func ffloor(x float64) int {
	if math.IsNaN(x) {
		return 0
	} else if math.IsInf(x, 1) {
		return maxInt
	} else if math.IsInf(x, -1) {
		return minInt
	} else {
		return int(math.Floor(x))
	}
}

const maxInt = int(^uint(0) >> 1)
const minInt = -maxInt - 1
