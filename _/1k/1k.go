package main

import (
	. "github.com/ktye/wg/module"
)

var tot, top int32

func init() {
	Memory(1)
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
	Functions(00, asn, mtc, ny2, cat, pow, rep, ny2, del, cal, ech, ovr, scn) //aa
	Functions(12, asn, nul, ny2, cat, cut, tak, drp, sfl, cal, ech, ovr, scn) //av
	Functions(24, asn, nul, ny2, cat, pov, tka, dra, fnd, ati, ecv, ovv, scn) //va
	Functions(36, asn, mtv, ny2, cat, ctv, tkv, drv, fnx, atv, ecv, ovv, scn) //vv
	Functions(48, add, sub, mul, div, min, max, les, mor, eql, mod)           //scalar

}
func main() {
	tot = 65536
}

type f1 = func(int32) int32
type f2 = func(int32, int32) int32

//dyadic
func add(x, y int32) int32 { return x + y }
func sub(x, y int32) int32 { return x - y }
func mul(x, y int32) int32 { return x * y }
func div(x, y int32) int32 { return x / y }
func mod(x, y int32) int32 { return x % y }
func min(x, y int32) int32 {
	if x < y {
		return x
	} else {
		return y
	}
}
func max(x, y int32) int32 {
	if x > y {
		return x
	} else {
		return y
	}
}
func les(x, y int32) int32 { return I32B(x < y) }
func mor(x, y int32) int32 { return I32B(x > y) }
func eql(x, y int32) int32 { return I32B(x == y) }

func ny2(x, y int32) int32 { return y }

func asn(x, y int32) int32 { return ny2(x, y) }
func mtc(x, y int32) int32 { return w(I32B(x == y)) }
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
func cat(x, y int32) int32 { // x,y
	x, y = el(x), el(y)
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
	return w(r)
}
func pov(x, y int32) int32 { return ec2(94, x, y) }               // v^a
func cut(x, y int32) int32 { return ny2(x, y) }                   // a^v
func ctv(x, y int32) int32 { return ny2(x, y) }                   // v^v
func rep(x, y int32) int32 { return atv(y, til(x)) }              // a#a
func tak(x, y int32) int32 { return atv(y, max(cnt(y), til(x))) } // a#v
func tka(x, y int32) int32 { return tka(x, enl(y)) }              // v#a (overload?)
func tkv(x, y int32) int32 { return atv(y, wer(inn(y, x))) }      // v#v
func drt(x, y int32) int32 { return drp(x, til(y)) }              // a_a   x_!y
func drp(x, y int32) int32 { // a_v
	yn := n(y)
	xv := v(x)
	if x < 0 {
		if xv < -yn {
			return cat(zer(w(yn+xv)), x)
		} else {
			return tak(x, w(-xv))
		}
	} else {
		if xv > y {
			return mk(0)
		} else {
			return add(x, til(w(yn-xv)))
		}
	}
}
func dra(x, y int32) int32 { return drv(x, til(cnt(x))) }         // v_a
func drv(x, y int32) int32 { return atv(y, wer(not(inn(y, x)))) } // v_v
func del(x, y int32) int32 { return ny2(x, y) }                   // a?a
func sfl(x, y int32) int32 { return ny2(x, y) }                   // a?v
func fnd(x, y int32) int32 { return fst(ec2(126, x, y)) }         // v?a
func fnx(x, y int32) int32 { return ec2(63, enl(x), y) }          // v?v   x?/:y
func cal(x, y int32) int32 { // a.a  a.v
	y = el(y)
	yn := n(y)
	i := int32(0)
	for ; i < 21; i++ {
		if x == I8(i) {
			break
		}
	}
	if i > 21 {
		return prs(lup(x), y)
	}
	x = fst(y)
	y = I32(4 + v(y))
	xa := I32B(n(x) < 0)
	ya := I32B(n(y) < 0)
	mo := I32B(yn < 2)
	if i < 10 { //scalar
		i += 48
		if mo != 0 {
			return nd(i, x, y, xa+ya)
		}
		return nm(i+10, x, xa)
	}
	i += 24*ya + 12*xa
	if mo != 0 {
		return Func[i+48].(f1)(x)
	}
	return Func[i].(f2)(x, y)
}
func nm(f, x, a int32) int32 {
	if a != 0 {
		return ech(46, cat(f, enl(x)))
	}
	return w(Func[f].(f1)(v(x)))
}
func nd(f, x, y, a int32) int32 {
	if a != 0 {
		return ec2(I8(f), x, y)
	}
	return w(Func[f].(f2)(v(x), v(y)))
}
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
		SetI32(p, cal(x, ati(y, w(i))))
		p += 4
	}
	return r
}
func ecv(x, y int32) int32 { return ec2(39, x, enl(y)) } // v'a  v'a
func ovr(x, y int32) int32 { // a/a  a/v
	yn := v(cnt(y))
	r := fst(x)
	for i := int32(1); i < yn; i++ {
		r = cal(x, ati(y, w(i)))
	}
	return r
}
func ovv(x, y int32) int32 { return ec2(47, x, enl(y)) } // v/a  v/v
func scn(x, y int32) int32 { // a\a  a\v
	yn := v(cnt(y))
	r := mk(yn)
	p := v(r)
	f := fst(x)
	SetI32(p, f)
	for i := int32(0); i < yn; i++ {
		p += 4
		f = cal(x, cat(enl(f), enl(ati(y, w(i)))))
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
	for i := int32(0); i < x; i++ {
		SetI32(p+4*i, w(i))
	}
	return r
}
func cnt(x int32) int32 {
	xn := n(x)
	if xn < 0 {
		return 2
	} else {
		return w(xn)
	}
}
func zer(x int32) int32 { return mk(v(x)) } // &a
func wer(x int32) int32 { // &v
	r := mk(0)
	xn := n(x)
	for i := int32(0); i < xn; i++ {
		j := w(i)
		r = cat(r, rep(j, ati(x, j)))
	}
	return r
}
func not(x int32) int32 { return w(I32B(v(1) == 0)) } // ~a
func ntv(x int32) int32 { return ech(252, x) }        // ~v
func fst(x int32) int32 { return ati(x, 0) }          // *a  *v

func ec2(f, x, y int32) int32 {
	rn := max(cnt(x), cnt(y)) >> 1
	r := mk(rn)
	p := v(r)
	for i := int32(0); i < rn; i++ {
		SetI32(p+4*i, cal(f, cat(enl(ati(x, i)), enl(ati(y, i)))))
	}
	return r
}

func inn(x, y int32) int32 {
	yn := w(n(y))
	if n(x) < 0 {
		return mor(yn, fnx(y, x))
	} else {
		return ovr(124, mor(yn, fnd(y, x))) // |/(#y)>y?x
	}
}

func el(x int32) int32 {
	if n(x)&1 != 0 {
		return x
	}
	return enl(x)
}

func v(x int32) int32 { return x >> 1 }
func w(x int32) int32 { return x << 1 }
func n(x int32) int32 {
	if x&1 != 0 {
		return -1
	}
	return I32(v(x) - 4)
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
	return 1 + w(r)
}
func lup(x int32) int32 {
	return x //nyi
}
func prs(x, y int32) int32 {
	return x //nyi
}
