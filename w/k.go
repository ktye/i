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
	atom, kkey, kval, dyad k = 0x0fffffff, 0x30, 0x31, 50
	NaI                    i = -2147483648
)

type (
	f1 func(k, k)
	f2 func(k, k, k)
	fc func(k, k) bool
)
type slice struct {
	p uintptr
	l int
	c int
}

//                 C  I  F   Z  S  L  D  0  1  2  3  4
var lns = [13]k{0, 1, 4, 8, 16, 8, 4, 8, 4, 4, 4, 4, 4}
var m struct { // linear memory (slices share underlying arrays)
	c []c
	k []k
	f []f
	z []z
}

var cpx = []f1{nil, cpC, cpI, cpF, cpZ, cpF, cpL}      // copy
var swx = []f1{nil, swC, swI, swF, swZ, swF, swI, swI} // swap
var nax = []func(k){nil, naC, naI, naF, naZ, naS}      // set missing/nan
var eqx = []fc{nil, eqC, eqI, eqF, eqZ, eqS, nil}      // equal
var ltx = []fc{nil, ltC, ltI, ltF, ltZ, ltS}           // less than
var gtx = []fc{nil, gtC, gtI, gtF, gtZ, gtS}           // greater than
var stx = []func(k, k) k{nil, nil, stI, stF, stZ, stS} // tostring (assumes 56 bytes space at dst)
var tox = []f1{nil, func(r, x k) { m.k[r] = k(i(m.c[x])) }, func(r, x k) { m.f[r] = f(m.c[x]) }, func(r, x k) { m.z[r] = complex(f(m.c[x]), 0) }, func(r, x k) { m.c[r] = c(m.k[x]) }, nil, func(r, x k) { m.f[r] = f(i(m.k[x])) }, func(r, x k) { m.z[r] = complex(f(m.k[x]), 0) }, func(r, x k) { m.c[r] = c(m.f[x]) }, func(r, x k) { m.k[r] = k(i(f(m.f[x]))) }, nil, func(r, x k) { m.z[r] = complex(m.f[x], 0) }, func(r, x k) { m.c[r] = c(m.f[x<<1]) }, func(r, x k) { m.k[r] = k(i(m.f[x<<1])) }, func(r, x k) { m.f[r] = m.f[x<<1] }}

func cpC(dst, src k)  { m.c[dst] = m.c[src] }
func cpI(dst, src k)  { m.k[dst] = m.k[src] }
func cpF(dst, src k)  { m.f[dst] = m.f[src] }
func cpZ(dst, src k)  { m.z[dst] = m.z[src] }
func cpL(dst, src k)  { inc(m.k[src]); cpI(dst, src) }
func swC(dst, src k)  { m.c[dst], m.c[src] = m.c[src], m.c[dst] }
func swI(dst, src k)  { m.k[dst], m.k[src] = m.k[src], m.k[dst] }
func swF(dst, src k)  { m.f[dst], m.f[src] = m.f[src], m.f[dst] }
func swZ(dst, src k)  { m.z[dst], m.z[src] = m.z[src], m.z[dst] }
func naC(dst k)       { m.c[dst] = 32 }
func naI(dst k)       { m.k[dst] = 0x80000000 }
func naF(dst k)       { u := uint64(0x7FF8000000000001); m.f[dst] = *(*f)(unsafe.Pointer(&u)) }
func naZ(dst k)       { naF(dst << 1); naF(1 + dst<<1) }
func naS(dst k)       { mys(dst<<2, uint64(' ')<<(56)) }
func eqC(x, y k) bool { return m.c[x] == m.c[y] }
func eqI(x, y k) bool { return i(m.k[x]) == i(m.k[y]) }
func eqF(x, y k) bool { return m.f[x] == m.f[y] }
func eqZ(x, y k) bool { return eqF(x<<1, y<<1) && eqF(1+x<<1, 1+y<<1) }
func eqS(x, y k) bool { return m.k[x<<1] == m.k[y<<1] && m.k[1+x<<1] == m.k[1+y<<1] }
func ltC(x, y k) bool { return m.c[x] < m.c[y] }
func gtC(x, y k) bool { return m.c[x] > m.c[y] }
func ltI(x, y k) bool { return i(m.k[x]) < i(m.k[y]) }
func gtI(x, y k) bool { return i(m.k[x]) > i(m.k[y]) }
func ltF(x, y k) bool { return m.f[x] < m.f[y] }
func gtF(x, y k) bool { return m.f[x] > m.f[y] }
func ltZ(x, y k) bool { // real than imag
	if ltF(x<<1, y<<1) {
		return true
	} else if eqF(x<<1, y<<1) {
		return ltF(1+x<<1, 1+y<<1)
	}
	return false
}
func gtZ(x, y k) bool { // real than imag
	if gtF(x<<1, y<<1) {
		return true
	} else if eqF(x<<1, y<<1) {
		return gtF(1+x<<1, 1+y<<1)
	}
	return false
}
func ltS(x, y k) bool { return sym(x<<3) < sym(y<<3) }
func gtS(x, y k) bool { return sym(x<<3) > sym(y<<3) }
func stI(dst, x k) k { // TODO remove strconv
	if m.k[x] == 0x80000000 {
		m.c[dst] = '0'
		m.c[dst+1] = 'N'
		return 2
	}
	s := strconv.Itoa(int(i(m.k[x])))
	n := k(len(s))
	copy(m.c[dst:dst+n], []byte(s))
	return n
}
func stF(dst, x k) k { // TODO remove strconv
	s := strconv.FormatFloat(m.f[x], 'g', 6, 64)
	n := k(len(s))
	copy(m.c[dst:dst+n], []byte(s))
	return n
}
func stZ(dst, x k) k {
	n := stF(dst, x<<1)
	m.c[dst+n] = 'i'
	return 1 + n + stF(dst+1+n, 1+x<<1)
}
func stS(dst, x k) k {
	u := sym(x << 3)
	for i := k(0); i < 8; i++ {
		if c := c(u >> (8 * (7 - i))); c == 0 {
			return i
		} else {
			m.c[dst+i] = c
		}
	}
	return 8
}
func ptr(x, t k) k { // convert k address to type dependend index of data section
	switch t {
	case C:
		return (2 + x) << 2
	case I, L:
		return 2 + x
	case F, S:
		return (2 + x) >> 1
	case Z:
		return (4 + x) >> 2
	}
	panic("type")
}
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
	copy(m.c[136:169], []c(`:+-*%&|<>=!~,^#_$?@.0123456789'/\`))
	copy(m.c[169:181], []c{0, 'c', 'i', 'f', 'z', 'n', '.', 'a', 0, '1', '2', '3', '4'})
	m.k[kkey] = mk(S, 0) // k-tree keys
	m.k[kval] = mk(L, 0) // k-tree values
	m.k[3] = mk(D, atom)
	m.k[2+m.k[3]] = mk(S, 0)
	m.k[3+m.k[3]] = mk(C, 0)
	o := c(39) // monads
	builtin(o+0, "lsv")
	builtin(o+1, "clv")
	builtin(o+5, "help")
	o += c(dyad) // dyads
	builtin(o+0, "in")
	builtin(o+1, "within")
	builtin(o+2, "bin")
	builtin(o+3, "like")
	builtin(o+4, "del")
}
func builtin(code c, s string) {
	key, val := mk(S, atom), mk(C, 1)
	mys(8+key<<2, btou([]byte(s)))
	m.c[8+val<<2] = code
	m.k[2+m.k[3]] = cat(m.k[2+m.k[3]], key)
	m.k[3+m.k[3]] = cat(m.k[3+m.k[3]], val)
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
	return a
}
func mki(i k) (r k)      { r = mk(I, atom); m.k[2+r] = i; return r }
func mkc(c c) (r k)      { r = mk(C, atom); m.c[8+r<<2] = c; return r }
func mks(s s) (r k)      { return mku(btou([]c(s))) }
func mku(u uint64) (r k) { r = mk(S, atom); mys(8+r<<2, u); return r }
func typ(x k) (k, k) { // type and length at addr
	return m.k[x] >> 28, m.k[x] & atom
}
func typs(x, y k) (xt, yt, xn, yn k) { xt, xn = typ(x); yt, yn = typ(y); return }
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
	case N + 1, N + 2, N + 3, N + 4:
		if n == atom && m.k[2+x] > 255 { // derived
			inc(m.k[3+x])
		} else if n == 0 { // lambda
			inc(m.k[2+x])
			inc(m.k[3+x])
		} else if n != atom { // projection
			if n == 1 { // lambda-projection
				inc(m.k[2+x])
			}
			inc(m.k[3+x])
		}
	}
	m.k[1+x]++
	return x
}
func decr(x, r k) k     { dec(x); return r }
func decr2(x, y, r k) k { dec(x); dec(y); return r }
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
	case N + 1, N + 2, N + 3, N + 4:
		if n == atom && m.k[2+x] > 255 { // derived
			dec(m.k[3+x])
		} else if n == 0 { // lambda
			dec(m.k[2+x])
			dec(m.k[3+x])
		} else if n != atom { // projection n: 1 or 2
			if n == 1 { // lambda-projection
				dec(m.k[2+x])
			}
			dec(m.k[3+x])
		}
	}
	m.k[1+x]--
	if m.k[1+x] == 0 {
		free(x)
	}
}
func free(x k) {
	t, n := typ(x)
	bt := bk(t, n)
	m.k[x] = bt
	m.k[x+1] = m.k[bt]
	m.k[bt] = x
}
func srk(x, t, n, nn k) (r k) { // shrink bucket
	if m.k[x]>>28 != t {
		panic("type")
	}
	if bk(t, nn) < bk(t, n) { // alloc not split: prevent small object accumulation
		r = mk(t, nn)
		ln := nn * lns[t]
		if t == Z {
			ln += 8
		}
		rc, xc := 8+r<<2, 8+x<<2
		copy(m.c[rc:rc+ln], m.c[xc:xc+ln])
		return decr(x, r)
	}
	m.k[x] = t<<28 | nn
	return x
}
func to(x, rt k) (r k) { // numeric conversions for types CIFZ
	if rt == 0 || rt >= L {
		return x
	}
	t, n := typ(x)
	if rt == t {
		return x
	}
	var g func(k, k)
	if t == S && rt == I { // for symbol comparison to bool
		g = func(r, x k) {
			if m.f[x] == 0 {
				m.k[r] = 0
			} else {
				m.k[r] = 1
			}
		}
	} else {
		g = tox[4*(t-1)+rt-1]
	}
	r = mk(rt, n)
	n = atm1(n)
	xp, rp := ptr(x, t), ptr(r, rt)
	for i := k(0); i < k(n); i++ {
		g(rp+i, xp+i)
	}
	return decr(x, r)
}
func nm(x, rt k, fx []f1) (r k) { // numeric monad
	t, n := typ(x)
	min := C
	if fx[C] == nil {
		min = I
	}
	if fx[I] == nil {
		min = F
	} // TODO: Z only for ff == nil ?
	if min > t { // uptype x
		x, t = to(x, min), min
	}
	if t == Z && fx[Z] == nil { // e.g. real functions
		x, t = to(x, F), F
	}
	r = mk(t, n)
	n = atm1(n)
	switch t {
	case L:
		for j := k(0); j < k(n); j++ {
			m.k[r+2+j] = nm(inc(m.k[j+2+x]), rt, fx)
		}
	case D:
		if r != x {
			m.k[2+r] = inc(m.k[2+x])
		}
		m.k[3+r] = nm(m.k[3+x], rt, fx)
	case C, I, F, Z:
		rp, xp, f := ptr(r, t), ptr(x, t), fx[t]
		for i := k(0); i < k(n); i++ {
			f(rp+i, xp+i)
		}
	default:
		panic("type")
	}
	if rt != 0 && t != rt {
		r = to(r, rt)
	}
	return decr(x, r)
}
func ntyps(xt, yt k, fx []f2, fc []fc) (t k) {
	if yt > xt {
		xt = yt
	}
	if xt == C && fc == nil && fx[C] == nil {
		xt = I
	}
	if (xt == I && fc == nil && fx[I] == nil) || (xt == Z && fc == nil && fx[Z] == nil) {
		xt = F
	}
	return xt
}
func nd(x, y, rt k, fx []f2, fc []fc) (r k) { // numeric dyad
	xt, yt, xn, yn := typs(x, y)
	if xt == D {
		r = mk(D, atom)
		m.k[2+r] = inc(m.k[2+x])
		m.k[3+r] = nd(inc(m.k[3+x]), y, rt, fx, fc)
		return decr(x, r)
	} else if yt == D {
		r = mk(D, atom)
		m.k[2+r] = inc(m.k[2+y])
		m.k[3+r] = nd(y, inc(m.k[3+y]), rt, fx, fc)
		return decr(y, r)
	}
	n, sc := xn, k(0)
	if xn == atom {
		n, sc = yn, 1
	} else if yn == atom {
		n, sc = xn, 2
	} else if xn != yn {
		panic("size")
	}
	if xt == L || yt == L {
		r = mk(L, n)
		if xn == atom {
			for i := k(0); i < n; i++ {
				m.k[2+r+i] = nd(inc(x), inc(m.k[2+y+i]), rt, fx, fc)
			}
		} else if yn == atom {
			for i := k(0); i < n; i++ {
				m.k[2+r+i] = nd(inc(m.k[2+x+i]), inc(y), rt, fx, fc)
			}
		} else {
			for i := k(0); i < n; i++ {
				m.k[2+r+i] = nd(inc(m.k[2+x+i]), inc(m.k[2+y+i]), rt, fx, fc)
			}
		}
		return decr2(x, y, uf(r))
	} else if xt > S || yt > S {
		panic("type")
	}
	t := ntyps(xt, yt, fx, fc)
	if xt != t {
		x, xt = to(x, t), t
	}
	if yt != t {
		y, yt = to(y, t), t
	}
	if fc == nil {
		r = mk(t, n)
	} else if t < L {
		r = mk(I, n)
	} else {
		r = mk(L, n)
	}
	n = atm1(n)
	if fc == nil {
		if xt == S {
			panic("type")
		}
		ns(ptr(r, t), ptr(x, t), ptr(y, t), t, n, n, sc, fx[t])
	} else {
		g := fc[t]
		f := func(rp, xp, yp k) {
			if x := g(xp, yp); x {
				m.k[rp] = 1
			} else {
				m.k[rp] = 0
			}
		}
		ns(ptr(r, I), ptr(x, t), ptr(y, t), t, n, n, sc, f)
	}
	if rt != 0 && t > rt {
		r = to(r, rt)
	}
	return decr2(x, y, r)
}
func ns(rp, xp, yp, t, xn, yn, c k, f f2) {
	switch c {
	case 0: // v f v
		for i := k(0); i < xn; i++ {
			f(rp+i, xp+i, yp+i)
		}
	case 1: // a f v
		for i := k(0); i < yn; i++ {
			f(rp+i, xp, yp+i)
		}
	case 2: // v f a
		for i := k(0); i < xn; i++ {
			f(rp+i, xp+i, yp)
		}
	case 3: // f/ x → a
		cpx[t](rp, yp)
		for i := k(1); i < yn; i++ {
			f(rp, rp, yp+i)
		}
	case 4: // x f/ y → v(xn)
		cp := cpx[t]
		for i := k(0); i < xn; i++ {
			rp++
			cp(rp, xp+i)
			for j := k(0); j < yn; j++ {
				f(rp, rp, yp+j)
			}
		}
	case 5: // f\ v → v
		cpx[t](rp, yp)
		for i := k(1); i < yn; i++ {
			f(rp+i, rp+i-1, yp+i)
		}
	case 6: // a f\ v
		f(rp, xp, yp)
		for i := k(1); i < yn; i++ {
			f(rp+i, rp+i-1, yp+i)
		}
	case 7: // v f\ v
		p := rp
		for j := k(0); j < yn; j++ {
			m.k[rp+j] = mk(t, xn)
		}
		for i := k(0); i < xn; i++ {
			rp = m.k[2+i+p]
			f(rp, xp+i, yp)
			for j := k(1); j < yn; j++ {
				f(rp+i, rp+i-1, yp+j)
			}
		}
	case 8: // f': v
		cpx[t](rp, yp)
		for i := k(1); i < yn; i++ {
			f(rp+i, yp+i-1, yp+i)
		}
	case 9: // a f': v (vector x is not handled)
		f(rp, xp, yp)
		for i := k(1); i < yn; i++ {
			f(rp+i, yp+i-1, yp+i)
		}
	case 10: // x f/: y
		for j := k(0); j < yn; j++ {
			m.k[rp+j] = mk(t, xn)
			p := ptr(m.k[rp+j], L)
			for i := k(0); i < xn; i++ {
				f(p+i, xp+i, yp+j)
			}
		}
	case 11: // x f\: y
		for i := k(0); i < xn; i++ {
			m.k[rp+i] = mk(t, yn)
			p := ptr(m.k[rp+i], L)
			for j := k(0); j < yn; j++ {
				f(p+j, xp+i, yp+j)
			}
		}
	}
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
func explode(x k) (r k) { // explode an array (or atom) to a list of atoms
	t, n := typ(x)
	if t == L {
		return x
	} else if t > L {
		panic("type")
	}
	n = atm1(n)
	cp, xp := cpx[t], ptr(x, t)
	r = mk(L, n)
	for i := k(0); i < n; i++ {
		rk := mk(t, atom)
		rp := ptr(rk, t)
		cp(rp, xp+i)
		m.k[2+r+i] = rk
	}
	return decr(x, r)
}
func uf(x k) (r k) { // unify lists if possible
	xt, xn := typ(x)
	if xt != L || xn == 0 {
		return x
	}
	ut := k(0)
	for j := k(0); j < xn; j++ {
		t, n := typ(m.k[2+x+j])
		switch {
		case t >= L || n != atom:
			return x
		case j == 0:
			ut = t
		case t != ut:
			return x
		}
	}
	r = mk(ut, xn)
	cp, rp := cpx[ut], ptr(r, ut)
	for i := k(0); i < xn; i++ {
		cp(rp+i, ptr(m.k[2+x+i], ut))
	}
	return decr(x, r)
}

func idn(x k) (r k) { return x } // :x
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
		switch {
		case j == 0:
			nr, tt = nj, tj
		case nj != nr:
			panic("size") // k7 does extends for atoms, and nan-fills short arrays
		case j > 0 && tt != L && tt != tj:
			tt = L
		}
	}
	if tt > L {
		panic("type")
	}
	nr = atm1(nr)
	cp := cpx[tt]
	r = mk(L, nr)
	for i := k(0); i < nr; i++ {
		rr := mk(tt, n)
		m.k[2+r+i] = rr
	}
	if tt == L {
		for k := k(0); k < n; k++ {
			col := explode(inc(m.k[2+x+k]))
			for i := uint32(0); i < nr; i++ {
				m.k[2+k+m.k[2+r+i]] = inc(m.k[2+col+i])
			}
			dec(col)
		}
		for i := k(0); i < nr; i++ {
			m.k[2+r+i] = uf(m.k[2+r+i])
		}
	} else {
		for i := uint32(0); i < nr; i++ { // Rik = +Xki (cdn.mos.cms.futurecdn.net/XTZkbu7r5c4LZQ5SMzJDbV-970-80.jpg)
			for k := k(0); k < n; k++ {
				cp(k+ptr(m.k[2+r+i], tt), i+ptr(m.k[2+x+k], tt))
			}
		}
	}
	return decr(x, r)
}
func neg(x k) k { // -x
	return nm(x, 0, []f1{nil, func(r, x k) { m.c[r] = -m.c[x] }, func(r, x k) { m.k[r] = k(-i(m.k[x])) }, func(r, x k) { m.f[r] = -m.f[x] }, func(r, x k) { m.z[r] = -m.z[x] }})
}
func fst(x k) (r k) { // *x
	t, n := typ(x)
	if t == D {
		inc(m.k[3+x])
		r = fst(m.k[3+x])
		return decr(x, r)
	}
	if n == atom {
		return x
	} else if n == 0 {
		panic("nyi: fst empty") // what to return? missing value? panic?
	}
	if t == L {
		r = m.k[2+x]
		inc(r)
		return decr(x, r)
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
	return decr(x, r)
}
func inv(x k) k { // %x
	return nm(x, 0, []f1{nil, nil, nil, func(r, x k) { m.f[r] = 1 / m.f[x] }, func(r, x k) { m.z[r] = 1 / m.z[x] }})
}
func wer(x k) (r k) { // &x
	t, n := typ(x)
	if t != I {
		panic("type")
	}
	n = atm1(n)
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
	return decr(x, r)
}
func rev(x k) (r k) { // |x
	t, n := typ(x)
	if n == atom || n < 2 {
		if t == D {
			r = mk(t, n)
			m.k[r+2] = rev(inc(m.k[x+2]))
			m.k[r+3] = rev(inc(m.k[x+3]))
			return decr(x, r)
		}
		return x
	}
	r = mk(t, n)
	if t < D {
		cp, m := cpx[t], n
		if t == L {
			cp = cpI
		}
		src, dst := ptr(x, t), ptr(r, t)
		if src == dst {
			cp, m = swx[t], n/2
		}
		for i := k(0); i < m; i++ {
			cp(dst+n-1-i, src+i)
		}
	} else {
		panic("nyi")
	}
	if t == L {
		for i := k(0); i < n; i++ {
			inc(m.k[2+i+r])
		}
	}
	return decr(x, r)
}
func asc(x k) (r k) { // <x
	t, n := typ(x)
	if n == atom || t >= L { // k7 also sorts lists of different numeric types
		panic("type")
	}
	r = til(mki(n))
	lt, sw := ltx[t], swI
	src, ind, dst := ptr(x, t), 2+r, ptr(r, I)
	for i := k(1); i < n; i++ { // insertion sort, should be replaced
		for j := k(i); j > 0 && lt(src+m.k[ind+j], src+m.k[ind+j-1]); j-- {
			sw(dst+j, dst+(j-1))
		}
	}
	return decr(x, r)
}
func dsc(x k) (r k) { return rev(asc(x)) } // >x
func grp(x k) (r k) { // =x
	t, n := typ(x)
	if n == atom {
		panic("value")
	}
	eq := eqx[t]
	r = mk(D, atom)
	u := unq(inc(x)) // TODO: keys are sorted in k7
	l, nu := mk(L, m.k[u]&atom), m.k[u]&atom
	m.k[2+r], m.k[3+r] = u, l
	up, xp := ptr(u, t), ptr(x, t)
	b := mk(C, n) // boolean
	bc := 8 + b<<2
	for j := k(0); j < nu; j++ { // over ?x
		nr := k(0)
		for jj := k(0); jj < n; jj++ { // over x
			m.c[bc+jj] = 0
			if eq(up+j, xp+jj) {
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
	return decr2(b, x, r)
}
func til(x k) (r k) { // !x
	t, n := typ(x)
	if n != atom {
		panic("nyi !a")
	} else if t == D {
		r = inc(m.k[2+x])
		return decr(x, r)
	} else if t > Z {
		panic("type")
	}
	if nn := idx(x, t); nn < 0 {
		return decr(x, eye(k(-nn)))
	} else {
		return decr(x, jota(k(nn)))
	}
}
func jota(n k) (r k) { // !n
	r = mk(I, n)
	for j := k(0); j < n; j++ {
		m.k[2+r+j] = j
	}
	return r
}
func eye(n k) (r k) { // !-n
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
func not(x k) (r k) { // ~x
	return nm(x, I, []f1{nil, func(r, x k) {
		if m.c[x] == 0 {
			m.c[r] = 1
		} else {
			m.c[r] = 0
		}
	}, func(r, x k) {
		if m.k[x] == 0 {
			m.k[r] = 1
		} else {
			m.k[r] = 0
		}
	}, func(r, x k) {
		if m.f[x] == 0 {
			m.f[r] = 1
		} else {
			m.f[r] = 0
		}
	}, func(r, x k) {
		if m.z[x] == 0 {
			m.z[r] = 1
		} else {
			m.z[r] = 0
		}
	}})
}
func enl(x k) (r k) { // ,x (collaps uniform)
	t, n := typ(x)
	if t < L && n == atom {
		r = mk(t, 1)
		cp := cpx[t]
		src, dst := ptr(x, t), ptr(r, t)
		cp(dst, src)
		return decr(x, r)
	}
	r = mk(L, 1)
	m.k[2+r] = x
	return r
}
func enlist(x k) (r k) { // dont unify
	r = mk(L, 1)
	m.k[2+r] = x
	return r
}
func srt(x k) (r k) { return atx(x, asc(inc(x))) } // ^x  TODO: replace with a sort implementation
func cnt(x k) (r k) { // #x
	t, n := typ(x)
	if t == D {
		_, n = typ(m.k[x+2])
	}
	n = atm1(n)
	return decr(x, mki(k(i(n))))
}
func flr(x k) k { // _x
	return nm(x, I, []f1{nil, func(r, x k) { m.c[r] = m.c[x] }, func(r, x k) { m.k[r] = m.k[x] }, func(r, x k) {
		y := f(i(m.f[x]))
		if m.f[x] < y {
			y -= 1.0
		}
		m.f[r] = y
	}, nil}) // TODO: k7 does not convert c to i
}
func str(x k) (r k) { // $x
	t, n := typ(x)
	if t == C {
		return x
	}
	if t < L {
		st, xp := stx[t], ptr(x, t)
		if n == atom {
			r = mk(C, 56)
			r = srk(r, C, 56, st(8+r<<2, xp))
		} else {
			r = mk(L, n)
			for i := k(0); i < n; i++ {
				y := mk(C, 56)
				m.k[2+r+i] = srk(y, C, 56, st(8+y<<2, xp+i))
			}
		}
	} else {
		switch t {
		case L:
			r = mk(L, n)
			for i := k(0); i < n; i++ {
				m.k[2+i+r] = str(inc(m.k[2+i+x]))
			}
		case D:
			r = mk(D, atom)
			m.k[2+r] = inc(m.k[2+x])
			m.k[3+r] = str(inc(m.k[3+x]))
		case N:
			r = mk(C, 0)
		case N + 1, N + 2, N + 3, N + 4:
			f := m.k[2+x]
			if n == 0 || n == 1 { // 0(lambda), 1(lambda projection)
				r = inc(m.k[2+x]) // `C
			} else if (f >= 39 && f < dyad) || (f >= 83 && f < 100) { // built-ins
				r = str(atx(inc(m.k[2+m.k[3]]), fst(wer(eql(mki(f), inc(m.k[3+m.k[3]]))))))
			} else if f < 33 { // monad +: 3: /:
				r = mk(C, 2)
				m.c[8+r<<2] = m.c[136+m.k[2+x]]
				m.c[9+r<<2] = ':'
			} else if f >= 70 && f < 80 { // dyadic ioverb 3:
				r = mk(C, 2)
				m.c[8+r<<2] = '0' + c(f-70)
				m.c[9+r<<2] = ':'
			} else if f > 49 && f < 83 { // dyad * /
				r = mk(C, 1)
				m.c[8+r<<2] = m.c[136+m.k[2+x]-dyad]
			}
			// TODO: derived?
			if n == 1 || n == 2 { // projection
				a := m.k[3+x]
				if n == 2 && f < 100 && m.k[m.k[3+a]]>>28 == N {
					r = cat(kst(inc(m.k[2+a])), r) // short form: 2+
				} else {
					a = kst(inc(a))   // arg list
					m.c[8+a<<2] = '[' // convert () to []
					m.c[7+(m.k[a]&atom)+a<<2] = ']'
					r = cat(r, a)
				}
			}

		default:
			panic("nyi")
		}
	}
	return decr(x, r)
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
	eq, cp := eqx[t], cpx[t]
	if t == L {
		eq = func(x, y k) bool { return match(m.k[x], m.k[y]) }
	}
	src, dst := ptr(x, t), ptr(r, t)
	nn := k(0)
	for i := k(0); i < n; i++ { // quadratic, should be improved
		u := true
		srci := src + i
		for j := k(0); j < nn; j++ {
			if eq(srci, dst+j) {
				u = false
				break
			}
		}
		if u {
			cp(dst+nn, srci)
			nn++
		}
	}
	return decr(x, srk(r, t, n, nn))
}
func tip(x k) (r k) { // @x
	r = mku(0)
	t, n := typ(x)
	s := m.c[169+t]
	if n != atom && t < L && s != 0 {
		s -= 32
	}
	mys(8+r<<2, uint64(s)<<56)
	return decr(x, r)
}
func evl(x k) (r k) { // .x
	t, n := typ(x)
	if t != L {
		if t == S && n == 1 {
			return fst(x)
		} else if t == S && n == atom {
			return lup(x)
		}
		return x
	}
	if n == 0 {
		panic("evl empty list?") // what TODO?
	}
	v := m.k[2+x]
	vt, vn := typ(v)
	if vt == S {
		if n == 1 { // ,`a`b → `a`b
			inc(v)
			return decr(x, v)
		}
		if m.k[2+v] == 0 && m.k[3+v] == 0 { // (`;…) → ex;ex…
			for i := k(1); i < n; i++ {
				if i > 1 {
					dec(r)
				}
				r = evl(inc(m.k[2+i+x]))
			}
			return decr(x, r)
		}
	}
	switch vt {
	case N: // (;…) → list
		r = mk(L, n-1)
		for i := int(n - 2); i >= 0; i-- {
			m.k[2+r+k(i)] = evl(inc(m.k[3+x+k(i)]))
		}
		return decr(x, uf(r))
	default:
		inc(v)
		if vt == S && vn == atom {
			v = lup(v)
			vt, vn = typ(v)
		}
		if vt > N && m.k[2+v] == dyad { // ':'
			if n != 3 {
				panic("nyi modified assignment")
			}
			name, val := inc(m.k[3+x]), evl(inc(m.k[4+x]))
			return decr2(v, x, asn(name, val))
		}
		r = mk(L, n-1)
		for i := int(n - 2); i >= 0; i-- {
			m.k[2+r+k(i)] = evl(inc(m.k[3+x+k(i)]))
		}
		dec(x)
		v = evl(v)
		vt, vn := typ(v)
		if vt > N && !(vn == atom && vt == N+1 && n-1 == 2 && m.k[2+v] > 255) { // allow dyadic derived
			if n-1 > vt-N {
				panic("args") // too many arguments
			}
			for i := n - 1; i < vt-N; i++ { // fill args, e.g. 2+
				r = lcat(r, mk(N, atom))
			}
			for i := k(0); i < m.k[r]&atom; i++ {
				if m.k[m.k[2+i+r]]>>28 == N {
					return prj(v, r)
				}
			}
		}
		return cal(v, r)
	}
	return x
}
func prj(x, y k) (r k) { // convert x to a projection
	t := m.k[x] >> 28
	r = mk(t, 2)
	if f := m.k[2+x]; f < 256 {
		m.k[2+r] = f // #1: function code if < 256
		dec(x)
	} else {
		m.k[2+r] = x // #1: pointer to lambda function if code >= 256
	}
	m.k[3+r] = y // #2: argument list with holes
	n := k(0)
	for i := k(0); i < m.k[y]&atom; i++ {
		if m.k[m.k[2+y+i]]>>28 == N {
			n++
		}
	}
	m.k[r] = k(N+n)<<28 | 2 // a projection has length 2
	return r
}
func kst(x k) (r k) { // `k@x
	t, n := typ(x)
	atm := n == atom
	if atm {
		n = 1
	}
	if n == 0 && t < N {
		r = mk(C, 0)
		rc, rn := 8+r<<2, k(0)
		switch t { // these could also be in the k-tree
		case C:
			rn = putb(rc, rn, []c(`""`))
		case I:
			rn = putb(rc, rn, []c("!0"))
		case F:
			rn = putb(rc, rn, []c("!0.0"))
		case Z:
			rn = putb(rc, rn, []c("!0i0"))
		case S:
			rn = putb(rc, rn, []c("0#`"))
		case L:
			rn = putb(rc, rn, []c("()"))
		case D:
			rn = putb(rc, rn, []c("()!()"))
		default:
			panic("nyi")
		}
		m.k[r] = C<<28 | rn
		return decr(x, r)
	}
	switch t {
	case C: // ,"a" "a" "ab" "a\nb" ,0x01 0x010203
		r = mk(C, 2+2*n) // for both "a\nb" or 0x01234 or ,"\n"(short enough)
		// no need to shrink: 2*(10+n) is never <= 10+2*n
		rc, rn, xc := 8+r<<2, k(0), 8+x<<2
		if n == 1 && !atm {
			rn = putc(rc, rn, ',')
		}
		rn = putc(rc, rn, '"')
		hex := false
		for i := k(0); i < n; i++ {
			c := m.c[xc+i]
			if c < 32 || c > 126 || c == '"' {
				if c, o := qt(c); o {
					rn = putc(rc, rn, '\\')
					rn = putc(rc, rn, c)
				} else {
					hex = true
					break
				}
			} else {
				rn = putc(rc, rn, c)
			}
		}
		rn = putc(rc, rn, '"')
		if hex {
			rn = 0
			if n == 1 && !atm {
				rn = 1
			}
			rn = putc(rc, rn, '0')
			rn = putc(rc, rn, 'x')
			for i := k(0); i < n; i++ {
				c1, c2 := hxb(m.c[xc+i])
				rn = putc(rc, rn, c1)
				rn = putc(rc, rn, c2)
			}
		}
		m.k[r] = C<<28 | rn
	case I, F, Z:
		r = mk(C, 0)
		if n == 1 && !atm {
			m.c[8+r<<2] = ','
			m.k[r] = C<<28 | 1
		}
		rr := mk(C, 56)
		st, xp, rrc := stx[t], ptr(x, t), 8+rr<<2
		sp := mk(C, 1)
		m.c[8+sp<<2] = ' '
		for i := k(0); i < n; i++ {
			rn := st(rrc, xp+i)
			m.k[rr] = C<<28 | rn
			r = cat(r, inc(rr))
			if i < n-1 {
				r = cat(r, inc(sp))
			}
		}
		dec(sp)
		m.k[rr] = C<<28 | 56
		dec(rr)
		if t == F {
			_, n = typ(r)
			rc, dot := 8+r<<2, false
			for i := k(0); i < n; i++ {
				if m.c[rc+i] == '.' {
					dot = true
					break
				}
			}
			if !dot {
				f := mk(C, 1)
				m.c[8+f<<2] = 'f'
				r = cat(r, f)
			}
		}
	case S:
		if atm || n == 1 {
			rr := mk(C, 0)
			sn, rrc, rn, q := stS(8+rr<<2, ptr(x, S)), 8+rr<<2, k(1), false
			for i := k(0); i < sn; i++ {
				c := m.c[rrc+i]
				if !(cr09(c) || craZ(c) || c == '.') {
					q = true
				}
				if _, o := qt(c); o {
					rn++
				}
				rn++
			}
			if q {
				rn += 2
			}
			if !atm {
				rn++
			}
			r = mk(C, rn)
			rc, rn := 8+r<<2, k(0)
			if !atm {
				rn = putc(rc, rn, ',')
			}
			rn = putc(rc, rn, '`')
			if q {
				rn = putc(rc, rn, '"')
			}
			for i := k(0); i < sn; i++ {
				c, o := qt(m.c[rrc+i])
				if o {
					rn = putc(rc, rn, '\\')
				}
				rn = putc(rc, rn, c)
			}
			if q {
				rn = putc(rc, rn, '"')
			}
			dec(rr)
		} else {
			r = mk(C, 0)
			ix := mk(I, atom)
			for i := k(0); i < n; i++ {
				m.k[2+ix] = i
				r = cat(r, kst(atx(inc(x), inc(ix))))
			}
			dec(ix)
		}
	case L:
		r = mk(C, 1)
		rc := 8 + r<<2
		m.c[rc] = '('
		if n == 1 {
			m.c[rc] = ','
		}
		y := mk(C, 1)
		m.c[8+y<<2] = ';'
		for i := k(0); i < n; i++ {
			r = cat(r, kst(inc(m.k[2+i+x])))
			if i < n-1 {
				r = cat(r, inc(y))
			}
		}
		if n != 1 {
			m.c[8+y<<2] = ')'
			r = cat(r, y)
		} else {
			dec(y)
		}
	case D:
		r = mk(C, 0)
		rr, encl := kst(inc(m.k[x+2])), false
		kt, nk := typ(m.k[x+2])
		if (kt < L && nk == 1) || (kt == D) || (kt > D) {
			encl = true
		}
		y := mk(C, 1)
		if encl {
			m.c[8+y<<2] = '('
			r = cat(r, inc(y))
			r = cat(r, rr)
			m.c[8+y<<2] = ')'
			r = cat(r, inc(y))
		} else {
			r = cat(r, rr)
		}
		m.c[8+y<<2] = '!'
		r = cat(r, y)
		r = cat(r, kst(inc(m.k[x+3])))
	case N, N + 1, N + 2, N + 3, N + 4:
		r = str(inc(x))
	default:
		panic("nyi")
	}
	return decr(x, r)
}
func putc(rc, rn k, c c) k { // assumes enough space
	m.c[rc+rn] = c
	return rn + 1
}
func putb(rc, rn k, b []c) k {
	rc += rn
	copy(m.c[rc:rc+k(len(b))], b)
	return rn + k(len(b))
}
func qt(c c) (c, bool) { // quote
	if c == '"' {
		return c, true
	} else if c == '\n' {
		return 'n', true
	} else if c == '\t' {
		return 't', true
	} else if c == '\r' {
		return 'r', true
	}
	return c, false
}
func sym(x k) uint64 { return *(*uint64)(unsafe.Pointer(&m.c[x])) }
func mys(x k, u uint64) k {
	var b [8]c
	b = *(*[8]c)(unsafe.Pointer(&u))
	copy(m.c[x:x+8], b[:])
	return x
}
func btou(b []c) uint64 {
	if len(b) > 8 {
		panic("size")
	}
	var u uint64
	for i := k(0); i < k(len(b)); i++ {
		u |= uint64(b[i]) << (8 * c(7-i))
	}
	return u
}
func op2(code k) ([]f2, []fc) {
	switch code - dyad {
	case 1: // +
		return []f2{nil, adC, adI, adF, adZ, nil}, nil
	case 2: // -
		return []f2{nil, sbC, sbI, sbF, sbZ, nil}, nil
	case 3: // *
		return []f2{nil, muC, muI, muF, muZ, nil}, nil
	case 4: // %
		return []f2{nil, nil, nil, diF, diZ, nil}, nil
	case 5: // &
		return []f2{nil, miC, miI, miF, miZ, miS}, nil
	case 6: // |
		return []f2{nil, maC, maI, maF, maZ, maS}, nil
	case 7: // <
		return nil, ltx
	case 8: // >
		return nil, gtx
	case 9: // =
		return nil, eqx
	default:
		return nil, nil
	}
}
func add(x, y k) (r k) { f, g := op2(1 + dyad); return nd(x, y, 0, f, g) } // x+y
func sub(x, y k) (r k) { f, g := op2(2 + dyad); return nd(x, y, 0, f, g) } // x-y
func mul(x, y k) (r k) { f, g := op2(3 + dyad); return nd(x, y, 0, f, g) } // x*y
func div(x, y k) (r k) { f, g := op2(4 + dyad); return nd(x, y, 0, f, g) } // x%y
func min(x, y k) (r k) { f, g := op2(5 + dyad); return nd(x, y, 0, f, g) } // x&y
func max(x, y k) (r k) { f, g := op2(6 + dyad); return nd(x, y, 0, f, g) } // x|y
func les(x, y k) (r k) { return nd(x, y, I, nil, ltx) }                    // x<y
func mor(x, y k) (r k) { return nd(x, y, I, nil, gtx) }                    // x>y
func eql(x, y k) (r k) { return nd(x, y, I, nil, eqx) }                    // x=y
func adC(r, x, y k)    { m.c[r] = m.c[x] + m.c[y] }
func adI(r, x, y k)    { m.k[r] = m.k[x] + m.k[y] }
func adF(r, x, y k)    { m.f[r] = m.f[x] + m.f[y] }
func adZ(r, x, y k)    { m.z[r] = m.z[x] + m.z[y] }
func sbC(r, x, y k)    { m.c[r] = m.c[x] - m.c[y] }
func sbI(r, x, y k)    { m.k[r] = m.k[x] - m.k[y] }
func sbF(r, x, y k)    { m.f[r] = m.f[x] - m.f[y] }
func sbZ(r, x, y k)    { m.z[r] = m.z[x] - m.z[y] }
func muC(r, x, y k)    { m.c[r] = m.c[x] * m.c[y] }
func muI(r, x, y k)    { m.k[r] = m.k[x] * m.k[y] }
func muF(r, x, y k)    { m.f[r] = m.f[x] * m.f[y] }
func muZ(r, x, y k)    { m.z[r] = m.z[x] * m.z[y] }
func diF(r, x, y k)    { m.f[r] = m.f[x] / m.f[y] }
func diZ(r, x, y k)    { m.z[r] = m.z[x] / m.z[y] }
func miC(r, x, y k)    { m.c[r] = m.c[ter(m.c[x] < m.c[y], x, y)] }
func miI(r, x, y k)    { m.k[r] = m.k[ter(i(m.k[x]) < i(m.k[y]), x, y)] }
func miF(r, x, y k)    { m.f[r] = m.f[ter(m.f[x] < m.f[y], x, y)] }
func miZ(r, x, y k)    { m.z[r] = m.z[ter(ltZ(x, y), x, y)] }
func miS(r, x, y k)    { m.f[r] = m.f[ter(ltS(x, y), x, y)] }
func maC(r, x, y k)    { m.c[r] = m.c[ter(m.c[x] > m.c[y], x, y)] }
func maI(r, x, y k)    { m.k[r] = m.k[ter(i(m.k[x]) > i(m.k[y]), x, y)] }
func maF(r, x, y k)    { m.f[r] = m.f[ter(m.f[x] > m.f[y], x, y)] }
func maZ(r, x, y k)    { m.z[r] = m.z[ter(gtZ(x, y), x, y)] }
func maS(r, x, y k)    { m.f[r] = m.f[ter(gtS(x, y), x, y)] }
func ter(b bool, x, y k) k {
	if b {
		return x
	} else {
		return y
	}
}
func key(x, y k) (r k) { // x!y
	_, yt, xn, yn := typs(x, y)
	if xn == atom {
		x, xn = enl(x), 1
	}
	if yn == atom {
		y, yn = ext(y, yt, xn), xn
	}
	if xn == 1 && yn > 1 {
		y, yn = enl(y), 1
	}
	if xn != yn {
		panic("length")
	}
	r = mk(D, atom)
	m.k[2+r] = x
	m.k[3+r] = y
	return r
}
func ext(x, t, n k) (r k) { // scalar extension
	r = mk(t, n)
	xp, rp, cp := ptr(x, t), ptr(r, t), cpx[t]
	for i := k(0); i < n; i++ {
		cp(rp+i, xp)
	}
	return decr(x, r)
}
func mch(x, y k) (r k) { // x~y
	r = mkc(0)
	if match(x, y) {
		m.c[8+r<<2] = 1
	}
	return decr2(x, y, r)
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
	n = atm1(n)
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
		eq := eqx[t]
		if eq == nil {
			panic("type")
		}
		x, y = ptr(x, t), ptr(y, t)
		for j := k(0); j < n; j++ {
			if eq(x+j, y+j) == false {
				return false
			}
		}
		return true
	}
	return false
}
func cat(x, y k) (r k) { // x,y
	xt, yt, xn, yn := typs(x, y)
	switch {
	case xt < L && yt == xt:
		return ucat(x, y, xt, xn, yn)
	case xt > L || yt > L: // TODO: cat D?
		panic("type")
	case xt != L:
		x = explode(x)
		xt, xn = typ(x)
	case yt != L:
		y = explode(y)
		xt, yn = typ(y)
	}
	if m.k[xt+1] > 1 || bk(L, xn+yn) > bk(L, xn) {
		r = mk(L, xn+yn)
		for j := k(0); j < xn; j++ {
			m.k[2+r+j] = inc(m.k[2+x+j])
		}
		dec(x)
	} else {
		r = x
	}
	for j := k(0); j < yn; j++ {
		m.k[2+r+xn+j] = inc(m.k[2+y+j])
	}
	return decr(y, r)
}
func ucat(x, y, t, xn, yn k) (r k) { // x, y same type < L
	xn, yn = atm1(xn), atm1(yn)
	cp, xp := cpx[t], ptr(x, t)
	if m.k[x+1] > 1 || bk(t, xn+yn) != bk(t, xn) {
		r = mk(t, xn+yn)
		rp := ptr(r, t)
		for i := k(0); i < xn; i++ {
			cp(rp+i, xp+i)
		}
	} else {
		r = x
		m.k[r] = t<<28 | (xn + yn)
	}
	rp, yp := xn+ptr(r, t), ptr(y, t)
	for i := k(0); i < yn; i++ {
		cp(rp+i, yp+i)
	}
	if r != x {
		dec(x)
	}
	return decr(y, r)
}
func lcat(x, y k) (r k) { // append anything to a list; no unify
	_, nl := typ(x)
	if m.k[x+1] == 1 && bk(L, nl) == bk(L, nl+1) {
		m.k[2+x+nl] = y
		m.k[x] = L<<28 | (nl + 1)
		return x
	}
	r = mk(L, nl+1)
	for i := k(0); i < nl; i++ {
		m.k[2+i+r] = inc(m.k[2+i+x])
	}
	m.k[2+nl+r] = y
	return decr(x, r)
}
func ept(x, y k) (r k) { // x^y
	t, yt, n, yn := typs(x, y)
	if t != yt || t > L || n == atom {
		panic("type")
	} else if yn == atom {
		y, yn = enl(y), 1
	}
	eq, b, xp, yp := eqx[t], mk(I, n), ptr(x, t), ptr(y, t)
	if t == L {
		eq = match
	}
	all := true
	for i := k(0); i < n; i++ { // TODO: quadratic
		m.k[2+i+b] = 1
		for j := k(0); j < yn; j++ {
			if eq(xp+i, yp+j) {
				m.k[2+i+b], all = 0, false
				break
			}
		}
	}
	if all {
		return decr2(b, y, x)
	}
	return decr(y, atx(x, wer(b)))
}
func tak(x, y k) (r k) { // x#y
	xt, yt, xn, yn := typs(x, y)
	if yt == D {
		return key(x, atx(y, inc(x)))
	}
	if xt != I {
		panic("type")
	}
	if xn == atom {
		xn = 1
	} else if xn == 0 {
		dec(x)
		return fst(y)
	}
	yn = atm1(yn)
	n, o := m.k[2+x+xn-1], k(0) // n:-1#x
	if i(n) < 0 {
		if yn != 0 {
			o = k(i(yn) + ((i(yn) + i(n)) % i(yn)))
		}
		n = k(-i(n))
	}
	if xn == 1 {
		return decr(x, take(n, o, y))
	}
	r, o = rsh(2+x, xn-1, n, o, y, yn)
	return decr2(x, y, r)
}
func rsh(xp, xn, n, o, y, yn k) (r, oo k) { // reshape (with offset): (x,n)#y
	a := m.k[xp]
	if i(a) < 0 {
		panic("domain")
	}
	r = mk(L, a)
	for i := k(0); i < a; i++ {
		if xn > 1 {
			m.k[2+i+r], o = rsh(xp+1, xn-1, n, o, y, yn)
		} else {
			m.k[2+i+r] = take(n, o, inc(y))
			o = (o + n) % yn
		}
	}
	return r, o
}
func take(n, o, y k) (r k) { // integer index and offset
	t, yn := typ(y)
	cp, yp := cpx[t], ptr(y, t)
	if yn == 0 {
		r = mk(t, n)
		rp, na := ptr(r, t), nax[t]
		for i := k(0); i < n; i++ {
			na(rp + i)
		}
		return decr(y, r)
	}
	r = mk(t, n)
	rp := ptr(r, t)
	for i := k(0); i < n; i++ {
		cp(rp+i, yp+((i+o)%yn))
	}
	return decr(y, uf(r))
}
func drp(x, y k) (r k) { // x_y
	xt, t, xn, yn := typs(x, y)
	if t == D { // x_d: del
		u := ept(inc(m.k[y+2]), x)
		if r == m.k[y+2] {
			return y
		}
		return key(u, atx(y, inc(u)))
	} else if xt != I {
		panic("type")
	} else if yn == atom {
		panic("rank")
	} else if xn != atom {
		return cut(x, y)
	}
	n := m.k[2+x]
	return decr(x, uf(drop(i(n), y)))
}
func drop(x i, y k) (r k) { // integer index; does not unify
	t, yn := typ(y)
	if yn == atom {
		panic("type")
	}
	n, neg, o := k(x), false, k(x)
	if x < 0 {
		n, neg, o = k(-x), true, 0
	}
	yp, cp := ptr(y, t), cpx[t]
	if m.k[1+y] == 1 && t != L {
		if neg {
			return srk(y, t, yn, yn-n)
		}
		for i := k(0); i < yn-n; i++ {
			cp(yp+i, yp+o+i)
		}
		return uf(srk(y, t, yn, yn-n)) // uf? TODO rm
	}
	r = mk(t, yn-n)
	rp := ptr(r, t)
	for i := k(0); i < yn-n; i++ {
		cp(rp+i, yp+o+i)
	}
	return decr(y, r)
}
func cut(x, y k) (r k) { // x_y
	xt, yt, xn, yn := typs(x, y)
	if xt != I || yn == atom {
		panic("type")
	}
	for i := k(0); i < xn; i++ {
		if a := m.k[2+x+i]; int32(a) < 0 || (i > 0 && m.k[1+x+i] > a) || a > yn {
			panic("domain")
		}
	}
	r = mk(L, xn)
	cp, yp := cpx[yt], ptr(y, yt)
	for i := k(0); i < xn; i++ {
		nn := yn
		if i < xn-1 {
			nn = m.k[3+i+x]
		}
		ln := nn - m.k[2+i+x]
		a := mk(yt, ln)
		yp, ap := yp+m.k[2+i+x], ptr(a, yt)
		for j := k(0); j < ln; j++ {
			cp(ap+j, yp+j)
		}
		yp += ln
		m.k[2+i+r] = uf(a)
	}
	return decr2(x, y, r)
}
func cst(x, y k) (r k) { // x$y
	xt, yt, xn, yn := typs(x, y)
	if xt != S || xn != atom {
		panic("type")
	}
	s := c(sym(8+x<<2) >> 56)
	t, o := k(0), k(169)
	for i := o; i < o+15; i++ {
		if s == m.c[i] {
			t = i - o
		}
	}
	if t < 1 || t >= L || yt >= L {
		panic("type")
	} else if t == S && yt == C { // TODO: or `$x
		yn = atm1(yn)
		yc := 8 + y<<2
		r = mku(btou(m.c[yc : yc+yn]))
		return decr2(x, y, r)
	}
	return decr(x, to(y, t)) // TODO other conversions?
}
func fnd(x, y k) (r k) { // x?y
	t, yt, xn, yn := typs(x, y)
	if xn == atom || t != yt {
		panic("type")
	}
	r = mk(I, yn)
	yn = atm1(yn)
	eq, xp, yp := eqx[t], ptr(x, t), ptr(y, t)
	if t == L {
		eq = match
	}
	for j := k(0); j < yn; j++ {
		n := xn // TODO: or 0N?
		for i := k(0); i < xn; i++ {
			if eq(xp+i, yp+j) {
				n = i
				break
			}
		}
		m.k[2+j+r] = n
	}
	return decr2(x, y, r)
}
func atx(x, y k) (r k) { // x@y
	xt, yt, xn, yn := typs(x, y)
	if xn == atom && xt == S {
		x = str(x)
		if m.k[x]&atom != 1 {
			panic("class")
		}
		dec(x)
		switch m.c[8+x<<2] {
		case 'p':
			return prs(y)
		case 'k':
			return kst(y)
		default:
			panic("class")
		}
	} else if xt > N {
		return cal(x, enl(y))
	} else if xn == atom && xt != D {
		panic("type")
	}
	switch {
	case xt < L && yt == I:
		cp, na, xp := cpx[xt], nax[xt], ptr(x, xt)
		r = mk(xt, yn)
		yn = atm1(yn)
		rp, yp := ptr(r, xt), 2+y
		for i := k(0); i < yn; i++ {
			if i >= xn {
				na(rp + i)
			} else {
				cp(rp+i, xp+m.k[yp+i])
			}
		}
		return decr2(x, y, r)
	case xt == L && yt == I:
		if yn == atom {
			r = inc(m.k[2+x+m.k[2+y]])
		} else {
			r = mk(L, yn)
			for i := k(0); i < yn; i++ {
				m.k[2+r+i] = inc(m.k[2+x+m.k[2+y+i]])
			}
			r = uf(r)
		}
		return decr2(x, y, r)
	case xt == D:
		keys := m.k[2+x]
		kt, nk := typ(keys)
		vt, _ := typ(m.k[3+x])
		if kt != yt {
			panic("type")
		}
		r = mk(vt, yn)
		yn = atm1(yn)
		cp, na, eq, kp, vp, rp, yp := cpx[vt], nax[vt], eqx[kt], ptr(keys, kt), ptr(m.k[3+x], vt), ptr(r, vt), ptr(y, yt)
		for i := k(0); i < yn; i++ {
			na(rp + i)
			for j := k(0); j < nk; j++ {
				if eq(kp+j, yp+i) {
					cp(rp+i, vp+j)
					break
				}
			}
		}
		return decr2(x, y, r)
	// case xt == L:
	//	missing element for a list is nax[type of first element]
	default:
		panic("nyi atx")
	}
}
func cal(x, y k) (r k) { // x.y
	xt, yt, xn, yn := typs(x, y)
	if xt <= D { // TODO dict
		if yt == L {
			if yn == 0 {
				return decr(y, x)
			}
			return cal(cal(x, fst(inc(y))), drop(1, y)) // at depth
		}
		return atx(x, y)
	}
	y = explode(y)
	if xn == 1 || xn == 2 { // convert projected to full call
		l := m.k[x+3] // arg list with holes
		n := m.k[l] & atom
		if n != xt-N+yn {
			panic("valence")
		}
		a, l, yi := mk(L, n), m.k[x+3], k(0) // a: full arg vector
		for i := k(0); i < n; i++ {
			if v := m.k[2+l+i]; m.k[v]>>28 == N {
				m.k[2+a+i] = inc(m.k[2+y+yi])
				yi++
			} else {
				m.k[2+a+i] = inc(v)
			}
		}
		dec(y)
		r = m.k[x+2]
		if f := m.k[x+2]; xn == 1 { // lambda projection
			r, xn = inc(f), 0
		} else {
			r = mk(N+1, atom)
			m.k[2+r] = f
			if f >= dyad {
				m.k[r] = (N+2)<<28 | atom
			}
		}
		dec(x)
		x, y, xt, yn = r, a, N+n, n
	}
	if xn == 0 {
		return lambda(x, y)
	}
	code := m.k[2+x]
	if code > 255 { // derived
		code >>= 8
		switch yn {
		case 1:
			f := table[code].(func(k, k) k)
			a, b := inc(m.k[3+x]), inc(m.k[2+y])
			r = f(a, b)
		case 2:
			f := table[code+dyad].(func(k, k, k) k)
			a, b, c := inc(m.k[3+x]), inc(m.k[2+y]), inc(m.k[3+y])
			r = f(a, b, c)
		default:
			panic("valence")
		}
		return decr2(x, y, r)
	}
	switch xt {
	case N + 1:
		if yn != 1 {
			panic("valence") // TODO projection
		}
		f := table[code].(func(k) k)
		r = f(inc(m.k[2+y]))
	case N + 2:
		if yn != 2 {
			panic("valence") // TODO projection
		}
		f := table[code].(func(k, k) k)
		r = f(inc(m.k[2+y]), inc(m.k[3+y]))
	default:
		panic("nyi")
	}
	return decr2(x, y, r)
}
func cal2(f, x, y k) (r k) {
	l := mk(L, 2)
	m.k[2+l] = x
	m.k[3+l] = y
	return cal(f, l)
}
func lambda(x, y k) (r k) { // call lambda
	v := (m.k[x] >> 28) - N
	if v < 1 || v > 3 {
		panic("valence")
	}
	if yt, yn := typ(y); yt != L || yn != v {
		panic("args")
	}
	l := m.k[3+x] // lambda tree
	lt, nl := typ(l)
	if nl == 0 {
		return decr2(x, y, mk(N, atom))
	} else if lt != L {
		panic("type")
	}
	var xx, vv [3]k // `x `y `z, x  y  z (save old values)
	for i := k(0); i < v; i++ {
		xx[i] = mku(uint64('x'+i) << 56)
		vv[i] = lupo(inc(xx[i]))
		dec(asn(inc(xx[i]), inc(m.k[2+y+i])))
	}
	dec(y)
	// TODO assign f
	for i := k(0); i < nl; i++ {
		if r > 0 {
			dec(r)
		}
		r = evl(inc(m.k[2+i+l]))
	}
	for i := k(0); i < v; i++ {
		if vv[i] != 0 { // reassign old value
			dec(asn(xx[i], vv[i]))
		} else { // delete variable
			dec(del(xx[i]))
		}
	}
	return decr(x, r)
}
func qot(x k) (r k) { return drv(0, x) } // '
func sla(x k) (r k) { // /
	// TODO: as verb: div
	t, _ := typ(x)
	if t < N+1 {
		panic("nyi div")
	}
	return drv(1, x)
}
func bsl(x k) (r k) { return drv(2, x) } // \
func qtc(x k) (r k) { return drv(3, x) } // ':
func slc(x k) (r k) { return drv(4, x) } // /:
func bsc(x k) (r k) { return drv(5, x) } // \:
func drv(op k, x k) (r k) { // derived function
	r = mk(N+1, atom)
	m.k[2+r] = (33 + op) << 8 // op: 0(') 1(/) 2(\) 3(':) 4(/:) 5(\:)
	m.k[3+r] = x
	return r
}
func ech(f, x k) (r k) { // f'x
	t, n := typ(x)
	if t == D {
		r = mk(D, atom)
		m.k[2+r] = inc(m.k[2+x])
		m.k[3+r] = ech(f, inc(m.k[3+x]))
		return decr(x, r)
	}
	if n == atom {
		return atx(f, x)
	}
	r = mk(L, n)
	for i := k(0); i < n; i++ {
		m.k[2+r+i] = atx(inc(f), atx(inc(x), mki(i)))
	}
	return decr2(f, x, uf(r))
}
func ecd(f, x, y k) (r k) { // x f'y (each pair)
	xt, yt, xn, yn := typs(x, y)
	if xn == atom || yn == atom {
		return cal2(f, x, y)
	}
	if xn == atom {
		x, xn = ext(x, xt, yn), yn
	} else if yn == atom {
		y, yn = ext(y, yt, xn), xn
	}
	if xn != yn {
		panic("length")
	}
	r = mk(L, xn)
	for i := k(0); i < xn; i++ {
		m.k[2+i+r] = cal2(inc(f), atx(inc(x), mki(i)), atx(inc(y), mki(i)))
	}
	dec(f)
	return decr2(x, y, uf(r))
}
func ecp(f, x k) (r k) { // f':x (each prior)
	t, n := typ(x)
	if t == D {
		r = mk(D, atom)
		m.k[2+r] = inc(m.k[2+x])
		m.k[3+r] = ecp(f, inc(m.k[3+x]))
		return decr(x, r)
	} else if t > L || n == atom {
		panic("class")
	} else if n == 1 {
		dec(f)
		return x
	}
	if x, t, op := scop1(f, x); op != nil {
		r = mk(t, n)
		rp, xp := ptr(r, t), ptr(x, t)
		cpx[t](rp, xp)
		for i := k(1); i<n; i++ {
			op(rp+i, xp+i, xp+i-1)
		}
		return decr(x, r)
	}
	r = mk(L, n)
	m.k[2+r] = atx(inc(x), mki(0))
	for i := k(1); i < n; i++ {
		m.k[2+r+i] = cal2(inc(f), atx(inc(x), mki(i)), atx(inc(x), mki(i-1)))
	}
	return decr2(f, x, uf(r))
}
func epi(f, x, y k) (r k) { // x f':y (each prior initial)
	n := m.k[y] & atom
	if m.k[x]&atom == atom {
		if x, t, op := scop1(f, x); op != nil {
			r = mk(t, n)
			rp, xp, yp := ptr(r, t), ptr(x, t), ptr(y, t)
			op(rp, yp, xp)
			for i := k(1); i<n; i++ {
				op(rp+i, yp+i, yp+i-1)
			}
			return decr2(x, y, r)
		}
	}
	r = mk(L, n)
	m.k[2+r] = cal2(inc(f), x, atx(inc(x), mki(0)))
	for i := k(1); i < n; i++ {
		m.k[2+r+i] = cal2(inc(f), atx(inc(x), mki(i)), atx(inc(x), mki(i-1)))
	}
	return decr2(f, x, uf(r))
}
func ecr(f, x, y k) (r k) { // x f/: y
	xt, yt, _, yn := typs(x, y)
	if xt > L || yt > L { // TODO D
		panic("type")
	} else if yn == atom {
		return cal2(f, x, y)
	}
	r = mk(L, yn)
	for i := k(0); i < yn; i++ {
		m.k[2+r+i] = cal2(inc(f), inc(x), atx(inc(y), mki(i)))
	}
	dec(f)
	return decr2(x, y, r)
}
func ecl(f, x, y k) (r k) { // x f\: y
	xt, yt, xn, _ := typs(x, y)
	if xt > L || yt > L { // TODO D
		panic("type")
	} else if xn == atom {
		return cal2(f, x, y)
	}
	r = mk(L, xn)
	for i := k(0); i < xn; i++ {
		m.k[2+r+i] = cal2(inc(f), atx(inc(x), mki(i)), inc(y))
	}
	dec(f)
	return decr2(x, y, r)
}
func ovr(f, x k) (r k) { // f/x
	if x, t, op := scop1(f, x); op != nil {
		n := m.k[x] & atom
		if n == atom {
			return x
		}
		r := mk(t, atom)
		rp, xp, cp := ptr(r, t), ptr(x, t), cpx[t]
		cp(rp, xp)
		for i := k(1); i < n; i++ {
			op(rp, rp, xp+i)
		}
		return decr(x, r)
	}
	return ovsc(f, x, false)
}
func scn(f, x k) (r k) { // f\x
	if x, t, op := scop1(f, x); op != nil {
		n := m.k[x] & atom
		if n == atom {
			return x
		}
		r := mk(t, n)
		rp, xp, cp := ptr(r, t), ptr(x, t), cpx[t]
		cp(rp, xp)
		for i := k(1); i < n; i++ {
			op(rp+i, rp+i-1, xp+i)
		}
		return decr(x, r)
	}
	return ovsc(f, x, true)
}
func ovsc(f, x k, scan bool) (r k) {
	t, n := typ(x)
	if t == D {
		if scan {
			r = mk(D, atom)
			m.k[2+r] = inc(m.k[2+x])
			m.k[3+r] = ovsc(f, inc(m.k[3+x]), scan)
			return decr(x, r)
		} else {
			r = inc(m.k[3+x])
			dec(x)
			x = r
			t, n = typ(x)
		}
	}
	if t > L {
		panic("type")
	}
	if n == atom {
		x = enl(x)
	}
	if n == 0 { // TODO: k7 returns special values for empty arrays
		return decr(f, x)
	} else if n == 1 {
		return decr(f, fst(x))
	}
	if scan {
		r = mk(L, n)
		m.k[2+r] = atx(inc(x), mki(0))
		p := m.k[2+r]
		for i := k(1); i < n; i++ {
			p = cal2(inc(f), inc(p), atx(inc(x), mki(i)))
			m.k[2+i+r] = p
		}
		return decr2(x, f, uf(r))
	} else {
		r = atx(inc(x), mki(0))
		for i := k(1); i<n; i++ {
			r = cal2(inc(f), r, atx(inc(x), mki(i)))
		}
		return decr2(x, f, r)
	}
}
func ovi(f, x, y k) (r k) { // x f/y
	xn, yn := m.k[x]&atom, m.k[y]&atom
	if yn == atom {
		panic("class")
	}
	if x, y, t, op := scop2(f, x, y); op != nil {
		r := mk(t, xn)
		xn, yn = atm1(xn), atm1(yn)
		rp, xp, yp, cp := ptr(r, t), ptr(x, t), ptr(y, t), cpx[t]
		for i := k(0); i < xn; i++ {
			cp(rp, xp+i)
			for j := k(0); j < yn; j++ {
				op(rp, rp, yp+j)
			}
			rp++
		}
		return decr2(x, y, r)
	}
	r = x
	for i := k(0); i < yn; i++ {
		r = cal2(inc(f), r, atx(inc(y), mki(i)))
	}
	return decr2(y, f, r)
}
func sci(f, x, y k) (r k) { // x f\y
	xn, yn := m.k[x]&atom, m.k[y]&atom
	if yn == atom {
		panic("class")
	}
	if xn == atom {
		if x, y, t, op := scop2(f, x, y); op != nil {
			r = mk(t, yn)
			xp, yp, rp := ptr(x, t), ptr(y, t), ptr(r, t)
			op(rp, xp, yp)
			for i := k(1); i<yn; i++ {
				op(rp+i, rp+i-1, yp+i)
			}
			return decr2(x, y, r)
		}
	}
	r = mk(L, yn)
	for i := k(0); i < yn; i++ {
		x = cal2(inc(f), x, atx(inc(y), mki(i)))
		m.k[2+i+r] = inc(x)
	}
	dec(f)
	return decr2(x, y, uf(r))
}
func scop1(f, x k) (k, k, f2) {
	if m.k[f]&atom != atom || m.k[2+f] > 255 {
		return x, 0, nil
	}
	t := m.k[x] >> 28
	if t >= L {
		return x, 0, nil
	}
	opx, _ := op2(m.k[2+f])
	if opx == nil {
		return x, 0, nil
	}
	dec(f)
	op := opx[t]
	if t == C && opx[C] == nil {
		t = I
	}
	if (t == I && opx[I] == nil) || (t == Z && opx[Z] == nil) {
		t = F
	}
	x = to(x, t)
	return to(x, t), t, op
}
func scop2(f, x, y k) (k, k, k, f2) {
	if m.k[f]&atom != atom || m.k[2+f] > 255 {
		return x, y, 0, nil
	}
	xt, yt := m.k[x]>>28, m.k[y]>>28
	if xt >= L || yt >= L {
		return x, y, 0, nil
	}
	opx, _ := op2(m.k[2+f])
	if opx == nil {
		return x, y, 0, nil
	}
	dec(f)
	t := ntyps(xt, yt, opx, nil)
	op := opx[t]
	if xt < t {
		x, xt = to(x, t), t
	}
	if yt < t {
		y, yt = to(y, t), t
	}
	return x, y, t, op
}
func bin(x, y k) (r k) { // x bin y
	xt, yt, xn, yn := typs(x, y)
	if xt != yt || xt > S || xn == atom {
		panic("type")
	}
	r = mk(I, yn)
	yn = atm1(yn)
	lt, xp, yp := ltx[yt], ptr(x, xt), ptr(y, yt)
	for i := k(0); i < yn; i++ {
		m.k[2+i+r] = ibin(xp, xt, xn, yp+i, lt)
	}
	return decr2(x, y, r)
}
func ibin(xp, t, n, yp k, lt func(x, y k) bool) (r k) {
	if n == 0 {
		return 0
	}
	if lt(yp, xp) {
		n := i(-1)
		return k(n)
	}
	i, j, h := k(0), k(n), k(0)
	for i < j {
		h = (i + j) >> 1
		if lt(xp+h, yp) {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}
func insert(x, y, idx k) (r k) { // insert y into x at k
	xt, yt, xn, yn := typs(x, y)
	if xt > L || xn == atom || (xt != L && (xt != yt || yn != atom)) {
		panic("type")
	}
	if xt == L {
		x = lcat(x, y)
	} else {
		x = ucat(x, y, xt, xn, yn)
	}
	if idx != xn {
		sw, xp := swx[xt], ptr(x, xt)
		sw(xp+idx, xp+idx+xn)
	}
	return x
}
func unsert(x, idx k) (r k) { // delete index from x
	t, n := typ(x)
	if t == atom || n > L || idx >= n {
		panic("type")
	}
	cp, xp := cpx[t], ptr(x, t)
	if m.k[x+1] == 1 && bk(t, n-1) == bk(t, n) {
		if t == L {
			cp = cpI
			dec(m.k[2+x+idx])
		}
		for i := k(idx); i < n-1; i++ {
			cp(xp+i, xp+i+1)
		}
		m.k[x] = t<<28 | (n - 1)
		return x
	}
	r = mk(t, n-1)
	rp := ptr(r, t)
	for i := k(0); i < idx; i++ {
		cp(rp+i, xp+i)
	}
	for i := k(idx) + 1; i < n; i++ {
		cp(rp+i-1, xp+i)
	}
	return decr(x, r)
}
func asn(x, y k) (r k) { // `x:y
	keys, vals := m.k[kkey], m.k[kval]
	if ix, exists := varn(ptr(x, S)); exists {
		dec(m.k[2+vals+ix])
		m.k[2+vals+ix] = inc(y)
		return decr(x, y)
	} else {
		m.k[kkey] = insert(keys, x, ix)
		m.k[kval] = insert(vals, inc(y), ix)
		return y
	}
}
func lup(x k) (r k) { // lookup
	r = lupo(x)
	if r == 0 {
		panic("undefined")
	}
	return r
}
func lupo(x k) (r k) { // lup, 0 on undefined
	ix, o := varn(ptr(x, S))
	if !o {
		return decr(x, 0)
	}
	vals := m.k[kval]
	r = inc(m.k[2+vals+ix])
	return decr(x, r)
}
func varn(xp k) (idx k, exists bool) {
	keys := m.k[kkey]
	kp := ptr(keys, S)
	kn := m.k[keys] & atom
	ix := ibin(kp, S, kn, xp, ltS)
	if i(ix) < 0 {
		ix = 0
	}
	return ix, ix < kn && eqS(kp+ix, xp)
}
func vars(dummy k) (r k) { dec(dummy); return inc(m.k[kkey]) }
func del(x k) (r k) { // delete variable
	t, n := typ(x)
	if t != S {
		panic("type")
	}
	n = atm1(n)
	xp := ptr(x, S)
	for i := k(0); i < n; i++ {
		if idx, o := varn(xp + i); o {
			m.k[kkey] = unsert(m.k[kkey], idx)
			m.k[kval] = unsert(m.k[kval], idx)
		}
	}
	return decr(x, mk(N, atom))
}
func clear() { // clear variables
	dec(m.k[kkey])
	dec(m.k[kval]) // m.k[kkey] and m.k[kval] should not be free blocks
	m.k[kkey] = mk(S, 0)
	m.k[kval] = mk(L, 0)
}
func lsv(x k) (r k) { dec(x); return inc(m.k[kkey]) }       // \v (list variables)
func clv(x k) (r k) { dec(x); clear(); return mk(N, atom) } // \c (clear variables)
func hlp(x k) (r k) {
	dec(x)
	n := k(168 - 136)
	r = mk(C, n)
	rc := 8 + r<<2
	copy(m.c[rc:rc+n], m.c[136:168])
	return cat(r, kst(inc(m.k[2+m.k[3]])))
}

func hxb(x c) (c, c) { h := "0123456780abcdef"; return h[x>>4], h[x&0x0F] }
func hxk(x k) s {
	b := []c{'0', 'x', '0', '0', '0', '0', '0', '0', '0', '0'}
	for j := k(0); j < 4; j++ {
		n := 8 * (3 - j)
		b[2+2*j], b[3+2*j] = hxb(c((x & (0xFF << n)) >> n))
	}
	return s(b)
}

func isnan(x f) bool { return x != x }
func atm1(n k) k {
	if n == atom {
		return 1
	}
	return n
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

var table [100]interface{} // function table :+-*%&|<>=!~,^#_$?@.0123456789'/\
func init() {
	table = [100]interface{}{
		//   1                   5                        10                       15
		idn, flp, neg, fst, inv, wer, rev, asc, dsc, grp, til, not, enl, srt, cnt, flr, str, unq, tip, evl,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, qtc, slc, bsc, ech, ovr, scn, ecp, nil, nil,
		lsv, clv, nil, nil, nil, hlp, nil, nil, nil, nil, nil,
		nil, add, sub, mul, div, min, max, les, mor, eql, key, mch, cat, ept, tak, drp, cst, fnd, atx, cal,
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, qot, sla, bsl, ecd, ovi, sci, epi, ecr, ecl,
		nil, nil, bin, nil, del, nil, nil, nil, nil, nil, nil,
	}
}
