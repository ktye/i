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
		// flp TODO
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
		// fst: TODO dict
		{"sqr", sqr, 4, 2},
		// {"sqr", sqr, -1.0, math.NaN()}, not comparable
		{"sqr", sqr, -7 + 24i, c(3, 4)},
		{"til", til, 3.0, fv{0, 1, 2}},
		{"til", til, 3, iv{0, 1, 2}},
		// TODO til l (odometer)
		// odo TODO
		// wer TODO
		// rev TODO
		// asc TODO
		// dsc TODO
		// eye TODO
		// grp TODO
		{"not", not, 1, 0},
		{"not", not, 1 + 2i, 0 + 0i},
		{"not", not, 0 + 0i, 1 + 0i},
		{"not", not, c(math.Inf(1), 0), 0 + 0i},
		{"not", not, c(math.NaN(), 0), 0 + 0i},

		{"enl", enl, 1.2, fv{1.2}},
		{"enl", enl, iv{1, 2}, l{iv{1, 2}}},
		{"enl", enl, IV{4, 5, 6}, l{IV{4, 5, 6}}},
		// is0 TODO
		// cnt TODO
		{"flr", flr, 3.5, 3.0},
		// fmt TODO
		// fng TODO
		// unq TODO
		// evl TODO
	}

	for _, tc := range testCases {
		r := tc.f(tc.x)
		if m, ok := r.(map[v]v); ok {
			delete(m, "_")
			r = m
		}
		printf("%s %+v: %+v\n", tc.s, tc.x, r)
		if reflect.DeepEqual(r, tc.r) == false {
			printf("exp:\n%+#v\ngot:\n%+#v\n", tc.r, r)
			t.Fatalf("%s %+v: exp: %+v got %+v\n", tc.s, tc.x, tc.r, r)
		}
	}
}

func TestDV(t *testing.T) {
	type IV []int
	testCases := []struct {
		s       string
		f       func(v, v) v
		x, y, r v
	}{
		{"add", add, 1, 2, 3},
		{"add", add, 1, 2.0, 3.0},
		{"add", add, 1.0, uint(2), 3.0},
		{"add", add, iv{1, 2}, 3, iv{4, 5}},
		{"add", add, iv{1, 2}, iv{3}, iv{4, 5}},
		{"add", add, 1, iv{2, 3}, iv{3, 4}},
		{"add", add, iv{1}, iv{2, 3}, iv{3, 4}},
		{"add", add, iv{1, 2, 3}, iv{4, 5, 6}, iv{5, 7, 9}},
		{"add", add, l{1, 2.0, 3}, 1, l{2, 3.0, 4}},
		//{"add", add,
		//	map[v]v{"a": 1, "b": 2.0},
		//	map[v]v{"b": fv{3, 4}},
		//	map[v]v{"a": 1, "b": fv{5.0, 6.0}}}, // eql but fail?
		{"add", add, map[v]v{"a": false}, map[v]v{"a": true}, map[v]v{"a": true}},
		//{"add", add, map[v]v{"a": []bool{false}}, map[v]v{"a": []bool{false, true}}, map[v]v{"a": []bool{false, true}}}, // eql but fail?
		{"add", add, map[v]v{"a": 1, "b": fv{2, 3}}, 3, map[v]v{"a": 4, "b": fv{5, 6}}},
		{"add", add, mystruct{}, mystruct{true, 2, nil}, mystruct{true, 2, nil}},
		{"add", add, mystruct{false, 1, []myint{1, 2}}, mystruct{true, 2, []myint{3, 4}}, mystruct{true, 3, []myint{4, 6}}},
		// {"add", add, mystruct{true, 1, []myint{1, 2}}, map[v]v{"B": 3, "I": 1 + 1i, "V": fv{3, 4}}, map[v]v{"B": 4.0, "F": 1, "I": 2 + 1i, "V": fv{4, 6}}}, // eql but fail?
		{"sub", sub, 1, 2, -1},
		{"mul", mul, 2, 3, 6},
		{"div", div, 1.0, 0, math.Inf(1)},
		//{"div", div, c(1, 0), 0, c(math.Inf(1), math.NaN())}, // cannot be compared
		// mod TODO
		// mkd TODO
		{"min", min, 2, 3, 2},
		{"min", min, iv{1, 2, 3}, 2, iv{1, 2, 2}},
		{"max", max, 2, 3, 3},
		{"les", les, 2, 3, 1},
		{"les", les, 2, c(4, 0), c(1, 0)},
		{"mor", mor, 2, 3, 0},
		{"mor", mor, 2, c(3, 3), c(0, 0)},
		{"eql", eql, fv{1, 2, math.NaN(), math.Inf(1)}, iv{5, 2, 7, 8}, fv{0, 1, 0, 0}},
		// mch TODO
		{"cat", cat, 1, 2, iv{1, 2}},
		{"cat", cat, 1, iv{2, 3}, iv{1, 2, 3}},
		{"cat", cat, iv{2, 3}, 1, iv{2, 3, 1}},
		{"cat", cat, iv{2, 3}, fv{4, 5}, l{2, 3, 4.0, 5.0}},
		// cat: TODO dict
		// ept: TODO
		// tak: TODO
		// rsh: TODO
		// fil: TODO
		// drp: TODO
		// cut: TODO
		// cst: TODO
		// rnd: TODO
		// fnd: TODO
		// pik: TODO
		// rfd: TODO
		// atx: TODO
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
		if m, ok := r.(map[v]v); ok {
			delete(m, "_")
			r = m
		}
		printf("%s %+v %+v: %+v\n", tc.s, tc.x, tc.y, r)
		if reflect.DeepEqual(r, tc.r) == false {
			printf("exp:\n%+#v\ngot:\n%+#v\n", tc.r, r)
			t.Fatalf("%s %+v %+v: exp: %+v got %+v\n", tc.s, tc.x, tc.y, tc.r, r)
		}
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

// mymap is a custom dict type defined as a map.
type mymap map[string]int

func printf(v ...v) {
	if !testing.Verbose() { // temporarily switched
		s := v[0].(string)
		_fmt.Printf(s, v[1:]...)
	}
}
