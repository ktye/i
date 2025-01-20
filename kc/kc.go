package main

import (
	"fmt"
	"os"
	"strings"
)

const C = "aaaaaaaaaanaaaaaaaaaaaaaaaaaaaaaadhddddebcdddjgmggggggggggdbdddddffffffffffffffffffffffffffblcddiffffkfffffffffffffffffffffbdcdaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
const T = "abcdefghijfekbabcdefghijfekbabcdefghidfeebabcdefghijfeebabcdefghijfeebabcdemmhidmeebabcdennhidoeebppppppprpppqppabcdemghidmeebabcdefnhijfeeblllllllllllllblllllllllllllbabcdemmhidmeebabcdennhidoeebabcdennhinneebppppppprpppqppppppppppppppppabcdefghidfeeb"

var tkst, sp, pp int
var a, b = -1, 0
var src []byte
var glo map[string]byte
var loc map[string]byte
var fun map[string][]byte // "xyzr"

func main() {
	if len(os.Args) > 1 {
		src = []byte(os.Args[1])
		fun = map[string][]byte{
			"f": []byte("ii"),
			"g": []byte("iii"),
		}
		glo = map[string]byte{"G": 'I'}
		loc = map[string]byte{"x": 'i', "y": 'i', "X": 'I', "Y": 'I', "a": 'f', "b": 'f', "A": 'F', "B": 'F'}
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
	ty := y[len(y)-1]
	p := v[len(v)-1].p
	s := v[len(v)-1].s
	if len(v) != 1 {
		E(p, "assign")
	} else if tx.s == "@" {
		//indexed assign
	} else if len(x) != 1 || 5 != cl(x[0].s[0]) {
		E(p, "must assign to variable")
	} else if s == ":" {
		if l, o := loc[x[0].s]; o {
			if l != ty.t {
				E(p, "reassign type")
			}
		} else {
			loc[x[0].s] = ty.t
		}
	} else {
		E(p, "nyi modified")
	}
	v[len(v)-1].n = 2
	v[len(v)-1].v = false
	v[len(v)-1].t = ty.t
	return append(append(y, x...), v...)
}
func dy(x, o, y []token) []token {
	if len(o) == 1 {
		o[len(o)-1].n = 2
	}
	o[len(o)-1].v = false
	o[len(o)-1].t = dyt(x, o, y)
	return append(append(y, x...), o...)
}
func dyt(x, o, y []token) (t byte) {
	d := o[len(o)-1]
	c := d.s[0]
	tx, ty := x[len(x)-1].t, y[len(y)-1].t
	rx, ry := rank(tx), rank(ty)
	mxr := rx
	if ry > mxr {
		mxr = ry
	}
	t = '?'
	switch c { // :+-*%&|<>=^!~,#_$?@.
	case ':':
		t = ty
	case '+', '-', '*', '%', '&', '|':
		t = maxtype(base(tx), base(ty)) - byte(32*mxr)
	case '<', '>', '=':
		t = []byte{'i', 'I', 'I' - 32}[mxr]
	case '^': // 3^v
		if tx == 'i' || tx == 'I' {
			if ry != 1 {
				E(d.p, "rank")
			}
			t = ty - 32
		} else {
			E(d.p, "type")
		}
	case '!':
		if tx != 'i' || base(ty) != 'i' {
			E(d.p, "type")
		}
		t = ty
	case '~':
		t = 'i'
	case ',':
		if base(tx) != base(ty) {
			E(d.p, "type")
		}
		if rx == 0 && ry == 0 {
			t = tx - 32
		}
	case '#':
		if tx != 'i' {
			E(d.p, "type")
		}
		d.t = ty
		if ry == 0 {
			t -= 32
		}
	case '_':
		if tx != 'i' {
			E(d.p, "type")
		}
		if ry == 0 {
			E(d.p, "rank")
		}
		t = ty
	case '$':
		if tx != 'c' || tx != 'C' {
			E(d.p, "type")
		}
		t = 'C'
	case '?':
		if base(tx) != base(ty) {
			E(d.p, "type")
		}
		if ry > rx {
			E(d.p, "rank")
		}
		t = []byte{'i', 'I', 'I' - 32}[ry]
	case '@': //todo function calls
		t = calt(x, [][]token{y})
	case '.':
		//todo typecheck function arguments
		t = tx
	case '/': // each-right
		if ry == 0 {
			E(d.p, "rank")
		}
		ly := y[len(y)-1]
		ly.t += 32
		t = dyt(x, o[:len(o)-1], []token{ly}) - 32
	case '\\': //each-left
		if rx == 0 {
			E(d.p, "rank")
		}
		lx := x[len(x)-1]
		lx.t += 32
		t = dyt([]token{lx}, o[:len(o)-1], y) - 32
	case '\'': //each2
		if rx == 0 || ry == 0 {
			E(d.p, "rank")
		}
		lx := x[len(x)-1]
		lx.t += 32
		ly := x[len(y)-1]
		ly.t += 32
		t = dyt([]token{lx}, o[:len(o)-1], []token{ly}) - 32
	default:
		E(d.p, "unknown dyadic primitive")
	}
	return t
}
func mo(o, x []token) []token {
	mot(o, x)
	o[len(o)-1].t = mot(o, x)
	if len(o) == 1 {
		o[len(o)-1].n = 1
	}
	o[len(o)-1].v = false
	return append(x, o...)
}
func mot(o, x []token) (t byte) {
	m := o[len(o)-1]
	m.n = 1
	m.v = false
	tx := x[len(x)-1]
	t = tx.t // :+-*
	c := m.s[0]
	rk := rank(tx.t)
	ops := `:+-*%&|<>=^!~,#_$?@./\'`
	mir := "00000211111000000000111"
	mar := "22222221111121222120222"
	i := strings.IndexByte(ops, c)
	if i < 0 {
		E(m.p, "unknown monadic primitive")
	} else if rk < int(mir[i]-'0') || rk > int(mar[i]-'0') {
		E(m.p, "rank")
	}
	switch c {
	case ':', '+', '-', '*':
	case '%':
		t = []byte{'f', 'F', 'F' - 32}[rk]
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
			E(m.p, "type must be i or I")
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
			E(m.p, "type must be i or rank 1")
		}
	case '@':
		if rk > 0 {
			t += 32
		}
	case '.':
		E(m.p, "not implemented")
	case '/':
		tx.t += 32
		t = dyt([]token{tx}, o[:len(o)-1], []token{tx})
	case '\\':
		tx.t += 32
		t = dyt([]token{tx}, o[:len(o)-1], []token{tx}) - 32
	case '\'':
		tx.t += 32
		t = mot(o[:len(o)-1], []token{tx}) - 32
	default:
		E(m.p, "unknown monadic primitive")
	}
	return t
}
func at(x, y []token, p int) []token {
	t := x[len(x)-1]
	t.s = "@"
	t.v = true
	t.p = p
	t.n = 2
	return dy(x, []token{t}, y)
}
func call(f []token, x [][]token, p int) (r []token) {
	t := calt(f, x)
	if len(x) == 1 {
		return append(append(x[0], f...), token{s: "@", p: p, t: t, n: 2, v: false})
	} else {
		for i := range x {
			r = append(r, x[len(x)-i-1]...)
		}
		return append(r, token{s: ".", p: p, t: t, v: false, n: len(x)})
	}
}
func calt(f []token, x [][]token) byte {
	fn := f[len(f)-1]
	s := fn.s
	if 5 == cl(s[0]) {
		if sig := fun[s]; sig != nil && loc[s] == 0 {
			if len(sig) != 1+len(x) {
				E(fn.p, "valence")
			}
			for i, xi := range x {
				if sig[len(sig)-i-2] != xi[len(xi)-1].t {
					E(fn.p, "type")
				}
			}
			return sig[len(sig)-1]
		}
	}
	rk := rank(fn.t)
	b := base(fn.t)
	m := map[string]int{"i": rk - 1, "I": rk, "ii": 0, "iI": 1, "Ii": 1, "II": 2}
	if rk < len(x) || len(x) == 0 {
		E(fn.p, "rank")
	}
	x0 := x[0]
	it := []byte{x0[len(x0)-1].t}
	if len(x) > 1 {
		x1 := x[1]
		it = append(it, x1[len(x1)-1].t)
	}
	r, o := m[string(it)]
	if !o {
		E(fn.p, "type")
	}
	return b - byte(32*r)
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

func t() []token {
	var x token
	var c byte
	for {
		c = peak()
		if c == 59 || 2 == cl(c) {
			return nil //;])}
		}
		x = token{s: tok()}
		if x.s == "" {
			return nil
		}
		if len(x.s) > 1 && x.s[0] == '/' {
			tok()
			continue
		} else {
			break
		}
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
		x.t = '*'
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
			for i := 0; i < len(y); i++ {
				te := y[i][len(y[i])-1].t
				if i == 0 {
					tp = te
				} else if tp != te {
					E(p, "mixed type list")
				}
				r = append(r, y[i]...)
			}
			r = append(r, token{p: p, s: "enlist", t: tp - 32, n: len(y)})
		}
	} else if c == '{' {
		r = es()
		if peak() != '}' {
			E(1+pp, "} expected")
		}
		tok()
		return append(append([]token{x}, r...), token{s: "}", p: pp})
	}
	for {
		y := peak()
		if strings.IndexByte(`'/\`, y) >= 0 {
			y := token{s: tok()}
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
			if len(y) == 1 {
				r = at(r, y[0], p)
			} else {
				r = call(r, y, p)
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
		if ty.s[len(ty.s)-1] == ':' {
			return as(x, y, r)
		}
		return dy(x, y, r)
	}
	r := e(y)
	if !tx.v {
		return at(x, r, tx.p)
	} else if len(r) == len(y) && ty.v && tx.v {
		E(ty.p, "no composition")
	}
	return mo(x, r)
}
func es() (r []token) {
	for {
		if s := peak(); s == '\n' || s == ';' {
			if len(r) > 0 && r[len(r)-1].s != ";" {
				r = append(r, token{s: ";", p: pp})
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
