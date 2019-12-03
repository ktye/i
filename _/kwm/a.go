package main

import "syscall/js"

// max memory, heap size must be set explicitly when compiling with tinygo
const maxmem = 1 << 23 // number of floats (64 MB)
var mem []f
var obuf k
var imgp k

func main() {
	mem = make([]f, maxmem)
	ini(mem[:1<<13])
	table[21] = red
	table[21+dyad] = wrt
	table[39] = trp
	table[29+dyad] = drw
	mkk(".k", "{(`key;x)}")
	mkk(".m", "{(`mouse;x)}")
	obuf = mk(C, 0)
	asn(mks(".fs"), key(mk(S, 0), mk(L, 0)), inc(null))
	js.Global().Set("kio", maxmem>>17)
	select {}
}
func grw(c k) {
	if len(m.f) >= maxmem {
		panic("ws")
	}
	m.f = mem[:2*len(m.f)]
}
func red(x k) (r k) { // 1:x (from .fs)
	t, n := typ(x)
	if t == C {
		x, t, n = c2s(x), S, atom
	} else if t != S || n != atom {
		panic("type")
	}
	return atx(lup(mks(".fs")), x)
}
func wrt(x, y k) (r k) { // x 1:y
	xt, yt, xn, yn := typs(x, y)
	if yt != C {
		panic("type")
	} else if yn == atom {
		y, yn = enl(y), 1
	}

	if xt == C {
		x, xt, xn = c2s(x), S, atom
	}
	if xt != S || xn != atom {
		panic("type")
	}
	if m.k[2+x] == 0 {
		obuf = cat(obuf, y)
		return x
	}
	f := mk(N+2, atom)
	m.k[2+f] = 12 + dyad // ,
	dec(asn(mks(".fs"), key(inc(x), y), f))
	return x
}
func trp(x, y k) (r k) { panic("nyi") }
func drw(x, y k) (r k) { // x 9:y (draw)
	xt, yt, xn, yn := typs(x, y)
	if xt != I || yt != I || xn != atom || yn == atom {
		panic("type")
	}
	for i := k(0); i < yn; i++ { // opaque
		m.k[2+i+y] |= 0xFF000000
	}
	w := m.k[2+x]
	h := yn / w
	if w*h != yn {
		panic("length")
	}
	p := ptr(y, C)
	imgp = p
	js.Global().Set("imgw", w)
	js.Global().Set("imgh", h)
	return decr(x, y, inc(null))
}

//go:export K
func K() { // execute k string via js variable kio
	if m.k[obuf]&atom != 0 {
		obuf = take(0, 0, obuf)
	}
	x := mkb([]c(js.Global().Get("kio").String()))
	js.Global().Set("kio", "")
	evp(x)
	if t, n := typ(obuf); n != 0 {
		if t != C || n == atom {
			panic("type")
		}
		p := 8 + obuf<<2
		js.Global().Set("kio", s(m.c[p:p+n]))
	}
}

//go:export P
func P() *byte { return &m.c[0] } // k-memory offset in wasm memory buffer

//go:export Img
func Img() *byte { return &m.c[imgp] } // pointer to current image data

//go:export Srcp
func Srcp() int { return int(m.k[srcp]) } // source pointer (error indicator)

//go:export Us
func Us(w, h int) { // store canvas size
	dec(asn(mks("w"), mki(k(w)), inc(null)))
	dec(asn(mks("h"), mki(k(h)), inc(null)))
}

//go:export Ui
func Ui(t, b, x0, x1, y0, y1, mod int) { // mouse event
	obuf = take(0, 0, obuf)
	js.Global().Set("kio", "")
	var r k
	if t == 0 { // key
		r = mk(I, 2)
		m.k[2+r] = k(b)
		m.k[3+r] = k(mod)
		r = kx(mks(".k"), r)
	} else { // mouse
		r = mk(L, 4)
		m.k[2+r] = mki(k(b))
		x := mk(I, 2)
		m.k[2+x] = k(x0)
		m.k[3+x] = k(x1)
		m.k[3+r] = x
		y := mk(I, 2)
		m.k[2+y] = k(y0)
		m.k[3+y] = k(y1)
		m.k[4+r] = y
		m.k[5+r] = mki(k(mod))
		r = kx(mks(".m"), r)
	}
	if match(r, null) {
		dec(r)
		return
	}
	out(r)
	if t, n := typ(obuf); n != 0 {
		if t != C || n == atom {
			panic("type")
		}
		p := 8 + obuf<<2
		js.Global().Set("kio", s(m.c[p:p+n]))
	}
}

//go:export Store
func Store(n int) *byte { // create entry and allocate memory for dropped file
	if n <= 0 || n > maxmem/4 {
		panic("size")
	}
	s := c2s(mkb([]c(js.Global().Get("kio").String())))
	c := mk(C, k(n))
	p := 8 + c<<2
	r := key(s, c)
	f := mk(N+2, atom)
	m.k[2+f] = 12 + dyad // ,
	dec(asn(mks(".fs"), r, f))
	return &m.c[p]
}
