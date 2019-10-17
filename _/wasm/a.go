package main

import "syscall/js"

// max memory, heap size must be set explicitly when compiling with tinygo
const maxmem = 1 << 21 // number of floats (16 MB)

var mem []float64

func main() {
	mem = make([]float64, maxmem)
	ini(mem[:1<<13])
	table[21] = red
}
func grw() {
	if len(m.f) >= maxmem {
		panic("mem")
	}
	m.f = mem[:2*len(m.f)]
}

//go:export ex
func ex() {
	s := js.Global().Get("input").String()
	dec(evl(prs(mkb([]c(s)))))
	/*
		if m.k[r]>>28 == N {
			dec(r)
			return
		}
		r = kst(r)
		n, p := m.k[r]&atom, 8+r<<2
		s = string(m.c[p : p+n])
		dec(r)
	*/
}

//go:export ui
func ui(a, b, c, d, e, f, g, h int) {
	println("ui:", a, b, c, d, e, f, g, h)
	var r k
	switch a {
	case 0: // key, shift, alt, cntrl
		r = cal(lup(mks("uk")), l2(mki(k(b)), modifiers(c, d, e)))
	// case 1: // TODO mouse
	case 2: // width, height
		println("call us")
		r = cal2(lup(mks("us")), mki(k(b)), mki(k(c)))
	default:
		return
	}
	t, n := typ(r)
	if t != I || n == atom || n == 0 {
		dec(r)
		return
	}
	for i := k(0); i < n; i++ {
		m.k[2+i+r] |= 0xFF000000
	}
	p := 8 + r<<2
	js.CopyBytesToJS(js.Global().Get("screen"), m.c[p:p+4*n])
	js.Global().Get("draw").Invoke()
	dec(r)
}
func modifiers(a, b, c int) (r k) {
	r = mk(I, 3)
	m.k[2+r] = k(a)
	m.k[3+r] = k(b)
	m.k[4+r] = k(c)
	return r
}

func red(x k) (r k) { // 1:x read from js variable
	t, n := typ(x)
	if t == S && n == atom {
		x = str(x)
		t, n = typ(x)
	}
	if t != C || n == atom {
		panic("type")
	}
	p := 8 + x<<2
	s := string(m.c[p : p+n])
	b := js.Global().Get("f").Get(s).String()
	return decr(x, mkb([]c(b)))
}
