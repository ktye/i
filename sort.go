package main

import (
	. "github.com/ktye/wg/module"
)

func Srt(x K) (r K) { // ^x
	xt := tp(x)
	if xt < 16 {
		trap(Type)
	}
	if xt == Dt {
		r = x0(x)
		x = r1(x)
		i := rx(Asc(rx(x)))
		return Key(atv(r, i), atv(x, i))
	}
	xn := nn(x)
	if xn < 2 {
		return x
	}
	switch xt - 18 {
	case 0:
		r = srtC(x, xn)
	case 1:
		r = srtI(x, xn)
	case 2:
		r = srtI(x, xn)
	case 3:
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
	if xt < 16 {
		trap(Type)
	}
	if xt == Dt {
		r = x0(x)
		return Atx(r, grade(r1(x), f))
	}
	n := nn(x)
	if xt == Tt {
		return kxy(104, x, Ki(I32B(f == 336))) //gdt ngn:{(!#x){x@<y x}/|.+x}
	}
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
func srtI(x K, n int32) K {
	if n < 32 {
		x = use(x)
		isrtsI(int32(x), n)
		return x
	}
	// return atv(x, Asc(rx(x))) // merge?
	return radixI(x, n)
}
func srtF(x K, n int32) K {
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
L:
	for ko := int32(0); ko < 32; ko += 8 {
		s := int32(1)
		var prev int32 = nai
		Memoryfill(op, 0, 1024)
		for i := int32(0); i < n; i += 4 {
			e := I32(fr + i)
			ok := op + 4*(255&(e>>ko))
			SetI32(ok, 1+I32(ok))
			if s != 0 {
				s = I32B(e >= prev)
				prev = e
			}
			continue
		}
		if s != 0 {
			if (ko>>3)%2 == 1 {
				Memorycopy(to, fr, n)
			}
			break L
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
		t := to
		to = fr
		fr = t
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
		if v != v {
			SetF64(xp+i, F64(xp+na))
			SetF64(xp+na, v)
			na += 8
		}
	}
	fr := int32(x) + na
	to := int32(b)
	n -= na

	var u uint64
L:
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
				s = I32B(v >= prev)
				prev = v
			}
		}
		if s != 0 {
			if (ko>>3)%2 == 1 {
				Memorycopy(to, fr, n)
			}
			break L
		}
		radixp(op, 0, 256, 0)
		for i := int32(0); i < n; i += 8 {
			v := I64(fr + i)
			u = floatflp(uint64(v))
			ok := op + 4*int32(255&(u>>uint64(ko)))
			SetI64(to+8*I32(ok), v)
			SetI32(ok, 1+I32(ok))
		}
		t := to
		to = fr
		fr = t
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

func guC(xp, yp int32) int32 { return I32B(I8(xp) < I8(yp)) }
func guI(xp, yp int32) int32 { return I32B(I32(xp) < I32(yp)) }
func guF(xp, yp int32) int32 { return ltf(F64(xp), F64(yp)) }
func guZ(xp, yp int32) int32 { return ltz(F64(xp), F64(xp+8), F64(yp), F64(yp+8)) }
func guL(xp, yp int32) int32 { return ltL(K(I64(xp)), K(I64(yp))) }

func gdC(xp, yp int32) int32 { return I32B(I8(xp) > I8(yp)) }
func gdI(xp, yp int32) int32 { return I32B(I32(xp) > I32(yp)) }
func gdF(xp, yp int32) int32 { return guF(yp, xp) }
func gdZ(xp, yp int32) int32 { return guZ(yp, xp) }
func gdL(xp, yp int32) int32 { return guL(yp, xp) }

func ltL(x, y K) (r int32) { // sort lists lexically
	xt := tp(x)
	if xt != tp(y) {
		return I32B(xt < tp(y))
	}
	if xt < 16 {
		return int32(Les(rx(x), rx(y)))
	}
	xp, yp := int32(x), int32(y)
	if xt > Lt {
		a, b := K(I64(xp)), K(I64(yp))
		if match(a, b) == 0 {
			return ltL(a, b)
		}
		return ltL(K(I64(xp+8)), K(I64(yp+8)))
	}
	xn, yn := nn(x), nn(y)
	n := mini(xn, yn)
	switch sz(xt) >> 2 {
	case 0:
		r = taoC(xp, yp, n)
	case 1:
		r = taoI(xp, yp, n)
	case 2:
		if xt == Lt {
			r = taoL(xp, yp, n)
		} else {
			r = taoF(xp, yp, n)
		}
	default:
		r = taoZ(xp, yp, n)
	}
	if r == 2 {
		return I32B(xn < yn)
	} else {
		return r
	}
}
func taoC(xp, yp, n int32) int32 {
	e := xp + n
	for xp < e {
		if I8(xp) != I8(yp) {
			return I32B(I8(xp) < I8(yp))
		}
		yp++
		xp++
	}
	return 2
}
func taoI(xp, yp, n int32) int32 {
	e := xp + 4*n
	for xp < e {
		if I32(xp) != I32(yp) {
			return I32B(I32(xp) < I32(yp))
		}
		yp += 4
		xp += 4
	}
	return 2
}
func taoL(xp, yp, n int32) int32 {
	e := xp + 8*n
	for xp < e {
		x, y := K(I64(xp)), K(I64(yp))
		if match(x, y) == 0 {
			return ltL(x, y)
		}
		yp += 8
		xp += 8
	}
	return 2
}
func taoF(xp, yp, n int32) int32 {
	e := xp + 8*n
	for xp < e {
		x, y := F64(xp), F64(yp)
		if eqf(x, y) == 0 {
			return ltf(x, y)
		}
		yp += 8
		xp += 8
	}
	return 2
}
func taoZ(xp, yp, n int32) int32 {
	e := xp + 16*n
	for xp < e {
		xr, xi, yr, yi := F64(xp), F64(xp+8), F64(yp), F64(yp+8)
		if eqz(xr, xi, yr, yi) == 0 {
			return ltz(xr, xi, yr, yi)
		}
		yp += 16
		xp += 16
	}
	return 2
}
