package k

//@ 4 3 2 1@0 /4
//@ 4 3 2 1@0 3 /4 1
//@ 1 2 3[1] /2
//@ 1 2 4[3-1] /4
//@ 1 2 3  1 /2
//@ (0.+!10)(1;2 3) /(1.;2 3.)
func (k *K) atx(x, y T) T {
	switch v := x.(type) {
	case string:
		return encode(v, y)
	case F1:
		return v.call1(y)
	case λ:
		if v.ary == 1 {
			k.push(y)
			return v.call(k)
		} else {
			panic("rank")
		}
	case vector:
		return atv(v, y)
	default:
		if yv, o := y.(vector); o {
			dx(y)
			return ntake(yv.ln(), x)
		} else {
			return x
		}
	}
}
func atv(x vector, y T) T {
	defer dx(x)
	defer dx(y)
	switch v := y.(type) {
	case B:
		if len(v.v) != x.ln() {
			panic("length")
		}
		x.ref()
		return x.atv(whereb(v.v))
	case int:
		i, o := y.(int)
		if o == false {
			panic("type")
		}
		return x.at(i)
	case I:
		return x.atv(v.v)
	case L:
		r := make([]T, len(v.v))
		for i := range r {
			x.ref()
			r[i] = atv(x, rx(v.v[i]))
		}
		return KL(r)
	default:
		panic("type")
	}
}

//. (1 2 3).(1) /2
//. (1 2 3;4 5 6).(1 2) /6
//. (1 2 3;4 5 6)[1;2] /6
//. (1 2 3;4 5 6)[1 0;1] /5 2
//. (1 2 3;4 5 6)[0 1;1 0] /(2 1;5 4)
func (k *K) atdepth(x vector, y T) T {
	switch v := y.(type) {
	case int:
		return atv(x, v)
	case I:
		var r T = x
		for _, i := range v.v {
			r = k.atx(r, i)
		}
		dx(y)
		return r
	case L:
		defer dx(y)
		if len(v.v) == 0 {
			return mk(x.zero(), 0)
		} else if len(v.v) != 2 {
			panic("rank")
		}
		y0, y1 := rx(v.v[0]), rx(v.v[1])
		r := atv(x, y0)
		return eachleft(f2(k.atx), r, y1)
	default:
		panic("type")
	}
}

//@ `k("a\nb c";1 2;3.0 4) /"(\"a\\nb c\";1 2;3 4.)"
func encode(s string, x T) T {
	switch s {
	case "k":
		return kst(x)
	case "json":
		return encodeJson(x)
	default:
		panic("value")
	}
}

//& &2 0 1 /0 0 2
func where(x T) T {
	switch v := x.(type) {
	case B:
		dx(x)
		return KI(whereb(v.v))
	case I:
		dx(x)
		return KI(repeat(v.v))
	default:
		panic("type")
	}
}
func whereb(b []bool) []int {
	n := 0
	for _, v := range b {
		if v {
			n++
		}
	}
	r := make([]int, n)
	k := 0
	for i, v := range b {
		if v {
			r[k] = i
			k++
		}
	}
	return r
}
func repeat(v []int) []int {
	n := 0
	for _, u := range v {
		n += u
	}
	r := make([]int, n)
	k := 0
	for i, u := range v {
		if u < 0 {
			panic("domain")
		} else if u > 0 {
			for j := 0; j < u; j++ {
				r[k] = i
				k++
			}
		}
	}
	return r
}
