package main

import (
	. "github.com/ktye/wg/module"
)

func Srt(x K) (r K) { // ^x
	xt := tp(x)
	if xt < 16 {
		trap(Type)
	}
	xn := nn(x)
	if xn < 2 {
		return x
	}
	switch xt - 17 {
	case 0:
		r = srtB(x, xn)
	case 1:
		r = srtC(x, xn)
	case 2:
		r = srtI(x, xn)
	case 3:
		r = srtI(x, xn)
	default: // todo radix-F
		r = atv(x, Asc(rx(x)))
	}
	return r
}
func Asc(x K) K { return grade(x, 343) } // <x
func Dsc(x K) K { return grade(x, 336) } // >x
func grade(x K, f int32) (r K) { // <x >x
	xt := tp(x)
	if xt < 16 || xt > Lt {
		trap(Type)
	}
	n := nn(x)
	if n < 2 {
		dx(x)
		return seq(n)
	}
	r = seq(n)
	rp := int32(r)
	xp := int32(x)
	if n < 16 {
		igrd(rp, xp, n, sz(xt), f+int32(xt))
	} else {
		w := mk(It, n)
		wp := int32(w)
		Memorycopy(wp, rp, 4*n)
		msrt(wp, rp, 0, n, xp, sz(xt), f+int32(xt))
		dx(w)
	}
	dx(x)
	return r
}
func srtB(x K, n int32) (r K) { // ^B
	r = mk(Bt, n)
	rp := int32(r)
	s := n - sumb(int32(x), n)
	Memoryfill(rp, 0, s)
	Memoryfill(rp+s, 1, n-s)
	dx(x)
	return r
}
func srtC(x K, n int32) (r K) {
	r = mk(Ct, n)
	y := ntake(256, Ki(0))
	yp := int32(y)
	xp := int32(x)
	for i := int32(0); i < n; i++ {
		p := yp + 4*I8(xp+i)
		SetI32(p, 1+I32(p))
		continue
	}
	rp := int32(r)
	for i := int32(0); i < 256; i++ {
		s := I32(yp + 4*i)
		Memoryfill(rp, i, s)
		rp += s
		continue
	}
	dx(x)
	dx(y)
	return r
}
func srtI(x K, n int32) (r K) {
	xp := int32(x)
	if n < 16 {
		if I32(xp-4) == 1 {
			isrtsI(xp, n)
			return x
		} else {
			r = mk(tp(x), n)
			Memorycopy(int32(r), xp, 4*n)
			isrtsI(int32(r), n)
			dx(x)
			return r
		}
	}
	return atv(x, Asc(rx(x))) // todo radix-I
}
func isrtsI(xp, n int32) { // insertion sort ints inplace
	for i := int32(1); i < n; i++ {
		x := I32(xp + 4*i)
		j := i - 1
		for j >= 0 && I32(xp+4*j) > x {
			jj := xp + 4*j
			SetI32(4+jj, I32(jj))
			j--
		}
		SetI32(xp+4+4*j, x)
		continue
	}
}

func msrt(x, r, a, b, p, s, f int32) {
	if b-a < 2 {
		return
	}
	c := (a + b) >> 1
	msrt(r, x, a, c, p, s, f)
	msrt(r, x, c, b, p, s, f)
	mrge(x, r, 4*a, 4*b, 4*c, p, s, f)
}
func mrge(x, r, a, b, c, p, s, f int32) {
	i, j := a, c
	var q int32
	for k := a; k < b; k += 4 {
		if i < c && j < b {
			q = Func[f].(f2i)(p+s*I32(x+i), p+s*I32(x+j))
		} else {
			q = 0
		}
		if i >= c || q != 0 {
			SetI32(r+k, I32(x+j))
			j += 4
		} else {
			SetI32(r+k, I32(x+i))
			i += 4
		}
	}
}

func igrd(rp, xp, n, s, f int32) { // insertion grade with comparison
	for i := int32(1); i < n; i++ { // f: gt(<) lt(>)
		x := I32(rp + 4*i)
		j := i - 1
		for j >= 0 {
			if Func[f].(f2i)(xp+s*I32(rp+4*j), xp+s*x) == 0 {
				break
			}
			jj := rp + 4*j
			SetI32(4+jj, I32(jj))
			j--
		}
		SetI32(rp+4+4*j, x)
		continue
	}
}

func guC(xp, yp int32) int32 { return ib(I8(xp) < I8(yp)) }
func guI(xp, yp int32) int32 { return ib(I32(xp) < I32(yp)) }
func guF(xp, yp int32) int32 { return ltf(F64(xp), F64(yp)) }
func guZ(xp, yp int32) int32 { return ltz(F64(xp), F64(xp+8), F64(yp), F64(yp+8)) }
func guS(xp, yp int32) int32 { return ltS(K(I64(xp)), K(I64(yp))) }

func gdC(xp, yp int32) int32 { return ib(I8(xp) > I8(yp)) }
func gdI(xp, yp int32) int32 { return ib(I32(xp) > I32(yp)) }
func gdF(xp, yp int32) int32 { return gtf(F64(xp), F64(yp)) }
func gdZ(xp, yp int32) int32 { return gtz(F64(xp), F64(xp+8), F64(yp), F64(yp+8)) }
func gdS(xp, yp int32) int32 { return gtS(K(I64(xp)), K(I64(yp))) }

func ltS(x, y K) int32 {
	xp, yp := int32(x), int32(y)
	xn, yn := nn(x), nn(y)
	nn(y)
	n := mini(xn, yn)
	for i := int32(0); i < n; i++ {
		xi, yi := I8(xp+i), I8(yp+i)
		if xi != yi {
			return ib(xi < yi)
		}
	}
	return ib(xn < yn)
}
func gtS(x, y K) int32 {
	xp, yp := int32(x), int32(y)
	xn, yn := nn(x), nn(y)
	nn(y)
	n := mini(xn, yn)
	for i := int32(0); i < n; i++ {
		xi, yi := I8(xp+i), I8(yp+i)
		if xi != yi {
			return ib(xi > yi)
		}
	}
	return ib(xn > yn)
}
