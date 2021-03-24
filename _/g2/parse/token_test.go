package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

type testCase struct {
	in, exp string
}

func tcFile(t *testing.T, f string) (tc []testCase) {
	b, e := ioutil.ReadFile(f)
	if e != nil {
		t.Fatal(e)
	}
	v := bytes.Split(b, []byte{10})
	for _, u := range v {
		s := string(u)
		if strings.HasPrefix(s, "// ") {
			s = s[3:]
			arrow := " → "
			if idx := strings.Index(s, arrow); idx > 0 {
				tc = append(tc, testCase{s[:idx], s[idx+len(arrow):]})
			}
		}
	}
	if tc == nil {
		t.Fatalf("%s has not test cases", f)
	}
	return tc
}
func TestToken(t *testing.T) {
	for _, tc := range tcFile(t, "token.go") {
		tok, _ := token([]byte(tc.in))
		fmt.Printf("%s → ", tc.in)
		var r string
		if len(tok) == 1 {
			r = fmt.Sprint(tok[0])
		} else {
			r = fmt.Sprint(tok)
		}
		if r != tc.exp {
			t.Fatalf("expected %s got %s", tc.exp, r)
		}
		fmt.Println(r)
	}
}

func TestFloat(t *testing.T) {
	testCases := []struct {
		s string
		n int
	}{
		{"1.23", 4},
		{"0.a", 2},
		{"0.", 2},
		{".0", 2},
		{".", 0},
		{"a123", 0},
		{"1e2", 3},
		{"1ea", 2},
		{"1a30", 1},
	}
	for _, tc := range testCases {
		fmt.Println(tc.s)
		if n := sFloat([]byte(tc.s)); n != tc.n {
			t.Fatalf("%s: expected %d got %d", tc.s, tc.n, n)
		}
	}
}
