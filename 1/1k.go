package main

import (
	. "github.com/ktye/wg/module"
	"strconv"
	"strings"
)

var tot, top, loc int32

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
	//Data(0, "\x57\x5b\x55\x4b\x4d\xf9\x79\x7d\x7b\x43\x75\xfd\x59\xbd\x47\xbf\x7f\x5d\x41\x4f\x5f\xb9") // (1+2*)  +-*%&|<>=!:~,^#_?. '/\  (22)
	Data(0, "\x2b\x2d\x2a\x25\x26\x7c\x3c\x3e\x3d\x21\x3a\x7e\x2c\x5e\x23\x5f\x3f\x2e\x20\x27\x2f\x5c") // +-*%&|<>=!:~,^#_?. '/\  (22)

	//              0    1    2    3    4    5    6    7    8    9   10   11
	//              :    ~    ,    ^    #    _    ?    .         '    /    \
	Functions(00, asn, mtv, cat, cts, tkv, dpv, fnd, atx, spc, ech, ovr, scn)                                    //vy
	Functions(12, asn, mtc, cat, cut, tak, drp, fna, cal, spc, ech, ovr, scn)                                    //ay
	Functions(24, add, sub, mul, div, min, max, les, mor, eql, mod)                                              //scalar dyadic
	Functions(34, flp, neg, idn, rot, wer, rev, gup, gdn, grp, til, idn, not, enl, tok, cnt, las, unq, val, out) //monadic
	//              +    -    *    %    &    |    <    >    =    !    :    ~    ,    ^    #    _    ?    .  spc    '   /  \
	//             43   45   42   37   38  124   60   62   61   33   58  126   44   94   35   95   63   46   32   39  47  92
}
func main() {
	k1()
	for {
		rep()
	}
}
func k1() {
	tot = 65536
	top = 24
	loc = 120
	rm(22) //22 primitive symbols at 28 +-*%&|<>=!:~,^#_?. '/\
	for n(28) < 22 {
		c1(28, w(I8(n(28))))
	}
	rm(63) //63 globals at 120
	//top is 372
}
func rep() {
	r := rm(126)
	x := mk(1)
	for ReadIn(x, 1) != 0 {
		c1(r, w(I8(x)))
	}
	c1(r, 21)
	//out(exe(r, 0))
	//o(10)
	exe(r, 0)
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
	sn(r, 0)
	return r
}
func sn(x, y int32) { set(x, -1, y) }
func c1(x, y int32) {
	p := I32(x - 4)
	sn(x, 1+p)
	set(x, p, y)
}
func cp(x, y, n int32) {
	for i := int32(0); i < 4*n; i += 4 {
		set(x, i, get(y, i))
	}
}
func l2(x, y int32) int32 { return cat(enl(x), enl(y)) }
func el(x int32) int32 {
	if x&1 != 0 {
		return enl(x)
	}
	return x
}
func fs(x int32) int32 {
	if x&1 == 0 {
		if x == 0 {
			return 0
		}
		return fs(I32(x))
	}
	return x
}
func fl(x int32) int32 {
	a := int32(1)
	for i := int32(0); i < n(x); i++ {
		a = min(a, I32B(uint32(v(get(x, i))-10) < 117))
	}
	return a
}
func iv(x int32) int32     { return v(fnd(28, x)) }
func set(x, i, y int32)    { SetI32(x+4*i, y) }
func get(x, i int32) int32 { return I32(x + 4*i) }
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

func asn(x, y int32) int32 { // a:y  (a i):y
	if n(x) < 0 {
		x = v(x) - 65
		if uint32(x) < 63 {
			set(loc, x, y) //asn is always local
		}
		return y
	}
	i := drp(3, x)
	x = fs(x)
	return asn(x, amd(val(x), i, y))
}
func amd(x, i, y int32) int32 { // @[x;i;y]
	x = el(x)
	xn := n(x)
	r := mk(xn) //amd is the only place that needs to copy
	if xn > 0 {
		Memorycopy(r, x, 4*xn)
		if n(i) < 0 {
			i = enl(i)
			y = enl(y)
		}
		y = tak(cnt(i), y)
		for j := int32(0); j < n(i); j++ {
			set(r, mod(v(get(i, j)), xn), get(y, j))
		}
	}
	return r
}
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
			if mtc(get(x, i), y) == 3 {
				break
			}
		}
		return w(i)
	}
	return ec2(127, enl(x), y)
}
func cal(x, y int32) int32 { // a.a  a.v    +derived (+ /) y
	//println("cal", tostring(x), tostring(y), "arity", n(y))
	if x&1 == 0 { //derived  x:(/ +) y => x:/ y:(+ y)
		y = l2(atx(x, 3), fst(y))
		x = fst(x)
	}
	y = el(y)
	yn := n(y)
	i := iv(x)
	//	println("cali", i)
	if i == 22 { // not a primitive
		//println("lambda")
		return exe(val(x), y)
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
		return ec2((get(28, i)), x, y)
	}
	if mo != 0 {
		return Func[i+34].(f1)(x) //monadic
	}
	i = (i - 10) + 12*xa
	//println("dyadic", i)
	return Func[i].(f2)(x, y) //dyadic
}
func atx(x, y int32) int32 { // v.a  (also a.v)
	// println("atx", tostring(x), tostring(y))
	if y&1 != 0 {
		xn := n(x)
		if xn > 0 {
			return get(x, mod(v(y), xn))
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
		c1(r, cal(x, enl(get(y, i))))
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
		r = cal(x, l2(r, get(y, i)))
	}
	return r
}

// func spl(x, y int32) int32 { return cut(cat(0, wer(ec2(126, enl(x), fst(y)))), x) } // v/a    (0,&(,x)~'*y)^x
func scn(x, y int32) int32 { // a\a  a\v
	yn := v(cnt(y))
	r := rm(yn)
	if yn != 0 {
		c1(r, fst(y))
		for i := int32(1); i < yn; i++ {
			c1(r, cal(x, l2(las(r), get(y, i))))
		}
	}
	return r
}

//func jon(x, y int32) int32 { return ovr(44, ec2(44, x, enl(y))) } // v\a   ,/x,',y

func ecv(x, y int32) int32 { return ec2(39, x, enl(y)) } // v'a  v'a

func ovv(x, y int32) int32 { return ec2(47, x, enl(y)) } // v/a  v/v

// monadic
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
		ri := get(r, i)
		j := i - 1
		for j >= 0 {
			if Func[c].(f2)(get(x, v(get(r, j))), get(x, v(ri))) == 0 {
				break
			}
			jj := r + 4*j
			SetI32(4+jj, I32(jj))
			j--
		}
		set(r+4, j, ri)
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
		c1(v, wer(ec2(253, x, get(k, i))))
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
func cnt(x int32) int32 { // #x
	xn := n(x)
	if xn < 0 {
		return 3
	}
	return w(xn)
}
func las(x int32) int32 { // _x  last
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
		xi := get(x, i)
		if v(fnd(r, xi)) == n(r) {
			c1(r, xi)
		}
	}
	return r
}
func val(x int32) int32 { // .x
	xn := n(x)
	if xn < 0 { // lup
		r := lup(x, loc) //try local
		if r != 0 {
			return r
		}
		r = lup(x, 120) //try global
		if r != 0 {
			return r
		}
		return 1
	}
	return exe(x, 0)
}
func lup(x, env int32) int32 { return get(env, v(x)-65) }
func tok(x int32) int32 { // ^x tokenize, enlists (123 and "abc")
//println("tok", tostring(x))
	t := int32(0)
	x = el(x)
	xn := n(x)
	r := rm(xn)
	q := rm(xn)
	xn = x + 4*xn
	for x < xn {
		c := I32(x)
		x += 4
		if 0 < n(q) {
			if c == 69 { // "
				c1(r, enl(drp(3, q)))
				sn(q, 0)
				continue
			}
			c1(q, c)
			continue
		} else if c == 69 {
			c1(q, 1)
			continue
		}
		if c >= 97 && c <= 115 {
			t *= 10
			t += v(c) - 48
			if t == 0 {
				c1(r, enl(1))
			}
			continue
		}
		if t != 0 {
			c1(r, enl(w(t)))
			t = 0
		}
		if c > 21 {
			c1(r, c)
		}
	}
	return r
}

func exe(x, args int32) int32 { //parse and execute
//println("exe", tostring(x))
	sv, sp := loc, top

	if args != 0 {
		a := wer(127)
		for i := int32(0); i < n(args); i++ {
			set(a, i+55, get(args, i)) // xyz..
			println("setloc", a, i+55, tostring(get(args, i)))
		}
		loc = a
	}

	x = rev(tok(x))
//println("rev tok", tostring(x))
	x = fst(e(t(x), x))

	if args != 0 {
		loc, top = sv, sp
		xn := n(x)
		if xn >= 0 {
			rm(xn) //sp+4
			cp(sp+4, x, xn)
		}
	}
	return x
}
func e(x, b int32) int32 {
	if x == 0 {
		return 0
	}
	y := t(b)
	if y == 0 {
		return enl(vau(x))
	}
	var r int32
	if ver(y) != 0 && ver(x) == 0 {
		r = vau(e(t(b), b))
		if y == 117 {
			return enl(asn(x, r))
		}
		//println("dyadic", tostring(x), tostring(y), tostring(r))

		//r = vau(r)
		return enl(cal(y, l2(vau(x), r))) //dyadic
	}

	r = vau(e(y, b))
	if ver(x) == 0 { // juxtaposition
		//println("jux")
		//x=vau(x)
		return enl(cal(93, l2(x, r)))
	}
	println("monadic", tostring(x), tostring(r))
	if las(x) == 0 {
		x = tak(5, x) //unmark derived
	}
	//println("monadic", tostring(r), tostring(x))
	return enl(cal(x, enl(r)))
}
func t(x int32) int32 {
	r := nxt(x)
	if r == 0 || r == 83 {
		return 0
	}
	if r == 81 { // (
		return e(t(x), x)
	}
	d := int32(0)
	for {
		y := nxt(x)
		if y == 0 {
			break
		}
		if adv(y) != 0 {
			r = l2(y, r)
			d = 1
		} else {
			sn(x, n(x)+1)
			break
		}
	}
	if d != 0 {
		return cat(r, enl(0)) //mark derived
	}
	return r
}
func nxt(x int32) int32 {
	if n(x) == 0 {
		return 0
	}
	sn(x, n(x)-1)
	return get(x, n(x))
}
func ver(x int32) int32 {
	if x&1 == 0 {
		return I32B(las(x) == 0) // derived
		//return 0
	}
	return I32B(iv(x) < 22)
}
func adv(x int32) int32 {
	if ver(x) != 0 {
		return I32B(iv(x) > 18) // 19 20 21
	}
	return 0
}

func vau(x int32) int32 { // combien vaut-il?
	if x&1 != 0 {
		return val(x)
	}
	return fst(x)
}

// print
func out(x int32) int32 { 
	println("out", tostring(x))
	ou(x); o(10); return x }
func ov(x int32) {
	for i := int32(0); i < n(x); i++ {
		o(v(get(x, i)))
	}
}
func oi(x int32) {
	x = v(x)
	if x == 0 {
		o(48)
		return
	}
	if x < 0 {
		o(45)
		oi(w(-x))
		return
	}
	r := rm(10)
	for x != 0 {
		c1(r, w(48+x%10))
		x /= 10
	}
	ov(rev(r))
}
func ou(x int32) {
	if x&1 != 0 {
		oi(x)
	} else {
		a := fl(x)
		if a == 1 {
			o(34)	
			ov(x)
			o(34)	
		} else {
			o(40)
			for i := int32(0); i < n(x); i++ {
				if i > 0 {
					o(32)
				}
				ou(get(x, i))
			}
			o(41)
		}
	}
}
func o(x int32) { print(string(x)) }

func tostring(x int32) string { //todo rm
	if x == 0 {
		return "null"
	} else if x&1 != 0 {
		return strconv.Itoa(int(x >> 1))
	} else {
		xn := n(x)
		if xn < 0 || xn > 30 {
			panic("xn")
		}
		a := make([]byte, xn)
		for i := int32(0); i<xn; i++ {
			if xi := v(get(x, i)); xi > 31 && xi < 128 {
				a[i] = byte(xi)
			} else {
				a = nil
				break
			}
		}
		if a != nil {
			return `"` + string(a) + `"`
		}
		u := make([]string, xn)
		for i := int32(0); i < xn; i++ {
			u[i] = tostring(I32(x + 4*i))
		}
		return "(" + strings.Join(u, " ") + ")"
	}
}
