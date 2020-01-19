package main

import (
	"bytes"
	"testing"
)

func TestW(t *testing.T) {
	testCases := [][2]string{
		{"add:I:II::2000 2001 6a\n", "0061736d0100000001070160027f7f017f0302010005030100010707010361646400000a09010700200020016a0b"},
	}
	for n, tc := range testCases {
		i, e := tc[0], tc[1]
		o := string(hex(run(bytes.NewReader([]c(i)))))
		if o != e {
			t.Fatalf("%d: expected/got:\n%q\n%q\n", n, e, o)
		}
	}
}
func hex(a []c) []c {
	var r bytes.Buffer
	for _, b := range a {
		hi, lo := hxb(b)
		r.WriteByte(hi)
		r.WriteByte(lo)
	}
	return r.Bytes()
}
