package main

import (
	"fmt"
	"os"
	"strings"
)

const C = "aaaaaaaaaanaaaaaaaaaaaaaaaaaaaaaadhddddebcdddjgmggggggggggdbdddddffffffffffffffffffffffffffblcddiffffkfffffffffffffffffffffbdcdaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
const T = "abcdefghijfekbabcdefghijfekbabcdefghidfeebabcdefghijfeebabcdefghijfeebabcdemmhidmeebabcdennhidoeebppppppprpppqppabcdemghidmeebabcdefnhijfeeblllllllllllllblllllllllllllbabcdemmhidmeebabcdennhidoeebabcdennhinneebppppppprpppqppppppppppppppppabcdefghidfeeb"
const prims = ":+-*%&|<>=^!~,#_$?@."

var tkst, sp, pp int
var a, b = -1, 0
var src []byte
var glo map[string]byte
var loc map[string]byte
var fun map[string][]byte // "xyzr"

func main() {
	//do(`1+2`)
	//do(`1+abc`)
	//do(`1-+/!1+2`)
	//do(`(a*b)+c`)
	//do(`f[a]:x`)
	//do(`1+f[2*a;3*b]`)
	do(`2.3`)
	do(`a:2 3. 4`)
}
func do(x string) {
	r := parse(x)
	r.print(0)
	fmt.Println()
}
func parse(x string) (r node) {
	a, b, tkst, sp, pp = -1, 0, 0, 0, 0
	src = []byte(x)
	loc = make(map[string]byte)
	c := es()
	if len(c) == 1 {
		r = c[0]
	} else {
		r = node{t:token{s:";"}, c:c}
	}
	for _, f := range []func(x node) node{applyTyp} {
		r = f(r)
	}
	return r
}


func applyTyp(x node) node {
	if x.t.s == ":" && len(x.c) == 3 {
		x.c[2] = applyTyp(x.c[2])
		t := x.c[2].t.t
		x.c[1] = applyTyp(x.c[1])
		if x.c[1].nil() == false {
			if x.c[1].t.t == 'i' {
				t -= 32 // F[i]:f
			} else if x.c[1].t.t != 'I' {
				E(x.t.p, "index type")
			}
		}
		x.c[0].t.t = t
		x.t.t = t
		return x
	}
	for i := range x.c {
		j := len(x.c) - i - 1
		if x.c[j].t.t != 0 {
			x.c[j] = applyTyp(x.c[j])
		}
	}
	s, p := x.t.s, x.t.p
	ft, fn := fun[s]
	_, gl := glo[s]
	if s == "" { //(: v nil rhs)
	} else if contains(prims, s[0]) {
		if len(x.c) == 1 {
			x.t.t = mot(s[0], x.c[0].t.t, p)
		} else if len(x.c) == 2 {
			x.t.t = dyt(s[0], x.c[0].t.t, x.c[1].t.t, p)
		} else  { //(: a nil rhs)
			fmt.Println("s?", s, strings.IndexByte(prims, s[0]))
			E(p, "typ!")
		
			E(p, "typ?")
		}
	} else if fn {
		x.t.t = ft[len(ft)-1]
	} else if gl {
		x.t.t = glo[s]
	} else if isnum(s) && len(x.c) == 0 {
		x.t.t = 'i'
		if contains(s, '.') {
			x.t.t = 'f'
		}
	} else if s == "vlit" {
		x.t.t = 'I'
		for _, c := range x.c {
			if contains(c.t.s, '.') {
				x.t.t = 'F'
			}
		}
	} else {
		E(p, "typ")
	}
	return x
}
func mot(o, x byte, p int) byte {
	rk := rank(x)
	ops := `:+-*%&|<>=^!~,#_$?@.`
	mir := "00000211111000000000"
	mar := "22222221111121222120"
	i := strings.IndexByte(ops, o)
	if i < 0 {
		E(p, "unknown monadic primitive")
	} else if rk < int(mir[i]-'0') || rk > int(mar[i]-'0') {
		E(p, "rank")
	}
	t := x
	switch o {
	case ':', '+', '-', '*':
	case '%':
		x = []byte{'f', 'F', 'F' - 32}[rk]
	case '&':
	case '|':
	case '<', '>':
		t = 'I'
	case '=':
		t = 'I' - 32
	case '^':
	case '!':
		if t == 'i' {
			t = 'I'
		} else if t != 'I' {
			E(p, "type must be i or I")
		}
	case '~':
		t = []byte{'i', 'I', 'I' - 32}[rk]
	case ',':
		t -= 32
	case '#':
		t = 'i'
	case '_':
		t = []byte{'i', 'I', 'I' - 32}[rk]
	case '$':
		t = 'C'
	case '?':
		if t == 'i' {
			t = 'F'
		} else if rk == 1 {
			E(p, "type must be i or rank 1")
		}
	case '@':
		if rk > 0 {
			t += 32
		}
	case '.':
		E(p, "monadic .")
	default:
		E(p, "monadic")
	}
	return t
}
func dyt(o, x, y byte, p int) byte {
	rx, ry := rank(x), rank(y)
	mxr := rx
	if ry > mxr {
		mxr = ry
	}
	t := y
	switch o {
	case ':':
		t = y
	case '+', '-', '*', '%', '&', '|':
		t = maxtype(base(x), base(y)) - byte(32*mxr)
	case '<', '>', '=':
		t = []byte{'i', 'I', 'I' - 32}[mxr]
	case '^':
		if x == 'i' || x == 'I' {
			if ry != 1 {
				E(p, "rank")
			}
			t = y - 32
		} else {
			E(p, "type")
		}
	case '!':
		if x != 'i' || base(y) != 'i' {
			E(p, "type")
		}
	case '~':
		t = 'i'
	case ',':
		if base(x) != base(y) {
			E(p, "type")
		}
		if rx == 0 && ry == 0 {
			t = x - 32
		}
	case '#':
		if x != 'i' {
			E(p, "type")
		}
		if ry == 0 {
			t -= 32
		}
	case '_':
		if x != 'i' || ry == 0 {
			E(p, "type")
		}
	case '$':
		if x != 'c' || x != 'C' {
			E(p, "type")
		}
		t = 'C'
	case '?':
		if base(x) != base(y) {
			E(p, "type")
		}
		if ry > rx {
			E(p, "rank")
		}
		t = []byte{'i', 'I', 'I' - 32}[ry]
	case '@': //todo function calls
		if y == 'i' {
			t = x + 32
		} else if y != 'I' {
			E(p, "index-type")
		}
	case '.':
		E(p, "dyadic .")
	default:
		E(p, "unknown dyadic")
	}
	return t
}

type token struct {
	s string //string
	p int    //src position
	t byte   //type iIfF..
	v bool   //verb
	n int    //arity
}
type node struct {
	t token
	c []node
}
func (n node) nil() bool { return n.t.s == "" }
func (n node) print(l int) {
	//fmt.Print(strings.Repeat(" ", l))
	if n.c == nil {
		fmt.Printf("%s", n.t.s)
		if n.t.t != 0 {
			fmt.Printf("%c", n.t.t)
		}
		return
	}
	fmt.Printf("(%s%c", n.t.s, n.t.t)
	for i := range n.c {
		fmt.Printf(" ")
		n.c[i].print(1+l)
	}
	fmt.Printf(")")
}
func (a node) eql(b node) bool {
	if a.t != b.t {
		return false
	}
	if len(a.c) != len(b.c) {
		return false
	}
	for i := range a.c {
		if a.c[i].t != b.c[i].t || len(a.c[i].c) != len(b.c[i].c) {
			return false
		}
	}
	return true
}

func t() (r node) {
	var x token
	var c byte
	for {
		c = peak()
		if c == 59 || 2 == cl(c) {
			return r //;])}
		}
		x = token{s: tok()}
		if x.s == "" {
			return r
		}
		if len(x.s) > 1 && x.s[0] == '/' {
			tok()
			continue
		} else {
			break
		}
	}
	c = x.s[0]
	if 2 == cl(c) {
		return r //)]};
	}
	x.p = pp
	if 1 == len(x.s) && (3 == cl(c) || c == '-') {
		x.v = true
		x.t = 0
		if peak() == ':' {
			tok()
			x.s += ":"
		}
	}
	r.t = x

	if c == '(' {
		y := ls()
		if 1 == len(y) {
			r = y[0]
		} else {
			r.t.s = "enlist"
			r.t.n = len(y)
			r.c = y //rev? typecheck children
		}
	} else if c == '{' {
		//todo r = es(); peak()=='}' ..
		E(x.p, "{}")
	}
	for {
		y := peak()
		if strings.IndexByte(`'/\`, y) >= 0 {
			y := token{s: tok()}
			y.p = pp
			y.n = 1
			y.v = true
			r = node{t:y, c:[]node{r}}
		} else if y == '[' {
			if len(r.c) > 0 || cl(r.t.s[0]) != 5 {
				E(pp, "[ expected after variable")
			}
			p := pp
			tok()
			y := rev(ls())
			if len(y) == 1 {
				r = node{t:token{s:"@", p:p, n:2}, c:[]node{r, y[0]}}
			} else {
				r = node{t:token{s:".", p:p, n:len(y)}, c:append([]node{node{t:x}}, y...)}
			}
		} else if (isnum(r.t.s) && len(r.c) == 0) || r.t.s == "vlit" { //vector literal
			if 6 == cl(y) || (y == '-' && 6 == cl(src[1+b])) {
				b := node{t:token{s: tok()}}
				if r.t.s == "vlit" {
					r.c = append(r.c, b)
				} else {
					r = node{t:token{s:"vlit", p:r.t.p}, c:[]node{node{t:r.t}, b}}
				}
			} else {
				break
			}
		} else {
			break
		}
	}
	return r
}
func e(x node) node {
	if x.nil() {
		return x
	}
	y := t()
	if y.nil() {
		return x
	}
	tx := x.t
	ty := y.t
	if ty.v && !tx.v {
		r := e(t())
		if r.nil() || r.t.v {
			E(ty.p, "no projection")
		}
		if ty.s[len(ty.s)-1] == ':' {
			sy := tx
			var ix node
			if x.t.s == "@" && len(x.c) == 2 {
				sy = x.c[0].t
				ix = x.c[1]
			}
			if len(y.c) != 0 || 5 != cl(sy.s[0]) {
				E(ty.p, "assign")
			}
			ty.v = false
			return node{t:ty, c:[]node{node{t:sy}, ix, r}} //assign
		}
		ty.v = false
		return node{t: ty, c: []node{x, r}} //dy(x, y, r)
	}
	r := e(y)
	if !tx.v { //juxtaposition
		return node{t: token{s:"@", n:2, p: tx.p}, c: []node{x, r}}
	} else if ty.v && tx.v && r.eql(y) {
		E(ty.p, "no composition")
	}
	tx.v = false
	return node{t:tx, c:[]node{r}} //mo(x, r)
}
func es() (r []node) {
	for {
		if s := peak(); s == '\n' || s == ';' {
			tok()
			continue
		}
		x := e(t())
		if x.nil() {
			return r
		}
		r = append(r, x)
	}
}
func ls() (r []node) {
	for {
		if s := peak(); s == 0 || 2 == cl(s) {
			tok()
			return r
		} else if s == 59 {
			tok()
			continue
		}
		x := e(t())
		r = append(r, x)
	}
}
func rev(x []node) []node {
	if len(x) < 2 {
		return x
	}
	r := make([]node, len(x))
	for i := range r {
		r[i] = x[len(x)-1-i]
	}
	return r
}
func isnum(s string) bool { return (len(s) > 1 && s[0] == '-') || 6 == cl(s[0]) }
func contains(s string, b byte) bool { return strings.IndexByte(s, b) >= 0 }


func E(p int, s string) {
	fmt.Println("p", p)
	a, b := 0, len(src)
	if p >= a && p <= b {
		for i := 0; i < 10; i++ {
			if p-i == 0 {
				break
			}
			if src[p-i] == '\n' {
				a = 1 + p - i
			}
		}
		for i := 0; i < 20; i++ {
			if p+i == b {
				break
			}
			if src[p+i] == '\n' {
				b = p + i
			}
		}
		fmt.Println(string(src[a:b]))
	}
	for i := 0; i < p-a; i++ {
		fmt.Print(" ")
	}
	fmt.Println("^")
	fmt.Println(linepos(p), s)
	os.Exit(1)
}
func linepos(x int) string {
	l, c := 0, 0
	if x < 0 || x >= len(src) {
		return "file?:"
	}
	for i := 0; i < x; i++ {
		c++
		if src[i] == 10 {
			l++
			c = 0
		}
	}
	return fmt.Sprintf("file:%d:%d", 1+l, 1+c)
}

func rank(c byte) int {
	if c >= 'a' {
		return 0
	} else if c >= 'A' {
		return 1
	} else {
		return 2
	}
}
func base(c byte) byte { return byte(int(c) + 32*rank(c)) }
func maxtype(a, b byte) byte {
	ia, ib := strings.IndexByte("cifz", a), strings.IndexByte("cifz", b)
	if ia > ib {
		return a
	}
	return b
}
func cl(c byte) int { return int(C[c]) - 97 }
func nxt() int {
	if sp == len(src) {
		return -1
	}
	for sp < len(src) {
		tkst = int(T[14*tkst+cl(src[sp])] - 97)
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
	pp = a
	if a < 0 {
		return ""
	}
	for a >= 0 && cl(src[a]) == 0 {
		a = nxt()
		pp = a
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
func peak() byte {
	for b >= 0 && (src[b] == 32 || src[b] == '\t') {
		b++
	}
	if b < 0 {
		return 0
	}
	return src[b]
}
