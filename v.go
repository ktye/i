package i

import (
	"math"
	"math/cmplx"
	"math/rand"
	"reflect"
	"sort"
	"strconv"
)

//    +m  ⍉ flip, transpose    /  ⍉3 2⍴⍳6           → (0 2 4;1 3 5)
func flp(x v) v {
	n, m := ln(x), -1
	if n < 0 {
		return x
	}
	a, _ := ls(x)
	var r []l
	for i := range a {
		if na := ln(a[i]); na < 0 {
			return x
		} else if i == 0 {
			m = na
			r = make([]l, m)
			for k := range r {
				r[k] = make(l, len(a))
			}
		} else if na != m {
			return e("length")
		}
		b, _ := ls(a[i])
		for k := range b {
			r[k][i] = cp(b[k])
		}
	}
	ul := make(l, len(r))
	for i := range r {
		ul[i] = uf(r[i])
	}
	return ul
}

//    -x  negate               / -(1;2;(3;(4;5)))   → (-1;-2;(-3;(-4 -5)))
func neg(x v) v { return nm(x, func(x z) z { return -x }) }

//    *l  first                / *2 3 4             → 2
//    *d  first                / *[a:1;b:2]         → 1
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

// con[x] conjugate complex    / con 1i3            → 1i-3
func con(x v) v { return nm(x, func(x z) z { return cmplx.Conj(x) }) }

//    √x  sqrt[x] square root  / √- 4 9             → 0i-2 0i-3
func sqr(x v) v { return nm(x, func(x z) z { return cmplx.Sqrt(x) }) }

//    %x  ÷ inverse            / %2 5               → 0.5 0.2
func inv(x v) v { return nm(x, func(x z) z { return 1 / x }) }

//    ‖x  abs[x] magnitue      / ‖3.2a30            → 3.2
func abs(x v) v { return nm(x, func(x z) z { return complex(cmplx.Abs(x), 0) }) }

//    𝜑x  rad[x] phase         / 𝜑-1                → π
func rad(x v) v { return nm(x, func(x z) z { return complex(cmplx.Phase(x), 0) }) }

//    °x  deg[x] angle         / °1i1               → 45
func deg(x v) v {
	return nm(x, func(x z) z {
		p := cmplx.Phase(x) / math.Pi * 180.0
		if p < 0 {
			p += 360.0
		}
		return complex(p, 0)
	})
}

//    ℜx  re[x] real part      / ℜ1a90              → 0
func zre(x v) v { return nm(x, func(x z) z { return complex(real(x), 0) }) }

//    ℑx  im[x] imag part      / ℑ1a90              → 1
func zim(x v) v { return nm(x, func(x z) z { return complex(imag(x), 0) }) }

//    !n  ⍳ til, iota          / ⍳3                 → 0 1 2
//    !l  ⍳ keys               / !5 3 1             → 0 1 2
//    !d  ⍳ keys               / ![a:1;b:2]         → (`a;`b)
func til(x v) v { // !x ⍳x iota
	if d, ok := md(x); ok {
		return d.k
	}
	if n := ln(x); n >= 0 {
		return til(zi(n))
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

// odo[l] odometer             / odo[3 2]           → (0 0 1 1 2 2;0 1 0 1 0 1)
func odo(x v) v {
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

//    &l  where, repeat        / &3 0 4             → 0 0 0 2 2 2 2
func wer(x v) v { // &x where
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

//    |x  ⌽ reverse            / ⌽1 2 3             → 3 2 1
//    |d  ⌽ reverse            / ⌽[a:1;b:2]         → [b:2;a:1]
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

//    <x  ⍋ grade up           / ⍋8 2 9 1           → 3 1 0 2
func asc(x v) v { return grade(true, x) }

//    >x  ⍒ grade down         / ⍒`alpha`beta       → 1 0
func dsc(x v) v { return grade(false, x) }

//    =n  unit matrix          / =3                 → (1 0 0;0 1 0;0 0 1)
func eye(x v) v { // =x unit matrix
	n := pidx(x)
	l := make(l, n)
	for i := range l {
		r := make(zv, n)
		r[i], l[i] = 1, r
	}
	return l
}

//    =l  group                / =(3;"a";5;3;"a";3) → (3;"a";5)!(0 3 5;1 4;,2)
func grp(x v) v {
	n := ln(x)
	if n < 0 {
		return eye(x)
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

//    ~x  not                  / ~3 2 -1 0          → 0 0 0 1
func not(x v) v { return nm(x, func(x z) z { return zter(x == 0, 1, 0) }) }

//    ,x  enlist               / ,1                 → ,1
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

//    ^x  isnil isnan          / ^(0%0;0;ø)         → 1 0 1
func is0(x v) v {
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

//    ⍣x  exp[x] exponential   / 𝜀>‖1+⍣0i1*π        → 1
func exp(x v) v { return nm(x, func(x z) z { return cmplx.Exp(x) }) }

//    ⍟x  log[x] logarithm     / ⍟⍣1                → 1
func log(x v) v { return nm(x, func(x z) z { return cmplx.Log(x) }) }

//    #x  ⍴ count, length      / ⍴2 3⍴⍳6            → 2
func cnt(x v) v { // #x ⍴x count, length
	if d, o := md(x); o {
		return zi(len(d.k))
	} else if n := ln(x); n >= 0 {
		return zi(n)
	}
	return zi(1)
}

//    _x  ⌊ floor              / ⌊1.23              → 1
func flr(x v) v { return nm(x, func(x z) z { return complex(math.Floor(real(x)), 0) }) }

//    ⌈x  ceil                 / ⌈1.23              → 2
func cil(x v) v { return nm(x, func(x z) z { return complex(math.Ceil(real(x)), 0) }) }

//    $x  format, tostring     / $(1;2;3)           → "(1;2;3)"
func fmt(x v) v { return cst(nil, x) } // $x format to string
func num(x v) v { // num s parse number // TODO: move to 0$"1.23" (needed in ⍳/io.go)
	t, o := x.(s)
	if !o {
		return e("type")
	}
	return (&p{}).num(t)
}

//    ?n  random uniform       / +/1>?1000          → 1000
//    ?-n random normal        / 900 > +/1>?-1000   → 1
//    ?i  random binormal      / (+/‖?0i1000)<1300  → 1 / √π÷2
func rng(x v) v {
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

//    ?x  unique               / ?2 3 3 0 4 0       → 2 3 0 4
//    ?s  parse                / ?"1+2"             → ("+";1;2)
func unq(x v) v {
	n := ln(x)
	if n < 0 {
		if s, o := x.(s); o {
			return prs(s)
		}
		return rng(x)
	}
	w, t := ls(x)
	r := make(l, 0)
	nw := func(u v) bool {
		for i := range r {
			if mch(u, r[i]) == zi(1) {
				return false
			}
		}
		return true
	}
	for i := range w {
		if nw(w[i]) {
			r = append(r, cp(w[i]))
		}
	}
	return sl(r, t)
}

//    @x  type of              / (@5)≡@-1.2i5       → 1
func typ(x v) v { return rtyp(x) }

//    .s  parse and eval       / ."1+2"             → 3
//    .l  evaluate             / .(+;1;2)           → 3
func evl(x v, a map[v]v) v {
	var b l
	if s, o := x.(s); o {
		b = prs(s).(l)
	} else if l, o := x.(l); o {
		b = l
	} else {
		return e("type")
	}
	return eva(b, a)
}

//   x+y  add                  / 1+2 3 4            → 3 4 5
func add(x, y v) v { return nd(x, y, func(x, y z) z { return x + y }) }

//   x-y  substract            / 2 3 4-1            → 1 2 3
func sub(x, y v) v { return nd(x, y, func(x, y z) z { return x - y }) }

//   x*y  × multiply           / 2 0i1 1a270*0i1    → 0i2 -1 1
func mul(x, y v) v { return nd(x, y, func(x, y z) z { return x * y }) }

//   x%x  ÷ divide             / 1÷2                → 0.5
func div(x, y v) v { return nd(x, y, func(x, y z) z { return x / y }) }

//mod[x;y]  modulo             / mod[2;5]           → 1
func mod(x, y v) v { return nd(x, y, func(x, y z) z { return complex(math.Mod(real(y), real(x)), 0) }) } // x!y modulo

//   x!y  make dictionary      / `a`b`c!(10;2 3;`f) → [a:10;b:2 3;c:`f]
func mkd(x, y v) v {
	nx, ny := ln(x), ln(y)
	if nx < 0 && ny < 0 {
		x, y, nx, ny = l{x}, l{y}, 1, 1
	} else if nx < 0 {
		return dict{k: l{x}, v: l{y}}.mp()
	} else if ny < 0 {
		yl := make(l, nx)
		for i := range yl {
			yl[i] = cp(y)
		}
		y = yl
	}
	a, _ := ls(x)
	b, _ := ls(y)
	if len(a) != len(b) {
		return e("length")
	}
	return dict{k: a, v: b}.mp()
}

//   x&y  ⌊ min                / 1 2 3 4&4 3 2 1    → 1 2 2 1
func min(x, y v) v { return nd(x, y, func(x, y z) z { return zter(real(x) < real(y), x, y) }) }

//   x|y  ⌈ max                / 1 2 3 4|4 3 2 1    → 4 3 3 4
func max(x, y v) v { return nd(x, y, func(x, y z) z { return zter(real(x) > real(y), x, y) }) }

//   x<y  less than            / 5<8 1 5            → 1 0 0
func les(x, y v) v {
	x, y = sn2(x, y)
	return nd(x, y, func(x, y z) z { return zter(real(x) < real(y), 1, 0) })
}

//   x>y  more than            / 5>8 1 5            → 0 1 0
func mor(x, y v) v {
	x, y = sn2(x, y)
	return nd(x, y, func(x, y z) z { return zter(real(x) > real(y), 1, 0) })
}

//   x=y  equals               / 1 ø ∞=(1a0;0%0;1%0)→ 1 1 1
func eql(x, y v) v {
	x, y = sn2(x, y)
	return nd(x, y, func(x, y z) z {
		// all these are Inf not NaN: c(∞,ø) c(-∞,ø) c(ø, ∞) c(ø, -∞)
		if (cmplx.IsNaN(x) && cmplx.IsNaN(y)) || (cmplx.IsInf(x) && cmplx.IsInf(y)) {
			return 1
		}
		return zter(x == y, 1, 0)
	})
}

//   x‖y  rct parts to complex / 2‖!4               → 2 2i1 2i2 2i3
func rct(x, y v) v { return nd(x, y, func(x, y z) z { return complex(real(x), real(y)) }) }

//   x°y  pol polar to complex / 1 2 3°0 90 180     → 1 0i2 -3
func pol(x, y v) v { // x°y complex from abs and deg
	return nd(x, y, func(x, y z) z {
		r := cmplx.Rect(real(x), real(y)*math.Pi/180.0)
		if y == 0 || y == 180 {
			r = complex(real(r), 0)
		} else if y == 90 || y == 270 {
			r = complex(0, imag(r))
		}
		return r
	})
}

//   x𝜑y  prd polar to complex / 1𝜑0 π -π           → 1 -1 -1
func prd(x, y v) v {
	return nd(x, y, func(x, y z) z {
		switch real(y) {
		case 0:
			return complex(real(x), 0)
		case math.Pi, -math.Pi:
			return complex(-real(x), 0)
		}
		return cmplx.Rect(real(x), real(y))
	})
}

//   x⍣y  power                / 2⍣3                → 8
func pow(x, y v) v { return nd(x, y, func(x, y z) z { return cmplx.Pow(x, y) }) }

//   x√y  nrt[x;y] nth root    / 𝜀>‖2-3√8           → 1
func nrt(x, y v) v { return nd(x, y, func(x, y z) z { return cmplx.Pow(y, 1/x) }) }

//   x⍟y  lgn[x;y] base n log  / 𝜀>‖3-10⍟1000       → 1
func lgn(x, y v) v { return nd(x, y, func(x, y z) z { return cmplx.Log(y) / cmplx.Log(x) }) }

//   x~y  ≡ match, deep equal  / 1a90 2.0 3≡0i1 2 3 → 1
func mch(x, y v) v {
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

//   x,y  catentate            / (1;2),3            → (1;2;3)
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

//   x^y  except               / (!10)^!7           → 7 8 9
func ept(x, y v) v {
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

//   a#d  ↑ take               / 2↑[a:1;b:2;c:3]    → [a:1;b:2]
//   a#l  ↑ take               / -2↑⍳10             → 8 9
func tak(x, y v) v {
	// nyi: 5,8,9: function, verb, adverb
	if ln(x) >= 0 {
		return rsh(x, y)
	}
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

//   l#d  ⍴ select             / `a`c#[a:1;b:2;c:3] → [a:1;c:3]
//   l#y  ⍴ reshape            / 2 3⍴⍳6             → (0 1 2;3 4 5)
func rsh(x, y v) v { // x#y x⍴y reshape
	if yd, o := md(y); o { // select from dict
		yd.v, _ = ls(atx(y, x, nil))
		yd.k, _ = ls(x)
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
	rng := func(n int, f func(int) v) v {
		l := make(l, n)
		for i := range l {
			l[i] = f(i)
		}
		return uf(l)
	}
	var rshr func(x, y v, i int) v
	rshr = func(x, y v, i int) v {
		nx, ny := ln(x), ln(y)
		return rng(xv[i], func(z int) v {
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
		return cut(rng(n, func(z int) v { return z * b }), y)
	} else if b < 0 {
		return cut(rng(idx(a), func(z int) v { return z * ny / a }), y)
	}
	return rshr(x, y, 0)
}

//   x_d  ↓ delete             / `a`b_[a:1;b:2;c:3] → [c:3]
//   x_y  ↓ drop               / (1_1 2;-1_!3;5_,1) → (,2;0 1;0#,0)
func drp(x, y v) v {
	if d, o := md(y); o {
		nx := ln(x)
		if nx < 0 {
			d.k, d.v = drp(x, d.k).(l), drp(x, d.v).(l)
		} else {
			d.k = ept(d.k, x).(l)
			d.v, _ = ls(atx(y, d.k, nil))
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

//  xl_yl cut                  / 3 5_!8             → (3 4;5 6 7)
func cut(x, y v) v {
	if ln(x) < 0 || ln(y) < 0 {
		return drp(x, y)
	}
	xl, _ := ls(x)
	r := make(l, len(xl))
	ll, _ := ls(cp(y))
	for i := range xl {
		a := pidx(xl[i])
		b := len(ll)
		if i < len(r)-1 {
			b = pidx(xl[i+1])
		}
		if b < a {
			return e("domain")
		}
		r[i] = uf(ll[a:b])
	}
	return r
}

//   x$y  convert to typeof    / (int@8)$128        → (int@8)$-128
//   d$y  format               / [t:1]$[a:,1;b:,3]  → "a b\n---\n1 3"
func cst(x, y v) v { // x$y cast
	type cvt interface {
		ConvertTo(v) v
	}
	type cfm interface {
		Format(s) s
	}
	var m dict
	if x == "c" {
		x = byte(0)
		if s, o := y.(s); o {
			r := []byte(s)
			if len(r) == 1 {
				return r[0]
			}
			return r
		}
	}
	if x == nil { // nil argument: tostring, same as fmt.
		m, _ = md(map[v]v{0: 0})
	} else if c, o := x.(cvt); o { // use ConvertTo method of x
		return c.ConvertTo(y)
	} else if f, o := y.(cfm); o && iss(x) { // custom Format method
		return f.Format(x.(s))
	} else if d, o := md(x); o { // dict controls formatting
		m = d
	} else if rval(x).Kind() < reflect.Array { // convert to any numeric type
		z, vec, _ := nv(y) // (int$-8)$255
		return vn(z, vec, rtyp(x))
	}
	getf := func(a s) int {
		if k, u := m.at(a); k == -1 {
			return -1
		} else {
			return pidx(u)
		}
	}
	mat := func(x v) int {
		if n := ln(x); n <= 0 {
			return 0
		} else {
			m := 0
			for i := 0; i < n; i++ {
				mm := ln(at(x, i))
				if i == 0 {
					if mm <= 0 {
						return 0
					}
					m = mm
				} else if mm != m {
					return 0
				}
			}
			return n
		}
	}
	compact := func() dict {
		mm := cp(m).(dict)
		mm.set("l", 0)
		mm.set("t", 0)
		mm.set("d", 0)
		return mm
	}
	tab := func(hdr, cols l, n int) s {
		mm := compact()
		m := make([]sv, n+2)
		for i := range m {
			m[i] = make(sv, len(cols))
		}
		mm.set("q", 1)
		for k := range hdr {
			m[0][k] = ""
			m[1][k] = ""
			if hdr != nil {
				m[0][k] = cst(mm, hdr[k]).(s)
			}
		}
		mm.set("q", 0)
		for k := range cols {
			mx := len(m[0][0])
			c, _ := ls(cols[k])
			for i := 0; i < n; i++ {
				t := cst(mm, c[i]).(s)
				m[2+i][k] = t
				if n := len(rv(t)); n > mx {
					mx = n
				}
			}
			for i := range m {
				f := rv(m[i][k])
				w := make(rv, mx)
				for j := range w {
					if i == 1 {
						w[j] = '-'
					} else if j < len(f) {
						w[j] = f[j]
					} else {
						w[j] = ' '
					}
				}
				m[i][k] = string(w)
			}
		}
		if hdr == nil {
			m = m[2:]
		}
		var b []byte
		for i := range m {
			for k := range m[i] {
				b = append(b, []byte(m[i][k])...)
				if k < len(m[i])-1 {
					c := byte(' ')
					if i == 1 && hdr != nil {
						c = '-'
					}
					b = append(b, c)
				}
			}
			if i < len(m)-1 {
				b = append(b, '\n')
			}
		}
		return string(b)
	}
	dct := func(d dict) string {
		mm := compact()
		hb := make(sv, len(d.k))
		var b []byte
		mx := 0
		for i := range d.k {
			t := cst(mm, d.k[i]).(s)
			if n := len(rv(t)); n > mx {
				mx = n
			}
			hb[i] = t
		}
		for i := range d.k {
			t := rv(hb[i])
			b = append(b, []byte(hb[i])...)
			for k := 0; k < mx-len(t); k++ {
				b = append(b, ' ')
			}
			b = append(b, '|', ' ')
			b = append(b, []byte(cst(mm, d.v[i]).(s))...)
			if i < len(d.k)-1 {
				b = append(b, '\n')
			}
		}
		return string(b)
	}
	if y == nil {
		return "::"
	}
	type stringer interface {
		String() string
	}
	if s, o := y.(stringer); o {
		return s.String()
	}
	if b, o := y.(byte); o {
		return cst(nil, []byte{b})
	}
	n := ln(y)
	if n >= 0 {
		if b, o := y.([]byte); o {
			s := make([]byte, 2+2*len(b))
			s[0] = '0'
			s[1] = 'x'
			t := "0123456789abcdef"
			for i := range b {
				s[2+2*i] = byte(t[b[i]>>4])
				s[3+2*i] = byte(t[b[i]&15])
			}
			return string(s)
		}
		if getf("m") > 0 {
			if n := mat(y); n > 0 {
				cols := flp(y).(l)
				return tab(nil, cols, n)
			}
		}
		r, t := ls(y)
		vs := make(sv, len(r))
		for i := range vs {
			st := cst(m, r[i]).(s)
			if getf("l") > 0 {
				st = jon("\n ", spl("\n", st)).(s)
			}
			vs[i] = st
		}
		if len(vs) == 1 {
			return "," + vs[0]
		}
		if t == nil {
			if getf("l") > 0 && len(vs) > 0 {
				vs[0] = "(" + vs[0]
				vs[len(vs)-1] += ")"
				return jon("\n ", vs)
			}
			return "(" + jon(";", vs).(s) + ")"
		}
		return jon(" ", vs).(s)
	} else if d, o := md(y); o {
		if getf("t") > 0 {
			n := 0
			for i := range d.v {
				m := ln(d.v[i])
				if m < 0 {
					n = 0
					break
				} else if i == 0 {
					n = m
				} else if m != n {
					break
				}
			}
			if n > 0 {
				return tab(d.k, d.v, n)
			}
		}
		if getf("d") > 0 {
			return dct(d)
		}
		return "(" + cst(m, d.k).(s) + "!" + cst(m, d.v).(s) + ")"
	} else if y, o := y.(z); o {
		if cmplx.IsNaN(y) {
			return "ø"
		} else if cmplx.IsInf(y) {
			return "∞"
		}
		p := getf("p")
		re := strconv.FormatFloat(real(y), 'g', p, 64)
		if imag(y) == 0 {
			return re
		}
		if d := getf("a"); d < 0 {
			im := strconv.FormatFloat(imag(y), 'g', -1, 64)
			return re + "i" + im
		} else {
			r, phi := cmplx.Polar(y)
			phi *= 180.0 / math.Pi
			if phi < 0 {
				phi += 360
			}
			return strconv.FormatFloat(r, 'g', p, 64) + "a" + strconv.FormatFloat(phi, 'f', d, 64)
		}
	}
	u := rval(y)
	if k := u.Kind(); k < reflect.Array && k != reflect.Uintptr {
		r, _, _ := nv(y)
		return cst(m, r[0]).(s)
	} else if k == reflect.String {
		u := u.String()
		t := strconv.Quote(u)
		if getf("q") == 1 {
			if u == t[1:len(t)-1] {
				return t[1 : len(t)-1]
			}
		}
		return t
	}
	return "(?" + u.Type().String() + ")"
}

//   n?m  roll                 / 100>#?100?100      → 1
//  -n?m  deal                 / #?-10?10           → 10
//   n?l  random select        / 100>#?100?!100     → 1
//  -n?l  deal shuffled        / #?-100?!100        → 100
func rnd(x, y v) v {
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
			idx := int(math.Round(float64(ny-1) * rand.Float64()))
			r[i] = ll[idx]
		}
	}
	return sl(r, rT)
}

//  xl?y  find                 / 3 5?⍳7             → 2 2 2 0 2 1 2
func fnd(x, y v) v {
	nx := ln(x)
	if nx < 0 || isd(x) {
		return rnd(x, y)
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

//   l@y  at, index            / 2 5 6@0 2          → 2 6
//   d@y  at, index            / [a:1;b:2;c:3]@`a`c → 1 3
//   f@y  monadic call         / {-x}@2 3           → -2 -3
func atx(x, y v, a map[v]v) v {
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
		ll, _ := ls(y)
		r := make(l, len(ll))
		for i := range r {
			r[i] = atx(x, ll[i], nil)
		}
		return uf(r)
		// return kmap(y, func(z v, i int) v { return atx(x, z, nil) })
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

//   l.y  depth list index     / (1;(2;(3;4))).1 1 0→ 3
//   d.y  depth dict index     / [a:1;b:[c:2]].`b`c → 2
//   x.m  method               / (myf`F0)[]+0       → 2
//   x.y  call                 / {x+y}.(3;4 5)      → 7 8
//   x.y  curry                / ({x+y+z}.(1;;3)) 2 → 6
func cal(x, y v, a map[v]v) v {
	if x == nil {
		return e("call nil")
	}
	if y != nil {
		if yl, o := y.(l); o && len(yl) == 1 {
			if m, o := yl[0].(s); o {
				var zero rV
				if f := rval(x).MethodByName(m); f != zero {
					return f.Interface()
				}
			}
		}
	}
	if sx, o := x.(s); o {
		f := lupr(a, sx)
		if f == nil {
			return e("nil:" + sx)
		}
		return cal(f, y, a)
	}
	f := rval(x)
	if f.Kind() != reflect.Func {
		if _, o := md(x); o || ln(x) > 0 {
			if yn := ln(y); yn > 0 {
				return atd(x, y, a)
			}
			return atx(x, y, a)
		}
		return e("type:" + f.Kind().String())
	}
	var in, r []rV
	var cur []int
	yl, _ := ls(y)
	in = make([]rV, len(yl))
	for i := range in {
		if yl[i] == nil {
			cur = append(cur, i)
		} else {
			in[i] = rval(yl[i])
		}
	}
	if cur == nil {
		// TODO: functions might need a. Test with t:=f.Type();t.IsVariadic() == false && t.NumIn()... and append a.
		r = f.Call(in)
		if len(r) == 0 {
			return nil
		}
		return r[0].Interface()
	}
	return curry(func(w ...v) v { // curry
		if len(w) == 0 { // report number of arguments
			return len(cur)
		}
		if len(w) != len(cur) {
			return e("args")
		}
		yl := y.(l)
		for i := range w {
			if w[i] != nil {
				yl[cur[i]] = w[i]
			}
		}
		return cal(x, yl, a)
	})
}
func atd(x, y v, a map[v]v) v { // at depth
	n := ln(y)
	if n < 0 {
		return atx(x, y, a)
	} else if n == 1 {
		return atx(x, at(y, 0), a)
	}
	var ata func(v) v
	var all func() v
	if d, o := md(x); o {
		ata = func(x v) v { _, r := d.at(x); return r }
		all = func() v { return d.k }
	} else {
		ata = func(z v) v { return at(x, pidx(z)) }
		all = func() v { return til(cnt(x)) }
	}
	y0 := at(y, 0)
	if y0 == nil {
		y0 = all()
	}
	yp := drp(1, y)
	m := ln(y0)
	if m < 0 {
		return atd(ata(y0), yp, a)
	}
	r := make(l, m)
	for i := range r {
		yi := at(y0, i)
		r[i] = atd(ata(yi), yp, a)
	}
	return uf(r)
}

//   s/y  join                 / ";"/`alpha`beta    → "alpha;beta"
func jon(x, y v) v { // a/l join
	xs, xo := x.(s)
	yy, yo := y.(sv)
	if !xo || !yo {
		return e("type") // TODO custom string types
	}
	if len(yy) == 0 {
		return ""
	}
	n := 0
	for i := range yy {
		n += len(yy[i])
	}
	n += len(xs)*len(yy) - 1
	r := make([]byte, 0, n)
	sep := []byte(xs)
	for i := range yy {
		r = append(r, []byte(yy[i])...)
		if i != len(yy)-1 {
			r = append(r, sep...)
		}
	}
	return string(r)
}

/*   TODO encode */
func enc(x, y v) v { return e("nyi") } // l/a encode, pack

//   s\y  split                / ";"\"a;b;;c;d"     → `a`b``c`d
func spl(x, y v) v { // a\x split, decode?
	eq := func(a, b []rune) bool {
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}
	xs, xo := x.(s)
	ys, yo := y.(s)
	if !xo || !yo {
		return e("type")
	}
	xr, yr := []rune(xs), []rune(ys)
	var r sv
	if len(xs) == 0 {
		if len(ys) == 0 {
			return sv{""}
		}
		r = make(sv, len(yr))
		for i := range yr {
			r[i] = string(yr[i])
		}
		return r
	}
	l := 0
	for i := 0; i < len(yr); i++ {
		if len(yr)-i < len(xr) {
			break
		}
		if eq(yr[i:i+len(xr)], xr) {
			r = append(r, string(yr[l:i]))
			i += len(xr) - 1
			l = i + 1
		}
	}
	return append(r, string(yr[l:]))
}

/*   TODO decode */
func dec(x, y v) v { return e("nyi") }

//   f'x  ¨ each               / -:¨1 2             → -1 -2
func ech(f, x v, a map[v]v) v {
	if d, o := md(x); o {
		r, _ := ls(ech(f, d.v, a))
		if len(r) != len(d.k) {
			return e("length")
		}
		d.v = r
		return d.mp()
	}
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

// x g'y  ¨ each pair          / 2 3*¨4 5           → 8 15
func ecd(f, x, y v, a map[v]v) v {
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

//   g':x ⍨ each prior         / -⍨1 5 3            → 1 4 -2
func ecp(f, x v, a map[v]v) v {
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

// x g':y ⍨ each prior initial / 7-⍨1 5 3           → -6 4 -2
func eci(f, x, y v, a map[v]v) v {
	yl, t := ls(y)
	r := make(l, len(yl))
	r[0] = cal(f, l{yl[0], x}, a)
	for i := 1; i < len(r); i++ {
		r[i] = cal(f, l{yl[i], yl[i-1]}, a)
	}
	return sl(r, t)
}

// x g/:y ⌿ each right         / 1 2+⌿3 4 5         → (4 5;5 6;6 7)
func ecr(f, x, y v, a map[v]v) v {
	yl, t := ls(y)
	r := make(l, len(yl))
	for i := range r {
		r[i] = cal(f, l{x, yl[i]}, a)
	}
	return sl(r, t)
}

// x g\:  ⍀ each left          / 1 2+⍀3 4 5         → (4 5 6;5 6 7)
func ecl(f, x, y v, a map[v]v) v {
	xl, t := ls(x)
	r := make(l, len(xl))
	for i := range r {
		r[i] = cal(f, l{xl[i], y}, a)
	}
	return sl(r, t)
}

//   g/y  over, reduce         / +/1 2 3            → 6
func ovr(f, x v, a map[v]v) v {
	nx := ln(x)
	if nx <= 0 { // no default values, but empty list, like k4
		return x
	}
	w, _ := ls(x)
	if nx == 1 {
		return w[0]
	}
	return ovi(f, w[0], w[1:], a)
}

//   g\y  scan                 / +\1 2 3            → 1 3 6
func scn(f, x v, a map[v]v) v { // f2\x scan
	w, t := ls(x)
	r := make(l, len(w))
	r[0] = w[0]
	for i, u := range w[1:] {
		r[i+1] = cal(f, l{r[i], u}, a)
	}
	return sl(r, t)
}

// x g/y  over initial         / 5+/1 2 3           → 11
func ovi(f, x, y v, a map[v]v) v { // x f2/y over initial
	w, _ := ls(y)
	for _, u := range w {
		x = cal(f, l{x, u}, a)
	}
	return x
}

// x g\y  scan initial         / 5+\1 2 3           → 6 8 11
func sci(f, x, y v, a map[v]v) v {
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

// n f/y  for, repeat          / 3 (2⍣)/2           → 65536
// t f/y  while                / {x<100}{x*2}/1     → 128
func whl(f, x, y v, a map[v]v) v {
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

// x f\y  scan for             / 2 √:\81            → 81 9 3
// t f\y  scan while           / {x<100}(2*)\1      → 1 2 4 8 16 32 64 128
func swl(f, x, y v, a map[v]v) v { // x f1\y scan for, g1 f1\y scan while
	r := l{cp(y)}
	if rval(x).Kind() == reflect.Func {
		for {
			if b := cal(x, l{y}, a); idx(b) != 1 {
				return uf(r)
			}
			y = cal(f, l{y}, a)
			r = append(r, cp(y))
		}
	}
	n := pidx(x)
	for i := 0; i < n; i++ {
		y = cal(f, l{y}, a)
		r = append(r, cp(y))
	}
	return uf(r)
}

//   f/y  fixed point          / √:/2               → 1
func fix(f, x v, a map[v]v) v {
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

//   f\y  scan fixed           / -:\3               → 3 -3
func sfx(f, x v, a map[v]v) v { // f1\x scan fixed
	x0 := cp(x)
	y := cp(x)
	z1 := zi(1)
	u := l{x0}
	for {
		r := cal(f, l{y}, a)
		if mch(r, x0) == z1 || mch(r, y) == z1 {
			break
		}
		u = append(u, r)
		y = r
	}
	return uf(u)
}

//  πø∞𝜀  numeric constants    / (π;ø;∞;𝜀)=π ø ∞ 𝜀  → 1 1 1 1
//    $[x;z;y;…] if, switch    / $[1>2;∞;𝜀]         → 1e-14
//   x∇y  tail call            / {$[x>100;x;∇x+1]}1 → 101

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
		x = uf(l)
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
