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
		if len(b) == 0 {
			continue
		}
		fmt.Println(string(b))
		i := bytes.IndexByte(b, '(')
		if i < 0 {
			t.Fatal("no (")
		}
		r := runtest(t, b[:i])
		if r != string(b[i:]) {
			t.Fatalf("got %s", r)
		}
	}
}
func runtest(t *testing.T, b []byte) string {
	ini()
	defer Leak()
	for _, c := range b {
		if Step(uint32(c)) != 0 {
			t.Fatal("early value")
		}
	}
	r := Step(10)
	if r == 0 {
		t.Fatal("zero result")
	}
	return "(" + X(r) + ")"
}
