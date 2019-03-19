package i

import "testing"

func TestKmap(t *testing.T) {
	testCases := []struct {
		x v
		f func(v) v
		r v
	}{
		{l{1.0, 2.0}, func(x v) v { return 1 + x.(f) }, l{2.0, 3.0}},
		{l{1, 2}, func(x v) v { return 1 + x.(i) }, l{2, 3}},
		{iv{1, 2}, func(x v) v { return 1 + x.(i) }, iv{2, 3}},
		{fv{1, 2}, func(x v) v { return 1 + x.(f) }, fv{2, 3}},
		{zv{c(1, 2)}, func(x v) v { return c(3, 4) + x.(z) }, zv{(c(4, 6))}},
		{[]mynum{"a", "b"}, func(x v) v { return string(x.(mynum)) + "/" }, sv{"a/", "b/"}},
	}
	for _, tc := range testCases {
		r := kmap(tc.x, tc.f)
		tt(t, tc.r, r, "kmap %+v (f): %+v\n", tc.x, r)
	}
}
