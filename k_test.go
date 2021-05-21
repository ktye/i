package k

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"

	. "github.com/ktye/wg/module"
)

func newtest() {
	Bytes = make([]byte, 64*1024)
	kinit()
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
	if n := intvalue(Count(xi)); n != 1 {
		t.Fatalf("got n=%d expected %d\n", n, 1)
	}
	xf := Kf(math.Pi)
	if f := floatvalue(xf); f != math.Pi {
		t.Fatalf("got f=%v expected %v\n", f, math.Pi)
	}
	if r := tp(xf); r != ft {
		t.Fatalf("got t=%d expected %d\n", r, ft)
	}
	if s := sK(cvb); s != `":+-*%!&|<>=~,^#_$?@.'/\\"` {
		t.Fatal(s)
	}
}
func TestBucket(t *testing.T) {
	tc := []struct{ in, exp int32 }{
		{0, 4},
		{4, 4},
		{8, 4},
		{9, 5},
		{24, 5},
		{25, 6},
	}
	for _, tc := range tc {
		if got := bucket(tc.in); got != tc.exp {
			t.Fatalf("bucket %d => %d (exp %d)\n", tc.in, got, tc.exp)
		}
	}
}
func TestVerbs(t *testing.T) {
	newtest()
	x := Til(Ki(3))
	if r := iK(Count(x)); r != 3 {
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
		got := sK(tok(mkchars([]byte(tc.in))))
		if got != tc.exp {
			t.Fatalf("got %s expected %s", got, tc.exp)
		}
	}
}

func rc(x K) int32 { return I32(int32(x) - 8) }
func sK(x K) string {
	xp := int32(x)
	switch tp(x) {
	case 0:
		if x == 0 {
			return "<null>"
		}
		if x > 23 {
			panic("verb " + strconv.FormatInt(int64(x), 10))
		} else {
			s := []byte("0:+-*%!&|<>=~,^#_$?@.'/\\")
			return string(s[x])
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
		panic("nyi-st")
	case ft:
		return strconv.FormatFloat(F64(xp), 'g', -1, 64)
	case zt:
		panic("nyi-zt")
	case Bt:
		r := bytes.Repeat([]byte{'0'}, int(nn(x)))
		for i := range r {
			if I8(xp+int32(i)) != 0 {
				r[i] = '1'
			}
		}
		return string(r) + "b"
	case Ct:
		return strconv.Quote(string(Bytes[xp : xp+nn(x)]))
	case It:
		r := make([]string, nn(x))
		for i := range r {
			r[i] = strconv.Itoa(int(I32(xp + 4*int32(i))))
		}
		return strings.Join(r, " ")
	case St:
		panic("nyi-St")
	case Ft:
		r := make([]string, nn(x))
		for i := range r {
			r[i] = strconv.FormatFloat(F64(xp+8*int32(i)), 'g', -1, 64)
		}
		return strings.Join(r, " ")
	case Zt:
		panic("nyi-Zt")
	case Lt:
		r := make([]string, nn(x))
		for i := range r {
			r[i] = sK(K(I64(xp + 8*int32(i))))
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
