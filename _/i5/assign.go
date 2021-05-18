package k

func (k *K) Lookup(s string) T {
	v, o := k.Var[s]
	if !o {
		panic("undefined")
	}
	return rx(v)
}
func (k *K) Assign(s string, v T) T {
	if a, o := k.Var[s]; o {
		dx(a)
	}
	k.Var[s] = rx(v)
	return v
}

type copula struct {
	s string
	v []string
	g bool
	i bool
	d bool
	f F2
}

//: x:3 /3
//: (x;y):1 /1
//: (x;y):1 2 /1 2  (destructing assign)
//: (x;y):1 2 3 /1 2 3
//: z:(x;y):1 2 3 /1 2 3
//: a+:a:3 /6  (modified assign)
//: a:!5;a[!2]:0 /0 0 2 3 4  (indexed assign)
//: x:!5;x[2 3]+:1 /0 1 3 4 4  (indexed modified)
//: x:2^!6;x[1;2]:9 /(0 1 2;3 4 9)  (at-depth)
//: x:2^!6;x[!2;!2]:9 /(9 9 2;9 9 5)  (matrix-assign)
//: x:2^!6;x[;1]*:10 /(0 10 2;3 40 5)  (column)
func (k *K) assign(x copula) {
	var ix T
	if x.i {
		v, o := k.Var[x.s].(vector)
		if o == false {
			panic("type")
		}
		ix = k.pop()
		y := k.pop()
		var u T
		if x.d {
			u = k.dmend(v, ix, x.f, y)
		} else {
			u = k.amend(v, ix, x.f, y)
		}
		k.Var[x.s] = rx(u)
		k.push(u)
		return
	}
	y := k.pop()
	if x.v == nil {
		if x.f != nil {
			r := x.f.call2(k.Var[x.s], y)
			k.Var[x.s] = r
			k.push(rx(r))
		} else {
			k.push(k.Assign(x.s, y))
		}
	} else {
		n := count(rx(y)).(int)
		yv, isyv := y.(vector)
		for i, s := range x.v {
			if i == len(x.v)-1 && n > len(x.v) {
				rx(k.Assign(s, ndrop(i, rx(y))))
			} else if isyv {
				yv.ref()
				rx(k.Assign(s, atv(yv, i)))
			} else {
				rx(k.Assign(s, rx(y)))
			}
		}
		k.push(y)
	}
}

//@ @[1 2 3;1;0] /1 0 3
//@ @[1 2 3;1 2;0] /1 0 0
//@ @[1 2 3;1 2;4 5] /1 4 5
//@ @[1 2 3;1;0 0] /(1;0 0;3)
//@ @[1 2 3;1;2.] /(1;2.;3)
//@ @[1 2 3;1;+;2] /1 4 3
func (k *K) amend(x, i, f, y T) T {
	xv := x.(vector)
	if f != nil {
		xv.ref()
		y = k.call(f, l2(atv(xv, rx(i)), y))
	}
	xt, yt := numtype(xv), numtype(y)
	yv, isyv := y.(vector)
	iv, isiv := i.(I)
	if xt != yt || xt == 0 || (isiv == false && isyv) {
		xv = explode(xv)
	}
	xv = use(xv)
	if isiv && isyv {
		if iv.ln() != yv.ln() {
			panic("length")
		}
		yv.unref()
		iv.unref()
		xv.setv(i.(I).v, yv)
	} else if isiv && isyv == false {
		iv.unref()
		yv = ntake(iv.ln(), y).(vector)
		xv.setv(iv.v, yv)
	} else {
		xv.set(i.(int), y)
	}
	return xv
}

//. .[(1;2 3);1 0;5] /(1;5 3)
//. .[(1;2 3);1 0;+;5] /(1;7 3)
//. .[(1 2 3;4 5 6);(1;1 2);+;5] /(1 2 3;4 10 11)
//. .[(1 2 3;4 5 6);(0 1;1 2);+;5] /(1 7 8;4 10 11)
//. .[(1 2 3;4 5 6);(;2);+;1] /(1 2 4;4 5 7)
func (k *K) dmend(x, i, f, y T) T {
	switch v := i.(type) {
	case int:
		return k.amend(x, i, f, y)
	case vector:
		n := v.ln()
		if n == 0 {
			panic("length")
		} else if n == 1 {
			return k.amend(x, first(i), f, y)
		}
		ii := first(rx(i))
		if ii == nil {
			ii = seq(0, count(rx(x)).(int))
		}
		i = ndrop(1, i)
		if iv, o := ii.(vector); o {
			xv, o := x.(vector)
			if o == false {
				panic("type")
			}
			xn := xv.ln()
			for j := 0; j < xn; j++ {
				jj := iv.at(j)
				x = k.amend(x, jj, nil, k.dmend(k.atx(rx(x), jj), rx(i), f, y))
			}
			dx(i)
			return x
		} else {
			return k.amend(x, ii, nil, k.dmend(k.atx(rx(x), ii), i, f, y))
		}
	default:
		panic("type")
	}
}
