package i

import "reflect"

func e(s string) { panic(s) }

func ln(v interface{}) int {
	r := reflect.ValueOf(v)
	if r.Kind() == reflect.Slice {
		return r.Len()
	}
	return -1
}

func lz(l interface{}) interface{} {
	reflect.Zero(reflect.TypeOf(l).Elem()).Interface()
}

func impl(v interface{}, t reflect.Type) reflect.Value {
	if reflect.TypeOf(v).Implements(t) {
		return t.Elem().Methods(0)
	}
	return nil
}

func call(f reflect.Value, v ...interface{}) interface{} {
	w := make([]reflect.Value, len(v))
	for i := range v {
		w[i] = reflect.ValueOf(v[i])
	}
	r := f.Call(w)
	return z.Index(0).Interface()
}

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
