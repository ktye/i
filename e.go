package i

import (
	"reflect"
)

func eva(x v, a map[v]v) v {
	if sy, o := x.(s); o {
		return lup(a, sy)
	}
	l, o := cp(x).(l)
	if !o || len(l) == 0 {
		return x
	}
	if len(l) == 3 && l[0] == "." {
		if l[1] == "$" && len(l[2].([]v)) > 2 {
			return cnd(l[2].([]v)[1:], a) // (.;$;(::;A;B;C;...))
		} else if l[1] == "∇" {
			return tail(l[2].([]v)[1:]) // (.;∇;(::;X;Y;Z;...))
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
		if len(l) != 3 {
			return e("nyi:modified assignment")
		}
		s, o := l[1].(s)
		if !o {
			return e("assign:type")
		}
		y := eva(l[2], a)
		if l[0] == "::" {
			a = ktr(a)
		}
		a[s] = cp(y)
		return y
	case "∇":
		return tail(l[1:])
	case "λ":
		return ead(l, a)
	default:
		f := l[0]
		if u, o := f.(s); o {
			f = lup(a, u)
		} else if u, o := f.([]v); o {
			f = ead(u, a)
		}
		if k := rval(f).Kind(); k != reflect.Func {
			return e("type:func?" + k.String())
		}
		if len(l) == 3 && l[2] == nil { // curry
			x := eva(l[1], a)
			return func(y v) v {
				return cal(f, []v{x, y}, a)
			}
		}
		for i := len(l) - 1; i > 0; i-- { // right to left
			l[i] = eva(l[i], a)
		}
		return cal(f, l[1:], a)
	}
	return e("impossible")
}
func ead(u l, a map[v]v) v { // evaluate adverb expr
	// TODO: verb trains
	if u[0].(s) == "λ" {
		return λ(u[1].(l), a)
	}
	af := lup(a, u[0].(s)).(func(v) v)
	w := eva(u[1], a)
	return af(w) // func(...v)v
}
func lup(a map[v]v, s s) v { // lookup
	if r := a[s]; r != nil {
		return r
	}
	if p, o := a[".."]; o {
		pp, o := p.(*map[v]v)
		if !o {
			return e("type")
		}
		return lup(*pp, s)
	}
	return e("undefined:" + s)
}
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

type tail l
