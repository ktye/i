package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestB(t *testing.T) {
	testCases := []struct {
		sig string
		b   string
		e   string
	}{
		{"I:II", "x+y", "20002001 6a"},
		{"F:FF", "(x*y)", "20002001 a2"},
		{"F:FF", "x-y", "20002001 a1"},
		{"F:FF", "3.*x+y", "44 0000000000000840 20002001 a0 a2"},
	}
	for n, tc := range testCases {
		f := newfn(tc.sig, tc.b)
		e := f.parse()
		b := string(hex(e.bytes()))
		s := trim(tc.e)
		if b != s {
			t.Fatalf("%d: expected/got:\n%s\n%s", n+1, s, b)
		}
		fmt.Println(b)
	}
}

/*
func TestW(t *testing.T) {
	testCases := [][2]string{
		{"add:I:II::2000 2001 6a\n", "0061736d0100000001070160027f7f017f0302010005030100010707010361646400000a09010700200020016a0b"},
	}
	for n, tc := range testCases {
		i, e := tc[0], tc[1]
		o := string(hex(run(bytes.NewReader([]c(i))).wasm()))
		if o != e {
			t.Fatalf("%d: expected/got:\n%q\n%q\n", n+1, e, o)
		}
	}
}
*/
func hex(a []c) []c {
	var r bytes.Buffer
	for _, b := range a {
		hi, lo := hxb(b)
		r.WriteByte(hi)
		r.WriteByte(lo)
	}
	return r.Bytes()
}
func newfn(sig string, body string) fn {
	var buf bytes.Buffer
	buf.WriteString(body)
	buf.WriteByte('}')
	v := strings.Split(sig, ":")
	if len(v) != 2 {
		panic("signature")
	}
	f := fn{src: [2]int{1, 0}, Buffer: buf}
	f.t = typs[v[0][0]]
	for _, c := range v[1] {
		f.args = append(f.args, typs[byte(c)])
	}
	return f
}
func trim(s string) string { return strings.Replace(s, " ", "", -1) }
