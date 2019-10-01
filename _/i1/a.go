// i interpret
package i

import (
	"math"
	"reflect"
)

func P(s s) v            { return prs(s) }
func E(l v, a map[v]v) v { return eva(l, kinit(a)) }

type (
	z  = complex128
	zv = []z
	s  = string
	sv = []s
	v  = interface{}
	l  = []v
)
type rV = reflect.Value
type rT = reflect.Type

func rval(x v) rV { return reflect.ValueOf(x) }
func rtyp(x v) rT { return reflect.TypeOf(x) }

var rTb = rtyp(true)
var rTf = rtyp(0.0)
var rTz = rtyp(complex(0, 0))
var rTs = rtyp("")
var rTv = rtyp(l{}).Elem()

type cpr interface {
	Copy() v
}

func cp(x v) v {
	if k := rval(x).Kind(); k < reflect.Array || k == reflect.String {
		return x
	}
	switch t := x.(type) {
	case l:
		r := make(l, len(t))
		for i := range r {
			r[i] = cp(t[i])
		}
		return r
	case dict:
		r := t
		r.k, r.v = cp(t.k).(l), cp(t.v).(l)
		return r
	}
	if v, ok := x.(cpr); ok {
		return v.Copy()
	}
	v := rval(x)
	switch v.Kind() {
	case reflect.Slice:
		if v.IsNil() {
			return v.Interface()
		}
		n := v.Len()
		r := reflect.MakeSlice(v.Type(), n, n)
		for i := 0; i < n; i++ {
			y := cp(v.Index(i).Interface())
			r.Index(i).Set(rval(y))
		}
		return r.Interface()
	case reflect.Map, reflect.Struct:
		d, o := md(x)
		if !o {
			e("type")
		}
		return d.mp()
	}
	return x // Ptr, Chan, Interface, UnsafePointer are not copied
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
		if rv := rval(l[i]); rv.Type().ConvertibleTo(et) == false {
			return l
		} else {
			r.Index(i).Set(rv.Convert(et))
		}
	}
	return r.Interface()
}
func uf(l l) v { // convert a list to a uniform vector if possible
	if len(l) == 0 {
		return l
	}
	t := rtyp(l[0])
	for i := range l {
		if rtyp(l[i]) != t || ln(l[i]) >= 0 {
			return l
		}
	}
	if t.Kind() == reflect.Slice {
		return l
	}
	return sl(l, t)
}
func ms(eT rT, n int) rV { // make slice from element type, but lists from slices
	if eT == nil || eT.Kind() == reflect.Slice {
		return rval(make(l, n))
	}
	return reflect.MakeSlice(reflect.SliceOf(eT), n, n)
}
func isl(x v) bool { _, o := x.(l); return o } // is list
func iss(x v) bool { _, o := x.(s); return o } // is string
func isd(x v) bool { // is dict
	switch x.(type) {
	case dict:
		return true
	case [2]l:
		return true
	}
	if k := rval(x).Kind(); k == reflect.Struct || k == reflect.Map {
		return true
	}
	return false
}
func ex(x v, dst rT) v { // export to dst type
	if t := rtyp(x); t == dst {
		return x
	} else if dst == rTv {
		return x
	}
	switch dst.Kind() {
	case reflect.Slice:
		eT := dst.Elem()
		if k := eT.Kind(); k <= reflect.Complex128 {
			c, vec, _ := nv(x)
			if !vec {
				return e("type")
			}
			return vn(c, true, eT)
		} else if k == reflect.String {
			n := ln(x)
			sv := reflect.MakeSlice(dst, n, n)
			for i := 0; i < n; i++ {
				sv.Index(i).Set(rval(x).Index(i))
			}
			return sv.Interface()
		} else if k == reflect.Struct { // []struct{...} ←→ l{dict}
			n := ln(x)
			if n < 0 {
				return nil
			}
			sv := reflect.MakeSlice(dst, n, n)
			for i := 0; i < n; i++ {
				sv.Index(i).Set(rval(ex(at(x, i), eT)))
			}
			return sv.Interface()
		}
		return e("type")
	case reflect.Map, reflect.Struct:
		d, o := md(x)
		if !o {
			return e("type")
		}
		var r rV
		if dst.Kind() == reflect.Struct { // TODO: error on extra fields?
			r = reflect.New(dst).Elem()
			for i := 0; i < dst.NumField(); i++ {
				f := dst.Field(i)
				k, u := d.at(f.Name)
				if k < 0 {
					continue
				}
				u = ex(u, f.Type)
				r.Field(i).Set(rval(u))
			}
		} else {
			r = reflect.MakeMap(dst)
			kT := dst.Key()
			eT := dst.Elem()
			for i := range d.k {
				r.SetMapIndex(rval(ex(d.k[i], kT)), rval(ex(d.v[i], eT)))
			}
		}
		return r.Interface()
	default:
		if dst.Kind() <= reflect.Complex128 {
			c, vec, _ := nv(x)
			if vec {
				e("type")
			}
			return vn(c, false, dst)
		}
		return rval(x).Convert(dst).Interface()
	}
}

type dict struct {
	k, v l
	t    reflect.Type // orig type
}

func md(x v) (dict, bool) { // import maps and structs as dicts
	if m, o := x.([2]l); o {
		return dict{cp(m[0]).(l), cp(m[1]).(l), nil}, true
	} else if d, o := x.(dict); o {
		return d, true
	}
	var d dict
	v := rval(x)
	d.t = v.Type()
	if kind := v.Kind(); kind == reflect.Map || kind == reflect.Struct {
		n := 0
		if kind == reflect.Map {
			n = v.Len()
		} else {
			n = v.NumField()
		}
		d.k, d.v = make(l, n), make(l, n)
		if kind == reflect.Map {
			keys := v.MapKeys()
			for i, k := range keys {
				d.k[i], d.v[i] = cp(k.Interface()), cp(v.MapIndex(k).Interface())
			}
		} else {
			t := v.Type()
			j := 0
			for i := 0; i < n; i++ {
				u := v.Field(i).Interface()
				if rv := rval(u); rv.Kind() == reflect.Slice && rv.IsNil() {
					continue // skip nil slices
				}
				if d, o := md(u); o {
					u = d
				} else {
					u = cp(u)
				}
				d.k[j], d.v[j] = t.Field(i).Name, u
				j++
			}
			d.k, d.v = d.k[:j], d.v[:j]
		}
		return d, true
	}
	return d, false
}
func md2(x, y v) (dict, dict, bool) {
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
func (d dict) mp() v { // convert dict back to original type
	if d.t == nil {
		return [2]l{d.k, d.v}
	}
	return ex([2]l{d.k, d.v}, d.t)
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
		d.k, d.v = append(d.k, key), append(d.v, val)
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
	case l:
		u := uf(t)
		if sv, o := u.(sv); o {
			return sv, len(sv), nil, true
		}
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
func idx(v v) int {
	switch w := v.(type) {
	case z:
		return int(real(w))
	}
	r := rval(v)
	if k := r.Kind(); k == reflect.Bool {
		if r.Bool() {
			return 1
		}
		return 0
	} else if k < reflect.Uint {
		return int(r.Int())
	} else if k < reflect.Uintptr {
		return int(r.Uint())
	}
	return int(r.Float())
}
func pidx(v v) int { // to positive int
	n := idx(v)
	if n < 0 {
		e("range")
	}
	return n
}

func at(L v, i int) v {
	switch t := L.(type) {
	case l:
		return cp(t[i])
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
	case zv:
		t[i] = x.(complex128)
		return
	}
	rval(L).Index(i).Set(rval(x))
}

type curry func(...v) v // curry() reports it number of arguments
//type nfn func(...v) v        // func that knows it's name
//func (f nfn) String() string { v := f(); return v.(s) }
func ni1(x v) v    { return e("undefined monad") }
func ni2(x, y v) v { return e("undefined dyad") }

func kinit(a map[v]v) map[v]v {
	if len(a) > 0 {
		return a
	} else if a == nil {
		a = make(map[v]v)
	}
	a["doc"] = doc
	type v2 [2]v
	vtab := map[s]v2{
		"+": v2{flp, add}, "⍉": v2{flp, ni2},
		"-": v2{neg, sub},
		"*": v2{fst, mul}, "×": v2{ni1, mul},
		"%": v2{inv, div}, "÷": v2{inv, div},
		"√": v2{sqr, nrt},
		"ℜ": v2{zre, ni2}, "ℑ": v2{zim, ni2},
		"‖": v2{abs, rct},
		"°": v2{deg, pol}, "𝜑": v2{rad, prd},
		"⍣": v2{exp, pow}, "⍟": v2{log, lgn},
		"!": v2{til, mkd}, "⍳": v2{til, ni2},
		"&": v2{wer, min},
		"⌊": v2{flr, min}, "⌈": v2{cil, max},
		"|": v2{rev, max}, "⌽": v2{rev, ni2},
		"<": v2{asc, les}, ">": v2{dsc, mor},
		"⍋": v2{asc, ni2}, "⍒": v2{dsc, ni2},
		"=": v2{grp, eql}, "^": v2{is0, ept},
		"~": v2{not, mch}, "≡": v2{ni1, mch},
		",": v2{enl, cat}, "_": v2{flr, cut},
		"#": v2{cnt, tak}, "⍴": v2{cnt, rsh},
		"↑": v2{ni1, tak}, "↓": v2{ni1, cut},
		"$": v2{fmt, cst}, ".": v2{evl, cal},
		"?": v2{unq, fnd}, "@": v2{typ, atx},
	}
	for _s, _u := range vtab {
		s, u := _s, _u
		var monad func(v) v
		if rtyp(u[0]).NumIn() == 2 {
			f := u[0].(func(v, map[v]v) v) // evl
			monad = func(x v) v { return f(x, a) }
		} else {
			monad = u[0].(func(v) v)
		}
		var dyad func(v, v) v
		if rtyp(u[1]).NumIn() == 3 {
			f := u[1].(func(v, v, map[v]v) v) // atx, cal
			dyad = func(x, y v) v { return f(x, y, a) }
		} else {
			dyad = u[1].(func(v, v) v)
		}
		a[s+":"] = monad        // monad: "+:" is always func(v) v
		a[s+s] = dyad           // dyad: "++" is always func(v, v) v
		a[s] = func(w ...v) v { // ambivalent 1 or 2 args
			if len(w) == 1 {
				return monad(w[0])
			} else if len(w) == 2 {
				return dyad(w[0], w[1])
			}
			return e("args")
		}
	}

	type v4 [4]v
	atab := map[s]v4{
		//       m|x  d|x  xm|y xd|y
		"'":  v4{ech, ecd, ecd, ecd},
		"¨":  v4{ech, ecd, ecd, ecd},
		"':": v4{nil, ecp, nil, eci},
		"⍨":  v4{nil, ecp, nil, eci},
		"/:": v4{nil, nil, ecr, ecr},
		"⌿":  v4{nil, nil, ecr, ecr},
		`\:`: v4{nil, nil, ecl, ecl},
		`⍀`:  v4{nil, nil, ecl, ecl},
		"/":  v4{fix, ovr, whl, ovi},
		`\`:  v4{sfx, scn, swl, sci},
	}
	for _s, _u := range atab {
		s, u := _s, _u
		a[s] = func(f v) v {
			if rval(f).Kind() != reflect.Func {
				if s == "/" { // jon and spl are special, f is a string.
					return func(x v) v { return jon(f, x) }
				} else if s == `\` {
					return func(x v) v { return spl(f, x) }
				}
				return e("type")
			}
			return func(w ...v) v {
				cs := 1
				if len(w) < 1 || len(w) > 2 {
					return e("nargs")
				} else if len(w) == 2 {
					cs += 2
				}
				if cf, o := f.(curry); o && cf().(int) == 1 {
					cs--
				} else if t := rval(f).Type(); t.IsVariadic() == false && t.NumIn() == 1 {
					cs--
				}
				g := u[cs]
				if g == nil {
					e(s + ":argtype")
				}
				in := make([]rV, len(w)+2)
				in[0] = rval(f)
				for i := range w {
					in[i+1] = rval(w[i])
				}
				in[len(in)-1] = rval(a)
				r := rval(g).Call(in)
				return r[0].Interface()
			}
		}
	}
	for k, u := range map[s]v{
		"pi": math.Pi, "π": complex(math.Pi, 0), "𝜀": complex(1E-14, 0),
		"jon": jon, "num": num, "odo": odo,
		"inf": complex(math.Inf(1), 0), "∞": complex(math.Inf(1), 0), "nan": complex(math.NaN(), 0), "ø": complex(math.NaN(), 0),
		"mod": mod, "sqr": sqr, "pow": pow, "exp": exp, "log": log, "lgn": lgn, "nrt": nrt,
		"abs": abs, "deg": deg, "rad": rad, "re": zre, "im": zim, "con": con, "pol": pol, "prd": prd, "rct": rct,
		"ln": ln,
	} {
		a[k] = u
	}
	a["int"] = map[z]v{
		0: int(0), 1: false, 8: int8(0), 16: int16(0), 32: int32(0), 64: int64(0),
		-1: uint(0), -8: uint8(0), -16: uint16(0), -32: uint32(0), -64: uint64(0),
	}
	return a
}

const doc = `
+⍉  flp     add      ⍣exp|pow 
-   neg     sub      ⍟log|lgn
*×  fst     mul      √sqr|nrt        
%÷  inv     div       sin cos tan
!⍳  til,odo mkd      ‖abs|rct 
&⌊  wer     min      𝜑rad|prd
|⌽⌈ rev     max      °deg|pol
<⍋  asc     les      ℜre ℑim  con
>⍒  dsc     mor       odo jon num    
=   grp,eye eql
~≡  not     mch      ∞inf ønan πpi  𝜀
,   enl     cat
^   is0     ept
#⍴↑ cnt     tak,rsh
_⌊↓ flr     drp,cut
$   fmt     dst
?   unq,rng fnd,rng
@   typ     atx
.   evl     cal      
                                                                     
    mv/nv dv    l-mv  l-dv        
'   ech   ecd   ecd   ecd   ¨     
':   -    ecp    -    eci   ⍨     
/:   -     -    ecr   ecr   ⌿     
\:   -     -    ecl   ecl   ⍀     
/   fidx  ovr   whl   ovi         
\   sfx   scn   swl   sci`
