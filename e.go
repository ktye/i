package i

func eva(x v, a map[v]v) v {
	if sy, o := x.(s); o {
		return lupr(a, sy)
	}
	l, o := cp(x).(l)
	if !o || len(l) == 0 {
		return x
	}
	if len(l) > 3 && l[0] == "$" {
		return cnd(l[1:], a) // $[A;B;C;…]
	}
	if p, o := l[0].(s); o && len(l) == 3 && len(p) > 1 && p != "::" && p[len(p)-1] == ':' {
		l = []v{":", p[:len(p)-1], l[1], l[2]}               // mod assignment
		if p := l[1].(s); len(p) > 1 && p[len(p)-1] == ':' { // ::
			l[0], l[1] = "::", p[:len(p)-1]
		}
	}
	switch l[0] {
	case nil:
		for i := len(l) - 1; i > 0; i-- { // right to left
			l[i] = eva(l[i], a)
		}
		return l[1:]
	case "`":
		return l[1]
	case ";":
		l = l[1:]
		var r interface{}
		for i := range l { // left to right
			r = eva(l[i], a)
		}
		return r
	case ":", "::":
		return asn(l[1:], l[0] == "::", a)
	case "∇":
		return tail(l[1:])
	case "λ":
		return λ(l[1:], a)
	default:
		if f, o := adv(l, a); o {
			return f
		}
		for i := len(l) - 1; i >= 0; i-- { // right to left
			l[i] = eva(l[i], a)
		}
		if len(l) == 1 {
			return l[0]
		}
		if _, o := l[1].(nilad); o {
			return cal(l[0], []v{}, a)
		}
		return cal(l[0], l[1:], a)
	}
}
func adv(u l, a map[v]v) (v, bool) { // evaluate adverb expr
	// TODO: verb trains
	if len(u) != 2 || u[0] == nil || !iss(u[0]) {
		return u, false
	}
	s := u[0].(s)
	if sAdv(rv(s)) == 0 {
		return u, false
	}
	af := lupr(a, s).(func(v) v)
	w := eva(u[1], a)
	return af(w), true // func(...v)v
}
func lup(a map[v]v, s s) (v, map[v]v) { // lookup
	if r := a[s]; r != nil {
		return r, a
	}
	if p, o := a[".."]; o {
		pp, o := p.(*map[v]v)
		if !o {
			return e("type"), nil
		}
		return lupr(*pp, s), *pp
	}
	return e("undefined:" + s), nil
}
func lupr(a map[v]v, s s) v { r, _ := lup(a, s); return r }
func ktr(a map[v]v) map[v]v { // k-tree root
	for {
		p, o := a[".."]
		if !o {
			break
		}
		r := p.(*map[v]v)
		a = *r
	}
	return a
}
func λ(b l, a map[v]v) v {
	var λn func(w ...v) v
	λn = func(w ...v) v {
		env := map[v]v{"..": &a} // lexical scoping
		env["o"] = λn
	tail:
		for i, r := range "xyz" {
			if len(w) > i {
				env[string(r)] = cp(w[i])
			}
		}
		for i, ex := range b {
			r := eva(ex, env)
			if i == len(b)-1 {
				if t, o := r.(tail); o {
					w = t
					for i := len(w) - 1; i >= 0; i-- {
						w[i] = eva(w[i], env)
					}
					goto tail
				}
				return r
			}
		}
		return nil
	}
	m := make(map[rune]bool)
	var cnt func(v)
	cnt = func(b v) {
		switch t := b.(type) {
		case string:
			if r := rv(t); len(r) == 1 && any(r[0], "xyz") {
				m[r[0]] = true
			}
			return
		case l:
			if len(t) > 0 {
				if s, o := t[0].(s); o && s == "λ" {
					return
				}
			}
			for i := range t {
				cnt(t[i])
			}
		}
	}
	cnt(b)
	if m['z'] {
		return func(x, y, z v) v { return λn(x, y, z) }
	} else if m['y'] {
		return func(x, y v) v { return λn(x, y) }
	} else if m['x'] {
		return func(x v) v { return λn(x) }
	}
	return func() v { return λn() }
}
func cnd(x l, a map[v]v) v { // conditional, case
	if len(x)%2 == 0 {
		return e("args")
	}
	def := x[len(x)-1]
	x = x[:len(x)-1]
	for i := 0; i < len(x); i += 2 {
		if pidx(eva(x[i], a)) > 0 {
			return eva(x[i+1], a)
		}
	}
	return eva(def, a)
}
func asn(x l, g bool, a map[v]v) v { // assignment
	var mod v
	if len(x) == 3 {
		mod = x[0]
		x = x[1:]
	}
	if len(x) != 2 {
		return e("length")
	}
	var idx l
	s, o := x[0].(s)
	if !o {
		if xl, o := x[0].(l); o && len(xl) > 1 && iss(xl[0]) {
			s = xl[0].(string)
			idx = make(l, len(xl)-1)
			for i, u := range xl[1:] {
				idx[i] = eva(u, a)
			}
		} else {
			return e("assign:type")
		}
	}
	y := eva(x[1], a)
	if g {
		a = ktr(a)
	}
	var r v = y
	if mod != nil {
		x := lupr(a, s)
		if idx != nil {
			x = cal(x, idx, a)
		}
		r = cal(mod, []v{x, y}, a)
	}
	if idx == nil {
		a[s] = cp(r)
		return r
	}
	var u v
	u, a = lup(a, s)
	y = xas(u, r, idx)
	a[s] = cp(y)
	return y
}
func xas(u, y v, idx l) (r v) {
	var get func(v) v
	var set func(x, y v)
	if d, o := md(u); o {
		get = func(x v) v { _, r := d.at(x); return r } // maybe nil
		set = d.set
		defer func() {
			r = d.mp()
		}()
	} else {
		ul, t := ls(u)
		defer func() {
			r = sl(ul, t)
			if l, o := r.(l); o {
				r = uf(l)
			}
		}()
		get = func(x v) v { return ul[pidx(x)] }
		set = func(x, y v) { ul[pidx(x)] = y }
	}
	if len(idx) == 0 {
		idx = l{nil}
	}
	d := len(idx)
	i0 := idx[0]
	idx = idx[1:]
	if i0 == nil {
		i0 = til(cnt(u))
	}
	n := ln(i0)
	if n < 0 {
		i0, n, y = l{i0}, 1, l{y}
	}
	m := ln(y)
	if m < 0 {
		q := make(l, n)
		for i := range q {
			q[i] = cp(y)
		}
		y = q
		m = n
	}
	if m != n {
		return e("size")
	}
	yl, _ := ls(y)
	for i := range yl {
		var w v
		k := at(i0, i)
		if d > 1 {
			w = xas(get(k), yl[i], idx)
		} else {
			w = at(y, i)
		}
		set(k, w)
	}
	return u
}

type tail l
