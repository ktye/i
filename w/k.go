package w

import "unsafe"

type c = byte
type k = uint32
type i = int32
type f = float64
type s = string

const (
	C, I, F, Z, S, G, L, D byte = 1, 2, 3, 4, 5, 6, 7, 8
)

type (
	fc1 func(c) c
	fi1 func(i) i
	ff1 func(f) f
	fz1 func(f, f) (f, f)
)

//                c  i  f   z  s  g  l  d
var lns = [9]k{0, 1, 4, 8, 16, 8, 4, 4, 8}

var m struct { // linear memory (slices share underlying arrays)
	c []c
	k []k
	f []f
}

type slice struct {
	p uintptr
	l int
	c int
}

func ini() { // start function
	m.f = make([]f, 1<<13)
	msl()
	p := k(1 << 16)
	for i := 15; i > 6; i-- {
		p >>= 1
		m.c[p] = c(i)
		m.k[i] = p
	}
	m.c[0] = 7
	m.c[4] = 16 // total memory (log2)
	println("ini mc4 mk1", m.c[4], m.k[1])
	// TODO: pointer to k-tree at 8
	m.k[9] = 0 // no free bucket 9
	a := 1 << 7
	m.k[a] = k(73) // 73: 1<<6|9 (type i, bucket 9), length is ignored
	for i := range lns {
		m.k[a+i+2] = k(lns[i])
	}
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
	s := m.k[1]
	if 1<<k(s) != len(m.c) {
		println(s, len(m.c), 1<<k(s))
	}
	m.k[s] = k(len(m.c))
	m.f = append(m.f, make([]f, len(m.f))...)
	msl()
	m.k[1] = s + 1
	m.c[1<<s] = c(s)
}
func mk(t c, n int) k { // make type t of len n (-1:atom)
	sz := lns[t]
	if n >= 0 {
		sz *= k(n)
	}
	sz += 8 // size needed including header
	bs := k(16)
	bt := 0
	for i := 4; i < 30; i++ { // calculate bucket bt from size sz (clz)
		if sz <= bs {
			bt = i
			break
		}
		bs <<= 1
	}
	if bt == 0 {
		panic("memory")
	}
	fb, a := 0, k(0)
	for i := bt; i < 30; i++ { // find next free bucket >= bt
		if m.k[5] == 1 {
			println("i", i)
			panic("wrong 0")
		}
		if k(i) >= m.k[1] {
			grw()
		}
		if m.k[i] != 0 {
			fb, a = i, m.k[i]
			break
		}
	}
	m.k[fb] = m.k[1+a>>2]           // occupy
	for i := fb - 1; i >= bt; i-- { // split large buckets
		m.k[1+a>>2] = m.k[i]
		m.k[i] = a
		m.k[a>>2] = k(i)
		a += k(1) << c(i)
		m.k[a>>2] = k(i)
	}
	if n < 0 { // set header
		m.c[int(a+1)] = t
	} else {
		//set(k(a), k(m[int(a)]|t<<5)|k(n)<<8)
		m.k[a>>2] = k(m.c[int(a)]|t<<5) | k(n)<<8
	}
	m.k[1+a>>2] = 1 // refcount
	return a
}
func typ(a k) (c, int) { // type and length at addr
	i := int(a)
	t := m.c[i] >> 5
	if t == 0 {
		return m.c[int(i+1)], -1
	}
	return t, int(m.k[k(i)>>2]) >> 8
}
func inc(x k) k { m.k[1+x>>2]++; return x }
func free(x k) {
	bt := 0x3f & m.c[x]
	b := x ^ k(1<<bt) // buddy address
	if m.c[b]&0x3f != bt {
		panic("buddy has wrong bt") // TODO rm
	}
	m.k[x>>2] = k(bt) << 24
	m.k[1+x>>2] = 0
	if m.k[b>>2] == k(bt)<<24 { // buddy is free: merge
		p := m.k[bt] // find free list parent of b
		n := p
		for n != b {
			p = n
			n = m.k[1+p>>2]
		}
		if p == b {
			m.k[bt] = m.k[1+p>>2]
		} else {
			m.k[1+p>>2] = m.k[1+b>>2] // delete from within the free list
		}
		if b < x {
			x = b
		}
		m.c[x] = bt + 1
		free(x)
	} else {
		m.k[1+x>>2] = m.k[bt]
		m.k[bt] = x
	}
}
func dec(x k) {
	rc := m.k[1+x>>2] - 1
	m.k[1+x>>2] = rc
	if rc == 0 {
		free(x)
	}
}
func to(x k, rt c) (r k) { // numeric conversions for types CIFZ
	if rt == 0 {
		return x
	}
	t, n := typ(x)
	if rt == t {
		return x
	}
	r = mk(t, n)
	if n < 0 {
		n = 1
	}
	var g func(k, k)
	switch {
	case t == C && rt == I:
		g = func(x, y k) { m.k[y>>2] = k(i(m.c[int(x)])) }
	case t == C && rt == F:
		g = func(x, y k) { m.f[y>>3] = f(m.c[int(x)]) }
	case t == I && rt == C:
		g = func(x, y k) { m.c[int(y)] = c(i(m.k[x>>2])) }
	case t == I && rt == F:
		g = func(x, y k) {
			println("I2F", i(m.k[x>>2]))
			m.f[y>>3] = f(i(m.k[x>>2]))
		}
	case t == F && rt == C:
		g = func(x, y k) { m.c[int(y)] = c(m.f[x>>3]) }
	case t == F && rt == I:
		g = func(x, y k) { m.k[y>>2] = k(i(m.f[x>>3])) }
		// TODO Z
	}
	xs, rs := lns[t], lns[rt]
	for j := k(0); j < k(n); j++ {
		a := x + 8 + j*xs
		b := r + 8 + j*rs
		g(a, b)
	}
	dec(x)
	return r
}
func cl1(x, r k, n k, op fc1) { // C vector r=f(x)
	o := r - x
	for j := 8 + x; j < 8+x+n; j++ {
		m.c[j] = op(m.c[o+j])
	}
}
func il1(x, r k, n k, op fi1) { // I vector r=f(x)
	o := (r - x) >> 2
	for j := 2 + x>>2; j < 2+n+x>>2; j++ {
		m.k[o+j] = k(op(i(m.k[j])))
	}
}
func fl1(x, r k, n k, op ff1) { // F vector r=f(x)
	o := (r - x) >> 3
	for j := 1 + x>>3; j < 1+n+x>>3; j++ {
		m.f[o+j] = op(m.f[j])
	}
}
func zl1(x, r k, n k, op fz1) { // Z vector r=f(x)
	o := (r - x) >> 3
	for j := 1 + x>>3; j < 1+2*n+x>>2; j += 2 {
		m.f[o+j], m.f[o+j+1] = op(m.f[j], m.f[j+1])
	}
}
func nm(x k, fc fc1, fi fi1, ff ff1, fz fz1, rt c) (r k) { // numeric monad
	t, n := typ(x)
	min := C
	if fc == nil {
		min = I
	}
	if fi == nil {
		min = F
	} // TODO: Z only for ff == nil ?
	if min > t { // uptype x
		println("min", min, t)
		x, t = to(x, min), min
		println("=>", x, t)
		println("@x", m.c[0], m.c[1], m.c[2], m.c[3], m.c[4], m.f[1+x>>3])
	}
	if t == Z && fz == nil { // e.g. real functions
		x, t = to(x, F), F
	}
	if m.k[1+x>>2] == 1 {
		r = inc(x)
	} else {
		r = mk(t, n)
	}
	if n < 0 {
		n = 1
	}
	switch t {
	case L:
		r >>= 2
		x >>= 2
		for j := k(0); j < k(n); j++ {
			//set(r+8+4*j, nm(inc(get(x+8+j*8)), fc, fi, ff, fz, rt))
			m.k[r+2+j] = nm(inc(m.k[j+2+x]), fc, fi, ff, fz, rt)
		}
		r <<= 2
		x <<= 2
	case D:
		if r != x {
			m.k[2+r>>2] = inc(m.k[2+x>>2])
		}
		m.k[3+r>>2] = nm(m.k[3+x>>2], fc, fi, ff, fz, rt)
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
	if r != 0 && t > rt {
		r = to(r, rt) // downtype, e.g. floor
	}
	return r
}
func kdx(x k, t c) k { // x must be a numeric scalar
	switch t {
	case I:
		return m.k[2+x>>2]
	case F, Z:
		i := int(m.f[1+x>>3])
		if i < 0 {
			panic("domain")
		}
		return k(i)
	}
	panic("type")
}
func til(x k) k { // !n
	t, n := typ(x)
	if n < 0 {
		if t > Z {
			panic("type")
		}
		n := kdx(x, t)
		r := mk(I, int(n))
		a := r >> 2
		for i := k(0); i < n; i++ {
			m.k[int(2+i+a)] = k(i)
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
