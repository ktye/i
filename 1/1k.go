package main

import (
	. "github.com/ktye/wg/module"
)

var tot, top int32

func init() {
	Memory(1)
	// no @ $ /: \: ': ` [ ]  { }  may be assigned to as user functions but not infix

	// x(space)y          is (,x),y
	// e.g.:   l:4 5 6    is (4;5;6)
	//         1 l 2      is (1;4 5 6;2)  or as:  1 (4 5 6) 2      (space matters)
	//
	// "abc" is quotation/literal, also used for lambdas: f:"(+/x)%#x"
	//
	// juxtaposition
	// ay cal
	// vy idx
	Data(0, "\x57\x5b\x55\x4b\x4d\xf9\x79\x7d\x7b\x43\x75\xfd\x59\xbd\x47\xbf\x7f\x5d\x41\x4f\x5f\xb9") // (1+2*)  +-*%&|<>=!:~,^#_?. '/\  (22)

	//              0    1    2    3    4    5    6    7    8    9   10   11
	//              :    ~    ,    ^    #    _    ?    .         '    /    \
	Functions(00, asn, mtc, cat, cts, tkv, dpv, fnd, atx, spc, inn, spl, jon)                                    //vy
	Functions(12, asn, mtc, cat, cut, tak, drp, rol, cal, spc, ech, ovr, scn)                                    //ay
	Functions(24, add, sub, mul, div, min, max, les, mor, eql, mod)                                              //scalar dyadic
	Functions(34, flp, neg, idn, rot, wer, rev, grd, gdn, grp, til, idn, not, enl, str, cnt, lst, unq, val, enl) //monadic
	//              +    -    *    %    &    |    <    >    =    !    :    ~    ,    ^    #    _    ?    .  spc    '   /  \
	//             43   45   42   37   38  124   60   62   61   33   58  126   44   94   35   95   63   46   32   39  47  92
}
func main() {
	k1()
}
func k1() {
	tot = 65536
	top = 24
	rm(256) //keys at 24
	rm(256) //vals at 270
}

func v(x int32) int32 { return x >> 1 }
func w(x int32) int32 { return (x << 1) + 1 }
func n(x int32) int32 {
	if x&1 != 0 {
		return -1
	}
	return I32(x - 4)
}
func mk(x int32) int32 {
	r := top
	top += 4 + 4*x
	for tot < top {
		Memorygrow(1)
		tot += 65536
	}
	SetI32(r, x)
	return 4 + r
}
func rm(x int32) int32 { // reset make, use with c1
	r := mk(x)
	SetI32(r-4, 0)
	return r
}
func c1(x, y int32) {
	p := I32(x - 4)
	SetI32(x-4, 1+p)
	SetI32(x+4*p, y)
}
func l2(x, y int32) int32 { return cat(enl(x), enl(y)) }
func el(x int32) int32 {
	if x&1 != 0 {
		return enl(x)
	}
	return x
}
func ec2(f, x, y int32) int32 {
	rn := max(cnt(x), cnt(y)) >> 1
	r := rm(rn)
	for i := int32(0); i < rn; i++ {
		c1(r, cal(f, cat(enl(atx(x, w(i))), enl(atx(y, w(i))))))
	}
	return r
}

/*
	func seq(x, o, m int32) int32 {
		r := rm(x)
		for i := int32(0); i < x; i++ {
			c1(r, w((i+o)%m))
		}
		return r
	}
*/
func seq(x, o int32) int32 {
	r := rm(abs(x))
	for i := int32(0); i < abs(x); i++ {
		c1(r, w(mod(i+o, x)))
	}
	return r
}
func lup(x int32) int32 {
	return x //nyi
}

type f1 = func(int32) int32
type f2 = func(int32, int32) int32

// dyadic
func add(x, y int32) int32 { return x + y }
func sub(x, y int32) int32 { return x - y }
func mul(x, y int32) int32 { return x * y }
func div(x, y int32) int32 { return x / y }
func mod(x, y int32) int32 {
	x = x % y           // some targets need return 0 for y=0
	if x < 0 || y < 0 { // simplified (no 0>y)
		return x + y
	}
	return x
}
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
	Memorycopy(r, x, 4*xn)
	Memorycopy(r+4*xn, y, 4*yn)
	return r
}
func cut(x, y int32) int32 { return ny2(x, el(y)) }              // a^y
func cts(x, y int32) int32 {              return ny2(x, y) } // v^y
func tak(x, y int32) int32 { return atx(el(y), til(x)) }         // a#v
func tkv(x, y int32) int32 { return atx(y, wer(inn(el(y), x))) } // v#y
func drp(x, y int32) int32 { // a_y
/*
	x = v(x)
	yn := n(y)
	rn := max(0, yn-abs(x))
	o := mod(x+yn, yn)
	r := mk(rn)
	Memorycopy(r, y+o, 4*rn)
	return r
*/	

	//y = el(y)
	// n(y)-v(x)
	// y:0 1 2
	// 2_y
	// o:(yn-x)
	// n:x
	//println("x/vx",x,v(x))
	//println("seq",tostring(seq(n(y)-v(x),v(x))))


//	1+!yn   1 2 3
//	(yn-x)# 1 2

//	fst(cts(enl(x),til(y)))  // *(,x)^!y

//	rev(cut(x,til(y)))

	y = el(y)
	return atx(y, tak(w(max(0,n(y)-v(x))), seq(n(y), v(x))))
}
/*
	//println("drp", tostring(x), tostring(y))
	y = el(y)
	yn := n(y)
	xv := v(x)
	if x < 0 {
		if xv < -yn {
			panic("aaa")
			//			return atx(x, seq(-xv, xv+yn, yn))
		}
		panic("bbb")
		return tak(x, w(-xv))
	} else {
		panic("drop")
		//		return atx(y, seq(max(0,yn-xv), xv, yn))
	}
}
*/
func dpv(x, y int32) int32 { return atx(y, wer(not(inn(el(y), x)))) } // v_y
func rol(x, y int32) int32 { return ny2(x, y) }                       // a?y
func fnd(x, y int32) int32 { return fst(ec2(126, x, y)) }             // v?a
func fnx(x, y int32) int32 { return ec2(63, enl(x), y) }              // v?v   x?/:y
func cal(x, y int32) int32 { // a.a  a.v
	//	println("cal", tostring(x), tostring(y))
	y = el(y)
	yn := n(y)
	i := int32(0)
	for ; i < 21; i++ {
		if x == I8(i) {
			break
		}
	}
	//	println("cali", i)
	if i > 20 {
		panic("exe lup")
		return exe(lup(x), y)
	}
	x = fst(y)
	y = atx(y, 3)
	xa := I32B(n(x) < 0)
	ya := I32B(n(y) < 0)
	mo := I32B(yn < 2)
	if i < 10 { //scalar
		if mo != 0 {
			return Func[i+34].(f1)(x)
		}
		if xa&ya != 0 {
			return w(Func[i+24].(f2)(v(x), v(y)))
		}
		return ec2(I8(i), x, y)
	}
	if mo != 0 {
		panic("mo")
		return Func[i+48].(f1)(x)
	}
	i = (i - 10) + 12*xa
	return Func[i].(f2)(x, y)
}
func atx(x, y int32) int32 { // v.a  (also a.v)
	//println("atx", tostring(x), tostring(y))
	if y&1 != 0 {
		xn := n(x)
		if xn > 0 {
			return I32(x + 4*mod(v(y), xn))
		}
		// todo xn<0: does this happen?
		return x*I32B(xn < 0) + I32B(xn == 0)
	}
	return ec2(93, enl(x), y)
}
func spc(x, y int32) int32 { return cat(enl(x), y) } // x(space)y
func ech(x, y int32) int32 { // a'a  a'v
	yn := v(cnt(y))
	r := rm(yn)
	for i := int32(0); i < yn; i++ {
		c1(r, cal(x, atx(y, w(i))))
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
func ovr(x, y int32) int32 { // a/a  a/v
	yn := v(cnt(y))
	r := fst(x)
	for i := int32(1); i < yn; i++ {
		r = cal(x, atx(y, w(i)))
	}
	return r
}
func spl(x, y int32) int32 { return cut(cat(0, wer(ec2(126, enl(x), fst(y)))), x) } // v/a    (0,&(,x)~'*y)^x
func scn(x, y int32) int32 { // a\a  a\v
	yn := v(cnt(y))
	r := mk(yn)
	p := v(r)
	f := fst(x)
	SetI32(p, f)
	for i := int32(0); i < yn; i++ {
		p += 4
		f = cal(x, cat(enl(f), enl(atx(y, w(i)))))
		SetI32(p, f)
	}
	return r
}
func jon(x, y int32) int32 { return ovr(44, ec2(44, x, enl(y))) } // v\a   ,/x,',y

func ecv(x, y int32) int32 { return ec2(39, x, enl(y)) } // v'a  v'a

func ovv(x, y int32) int32 { return ec2(47, x, enl(y)) } // v/a  v/v

// monadic
func nyi(x int32) int32 { return x }
func abs(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}
func flp(x int32) int32 { return nyi(x) } // +x
func neg(x int32) int32 { // -x
	if n(x) < 0 {
		return w(-v(x))
	}
	return ech(91, x)
}
func fst(x int32) int32 { return atx(x, 1) }              // *x
func rot(x int32) int32 { return cat(drp(2, x), fst(x)) } // %x
func wer(x int32) int32 { // &x
	x = el(x)
	r := mk(0)
	xn := n(x)
	for i := int32(0); i < xn; i++ {
		j := w(i)
		r = cat(r, tak(j, atx(x, j)))
	}
	return r
}
func rev(x int32) int32 { // |x
	x = el(x)
	xn := n(x)
	r := rm(xn)
	for i := int32(0); i < xn; i++ {
		c1(r, xn-i-1)
	}
	return atx(x, r)
}
func grd(x int32) int32 { return nyi(el(x)) }   // <x
func gdn(x int32) int32 { return nyi(el(x)) }   // >x
func grp(x int32) int32 { return nyi(el(x)) }   // =x
func til(x int32) int32 { // !x  !v(domain)
	if x&1 != 0 {
		return seq(v(x), 0) 
	}
	return til(cnt(x))
}
func idn(x int32) int32 { return x }            // :x
func not(x int32) int32 { // ~a
	if n(x) < 0 {
		return w(I32B(v(x) != 0))
	}
	return ech(252, x)
}
func enl(x int32) int32 { r := mk(1); SetI32(r, x); return r } // ,x
func str(x int32) int32 { // ^x
	if n(x) < 0 {
		x = v(x)
		if x < 0 {
			return cat(90, str(w(-x)))
		} else if x == 0 {
			return enl(0)
		}
		r := rm(10)
		i := int32(0)
		for x > 0 {
			c1(r, x%10)
			x -= 10
			x /= 10
			i++
		}
		return tak(w(i), r)
	}
	return ech(188, x)
}
func cnt(x int32) int32 { // #x
	xn := n(x)
	if xn < 0 {
		return 3
	} else {
		return w(xn)
	}
}
func lst(x int32) int32 { return atx(el(x), w(n(x)-1)) } // _x
func unq(x int32) int32 { return nyi(x) }                // ?x
func val(x int32) int32 { // .x
	xn := n(x)
	if xn < 0 {
		return lup(x)
	}
	return exe(x, 0)
}
func tok(x int32) int32 {
	xn := n(x)
	r := rm(xn)
	x = v(x)
	xe := x + 4*xn
	t := int32(0)
	q := rm(xn)
	for x < xe {
		c := I32(x)
		if 0 < n(q) {
			if c == 68 { // "
				c1(r, q)
				q = rm(xn)
				continue
			}
			c1(q, c)
			continue
		}
		if c >= 96 && c <= 114 {
			t *= 10
			t += v(c - 36)
			continue
		}
		c1(r, x)
		x += 4
	}
	return r
}

// func prs(r, x, t int32) int32 {
// 11-(1+1)+-1
//   - 1)     >   A       ? /    >d
//   - -1)     m   B       - -    m>
//     )  +1)     >   C       1 -    >d
//     1  )+1)    >   D       1 1    c
//     ( ?    <>
//   - 1)+1)   >   A     else     >
//     1  +1)+1)  d   E
//     (  1)+1)   b<  F      -/1()
//     1    +1)   d   E
//   - 1)    >   A
//     1    -1)   d   E
//     1    1)    j   G
//     1)
//   - +)    m
//
// x == (      => move<y over (, remove )  next
// y == /      => move> dyadic
// y == -
//
//	| x == 1  => move> dyadic
//	| x == -  => monadic, move>
//
// else        => move
// }
func exe(x, a int32) int32 {
	x = rev(tok(x))
	return x
}
