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
	r := mk(It, 2)
	if rc := I32(int32(r) - 16); rc != 1 {
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
	newtest()
	x := Til(Ki(3))
	if r := iK(Cnt(x)); r != 3 {
		t.Fatalf("got %d expected %d", r, 3)
	}
}
func TestTok(t *testing.T) {
	tc := []struct {
		in, exp string
	}{
		{"1234567", "(1234567)"},
		{"-1234567", "(-1234567)"},
		{"*", "(*)"},
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
func TestK(t *testing.T) {
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
		x := mkchars([]byte(a[0]))
		got := sK(exec(parse(x)))
		if got != exp {
			t.Fatalf("%s:\nexp: %s\ngot: %s", in, exp, got)
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

func rc(x K) int32 { return I32(int32(x) - 16) }
func sK(x K) string {
	xp := int32(x)
	switch tp(x) {
	case 0:
		if x == 0 {
			return ""
		}
		if xp < 23 {
			s := []byte("0:+-*%!&|<>=~,^#_$?@.'/\\")
			return string(s[xp])
		} else {
			return fmt.Sprintf("<%d>", x)
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
		return "<drv>"
	case pf:
		f := K(I64(xp))
		l := K(I64(xp + 8))
		i := K(I64(xp + 16))
		if tp(f) == 0 && nn(i) == 1 && I32(int32(i)) == 1 {
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
		return "(" + strings.Join(r, ";") + ")"
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
