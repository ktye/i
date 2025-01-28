package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
)

const C = "aaaaaaaaaanaaaaaaaaaaaaaaaaaaaaaadhddddebcdddjgmggggggggggdbdddddffffffffffffffffffffffffffblcddiffffkfffffffffffffffffffffbdcdaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
const T = "abcdefghijfekbabcdefghijfekbabcdefghidfeebabcdefghijfeebabcdefghijfeebabcdemmhidmeebabcdennhidoeebppppppprpppqppabcdemghidmeebabcdefnhijfeeblllllllllllllblllllllllllllbabcdemmhidmeebabcdennhidoeebabcdennhinneebppppppprpppqppppppppppppppppabcdefghidfeeb"
const prims = ":+-*%&|<>=^!~,#_$?@."

var tkst, sp, pp int
var a, b = -1, 0
var src []byte
var lle, ssax = 0, 0
var glo map[string]byte
var loc map[string]byte
var fun map[string][]byte // "xyzr"

func main() {
	//do(`1+2`)
	//do(`1+abc:2`)
	//do(`1-+/!1+2`)
	//do(`(a*b)+c`)
	//do(`f[a]:x`)
	//do(`1+f[2*a;3*b]`)
	//do(`2.3`)
	//do(`a:2 3. 4`)
	//do(`f:"ii":{a:x;a+y}`)
	//do(`f:"i":{x*x:x+3}`) // x:x+3;x*x
	//do(`f:"i":{%-!x}`)
	do(`f:"i":{r:1+!x}`)
}
func do(x string) {
	r := parse(x)
	if len(r.c) != 0 {
		E(r.p, "parse does not complete")
	}
	fmt.Println(r.s)
}
func parse(x string) (r node) {
	a, b, tkst, sp, pp = -1, 0, 0, 0, 0
	src = []byte(x)
	loc = make(map[string]byte)
	fun = make(map[string][]byte)
	c := es()
	if len(c) == 1 {
		r = c[0]
	} else {
		r = node{s:";", c:c}
	}
	for _, f := range []func(x node) node{applyTyp, splitAsn, inlambda(loop), liftloop, retLast, emtC} {
		r = f(r)
	}
	//r.s = chead + r.s
	return r
}

func ssa(t byte, p int) node { s := fmt.Sprintf("$%d", ssax); ssax++; loc[s]=t; return node{s:s, t:t, p:p} }
func ctyp(b byte) string { return map[byte]string{'i':"int",'f':"double",'c':"char",'I':"int*",'F':"double*",'C':"char*"}[b] }
func ctypname(b byte, s string) string {
	r := ctyp(b)
	if strings.HasSuffix(r, "*") || strings.HasPrefix(s, "*") {
		return r + s
	}
	return r + " " + s
}
func emtC(x node) node {
	for i := range x.c {
		x.c[i] = emtC(x.c[i])
	}
	s, p, t := x.s, x.p, ""
	if isasn(x) {
		t = x.c[0].s
		if s := x.c[1].s; s != "" {
			t += "[" + s + "]"
		}
		t += "=" + x.c[2].s
	} else if len(s) == 1 && contains(prims, s[0]) {
		if len(x.c) == 1 {
			s0 := x.c[0].s
			if i := strings.IndexByte("+*%", s[0]); i >= 0 {
				c := []string{"abs", "sqr", "sqrt"}[i]
				t = c + "(" + s0 + ")"
			} else if s == "-" {
				t = "-" + s0
			} else {
				E(p, "emit monadic")
			}
		} else if len(x.c) == 2 {
			s0, s1 := x.c[0].s, x.c[1].s
			if i := strings.IndexByte("+-*%&|", s[0]); i >= 0 {
				c := []string{"+", "-", "*", "/", "min", "max"}[i]
				if len(c) == 1 {
					t = "(" + s0 + c + s1 + ")"
				} else {
					t = c + string(x.t) + "(" + s0 + "," + s1 + ")"
				}
			} else if s == "@" {
				t = s0 + "[" + s1 + "]"
			} else {
				E(p, "emit dyadic")
			}
		} else {
			E(p, "emit primitive")
		}
	} else if len(x.c) == 0 {
		t = x.s
	} else if s == ";" {
		for _, c := range x.c {
			s := emtC(c).s
			if !strings.HasSuffix(s, "}\n") {
				s += ";\n"
			}
			t += " " + s
		}
	} else if s == "_R" {
		r := x.c[0].s
		for s, b := range loc {
			if s != r && b < 'a' {
				t += "_("+s+");"
			}
		}
		t += "return " + r
	} else if s == "_args" {
		a := make([]string, len(x.c))
		for i, c := range x.c {
			a[i] = ctypname(c.t, c.s)
		}
		t = "(" + strings.Join(a, ",") + ")"
	} else if islam(s) {
		t = ctypname(x.t,x.c[0].s) + x.c[1].s + "{\n" + decl(x) + "\n" + x.c[2].s + "}\n"
	} else if s == "_N" { // todo pre
		n := x.c[0].s
		t = "{int $n=" + n + ";"
		if len(x.c) > 2 {
			vt := x.c[2].t
			st := string(vt)
			if vt < 'A' {
				st = "L"+string(vt+32)
			}
			t += ctypname(x.c[2].t, x.c[2].s) + "=$"+st+"($n);\n"
		}
		t += " for(int $i=0;$i<$n;$i++){\n"
		t += "  " + x.c[1].s + ";\n"
		t += " }}\n"
	} else if s == "_n" {
		t = "_n(" + x.c[0].s + ")"
	} else {
		fmt.Println(x)
		E(p, "emit")
	}
	x.s = t
	x.c = nil
	return x
}
func retLast(x node) node {
	if islam(x.s) {
		c := x.c[2]
		if len(c.c) > 0 {
			lc := c.c[len(c.c)-1]
			lc = node{s:"_R", p:lc.p, t:lc.t, c:[]node{lc}}
			c.c[len(c.c)-1] = lc
			x.c[2] = c
		}
	} else {
		for i := range x.c {
			x.c[i] = retLast(x.c[i])
		}
	}
	return x
}
func liftloop(x node) node {
	if islam(x.s) {
		ssax = 0
		x.c[2] = liftloops([]node{x.c[2]})[0]
	} else {
		for i, c := range x.c {
			x.c[i] = liftloop(c)
		}
	}
	return x
}
func liftloops(x []node) []node {
	x0, p := x[0], x[0].p
	if (isasn(x0) && x0.c[2].s == "_N") || x0.s == "_N" {
		var r node
		if x0.s == "_N" {
			r = ssa(x0.t, p)
		} else {
			r = x0.c[0]
			x0 = x0.c[2]
		}
		x0.c = append(x0.c, r)
		x0.c[1] = node{s:":", p:p, t:x0.t + 32, c:[]node{r, node{s:"$i", t:'i', p:p}, x0.c[1]}}
		return append([]node{r, x0}, x[1:]...)
/*rm
	} else if x0.s == "_N" {
		var r node
		fmt.Println("lifloops:", isasn(x0.c[1]), x0)
		if isasn(x0.c[1]) {
			panic("unreachable")
			r = x0.c[1].c[0]
			r.c = nil
		} else if x0.c[1].s == ";" {
			panic("unreachable")
			c := x0.c[1].c
			l := c[len(c)-1]
			if isasn(l) {
				r = l
				r.c = nil
			}
		}
		if r.s == "" {
			r = ssa(x0.t, p)
			if x0.c[1].s == ";" {
				E(p, "lambda ;")
			} else {
				x0.c = append(x0.c, r) //c2
				x0.c[1] = node{s:":", p:p, t:x0.t + 32, c:[]node{r, node{s:"$i", t:'i', p:p}, x0.c[1]}}
			}
		}
		return append([]node{r, x0}, x[1:]...)
*/
	} else if x0.s == ";" {
		var r []node
		for i := range x0.c {
			c := liftloops([]node{x0.c[i]})
			r = append(r, c[1:]...)
			r = append(r, c[0])
		}
		x0.c = r
		x[0] = x0
		return x
	}
	for j := len(x0.c)-1; j>=0; j-- {
		c := liftloops([]node{x0.c[j]})
		x[0].c[j] = c[0]
		x = append(x, c[1:]...)
	}
	return x
}
func splitAsn(x node) node { // x*x:3+x => x:3+x;x*x // (* x (: x (+ 3 x))) => (; (: x (+ 3 x) (* x x))
	if islam(x.s) {
		x.c[2] = asplit([]node{x.c[2]})[0]
	} else {
		for i, c := range x.c {
			x.c[i] = splitAsn(c)
		}
	}
	return x
}
func asplit(x []node) []node {
	x0 := x[0]
	if isasn(x0) {
		rhs := asplit([]node{x0.c[2]})
		idx := asplit([]node{x0.c[1]})
		x[0] = x0.c[0]
		x = append(x, rhs[1:]...)
		x = append(x, idx[1:]...)
		x = append(x, node{s:x0.s, p:x0.p, t:x0.t, c:[]node{x0.c[0], idx[0], rhs[0]}})
		return x
	} else if x0.s == ";" { //insert assignments into innermost statement list
		var r []node
		for i := range x0.c {
			c := asplit([]node{x0.c[i]})
			if len(c) == 2 && strings.HasSuffix(c[1].s, ":") && c[1].c[0].s == c[0].s { //dont split x:y into x:y;y at 1st level
				c = c[1:]
			}
			r = append(r, c[1:]...)
			r = append(r, c[0])
		}
		x0.c = r
		x[0] = x0
		return x
	}
	for j := len(x0.c) - 1; j>=0; j-- {
		c := asplit([]node{x0.c[j]})
		x[0].c[j] = c[0]
		x = append(x, c[1:]...)
	}
	return x
}
func inlambda(f func(x node) node) func(x node) node {
	var g func(x node) node
	g = func(x node) node {
		if islam(x.s) {
			x.c[2] = f(x.c[2])
		} else {
			for i := range x.c {
				x.c[i] = g(x.c[i])
			}
		}
		return x
	}
	return g
}

func applyTyp(x node) node {
	if isasn(x) {
		if islam(x.c[2].s) {
			return lambdaTyp(x.c[0], x.c[2])
		}
		x.c[2] = applyTyp(x.c[2])
		t := x.c[2].t
		x.c[1] = applyTyp(x.c[1])
		if x.c[1].nil() == false {
			if x.c[1].t == 'i' {
				t -= 32 // F[i]:f
			} else if x.c[1].t != 'I' {
				E(x.p, "index type")
			}
		}
		s := x.c[0].s
		if lt, o := loc[s]; o && lt != t {
			E(x.p, "reassign type")
		}
		loc[s] = t
		x.c[0].t = t
		x.t = t
		return x
	}
	if x.s == ";" {
		for i := range x.c {
			x.c[i] = applyTyp(x.c[i])
		}
	} else {
		for i := range x.c {
			j := len(x.c) - i - 1
			if x.c[j].t == 0 {
				x.c[j] = applyTyp(x.c[j])
			}
		}
	}
	s, p := x.s, x.p
	lt, lo := loc[s]
	ft, fn := fun[s]
	_, gl := glo[s]
	if s == "" { //(: v nil rhs)
	} else if contains(prims, s[0]) {
		if len(x.c) == 1 {
			x.t = mot(s[0], x.c[0].t, p)
		} else if len(x.c) == 2 {
			x.t = dyt(s[0], x.c[0].t, x.c[1].t, p)
		} else  { //(: a nil rhs)
			E(p, "typ?")
		}
	} else if s == ";" {
		x.t = x.c[len(x.c)-1].t
	} else if lo {
		x.t = lt
	} else if fn {
		x.t = ft[len(ft)-1]
	} else if gl {
		x.t = glo[s]
	} else if isnum(s) && len(x.c) == 0 {
		x.t = 'i'
		if contains(s, '.') {
			x.t = 'f'
		}
	} else if s == "_vlit" {
		x.t = 'I'
		for _, c := range x.c {
			if contains(c.s, '.') {
				x.t = 'F'
			}
		}
	} else {
		fmt.Println(x)
		E(p, "typ")
	}
	return x
}
func lambdaTyp(sy, x node) node {
	name := sy.s
	p := x.p
	if !contains(x.s, ':') {
		E(p, "untyped lambda")
	}
	v := strings.Split(x.s, ":")
	var rt byte
	if len(v) == 3 { // f:"i:ii":{..}
		rt = v[1][0]
		fun[name] = append([]byte(v[1]), rt)
	} else if len(v) != 2 { // f:"ii":{..}
		E(p, "untyped type annotation")
	}
	args := v[1]
	loc = make(map[string]byte)
	sym := "xyzabcdefghijklmnopqrstuvwxyz"
	var xy []node
	for i := range args {
		loc[string(sym[i])] = args[i]
		xy = append(xy, node{s:string(sym[i]), t:args[i]})
	}
	for i := range x.c {
		if x.c[i].t == 0 {
			x.c[i] = applyTyp(x.c[i])
		}
	}
	if rt == 0 {
		rt = x.c[len(x.c)-1].t
		fun[name] = append([]byte(args), rt)
	}
	x.t = rt
	sy.t = rt
	x.c = append([]node{sy, node{s:"_args", t:rt, c:xy}}, x.c...)
	return x //(symbol;args;body)
}
func decl(x node) (r string) {
	a := make(map[string]bool)
	s := x.c[1].s // (int x,int y)
	for _, c := range strings.Split(s[1:len(s)-1], ",") {
		for _, t := range []string{"char ", "char**", "char*", "int ","int**", "int*", "double complex ","double complex**", "double complex*", "double ", "double**", "double*"} {
			c = strings.TrimPrefix(c, t)
		}
		a[c] = true
	}
	b := make(map[byte][]string)
	for s, t := range loc {
		if !a[s] {
			p := ""
			if t < 'A' {
				t += 64
				p = "**"
			} else if t < 'a' {
				t += 32
				p = "*"
			} 
			b[t] = append(b[t], p+s)
		}
	}
	for _, t := range []byte("cifz") {
		if v, o := b[t]; o {
			r += " " + ctypname(t,strings.Join(v, ",")) + ";"
		}
	}
	return r
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
func loop(x node) (r node) {
	for i := range x.c {
		x.c[i] = loop(x.c[i])
	}
	p := x.p
	if x.s == "!" && len(x.c) == 1 {
		x.s = "_N"
		x.c = append(x.c, node{s:"$i", p:p}) //N(n;body;[pre])
		return x
	} else if len(x.c) == 1 && len(x.s) == 1 && contains("-*%", x.s[0]) {
		t := x.t
		if x.c[0].s == "_N" {
			r := x.c[0]
			c := r.c[1]
			x.c[0] = c
			r.c[1] = x
			return r
		} else if isvec(t) {
			return node{s:"_N", t:t, p:x.p, c: []node{count(x), ati(x)}}
		}
	} else if len(x.c) == 2 && len(x.s) == 1 && contains("+-*%&|<=>", x.s[0]) {
		t := x.t
		if x.c[0].s == "_N" && x.c[1].s == "_N" {
			l := x.c[0]
			r := x.c[1]
			x.c[0] = l.c[1]
			x.c[1] = r.c[1]
			r.c[1] = x
			r.t = t
			r.p = x.p
			return r
		} else if x.c[0].s == "_N" {
			r := x.c[0]
			x.c[0] = r.c[1]
			return x
		} else if x.c[1].s == "_N" {
			r := x.c[1]
			x.c[1] = r.c[1]
			r.c[1] = x
			return r
		}
	} else if isasn(x) && x.c[0].s == "_N" { //unloop lhs l[i]:X
		r := x.c[0]
		x.c[0] = r.c[1].c[0]
		return x
	} else if isasn(x) && x.c[2].s == "_N" { //todo x[i]:..
		r := x.c[2]
		x.c[1] = node{s:"$i", t:'i', p:x.p}
		x.c[2] = r.c[1]
		r.c[1] = x
		return r
	} else if isvec(x.t) && isnam(x.s) {
		return node{s:"_N", p:x.p, t:x.t, c:[]node{count(x), ati(x)}}
	}
	return x
}
func ati(x node) node {	return node{s:"@", p:x.p, t:x.t - 32, c:[]node{x, node{s:"$i", t:'i', p:x.p}}} }
func count(x node) (n node) {
	if x.s == "_vlit" {
		n = node{s:strconv.Itoa(len(x.c)), p:x.p, t:'i'}
	} else if isnam(x.s) {
		n = node{s:"_n", p:x.p, t:'i', c:[]node{x}}
	} else {
		E(x.p, "cannot count")
	}
	return n
}

type node struct {
	s string //string
	p int    //src position
	t byte   //type iIfF..
	v bool   //verb
	n int    //arity
	c []node
}
func (n node) nil() bool { return n.s == "" }
func (n node) print(l int) {
	//fmt.Print(strings.Repeat(" ", l))
	if n.c == nil {
		fmt.Printf("%s", n.s)
		if n.t != 0 {
			fmt.Printf("%c", n.t)
		}
		return
	}
	fmt.Printf("(%s%c", n.s, n.t)
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
	var x node
	var c byte
	for {
		c = peak()
		if c == 59 || 2 == cl(c) {
			return r //;])}
		}
		x = node{s: tok()}
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
	r = x

	if c == '(' {
		y := ls()
		if 1 == len(y) {
			r = y[0]
		} else {
			r.s = "enlist"
			r.n = len(y)
			r.c = y //rev?
		}
	} else if c == '{' {
		if lle > 0 {
			E(x.p, "block")
		}
		lle++
		l := es()
		if peak() != 125 {
			E(pp, "missing "+string(125))
		}
		tok()
		lle--
		r.s = "_lambda"
		r.c = []node{node{s:";", p:r.p, c:l}}
	}
	for {
		y := peak()
		if strings.IndexByte(`'/\`, y) >= 0 {
			y := node{s: tok()}
			y.p = pp
			y.n = 1
			y.v = true
			y.c = []node{r}
			r = y
		} else if y == '[' {
			if len(r.c) > 0 || cl(r.s[0]) != 5 {
				E(pp, "[ expected after variable")
			}
			p := pp
			tok()
			y := rev(ls())
			if len(y) == 1 {
				r = node{s:"@", p:p, n:2, c:[]node{r, y[0]}}
			} else {
				r = node{s:".", p:p, n:len(y), c:append([]node{x}, y...)}
			}
		} else if (isnum(r.s) && len(r.c) == 0) || r.s == "_vlit" { //vector literal
			if 6 == cl(y) || (y == '-' && 6 == cl(src[1+b])) {
				b := node{s: tok()}
				if r.s == "_vlit" {
					r.c = append(r.c, b)
				} else {
					r = node{s:"_vlit", p:r.p, c:[]node{r, b}}
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
	if y.v && !x.v {
		r := e(t())
		if r.nil() || r.v {
			E(y.p, "no projection")
		}
		if y.s[len(y.s)-1] == ':' {
			sy := x
			var ix node
			if x.s == "@" && len(x.c) == 2 {
				sy = x.c[0]
				ix = x.c[1]
			}
			if x.s[0] == '"' && r.s == "_lambda" {
				r.s = "_lambda:" + x.s[1:len(x.s)-1]
				return r
			}
			if len(y.c) != 0 || 5 != cl(sy.s[0]) {
				E(y.p, "assign")
			}
			y.v = false
			y.c = []node{sy, ix, r}
			return y //assign
		}
		y.v = false
		y.c = []node{x, r}
		return y //dy
	}
	r := e(y)
	if !x.v { //juxtaposition
		return node{s:"@", n:2, p:x.p, c: []node{x, r}}
	} else if y.v && x.v && r.eql(y) {
		E(y.p, "no composition")
	}
	x.v = false
	x.c = []node{r}
	return x //mo
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
func isnam(s string) bool { return len(s) > 0 && contains("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",s[0]) }
func islam(s string) bool { return strings.HasPrefix(s, "_lambda") }
func isvec(b byte) bool { return b < 'a' }
func isasn(x node) bool { return len(x.c) == 3 && strings.HasSuffix(x.s, ":") }
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

const chead = `#include<stdlib.h>
#define _(x) if(!--((int*)x)[-2])free(x)
char   *$C(int n){int*r=malloc(8+  n);r[0]=0;r[1]=n;return(char  *)(2+r);}
#define $I(n)  (int            *)$C(n<<2)
#define $F(n)  (double         *)$C(n<<3)
#define $Z(n)  (double complex *)$C(n<<4)
#define $LC(n) (char          **)$C(n<<3)
#define $LI(n) (int           **)$C(n<<3)
#define $LF(n) (double        **)$C(n<<3)
#define $LZ(n) (double complex**)$C(n<<3)
`
