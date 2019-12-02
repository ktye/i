package main

import "syscall/js"

var obuf k

// max memory, heap size must be set explicitly when compiling with tinygo
const maxmem = 1 << 21 // number of floats (16 MB)
var mem []f

func main() {
	mem = make([]f, maxmem)
	ini(mem[:1<<13])
	table[21] = red
	table[21+dyad] = wrt
	table[39] = trp
	table[29+dyad] = drw
	js.Global().Set("kio", maxmem>>17)
	select {}
}
func grw(c k) {
	if len(m.f) >= maxmem {
		panic("ws")
	}
	m.f = mem[:2*len(m.f)]
}
func red(x k) (r k) { panic("nyi") } // 1:x read from js f[file]
func wrt(x, y k) (r k) { // x 1:y
	if m.k[y]>>28 != C {
		panic("type")
	}
	if obuf != 0 {
		dec(obuf)
	}
	obuf = y
	return x
}
func trp(x, y k) (r k) { panic("nyi") }
func drw(x, y k) (r k) { // x 9:y (draw)
	xt, yt, xn, yn := typs(x, y)
	if xt != I || yt != I || xn != atom || yn == atom {
		panic("type")
	}
	w := m.k[2+x]
	h := yn / w
	if w*h != yn {
		panic("length")
	}
	p := ptr(y, C)
	js.Global().Set("img", m.c[p:p+4*yn]) // TODO: invalid value (CopyBytesToJS is missing from tinygo)
	return decr(x, y, inc(null))
}

//go:export K
func K() {
	x := mkb([]c(js.Global().Get("kio").String()))
	js.Global().Set("kio", "")
	evp(x)
	if obuf != 0 {
		t, n := typ(obuf)
		if t != C {
			panic("type")
		}
		p := 8 + obuf << 2
		js.Global().Set("kio", s(m.c[p:p+n]))
	}
}

