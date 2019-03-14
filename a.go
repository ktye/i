// i interpret
package i

import (
	"math"
	"math/cmplx"
	"reflect"
)

// Parse
func P(s string) interface{} {
	return prs(s)
}

// Eval
func E(a map[interface{}]interface{}, l interface{}) interface{} {
	return eva(a, l)
}

func cpy(v interface{}) interface{} {
	return e("TODO")
}

func e(s string) interface{} { panic(s); return nil }

func ln(v interface{}) int {
	r := reflect.ValueOf(v)
	if r.Kind() == reflect.Slice {
		return r.Len()
	}
	return -1
}

func lz(l interface{}) interface{} {
	return reflect.Zero(reflect.TypeOf(l).Elem()).Interface()
}

func mk(l interface{}, n int) interface{} {
	switch l.(type) {
	case float64:
		return make([]float64, n)
	case complex128:
		return make([]complex128, n)
	case string:
		return make([]string, n)
	}
	return reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(l)), n, n).Interface()
}

// function krange(x, f) { var r=[]; for(var z=0;z<x;z++) { r.push(f(z)); } return k(3,r); }
func krange(n int, f func(int) interface{}) []interface{} {
	l := make([]interface{}, n)
	for i := 0; i < n; i++ {
		l[i] = f(i)
	}
	return l
}

// function kmap (x, f) { return k(3, l(x).v.map(f)); }
func kmap(x interface{}, f func(interface{}, int) interface{}) interface{} {
	n := ln(x)
	if n < 0 {
		e("type")
	}
	for i := 0; i < n; i++ {
		set(x, i, f(at(x, i), i))
	}
	return x
}

// function kzip (x, y, f) { return kmap(sl(x,y), function(z, i) { return f(z, y.v[i]); }); }
func kzip(x, y interface{}, f func(interface{}, interface{}) interface{}) interface{} {
	return kmap(sl, func(v interface{}, i int) interface{} {
		return f(v, at(y, i))
	})
}

func sl(x, y interface{}) interface{} {
	if ln(x) != ln(y) {
		e("len")
	}
	return x
}

func na(v interface{}) bool {
	if ln(v) < 0 {
		return false
	}
	return math.IsNaN(re(v))
}

func impl(v interface{}, t reflect.Type) reflect.Method {
	if reflect.TypeOf(v).Implements(t) {
		return t.Elem().Method(0)
	}
	return reflect.Method{}
}

/*
func icall1(f reflect.Value, x reflect.Value) reflect.Value {
	r := f.Call([]reflect.Value{x})
	return r[0].Interface()
}
*/

func c1(v interface{}) (complex128, bool) {
	if z, ok := v.(complex128); ok {
		return z, true
	}
	if r := reflect.ValueOf(v); r.Kind() == reflect.Complex128 {
		return r.Complex(), true
	}
	if r := reflect.ValueOf(v); r.Kind() == reflect.Complex128 {
		return r.Complex(), true
	}
	return 0, false
}

var boolT = reflect.TypeOf(true)
var floatT = reflect.TypeOf(0.0)
var complexT = reflect.TypeOf(complex(0, 0))

func n1(x interface{}) (float64, complex128, bool) {
	v := reflect.ValueOf(x)
	if v.Type().ConvertibleTo(boolT) {
		b := v.Convert(boolT).Bool()
		if b {
			return 1, 0, false
		}
		return 0, 0, false
	} else if v.Type().ConvertibleTo(floatT) {
		return v.Convert(floatT).Float(), 0, false
	} else if v.Type().ConvertibleTo(complexT) {
		return 0, v.Convert(complexT).Complex(), true
	}
	e("type")
	return 0, 0, false
}

func n2(x, y interface{}) (float64, float64, complex128, complex128, bool) {
	a, aok := c1(x)
	b, bok := c1(y)
	if aok && bok {
		return 0, 0, a, b, true
	} else if aok {
		return 0, 0, a, complex(re(y), 0), true
	} else if bok {
		return 0, 0, complex(re(x), 0), b, true
	}
	return re(x), re(y), 0, 0, false
}

func idx(v interface{}) int { return int(re(v)) }

func re(v interface{}) float64 {
	switch w := v.(type) {
	case float64:
		return w
	case bool:
		if w {
			return 1
		}
		return 0
	case int:
		return float64(w)
	case complex128:
		if cmplx.IsNaN(w) {
			return math.NaN()
		}
		return real(w)
	}
	r := reflect.ValueOf(v)
	if k := r.Kind(); k == reflect.Bool {
		if r.Bool() {
			return 1
		}
		return 0
	} else if k < reflect.Uint {
		return float64(r.Uint())
	}
	return float64(r.Int()) // panics
}

func at(l interface{}, i int) interface{} {
	switch v := l.(type) {
	case []interface{}:
		return v[i]
	case []float64:
		return v[i]
	case []complex128:
		return v[i]
	}
	v := reflect.ValueOf(l)
	return v.Index(i).Interface()
}

func set(l interface{}, i int, v interface{}) {
	switch t := l.(type) {
	case []interface{}:
		t[i] = v
		return
	case []float64:
		t[i] = v.(float64)
		return
	case []complex128:
		t[i] = v.(complex128)
		return
	}
	reflect.ValueOf(l).Index(i).Set(reflect.ValueOf(v))
}
