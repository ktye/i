package main

import (
	"strconv" // to be removed
	"unsafe"
)

type c = byte
type k = uint32
type i = int32
type f = float64
type z = complex128
type s = string

const (
	C, I, F, Z, S, L, D, N k = 1, 2, 3, 4, 5, 6, 7, 8
	atom                   k = 0x0fffffff
	NaI                    i = -2147483648
)

type (
	fc1 func(c) c
	fi1 func(i) i
	ff1 func(f) f
	fz1 func(z) z
)
type slice struct {
	p uintptr
	l int
	c int
}

//                 C  I  F   Z  S  L  D  0  1  2  3  4
var lns = [13]k{0, 1, 4, 8, 16, 8, 4, 8, 4, 4, 4, 4, 4}
var ofs = [13]k{0, 0, 0, 0, 8, 0, 0, 0, 0, 0, 0, 0, 0}

var m struct { // linear memory (slices share underlying arrays)
	c []c
	k []k
	f []f
	z []z
}
var cpx = []func(k, k){nil, cpC, cpI, cpF, cpZ, cpF, cpL}      // copy (arguments are byte addresses)
var swx = []func(k, k){nil, swC, swI, swF, swZ, swF, swI, swI} // swap
var nax = []func(k){nil, naC, naI, naF, naZ, naS}              // set missing/nan
var eqx = []func(k, k) bool{nil, eqC, eqI, eqF, eqZ, eqS, nil} // equal
var ltx = []func(k, k) bool{nil, ltC, ltI, ltF, ltZ, ltS}      // less than
var gtx = []func(k, k) bool{nil, gtC, gtI, gtF, gtZ, gtS}      // greater than

func ini() { // start function
	m.f = make([]f, 1<<13)
	msl()
	m.k[2] = 16
	p := k(64)
	m.k[8] = p
	m.k[p] = 8
	for i := 9; i < 16; i++ {
		p *= 2
		m.k[i] = p
		m.k[p] = k(i)
	}
	m.k[0] = (I << 28) | 31
	// TODO: K tree
	// TODO: size vector
}
func msl() { // update slice header after increasing m.f
	f := *(*slice)(unsafe.Pointer(&m.f))
	i := *(*slice)(unsafe.Pointer(&m.k))
	i.l, i.c, i.p = f.l*2, f.c*2, f.p
	m.k = *(*[]k)(unsafe.Pointer(&i))
	b := *(*slice)(unsafe.Pointer(&m.c))
	b.l, b.c, b.p = f.l*8, f.c*8, f.p
	m.c = *(*[]c)(unsafe.Pointer(&b))
	zz := *(*slice)(unsafe.Pointer(&m.z))
	zz.l, zz.c, zz.p = f.l/2, f.l/2, f.p
	m.z = *(*[]z)(unsafe.Pointer(&zz))
}
func cpC(dst, src k)  { m.c[dst] = m.c[src] }
func cpI(dst, src k)  { m.k[dst>>2] = m.k[src>>2] }
func cpF(dst, src k)  { m.f[dst>>3] = m.f[src>>3] }
func cpZ(dst, src k)  { m.z[dst>>4] = m.z[src>>4] }
func cpL(dst, src k)  { inc(m.k[src>>2]); cpI(dst, src) }
func swC(dst, src k)  { m.c[dst], m.c[src] = m.c[src], m.c[dst] }
func swI(dst, src k)  { m.k[dst>>2], m.k[src>>2] = m.k[src>>2], m.k[dst>>2] }
func swF(dst, src k)  { m.f[dst>>3], m.f[src>>3] = m.f[src>>3], m.f[dst>>3] }
func swZ(dst, src k)  { m.z[dst>>4], m.z[src>>4] = m.z[src>>4], m.z[dst>>4] }
func naC(dst k)       { m.c[dst] = 32 }
func naI(dst k)       { m.k[dst>>2] = 0x80000000 }
func naF(dst k)       { u := uint64(0x7FF8000000000001); m.f[dst>>3] = *(*f)(unsafe.Pointer(&u)) }
func naZ(dst k)       { naF(dst); naF(8 + dst) }
func naS(dst k)       { mys(dst, uint64(' ')<<(56)) }
func eqC(x, y k) bool { return m.c[x] == m.c[y] }
func ltC(x, y k) bool { return m.c[x] < m.c[y] }
func gtC(x, y k) bool { return m.c[x] > m.c[y] }
func eqI(x, y k) bool { return i(m.k[x>>2]) == i(m.k[y>>2]) }
func ltI(x, y k) bool { return i(m.k[x>>2]) < i(m.k[y>>2]) }
func gtI(x, y k) bool { return i(m.k[x>>2]) > i(m.k[y>>2]) }
func eqF(x, y k) bool { return m.f[x>>3] == m.f[y>>3] }
func ltF(x, y k) bool { return m.f[x>>3] < m.f[y>>3] }
func gtF(x, y k) bool { return m.f[x>>3] > m.f[y>>3] }
func eqZ(x, y k) bool { return eqF(x, y) && eqF(8+x, 8+y) }
func ltZ(x, y k) bool { // real than imag
	if ltF(x, y) {
		return true
	} else if eqF(x, y) {
		return ltF(8+x, 8+y)
	}
	return false
}
func gtZ(x, y k) bool { // real than imag
	if gtF(x, y) {
		return true
	} else if eqF(8+x, 8+y) {
		return gtF(8+x, 8+y)
	}
	return false
}
func eqS(x, y k) bool { return m.k[x>>2] == m.k[y>>2] && m.k[1+x>>2] == m.k[1+y>>2] }
func ltS(x, y k) bool { return sym(x) < sym(y) }
func gtS(x, y k) bool { return sym(x) > sym(y) }
func mv(dst, src k) {
	t, n := typ(src)
	ln := k(1 << bk(t, n))
	rc := m.k[1+dst]
	dst, src = dst<<2, src<<2
	copy(m.c[dst:dst+ln], m.c[src:src+ln]) // copy bucket
	dst >>= 2
	m.k[dst] = t<<28 | n // restore header
	m.k[1+dst] = rc
}

func grw() { // double memory
	s := m.k[2]
	if 1<<k(s) != len(m.c) {
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
	return buk(sz + 8) // complex values have an additional 8 byte padding after the header (does not change bucket type)
}
func mk(t, n k) k { // make type t of len n (-1:atom)
	bt := bk(t, n)
	fb, a := k(0), k(0)
	for i := bt; i < 31; i++ { // find next free bucket >= bt
		if k(i) >= m.k[2] {
			panic("grow")
			grw() // TODO: run a gc cycle (merge blocks) before growing?
		}
		if m.k[i] != 0 {
			fb, a = i, m.k[i]
			break
		}
	}
	m.k[fb] = m.k[1+a] // occupy
	if p := m.k[fb]; p > 0 && p < 40 {
		panic("illegal free pointer")
	}
	for i := fb - 1; i >= bt; i-- { // split large buckets
		u := a + 1<<(i-2) // free upper half
		m.k[1+u] = m.k[i]
		m.k[i] = u
		m.k[u] = i
	}
	m.k[a] = n | t<<28 // ok for atoms
	m.k[a+1] = 1       // refcount
	// println("mk", t, n, hxk(a))
	return a
}

func typ(a k) (k, k) { // type and length at addr
	return m.k[a] >> 28, m.k[a] & atom
}
func inc(x k) k {
	t, n := typ(x)
	switch t {
	case L:
		if n == atom {
			panic("type")
		}
		for i := k(0); i < n; i++ {
			inc(m.k[2+x+i])
		}
	case D:
		if n != atom {
			panic("type")
		}
		inc(m.k[2+x])
		inc(m.k[3+x])
	}
	m.k[1+x]++
	return x
}
func use(x, t, n k) k {
	if m.k[1+x] == 1 && bk(typ(x)) == bk(t, n) {
		m.k[x] = t<<28 | n
		return x
	} else {
		return mk(t, n)
	}
}
func decret(x, r k) k {
	if r != x {
		dec(x)
	}
	return r
}
func dec(x k) {
	if m.k[x]>>28 == 0 || m.k[1+x] == 0 {
		//xxd()
		panic("unref " + hxk(x))
	}
	t, n := typ(x)
	switch t {
	case L:
		if n == atom {
			panic("type")
		}
		for i := k(0); i < n; i++ {
			dec(m.k[2+x+i])
		}
	case D:
		if n != atom {
			panic("type")
		}
		dec(m.k[2+x])
		dec(m.k[3+x])
	}
	m.k[1+x]--
	if m.k[1+x] == 0 {
		free(x)
	}
}
func free(x k) {
	// println("free", hxk(x))
	t, n := typ(x)
	bt := bk(t, n)
	m.k[x] = bt
	m.k[x+1] = m.k[bt]
	m.k[bt] = x
}
func srk(r, t, n, nn k) { // shrink bucket
	if m.k[r]>>28 != t {
		panic("type")
	}
	m.k[r] = t<<28 | nn
	big, small := bk(t, n), bk(t, nn)
	s := k(1 << (big - 3))
	for j := big; j > small; j-- {
		m.k[r+s] = I<<28 | (1 << (j - 3)) - 2
		free(r + s)
		s >>= 1
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
	case t == Z && rt == F:
		g = func(x, y k) { m.f[y] = m.f[2+x<<1] }
	case t == Z && rt == C:
		g = func(x, y k) { m.c[y] = c(i(m.f[2+x<<1])) }
	default:
		panic("nyi")
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
	o := r<<2 - x<<2
	for j := 8 + x<<2; j < 8+n+x<<2; j++ {
		m.c[o+j] = op(m.c[j])
	}
}
func il1(x, r k, n k, op fi1) { // I vector r=f(x)
	o := r - x
	for j := 2 + x; j < 2+n+x; j++ {
		m.k[o+j] = k(op(i(m.k[j])))
	}
}
func fl1(x, r k, n k, op ff1) { // F vector r=f(x)
	o := r>>1 - x>>1
	for j := 1 + x>>1; j < 1+n+x>>1; j++ {
		m.f[o+j] = op(m.f[j])
	}
}
func zl1(x, r k, n k, op fz1) { // Z vector r=f(x)
	o := r>>2 - x>>2
	for j := 1 + x>>2; j < 1+n+x>>2; j++ {
		m.z[o+j] = op(m.z[j])
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
	r = use(x, t, n)
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
	decret(x, r)
	if rt != 0 && t > rt {
		r = to(r, rt) // downtype, e.g. floor
	}
	return r
}
func idx(x, t k) i { // int from a numeric scalar (trunc, ignore imag)
	switch t {
	case C:
		return i(m.c[8+x])
	case I:
		return i(m.k[2+x])
	case F:
		return i(m.f[1+x>>1])
	case Z:
		return i(m.f[2+x])
	}
	panic("type")
}
func til(x k) (r k) { // !n
	t, n := typ(x)
	if n != atom {
		panic("nyi !a")
	} else if t == D {
		r = inc(m.k[2+x])
		dec(x)
		return r
	} else if t > Z {
		panic("type")
	}
	if nn := idx(x, t); nn < 0 {
		dec(x)
		return eye(k(-nn))
	} else {
		dec(x)
		return jota(k(nn))
	}
}
func jota(n k) (r k) {
	r = mk(I, n)
	for j := k(0); j < n; j++ {
		m.k[2+r+j] = j
	}
	return r
}
func eye(n k) (r k) {
	r = mk(L, n)
	for j := k(0); j < n; j++ {
		rj := mk(I, n)
		m.k[2+r+j] = rj
		for jj := k(0); jj < n; jj++ {
			if j == jj {
				m.k[2+rj+jj] = 1
			} else {
				m.k[2+rj+jj] = 0
			}
		}
	}
	return r
}
func flp(x k) (r k) { // +x
	t, n := typ(x)
	if t > L || n == atom { // tables are not implemented
		panic("type")
	} else if t < L {
		return x
	}
	nr, tt := k(0), L
	for j := k(0); j < n; j++ {
		tj, nj := typ(m.k[2+x+j])
		if j == 0 {
			nr = nj
		} else if nj != nr { // k7 does scalar extension for atoms, and nan-fills for short arrays
			panic("size")
		}
		if j == 0 {
			tt = tj
		} else if tt != L && tt != tj {
			tt = L
		}
	}
	if tt > L {
		panic("type")
	}
	if nr == atom {
		nr = 1
	}
	cp, sz, o := cpx[tt], lns[tt], ofs[tt]
	r = mk(L, nr)
	for j := k(0); j < nr; j++ {
		rr := mk(tt, n)
		m.k[2+r+j] = rr
	}
	if tt == L {
		for k := k(0); k < n; k++ {
			col := explode(2 + x + k)
			for i := uint32(0); i < nr; i++ {
				m.k[2+k+m.k[2+r+i]] = inc(m.k[2+col+i])
			}
			dec(col)
		}
	} else {
		for i := uint32(0); i < nr; i++ { // Rik = +Xki (cdn.mos.cms.futurecdn.net/XTZkbu7r5c4LZQ5SMzJDbV-970-80.jpg)
			for k := k(0); k < n; k++ {
				cp(8+o+sz*k+m.k[2+r+i]<<2, 8+o+sz*i+m.k[2+x+k]<<2)
			}
		}
	}
	dec(x)
	return r
}
func explode(x k) (r k) { // explode an array to a list of atoms
	t, n := typ(x)
	if t == L {
		return x
	} else if t > L {
		panic("type")
	}
	if n == atom {
		n = 1
	}
	cp, sz, o := cpx[t], lns[t], ofs[t]
	r = mk(L, n)
	for j := k(0); j < n; j++ {
		dst := mk(t, 1) << 2
		cp(8+o+dst, 8+o+x+sz*j+x<<2)
		m.k[2+r+j] = dst >> 2
	}
	dec(x)
	return r
}
func neg(x k) k { // -x
	return nm(x, func(x c) c { return -x }, func(x i) i { return -x }, func(x f) f { return -x }, func(x z) z { return -x }, 0) // TODO Z
}
func inv(x k) k { // %x
	return nm(x, nil, nil, func(x f) f { return 1.0 / x }, func(x z) z { return 1 / x }, 0)
}
func fst(x k) (r k) { // *x
	t, n := typ(x)
	if t == D {
		inc(m.k[3+x])
		r = fst(m.k[3+x])
		dec(x)
		return r
	}
	if n == atom {
		return x
	} else if n == 0 {
		panic("nyi: fst empty") // what to return? missing value? panic?
	}
	if t == L {
		r = m.k[2+x]
		inc(r)
		dec(x)
		return r
	}
	r = mk(t, atom)
	switch t {
	case C:
		m.c[8+r<<2] = m.c[8+x<<2]
	case I:
		m.k[2+r] = m.k[2+x]
	case F, S:
		m.f[1+r>>1] = m.f[1+x>>1]
	case Z:
		m.z[1+r>>2] = m.z[1+x>>2]
	default:
		panic("nyi")
	}
	dec(x)
	return r
}
func rev(x k) (r k) { // |x
	t, n := typ(x)
	if n == atom || n < 2 {
		if t == D {
			r = use(x, t, n)
			if r == x {
				m.k[r+2] = rev(m.k[x+2])
				m.k[r+3] = rev(m.k[x+3])
			} else {
				m.k[r+2] = rev(inc(m.k[x+2]))
				m.k[r+3] = rev(inc(m.k[x+3]))
			}
			return decret(x, r)
		}
		return x
	}
	r = use(x, t, n)
	if t < D {
		sz, cp, m, o := lns[t], cpx[t], n, ofs[t]
		if t == L {
			cp = cpI // copy pointer, don't inc.
		}
		src, dst := o+8+x<<2, o+8+r<<2
		if src == dst {
			cp, m = swx[t], n/2
		}
		for j := k(0); j < m; j++ {
			cp(dst+(n-1-j)*sz, src+j*sz)
		}
	} else {
		panic("nyi")
	}
	if t == L && x != r {
		for i := k(0); i < n; i++ {
			inc(m.k[2+i+r])
		}
	}
	return decret(x, r)
}
func wer(x k) (r k) { // &x
	t, n := typ(x)
	if t != I {
		panic("type")
	} else if n == atom {
		n = 1
	}
	nn := k(0)
	for j := k(0); j < n; j++ {
		if p := i(m.k[2+x+j]); p < 0 {
			panic("domain")
		} else {
			nn += k(p)
		}
	}
	r = mk(I, nn)
	jj := k(0)
	for j := k(0); j < n; j++ {
		for p := k(0); p < m.k[2+x+j]; p++ {
			m.k[2+r+jj] = j
			jj++
		}
	}
	dec(x)
	return r
}
func asc(x k) (r k) { // <x
	t, n := typ(x)
	if n == atom || t >= L { // k7 also sorts lists of different numeric types
		panic("type")
	}
	a := mk(I, atom)
	m.k[2+a] = n
	r = til(a)
	sz, lt, sw, o := lns[t], ltx[t], swI, ofs[t]
	src, ind, dst := 8+o+x<<2, 2+r, 8+r<<2
	for i := k(1); i < n; i++ { // insertion sort, should be replaced
		for j := k(i); j > 0 && lt(src+sz*m.k[ind+j], src+sz*m.k[ind+j-1]); j-- {
			sw(dst+4*j, dst+4*(j-1))
		}
	}
	dec(x)
	return r
}
func dsc(x k) (r k) { return rev(asc(x)) } // >x
func grp(x k) (r k) { // =x
	t, n := typ(x)
	if n == atom {
		panic("value")
	}
	eq, sz, o := eqx[t], lns[t], ofs[t]
	r = mk(D, atom)
	u := unq(inc(x)) // TODO: keys are sorted in k7
	l, nu := mk(L, m.k[u]&atom), m.k[u]&atom
	m.k[2+r], m.k[3+r] = u, l
	uc, kc := 8+o+u<<2, 8+o+x<<2
	b := mk(C, n) // boolean
	bc := 8 + b<<2
	for j := k(0); j < nu; j++ { // over ?x
		nr := k(0)
		for jj := k(0); jj < n; jj++ { // over x
			m.c[bc+jj] = 0
			if eq(uc+sz*j, kc+sz*jj) {
				m.c[bc+jj] = 1
				nr++
			}
		}
		lj, p := mk(I, nr), k(0)
		for jj := k(0); jj < n; jj++ { // over x
			if m.c[bc+jj] == 1 {
				m.k[2+lj+p] = jj
				p++
			}
		}
		m.k[2+l+j] = lj
	}
	dec(b)
	dec(x)
	return r
}
func not(x k) (r k) { // ~x
	return nm(x, func(x c) (r c) {
		if x == 0 {
			r = 1
		}
		return r
	}, func(x i) (r i) {
		if x == 0 {
			r = 1
		}
		return r
	}, func(x f) (r f) {
		if x == 0 {
			r = 1
		}
		return r
	}, func(x z) (r z) {
		if x == 0 {
			r = 1
		}
		return r
	}, C)
}
func enl(x k) (r k) { // ,x
	t, n := typ(x)
	if t < L && n == atom {
		r = use(x, t, 1)
		if r == x {
			m.k[r] = t<<28 | 1
			return r
		}
		cp, o := cpx[t], ofs[t]
		src, dst := o+8+x<<2, o+8+r<<2
		cp(dst, src)
		dec(x)
		return r
	}
	r = mk(L, 1)
	m.k[2+r] = x
	return r
}
func srt(x k) (r k) { return atx(x, asc(inc(x))) } // ^x  TODO: replace with a sort implementation
func cnt(x k) (r k) { // #x
	t, n := typ(x)
	r = mk(I, atom)
	if t == D {
		_, n = typ(m.k[x+2])
	} else if n == atom {
		n = 1
	}
	m.k[2+r] = k(i(n))
	dec(x)
	return r
}
func flr(x k) k { // _x
	return nm(x, func(x c) c { return x }, func(x i) i { return x }, func(x f) f {
		y := float64(int32(x))
		if x < y {
			y -= 1.0
		}
		return y
	}, nil, I)
}
func fms(x k) (r k) { // $x
	lim := k(120)
	r = mk(C, lim)
	t, n := typ(x)
	dd := false
	pcnt := k(0)
	push := func(s s) bool {
		ln := k(len(s))
		if pcnt+ln > lim {
			dd, ln = true, lim-pcnt
		}
		for j := k(0); j < ln; j++ {
			m.c[8+pcnt+r<<2] = c(s[j])
			pcnt++
		}
		return dd
	}
	pushr := func(p k) (o bool) {
		a := fms(inc(p))
		as, ln := a<<2, m.k[a]&atom
		if at, an := typ(a); at != C || an > lim {
			panic("type")
		}
		if push(s(m.c[8+as : 8+as+ln])) {
			o = true
		}
		dec(a)
		return o
	}
	vecs := func(fs func(jj k) s) {
		if n == 0 {
			push("[-]") // TODO: format empty arrays?
		} else if n == atom {
			n = 1
		}
		sep := ""
		for j := k(0); j < n; j++ {
			if push(sep + fs(j)) {
				break
			}
			sep = " "
		}
	}
	if n == 1 {
		push(",")
	}
	switch t {
	case C:
		nn, xc, o := n, x<<2, true
		if nn == atom {
			nn = 1
		} else if nn > lim {
			nn = lim
		}
		ispr := func(b c) bool { return b >= 0x20 && b < 0x7E }
		for j := k(0); j < nn; j++ {
			if !ispr(m.c[8+xc+j]) {
				o = false
				break
			}
		}
		if o {
			push(`"`)
			push(s(m.c[8+xc : 8+xc+nn]))
			push(`"`)
		} else {
			push("0x")
			for j := k(0); j < nn; j++ {
				c1, c2 := hxb(m.c[8+xc+j])
				if push(s([]c{c1, c2})) {
					break
				}
			}
		}
	case I: // TODO no strconv
		vecs(func(j k) s { return strconv.Itoa(int(m.k[2+x+j])) })
	case F: // TODO no strconv
		vecs(func(j k) s { return strconv.FormatFloat(m.f[1+j+x>>1], 'g', -1, 64) })
		o, af := 8+r<<2, true
		for j := k(0); j < pcnt; j++ {
			if m.c[o+j] == '.' {
				af = false
				break
			}
		}
		if af {
			push("f")
		}
	case Z: // TODO no strconv
		vecs(func(j k) s {
			return strconv.FormatFloat(m.f[2+2*j+x>>1], 'g', -1, 64) + "i" + strconv.FormatFloat(m.f[3+2*j+x>>1], 'g', -1, 64)
		})
	case S:
		if n == atom {
			n = 1
		}
		for j := k(0); j < n; j++ {
			if push("`" + str(sym(8+8*j+x<<2))) {
				break
			}
		}
	case L:
		if n != 1 {
			push("(")
		}
		for j := k(0); j < n; j++ {
			if j > 0 {
				push(";")
			}
			if pushr(m.k[2+x+j]) {
				break
			}
		}
		if n != 1 {
			push(")")
		}
	case D:
		push("(")
		pushr(m.k[2+x])
		push("!")
		pushr(m.k[3+x])
		push(")")
	default:
		println("fms t=", t)
		panic("nyi")
	}
	if dd {
		m.c[8+r+lim-1], m.c[8+r+lim-2] = '.', '.'
	}
	m.k[r] = C<<28 | pcnt
	if nn := m.k[r] & atom; nn < (8+lim)>>1 {
		srk(r, C, lim, nn)
	}
	dec(x)
	return r
}
func sym(x k) uint64 { return *(*uint64)(unsafe.Pointer(&m.c[x])) }
func mys(x k, u uint64) {
	var b [8]c
	b = *(*[8]c)(unsafe.Pointer(&u))
	copy(m.c[x:x+8], b[:])
}
func str(u uint64) s {
	var b [8]c
	n := 8
	for i := k(0); i < 8; i++ {
		j := 8 * (7 - i)
		b[i] = c((u & (0xFF << j)) >> j)
		if b[i] == 0 {
			n = int(i)
			break
		}
	}
	return s(b[:n])
}
func unq(x k) (r k) { // ?x
	t, n := typ(x)
	if n == atom {
		panic("nyi") // overloads, random numbers?
	} else if t == D { // what does ?d do?
		panic("type")
	} else if n < 2 {
		return x
	}
	r = mk(t, n)
	eq, cp, o := eqx[t], cpx[t], ofs[t]
	if t == L {
		eq = func(x, y k) bool { return match(m.k[x>>2], m.k[y>>2]) }
	}
	sz := lns[t]
	src, dst := o+8+x<<2, o+8+r<<2
	nn := k(0)
	for i := k(0); i < n; i++ { // quadratic, should be improved
		u := true
		srci := src + i*sz
		for j := k(0); j < nn; j++ {
			if eq(srci, dst+j*sz) {
				u = false
				break
			}
		}
		if u {
			cp(dst+nn*sz, srci)
			nn++
		}
	}
	srk(r, t, n, nn)
	dec(x)
	return r
}
func tip(x k) (r k) { // @x
	r = mk(S, atom)
	m.k[2+r] = 0
	m.k[3+r] = 0
	t, n := typ(x)
	if t == L {
		dec(x)
		return r // empty symbol
	}
	tns := "_cifzn a" // TODO k7 compatibility, function types 1..4?
	s := tns[t]
	if n != atom {
		s -= 32
	}
	mys(8+r<<2, uint64(s)<<56)
	dec(x)
	return r
}
func evl(x k) (r k) { // .x
	t, n := typ(x)
	if t != L {
		return x
	}
	if n == 0 {
		panic("evl empty list?") // what TODO?
	}
	v := m.k[2+x]
	vt, vn := typ(v)
	switch vt {
	case S:
		if vn != atom || n != 2 {
			panic("nyi")
		}
		c := c(sym(8+v<<2) >> 56) // TODO: this is only 1 char
		s := "+-%*|&<>=~,^#_$?@."
		h := []func(k) k{flp, neg, inv, fst, rev, wer, asc, dsc, grp, not, enl, srt, cnt, flr, fms, unq, tip, evl}
		var g func(k) k
		for i := range s {
			if c == s[i] {
				g = h[i]
			}
		}
		if g == nil {
			panic("nyi")
		}
		a := inc(m.k[3+x])
		at, _ := typ(a)
		if at == L {
			a = evl(a)
		}
		r = g(a)
		dec(x)
		return r
	}
	panic("nyi")
	return x
}

func atx(x, y k) (r k) { // x@y
	xt, xn := typ(x)
	yt, yn := typ(y)
	switch {
	case xt < L && yt == I:
		cp, sz, na, o := cpx[xt], lns[xt], nax[xt], ofs[xt]
		xc, idx := 8+o+x<<2, 2+y
		r = mk(xt, yn)
		if xn == atom {
			xn = 1
		}
		if yn == atom {
			yn = 1
		}
		rc := 8 + o + r<<2
		for j := k(0); j < yn; j++ {
			if j >= xn {
				na(rc + sz*j)
			} else {
				cp(rc+sz*j, xc+sz*m.k[idx+j])
			}
		}
		dec(x)
		dec(y)
		return r
	// case xt == L:
	//	missing element for a list is nax[type of first element]
	default:
		panic("nyi")
	}
}

func match(x, y k) (rv bool) { // recursive match
	if x == y {
		return true
	}
	t, n := typ(x)
	tt, nn := typ(y)
	if tt != t || nn != n {
		return false
	}
	if n == atom {
		n = 1
	}
	switch t {
	case L:
		for j := k(0); j < n; j++ {
			if match(m.k[2+x+j], m.k[2+y+j]) == false {
				return false
			}
		}
		return true
	case D:
		if match(m.k[2+x], m.k[2+y]) == false || match(m.k[3+x], m.k[3+y]) == false {
			return false
		}
		return true
	default:
		eq, sz, o := eqx[t], lns[t], ofs[t]
		if eq == nil {
			panic("type")
		}
		x, y = 8+o+x<<2, 8+o+y<<2
		for j := k(0); j < n; j++ {
			if eq(x+j*sz, y+j*sz) == false {
				return false
			}
		}
		return true
	}
	return false
}

func hxb(x c) (c, c) { return hexs[x>>4], hexs[x&0x0F] }
func hxk(x k) s {
	b := []c{'0', 'x', '0', '0', '0', '0', '0', '0', '0', '0'}
	for j := k(0); j < 4; j++ {
		n := 8 * (3 - j)
		b[2+2*j], b[3+2*j] = hxb(c((x & (0xFF << n)) >> n))
	}
	return s(b)
}

func isnan(x f) bool { return x != x }

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

const hexs = "0123456780abcdef"
