package w

type c = byte
type k = uint32
type i = int32
type f = float64
type z = complex128
type s = string

const (
	C, I, F, Z, S, G, L, D byte = 1, 2, 3, 4, 5, 6, 7, 8
)

type (
	fc1 func(c) c
	fi1 func(i) i
	ff1 func(f) f
	fz1 func(z) z
)

//                c  i  f   z  s  g  l  d
var lns = [9]k{0, 1, 4, 8, 16, 4, 4, 4, 8}
var e k = 0xFFFFFFFF

var m []c

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
		return e
	}
	fb, a := 0, k(0)
	for i := bt; i < 30; i++ { // find next free bucket >= bt
		if k(i) >= get(4) {
			grw()
		}
		if get(k(4*i)) != 0 {
			fb, a = i, get(k(4*i))
			break
		}
	}
	if fb == 0 {
		return e
	}
	for i := fb - 1; i >= bt; i-- { // split large buckets
		put(k(4*i), a)
		m[a] = c(i)
		a += k(1) << c(i)
		m[a] = c(i)
	}
	if n < 0 { // set header
		m[int(a+1)] = t
	} else {
		put(k(a), k(m[int(a)]|t<<5)|k(n)<<8)
	}
	put(a+4, 1) // refcount
	return a
}
func typ(a k) (c, int) { // type and length at addr
	i := int(a)
	t := m[i] >> 5
	if t == 0 {
		return m[int(i+1)], -1
	}
	return t, int(get(k(i)) >> 8)
}
func ini() { // start function
	m = make([]c, 1<<16)
	p := k(len(m))
	for i := 15; i > 6; i-- {
		p >>= 1
		m[p] = c(i)
		put(k(4*i), p)
	}
	m[0] = 7
	put(4, 16) // total memory (log2)
	// TODO: pointer to k-tree at 8
	put(k(4*9), 0)   // no free bucket 9
	put(1<<9, k(73)) // 73: 1<<6|9 (type i, bucket 9), length is ignored
	for i := range lns {
		put(k(4*i+8)+1<<9, k(lns[i]))
	}
}
func put(a, x k) {
	i := int(a)
	m[i] = c(x)
	m[i+1] = c(x >> 8)
	m[i+2] = c(x >> 16)
	m[i+3] = c(x >> 24)
}
func putf(a k, x f) {
	panic("TODO putf")
}
func get(a k) k  { i := int(a); return k(m[i]) | k(m[i+1])<<8 | k(m[i+2])<<16 | k(m[i+3])<<24 }
func getf(a k) f { panic("TODO: getf"); return 0.0 }
func grw() {
	s := m[4]
	if 1<<k(s) != len(m) {
		panic("grw")
	}
	put(k(4*s), k(len(m)))
	m = append(m, make([]c, len(m))...)
	m[4] = s + 1
	m[1<<s] = s
}
func inc(x k) k { put(x+4, get(x+4)+1); return x }
func dec(x k) {
	rc := get(x + 4)
	rc--
	if rc == 0 {
		panic("TODO free")
	}
	put(x+4, rc)
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
		g = func(x, y k) { put(y, k(i(m[int(x)]))) }
	case t == C && rt == F:
		g = func(x, y k) { putf(y, f(m[int(x)])) }
	case t == I && rt == C:
		g = func(x, y k) { m[int(y)] = c(i(get(x))) }
	case t == I && rt == F:
		g = func(x, y k) { putf(y, f(i(get(x)))) }
	case t == F && rt == C:
		g = func(x, y k) { m[int(y)] = c(getf(x)) }
	case t == F && rt == I:
		g = func(x, y k) { put(y, k(i(getf(x)))) }
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
func cl1(x, r k, n k, op fc1) {
	o := r - x
	for j := 8 + x; j < 8+x+n; j++ {
		m[j] = op(m[o+j])
	}
}
func il1(x, r k, n k, op fi1) {
	o := r - x
	for j := 8 + x; j < 8+x+4*n; j += 4 {
		put(o+j, k(op(i(get(j)))))
	}
}
func fl1(x, r k, n k, op ff1) {
	o := r - x
	for j := 8 + x; j < 8+x+8*n; j += 8 {
		putf(o+j, op(getf(j)))
	}
}
func il2(x, y, r k, n k, op func(i, i) i) {
	for j := k(0); j < n; j++ {
		u := get(8 + x + 4*j)
		v := get(8 + y + 4*j)
		put(8+r+4*j, k(op(i(u), i(v))))
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
		x, t = to(x, min), min
	}
	if t == Z && fz == nil { // e.g. real functions
		x, t = to(x, F), F
	}
	if get(x+4) == 1 {
		r = inc(x)
	} else {
		r = mk(t, n)
	}
	if n < 0 {
		n = 1
	}
	switch t {
	case L:
		for j := k(0); j < k(n); j++ {
			put(r+8+4*j, nm(inc(get(x+8+j*8)), fc, fi, ff, fz, rt))
		}
	case D:
		if r != x {
			put(r+8, inc(get(x+8)))
		}
		put(r+12, nm(get(x+12), fc, fi, ff, fz, rt))
	case C:
		cl1(x, r, k(n), fc)
	case I:
		il1(x, r, k(n), fi)
	case F:
		fl1(x, r, k(n), ff)
	// TODO Z
	default:
		panic("type")
	}
	dec(x)
	if r != 0 && t > rt {
		r = to(r, rt) // downtype, e.g. floor
	}
	return r
}
func neg(x k) k {
	return nm(x, func(x c) c { return -x }, func(x i) i { return -x }, func(x f) f { return -x }, nil, 0) // TODO Z
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
