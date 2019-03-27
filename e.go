package i

func eva(v v, a map[v]v) v {
	l, o := v.(l)
	if !o {
		return v
	}
	switch l[0] {
	case "`":
		return l[1:]
	case "$":
		if len(l) == 4 {
			return e("nyi: $cond") // do not eval all args
		}
	default:
		for i := len(l) - 1; i > 0; i-- {
			l[i] = eva(l[i], a)
		}
		return cal(l[0], l[1:], a)
	}
	return e("nyi")
}

func lup(a map[v]v, s s) v { // lookup `.a.b.c
	return e("nyi")
}
