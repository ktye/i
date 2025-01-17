package main

import (
	"fmt"
	"os"
)

const C = "aaaaaaaaaaoaaaaaaaaaaaaaaaaaaaaaadhddddebcdddjgnggggggggggkbdddddffffffffffffffffffffffffffbmcddifffflfffffffffffffffffffffbdcdaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
const T = "abcdefghijdfekbabcdefghijdfekbabcdefghiddfeebabcdefghijdfeebabcdefghijmfeebabcdennhiddneebabcdeoohidmpeebqqqqqqqsqqqqrqqabcdenghiddneebabcdefohijdfeebllllllllllllllbllllllllllllllbabcdefghijdfeebabcdennhiddneebabcdeoohiddpeebabcdeoohiodoeebqqqqqqqsqqqqrqqqqqqqqqqqqqqqqqabcdefghiddfeeb"

var tkst, sp int
var a, b = -1, 0
var src []byte

func main() {
	if len(os.Args) > 1 {
		src = []byte(os.Args[1])
		//for s := tok(); s != ""; s = tok() {
		//	fmt.Println(s)
		//}
		r := e(t())
		fmt.Println(r)
	}
}

type token struct {
	s string //string
	p int    //src position
	t byte   //type iIfF..
	v bool   //verb
	n int    //arguments
}

func dy(x, o, y []token) []token {
	o[len(o)-1].n = 2
	return append(append(y, x...), o...)
}
func mo(o, x []token) []token {
	o[len(o)-1].n = 1
	return append(x, o...)
}
func at(x, y []token) []token {
	t := x[len(x)-1]
	t.s = "@"
	t.v = true
	t.n = 2
	return append(append(y, x...), t)
}

func t() []token {
	x := token{s: tok()}
	if x.s == "" {
		return nil
	}
	x.p = sp
	c := x.s[0]
	if 2 == cl(c) {
		return nil //)]};
	}
	if c == '(' {
		return e(t())
	}
	if 1 == len(x.s) && (3 == cl(c) || c == '-') {
		x.v = true
	}
	x.t = 'i'
	return []token{x}
}
func e(x []token) []token {
	if x == nil {
		return x
	}
	y := t()
	if y == nil {
		return x
	}
	tx := x[len(x)-1]
	ty := y[len(y)-1]
	if ty.v && !tx.v {
		r := e(t())
		return dy(x, y, r)
	}
	r := e(y)
	if tx.v {
		return mo(x, r)
	} else {
		return at(x, r)
	}
}
func cl(c byte) int { return int(C[c]) - 97 }
func nxt() int {
	if sp == len(src) {
		return -1
	}
	for sp < len(src) {
		tkst = int(T[15*tkst+cl(src[sp])] - 97)
		sp++
		if 11 > tkst {
			return sp - 1
		}
	}
	return -1
}
func tok() string {
	if a < 0 {
		a = nxt()
	} else {
		a = b
	}
	if a < 0 {
		return ""
	}
	for a >= 0 && cl(src[a]) == 0 {
		a = nxt()
		if a < 0 {
			return ""
		}
	}
	if a >= 0 {
		b = nxt()
	}
	if b > 0 {
		return string(src[a:b])
	}
	return string(src[a:])
}
