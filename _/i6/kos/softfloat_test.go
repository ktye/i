// from: go.dev/src/runtime/softfloat64_test.go

package main

import (
	//"fmt"
	"math"
	"math/rand"
	"testing"
)

// from: github.com/ktye/wg/module
func I32B(b bool) int32 {
	if b {
		return 1
	}
	return 0
}
func F64reinterpret_i64(x uint64) float64 { return math.Float64frombits(x) }
func I64reinterpret_f64(x float64) uint64 { return math.Float64bits(x) }

func add(x, y float64) float64 { return x + y }
func sub(x, y float64) float64 { return x - y }
func mul(x, y float64) float64 { return x * y }
func div(x, y float64) float64 { return x / y }
func fop(f func(x, y uint64) uint64) func(x, y float64) float64 {
	return func(x, y float64) float64 {
		bx := math.Float64bits(x)
		by := math.Float64bits(y)
		return math.Float64frombits(f(bx, by))
	}
}

func TestFloat64(t *testing.T) {
	base := []float64{
		0,
		math.Copysign(0, -1),
		-1,
		1,
		math.NaN(),
		math.Inf(+1),
		math.Inf(-1),
		0.1,
		1.5,
		1.9999999999999998,     // all 1s mantissa
		1.3333333333333333,     // 1.010101010101...
		1.1428571428571428,     // 1.001001001001...
		1.112536929253601e-308, // first normal
		2,
		4,
		8,
		16,
		32,
		64,
		128,
		256,
		3,
		12,
		1234,
		123456,
		-0.1,
		-1.5,
		-1.9999999999999998,
		-1.3333333333333333,
		-1.1428571428571428,
		-2,
		-3,
		1e-200,
		1e-300,
		1e-310,
		5e-324,
		1e-105,
		1e-305,
		1e+200,
		1e+306,
		1e+307,
		1e+308,
	}
	all := make([]float64, 200)
	copy(all, base)
	for i := len(base); i < len(all); i++ {
		all[i] = rand.NormFloat64()
	}
	test(t, "+", add, fop(fadd), all)
	test(t, "-", sub, fop(fsub), all)
	test(t, "*", mul, fop(fmul), all)
	test(t, "/", div, fop(fdiv), all)
	test(t, "cps", math.Copysign, fop(fcps), all)

	//fmin/fmax is like libm and differs from go on 0n
	//test(t, "max", math.Max, fop(fmax), all)
	//test(t, "min", math.Min, fop(fmin), all)
}

func test(t *testing.T, op string, hw, sw func(float64, float64) float64, all []float64) {
	for _, f := range all {
		for _, g := range all {
			//fmt.Printf("%v %s %v\n", f, op, g)
			h := hw(f, g)
			s := sw(f, g)
			if !same(h, s) {
				t.Fatalf("%g %s %g = sw %g, hw %g\n", f, op, g, s, h)
			}
			testu(t, "jfcst", hwff, swjf, h)
			testu(t, "fjcst", hwff, swfj, h)
			testcmp(t, f, h)
			testcmp(t, h, f)
			testcmp(t, g, h)
			testcmp(t, h, g)
			testfabs(t, h)
			testfloor(t, h)
			testsqrt(t, h)
		}
	}
}

func testu(t *testing.T, op string, hw, sw func(float64) float64, v float64) {
	h := hw(v)
	s := sw(v)
	if !same(h, s) {
		t.Fatalf("%s %g = sw %g, hw %g\n", op, v, s, h)
	}
}
func hwff(f float64) float64 { return float64(int64(f)) }
func swfj(f float64) float64 { return math.Float64frombits(fjcst(int64(f))) }
func swjf(f float64) float64 {
	u := math.Float64bits(f)
	i := jfcst(math.Float64bits(f))

	if _, ok := f64toint(u); !ok {
		i = int64(f)
	}
	return float64(i)
}

func testcmp(t *testing.T, f, g float64) {
	x, y := math.Float64bits(f), math.Float64bits(g)
	eqh := f == g
	neh := f != g
	gth := f > g
	geh := f >= g
	lth := f < g
	leh := f <= g
	eqs := feql(x, y) != 0
	nes := fneq(x, y) != 0
	gts := fmor(x, y) != 0
	ges := fgte(x, y) != 0
	lts := fles(x, y) != 0
	les := flte(x, y) != 0
	if eqh != eqs {
		t.Fatalf("(%g == %g) = sw %v, hw %v\n", f, g, eqh, eqs)
	}
	if neh != nes {
		t.Fatalf("(%g != %g) = sw %v, hw %v\n", f, g, neh, nes)
	}
	if gth != gts {
		t.Fatalf("(%g > %g) = sw %v, hw %v\n", f, g, gth, gts)
	}
	if geh != ges {
		t.Fatalf("(%g >= %g) = sw %v, hw %v\n", f, g, geh, ges)
	}
	if leh != les {
		t.Fatalf("(%g < %g) = sw %v, hw %v\n", f, g, lth, lts)
	}
	if leh != les {
		t.Fatalf("(%g <= %g) = sw %v, hw %v\n", f, g, leh, les)
	}
}
func testfabs(t *testing.T, f float64) {
	u := math.Float64bits(f)
	h := math.Abs(f)
	s := fabs(u)
	if math.Float64bits(h) != s {
		sf := math.Float64frombits(s)
		t.Fatalf("fabs %v (%x) = (%v hw, %v sw)", f, u, h, sf)
	}
}
func testfloor(t *testing.T, f float64) {
	u := math.Float64bits(f)
	h := math.Floor(f)
	s := flor(u)
	if math.Float64bits(h) != s {
		sf := math.Float64frombits(s)
		t.Fatalf("floor %v (%x) = (%v hw, %v sw)", f, u, h, sf)
	}
}
func testsqrt(t *testing.T, f float64) {
	u := math.Float64bits(f)
	h := math.Sqrt(f)
	s := fsqr(u)
	sf := math.Float64frombits(s)
	if same(h, sf) == false {
		t.Fatalf("sqrt %v (%x) = (%v hw, %v sw)", f, u, h, sf)
	}
}

func same(f, g float64) bool {
	if math.IsNaN(f) && math.IsNaN(g) {
		return true
	}
	if math.Copysign(1, f) != math.Copysign(1, g) {
		return false
	}
	return f == g
}

/*
func funpack64(f uint64) (sign, mant uint64, exp int, inf, nan bool) {
	sign = f & (1 << (mantbits64 + expbits64))
	mant = f & (1<<mantbits64 - 1)
	exp = int(f>>mantbits64) & (1<<expbits64 - 1)
	switch exp {
	case 1<<expbits64 - 1:
		if mant != 0 {
			nan = true
			return
		}
		inf = true
		return
	case 0:
		// denormalized
		if mant != 0 {
			exp += bias64 + 1
			for mant < 1<<mantbits64 {
				mant <<= 1
				exp--
			}
		}
	default:
		// add implicit top bit
		mant |= 1 << mantbits64
		exp += bias64
	}
	return
}
const (
	mantbits64 uint = 52
	expbits64  uint = 11
	bias64          = -1<<(expbits64-1) + 1

	nan64 uint64 = (1<<expbits64-1)<<mantbits64 + 1<<(mantbits64-1) // quiet NaN, 0 payload
	inf64 uint64 = (1<<expbits64 - 1) << mantbits64
	neg64 uint64 = 1 << (expbits64 + mantbits64)

	mantbits32 uint = 23
	expbits32  uint = 8
	bias32          = -1<<(expbits32-1) + 1

	nan32 uint32 = (1<<expbits32-1)<<mantbits32 + 1<<(mantbits32-1) // quiet NaN, 0 payload
	inf32 uint32 = (1<<expbits32 - 1) << mantbits32
	neg32 uint32 = 1 << (expbits32 + mantbits32)
)
*/
