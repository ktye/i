// +build ignore

// standalone k.w parser (w/o vectors)
// E:E;e|e e:nve|te| t:n|v|{E} v:tA|V n:t[E]|(E)|N
//
// go run p.go      /run all tests
// go run p.go -5   /run #5 traced
// go run p.go -c   /generate char map
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var tests = [][2]string{
	{"3", "3"},
	{"(1;2)", "(1;2)"},
	{"+", "+"},
	{`""`, `""`},
	{`"abc+-*;"`, `"abc+-*;"`},
	{"`abc", "(`abc)"},
	{"`\"abc\"", "(`abc)"},
	{"aBc", "`aBc"},
	{"A5", "`A5"},
	{"-:", "-:"},
	{"+/", "(/;+)"},
	{"*3", "(*;3)"},
	{"3*", "(*;3;)"},
	{"%*3", "(%;(*;3))"},
	{"1+2", "(+;1;2)"},
	{"(a;b;`c)", "(`a;`b;(`c))"},
	{"(a;(1;2);3)", "(`a;(1;2);3)"},
	{"x;y", "(`x;`y)"},
	{"(x;y)", "(`x;`y)"},
	{"f[x;y]", "(`f;`x;`y)"},
	{"f(a;b)", "(`f;(`a;`b))"},
	{"x+y", "(+;`x;`y)"},
	{"x+*y", "(+;`x;(*;`y))"},
	{"1+3*x", "(+;1;(*;3;`x))"},
	{"(+x)%y", "(%;(+;`x);`y)"},
	{"(+/x)%#x", "(%;((/;+);`x);(#;`x))"},
	{"x+m[*i]/y", "(+;`x;((/;(`m;(*;`i)));`y))"},
	{"#'=x", "((';#);(=;`x))"},
	{"x~|x", "(~;`x;(|;`x))"},
	{"(x+y)", "(+;`x;`y)"},
	{"{x;y}", "{`x;`y}"},
	{"{x+y}", "{+;`x;`y}"},
	{"()", ""},
	{"(1)", "1"},
	{"{}", "{}"},
	{"{1}", "{1}"},
	{"{x+y}[1;2]", "({+;`x;`y};1;2)"},
	{"1{x+y}2", "({+;`x;`y};1;2)"},
	{"x:3", "(:;`x;3)"},
	{"x::3", "(::;`x;3)"},
	{"x+:3", "(+:;`x;3)"},
	{"x+:", "(+:;`x;)"},
	{"x[1]:5", "(:;(`x;1);5)"},
	{"x[1;2]+:", "(+:;(`x;1;2);)"},
	{"x[1;2]+:3", "(+:;(`x;1;2);3)"},
	{"x[1]:5", "(:;(`x;1);5)"},
	{"+-", "(+;-)"},
	{"(+-)x", "((+;-);`x)"},
	{"+-*", "(+;(-;*))"},
	{"(+-*)3", "((+;(-;*));3)"},
	{"(+-*)[1;3]", "((+;(-;*));1;3)"},
	{"+-*%", "(+;(-;(*;%)))"},
	{"(+/2*)5", "(((/;+);(*;2;));5)"},
}

func main() {
	if len(os.Args) == 2 {
		if os.Args[1] == "-c" {
			cmap() // print cmap for k.w
			return
		}
		if i, err := strconv.Atoi(os.Args[1]); err == nil {
			if os.Args[1][0] == '-' {
				i = -i
				trace = fmt.Printf
			}
			run(i)
			return
		}
	}
	for i := range tests {
		run(i)
	}
}

// simplified K type system
type K interface{} // any K value
type L []K         // list (a;b;..)
type λ []K         // {x+y}
type C []byte      // "abc"
type S string      // `abc
type I int         // 3
type V string      // verb + -:
type A string      // adverb / ':

var s C   // parser input
var p int // current index in s

func run(i int) {
	x := tests[i]
	a, b := x[0], x[1]
	s, p = C(x[0]), 0
	r := E()
	g := o(r)
	if len(r) == 1 {
		g = o(r[0])
	}
	fmt.Printf("%-3d %-14s / %s\n", i, a, g)
	if g != b {
		fmt.Printf("expected: %s\ngot:      %s\n", b, g)
		os.Exit(1)
	}
}
func e(x K) (r K) {
	defer func() { trace("e(%s)->%s\n", o(x), o(r)) }()
	if x == nil || w() || is(s[p], TE) {
		return x
	}
	y := t() // nil?
	if isverb(y) && !isverb(x) {
		return L{y, x, e(t())}
	}
	return L{x, e(y)}
}
func isverb(x K) bool {
	if _, ok := x.(V); ok {
		return true // verb
	}
	if _, ok := x.(λ); ok {
		return true // lambda (not in other K)
	}
	if l, ok := x.(L); ok && len(l) == 2 {
		if _, ok := l[0].(A); ok {
			return true // adverb derived
		}
	}
	return false // noun
}
func t() (r K) {
	defer func() { trace("t->%s\n", o(r)) }()
	x := tok()
	if x == nil {
		if p == len(s) {
			return nil
		}
		if s[p] == '(' || s[p] == '{' {
			p++
			l := E()
			x = l
			if len(l) == 1 {
				x = l[0]
			}
			if s[p-1] == '}' {
				if l, ok := x.(L); ok {
					x = λ(l)
				} else {
					x = λ{x}
				}
			}
		}
	}
	for {
		if p < len(s) && (is(s[p], AD) || s[p] == '[') {
			if s[p] == '[' {
				p++
				x = append(L{x}, E()...) //prepend x to E(list)
			} else {
				x = L{adverb(), x}
			}
		} else {
			return x
		}
	}
}
func E() (r L) {
	defer func() { trace("E->%s\n", o(r)) }()
	r = L{e(t())}
	for {
		if w() || s[p] != ';' { // or newline
			if p < len(s) {
				p++
			}
			return r
		} else {
			p++
			r = append(r, e(t()))
		}
	}
}
func w() bool {
	for {
		if p == len(s) { // EOF
			return true
		}
		b := s[p]
		if b < '!' || b > 126 { // whitespace
			p++
		} else {
			break
		}
	}
	return false
}
func tok() (r K) { // next token
	if w() {
		return nil
	}
	if is(s[p], TE) {
		return nil
	}
	if chars(&r) || number(&r) || name(&r) || symbol(&r) || verb(&r) {
		trace("Token: %s %T\n", o(r), r)
	}
	return r
}
func chars(r *K) bool {
	if s[p] != '"' {
		return false
	}
	a := p + 1
	for {
		p++
		if p == len(s) {
			return xx("chars: unterminated")
		} else if s[p] == '"' {
			*r = s[a:p]
			p++
			return true
		}
	}
}
func number(r *K) (v bool) {
	atoi := func(b []byte) I { x, _ := strconv.Atoi(string(b)); return I(x) }
	for i := 0; ; i++ {
		if p+i == len(s) || !is(s[p+i], NM) {
			if i > 0 {
				*r = atoi(s[p : p+i])
				p += i
				return true
			}
			return false
		}
	}
}
func name(r *K) (v bool) { // abc A3
	if !is(s[p], az+AZ) {
		return false
	}
	for i := 0; ; i++ {
		if p+i == len(s) || !is(s[p+i], az+AZ+NM) {
			*r = S(s[p : p+i])
			p += i
			return true
		}
	}
}
func symbol(r *K) (v bool) { // `abc `"abc"
	if s[p] != '`' {
		return false
	}
	p++
	if chars(r) {
		*r = S((*r).(C))
	} else if name(r) == false {
		*r = L{S("")}
	}
	*r = L{*r} // enlist
	return true
}
func verb(r *K) (v bool) {
	if !is(s[p], VB) {
		return false
	}
	x := string(s[p])
	p++
	if p < len(s) && s[p] == ':' {
		x += ":"
		p++
	}
	*r = V(x)
	return true
}
func adverb() A {
	x := string(s[p])
	p++
	if p < len(s) && s[p] == ':' {
		x += ":"
		p++
	}
	return A(x)
}
func xx(e S) bool {
	sp := ""
	if p > 1 {
		sp = strings.Repeat(" ", p-1)
	}
	fmt.Printf("%s\n%s^%s\n", s, sp, e)
	os.Exit(1)
	return false
}

func o(x K) string {
	switch u := x.(type) {
	case L:
		v := make([]string, len(u))
		for i, e := range u {
			v[i] = o(e)
		}
		return "(" + strings.Join(v, ";") + ")"
	case λ:
		s := o(L(u))
		return "{" + s[1:len(s)-1] + "}"
	case C:
		return `"` + string(u) + `"`
	case S:
		return "`" + string(u)
	case I:
		return strconv.Itoa(int(u))
	case V:
		return string(u)
	case A:
		return string(u)
	case nil:
		return ""
	default:
		return fmt.Sprintf("[unknown K type: %T]", x)
	}
}

var trace = func(s string, v ...interface{}) (int, error) { return 0, nil }

func is(x, m byte) bool { return (m & c_[x]) != 0 }

const (
	az = 1 << iota //  1 a-z
	AZ             //  2 A-Z
	NM             //  4 numbers     0123456789
	VB             //  8 verbs       :+-*%!&|<>=~,^#_$?@.
	AD             // 16 adverbs     '/\
	TE             // 32 terminators ;)]} (space)
	NW             // 64 nonwhite    33..126
)

var c_ [128]byte // class map (constant)

func cmap() {
	fmt.Printf("x")
	for _, c := range c_[32:] {
		fmt.Printf("%02x", c)
	}
	fmt.Printf("\n")
}
func init() {
	m := func(s string, b byte) {
		for i := range s {
			c_[s[i]] |= b
		}
	}
	m("abcdefghijklmnopqrstuvwxyz", az)
	m("ABCDEFGHIJKLMNOPQRSTUVWXYZ", AZ)
	m("0123456789", NM)
	m(":+-*%!&|<>=~,^#_$?@.", VB)
	m("'/\\", AD)
	m(";)]} ", TE)
	for i := 33; i < 127; i++ {
		c_[i] |= NW
	}
}
