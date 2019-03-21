// i interpret
package i

import (
	"math"
	"math/cmplx"
	"reflect"
)

func P(s s) v            { return prs(s) }
func E(l v, a map[v]v) v { return eva(l, kinit(a)) }

type (
	i  = int
	f  = float64
	fv = []f
	z  = complex128
	zv = []z
	s  = string
	sv = []s
	v  = interface{}
	l  = []v
	d  = dict
)
type rV = reflect.Value
type rT = reflect.Type

func rval(x v) rV { return reflect.ValueOf(x) }
func rtyp(x v) rT { return reflect.TypeOf(x) }

var rTb = rtyp(true)
var rTf = rtyp(0.0)
var rTz = rtyp(complex(0, 0))
var rTs = rtyp("")

type cpr interface {
	Copy() v
}

func cp(x v) v {
	switch t := x.(type) {
	case f:
		return x
	case z:
		return x
	case s:
		return x
	case l:
		r := make(l, len(t))
		for i := range r {
			r[i] = cp(t[i])
		}
		return r
	case d:
		r := t
		r.k = cp(t.k).(l)
		r.v = cp(t.v).(l)
		return r
	}
	if v, ok := x.(cpr); ok {
		return v.Copy()
	}
	v := rval(x)
	switch v.Kind() {
	case reflect.Slice:
		n := v.Len()
		r := reflect.MakeSlice(v.Type(), n, n)
		for i := 0; i < n; i++ {
			y := cp(v.Index(i).Interface())
			r.Index(i).Set(rval(y))
		}
		return r.Interface()
	case reflect.Chan, reflect.Interface, reflect.Ptr, reflect.UnsafePointer:
		// TODO: allow pointer, with or without deep copy, depending on type?
		e("type") // TODO: these types should be returned verbatim
	case reflect.Map, reflect.Struct:
		e("assert") // already converted to dict
	}
	return x
}

func e(s string) v { panic(s); return nil }

func ln(v v) int { // slice len
	r := rval(v)
	if r.Kind() == reflect.Slice {
		return r.Len()
	}
	return -1
}
func lz(l v) v { // zero element of a slice (returns nil for type l)
	return reflect.Zero(rtyp(l).Elem()).Interface()
}
func ls(v v) (l, rT) { // import list from any slice
	if v, ok := v.(l); ok {
		return v, nil
	}
	r := rval(v)
	if r.Kind() != reflect.Slice {
		e("type")
	}
	l := make(l, r.Len())
	for i := range l {
		l[i] = cp(r.Index(i).Interface())
	}
	return l, r.Type().Elem()
}
func sl(l l, et rT) v { // convert list back to slice with original element type
	if et == nil {
		return l
	}
	r := reflect.MakeSlice(reflect.SliceOf(et), len(l), len(l))
	for i := 0; i < len(l); i++ {
		r.Index(i).Set(rval(l[i]).Convert(et))
	}
	return r.Interface()
}
func uf(l l) (v, bool) { // convert a list to a uniform vector if possible
	if len(l) == 0 {
		return l, false
	}
	t := rtyp(l[0])
	for i := range l {
		if rtyp(l[i]) != t || ln(l[i]) >= 0 {
			return l, false
		}
	}
	return sl(l, t), true
}

type dict struct {
	k, v          l
	f, s, u, p, g bool // flipped, sorted, uniq, parted, grouped
	t             reflect.Type
}

func md(x interface{}) (d, bool) { // import maps and structs as dicts
	var d d
	switch m := x.(type) {
	case map[v]v:
		off := 0
		if h, ok := m["_"]; ok {
			hdr := h.(dict)
			d.f, d.s, d.u, d.p, d.g = hdr.f, hdr.s, hdr.u, hdr.p, hdr.g
			d.k = cp(hdr.k).(l)
			off = 1
		}
		n := len(m) - off
		if off == 0 {
			d.k = make(l, n)
			d.v = make(l, n)
			i := 0
			for k, v := range m {
				d.k[i] = cp(k)
				d.v[i] = cp(v)
				i++
			}
		} else {
			d.v = make(l, n)
			for i, k := range d.k {
				d.v[i] = cp(m[k])
			}
		}
		return d, true
	}
	// convert from user supplied map or struct

	v := rval(x)
	d.t = v.Type()
	if kind := v.Kind(); kind == reflect.Map || kind == reflect.Struct {
		n := 0
		if kind == reflect.Map {
			n = v.Len()
		} else {
			n = v.NumField()
		}
		d.k = make(l, n)
		d.v = make(l, n)
		if kind == reflect.Map {
			keys := v.MapKeys()
			for i, k := range keys {
				d.k[i] = cp(k.Interface())
				d.v[i] = cp(v.MapIndex(k).Interface())
			}
		} else {
			t := v.Type()
			for i := 0; i < n; i++ {
				d.k[i] = t.Field(i).Name
				d.v[i] = cp(v.Field(i).Interface())
			}
		}
		return d, true
	}
	return d, false
}
func md2(x, y interface{}) (dict, dict, bool) {
	dx, o := md(x)
	if !o {
		return dict{}, dict{}, false
	}
	dy, o := md(y)
	if !o {
		return dict{}, dict{}, false
	}
	return dx, dy, true
}
func (d dict) mp() interface{} { // convert dict back to original type
	if d.t == nil {
		r := make(map[v]v)
		r["_"] = dict{d.k, nil, d.f, d.s, d.u, d.p, d.g, nil}
		for i, k := range d.k {
			r[k] = d.v[i]
		}
		return r
	}

	// convert back to original map or struct type.
	v := reflect.New(d.t)
	v = v.Elem()
	if v.Kind() == reflect.Map {
		v = reflect.MakeMap(d.t)
		keytype := v.Type().Key()
		valtype := v.Type().Elem()
		for i, k := range d.k {
			rk := rval(k)
			if t := rk.Type(); t != keytype {
				rk = rk.Convert(t)
			}
			rv := rval(d.v[i])
			if t := rv.Type(); t != valtype {
				rv = rv.Convert(t)
			}
			v.SetMapIndex(rk, rv)
		}
		return v.Interface()
	} else if v.Kind() == reflect.Struct {
		for i, k := range d.k {
			f := v.FieldByName(rval(k).String())
			if f.Kind() == reflect.Slice {
				w := rval(d.v[i])
				if w.IsValid() == false {
					continue
				}
				sv := reflect.MakeSlice(f.Type(), w.Len(), w.Len())
				reflect.Copy(sv, w)
			} // TODO: make other types, that need it.
			f.Set(rval(d.v[i]))
		}
		return v.Interface()
	}
	return e("type")
}
func (d dict) at(key v) (int, v) {
	for i, k := range d.k {
		if k == key {
			return i, d.v[i]
		}
	}
	return -1, nil
}
func (d *dict) set(key, val v) {
	if i, _ := d.at(key); i < 0 {
		d.k = append(d.k, key)
		d.v = append(d.v, val)
	} else {
		d.v[i] = val
	}
}

func sy(v v) (sv, int, rT, bool) { // import any string or string slice to symbols
	switch t := v.(type) {
	case s:
		return sv{t}, -1, nil, true
	case sv:
		return t, len(t), nil, true
	}
	r := rval(v)
	if r.Kind() == reflect.String {
		return sv{r.String()}, -1, r.Type(), true
	} else if r.Kind() == reflect.Slice && reflect.Zero(r.Type().Elem()).Kind() == reflect.String {
		n := r.Len()
		u := make(sv, n, n)
		for i := range u {
			u[i] = r.Index(i).String()
		}
		return u, n, r.Type().Elem(), true
	}
	return nil, 0, nil, false
}
func ys(x sv, vec bool, eT rT) v { // convert strings back to orig type
	if !vec {
		return rval(x[0]).Convert(eT).Interface()
	}
	r := reflect.MakeSlice(reflect.SliceOf(eT), len(x), len(x))
	for i := range x {
		r.Index(i).Set(rval(x[i]).Convert(eT))
	}
	return r.Interface()
}

/*
function ktos(x, esc) {
	if (x.t != 3) { x = enlist(x); }
	var h = x.v.some(function(v){ return (v.v<32||v.v>127)&v.v!=9&v.v!=10; });
	if (h) { return "0x"+x.v.map(h2).join(""); }
	var r = x.v.map(function(k) { return String.fromCharCode(k.v); }).join("");
	return esc ? '"'+EC.reduce(function(r,p) { return r.split(p[0]).join(p[1]); }, r)+'"' : r;
}
func ktos(x, esc v) v {
	if ln(x) < 0 {
		x = enl(x)
	}
	r := rval(x)
	if r.Kind() != reflect.Slice {
	}
}
*/

/*
// function krange(x, f) { var r=[]; for(var z=0;z<x;z++) { r.push(f(z)); } return k(3,r); }
func krange(n int, f func(int) v) l {
	l := make(l, n)
	for i := 0; i < n; i++ {
		l[i] = f(i)
	}
	return l
}

*/
// function kmap (x, f) { return k(3, l(x).v.map(f)); }
func kmap(x v, f func(v) v) v {
	n := ln(x)
	if n < 0 {
		e("type")
	}
	var it, ot rT
	if t := rtyp(lz(x)); t != nil && t.Kind() != reflect.Interface {
		it = t
	}
	in, _ := ls(x) // rT is determined by result of f(x)
	l := make(l, n)
	for i := 0; i < n; i++ {
		l[i] = f(in[i])
		t := rval(l[i]).Type()
		if t != nil && i == 0 {
			ot = t
		} else if t != nil && ot != nil && t != ot {
			e("type") // f returns non-uniform on uniform input
		}
	}
	if it == nil || ot == nil {
		return l
	}
	r := reflect.MakeSlice(reflect.SliceOf(ot), n, n)
	for i := 0; i < n; i++ {
		r.Index(i).Set(rval(l[i]).Convert(ot))
	}
	return r.Interface()
}

/*
// function kzip (x, y, f) { return kmap(sl(x,y), function(z, i) { return f(z, y.v[i]); }); }
func kzip(x, y v, f func(v, v) v) v {
	return kmap(xxx, func(v v, i int) v {
		return f(v, at(y, i))
	})
}
*/

func some(l l, f func(v v) bool) bool {
	for _, i := range l {
		if f(i) {
			return true
		}
	}
	return false
}

func impl(v v, t reflect.Type) reflect.Method {
	if rtyp(v).Implements(t) {
		return t.Elem().Method(0)
	}
	return reflect.Method{}
}

func idx(v v) int {
	var n int
	var f float64
	switch w := v.(type) {
	case int:
		return w
	case float64:
		f, n = w, int(w)
	default:
		f = re(v)
		n = int(f)
	}
	if float64(n) != f {
		e("type") // rounding
	}
	return n
}

func re(v v) float64 {
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
	r := rval(v)
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

func at(L v, i int) v {
	switch t := L.(type) {
	case l:
		return cp(t[i])
	case fv:
		return t[i]
	case zv:
		return t[i]
	}
	if r := rval(L); r.Kind() != reflect.Slice {
		return e("type")
	} else {
		return cp(r.Index(i).Interface())
	}
}
func set(L v, i int, x v) {
	switch t := L.(type) {
	case l:
		t[i] = x
		return
	case fv:
		t[i] = x.(float64)
		return
	case zv:
		t[i] = x.(complex128)
		return
	}
	rval(L).Index(i).Set(rval(x))
}

type kt map[v]v

func (a kt) at(s s) v { return e("nyi") }
func kinit(a kt) kt {
	return kt(a)
}
