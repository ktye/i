package i

import (
	"math"
	"math/cmplx"
	"reflect"
)

// Monads
func rneg(a float64) float64       { return -a }
func zneg(a complex128) complex128 { return -a }
func rflr(a float64) float64       { return math.Floor(a) }
func zflr(a complex128) complex128 { return complex(math.Floor(cmplx.Abs(a)), 0) }
func rsqr(a float64) float64       { return math.Sqrt(a) }
func zsqr(a complex128) complex128 { return cmplx.Sqrt(a) }
func rnot(a float64) float64 {
	if a == 0 {
		return 1
	}
	return 0
}
func znot(a complex128) complex128 { // or should that not exist?
	if a == 0 {
		return complex(1, 0)
	}
	return complex(1, 0)
}

// Dyads
func radd(a, b float64) float64       { return a + b }
func zadd(a, b complex128) complex128 { return a + b }
func rsub(a, b float64) float64       { return a - b }
func zsub(a, b complex128) complex128 { return a - b }
func rmul(a, b float64) float64       { return a * b }
func zmul(a, b complex128) complex128 { return a * b }
func rdiv(a, b float64) float64       { return a / b }
func zdiv(a, b complex128) complex128 { return a / b }
func rmin(a, b float64) float64       { return rter(a < b, a, b) }
func zmin(a, b complex128) complex128 { return zter(cmplx.Abs(a) < cmplx.Abs(b), a, b) }
func rmax(a, b float64) float64       { return rter(a > b, a, b) }
func zmax(a, b complex128) complex128 { return zter(cmplx.Abs(a) > cmplx.Abs(b), a, b) }
func rles(a, b float64) float64       { return rter(a < b, 1, 0) }
func zles(a, b complex128) complex128 { return zter(cmplx.Abs(a) < cmplx.Abs(b), 1, 0) }
func rmor(a, b float64) float64       { return rter(a > b, 1, 0) }
func zmor(a, b complex128) complex128 { return zter(cmplx.Abs(a) > cmplx.Abs(b), 1, 0) }
func reql(a, b float64) float64       { return rter(a == b, 1, 0) }
func zeql(a, b complex128) complex128 { return zter(cmplx.Abs(a) == cmplx.Abs(b), 1, 0) }

func rter(c bool, a, b float64) float64 {
	if c {
		return a
	}
	return b
}
func zter(c bool, a, b complex128) complex128 {
	if c {
		return a
	}
	return b
}

type f1 func(interface{}) interface{}
type f2 func(interface{}, interface{}) interface{}
type fr1 func(float64) float64
type fr2 func(float64, float64) float64
type fz1 func(complex128) complex128
type fz2 func(complex128, complex128) complex128

func nm(x interface{}, fr fr1, fz fz1, m method) interface{} {
	switch t := x.(type) {
	case float64:
		return fr(t)
	case complex128:
		return fz(t)
	case []float64:
		r := make([]float64, len(t))
		for i := range r {
			r[i] = fr(t[i])
		}
		return r
	case []complex128:
		r := make([]complex128, len(t))
		for i := range r {
			r[i] = fz(t[i])
		}
		return r
	case []interface{}:
		r := make([]interface{}, len(t))
		for i := range r {
			r[i] = nm(t[i], fr, fz, m)
		}
		return r
	}

	if r, ok := m.call1(x); ok {
		return r
	}

	v := reflect.ValueOf(x)
	if v.Kind() == reflect.Slice {
		n := v.Len()
		r := make([]interface{}, n)
		for i := range r {
			r[i] = nm(v.Index(i).Interface(), fr, fz, m)
		}
		return r
	}
	r, z, isz := n1(x) // panics
	if isz {
		return fz(z)
	}
	return fr(r)
}

func nd(x, y interface{}, fr fr2, fz fz2, mt string) f2 {
	e("TODO")
	return nil
}

type method string

func (m method) call1(x interface{}) (interface{}, bool) {
	v := reflect.ValueOf(x)
	if v.Kind() == reflect.Slice {
		n := v.Len()
		if m, ok := v.Type().Elem().MethodByName(string(m)); ok {
			r := make([]interface{}, n)
			for i := range r {
				r[i] = m.Func.Call([]reflect.Value{v.Index(i)})[0].Interface()
			}
			return r, true
		}
		return nil, false
	} else if m, ok := v.Type().MethodByName(string(m)); ok {
		return m.Func.Call([]reflect.Value{v})[0].Interface(), true
	}
	return nil, false
}

/*
func nd(fr func(float64, float64) float64, fz func(complex128, complex128) complex128, it reflect.Type) func(v ...interface{}) interface{} {
	return func(v ...interface{}) interface{} {
		x, y := v[0], v[1]
		xn, yn := ln(x, ln(y))
		if xn < 0 && yn < 0 {
			if f := impl(x, it); f != nil {
				return ical(f, x, y)
			}
			a, b, c, d, z := n2(x, y)
			if z {
				return fz(c, d)
			}
			return fr(a, b)
		} else if xn == 0 && yn == 0 {
			return nil
		} else if xn != yn {
			if xn < 0 || xn == 1 {
				x = rsh([]int{yn}, x)
				xn = yn
			} else if yn < 0 || yn == 1 {
				y = rsh([]int{xn}, y)
				yn = xn
			}
		}
		x0, y0 := lz(x), lz(y)

		TODO
		if _, _, _, _, z := n2(x0, y0); z {

		}

		} else if (an < 0 || an == 1) {
			// reshape a to bn
		} else if (bn < 0 || bn < 0)

		b := v[1]

	}
}
*/
