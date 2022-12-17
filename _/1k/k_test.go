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
	readtests("readme")
	k1()
}
func tostring(x int32) string {
	if x&1 != 0 {
		xn := n(x)
		if xn < 0 || xn > 20 {
			panic("xn")
		}
		u := make([]string, xn)
		for i := int32(0); i < xn; i++ {
			u[i] = tostring(I32(v(x) + 4*i))
		}
		return "(" + strings.Join(u, " ") + ")"
	} else {
		return strconv.Itoa(int(x >> 1))
	}
}
func readtests(file string) (r [][2]string) {
	b, e := os.ReadFile(file)
	fatal(e)
	v := bytes.Split(b, []byte("\n"))
	for i := range v {
		s := string(v[i])
		if s[i] == ' ' {
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
