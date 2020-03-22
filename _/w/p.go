// +build ignore

// E:E;e|e e:nve|te| t:n|v v:tA|V n:t[E]|(E)|{E}|N
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
	{"-:", "-:"},
	{"+/", "(/;+)"},
	{"*3", "(*;3)"},
	{"3*", "(*;3;)"},
	{"%*3", "(%;(*;3))"},
	{"1+2", "(+;1;2)"},
	{"(a;b;`c)", "(`a;`b;,`c)"},
	{"(a;(1;2);3)", "(`a;(1;2);3)"},
	{"x;y", "(`x;`y)"},
	{"(x;y)", "(`x;`y)"},
	{"f[x;y]", "(`f;`x;`y)"},
	{"x+y", "(+;`x;`y)"},
	{"x+*y", "(+;`x;(*;`y))"},
	{"1+3*x", "(+;1;(*;3;`x))"},
	{"(+x)%y", "(%;(+;`x);`y)"},
	{"(+/x)%#x", "(%;((/;+);`x);(#;`x))"},
	{"x+m[*i]/y", "(+;`x;((/;(`m;(*;`i)));`y))"},
	{"#'=x", "((';#);(=;`x))"},
	{"x~|x", "(~;`x;(|;`x))"},
}

func main() {
	if len(os.Args) == 2 {
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

// K type system
type K interface{} // any K value
type L []K         // list (a;b;..)
type F []K         // lambda {x+y}
type C []byte      // "abc"
type S string      // `abc (a-z)
type I = int       // 3
type V string      // verb + -:
type A string      // adverb / ':

// classes a v ; (lookup table)
// a: ['/\                  adverbs and [
// v: :+-*%!&|<>=~,^#_$?@.  verbs
// ;: ;)]                   terminators (including space)
const c = "                                ;v vvvva ;vvvvva          v;vvvvv                          aa;vv                            v v"

var s C   // parser input
var p int // current index in s

func run(i int) {
	x := tests[i]
	a, b := x[0], x[1]
	s, p = C(x[0]), 0
	r := E(0)
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
	if x == nil || w() || c[s[p]] == ';' {
		return x
	}
	y := t() // nil?
	if isverb(y) && !isverb(x) {
		r = L{y, x, e(t())}
		return r
	}
	l := L{x, e(y)}
	return l
}
func isverb(x K) bool {
	if _, ok := x.(V); ok {
		return true // verb
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
		if s[p] == '(' {
			p++
			l := E(')')
			x = l
			if len(l) == 1 {
				x = l[0]
			}
		} else if s[p] == '{' {
			p++
			x = F(E('}'))
		}
	}
	for {
		if p < len(s) && c[s[p]] == 'a' {
			if s[p] == '[' {
				p++
				x = append(L{x}, E(']')...) //prepend x to E(list)
			} else {
				x = L{adverb(), x}
			}
		} else {
			return x
		}
	}
}
func E(c byte) (r L) {
	defer func() { trace("E->%s\n", o(r)) }()
	r = L{e(t())}
	for {
		if w() || s[p] != ';' { // or newline
			if c != 0 && (p == len(s) || s[p] != c) {
				xx("expected terminating " + S(c))
			}
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
	var b byte
	if w() {
		return nil
	}
	if c[b] == ';' { // terminator ;)]
		return nil
	}
	if chars(&r) || number(&r) || symbol(&r) || quote(&r) || verb(&r) {
	}
	return r
}
func chars(r *K) (v bool) {
	//defer func() { s := iff(v, r); fmt.Printf("chars? %v %s\n", v, s) }()
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
	//defer func() { s := iff(v, r); fmt.Printf("number? %v %s\n", v, s) }()
	for i := 0; ; i++ {
		if p+i == len(s) || s[p+i] < '0' || s[p+i] > '9' {
			if i > 0 {
				*r, _ = strconv.Atoi(string(s[p : p+i]))
				p += i
				return true
			}
			return false
		}
	}
}
func symbol(r *K) (v bool) { // abc
	if s[p] < 'a' || s[p] > 'z' {
		return false
	}
	for i := 0; ; i++ {
		if p+i == len(s) || s[p+i] < 'a' || s[p+i] > 'z' {
			*r = S(s[p : p+i])
			p += i
			return true
		}
	}
}
func quote(r *K) (v bool) { // `abc
	if s[p] != '`' {
		return false
	}
	p++
	if symbol(r) == false {
		*r = S("")
	}
	*r = L{*r}
	return true
}
func verb(r *K) (v bool) {
	//defer func() { s := iff(v, r); fmt.Printf("verb? %v %s\n", v, s) }()
	if c[s[p]] != 'v' {
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
		if len(u) == 1 {
			return "," + o(u[0])
		}
		v := make([]string, len(u))
		for i, e := range u {
			v[i] = o(e)
		}
		return "(" + strings.Join(v, ";") + ")"
	case F:
		s := o(L(u))
		return "{" + s[1:len(s)-2] + "}"
	case C:
		return `"` + string(u) + `"`
	case S:
		return "`" + string(u)
	case I:
		return strconv.Itoa(u)
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
