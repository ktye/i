package j

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestJ(t *testing.T) {
	b, err := ioutil.ReadFile("t")
	if err != nil {
		t.Fatal(err)
	}
	v := bytes.Split(b, []byte{10})
	for _, b := range v {
		if len(b) == 0 || b[0] == '(' { // (section)
			continue
		}
		fmt.Println(string(b))
		i := bytes.IndexByte(b, '(')
		if i < 0 {
			t.Fatal("no (")
		}
		e := bytes.IndexByte(b, ')') // first) (other comments possible)
		if e < 0 {
			t.Fatal("no )")
		}
		exp := string(b[i : 1+e])
		r := runtest(t, b[:i])
		if r != exp {
			t.Fatalf("got %q\nexp %q\n", r, exp)
		}
	}
}
func runtest(t *testing.T, b []byte) string {
	for _, c := range b {
		if J(uint32(c)) != 0 {
			t.Fatal("early value")
		}
	}
	r := J(10)
	if r == 0 {
		t.Fatal("zero result")
	}
	s := X(I(4))
	s = "(" + s[1:len(s)-1] + ")"
	Leak()
	return s
}
