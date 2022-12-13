package main

import (
	. "github.com/ktye/wg/module"
)

var M []uint32 // heap
var F []func() // function table

func init() {
	// no @ $ /: \: ':

	//
	// juxtaposition
	// aa cal
	// av cal
	// va idx
	// vv idv
	Data(0, "\x3a\x2b\x2d\x2a\x25\x26\x7c\x3c\x3e\x3d\x7e\x21\x2c\x5e\x23\x5f\x3f\x2e\x27\x2f\x5c")

	//      a     ???  abs  neg  fst  ???  zer  ???  ???  ???  ???  not  til  enl  ???  cnt  ???  ???  lup  ---  ---  ---
	//      v     ???  flp  ngv  fst  rot  wer  rev  grd  gdn  grp  ntv  odo  enl  srt  cnt  dp1  unq  exe  ---  ---  ---
	//              :    +    -    *    %    &    |    <    >    =    ~   !    ,     ^    #    _    ?    .    '    /    \`
	//             58   43   45   42   37   38  124   60   62   61  126  33   44    94   35   95   63   46   39   47   92
	//              0    1    2    3    4    5    6    7    8    9   10  11   12    13   14   15   16   17   18   19   20
	Functions(00, asn, add, sub, mul, div, min, max, les, mor, eql, mtc, mod, caa, pow, rep, ny2, del, cla, ech, ovr, scn) //aa
	Functions(21, asn, adx, sbx, mlx, dvx, mnx, mxx, lsx, mrx, eqx, nul, mdx, cav, cut, tak, drp, sfl, cal, ech, ovr, scn) //av
	Functions(42, asn, adx, sbx, mlx, dvx, mnx, mxx, lsx, mrx, eqx, nul, mdx, cva, pov, tka, dra, fnd, ati, ecv, ovv, scn) //va
	Functions(63, asn, adx, sbx, mlx, dvx, mnx, mxx, lsx, mrx, eqx, mtv, mdx, cvv, ctv, tkv, drv, fnx, atv, ecv, ovv, scn) //vv
}

type f1 = func(int32) int32
type f2 = func(int32, int32) int32

//dyadic
func ny2(x, y int32) int32 { return y }
func asn(x, y int32) int32 { return ny2(x, y) }
func add(x, y int32) int32 { return ((x >> 1) + (y >> 1)) << 1 }
func adx(x, y int32) int32 { return ec2(43, x, y) }
func sub(x, y int32) int32 { return ((x >> 1) - (y >> 1)) << 1 }
func sbx(x, y int32) int32 { return ec2(45, x, y) }
func mul(x, y int32) int32 { return ((x >> 1) * (y >> 1)) << 1 }
func mlx(x, y int32) int32 { return ec2(42, x, y) }
func div(x, y int32) int32 { return ((x >> 1) / (y >> 1)) << 1 }
func dvx(x, y int32) int32 { return ec2(37, x, y) }
func min(x, y int32) int32 { x >>= 1; y >>= 1; return tern(I32B(x < y), x, y) << 1 }
func mnx(x, y int32) int32 { return ec2(38, x, y) }
func max(x, y int32) int32 { x >>= 1; y >>= 1; return tern(I32B(x > y), x, y) << 1 }
func mxx(x, y int32) int32 { return ec2(124, x, y) }
func les(x, y int32) int32 { return I32B(x < y) << 1 }
func lsx(x, y int32) int32 { return ec2(60, x, y) }
func mor(x, y int32) int32 { return I32B(x > y) << 1 }
func mrx(x, y int32) int32 { return ec2(62, x, y) }
func eql(x, y int32) int32 { return I32B(x == y) << 1 }
func eqx(x, y int32) int32 { return ec2(61, x, y) }
func mtc(x, y int32) int32 { return I32B(x == y) << 1 }
func nul(x, y int32) int32 { return 0 }
func mtv(x, y int32) int32 {
	xn := n(x)
	if x == y {
		return 2
	}
	if xn != n(y) {
		return 0
	}
	return ovr(38, ec2(126, x, y)) // &/x~'y
}
func mod(x, y int32) int32 { return ((x >> 1) % (y >> 1)) << 1 }
func mdx(x, y int32) int32 { return ec2(33, x, y) }
func caa(x, y int32) int32 { return cvv(enl(x), enl(y)) } // a,a
func cav(x, y int32) int32 { return cvv(enl(x), y) }      // a,v
func cva(x, y int32) int32 { return cvv(x, enl(y)) }      // v,a
func cvv(x, y int32) int32 { // v,v
	xn, yn := n(x), n(y)
	r := mk(xn + yn)
	Memorycopy(v(r), v(x), 4*xn)
	Memorycopy(v(r)+4*xn, v(y), 4*yn)
	return r
}
func pow(x, y int32) int32 { // a^a
	x, y = v(x), v(y)
	r := int32(1)
	for {
		if y&1 == 1 {
			r *= x
		}
		y >>= 1
		if y == 0 {
			break
		}
		x *= x
	}
	return r << 1
}
func pov(x, y int32) int32 { return ec2(94, x, y) }               // v^a
func cut(x, y int32) int32 { return ny2(x, y) }                   // a^v
func ctv(x, y int32) int32 { return ny2(x, y) }                   // v^v
func rep(x, y int32) int32 { return atv(y, til(x)) }              // a#a
func tak(x, y int32) int32 { return atv(y, mdx(cnt(y), til(x))) } // a#v
func tka(x, y int32) int32 { return tka(x, enl(y)) }              // v#a (overload?)
func tkv(x, y int32) int32 { return atv(y, wer(inn(y, x))) }      // v#v
func drt(x, y int32) int32 { return drp(x, til(y)) }              // a_a   x_!y
func drp(x, y int32) int32 { // a_v
	yn := n(y)
	xv := v(x)
	if x < 0 {
		if xv < -yn {
			return cvv(zer((yn+xv)<<1), x)
		} else {
			return tak(x, (-xv)<<1)
		}
	} else {
		if xv > y {
			return mk(0)
		} else {
			return adx(x, til((yn-xv)<<1))
		}
	}
}
func dra(x, y int32) int32 { return drv(x, til(cnt(x))) }                       // v_a
func drv(x, y int32) int32 { return atv(y, wer(not(inn(y, x)))) }               // v_v
func del(x, y int32) int32 { return ny2(x, y) }                                 // a?a
func sfl(x, y int32) int32 { return ny2(x, y) }                                 // a?v
func fnd(x, y int32) int32 { return fst(ec2(126, x, y)) }                       // v?a
func fnx(x, y int32) int32 { return ec2(63, enl(x), y) }                        // v?v   x?/:y
func cla(x, y int32) int32 { return cal(x, enl(y)) }                            // a.a
func cal(x, y int32) int32 { t := loc; loc = y; x = exe(x); loc = t; return x } // a.v
func ati(x, y int32) int32 { // v.a  (also a.v)
	xn := n(x)
	if xn < 0 {
		return x
	}
	yv := v(y)
	if yv < 0 || xn > yv {
		return 0
	} else {
		return I32(v(x) + 4*yv)
	}
}
func atv(x, y int32) int32 { return ec2(46, enl(x), y) } // v.v  x./:y
func ech(x, y int32) int32 { // a'a  a'v
	yn := v(cnt(y))
	r := mk(yn)
	p := v(r)
	for i := int32(0); i < yn; i++ {
		SetI32(p, cal(x, ati(y, i<<1)))
		p += 4
	}
	return r
}
func ecv(x, y int32) int32 { return ec2(39, x, enl(y)) } // v'a  v'a
func ovr(x, y int32) int32 { // a/a  a/v
	yn := v(cnt(y))
	r := fst(x)
	for i := int32(1); i < yn; i++ {
		r = cal(x, ati(y, i<<1))
	}
	return r
}
func ovv(x, y int32) int32 { ec2(47, x, enl(y)) } // v/a  v/v
func scn(x, y int32) int32 { // a\a  a\v
	yn := v(cnt(y))
	r := mk(yn)
	p := v(r)
	f := fst(x)
	SetI32(p, f)
	for i := int32(0); i < yn; i++ {
		p += 4
		f = cal(x, cvv(enl(f), enl(ati(y, i<<1))))
		SetI32(p, f)
	}
	return r
}

//monadic
func nyi(x int32) int32 { return x }
func enl(x int32) int32 { r := mk(1); SetI32(r, x); return r }
func til(x int32) int32 {
	x = v(x)
	r := mk(x)
	p := v(r)
	for i := int(0); i < x; i++ {
		SetI32(p+4*i, i<<1)
	}
	return r
}
func cnt(x int32) int32 {
	xn := n(x)
	if xn < 0 {
		return 2
	} else {
		return xn << 1
	}
}
func zer(x int32) int32 { return mk(v(x)) } // &a
func wer(x int32) int32 { // &v
	r := mk(0)
	xn := n(x)
	for i := int32(0); i < xn; i++ {
		r = cvv(r, rep(ati(x, i<<1)))
	}
	return r
}
func not(x int32) int32 { return I32B(v(1) == 0) << 1 } // ~a
func ntv(x int32) int32 { return ech(252, x) }          // ~v
func fst(x int32) int32 { return ati(x, 0) }            // *a  *v

func ec2(f, x, y int32) int32 {
	rn := max(cnt(x), cnt(y)) >> 1
	r := mk(rn)
	p := v(r)
	for i := int32(0); i < rn; i++ {
		SetI32(p+4*i, Functions[f].(f2)(ati(x, i), ati(y, i)))
	}
	return r
}

func inn(x, y int32) int32 {
	xn := n(x)
	yn := n(y) << 1
	if n(x) < 0 {
		return mor(yn, fnx(y, x))
	} else {
		return ovr(124, gtx(yn, fnd(y, x))) // |/(#y)>y?x
	}
}

func v(x int32) int32 { return x >> 1 }
func n(x int32) int32 {
	r := int32(-1)
	if x&1 == 0 {
		r = I32(v(x) - 4)
	}
	return r << 1
}

func tern(c, x, y int32) int32 {
	if c == 0 {
		return y
	} else {
		return x
	}
}
func mk(x int32) int32 {
	x = x*4 + 4
	for tot < top+x {
		Memorygrow(1)
		tot += 65536
	}
	r := top
	SetI32(r, x)
	top += x
	return 1 + (r << 1)
}
