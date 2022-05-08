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
	"unicode/utf8"

	. "github.com/ktye/wg/module"
)

var save []byte

func newtest() {
	Stdout = os.Stdout
	rand_ = 1592653589
	if save == nil {
		kinit()
		save = make([]byte, len(Bytes))
		copy(save, Bytes)
	} else {
		Bytes = make([]byte, len(save))
		copy(Bytes, save)
		pp, pe, sp = 0, 0, 256
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
	//t.Skip()
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
	//t.Skip()
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
	//t.Skip()
	newtest()
	dx(val(Ku(35382781554225)))
	reset()
}
func TestVerbs(t *testing.T) {
	//t.Skip()
	newtest()
	x := Til(Ki(3))
	if r := int32(Cnt(x)); r != 3 {
		t.Fatalf("got %d expected %d", r, 3)
	}
}
func TestTok(t *testing.T) {
	//t.Skip()
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
	//t.Skip()
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
	//t.Skip()
	newtest()
	dx(shuffle(seq(8), 5))
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
		if len(v[i]) > 1 {
			test(mkchars(append(v[i], 10)))
			reset()
		}
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
		Stdout = &buf
		e := tryf(func() { test(mkchars(append(v[i], 10))) })
		exp := parseError(strings.Split(strings.Split(string(v[i]), " /")[1], " ")[0])
		fmt.Println(string(v[i]))
		if e != exp {
			t.Fatalf("expected error %d got %d", exp, e)
		}
	}
}
func TestTraps(t *testing.T) {
	//t.Skip()
	testCases := []struct {
		f func()
		e int32
	}{
		{func() { x := mk(It, 0); dx(x); dx(x) }, Unref},
		{func() { mk(2, 3) }, Value},
		{func() { cal(mk(It, 0), seq(3)) }, Type},
		{func() { test(seq(3)) }, Type},
		{func() { test(mkchars([]byte("1 2 /1 /2\n"))) }, Length},
		{func() { test(mkchars([]byte("1 /2\n"))) }, Err},
		{func() { sp = 0; reset() }, Stack},
		{func() { mk(It, 0); reset() }, Err},
		{func() { use(Key(seq(2), seq(2))) }, Type},
		{func() { nyi(Ki(0)) }, Nyi},
		{func() { ndrop(5, Key(seq(2), seq(2))) }, Type},
		{func() { mfree(int32(seq(1)), 5) }, Unref},
	}
	for i, tc := range testCases {
		newtest()
		Stdout = io.Discard
		e := tryf(tc.f)
		if e != tc.e {
			t.Fatalf("tc %d: expected %d got %d", i, tc.e, e)
		}
	}
}
func tryf(f func()) (err int32) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(int32)
		}
	}()
	f()
	return -1
}
func parseError(s string) int32 {
	errs := []string{"Err", "Type", "Value", "Index", "Length", "Rank", "Parse", "Stack", "Grow", "Unref", "Io", "Nyi"} // err.go
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
func TestRepl(t *testing.T) {
	testCases := [][2]string{
		[2]string{"1+1", "2\n"},
		[2]string{"x:!10", ""},
		[2]string{"\\c", ""},
		[2]string{"\\m", "18\n"},
	}
	for _, tc := range testCases {
		newtest()
		var buf bytes.Buffer
		Stdout = &buf
		repl(mkchars([]byte(tc[0])))
		if r := string(buf.Bytes()); r != tc[1] {
			t.Fatalf("expected %q got %q\n", tc[1], r)
		}
		reset()
	}
}
func TestIndex(t *testing.T) {
	newtest()
	x := Flr(seq(5))
	dx(x)
	if idx(2, int32(x), int32(x)+5) != 2 {
		t.Fatal()
	}
	if idx(6, int32(x), int32(x)+5) != -1 {
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
func Test360(t *testing.T) {
	//t.Skip()
	Tx := []struct {
		s     string
		shape []int
		r     interface{}
	}{
		{"-7", []int{}, []float64{-7}},
		{"'ALPHA⍴⍳'", []int{7}, "ALPHA⍴⍳"},
		{"1+2", []int{}, []float64{3}},
		{"3≤7", []int{}, []bool{true}},
		{"7≤3", []int{}, []bool{false}},
		{"1 2 3 4×4 3 2 1", []int{4}, []float64{4, 6, 6, 4}},
		{"2+1 2 3 4", []int{4}, []float64{3, 4, 5, 6}},
		{"1 2 3 4⌈2", []int{4}, []float64{2, 2, 3, 4}},
		{"1 2 3", []int{3}, []float64{1, 2, 3}},
		{"⍳3", []int{3}, []int{1, 2, 3}},
		{"0)⍳3", []int{3}, []int{0, 1, 2}},
		{"⍳0", []int{}, []int{}},
		{"6-⍳6", []int{6}, []float64{5, 4, 3, 2, 1, 0}},
		{"2×⍳0", []int{}, []float64{}},
		{"2×⍳6", []int{6}, []float64{2, 4, 6, 8, 10, 12}},
		{"X,X←2 3 5 7 11", []int{10}, []float64{2, 3, 5, 7, 11, 2, 3, 5, 7, 11}},
		{"1,2 3", []int{3}, []float64{1, 2, 3}},
		{"1 2,3", []int{3}, []float64{1, 2, 3}},
		{"1 2,⍴⍳5", []int{3}, []float64{1, 2, 5}},
		{",2 3⍴⍳6", []int{6}, []int{1, 2, 3, 4, 5, 6}},
		{"⍴1", []int{0}, []int{}},
		{"⍴⍴1", []int{1}, []int{0}},
		{"⍴⍴⍴1", []int{1}, []int{1}},
		{"⍴1 2 3", []int{1}, []int{3}},
		{"⍴⍴1 2 3", []int{1}, []int{1}},
		{"⍴⍴⍴1 2 3", []int{1}, []int{1}},
		{"⍴2 3⍴⍳6", []int{2}, []int{2, 3}},
		{"⍴⍴2 3⍴⍳6", []int{1}, []int{2}},
		{"⍴⍴⍴2 3⍴⍳6", []int{1}, []int{1}},
		{"⍴4 2 3⍴⍳6", []int{3}, []int{4, 2, 3}},
		{"⍴⍴4 2 3⍴⍳6", []int{1}, []int{3}},
		{"⍴⍴⍴4 2 3⍴⍳6", []int{1}, []int{1}},
		{"2 3⍴1 2", []int{2, 3}, []float64{1, 2, 1, 2, 1, 2}},
		{"⍴''", []int{1}, []int{0}},
		{"''", []int{0}, ""},
		{"'X'", []int{1}, "X"},
		{"'CAN''T'", []int{5}, "CAN'T"},
		{"A←'ABCDEFG'", []int{7}, "ABCDEFG"},
		{"M←4 3⍴3 1 4 2 1 4 4 1 2 4 1 4", []int{4, 3}, []float64{3, 1, 4, 2, 1, 4, 4, 1, 2, 4, 1, 4}},
		{"A[M]", []int{4, 3}, "CADBADDABDAD"},
		{"(3 4⍴⍳12)[2;3]", []int{}, []int{7}},
		{"(M←3 4⍴⍳12)[1 3;2 3 4]", []int{2, 3}, []int{2, 3, 4, 10, 11, 12}},
		{"M[2;]", []int{4}, []int{5, 6, 7, 8}},
		{"M[;2 1]", []int{3, 2}, []int{2, 1, 6, 5, 10, 9}},
		{"M[M←4 3⍴3 1 4 2 1 4 4 1 2 4 1 4;]", []int{4, 3, 3}, []float64{4, 1, 2, 3, 1, 4, 4, 1, 4, 2, 1, 4, 3, 1, 4, 4, 1, 4, 4, 1, 4, 3, 1, 4, 2, 1, 4, 4, 1, 4, 3, 1, 4, 4, 1, 4}},
		{"(3 3⍴⍳9)[;-3-4]", []int{3}, []int{1, 4, 7}},
		{"3 3⍴1 0 0 0", []int{3, 3}, []float64{1, 0, 0, 0, 1, 0, 0, 0, 1}},
		{"3⍴1", []int{3}, []float64{1, 1, 1}},
		{"2 3⍴⍳6", []int{2, 3}, []int{1, 2, 3, 4, 5, 6}},
		{"2 3 4⍴⍳6", []int{2, 3, 4}, []int{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6}},
		{"(1-2)", []int{}, []float64{-1}},
		{"(⍳3)+4", []int{3}, []float64{5, 6, 7}},
		{"X←⍳6", []int{6}, []int{1, 2, 3, 4, 5, 6}},
		{"X[2 1]←8 9", []int{6}, []float64{9, 8, 3, 4, 5, 6}},
		{"X[2 1 3]←7", []int{6}, []float64{7, 7, 7, 4, 5, 6}},
		{"3 4⍳4", []int{}, []int{2}},
		{"1 2 3 4⍳2 3⍴⍳6", []int{2, 3}, []int{1, 2, 3, 4, 5, 5}},
		{"'ABCDEFGH'⍳'GAFFE'", []int{5}, []int{7, 1, 6, 6, 5}},
		{"3 4 7∊⍳5", []int{3}, []bool{true, true, false}},
		{"4∊⍳5", []int{}, []bool{true}},
		{"2 3 5 7∊⍳4", []int{4}, []bool{true, true, false, false}},
		{"(3 4⍴⍳12)∊2 3 5 7", []int{3, 4}, []bool{false, true, true, false, true, false, true, false, false, false, false, false}},
		{"0↑3 4 5 6", []int{0}, []float64{}},
		{"2↑3 4 5 6", []int{2}, []float64{3, 4}},
		{"2 3↑3 4⍴⍳12", []int{2, 3}, []int{1, 2, 3, 5, 6, 7}},
		{"2 ¯3↑3 4⍴⍳12", []int{2, 3}, []int{2, 3, 4, 6, 7, 8}},
		{"2↓3 4 5 6", []int{2}, []float64{5, 6}},
		{"1 2↓3 4⍴⍳12", []int{2, 2}, []int{7, 8, 11, 12}},
		{"¯1 1↓3 4⍴⍳12", []int{2, 3}, []int{2, 3, 4, 6, 7, 8}},
		{"⍋3 5 3 2", []int{4}, []int{4, 1, 3, 2}},
		{"0)⍋3 5 3 2", []int{4}, []int{3, 0, 2, 1}},
		{"⍒3 5 3 2", []int{4}, []int{2, 1, 3, 4}},
		{"P←2 3 5 7", []int{4}, []float64{2, 3, 5, 7}},
		{"E←3 4⍴⍳12", []int{3, 4}, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}},
		{"1 0 1 0/P", []int{2}, []float64{2, 5}},
		{"1 0 1 0/E", []int{3, 2}, []int{1, 3, 5, 7, 9, 11}},
		{"1 0 1/[1]E", []int{2, 4}, []int{1, 2, 3, 4, 9, 10, 11, 12}},
		{"1 0 1⌿E", []int{2, 4}, []int{1, 2, 3, 4, 9, 10, 11, 12}},
		{"0 1 0 1/⍳4", []int{2}, []int{2, 4}},
		{"1 0 1\\⍳2", []int{3}, []int{1, 0, 2}},
		{"1 0 1 1 1\\3 4⍴'ABCDEFGHIJKL'", []int{3, 5}, "A BCDE FGHI JKL"},
		{"1 0 1⍀2 3⍴⍳6", []int{3, 3}, []int{1, 2, 3, 0, 0, 0, 4, 5, 6}},
		{"⍉2 3⍴⍳6", []int{3, 2}, []int{1, 4, 2, 5, 3, 6}},
		{"1 1⍉3 3⍴⍳9", []int{3}, []int{1, 5, 9}},
		{"⍴⍴(⍳0)⍉5", []int{1}, []int{0}},
		{"⌽X←3 4⍴'ABCDEFGHIJKL'", []int{3, 4}, "DCBAHGFELKJI"},
		{"⌽X", []int{3, 4}, "DCBAHGFELKJI"},
		{"⌽[1]3 4⍴'ABCDEFGHIJKL'", []int{3, 4}, "IJKLEFGHABCD"},
		{"⊖X", []int{3, 4}, "IJKLEFGHABCD"},
		{"3⌽2 3 5 7", []int{4}, []float64{7, 2, 3, 5}},
		{"¯1⌽2 3 5 7", []int{4}, []float64{7, 2, 3, 5}},
		{"¯7⌽'ABCDEF'", []int{6}, "FABCDE"},
		{"1⊖3 4⍴⍳12", []int{3, 4}, []int{5, 6, 7, 8, 9, 10, 11, 12, 1, 2, 3, 4}},
		{"1 0 ¯1⌽3 4⍴'ABCDEFGHIJKL'", []int{3, 4}, "BCDAEFGHLIJK"},
		{"(2 2⍴0 ¯1 3 1)⌽[2]2 4 2⍴⍳16", []int{2, 4, 2}, []int{1, 8, 3, 2, 5, 4, 7, 6, 15, 12, 9, 14, 11, 16, 13, 10}},
		{"(2 4⍴0 1 ¯1 0 0 3  2 1)⌽[1]2 2 4⍴⍳16", []int{2, 2, 4}, []int{1, 10, 11, 4, 5, 14, 7, 16, 9, 2, 3, 12, 13, 6, 15, 8}},
		{"10⊥1 7 7 6", []int{}, []float64{1776}},
		{"24 60 60⊥1 2 3", []int{}, []float64{3723}},
		{"24 60 60⊤3723", []int{3}, []int{1, 2, 3}},
		{"60 60⊤3723", []int{2}, []int{2, 3}},
		{"1+-3", []int{}, []float64{-2}},
		{"A←4", []int{}, []float64{4}},
		{"A+2+A←4", []int{}, []float64{10}},
		{"3 4[2]", []int{}, []float64{4}},
		{"+/1 2 3", []int{}, []float64{6}},
		{"+/⍳3", []int{}, []int{6}},
		{"+/2 3⍴⍳6", []int{2}, []int{6, 15}},
		{"+/[1]2 3⍴⍳6", []int{3}, []int{5, 7, 9}},
		{"+⌿2 3⍴⍳6", []int{3}, []int{5, 7, 9}},
		{"+/[2]2 4 3⍴⍳24", []int{2, 3}, []int{22, 26, 30, 70, 74, 78}},
		{"2 3 4+.×5 6 7", []int{}, []float64{56}},
		{"1 2 3∘.×1 2 3 4", []int{3, 4}, []float64{1, 2, 3, 4, 2, 4, 6, 8, 3, 6, 9, 12}},
		{"(⍳3)∘.×⍳4", []int{3, 4}, []int{1, 2, 3, 4, 2, 4, 6, 8, 3, 6, 9, 12}},
	}
	nt := func(x K, y int, xt T) {
		if n := nn(x); int(n) != y {
			t.Fatalf("length is %d should be %d", n, y)
		}
		if tp(x) != xt {
			t.Fatalf("type is %v not %v", tp(x), xt)
		}
	}
	chars := func(x K, y []byte) {
		//fmt.Println("chars", sK(x), "|", y)
		nt(x, len(y), Ct)
		for i := 0; i < len(y); i++ {
			if xi := I8(int32(x) + int32(i)); xi != int32(int8(y[i])) {
				t.Fatalf("x[%d] is %d not %d", i, xi, y[i])
			}
		}
	}
	runes := func(x K, y string) {
		//fmt.Println("runes", sK(x), "|", y)
		ny := 0
		for range y {
			ny++
		}
		buf := make([]byte, 4)
		nt(x, ny, Lt)
		i := 0
		for _, r := range y {
			nr := utf8.EncodeRune(buf, r)
			chars(K(I64(int32(x)+int32(8*i))), buf[:nr])
			i++
		}
	}
	bools := func(x K, y []bool) {
		//fmt.Println("bools", sK(x), "|", y)
		nt(x, len(y), Bt)
		for i := 0; i < len(y); i++ {
			if xi := I8(int32(x) + int32(i)); xi != I32B(y[i]) {
				t.Fatalf("x[%d] is %d not %v", i, xi, y[i])
			}
		}
	}
	ints := func(x K, y []int) {
		//fmt.Println("ints", sK(x), "|", y)
		nt(x, len(y), It)
		for i := 0; i < len(y); i++ {
			if xi := I32(int32(x) + 4*int32(i)); xi != int32(y[i]) {
				t.Fatalf("x[%d] is %d not %d", i, xi, y[i])
			}
		}
	}
	floats := func(x K, y []float64) {
		//fmt.Println("floats", sK(x), "|", y)
		nt(x, len(y), Ft)
		for i := 0; i < len(y); i++ {
			if xi := F64(int32(x) + 8*int32(i)); xi != y[i] {
				t.Fatalf("x[%d] is %v not %v", i, xi, y[i])
			}
		}
	}
	newtest()
	x := mkchars([]byte("apl/apl.k"))
	dofile(x, readfile(rx(x)))
	def := val(mkchars([]byte("APL")))
	apl := val(mkchars([]byte("{RUN TOK x}"))) //lup(sc(mkchars([]byte("APL"))))
	o := sc(mkchars([]byte{'O'}))
	for _, tc := range Tx {
		fmt.Println(tc.s)
		Asn(o, Ki(1))
		if strings.HasPrefix(tc.s, "0)") {
			tc.s = tc.s[2:]
			Asn(o, Ki(0))
		}
		if tc.shape == nil && tc.r == nil {
			r := Atx(rx(def), mkchars([]byte(tc.s)))
			if r != Ki(0) {
				t.Fatalf("expected nil %s", sK(r))
			}
			continue
		}
		r := Atx(rx(apl), mkchars([]byte(tc.s)))
		nt(r, 3, Lt)
		shape := x1(int32(r))
		ravel := x2(int32(r))
		ints(shape, tc.shape)
		switch v := tc.r.(type) {
		case string:
			runes(ravel, v)
		case []bool:
			bools(ravel, v)
		case []int:
			ints(ravel, v)
		case []float64:
			floats(ravel, v)
		default:
			t.Fatal("wrong r type")
		}
		dx(shape)
		dx(ravel)
		dx(r)
	}
	dx(apl)
	dx(def)
	reset()
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
		total := (int32(1) << I32(128)) - 8192
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
	p := int32(8192)
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
