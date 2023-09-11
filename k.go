package main

import (
	. "github.com/ktye/wg/module"
)

const nai int32 = -2147483648 // 0N
var loc, xyz K
var na, inf float64
var pp, pe, sp, srcp, rand_ int32 //parse or execution position/end, stack position, src pointer

//   0....7  key
//   8...15  val
//  16...19
//  20..127  free list
// 128..131  memsize log2
// 132..227  char map (starts at 100)
// 228..253  verbs :+-*%!&|<>=~,^#_$?@.':/:\:
// 256..511  stack
// 512..519  wasi iovec
// 520..545  "vbcisfzldtcdpl000BCISFZLDT"
// 552..559  src (aligned)

func kinit() {
	minit(12, 16) //4k..64k
	sp = 256
	SetI64(552, int64(mk(Ct, 0))) //src
	loc = 0
	na = F64reinterpret_i64(uint64(0x7FF8000000000001))
	inf = F64reinterpret_i64(uint64(0x7FF0000000000000))
	rand_ = 1592653589
	SetI64(0, int64(mk(Lt, 0)))
	SetI64(8, int64(mk(Lt, 0)))
	sc(Ku(0))        // `   0
	x := sc(Ku(120)) // `x  8
	y := sc(Ku(121)) // `y 16
	z := sc(Ku(122)) // `z 24
	xyz = cat1(Cat(x, y), z)
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

	dt T = 8  // dict   (8)           24
	tt T = 9  // table  (8)           25
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
// lambda      code    locals   string
// native      ptr(Ct) string

// ptr: int32(x)
//  p-12    p-4 p
// [length][rc][data]

func Kc(x int32) K { return K(uint32(x)) | K(ct)<<59 }
func Ki(x int32) K { return K(uint32(x)) | K(it)<<59 }
func Ks(x int32) K { return K(uint32(x)) | K(st)<<59 }
func Kf(x float64) K {
	r := mk(Ft, 1)
	SetF64(int32(r), x)
	return K(int32(r)) | K(ft)<<59
}
func Kz(x, y float64) K {
	r := mk(Zt, 1)
	rp := int32(r)
	SetF64(rp, x)
	SetF64(rp+8, y)
	return K(rp) | K(zt)<<59
}
func l1(x K) K {
	r := mk(Lt, 1)
	SetI64(int32(r), int64(x))
	return r
}
func l2t(x, y K, t T) K {
	r := mk(Lt, 2)
	SetI64(int32(r), int64(x))
	SetI64(8+int32(r), int64(y))
	return K(uint32(r)) | K(t)<<59
}
func l2(x, y K) K    { return l2t(x, y, Lt) }
func l3(x, y, z K) K { return cat1(l2(x, y), z) }
func r0(x K) K       { r := x0(x); dx(x); return r }
func r1(x K) K       { r := x1(x); dx(x); return r }
func r3(x K) K       { r := x3(x); dx(x); return r }
func x0(x K) K       { return rx(K(I64(int32(x)))) }
func x1(x K) K       { return x0(x + 8) }
func x2(x K) K       { return x0(x + 16) }
func x3(x K) K       { return x0(x + 24) }
func Ku(x uint64) K { // Ct
	r := mk(Ct, 0)
	p := int32(r)
	for x != 0 {
		SetI8(p, int32(x))
		x >>= uint64(8)
		p++
	}
	SetI32(int32(r)-12, p-int32(r))
	return r
}
func kx(u int32, x K) K     { return cal(Val(Ks(u)), l1(x)) }
func kxy(u int32, x, y K) K { return cal(Val(Ks(u)), l2(x, y)) }

/* encode bytes with: https://play.golang.org/p/4ethx6OEVCR
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

func sc(c K) K {
	s := K(I64(0))
	sp := int32(s)
	sn := nn(s)
	for i := int32(0); i < sn; i++ {
		if match(c, K(I64(sp))) != 0 {
			dx(c)
			return K(sp-int32(s)) | K(st)<<59
		}
		sp += 8
	}
	SetI64(0, int64(cat1(s, c)))
	SetI64(8, int64(cat1(K(I64(8)), 0)))
	return K(8*sn) | K(st)<<59
}
func cs(x K) K { return x0(K(I32(0)) + x) }
func td(x K) K { // table from dict
	r := x0(x)
	x = r1(x)
	if tp(r) != St || tp(x) != Lt {
		trap(Type)
	}
	m := maxcount(int32(x), nn(x))
	x = Ech(15, l2(Ki(m), x)) // (|/#'x)#'x
	r = l2(r, x)
	SetI32(int32(r)-12, m)
	return K(int32(r)) | K(Tt)<<59
}
func missing(t T) K {
	switch t - 2 {
	case 0: // ct
		return Kc(32)
	case 1: // it
		return Ki(nai)
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
