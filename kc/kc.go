package main

import (
	"fmt"
	"os"
	"strings"
)

const C="aaaaaaaaaanaaaaaaaaaaaaaaaaaaaaaadhddddebcdddjgmggggggggggdbdddddffffffffffffffffffffffffffblcddiffffkfffffffffffffffffffffbdcdaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
const T="abcdefghijfekbabcdefghijfekbabcdefghidfeebabcdefghijfeebabcdefghijfeebabcdemmhidmeebabcdennhidoeebppppppprpppqppabcdemghidmeebabcdefnhijfeeblllllllllllllblllllllllllllbabcdemmhidmeebabcdennhidoeebabcdennhinneebppppppprpppqppppppppppppppppabcdefghidfeeb"
var tkst, sp, pp int
var a, b = -1, 0
var src []byte

func main() {
	if len(os.Args) > 1 {
		src = []byte(os.Args[1])
		//for s := tok(); s != ""; s = tok() {
		//	fmt.Println(s)
		//}
		r := es()
		fmt.Println(r)
		for _, t := range r {
			fmt.Printf("%c", t.s[0])
		}
		fmt.Println()
		for _, t := range r {
			fmt.Printf("%c", t.t)
		}	
		fmt.Println()
		for _, t := range r {
			fmt.Printf("%d", t.n)
		}	
		fmt.Println()
	}
}
func E(p int, s string) {
	fmt.Println("p",p)
	a, b  := 0, len(src)
	if p >= a && p <= b {
		for i := 0; i<10; i++ {
			if p-i == 0 {
				break
			}
			if src[p-i] == '\n' {
				a = 1+p-i
			}
		}
		for i := 0; i<20; i++ {
			if p+i == b {
				break
			}
			if src[p+i] == '\n' {
				b = p+i
			}
		}
		fmt.Println(string(src[a:b]))
	}
	for i:=0; i<p-a; i++ {
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
	for i := 0; i<x; i++ {
		c++
		if src[i] == 10 {
			l++
			c = 0
		}
	}
	return fmt.Sprintf("file:%d:%d", 1+l, 1+c)
}

type token struct {
	s string //string
	p int    //src position
	t byte   //type iIfF..
	v bool   //verb
	n int    //arguments
}

func as(x, v, y []token) []token {
	fmt.Println("as x:", x, "y:", y, "v:", v)
	tx := x[len(x)-1]
	if tx.s == "@" {
		//indexed assign
	} else if len(x) != 1 || 5 != cl(x[0].s[0]) {
		E(v[len(v)-1].p, "must assign to variable")
	}
	v[len(v)-1].n = 2
	v[len(v)-1].v = false
	return append(append(y, x...), v...)
}
func dy(x, o, y []token) []token {
	o[len(o)-1].n = 2
	o[len(o)-1].v = false
	return append(append(y, x...), o...)
}
func mo(o, x []token) []token {
	o[len(o)-1].n = 1
	o[len(o)-1].v = false
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
	c := peak()
	if c == 59 || 2 == cl(c) {
		return nil //;])}
	}
	x := token{s: tok()}
	if x.s == "" {
		return nil
	}
	c = x.s[0]
	fmt.Printf("c %c %d\n", c, cl(c))
	if 2 == cl(c) {
		return nil //)]};
	}
	x.p = pp
	x.t = 'i'
	if len(x.s) > 1 && strings.Contains(x.s, ".") {
		x.t = 'f'
	}
	if 1 == len(x.s) && (3 == cl(c) || c == '-') {
		x.v = true
		if peak() == ':' {
			tok()
			x.s += ":"
		}
	}
	r := []token{x}

	if c == '(' {
		y := ls()
		if 1 == len(y) {
			r = y[0]
		} else {
			p := pp
			r = []token{}
			y = rev(y)
			tp := byte(0)
			for i := 0; i<len(y); i++ {
				te := y[i][len(y[i])-1].t
				if i == 0 {
					tp = te
				} else if tp != te {
					E(p, "mixed type list")
				}
				r = append(r, y[i]...)
			}
			r = append(r, token{p:p, s:"enlist", t: tp-32, n:len(y)})
		}
	} else if c == '{' {
		r = es()
		if peak() != '}' {
			E(1+pp, "} expected")
		}
		tok()
		return append(append([]token{x}, r...),token{s:"}", p:pp})
	}
	for {
		y := peak()
		if strings.IndexByte(`'/\`, y) >= 0 {
			y := token{s:tok()}
			y.p = pp
			y.n = 1
			y.v = true
			r = append(r, y)
		} else if y == '[' {
			if len(r) != 1 || cl(r[0].s[0]) != 5 {
				E(pp, "[ expected after variable")
			}
			p := pp
			tok()
			y := rev(ls())
			r = []token{}
			for i := range y {
				r = append(r, y[i]...)
			}
			//x.n = len(y)
			r = append(r, x)
			fmt.Println("y", y, "#y", len(y))
			if len(y) == 1 {
				r = append(r, token{p:p, s:"@", n:2})
			} else {
				r = append(r, token{p:p, s:".", n:1+len(y)})
			}
		} else {
			break
		}
	}
	return r
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
		if r == nil || r[len(r)-1].v {
			fmt.Println("r", r)
			E(ty.p, "no projection")
		}
		if ty.s[len(ty.s)-1] ==':' {
			return as(x, y, r)
		}
		return dy(x, y, r)
	}
	r := e(y)
	if !tx.v {
		return at(x, r)
	} else if len(r) == len(y) && ty.v && tx.v {
		E(ty.p, "no composition")
	}
	return mo(x, r)
}
func es() (r []token) {
	for {
		if s := peak(); s == '\n' || s == ';' {
			if len(r) > 0 && r[len(r)-1].s != ";" {
				r = append(r, token{s:";",p:pp})
			}
			tok()
			continue
		}

		x := e(t())
		if x == nil {
			return r
		}
		r = append(r, x...)
	}
}
func ls() (r [][]token) {
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
func rev(x [][]token) [][]token {
	if len(x) < 2 {
		return x
	}
	r := make([][]token, len(x))
	for i := range r {
		r[i] = x[len(x)-1-i]
	}
	return r
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
	if b < 0 {
		return 0
	}
	return src[b] 
}
