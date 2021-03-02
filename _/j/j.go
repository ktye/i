package j

import (
	"fmt"
	"j/x"
	jx "j/x"
	"math/bits"
)

var M []uint32 // heap
var F []func() // function table

func J(x uint32) uint32 {
	if M == nil {
		return ii(x)
	}
	if 1 == lu(C) {
		if 41 == x { // ) end comment
			as(C, 0)
		}
		return 0
	}
	if 40 == x { // ( start comment
		as(C, 1)
		return 0
	}
	n := lu(N)
	if 47 < x { // x >= '0'
		if 58 > x { // x <= '9' parse number
			x -= '0'
			if x|n == 0 {
				pc(1) // number 0
				return 0
			}
			n *= 10
			n += x
			as(N, n)
			return 0
		}
	}
	if n != 0 {
		pc(1 | n<<1)
		as(N, 0)
	}
	y := lu(Y)
	if 96 < x { // x >= 'a'
		if 123 > x { // x <= 'z' parse symbol
			y *= 32
			y += x - '`' // 96
			as(Y, y)
			return 0
		}
	}
	if y != 0 {
		pc(2 | y<<2)
		as(Y, 0)
	}
	if 33 > x {
		if 10 == x {
			exe()
			push(mk(0))
			return 1
		}
	} else if 91 == x { // [
		push(mk(0))
	} else if 93 == x { // ]
		pc(lp())
	} else {
		pc(4 | (x-33)<<3) // operator
	}
	return 0
}
func exe() { // [q]. exec
	x := lp()
	p := x + 8
	l := x
	if nn(x) != 0 {
		l = pl(x)
	}
	for p <= l {
		c := I(p)
		if 2 == c&3 { // symbol
			push(rx(lu(c)))
			exe()
		} else if 4 != c&7 { // list or number
			push(rx(c))
		} else if 740 == c { // } push to swap
			pp()
			sw()
		} else if 724 == c { // { pop from swap
			sw()
			pp()
		} else {
			F[c>>3]()
		}
		p += 4
	}
	dx(x)
}
func pp() {
	t := pop()
	sw()
	push(t)
}
func sw() { // swap stacks
	t := I(4)
	sI(4, I(8))
	sI(8, t)
}
func ife() { // c[t][e]?  (if c then t else e)
	e := pop()
	t := pop()
	if ip() == 0 {
		dx(t)
		push(e)
	} else {
		dx(e)
		push(t)
	}
	exe()
}
func I(x uint32) uint32 { return M[x>>2] }
func sI(x, y uint32)    { M[x>>2] = y }
func bk(n uint32) (r uint32) { // bucket type
	r = uint32(32 - bits.LeadingZeros32(7+4*n))
	if r < 4 {
		return 4
	}
	return r
}
func mk(x uint32) (r uint32) { // allocate
	t := bk(x)
	i := 4 * t
	m := 4 * I(0)
	for I(i) == 0 {
		if i >= m {
			panic("memory")
		}
		i += 4
	}
	a := I(i)
	sI(i, I(a))
	for j := i - 4; j >= 4*t; j -= 4 {
		u := a + 1<<(j>>2)
		sI(u, I(j))
		sI(j, u)
	}
	sI(a, 1)
	sI(a+4, x)
	return a
}
func rx(x uint32) uint32 {
	if x&7 == 0 {
		sI(x, I(x)+1)
	}
	return x
}
func dx(x uint32) {
	if x != 0 && x&7 == 0 {
		if I(x) == 0 {
			panic("dx on free")
		}
		r := I(x) - 1
		sI(x, r)
		if r == 0 {
			n := I(x + 4)
			p := x + 8
			for i := uint32(0); i < n; i++ {
				dx(I(p))
				p += 4
			}
			fr(x)
		}
	}
}
func fr(x uint32) {
	p := 4 * bk(I(4+x))
	sI(x, I(p))
	sI(p, x)
}
func nn(x uint32) uint32 { return I(4 + x) }
func lc(x uint32, y uint32) (r uint32) { // list cat (append a single element)
	n := nn(x)
	if I(x) == 1 && bk(n) == bk(1+n) {
		sI(8+x+4*n, y)
		sI(4+x, 1+n)
		return x
	}
	r = mk(1 + n)
	sI(cp(x, r, n), y)
	dx(x)
	return r
}
func cp(x, r, n uint32) (rp uint32) {
	x += 8
	r += 8
	for i := uint32(0); i < n; i++ {
		sI(r, rx(I(x)))
		r += 4
		x += 4
	}
	return r
}
func pc(x uint32) { push(lc(lp(), x)) } // cat to top list
func pl(x uint32) uint32 { // pointer of last element in list
	n := nn(x)
	if n == 0 {
		panic("empty")
	}
	return 4 + x + 4*n
}
func fi(x uint32) (r uint32) {
	if nn(x) == 0 {
		panic("empty")
	}
	r = rx(I(x + 8))
	dx(x)
	return r
}
func us(x uint32) uint32 {
	if I(x) == 1 {
		//fmt.Println("reuse")
		return x
	}
	n := nn(x)
	r := mk(n)
	cp(x, r, n)
	dx(x)
	return r
}
func ip() int32 {
	x := pop()
	if x&1 == 0 {
		panic("int expected")
	}
	return int32(x) >> 1
}
func lp() uint32 {
	x := pop()
	if x&7 != 0 {
		panic("list expected")
	}
	return x
}
func add()       { pi(ip() + ip()) }
func sub()       { pi(-ip() + ip()) }
func mul()       { pi(ip() * ip()) }
func div()       { swp(); pi(ip() / ip()) }
func mod()       { swp(); pi(ip() % ip()) }
func eql()       { pb(ip() == ip()) }
func gti()       { pb(ip() < ip()) }
func lti()       { pb(ip() > ip()) }
func pi(i int32) { push(1 + 2*uint32(i)) }
func pb(b bool) {
	if b {
		push(3)
	} else {
		push(1)
	}
}
func stk() { fmt.Println(x.X(M, I(4)) + " -- " + x.X(M, I(8))) } // !
func dup() { x := pop(); push(rx(x)); push(x) }                  // "
func drp() { dx(pop()) }                                         // x _ -- (pop)
func swp() { // ~
	x := pop()
	y := pop()
	push(x)
	push(y)
}
func rol() { // |
	a := pop()
	b := pop()
	c := pop()
	push(a)
	push(c)
	push(b)
}
func cnt() { // #
	x := pop()
	r := uint32(0xffffffff)
	if x&7 == 0 {
		r = 1 + 2*nn(x)
	}
	dx(x)
	push(r)
}
func amd() { // [a]i v$ amend (set array at index i to v)
	v := pop()
	i := ip()
	a := us(lp())
	n := int32(nn(a))
	if i == n {
		a = lc(a, v)
	} else if i < 0 || i > n {
		panic("amd: range")
	} else {
		ap := 8 + a + 4*uint32(i)
		rx(I(ap))
		sI(ap, v)
	}
	push(a)
}
func atx() { // [..]i@
	i := ip()
	l := lp()
	if i < 0 || i >= int32(nn(l)) {
		panic("atx: range")
	}
	push(rx(I(8 + 4*uint32(i) + l)))
	dx(l)
}
func cat() { // ,
	y := pop()
	x := pop()
	if x&7 != 0 {
		x = lc(mk(0), x)
	}
	if y&7 != 0 {
		x = lc(x, y)
	} else {
		yp := y + 8
		for i := uint32(0); i < nn(y); i++ {
			x = lc(x, rx(I(yp)))
			yp += 4
		}
		dx(y)
	}
	push(x)
}
func asn() { // [q][s]: -- (assign)
	y := fi(lp())
	if y&3 != 2 {
		panic("asn: not a symbol")
	}
	v := pop()
	if v&7 != 0 {
		v = lc(mk(0), v) // enlist atoms
	}
	p := ps(y)
	dx(I(p))
	sI(p, v)
}
func ps(x uint32) uint32 {
	s := I(12)
	p := fn(x)
	if p == 0 {
		s = lc(s, x)
		s = lc(s, 1)
		p = pl(s)
	}
	sI(12, s)
	return p
}
func as(x, y uint32) { // no dx
	p := ps(x)
	sI(p, y)
}
func drw() { dx(lp()); dx(lp()) }
func lu(x uint32) uint32 {
	p := fn(x)
	if p == 0 {
		panic("undefined: " + jx.X(M, x))
	}
	return I(p)
}
func fn(x uint32) uint32 {
	s := I(12)
	n := nn(s) / 2
	p := s + 8
	for i := uint32(0); i < n; i++ {
		if I(p) == x {
			return p + 4
		}
		p += 8
	}
	return 0
}
func push(x uint32) {
	s := I(4)
	n := nn(s)
	if n == jx.SZ {
		panic("stack overflow")
	}
	sI(s+4, 1+n)
	sI(pl(s), x)
}
func pop() (r uint32) {
	s := I(4)
	n := nn(s)
	if n == 0 {
		panic("stack underflow")
	}
	r = I(pl(s))
	sI(s+4, n-1)
	return r
}
func finit() {
	f := func(c byte, g func()) { F[c-33] = g }
	F = make([]func(), 128)
	f('!', stk) // 0
	f('"', dup) // 1
	f('#', cnt) // 2
	f('$', amd) // 3
	f('%', mod) // 4
	f('&', drw) // 5
	f('*', mul) // 9
	f('+', add) // 10
	f(',', cat) // 11
	f('-', sub) // 12
	f('.', exe) // 13
	f('/', div) // 14
	f(':', asn) // 25
	f('<', lti) // 27
	f('=', eql) // 28
	f('>', gti) // 29
	f('?', ife) // 30
	f('@', atx) // 31
	f('_', drp) // 62
	f('|', rol) // 91
	f('~', swp) // 93
}
func ii(x uint32) uint32 { // ini(16): 64kB
	finit()
	M = make([]uint32, 1<<(x-2))
	M[0] = uint32(x)
	p := uint32(128)
	for i := uint32(7); i < x; i++ {
		sI(4*i, p) // free pointer
		p *= 2
	}
	sI(4, mk(jx.SZ)) // stack
	sI(8, mk(jx.SZ)) // swap stack
	sI(4+I(4), 0)
	sI(4+I(8), 0)
	push(mk(0))
	sI(12, mk(0)) // key/value list

	// parse state is stored in hidden symbols
	// values 110 114 118 are symbols but cannot be entered
	as(N, 0) // number
	as(Y, 0) // symbol
	as(C, 0) // comment
	return x
}

var N, Y, C uint32 = jx.N, jx.Y, jx.C // 110,114,118

//const sz uint32 = 126
