package i

import (
	"reflect"
)

func eva(x v, a map[v]v) v {
	if sy, o := x.(s); o {
		return lup(a, sy)
	}
	l, o := x.(l)
	if !o || len(l) == 0 {
		return x
	}
	switch l[0] {
	case nil:
		for i := len(l) - 1; i > 0; i-- { // right to left
			l[i] = eva(l[i], a)
		}
		return l[1:]
	case "`":
		return l[1]
	/* TODO cond
	case "$":
		if len(l) == 4 {
			return e("nyi: $cond") // do not eval all args
		}
	*/
	case ";":
		l = l[1:]
		var r interface{}
		for i := range l { // left to right
			r = eva(l[i], a)
		}
		return r
	case ":":
		if len(l) != 3 {
			return e("nyi:modified assignment")
		}
		s, o := l[1].(s)
		if !o {
			return e("assign:type")
		}
		a[s] = eva(cp(l[2]), a) // TODO namespace
		return l[2]
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
