package k

import (
	. "github.com/ktye/wg/module"
)

var src K
var pp, pe, sp int32 //parse or execution position/end, stack position

func init() {
	Memory(1)
	Data(132, "\x10A`AAAAI@pAAAAAIddddddddddAPAAAAAbbbbbbbbbbbbbbbbbbbbbbbbbb@IpAA@bbbbbbbbbbbbbbbbbbbbbbbbbb@ApA")
	Data(228, ":+-*%!&|<>=~,^#_$?@.'/\\")
	Export(ktest, kinit, Ki, iK, Til, Count, At)
	//            :    +    -    *    %    !    &    |  <>=~,^#_$?@.'/\\
	Functions(41, nil, Flp, Neg, Fst, Sqr, Til, Wer, Rev)
	Functions(81, Dex, Add, Sub, Mul, Div, Key, Min, Max)
	Functions(120, tnum, tvrb)
}

//  16..127  free list
// 128..131  memsize log2
// 132..227  char map (starts at 100)
// 228..250  verbs :+-*%!&|<>=~,^#_$?@.'/\
// 256..511  stack
func kinit() {
	minit(10, 16)
	sp = 256
}
func ktest(x int32) int32 {
	minit(10, 16)
	return iK(Count(Til(Ki(x))))
}

type K uint64
type T uint32

// typeof(x K): t=x>>59
// isatom:      t<16
// isvector:    t>16
// isflat:      t<22
// basetype:    t&15  0..9
// istagged:    t<5
// haspointers: t>5   (recursive unref)
// elementsize: $[t<19;1;t<21;4;8]
const ( //base t&15          bytes  atom  vector
	bt T = 1  // bool    1      1     17
	ct T = 2  // char    1      2     18
	it T = 3  // int     4      3     19
	st T = 4  // symbol  4      4     20
	ft T = 5  // float   8      5     21
	zt T = 6  // complex(8)     6     22
	lt T = 7  // list    8            23
	dt T = 8  // dict   (8)           24
	tt T = 9  // table  (8)           25
	cf T = 10 // comp   (8)    10
	df T = 11 // derived(8)    11
	pf T = 12 // proj   (8)    12
	lf T = 13 // lambda (8)    13
	Bt T = bt + 16
	Ct T = ct + 16
	It T = it + 16
	St T = st + 16
	Ft T = ft + 16
	Zt T = zt + 16
	Lt T = lt + 16
	Dt T = dt + 16
	Tt T = tt + 16
)

// func t=0
// basic x < 64 (triadic/tetradic)
// xn=2: composition string funclist
// xn=3: derived     string func symb
// xn=4: projection  string func arglist emptylist
// xn=5: lambda      string code locals arity save

// ptr: int32(x)
// p-8       p-4     p
// [refcount][length][data]

func Kb(x bool) K  { return K(ib(x)) | K(bt)<<59 }
func Kc(x int32) K { return K(x) | K(ct)<<59 }
func Ki(x int32) K { return K(uint32(x)) | K(it)<<59 }
func iK(x K) int32 { return int32(x) }
func Kf(x float64) (r K) {
	r = mk(ft+16, 1)
	SetF64(int32(r), x)
	return K(int32(r)) | K(ft)<<59
}
func Kz(x, y K) (z K) {
	z = l2(x, y)
	return K(int32(z)) | K(zt)<<59
}
func l1(x K) (r K) {
	r = mk(Lt, 1)
	SetI64(int32(r), int64(x))
	return r
}
func l2(x, y K) (r K) {
	r = mk(lt+16, 2)
	SetI64(int32(r), int64(x))
	SetI64(8+int32(r), int64(y))
	return r
}
func Ku(x uint64) (r K) { // Ct
	r = mk(Ct, 0)
	p := int32(r)
	for x != 0 {
		SetI8(p, int32(x))
		x >>= uint64(8)
		p++
	}
	SetI32(int32(r)-4, p-int32(r))
	return r
}

/* encode bytes with
func enc(x []byte) (r uint64) {
	var o uint64 = 1
	for _, b := range x {
		r += o * uint64(b)
		o <<= 8
	}
	return r
}
*/

func ib(x bool) int32 {
	if x {
		return 1
	}
	return 0
}
