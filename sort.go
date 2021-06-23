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
	case 4:
		r = srtF(x, xn)
	default:
		r = atv(x, Asc(rx(x)))
	}
	return r
}
func Asc(x K) K { // <x  <`file
	if tp(x) == st {
		return readfile(cs(x))
	}
	return grade(x, 343)
}
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
	if n < 32 {
		x = use(x)
		isrtsI(int32(x), n)
		return x
	}
	// return atv(x, Asc(rx(x))) // merge?
	return radixI(x, n)
}
func srtF(x K, n int32) (r K) {
	if n < 32 { // todo limit?
		return atv(x, Asc(rx(x)))
	}
	return radixF(x, n)
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

func radixI(x K, n int32) K { // ^I  see:shawnsmithdev/zermelo
	x = use(x)
	b := ntake(n, Ki(0))
	o := mk(It, 256)
	n *= 4
	op := int32(o)
	fr := int32(x)
	to := int32(b)
	for ko := int32(0); ko < 32; ko += 8 {
		s := int32(1)
		var prev int32 = -2147483648
		Memoryfill(op, 0, 1024)
		for i := int32(0); i < n; i += 4 {
			e := I32(fr + i)
			ok := op + 4*(255&(e>>ko))
			SetI32(ok, 1+I32(ok))
			if s != 0 {
				s = ib(e >= prev)
				prev = e
			}
			continue
		}
		if s != 0 {
			if (ko>>3)%2 == 1 {
				Memorycopy(to, fr, n)
			}
			break
		}
		w := int32(0)
		if ko == 24 {
			w = radixp(op, 0, 128, radixp(op, 128, 256, w))
		} else {
			w = radixp(op, 0, 256, w)
		}
		for i := int32(0); i < n; i += 4 {
			e := I32(fr + i)
			ok := op + 4*(255&(e>>ko))
			SetI32(to+4*I32(ok), e)
			SetI32(ok, 1+I32(ok))
			continue
		}
		to, fr = swap(to, fr)
	}
	dx(o)
	dx(b)
	return x
}
func radixp(op, a, b, w int32) int32 {
	op += 4 * a
	for i := a; i < b; i++ {
		c := I32(op)
		SetI32(op, w)
		w += c
		op += 4
		continue
	}
	return w
}
func radixF(x K, n int32) K { // ^F
	x = use(x)
	b := ntake(n, Kf(0))
	o := mk(It, 256)
	n *= 8
	op := int32(o)

	xp := int32(x)
	na := int32(0)
	for i := int32(0); i < n; i += 8 {
		v := F64(xp + i)
		if isnan(v) {
			SetF64(xp+i, F64(xp+na))
			SetF64(xp+na, v)
			na += 8
		}
	}
	fr := int32(x) + na
	to := int32(b)
	n -= na

	var u uint64
	for ko := int32(0); ko < 64; ko += 8 {
		s := int32(1)
		prev := float64(0)
		Memoryfill(op, 0, 1024)
		for i := int32(0); i < n; i += 8 {
			u = floatflp(uint64(I64(fr + i)))
			ok := op + 4*int32(255&(u>>uint64(ko)))
			SetI32(ok, 1+I32(ok))
			if s != 0 {
				v := F64(fr + i)
				s = ib(v >= prev)
				prev = v
			}
		}
		if s != 0 {
			if (ko>>3)%2 == 1 {
				Memorycopy(to, fr, n)
			}
			break
		}
		radixp(op, 0, 256, 0)
		for i := int32(0); i < n; i += 8 {
			v := I64(fr + i)
			u = floatflp(uint64(v))
			ok := op + 4*int32(255&(u>>uint64(ko)))
			SetI64(to+8*I32(ok), v)
			SetI32(ok, 1+I32(ok))
		}
		to, fr = swap(to, fr)
	}
	dx(o)
	dx(b)
	return x
}
func floatflp(x uint64) uint64 {
	if (x & 0x8000000000000000) == 0x8000000000000000 {
		return x ^ 0xFFFFFFFFFFFFFFFF
	}
	return x ^ 0x8000000000000000
}

func swap(x, y int32) (int32, int32) { return y, x }

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
