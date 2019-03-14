package i

import (
	_fmt "fmt"
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
		// idn TODO
		// flp TODO
		// neg TODO
		{"neg", neg, 1.0, -1.0},
		{"neg", neg, complex(1, 2), complex(-1, -2)},
		{"neg", neg, []float64{1, 2, 3}, []float64{-1, -2, -3}},
		{"neg", neg, []complex128{1, complex(2, 3)}, []complex128{-1, complex(-2, -3)}},
		{"neg", neg, true, -1.0},
		{"neg", neg, uint16(4), -4.0},
		{"neg", neg, 1, -1.0},
		{"neg", neg, mynum("33"), mynum("-33")},
		{"neg", neg, []mynum{"a", "b", "c"}, []interface{}{mynum("-a"), mynum("-b"), mynum("-c")}},
		{"fst", fst, []int{5, 6, 7}, 5},
		{"fst", fst, []float64{5, 6, 7}, 5.0},
		{"fst", fst, []complex128{complex(2, 3), 0, 0}, complex(2, 3)},
		{"fst", fst, IV{5, 6}, 5},
		{"fst", fst, IV{}, nil},
		// fst: TODO dict
		// sqr: TODO
		// iot TODO
		// odo TODO
		// wer TODO
		// rev TODO
		// asc TODO
		// dsc TODO
		// eye TODO
		// grp TODO
		// not TODO
		{"enl", enl, 1.2, []float64{1.2}},
		{"enl", enl, []int{1, 2}, []interface{}{[]int{1, 2}}},
		{"enl", enl, IV{4, 5, 6}, []interface{}{IV{4, 5, 6}}},
		// is0 TODO
		// cnt TODO
		// flr TODO
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
		// add TODO
		// sub TODO
		// mul TODO
		// div TODO
		// mod TODO
		// mkd TODO
		// min TODO
		// max TODO
		// les TODO
		// mor TODO
		// eql TODO
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
