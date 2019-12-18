package main

// removal of syscall/js -23k (disabling math would -7.7k)

// max memory, heap size must be set explicitly when compiling with tinygo
const maxmem = 1 << 23 // number of floats (64 MB)
var mem, bak []f
var ibuf k // js to k (input char buffer)
var obuf k // k to js
var imgp k // pointer to image data (k index)
var imgs k // image size width<<16|height
var file k // file pointer for Get

func main() {
	mem = make([]f, maxmem>>1)
	bak = make([]f, maxmem>>1)
	ini(mem[:1<<13])
	save(bak)
	table[21] = red
	table[21+dyad] = wrt
	table[39] = trp
	table[29+dyad] = drw
	mkk(".k", "{(`key;x)}")
	mkk(".m", "{(`mouse;x)}")
	dec(asn(mks(".fs"), key(mk(S, 0), mk(L, 0)), inc(null)))
	dec(asn(mks(".size"), mki(maxmem>>17), inc(null))) // banner MB
	obuf = mk(S, 0)
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
	fs := lup(mks(".fs"))
	ky := m.k[2+fs]
	ix := fnd(inc(ky), x)
	if i := m.k[2+ix]; i < m.k[ky]&atom {
		return decr(ix, fs, inc(m.k[2+i+m.k[3+fs]]))
	}
	return decr(ix, fs, mk(C, 0))
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
	imgs = w<<16 | h
	return decr(x, y, inc(null))
}
func jr() { // js reset output
	imgs = 0
}

//go:export E
func E() int { // execute k string via js variable kio
	imgs = 0
	evp(ibuf)
	n := m.k[obuf] & atom
	return int(n)
}

//go:export In
func In(n int) *byte { // request new input buffer
	ibuf = mk(C, k(n))
	return &m.c[8+ibuf<<2]
}

//go:export Out
func Out() *byte {
	return &m.c[8+obuf<<2]
}

//go:export P
func P() *byte { return &m.c[0] } // k-memory offset in wasm memory buffer

//go:export Img
func Img() *byte { return &m.c[imgp] } // pointer to current image data

//go:export Isz
func Isz() uint32 { return imgs } // w<<16|h

//go:export Srcp
func Srcp() int { return int(m.k[srcp]) } // source pointer (error indicator)

//go:export Ui
func Ui(t, b, x0, x1, y0, y1, mod int) int { // mouse event
	imgs = 0
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
		return 0
	}
	out(r)
	r = m.k[obuf] & atom
	return int(r)
}

//go:export Store
func Store(n int) *byte { // create entry and allocate memory for dropped file
	imgs = 0
	if n <= 0 || n > maxmem/4 {
		panic("size")
	}
	s := c2s(ibuf) // file name
	c := mk(C, k(n))
	p := 8 + c<<2
	r := key(s, c)
	f := mk(N+2, atom)
	m.k[2+f] = 12 + dyad // ,
	dec(asn(mks(".fs"), r, f))
	return &m.c[p]
}

//go:export Get
func Get() int { // get file (length)
	imgs = 0
	s := c2s(ibuf)
	kws := mks("k.ws")
	dec(kws)
	if m.k[2+s] == m.k[2+kws] {
		dec(s)
		file = 0
		return 8 * maxmem >> 1
	}
	file = red(s)
	t, n := typ(file)
	if t != C || n == 0 || n == atom {
		dec(file)
		return 0
	}
	return int(m.k[file] & atom)
}

//go:export File
func File() *byte {
	if file == 0 {
		return &m.c[0]
	}
	dec(file)
	return &m.c[8+file<<2]
}

//go:export Save
func Save() {
	if m.k[obuf]&atom != 0 {
		obuf = take(0, 0, obuf)
	}
	save(bak)
}

//go:export Swap
func Swap() { swap(bak) }
