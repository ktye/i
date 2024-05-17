package main

import (
	. "github.com/ktye/wg/module"
)

func Srt(x K) K { // ^x
	var r K
	xt := tp(x)
	if xt < 16 {
		trap() //type
	}
	if xt == Dt {
		r = x0(x)
		x = r1(x)
		i := rx(Asc(rx(x)))
		return Key(atv(r, i), atv(x, i))
	}
	if nn(x) < 2 {
		return x
	}
	return atv(x, Asc(rx(x)))
}
func Asc(x K) K { // <x  <`file
	if tp(x) == st {
		return readfile(cs(x))
	}
	return grade(x, 1)
}
func Dsc(x K) K { return grade(x, -1) } //254 // >x
func grade(x K, f int32) K { // <x >x
	var r K
	xt := tp(x)
	if xt < 16 {
		trap() //type
	}
	if xt == Dt {
		r = x0(x)
		return Atx(r, grade(r1(x), f))
	}
	n := nn(x)
	if xt == Tt {
		return cal(lup(Ks(88)), l2(x, Ki(I32B(f == -1)))) //gdt ngn:{(!#x){x@<y x}/|.+x}
	}
	if n < 2 {
		dx(x)
		return seq(n)
	}
	r = seq(n)
	rp := int32(r)
	xp := int32(x)
	w := mk(It, n)
	wp := int32(w)
	Memorycopy(wp, rp, 4*n)
	msrt(wp, rp, 0, n, xp, int32(xt), f)
	dxy(w, x)
	return r
}

func msrt(x, r, a, b, p, t, f int32) {
	if b-a < 2 {
		return
	}
	c := (a + b) >> 1
	msrt(r, x, a, c, p, t, f)
	msrt(r, x, c, b, p, t, f)
	mrge(x, r, 4*a, 4*b, 4*c, p, t, f)
}
func mrge(x, r, a, b, c, p, t, f int32) {
	var q int32
	i, j := a, c
	s := sz(T(t))
	for k := a; k < b; k += 4 {
		if i < c && j < b {
			q = I32B(f == Func[234+t].(f2i)(p+s*I32(x+i), p+s*I32(x+j)))
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
func cmL(xp, yp int32) int32 { // compare lists lexically
	var r int32
	x, y := K(I64(xp)), K(I64(yp))
	xt, yt := tp(x), tp(y)
	if xt != yt {
		return I32B(xt > yt) - I32B(xt < yt)
	}
	if xt < 16 { // 11(derived), 12(proj), 13(lambda), 14(native)?
		xp, yp := int32(x), int32(y)
		return Func[245+xt].(f2i)(xp, yp)
	}
	if xt > Lt {
		xp, yp := int32(x), int32(y)
		r = cmL(xp, yp)
		if r == 0 {
			r = cmL(xp+8, yp+8)
		}
		return r
	}
	xn, yn := nn(x), nn(y)
	xp = int32(x)
	yp = int32(y) - xp
	n := mini(xn, yn)
	s := sz(xt)
	e := xp + n*s
	for xp < e {
		r = Func[234+xt].(f2i)(xp, xp+yp)
		if r != 0 {
			return r
		}
		xp += s
	}
	return I32B(xn > yn) - I32B(xn < yn)
}
