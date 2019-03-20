package i

import (
	"reflect"
	"sort"
)

// monadic verbs
func flp(x v) v { return e("nyi") }
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
func odo(x v) v { return e("nyi") } // →impl
func wer(x v) v { return e("nyi") }
func rev(x v) v { return e("nyi") }
func asc(x v) v { return grade(true, x) }
func dsc(x v) v { return grade(false, x) }
func eye(x v) v { return e("nyi") }
func grp(x v) v { return e("nyi") }
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
func is0(x v) v { return e("nyi") }
func cnt(x v) v { return e("nyi") }
func flr(x v) v { return nm(x, rflr, zflr, "Flr") }
func fmt(x v) v { return e("nyi") }
func fgn(x v) v { return e("nyi") }
func unq(x v) v { return e("nyi") }
func evl(x v) v { return e("nyi") }

// dyadic verbs
func add(x, y v) v { return nd(x, y, radd, zadd, "Add") }
func sub(x, y v) v { return nd(x, y, rsub, zsub, "Sub") }
func mul(x, y v) v { return nd(x, y, rmul, zmul, "Mul") }
func div(x, y v) v { return nd(x, y, rdiv, zdiv, "Div") }
func mod(x, y v) v { return e("nyi") }
func mkd(x, y v) v { return e("nyi") }
func min(x, y v) v { return nd(x, y, rmin, zmin, "Min") }
func max(x, y v) v { return nd(x, y, rmax, zmax, "Max") } // cast to bool?
func les(x, y v) v { return nd(x, y, rles, zles, "Les") } // ?
func mor(x, y v) v { return nd(x, y, rmor, zmor, "Mor") } // ?
func eql(x, y v) v { return nd(x, y, reql, zeql, "Eql") } // ?
func mch(x, y v) v { return e("nyi") }
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
func tak(x, y v) v { return e("nyi") }
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
func fil(x, y v) v { return e("nyi") }
func drp(x, y v) v { return e("nyi") }
func cut(x, y v) v { return e("nyi") }
func cst(x, y v) v { return e("nyi") }
func rnd(x, y v) v { return e("nyi") }
func fnd(x, y v) v { return e("nyi") }
func pik(x, y v) v { return e("nyi") }
func rfd(x, y v) v { return e("nyi") }
func atx(x, y v, a kt) v {
	if s, o := x.(s); o {
		return atx(a.at(s), y, a) // 1
	}
	if y == nil {
		return x // 2
	}
	xl, yl := false, false
	nx, ny := ln(x), ln(y)
	if nx > 0 {
		xl = true
	}
	if ny > 0 {
		yl = true
	}
	xdict, xd := md(x)
	ydict, yd := md(y)
	_ = xdict
	_ = ydict
	switch {
	case xl && yd: // 3
		ydict.v = atx(x, cp(ydict.v), a).(l)
		return yd
	case rval(x).Kind() == reflect.Func && rval(y).Kind() == reflect.Func:
		return e("nyi") // x, verb, y adverb: // 4
	case (xl || xd) && yl: // 5
		return kmap(y, func(z v) v { return atx(x, z, nil) })
	case xl: // 6
		// TODO other checks: (y.t > 1 || y.v < 0 || y.v >= len(x) || y.v%1 != 0) ? NA
		i := idx(y)
		if i < 0 || i >= nx {
			return e("range")
		}
		return at(x, i)
	case xd: // 7
		_, r := xdict.at(y)
		return cp(r)
	}
	// TODO call // 8
	return e("nyi")
}
func cal(x, y v) v { return e("nyi") }
func bin(x, y v) v { return e("nyi") }
func rbn(x, y v) v { return e("nyi") }
func pak(x, y v) v { return e("nyi") }
func upk(x, y v) v { return e("nyi") }
func spl(x, y v) v { return e("nyi") }
func win(x, y v) v { return e("nyi") }

// adverbs
// nyi

func grade(up bool, x v) v {
	if d, o := md(x); o {
		return atx(d.k, grade(up, d.v), nil)
	}
	x = cp(x)
	switch t := x.(type) {
	case fv:
		x = sort.Float64Slice(t)
	case sv:
		x = sort.StringSlice(t)
	}
	if d, o := x.(sort.Interface); o {
		if !up {
			d = sort.Reverse(d)
		}
		i := make(fv, d.Len())
		for n := range i {
			i[n] = f(n)
		}
		sort.Sort(grades{d, i})
		return i
	}
	if l, o := x.(l); o {
		u, o := uf(l)
		if !o {
			return e("type")
		}
		x = u
	}
	fv, zv, vec, _ := nv(x)
	switch {
	case vec && fv != nil:
		return grade(up, fv)
	case vec && zv != nil:
		return grade(up, zv)
	}

	println(rtyp(x).String())
	return e("type")
}

type grades struct {
	sort.Interface
	idx fv
}

func (s grades) Swap(i, j int) {
	s.Interface.Swap(i, j)
	s.idx[i], s.idx[j] = s.idx[j], s.idx[i]
}
