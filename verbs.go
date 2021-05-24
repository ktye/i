package k

import (
	. "github.com/ktye/wg/module"
)

func nyi(x K) K { trap(Nyi); return x }
func Dex(x, y K) K { // x:y
	dx(x)
	return y
}
func Flp(x K) K { // +x
	trap(Nyi)
	return x
}
func Fst(x K) K { // *x
	t := tp(x)
	if t < 16 {
		return x
	}
	n := nn(x)
	if n == 0 {
		trap(Nyi)
	}
	return ati(x, 0)
}
func Lst(x K) K { // *|x
	t := tp(x)
	if t < 16 {
		return x
	}
	n := nn(x)
	if n == 0 {
		return Fst(x)
	}
	return ati(x, n-1)
}

func Cnt(x K) (r K) { // #x
	t := tp(x)
	dx(x)
	if t < 16 {
		return Ki(1)
	}
	return Ki(nn(x))
}
func Til(x K) (r K) {
	if tp(x) != it {
		trap(Type)
	}
	return seq(int32(x))
}
func seq(n int32) (r K) {
	n = maxi(n, 0)
	r = mk(It, n)
	p := int32(r)
	for i := int32(0); i < n; i++ {
		SetI32(p, i)
		p += 4
	}
	return r
}

func Key(x, y K) (r K) { // x!y
	xt, yt := tp(x), tp(y)
	if xt < 16 {
		x = enl(x)
		xt = tp(x)
	}
	if yt < 16 {
		y = enl(y)
		yt = tp(y)
	}
	nx, ny := nn(x), nn(y)
	if nx != ny {
		if nx == 1 && ny > 1 {
			x = ntake(ny, x)
		} else if ny == 1 && nx > 1 {
			y = ntake(ny, x)
		} else {
			trap(Length)
		}
	}
	r = l2(x, y)
	return K(int32(r)) | K(dt)<<59
}

func Tak(x, y K) (r K) { // x#y
	xt := tp(x)
	if xt == it {
		return ntake(int32(x), y)
	}
	trap(Nyi) // set take
	return 0
}
func ntake(n int32, y K) (r K) {
	if n < 0 {
		trap(Value)
	}
	t := tp(y)
	if t < 5 {
		t += 16
		r = mk(t, n)
		s := sz(t)
		rp := int32(r)
		yp := int32(y)
		if s == 1 {
			Memoryfill(rp, yp, n)
		} else {
			for i := int32(0); i < n; i++ {
				SetI32(rp, yp)
				rp += 4
			}
		}
		return r
	} else if t == ft {
		r = mk(Ft, n)
		rp := int32(r)
		f := F64(int32(y))
		for i := int32(0); i < n; i++ {
			SetF64(rp, f)
			rp += 8
		}
		dx(y)
		return r
	} else if t == zt {
		r = ntake(n, rx(K(I64(int32(y)))))
		im := K(I64(8 + int32(y)))
		if im != 0 {
			im = ntake(n, im)
		}
		dx(y)
		return Kz(r, im)
	} else if t < 16 {
		r = mk(Lt, n)
		rp := int32(r)
		for i := int32(0); i < n; i++ {
			SetI64(rp, int64(rx(y)))
			rp += 8
		}
		dx(y)
		return r
	}
	return Atx(y, seq(nn(y)))
}

func Rev(x K) (r K) { // |x
	t := tp(x)
	if t < 16 {
		return x
	}
	if t >= Ft && t != Lt {
		panic(Nyi)
	}
	xn := nn(x)
	if xn < 2 {
		return x
	}
	r = mk(It, xn)
	rp := int32(r) + 4*xn
	for i := int32(0); i < xn; i++ {
		rp -= 4
		SetI32(rp, i)
	}
	return atv(x, r)
}

func Wer(x K) (r K) { // &x
	t := tp(x)
	if t < 16 {
		x = enl(x)
		t = tp(x)
	}
	var n, rp int32
	xn := nn(x)
	xp := int32(x)
	if t == Bt {
		n = sumb(x)
		r = mk(It, n)
		rp = int32(r)
		for i := int32(0); i < xn; i++ {
			if I8(xp) != 0 {
				SetI32(rp, i)
				rp += 4
			}
			xp++
		}
	} else if t == It {
		n = sumi(x)
		r = mk(It, n)
		rp = int32(r)
		for i := int32(0); i < xn; i++ {
			j := I32(xp)
			for k := int32(0); k < j; k++ {
				SetI32(rp, i)
				rp += 4
			}
			xp += 4
		}
	} else {
		trap(Type)
	}
	dx(x)
	return r
}

func Typ(x K) (r K) { // @x
	r = Ki(int32(tp(x)))
	dx(x)
	return r
}

func mini(x, y int32) int32 {
	if x < y {
		return x
	}
	return y
}
func maxi(x, y int32) int32 {
	if x > y {
		return x
	}
	return y
}
func sumb(x K) (r int32) {
	p := int32(x)
	e := p + nn(x)
	for p < e {
		r += I8(p)
		p++
	}
	return r
}
func sumi(x K) (r int32) {
	p := int32(x)
	e := p + 4*nn(x)
	for p < e {
		r += I8(p)
		p += 4
	}
	return r
}
