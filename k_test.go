package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
	"testing"

	. "github.com/ktye/wg/module"
)

var save []byte

func newtest() {
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
func TestVerbs(t *testing.T) {
	//t.Skip()
	newtest()
	x := Til(Ki(3))
	if r := iK(Cnt(x)); r != 3 {
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
func TestKT(t *testing.T) {
	//t.Skip()
	newtest()
	b, err := ioutil.ReadFile("t")
	if err != nil {
		t.Fatal(err)
	}
	x := mkchars(b)
	test(x, 0)
}
func TestK(t *testing.T) {
	t.Skip()
	b, err := ioutil.ReadFile("t")
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range bytes.Split(b, []byte("\n")) {
		s := string(v)
		if len(s) == 0 {
			continue
		}
		a := strings.Split(s, " /")
		in := a[0]
		exp := a[1]
		fmt.Println(in, "/"+exp)
		newtest()
		//fmt.Println("newtest")
		x := mkchars([]byte(a[0]))
		x = val(x)
		got := sK(x)
		if got != exp {
			t.Fatalf("%s:\nexp: %s\ngot: %s", in, exp, got)
		}
		dx(x)
		//check(t)
		reset()
		reset()
	}
}
func check(t *testing.T) {
	if sp != 256 {
		t.Fatalf("nonempty stack: sp=%d", sp)
	}
	if n := rc(src); n > 1 {
		t.Fatalf("src: rc %d\n", n)
	}
	dx(src)
	src = 0
	dx(xyz)
	dx(K(I64(0)))
	dx(K(I64(8)))
	m := mcount()
	u := (uint32(1) << uint32(I32(128))) - 1024
	d := int32(u - m)
	mark()
	scan()
	if d != 0 {
		t.Fatalf("m %d, diff %d", m, int32(u-m))
	}
}
func mark() {
	for i := int32(5); i < 31; i++ {
		marki(i)
	}
}
func marki(i int32) {
	p := I32(4 * i)
	l := int32(0)
	for p != 0 {
		if p&31 != 0 {
			panic(fmt.Errorf("illegal block in free list: %d type %d (last=%d)", p, i, l))
		}
		n := I32(p)
		SetI32(p, i)
		SetI32(4+p, 1<<uint32(i))
		l = p
		p = n
	}
}
func scan() {
	p := int32(1024)
	a := I32(128)
	if a < 10 || a > 32 {
		panic(fmt.Errorf("illegal total alloc %d", a))
	}
	e := int32(1) << I32(128)
	for {
		if p == e {
			return
		}
		k := p + 16
		t := I32(p)
		if t < 5 || t > 31 {
			panic(fmt.Errorf("%d: illegal block type=%d at %d", k, t, p))
		}
		s := I32(4 + p)
		if s&31 != 0 {
			panic(fmt.Errorf("%d: illegal size=%d type=%d at %d", k, s, t, p))
		}
		p += s
		if p == e {
			return
		} else if p > e {
			panic(fmt.Errorf("%d: illegal block %d > e(%d)", k, p, e))
		}
	}
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
	//fmt.Printf("%q\n", string(c[32:]))
}

func rc(x K) int32 {
	xt := tp(x)
	if xt < 16 {
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
				r = itoa(xp)
			}
			return r
		case xp < 128:
			if xp-64 < 23 {
				r = string(s[xp-64])
			} else {
				r = itoa(xp)
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
		x = K(I64(xp + 24))
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
		panic("nyi-Dt")
	case Tt:
		panic("nyi-Tt")
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
