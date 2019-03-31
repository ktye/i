package i

import "testing"

func TestKmap(t *testing.T) {
	testCases := []struct {
		x v
		f func(v, int) v
		r v
	}{
		{l{1.0, 2.0}, func(x v, j int) v { return 1 + x.(f) }, l{2.0, 3.0}},
		{l{1, 2}, func(x v, j int) v { return 1 + x.(i) }, l{2, 3}},
		{iv{1, 2}, func(x v, j int) v { return 1 + x.(i) }, iv{2, 3}},
		{fv{1, 2}, func(x v, j int) v { return 1 + x.(f) }, fv{2, 3}},
		{zv{c(1, 2)}, func(x v, j int) v { return c(3, 4) + x.(z) }, zv{(c(4, 6))}},
		{l{1, 2}, func(x v, j int) v { return j + x.(i) }, l{1, 3}},
	}
	for _, tc := range testCases {
		r := kmap(tc.x, tc.f)
		tt(t, tc.r, r, "kmap %+v (f): %+v\n", tc.x, r)
	}
}
