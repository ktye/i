package i

import "testing"

func TestSn(t *testing.T) {
	testCases := []struct {
		n s
		x v
		c bool
		r zv
	}{
		{"sn", "", false, zv{0}},
		{"sn", "a", false, zv{0}},
		{"sn", sv{"a", "b"}, true, zv{0, 1}},
		{"sn", sv{"b", "a"}, true, zv{1, 0}},
		{"sn", sv{"x", "x", "", "b", "w", "b"}, true, zv{3, 3, 0, 1, 2, 1}},
	}
	for _, tc := range testCases {
		r, c, ok := sn(tc.x)
		if !ok {
			t.Fatal()
		} else if c != tc.c {
			t.Fatalf("%s: scalar/vector test failed", tc.n)
		}
		tt(t, tc.r, r, "sn %+v: %+v\n", tc.x, r)
	}
}
func TestSn2(t *testing.T) {
	testCases := []struct {
		n      s
		x, y   v
		rx, ry v
	}{
		{"sn2", "", "", z0, z0},
		{"sn2", "a", 13, "a", 13},
		{"sn2", "a", "a", z0, z0},
		{"sn2", "a", "b", z0, z1},
		{"sn2", sv{"b", "a"}, "a", zv{1, 0}, zv{0}},
		{"sn2", sv{"x", "x", "", "b", "w", "b"}, sv{"b", "a", "c", "z", "z", "z"}, zv{5, 5, 0, 2, 4, 2}, zv{2, 1, 3, 6, 6, 6}},
	}
	for _, tc := range testCases {
		rx, ry := sn2(tc.x, tc.y)
		tt(t, tc.rx, rx, "sn2x %+v: %+v\n", tc.x, rx)
		tt(t, tc.rx, rx, "sn2y %+v: %+v\n", tc.y, ry)
	}
}
