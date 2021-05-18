package k

import "fmt"

/*
//. (1 2;0)[0][1] /2
 //. +- /+-
 //. (*-)@2 3 /-2
 //. (*-)2 3 /-2
 //. ++/ /++/
//. 3+ /3+  (projection)
//. +[3;] /3+
//. +[;3] /+[;3]
//. 3 imag /3 imag
//. + /+
//. -[1] /-1
//. +/[1 2 3] /6
//. -[5;3] /2
 //. $[1;2;3] /2
*/
type monadic interface{ call1(x T) T }
type dyadic interface{ call2(x, y T) T }
type projection struct {
	stringer
	f    T
	x    L
	zero []int
}

func funtype(x T) int {
	switch x.(type) {
	case verb:
		return 3
	case derived:
		return 4
	case projection:
		return 5
	case train:
		return 6
	case F1:
		return 1
	case F2:
		return 2
	case λ:
		return 7
	default:
		return 0
	}
}

func conform(x, y T) int { // -1 atoms
	xv, isxv := x.(vector)
	yv, isyv := y.(vector)
	if isxv == false && isyv == false {
		return -1 // atom
	} else if isxv == true && isyv == false {
		return xv.ln()
	} else if isxv == false && isyv == true {
		return yv.ln()
	} else if xv.ln() != yv.ln() {
		panic("length")
	} else {
		return xv.ln()
	}
}

//. {x+y}[1;2] /3
//. {x+y}.1 2 /3
func (k *K) call(x, y T) T {
	if xv, o := x.(vector); o {
		return k.atdepth(xv, y)
	}
	yv, o := y.(vector)
	if o == false {
		panic("type")
	}
	if l, o := y.(L); o {
		if z := l.nulls(); z != nil {
			return proj(x, l, z)
		}
	}
	n := yv.ln()
	switch f := x.(type) {
	case F12:
		switch n {
		case 1:
			return f.call1(first(y))
		case 2:
			a, b := lsplit(yv)
			return f.call2(a, b)
		case 3, 4:
			return k.special(f, yv)
		default:
			panic("rank")
		}
	case F1:
		return f.call1(first(y))
	case F2:
		if n == 2 {
			a, b := lsplit(yv)
			return f.call2(a, b)
		} else if yv.ln() == 1 {
			return proj(x, explode(yv), []int{1})
		} else {
			panic("rank")
		}
	case λ:
		if n < f.ary {
			z := make([]int, f.ary-n)
			for i := range z {
				z[i] = i + n
			}
			return proj(x, explode(yv), z)
		}
		if f.ary == n {
			for i := 0; i < n; i++ {
				yv.ref()
				k.push(atv(yv, i))
			}
			return f.call(k)
		} else {
			panic("rank")
		}
	default:
		fmt.Printf("%T %#v\n", x, x)
		panic("call")
	}
}

func lsplit(v vector) (x, y T) { v.ref(); return atv(v, 0), atv(v, 1) }

func link(x, y T) train {
	g := y.(F12)
	switch t := x.(type) {
	case train:
		return append(t, g)
	default:
		return train{x.(F12), g}
	}
}
func (t train) call1(x T) T {
	for _, f := range t {
		x = f.call1(x)
	}
	return x
}
func (t train) call2(x, y T) T {
	x = t[0].call2(x, y)
	for _, f := range t[1:] {
		x = f.call1(x)
	}
	return x
}

func (l L) nulls() (r []int) {
	for i, u := range l.v {
		if u == nil {
			r = append(r, i)
		}
	}
	return r
}

func proj(f T, x L, zero []int) T {
	switch v := f.(type) {
	case F12:
		var s string
		if len(zero) == 1 {
			s = v.String()
			if zero[0] == 1 {
				sp := ""
				if len(s) > 1 {
					sp = " "
				}
				s = sstr(x.v[0]) + sp + v.String()
			} else {
				s = v.String() + x.brackets()
			}
		}
		return projection{
			stringer: stringer{s},
			f:        v,
			x:        x,
			zero:     zero,
		}
	default:
		panic("proj")
	}
}

func (k *K) special(xf F12, yv vector) T {
	var ff func(T, T, T, T) T
	v := xf.(Verb)
	switch v.s {
	case "@":
		ff = k.amend
	case ".":
		ff = k.dmend
	default:
		panic("rank")
	}
	x := yv.at(0)
	i := yv.at(1)
	y := yv.at(2)
	var f T
	if yv.ln() == 4 {
		f = yv.at(3)
		f, y = y, f
	}
	dx(yv)
	return ff(x, i, f, y)
}
