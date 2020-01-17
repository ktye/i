package main

import (
	"testing"
)

func TestW(t *testing.T) {
	testCases := [][2]string{
		{" .626363 \n636261\n\n", "616263636261"},
	}
	for _, tc := range testCases {
		in := tc[0]
		var r []c
		v := strings.Split(in, "\n")
		for _, s := range v {
			r = append(r, run([]c(s)))
		}
		if s := string(hex(r)); s != tc[1] {
			t.Fatalf("expected %q got %q\n", tc[1], s)
		}
	}
}
func hex(b []c) (r []c) {
	r = make(2 * len(b))
	for i, z := range b {
		r[2*i], r[2*i+1] = hxb(z)
	}
	return r
}
