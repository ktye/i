package i

import "reflect"

// monadic verbs
func flp(x v) v { return e("TODO") }
func neg(x v) v { return nm(x, rneg, zneg, "Neg") }
func fst(v v) v {
	//function first (x) { return (x.t == 4) ? first(x.v) : (x.t != 3) ? x : len(x) ? x.v[0]:k(3,[]); }
	// TODO dict
	if n := ln(v); n < 0 {
		return v
	} else if n == 0 {
		return nil
	}
	return at(v, 0)
}
func sqr(x v) v { return nm(x, rsqr, zsqr, "Sqr") }
func til(x v) v {
	if d, ok := md(x); ok {
		return d.k
	}
	f, z, vec, t := nv(x)
	if z != nil || vec {
		e("type") // !z
	}
	if f[0] < 0 {
		return e("domain") // !-n
	}
	r := make(fv, int(f[0]))
	for i := range r {
		r[i] = float64(i)
	}
	return vn(r, nil, true, t)
}
func odo(x v) v { return e("TODO") } // →impl
func wer(x v) v { return e("TODO") }
func rev(x v) v { return e("TODO") }
func asc(x v) v { return e("TODO") }
func dsc(x v) v { return e("TODO") }
func eye(x v) v { return e("TODO") }
func grp(x v) v { return e("TODO") }
func not(x v) v { return nm(x, rnot, znot, "Not") }
func enl(x v) v {
	v := rval(x)
	switch v.Kind() {
	case reflect.Func, reflect.Slice, reflect.Map:
		return l{x}
	}
	l := reflect.MakeSlice(reflect.SliceOf(v.Type()), 1, 1)
	l.Index(0).Set(v)
	return l.Interface()
}
func is0(x v) v { return e("TODO") }
func cnt(x v) v { return e("TODO") }
func flr(x v) v { return nm(x, rflr, zflr, "Flr") }
func fmt(x v) v { return e("TODO") }
func fgn(x v) v { return e("TODO") }
func unq(x v) v { return e("TODO") }
func evl(x v) v { return e("TODO") }

// dyadic verbs
func add(x, y v) v { return nd(x, y, radd, zadd, "Add") }
func sub(x, y v) v { return nd(x, y, rsub, zsub, "Sub") }
func mul(x, y v) v { return nd(x, y, rmul, zmul, "Mul") }
func div(x, y v) v { return nd(x, y, rdiv, zdiv, "Div") }
func mod(x, y v) v { return e("TODO") }
func mkd(x, y v) v { return e("TODO") }
func min(x, y v) v { return nd(x, y, rmin, zmin, "Min") }
func max(x, y v) v { return nd(x, y, rmax, zmax, "Max") } // cast to bool?
func les(x, y v) v { return nd(x, y, rles, zles, "Les") } // ?
func mor(x, y v) v { return nd(x, y, rmor, zmor, "Mor") } // ?
func eql(x, y v) v { return nd(x, y, reql, zeql, "Eql") } // ?
func mch(x, y v) v { return e("TODO") }
func cat(x, y v) v {
	// TODO dict
	nx := ln(x)
	if nx < 0 {
		x = enl(x)
		nx = 1
	}
	ny := ln(y)
	if ny < 0 {
		y = enl(y)
		ny = 1
	}
	if t := rtyp(x); t == rtyp(y) {
		var l reflect.Value
		l = reflect.MakeSlice(t, nx+ny, nx+ny)
		for i := 0; i < nx; i++ {
			l.Index(i).Set(rval(at(x, i)))
		}
		for i := 0; i < ny; i++ {
			l.Index(nx + i).Set(rval(at(y, i)))
		}
		return l.Interface()
	}
	l := make(l, nx+ny)
	for i := 0; i < nx; i++ {
		l[i] = at(x, i)
	}
	for i := 0; i < ny; i++ {
		l[i+nx] = at(y, i)
	}
	return l
}
func tak(x, y v) v { return e("TODO") }
func rsh(x, y v) v {
	// TODO temporarily only rsh(int, []interface{}) is supported.
	n := x.(int)
	v := y.(l)
	m := len(v)
	if n == m {
		return v
	}
	r := make(l, n)
	for i := range r {
		r[i] = cp(v[i%m])
	}
	return r
	// TODO
}
func fil(x, y v) v { return e("TODO") }
func drp(x, y v) v { return e("TODO") }
func cut(x, y v) v { return e("TODO") }
func cst(x, y v) v { return e("TODO") }
func rnd(x, y v) v { return e("TODO") }
func fnd(x, y v) v { return e("TODO") }
func pik(x, y v) v { return e("TODO") }
func rfd(x, y v) v { return e("TODO") }
func atx(x, y v, a kt) v {
	if s, o := x.(s); o {
		return atx(a.at(s), y, a)
	}
	if y == nil {
		return x
	}
	xl, yl := false, false
	if ln(x) > 0 {
		xl = true
	}
	if ln(y) > 0 {
		yl = true
	}
	xdict, xd := md(x)
	ydict, yd := md(y)
	_ = xdict
	_ = ydict
	switch {
	case xl && yd:
		// TODO
	// TODO x verb, y adverb
	case xl || xd && yl:
		// TODO
	case xl:
		// TODO nested
	case xd:
		// TODO
	}
	// TODO call
	return e("TODO")
}
func cal(x, y v) v { return e("TODO") }
func bin(x, y v) v { return e("TODO") }
func rbn(x, y v) v { return e("TODO") }
func pak(x, y v) v { return e("TODO") }
func upk(x, y v) v { return e("TODO") }
func spl(x, y v) v { return e("TODO") }
func win(x, y v) v { return e("TODO") }
