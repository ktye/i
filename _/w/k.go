// +build ignore

// prototype implementation for k.w (same memory layout)
// go run k.go t      /all tests
// go run k.go EXPR

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"math/bits"
	"math/cmplx"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

const trace = false
const parse = false

type c = byte
type s = string
type i = uint32
type j = uint64
type f = float64

var MC []c // MC, MI, MJ, MF share array (see msl)
var MI []i
var MJ []j
var MF []f
var MT [256]interface{}

type vt1 func(i) i
type vt2 func(i, i) i
type slice struct {
	p uintptr
	l int
	c int
}

const naI i = 2147483648
const naJ j = 9221120237041090561
const pp, kkey, kval, xyz, cmap = 8, 132, 136, 148, 160

func main() {
	if len(os.Args) == 2 && os.Args[1] == "t" {
		runtest()
	} else {
		fmt.Println(run(strings.Join(os.Args[1:], " ")))
	}
}
func runtest() {
	b, e := ioutil.ReadFile("t")
	if e != nil {
		panic(e)
	}
	v := strings.Split(strings.TrimSpace(string(b)), "\n")
	for i := range v {
		if len(v[i]) == 0 {
			fmt.Println("skip rest")
			os.Exit(0)
		}
		if len(v[i]) == 0 || v[i][0] == '/' {
			continue
		}
		vv := strings.Split(v[i], " /")
		if len(vv) != 2 {
			panic("test file")
		}
		in := strings.TrimRight(vv[0], " \t\r")
		exp := strings.TrimSpace(vv[1])
		got := run(in)
		fmt.Println(in, "/", got)
		if exp != got {
			fmt.Println("expected:", exp)
			os.Exit(1)
		}
	}
}
func run(s string) string {
	m0 := 16
	MJ = make([]j, (1<<m0)>>3)
	msl()
	ini(16)
	x := mk(1, i(len(s)))
	copy(MC[x+8:], s)
	if parse {
		x = prs(x)
	} else {
		x = prs(x)
		x = evl(x, 0)
	}
	s = kst(x)
	dx(x)
	leak()
	return s
}
func ini(x i) i {
	copy(MT[0:], []interface{}{
		//   1    2    3    4    5    6    7    8    9    10   11   12   13   14   15
		nil, gtc, gti, gtf, gtl, gtl, nil, nil, nil, eqc, eqi, eqf, eqz, eqL, eqL, nil, abc, abi, abf, abz, nec, nei, nef, nez, nil, nil, nil, nil, sqc, sqi, sqf, sqz, // 000..031
		nil, mkd, nil, rsh, cst, diw, min, ecv, ecd, epi, mul, add, cat, sub, cal, ovv, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, dex, nil, les, eql, mor, fnd, // 032..063
		atx, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, ecl, scv, nil, exc, cut, // 064..095
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, ecr, max, nil, mtc, nil, // 096..127
		nil, nil, nil, nil, nil, nil, nil, nil, nms, vrb, chr, nam, sms, nil, nil, nil, adc, adi, adf, adz, suc, sui, suf, suz, muc, mui, muf, muz, dic, dii, dif, diz, // 128..159
		nil, til, nil, cnt, str, sqr, wer, epv, ech, ecp, fst, abs, enl, neg, val, riv, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, lst, nil, grd, grp, gdn, unq, // 160..191
		typ, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, scn, liv, spl, srt, flr, // 192..223
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, prs, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, ovr, rev, jon, not, nil, // 224..255
	})
	sJ(0, 289360742959022340) // type sizes uint64(0x0404041008040104)
	sI(128, x)                // alloc
	p := i(256)
	for i := i(8); i < x; i++ {
		sI(4*i, p) // free pointer
		p *= 2
	}
	sI(kkey, mk(5, 0))
	sI(kval, mk(6, 0))
	sI(xyz, cat(cat(mks(120), mks(121)), mks(122)))
	return x
}
func msl() { // update slice headers after set/inc MJ
	cp := *(*slice)(unsafe.Pointer(&MC))
	ip := *(*slice)(unsafe.Pointer(&MI))
	jp := *(*slice)(unsafe.Pointer(&MJ))
	fp := *(*slice)(unsafe.Pointer(&MF))
	fp.l, fp.c, fp.p = jp.l, jp.c, jp.p
	ip.l, ip.c, ip.p = jp.l*2, jp.c*2, jp.p
	cp.l, cp.c, cp.p = ip.l*4, ip.c*4, ip.p
	MF = *(*[]f)(unsafe.Pointer(&fp))
	MI = *(*[]i)(unsafe.Pointer(&ip))
	MC = *(*[]c)(unsafe.Pointer(&cp))
	cmake() // copy character map data
}
func bk(t, n i) (r i) {
	r = i(32 - bits.LeadingZeros32(7+n*i(C(t))))
	if r < 4 {
		return 4
	}
	return r
}
func mk(x, y i) (r i) {
	t := bk(x, y)
	i := 4 * t
	for I(i) == 0 {
		i += 4
	}
	if i == 128 {
		panic("Ω")
	}
	a := I(i)
	sI(i, I(a))
	for j := i - 4; j >= 4*t; j -= 4 {
		u := a + 1<<(j>>2)
		sI(u, I(j))
		sI(j, u)
	}
	sI(a, y|x<<29)
	sI(a+4, 1)
	return a
}
func fr(x i) {
	xt, xn, _ := v1(x)
	t := 4 * bk(xt, xn)
	sI(x, I(t))
	sI(t, x)
}
func dx(x i) {
	if x > 255 {
		xr := I(x + 4)
		sI(x+4, xr-1)
		if xr == 1 {
			xt, xn, xp := v1(x)
			if xt == 0 || xt > 4 {
				for i := i(0); i < xn; i++ {
					dx(I(xp + 4*i))
				}
			}
			fr(x)
		}
	}
}
func rx(x i) { rxn(x, 1) }
func rxn(x, y i) {
	if x > 255 {
		MI[1+x>>2] += y
	}
}
func rl(x i) {
	_, xn, xp := v1(x)
	for i := i(0); i < xn; i++ {
		rx(I(xp))
		xp += 4
	}
}
func rld(x i)          { rl(x); dx(x) }
func dxr(x, r i) i     { dx(x); return r }
func dxyr(x, y, r i) i { dx(x); dx(y); return r }
func mki(i i) (r i)    { r = mk(2, 1); sI(r+8, i); return r }
func mkc(c i) (r i)    { r = mk(1, 1); sC(r+8, byte(c)); return r }
func mks(c i) (r i)    { return sc(mkc(c)) }
func mkd(x, y i) (r i) {
	xt, yt, xn, yn, _, _ := v2(x, y)
	if xn != yn {
		trap()
	}
	if xt != 5 {
		trap()
	}
	if yt != 6 {
		y = lx(y) // explode
	}
	r = l2(x, y) // todo ext?
	sI(r, 2|7<<29)
	return r
}
func l2(x, y i) (r i) {
	r = mk(6, 2)
	sI(r+8, x)
	sI(r+12, y)
	return r
}
func l3(x, y, z i) (r i) {
	r = mk(6, 3)
	sI(r+8, x)
	sI(r+12, y)
	sI(r+16, z)
	return r
}
func tp(x i) i {
	if x < 256 {
		return 0
	}
	return I(x) >> 29
}
func nn(x i) (xn i) {
	if x < 256 {
		return 1
	}
	return I(x) & 536870911
}
func v1(x i) (xt, xn, xp i) { return tp(x), nn(x), 8 + x }
func v2(x, y i) (xt, yt, xn, yn, xp, yp i) {
	xt, xn, xp = v1(x)
	yt, yn, yp = v1(y)
	return
}
func ary(x i) (r i) { // arity
	if x < 128 {
		return 2
	}
	if x < 256 {
		return 1
	}
	n := nn(x)
	if n == 2 {
		return 1 // derived
	}
	if n != 4 {
		fmt.Println("n", n)
		panic("ary")
	}
	return nn(I(x + 16)) // lambda
}
func use(x i) (r i) {
	if I(x+4) == 1 {
		return x
	}
	xt, xn, xp := v1(x)
	r = mk(xt, xn)
	mv(r+8, xp, xn*i(C(xt)))
	dx(x)
	return r
}
func mv(dst, src, n i) { copy(MC[dst:dst+n], MC[src:src+n]) }
func ext(x, y i) (rx, ry i) {
	_, _, xn, yn, _, _ := v2(x, y)
	if xn == yn {
		return x, y
	}
	if xn == 1 && yn > 1 {
		return take(x, yn), y
	}
	if xn > 1 && yn == 1 {
		return x, take(y, xn)
	}
	panic("length")
}
func upxy(x, y i) (i, i) { x = upx(x, y); y = upx(y, x); return x, y }
func upx(x, y i) (r i) {
	xt := tp(x)
	yt := tp(y)
	if xt == yt {
		return x
	}
	if xt == 7 || yt == 7 {
		trap()
	}
	if yt == 6 {
		return lx(x)
	}
	xn := nn(x)
	for xt < yt {
		x = up(x, xt, xn)
		xt++
	}
	return x
}
func up(x, t, n i) (r i) {
	r = mk(t+1, n)
	xp, rp := x+8, r+8
	switch t {
	case 1:
		for i := i(0); i < n; i++ {
			sI(rp, uint32(C(xp+i)))
			rp += 4
		}
	case 2:
		for i := i(0); i < n; i++ {
			sF(rp, float64(int32(I(xp))))
			xp += 4
			rp += 8
		}
	default:
		trap()
	}
	return dxr(x, r)
}
func lx(x i) (r i) { // explode
	xt, n, _ := v1(x)
	if xt == 6 {
		return x
	}
	r = mk(6, n)
	rp := r + 8
	rxn(x, n)
	for i := i(0); i < n; i++ {
		sI(rp, atx(x, mki(i)))
		rp += 4
	}
	return dxr(x, r)
}
func til(x i) (r i) {
	xt, _, xp := v1(x)
	if xt != 2 {
		trap()
	}
	n := I(xp)
	dx(x)
	if ii := int32(n); ii < 0 {
		return tir(i(-ii))
	}
	return seq(0, n, 1)
}
func seq(a, n, s i) (r i) {
	r = mk(2, n)
	rp := r + 8
	for i := i(0); i < n; i++ {
		sI(rp, s*(a+i))
		rp += 4
	}
	return r
}
func tir(n i) (r i) {
	r = mk(2, n)
	rp := 8 + r + 4*(n-1)
	for i := i(0); i < n; i++ {
		sI(rp, i)
		rp -= 4
	}
	return r
}
func rev(x i) (r i) {
	n := nn(x)
	if n == 0 {
		return x
	}
	return atx(x, tir(n))
}
func fst(x i) (r i) {
	xt, _, _ := v1(x)
	if xt == 7 {
		return fst(val(x))
	}
	return atx(x, mki(0))
}
func dex(x, y i) (r i) { return dxr(x, y) } // :[x;y]
func lst(x i) (r i) {
	if tp(x) == 7 {
		return lst(val(x))
	}
	return atx(x, mki(nn(x)-1)) /* TODO k.w differs */
} // ::x
func drop(x, n i) (r i) {
	xt, xn, _ := v1(x)
	if n > xn {
		n = xn
	}
	if xt == 6 && xn-n == 1 {
		return enl(lst(x))
	}
	return atx(x, seq(n, xn-n, 1))
}
func cut(x, y i) (r i) {
	xt, _, xn, yn, xp, _ := v2(x, y)
	if xt != 2 {
		panic("type")
	}
	if xn == 1 {
		n := I(xp)
		return dxr(x, drop(y, n))
	}
	r = mk(6, xn)
	rp := r + 8
	for i := i(0); i < xn; i++ {
		a := I(xp)
		b := I(xp + 4)
		if i == xn-1 {
			b = yn
		}
		if b < a {
			panic("domain")
		}
		rx(y)
		sI(rp, atx(y, seq(a, b-a, 1)))
		xp += 4
		rp += 4
	}
	return dxyr(x, y, r)
}
func rsh(x, y i) (r i) {
	xt, _, xn, _, xp, _ := v2(x, y)
	if xt != 2 {
		panic("type")
	}
	n := prod(xp, xn)
	r = take(y, n)
	if xn == 1 {
		return dxr(x, r)
	}
	xn--
	xe := xp + 4*xn
	for i := i(0); i < xn; i++ {
		m := I(xe)
		n /= m
		n = prod(xp, xn-i)
		r = cut(seq(0, n, m), r)
		xe -= 4
	}
	return dxr(x, r)
}
func prod(xp, n i) (r i) {
	r = 1
	for i := i(0); i < n; i++ {
		r *= I(xp)
		xp += 4
	}
	return r
}
func take(x, n i) (r i) {
	xn := nn(x)
	r = seq(0, n, 1)
	if xn < n {
		rp := 8 + r
		for i := i(0); i < n; i++ {
			sI(rp, I(rp)%xn)
			rp += 4
		}
	}
	return atx(x, r)
}
func atm(x, y i) (r i) { // {$[0~#y;:0#x;];}
	if 0 == nn(y) {
		return dxr(x, y)
	}
	rx(y)
	f := fst(y)
	t := drop(y, 1)
	if tp(t) == 6 && nn(t) == 1 {
		t = fst(t)
	}
	nf := nn(f)
	if f == 0 || nf != 1 { // matrix index
		if f != 0 {
			x = atx(x, f)
		}
		return ecl(x, t, 64)
	}
	return atx(atx(x, f), t)
}
func atd(x, y, yt i) (r i) {
	k := I(x + 8)
	v := I(x + 12)
	if yt == 5 {
		rx(k)
		y = fnd(k, y)
		yt = 2
	}
	rx(v)
	dx(x)
	return atx(v, y)
}
func atx(x, y i) (r i) {
	xt, yt, xn, yn, xp, yp := v2(x, y)
	if xt == 0 {
		return cal(x, enl(y))
	}
	if xt == 7 { // d@..
		return atd(x, y, yt)
	}
	if yt > 5 { // at-list at-dict
		return ecr(x, y, '@') //64
	}
	if yt != 2 {
		fmt.Printf("atx x=%s y=%s\n", kst(x), kst(y))
		panic("atx yt~I")
	}
	r = mk(xt, yn)
	rp := r + 8
	w := i(C(xt))
	for i := i(0); i < yn; i++ {
		yi := I(yp)
		if xn <= yi {
			trap()
		}
		mv(rp, xp+w*yi, w)
		rp += w
		yp += 4
	}
	if xt > 4 {
		rl(r)
	}
	if xt == 6 && yn == 1 {
		rx(I(r + 8))
		dx(r)
		r = I(r + 8)
	}
	return dxyr(x, y, r)
}
func cal(x, y i) (r i) {
	xt, yt, xn, yn, xp, yp := v2(x, y)
	if xt != 0 {
		return atm(x, y)
	}
	if yt != 6 {
		panic("type")
	}
	if yn == 1 {
		if x == '\'' || x == '/' || x == '\\' || x == 128+'\'' || x == 128+'/' || x == 128+'\\' {
			f := MT[x].(func(i) i)
			return f(fst(y))
		}
		if x < 128 {
			x += 128
		}
	}
	if x < 128 {
		if yn != 2 {
			panic("arity")
		}
		rld(y)
		f := MT[x].(func(i, i) i)
		return f(I(yp), I(yp+4))
	} else if x < 256 {
		if yn != 1 {
			panic("arity")
		}
		f := MT[x].(func(i) i)
		return f(fst(y))
	}
	if xn == 2 { // derived
		rld(x)
		a := I(xp)
		if yn == 2 {
			rld(y)
			f := MT[a].(func(i, i, i) i)
			return f(I(yp), I(yp+4), I(xp+4))
		} else if yn != 1 {
			panic("arity")
		}
		f := MT[a+128].(func(i, i) i)
		return f(fst(y), I(xp+4))
	}
	if xn == 4 { // lambda
		return lcl(x, y)
	}
	panic("nyi")
}
func lcl(x, y i) (r i) { // call lambda
	fn := I(x + 20)
	if nn(y) != fn {
		panic("arity")
	}
	a := I(x + 16)
	rx(a)
	t := I(x + 12)
	rx(t)
	an := nn(I(x + 16))
	if fn < an {
		y = lcat(y, take(enl(0), an-fn))
	}
	d := mkd(a, y)
	r = lst(lev(t, d))
	dx(x)
	dx(d)
	return r
}
func cat(x, y i) (r i) {
	xt, yt, _, _, _, _ := v2(x, y)
	if xt == 0 {
		x, xt = enl(x), 6
	}
	if xt == yt {
		return ucat(x, y)
	}
	if xt == 6 {
		return ucat(x, lx(y))
	}
	if yt == 6 {
		return ucat(lx(x), y)
	}
	fmt.Printf("cat x=%s y=%s\n", kst(x), kst(y))
	panic("nyi cat")
}
func ucat(x, y i) (r i) {
	xt, _, xn, yn, xp, yp := v2(x, y)
	if xt > 4 {
		rl(x)
		rl(y)
	}
	if xt == 7 {
		r = mkd(ucat(x+8, y+8), ucat(x+12, y+12))
		return dxyr(x, y, r)
	}
	r = mk(xt, xn+yn)
	w := i(C(xt))
	mv(r+8, xp, w*xn)
	mv(r+8+w*xn, yp, w*yn)
	return dxyr(x, y, r)
}
func cc(x, y i) (r i) { // cat char
	n := nn(x)
	if bk(1, n) < bk(1, n+1) {
		return ucat(x, mkc(y))
	}
	sC(x+8+n, byte(y))
	sI(x, I(x)+1)
	return x
}
func lcat(x, y i) (r i) { // list append
	x = use(x)
	xt, xn, xp := v1(x)
	if bk(xt, xn) < bk(xt, xn+1) {
		r = mk(xt, xn+1)
		rld(x)
		mv(r+8, xp, 4*xn)
		x, xp = r, r+8
	}
	sI(x, (xn+1)|6<<29)
	sI(xp+4*xn, y)
	return x
}
func enl(x i) (r i) { r = mk(6, 1); sI(8+r, x); return r }
func cnt(x i) (r i) { dx(x); return mki(nn(x)) }
func typ(x i) (r i) {
	xt, _, _ := v1(x)
	r = mk(2, 1)
	sI(8+r, xt)
	return dxr(x, r)
}
func wer(x i) (r i) {
	xt, xn, xp := v1(x)
	if xt == 1 {
		return prs(x)
	}
	if xt == 6 {
		return flp(x)
	}
	if xt != 2 {
		panic("type")
	}
	n := i(0)
	for i := i(0); i < xn; i++ {
		n += I(xp + 4*i)
	}
	r = mk(2, n)
	rp := 8 + r
	for i := i(0); i < xn; i++ {
		nj := I(xp)
		for j := uint32(0); j < nj; j++ {
			sI(rp, i)
			rp += 4
		}
		xp += 4
	}
	return dxr(x, r)
}
func mtc(x, y i) (r i) { // x~y
	r = mk(2, 1)
	sI(r+8, match(x, y))
	return dxyr(x, y, r)
}
func match(x, y i) (r i) { // x~y
	if x == y {
		return 1
	}
	if I(x) != I(y) {
		return 0
	}
	xt, xn, xp := v1(x)
	yp, nn := y+8, i(0)
	switch xt {
	case 0:
		return 1 // todo
	case 1:
		nn = xn
	case 2:
		nn = xn << 2
	case 3:
		nn = xn << 3
	case 4:
		nn = xn << 4
	default:
		for i := i(0); i < xn; i++ {
			if match(I(xp), I(yp)) == 0 {
				return 0
			}
			xp += 4
			yp += 4
		}
		return 1
	}
	for i := i(0); i < nn; i++ {
		if C(xp+i) != C(yp+i) {
			return 0
		}
	}
	return 1
}
func not(x i) (r i) { return eql(mki(0), x) }
func tru(x i) (r i) {
	xt, xn, xp := v1(x)
	dx(x)
	if xt != 2 {
		trap()
	}
	if xn == 0 {
		return 0
	}
	return I(xp)
}
func cmp(x, y, eq i) (r i) {
	x, y = upxy(x, y)
	x, y = ext(x, y)
	if tp(x) == 6 {
		return ecd(x, y, 62-eq)
	}
	t, _, n, _, xp, yp := v2(x, y)
	cm := MT[t].(func(i, i) i)
	if eq == 1 {
		cm = MT[t+8].(func(i, i) i)
	}
	w := uint32(C(t))
	r = mk(2, n)
	rp := r + 8
	for i := i(0); i < n; i++ {
		sI(rp, cm(xp, yp))
		xp += w
		yp += w
		rp += 4
	}
	return dxyr(x, y, r)
}
func eql(x, y i) (r i) { return cmp(x, y, 1) }
func mor(x, y i) (r i) { return cmp(x, y, 0) }
func les(x, y i) (r i) { return cmp(y, x, 0) }
func fnd(x, y i) (r i) { // x?y
	xt, yt, _, yn, _, yp := v2(x, y)
	if xt != yt {
		trap()
	}
	r = mk(2, yn)
	rp := r + 8
	w := i(C(yt))
	for i := i(0); i < yn; i++ {
		sI(rp, fnx(x, yp))
		rp += 4
		yp += w
	}
	return dxyr(x, y, r)
}
func fnx(x, yp i) (r i) {
	xt, xn, xp := v1(x)
	eq := MT[8+xt].(func(i, i) i)
	w := uint32(C(xt))
	for i := i(0); i < xn; i++ {
		if eq(xp, yp) == 1 {
			return i
		}
		xp += w
	}
	return xn
}
func fnc(xp, xn i, c c) (r i) {
	for r = 0; r < xn; r++ {
		if c == C(xp+r) {
			return r
		}
	}
	return r
}
func jon(x, y i) (r i) { // y/:x (join)
	xt, xn, xp := v1(x)
	if xt != 6 || xn == 0 {
		trap()
	}
	rl(x)
	r = I(xp)
	rxn(y, xn-2)
	for i := i(0); i < xn-1; i++ {
		xp += 4
		r = cat(cat(r, y), I(xp))
	}
	return dxr(x, r)
}
func spl(x, y i) (r i) { // y\:x (split)
	rx(x)
	yn := nn(y)
	r = cut(cat(mki(0), fds(x, y)), x)
	rn := nn(r) - 1
	rp := r + 8
	for i := i(0); i < rn; i++ {
		rp += 4
		sI(rp, drop(I(rp), yn))
	}
	return r
}
func fds(x, y i) (r i) { // find subarray y in x
	xt, yt, xn, yn, xp, yp := v2(x, y)
	if xt != yt || xt > 5 {
		trap()
	}
	if yn == 0 || xn < yn {
		return dxyr(x, y, mk(2, 0))
	}
	r = mk(2, 0)
	eq := MT[8+xt].(func(i, i) i)
	w := i(C(xt))
	for i := i(0); i < xn-yn; i++ {
		a := uint32(0)
		for j := uint32(0); j < yn; j++ {
			k := w * j
			a += eq(xp+k, yp+k)
		}
		if a == yn {
			r = ucat(r, mki(i))
			i += yn - 1
			xp += w * (yn - 1)
		}
		xp += w
	}
	return dxyr(x, y, r)
}
func exc(x, y i) (r i) { rx(x); return atx(x, wer(eql(mki(nn(y)), fnd(y, x)))) } // x^y
func grd(x i) (r i) { // <x
	xt, xn, xp := v1(x)
	r = seq(0, xn, 1)
	y := seq(0, xn, 1)
	msrt(y+8, r+8, 0, xn, xp, xt) // xt:1,2,3,4,5
	return dxyr(x, y, r)
}
func gdn(x i) (r i) { return rev(grd(x)) }           // >x
func srt(x i) (r i) { rx(x); return atx(x, grd(x)) } // ^x
func msrt(x, y, z, x3, x4, x5 i) { // merge sort
	if x3-z < 2 {
		return
	}
	c := (z + x3) / 2
	msrt(y, x, z, c, x4, x5)
	msrt(y, x, c, x3, x4, x5)
	mrge(x, y, z, x3, c, x4, x5)
}
func mrge(x, y, z, x3, x4, x5, x6 i) {
	k, j, a := z, x4, i(0)
	gt := MT[x6].(func(i, i) i)
	w := uint32(C(x6))
	for i := z; i < x3; i++ {
		if k >= x4 || (j < x3 && gt(x5+w*I(x+k<<2), x5+w*I(x+j<<2)) == 1) {
			a = j
			j++
		} else {
			a = k
			k++
		}
		sI(y+i<<2, I(x+a<<2))
	}
}
func uqg(x, y i) (r i) { // ?x =x uniq/group
	xt, xn, xp := v1(x)
	r = mk(xt, 0)
	n := i(0)
	w := i(C(xt))
	for i := i(0); i < xn; i++ {
		m := fnx(r, xp)
		if m == n {
			rx(x)
			r = cat(r, atx(x, mki(i)))
			if y != 0 {
				y = lcat(y, mk(2, 0))
			}
			n += 1
		}
		if y != 0 {
			yi := y + 8 + 4*m
			sI(yi, cat(I(yi), mki(i)))
		}
		xp += w
	}
	if y != 0 {
		r = l2(r, y)
	}
	return dxr(x, r)
}
func unq(x i) (r i) { return uqg(x, 0) }        // ?x (uniq)
func grp(x i) (r i) { return uqg(x, mk(6, 0)) } // =x (group)
func flr(x i) (r i) { panic("nyi") }
func flp(x i) (r i) { panic("nyi") }

func str(x i) (r i) {
	xt, xn, xp := v1(x)
	if xt == 1 {
		return x
	}
	if xt > 5 || xn != 1 {
		if xt != 0 {
			return ech(x, 164)
		}
	}
	switch xt {
	case 0:
		r = cg(x, xn)
	case 2:
		n := I(xp)
		r = ci(n, 0)
	case 3:
		r = cf(F(xp))
	case 4:
		r = cf(F(xp))
		r = cc(r, 'i') //105
		r = ucat(r, cf(F(xp+8)))
	case 5:
		r = I(xp)
		rx(r)
	default:
		panic("nyi")
	}
	return dxr(x, r)
}
func cg(x i, xn i) (r i) {
	if x < 127 {
		r = mkc(x)
	} else if x < 256 {
		r = cc(mkc(x-128), ':') //58
	} else if xn == 4 {
		r = I(x + 8)
		rx(r)
	} else {
		panic("nyi str f")
	}
	return r
}
func ng(x, y i) (r i) {
	if y != 0 {
		x = ucat(mkc('-'), x) //45
	}
	return x
}
func ci(n, t i) (r i) {
	if n == 0 {
		return mkc('0') //48
	}
	m := i(0)
	if int32(n) < 0 {
		n = uint32(-int32(n))
		m = 1
	}
	r = mk(1, 0)
	for n != 0 {
		c := n % 10
		if c != 0 {
			t = 0
		}
		if t != 1 {
			r = cc(r, i('0'+c))
		}
		n /= 10
	}
	if nn(r) == 0 {
		r = cc(r, '0')
	}
	return ng(rev(r), m)
}
func cf(f float64) (r i) {
	if f != f {
		return cc(mkc('0'), 'n') //48 110
	}
	m := i(0)
	if f < 0 {
		m = 1
		f = -f
	}
	if f > 1.7976931348623157e+308 {
		return ng(cc(mkc('0'), 'w'), m) //119
	}
	e := i(0)
	for f > 1000 {
		e += 3
		f /= 1000.0
	}
	for f < 1 {
		e -= 3
		f *= 1000.0
	}
	fmt.Println("f/e", f, e)
	n := int32(f)
	r = ci(uint32(n), 0)
	f -= float64(n)
	d := 6 - nn(r)
	if int32(d) < 1 {
		d = 1
	}
	for i := i(0); i < d; i++ {
		f *= 10
	}
	r = ucat(cc(r, '.'), ci(uint32(f), 1))
	if e != 0 {
		r = ucat(cc(r, 'e'), ci(e, 0))
	}
	return ng(r, m)
}

func cst(x, y i) (r i) { panic("nyi") }
func sc(x i) (r i) {
	r = enl(x)
	sI(r, 1|5<<29)
	return r
}
func cs(x i) (r i) {
	r = x + 8
	rx(r)
	return dxr(x, r)
}
func mia(x, y, z i) (r i) { // minmax
	x, y = upxy(x, y)
	x, y = ext(x, y)
	if tp(x) == 6 {
		return ecd(x, y, z)
	}
	rx(x)
	rx(y)
	var a i
	if z == 38 {
		a = les(x, y)
	} else {
		a = mor(x, y)
	}
	a = wer(a)
	rx(a)
	return asi(y, a, atx(x, a))
}
func min(x, y i) (r i) { return mia(x, y, 38) }
func max(x, y i) (r i) { return mia(x, y, 124) }
func nm(x, f, h i) (r i) { // numeric monad f:scalar index (nec..), h:original func (e.g. -: neg)
	if tp(x) > 5 {
		return ech(x, h)
	}
	r = use(x)
	t, n, rp := v1(r)
	xp := x + 8
	w := uint32(C(t))
	g := MT[f+t].(func(i, i))
	for i := i(0); i < n; i++ {
		g(xp, rp)
		xp += w
		rp += w
	}
	if t == 4 && f == 19 { // +z
		return zre(r)
	}
	return r
}
func nd(x, y, f, h i) i {
	x, y = upxy(x, y)
	x, y = ext(x, y)
	t, _, n, _, xp, yp := v2(x, y)
	if t == 6 {
		return ecd(x, y, h)
	}
	w := uint32(C(t))
	g := MT[f+t].(func(i, i, i))
	r := mk(t, n)
	rp := r + 8
	for i := i(0); i < n; i++ {
		g(xp, yp, rp)
		xp += w
		yp += w
		rp += w
	}
	return dxyr(x, y, r)
}
func gtc(x, y i) i { return boolvar(C(x) > C(y)) }
func eqc(x, y i) i { return boolvar(C(x) == C(y)) }
func gti(x, y i) i { return boolvar(int32(I(x)) > int32(I(y))) }
func eqi(x, y i) i { return boolvar(int32(I(x)) == int32(I(y))) }
func gtf(x, y i) i { return boolvar(F(x) > F(y)) }
func eqf(x, y i) i { return boolvar(F(x) == F(y)) }
func eqz(x, y i) i { return boolvar(F(x) == F(y) && F(x+8) == F(y+8)) }
func gtl(x, y i) i {
	x, y = I(x), I(y)
	xt, yt, xn, yn, xp, yp := v2(x, y)
	if xt != yt {
		return boolvar(xt > yt)
	}
	n := xn
	if yn < xn {
		n = yn
	}
	gt := MT[xt].(func(i, i) i)
	w := uint32(C(xt))
	for i := i(0); i < n; i++ {
		a, b := xp+w*i, yp+w*i
		if gt(a, b) == 1 {
			return 1
		}
		if gt(b, a) == 1 {
			return 0
		}
	}
	return boolvar(xn > yn)
}
func eqL(x, y i) i  { return match(I(x), I(y)) }
func adc(x, y, r i) { sC(r, C(x)+C(y)) }
func adi(x, y, r i) { sI(r, I(x)+I(y)) }
func adf(x, y, r i) { sF(r, F(x)+F(y)) }
func adz(x, y, r i) { sZ(r, Z(x)+Z(y)) }
func suc(x, y, r i) { sC(r, C(x)-C(y)) }
func sui(x, y, r i) { sI(r, I(x)-I(y)) }
func suf(x, y, r i) { sF(r, F(x)-F(y)) }
func suz(x, y, r i) { sZ(r, Z(x)-Z(y)) }
func muc(x, y, r i) { sC(r, C(x)*C(y)) }
func mui(x, y, r i) { sI(r, I(x)*I(y)) }
func muf(x, y, r i) { sF(r, F(x)*F(y)) }
func muz(x, y, r i) { sZ(r, Z(x)*Z(y)) }
func dic(x, y, r i) { sC(r, C(x)/C(y)) }
func dii(x, y, r i) { sI(r, I(x)/I(y)) }
func dif(x, y, r i) { sF(r, F(x)/F(y)) }
func diz(x, y, r i) { sZ(r, Z(x)/Z(y)) }
func add(x, y i) i  { return nd(x, y, 15+128, 43) }
func sub(x, y i) i  { return nd(x, y, 19+128, 45) }
func mul(x, y i) i  { return nd(x, y, 23+128, 42) }
func diw(x, y i) i  { return nd(x, y, 27+128, 37) }
func abs(x i) i     { return nm(x, 15, 171) }
func neg(x i) i     { return nm(x, 19, 173) }
func sqr(x i) i     { return nm(x, 27, 165) }
func abc(x, r i) { // +c (toupper)
	if c := C(x); is(c, az) {
		sC(r, c-32)
	} else {
		sC(r, c)
	}
}
func abi(x, r i) {
	if c := int32(I(x)); c < 0 {
		sI(r, i(-c))
	} else {
		sI(r, i(c))
	}
}
func abf(x, r i) { sF(r, math.Abs(F(x))) }
func abz(x, r i) { sF(r, cmplx.Abs(Z(x))) }
func nec(x, r i) { // -c (tolower)
	if c := C(x); is(c, AZ) {
		sC(r, c+32)
	} else {
		sC(r, c)
	}
}
func nei(x, r i) { sI(r, i(-int32(I(x)))) }
func nef(x, r i) { sF(r, -F(x)) }
func nez(x, r i) { sZ(r, -Z(x)) }
func sqc(x, r i) { panic("%c") } // %c ?
func sqi(x, r i) { panic("%i") } // %i ?
func sqf(x, r i) { sF(r, math.Sqrt(F(x))) }
func sqz(x, r i) { sZ(r, cmplx.Conj(Z(x))) } // %z complex conjugate
func zri(x i, o i) (r i) {
	t, n, xp := v1(x)
	if t != 4 {
		panic("type")
	}
	r = mk(3, n)
	rp := r + 8
	xp += o
	for i := i(0); i < n; i++ {
		sF(rp, F(xp))
		rp += 8
		xp += 16
	}
	return dxr(x, r)
}
func zre(x i) (r i) { return zri(x, 0) }
func zim(x i) (r i) { return zri(x, 8) }

func drv(x, y i) (r i) { // x(adv) y(verb), e.g. ech +
	r = mk(0, 2)
	sI(8+r, x)
	sI(12+r, y)
	return r
}
func ecv(x i) (r i) { return drv(40, x) }  // '  ech(168) ecd(40)
func epv(x i) (r i) { return drv(41, x) }  // ': ecp(169) epi(41)
func ovv(x i) (r i) { return drv(123, x) } // /  ovr(251) ecr(123)
func riv(x i) (r i) { return drv(125, x) } // /: jon(253) ?(125)
func scv(x i) (r i) { return drv(91, x) }  // \  scn(219) ecl(91)
func liv(x i) (r i) { return drv(93, x) }  // \: spl(221) ?(93)
func ech(x, y i) (r i) { // f'x (each)
	if tp(y) != 0 {
		return bin(y, x)
	}
	if tp(x) == 7 {
		rld(x)
		k := I(x + 8)
		v := I(x + 12)
		return mkd(k, ech(v, y))
	}
	x = lx(x)
	_, xn, xp := v1(x)
	r = mk(6, xn)
	rp := r + 8
	rl(x)
	if y < 128 { // force monad
		y += 128
	}
	for i := i(0); i < xn; i++ {
		rx(y)
		sI(rp, atx(y, I(xp)))
		xp += 4
		rp += 4
	}
	return dxyr(x, y, r)
}
func ecp(x, y i) (r i) { // f':x (each-prior)
	rx(x)
	p := fst(x)
	return epi(p, x, y)
}
func epi(x, y, z i) (r i) { // x f':y (each-prior-initial)
	n := nn(y)
	if n == 0 {
		return dxyr(x, z, y)
	}
	rxn(y, n)
	rxn(z, n)
	r = mk(6, n)
	rp := r + 8
	var yi i
	for i := i(0); i < n; i++ {
		yi = atx(y, mki(i))
		rx(yi)
		sI(rp, cal(z, l2(yi, x)))
		x = yi
		rp += 4
	}
	dx(yi)
	return dxyr(y, z, r)
}
func ovr(x, y i) (r i) { return ovs(x, y, 0) }             // y/x (over/reduce)
func scn(x, y i) (r i) { return ovs(x, y, enl(mk(6, 0))) } // y\x (scan)
func ovs(x, y, z i) (r i) { // over/scan
	if ary(y) == 1 {
		return fxp(x, y, z)
	}
	n := nn(x)
	rxn(x, n)
	r = fst(x) // panics on n~0
	rxn(y, n-1)
	scl(z, r)
	for i := i(0); i < n-1; i++ {
		r = cal(y, l2(r, atx(x, mki(i+1))))
		scl(z, r)
	}
	dx(x)
	dx(y)
	if z == 0 {
		return r
	}
	dx(r)
	return fst(z)
}
func fxp(x, y, z i) (r i) { // fixed point/converge
	t := x
	rx(x)
	for {
		rx(x)
		rx(y)
		r = atx(y, x)
		if match(r, x)+match(r, t) != 0 {
			dx(x)
			dx(y)
			dx(t)
			if z != 0 {
				r = lcat(fst(z), r)
			}
			return r
		}
		scl(z, x)
		dx(x)
		x = r
	}
}
func scl(x, y i) {
	if x != 0 {
		xp := x + 8
		rx(y)
		sI(xp, lcat(I(xp), y))
	}
}
func ecr(x, y, f i) (r i) { // x f/ y (each-right)
	if ary(f) == 1 {
		return whl(x, y, f, 0)
	}
	if tp(y) == 7 {
		rld(y)
		k := I(y + 8)
		v := I(y + 12)
		return mkd(k, ecr(x, v, f))
	}
	n := nn(y)
	r = mk(6, n)
	rp := r + 8
	rxn(x, n)
	rxn(y, n)
	rxn(f, n)
	for i := i(0); i < n; i++ {
		sI(rp, cal(f, l2(x, atx(y, mki(i)))))
		rp += 4
	}
	dx(f)
	return dxyr(x, y, r)
}
func whl(x, y, f, s i) (r i) {
	xt := tp(x)
	if xt != 0 {
		if xt != 2 || nn(x) != 1 {
			trap()
		}
		dx(x)
		return nlp(y, f, s, I(x+8))
	}
	r = y
	scl(s, r)
	n := mki(0)
	for {
		rx(x)
		rx(f)
		r = atx(f, r)
		scl(s, r)
		rx(r)
		t := atx(x, r)
		if 1 == match(t, n) {
			dx(t)
			dx(n)
			dx(f)
			dx(x)
			if s != 0 {
				dx(r)
				r = fst(s)
			}
			return r
		}
		dx(t)
	}
}
func nlp(x, f, s, n i) (r i) { // n f/y (for)  n f\y (scan-for)
	if int32(n) < 0 {
		trap()
	}
	r = x
	rxn(f, n)
	scl(s, x)
	for i := i(0); i < n; i++ {
		r = atx(f, r)
		scl(s, r)
	}
	dx(f)
	if s != 0 {
		dx(r)
		r = fst(s)
	}
	return r
}
func ecl(x, y, f i) (r i) { // x f\ y (each-left)
	if ary(f) == 1 {
		return whl(x, y, f, enl(mk(6, 0)))
	}
	if tp(x) == 7 {
		rld(x)
		k := I(x + 8)
		v := I(x + 12)
		return mkd(k, ecl(v, y, f))
	}
	n := nn(x)
	r = mk(6, n)
	rp := r + 8
	rxn(x, n)
	rxn(y, n)
	rxn(f, n)
	for i := i(0); i < n; i++ {
		sI(rp, cal(f, l2(atx(x, mki(i)), y)))
		rp += 4
	}
	dx(f)
	return dxyr(x, y, r)
}
func ecd(x, y, f i) (r i) { // x f' y
	n := nn(x)
	if n != nn(y) {
		panic("ecd length")
	}
	r = mk(6, n)
	rp := r + 8
	rxn(x, n)
	rxn(y, n)
	rxn(f, n)
	for i := i(0); i < n; i++ {
		c := mki(i)
		rx(c)
		sI(rp, cal(f, l2(atx(x, c), atx(y, c))))
		rp += 4
	}
	dx(f)
	return dxyr(x, y, r)
}
func bin(x, y i) (r i) { // x'y binary search
	xt, yt, xn, yn, xp, yp := v2(x, y)
	if xt != yt {
		trap()
	}
	r = mk(2, yn)
	rp := r + 8
	w := i(C(xt))
	for i := i(0); i < yn; i++ {
		sI(rp, ibin(xp, yp, xn, xt))
		rp += 4
		yp += w
	}
	return dxyr(x, y, r)
}
func ibin(x, y, n, t i) (r i) {
	k, j, h := i(0), n-1, i(0)
	gt := MT[t].(func(i, i) i)
	w := i(C(t))
	for {
		if int32(k) > int32(j) {
			return k - 1
		}
		h = (k + j) >> 1
		if gt(x+w*h, y) != 0 {
			j = h - 1
		} else {
			k = h + 1
		}
	}
}

func val(x i) (r i) {
	xt, _, _ := v1(x)
	switch xt {
	case 1:
		return evl(prs(x), 0)
	case 6:
		return evl(x, 0)
	case 7:
		rx(x + 12)
		return dxr(x, x+12)
	default:
		fmt.Printf("val xt=%d\n", xt)
		panic("nyi")
	}
}
func kv(x, loc i) (k, v, n, m i) {
	k = I(kkey)
	v = I(kval)
	if loc != 0 {
		k = I(loc + 8)
		v = I(loc + 12)
	}
	n = I(k) & 536870911
	m = fnx(k, 8+x)
	return k, v, n, m
}
func lup(x, loc i) (r i) {
	_, v, n, m := kv(x, loc)
	if m == n {
		if loc != 0 {
			return lup(x, 0)
		}
		return dxr(x, 0)
	}
	r = I(v + 8 + 4*m)
	rx(r)
	return dxr(x, r)
}
func asn(x, loc, u i) (r i) {
	xt, _, _ := v1(x)
	if xt != 5 {
		trap()
	}
	rx(u)
	k, v, n, m := kv(x, loc)
	if n == m {
		if loc != 0 {
			panic("asn to undefined local")
		}
		sI(kkey, cat(k, x))
		sI(kval, lcat(v, u))
		return u
	}
	vp := 8 + v + 4*m
	dx(I(vp))
	sI(vp, u)
	return dxr(x, u)
}
func asi(x, y, z i) (r i) { //x[..y..]:z
	xt, yt, xn, yn, xp, yp := v2(x, y)
	if xt == 7 && yt == 5 {
		rld(x)
		k := I(xp)
		v := I(xp + 4)
		rx(k)
		y = fnd(k, y)
		return mkd(k, asi(v, y, z))
	}
	if yt == 6 {
		if xt == 7 { // (:;d;y;z)  {k:!x;v:. x;y:(,k?*y),1_y;k!.(:;v;y;z)}
			rld(x)
			k := I(xp)
			v := I(xp + 4)
			rx(y)
			f := fst(y)
			if f != 0 {
				rx(k)
				f = fnd(k, f)
			} else {
				f = seq(0, nn(k), 1)
			}
			return mkd(k, asi(v, cat(enl(f), drop(y, 1)), z))
		}
		if xt != 6 || yt != 6 {
			trap()
		}
		r = take(x, xn)
		rp := r + 8
		rx(y)
		a := fst(y)
		y = drop(y, 1)
		if nn(y) == 1 {
			y = fst(y)
		}
		if a == 0 {
			a = seq(0, xn, 1)
		}
		at, an, ap := v1(a)
		if at != 2 {
			panic("index-type")
		}
		if an == 1 { // depth-assign
			dx(a)
			ri := rp + 4*I(ap)
			sI(ri, asi(I(ri), y, z))
			return r
		}
		if yn != 2 {
			panic("matrix-assign")
		} // matrix-assign
		zt := tp(z)
		if zt != 6 {
			z = take(enl(z), an)
		}
		if nn(z) != an {
			trap()
		}
		rxn(y, an-1)
		rl(z)
		zp := z + 8
		for i := i(0); i < an; i++ {
			ri := rp + 4*I(ap)
			sI(ri, asi(I(ri), y, I(zp)))
			ap += 4
			zp += 4
		}
		dx(a)
		dx(z)
		return r
	}
	zt, zn, zp := v1(z)
	if yn > 1 && zn == 1 {
		if zn != yn {
			if zn != 1 {
				trap()
			}
			z = take(z, yn)
			zn = yn
			zp = z + 8
		}
	}
	if xt < 5 {
		if zt != xt {
			trap()
		}
		r = mk(xt, xn)
		rp := r + 8
		w := i(C(xt))
		mv(rp, xp, w*xn)
		for i := i(0); i < yn; i++ {
			k := I(yp)
			mv(rp+w*k, zp, w)
			yp += 4
			zp += w
		}
		dx(z)
		return dxyr(x, y, r)
	}
	if xt == 6 {
		r = take(x, xn)
		rp := r + 8
		z = lx(z) // explode
		zp = z + 8
		rl(z)
		for i := i(0); i < yn; i++ {
			k := I(yp)
			t := rp + 4*k
			dx(I(t))
			sI(t, I(zp))
			yp += 4
			zp += 4
		}
		return dxyr(y, z, r)
	}
	trap()
	return x
}
func asd(x, loc i) (r i) { // (+;`x;a;y)
	rld(x)
	v, s, a, u := I(x+8), I(x+12), I(x+16), I(x+20)
	if v != ':' { //58
		rx(s)
		r = lup(s, loc)
		if a != 0 {
			rx(a)
			r = atx(r, a)
		}
		u = cal(v, l2(r, u))
	}
	if a != 0 {
		rx(s)
		u = asi(lup(s, loc), a, u)
	}
	rx(s)
	r = asn(s, loc, u)
	return dxr(s, r)
}
func swc(x, loc i) (r i) { // ($;a;b;...)
	_, xn, xp := v1(x)
	for i := i(1); i < xn; {
		r = I(xp + 4*i)
		rx(r)
		r = evl(r, loc)
		if i%2 == 0 || i == xn-1 {
			return dxr(x, r)
		}
		dx(r)
		i++
		if I(r+8) == 0 {
			i++
		}
	}
	return dxr(x, 0)
}
func lev(x, loc i) (r i) {
	xt, xn, xp := v1(x)
	if xt != 6 {
		return x
	}
	rl(x)
	r = mk(6, xn)
	rp := r + 8
	for i := i(0); i < xn; i++ {
		sI(rp, evl(I(xp), loc))
		rp += 4
		xp += 4
	}
	return dxr(x, r)
}
func ras(x, xn, loc i) (r i) { // rewrite assignments x[i]+:y  (+:;(`x;i);y)→(+;,`x;,i;y)  (and collect locals)
	// defer func() { fmt.Printf("ras r=%s\n", kst(r)) }()
	v := I(x + 8)
	if xn == 3 && v < 256 && (v == ':' || v > 128) { //58
		if v > 128 {
			v -= 128
		}
		r = I(x + 12)
		rxn(r, 2)
		s := fst(r)
		a := drop(r, 1)
		if nn(a) == 0 {
			dx(a)
			a = 0
		} else {
			a = lev(a, loc)
			an := nn(a)
			if an == 1 {
				a = fst(a)
			}
		}
		u := I(x + 16)
		rx(u)
		dx(x)
		return lcat(l3(v, s, a), evl(u, loc))
	}
	return 0
}
func evl(x, loc i) (r i) {
	xt, xn, xp := v1(x)
	if xt != 6 {
		if xt == 5 && xn == 1 {
			r = lup(x, loc)
			if r == 0 {
				panic("name does not exist")
			}
			return r
		}
		return x
	} else if xn == 0 {
		return x
	} else if xn == 1 {
		return lev(fst(x), loc)
	}
	v := I(xp)
	if v == '$' && xn > 3 { // 36 ($;a;b;..) switch $[a;b;..]
		return swc(x, loc)
	}
	r = ras(x, xn, loc)
	if r != 0 {
		return asd(r, loc)
	}
	x = lev(x, loc)
	xn = nn(x)
	xp = x + 8
	if v == 128 {
		if xn > 2 { // 128 (,:;a;b;c) sequence
			return lst(x)
		}
	}
	if xn == 2 {
		rl(x)
		r = atx(I(xp), I(xp+4))
		return dxr(x, r)
	}
	rx(I(xp))
	return cal(I(xp), drop(x, 1))
}
func prs(x i) (r i) { // parse (k.w) E:E;e|e e:nve|te| t:n|v|{E} v:tA|V n:t[E]|(E)|N
	xt, xn, xp := v1(x)
	if xt != 1 {
		trap()
	}
	sI(pp, xp)
	r = sq(xp + xn)
	if nn(r) == 1 {
		r = fst(r)
	} else {
		r = cat(128, r) // :::
	}
	return dxr(x, r)
}
func sq(s i) (r i) { // E
	r = enl(ex(pt(s), s))
	for {
		v := ws(s)
		p := I(pp)
		if !v {
			v = C(p) != 59 // ;
		}
		if v {
			if p < s {
				sI(pp, p+1)
			}
			return r
		}
		sI(pp, p+1)
		r = lcat(r, ex(pt(s), s))
	}
}
func ex(x, s i) (r i) { // e
	//defer func() { fmt.Printf("ex %d %q\n", r, kst(r)) }()
	if x == 0 || ws(s) || is(C(I(pp)), TE) { //32
		return x
	}
	r = pt(s) // nil?
	if isv(r) && !isv(x) {
		return l3(r, x, ex(pt(s), s))
	}
	return l2(x, ex(r, s))
}
func pt(s i) (r i) { // t
	//defer func() { fmt.Printf("pt r=%d %q\n", r, kst(r)) }()
	r = tok(s)
	if r == 0 {
		p := I(pp)
		if p == s {
			return 0
		}
		λ := C(p) == 123     // {
		if λ || C(p) == 40 { // (
			sI(pp, p+1)
			r = sq(s)
			if λ {
				r = lam(p, I(pp), r)
			} else {
				if n := nn(r); n == 1 {
					r = fst(r)
				} else if n > 1 {
					r = enl(r)
				}
			}
		}
	}
	for {
		p := I(pp)
		b := C(p)
		if p == s {
			return r
		}
		if is(b, AD) { // 16
			r = l2(tok(s), r) // adverb
		} else if b == '[' {
			sI(pp, p+1)
			r = cat(enl(r), sq(s))
		} else {
			return r
		}
	}
}
func isv(x i) (r bool) { // is verb or (adverb;_)
	xt, xn, xp := v1(x)
	if xt == 0 {
		return true // function
	}
	if xt == 6 && xn == 2 {
		a := I(xp)
		if a < 256 {
			if is(c(a), AD) || is(c(a-128), AD) { // AD(16)
				return true // adverb
			}
		}
	}
	return false
}
func lac(x, a i) (r i) { // lambda arity from tree {x+z}->3
	xt, xn, xp := v1(x)
	if xt == 6 {
		for i := i(0); i < xn; i++ {
			a = lac(I(xp), a)
			xp += 4
		}
	}
	if xt == 5 && xn == 1 {
		p := I(xp)
		if nn(p) == 1 {
			r = i(C(8+p) - 'w') //119
			if r > a {
				if r < 4 {
					return r
				}
			}
		}
	}
	return a
}
func loc(x, y i) (r i) {
	xt, xn, xp := v1(x)
	if xt != 6 {
		return y
	}
	for i := i(0); i < xn; i++ {
		y = loc(I(xp), y)
		xp += 4
	}
	xp = x + 8
	if xn == 3 && I(xp) == 58 { // :
		r = I(xp + 4)
		rx(r)
		s := fst(r)
		n := nn(y)
		if fnx(y, s+8) == n {
			rx(s)
			y = cat(y, s)
		}
		dx(s)
	}
	return y
}
func lam(p, s, z i) (r i) {
	var a i
	if C(1+p) == '[' { //91 {[a;b]a..} -> ((;`a;`b);(..))
		rx(z)
		a = ovr(drop(fst(z), 1), 44) // ,/1_*z
		z = drop(z, 1)
	} else {
		r = I(xyz)
		rx(r)
		a = take(r, lac(z, 0))
	}
	v := nn(a) // arity (<256)
	a = loc(z, a)
	n := s - p
	t := mk(1, n)
	mv(t+8, p, n)
	r = mk(0, 4)
	sI(r+8, t)  // string
	sI(r+12, z) // tree
	sI(r+16, a) // args
	sI(r+20, v) // arity
	return r
}
func ws(s i) bool { // skip whitespace
	p := I(pp)
	for {
		if p == s {
			sI(pp, p)
			return true // EOF
		}
		b := C(p)
		if is(b, NW) || b == 10 { //64
			sI(pp, p)
			return false
		}
		p++
	}
}
func tok(s i) (r i) { // next token
	if ws(s) {
		return 0
	}
	p := I(pp)
	b := C(p)
	if is(b, TE) { //32
		return 0
	}
	n := 0
	if is(C(p-1), NM) { //4 (0-1: parse - as verb, skip nms); C(p-1) maybe refcount's last byte
		n = 1
	}
	for j := 0; j < 5-n; j++ { // nms vrb chr nam sms
		r = MT[j+136+n].(func(c, i, i) i)(b, p, s)
		if r != 0 {
			return r
		}
	}
	return 0
}
func pun(b c, p, s i) (r i) { // parse unsigned int
	if !is(b, NM) { // 4
		return 0
	}
	for is(b, NM) && p < s {
		r *= 10
		r += uint32(b) - '0' //48
		p++
		b = C(p)
	}
	sI(pp, p)
	return mki(r)
}
func pin(b c, p, s i) (r i) { // parse signed int
	r = pun(b, p, s)
	if r != 0 {
		return r
	}
	if b == '-' { //45
		p++
		if p < s {
			b = C(p)
			sI(pp, p)
			r = pun(b, p, s)
			if r == 0 {
				sI(pp, p-1)
				return 0
			}
			sI(8+r, -I(8+r))
			return r
		}
	}
	return 0
}
func pfl(b c, p, s i) (r i) { // parse float (-)(u32).(u32) parts may overflow, no exp
	r = pin(b, p, s)
	if r == 0 {
		return r
	}
	p = I(pp)
	if C(p) == '.' { //46
		r = up(r, 2, 1)
		p++
		sI(pp, p)
		if p < s {
			b = C(p)
			q := pun(b, p, s)
			if q != 0 {
				q = up(q, 2, 1)
				f := 1.0
				n := I(pp) - p
				for i := i(0); i < n; i++ {
					f *= 10
				}
				rp := 8 + r
				sF(rp, F(rp)+F(8+q)/f)
				dx(q)
			}
		}
	}
	return r
}
func num(b c, p, s i) (r i) { // parse single number
	return pfl(b, p, s)
}
func nms(b c, p, s i) (r i) { // parse numeric vector
	r = num(b, p, s)
	if r == 0 {
		return r
	}
	for {
		p = I(pp)
		b = C(p)
		if p+2 > s {
			return r
		}
		if b != 32 {
			return r
		}
		p++
		sI(pp, p)
		q := num(C(p), p, s)
		if q == 0 {
			return r
		}
		r = upx(r, q)
		q = upx(q, r)
		r = cat(r, q)
	}
}
func vrb(b c, p, s i) (r i) { // verb or adverb + -: ':
	if !is(b, VB|AD) { // 24
		return 0
	}
	if b == 39 { // (space)'c  spacy verb
		if C(p-1) == 32 {
			p++
		}
	}
	r = i(C(p))
	if s > p+1 {
		if C(p+1) == ':' { // 58
			p++
			r += 128
		}
	}
	sI(pp, p+1)
	return r
}
func nam(b c, p, s i) (r i) { // abc  A3 (as `abc)
	if !is(b, az|AZ) { //3
		return 0
	}
	a := p
	for {
		p++
		if p == s || !is(C(p), az|AZ|NM) { //7
			n := p - a
			r = mk(1, n)
			mv(r+8, a, n)
			sI(pp, p)
			return sc(r)
		}
	}
}
func sym(b c, p, s i) (r i) { // `abc `"abc"
	if b != '`' { //96
		return 0
	}
	p++
	b = C(p)
	sI(pp, p)
	if p < s {
		if r = nam(b, p, s); r != 0 {
			return r
		}
		if r = chr(b, p, s); r != 0 {
			return sc(r)
		}
	}
	r = mk(5, 1)
	sI(r+8, mk(1, 0))
	return r
}
func sms(b c, p, s i) (r i) { // `a`b as ,`a`b
	r = sym(b, p, s)
	if r == 0 {
		return r
	}
	for {
		p = I(pp)
		q := sym(C(p), p, s)
		if q == 0 {
			return enl(r)
		}
		r = cat(r, q)
	}
	return r
}
func chr(b c, p, s i) (r i) { // "abc"
	if b != '"' { //34
		return 0
	} // todo hex
	a := p + 1
	for {
		p++
		if p == s {
			panic("chr/eof")
		}
		if C(p) == '"' { // todo quote
			n := p - a
			r = mk(1, n)
			mv(r+8, a, n)
			sI(pp, p+1)
			return r
		}
	}
}

const az, AZ, NM, VB, AD, TE, NW = 1, 2, 4, 8, 16, 32, 64 // see p.go
func is(x, m c) bool                                      { return (m & cla(x)) != 0 }
func cla(b c) c {
	if 128 < (b - 32) {
		return 0
	}
	return C(i(128 + b))
}
func cmake() { // init character token map (generated by go run p.go -c)
	n := 128 - 32
	s := "204840484848485040604848484848504444444444444444444448604848484848424242424242424242424242424242424242424242424242424240506048484041414141414141414141414141414141414141414141414141414048604800"
	if len(s) != 2*n {
		panic("cmap")
	}
	for i := 0; i < n; i++ {
		n, _ := strconv.ParseUint(s[2*i:2+2*i], 16, 8)
		MC[cmap+i] = c(n)
	}
}

func boolvar(b bool) i {
	if b {
		return 1
	}
	return 0
}
func trap() { panic("trap") }
func dump(a, n i) i { // type: cifzsld -> 2468ace
	p := a >> 2
	fmt.Printf("%.8x ", a)
	for i := i(0); i < n; i++ {
		x := MI[p+i]
		fmt.Printf(" %.8x", x)
		if i > 0 && (i+1)%8 == 0 {
			fmt.Printf("\n%.8x ", a+4*i+4)
		} else if i > 0 && (i+1)%4 == 0 {
			fmt.Printf(" ")
		}
	}
	fmt.Println()
	return 0
}
func C(a i) c              { return MC[a] } // global get, e.g. I i
func I(a i) i              { return MI[a>>2] }
func J(a i) j              { return MJ[a>>3] }
func F(a i) f              { return MF[a>>3] }
func Z(a i) complex128     { return complex(MF[a>>3], MF[1+a>>3]) }
func sC(a i, v c)          { MC[a] = v } // global set, e.g. i::v
func sI(a i, v i)          { MI[a>>2] = v }
func sJ(a i, v j)          { MJ[a>>3] = v }
func sF(a i, v f)          { MF[a>>3] = v }
func sZ(a i, v complex128) { MF[a>>3] = real(v); MF[1+a>>3] = imag(v) }
func atoi(s string) i {
	if x, e := strconv.Atoi(s); e == nil {
		return i(x)
	}
	panic("atoi")
}
func mark() { // mark bucket type within free blocks
	//dump(0, 200)
	for t := i(4); t < 32; t++ {
		p := MI[t]
		for p != 0 {
			MI[1+(p>>2)] = 0
			MI[2+(p>>2)] = t
			p = MI[p>>2]
		}
	}
}
func leak() {
	//dump(0, 200)
	dx(I(kkey))
	dx(I(kval))
	dx(I(xyz))
	//dump(0, 300)
	mark()
	p := i(64)
	for p < i(len(MI)) {
		if MI[p+1] != 0 {
			panic(fmt.Errorf("non-free block at %d(%x)", p<<2, p<<2))
		}
		t := MI[p+2]
		if t < 4 || t > 31 {
			panic(fmt.Errorf("illegal bucket type %d at %d(%x)", t, p<<2, p<<2))
		}
		dp := i(1) << t
		p += dp >> 2
	}
}
func kst(x i) s {
	if x == 0 || x == 128 {
		return ""
	}
	t, n, _ := v1(x)
	var f func(i i) s
	var tof func(s) s = func(s s) s { return s }
	istr := func(i i) s {
		if n := int32(MI[i+2+x>>2]); n == -2147483648 {
			return "0N"
		} else {
			return strconv.Itoa(int(n))
		}
	}
	fstr := func(i i) s {
		if f := MF[i+1+x>>3]; math.IsNaN(f) {
			return "0n"
		} else {
			return strconv.FormatFloat(f, 'g', -1, 64)
		}
	}
	zstr := func(i i) s {
		if z := Z(x + 8 + 16*i); cmplx.IsNaN(z) {
			return "0ni0n"
		} else {
			return strconv.FormatFloat(real(z), 'g', -1, 64) + "i" + strconv.FormatFloat(imag(z), 'g', -1, 64)
		}
	}
	sstr := func(i i) s {
		r := I(x + 8 + 4*i)
		rn := I(r) & 536870911
		return string(MC[r+8 : r+8+rn])
	}
	sep := " "
	switch t {
	case 0:
		fc := []byte(":+-*%&|<>=!~,^#_$?@.'/\\({[)})")
		if x < 128 && bytes.Index(fc, []byte{byte(x)}) != -1 {
			return string(byte(x))
		} else if x < 256 && bytes.Index(fc, []byte{byte(x - 128)}) != -1 {
			return string(byte(x-128)) + ":"
		} else if n == 2 { // todo: ({[)}) => ' / \ ': /: \:
			return kst(MI[(x+12)>>2]) + kst(MI[(x+8)>>2])
		} else if n == 4 { // lambda
			f, n = sstr, 1
		} else {
			return fmt.Sprintf(" '(%d)", x)
		}
	case 1:
		return `"` + string(MC[x+8:x+8+n]) + `"`
	case 2:
		f = istr
	case 3:
		f = fstr
		tof = func(s s) s {
			if strings.Index(s, ".") == -1 {
				return s + ".0"
			}
			return s
		}
	case 4:
		f = zstr
	case 5:
		if n == 0 {
			return "0#`"
		}
		f = sstr
		sep = "`"
		tof = func(s s) s { return "`" + s }
	case 6:
		if n == 1 {
			return "," + kst(I(8+x))
		}
		f = func(i i) s { return kst(MI[2+i+x>>2]) }
		sep = ";"
		tof = func(s s) s { return "(" + s + ")" }
	case 7:
		return kst(I(x+8)) + "!" + kst(I(x+12))
	default:
		panic(fmt.Sprintf("nyi: kst: t=%d", t))
	}
	r := make([]s, n)
	for k := range r {
		r[k] = f(i(k))
	}
	return tof(strings.Join(r, sep))
}
