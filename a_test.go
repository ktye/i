package i

import (
	"testing"
)

func TestEx(t *testing.T) {
	testCases := []struct {
		i v
		t rT
		r v
	}{
		{zv{1, 2, 3}, rtyp([]bool{}), []bool{true, true, true}},
		{iv{1, 2, 3}, rtyp([]uint32{}), []uint32{1, 2, 3}},
		{zv{1, 2 + 1i, 3}, rtyp([]float64{}), fv{1, 2, 3}},
		{[2]l{l{"B", "F", "V"}, l{-0.0, 2.3, zv{3, 4, 5}}}, rtyp(Mystruct{}), Mystruct{false, 2.3, []myint{3, 4, 5}}},
		{
			[2]l{l{"A", "S"}, l{c(3, 0), [2]l{l{"B", "F", "V"}, l{true, 2.3, zv{1, 2, 3}}}}},
			rtyp(nested{}),
			nested{3, Mystruct{true, 2.3, []myint{1, 2, 3}}},
		},
		{
			[2]l{l{"Mystruct"}, l{[2]l{l{"F", "V"}, l{2.3, zv{1, 2, 3}}}}},
			rtyp(embed{}),
			embed{Mystruct: Mystruct{F: 2.3, V: []myint{1, 2, 3}}},
		},
		{
			[2]l{l{"alpha", "beta"}, l{1.0, 2.0}},
			rtyp(mymap{}),
			mymap{"alpha": 1, "beta": 2},
		},
		{
			[2]l{l{"Type", "Lines"}, l{"abc", l{[2]l{l{"A"}, l{1}}, [2]l{l{"A"}, l{2.0}}}}},
			rtyp(Plot{}),
			Plot{Type: "abc", Lines: []Line{Line{1}, Line{2}}},
		},
	}
	for _, tc := range testCases {
		r := ex(tc.i, tc.t)
		tt(t, tc.r, r, "ex %+v %s: %+v\n", tc.i, tc.t, tc.r)
	}
}

type Plot struct {
	Type  string
	Lines []Line
}
type Line struct {
	A int
}
