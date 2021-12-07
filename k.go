package main

import (
	. "github.com/ktye/wg/module"
)

const nai int32 = -2147483648 // 0N
var src, loc, xyz K
var na, inf float64
var pp, pe, sp, srcp, rand_ int32 //parse or execution position/end, stack position, src pointer

/* debug c vs go
func kprint(x K) { //rm
	h := int32(x)
	l := int32(x >> 32)
	Printf("%d %d\n", h, l)
}
func printlist(x K) { //rm
	n := nn(x)
	var h, l int32
	Printf("list:#%d ", n)
	for i := int32(0); i < n; i++ {
		xi := I64(int32(x) + 8*i)
		h = int32(xi)
		l = int32(xi >> 32)
		Printf(" (%d,%d)\n", h, l)
	}
	Printf("\n")
}
func printints(x K) { //rm
	n := nn(x)
	Printf("ints[%d]", n)
	var xi int32
	for i := int32(0); i < n; i++ {
		xi = I32(int32(x) + 4*i)
		Printf(" %d", xi)
	}
	Printf("\n")
}
func printfloats(x K) { //rm
	n := nn(x)
	Printf("floats[%d]", n)
	var xi float64
	for i := int32(0); i < n; i++ {
		xi = F64(int32(x) + 8*i)
		Printf(" %f", xi)
	}
	Printf("\n")
}
func printchars(x K) { //rm
	n := nn(x)
	Printf("bools/chars[%d]: ", n)
	var xi int32
	for i := int32(0); i < n; i++ {
		xi = I8(int32(x) + i)
		Printf(" %d", xi)
	}
	Printf("\n")
}
*/

func init() {
	Memory(1)
	Memory2(1)
	Data(132, "\x00\x01@\x01\x01\x01\x01\t\x10`\x01\x01\x01\x01\x01\t\xc4\xc4\xc4\xc4\xc4\xc4\xc4\xc4\xc4\xc4\x01 \x01\x01\x01\x01\x01BBBBBBBBBBBBBBBBBBBBBBBBBB\x10\t`\x01\x01\x00\xc2\xc2\xc2\xc2\xc2\xc2BBBBBBBBBBBBBBBBBBBB\x10\x01`\x01") // k_test.go: TestClass
	Data(228, ":+-*%!&|<>=~,^#_$?@.':/:\\:")
	Data(520, "vbcisfzldtmdplx00BCISFZLDT") //546
	Export(kinit, mk, nn, Val, Kst)
	ExportAll()
	//            0    :    +    -    *    %    !    &    |    <    >10  =    ~    ,    ^    #    _    $    ?    @    .20  '    ':   /    /:   \    \:                  30                       35                       40                       45
	Functions(00, nul, Idy, Flp, Neg, Fst, Sqr, Til, Wer, Rev, Asc, Dsc, Grp, Not, Enl, Srt, Cnt, Flr, Str, Unq, Typ, Val, ech, ecp, rdc, ecr, scn, ecl, lst, Kst, Out, Any, Fin, Abs, Img, Cnj, Ang, nyi, Uqs, Xpt, Cos, Fwh, Las, Exp, Log, Sin, Tok, Prs)
	Functions(64, Asn, Dex, Add, Sub, Mul, Div, Key, Min, Max, Les, Mor, Eql, Mtc, Cat, Cut, Tak, Drp, Cst, Fnd, Atx, Cal, Ech, Ecp, Rdc, Ecr, Scn, Ecl, com, prj, Otu, In, Find, Hyp, Cpx, nyi, Rot, Enc, Dec, Fil, nyi, Bin, Mod, Pow, Lgn, nyi, nyi, Rtp)
	Functions(192, tbln, tchr, tnms, tvrb, tpct, tvar, tsym) //                                                                                                                             ^ccal(94)
	Functions(211, Amd, Dmd)

	Functions(220, negi, negf, negz, negC, negI, negF, negZ)
	Functions(227, absi, absf, nyi, absC, absI, absF, absZ)
	Functions(234, addi, addf, addz, addcC, addiI, addfF, addzZ, addC, addI, addF, addZ)
	Functions(245, subi, subf, nyi, subcC, subiI, subfF, subzZ, subC, subI, subF, subZ)
	Functions(256, muli, mulf, mulz, mulcC, muliI, mulfF, mulzZ, mulC, mulI, mulF, mulZ)
	Functions(267, divi, divf, divz, nyi, nyi, divfF, divzZ, nyi, nyi, divF, divZ)
	Functions(278, mini, minf, minz, mincC, miniI, minfF, minzZ, minC, minI, minF, minZ)
	Functions(289, maxi, maxf, maxz, maxcC, maxiI, maxfF, maxzZ, maxC, maxI, maxF, maxZ)
	Functions(300, nyi, sqrf, nyi, nyi, nyi, sqrF, nyi)

	Functions(308, lti, ltf, ltz, ltcC, ltiI, ltfF, ltzZ, ltCc, ltIi, ltFf, ltZz, ltC, ltI, ltF, ltZ)
	Functions(323, gti, gtf, gtz, gtcC, gtiI, gtfF, gtzZ, gtCc, gtIi, gtFf, gtZz, gtC, gtI, gtF, gtZ)
	Functions(338, eqi, eqf, eqz, eqcC, eqiI, eqfF, eqzZ, eqCc, eqIi, eqFf, eqZz, eqC, eqI, eqF, eqZ)

	Functions(353, guC, guC, guI, guI, guF, guZ, guL, gdC, gdC, gdI, gdI, gdF, gdZ, gdL)

	Functions(367, sum, rd0, prd, rd0, rd0, min, max)
	Functions(374, sums, rd0, prds, rd0, rd0, mins, maxs)
}

//   0....7  key
//   8...15  val
//  20..127  free list
// 128..131  memsize log2
// 132..227  char map (starts at 100)
// 228..253  verbs :+-*%!&|<>=~,^#_$?@.':/:\:
// 256..511  stack
// 512..519  wasi iovec
// 520..545  "vbcisfzldtcdpl000BCISFZLDT"

func kinit() {
	minit(13, 16) //8k..64k
	sp = 256
	src = mk(Ct, 0)
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
func reset() {
	if sp != 256 {
		panic(Stack)
	}
	dx(src)
	dx(xyz)
	dx(K(I64(0)))
	dx(K(I64(8)))
	//check() // k_test.go
	if (uint32(1)<<uint32(I32(128)))-(8192+mcount()) != 0 {
		trap(Err)
	}
	for i := int32(5); i < 31; i++ {
		SetI32(4*i, 0)
	}
	kinit()
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
	xf T = 14 // native (8)    14
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
// composition .. f2 f1 f0
// derived     func    symb
// projection  func    arglist  emptylist
// lambda      code    locals   string
// native      ptr(Ct) string

// ptr: int32(x)
//  p-12    p-4 p
// [length][rc][data]

func Kb(x int32) K { return K(uint32(x)) | K(bt)<<59 }
func Kc(x int32) K { return K(uint32(x)) | K(ct)<<59 }
func Ki(x int32) K { return K(uint32(x)) | K(it)<<59 }
func Ks(x int32) K { return K(uint32(x)) | K(st)<<59 }
func Kf(x float64) (r K) {
	r = mk(Ft, 1)
	SetF64(int32(r), x)
	return K(int32(r)) | K(ft)<<59
}
func Kz(x, y float64) (r K) {
	r = mk(Zt, 1)
	rp := int32(r)
	SetF64(rp, x)
	SetF64(rp+8, y)
	return K(rp) | K(zt)<<59
}
func l1(x K) (r K) {
	r = mk(Lt, 1)
	SetI64(int32(r), int64(x))
	return r
}
func l2t(x, y K, t T) (r K) {
	r = mk(Lt, 2)
	SetI64(int32(r), int64(x))
	SetI64(8+int32(r), int64(y))
	return K(int32(r)) | K(t)<<59
}
func l2(x, y K) (r K)    { return l2t(x, y, Lt) }
func l3(x, y, z K) (r K) { return cat1(l2(x, y), z) }
func x0(x int32) K       { return rx(K(I64(x))) }
func x1(x int32) K       { return x0(x + 8) }
func x2(x int32) K       { return x0(x + 16) }
func x3(x int32) K       { return x0(x + 24) }
func Ku(x uint64) (r K) { // Ct
	r = mk(Ct, 0)
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
func enc(x []byte) (r uint64) {
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
func cs(x K) (r K) { return x0(I32(0) + int32(x)) }
func td(x K) (r K) { // table from dict
	r, x = spl2(x)
	if tp(r) != St || tp(x) != Lt {
		trap(Type)
	}
	m := maxcount(int32(x), nn(x))
	x = Ech(15, l2(Ki(m), x)) // (|/#'x)#'x
	r = l2(r, x)
	SetI32(int32(r)-12, m)
	return K(int32(r)) | K(Tt)<<59
}
func missing(t T) (r K) {
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
	case 5: // lt
		return mk(Ct, 0) //Kb(0)
	default:
		return K(t) << 59
	}
}
