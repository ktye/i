package k

func (k *K) derive(a adverb, f T, s string) (r derived) {
	//fmt.Printf("derive: a=%v f=%#v\n", a, f)
	r.s = s
	switch a {
	case "'":
		r.f1 = func(x T) T { return each(f, x) }
		r.f2 = func(x, y T) T { return each2(f, x, y) }
	case "/":
		r.f1 = func(x T) T { return over(f, x) }
		r.f2 = func(x, y T) T { return over2(f, x, y) }
	case "\\":
		r.f1 = func(x T) T { return k.scan(f, x) }
		r.f2 = func(x, y T) T { return scan2(f, x, y) }
	case "':":
		r.f1 = func(x T) T { return pairs(f, x) }
		r.f2 = func(x, y T) T { return pairs2(f, x, y) }
	case "/:":
		r.f1 = func(x T) T { return fix(f, x) }
		r.f2 = func(x, y T) T { return eachright(f, x, y) }
	case "\\:":
		r.f1 = func(x T) T { return scanfix(f, x) }
		r.f2 = func(x, y T) T { return eachleft(f, x, y) }
	default:
		panic("value")
	}
	return r
}

//' -'1 2 3 /-1 -2 -3
func each(f, x T) (r T) {
	g := f.(F1)
	v, o := x.(vector)
	if o == false {
		return g.call1(x)
	}
	defer dx(f)
	defer dx(x)
	n := v.ln()
	l := make([]T, n)
	for i := range l {
		v.ref()
		l[i] = g.call1(atv(v, i))
	}
	return KL(l).uf()
}

//' (1;2 3)+'(2;4 5) /(3;6 8)
//' (1;2 3)+'5 /(6;7 8)
func each2(f, x, y T) (r T) {
	g := f.(F2)
	n := conform(x, y)
	if n == 0 {
		dx(f)
		dx(x)
		return ntake(0, y)
	} else if n < 0 {
		return g.call2(x, y)
	}
	xv, o := x.(vector)
	if o == false {
		xv = ntake(n, x).(vector)
	}
	yv, o := y.(vector)
	if o == false {
		yv = ntake(n, y).(vector)
	}
	l := make([]T, n)
	for i := range l {
		l[i] = g.call2(xv.at(i), yv.at(i))
	}
	dx(f)
	dx(xv)
	dx(yv)
	return KL(l).uf()
}

/// +/1 2 3 /6
/// (+)/1 2 3 /6
func over(f, x T) (r T) {
	g, o := f.(F2)
	if o == false {
		return join(f, x)
	}
	xv, o := x.(vector)
	if o == false {
		dx(f)
		return x
	}
	n := xv.ln()
	if n == 0 {
		r = xv.zero()
	} else {
		r = first(rx(x))
		for i := 1; i < n; i++ {
			xv.ref()
			r = g.call2(r, atv(xv, i))
		}
	}
	dx(f)
	dx(x)
	return r
}
func overI(f func(int, int) int, x []int) (r int) {
	r = x[0]
	for _, i := range x {
		r = f(r, i)
	}
	return r
}

func over2(f, x, y T) T {
	g, o := f.(F2)
	if o == false {
		return replace(x, f, y)
	}
	yv, o := y.(vector)
	if o == false {
		return g.call2(x, y)
	}
	n := yv.ln()
	for i := 0; i < n; i++ {
		yv.ref()
		x = g.call2(x, atv(yv, i))
	}
	dx(f)
	dx(y)
	return x
}

//\ -\1 2 3 /1 -1 -4
func (k *K) scan(f, x T) (r T) {
	g, o := f.(F2)
	if o == false {
		return k.split(f, x)
	}
	xv, o := x.(vector)
	if o == false {
		dx(f)
		return x
	}
	t := first(rx(x))
	r = enlist(rx(t))
	n := xv.ln()
	for i := 1; i < n; i++ {
		xv.ref()
		t = g.call2(t, atv(xv, i))
		r = cat(r, rx(t))
	}
	dx(f)
	dx(x)
	dx(t)
	return r
}
func scan2(f, x, y T) (r T) {
	g := f.(F2)
	yv, o := y.(vector)
	if o == false {
		return g.call2(x, y)
	}
	n := yv.ln()
	if n == 0 {
		dx(f)
		dx(x)
		return y
	}
	t := g.call2(rx(x), first(rx(y)))
	r = enlist(rx(t))
	for i := 1; i < n; i++ {
		yv.ref()
		t = g.call2(t, atv(yv, i))
		r = cat(r, rx(t))
	}
	dx(f)
	dx(x)
	dx(y)
	dx(t)
	return r
}
func eachleft(f, x, y T) T {
	g := f.(F2)
	xv, o := x.(vector)
	if o == false {
		return g.call2(x, y)
	}
	n := xv.ln()
	r := make([]T, n)
	for i := range r {
		xv.ref()
		r[i] = g.call2(atv(xv, i), rx(y))
	}
	dx(f)
	dx(x)
	dx(y)
	return KL(r).uf()
}
func eachright(f, x, y T) T {
	g := f.(F2)
	yv, o := y.(vector)
	if o == false {
		return g.call2(x, y)
	}
	n := yv.ln()
	r := make([]T, n)
	for i := range r {
		yv.ref()
		r[i] = g.call2(rx(x), atv(yv, i))
	}
	dx(f)
	dx(x)
	dx(y)
	return KL(r).uf()
}

//': <':4 2 1 0 2 /01110b
func pairs(f, x T) (r T) {
	xv, o := x.(vector)
	if o == false {
		panic("type")
	}
	return pairs2(f, xv.zero(), x)
}
func pairs2(f, x, y T) T {
	g := f.(F2)
	yv, o := y.(vector)
	if o == false {
		return g.call2(x, y)
	}
	n := yv.ln()
	if n == 0 {
		dx(f)
		dx(x)
		return y
	}
	yv.ref()
	var r T
	for i := 0; i < n; i++ {
		yv.ref()
		t := atv(yv, i)
		r = cat(r, g.call2(rx(t), x))
		x = t
	}
	dx(x)
	dx(f)
	dx(y)
	return r
}

func fix(f, x T) T {
	panic("nyi")
}

func scanfix(f, x T) T {
	panic("nyi")
}
