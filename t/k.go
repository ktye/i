// k bootstrap interpreter
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/bits"
	"math/cmplx"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"unsafe"
)

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

const pp, kkey, kval, xyz, kcon, cmap = 8, 132, 136, 148, 156, 160

var version = "k.go"
var argvParsers []func(string, []string) ([]string, bool)
var replParsers []func(string) bool
var kiniRunners []func()
var exit = os.Exit
var Out func(x uint32)

func init() {
	// -t(test)
	t := func(a string, tail []string) ([]string, bool) {
		if a != "-t" {
			return tail, false
		}
		multitest()
		runtest("t")
		os.Exit(0)
		return tail, true
	}
	// -f file  -fvar file  (load bytes + assign)
	f := func(a string, tail []string) ([]string, bool) {
		if strings.HasPrefix(a, "-f") == false {
			return tail, false
		}
		b, e := ioutil.ReadFile(tail[0])
		fatal(e)
		name := a[2:]
		if name == "" {
			name = tail[0]
		}
		dx(asn(sc(kC([]byte(name))), kC(b)))
		return tail[1:], true
	}
	// file.k (load)
	k := func(a string, tail []string) ([]string, bool) {
		if strings.HasSuffix(a, ".k") == false {
			return tail, false
		}
		load(a)
		return tail, true
	}
	// -leak (test)
	lk := func(a string, tail []string) ([]string, bool) {
		if a == "-leak" {
			leak()
			fmt.Println("no leak")
		}
		return tail, false
	}
	// -e [EXPR] then exit
	e := func(a string, tail []string) ([]string, bool) {
		if a != "-e" {
			return tail, false
		}
		if len(tail) > 0 {
			s := strings.Join(tail, " ")
			dx(out(val(kC([]byte(s)))))
		}
		leak()
		interactive = true // always exit from argv -e
		exit(0)
		return tail, false
	}
	// -d (debug)
	deb := func(a string, tail []string) ([]string, bool) {
		if a == "-d" == false {
			return tail, false
		}
		ddd = true
		return tail, true
	}
	argvParsers = append(argvParsers, t, f, k, lk, e, deb)

	// \leak
	le := func(a string) bool {
		if strings.HasPrefix(a, `\leak`) == false {
			return false
		}
		b := make([]uint64, len(MJ))
		copy(b, MJ)
		leak()
		copy(MJ, b)
		msl()
		fmt.Println("no leak")
		return true
	}
	// \v(vars)
	v := func(a string) bool {
		if a != `\v` {
			return false
		}
		x := I(kkey)
		rx(x)
		fmt.Printf("kkey[%x] %d/%d\n", x, tp(x), nn(x))
		dx(out(jon(x, mkc('`'))))
		return true
	}
	// \d(toggle debug)
	d := func(a string) bool {
		if a != `\d` {
			return false
		}
		ddd = !ddd
		fmt.Println("debug", ddd)
		return true
	}
	// \s(stack)
	s := func(a string) bool {
		if a != `\s` {
			return false
		}
		os.Stdout.Write(lastStack)
		return true
	}
	// \\(exit)
	ex := func(a string) bool {
		if a != `\\` {
			return false
		}
		exit(0)
		return true
	}
	replParsers = append(replParsers, le, v, d, s, ex)
}
func main() { Main(os.Args[1:]) }
func Main(args []string) {
	kinit()
	defer indicate()
	for {
		if len(args) == 0 {
			break
		}
		no := true
		for _, p := range argvParsers {
			if tail, ok := p(args[0], args[1:]); ok {
				args = tail
				no = false
				break
			}
		}
		if no {
			panic("arg: " + args[0])
		}
	}
	repl()
}

var interactive bool

func repl() {
	interactive = true
	s := bufio.NewScanner(os.Stdin)
	fmt.Printf("%s\n ", version)
	o := Out
	for s.Scan() {
		Out = o
		t := s.Text()
		if len(t) > 0 && interactive && t[0] == ' ' {
			Out = nil
		}
		c := strings.TrimSpace(t)
		no := true
		for _, p := range replParsers {
			if p(c) == true {
				no = false
				break
			}
		}
		if no {
			ktry(t)
		}
		fmt.Printf(" ")
	}
}
func fatal(e error) {
	if e != nil {
		panic(e)
	}
}
func load(f string) {
	b, e := ioutil.ReadFile(f)
	if e != nil {
		panic(e)
	}
	if n := bytes.Index(b, []byte("\n\\")); n != -1 {
		b = b[:n+1]
	}
	dx(out(val(kC(b))))
}

var ddd = false
var lastStack []byte

func ktry(s string) {
	b := make([]uint64, len(MJ))
	copy(b, MJ)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			lastStack = debug.Stack()
			ics(I(140), I(144), os.Stdout)
			MJ = b
			msl()
		}
	}()
	dx(out(val(kC([]byte(s)))))
}
func indicate() {
	if r := recover(); r != nil {
		ics(I(140), I(144), os.Stdout)
		panic(r)
	}
}

func multitest() { //multiline, spaces, comments
	tc := []struct{ in, exp string }{
		{"1+2", "3"},
		{"1 +2", "3"},
		{"1+ 2", "3"},
		{"(1;2)", "(1;2)"},
		{"(1\n2)", "(1;2)"},
		{"(1 \n2)", "(1;2)"},
		{"(1 /xx\n2)", "(1;2)"},
		{"(1\n2 3;4 ) /comment", "(1;2 3;4)"},
		{`1+2 /alpha`, "3"},
		//{`&"1+2 /alpha"`, "(+;1;2)"},
		{"/x\n2", "2"},
	}
	for _, t := range tc {
		got := run1(t.in)
		fmt.Printf("%q /%s\n", t.in, got)
		if got != t.exp {
			fmt.Println("expected:", t.exp)
			exit(1)
		}
	}
}
func runtest(file string) {
	b, e := ioutil.ReadFile(file)
	if e != nil {
		panic(e)
	}
	v := strings.Split(string(b), "\n")
	for i := range v {
		if len(v[i]) == 0 {
			fmt.Println("skip rest")
			exit(0)
		}
		if len(v[i]) == 0 || v[i][0] == '/' {
			continue
		}
		if v[i][0] == '&' && v[i][1] == '"' {
			fmt.Printf("skip %s\n", v[i])
			continue
		}
		vv := strings.Split(v[i], " /")
		if len(vv) != 2 {
			panic("test file")
		}
		in := strings.TrimRight(vv[0], " \t\r")
		exp := strings.TrimSpace(vv[1])
		got := run1(in)
		fmt.Println(in, "/", got)
		if exp != got {
			fmt.Println("expected:", exp)
			exit(1)
		}
	}
}
func run1(s string) string {
	kinit()
	x := kC([]byte(s))
	x = kst(val(x))
	s = string(MC[x+8 : x+nn(x)+8])
	dx(x)
	leak()
	return s
}
func ics(x, y i, w io.Writer) {
	xt, xn, xp := v1(x)
	if xt != 1 || xn == 0 {
		fmt.Fprintf(w, "src(140): %d/%d\n", xt, xn)
		return
	}
	y -= 8
	p := xp
	q := xp + xn
	for i := i(0); i < xn; i++ {
		if i < y && C(xp) == 10 {
			p = xp + 1
		} else if i > y && C(xp) == 10 && q > xp {
			q = xp
		}
		xp++
	}
	y -= p - (x + 8)
	w.Write(MC[p:q])
	w.Write([]byte{10})
	if y > 0 && y < 100 {
		w.Write(bytes.Repeat([]byte{' '}, int(y)))
	}
	w.Write([]byte{'^', 10})
}
func kinit() {
	m0 := 16
	MJ = make([]j, (1<<m0)>>3)
	msl()   // pointers MC, MI, ..
	cmake() // char maps
	ini(16)
	if len(kiniRunners) > 0 {
		for i := len(kiniRunners) - 1; i >= 0; i-- {
			kiniRunners[i]()
		}
	}
}
func ini(x i) i {
	sin, cos, exp, log := math.Sin, math.Cos, math.Exp, math.Log
	copy(MT[0:], []interface{}{
		//   1    2    3    4    5    6    7    8    9    10   11   12   13   14   15
		nil, gtc, gti, gtf, nil, gti, gtl, mod, nil, eqc, eqi, eqf, eqz, eqi, eqL, nil, abc, abi, abf, abz, nec, nei, nef, nez, nil, moi, nil, nil, sqc, sqi, sqf, sqz, // 000..031
		nil, mkd, nil, rsh, cst, diw, min, ecv, ecd, epi, mul, add, cat, sub, cal, ovv, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, asn, nil, les, eql, mor, fnd, // 032..063
		atx, nil, nil, nil, nil, nil, nmf, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, sci, scv, ecl, exc, cut, // 064..095
		nil, nil, nil, nil, drw, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, ovi, max, ecr, mtc, nil, // 096..127
		nil, sin, cos, exp, log, nil, nil, nil, chr, nms, vrb, nam, sms, nil, nil, nil, adc, adi, adf, adz, suc, sui, suf, suz, muc, mui, muf, muz, dic, dii, dif, diz, // 128..159
		out, til, nil, cnt, str, sqr, wer, epv, ech, ecp, fst, abs, enl, neg, val, riv, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, lst, nil, grd, grp, gdn, unq, // 160..191
		typ, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, scn, liv, spl, srt, flr, // 192..223
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, kst, nil, nil, nil, odf, nil, nil, rnd, nil, tok, unf, nil, nil, nil, nil, nil, ovr, rev, jon, not, nil, // 224..255
	})
	sJ(0, 289360742959022340) // type sizes uint64(0x0404041008040104)
	sI(12, 0x70881342)        // rng state
	sI(128, x)                // alloc
	p := i(256)
	for i := i(8); i < x; i++ {
		sI(4*i, p) // free pointer
		p *= 2
	}
	sI(kkey, enl(mk(1, 0)))
	sI(kval, enl(0))
	sI(xyz, cat(cat(mks(120), mks(121)), mks(122)))

	r := mk(1, 3)
	sI(r+8, 7370084)    // "dup"
	dx(sc(r))           // -> offset 0x18
	r = mk(1, 4)        //
	sI(r+8, 1751349349) // "exch"
	dx(sc(r))           // -> offset 0x1c

	sI(kcon, mk(6, 0))
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
}
func kC(x []byte) (r uint32) {
	r = mk(1, uint32(len(x)))
	copy(MC[r+8:], x)
	return r
}
func grow(x i) (r i) {
	if x > 31 {
		panic("Ω")
	}
	c := make([]uint64, 1<<(x-3))
	copy(c, MJ)
	MJ = c
	msl()
	return x
}
func bk(t, n i) (r i) {
	r = i(32 - bits.LeadingZeros32(7+n*i(C(t))))
	if r < 4 {
		return 4
	}
	return r
}
func mk(x, y i) (r i) {
	// defer func() { fmt.Printf("mk r=%d\n", r) }()
	t := bk(x, y)
	i := 4 * t
	m := 4 * I(128)
	for I(i) == 0 {
		if i >= m {
			sI(128, grow(1+i/4))
			for j := m; j <= i; j += 4 { // todo: not in w.k/k.k (mark intermediates for large allocations)
				sI(j, 1<<(j>>2))
			}
			m = i
			i -= 4
		}
		i += 4
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
		if xr == 0 {
			if leaktest {
				return // interned constants may already be released
			}
			fmt.Printf("x=%d(%x) %d/%d\n", x, x, tp(x), nn(x))
			dump(x, 200)
			panic("dx free")
		}
		sI(x+4, xr-1)
		if xr == 1 {
			xt, xn, xp := v1(x)
			if xt == 0 || xt > 5 {
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
	if y < 0 {
		panic("rxn<0")
	}
	if x > 255 {
		MI[1+x>>2] += y
	}
}
func rl(x i) {
	xt, xn, xp := v1(x)
	//fmt.Printf("rl %d/%d\n", xt, xn)
	if xt > 0 && xt < 6 {
		panic(fmt.Sprintf("rl on %d/%d", xt, xn))
	}
	for i := i(0); i < xn; i++ {
		rx(I(xp))
		xp += 4
	}
}
func rld(x i)          { rl(x); dx(x) }
func kvd(x i) (i, i)   { rld(x); return I(x + 8), I(x + 12) }
func dxr(x, r i) i     { dx(x); return r }
func dxyr(x, y, r i) i { dx(x); dx(y); return r }
func mki(i i) (r i)    { r = mk(2, 1); sI(r+8, i); return r }
func mkc(c i) (r i)    { r = mk(1, 1); sC(r+8, byte(c)); return r }
func mks(c i) (r i)    { return sc(mkc(c)) }
func mkd(x, y i) (r i) {
	if tp(x) != 5 {
		r = mk(4, 1)
		sF(r+8, float64(0))
		sF(r+16, float64(1))
		return add(x, mul(y, r))
	}
	xn := nn(x)
	if xn == 1 && nn(y) != 1 {
		y = enl(y)
	}
	y = lx(y)
	if xn != nn(y) {
		trap()
	}
	r = l2(x, y)
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
	if xn == 1 {
		return take(x, yn), y
	}
	if yn == 1 {
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
	case 3:
		for i := i(0); i < n; i++ {
			sF(rp, F(xp))
			sF(rp+8, 0.0)
			xp += 8
			rp += 16
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
	if xt == 7 || n == 1 {
		return enl(x)
	}
	if xt == 0 {
		trap()
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
	if xt == 4 {
		return zim(x)
	}
	if xt == 6 {
		return ech(x, 161)
	}
	if xt == 7 {
		r = I(xp)
		rx(r)
		return dxr(x, r)
	}
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
	if tp(x) == 7 {
		k, v := kvd(x)
		return mkd(rev(k), rev(v))
	}
	n := nn(x)
	if n < 2 {
		return x
	}
	return atx(x, tir(n))
}
func fst(x i) (r i) {
	xt, xn, _ := v1(x)
	if xn == 0 {
		dx(x)
		if xt == 0 {
			return 0
		}
		if xt == 5 {
			return sc(mk(1, 0))
		}
		if xt > 5 {
			return mk(6, 0)
		}
		return cst(mki(xt), mkc(0))
	}
	if xt == 0 {
		return x
	}
	if xt == 7 {
		return fst(val(x))
	}
	return atx(x, mki(0))
}
func lst(x i) (r i) { // *|
	if tp(x) == 7 {
		return lst(val(x))
	}
	if nn(x) == 0 {
		return fst(x)
	}
	return atx(x, mki(nn(x)-1))
}
func drop(x, n i) (r i) {
	xt, xn, _ := v1(x)
	a := n
	if xn == 0 {
		return x
	}
	if int32(n) < 0 {
		n = -n
		a = 0
	}
	if n > xn {
		return dxr(x, mk(xt, 0))
	}
	x = atx(x, seq(a, xn-n, 1))
	if xt == 6 && xn-n == 1 {
		x = enl(x)
	}
	return x
}
func cut(x, y i) (r i) {
	xt, yt, xn, yn, xp, _ := v2(x, y)
	if yt == 7 {
		if xt == 2 {
			if xn != 1 {
				trap()
			}
			k, v := kvd(y)
			rx(x)
			return mkd(cut(x, k), cut(x, v))
		}
		rx(y)
		return tkd(exc(til(y), x), y)
	}
	if xt == 0 {
		rx(y)
		return atx(y, wer(not(atx(x, y))))
	}
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
	if y == 0 {
		panic("rsh")
	}
	xt, yt, xn, _, xp, _ := v2(x, y)
	if yt == 7 {
		if xt == 0 && nn(x) == 5 {
			rx(y)
			x = expr(x, y)
			if tp(x) == 2 {
				return ecl(y, wer(x), 64)
			} else {
				return dxr(y, x)
			}
		} else if xt == 2 {
			if xn != 1 {
				trap()
			}
			k, v := kvd(y)
			rx(x)
			return mkd(rsh(x, k), rsh(x, v))
		}
		return tkd(x, y)
	}
	if xt == 0 {
		rx(y)
		return atx(y, wer(atx(x, y)))
	}
	if xt != 2 {
		panic("type")
	}
	n := prod(xp, xn)
	r = take(y, n)
	if xn == 1 {
		if yt == 6 && n == 1 {
			r = enl(r)
		}
		return dxr(x, r)
	}
	xn--
	xe := xp + 4*xn
	for i := i(0); i < xn; i++ {
		m := I(xe)
		n /= m
		n = prod(xp, xn-i)
		r = cut(seq(0, n, m), r)
		if m == 1 && i > 0 {
			if tp(r) != 6 {
				r = enl(r)
			} else {
				r = ech(r, 172)
			}
		}
		xe -= 4
	}
	if I(xe) == 1 {
		r = enl(r)
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
	if 0 == nn(x) {
		x = fst(x)
	}
	xn := nn(x)
	o := i(0)
	if int32(n) < 0 {
		o = xn + n
		n = -n
		if int32(o) < 0 {
			return x
		}
	}
	r = seq(o, n, 1)
	if xn < n {
		rp := 8 + r
		for i := i(0); i < n; i++ {
			sI(rp, I(rp)%xn)
			rp += 4
		}
	}
	return atx(x, r)
}
func tkd(x, y i) (r i) {
	t := tp(x)
	k, v := kvd(y)
	if t != 5 {
		trap()
	}
	rx(k)
	x = fnd(k, x)
	rx(x)
	v = atx(v, x)
	if nn(x) == 1 {
		v = enl(v)
	}
	return mkd(atx(k, x), v)
}
func phi(x, y i) (r i) {
	n := nn(y)
	r = mk(4, n)
	rp := r + 8
	yp := y + 8
	for i := i(0); i < n; i++ {
		p := 0.017453292519943295 * F(yp)
		sF(rp, math.Cos(p))
		sF(rp+8, math.Sin(p))
		rp += 16
		yp += 8
	}
	dx(y)
	return mul(x, r)
}
func atm(x, y i) (r i) { // {$[0~#y;:0#x;];}
	if 0 == nn(y) {
		return dxr(x, y)
	}
	rx(y)
	f := fst(y)
	t := drop(y, 1)
	if nn(t) > 1 {
		return atm(atx(x, f), t)
	}
	t = fst(t)
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
	if yt == 0 {
		return rsh(y, x)
	}
	if xt == 7 { // d@..
		return atd(x, y, yt)
	}
	if yt > 5 { // at-list at-dict
		return ecr(x, y, '@') //64
	}
	if yt == 3 && xt < 5 {
		return phi(x, y)
	}
	if yt != 2 {
		fmt.Printf("x=%s y=%s\n", X(x), X(y))
		fmt.Printf("atx xt/xn %d/%d yt~I/yn %d/%d\n", xt, xn, yt, yn)
		panic("atx yt~I")
	}
	r = mk(xt, yn)
	rp := r + 8
	w := i(C(xt))
	for i := i(0); i < yn; i++ {
		yi := I(yp)
		if xn <= yi {
			fmt.Printf("atx xt/xn %d/%d @%d\n", xt, xn, yi)
			trap()
		}
		mv(rp, xp+w*yi, w)
		rp += w
		yp += 4
	}
	if xt > 5 {
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
	y = lx(y)
	xt, _, xn, yn, xp, yp := v2(x, y)
	if xt != 0 {
		return atm(x, y)
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
		r = f(I(yp), I(yp+4))
		return r
	} else if x < 256 {
		if yn != 1 {
			panic("arity")
		}
		f, ok := MT[x].(func(i) i)
		if ok == false {
			fmt.Printf("cal? x=%d (%s) y=%s\n", x, X(x), X(y))
			panic("!call")
		}
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
	if xn == 3 { // proj
		rl(x)
		if yn == 1 {
			y = fst(y)
		}
		r = asi(I(x+12), I(x+16), y)
		v := I(x + 8)
		dx(x)
		return cal(v, r)
	}
	if xn == 5 { // lambda
		a := I(x + 20)
		if a > yn {
			a -= yn
			for i := i(0); i < a; i++ {
				y = lcat(y, 0)
			}
			return prj(x, y, seq(yn, a, 1))
		}
		return lcl(x, y)
	}
	panic("nyi")
}
func lcl(x, y i) (r i) {
	fn := I(x + 20)
	if fn == 0 {
		dx(y)
		y = mk(6, 0)
	}
	if nn(y) != fn {
		panic("arity")
	}
	yp := y + 8
	a := I(x + 16)
	ap := a + 8
	an := nn(a)
	l := mk(2, an) // stores reference which should not be modified
	lp := l + 8
	sp := I(kval)
	for i := i(0); i < an; i++ {
		d := sp + I(ap)
		sI(lp, I(d))
		v := uint32(0)
		if i < fn {
			v = I(yp)
			rx(v)
			yp += 4
		}
		sI(d, v)
		ap += 4
		lp += 4
	}
	dx(y)
	b := I(x + 12)
	s := I(x + 24)
	t := I(x + 8)
	rx(b)
	rx(t)
	rx(s)
	r = run(l3(b, s, t))
	sp = I(kval) // could have changed
	lp = l + 8
	ap = a + 8
	for i := i(0); i < an; i++ {
		d := sp + I(ap)
		dx(I(d))
		sI(d, I(lp))
		ap += 4
		lp += 4
	}
	dx(l)
	return dxr(x, r)
}
func expr(x, y i) (r i) { // expr(lambda, dict)
	if tp(x) != 0 || nn(x) != 5 || MI[2+3+x>>2] != 0 {
		panic("expr")
	}
	k, v := kvd(y)
	rxn(k, 2)
	env := ech(k, 46)   // save    env: .'!t
	dx(ecd(k, v, 58))   // (!t):'(.t)
	r = atx(x, 0)       // r:x[]
	dx(ecd(k, env, 58)) // restore (!t):'env
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
	if xt == 7 {
		return dcat(x, y)
	}
	if yt == 6 {
		return ucat(lx(x), y)
	}
	return ucat(lx(x), lx(y))
}
func ucat(x, y i) (r i) {
	xt, yt, xn, yn, xp, yp := v2(x, y)
	if xt != yt {
		fmt.Printf("ucat %x,%x %d/%d %d/%d %s %s\n", x, y, xt, xn, yt, yn, X(x), X(y))
		panic("!ucat")
	}
	if xt > 5 {
		rl(x)
		rl(y)
	}
	if xt == 7 {
		r = mkd(ucat(I(x+8), I(y+8)), ucat(I(x+12), I(y+12)))
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
func ic(x, y i) (r i) { return ucat(x, mki(y)) }
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
func dcat(x, y i) (r i) { // d,v
	k, v := kvd(x)
	return mkd(k, ecd(v, y, 44))
}
func enl(x i) (r i) { r = mk(6, 1); sI(8+r, x); return r }
func cnt(x i) (r i) {
	if tp(x) == 7 {
		x = til(x)
	}
	r = mki(nn(x))
	return dxr(x, r)
}
func typ(x i) (r i) {
	xt, _, _ := v1(x)
	r = mk(2, 1)
	sI(8+r, xt)
	return dxr(x, r)
}
func wer(x i) (r i) {
	xt, xn, xp := v1(x)
	//if xt == 1 {
	//	return prs(x)
	//}
	if xt == 4 {
		return zan(x, xn, xp)
	}
	if xt == 6 {
		return flp(x)
	}
	if xt != 2 {
		panic("type")
	}
	rn := i(0)
	for i := i(0); i < xn; i++ {
		rn += I(xp + 4*i)
	}
	r = mk(2, rn)
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
	case 1:
		nn = xn
	case 2:
		nn = xn << 2
	case 3:
		nn = xn << 3
	case 4:
		nn = xn << 4
	case 5:
		nn = xn << 2
	default:
		if xt == 0 && x < 256 {
			return boolvar(x == y)
		}
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
func not(x i) (r i) {
	t := tp(x)
	if t > 5 {
		return ech(x, 126)
	}
	if t == 0 {
		if x == 0 {
			return mki(1)
		}
		dx(x)
		return mki(0)
	}
	return eql(mki(0), x)
}
func cmp(x, y, eq i) (r i) {
	x, y = upxy(x, y)
	x, y = ext(x, y)
	if tp(x) == 6 {
		return ecd(x, y, 62-eq)
	}
	t, _, n, _, xp, yp := v2(x, y)
	var cm func(i, i) i
	if eq == 1 {
		cm = MT[t+8].(func(i, i) i)
	} else {
		cm = MT[t].(func(i, i) i) //nil for Z
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
func unf(x i) (r i) { // unify list 'ux
	xt, xn, xp := v1(x)
	if xt != 6 || xn == 0 {
		return x
	}
	t := tp(I(xp))
	for i := i(0); i < xn; i++ {
		r = I(xp)
		if tp(r) != t || nn(r) != 1 {
			return x
		}
		xp += 4
	}
	return ovr(x, 44) // ,/x
}
func eql(x, y i) (r i) {
	xt, xn, _ := v1(x)
	yt := tp(y)
	if xt == 5 && yt == 7 { // `g=t table group
		rx(y)
		rxn(x, 2)
		r = atx(y, x)
		if xn > 1 {
			r = flp(r)
		}
		r = grp(r)
		y = cut(x, y)
		u, q := kvd(r)
		k, v := kvd(y)
		v = ecl(v, q, 64) // v@\:x
		if xn > 1 {
			u = flp(u)
			u = ech(u, 'u')
		}
		return l2(mkd(x, u), mkd(k, v))
	}
	return cmp(x, y, 1)
}
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
func rnd(x i) (r i) { // 'r x
	if I(x) != 1073741825 { // int#1
		trap()
	}
	dx(x)
	x = I(x + 8)
	r = mk(2, x)
	rp := r + 8
	for i := i(0); i < x; i++ {
		sI(rp, rng())
		rp += 4
	}
	return r
}
func rng() (r i) {
	r = I(12)
	r ^= (r << 13)
	r ^= (r >> 17)
	r ^= (r << 5)
	sI(12, r)
	return r
}
func lop(x, y, l i) (r i) {
	t := tp(y)
	if t == 0 {
		return fxp(x, y, l)
	}
	if t == 6 {
		y, f := kvd(y) // not a dict
		return whl(y, x, f, l)
	}
	dx(l)
	return 0
}
func jon(x, y i) (r i) { // y/:x (join)
	r = lop(x, y, 0)
	if r != 0 {
		return r
	}
	xt, xn, xp := v1(x)
	if xn == 0 {
		r = tp(y)
		dx(y)
		return dxr(x, mk(r, 0))
	} else if xn == 1 {
		return dxr(y, fst(x))
	}
	if xt != 6 {
		return dxr(y, x) // allow ","/:"abc" -> "abc"
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
	r = lop(x, y, enl(mk(6, 0)))
	if r != 0 {
		return r
	}
	rx(x)
	yn := nn(y)
	r = fds(x, y)
	if nn(r) == 0 {
		dx(r)
		return enl(x)
	}
	r = cut(cat(mki(0), r), x)
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
	if yn == 0 {
		return dxyr(x, y, drop(seq(0, xn, 1), 1))
	}
	if xn < yn {
		return dxyr(x, y, mk(2, 0))
	}
	r = mk(2, 0)
	eq := MT[8+xt].(func(i, i) i)
	w := i(C(xt))
	for i := i(0); i < xn; i++ {
		a := uint32(0)
		for j := uint32(0); j < yn; j++ {
			k := w * j
			a += eq(xp+k, yp+k)
		}
		if a == yn {
			r = ic(r, i)
			i += yn - 1
			xp += w * (yn - 1)
		}
		xp += w
	}
	return dxyr(x, y, r)
}
func exc(x, y i) (r i) { n := mki(nn(y)); rx(x); return atx(x, wer(eql(n, fnd(y, x)))) } // x^y
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
func unq(x i) (r i) { // ?x uniq  x@&(x?x)=!#x
	xt, xn, xp := v1(x)
	r = mk(2, 0)
	w := i(C(xt))
	for i := i(0); i < xn; i++ {
		m := fnx(x, xp)
		if m == i {
			r = ic(r, i)
		}
		xp += w
	}
	return atx(x, r)
}
func grp(x i) (r i) { // =x (group)  &'u=\:x
	xt, xn, xp := v1(x)
	rx(x)
	u := unq(x)
	l := mk(6, 0)
	n := i(0)
	w := i(C(xt))
	for i := i(0); i < xn; i++ {
		m := fnx(u, xp)
		if m == n {
			l = lcat(l, mki(i))
			n++
		} else {
			p := l + 8 + 4*m
			sI(p, ic(I(p), i))
		}
		xp += w
	}
	return dxr(x, l2(u, l))
}
func flr(x i) (r i) {
	xt, xn, xp := v1(x)
	if xt > 5 {
		return ech(x, 223)
	}
	if xt == 0 {
		return dxr(x, mki(x))
	}
	if xt == 1 {
		dx(x)
		return i(C(xp))
	}
	if xt == 2 {
		r = mk(1, xn)
		rp := r + 8
		for i := i(0); i < xn; i++ {
			sC(rp+i, c(I(xp)))
			xp += 4
		}
		return dxr(x, r)
	}
	if xt == 3 {
		r = mk(2, xn)
		rp := r + 8
		for i := i(0); i < xn; i++ {
			sI(rp, uint32(int32(F(xp))))
			xp += 8
			rp += 4
		}
		return dxr(x, r)
	}
	if xt == 4 {
		return zre(x)
	}
	trap()
	return x
}
func ang(x, y f) float64 {
	p := 57.29577951308232 * math.Atan2(y, x)
	if p < 0 {
		p += 360.0
	}
	return p
}
func flp(x i) (r i) { // flip/transpose {n:#*x;(,/x)(n*!#x)+/:!n}
	n := nn(I(x + 8))
	m := nn(x)
	return atx(ovr(x, ','), ecr(mul(mki(n), seq(0, m, 1)), seq(0, n, 1), '+')) //44 43
}
func drw(x, y i) (r i) { // x 'd y
	return dxyr(x, y, 0)
}
func out(x i) (r i) {
	if Out != nil {
		Out(x)
		return x
	}
	rx(x)
	r = x
	if tp(x) != 1 {
		r = kst(x)
	}
	n := nn(r)
	if n > 0 {
		fmt.Printf("%s\n", string(MC[r+8:r+8+n]))
	}
	dx(r)
	return x
}
func kst(x i) (r i) {
	t := tp(x)
	if nn(x) == 0 && t > 1 && t < 6 {
		dx(x)
		r = cc(cc(mkc('0'), '#'), '0') //48 35
		if t == 3 {
			r = cc(r, '.') //46
		}
		if t == 4 {
			r = cc(r, 'a')
		}
		if t == 5 {
			sC(r+10, '`') //96
		}
		return r
	}
	if t == 7 {
		r, t := kvd(x)
		r = cc(kst(r), '!') //33
		if nn(t) == 0 {     // 0#`!0
			dx(t)
			t = mki(0)
		}
		return ucat(r, kst(t))
	}
	if t == 6 {
		if nn(x) == 1 {
			return ucat(mkc(','), kst(fst(x))) //44
		}
		x = ech(x, 235) // 'k-each
	} else {
		x = str(x)
	}
	switch t {
	case 0:
		r = x
	case 1:
		r = cc(ucat(mkc('"'), x), '"') //34   todo quote
	case 5:
		r = ucat(mkc('`'), jon(x, mkc('`'))) //96
	case 6:
		r = cc(ucat(mkc('('), jon(x, mkc(';'))), ')') //40 59 41
	default:
		r = jon(x, mkc(' ')) //32  2 3 4
	}
	return r
}
func str(x i) (r i) {
	xt, xn, xp := v1(x)
	if xt == 1 {
		return x
	}
	if xt == 0 {
		return cg(x, xn)
	}
	if xt > 5 || xn != 1 {
		return ech(x, 164)
	}
	switch xt {
	case 2:
		n := I(xp)
		r = ci(n)
	case 3:
		r = cf(F(xp))
	case 4:
		r = cz(F(xp), F(xp+8))
	case 5:
		rx(x)
		r = cs(x)
	default:
		panic("nyi")
	}
	return dxr(x, r)
}
func cg(x i, xn i) (r i) {
	if x == 0 || x == 128 {
		return mk(1, 0)
	}
	if x < 127 {
		r = mkc(x)
		// {[()]} should be replaced with adverbs
	} else if x < 256 {
		r = cc(mkc(x-128), ':') //58
	} else if xn == 2 {
		rl(x)
		r = cat(str(I(x+12)), str(I(x+8)))
	} else if xn == 3 {
		rl(x)
		dx(I(x + 16))
		r = kst(I(x + 12))
		sC(r+8, '[')       //91
		sC(r+7+nn(r), ']') //93
		r = ucat(str(I(x+8)), r)
	} else if xn == 5 { // todo arity
		r = I(x + 8)
		rx(r)
	} else {
		panic("nyi cg")
	}
	return dxr(x, r)
}
func ng(x, y i) (r i) {
	if y != 0 {
		x = ucat(mkc('-'), x) //45
	}
	return x
}
func ci(n i) (r i) {
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
		r = cc(r, i('0'+c))
		n /= 10
	}
	if nn(r) == 0 {
		r = cc(r, '0')
	}
	return ng(rev(r), m)
}

/*
func hf(f float64) string {
	b := make([]byte, 8)
	u := math.Float64bits(f)
	binary.LittleEndian.PutUint64(b, u)
	return "0x"+hex.EncodeToString(b)
}
*/
func cf(f float64) (r i) {
	if f != f {
		return cc(mkc('0'), 'n') //48 110
	}
	if f == 0 {
		return cc(cc(mkc('0'), '.'), '0') //48 46 48
	}
	m := i(0)
	if f < 0 {
		m = 1
		f = -f
	}
	if f > 1.7976931348623157e+308 { // 0xffffffffffffef7f, see hf()
		return ng(cc(mkc('0'), 'w'), m) //119
	}
	e := i(0)
	for f > 1000 {
		e += 3
		f /= 1000.0
	}
	d := i(7)
	if f < 1 {
		d++
		if f < 0.1 {
			d++
			if f < 0.01 { // 0x7b14ae47e17a843f
				d++
				if f < 0.001 { // 0xfca9f1d24d62503f
					d = 7
					for f < 1 {
						e -= 3
						f *= 1000.0
					}
				}
			}
		}
	}

	n := int32(f)
	r = ci(uint32(n))
	f -= float64(n)
	d -= nn(r)
	if int32(d) < 1 {
		d = 1
	}
	r = cc(r, '.') //46
	t := i(0)
	for i := i(0); i < d; i++ {
		f *= 10
		n = int32(f)
		r = cc(r, uint32('0'+n))
		f -= float64(n)
		t = (1 + t) * boolvar(i > 0 && n == 0)
	}
	r = drop(r, -t)
	if e != 0 {
		r = ucat(cc(r, 'e'), ci(e))
	}
	return ng(r, m)
}
func cz(x, y f) (r i) {
	a := math.Hypot(x, y)
	p := uint32(0.5 + ang(x, y))
	return ucat(cc(cf(a), 'a'), ci(p))
}
func cst(x, y i) (r i) { // x$y
	xt, yt, xn, yn, _, _ := v2(x, y)
	if xt == 5 {
		if yt == 1 { // `$"abc" /`abc
			dx(x)
			return sc(y)
		} else if yt == 6 {
			return unf(ecr(x, y, 36))
		}
	}
	if xt != 2 || xn != 1 {
		trap()
	}
	dx(x)
	x = I(x + 8)
	if int32(x) < 0 && yt == 1 { // -3$0x123456.. (raw bytes to ifz)
		x = -x
		n := yn / i(C(x))
		if i(C(x))*n != yn {
			trap()
		}
		r = use(y)
		sI(r, n|x<<29)
		return r
	}
	if yn == 0 {
		dx(y)
		if x == 7 {
			return mkd(mk(5, 0), mk(6, 0))
		}
		return mk(x, 0)
	}
	if yt > x || yt > 4 { // flr?
		trap()
	}
	if x == 8 { // 8$ifz (raw bytes)
		n := yn * i(C(yt))
		r = use(y)
		sI(r, n|1<<29)
		return r
	}
	for yt < x {
		y = up(y, yt, yn)
		yt++
	}
	return y
}
func sc(x i) (r i) {
	k := I(kkey)
	n := nn(k)
	x = enl(x)
	r = fnx(k, x+8)
	if r < n {
		dx(x)
	} else {
		sI(kkey, cat(k, x))
		sI(kval, lcat(I(kval), 0))
	}
	r = mki(8 + 4*r)
	sI(r, 1|5<<29)
	return r
}
func cs(x i) (r i) {
	r = I(8 + x)
	r = I(I(kkey) + r)
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
func nmf(x, y i) (r i) { // 'F (sin cos exp log)
	dx(x)
	x = I(x + 8)
	y = use(cst(mki(3), y))
	n := nn(y)
	yp := y + 8
	g := MT[x].(func(float64) float64)
	for i := i(0); i < n; i++ {
		sF(yp, g(F(yp)))
		yp += 8
	}
	return y
}
func nm(x, f, h i) (r i) { // numeric monad f:scalar index (nec..), h:original func (e.g. -: neg)
	if tp(x) > 5 {
		return ech(x, h)
	}
	r = use(x)
	t, n, rp := v1(r)
	xp := x + 8
	w := uint32(C(t))
	f += t
	g := MT[f].(func(i, i))
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
func eqf(x, y i) i { return boolvar(I(x) == I(y) && I(x+4) == I(y+4)) }
func eqz(x, y i) i { return eqf(x, y) * eqf(x+8, y+8) }
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
func imod(x, y int32) (r i) { return i((x%y + y) % y) }
func eqL(x, y i) i          { return match(I(x), I(y)) }
func adc(x, y, r i)         { sC(r, C(x)+C(y)) }
func adi(x, y, r i)         { sI(r, I(x)+I(y)) }
func adf(x, y, r i)         { sF(r, F(x)+F(y)) }
func adz(x, y, r i)         { sZ(r, Z(x)+Z(y)) }
func suc(x, y, r i)         { sC(r, C(x)-C(y)) }
func sui(x, y, r i)         { sI(r, I(x)-I(y)) }
func suf(x, y, r i)         { sF(r, F(x)-F(y)) }
func suz(x, y, r i)         { sZ(r, Z(x)-Z(y)) }
func muc(x, y, r i)         { sC(r, C(x)*C(y)) }
func mui(x, y, r i)         { sI(r, I(x)*I(y)) }
func muf(x, y, r i)         { sF(r, F(x)*F(y)) }
func muz(x, y, r i)         { sZ(r, Z(x)*Z(y)) }
func dic(x, y, r i)         { sC(r, C(x)/C(y)) }
func dii(x, y, r i)         { sI(r, i(int32(I(x))/int32(I(y)))) }
func dif(x, y, r i)         { sF(r, F(x)/F(y)) }
func diz(x, y, r i)         { sZ(r, Z(x)/Z(y)) }
func moi(x, y, r i)         { sI(r, imod(int32(I(x)), int32(I(y)))) }
func add(x, y i) i          { return nd(x, y, 15+128, 43) }
func sub(x, y i) i          { return nd(x, y, 19+128, 45) }
func mul(x, y i) i          { return nd(x, y, 23+128, 42) }
func diw(x, y i) i          { return nd(x, y, 27+128, 37) }
func mod(x, y i) i          { return nd(x, y, 23, 7) }
func abs(x i) i             { return nm(x, 15, 171) }
func neg(x i) i             { return nm(x, 19, 173) }
func sqr(x i) (r i) {
	xt, xn, xp := v1(x)
	if xt == 2 && xn == 1 {
		r = I(xp)
		m := r & 0xf
		if m == 0 {
			if tp(r) < 8 && I(r+4) > 0 {
				rx(r)
				return dxr(x, r)
			}
		} else if m == 5 {
			r >>= 4
			k := I(kkey)
			if r%4 == 0 && r-8 < 4*nn(k) {
				dx(x)
				x = mk(5, 1)
				sI(x+8, r)
				return x
			}
		}
	}
	return nm(x, 27, 165)
}
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
func sqi(x, r i) { panic("%i") }
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
func zan(x, xn, xp i) (r i) {
	r = mk(3, xn)
	rp := r + 8
	for i := i(0); i < xn; i++ {
		sF(rp, ang(F(xp), F(xp+8)))
		xp += 16
		rp += 8
	}
	return dxr(x, r)
}
func agg(x, y i) (r i) { // (*;+/)'d list-each dict-each
	var k, v, a, b i
	if tp(x) == 7 {
		k, v = kvd(x)
	} else { // 6
		rx(x)
		v = x
		k = ovr(ecd(sc(mk(1, 0)), str(x), 36), 44) // ,/`$'$x
	}
	a, b = kvd(y)
	a = ucat(sc(mk(1, 0)), a) // ``a`b`c
	n := nn(v)
	rxn(b, n)
	rld(v)
	p := 8 + v
	t := mk(6, 0)
	for i := i(0); i < n; i++ {
		t = lcat(t, ech(b, I(p)))
		p += 4
	}
	dx(b)
	return mkd(a, cat(enl(k), flp(t)))
}

func drv(x, y i) (r i) { // x(adv) y(verb), e.g. ech +
	r = mk(0, 2)
	sI(8+r, x)
	sI(12+r, y)
	return r
}
func ecv(x i) (r i) { return drv(40, x) }  // '  ech(168) ecd(40)
func epv(x i) (r i) { return drv(41, x) }  // ': ecp(169) epi(41)
func ovv(x i) (r i) { return drv(123, x) } // /  ovr(251) ovi(123)
func riv(x i) (r i) { return drv(125, x) } // /: jon(253) ecr(125)
func scv(x i) (r i) { return drv(91, x) }  // \  scn(219) sci(91)
func liv(x i) (r i) { return drv(93, x) }  // \: spl(221) ecl(93)
func ech(x, y i) (r i) { // f'x (each)
	if tp(y) != 0 {
		if tp(y) >= 6 && tp(x) == 7 { // (*;+/)'t
			return agg(y, x)
		}
		return bin(y, x)
	}
	if tp(x) == 7 {
		k, v := kvd(x)
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
func ovr(x, y i) (r i) { // y/x (over/reduce)
	t := tp(y)
	if t == 2 {
		return mod(x, y)
	}
	return ovs(x, y, 0, 0)
}
func scn(x, y i) (r i) { // y\x (scan)
	t := tp(y)
	if t != 0 && t < 5 {
		return diw(x, y) // y%x (flipped)
	}
	return ovs(x, y, enl(mk(6, 0)), 0)
}
func ovi(x, y, z i) (r i) { return ovs(y, z, 0, x) }             // z y/ x (over initial)
func sci(x, y, z i) (r i) { return ovs(y, z, enl(mk(6, 0)), x) } // z y/ x (scan initial)
func ovs(x, y, z, l i) (r i) { // over/scan
	n := nn(x)
	if n == 0 && l == 0 {
		x = enl(x)
		n = 1
	}
	rxn(x, n)
	r = l
	o := i(1)
	if r == 0 {
		r = fst(x)
		o = 0
		n--
		scl(z, r)
	}
	rxn(y, n)
	for i := i(0); i < n; i++ {
		r = cal(y, l2(r, atx(x, mki(i+1-o))))
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
func fxp(x, y, z i) (r i) { // f/:x f\:x fixed point/converge
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
func ecr(x, y, f i) (r i) { // x f/: y (each-right)
	if tp(y) == 7 {
		k, v := kvd(y)
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
		rx(f)
		r = atx(f, r)
		scl(s, r)
	}
}
func nlp(x, f, s, n i) (r i) { // (n;f)/:y (for)  (n;f)\:y (scan-for)
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
func ecl(x, y, f i) (r i) { // x f\: y (each-left)
	if tp(x) == 7 {
		k, v := kvd(x)
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
func sjn(xp, yp i) (r i) {
	xp = I(I(kkey) + xp)
	yp = I(I(kkey) + yp)
	rx(xp)
	rx(yp)
	r = sc(ucat(xp, yp))
	return dxr(r, I(r+8))
}
func ecd(x, y, f i) (r i) { // x f' y
	if tp(x) == 5 && tp(y) == 7 {
		if tp(f) == 0 {
			f = enl(f)
		}
		if tp(f) == 6 {
			if nn(f) > 1 {
				rx(f)
				r = cst(mk(5, 0), str(f))
			} else {
				r = sc(mk(1, 0))
			}
			f = mkd(r, f)
		}
		var s i
		s, f = kvd(f)
		m := nn(f)
		u := eql(x, y)
		x, y = kvd(u)
		x, u = kvd(x)
		r, y = kvd(y)
		n := nn(y)
		yp := y + 8
		k := mk(5, m*n)
		kp := k + 8
		rp := r + 8
		for i := i(0); i < n; i++ {
			fp := f + 8
			sp := s + 8
			for j := uint32(0); j < m; j++ {
				rx(I(fp))
				rx(I(yp))
				u = lcat(u, ovr(ech(I(yp), I(fp)), 44))
				sI(kp, sjn(I(sp), I(rp)))
				fp += 4
				kp += 4
				sp += 4
			}
			yp += 4
			rp += 4
		}
		dx(y)
		dx(f)
		dx(r)
		dx(s)
		return mkd(ucat(x, k), u)
	}
	x, y = ext(x, y)
	n := nn(x)
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
	xt, xn, _ := v1(x)
	switch xt {
	case 0:
		if x < 256 {
			return x
		}
		rl(x)
		r = mk(6, xn)
		mv(r+8, x+8, 4*xn)
		if xn == 5 {
			sI(r+20, mki(I(r+20)))
		}
		dx(x)
	case 1:
		x = prs(x)
		rx(x)
		r = lst(fst(x))
		dx(r)
		n := I(r+8) == 930
		r = run(x)
		if n {
			dx(r)
			return 0
		}
	case 2:
		// r = run(l2(x, take(mki(0), nn(x))))
	case 5:
		r = lup(x)
	case 6:
		return run(x)
	case 7:
		r = I(x + 12)
		rx(r)
		dx(x)
	default:
		fmt.Printf("val xt=%d\n", xt)
		panic("nyi")
	}
	return r
}
func lup(x i) (r i) {
	if I(x+8) == 0 {
		//panic("lup#0!!")
	}
	r = I(I(kval) + I(x+8))
	rx(r)
	return dxr(x, r)
}

func con(x i) (r i) { // intern constant
	l := I(kcon)
	lp := l + 8
	for i := i(0); i < nn(l); i++ {
		r = I(lp)
		if match(r, x) != 0 {
			dx(x)
			return r
		}
		lp += 4
	}
	sI(kcon, lcat(l, x))
	return x
}
func asn(x, y i) (r i) {
	//fmt.Printf("ASN %s <- %s\n", X(x), X(y))
	xt, _, _ := v1(x)
	if xt != 5 {
		trap()
	}
	if I(x+8) == 0 {
		panic("asn#0!!")
	}
	p := I(kval) + I(x+8)
	dx(I(p))
	sI(p, y)
	dx(x)
	rx(y)
	return y
}
func asi(x, y, z i) (r i) { //x[..y..]:z
	xt, yt, xn, yn, _, yp := v2(x, y)
	if xt == 7 && yt < 6 {
		k, v := kvd(x)
		if yt == 5 {
			rx(k)
			y = fnd(k, y)
		}
		return mkd(k, asi(v, y, z))
	}
	if yt == 6 {
		if xt == 7 { // (:;d;y;z)  {k:!x;v:. x;y:(,k?*y),1_y;k!.(:;v;y;z)}
			k, v := kvd(x)
			rx(y)
			f := fst(y)
			if f != 0 {
				if tp(f) == 5 {
					rx(k)
					f = fnd(k, f)
				}
			} else {
				f = seq(0, nn(k), 1)
			}
			return mkd(k, asi(v, cat(enl(f), drop(y, 1)), z))
		}
		if xt != 6 || yt != 6 {
			trap()
		}
		r = take(x, xn)
		if xn == 1 {
			r = enl(r)
		}
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
	if xt < 6 {
		if zt != xt {
			trap()
		}
		r = use(x)
		rp := r + 8
		w := i(C(xt))
		for i := i(0); i < yn; i++ {
			k := I(yp)
			mv(rp+w*k, zp, w)
			yp += 4
			zp += w
		}
		return dxyr(y, z, r)
	}
	if xt == 6 {
		r = take(x, xn)
		if xn == 1 {
			r = enl(r)
		}
		if yn == 1 {
			z = enl(z)
			zn = 1
			zt = 6
		}
		if zt != 6 {
			z = lx(z)
		}
		rp := r + 8
		if yn != zn {
			trap()
		}
		zp = z + 8
		rl(z)
		for i := i(0); i < yn; i++ {
			k := I(yp)
			if k >= xn {
				trap()
			}
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
func asd(v, s, a, u i) (r i) { // (+;`x;a;y)
	if v != ':' && v != 186 { //58
		rx(s)
		r = lup(s)
		if a != 0 {
			rx(a)
			r = atx(r, a)
		}
		u = cal(v, l2(r, u))
	}
	r = u
	rx(r)
	if a != 0 {
		rx(s)
		u = asi(lup(s), a, u)
	}
	dx(asn(s, u))
	return r
}

func odf(x i) (r i) {
	fmt.Printf("%s\n", X(x))
	if tp(x) != 0 || nn(x) != 5 {
		panic("odf: no-lambda")
	}
	s := I(x + 8)  // string
	b := I(x + 12) // tree
	m := I(x + 24) // src
	rx(s)
	rx(b)
	rx(m)
	r = l3(b, m, s)
	od(r)
	dx(r)
	return x
}

// nop            0=x      ignore
// const(k-value) 0=x&0xf  push x
// monad      n   1=x&0xf  n=x>>4  0..255
// dyad       n   2=x&0xf  n=x>>4  0..255
// tripple@       3=x&0x7          @[a;b;c]
// quote verb n   4=x&0xf  n=x>>4  0..255 push verb
// var        n   5=x&0xf  n=x>>4  var index
// mkl        n   6=x&0xf  n=x>>4  length
// assign         7=x&0x7  n=4     asd(top)
// jif0       n   8=x&0x7  n=x>>4  rel offset
// j          n   9=x&0x7  n=x>>4  rel offset
// prj           10=x&0x7          project
// drop          15=x&0x7          drop 1
func od(x i) {
	y := I(x + 12)
	z := I(x + 16)
	x = I(x + 8)
	_, _, xn, yn, xp, yp := v2(x, y)
	if xn != yn {
		panic(fmt.Sprintf("length %d %d", xn, yn))
	}
	n, m := i(0), i(0)
	for i := i(0); i < xn; i++ {
		r := I(xp)
		m = r & 0xf
		n = r >> 4
		s := fmt.Sprintf("%04d %6q  ", I(yp), string(C(z+I(yp))))
		if m == 0 {
			ds := X(r)
			if len(ds) > 20 {
				ds = ds[:20] + ".."
			}
			s += ds
		} else if m < 3 {
			s += X(n)
		} else if m == 4 {
			s += "(" + X(n) + ")"
		} else if m == 5 {
			ds := X(I(I(kkey) + n))
			s += ds[1 : len(ds)-1]
		} else if m == 8 {
			s += "jif +" + strconv.Itoa(int(n))
		} else if m == 9 {
			s += "j +" + strconv.Itoa(int(n))
		}
		fmt.Printf(" %8d [%2d 0x%08x] %s\n", r, m, n, s)
		xp += 4
		yp += 4
	}
}
func sum(x i) (r i) {
	_, xn, xp := v1(x)
	r = i(0)
	for i := i(0); i < xn; i++ {
		r += I(xp)
		xp += 4
	}
	return r
}
func run(x i) (r i) { // execute byte code (run don't walk)
	if ddd {
		sum := sum(I(x + 8))
		fmt.Println("RUN", x, sum)
		defer func(x, y i) { fmt.Println("RET", x, y) }(x, sum)
		od(x)
	}
	rl(x)
	y := I(x + 12)
	z := I(x + 16)
	ss := I(140)
	sI(140, z) // store src
	dx(x)
	x = I(x + 8)

	xt, yt, xn, yn, xp, yp := v2(x, y)
	if xt != 2 || yt != 2 {
		panic("type")
	}
	if xn != yn {
		panic("length")
	}
	if xn == 0 {
		dx(x)
		dx(y)
		dx(z)
		return 0
	}
	s := mk(2, 1022)
	se := s + 8 + 4*nn(s)
	sp := s + 4
	t := xp + 4*xn
	a, b := i(0), i(0)
	n, m := i(0), i(0)
	for xp < t {
		r = I(xp)
		m = r & 0xf
		n = r >> 4
		if I(yp) != 0 {
			sI(144, I(yp))
		}
		if r == 0 { //nop
		} else if m == 0 { //k-const
			sp += 4
			sI(sp, a)
			a = r
			rx(a)
		} else if m == 1 { //monad
			f := MT[n].(func(i) i)
			a = f(a)
		} else if m == 2 { //dyad
			b = I(sp)
			sp -= 4
			f := MT[n].(func(i, i) i)
			//fmt.Println(X(a), X(b))
			a = f(a, b)
		} else if m == 3 { //@[x;y;z]
			a = asi(a, I(sp), I(sp-4))
			sp -= 8
		} else if m == 4 { //(+)
			sp += 4
			sI(sp, a)
			a = n
		} else if m == 5 { //var
			sp += 4
			sI(sp, a)
			switch n {
			case 0x18: // dup
			case 0x1c: // exch
				sp -= 4
				a = I(sp)
				sI(sp, I(sp+4))
				if I(xp+4) == 1026 {
					xp += 4 // skip @
					yp += 4
				}
			default:
				a = I(I(kval) + n) //136
			}
			rx(a)
		} else if m == 6 { //mkl
			r = mk(6, n)
			rp := r + 8
			sI(rp, a)
			for i := i(0); i < n-1; i++ {
				rp += 4
				sI(rp, I(sp))
				sp -= 4
			}
			a = r
		} else if m == 7 { //assign
			a = asd(a, I(sp), I(sp-4), I(sp-8))
			sp -= 12
		} else if m == 8 { //rel jump if 0
			if a == 0 || I(8+a) == 0 {
				xp += n
				yp += n
			}
			dx(a)
			a = I(sp)
			sp -= 4
		} else if m == 9 { //rel jump
			xp += n
			yp += n
		} else if m == 10 { // prj
			a = prj(a, I(sp-4), I(sp))
			sp -= 8
		} else if m == 15 { //drop
			dx(a)
			a = I(sp)
			sp -= 4
		}
		xp += 4
		yp += 4

		if sp < s+4 {
			fmt.Printf("stack underflow\n")
			panic("!stack")
		} else if sp > se {
			fmt.Printf("stack overflow: %d > %d\n", sp, se)
			panic("!stack")
		}
	}
	if sp != s+8 {
		fmt.Println("unbalanced stack", sp-s-8)
		panic("!stack")
	}
	dx(x)
	dx(y)
	dx(z)
	dx(s)
	sI(140, ss)
	return a
}
func swc(x, y i) (r i) { // $[a;b;..]
	_, _, xn, _, xp, yp := v2(x, y)
	if 0 == xn%2 {
		fmt.Println("$[even]")
		panic("cond")
	}
	rl(x)
	rl(y)
	s := mk(2, 0)
	r = mk(2, 0)
	a, b, t, w := i(0), i(0), i(0), i(0)
	a -= 4
	for i := i(0); i < xn; i++ {
		t = I(xp)
		w = I(yp)
		b = 4 * nn(t)
		a += 4 + b
		if 0 == i%2 {
			b = op(9, a)
		} else {
			b = op(8, 4+b)
		}
		r = ucat(ucat(mki(b), t), r)
		s = ucat(ucat(mki(0), w), s)
		xp += 4
		yp += 4
	}
	dx(x)
	dx(y)
	return l2(drop(r, 1), drop(s, 1))
}
func op(x, y i) (r i) { return x | y<<4 }
func prj(x, y, z i) (r i) {
	r = mk(0, 3)
	sI(r+8, x)
	sI(r+12, y)
	sI(r+16, z)
	return r
}
func prs(x i) (r i) {
	// fmt.Printf("prs %s\n", X(x))
	rx(x)
	r = tok(x)
	rld(r)
	a := I(r + 12) //pos
	r = I(r + 8)   //tokens
	sI(152, x)     //src string
	sI(144, a-r)   //offset

	rp := r + 8
	rn := nn(r)
	k, v := kvd(sq(rp, rp+4*rn))
	k = jon(k, mki(15))
	v = jon(v, mki(0))
	dx(a)
	dx(r)
	return l3(k, v, x)
}
func sq(x, y i) (r i) {
	a := mk(6, 0)
	s := mk(6, 0)
	n := 0
	for {
		k, v := kvd(ex(pt(x, y), y))
		if k < 128 {
			dx(v)
			if k == ';' || nn(a) > 0 {
				a = lcat(a, mki(4))
				s = lcat(s, mki(0))
			}
			if k != ';' {
				return l2(a, s)
			}
		} else if k > 128 {
			a = lcat(a, k)
			s = lcat(s, v)
			x = I(pp)
			if I(x) != ';' {
				return l2(a, s)
			}
		}
		x += 4
		n++
	}
}
func fnl(x i) (r i) {
	rx(x)
	r = wer(ovr(ecr(mki(4), x, 126), 44)) // &,/4~/:x
	if nn(r) == 0 {
		dx(r)
		return 0
	}
	return r
}
func vb(x i) (r i) {
	n := nn(x)
	x += 8
	r = I(x)
	if n == 1 && r&0x0f == 4 {
		return 4 //verb
	}
	if n == 1 && r&0x0f == 0 && tp(r) == 0 && nn(r) == 2 {
		return 3 //applied adverb
	}
	if I(x+4*(n-1))&0x0f == 1 {
		return 1 //adverb
	}
	return 0 //noun
}
func quo(x i) (r i) {
	r = mk(5, 1)
	sI(r+8, x>>4)
	return mki(con(r))
}
func ras(x, y, z, u, v, w i) (r i) { // assign
	t := i(0)
	if vb(z) == 4 {
		p := I(z + 8)
		if p == 932 || p > 2048 {
			a := i(0)
			n := nn(x)
			s := I(x + 4*n)
			if n == 1 && I(x+8)&0xf == 5 { // x:y x+:y
				dx(x)
				x = quo(I(x + 8))
				if p == 932 {
					dx(z)
					dx(w)
					return l2(ucat(ucat(y, x), mki(2+':'<<4)), ucat(ucat(v, u), mki(0))) // y x :
				}
				a = mki(4)
				t = mki(0)
			} else if n > 2 && s&0xf == 5 {
				a = take(x, n-2)
				t = take(u, n-2)
				x = quo(s)
				u = mki(0)
			} else {
				return 0
			}
			if p > 2048 {
				p -= 2048
			}
			dx(z)
			dx(w)
			r = ucat(ucat(ucat(ucat(y, a), x), mki(p)), mki(7))
			t = ucat(ucat(ucat(ucat(v, t), u), mki(0)), mki(0))
			return l2(r, t)
		}
	}
	return 0
}
func src(x i) (r i) { return mki(I(x+I(144)) - I(152)) }
func ex(x, y i) (r i) {
	x, s := kvd(x)
	if x < 128 {
		return l2(x, s)
	}
	p := vb(x)
	t := i(0)
	r, t = kvd(pt(I(pp), y))
	if r < 128 {
		dx(t)
		return l2(x, s)
	}
	q := vb(r)
	w := i(0)
	if q != 0 && p == 0 {
		y, w = kvd(ex(pt(I(pp), y), y))
		v := i(2 + '.'<<4)
		if y < 256 {
			y = mki(4)
			r = ucat(mki(con(mki(1))), r)
			t = ucat(mki(0), t)
			v = 10
		}
		p = ras(x, y, r, s, w, t)
		if p != 0 {
			return p
		}
		if q == 4 && v != 10 {
			dx(r)
			r = mki(I(r+8) - 2)
			r = ucat(ucat(y, x), r)
			return l2(r, ucat(ucat(w, s), t)) // x y +
		}
		r = ucat(ucat(ucat(ucat(y, x), mki(6+2<<4)), r), mki(v)) // x y l2 r .
		t = ucat(ucat(ucat(ucat(w, s), mki(0)), t), mki(0))
		return l2(r, t)
	}
	r, t = kvd(ex(l2(r, t), y))
	if p != 4 {
		x = ucat(x, mki(1026))
		s = ucat(s, mki(0))
	} else {
		p := r + 4*(1+nn(r))
		v := I(8 + x)
		dx(x)
		if v == 932 { // :x
			x = mki(4294967289) //jump far (return)
		} else if v == 676 && I(p) == 4033 { // *|
			dx(s)
			sI(p, 2977)
			return l2(r, t)
		} else {
			x = mki(I(x+8) + 2045) //monadic
		}
	}
	return l2(ucat(r, x), ucat(t, s))
}
func pt(x, y i) (r i) {
	if x >= y {
		sI(pp, x)
		return l2(0, mki(0))
	}
	r = I(x)
	t := i(0)
	if r == '(' {
		r, t = kvd(sq(x+4, y))
		n := nn(r)
		if n == 0 {
			dx(t)
			r = mki(con(r)) // ()
			t = src(x)
		} else if n == 1 {
			r = ic(fst(r), 0)
			t = ic(fst(t), 0)
			//r = ic(ic(fst(r), 4), 15) //not infix
			//t = ic(ic(fst(t), 0), 0)
		} else {
			r = ic(jon(rev(r), mk(2, 0)), op(6, n))
			t = ic(jon(rev(t), mk(2, 0)), 0)
		}
		x = I(pp)
	} else if r == '{' { // {
		r = lam(x, y)
		t = src(x)
		x = I(pp)
	} else if r < 128 && r != '[' {
		sI(pp, x)
		return l2(r, mki(0))
	} else {
		r = mki(r)
		t = src(x)
	}
	x += 4
	for {
		p := I(x)
		if x < y && p == '[' { // [
			a, w := kvd(sq(x+4, y))
			n := nn(a)
			if n == 1 {
				a = fst(a)
				w = fst(w)
				r = ucat(a, r)
				t = ucat(w, t)
				r = ucat(r, mki(2+'@'<<4))
				t = ucat(t, mki(0))
			} else if n > 2 && I(r+8) == 580 { // $[..] cond
				dx(r)
				dx(t)
				r, t = kvd(swc(rev(a), rev(w)))
			} else if n == 3 && I(r+8) == 1028 { // @[x;y;z]
				dx(r)
				r = ic(jon(rev(a), mk(2, 0)), 3)
				t = ucat(jon(rev(w), mk(2, 0)), t)
			} else {
				l := i(0)
				if n == 0 {
					dx(a)
					a = mki(con(mk(6, 0)))
					dx(w)
					w = mki(0)
				} else {
					l = fnl(a)
					a = ucat(jon(rev(a), mk(2, 0)), mki(op(6, n)))
					w = ucat(jon(rev(w), mk(2, 0)), mki(0))
					if l != 0 {
						a = ucat(a, mki(con(l)))
						w = ucat(w, mki(0))
						l = 10
					}
				}
				if l == 0 {
					l = 2 + '.'<<4
				}
				r = ucat(ucat(a, r), mki(l))
				t = ucat(ucat(w, t), mki(0))
			}
			x = 4 + I(pp)
		} else if x < y && p&0xf == 1 { // adverb
			s := src(x)
			v := vb(r)
			if v > 2 {
				dx(r)
				r = I(r + 8)
				if v == 4 {
					r >>= 4
				}
				rx(r)
				r = mki(con(atx(p>>4, r))) // apply adv
				t = ucat(take(t, nn(t)-1), s)
			} else {
				r = ucat(r, mki(p))
				t = ucat(t, s)
			}
			x += 4
		} else {
			sI(pp, x)
			return l2(r, t)
		}
	}
}
func lam(x, y i) (r i) {
	ds := I(x+I(144)) - I(152) - 8
	s0 := I(x + I(144))
	a := i(0)
	if I(x+4) == '[' { // {[args
		a = fst(sq(x+8, y))
		n := nn(a)
		if n == 1 {
			a = fst(a)
		} else {
			a = jon(a, mk(2, 0))
		}
		a = diw(a, mki(16))
		sI(a, I(a)+3<<29) // 5<<29
		x = I(pp)
	}
	t := i(0)
	r, t = kvd(sq(x+4, y)) //tree
	n := nn(r)
	if n == 1 {
		r = fst(r)
		t = fst(t)
	} else {
		r = jon(r, mki(15))
		t = jon(t, mki(0))
	}
	p := t + 8
	for i := i(0); i < nn(t); i++ {
		if I(p) != 0 {
			sI(p, I(p)-ds)
		}
		p += 4
	}
	s1 := 1 + I(I(pp)+I(144))
	s := mk(1, s1-s0)
	mv(s+8, s0, s1-s0)
	if a != 0 {
		x = a
		a = nn(x)
	} else { // detect xyz
		n = nn(r)
		x = 8 + mki(197)
		for i := i(0); i < 3; i++ {
			if fnx(r, x) < n {
				a = 1 + i
			}
			sI(x, I(x)+64)
		}
		dx(x - 8)
		x = I(xyz)
		rx(x)
		x = take(x, a)
	}
	rp := r + 8
	u := i(0)
	for i := i(0); i < nn(r); i++ {
		if I(rp) == 930 && u&0xf == 0 && tp(u) == 5 {
			rx(u)
			x = ucat(x, u) // locals
		}
		u = I(rp)
		rp += 4
	}
	x = unq(x)
	f := mk(0, 5)
	sI(f+8, s)  // string
	sI(f+12, r) // tree
	sI(f+16, x) // args
	sI(f+20, a) // arity
	sI(f+24, t) // src
	return mki(con(f))
}
func cmt(p, s i) (r i) {
	for p < s {
		if C(p) == 10 {
			return p
		}
		p++
	}
	return p
}
func tok(x i) (r i) { // tokenize
	xt, xn, xp := v1(x)
	if xt != 1 {
		trap()
	}
	r = mk(2, 0)
	a := mk(2, 0)
	xn += xp
	for xp < xn && C(xp) == 47 {
		xp = cmt(xp, xn)
	}
	sI(pp, xp)
	for {
		xp = ws(I(pp), xn)
		if xp >= xn {
			dx(x)
			return l2(r, a)
		}
		t := ntk(xp, xn)
		if t == 0 {
			dx(x)
			return l2(r, a)
		}
		r = ucat(r, mki(t))
		a = ucat(a, mki(xp))
	}
}
func ws(x, y i) (r i) {
	for {
		if x >= y {
			return x
		}
		if C(x) == 47 {
			c := C(x - 1)
			if c == 32 || c == 10 {
				x = cmt(x, y)
			}
		}
		b := C(x)
		if is(b, NW) || b == 10 { //64
			return x
		}
		x++
		//if C(x) == 47 { // /
		//	x = cmt(x, y)
		//	fmt.Printf("x=%d C(x)=%d\n", x, C(x))
		//}
	}
}
func ntk(x, y i) (r i) { // next token
	b := C(x)
	for j := 0; j < 5; j++ {
		r = MT[j+136].(func(c, i, i) i)(b, x, y)
		if r != 0 {
			if r&0xf == 0 {
				if I(r) == 1+5<<29 {
					dx(r)
					r = op(5, I(r+8))
				} else if I(r) == 1+6<<29 {
					r = con(fst(r))
				} else {
					r = con(r)
				}
			}
			return r
		}
	}
	sI(pp, 1+x)
	if b == 10 {
		return ';'
	}
	return i(b)
}
func pui(b c, p, s i) (r uint32) {
	if !is(b, NM) { // 4
		return 0
	}
	for is(b, NM) && p < s {
		r *= 10
		r += uint32(b) - '0' //48
		p++
		b = C(p)
	}
	if r == 0 && b == 'x' && p < s { // 120 (0x)
		panic("pui 0x ? chr should be first")
		return 0
	}
	sI(pp, p)
	return r
}
func pin(b c, p, s i) (r i) { // parse signed int
	sI(pp, p)
	u := pui(b, p, s)
	if I(pp) != p {
		return mki(i(u))
	}
	if b == '-' { //45
		p++
		if p < s {
			b = C(p)
			sI(pp, p)
			u = pui(b, p, s)
			if I(pp) != p {
				return mki(-i(u))
			}
			sI(pp, p-1)
			return 0
		}
	}
	return 0
}
func pfd(p, s i, f float64) float64 {
	g := float64(1)
	if f < 0 {
		g = -g
	}
	for {
		b := C(p)
		if p < s && is(b, NM) { //4
			g *= 0.1
			f += g * float64(b-'0') //48
		} else {
			sI(pp, p)
			return f
		}
		p++
	}
}
func pfl(b c, p, s i) (r i) { // parse float
	m := i(0)
	if b == '-' { //45
		t := C(p - 1)
		if is(t, az|AZ|NM) || t == ')' || t == ']' || t == '"' { // C(p-1) maybe refcount's last byte (k2 reflite p27)
			return 0
		}
		m = 1
	}
	r = pin(b, p, s) //overflows uint32
	p = I(pp)
	if r == 0 || p == s {
		return r
	}
	if C(p) == '.' { //46
		r = up(r, 2, 1)
		rp := r + 8
		sF(rp, pfd(p+1, s, F(rp)))
	}
	p = I(pp)
	if p < s && C(p) == 'e' { //101
		sI(pp, p+1)
		q := pin(C(p+1), p+1, s)
		if q == 0 {
			sI(pp, p)
			return r
		}
		e := I(q + 8)
		dx(q)
		f := F(r + 8)
		for int32(e) < 0 {
			f /= 10.0
			e++
		}
		for e > 0 {
			f *= 10.0
			e--
		}
		sF(r+8, f)
	}
	if m == 1 { // fix lost -0.
		f := F(r + 8)
		if f > 0 {
			sF(r+8, -f)
		}
	}
	return r
}
func num(b c, p, s i) (r i) { // parse single number
	r = pfl(b, p, s)
	if r == 0 {
		return r
	}
	p = I(pp)
	b = C(p)
	if p < s {
		if b == 'a' || b == 'p' || b == 'n' || b == 'w' { //97 112 110 119
			if tp(r) == 2 {
				r = up(r, 2, 1)
			}
			p++
			sI(pp, p)
			if b != 'a' {
				var f float64
				if b == 'p' {
					f = math.Pi * F(r+8)
				}
				if b == 'n' {
					f = math.NaN() // 0x7FF8000000000001 (see also go #44016)
					// use 0n@0n for complex
				}
				if b == 'w' {
					f = math.Inf(1) // todo -0w
				}
				sF(r+8, f)
				return r
			}

			r = up(r, 3, 1)
			a := pfl(C(p), p, s)
			if a == 0 {
				a = mki(0)
			}
			if tp(a) == 2 {
				a = up(a, 2, 1)
			}
			r = atx(r, a)
		}
	}
	return r
}
func nms(b c, p, s i) (r i) { // parse numeric vector
	r = num(b, p, s)
	if r == 0 {
		return 0
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
		q := num(C(p), p, s)
		if q == 0 {
			sI(pp, p-1) // keep space for " /comment" (todo?)
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
	if C(p-1) == 32 {
		if b == 92 { // space\.. (out)
			sI(pp, p+1)
			return op(4, 32)
		}
		if b == 39 { // (space)'c  spacy verb
			b = '+' //no adverb
			p++
		}
	}
	r = i(C(p))
	//n := i(2) - boolvar(is(byte(r), AD))
	if s > p+1 {
		if C(p+1) == ':' { // 58
			p++
			r += 128
		}
	}
	sI(pp, p+1)
	if is(b, AD) {
		return op(1, r)
	}
	return op(4, r)
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
	if b != '`' || p == s { //96
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
	sI(r+8, 8)
	return r
}
func sms(b c, p, s i) (r i) { // `a`b
	r = sym(b, p, s)
	if r == 0 {
		return 0
	}
	for {
		p = I(pp)
		q := sym(C(p), p, s)
		if q == 0 {
			return enl(r)
		}
		r = ucat(r, q)
	}
	return 0
}
func chr(b c, p, s i) (r i) { // "abc"
	if b == '0' {
		if 1 < s-p {
			if C(p+1) == 'x' {
				return phx(2+p, s)
			}
		}
	}
	if b != '"' { //34
		return 0
	}
	a := p + 1
	for {
		p++
		if p == s {
			panic("chr/eof")
		}
		if C(p) == '"' {
			n := p - a
			r = mk(1, n)
			mv(r+8, a, n)
			sI(pp, p+1)
			return r
		}
	}
}
func phx(p, s i) (r i) { // 0xab12
	h, q := true, i(0)
	r = mk(1, 0)
	for {
		c := C(p)
		if s <= p || !is(c, NM+az) {
			sI(pp, p)
			return r
		}
		if c < 58 { //':'
			c -= 48 //'0'
		}
		if c > 96 {
			c -= 87
		}
		h = !h
		if h {
			r = cc(r, i(c)+(q<<4))
		}
		q = i(c)
		p++
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
func Z(a i) complex128     { return complex(MF[a>>3], MF[1+(a>>3)]) }
func sC(a i, v c)          { MC[a] = v } // global set, e.g. i::v
func sI(a i, v i)          { MI[a>>2] = v }
func sJ(a i, v j)          { MJ[a>>3] = v }
func sF(a i, v f)          { MF[a>>3] = v }
func sZ(a i, v complex128) { MF[a>>3] = real(v); MF[1+(a>>3)] = imag(v) }

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

var leaktest = false

var leak = _leak

func _leak() {
	dx(I(kkey))
	dx(I(kval))
	dx(I(xyz))
	leaktest = true
	dx(I(kcon))
	leaktest = false
	mark()
	p := i(64)
	for p < i(len(MI)) {
		if MI[p+1] != 0 {
			panic(fmt.Errorf("non-free block at %d(%x) %d/%d", p<<2, p<<2, tp(p<<2), nn(p<<2)))
		}
		t := MI[p+2]
		if t < 4 || t > 31 {
			panic(fmt.Errorf("illegal bucket type %d at %d(%x)", t, p<<2, p<<2))
		}
		dp := i(1) << t
		p += dp >> 2
	}
}
func O(s s, x i) i { fmt.Printf("%s: %x(%d/%d) %s\n", s, x, tp(x), nn(x), X(x)); return x }
func X(x i) s {
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
	ystr := func(i i) s {
		n := I(x + 8 + 4*i)
		s := X(I(I(kkey) + n))
		return s[1 : len(s)-1]
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
			return X(MI[(x+12)>>2]) + X(MI[(x+8)>>2])
		} else if n == 3 {
			r := X(I(x + 12))
			return X(I(x+8)) + "[" + r[1:len(r)-1] + "]"
		} else if n == 5 { // lambda
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
		f = ystr
		sep = "`"
		tof = func(s s) s { return "`" + s }
	case 6:
		if n == 1 {
			return "," + X(I(8+x))
		}
		f = func(i i) s { return X(MI[2+i+x>>2]) }
		sep = ";"
		tof = func(s s) s { return "(" + s + ")" }
	case 7:
		return X(I(x+8)) + "!" + X(I(x+12))
	default:
		panic(fmt.Sprintf("nyi: t=%d", t))
	}
	r := make([]s, n)
	for k := range r {
		r[k] = f(i(k))
	}
	return tof(strings.Join(r, sep))
}
