package i

import (
	"math"
	"math/cmplx"
	"reflect"
)

// atomic numeric monads
func rneg(a float64) float64       { return -a }
func zneg(a complex128) complex128 { return -a }
func rflr(a float64) float64       { return math.Floor(a) }
func zflr(a complex128) complex128 { return complex(math.Floor(cmplx.Abs(a)), 0) }
func rsqr(a float64) float64       { return math.Sqrt(a) }
func zsqr(a complex128) complex128 { return cmplx.Sqrt(a) }
func rnot(a float64) float64       { return rter(a == 0, 1, 0) }
func znot(a complex128) complex128 { return zter(a == 0, 1, 0) }

// atomic numeric dyads
func radd(a, b float64) float64       { return a + b }
func zadd(a, b complex128) complex128 { return a + b }
func rsub(a, b float64) float64       { return a - b }
func zsub(a, b complex128) complex128 { return a - b }
func rmul(a, b float64) float64       { return a * b }
func zmul(a, b complex128) complex128 { return a * b }
func rdiv(a, b float64) float64       { return a / b }
func zdiv(a, b complex128) complex128 { return a / b }
func rmin(a, b float64) float64       { return rter(a < b, a, b) }
func zmin(a, b complex128) complex128 { return zter(cmplx.Abs(a) < cmplx.Abs(b), a, b) } // what about equal abs? compare angle?
func rmax(a, b float64) float64       { return rter(a > b, a, b) }
func zmax(a, b complex128) complex128 { return zter(cmplx.Abs(a) > cmplx.Abs(b), a, b) }
func rles(a, b float64) float64       { return rter(a < b, 1, 0) }
func zles(a, b complex128) complex128 { return zter(cmplx.Abs(a) < cmplx.Abs(b), 1, 0) }
func rmor(a, b float64) float64       { return rter(a > b, 1, 0) }
func zmor(a, b complex128) complex128 { return zter(cmplx.Abs(a) > cmplx.Abs(b), 1, 0) }
func reql(a, b float64) float64       { return rter(a == b, 1, 0) } // tolerance?
func zeql(a, b complex128) complex128 { return zter(a == b, 1, 0) }

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

type fr1 func(float64) float64
type fr2 func(float64, float64) float64
type fz1 func(complex128) complex128
type fz2 func(complex128, complex128) complex128

func nm(x interface{}, fr fr1, fz fz1, m method) interface{} {
	if r, ok := m.call1(x); ok {
		return r
	}

	switch t := x.(type) {
	case []interface{}:
		r := make([]interface{}, len(t))
		for i := range r {
			r[i] = nm(t[i], fr, fz, m)
		}
		return r
	case map[interface{}]interface{}:
		r := make(map[interface{}]interface{})
		for k, v := range t {
			r[cpy(k)] = nm(v, fr, fz, m)
		}
		return r
	}
	// TODO: custom map and struct types.

	y, z, vec, t := nv(x)
	if y != nil {
		for i, x := range y {
			y[i] = fr(x)
		}
		return vn(y, nil, vec, t)
	}
	for i, x := range z {
		z[i] = fz(x)

	}
	return vn(nil, z, vec, t)
}

func nd(x, y interface{}, fr fr2, fz fz2, m method) interface{} {
	if r, ok := m.call2(x, y); ok {
		return r
	}
	if r, ok := ndic(x, y, fr, fz, m); ok {
		return r
	}

	xi, isxf := x.([]interface{})
	yi, isyf := y.([]interface{})
	if isxf && !isyf {
		yi = nl(y)
		isyf = true
	} else if !isxf && isyf {
		xi = nl(x)
		isxf = true
	}
	if isxf && isyf {
		if len(xi) == 0 || len(yi) == 0 {
			return nil
		}
		if len(xi) == 1 || len(yi) > 1 {
			xi = rsh(len(yi), xi).([]interface{})
		} else if len(xi) > 1 || len(yi) == 1 {
			yi = rsh(len(xi), yi).([]interface{})
		}
		r := make([]interface{}, len(xi))
		for i := range r {
			r[i] = nd(xi[i], yi[i], fr, fz, m)
		}
		return r
	}

	// TODO: any is an interface of a custom type

	xr, xz, xvec, xt := nv(x)
	yr, yz, yvec, yt := nv(y)
	if xz != nil || yz != nil {
		if xr != nil {
			xz = toZ(xr)
			xr = nil
		} else if yr != nil {
			yz = toZ(yr)
			yr = nil
		}
	}
	n1, n2 := len(xr), len(yr)
	if xr == nil {
		n1, n2 = len(xz), len(yz)
	}
	if n1 == 0 || n2 == 0 {
		return nil
	}
	if n1 == 1 && n2 > 1 {
		xr, xz = nrsh(xr, xz, n2)
		n1 = n2
	} else if n1 > 1 && n2 == 1 {
		yr, yz = nrsh(yr, yz, n1)
		n1 = n2
	}
	if n1 != n2 {
		e("size")
	}
	if xr != nil {
		for i := range xr {
			xr[i] = fr(xr[i], yr[i])
		}
	} else {
		for i := range xz {
			xz[i] = fz(xz[i], yz[i])
		}
	}
	vec := false
	if xvec || yvec {
		vec = true
	}
	if xt == yt {
		return vn(xr, xz, vec, xt)
	}
	return vn(xr, xz, vec, nil)
}
func ndic(x, y interface{}, fr fr2, fz fz2, me method) (interface{}, bool) {
	// TODO custom dicts and structs
	xd, isx := x.(map[interface{}]interface{})
	yd, isy := y.(map[interface{}]interface{})
	if isx == false && isy == false {
		return nil, false
	}
	if isx && isy {
		r := cpy(xd).(map[interface{}]interface{})
		var ok bool
		for k, v := range yd {
			if x, ok = xd[k]; ok {
				r[k] = nd(x, v, fr, fz, me)
			} else {
				r[cpy(k)] = cpy(v)
			}
		}
		return r, true
	}
	d := xd
	a := y
	flip := false
	if isx == false {
		d = yd
		a = x
		flip = true
	}
	r := cpy(d).(map[interface{}]interface{})
	n := len(r)
	m := ln(a)
	if m < 0 {
		a = []interface{}{a}
		m = 1
	}
	if m == 1 && n > 1 {
		a = rsh(n, m)
		m = n
	}
	if n != m {
		e("size")
	}
	for i, k := range keys(r) {
		x, y = r[k], at(a, i)
		if flip {
			x, y = y, x
		}
		r[k] = nd(x, y, fr, fz, me)
	}
	return r, true
}

type method string

func (m method) call1(x interface{}) (interface{}, bool) {
	v := reflect.ValueOf(x)
	t := v.Type()
	if t.NumMethod() > 0 {
		if m, ok := t.MethodByName(string(m)); ok {
			return m.Func.Call([]reflect.Value{v})[0].Interface(), true
		}
	}
	if v.Kind() == reflect.Slice {
		t := t.Elem()
		if t.NumMethod() == 0 {
			return nil, false
		}
		n := v.Len()
		if m, ok := t.MethodByName(string(m)); ok {
			r := reflect.MakeSlice(v.Type(), n, n)
			for i := 0; i < n; i++ {
				y := m.Func.Call([]reflect.Value{v.Index(i)})[0].Interface()
				r.Index(i).Set(reflect.ValueOf(y))
			}
			return r.Interface(), true
		}
		return nil, false
	}
	return nil, false
}

func (m method) call2(x, y interface{}) (interface{}, bool) {
	v := reflect.ValueOf(x)
	t := v.Type()
	if t.NumMethod() > 0 {
		if m, ok := t.MethodByName(string(m)); ok {
			return m.Func.Call([]reflect.Value{v, reflect.ValueOf(y), reflect.ValueOf(true)})[0].Interface(), true
		}
	}
	// TODO: apply to slices of values that implement the method.

	v = reflect.ValueOf(y)
	t = v.Type()
	if t.NumMethod() > 0 {
		if m, ok := t.MethodByName(string(m)); ok {
			return m.Func.Call([]reflect.Value{v, reflect.ValueOf(x), reflect.ValueOf(false)})[0].Interface(), true
		}
	}
	return nil, false
}

func nv(x interface{}) ([]float64, []complex128, bool, reflect.Type) {
	switch t := x.(type) {
	case float64:
		return []float64{t}, nil, false, floatT
	case complex128:
		return nil, []complex128{t}, false, complexT
	case []float64:
		r := make([]float64, len(t))
		copy(r, t)
		return r, nil, true, floatT
	case []complex128:
		r := make([]complex128, len(t))
		copy(r, t)
		return nil, t, true, complexT
	}
	v := reflect.ValueOf(x)
	if v.Kind() == reflect.Slice {
		n := v.Len()
		z := reflect.Zero(v.Type().Elem())
		if z.Type().ConvertibleTo(floatT) {
			r := make([]float64, n)
			for i := range r {
				r[i] = v.Index(i).Convert(floatT).Float()
			}
			return r, nil, true, z.Type()
		} else if z.Type().ConvertibleTo(complexT) {
			r := make([]complex128, n)
			for i := range r {
				r[i] = v.Index(i).Convert(complexT).Complex()
			}
			return nil, r, true, z.Type()
		} else if z.Type().ConvertibleTo(boolT) {
			r := make([]float64, n)
			for i := range r {
				b := v.Index(i).Convert(boolT).Bool()
				if b {
					r[i] = 1
				}
			}
			return r, nil, true, z.Type()
		}
		e("type")
	}

	if v.Type().ConvertibleTo(floatT) {
		return []float64{v.Convert(floatT).Float()}, nil, false, v.Type()
	}
	if v.Type().ConvertibleTo(complexT) {
		return nil, []complex128{v.Convert(complexT).Complex()}, false, v.Type()
	}
	if v.Type().ConvertibleTo(boolT) {
		b := v.Convert(boolT).Bool()
		r := 0.0
		if b {
			r = 1.0
		}
		return []float64{r}, nil, false, v.Type()
	}
	e("type")
	return nil, nil, false, floatT
}
func vn(x []float64, z []complex128, vec bool, t reflect.Type) interface{} {
	if x != nil && (t == floatT || t == nil) {
		if vec {
			return x
		}
		return x[0]
	}
	if z != nil && (t == complexT || t == nil) {
		if vec {
			return z
		}
		return z[0]
	}
	if vec == false {
		if x != nil {
			if t.ConvertibleTo(boolT) {
				b := false
				if x[0] != 0 {
					b = true
				}
				return reflect.ValueOf(b).Convert(t).Interface()
			}
			return reflect.ValueOf(x[0]).Convert(t).Interface()
		}
		return reflect.ValueOf(z[0]).Convert(t).Interface()
	}
	n := len(x)
	if x == nil {
		n = len(z)
	}
	r := reflect.MakeSlice(reflect.SliceOf(t), n, n)
	for i := 0; i < n; i++ {
		if x != nil {
			if t.ConvertibleTo(boolT) {
				b := false
				if x[0] != 0 {
					b = true
				}
				r.Index(i).Set(reflect.ValueOf(b).Convert(t))
			} else {
				r.Index(i).Set(reflect.ValueOf(x[i]).Convert(t))
			}
		} else {
			r.Index(i).Set(reflect.ValueOf(z[i]).Convert(t))
		}
	}
	return r.Interface()
}

func toZ(x []float64) []complex128 {
	z := make([]complex128, len(x))
	for i, r := range x {
		z[i] = complex(r, 0)
	}
	return z
}
func nrsh(x []float64, z []complex128, n int) ([]float64, []complex128) {
	if x == nil {
		r := make([]complex128, n)
		for i := range r {
			r[i] = z[0]
		}
		return nil, r
	}
	r := make([]float64, n)
	for i := range r {
		r[i] = x[0]
	}
	return r, nil
}
func nl(x interface{}) []interface{} {
	v := reflect.ValueOf(x)
	if v.Kind() == reflect.Slice {
		r := make([]interface{}, v.Len())
		for i := range r {
			reflect.ValueOf(x).Index(i).Set(v.Index(i))
		}
		return r
	} else {
		return []interface{}{x}
	}
}
