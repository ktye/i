package main

import (
	"bytes"
	. "github.com/ktye/wg/module"
	"os"
	"strconv"
	"strings"
	"testing"
)

// var tot, top int32

func Test1K(t *testing.T) {
	k1()
	T(tostring(w(-2)), "-2")
	T(tostring(til(w(5))),"(0 1 2 3 4)")
	T(tostring(cnt(til(w(5)))),"5")
	T(tostring(cnt((w(3)))),"1")

	//readtests("readme")
}
func T(a, b string) {
	if a != b {
		panic("got:" + a + "\nnot:" + b)
	}
	println(a)
}
func tostring(x int32) string {
	if x&1 != 0 {
		return strconv.Itoa(int(x >> 1))
	} else {
		xn := n(x)
		if xn < 0 || xn > 20 {
			panic("xn")
		}
		u := make([]string, xn)
		for i := int32(0); i < xn; i++ {
			u[i] = tostring(I32(x + 4*i))
		}
		return "(" + strings.Join(u, " ") + ")"
	}
}
func readtests(file string) (r [][2]string) {
	b, e := os.ReadFile(file)
	fatal(e)
	v := bytes.Split(b, []byte("\n"))
	for i := range v {
		s := string(v[i])
		if len(s) > 0 && s[0] == ' ' {
			r = append(r, readcase(s[1:]))
		}
	}
	return r
}
func readcase(s string) (r [2]string) {
	var o bool
	r[0], r[1], o = strings.Cut(s, " /")
	if !o {
		panic("testcase: " + s)
	}
	return
}
func fatal(e error) {
	if e != nil {
		panic(e)
	}
}
