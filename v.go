package i

import (
	"math"
	"reflect"
	"sort"
)

// monadic verbs
func flp(x v) v {
	if d, o := md(x); o {
		d.f = !d.f
		return d.mp()
	}
	n, m := ln(x), -1
	if n < 0 {
		return x
	}
	var rl l
	for i := 0; i < n; i++ {
		rw := at(x, i)
		if nn := ln(rw); nn < 0 {
			return x
		} else if i == 0 {
			m = nn
			rl = make(l, m*n)
		} else if nn != m {
			return e("length")
		}
		for k := 0; k < m; k++ {
			rl[k*n+i] = at(rw, k)
		}
	}
	return cut(mul(n, (til(n-1))), rl)
}
func neg(x v) v { return nm(x, rneg, zneg, "Neg") }
func fst(v v) v {
	if d, o := md(v); o {
		return fst(d.v)
	}
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
func wer(x v) v { return e("nyi") } // →take, over asverb
func rev(x v) v {
	if d, ok := md(x); ok {
		k, u := rev(d.k).(l), rev(d.v).(l)
		d.k, d.v = k, u
		return d.mp()
	}
	x = cp(x)
	r := rval(x)
	if r.Kind() != reflect.Slice {
		return e("type")
	}
	n := r.Len()
	tmp := reflect.New(r.Type().Elem()).Elem()
	for i := n/2 - 1; i >= 0; i-- {
		j := n - 1 - i
		tmp.Set(r.Index(i))
		r.Index(i).Set(r.Index(j))
		r.Index(j).Set(tmp)
	}
	return r.Interface()
}
func asc(x v) v { return grade(true, x) }
func dsc(x v) v { return grade(false, x) }
func eye(x v) v {
	f, _, vec, _ := nv(x)
	if vec {
		return e("rank")
	}
	if f == nil {
		return e("type")
	}
	n := int(f[0])
	l := make(l, n)
	for i := range l {
		r := make(fv, n)
		r[i], l[i] = 1, r
	}
	return l
}
func grp(x v) v { return e("nyi") }
func not(x v) v { return nm(x, rnot, znot, "Not") }
func enl(x v) v {
	if d, o := x.(dict); o {
		return l{d}
	}
	v := rval(x)
	switch v.Kind() {
	case reflect.Func, reflect.Slice, reflect.Map:
		return l{x}
	}
	l := ms(v.Type(), 1)
	l.Index(0).Set(v)
	return l.Interface()
}
func is0(x v) v { return e("nyi") }
func cnt(x v) v {
	if d, o := md(x); o {
		return f(len(d.k))
	} else if n := ln(x); n >= 0 {
		return f(n)
	}
	return f(1)
}
func flr(x v) v { return nm(x, rflr, zflr, "Flr") }
func fmt(x v) v { return e("nyi") }
func fgn(x v) v { return e("nyi") }
func unq(x v) v {
	w, t := ls(x)
	r := make(l, 0)
	for i := range w {
		if !some(r, func(x v) bool { return mch(w[i], x) == 1.0 }) {
			r = append(r, cp(w[i]))
		}
	}
	return sl(r, t)
}
func evl(x v) v { return e("nyi") }

// dyadic verbs
func add(x, y v) v { return nd(x, y, radd, zadd, "Add") }
func sub(x, y v) v { return nd(x, y, rsub, zsub, "Sub") }
func mul(x, y v) v { return nd(x, y, rmul, zmul, "Mul") }
func div(x, y v) v { return nd(x, y, rdiv, zdiv, "Div") }
func mod(x, y v) v { return e("nyi") }
func mkd(x, y v) v { return e("nyi") }
func min(x, y v) v { return nd(x, y, rmin, zmin, "Min") }
func max(x, y v) v { return nd(x, y, rmax, zmax, "Max") }
func les(x, y v) v { x, y = sn2(x, y); return nd(x, y, rles, zles, "Les") }
func mor(x, y v) v { x, y = sn2(x, y); return nd(x, y, rmor, zmor, "Mor") }
func eql(x, y v) v { x, y = sn2(x, y); return nd(x, y, reql, zeql, "Eql") }
func mch(x, y v) v {
	if rtyp(x) != rtyp(y) {
		return 0.0
	}
	if xd, o := md(x); o {
		if yd, o := md(y); o {
			return min(mch(xd.k, yd.k), mch(xd.v, yd.v))
		}
		return e("assert")
	}
	if reflect.DeepEqual(x, y) {
		return 1.0
	}
	return 0.0
}
func cat(x, y v) v {
	if xd, yd, o := md2(x, y); o {
		for i := range yd.k {
			xd.set(yd.k[i], yd.v[i])
		}
		if xd.t != yd.t {
			xd.t = nil
		}
		return xd.mp()
	}
	xd, yd := false, false
	if dx, o := md(x); o {
		x, xd = dx, true
	}
	if dy, o := md(y); o {
		y, yd = dy, true
	}
	nx := ln(x)
	if nx < 0 {
		x, nx = enl(x), 1
	}
	ny := ln(y)
	if ny < 0 {
		y, ny = enl(y), 1
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
	r := make(l, nx+ny)
	for i := 0; i < nx; i++ {
		if i == 0 && xd {
			r[i] = x.(l)[0].(dict).mp()
		} else {
			r[i] = at(x, i)
		}
	}
	for i := 0; i < ny; i++ {
		if i == 0 && yd {
			r[nx] = y.(l)[0].(dict).mp()
		}
		r[i+nx] = at(y, i)
	}
	return r
}
func tak(x, y v) v { return e("nyi") }
func rsh(x, y v) v {
	if yd, o := md(y); o { // select from dict
		yd.v = atx(y, x, nil).(l)
		xl, _ := ls(x)
		yd.k = xl
		return yd.mp()
	}
	nx, ny := ln(x), ln(y)
	if ny < 0 {
		y, ny = enl(y), 1
	}
	if nx <= 0 {
		x, nx = enl(x), 1
	}
	xv := make(fv, nx)
	for i := range xv {
		u := at(x, i)
		switch t := u.(type) {
		case f:
			xv[i] = t
		case int:
			xv[i] = f(t)
		default:
			e("type")
		}
	}
	a, b, c := xv[0], xv[len(xv)-1], 0
	var rshr func(x, y v, i int) v
	rshr = func(x, y v, i int) v {
		nx, ny := ln(x), ln(y)
		return krange(int(xv[i]), func(z int) v {
			if i == nx-1 {
				c++
				return at(y, (c-1)%ny)
			}
			return rshr(x, y, i+1)
		})
	}
	if math.IsNaN(a) {
		if ny == 0 {
			return y
		}
		return cut(krange(int(math.Ceil(f(ny)/b)), func(z int) v { return z * int(b) }), y)
	} else if math.IsNaN(b) {
		return cut(krange(int(a), func(z int) v { return z * ny / int(a) }), y)
	}
	return rshr(x, y, 0)
}
func fil(x, y v) v { return e("nyi") }
func drp(x, y v) v {
	if d, o := md(y); o {
		d.k, d.v = drp(x, d.k).(l), drp(x, d.v).(l)
		return d.mp()
	}
	n := ln(y)
	if n <= 0 {
		return y
	}
	if ln(x) >= 0 {
		return e("length")
	}
	j := int(re(x))
	y = cp(y)
	if (j < 0 && j+n <= 0) || (j > 0 && n-j <= 0) {
		return ms(rtyp(y).Elem(), 0).Interface()
	}
	if j < 0 {
		return rval(y).Slice(0, n+j).Interface()
	}
	return rval(y).Slice(j, n).Interface()
}
func cut(x, y v) v {
	p := func(v v) int {
		n := -1
		switch t := v.(type) {
		case f:
			n = int(t)
		case int:
			n = t
		default:
			e("type")
		}
		if n < 0 {
			e("domain")
		}
		return n
	}
	return kzip(x, cat(drp(1, x), cnt(y)), func(a, b v) v {
		pa, pb := p(a), p(b)
		r := make(l, pb-pa)
		for i := pa; i < pb; i++ {
			r[i-pa] = at(y, i)
		}
		u, _ := uf(r)
		return u
	})
}
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
	switch {
	case xl && yd: // 3
		ydict.v = atx(x, cp(ydict.v), a).(l)
		return yd
	case rval(x).Kind() == reflect.Func && rval(y).Kind() == reflect.Func:
		return e("nyi") // x, verb, y adverb: // 4
	case (xl || xd) && yl: // 5
		return kmap(y, func(z v, i int) v { return atx(x, z, nil) })
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
	if ln(x) < 0 {
		x = enl(x)
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
