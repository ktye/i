package main

import (
	"bytes"
	"fmt"
	"testing"
)

var save []byte

func newtest() {
	rand_ = 1592653589
	if save == nil {
		kinit()
		save = make([]byte, len(Bytes))
		copy(save, Bytes)
	} else {
		Bytes = make([]byte, len(save))
		copy(Bytes, save)
		pp, pe, sp = 0, 0, 256
	}
}

func TestCsv(t *testing.T) {
	testCases := []struct {
		fm, fi, exp, fstr string
	}{
		{"zz", file1, k1, ",0:22-2:22-4"},
		{"", file2, k2a, ";0:18-1:18-2:18-3:18-4"},
		{";i2ff", file2, k2b, ";0:19-2:21-3:21-4"},
		{";2h0i3f", file2, k2c, ";0:19-3:21-4"},
		{" csi", file3, k3, " 0:18-1:20-2:19-3"},
	}
	for i, tc := range testCases {
		newtest()
		var buf bytes.Buffer
		x, f := kcsv(tc.fm, []byte(tc.fi), "")
		if s := f.fstr(); s != tc.fstr {
			t.Fatalf("tc %d format: exp %s got %s", i, tc.fstr, s)
		}
		wCK(&buf, Kst(x))
		got := string(buf.Bytes())
		exp := tc.exp
		if got != exp {
			t.Fatalf("tc %d\nexp\n%s\ngot\n%s", i, exp, got)
		}
	}
}

func (f *format) fstr() string {
	s := string(f.s)
	for i, t := range f.t {
		s += fmt.Sprintf("%d:%d-", f.i[i], t)
	}
	return s + fmt.Sprintf("%d", f.columns)
}

const file1 = `
1,90,2,270
2,90,3,180
`
const k1 = `(1a90 2a90;2a270 3a180)`

const file2 = `
12;"alpha";2,5;4.5
23;"beta"; 3,0;5.5
34;"gamma";4.0;6.5`

const k2a = `(("12";"23";"34");("alpha";"beta";"gamma");("2,5";"3,0";"4.0");("4.5";"5.5";"6.5"))`
const k2b = `(12 23 34;2.5 3. 4.;4.5 5.5 6.5)`
const k2c = `(,34;,6.5)`

const file3 = `
alpha  hans   1 
beta   peter  2
gamma  jochen 3
`
const k3 = "((\"alpha\";\"beta\";\"gamma\");`hans`peter`jochen;1 2 3)"
