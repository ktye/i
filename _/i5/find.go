package k

import (
	"bytes"
	"fmt"
	"regexp"
)

//~ 1~2 /0b
//~ 0 1 2~!3 /1b
//~ (+)~(+) /0b
func match(x, y T) T {
	if r, o := x.(*regexp.Regexp); o {
		return rematch(r, y)
	}
	dx(x)
	dx(y)
	return eq(x, y)
}

//~ ~!3 /100b
//~ ~(1;010b) /(0b;101b)
func (k *K) not(x T) T {
	switch v := x.(type) {
	case bool:
		return v == false
	case byte:
		return v == 0
	case int:
		return v == 0
	case float64:
		return v == 0
	case complex128:
		return v == 0
	case string:
		return v == ""
	case L:
		return each(f1(k.not), x)
	}
	defer dx(x)
	xv, o := x.(vector)
	if o == false {
		return false
	}
	r := make([]bool, xv.ln())
	switch v := x.(type) {
	case B:
		for i, u := range v.v {
			r[i] = u == false
		}
	case C:
		for i, u := range v.v {
			r[i] = u == 0
		}
	case I:
		for i, u := range v.v {
			r[i] = u == 0
		}
	case F:
		for i, u := range v.v {
			r[i] = u == 0
		}
	case Z:
		for i, u := range v.v {
			r[i] = u == 0
		}
	case S:
		for i, u := range v.v {
			r[i] = u == ""
		}
	default:
		panic("type")
	}
	return KB(r)
}

//? 2 3 4?3 /1
//? 2 3 4?6 /3
//? 2 3 4?!5 /3 3 0 1 2
//? "a13b145"?(/"[0-9]+") /(1 2;4 5 6)
func fnd(x, y T) T {
	if re, o := y.(*regexp.Regexp); o {
		return refind(x, re)
	}
	if s, o := x.(string); o {
		switch s {
		case "json":
			return decodeJson(y)
		default:
			panic("value")
		}
	}
	xv, o := x.(vector)
	if o == false {
		panic("type")
	}
	defer xv.unref()
	yv, o := y.(vector)
	if o == false {
		return ifind(xv, y)
	}
	r := make([]int, yv.ln())
	for i := 0; i < yv.ln(); i++ {
		xv.ref()
		yv.ref()
		r[i] = ifind(xv, atv(yv, i))
	}
	yv.unref()
	return KI(r)
}

//? ?3 2 3 4 2 /2 3 4
//? ?(1;"ab";2;1) /(1;"ab";2)
func uniq(x T) T {
	xv, o := x.(vector)
	if o == false {
		panic("type")
	}
	if xv.ln() < 2 {
		return x
	}
	if l, o := x.(L); o {
		return uniql(l)
	}
	xv = sort(x).(vector)
	xv.ref()
	r := pairs(f2(match), xv).(B)
	for i := range r.v {
		r.v[i] = !r.v[i]
	}
	return atv(xv, r)
}
func uniql(x L) T {
	r := KL([]T{})
	for i := range x.v {
		if ifind(r, x.v[i]) == len(r.v) {
			r.v = append(r.v, rx(x.v[i]))
		}
	}
	dx(x)
	return r
}

//= =1 3 1 1 5 /1 3 5!(0 2 3;,1;,4)
func group(x T) T {
	key := uniq(rx(x)).(vector)
	n := key.(vector).ln()
	vv := make([][]int, n)

	key.ref()
	v := fnd(key, x).(I)
	for i := range v.v {
		m := v.v[i]
		vv[m] = append(vv[m], i)
	}
	val := make([]T, n)
	for i, u := range vv {
		val[i] = KI(u)
	}
	return dict(key, KL(val))
}

func eq(x, y T) bool {
	if identical(x, y) {
		return true
	}
	switch v := x.(type) {
	case bool:
		w, o := y.(bool)
		return o && v == w
	case byte:
		w, o := y.(byte)
		return o && v == w
	case int:
		w, o := y.(int)
		return o && v == w
	case float64:
		w, o := y.(float64)
		return o && feq(v, w)
	case complex128:
		w, o := y.(complex128)
		return o && zeq(v, w)
	case B:
		w, o := y.(B)
		if o && len(v.v) == len(w.v) {
			for i := range v.v {
				if v.v[i] != w.v[i] {
					return false
				}
			}
			return true
		}
	case C:
		w, o := y.(C)
		if o && len(v.v) == len(w.v) {
			for i := range v.v {
				if v.v[i] != w.v[i] {
					return false
				}
			}
			return true
		}
	case I:
		w, o := y.(I)
		if o && len(v.v) == len(w.v) {
			for i := range v.v {
				if v.v[i] != w.v[i] {
					return false
				}
			}
			return true
		}
	case F:
		w, o := y.(F)
		if o && len(v.v) == len(w.v) {
			for i := range v.v {
				if feq(v.v[i], w.v[i]) == false {
					return false
				}
			}
			return true
		}
	case Z:
		w, o := y.(Z)
		if o && len(v.v) == len(w.v) {
			for i := range v.v {
				if zeq(v.v[i], w.v[i]) == false {
					return false
				}
			}
			return true
		}
	case S:
		w, o := y.(S)
		if o && len(v.v) == len(w.v) {
			for i := range v.v {
				if v.v[i] != w.v[i] {
					return false
				}
			}
			return true
		}
	case L:
		w, o := y.(L)
		if o && len(v.v) == len(w.v) {
			for i := range v.v {
				if eq(v.v[i], w.v[i]) == false {
					return false
				}
			}
			return true
		}
	case D:
		w, o := y.(D)
		return o && v.tab == w.tab && eq(v.k, w.k) && eq(v.v, w.v)
	}
	return false
}
func identical(x, y T) bool {
	switch v := x.(type) {
	case B:
		w, o := y.(B)
		return o && len(v.v) == len(w.v) && (len(v.v) == 0 || &v.v[0] == &w.v[0])
	case C:
		w, o := y.(C)
		return o && len(v.v) == len(w.v) && (len(v.v) == 0 || &v.v[0] == &w.v[0])
	case I:
		w, o := y.(I)
		return o && len(v.v) == len(w.v) && (len(v.v) == 0 || &v.v[0] == &w.v[0])
	case F:
		w, o := y.(F)
		return o && len(v.v) == len(w.v) && (len(v.v) == 0 || &v.v[0] == &w.v[0])
	case Z:
		w, o := y.(Z)
		return o && len(v.v) == len(w.v) && (len(v.v) == 0 || &v.v[0] == &w.v[0])
	case S:
		w, o := y.(S)
		return o && len(v.v) == len(w.v) && (len(v.v) == 0 || &v.v[0] == &w.v[0])
	}
	return false
}

func in(x, y T) T {
	_, ix := x.(vector)
	yv, iy := y.(vector)
	if !(ix && iy) {
		panic("type")
	}
	i := fnd(y, x).(I)
	r := make([]bool, len(i.v))
	n := yv.ln()
	for i, j := range i.v {
		r[i] = j < n
	}
	dx(i)
	return KB(r)
}
func contains(x vector, y T) bool { return ifind(x, y) < x.ln() }
func ifind(x vector, y T) int {
	switch v := x.(type) {
	case B:
		if e, o := y.(bool); o {
			for i, u := range v.v {
				if e == u {
					return i
				}
			}
		}
	case C:
		if e, o := y.(byte); o {
			for i, u := range v.v {
				if e == u {
					return i
				}
			}
		}
	case I:
		if e, o := y.(int); o {
			for i, u := range v.v {
				if e == u {
					return i
				}
			}
		}
	case F:
		if e, o := y.(float64); o {
			for i, u := range v.v {
				if e == u {
					return i
				}
			}
		}
	case Z:
		if e, o := y.(complex128); o {
			for i, u := range v.v {
				if e == u {
					return i
				}
			}
		}
	case S:
		if e, o := y.(string); o {
			for i, u := range v.v {
				if e == u {
					return i
				}
			}
		}
	case L:
		for i, u := range v.v {
			if eq(u, y) {
				return i
			}
		}
	}
	return x.ln()
}

//\ "x"\"abxdexxg" /("ab";"de";"";,"g")  (split)
//\ ""\"ab  de f " /("ab";"de";,"f")  (fields)
//\ (/"[0-9]")\"ab3cd2cv4" /("ab";"cd";"cv";"")  (regexp-split)
func (k *K) split(x, y T) T {
	if yc, o := y.(C); o {
		switch v := x.(type) {
		case byte:
			return splitc([]byte{v}, yc)
		case C:
			defer dx(x)
			return splitc(v.v, yc)
		case *regexp.Regexp:
			return resplit(v, yc)
		default:
			panic("type")
		}
	}
	panic("type")
}
func (k *K) split1(x T, y T) T { return k.cut(ucat(mk(0, 1).(I), where(equal(rx(y), x)).(I)), y) }
func splitc(sep []byte, x C) L {
	var r [][]byte
	if len(sep) == 0 {
		r = bytes.Fields(x.v)
	} else {
		r = bytes.Split(x.v, sep)
	}
	l := make([]T, len(r))
	for i := range r {
		l[i] = KC(r[i])
	}
	x.unref()
	return KL(l)
}
func resplit(re *regexp.Regexp, x C) L {
	r := re.Split(string(x.v), -1)
	dx(x)
	l := make([]T, len(r))
	for i := range r {
		l[i] = KC([]byte(r[i]))
	}
	return KL(l)
}

/// "n"/("ab";"cde") /"abncde"
/// 1 2/(!3;!2) /0 1 2 1 2 0 1
/// (1+2)/(1 2;7 8) /1 2 3 7 8
func join(x, y T) T {
	yl, o := y.(L)
	if o == false {
		fmt.Printf("join: %#v\n", y)
		panic("type")
	}
	n := yl.ln()
	if n < 2 {
		dx(x)
		return y
	}

	r := rx(yl.v[0])
	for i := 1; i < n; i++ {
		r = cat(r, rx(x))
		r = cat(r, rx(yl.v[i]))
	}
	dx(x)
	dx(y)
	return r
}

func regex(x T) T { // parse: (/"re")
	c, o := x.(C)
	if o == false {
		panic("type")
	}
	r, e := regexp.Compile(string(c.v))
	c.unref()
	if e != nil {
		panic(e)
	}
	return r
}

//~ (/"[1-3]")~"alpha" /0b
//~ (/"[1-3]")~"alpha24" /1b
func rematch(re *regexp.Regexp, x T) T {
	switch v := x.(type) {
	case C:
		r := re.Match(v.v)
		dx(x)
		return r
	case L:
		return eachleft(f2(match), re, x)
	}
	panic("type")
}

func refind(x T, re *regexp.Regexp) L {
	c, o := x.(C)
	if o == false {
		panic("type")
	}
	ix := re.FindAllIndex(c.v, -1)
	dx(x)
	r := make([]T, len(ix))
	for i, u := range ix {
		r[i] = seq(u[0], u[1])
	}
	return KL(r)
}

/// "pq"("ab")/"alaba" /"alpqa"  (replace)
/// "bb"("a")/"alpha" /"bblphbb"
/// "$2$1"(/".([0-9])..([0-9]).")/"12345678" /"5278"  (regex-replace)
func replace(rep, pat, str T) T {
	defer dx(rep)
	defer dx(pat)
	defer dx(str)
	if l, o := str.(L); o {
		r := make([]T, l.ln())
		for i := range l.v {
			rx(rep)
			rx(pat)
			r[i] = replace(rep, pat, rx(l.v[i]))
		}
		return KL(r)
	}
	s, o := str.(C)
	if o == false {
		panic("type")
	}
	var r []byte
	switch v := rep.(type) {
	case byte:
		r = []byte{v}
	case C:
		r = v.v
	default:
		fmt.Printf("rep = %T\n", rep)
		panic("type")
	}
	switch v := pat.(type) {
	case byte:
		return KC(bytes.Replace(s.v, []byte{v}, r, -1))
	case C:
		return KC(bytes.Replace(s.v, v.v, r, -1))
	case *regexp.Regexp:
		return KC(v.ReplaceAll(s.v, r))
	default:
		panic("type")
	}
}
