package main

import . "github.com/ktye/wg/module"

var src, loc, xyz K
var nan float64
var pp, pe, sp, srcp int32 //parse or execution position/end, stack position, src pointer

func init() {
	Memory(1)
	Data(132, "\x00\x01@\x01\x01\x01\x01\t\x10`\x01\x01\x01\x01\x01\tDDDDDDDDDD\x01 \x01\x01\x01\x01\x01BBBBBBBBBBBBBBBBBBBBBBBBBB\x10\t`\x01\x01\x00BBBBBBBBBBBBBBBBBBBBBBBBBB\x10\x01`\x01")
	Data(228, ":+-*%!&|<>=~,^#_$?@.':/:\\:")
	Data(520, "vbcisfzldtcdpl000BCISFZLDT")
	Export(kinit, mk, nn, Val, Kst)
	ExportAll()
	//           0    :    +    -    *    %    !    &    |    <    >10  =    ~    ,    ^    #    _    $    ?    @    .20  '    ':   /    /:   \    \:                  30
	Functions(0, nul, nyi, Flp, Neg, Fst, Sqr, Til, Wer, Rev, Asc, Dsc, Grp, Not, Enl, Srt, Cnt, Flr, Str, Unq, Typ, Val, ech, ecp, rdc, ecr, scn, ecl, lst, Kst, Out, Any, nyi, Abs, Imag, Conj, Angle)
	Functions(64, Asn, Dex, Add, Sub, Mul, Div, Key, Min, Max, Les, Mor, Eql, Mtc, Cat, Cut, Tak, Drp, Cst, Fnd, Atx, Cal, Ech, Ecp, Rdc, Ecr, Scn, Ecl, compose, nyi, Otu, In, Find, Hypot, Cmpl, nyi, Rot)
	Functions(192, tbln, tnms, tvrb, tpct, tvar, tsym, tchr)
	Functions(211, Amd, Dmd)
	//                                                                   229                              235                                 241                           247
	//Functions(220, addi, addf, addz, subi, subf, subz, muli, mulf, mulz, divi, divf, divz, nyi, nyi, nyi, mini, minf, minz, maxi, maxf, maxz, lti, ltf, ltz, gti, gtf, gtz, eqi, eqf, eqz, nyi, nyi, rot)
	Functions(220, negc, negi, negf, negz, negC, negI, negF, negZ)
	Functions(228, absc, absi, absf, nyi, absF, absI, absF, absZ)
	Functions(236, addi, addi, addf, addz, addcC, addiI, addfF, addzZ, addC, addI, addF, addZ)
	Functions(248, subi, subi, subf, subz, subcC, subiI, subfF, subzZ, subC, subI, subF, subZ)
	Functions(260, muli, muli, mulf, mulz, mulcC, muliI, mulfF, mulzZ, mulC, mulI, mulF, mulZ)
	Functions(272, divi, divi, divf, divz, nyi, nyi, divfF, divzZ, nyi, nyi, divF, divZ)
	Functions(284, mini, mini, minf, minz, mincC, miniI, minfF, minzZ, minC, minI, minF, minZ)
	Functions(296, maxi, maxi, maxf, maxz, maxcC, maxiI, maxfF, maxzZ, maxC, maxI, maxF, maxZ)

	Functions(308, eqi, eqf, eqz, eqcC, eqiI, eqfF, eqzZ, eqCc, eqIi, eqFf, eqZz, eqC, eqI, eqF, eqZ)
	Functions(323, lti, ltf, ltz, ltcC, ltiI, ltfF, ltzZ, ltCc, ltIi, ltFf, ltZz, ltC, ltI, ltF, ltZ)
	Functions(338, gti, gtf, gtz, gtcC, gtiI, gtfF, gtzZ, gtCc, gtIi, gtFf, gtZz, gtC, gtI, gtF, gtZ)

	Functions(353, guC, guC, guI, guI, guF, guZ, guS, gdC, gdC, gdI, gdI, gdF, gdZ, gdS)

	Functions(369, sum, rd0, prd, rd0, rd0, min, max)
	Functions(376, sums, rd0, prds, rd0, rd0, mins, maxs)
	Functions(383, nyi, nyi, sqrf, nyi, nyi, nyi, negF, nyi)
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
	minit(10, 16)
	sp = 256
	src = 0
	loc = 0
	nan = F64reinterpret_i64(uint64(0x7FF8000000000001))
	SetI64(0, int64(mk(Lt, 0)))
	SetI64(8, int64(mk(Lt, 0)))
	sc(Ku(0))            // `   0
	x := sc(Ku(120))     // `x  8
	y := sc(Ku(121))     // `y 16
	z := sc(Ku(122))     // `z 24
	sc(Ku(107))          // `k 32
	sc(Ku(108))          // `l 40
	sc(Ku(435610544247)) // `while 48
	sc(Ku(28265))        // `in    56
	sc(Ku(1684957542))   // `find  64
	sc(Ku(7561825))      // `abs   72
	sc(Ku(1734438249))   // `imag  80
	sc(Ku(1785622371))   // `conj  88
	sc(Ku(435610414689)) // `angle 96

	xyz = cat1(Cat(x, y), z)
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
	if (uint32(1)<<uint32(I32(128)))-(1024+mcount()) != 0 {
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
// derived     func symb
// projection  func arglist emptylist
// lambda      code locals save string

// ptr: int32(x)
//  p-12    p-4 p
// [length][rc][data]

func Kb(x int32) K { return K(uint32(x)) | K(bt)<<59 }
func Kc(x int32) K { return K(uint32(x)) | K(ct)<<59 }
func Ki(x int32) K { return K(uint32(x)) | K(it)<<59 }
func Ks(x int32) K { return K(uint32(x)) | K(st)<<59 }
func iK(x K) int32 { return int32(x) }
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
func ZF(x K) K { xp := int32(x); SetI32(xp-12, I32(xp-12)/2); return K(uint32(x)) | K(Zt)<<59 }
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
	n := nn(x)
	m := int32(0)
	xp := int32(x)
	for i := int32(0); i < n; i++ {
		ni := nn(K(I64(xp)))
		if i == 0 {
			m = ni
		} else if m != ni {
			trap(Length)
		}
		xp += 8
	}
	r = l2(r, x)
	SetI32(int32(r)-12, m)
	return K(int32(r)) | K(Tt)<<59
}
func zero(t T) (r K) {
	if t == ft {
		return Kf(0)
	} else if t == zt {
		return Kz(0, 0)
	}
	return K(t) << 59
}
func missing(t T) (r K) {
	switch t - 2 {
	case 0: // ct
		return Kc(32)
	case 1: // it
		return Ki(-2147483648)
	case 2: // st
		return Ks(0)
	case 3: // ft
		return Kf(nan)
	case 4: // zt
		return Kz(nan, nan)
	case 5: // lt
		return Kb(0)
	default:
		return K(t) << 59
	}
}
