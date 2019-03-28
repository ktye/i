package i

import (
	"math"
	"math/rand"
	"reflect"
	"sort"
)

// monadic verbs
func flp(x v) v { // +x ⍉x flip
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
func neg(x v) v { return nm(x, rneg, zneg, "Neg") } // -x negate
func fst(v v) v { // *x first
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
func sqr(x v) v { return nm(x, rsqr, zsqr, "Sqr") } // √x sqrt
func inv(x v) v { return nm(x, rinv, zinv, "Inv") } // %x inverse
func abs(x v) v { return nm(x, rabs, zabs, "Abs") } // ¯x absolute value
func til(x v) v { // !x ⍳x iota
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
func odo(x v) v { // !l odometer
	inc := func(idx, shp []i) {
		for i := len(idx) - 1; i >= 0; i-- {
			idx[i]++
			if idx[i] < shp[i] {
				break
			}
			idx[i] = 0
		}
	}
	n := ln(x)
	if n < 0 {
		return e("domain")
	}
	shp, idx, r, m := make([]int, n), make([]int, n), make(l, n), 1
	for i := range shp {
		n = pi(at(x, i))
		if n == 0 {
			return e("domain")
		}
		shp[i] = n
		m *= n
	}
	for i := range shp {
		r[i] = make(fv, m)
	}
	for j := 0; j < m; j++ {
		for i := range r {
			r[i].(fv)[j] = f(idx[i])
		}
		inc(idx, shp)
	}
	return r
}
func wer(x v) v { // &x ⍸x where
	nx := ln(x)
	if nx < 0 {
		x, nx = enl(x), 1
	}
	xi := make([]int, nx)
	n := 0
	for i := range xi {
		xi[i] = pi(at(x, i))
		n += xi[i]
	}
	r := make(fv, n)
	j := 0
	for i := range xi {
		for k := 0; k < xi[i]; k++ {
			r[j] = f(i)
			j++
		}
	}
	return r
}
func rev(x v) v { // |x ⌽x reverse
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
func asc(x v) v { return grade(true, x) }  // <x ⍋x grade up
func dsc(x v) v { return grade(false, x) } // >x ⍒x grade down
func eye(x v) v { // =x unit matrix
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
func grp(x v) v { // =x ⌸x group
	n := ln(x)
	if n <= 0 {
		return e("type")
	}
	d := dict{}
	d.k, _ = ls(unq(x))
	d.v = make(l, len(d.k))
	m := make(map[v]int)
	for i := range d.v {
		d.v[i] = fv{}
		m[d.k[i]] = i
	}
	for i := 0; i < n; i++ {
		u := at(x, i)
		j := m[u]
		w := d.v[j].(fv)
		w = append(w, f(i))
		d.v[j] = w
	}
	return d.mp()
}
func not(x v) v { return nm(x, rnot, znot, "Not") } // ~x not
func enl(x v) v { // ,x enlist
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
func is0(x v) v { // ^x isnil, isnan
	if x == nil {
		return 1.0
	} else if s, n, _, o := sy(x); o {
		if n < 0 && s[0] == "" {
			return 1.0
		} else if n >= 0 {
			r := make(fv, len(s))
			for i := range s {
				if s[i] == "" {
					r[i] = math.NaN()
				}
			}
			x = r
		}
	}
	return nm(x, ris0, zis0, "Is0")
}
func exp(x v) v { return nm(x, rexp, zexp, "Exp") } // ⍣x exponential
func log(x v) v { return nm(x, rlog, zlog, "Log") } // ⍟x logarithm
func cnt(x v) v { // #x ⍴x count, length
	if d, o := md(x); o {
		return f(len(d.k))
	} else if n := ln(x); n >= 0 {
		return f(n)
	}
	return f(1)
}
func flr(x v) v { return nm(x, rflr, zflr, "Flr") } // _x ⌊x floor
func fmt(x v) v { return e("nyi") }                 // $x ⍕x format
func rng(x v) v { // ?x random uniform, ?-x normal ?z bi-normal
	xf, xz, vec, t := nv(x)
	if vec {
		return e("domain")
	}
	if xf == nil {
		r := make(zv, int(math.Floor(real(xz[0])))) // complex(n,0)
		for i := range r {
			r[i] = complex(rand.NormFloat64(), rand.NormFloat64())
		}
		if t == rTz {
			return r
		}
		s := ms(t, len(r))
		for i := range r {
			s.Index(i).Set(rval(r[i]).Convert(t))
		}
		return s.Interface()
	}
	f, norm := xf[0], false
	if f < 0 {
		f, norm = -f, true
	}
	r := make(fv, int(math.Floor(f)))
	for i := range r {
		if norm {
			r[i] = rand.NormFloat64()
		} else {
			r[i] = rand.Float64()
		}
	}
	if t == rTf {
		return r
	}
	s := ms(t, len(r))
	for i := range r {
		s.Index(i).Set(rval(r[i]).Convert(t))
	}
	return s.Interface()
}
func unq(x v) v { // ?x ∪x uniq
	w, t := ls(x)
	r := make(l, 0)
	for i := range w {
		if !some(r, func(x v) bool { return mch(w[i], x) == 1.0 }) {
			r = append(r, cp(w[i]))
		}
	}
	return sl(r, t)
}
func typ(x v) v { return e("nyi") } // @x type of
func evl(x v) v { return e("nyi") } // .x ⍎x evaluate

// dyadic verbs
func add(x, y v) v { return nd(x, y, radd, zadd, "Add") } // x+y add
func sub(x, y v) v { return nd(x, y, rsub, zsub, "Sub") } // x-y substract
func mul(x, y v) v { return nd(x, y, rmul, zmul, "Mul") } // x*x x×y multiply
func div(x, y v) v { return nd(x, y, rdiv, zdiv, "Div") } // x%y x÷y divide
func mod(x, y v) v { return nd(x, y, rmod, zmod, "Mod") } // x!y modulo
func mkd(x, y v) v { // xl!yl make dictionary
	a, _ := ls(x)
	b, _ := ls(y)
	if len(a) != len(b) {
		return e("length")
	}
	return dict{k: a, v: b}.mp()
}
func min(x, y v) v { return nd(x, y, rmin, zmin, "Min") }                   // x&y x⌊y minimum
func max(x, y v) v { return nd(x, y, rmax, zmax, "Max") }                   // x|y x⌈y maximum
func les(x, y v) v { x, y = sn2(x, y); return nd(x, y, rles, zles, "Les") } // x<y less than
func mor(x, y v) v { x, y = sn2(x, y); return nd(x, y, rmor, zmor, "Mor") } // x>y more than
func eql(x, y v) v { x, y = sn2(x, y); return nd(x, y, reql, zeql, "Eql") } // x=y equal
func pow(x, y v) v { return nd(x, y, rpow, zpow, "Pow") }                   // x⍣y power
func mch(x, y v) v { // x~y x≡y match
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
func cat(x, y v) v { // x,y catenate
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
func ept(x, y v) v { // x^y except
	nx, ny := ln(x), ln(y)
	if nx < 0 {
		x = til(x)
		nx = ln(x)
	}
	if ny < 0 {
		y, ny = enl(y), 1
	}
	r := make(l, 0, nx)
	m := make(map[v]bool)
	for i := 0; i < ny; i++ {
		m[at(y, i)] = true
	}
	for i := 0; i < nx; i++ {
		if u := at(x, i); !m[u] {
			r = append(r, u)
		}
	}
	if t := rtyp(x); t.Kind() == reflect.Slice {
		return sl(r, t.Elem())
	}
	return r
}
func tak(x, y v) v { // x#y take
	// nyi: 5,8,9: function, verb, adverb
	if d, o := md(y); o {
		k := tak(x, d.k)
		u := tak(x, d.v)
		d.k, d.v = k.(l), u.(l)
		return d.mp()
	}
	ny := ln(y)
	if ny <= 0 {
		y, ny = enl(y), 1
	}
	p := re(x)
	n := int(p)
	if f(n) != p {
		return e("type")
	}
	a := 0
	if n < 0 {
		a = ny + n
		n = -n
	}
	r := make(l, n)
	for i := range r {
		j := (a + i) % ny
		if j < 0 {
			j += ny
		}
		r[i] = at(y, j)
	}
	return sl(r, rtyp(y).Elem())
}
func rsh(x, y v) v { // x#y x⍴y reshape
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
func drp(x, y v) v { // x_y x↓y drop
	if d, o := md(y); o {
		nx := ln(x)
		if nx < 0 {
			d.k, d.v = drp(x, d.k).(l), drp(x, d.v).(l)
		} else {
			d.k = ept(d.k, x).(l)
			d.v = atx(y, d.k, nil).(l)
		}
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
func cut(x, y v) v { // x_y cut
	return kzip(x, cat(drp(1, x), cnt(y)), func(a, b v) v {
		pa, pb := pi(a), pi(b)
		r := make(l, pb-pa)
		for i := pa; i < pb; i++ {
			r[i-pa] = at(y, i)
		}
		u, _ := uf(r)
		return u
	})
}
func cst(x, y v) v { return e("nyi") } // x$y x⌶y cast
func rnd(x, y v) v { // x?y random, roll, -x?y deal
	ff := re(x)
	n := int(ff)
	if f(n) != ff {
		return e("type")
	}
	ny := ln(y)
	if ny < 0 {
		y = til(y)
		ny = ln(y)
	}
	var r l
	ll, rT := ls(cp(y))
	if n < 0 {
		if -n > ny {
			e("size")
		}
		rand.Shuffle(ny, func(i, j int) { ll[i], ll[j] = ll[j], ll[i] })
		r = make(l, -n)
		for i := range r {
			r[i] = ll[i]
		}
	} else {
		if n > ny {
			e("size")
		}
		r = make(l, n)
		for i := range r {
			r[i] = ll[int(math.Round(f(ny)*rand.Float64()))]
		}
	}
	return sl(r, rT)
}
func fnd(x, y v) v { // l?a xl?yl find
	nx := ln(x)
	if nx < 0 {
		return e("length")
	}
	if d, o := md(y); o {
		if d.f {
			for i := range d.k {
				d.v[i] = fnd(x, d.v[i])
			}
			return d.mp()
		}
		for i := range d.k {
			u := d.v[i]
			d.v[i] = f(nx)
			for j := 0; j < nx; j++ {
				if mch(at(x, j), u) == 1.0 {
					d.v[i] = f(j)
					break
				}
			}
		}
		return d.mp()
	}
	ny, vec := ln(y), true
	if ny < 1 {
		y, ny, vec = enl(y), 1, false
	}
	r := make(fv, ny)
	for i := range r {
		for j := 0; j < nx; j++ {
			r[i] = f(nx)
			if mch(at(x, j), at(y, i)) == 1.0 {
				r[i] = f(j)
				break
			}
		}
	}
	if !vec {
		return r[0]
	}
	return r // nyi: extension to rectangular arrays
}
func atx(x, y v, a map[v]v) v { // x@y at, index
	if s, o := x.(s); o {
		return atx(atx(a, s, nil), y, a) // 1
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
	return cal(x, enl(y), a)
}
func cal(x, y v, a map[v]v) v { // x.y call
	if x == nil {
		return e("call nil")
	}
	if sy, o := x.(s); o {
		f := lup(a, sy)
		if f == nil {
			return e("nil:" + sy)
		}
		return cal(f, y, a)
	}
	// TODO other cases
	f := rval(x)
	if f.Kind() != reflect.Func {
		return e("nyi:call:" + f.Kind().String())
	}
	var in, r []rV
	if yl, o := y.(l); o {
		in = make([]rV, len(yl))
		for i := range in {
			in[i] = rval(yl[i])
		}
	}
	r = f.Call(in)
	if len(r) == 0 {
		return nil
	}
	return r[0].Interface()
}
func bin(x, y v) v       { return e("nyi") }
func rbn(x, y v) v       { return e("nyi") }
func pak(x, y v) v       { return e("nyi") }
func upk(x, y v) v       { return e("nyi") }
func spl(x, y v) v       { return e("nyi") }
func amd(x, y, z, w v) v { return e("nyi") } // amend
func dmd(x, y, z, w v) v { return e("nyi") } // dmend
func win(x, y v) v       { return e("nyi") }

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
