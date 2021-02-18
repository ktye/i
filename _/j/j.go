//  go:embed j.j
//  var j []byte
package j

import (
	_ "embed"
	"fmt"
	"math/bits"
)

//go:generate go run h.go

var M []uint32 // heap
var F []func() // function table
func ini() {
	finit()
	x := uint32(16)
	M = make([]uint32, 1<<(x-2)) // 64kB
	M[0] = uint32(x)
	p := uint32(128)
	for i := uint32(7); i < x; i++ {
		sI(4*i, p) // free pointer
		p *= 2
	}
	sI(4, mk(0))     // stack
	sI(12, mk(0))    // value list
	s := mk(5)       // parse state
	sI(s+8, mk(0))   // p(root list)
	sI(s+12, I(s+8)) // t(top list)
	sI(s+16, 0)      // n(number)
	sI(s+20, 0)      // y(symbol)
	sI(s+24, 0)      // c(comment)
	sI(8, s)
}

func J(x uint32) uint32 {
	s := I(8)
	P, T := I(s+8), I(s+12)
	if I(s+24) == 1 {
		if x == ')' { // end comment
			sI(s+24, 0)
		}
		return 0
	}
	if x == '(' { // start comment
		sI(s+24, 1)
		return 0
	}
	n := I(s + 16)
	if x >= '0' && x <= '9' { // parse number
		x -= '0'
		if x|n == 0 {
			T = pcat(1)
			return 0
		}
		n *= 10
		n += x
		sI(s+16, n)
		return 0
	}
	if n != 0 {
		T = pcat(1 | n<<1)
		sI(s+16, 0)
	}
	y := I(s + 20)
	if x >= 'a' && x <= 'z' { // parse symbol
		y *= 32
		y += x - '`'
		sI(s+20, y)
		return 0
	}
	if y != 0 {
		T = pcat(2 | y<<2)
		sI(s+20, 0)
	}
	if x < 33 {
		if x == 10 {
			if T != P {
				panic("unclosed]")
			}
			Exec(P)
			sI(s+8, mk(0))
			sI(s+12, I(s+8))
			return 1
		}
		return 0
	}
	if x == 91 { // '['
		T = pcat(mk(0))
		sI(s+12, last(T))
		return 0
	}
	if x == 93 {
		T = parent(P, T)
		if T == 0 {
			panic("parse]")
		}
		sI(s+12, T)
		return 0
	}
	T = pcat(4 | (x-33)<<3)
	return 0
}
func Exec(q uint32) {
	//fmt.Println("Exec", q, XX(q))
	if nn(q) == 0 {
		dx(q)
		return
	}
	if q&7 != 0 {
		panic("exec: not a quotation")
	}
	//fmt.Println("rc stk", I(stk), "q", I(q))
	r := uint32(0)
	p := q + 8
	l := lastp(q)
	tailcall := func() { //fmt.Println("tailcall")
		dx(q)
		q = pop()
		p = q + 4
		l = lastp(q)
	}
	for p <= l {
		x := I(p)
		//fmt.Println("p", p, "l", l, "x", X(x), refcount(x), "s", X(I(4)))
		if tail := p == l; tail && x&3 == 2 { // symbol
			push(lup(x))
			tailcall()
		} else if tail && x == 93 { // .
			tailcall()
		} else if tail && x == 127 { // ?
			e := pop()
			t := pop()
			if ipo() == 0 {
				dx(t)
				push(e)
				tailcall()
			} else {
				dx(e)
				push(t)
				tailcall()
			}
		} else if x&3 == 2 { // symbol
			Exec(lup(x))
		} else if x&7 != 4 { // no operator but list or number
			push(rx(x))
		} else if x == 740 { // } push to reg
			t := pop()
			r = swap(r)
			push(t)
			r = swap(r)
		} else if x == 724 { // { pop from reg
			r = swap(r)
			t := pop()
			r = swap(r)
			push(t)
		} else {
			//fmt.Println("execop", x>>3)
			//fmt.Println(x)
			F[x>>3]()
		}
		p += 4
	}
	dx(q)
	dx(r)
}
func swap(r uint32) uint32 { // swap register and stack
	if r == 0 {
		r = mk(0)
	}
	s := I(4)
	sI(4, r)
	return s
}
func exe() { Exec(lpo()) } // [q]. exec
func ife() { // c[t][e]?  (if c then t else e)
	e := pop()
	t := pop()
	if ipo() == 0 {
		dx(t)
		Exec(e)
	} else {
		dx(e)
		Exec(t)
	}
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
	return x
}
func fr(x uint32) {
	p := 4 * bk(I(4+x))
	sI(x, I(p))
	sI(p, x)
}
func nn(x uint32) uint32 { return I(4 + x) }
func cat() { // ,
	y := pop()
	x := pop()
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
	push(x)
}
func lcat(x uint32, y uint32) (r uint32) {
	/*
		rs := fmt.Sprint("lcat", X(x), refcount(x), X(y), refcount(y))
		defer func() {
			fmt.Printf("%s => r=%s %s", rs, X(r), refcount(r))
			if nn(r) > 0 {
				fmt.Printf(" *%s", refcount(I(8+r)))
			}
			fmt.Println()

		}()
	*/
	n := nn(x)
	if n == 0 {
		dx(x)
		r = mk(1)
		sI(8+r, y)
		return r
	}
	if I(x) == 1 && bk(n) == bk(1+n) {
		sI(4+lastp(x), y)
		sI(4+x, 1+I(4+x))
		return x
	}
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
func pcat(y uint32) (r uint32) {
	s := I(8)
	p, t := I(8+s), I(12+s)
	q := parent(p, t)
	r = lcat(t, y)
	sI(12+s, r)
	if t == p {
		sI(8+s, r)
		return r
	}
	sI(lastp(q), r)
	return r
}
func parent(x, y uint32) (r uint32) {
	for {
		if x&7 != 0 {
			panic("parent")
		}
		if nn(x) == 0 {
			return x
		}
		l := last(x)
		if l == y || l == 0 || x == y {
			return x
		}
		x = l
	}
}
func lastp(x uint32) uint32 {
	n := nn(x)
	if n == 0 {
		panic("empty")
	}
	return 4 + x + 4*n
}
func last(x uint32) uint32 { return I(lastp(x)) }
func first(x uint32) (r uint32) {
	if nn(x) == 0 {
		panic("empty")
	}
	r = rx(I(x + 8))
	dx(x)
	return r
}

func stk() { fmt.Println("(stk) " + X(I(4))) } // !
func swp() { // ~
	x := pop()
	y := pop()
	push(x)
	push(y)
}
func dup() { x := pop(); push(x); push(x) } // "
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
	push(x)
	push(r)
}
func use(x uint32) uint32 {
	if I(x) == 1 {
		//fmt.Println("reuse")
		return x
	}
	//fmt.Println("use: new rc", I(x), X(x))
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
func atx() { // [..]i@
	i := ipo()
	l := lpo()
	if i < 0 || i >= int32(nn(l)) {
		panic("atx: range")
	}
	push(rx(I(8 + 4*uint32(i) + l)))
	dx(l)
}
func amd() { // [a]i v$ amend (set array at index i to v)
	v := pop()
	i := ipo()
	a := use(lpo())
	n := int32(nn(a))
	if i == n {
		a = lcat(a, v)
	} else if i < 0 || i > n {
		panic("amd: range")
	} else {
		ap := 8 + a + 4*uint32(i)
		rx(I(ap))
		sI(ap, v)
	}
	push(a)
}
func ipo() int32 {
	x := pop()
	if x&1 == 0 {
		panic("int expected")
	}
	return int32(x) >> 1
}
func lpo() uint32 {
	x := pop()
	if x&7 != 0 {
		panic("list expected")
	}
	return x
}
func add()       { pi(ipo() + ipo()) }
func sub()       { pi(-ipo() + ipo()) }
func mul()       { pi(ipo() * ipo()) }
func div()       { swp(); pi(ipo() / ipo()) }
func mod()       { swp(); pi(ipo() % ipo()) }
func eql()       { pb(ipo() == ipo()) }
func gti()       { pb(ipo() < ipo()) }
func lti()       { pb(ipo() > ipo()) }
func pi(i int32) { push(1 + 2*uint32(i)) }
func pb(b bool) {
	if b {
		push(3)
	} else {
		push(1)
	}
}
func drp() { pop() } // x _ -- (pop)
func asn() { // [q][s]: -- (assign)
	y := first(lpo())
	if y&3 != 2 {
		panic("asn: not a symbol")
	}
	v := pop()
	if v&7 != 0 {
		v = lcat(mk(0), v) // enlist atoms
	}
	s := I(12)
	p := fns(s, y)
	if p == 0 {
		s = lcat(s, y)
		s = lcat(s, 1)
		p = lastp(s)
	}
	dx(I(p))
	sI(p, v)
	sI(12, s)
}
func lup(x uint32) uint32 {
	p := fns(I(12), x)
	if p == 0 {
		panic("undefined: " + X(x))
	}
	return rx(I(p))
}
func fns(x, y uint32) uint32 {
	n := nn(x) / 2
	p := x + 8
	for i := uint32(0); i < n; i++ {
		if I(p) == y {
			return p + 4
		}
		p += 8
	}
	return 0
}
func pop() (r uint32) {
	s := I(4)
	n := nn(s)
	if n == 0 {
		panic("stack underflow")
	}
	if I(s) != 1 {
		panic("stack rc")
	}
	p := lastp(s)
	r = I(lastp(s))
	sI(p, 0)
	n--
	if bk(n) == bk(1+n) {
		sI(4+s, n)
		return r
	}
	q := mk(n)
	qp := q + 8
	sp := s + 8
	for i := uint32(0); i < n; i++ {
		sI(qp, rx(I(sp)))
		qp += 4
		sp += 4
	}
	dx(s)
	sI(4, q)
	return r
}
func push(x uint32) {
	s := I(4)
	if I(s) != 1 {
		panic("stack rc")
	}
	sI(4, lcat(s, x))
}
func finit() {
	f := func(c byte, g func()) { F[c-33] = g }
	F = make([]func(), 128)
	f('!', stk) // 0
	f('"', dup) // 1
	f('#', cnt) // 2
	f('$', amd) // 3
	f('%', mod) // 4
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
func init() { ini() }
