package k

import (
	"math"
	"testing"

	. "github.com/ktye/wg/module"
)

func newtest() {
	Bytes = make([]byte, 64*1024)
	minit(10, 16)
}
func intvalue(x K) int32     { return int32(x) }
func floatvalue(x K) float64 { return F64(int32(x)) }
func TestTypes(t *testing.T) {
	newtest()
	xi := Ki(5)
	if v := intvalue(xi); v != 5 {
		t.Fatalf("got v=%d expected %d\n", v, 5)
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
