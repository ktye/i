package k

//* *3 2 1 /3
//* *2 /2
func first(x T) T {
	xv, o := x.(vector)
	if o == false {
		return x
	}
	if xv.ln() == 0 {
		dx(x)
		return xv.zero()
	}
	return atv(xv, 0)
}

//# 3#1 /1 1 1
//# "abc"#"ab0dbb" /"abbb"
func take(x, y T) T {
	if n, o := x.(int); o {
		return ntake(n, y)
	}
	xv, ix := x.(vector)
	yv, iy := y.(vector)
	if ix && iy {
		return stake(xv, yv)
	}
	panic("nyi:take")
}

func ntake(n int, y T) T {
	if v, o := y.(vector); o {
		m := v.ln()
		var a, b int
		if n < 0 {
			a, b = maxi(0, m+n-1), m
		} else {
			a, b = 0, mini(n, m)
		}
		return atv(v, seq(a, b))
	} else {
		return mk(y, n)
	}
}
func stake(x, y vector) T {
	y.ref()
	return atv(y, in(y, x))
}

//# #222 /1
//# #!5 /5
//# #() /0
func count(x T) T {
	v, o := x.(vector)
	if o == false {
		return 1
	}
	dx(x)
	return v.ln()
}

//! !5 /0 1 2 3 4
//! !{a:3;b:4} /`a`b!3 4
func (k *K) til(x T) T {
	switch v := x.(type) {
	case int:
		return seq(0, v)
	case λ:
		return v.dict(k)
	default:
		panic("type")
	}
}
func seq(a, b int) I {
	n := b - a
	v := make([]int, n)
	for i := range v {
		v[i] = a + i
	}
	r := I{v: v}
	r.init()
	return r
}

//_ 1_2 3 /,3
//_ 2_1 2 /!0
//_ 3_1 2 /!0
//_ -1_2 3 /,2
//_ -5_2 3 /!0
//_ 2 3_!5 /0 1 4
func (k *K) drop(x, y T) T {
	if n, o := x.(int); o {
		return ndrop(n, y)
	}
	xv, ix := x.(vector)
	yv, iy := y.(vector)
	if ix && iy {
		return k.sdrop(xv, yv)
	}
	panic("nyi:drop")
}
func ndrop(n int, y T) T {
	if v, o := y.(vector); o {
		m := v.ln()
		var a, b int
		if n < 0 {
			a, b = 0, maxi(0, n+m)
		} else {
			a, b = mini(n, m), m
		}
		return atv(v, seq(a, b))
	} else {
		panic("type")
	}
}
func (k *K) sdrop(x, y vector) T {
	y.ref()
	return atv(y, k.not(in(y, x)))
}

//! 2!3 /(,2)!,3
func dict(x, y T) T {
	xv, o := x.(vector)
	if o == false {
		xv = ntake(count(rx(y)).(int), x).(vector)
	}
	yv, o := y.(vector)
	if o == false {
		yv = ntake(xv.ln(), y).(vector)
	}
	if xv.ln() != yv.ln() {
		panic("length")
	}
	d := D{k: xv, v: yv}
	d.init()
	return d
}

/* + +`alpha!1 2 /+(,`alpha)!,1 2 hangs */
func flip(x T) T {
	switch v := x.(type) {
	case L:
		return transpose(v)
	case D:
		if v.tab == false {
			if _, o := v.k.(S); o == false {
				panic("type")
			}
			if l, o := v.v.(L); o == false {
				panic("type")
			} else {
				n := 0
				for i, u := range l.v {
					if w, o := u.(vector); o == false {
						panic("type")
					} else if i == 0 {
						n = w.ln()
					} else if w.ln() != n {
						panic("length")
					}
				}
			}
		}
		r := D{tab: !v.tab, k: rx(v.k).(vector), v: rx(v.v).(vector)}
		r.init()
		dx(x)
		return r
	case vector:
		return transpose(explode(v))
	default:
		panic("type")
	}
}

//+ +(1 2 3;4 5 6) /(1 4;2 5;3 6)
//+ +(1 2;"ab") /((1;"a");(2;"b"))
func transpose(x vector) T {
	n := x.ln()
	if n == 0 {
		return x
	}
	x = over(f2(cat), x).(vector)
	nm := x.ln()
	m := nm / n
	if nm != n*m {
		panic("uniform")
	}
	ix := make([]int, n)
	for i := range ix {
		ix[i] = i * m
	}
	r := make([]T, m)
	for i := range r {
		r[i] = x.atv(ix)
		for i := range ix {
			ix[i]++
		}
	}
	dx(x)
	return KL(r)
}

//^ 2^!5 /(0 1;2 3 4)
//^ 2 5^"alphabeta" /("pha";"beta")
//^ "abe"^"albetq" /("al";,"b";"etq")
func (k *K) cut(x, y T) T {
	xv, iv := x.(vector)
	yv, iy := y.(vector)
	if iy == false {
		panic("type")
	}
	switch v := x.(type) {
	case int:
		return k.icut(v, yv)
	case I:
		defer dx(v)
		return cuts(v.v, yv)
	default:
		if iv && iy {
			return k.setcut(xv, yv)
		}
		panic("type") // nyi: f^
	}
}
func (k *K) icut(n int, y vector) T {
	if n < 1 {
		panic("domain")
	}
	r := make([]T, n)
	m := y.ln() / n
	s := seq(0, m)
	for i := range r {
		if i < len(r)-1 {
			y.ref()
			r[i] = atv(y, rx(s))
			for i := range s.v {
				s.v[i] += m
			}
		} else {
			r[i] = k.drop(i*m, y)
		}

	}
	dx(s)
	return KL(r)
}
func cuts(x []int, y vector) T {
	r := make([]T, len(x))
	defer y.unref()
	if len(x) == 0 {
		return KL(r)
	}
	x = append(x, y.ln())
	for i := range r {
		if x[i] < 0 || x[1+i] < x[i] || x[1+i] > y.ln() {
			panic("range")
		}
		y.ref()
		r[i] = atv(y, seq(x[i], x[1+i]))
	}
	return KL(r)
}
func (k *K) setcut(x, y vector) T {
	y.ref()
	return k.cut(where(in(y, x)), y)
}

//| |3 2 8 /8 2 3
func reverse(x T) T {
	xv, o := x.(vector)
	if o == false || xv.ln() < 2 {
		return x
	}
	n := xv.ln()
	v := make([]int, xv.ln())
	for i := range v {
		v[i] = n - 1 - i
	}
	return atv(xv, KI(v))
}

//. ."1+2" /3
func (k *K) val(x T) T {
	switch v := x.(type) {
	case C:
		v.unref()
		b := make([]byte, len(v.v))
		copy(b, v.v)
		return k.exec(k.fold(k.parse(b), b), b)
	default:
		panic("type")
	}
}

func identity(x T) T { return x }
func dex(x, y T) T   { dx(x); return y }
