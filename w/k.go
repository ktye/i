package w

import "unsafe"

type c = byte
type k = uint32
type i = int32
type f = float64
type s = string

const (
	C, I, F, Z, S, G, L, D k = 1, 2, 3, 4, 5, 6, 7, 8
	atom                   k = 0x0fffffff
)

type (
	fc1 func(c) c
	fi1 func(i) i
	ff1 func(f) f
	fz1 func(f, f) (f, f)
)
type slice struct {
	p uintptr
	l int
	c int
}

//                C  I  F   Z  S  G  L  D
var lns = [9]k{0, 1, 4, 8, 16, 8, 4, 4, 8}
var m struct { // linear memory (slices share underlying arrays)
	c []c
	k []k
	f []f
}

func ini() { // start function
	m.f = make([]f, 1<<13)
	msl()
	m.k[2] = 16
	p := k(256)
	m.k[8] = p
	for i := 9; i < 16; i++ {
		p *= 2
		m.k[i] = p >> 2
		m.k[p>>2] = k(i)
	}
	m.k[0] = (I << 28) | 31
	// TODO: K tree
	// TODO: size vector
}
func msl() { // update slice header after increasing m.f
	f := *(*slice)(unsafe.Pointer(&m.f))
	i := *(*slice)(unsafe.Pointer(&m.k))
	i.l = f.l * 2
	i.c = f.c * 2
	i.p = f.p
	m.k = *(*[]k)(unsafe.Pointer(&i))
	b := *(*slice)(unsafe.Pointer(&m.c))
	b.l = f.l * 8
	b.c = f.c * 8
	b.p = f.p
	m.c = *(*[]c)(unsafe.Pointer(&b))
}
func grw() { // double memory
	s := m.k[2]
	if 1<<k(s) != len(m.c) {
		println("grow", len(m.c), 1<<k(s))
		panic("grow")
	}
	m.k[s] = k(len(m.c)) >> 2
	m.f = append(m.f, make([]f, len(m.f))...)
	msl()
	m.k[2] = s + 1
	m.k[1<<(s-2)] = s // bucket type of new upper half
}
func bk(t, n k) k {
	sz := lns[t]
	if n != atom {
		sz *= n
	}
	if sz > 1<<31 {
		panic("size")
	}
	return buk(sz + 8)
}
func mk(t, n k) k { // make type t of len n (-1:atom)
	bt := bk(t, n)
	fb, a := k(0), k(0)
	for i := bt; i < 31; i++ { // find next free bucket >= bt
		if k(i) >= m.k[2] {
			grw()
		}
		if m.k[i] != 0 {
			fb, a = i, m.k[i]
			break
		}
	}
	m.k[fb] = m.k[1+a]              // occupy
	for i := fb - 1; i >= bt; i-- { // split large buckets
		u := a + 1<<(i-2) // free upper half
		m.k[1+u] = m.k[i]
		m.k[i] = u
	}
	m.k[a] = n | t<<28 // ok for atoms
	m.k[a+1] = 1       // refcount
	return a
}
func typ(a k) (k, k) { // type and length at addr
	return m.k[a] >> 28, m.k[a] & 0x0fffffff
}
func inc(x k) k { m.k[1+x]++; return x }
func free(x k) {
	t, n := typ(x)
	bt := bk(t, n)
	m.k[x] = bt
	m.k[x+1] = m.k[bt]
	m.k[bt] = x
}
func dec(x k) {
	m.k[1+x]--
	if m.k[1+x] == 0 {
		free(x)
	}
}
func to(x, rt k) (r k) { // numeric conversions for types CIFZ
	if rt == 0 {
		return x
	}
	t, n := typ(x)
	if rt == t {
		return x
	}
	r = mk(rt, n)
	if n == atom {
		n = 1
	}
	var g func(k, k)
	switch {
	case t == C && rt == I:
		g = func(x, y k) { m.k[y] = k(i(m.c[x])) }
	case t == C && rt == F:
		g = func(x, y k) { m.f[y] = f(m.c[x]) }
	case t == I && rt == C:
		g = func(x, y k) { m.c[y] = c(i(m.k[x])) }
	case t == I && rt == F:
		g = func(x, y k) {
			m.f[y] = f(i(m.k[x]))
		}
	case t == F && rt == C:
		g = func(x, y k) { m.c[y] = c(m.f[x]) }
	case t == F && rt == I:
		g = func(x, y k) { m.k[y] = k(i(m.f[x])) }
		// TODO Z
	}
	xs, rs := lns[t], lns[rt]
	as, bs := l8t[xs-1], l8t[rs-1]
	x <<= 2
	r <<= 2
	for j := k(0); j < k(n); j++ {
		a := (x + 8 + j*xs) >> as
		b := (r + 8 + j*rs) >> bs
		g(a, b)
	}
	dec(x >> 2)
	return r >> 2
}
func cl1(x, r k, n k, op fc1) { // C vector r=f(x)
	o := (r - x) << 2
	for j := 8 + x<<2; j < 8+n+x<<2; j++ {
		m.c[j] = op(m.c[o+j])
	}
}
func il1(x, r k, n k, op fi1) { // I vector r=f(x)
	o := r - x
	for j := 2 + x; j < 2+n+x; j++ {
		m.k[o+j] = k(op(i(m.k[j])))
	}
}
func fl1(x, r k, n k, op ff1) { // F vector r=f(x)
	o := (r - x) >> 1
	for j := 1 + x>>1; j < 1+n+x>>1; j++ {
		m.f[o+j] = op(m.f[j])
	}
}
func zl1(x, r k, n k, op fz1) { // Z vector r=f(x)
	o := (r - x) >> 1
	for j := 1 + x>>1; j < 1+2*n+x>>1; j += 2 {
		m.f[o+j], m.f[o+j+1] = op(m.f[j], m.f[j+1])
	}
}
func nm(x k, fc fc1, fi fi1, ff ff1, fz fz1, rt k) (r k) { // numeric monad
	t, n := typ(x)
	min := C
	if fc == nil {
		min = I
	}
	if fi == nil {
		min = F
	} // TODO: Z only for ff == nil ?
	if min > t { // uptype x
		x, t = to(x, min), min
	}
	if t == Z && fz == nil { // e.g. real functions
		x, t = to(x, F), F
	}
	if m.k[1+x] == 1 {
		r = inc(x) // reuse x
	} else {
		r = mk(t, n)
	}
	if n == atom {
		n = 1
	}
	switch t {
	case L:
		for j := k(0); j < k(n); j++ {
			m.k[r+2+j] = nm(inc(m.k[j+2+x]), fc, fi, ff, fz, rt)
		}
	case D:
		if r != x {
			m.k[2+r] = inc(m.k[2+x])
		}
		m.k[3+r] = nm(m.k[3+x], fc, fi, ff, fz, rt)
	case C:
		cl1(x, r, k(n), fc)
	case I:
		il1(x, r, k(n), fi)
	case F:
		fl1(x, r, k(n), ff)
	case Z:
		zl1(x, r, k(n), fz)
	default:
		panic("type")
	}
	dec(x)
	if rt != 0 && t > rt {
		r = to(r, rt) // downtype, e.g. floor
	}
	return r
}
func kdx(x, t k) k { // unsigned int from a numeric scalar
	switch t {
	case I:
		return m.k[2+x]
	case F, Z:
		i := int(m.f[1+x>>1]) // trunc and ignore imag part
		if i < 0 {
			panic("domain")
		}
		return k(i)
	}
	panic("type")
}
func til(x k) k { // !n
	t, n := typ(x)
	if n == atom {
		if t > Z {
			panic("type")
		}
		n := kdx(x, t) // TODO: handle negative
		r := mk(I, n)
		for i := k(0); i < n; i++ {
			m.k[2+i+r] = i
		}
		return r
	} else {
		panic("nyi !a")
	}
	return 0
}
func neg(x k) k {
	return nm(x, func(x c) c { return -x }, func(x i) i { return -x }, func(x f) f { return -x }, func(x, y f) (f, f) { return -x, -y }, 0) // TODO Z
}
func inv(x k) k { return nm(x, nil, nil, func(x f) f { return 1.0 / x }, nil, 0) } // TODO Z
func flr(x k) k {
	return nm(x, func(x c) c { return x }, func(x i) i { return x }, func(x f) f {
		y := float64(int32(x))
		if x < y {
			y -= 1.0
		}
		return y
	}, nil, I)
}

func buk(x uint32) (n k) { // from https://golang.org/src/math/bits/bits.go (Len32)
	x--
	if x >= 1<<16 {
		x >>= 16
		n = 16
	}
	if x >= 1<<8 {
		x >>= 8
		n += 8
	}
	n += k(l8t[x])
	if n < 4 {
		return 4
	}
	return n
}

var l8t = [256]c{
	0x00, 0x01, 0x02, 0x02, 0x03, 0x03, 0x03, 0x03, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06,
	0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07,
	0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
	0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
}
