package i

import "reflect"

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
	case "$":
		if len(l) == 4 {
			return e("nyi: $cond") // do not eval all args
		}
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
		a[s] = cp(l[2]) // TODO namespace
		return l[2]
	default:
		f := l[0]
		if u, o := f.(s); o {
			f = lup(a, u)
		}
		if k := rval(f).Kind(); k != reflect.Func {
			return e("type:func?" + k.String())
		}
		n := 0
		for i := len(l) - 1; i > 0; i-- {
			if l[i] != nil {
				break
			}
			n++
		}
		if n > 0 { // curry
			argv := cp(l).([]v)
			return func(u ...v) v {
				if len(u) != n {
					return e("args")
				}
				copy(argv[len(argv)-n:], u)
				return eva(argv, a)
			}
		}
		for i := len(l) - 1; i > 0; i-- { // right to left
			l[i] = eva(l[i], a)
		}
		return cal(l[0], l[1:], a)
	}
	return e("impossible")
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
