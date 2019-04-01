package i

import (
	"math"
	"math/cmplx"
	"math/rand"
	"reflect"
	"sort"
)

// monadic verbs
func flp(x v) v { // +x ⍉x flip
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
func neg(x v) v { return nm(x, func(x z) z { return -x }) } // -x negate
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
func con(x v) v { return nm(x, func(x z) z { return cmplx.Conj(x) }) }              // con x conjugate complex
func sqr(x v) v { return nm(x, func(x z) z { return cmplx.Sqrt(x) }) }              // √x sqrt
func inv(x v) v { return nm(x, func(x z) z { return 1 / x }) }                      // %x inverse
func abs(x v) v { return nm(x, func(x z) z { return complex(cmplx.Abs(x), 0) }) }   // ‖x absolute value
func ang(x v) v { return nm(x, func(x z) z { return complex(cmplx.Phase(x), 0) }) } // ang x complex phase [-π,π]
func deg(x v) v {
	return nm(x, func(x z) z { // ang x complex phase [0,360]
		p := cmplx.Phase(x) / math.Pi * 360.0
		if p < 0 {
			p += 360.0
		}
		return complex(p, 0)
	})
}
func zre(x v) v { return nm(x, func(x z) z { return complex(real(x), 0) }) } // ℜx real part
func zim(x v) v { return nm(x, func(x z) z { return complex(imag(x), 0) }) } // ℑx real part
func til(x v) v { // !x ⍳x iota
	if d, ok := md(x); ok {
		return d.k
	}
	z, vec, t := nv(x)
	if vec {
		e("type")
	}
	if real(z[0]) < 0 || imag(z[0]) != 0 {
		return e("domain") // !-n
	}
	r := make(zv, int(real(z[0])))
	for i := range r {
		r[i] = zi(i)
	}
	return vn(r, true, t)
}
func odo(x v) v { // !l odometer
	inc := func(idx, shp []int) {
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
		n = pidx(at(x, i))
		if n == 0 {
			return e("domain")
		}
		shp[i] = n
		m *= n
	}
	for i := range shp {
		r[i] = make(zv, m)
	}
	for j := 0; j < m; j++ {
		for i := range r {
			r[i].(zv)[j] = zi(idx[i])
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
		xi[i] = pidx(at(x, i))
		n += xi[i]
	}
	r := make(zv, n)
	j := 0
	for i := range xi {
		for k := 0; k < xi[i]; k++ {
			r[j] = zi(i)
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
	f, vec, _ := nv(x)
	if vec {
		return e("rank")
	}
	if imag(f[0]) != 0 || real(f[0]) < 0 {
		return e("type")
	}
	n := int(real(f[0]))
	l := make(l, n)
	for i := range l {
		r := make(zv, n)
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
		d.v[i] = zv{}
		m[d.k[i]] = i
	}
	for i := 0; i < n; i++ {
		u := at(x, i)
		j := m[u]
		w := d.v[j].(zv)
		w = append(w, zi(i))
		d.v[j] = w
	}
	return d.mp()
}
func not(x v) v { return nm(x, func(x z) z { return zter(x == 0, 1, 0) }) } // ~x not
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
		return zi(1)
	} else if s, n, _, o := sy(x); o {
		if n < 0 && s[0] == "" {
			return zi(1)
		} else if n >= 0 {
			r := make(zv, len(s))
			for i := range s {
				if s[i] == "" {
					r[i] = cmplx.NaN()
				}
			}
			x = r
		}
	}
	z0 := func(x z) z {
		if cmplx.IsNaN(x) {
			return zi(1)
		}
		return zi(0)
	}
	return nm(x, z0)
}
func exp(x v) v { return nm(x, func(x z) z { return cmplx.Exp(x) }) } // ⍣x exponential
func log(x v) v { return nm(x, func(x z) z { return cmplx.Log(x) }) } // ⍟x logarithm
func cnt(x v) v { // #x ⍴x count, length
	if d, o := md(x); o {
		return zi(len(d.k))
	} else if n := ln(x); n >= 0 {
		return zi(n)
	}
	return zi(1)
}
func flr(x v) v { return nm(x, func(x z) z { return complex(math.Floor(real(x)), 0) }) } // _x ⌊x floor
func cil(x v) v { return nm(x, func(x z) z { return complex(math.Ceil(real(x)), 0) }) }  // ⌈x ceil
func fmt(x v) v { return e("nyi") }                                                      // $x ⍕x format
func rng(x v) v { // ?x random uniform, ?-x normal ?z bi-normal
	xz, vec, _ := nv(x)
	if vec {
		return e("domain")
	}
	if imag(xz[0]) != 0 {
		r := make(zv, int(math.Floor(imag(xz[0])))) // complex(0,n)
		for i := range r {
			r[i] = complex(rand.NormFloat64(), rand.NormFloat64())
		}
		return r
	}
	f, norm := real(xz[0]), false
	if f < 0 {
		f, norm = -f, true
	}
	r := make(zv, int(math.Floor(f)))
	for i := range r {
		if norm {
			r[i] = complex(rand.NormFloat64(), 0)
		} else {
			r[i] = complex(rand.Float64(), 0)
		}
	}
	return r
}
func unq(x v) v { // ?x ∪x uniq
	w, t := ls(x)
	r := make(l, 0)
	for i := range w {
		if !some(r, func(x v) bool { return mch(w[i], x) == zi(1) }) {
			r = append(r, cp(w[i]))
		}
	}
	return sl(r, t)
}
func typ(x v) v { return e("nyi") } // @x type of
func evl(x v) v { return e("nyi") } // .x ⍎x evaluate

// dyadic verbs
func add(x, y v) v { return nd(x, y, func(x, y z) z { return x + y }) }                                  // x+y add
func sub(x, y v) v { return nd(x, y, func(x, y z) z { return x - y }) }                                  // x-y substract
func mul(x, y v) v { return nd(x, y, func(x, y z) z { return x * y }) }                                  // x*x x×y multiply
func div(x, y v) v { return nd(x, y, func(x, y z) z { return x / y }) }                                  // x%y x÷y divide
func mod(x, y v) v { return nd(x, y, func(x, y z) z { return complex(math.Mod(real(y), real(x)), 0) }) } // x!y modulo
func mkd(x, y v) v { // xl!yl make dictionary
	a, _ := ls(x)
	b, _ := ls(y)
	if len(a) != len(b) {
		return e("length")
	}
	return dict{k: a, v: b}.mp()
}
func min(x, y v) v { return nd(x, y, func(x, y z) z { return zter(real(x) < real(y), x, y) }) } // x&y x⌊y minimum
func max(x, y v) v { return nd(x, y, func(x, y z) z { return zter(real(x) > real(y), x, y) }) } // x|y x⌈y maximum
func les(x, y v) v {
	x, y = sn2(x, y)
	return nd(x, y, func(x, y z) z { return zter(real(x) < real(y), 1, 0) })
} // x<y less than
func mor(x, y v) v {
	x, y = sn2(x, y)
	return nd(x, y, func(x, y z) z { return zter(real(x) > real(y), 1, 0) })
}                  // x>y more than
func eql(x, y v) v { x, y = sn2(x, y); return nd(x, y, func(x, y z) z { return zter(x == y, 1, 0) }) } // x=y equal
func pow(x, y v) v { return nd(x, y, func(x, y z) z { return cmplx.Pow(x, y) }) }                      // x⍣y power
func lgn(x, y v) v { return nd(x, y, func(x, y z) z { return cmplx.Log(y) / cmplx.Log(x) }) }
func mch(x, y v) v { // x~y x≡y match
	if rtyp(x) != rtyp(y) {
		return zi(0)
	}
	if xd, o := md(x); o {
		if yd, o := md(y); o {
			return min(mch(xd.k, yd.k), mch(xd.v, yd.v))
		}
		return e("assert")
	}
	if reflect.DeepEqual(x, y) {
		return zi(1)
	}
	return zi(0)
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
	n := idx(x)
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
	a, b, c := idx(at(x, 0)), idx(at(x, nx-1)), 0
	xv := make([]int, nx)
	for i := range xv {
		xv[i] = idx(at(x, i))
	}
	var rshr func(x, y v, i int) v
	rshr = func(x, y v, i int) v {
		nx, ny := ln(x), ln(y)
		return krange(xv[i], func(z int) v {
			if i == nx-1 {
				c++
				return at(y, (c-1)%ny)
			}
			return rshr(x, y, i+1)
		})
	}
	if a < 0 {
		if ny == 0 {
			return y
		}
		n := ny / b
		if n*b < ny {
			n++
		}
		return cut(krange(n, func(z int) v { return z * b }), y)
	} else if b < 0 {
		return cut(krange(idx(a), func(z int) v { return z * ny / a }), y)
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
	j := idx(x)
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
		pa, pb := pidx(a), pidx(b)
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
	n := idx(x)
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
			r[i] = ll[int(math.Round(float64(ny)*rand.Float64()))]
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
		for i := range d.k {
			u := d.v[i]
			d.v[i] = zi(nx)
			for j := 0; j < nx; j++ {
				if mch(at(x, j), u) == zi(1) {
					d.v[i] = zi(j)
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
	r := make(zv, ny)
	for i := range r {
		for j := 0; j < nx; j++ {
			r[i] = zi(nx)
			if mch(at(x, j), at(y, i)) == zi(1) {
				r[i] = zi(j)
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
	} else if m, o := y.(s); o {
		var zero rV
		if f := rval(x).MethodByName(m); f != zero {
			return f.Interface()
		}
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
	// TODO: functions might need a. Test with t:=f.Type();t.IsVariadic() == false && t.NumIn()... and append a.
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
func ech(f, x v, a map[v]v) v { // f1'x  f1¨x each
	if n := ln(x); n < 0 {
		return atx(f, x, a)
	}
	xl, t := ls(x)
	r := make(l, len(xl))
	for i := range r {
		r[i] = cal(f, l{xl[i]}, a)
	}
	return sl(r, t)
}
func ecd(f, x, y v, a map[v]v) v { // x f2'y  x f2¨y each dyad
	nx, ny := ln(x), ln(y)
	if nx < 0 && ny < 0 {
		return cal(f, l{x, y}, a)
	}
	if nx < 0 {
		x, nx = rsh(ny, x), ny
	} else if ny < 0 {
		y, ny = rsh(nx, y), nx
	}
	if nx != ny {
		return e("size")
	}
	xl, xt := ls(x)
	yl, yt := ls(y)
	r := make(l, nx)
	for i := range r {
		r[i] = cal(f, l{xl[i], yl[i]}, a)
	}
	if xt == yt {
		return sl(r, xt)
	}
	return r
}
func ecp(f, x v, a map[v]v) v { // f2':x  f2⍨x each prior
	if xn := ln(x); xn < 1 {
		return x
	}
	xl, t := ls(x)
	r := make(l, len(xl))
	r[0] = cp(xl[0]) // always start with first. Some Ks have f dependend initials.
	for i := 1; i < len(r); i++ {
		r[i] = cal(f, l{xl[i], xl[i-1]}, a)
	}
	return sl(r, t)
}
func eci(f, x, y v, a map[v]v) v { // x f2':y  x f2⍨y each prior initial
	yl, t := ls(y)
	r := make(l, len(yl))
	r[0] = cal(f, l{yl[0], x}, a)
	for i := 1; i < len(r); i++ {
		r[i] = cal(f, l{yl[i], yl[i-1]}, a)
	}
	return sl(r, t)
}
func ecr(f, x, y v, a map[v]v) v { // x f/:y  x f⌿y each right
	yl, t := ls(y)
	r := make(l, len(yl))
	for i := range r {
		r[i] = cal(f, l{x, yl[i]}, a)
	}
	return sl(r, t)
}
func ecl(f, x, y v, a map[v]v) v { // x f\:y  x f⍀y each left
	xl, t := ls(x)
	r := make(l, len(xl))
	for i := range r {
		r[i] = cal(f, l{xl[i], y}, a)
	}
	return sl(r, t)
}
func fix(f, x v, a map[v]v) v { // f1/x fixed point
	x0 := cp(x)
	y := cp(x)
	z1 := zi(1)
	for {
		r := cal(f, l{y}, a)
		if mch(r, x0) == z1 || mch(r, y) == z1 {
			break
		}
		y = r
	}
	return y
}
func ovr(f, x v, a map[v]v) v { // f2/x
	nx := ln(x)
	if nx <= 0 { // no default values, but empty list, like k4
		return x
	}
	w, _ := ls(x)
	if nx == 1 {
		return w[0]
	}
	return ovd(f, w[0], w[1:], a)
}
func whl(f, x, y v, a map[v]v) v { // n f1/y for, g1 f1/y while
	if rval(x).Kind() == reflect.Func {
		for {
			if b := cal(x, l{y}, a); idx(b) != 1 {
				return y
			}
			y = cal(f, l{y}, a)
		}
	}
	n := pidx(x)
	for i := 0; i < n; i++ {
		y = cal(f, l{y}, a)
	}
	return y
}
func ovd(f, x, y v, a map[v]v) v { // x f2/y over initial
	w, _ := ls(y)
	for _, u := range w {
		x = cal(f, l{x, u}, a)
	}
	return x
}
func sfx(f, x v, a map[v]v) v { return e("nyi") } // f1\x scan fixed
func scn(f, x v, a map[v]v) v { // f2\x scan
	w, t := ls(x)
	r := make(l, len(w))
	r[0] = w[0]
	for i, u := range w[1:] {
		r[i+1] = cal(f, l{r[i], u}, a)
	}
	return sl(r, t)
}
func swl(f, x, y v, a map[v]v) v { return e("nyi") } // x f1\y scan for, g1 f1\y scan while
func sci(f, x, y v, a map[v]v) v { // x f2\x scan initial
	w, t := ls(y)
	nx := ln(x)
	r := make(l, len(w))
	for i, u := range w {
		x = cal(f, l{x, u}, a)
		r[i] = cp(x)
	}
	if nx < 0 {
		return sl(r, t)
	}
	return r
}

func grade(up bool, x v) v {
	if d, o := md(x); o {
		return atx(d.k, grade(up, d.v), nil)
	}
	if ln(x) < 0 {
		x = enl(x)
	}
	x = cp(x)

	switch t := x.(type) {
	case zv:
		f := make([]float64, len(t))
		for i := range t {
			f[i] = real(t[i])
		}
		x = sort.Float64Slice(f)
	case sv:
		x = sort.StringSlice(t)
	}
	if d, o := x.(sort.Interface); o {
		if !up {
			d = sort.Reverse(d)
		}
		i := make(zv, d.Len())
		for n := range i {
			i[n] = zi(n)
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
	zv, _, _ := nv(x)
	return grade(up, zv)
}

type grades struct {
	sort.Interface
	idx zv
}

func (s grades) Swap(i, j int) {
	s.Interface.Swap(i, j)
	s.idx[i], s.idx[j] = s.idx[j], s.idx[i]
}
