package j

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"j/jgo"
	"j/jwa"
	"j/x"
	"os"
	"testing"
)

type jj struct{}

func (o jj) J(x uint32) uint32 { return J(x) }
func (o jj) M() []uint32       { return M }

// this tests 3 implementations:
// jj:  native go from ./j.go
// jgo: generated go   by w -go < j.w > jgo/j_.go
// jwa: generated wasm by w     < j.w > jwa/j.wasm interpreted with wagon
func TestJ(t *testing.T) {
	ms := uint32(16)
	js := []x.J{jj{}, jgo.New(), jwa.New()}
	all := func(js []x.J, x uint32) {
		for _, j := range js {
			j.J(x)
		}
	}
	all(js, ms)
	b, err := ioutil.ReadFile("t")
	if err != nil {
		t.Fatal(err)
	}
	v := bytes.Split(b, []byte{10})
	for _, b := range v {
		if len(b) == 0 || b[0] == '(' { // (section)
			continue
		}
		if string(b) == `\` {
			break
		}
		i := bytes.IndexByte(b, '(')
		if i < 0 {
			t.Fatal("no (")
		}
		e := bytes.IndexByte(b, ')') // first) (other comments possible)
		if e < 0 {
			t.Fatal("no )")
		}
		in := b[:i]
		exp := string(b[i : 1+e])

		os.Stdout.Write(in)
		os.Stdout.WriteString(exp)
		for i, j := range js {
			runtest(t, j, in, exp)
			if i > 0 {
				memcompare(t, js[0], js[i])
			}
		}
		os.Stdout.WriteString("\n")
	}
}
func memcompare(t *testing.T, aj, bj x.J) {
	a, b := aj.M(), bj.M()
	if len(a) != len(b) {
		t.Fatalf("memory#: %d != %d", len(a), len(b))
	}
	for i, u := range a {
		if u != b[i] {
			fmt.Println()
			x.Dump(a, 100)
			x.Dump(b, 100)
			t.Fatalf("mem differs at %d (%x): %x %x\n", i, i, u, b[i])
		}
	}
}
func runtest(t *testing.T, j x.J, b []byte, exp string) {
	for _, c := range b {
		if j.J(uint32(c)) != 0 {
			t.Fatal("early value")
		}
	}
	r := j.J(10)
	if r == 0 {
		t.Fatal("zero")
	}
	m := j.M()
	s := x.X(m, m[1])
	s = "(" + s[1:len(s)-1] + ")"
	exp = exp[:len(exp)-1] + " [])" // add empty stack
	if exp == "( [])" {
		exp = "([])"
	}
	if s == exp {
		os.Stdout.WriteString(" ok")
	} else {
		t.Fatalf("got %q\nexp %q\n", s, exp)
	}
	x.Leak(j)
	n := ln(m, m[1]) - 1 // last is new parse list
	for i := uint32(0); i < n; i++ {
		j.J('_')
	}
	j.J(10)
}
func ln(m []uint32, x uint32) uint32 { return m[1+(x>>2)] }
