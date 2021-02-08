//  go:embed j.j
//  var j []byte
package j

import (
	_ "embed"
	"fmt"
	"math/bits"
)

//go:generate go run h.go

var N uint32                // number
var Y uint32                // symbol
var P uint32                // parse root list
var T uint32                // parse top list
var C uint32                // parse comment
var M []uint32              // heap
var F []func(uint32) uint32 // function table
func ini() {
	N, Y, P, T, C = 0, 0, 0, 0, 0
	finit()
	x := uint32(16)
	M = make([]uint32, 1<<(x-2)) // 64kB
	M[0] = uint32(x)
	p := uint32(128)
	for i := uint32(7); i < x; i++ {
		sI(4*i, p) // free pointer
		p *= 2
	}
	M[1] = mk(8) // max stack size
	M[2] = mk(0)
	M[3] = mk(0)
	P = mk(0)
	T = P
	//dump(127)
}

func J(x uint32) uint32 {
	if C == 1 {
		if x == ')' { // end comment
			C = 0
		}
		return 0
	}
	if x == '(' { // start comment
		C = 1
		return 0
	}
	if x >= '0' && x <= '9' { // parse number
		x -= '0'
		if x|N == 0 {
			T = pcat(T, 1)
			return 0
		}
		N *= 10
		N += x
		return 0
	}
	if N != 0 {
		T = pcat(T, 1|N<<1)
		N = 0
	}
	if x >= 'a' && x <= 'z' { // parse symbol
		Y *= 32
		Y += x - '`'
		return 0
	}
	if Y != 0 {
		T = pcat(T, 2|Y<<2)
		Y = 0
	}
	if x < 33 {
		if x == 10 {
			if T != P {
				panic("unclosed]")
			}
			v := Exec(T)
			P = mk(0)
			T = P
			return v
		}
		return 0
	}
	if x == 91 { // '['
		T = pcat(T, mk(0))
		T = last(T)
		return 0
	}
	if x == 93 {
		T = parent(P, T)
		if T == 0 {
			panic("parse]")
		}
		return 0
	}
	T = pcat(T, 4|(x-33)<<3)
	return 0
}
func Exec(stk, q uint32) uint32 {
	if nn(q) == 0 {
		dx(q)
		return stk
	} else if q&7 != 0 {
		panic("exec: not a quotation")
	}
	//fmt.Println("rc stk", I(stk), "q", I(q))
	r := uint32(0)
	p := q + 8
	l := lastp(q)
	tailcall := func(x uint32) {
		fmt.Println("tailcall")
		dx(q)
		q = x
		p = q + 4
		l = lastp(q)
	}
	for p <= l {
		x := I(p)
		if tail := p == l; tail && x&3 == 2 { // symbol
			tailcall(lup(x))
		} else if tail && x == 93 { // "."
			t := rx(last(stk))
			stk = pop(stk)
			tailcall(t)
		} else if tail && x == 127 { // "?" branch last
			e := rx(last(stk))
			stk = pop(stk)
			t := rx(last(stk))
			stk = pop(stk)
			c := lasti(stk)
			stk = pop(stk)
			if c == 0 {
				dx(t)
				tailcall(e)
			} else {
				dx(e)
				tailcall(t)
			}
		} else if x&3 == 2 { // symbol
			stk = Exec(stk, lup((x)))
		} else if x&7 != 4 { // no operator but (list or number)
			stk = lcat(stk, rx(x))
		} else if x == 740 { // } push to reg
			if r == 0 {
				r = mk(0)
			}
			r = lcat(r, rx(last(stk)))
			stk = pop(stk)
		} else if x == 724 { // { pop from reg
			x := rx(last(r))
			r = pop(r)
			stk = lcat(stk, x)
		} else {
			//fmt.Println(x)
			stk = F[x>>3](stk)
		}
		p += 4
	}
	dx(q)
	dx(r)
	return stk
}
func exe(s uint32) uint32 { // [q]. exec
	q := rx(last(s))
	return Exec(pop(s), q)
}
func ife(s uint32) uint32 { // c[t][e]?  (if c then t else e)
	e := rx(last(s))
	s = pop(s)
	t := rx(last(s))
	s = pop(s)
	c := lasti(s)
	if c == 0 {
		dx(t)
		return Exec(pop(s), e)
	} else {
		dx(e)
		return Exec(pop(s), t)
	}
}
func whl(s uint32) uint32 { // [c][d]' (while c do d)
	panic("while")
	d := rx(last(s))
	s = pop(s)
	c := rx(last(s))
	s = pop(s)
	for {
		s = Exec(s, rx(c))
		i := lasti(s)
		// fmt.Println("whl", X(s))
		s = pop(s)
		if i == 0 {
			dx(c)
			dx(d)
			return s
		} else {
			s = Exec(s, rx(d))
		}
	}
}

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
	m := 4 * M[0]
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
func dx(x uint32) uint32 {
	if x != 0 && x&7 == 0 {
		if I(x) == 0 {
			panic("dx on free")
		}
		sI(x, I(x)-1)
		if I(x) == 0 {
			n := I(x + 4)
			p := x + 8
			for i := uint32(0); i < n; i++ {
				dx(I(p))
				p += 4
			}
			fr(x)
		}
	}
	return x
}
func fr(x uint32) {
	p := 4 * bk(I(4+x))
	sI(x, I(p))
	sI(p, x)
}
func nn(x uint32) uint32 { return I(4 + x) }
func cat(s uint32) uint32 {
	y := rx(last(s))
	s = pop(s)
	p := lastp(s)
	x := I(p)
	if x&7 != 0 {
		x = lcat(mk(0), x)
	}
	if y&7 != 0 {
		x = lcat(x, y)
	} else {
		yp := y + 8
		for i := uint32(0); i < nn(y); i++ {
			x = lcat(x, rx(I(yp)))
			yp += 4
		}
		dx(y)
	}
	sI(p, x)
	return s
}
func lcat(x uint32, y uint32) (r uint32) {
	n := nn(x)
	r = mk(1 + n)
	xp, rp := x+8, r+8
	for i := uint32(0); i < n; i++ {
		sI(rp, rx(I(xp)))
		rp += 4
		xp += 4
	}
	sI(rp, y)
	dx(x)
	return r
}
func pcat(x, y uint32) (r uint32) {
	p := parent(P, x)
	r = lcat(x, y)
	if x == P {
		P = r
		return r
	}
	sI(lastp(p), r)
	return r
}
func lastp(x uint32) uint32 {
	n := nn(x)
	if n == 0 {
		panic("empty")
	}
	return 4 + x + 4*n
}
func lasti(x uint32) (r uint32) {
	r = last(x)
	if r&1 == 0 {
		panic("int expected")
	}
	return r >> 1
}
func last(x uint32) (r uint32) {
	n := nn(x)
	if n == 0 {
		return 0
	}
	return I(lastp(x))
}
func prev(x uint32) (r uint32) {
	n := nn(x)
	if n < 2 {
		panic("prev: underflow")
	}
	return I(x + 4*n)
}
func last2(x uint32) (a, b uint32) {
	n := nn(x)
	if n < 2 {
		panic("stack-underflow")
	}
	x += 4 * n
	return I(x), I(x + 4)
}
func parent(x, y uint32) (r uint32) {
	if x&7 != 0 {
		panic("parent")
	}
	l := last(x)
	if l == y || l == 0 || x == y {
		return x
	}
	return parent(l, y)
}

func I(x uint32) uint32 { return M[x>>2] }
func sI(x, y uint32)    { M[x>>2] = y }

func stk(s uint32) uint32 { // !
	fmt.Println("(stk) " + X(s))
	return s
}

func swp(s uint32) uint32 { // ~
	x := lastp(s)
	if x < s+12 {
		panic("swp underflow")
	}
	t := I(x)
	sI(x, I(x-4))
	sI(x-4, t)
	return s
}
func dup(s uint32) uint32 { p := last(s); return lcat(s, rx(p)) }
func rol(s uint32) uint32 {
	p := lastp(s)
	if p < s+16 {
		panic("rol underflow")
	}
	a := I(p)
	sI(p, I(p-4))
	sI(p-4, I(p-8))
	sI(p-8, a)
	return s
}
func cnt(s uint32) uint32 {
	x := last(s)
	r := uint32(0xffffffff)
	if x&7 == 0 {
		r = 1 + 2*nn(x)
	}
	return v1(s, r)
}
func use(x uint32) uint32 {
	if I(x) == 1 {
		return x
	}
	n := nn(x)
	r := mk(n)
	rp := r + 8
	xp := x + 8
	for i := uint32(0); i < n; i++ {
		sI(rp, rx(I(xp)))
		xp += 4
		rp += 4
	}
	dx(x)
	return r
}
func amd(s uint32) uint32 { // [a]i v$ amend (set array at index i to v)
	v := rx(last(s))
	s = pop(s)
	i := lasti(s)
	s = pop(s)
	p := lastp(s)
	a := use(I(p))
	n := nn(a)
	if i == n {
		a = lcat(a, v)
	} else if i < 0 || i > n {
		panic("amd: range")
	} else {
		ap := 8 + a + 4*i
		rx(I(ap))
		sI(ap, v)
	}
	sI(p, a)
	return s
}
func v1(s, x uint32) uint32 {
	sp := s + 4 + 4*nn(s)
	dx(I(sp))
	sI(sp, x)
	return s
}
func ints(s uint32) (j, k int32) {
	b := pop()
	a := pop()
	if a&1 == 0 || b&1 == 0 {
		panic("ints")
	}
	return int32(a) >> 1, int32(b) >> 1
}
func add() { push(pop() + pop()) }
func sub() { push(pop() - pop()) }
func mul() { push(pop() * pop()) }
func div() { push(pop() / pop()) }
func mod() { push(pop() % pop()) }
func eql() { push(ib(pop() == pop())) }
func gti() { push(ib(pop() > pop())) }
func lti() { push(ib(pop() < pop())) }
func max() {
	a, b := ints(s)
	if a > b {
		return i2(s, a)
	}
	return i2(s, b)
}
func min(s uint32) uint32 {
	a, b := ints(s)
	if a < b {
		return i2(s, a)
	}
	return i2(s, b)
}
func ib(b bool) int32 {
	if b {
		return 1
	}
	return 0
}
func pop() (r uint32) {
	s := I(4)
	p := I(s + 12)
	if p == s+12 {
		panic("stack underflow")
	}
	r = I(p)
	sI(p, 0)
	sI(s+12, p-4)
	return r
}
func push(x uint32) {
	s := I(4)
	p := I(s+12) + 4
	if p == s+4*nn(s) {
		panic("stack overflow")
	}
	sI(p, x)
	sI(s+12, p)
}
func drp(x uint32) { pop(x) } // x _ -- (pop)
func asn(s uint32) uint32 { // [q][s]: -- (assign)
	y := pop(s)
	if y&3 != 2 {
		panic("asn: not a symbol")
	}
	v := pop(s)
	if v&7 != 0 {
		v = lcat(mk(0), v) // enlist atoms
	}
	p := fns(I(8), y)
	if p == 0 {
		sI(8, lcat(I(8), y))
		sI(12, lcat(I(12), 1))
		p = 4 + 4*nn(I(8))
	}
	dx(I(I(12) + p))
	sI(I(12)+p, v)
	return pop(s)
}
func lup(x uint32) uint32 {
	p := fns(I(8), x)
	if p == 0 {
		panic("undefined: " + X(x))
	}
	return rx(I(I(12) + p))
}
func fns(x, y uint32) uint32 {
	n := nn(x)
	p := uint32(8)
	for i := uint32(0); i < n; i++ {
		if I(x+p) == y {
			return p
		}
		p += 4
	}
	return 0
}
func atx(s uint32) uint32 { // [..]i@
	l := prev(s)
	if l&7 != 0 {
		panic("atx: not a list")
	}
	i := lasti(s)
	s = pop(s)
	if i < 0 || i >= nn(l) {
		panic("atx: range")
	}
	sI(lastp(s), rx(I(8+4*i+l)))
	dx(l)
	return s
}

func finit() {
	f := func(c byte, g func(uint32) uint32) { F[c-33] = g }
	F = make([]func(uint32) uint32, 128)
	f('!', stk) // 0
	f('"', dup) // 1
	f('#', cnt) // 2
	f('$', amd) // 3
	f('%', mod) // 4
	f('&', min) // 5
	f(3_9, whl) // 6 `
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
	f('^', max) // 61
	f('_', drp) // 62
	f('|', rol) // 91
	f('~', swp) // 93
}
func init() { ini() }
