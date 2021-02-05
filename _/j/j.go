// ..000 pointer(list)
// ....1 int  x>>1
// ...10 symbol x>>2
// ..100 operator x>>3
//
// 0     total memory (log2)
// 1     symbol list
// 2     value list
// 3
// 4..32 free list
//
// abc   symbol (max 6)
// 123   int (max 31 bit)
// [..]  list/quote
// [..]. exec
// [..]; continue with (insert below in execution stack)
// #     length/non-list: -1
// [.][a]: assign
// a     lookup&exec
// +-*%\ arith(mod)
// <=>   compare
// &^    min max
// ,     cat                    [[]~,][enlist]:
// ~"_   swap dup pop
//
// nyi:
// 'each /over
// a i@ index
// a i v$store
// ;putc
// !trace
// `trap
// depth{  setpc at depth to start   [1{][self]:
// depth}  setpc at depth to end     [1}][return]:
// (comment)
//  go:embed j.j
//  var j []byte
package j

import (
	_ "embed"
	"fmt"
	"math/bits"
)

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
	M[1] = mk(0)
	M[2] = mk(0)
	P = mk(0)
	T = P
	//dump(127)
}

func Step(x uint32) uint32 {
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
		N *= 10
		N += x - '0'
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
			return dx(v)
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
func Exec(x uint32) uint32 {
	if nn(x) == 0 {
		return x
	}
	est := lcat(mk(0), x) // execution stack
	rst := lcat(mk(0), 1) // return stack
	stk := mk(0)          // value stack
	exe := uint32(0)
	p := x + 8
	l := lastp(x)
	for {
		//fmt.Println("p/l", p, l, X(stk))
		for p > l {
			if nn(rst) == 1 {
				dx(est)
				dx(rst)
				return stk
			}
			p = last(rst) >> 1
			est = pop(est)
			rst = pop(rst)
			x = last(est)
			l = lastp(x)
		}
		x := I(p)
		if x&3 == 2 {
			if x&3 == 2 {
				exe = lup(x)
			}
		} else if x&7 != 4 {
			stk = lcat(stk, rx(x))
		} else {
			// fmt.Println("x", x, X(x))
			if x == 108 { // . execute
				if p == x+8 {
					panic(". underflow")
				}
				exe = last(stk)
				stk = pop(stk)
			} else if x == 724 { // depth{ reset pc
				panic("resetpc")
				fmt.Println("ddpc")
				a := lasti(stk)
				stk = pop(stk)
				if a == 0 {
					p += a
				} else {
					panic("depty nyi")
				}
			} else if x == 740 { // depth} drop at depth
				a := lasti(stk)
				stk = pop(stk)
				if a == 0 {
					p = l - 4
				} else {
					//fmt.Println("} del@", a, X(est), X(rst))
					est = delat(est, a)
					rst = delat(rst, a-1)
					//fmt.Println("! del@", a, X(est), X(rst))
				}
			} else if x == 212 { // ; continue with
				x := rx(last(stk))
				xp := 1 + 2*(x+8)
				stk = pop(stk)
				est = prependlast(est, x)
				rst = lcat(rst, xp)
			} else {
				if x == 4 { // !
					fmt.Println(" rst", X(rst))
					fmt.Println(" est", X(est))
				}
				stk = F[x>>3](stk)
			}
		}
		if exe != 0 {
			x = rx(exe)
			est = lcat(est, x)
			rst = lcat(rst, 9+2*p)
			p = x + 4
			l = lastp(x)
			exe = 0
		}
		p += 4
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
	if x&7 == 0 {
		if I(x) == 0 {
			panic("dx")
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
		return x
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
	p := lastp(pop(s))
	x := I(p)
	if x&7 != 0 {
		x = lcat(mk(0), x)
	}
	sI(p, lcat(x, y))
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

func finit() {
	f := func(c byte, g func(uint32) uint32) { F[c-33] = g }
	F = make([]func(uint32) uint32, 128)
	f('!', stk)
	f('~', swp)
	f('"', dup)
	f('_', pop)
	f('|', rol)
	f('#', cnt)
	f('+', add)
	f('-', sub)
	f('*', mul)
	f('%', dif)
	f('\\', mod)
	f('=', eql)
	f('>', gti)
	f('<', lti)
	f('&', min)
	f('^', max)
	f(':', asn)
	f(',', cat)
	f('@', atx)
}
func stk(s uint32) uint32 { // !
	fmt.Println(" " + X(s))
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
func v1(s, x uint32) uint32 {
	sp := s + 4 + 4*nn(s)
	dx(I(sp))
	sI(sp, x)
	return s
}
func ints(s uint32) (j, k int32) {
	a, b := last2(s)
	if a&1 == 0 || b&1 == 0 {
		panic("ints")
	}
	return int32(a) >> 1, int32(b) >> 1
}
func add(s uint32) uint32 { a, b := ints(s); return i2(s, a+b) }
func sub(s uint32) uint32 { a, b := ints(s); return i2(s, a-b) }
func mul(s uint32) uint32 { a, b := ints(s); return i2(s, a*b) }
func dif(s uint32) uint32 { a, b := ints(s); return i2(s, a/b) }
func mod(s uint32) uint32 { a, b := ints(s); return i2(s, a%b) }
func eql(s uint32) uint32 { a, b := last2(s); return i2(s, ib(a == b)) }
func gti(s uint32) uint32 { a, b := ints(s); return i2(s, ib(a > b)) }
func lti(s uint32) uint32 { a, b := ints(s); return i2(s, ib(a < b)) }
func max(s uint32) uint32 {
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
func i2(s uint32, a int32) uint32 {
	s = pop(s)
	n := nn(s)
	sI(s+4+4*n, uint32(1|(a<<1)))
	return s
}
func ib(b bool) int32 {
	if b {
		return 1
	}
	return 0
}
func pop(x uint32) (r uint32) {
	n := nn(x)
	if n == 0 {
		panic("pop:underflow")
	}
	if bk(n) == bk(n-1) {
		dx(last(x))
		sI(x+4, n-1)
	} else {
		n--
		r = mk(n)
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
	return x
}
func asn(s uint32) uint32 { // :
	y := last(last(s))
	if y&3 != 2 {
		panic("asn: not a symbol")
	}
	v := rx(last(pop(s)))
	if v&7 != 0 {
		v = lcat(mk(0), v) // enlist atoms
	}
	p := fns(I(4), y)
	if p == 0 {
		sI(4, lcat(I(4), y))
		sI(8, lcat(I(8), 1))
		p = 4 + 4*nn(I(4))
	}
	dx(I(I(8) + p))
	sI(I(8)+p, v)
	return pop(s)
}
func lup(x uint32) uint32 {
	p := fns(I(4), x)
	if p == 0 {
		panic("undefined: " + X(x))
	}
	return I(I(8) + p) // does not ref
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
func prependlast(x, y uint32) uint32 {
	if nn(x) == 0 {
		return lcat(x, y)
	}
	x = lcat(x, 1)
	p := lastp(x)
	sI(p, I(p-4))
	sI(p-4, y)
	return x
}
func delat(x, y uint32) uint32 { // delete at y from tail
	n := nn(x)
	if y >= n {
		panic("delat: underflow")
	}
	p := lastp(x) - 4*y
	dx(I(p))
	for i := uint32(0); i < y; i++ {
		sI(p, I(p+4))
		p += 4
	}
	sI(p, 1)

	// if bk(n) == bk(n-1) { reuse
	r := mk(n - 1)
	rp := r + 8
	p = x + 8
	for i := uint32(0); i < n-1; i++ {
		sI(rp, rx(I(p)))
		rp += 4
		p += 4
	}
	dx(x)
	return r
}
func init() { ini() }
