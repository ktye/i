package k

import (
	"reflect"
)

//, ,1 /,1
//, ,1 2 /,1 2
//, () /()
//, (1;) /(1;)
//, (;1) /(;1)
//, (;;) /(;;)
//, (1;(2;3 4);5) /(1;(2;3 4);5)
func enlist(x T) T {
	_, o := x.(vector)
	if !o {
		return mk(x, 1)
	}
	var r L
	r.v = []T{x}
	r.init()
	return r
}
func l1(x T) L    { return KL([]T{x}) }
func l2(x, y T) L { return KL([]T{x, y}) }
func cat(x, y T) (r T) {
	if x == nil {
		return y
	}
	xv, o := x.(vector)
	if o == false {
		xv = mk(x, 1).(vector)
	}
	yv, o := y.(vector)
	if o == false {
		yv = mk(y, 1).(vector)
	}
	if xv.ln() == 0 {
		dx(xv)
		return yv
	} else if yv.ln() == 0 {
		dx(yv)
		return xv
	}
	if reflect.TypeOf(xv) == reflect.TypeOf(yv) {
		return ucat(xv, yv)
	}
	return catl(explode(xv), explode(yv))
}
func lcat(x L, y T) L {
	r := use(x).(L)
	r.v = append(r.v, y)
	return r
}
func ucat(x, y vector) vector {
	if x.unref() > 0 {
		r := mk(x.zero(), x.ln()+y.ln()).(vector)
		switch v := r.(type) {
		case B:
			copy(v.v, x.(B).v)
			copy(v.v[x.ln():], y.(B).v)
		case C:
			copy(v.v, x.(C).v)
			copy(v.v[x.ln():], y.(C).v)
		case I:
			copy(v.v, x.(I).v)
			copy(v.v[x.ln():], y.(I).v)
		case F:
			copy(v.v, x.(F).v)
			copy(v.v[x.ln():], y.(F).v)
		case Z:
			copy(v.v, x.(Z).v)
			copy(v.v[x.ln():], y.(Z).v)
		case S:
			copy(v.v, x.(S).v)
			copy(v.v[x.ln():], y.(S).v)
		default:
			panic("ucat")
		}
		dx(y)
		return r
	} else {
		rx(x)
		dx(y)
		switch v := x.(type) {
		case B:
			v.v = append(v.v, y.(B).v...)
			return v
		case C:
			v.v = append(v.v, y.(C).v...)
			return v
		case I:
			v.v = append(v.v, y.(I).v...)
			return v
		case F:
			v.v = append(v.v, y.(F).v...)
			return v
		case Z:
			v.v = append(v.v, y.(Z).v...)
			return v
		case S:
			v.v = append(v.v, y.(S).v...)
			return v
		default:
			panic("ucat")
		}
	}
}
func catl(x, y L) L {
	l := make([]T, x.ln()+y.ln())
	for i := 0; i < x.ln(); i++ {
		l[i] = rx(x.v[i])
	}
	k := x.ln()
	for i := 0; i < y.ln(); i++ {
		l[k+i] = rx(y.v[i])
	}
	r := L{v: l}
	r.init()
	dx(x)
	dx(y)
	return r
}
func explode(x vector) L {
	if l, o := x.(L); o {
		return l
	}
	r := make([]T, x.ln())
	switch v := x.(type) {
	case B:
		for i := range r {
			r[i] = v.v[i]
		}
	case C:
		for i := range r {
			r[i] = v.v[i]
		}
	case I:
		for i := range r {
			r[i] = v.v[i]
		}
	case F:
		for i := range r {
			r[i] = v.v[i]
		}
	case Z:
		for i := range r {
			r[i] = v.v[i]
		}
	case S:
		for i := range r {
			r[i] = v.v[i]
		}
	}
	x.unref()
	u := L{v: r}
	u.init()
	return u
}

func (l L) uf() vector {
	if len(l.v) < 2 {
		return l
	}
	t := numtype(l.v[0])
	for _, u := range l.v {
		if _, o := u.(vector); o {
			return l
		}
		if t != numtype(u) || t == 0 {
			return l
		}
	}
	n := len(l.v)
	defer l.unref()
	switch l.v[0].(type) {
	case bool:
		r := make([]bool, n)
		for i := range r {
			r[i] = l.v[i].(bool)
		}
		return KB(r)
	case byte:
		r := make([]byte, n)
		for i := range r {
			r[i] = l.v[i].(byte)
		}
		return KC(r)
	case int:
		r := make([]int, n)
		for i := range r {
			r[i] = l.v[i].(int)
		}
		return KI(r)
	case float64:
		r := make([]float64, n)
		for i := range r {
			r[i] = l.v[i].(float64)
		}
		return KF(r)
	case complex128:
		r := make([]complex128, n)
		for i := range r {
			r[i] = l.v[i].(complex128)
		}
		return KZ(r)
	case string:
		r := make([]string, n)
		for i := range r {
			r[i] = l.v[i].(string)
		}
		return KS(r)
	default:
		panic("type")
	}
	return l
}
