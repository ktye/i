// +build ignore

// reference implementation for k.w
// go run k.go t
// go run k.go 5 mki til rev

package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"math/bits"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

const trace = false

type c = byte
type s = string
type i = uint32
type j = uint64
type f = float64

var MC []c // MC, MI, MJ, MF share array (see msl)
var MI []i
var MJ []j
var MF []f
var WT []i

type vt1 func(i) i
type vt2 func(i, i) i
type slice struct {
	p uintptr
	l int
	c int
}

const naI i = 2147483648
const naJ j = 9221120237041090561

func main() {
	if len(os.Args) == 2 && os.Args[1] == "t" {
		runtest()
	} else {
		fmt.Println(run(os.Args[1:]))
	}
}
func runtest() {
	b, e := ioutil.ReadFile("t")
	if e != nil {
		panic(e)
	}
	v := strings.Split(strings.TrimSpace(string(b)), "\n")
	for i := range v {
		if len(v[i]) == 0 || v[i][0] == '/' {
			continue
		}
		vv := strings.Split(v[i], " /")
		if len(vv) != 2 {
			panic("test file")
		}
		in := strings.TrimSpace(vv[0])
		exp := strings.TrimSpace(vv[1])
		got := run(strings.Fields(in))
		fmt.Println(in, "/", got)
		if exp != got {
			fmt.Printf("!")
			os.Exit(1)
		}
	}
}
func parseVector(s string) i {
	if len(s) > 0 && s[0] == '"' { // "char"
		s = strings.Trim(s, `"`)
		b := []c(s)
		r := mk(1, i(len(b)))
		for i := 0; i < len(b); i++ {
			MC[8+int(r)+i] = b[i]
		}
		return r
	}
	if len(s) > 0 && s[0] == '`' { // `symbols`b`c
		v := strings.Split(s[1:], "`")
		sn := i(len(v))
		sv := mk(4, sn)
		for i := i(0); i < sn; i++ {
			b := v[i]
			rn := uint32(len(b))
			r := mk(1, rn)
			for k := uint32(0); k < rn; k++ {
				MC[8+r+k] = b[k]
			}
			sI(sv+8+4*i, r)
		}
		return sv
	}
	if len(s) > 0 && s[0] == '(' { // (`list;1;2)
		return parseList(s[1:])
	}
	f := strings.Index(s, ".") // 1.23,2.34 (float)
	v := strings.Split(s, ",") // 1,2,3 (int vector)
	n := uint32(len(v))
	iv := make([]int64, n)
	fv := make([]float64, n)
	var e error
	for i, s := range v {
		if f == -1 {
			iv[i], e = strconv.ParseInt(s, 10, 32)
		} else {
			fv[i], e = strconv.ParseFloat(s, 64)
		}
		if e != nil {
			panic("parse number: " + s)
		}
	}
	if f == -1 {
		x := mk(2, n)
		for i := uint32(0); i < n; i++ {
			MI[(x+8+i*4)>>2] = uint32(iv[i])
		}
		return x
	} else {
		x := mk(3, n)
		for i := uint32(0); i < n; i++ {
			MF[(x+8+i*8)>>3] = fv[i]
		}
		return x
	}
}
func parseList(s string) i {
	if len(s) == 0 || s[len(s)-1] != ')' {
		panic("parse list")
	} else if len(s) == 1 {
		return mk(5, 0)
	}
	r := make([]i, 0)
	s = s[:len(s)-1]
	l, a := 0, 0
	for i, c := range s {
		if c == '(' {
			l++
		} else if c == ')' {
			l--
			if l < 0 {
				panic(")")
			}
		} else if l == 0 && c == ';' {
			r = append(r, parseVector(s[a:i]))
			a = i + 1
		}
	}
	r = append(r, parseVector(s[a:]))
	x := mk(5, i(len(r)))
	for k := range r {
		MI[2+(x>>2)+i(k)] = r[k]
	}
	return x
}
func run(args []string) string {
	m0 := 16
	fn1 := map[string]vt1{"til": til, "rev": rev, "fst": fst, "enl": enl, "cnt": cnt, "tip": tip, "wer": wer, "not": not}
	fn2 := map[string]vt2{"mk": mk, "atx": atx, "cut": cut, "rsh": rsh, "cat": cat, "eql": eql}
	stack := make([]i, 0)
	MJ = make([]j, (1<<m0)>>3)
	msl()
	ini(16)
	for _, a := range args {
		if f1, o := fn1[a]; o {
			r := f1(stack[len(stack)-1])
			if trace {
				fmt.Printf("%s %d: x%x\n", a, stack[len(stack)-1], r)
			}
			stack[len(stack)-1] = r
			continue
		}
		if f2, o := fn2[a]; o {
			x, y := stack[len(stack)-2], stack[len(stack)-1]
			r := f2(x, y)
			if trace {
				fmt.Printf("%s %d %d: x%x\n", a, x, y, r)
			}
			stack = stack[:len(stack)-1]
			stack[len(stack)-1] = r
			continue
		}
		if strings.HasPrefix(a, "dump") {
			a = strings.TrimPrefix(a, "dump")
			if n, e := strconv.Atoi(a); e == nil {
				dump(0, uint32(n))
			} else {
				dump(0, 100)
			}
			continue
		}
		stack = append(stack, parseVector(a))
	}
	if len(stack) != 1 {
		panic("stack #" + strconv.Itoa(len(stack)))
	}
	r := kst(stack[len(stack)-1])
	dx(stack[len(stack)-1])
	leak()
	return r
}
func ini(x i) i {
	sJ(0, 289360691419414784) // uint64(0x0404040408040100)
	sI(128, x)
	p := i(256)
	for i := i(8); i < x; i++ {
		sI(4*i, p)
		p *= 2
	}
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
	// todo Z
}
func bk(t, n i) (r i) {
	r = i(32 - bits.LeadingZeros32(7+n*i(C(t))))
	if r < 4 {
		return 4
	}
	return r
}
func mk(x, y i) i {
	t := bk(x, y)
	i := 4 * t
	for I(i) == 0 {
		i += 4
	}
	if i == 128 {
		panic("oom")
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
			if xt > 3 {
				for i := i(0); i < xn; i++ {
					dx(I(xp + 4*i))
				}
			}
			fr(x)
		}
	}
}
func rx(x i) {
	if x > 255 {
		MI[1+x>>2]++
	}
}
func rl(x i) {
	_, xn, xp := v1(x)
	for i := i(0); i < xn; i++ {
		rx(I(xp))
		xp += 4
	}
}
func dxr(x, r i) i     { dx(x); return r }
func dxyr(x, y, r i) i { dx(x); dx(y); return r }
func mki(i i) (r i)    { r = mk(2, 1); sI(r+8, i); return r }
func mkd(x, y i) (r i) {
	xt, _, xn, yn, _, _ := v2(x, y)
	if xt != 5 {
		panic("type")
	} else if xn != yn {
		panic("length")
	}
	r = mk(7, 2)
	MI[2+r>>2] = x
	MI[3+r>>2] = y
	return r
}
func v1(x i) (xt, xn, xp i) { u := I(x); return u >> 29, u & 536870911, 8 + x }
func v2(x, y i) (xt, yt, xn, yn, xp, yp i) {
	xt, xn, xp = v1(x)
	yt, yn, yp = v1(y)
	return
}
func use(x i) (r i) {
	xt, xn, xp := v1(x)
	if I(x+4) != 1 {
		r = mk(xt, xn)
		mv(r+8, xp, i(C(xt)))
		dx(x)
		x = r
	}
	return x
}
func mv(dst, src, n i) { copy(MC[dst:dst+n], MC[src:src+n]) }
func ext(x, y i) (rx, ry i) {
	xt, yt, xn, yn, _, _ := v2(x, y)
	if xt != yt {
		trap()
	}
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
	_, n, _ := v1(x)
	if n == 0 {
		return x
	}
	return atx(x, tir(n))
}
func fst(x i) (r i) {
	xt, _, _ := v1(x)
	if xt == 7 {
		rx(12 + x)
		dx(x)
		return fst(12 + x)
	}
	return atx(x, mki(0))
}
func drop(x, n i) (r i) {
	_, xn, _ := v1(x)
	if n > xn {
		n = xn
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
	r = mk(5, xn)
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
	_, xn, _ := v1(x)
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
func atx(x, y i) (r i) {
	xt, yt, xn, yn, xp, yp := v2(x, y)
	if xt == 7 {
		panic("nyi atx d")
	}
	if yt != 2 {
		panic("atx yt~I")
	}
	r = mk(xt, yn)
	rp := r + 8
	switch xt {
	case 1: // yn/((rp+i)::C?32;yi:I yp;(yi<xn)?(rp+i)::C xp+yi;yp+:4)
		for i := i(0); i < yn; i++ {
			sC(rp+i, 32)
			yi := I(yp)
			if yi < xn {
				sC(rp+i, C(xp+yi))
			}
			yp += 4
		}
	case 2: // yn/(rp::naI;yi:I yp;(yi<xn)?rp::I xp+4*yi;rp+:4;yp+:4)
		for i := i(0); i < yn; i++ {
			sI(rp, naI)
			yi := I(yp)
			if yi < xn {
				sI(rp, I(xp+4*yi))
			}
			rp += 4
			yp += 4
		}
	case 3: // yn/(rp::naF;yi:I yp;(yi<xn)?rp::F xp+8*yi;rp+:8;yp+:4)
		naF := math.Float64frombits(naJ)
		for i := i(0); i < yn; i++ {
			sF(rp, naF)
			yi := I(yp)
			if yi < xn {
				sF(rp, F(xp+8*yi))
			}
			rp += 8
			yp += 4
		}
	case 4, 5:
		naS := mk(1, 0)
		for i := i(0); i < yn; i++ {
			sI(rp, naS)
			yi := I(yp)
			if yi < xn {
				sI(rp, I(xp+4*yi))
			}
			rp += 4
			yp += 4
		}
		rl(r)
		dx(naS)
	default:
		panic(fmt.Sprintf("nyi atx xt=%d", xt))
	}
	return dxyr(x, y, r)
}
func cat(x, y i) (r i) {
	xt, yt, _, _, _, _ := v2(x, y)
	if xt == 0 || yt == 0 {
		trap()
	}
	if xt == yt {
		return ucat(x, y)
	}
	if xt == 5 {
		return lcat(x, y)
	}
	panic("nyi cat")
}
func ucat(x, y i) (r i) {
	xt, _, xn, yn, xp, yp := v2(x, y)
	if xt > 4 {
		rl(x)
	}
	if xt > 5 {
		r = mkd(x+8, x+12)
		return dxyr(x, y, r)
	}
	r = mk(xt, xn+yn)
	w := i(C(xt))
	mv(r+8, xp, w*xn)
	mv(r+8+w*xn, yp, w*yn)
	return dxyr(x, y, r)
}
func lcat(x, y i) (r i) { // list append
	x = use(x)
	xt, xn, xp := v1(x)
	if bk(xt, xn) < bk(xt, xn+1) {
		r = mk(xt, xn+1)
		mv(r+8, xp, 4*xn)
		dx(x)
		x, xp = r, r+8
	}
	sI(x, (xn+1)|5<<29)
	sI(xp+4*xn, y)
	return x
}
func enl(x i) (r i) { return lcat(mk(5, 0), x) }
func cnt(x i) (r i) {
	_, xn, _ := v1(x)
	dx(x)
	return mki(xn)
}
func tip(x i) (r i) {
	xt, _, _ := v1(x)
	r = mk(2, 1)
	sI(8+r, xt)
	return dxr(x, r)
}
func wer(x i) (r i) {
	xt, xn, xp := v1(x)
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
func not(x i) (r i) { return eql(mki(0), x) }
func eql(x, y i) (r i) {
	x, y = ext(x, y)
	xt, _, xn, _, xp, yp := v2(x, y)

	switch xt {
	case 0, 2, 4:
		r = eqI(xp, yp, xn)
	case 1:
		r = eqC(xp, yp, xn)
	case 3:
		r = eqF(xp, yp, xn)
	default:
		panic("nyi")
	}
	return dxyr(x, y, r)
}
func eqC(xp, yp, n i) (r i) {
	r = mk(2, n)
	rp := r + 8
	for i := i(0); i < n; i++ {
		sI(rp+i, boolvar(C(xp+i) == C(yp+i)))
	}
	return r
}
func eqI(xp, yp, n i) (r i) {
	r = mk(2, n)
	rp := r + 8
	for i := i(0); i < n; i++ {
		sI(rp, boolvar(I(xp) == I(yp)))
		rp += 4
		xp += 4
		yp += 4
	}
	return r
}
func eqF(xp, yp, n i) (r i) {
	r = mk(2, n)
	rp := r + 8
	for i := i(0); i < n; i++ {
		sI(rp, boolvar(F(xp) == F(yp)))
		rp += 4
		xp += 8
		yp += 8
	}
	return r
}

func boolvar(b bool) i {
	if b {
		return 1
	}
	return 0
}
func trap() { panic("trap") }
func dump(a, n i) i {
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
func C(a i) c     { return MC[a] } // global get, e.g. I i
func I(a i) i     { return MI[a>>2] }
func J(a i) j     { return MJ[a>>3] }
func F(a i) f     { return MF[a>>3] }
func sC(a i, v c) { MC[a] = v } // global set, e.g. i::v
func sI(a i, v i) { MI[a>>2] = v }
func sJ(a i, v j) { MJ[a>>3] = v }
func sF(a i, v f) { MF[a>>3] = v }
func atoi(s string) i {
	if x, e := strconv.Atoi(s); e == nil {
		return i(x)
	}
	panic("atoi")
}
func mark() { // mark bucket type within free blocks
	for t := i(4); t < 32; t++ {
		for p := MI[t] >> 2; p != 0; p = MI[p] >> 2 {
			MI[1+p] = 0
			MI[2+p] = t
		}
	}
}
func leak() {
	mark()
	//dump(0, 200)
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
	a := MI[x>>2]
	t, n := a>>29, a&536870911
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
	sstr := func(i i) s {
		r := I(x + 8 + 4*i)
		rn := I(r) & 536870911
		return string(MC[r+8 : r+8+rn])
	}
	sep := " "
	switch t {
	case 1:
		return `"` + string(MC[x+8:x+8+n]) + `"`
	case 2:
		f = istr
	case 3:
		f = fstr
		tof = func(s s) s {
			if strings.Index(s, ".") == -1 {
				return s + "f"
			}
			return s
		}
	case 4:
		f = sstr
		sep = "`"
		tof = func(s s) s { return "`" + s }
	case 5:
		f = func(i i) s { return kst(MI[2+i+x>>2]) }
		sep = ";"
		tof = func(s s) s { return "(" + s + ")" }
	default:
		panic(fmt.Sprintf("nyi: kst: t=%d", t))
	}
	r := make([]s, n)
	for k := range r {
		r[k] = f(i(k))
	}
	return tof(strings.Join(r, sep))
}
