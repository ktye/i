package main

import (
	. "github.com/ktye/wg/module"
)

// const nai int32 = -2147483648 // 0N
var loc, xyz K
var na, inf float64
var pp, pe, sp, srcp, rand_ int32 //parse position/end, stack position, src pointer

//   0....7  key
//   8...15  val
//  16...19  src(int32)
//  20..127  free list
// 128..131  memsize log2
// 132..226  char map (starts at 100)    -+
// 227..252  :+-*%!&|<>=~,^#_$?@.':/:\:   | text
// 253..279  vbcisfzldtcdpl000BCISFZLDT   | section
// 280.....  z.k                         -+
// 2k....4k  stack

func kinit() {
	minit(12, 16) //4k..64k
	sp = 2048
	SetI32(16, int32(mk(Ct, 0))) //SetI64(512, int64(mk(Ct, 0))) //src
	na = F64reinterpret_i64(uint64(0x7FF8000000000001))
	inf = F64reinterpret_i64(uint64(0x7FF0000000000000))
	rand_ = 1592653589
	SetI64(0, int64(mk(Lt, 0)))
	SetI64(8, int64(mk(Lt, 0)))
	xyz = sc(Ku(0))
	xyz = Ech(17, l2(xyz, Ku(8026488))) //`$'"xyz": `x`y`z -> 8 16 24
	zk()
}

type K uint64
type T int32

// typeof(x K): t=x>>59
// isatom:      t<16
// isvector:    t>16
// isflat:      t<22
// basetype:    t&15  0..9
// istagged:    t<5
// haspointers: t>5   (recursive unref)
// elementsize: $[t<19;1;t<21;4;8]
const ( //base t&15          bytes  atom  vector
	ct T = 2 // char    1      2     18
	it T = 3 // int     4      3     19
	st T = 4 // symbol  4      4     20
	ft T = 5 // float   8      5     21
	zt T = 6 // complex(8)     6     22

	cf T = 10 // comp   (8)    10
	df T = 11 // derived(8)    11
	pf T = 12 // proj   (8)    12
	lf T = 13 // lambda (8)    13
	xf T = 14 // native (8)    14
	Ct T = 18
	It T = 19
	St T = 20
	Ft T = 21
	Zt T = 22
	Lt T = 23 // list
	Dt T = 24 // dict
	Tt T = 25 // table
)

// func t=0
// basic x < 64 (triadic/tetradic)
// composition .. f2 f1 f0
// derived     func    symb
// projection  func    arglist  emptylist
// lambda      code    string	locals
// native      ptr(Ct) string

// ptr: int32(x)
//  p-12    p-4 p
// [length][rc][data]

func ti(t T, i int32) K { return K(t)<<59 | K(uint32(i)) }
func Kc(x int32) K      { return ti(ct, x) }
func Ki(x int32) K      { return ti(it, x) }
func Ks(x int32) K      { return ti(st, x) }
func Kf(x float64) K {
	r := mk(Ft, 1)
	SetF64(int32(r), x)
	return ti(ft, int32(r))
}
func Kz(x, y float64) K {
	r := mk(Zt, 1)
	rp := int32(r)
	SetF64(rp, x)
	SetF64(rp+8, y)
	return ti(zt, rp)
}
func l1(x K) K {
	r := mk(Lt, 1)
	SetI64(int32(r), int64(x))
	return r
}
func l2(x, y K) K {
	r := mk(Lt, 2)
	rp := int32(r)
	SetI64(rp, int64(x))
	SetI64(8+rp, int64(y))
	return r
}
func l3(x, y, z K) K { return cat1(l2(x, y), z) }
func r0(x K) K       { r := x0(x); dx(x); return r }
func r1(x K) K       { r := x1(x); dx(x); return r }
func x0(x K) K       { return rx(K(I64(int32(x)))) }
func x1(x K) K       { return x0(x + 8) }
func x2(x K) K       { return x0(x + 16) }
func Ku(x uint64) K { // Ct
	r := mk(Ct, 8)
	p := int32(r)
	SetI64(p, int64(x))
	SetI32(p-12, idx(0, p, p+8)) //assume <8
	return r
}

/* encode bytes for Ku(..) with: https://play.golang.org/p/4ethx6OEVCR
func enc(x []byte) uint64 {
	r := uint32(0)
	var o uint64 = 1
	for _, b := range x {
		r += o * uint64(b)
		o <<= 8
	}
	return r
}
*/

func kx(u int32, x K) K { return cal(Val(Ks(u)), l1(x)) } //call k func from z.k
func sc(c K) K { //symbol from chars
	s := K(I64(0))
	sp := int32(s)
	sn := nn(s)
	for i := int32(0); i < sn; i++ {
		if match(c, K(I64(sp))) != 0 {
			dx(c)
			return ti(st, sp-int32(s))
		}
		sp += 8
	}
	SetI64(0, int64(cat1(s, c)))
	SetI64(8, int64(cat1(K(I64(8)), 0)))
	return ti(st, 8*sn)
}
func cs(x K) K { return x0(K(I32(0)) + x) } //chars from symbol
func missing(t T) K {
	switch t - 2 {
	case 0: // ct
		return Kc(32)
	case 1: // it
		return Ki(0)
	case 2: // st
		return Ks(0)
	case 3: // ft
		return Kf(na)
	case 4: // zt
		return Kz(na, na)
	default: // lt
		return mk(Ct, 0) //Kb(0)
	}
}
