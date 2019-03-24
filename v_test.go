package i

import (
	_fmt "fmt"
	"math"
	"reflect"
	"testing"
)

type (
	iv = []i
)

func c(r, i float64) complex128 { return complex(r, i) }

func TestMV(t *testing.T) {
	type IV []int
	testCases := []struct {
		s    string
		f    func(v) v
		x, r v
	}{
		{"flp", flp, 1.0, 1.0},
		{"flp", flp, iv{1, 2}, iv{1, 2}},
		{"flp", flp, l{fv{1, 2}, fv{3, 4}, fv{5, 6}}, l{fv{1, 3, 5}, fv{2, 4, 6}}},
		// {"flp", flp, l{fv{1, 2}, fv{3, 4}, l{5, l{6, 7}}}, l{l{1, 3, 5}, l{2, 4, l{6, 7}}}}, // eql but fail?
		{"neg", neg, 1.0, -1.0},
		{"neg", neg, c(1, 2), c(-1, -2)},
		{"neg", neg, fv{1, 2, 3}, fv{-1, -2, -3}},
		{"neg", neg, zv{1, c(2, 3)}, zv{-1, c(-2, -3)}},
		{"neg", neg, iv{1, 2}, iv{-1, -2}},
		{"neg", neg, true, true},
		{"neg", neg, uint16(4), uint16(65532)},
		{"neg", neg, 1, -1},
		{"neg", neg, mynum("33"), mynum("-33")},
		{"neg", neg, []mynum{"a", "b", "c"}, []mynum{"-a", "-b", "-c"}},
		{"neg", neg, myvec{"a", "b"}, myvec{"<a", "b>"}},
		{"neg", neg, map[v]v{"a": fv{1, 2}}, map[v]v{"a": fv{-1.0, -2.0}}},
		{"neg", neg, mystruct{true, 2.0, []myint{1, 2, 3}}, mystruct{true, -2.0, []myint{-1, -2, -3}}},
		{"fst", fst, iv{5, 6, 7}, 5},
		{"fst", fst, fv{5, 6, 7}, 5.0},
		{"fst", fst, l{c(2, 3), 0, 0}, c(2, 3)},
		{"fst", fst, IV{5, 6}, 5},
		{"fst", fst, IV{}, nil},
		{"fst", fst, [2]l{l{"d", "c"}, l{5, 6}}, 5},
		{"sqr", sqr, 4, 2},
		// {"sqr", sqr, -1.0, math.NaN()}, not comparable
		{"sqr", sqr, -7 + 24i, c(3, 4)},
		{"inv", inv, 4.0, 0.25},
		{"inv", inv, c(0, 0.5), c(0, -2)},
		{"abs", abs, -2, 2},
		{"abs", abs, fv{-2, -3, 4}, fv{2, 3, 4}},
		{"abs", abs, zv{c(3, 4), 5}, zv{5, 5}},
		{"til", til, 3.0, fv{0, 1, 2}},
		{"til", til, 3, iv{0, 1, 2}},
		{"til", til, 0, iv{}},
		{"til", til, 0.0, fv{}},
		{"odo", odo, fv{2, 3}, l{fv{0, 0, 0, 1, 1, 1}, fv{0, 1, 2, 0, 1, 2}}},
		{"odo", odo, l{true, 2, 1.0, c(3, 0)}, l{fv{0, 0, 0, 0, 0, 0}, fv{0, 0, 0, 1, 1, 1}, fv{0, 0, 0, 0, 0, 0}, fv{0, 1, 2, 0, 1, 2}}},
		{"wer", wer, 3, fv{0, 0, 0}},
		{"wer", wer, zv{3}, fv{0, 0, 0}},
		{"wer", wer, []bool{false, false, true, false, true, true}, fv{2, 4, 5}},
		{"wer", wer, l{false, 0, c(1, 0), 0.0, 1.0, myint(1)}, fv{2, 4, 5}},
		{"rev", rev, fv{1, 2, 3}, fv{3, 2, 1}},
		// {"rev", rev, dct(l{"a", "b"}, l{1, 2}), "→[b:2;a:1]"}, // cannot compare
		{"asc", asc, 3, fv{0}},
		{"asc", asc, fv{4, 5, 6}, fv{0, 1, 2}},
		{"asc", asc, sv{"be", "g", "a"}, fv{2, 0, 1}},
		{"asc", asc, map[v]f{"b": 3, "c": 2, "a": 5}, sv{"c", "b", "a"}},
		{"asc", asc, "a", fv{0}},
		{"asc", asc, sv{"b", "c", "alpha"}, fv{2, 0, 1}},
		{"dsc", dsc, fv{5, -1, 3}, fv{0, 2, 1}},
		{"dsc", dsc, sv{"b", "c", "alpha"}, fv{1, 0, 2}},
		{"eye", eye, 0, l{}},
		{"eye", eye, 2, l{fv{1, 0}, fv{0, 1}}},
		// {"grp", grp, fv{1, 3, 3, 3, 1, 2}, map[v]v{1: fv{0, 4}, 3: fv{1, 2, 3}, 2: fv{5}}}, // eql but fail?
		{"not", not, 1, 0},
		{"not", not, 1 + 2i, 0 + 0i},
		{"not", not, 0 + 0i, 1 + 0i},
		{"not", not, c(math.Inf(1), 0), 0 + 0i},
		{"not", not, c(math.NaN(), 0), 0 + 0i},
		{"enl", enl, 1.2, fv{1.2}},
		{"enl", enl, iv{1, 2}, l{iv{1, 2}}},
		{"enl", enl, IV{4, 5, 6}, l{IV{4, 5, 6}}},
		{"is0", is0, 0, 0},
		{"is0", is0, nil, 1.0},
		{"is0", is0, iv{}, iv{}},
		{"is0", is0, l{1, math.NaN(), c(1, 0), c(1, math.NaN())}, l{0, 1.0, c(0, 0), c(1, 0)}},
		{"exp", exp, l{0, c(1, 0)}, l{1, c(math.E, 0)}},
		{"log", log, fv{1, 1.0 / math.E}, fv{0, -1}},
		{"log", log, c(-1, 0), c(0, math.Pi)},
		{"cnt", cnt, l{}, 0.0},
		{"cnt", cnt, iv{1, 2, 3}, 3.0},
		{"cnt", cnt, 4, 1.0},
		{"cnt", cnt, "alpha", 1.0},
		{"cnt", cnt, map[v]v{"a": iv{1, 2, 3}, "b": iv{2, 3, 4}}, 2.0},
		{"flr", flr, 3.5, 3.0},
		// fmt TODO
		{"unq", unq, l{}, l{}},
		{"unq", unq, iv{1, 2, 2, 1}, iv{1, 2}},
		{"unq", unq, l{1.0, 1.0}, l{1.0}},
		// evl TODO
	}

	for _, tc := range testCases {
		r := tc.f(tc.x)
		tt(t, tc.r, r, "%s %+v: %+v\n", tc.s, tc.x, r)
	}
}

func TestDV(t *testing.T) {
	type IV []int
	testCases := []struct {
		s       string
		f       func(v, v) v
		x, y, r v
	}{
		{"add", add, 1, 2.0, 3.0},
		{"add", add, 1.0, uint(3), 4.0},
		{"add", add, iv{1, 2}, 3, iv{4, 5}},
		{"add", add, iv{1, 2}, 3, iv{4, 5}},
		{"add", add, 1, iv{2, 3}, iv{3, 4}},
		{"add", add, 1, iv{2, 3}, iv{3, 4}},
		{"add", add, iv{1, 2, 3}, iv{4, 5, 6}, iv{5, 7, 9}},
		{"add", add, l{1, 2.0, 3}, 1, l{2, 3.0, 4}},
		{"add", add, iv{1, 2}, l{1, iv{2, 3}}, l{2, iv{4, 5}}},
		//{"add", add, [2]l{l{"a", "b"}, l{1, 2.0}}, [2]l{l{"b"}, l{fv{3, 4}}}, [2]l{l{"a", "b"}, l{1, fv{5, 6}}}}, // eql but fail?
		{"add", add, map[v]v{"a": false}, map[v]v{"a": true}, map[v]v{"a": true}},
		{"add", add, [2]l{l{"a"}, l{false}}, [2]l{l{"a"}, l{[]bool{false, true}}}, [2]l{l{"a"}, l{[]bool{false, true}}}}, // eql but fail?
		{"add", add, map[v]v{"a": 1, "b": fv{2, 3}}, 3, map[v]v{"a": 4, "b": fv{5, 6}}},
		{"add", add, mystruct{}, mystruct{true, 2, nil}, mystruct{true, 2, nil}},
		{"add", add, mystruct{false, 1, []myint{1, 2}}, mystruct{true, 2, []myint{3, 4}}, mystruct{true, 3, []myint{4, 6}}},
		//{"add", add, mystruct{true, 1, []myint{1, 2}}, map[v]v{"B": 3, "I": 1 + 1i, "V": fv{3, 4}}, map[v]v{"B": 4.0, "F": 1, "I": 2 + 1i, "V": fv{4, 6}}}, // eql but fail?
		{"sub", sub, 1, 2, -1},
		{"mul", mul, 2, 3, 6},
		{"div", div, 1.0, 0, math.Inf(1)},
		//{"div", div, c(1, 0), 0, c(math.Inf(1), math.NaN())}, // cannot be compared
		{"mod", mod, 2, fv{1, 2, 3, 4, 5, 6}, fv{1, 0, 1, 0, 1, 0}},
		{"mod", mod, 3, l{1, 2, 3, fv{4, 5}}, l{1, 2, 0, fv{1, 2}}},
		{"mod", mod, c(3, 0), l{1, 2, 3, c(4, 0)}, l{c(1, 0), c(2, 0), c(0, 0), c(1, 0)}},
		{"mkd", mkd, iv{1, 2, 3}, fv{2, 3, 4}, map[v]v{1: 2.0, 2: 3.0, 3: 4.0}},
		{"min", min, 2, 3, 2},
		{"min", min, iv{1, 2, 3}, 2, iv{1, 2, 2}},
		{"max", max, 2, 3, 3},
		{"les", les, 2, 3, 1},
		{"les", les, 2, c(4, 0), c(1, 0)},
		{"les", les, "a", "b", 1.0},
		{"les", les, "a", sv{"b", "a"}, fv{1, 0}},
		{"mor", mor, 2, 3, 0},
		{"mor", mor, 2, c(3, 3), c(0, 0)},
		{"mor", mor, sv{"z", "a"}, sv{"g", "h"}, fv{1, 0}},
		{"eql", eql, fv{1, 2, math.NaN(), math.Inf(1)}, iv{5, 2, 7, 8}, fv{0, 1, 0, 0}},
		{"eql", eql, "a", "a", 1.0},
		{"eql", eql, sv{"a", "b"}, "a", fv{1.0, 0.0}},
		{"pow", pow, fv{2, 2}, fv{0.5, 2}, fv{math.Sqrt2, 4}},
		{"mch", mch, 1, 1, 1.0},
		{"mch", mch, 1, 0, 0.0},
		{"mch", mch, l{}, fv{}, 0.0}, // ()~!0
		{"mch", mch, iv{1, 2}, iv{1, 2}, 1.0},
		{"mch", mch, iv{1, 2}, fv{1, 2}, 0.0},
		{"mch", mch, "a", "a", 1.0},
		{"mch", mch, "alpha", "beta", 0.0},
		{"mch", mch, [2]l{l{"a", "b"}, l{1, 2}}, [2]l{l{"a", "b"}, l{1, 2}}, 1.0},
		{"mch", mch, [2]l{l{"a", "b"}, l{1, 2}}, [2]l{l{"b", "a"}, l{2, 1}}, 0.0},
		{"cat", cat, 1, 2, iv{1, 2}},
		{"cat", cat, 1, iv{2, 3}, iv{1, 2, 3}},
		{"cat", cat, iv{2, 3}, 1, iv{2, 3, 1}},
		{"cat", cat, iv{2, 3}, fv{4, 5}, l{2, 3, 4.0, 5.0}},
		{"cat", cat, [2]l{l{"a", "b"}, l{1, 2}}, [2]l{l{"a", "c"}, l{7, 6}}, [2]l{l{"a", "b", "c"}, l{7, 2, 6}}},
		{"cat", cat, [2]l{{"a"}, l{1}}, 3, l{[2]l{l{"a"}, l{1}}, 3}},
		{"ept", ept, iv{5, 6}, 3, iv{5, 6}},
		{"ept", ept, iv{5, 6}, 5, iv{6}},
		{"ept", ept, iv{5, 6}, l{1, 2, 3, 4, 5, 6, 7}, iv{}},
		{"ept", ept, l{5, 8, 13.0, 12, 8}, l{13.0, 8}, l{5, 12}},
		{"ept", ept, 8.0, l{1.0, 3.0, 5.0}, fv{0, 2, 4, 6, 7}},
		{"tak", tak, 2, iv{1, 2, 3, 4}, iv{1, 2}},
		{"tak", tak, -2.0, l{1, 2, 3, 4}, l{3, 4}},
		{"tak", tak, c(4, 0), l{1, 2.0}, l{1, 2.0, 1, 2.0}},
		{"tak", tak, -3, l{1, 2}, l{2, 1, 2}},
		{"tak", tak, 8, l{1, 2}, l{1, 2, 1, 2, 1, 2, 1, 2}},
		{"tak", tak, -8, l{1, 2}, l{1, 2, 1, 2, 1, 2, 1, 2}},
		{"rsh", rsh, 3, 3, iv{3, 3, 3}},
		{"rsh", rsh, l{2}, l{1, 2}, iv{1, 2}},
		{"rsh", rsh, l{3}, 2.0, fv{2, 2, 2}},
		{"rsh", rsh, l{3}, fv{1, 2}, fv{1, 2, 1}},
		{"rsh", rsh, l{2, 3}, fv{1, 2}, l{fv{1, 2, 1}, fv{2, 1, 2}}},
		{"rsh", rsh, l{3, 2}, l{1, 2, 3}, l{iv{1, 2}, iv{3, 1}, iv{2, 3}}},
		{"rsh", rsh, l{2, 2}, 5, l{iv{5, 5}, iv{5, 5}}},
		{"rsh", rsh, l{2, 0}, l{1, 2, 3, 4}, l{l{}, l{}}},
		{"rsh", rsh, l{2, 3}, l{1, 2, iv{3, 4}}, l{l{1, 2, iv{3, 4}}, l{1, 2, iv{3, 4}}}},
		{"rsh", rsh, l{math.NaN(), 3}, l{0, 1, 2, 3, 4, 5, 6}, l{iv{0, 1, 2}, iv{3, 4, 5}, iv{6}}},
		{"rsh", rsh, l{3, math.NaN()}, l{0, 1, 2, 3, 4, 5, 6}, l{iv{0, 1}, iv{2, 3}, iv{4, 5, 6}}},
		{"drp", drp, 1, l{1, 2, 3}, l{2, 3}},
		{"drp", drp, -1, l{1, 2, 3}, l{1, 2}},
		{"drp", drp, -3, l{1, 2, 3}, l{}},
		{"drp", drp, -4, iv{1, 2, 3}, iv{}},
		{"drp", drp, 4, l{1, 2, 3}, l{}},
		{"drp", drp, 1, fv{1, 2}, fv{2}},
		{"drp", drp, 1, [2]l{l{"a", "b", "c"}, l{1, 2, 3}}, [2]l{l{"b", "c"}, l{2, 3}}},
		{"drp", drp, sv{"a", "c"}, [2]l{l{"a", "b", "c"}, l{1, 2, 3}}, [2]l{l{"b"}, l{2}}},
		{"cut", cut, l{1}, iv{1, 2}, l{iv{2}}},
		{"cut", cut, l{0, 3}, l{0, 1, 2, 3, 4, 5}, l{iv{0, 1, 2}, iv{3, 4, 5}}},
		{"cut", cut, iv{1, 1, 3}, iv{0, 1, 2, 3, 4, 5}, l{l{}, iv{1, 2}, iv{3, 4, 5}}},
		// cst: TODO
		{"fnd", fnd, iv{3}, iv{1, 2, 3}, fv{1, 1, 0}},
		{"fnd", fnd, iv{3}, 3, 0.0},
		{"fnd", fnd, iv{3}, 1, 1.0},
		{"fnd", fnd, iv{3, 4}, [2]l{l{"a", "b"}, l{6, 3}}, [2]l{l{"a", "b"}, l{2.0, 0.0}}},
		{"fnd", fnd, iv{3}, flp(map[v]v{"a": iv{1, 2}, "b": iv{3, 4}}), map[v]v{"a": fv{1, 1}, "b": fv{0, 1}}},
		// cal: TODO
		// bin: TODO
		// rbn: TODO
		// pak: TODO
		// upk: TODO
		// spl: TODO
		// win: TODO
	}
	for _, tc := range testCases {
		r := tc.f(tc.x, tc.y)
		tt(t, tc.r, r, "%s %+v %+v: %+v\n", tc.s, tc.x, tc.y, r)
	}
}

func TestVKt(t *testing.T) {
	testCases := []struct {
		s       string
		f       func(v, v, kt) v
		a       kt
		x, y, r v
	}{
		{"atx", atx, nil, iv{4, 5, 6}, 1, 5},
		{"atx", atx, nil, iv{4, 5, 6}, iv{0, 2}, iv{4, 6}},
		{"atx", atx, nil, map[v]v{"a": -1.0, "b": -2.0}, "b", -2.0},
		{"atx", atx, nil, map[v]v{"a": -1.0, "b": -2.0}, sv{"b", "a"}, fv{-2, -1}},
		// TODO atx var adv
		// TODO atx verb *
	}
	for _, tc := range testCases {
		r := tc.f(tc.x, tc.y, nil)
		tt(t, tc.r, r, "%s %+v %+v: %+v\n", tc.s, tc.x, tc.y, r)
	}
}
func TestRng(t *testing.T) {
	testCases := []struct {
		s string
		f func(v) v
		x v
		n int
		t rT
	}{
		{"rng", rng, 5.0, 5, rTf},
		{"rng", rng, -5.0, 5, rTf},
		{"rng", rng, myfloat(5), 5, rtyp(myfloat(0))},
		{"rng", rng, myfloat(-5), 5, rtyp(myfloat(0))},
		{"rng", rng, c(3, 0), 3, rTz},
	}
	for _, tc := range testCases {
		r := tc.f(tc.x)
		printf("%s %+v: %+v\n", tc.s, tc.x, r)
		if n := ln(r); n != tc.n {
			t.Fatalf("exp len %d got %d", tc.n, n)
		}
		if tp := rtyp(r).Elem(); tp != tc.t {
			t.Fatalf("exp type: %s got %s", tc.t, tp)
		}
	}
}
func TestRnd(t *testing.T) {
	testCases := []struct {
		s    string
		f    func(v, v) v
		x, y v
		n    int
		t    rT
	}{
		{"rnd", rnd, 6, 49, 6, rtyp(0)},
		{"rnd", rnd, -6, 49, 6, rtyp(0)},
		{"rnd", rnd, -6, 7.0, 6, rTf},
	}
	for _, tc := range testCases {
		r := tc.f(tc.x, tc.y)
		printf("%s %+v %+v: %+v\n", tc.s, tc.x, tc.y, r)
		if n := ln(r); n != tc.n {
			t.Fatalf("exp len %d got %d", tc.n, n)
		}
		if tp := rtyp(r).Elem(); tp != tc.t {
			t.Fatalf("exp type: %s got %s", tc.t, tp)
		}
	}
}

func tt(t *testing.T, exp, got v, s string, a ...v) {
	if m, ok := got.(map[v]v); ok {
		delete(m, "_")
		got = m
	}
	printf(s, a...)
	if reflect.DeepEqual(exp, got) == false {
		_fmt.Printf("exp: %#v (%T)\n", exp, exp)
		_fmt.Printf("got: %#v (%T)\n", got, got)
		t.Fatal()
	}
}

// myint is a custom number type that is convertible.
type myint int8

// mynum is a custom type (a string), that implements numeric methods.
type mynum string

func (s mynum) Neg() v { return "-" + s }
func (x mynum) Add(y v, l bool) v {
	if l {
		return mynum(_fmt.Sprintf("%s+%v", x, y))
	}
	return mynum(_fmt.Sprintf("%v+%s", y, x))
}

// mynum is a custom vector type, that implements numeric methods.
type myvec []string

func (s myvec) Neg() v { s[0] = "<" + s[0]; s[len(s)-1] += ">"; return s }
func (s myvec) Add(y v, l bool) v {
	if l {
		s[0] = "~" + s[0]
	}
	s[len(s)-1] += _fmt.Sprintf("+%v", y)
	return s
}

// mystruct is a custom dict type defined as a struct.
type mystruct struct {
	B bool
	F float64
	V []myint
}

// myfloat is a custom float type
type myfloat float64

// mymap is a custom dict type defined as a map.
type mymap map[string]int

func printf(f s, v ...v) {
	if !testing.Verbose() { // temporarily switched
		_fmt.Printf(f, v...)
	}
}
