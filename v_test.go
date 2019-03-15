package i

import (
	_fmt "fmt"
	"math"
	"reflect"
	"testing"
)

func TestMV(t *testing.T) {
	type IV []int
	testCases := []struct {
		s    string
		f    func(interface{}) interface{}
		x, r interface{}
	}{
		// flp TODO
		{"neg", neg, 1.0, -1.0},
		{"neg", neg, complex(1, 2), complex(-1, -2)},
		{"neg", neg, []float64{1, 2, 3}, []float64{-1, -2, -3}},
		{"neg", neg, []complex128{1, complex(2, 3)}, []complex128{-1, complex(-2, -3)}},
		{"neg", neg, []int{1, 2}, []int{-1, -2}},
		{"neg", neg, true, true},
		{"neg", neg, uint16(4), uint16(65532)},
		{"neg", neg, 1, -1},
		{"neg", neg, mynum("33"), mynum("-33")},
		{"neg", neg, []mynum{"a", "b", "c"}, []mynum{"-a", "-b", "-c"}},
		{"neg", neg, myvec{"a", "b"}, myvec{"<a", "b>"}},
		{"neg", neg, map[interface{}]interface{}{"a": []float64{1, 2}}, map[interface{}]interface{}{"a": []float64{-1.0, -2.0}}},
		{"fst", fst, []int{5, 6, 7}, 5},
		{"fst", fst, []float64{5, 6, 7}, 5.0},
		{"fst", fst, []complex128{complex(2, 3), 0, 0}, complex(2, 3)},
		{"fst", fst, IV{5, 6}, 5},
		{"fst", fst, IV{}, nil},
		// {"fst", fst, map[interface{}]interface{}{"a": []int{1, 2}}, map[interface{}]interface{}{"a": []float64{-1, -2}}},
		// fst: TODO dict
		{"sqr", sqr, 4, 2},
		// {"sqr", sqr, -1.0, math.NaN()}, not comparable
		{"sqr", sqr, -7 + 24i, complex(3, 4)},
		// til TODO
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
		{"not", not, complex(math.Inf(1), 0), 0 + 0i},
		{"not", not, complex(math.NaN(), 0), 0 + 0i},

		{"enl", enl, 1.2, []float64{1.2}},
		{"enl", enl, []int{1, 2}, []interface{}{[]int{1, 2}}},
		{"enl", enl, IV{4, 5, 6}, []interface{}{IV{4, 5, 6}}},
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
		if reflect.DeepEqual(r, tc.r) == false {
			t.Fatalf("%s %+v: exp: %+v got %+v\n", tc.s, tc.x, tc.r, r)
		}
		_fmt.Printf("%s %+v: %+v\n", tc.s, tc.x, r)
	}
}

func TestDV(t *testing.T) {
	type IV []int
	testCases := []struct {
		s       string
		f       func(interface{}, interface{}) interface{}
		x, y, r interface{}
	}{
		{"add", add, 1, 2, 3},
		{"add", add, 1, 2.0, 3.0},
		{"add", add, 1.0, uint(2), 3.0},
		{"add", add, []int{1, 2}, 3, []int{4, 5}},
		{"add", add, []int{1, 2}, []int{3}, []int{4, 5}},
		{"add", add, 1, []int{2, 3}, []int{3, 4}},
		{"add", add, []int{1}, []int{2, 3}, []int{3, 4}},
		{"add", add, []int{1, 2, 3}, []int{4, 5, 6}, []int{5, 7, 9}},
		{"add", add, []interface{}{1, 2.0, 3}, 1, []interface{}{2, 3.0, 4}},
		{"add", add, map[interface{}]interface{}{"a": 1, "b": 2.0}, map[interface{}]interface{}{"b": []float64{3, 4}, 9: "x"}, map[interface{}]interface{}{"a": 1, "b": []float64{5, 6}, 9: "x"}},
		{"add", add, map[interface{}]interface{}{"a": 1, "b": 2.0}, []int{3, 4}, map[interface{}]interface{}{"a": 4, "b": 6.0}},
		{"sub", sub, 1, 2, -1},
		{"mul", mul, 2, 3, 6},
		{"div", div, 1.0, 0, math.Inf(1)},
		//{"div", div, complex(1, 0), 0, complex(math.Inf(1), math.NaN())}, // cannot be compared
		// mod TODO
		// mkd TODO
		{"min", min, 2, 3, 2},
		{"min", min, []int{1, 2, 3}, 2, []int{1, 2, 2}},
		{"max", max, 2, 3, 3},
		{"les", les, 2, 3, 1},
		{"les", les, 2, complex(4, 0), complex(1, 0)},
		{"mor", mor, 2, 3, 0},
		{"mor", mor, 2, complex(3, 3), complex(0, 0)},
		{"eql", eql, []float64{1, 2, math.NaN(), math.Inf(1)}, []int{5, 2, 7, 8}, []float64{0, 1, 0, 0}},
		// mch TODO
		{"cat", cat, 1, 2, []int{1, 2}},
		{"cat", cat, 1, []int{2, 3}, []int{1, 2, 3}},
		{"cat", cat, []int{2, 3}, 1, []int{2, 3, 1}},
		{"cat", cat, []int{2, 3}, []float64{4, 5}, []interface{}{2, 3, 4.0, 5.0}},
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
		if reflect.DeepEqual(r, tc.r) == false {
			t.Fatalf("%s %+v %+v: exp: %+v got %+v\n", tc.s, tc.x, tc.y, tc.r, r)
		}
		_fmt.Printf("%s %+v %+v: %+v\n", tc.s, tc.x, tc.y, r)
	}
}

// mynum is a custom type (a string), that implements numeric methods.
type mynum string

func (s mynum) Neg() interface{} { return "-" + s }
func (x mynum) Add(y interface{}, l bool) interface{} {
	if l {
		return mynum(_fmt.Sprintf("%s+%v", x, y))
	}
	return mynum(_fmt.Sprintf("%v+%s", y, x))
}

// mynum is a custom vector type, that implements numeric methods.
type myvec []string

func (s myvec) Neg() interface{} { s[0] = "<" + s[0]; s[len(s)-1] += ">"; return s }
func (s myvec) Add(y interface{}, l bool) interface{} {
	if l {
		s[0] = "~" + s[0]
	}
	s[len(s)-1] += _fmt.Sprintf("+%v", y)
	return s
}
