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
	Functions(00, asn, mtv, cat, cts, tkv, dpv, fnd, atx, spc, inn, spl, jon)                                    //vy
	Functions(12, asn, mtc, cat, cut, tak, drp, fna, cal, spc, ech, ovr, scn)                                    //ay
	Functions(24, add, sub, mul, div, min, max, les, mor, eql, mod)                                              //scalar dyadic
	Functions(34, flp, neg, idn, rot, wer, rev, gup, gdn, grp, til, idn, not, enl, srt, cnt, lst, unq, val, enl) //monadic
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
	if n(x)*n(y) == 0 {
		return mk(0)
	}
	rn := max(cnt(x), cnt(y)) >> 1
	r := rm(rn)
	for i := int32(0); i < rn; i++ {
		c1(r, cal(f, cat(enl(atx(x, w(i))), enl(atx(y, w(i))))))
	}
	return r
}
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
	}
	return y
}
func max(x, y int32) int32 {
	if x > y {
		return x
	}
	return y
}
func les(x, y int32) int32 { return I32B(x < y) }
func mor(x, y int32) int32 { return I32B(x > y) }
func eql(x, y int32) int32 { return I32B(x == y) }

func ny2(x, y int32) int32 { return y }

func asn(x, y int32) int32 { return ny2(x, y) }
func mtc(x, y int32) int32 { return w(I32B(x == y)) } // a~y
func mtv(x, y int32) int32 { // v~y
	if x == y {
		return 3
	}
	if n(x) != n(y) {
		return 1
	}
	return ovr(77, ec2(253, x, y)) // &/x~'y
}
func cat(x, y int32) int32 { // x,y
	x, y = el(x), el(y)
	xn, yn := n(x), n(y)
	r := mk(xn + yn)
	Memorycopy(r, x, 4*xn)
	Memorycopy(r+4*xn, y, 4*yn)
	return r
}
func cut(x, y int32) int32 { return l2(atx(y, til(x)), drp(x, y)) } // a^y
func cts(x, y int32) int32 { return ec2(189, x, enl(el(y))) }       // v^y
func tak(x, y int32) int32 { return atx(el(y), til(x)) }            // a#v
func tkv(x, y int32) int32 { return atx(y, wer(inn(el(y), x))) }    // v#y
func drp(x, y int32) int32 { // a_y
	y = el(y)
	return atx(y, tak(w(max(0, n(y)-abs(v(x)))), seq(n(y), max(0, v(x)))))
}
func dpv(x, y int32) int32 { return atx(y, wer(not(inn(el(y), x)))) } // v_y
func fna(x, y int32) int32 { return fnd(enl(x), y) }                  // a?y
func fnd(x, y int32) int32 { // v?y
	if n(y) < 0 {
		i := int32(0)
		for ; i < n(x); i++ {
			if mtc(I32(x+4*i), y) == 3 {
				break
			}
		}
		return w(i)
	}
	return ec2(127, enl(x), y)
}
func cal(x, y int32) int32 { // a.a  a.v
	//	println("cal", tostring(x), tostring(y))
	y = el(y)
	yn := n(y)
	i := int32(0)
	for ; i < 21; i++ {
		if x == I8(i)&0xff {
			break
		}
	}
	//	println("cali", i)
	if i > 20 {
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
		return Func[i+34].(f1)(x)
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
	yn := n(y)
	if yn < 0 {
		return cal(x, enl(y))
	}
	r := rm(yn)
	for i := int32(0); i < yn; i++ {
		c1(r, cal(x, enl(I32(y+4*i))))
	}
	return r
}
func inn(x, y int32) int32 {
	yn := w(n(y))
	if n(x) < 0 {
		return mor(yn, fnd(y, x))
	} else {
		return ovr(124, mor(yn, fnd(y, x))) // |/(#y)>y?x
	}
}
func ovr(x, y int32) int32 { // a/a  a/v
	y = el(y)
	r := fst(y)
	for i := int32(1); i < n(y); i++ {
		r = cal(x, l2(r, I32(y+4*i)))
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
func flp(x int32) int32 { // +x   m:|/#'x   r:(,/m#/:x)(!m)+\:m*!n
	x = el(x)
	xn := cnt(x)
	xm := ovr(249, ech(71, x))
	return atx(ovr(89, ec2(71, enl(xm), x)), ec2(87, til(xm), enl(cal(85, l2(xm, til(xn))))))
}
func neg(x int32) int32 { return cal(91, l2(1, x)) }               // -x
func fst(x int32) int32 { return atx(x, 1) }                       // *x
func rot(x int32) int32 { x = el(x); return atx(x, seq(n(x), 1)) } // %x
func wer(x int32) int32 { // &x
	x = el(x)
	r := mk(0)
	for i := int32(1); i < w(n(x)); i += 2 {
		r = cat(r, tak(atx(x, i), i))
	}
	return r
}
func rev(x int32) int32 { // |x
	x = el(x)
	xn := n(x)
	return atx(x, cal(91, l2(w(xn), seq(xn, 1))))
}
func grd(x, c int32) int32 { // <x  todo tao
	x = el(x)
	xn := n(x)
	r := til(x)
	for i := int32(1); i < xn; i++ {
		ri := I32(r + 4*i)
		j := i - 1
		for j >= 0 {
			if Func[c].(f2)(I32(x+4*v(I32(r+4*j))), I32(x+4*v(ri))) == 0 {
				break
			}
			jj := r + 4*j
			SetI32(4+jj, I32(jj))
			j--
		}
		SetI32(r+4+4*j, ri)
		continue
	}
	return r
}
func gup(x int32) int32 { return grd(x, 31) } // <x
func gdn(x int32) int32 { return grd(x, 30) } // >x
func grp(x int32) int32 { // =x
	x = el(x)
	k := unq(x)
	v := rm(n(k))
	for i := int32(0); i < n(k); i++ {
		c1(v, wer(ec2(253, x, I32(k+4*i))))
	}
	return l2(k, v)
}
func til(x int32) int32 { // !x  !v(domain)
	if x&1 != 0 {
		return seq(v(x), 0)
	}
	return til(cnt(x))
}
func idn(x int32) int32 { return x } // :x
func not(x int32) int32 { // ~a
	if n(x) < 0 {
		return w(I32B(v(x) == 0))
	}
	return ech(253, x)
}
func enl(x int32) int32 { r := mk(1); SetI32(r, x); return r } // ,x
func srt(x int32) int32 { return atx(x, gup(x)) } // ^x
func cnt(x int32) int32 { // #x
	xn := n(x)
	if xn < 0 {
		return 3
	} else {
		return w(xn)
	}
}
func lst(x int32) int32 { // _x  last
	xn := n(x)
	if xn < 0 {
		return x
	}
	return atx(x, w(n(x)-1)) 
}
func unq(x int32) int32 { // ?x
	x = el(x)
	r := rm(x)
	for i := int32(0); i < n(x); i++ {
		xi := I32(x + 4*i)
		if v(fnd(r, xi)) == n(r) {
			c1(r, xi)
		}
	}
	return r
}
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
