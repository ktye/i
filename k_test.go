package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"

	. "github.com/ktye/wg/module"
	"github.com/ktye/wg/wasi_unstable"
)

var save []byte

func newtest() {
	wasi_unstable.Stdout = os.Stdout
	rand = 1592653589
	if save == nil {
		kinit()
		save = make([]byte, len(Bytes))
		copy(save, Bytes)
	} else {
		Bytes = make([]byte, len(save))
		copy(Bytes, save)
		src, pp, pe, sp = 0, 0, 0, 256
	}
}
func mkchars(b []byte) (r K) {
	r = mk(Ct, int32(len(b)))
	copy(Bytes[int32(r):], b)
	return r
}
func intvalue(x K) int32     { return int32(x) }
func floatvalue(x K) float64 { return F64(int32(x)) }
func TestTypes(t *testing.T) {
	t.Skip()
	newtest()
	xi := Ki(-5)
	if v := intvalue(xi); v != -5 {
		t.Fatalf("got v=%d expected %d\n", v, -5)
	}
	if r := tp(xi); r != it {
		t.Fatalf("got t=%d expected %d\n", r, it)
	}
	if n := intvalue(Cnt(xi)); n != 1 {
		t.Fatalf("got n=%d expected %d\n", n, 1)
	}
	xf := Kf(math.Pi)
	if f := floatvalue(xf); f != math.Pi {
		t.Fatalf("got f=%v expected %v\n", f, math.Pi)
	}
	if r := tp(xf); r != ft {
		t.Fatalf("got t=%d expected %d\n", r, ft)
	}
	dx(xf)
	reset()
}
func TestBucket(t *testing.T) {
	tc := []struct{ in, exp int32 }{
		{0, 5},
		{4, 5},
		{8, 5},
		{16, 5},
		{17, 6},
		{25, 6},
	}
	for _, tc := range tc {
		if got := bucket(tc.in); got != tc.exp {
			t.Fatalf("bucket %d => %d (exp %d)\n", tc.in, got, tc.exp)
		}
	}
}
func TestMk(t *testing.T) {
	t.Skip()
	newtest()
	r := mk(It, 2)
	if rc := I32(int32(r) - 4); rc != 1 {
		t.Fatalf("rc is %d not 1\n", rc)
	}
	if n := I32(int32(r) - 12); n != 2 {
		t.Fatalf("n is %d not 2\n", n)
	}
	if tx := tp(r); tx != It {
		t.Fatalf("t is %d not %d\n", tx, It)
	}
	if n := nn(r); n != 2 {
		t.Fatalf("nn(x) is %d not 2\n", n)
	}
}
func TestFloat(t *testing.T) {
	t.Skip()
	newtest()
	dx(val(Ku(35382781554225)))
	reset()
}
func TestVerbs(t *testing.T) {
	t.Skip()
	newtest()
	x := Til(Ki(3))
	if r := iK(Cnt(x)); r != 3 {
		t.Fatalf("got %d expected %d", r, 3)
	}
}
func TestTok(t *testing.T) {
	t.Skip()
	tc := []struct {
		in, exp string
	}{
		{"1234567", ",1234567"},
		{"-1234567", ",-1234567"},
		{"*", ",*"},
	}
	for _, tc := range tc {
		newtest()
		//fmt.Println(tc.in)
		got := sK(tok(mkchars([]byte(tc.in))))
		if got != tc.exp {
			t.Fatalf("got %s expected %s", got, tc.exp)
		}
	}
}
func TestMultiline(t *testing.T) {
	t.Skip()
	tc := []struct {
		in, exp string
	}{
		{"1+2\n3*4", "12"},
	}
	for _, tc := range tc {
		newtest()
		//fmt.Println(tc.in)
		got := sK(Val(mkchars([]byte(tc.in))))
		if got != tc.exp {
			t.Fatalf("got %s expected %s", got, tc.exp)
		}
	}
}
func TestShuffle(t *testing.T) {
	t.Skip()
	newtest()
	shuffle(seq(8), 5)
	reset()
}
func TestKT(t *testing.T) {
	//t.Skip()
	newtest()
	b, err := ioutil.ReadFile("k.t")
	if err != nil {
		t.Fatal(err)
	}
	v := bytes.Split(b, []byte{10})
	for i := range v {
		test(mkchars(v[i]))
		reset()
	}
	dofile(mkchars([]byte("k.t")), mkchars(b))
	reset()
	reset()
}
func TestKE(t *testing.T) {
	//t.Skip()
	b, err := ioutil.ReadFile("k.e")
	if err != nil {
		t.Fatal(err)
	}
	v := bytes.Split(b, []byte{10})
	for i := range v {
		newtest()
		var buf bytes.Buffer
		wasi_unstable.Stdout = &buf
		e := try(func() { test(mkchars(v[i])) })
		exp := parseError(strings.Split(strings.Split(string(v[i]), " /")[1], " ")[0])
		fmt.Println(string(v[i]))
		if e != exp {
			t.Fatalf("expected error %d got %d", exp, e)
		}
	}
}
func TestTraps(t *testing.T) {
	testCases := []struct {
		f func()
		e int32
	}{
		{func() { x := mk(It, 0); dx(x); dx(x) }, Unref},
		{func() { mk(2, 3) }, Value},
		{func() { cal(mk(It, 0), seq(3)) }, Type},
	}
	for _, tc := range testCases {
		newtest()
		wasi_unstable.Stdout = io.Discard
		e := try(tc.f)
		if e != tc.e {
			t.Fatalf("expected %d got %d", tc.e, e)
		}
	}
}
func try(f func()) (err int32) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(int32)
		}
	}()
	f()
	return -1
}
func parseError(s string) int32 {
	errs := []string{"Err", "Type", "Value", "Length", "Rank", "Parse", "Stack", "Grow", "Unref", "Io", "Nyi"} // err.go
	for i, x := range errs {
		if s == x {
			return int32(i)
		}
	}
	panic("unknown error: " + s)
}
func TestSymbols(t *testing.T) { // list symbols
	newtest()
	s := K(I32(8)) | K(St)<<59
	n := nn(s)
	for i := int32(0); i < n; i++ {
		y := 8 * i
		_ = y
		//fmt.Printf("%d %s\n", y, sK(Ks(y)))
	}
	reset()
}
func TestIndex(t *testing.T) {
	newtest()
	x := Flr(seq(5))
	dx(x)
	if index(2, int32(x), int32(x)+5) != 2 {
		t.Fatal()
	}
	if index(6, int32(x), int32(x)+5) != -1 {
		t.Fatal()
	}
	reset()
}
func TestClass(t *testing.T) {
	c := make([]byte, 127)
	cl := func(s string, n byte) {
		for _, b := range []byte(s) {
			c[b] |= n
		}
	}
	cl(`:+-*%!&|<>=~,^#_$?@.'/\`, 1)
	cl(`abcdefghijklmnopqrstuvxwyzABCDEFGHIJKLMNOPQRSTUVWXYZ`, 2)
	cl(`0123456789`, 4)
	cl(`'/\\`, 8)
	cl(`([{`, 16)
	cl("\n;)]}", 32)
	cl(`abcdefghijklmnopqrstuvxwyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789)]}"`, 64)
	cl(`0123456789abcdef`, 128)
	//fmt.Printf("%q\n", string(c[32:]))
}

/*
func memck() {
	icheck := func(i int32) {
		p := I32(4 * i)
		for p != 0 {
			if p < 4096 {
				fmt.Println("memck ", i, p)
				panic("memck")
			}
			p = I32(p)
		}
	}
	for i := int32(5); i < 32; i++ {
		icheck(i)
	}
}
*/

func check() { // debug in reset()
	for i := int32(5); i < 31; i++ {
		//fmt.Printf("[%d %d]: %d\n", i, 4*i, I32(4*i))
	}
	t := int32(0)
	for i := int32(5); i < 31; i++ {
		t += mark(i) * (int32(1) << i)
	}
	if t != 1<<I32(128) {
		total := (int32(1) << I32(128)) - 4096
		if total-t != 0 {
			fmt.Printf("free %d of %d (+%d)\n", t, total, total-t)
		}
	}
	scan()
}
func mark(i int32) (r int32) {
	p := I32(4 * i)
	for p != 0 {
		r++
		SetI32(p+12, 0) // rc
		if r := I32(p + 12); r != 0 {
			//fmt.Printf("mark: p=%d rc=%d\n", p, r)
			//panic("mark")
		}
		SetI32(p+4, i) // bt
		p = I32(p)
	}
	return r
}
func scan() {
	total := int32(1) << I32(128)
	p := int32(4096)
	for {
		t := I32(p + 4)
		if t < 5 || t > 31 {
			fmt.Printf("illegal type at p+16=%d, bt=%d\n", p+16, t)
			fmt.Printf("p=%d Ip=%d Ip+4=%d Ip+8=%d Ip+12=%d\n", p, I32(p), I32(p+4), I32(p+8), I32(p+12))
			panic("scan")
		}
		if r := I32(p + 12); r != 0 {
			fmt.Printf("non-free block at p=%d, bt=%d rc=%d\n", p, t, r)
			panic("scan")
		}
		p += 1 << t
		if p == total {
			return
		}
		if p > total {
			panic("scan/p>total")
		}
	}
}
func rc(x K) int32 {
	xt := tp(x)
	if xt < ft {
		return -1
	}
	return I32(int32(x) - 4)
}
func sK(x K) string {
	xp := int32(x)
	switch tp(x) {
	case 0:
		if x == 0 {
			return ""
		}
		s := []byte("0:+-*%!&|<>=~,^#_$?@.'/\\")
		var r string
		itoa := func(x int32) string { return strconv.Itoa(int(x)) }
		switch {
		case xp < 64:
			if xp < 23 {
				r = string(s[xp])
			} else {
				r = "`" + itoa(xp)
			}
			return r
		case xp < 128:
			if xp-64 < 23 {
				r = string(s[xp-64])
			} else {
				r = "`" + itoa(xp)
			}
			return r
		case xp == 211:
			return "@"
		case xp == 212:
			return "."
		case xp >= 448 && xp-448 < 23:
			return string(s[xp-448])
		default:
			return "`" + itoa(xp)
		}
	case bt:
		if int32(x) != 0 {
			return "1b"
		} else {
			return "0b"
		}
	case ct:
		return strconv.Quote(string([]byte{byte(xp)}))
	case it:
		return strconv.Itoa(int(xp))
	case st:
		x = cs(x)
		dx(x)
		xp = int32(x)
		if nn(x) == 0 {
			return "`"
		}
		return "`" + string(Bytes[xp:xp+nn(x)])
	case ft:
		return sflt(F64(xp))
	case zt:
		return sflz(F64(xp), F64(xp+8))
	case cf:
		xn := nn(x)
		xp = int32(x) + 8*xn
		s := ""
		for i := int32(0); i < xn; i++ {
			xp -= 8
			s += sK(K(I64(xp)))
		}
		return s
	case df:
		a := []string{"'", "':", "/", "/:", "\\", "\\:"}
		r := sK(K(I64(xp)))
		p := I64(xp + 8)
		return r + a[int(p)]
	case pf:
		f := K(I64(xp))
		l := K(I64(xp + 8))
		i := K(I64(xp + 16))
		// if tp(f) == 0 && nn(i) == 1 && I32(int32(i)) == 1 {
		if nn(i) == 1 && I32(int32(i)) == 1 {
			return sK(K(I64(int32(l)))) + sK(f) // 1+
		}
		return "<prj>"
	case lf:
		x = K(I64(xp + 16))
		xp = int32(x)
		return string(Bytes[xp : xp+nn(x)])
	case Bt:
		r := bytes.Repeat([]byte{'0'}, int(nn(x)))
		for i := range r {
			if I8(xp+int32(i)) != 0 {
				r[i] = '1'
			}
		}
		return comma(1 == nn(x)) + string(r) + "b"
	case Ct:
		return comma(1 == nn(x)) + strconv.Quote(string(Bytes[xp:xp+nn(x)]))
	case It:
		if nn(x) == 0 {
			return "!0"
		}
		r := make([]string, nn(x))
		for i := range r {
			r[i] = strconv.Itoa(int(I32(xp + 4*int32(i))))
		}
		return comma(1 == nn(x)) + strings.Join(r, " ")
	case St:
		r := make([]string, nn(x))
		for i := range r {
			r[i] = sK(K(I32(xp)) | K(st)<<59)
			xp += 4
		}
		if nn(x) == 0 {
			return "0#`"
		}
		return comma(1 == nn(x)) + strings.Join(r, "")
	case Ft:
		r := make([]string, nn(x))
		for i := range r {
			r[i] = sflt(F64(xp + 8*int32(i)))
		}
		return comma(1 == nn(x)) + strings.Join(r, " ")
	case Zt:
		r := make([]string, nn(x))
		for i := range r {
			r[i] = sflz(F64(xp), F64(xp+8))
			xp += 16
		}
		return comma(1 == nn(x)) + strings.Join(r, " ")
	case Lt:
		r := make([]string, nn(x))
		for i := range r {
			r[i] = sK(K(I64(xp)))
			xp += 8
		}
		if len(r) == 1 {
			return "," + r[0]
		} else {
			return "(" + strings.Join(r, ";") + ")"
		}
	case Dt:
		return sK(K(I64(xp))) + "!" + sK(K(I64(xp+8)))
	case Tt:
		return "+" + sK(K(I64(xp))) + "!" + sK(K(I64(xp+8)))
	default:
		fmt.Println("type ", tp(x))
		panic("type")
	}
}
func sflt(x float64) string {
	s := strconv.FormatFloat(x, 'g', -1, 64)
	if strings.Index(s, ".") < 0 {
		s += "."
	}
	return s
}
func sflz(x, y float64) (s string) {
	phi := 180.0 / math.Pi * math.Atan2(y, x)
	r := math.Hypot(x, y)
	s = strconv.FormatFloat(r, 'g', -1, 64) + "a"
	if phi != 0 {
		s += sflt(phi)
	}
	return s
}
func comma(x bool) string {
	if x {
		return ","
	} else {
		return ""
	}
}
