package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

type c = byte
type T c
type fn struct { // name:I:IIF::body..
	name string
	src  [2]int // line, col
	t    T      // return type
	args []T
	locl []T
	lmap map[string]int
	sign int
	ast  expr
	bytes.Buffer
}

const (
	I = T(0x7f) // i32
	J = T(0x7e) // i64
	F = T(0x7c) // f64
)
const (
	sNewl = iota
	sFnam
	sRety
	sArgs
	sLocl
	sByte
	sBody
	sCmnt
)

var typs = map[c]T{'I': I, 'J': J, 'F': F}
var tnum = map[T]int{I: 0, J: 1, F: 2}
var ctyp = map[T]string{I: "int32_t", J: "int64_t", F: "double"}
var gtyp = map[T]string{I: "int32", J: "int64", F: "float64"}
var P = I // I(wasm32) J(wasm64)

type module []fn

//func main() {
//	f := fn{t: F, args: []T{F, F}, Buffer: *bytes.NewBuffer([]c("3*x+y}"))}
//	e := f.parse()
//	fmt.Printf("%#+v\n", e)
//}

func main() {
	var html, cout, gout bool
	flag.BoolVar(&html, "html", false, "html output")
	flag.BoolVar(&cout, "cout", false, "c output")
	flag.BoolVar(&gout, "gout", false, "go output")
	flag.Parse()
	m := run(os.Stdin)
	if html {
		os.Stdout.Write(page(m.wasm()))
	} else if cout {
		os.Stdout.Write(m.cout())
	} else if gout {
		os.Stdout.Write(m.gout())
	} else {
		os.Stdout.Write(m.wasm())
	}
}
func (t T) String() string {
	if c := map[T]c{I: 'I', J: 'J', F: 'F'}[t]; c == 0 {
		return "0"
	} else {
		return string(c)
	}
}
func run(r io.Reader) module {
	rd := bufio.NewReader(r)
	state := sNewl
	line, char, hi := 1, 0, true
	err := func(s string) { panic(fmt.Sprintf("%d:%d: %s", line, char, s)) }
	var m module
	var f fn
	var p c
	for {
		b, e := rd.ReadByte()
		if e == io.EOF {
			if f.name != "" {
				m = append(m, f)
			}
			return m
		} else if e != nil {
			panic(e)
		}
		char++
		if b == '\n' {
			line++
			char = 1
		}
		switch state {
		case sNewl:
			if b == '/' {
				state = sCmnt
			} else if b == ' ' || b == '\t' || b == '\n' {
				if f.name == "" {
					err("parse name")
				}
				state = sByte
			} else {
				if f.name != "" {
					m = append(m, f)
				}
				f = fn{name: string(b)}
				state = sFnam
			}
		case sFnam:
			if craZ(b) || (len(f.name) > 0 && cr09(b)) {
				f.name += string(b)
			} else if b == ':' {
				state = sRety
			} else {
				err("parse function name")
			}
		case sRety:
			if f.t == 0 {
				f.t = typs[b]
				if f.t == 0 {
					err("parse return type")
				}
			} else if b == ':' {
				state = sArgs
			} else {
				err("parse return type")
			}
		case sArgs:
			if t := typs[b]; t == 0 && f.args == nil {
				err("parse args")
			} else if t != 0 {
				f.args = append(f.args, t)
			} else if b == ':' {
				state = sLocl
			} else {
				err("parse args")
			}
		case sLocl:
			if t := typs[b]; t == 0 && len(f.locl) == 0 && b == '{' {
				f.src = [2]int{line, char}
				state = sBody
			} else if t != 0 {
				f.locl = append(f.locl, t)
			} else if b == ':' {
				state = sByte
			} else {
				err("parse locals")
			}
		case sBody:
			f.WriteByte(b)
			if b == '}' {
				state = sCmnt
				f.parse()
			}
		case sByte:
			if b == '/' && hi {
				state = sCmnt
			} else if (b == ' ' || b == '\t') && hi {
				continue
			} else if (b == '\n') && hi {
				state = sNewl
			} else if crHx(b) {
				if hi {
					p, hi = xtoc(b)<<4, false
				} else {
					p, hi = p|xtoc(b), true
					f.WriteByte(p)
					p = 0
				}
			} else {
				err("parse body")
			}
		case sCmnt:
			if b == '\n' {
				state = sNewl
			}
		default:
			err("internal parse state")
		}
	}
}
func hxb(x c) (c, c) { h := "0123456789abcdef"; return h[x>>4], h[x&0x0F] }
func cr09(c c) bool  { return c >= '0' && c <= '9' }
func craZ(c c) bool  { return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') }
func cr0Z(c c) bool  { return cr09(c) || craZ(c) }
func crHx(c c) bool  { return cr09(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F') }
func xtoc(x c) c {
	switch {
	case x < ':':
		return x - '0'
	case x < 'G':
		return 10 + x - 'A'
	default:
		return 10 + x - 'a'
	}
}
func boolvar(v bool) int {
	if v {
		return 1
	}
	return 0
}

type parser struct {
	*fn
	p   int
	b   []byte
	tok []byte
}

func (f *fn) parse() expr { // parse function body
	f.lmap = make(map[string]int)
	p := parser{fn: f, b: strip(f.Bytes())}
	e := p.seq('}')
	e = p.locals(e, 0)
	f.ast = p.validate(e)
	return f.ast
}
func strip(b []c) []c { // strip comments
	lines := bytes.Split(b, []c{'\n'})
	for i, l := range lines {
		space := false
		for k, c := range l {
			if c == ' ' {
				space = true
			} else if c == '/' && space {
				lines[i] = l[:k]
				break
			}
		}
	}
	return bytes.Join(lines, []c{'\n'})
}
func (p *parser) err(s string) expr {
	panic(s)
	return nil
}
func (p *parser) w() {
	for len(p.b) > 0 {
		if c := p.b[0]; c == ' ' || c == '\t' || c == '\n' {
			p.p++
			p.b = p.b[1:]
		} else {
			return
		}
	}
}
func (p *parser) t(f func([]c) int) bool { // test
	p.tok = nil
	if len(p.b) < 1 {
		return false
	}
	if n := f(p.b); n > 0 {
		p.tok = p.b[:n]
		p.b = p.b[n:]
		p.p += n
		return true
	}
	return false
}
func (p *parser) seq(term c) expr {
	var seq seq
	for {
		e := p.ex(p.noun())
		if e != nil {
			seq = append(seq, e)
		} else {
			p.w()
			if len(p.b) == 0 {
				p.err("missing " + string(term))
			}
			if p.b[0] == term {
				p.b = p.b[1:]
				break
			} else if p.b[0] != ';' {
				p.err("expected ;")
			} else {
				p.b = p.b[1:]
			}
		}
	}
	if seq == nil {
		return nil // empty?
	} else if len(seq) == 1 {
		return seq[0]
	}
	return seq
}
func (p *parser) ex(x expr) expr {
	if x == nil {
		return x
	}
	h := p.p
	if p.verb(x) {
		if y := p.ex(p.noun()); y == nil {
			return x // verb ?
		} else {
			return p.monadic(x, y, pos(h))
		}
	} else {
		v := p.noun()
		if v == nil {
			return x // noun
		} else if p.verb(v) {
			h = p.p
			if y := p.ex(p.noun()); y == nil {
				return p.err("verb-verb (missing noun)")
			} else {
				return p.dyadic(v, x, y, pos(h))
			}
		} else {
			return p.err("noun-noun (missing verb)")
		}
	}
}
func (p *parser) monadic(f, x expr, h pos) expr {
	switch v := f.(type) {
	case opx:
		return v1{s: string(v), argv: argv{x}, pos: h}
	default:
		panic("nyi")
	}
}
func (p *parser) dyadic(f, x, y expr, h pos) expr {
	switch v := f.(type) {
	case nlp:
		return nlp{pos: h, argv: argv{x, y}}
	case opx:
		if _, o := v2Tab[string(v)]; o {
			return v2{s: string(v), argv: argv{x, y}, pos: h}
		}
		return cmp{s: string(v), argv: argv{x, y}, pos: h}
	default:
		panic("nyi")
	}
}
func (p *parser) verb(v expr) bool {
	switch v.(type) {
	case opx, nlp: // todo: others
		return true
	}
	return false
}
func (p *parser) noun() expr {
	p.w()
	if len(p.b) == 0 {
		return nil
	}
	switch {
	case p.t(sArg):
		a := p.pArg(p.tok).(arg)
		if a.n >= len(p.args) {
			return p.err("arg %s does not exist")
		}
		a.t = p.args[a.n]
		return a
	case p.t(sSym):
		return p.pSym(p.tok)
	case p.t(sC('/')):
		return nlp{}
	case p.t(sCon):
		return p.pCon(p.tok)
	case p.t(sOp):
		return p.pOp(p.tok)
	case p.t(sC('(')):
		return p.seq(')')
	default:
		return nil
	}
}

func (p *parser) locals(e expr, lv int) expr {
	if av, o := e.(argvec); o == false {
		return e
	} else {
		v := av.args()
		for i, a := range v {
			if f, o := a.(nlp); o {
				c := string('i' + lv)
				n, o := p.fn.lmap[c]
				if o == false {
					n = len(p.fn.lmap)
					p.fn.lmap[c] = n
				}
				f.c = n // set loop counter
				f.argv[0] = p.locals(f.argv[0], lv+1)
				f.argv[1] = p.locals(f.argv[1], lv+1)
				v[i] = f
			} else {
				// TODO: detect locals from assignments
				v[i] = p.locals(v[i], lv)
			}
		}
		return e
	}
}
func (p *parser) validate(e expr) expr {
	if e.valid() == false {
		if i, o := e.(indicator); o {
			return p.indicate(i.indicate(), "invalid")
		} else {
			return p.err("invalid")
		}
	}
	if t := e.rt(); t != p.fn.t {
		return p.err(fmt.Sprintf("return type is %s not %s", t, p.fn.t))
	}
	return e
}
func (p *parser) indicate(pos int, e string) expr {
	s := string(p.fn.Bytes())
	lines := strings.Split(s, "\n")
	for _, l := range lines {
		if pos < len(l) {
			if pos > 0 {
				pos--
			}
			return p.err("\n" + l + "\n" + strings.Repeat(" ", pos) + "^" + e)
		}
		pos -= len(l) + 1
	}
	return p.err(e)
}

func sArg(b []c) int { // x y z x3 x4 x5 ..
	c := b[0]
	if c != 'x' && c != 'y' && c != 'z' {
		return 0
	}
	if len(b) < 2 || b[0] != 'x' || cr09(b[1]) == false {
		return 1 // x y z
	}
	return 2 // x1..x9
}
func (p *parser) pArg(b []c) expr {
	if len(b) == 1 {
		return arg{n: int(b[0] - 'x')}
	}
	return arg{n: (int(b[1] - '0'))}
}
func sSym(b []c) int { // [aZ][a9]*
	c := b[0]
	if craZ(c) == false {
		return 0
	}
	for i, c := range b {
		if craZ(c) == false || cr09(c) == false {
			return i
		}
	}
	return len(b)
}
func (p *parser) pSym(b []c) expr { return loc{s: string(b)} }
func sCon(b []c) int { // 123 123i 123j .123 123. -..
	dot, neg := false, 0
	if len(b) > 1 && b[0] == '-' {
		neg, b = 1, b[1:]
	}
	for i, c := range b {
		if cr09(c) {
			continue
		} else if dot == false && (c == 'i' || c == 'j') {
			return neg + i + 1
		} else if dot == false && c == '.' {
			dot = true
		} else {
			return neg*boolvar(i > 0) + i
		}
	}
	return neg*boolvar(len(b) > 0) + len(b)
}
func (p *parser) pCon(b []c) expr {
	var r con
	if bytes.IndexByte(b, '.') != -1 {
		if f, e := strconv.ParseFloat(string(b), 64); e != nil {
			return p.err(e.Error())
		} else {
			r.t = F
			r.f = f
			return r
		}
	}
	r.t = I
	if c := b[len(b)-1]; c == 'i' || c == 'j' {
		b = b[:len(b)-1]
		if c == 'j' {
			r.t = J
		}
	}
	if i, e := strconv.ParseInt(string(b), 10, 64); e != nil {
		return p.err(e.Error())
	} else {
		r.i = i
	}
	return r
}
func sOp(b []c) int {
	for _, n := range []int{3, 2, 1} { // longest match first
		if len(b) >= n && allops[string(b[:n])] {
			return n
		}
	}
	return 0
}
func (p *parser) pOp(b []c) expr { return opx(string(b)) }
func sC(x c) func(b []c) int     { return func(b []c) int { return boolvar(b[0] == x) } }

// intermediate representation for function bodies (typed expression tree)
type expr interface {
	rt() T // result type, maybe 0
	valid() bool
	bytes() []c
}
type argvec interface{ args() []expr }
type argv []expr
type cstringer interface{ cstr() string }
type gstringer interface{ gstr() string }
type seq argv    // a;b;..
type v2 struct { // x+y unitype
	pos
	argv
	s string // +-*%
}
type v1 struct { // -y
	pos
	argv
	s string
	p int
}
type cmp struct { // x<y..
	pos
	argv
	s string
}
type arg struct { // x y ..
	pos
	t T
	n int
}
type loc struct { // abc
	pos
	arg
	s string
}
type con struct { // numeric constant
	pos
	t T
	i int64
	f float64
}
type nlp struct { // x/y loop
	pos
	argv
	x, y expr
	c    int
}
type opx string // operator
type pos int    // src position
type indicator interface{ indicate() int }

func (p pos) indicate() int { return int(p) }

func getop(tab map[string]code, op string, t T) (r c) {
	ops, ok := tab[op]
	if !ok {
		panic("type")
	}
	switch t {
	case I:
		r = ops.I
	case J:
		r = ops.J
	case F:
		r = ops.F
	default:
		panic("type")
	}
	if r == 0 {
		panic("type")
	}
	return r
}
func cop(tab map[string]code, op string, t T) (o, u string) {
	ops, ok := tab[op]
	if !ok {
		panic("type")
	}
	o = ops.c
	if strings.Index(o, ";") != -1 {
		v := strings.Split(o, ";")
		o = v[tnum[t]]
	}
	if o[0] == 'U' {
		if t != F {
			u = "(u" + ctyp[t] + ")"
		}
		o = o[1:]
	}
	return o, u
}
func gop(tab map[string]code, op string, t T) (o, u string) {
	ops, ok := tab[op]
	if !ok {
		panic("type")
	}
	o = ops.g
	if strings.Index(o, ";") != -1 {
		v := strings.Split(o, ";")
		o = v[tnum[t]]
	}
	if o[0] == 'U' {
		if t != F {
			u = "u" + gtyp[t]
		}
		o = o[1:]
	}
	return o, u
}

func (a argv) args() []expr { return a }
func (a argv) x() expr      { return a[0] }
func (a argv) y() expr      { return a[1] }
func (s seq) rt() T         { return s[len(s)-1].rt() }
func (s seq) args() []expr  { return s }
func (s seq) valid() bool { // all but the last expressions in a sequence must have no return type
	for i, e := range s {
		if i == len(s)-1 {
			return e.rt() != 0
		} else if e.rt() != 0 {
			return false
		}
	}
	return true
}
func (s seq) bytes() (r []c) {
	for _, e := range s {
		r = append(r, e.bytes()...)
	}
	return r
}
func (v v2) rt() T       { return v.x().rt() }
func (v v2) valid() bool { return v.x().rt() == v.y().rt() && v.x().rt() != 0 }
func (v v2) bytes() []c {
	return append(append(v.x().bytes(), v.y().bytes()...), getop(v2Tab, v.s, v.rt()))
}
func (v v2) cstr() string { return c2str(v2Tab, v.s, v.rt(), v.x(), v.y()) }
func (v v2) gstr() string { return g2str(v2Tab, v.s, v.rt(), v.x(), v.y()) }
func (v v1) rt() T        { return v.x().rt() }
func (v v1) valid() bool  { return v.x().rt() != 0 }
func (v v1) bytes() []c   { return append(v.x().bytes(), getop(v1Tab, v.s, v.rt())) }
func (v v1) cstr() string { o, u := cop(v1Tab, v.s, v.rt()); return jn(o, "(", u, cstring(v.x()), ")") }
func (v v1) gstr() string {
	o, u := gop(v1Tab, v.s, v.rt())
	return jn(o, u, "((", gstring(v.x()), "))")
}
func (v cmp) rt() T       { return I }
func (v cmp) valid() bool { return v.x().rt() == v.y().rt() && v.x().rt() != 0 }
func (v cmp) bytes() []c {
	return append(append(v.x().bytes(), v.y().bytes()...), getop(cTab, v.s, v.rt()))
}
func (v cmp) cstr() string { return c2str(cTab, v.s, v.rt(), v.x(), v.y()) }
func (v cmp) gstr() string { return g2str(cTab, v.s, v.rt(), v.x(), v.y()) }
func (v arg) rt() T        { return v.t }
func (v arg) valid() bool  { return v.n >= 0 && v.t != 0 }
func (v arg) bytes() []c   { return append([]c{0x20}, lebu(int(v.n))...) }
func (v arg) cstr() string { return string("x") + string(c('0')+c(v.n)) }
func (v arg) gstr() string { return v.cstr() }
func (v con) rt() T        { return v.t }
func (v con) valid() bool  { return v.t != 0 }
func (v con) bytes() (r []c) {
	r = append([]c{0x41}, lebu(int(v.i))...)
	if v.t == J {
		r[0]++
	} else if v.t == F {
		b := make([]byte, 9)
		b[0] = 0x44
		binary.LittleEndian.PutUint64(b[1:], math.Float64bits(v.f))
		return b
	}
	return r
}
func (v con) cstr() string {
	if v.t == F {
		s := fmt.Sprintf("%v", v.f)
		if strings.Index(s, ".") == -1 {
			s += ".0"
		}
		return s
	}
	return fmt.Sprintf("%d", v.i)
}
func (v con) gstr() string { return v.cstr() }
func (v nlp) rt() T        { return 0 }
func (v nlp) valid() bool  { return v.x.rt() == I && v.y.rt() == 0 }
func (v nlp) bytes() []c   { return []byte{} } // TODO
func (v nlp) cstr() string {
	return fmt.Sprintf("for(x%d=0;x%d<(%s);x%d++){%s}", v.c, v.c, cstring(v.x), v.c, cstring(v.y))
}
func (v nlp) gstr() string {
	return fmt.Sprintf("for x%d=0;x%d<(%s);x%d++{%s}", v.c, v.c, gstring(v.x), v.c, gstring(v.y))
}
func (v opx) rt() T       { return 0 }
func (v opx) valid() bool { return false }
func (v opx) bytes() []c  { return nil }

type code struct {
	I, J, F c
	c, g    string
}

func c2str(tab map[string]code, op string, t T, x, y expr) string {
	o, u := cop(tab, op, t)
	if len(o) > 2 {
		return jn(u, o, "(", cstring(x), ",", cstring(y), ")")
	} else {
		return jn("((", u, cstring(x), ")", o, "(", u, cstring(y), "))")
	}
}
func g2str(tab map[string]code, op string, t T, x, y expr) string {
	o, u := cop(tab, op, t)
	u += "("
	if len(o) > 2 {
		return jn(u, o, "(", cstring(x), ",", cstring(y), "))")
	} else {
		return jn("((", u, cstring(x), "))", o, "(", u, cstring(y), ")))")
	}
}
func cstring(x expr) string { xs := x.(cstringer); return xs.cstr() }
func gstring(x expr) string { xs := x.(gstringer); return xs.gstr() }

var v1Tab = map[string]code{
	"-": code{0, 0, 0x9a, "-", "-"},                                                                            // neg (no neg for ints)
	"+": code{0, 0, 0x99, "__builtin_fabs", "math.Abs"},                                                        // abs
	"_": code{1, 1, 0x9c, ";;__builtin_floor", "math.Floor"},                                                   // floor (ceil, trunc, nearest?)
	"*": code{0x67, 0x79, 0, "__builtin_clz;__builtin_clzll;", "Ubits.LeadingZeros32;Ubits.LeadingZeros64;"},   // clz
	"|": code{0x68, 0x79, 0, "__builtin_ctz;__builtin_ctzll;", "Ubits.TrailingZeros32;Ubits;TrailingZeros64;"}, // ctz
	"%": code{0, 0, 0x9f, "__builtin_sqrt", "math.Sqrt"},                                                       // sqr
}
var v2Tab = map[string]code{
	`+`:   code{0x6a, 0x7c, 0xa0, "+", "+"},               // add
	`-`:   code{0x6b, 0x7d, 0xa1, "-", "-"},               // sub
	`*`:   code{0x6c, 0x7e, 0xa2, "*", "*"},               // mul
	`%`:   code{0x6e, 0x80, 0xa3, "U/", "U/"},             // div/div_u
	`%'`:  code{0x6d, 0x7f, 0xa3, "/", "/"},               // div_s
	`\`:   code{0x70, 0x82, 0, "U%", "%U"},                // rem_u
	`\'`:  code{0x6f, 0x81, 0, "%", "%"},                  // rem_s
	`&`:   code{0x71, 0x83, 0, "&", "&"},                  // and
	`|`:   code{0x72, 0x84, 0, "|", "|"},                  // or
	`^`:   code{0x73, 0x85, 0, "^", "^"},                  // xor
	`<<`:  code{0x74, 0x86, 0, "<<", "<<"},                // shl
	`>>`:  code{0x76, 0x88, 0, "U>>", "U>>"},              // shr_u
	`>>'`: code{0x75, 0x87, 0, ">>", ">>"},                // shl_s
	`<|'`: code{0x77, 0x89, 0, "", ""},                    // rotl
	`>|'`: code{0x78, 0x8a, 0, "", ""},                    // rotr
	`&'`:  code{0, 0, 0xa4, "__builtin_fmin", "math.Max"}, // min
	`|'`:  code{0, 0, 0xa5, "__builtin_fmax", "math.Min"}, // max
}
var cTab = map[string]code{
	"<":   code{0x49, 0x54, 0x63, "U<", "U<"},   // lt/lt_u
	"<'":  code{0x48, 0x53, 0x63, "<", "<"},     // lt_s
	">":   code{0x4b, 0x56, 0x64, "U>", "U>"},   // gt/gt_u
	">'":  code{0x4a, 0x55, 0x64, ">", ">"},     // gt_s
	"<=":  code{0x4d, 0x58, 0x65, "U<=", "U<="}, // le/le_u
	"<='": code{0x4c, 0x57, 0x65, "<=", "<="},   // le_s
	">=":  code{0x4f, 0x5a, 0x66, "U>=", "U>="}, // ge/ge_u
	">='": code{0x4e, 0x59, 0x66, ">=", ">="},   // ge/ge_s
	"~":   code{0x46, 0x51, 0x61, "==", "=="},   // eq
	"!":   code{0x47, 0x52, 0x62, "!=", "!="},   // ne
}
var allops map[string]bool

func init() {
	allops = make(map[string]bool)
	for _, t := range []map[string]code{v1Tab, v2Tab, cTab} {
		for s := range t {
			allops[s] = true
		}
	}
}

// emit wasm byte code
func (m module) wasm() []c {
	o := bytes.NewBuffer([]c{0, 0x61, 0x73, 0x6d, 1, 0, 0, 0}) // header
	// type section(1: function signatures)
	sec := NewSection(1)
	sigs, sigv := make(map[string]int), make([]string, 0)
	for i, f := range m {
		s := string(f.sig())
		if n, o := sigs[s]; o == false {
			n = len(sigs)
			sigs[s] = n
			sigv = append(sigv, s)
			m[i].sign = n
		}
	}
	sec.cat(lebu(len(sigv)))
	for _, s := range sigv {
		sec.cat([]c(s))
	}
	sec.out(o)
	// no import section(2)
	// function section(3: function signature indexes)
	sec = NewSection(3)
	sec.cat(lebu(len(m)))
	for _, f := range m {
		sec.cat(lebu(sigs[string(f.sig())]))
	}
	sec.out(o)
	// no table section(4)
	// linear memory section(5)
	sec = NewSection(5)
	sec.cat([]c{1, 0, 1}) // 1 initial memory segment, unshared, size 1 block
	sec.out(o)
	// no global section(6)
	// export section(7)
	sec = NewSection(7)
	sec.cat(lebu(len(m))) // number of exports (all)
	for i, f := range m {
		sec.cat(lebu(len(f.name)))
		sec.cat([]c(f.name))
		sec.cat1(0) // function-export
		sec.cat(lebu(i))
	}
	sec.out(o)
	// no start section(8)
	// no element section(9)
	// code section(10)
	sec = NewSection(10)
	sec.cat(lebu(len(m))) // number of functions
	for _, f := range m {
		b := f.code()
		sec.cat(lebu(len(b)))
		sec.cat(b)
	}
	sec.out(o)
	// no data section(11)
	return o.Bytes()
}

type section struct {
	t c
	b []c
}

func NewSection(t c) section { return section{t: t} }
func (s *section) cat(b []c) { s.b = append(s.b, b...) }
func (s *section) cat1(b c)  { s.b = append(s.b, b) }
func (s *section) out(w *bytes.Buffer) {
	w.WriteByte(s.t)
	w.Write(lebu(len(s.b)))
	w.Write(s.b)
}

func (f fn) sig() (r []c) {
	r = append(r, 0x60)
	r = append(r, lebu(len(f.args))...)
	for _, t := range f.args {
		r = append(r, c(t))
	}
	r = append(r, 1)
	r = append(r, c(f.t))
	return r
}
func (f fn) code() (r []c) {
	r = append(r, f.locs()...)
	r = append(r, f.Bytes()...)
	return append(r, 0x0b)
}
func (f fn) locs() (r []c) {
	var u []T
	var n []int
	for i, t := range f.locl {
		if i > 0 && t == f.locl[i-1] {
			n[len(n)-1]++
		} else {
			u, n = append(u, t), append(n, 1)
		}
	}
	r = lebu(len(u))
	for i, t := range u {
		r = append(r, lebu(n[i])...)
		r = append(r, c(t))
	}
	return r
}

func lebu(v int) []c { // encode unsigned leb128
	var b []c
	for {
		c := uint8(v & 0x7f)
		v >>= 7
		if v != 0 {
			c |= 0x80
		}
		b = append(b, c)
		if c&0x80 == 0 {
			break
		}
	}
	return b
}

/*
func lebs(b []b, v int64) []b { // encode signed leb128
	for {
		c := uint8(v & 0x7f)
		s := uint8(v & 0x40)
		v >>= 7
		if (v != -1 || s == 0) && (v != 0 || s != 0) {
			c |= 0x80
		}
		b = append(b, c)
		if c&0x80 == 0 {
			break
		}
	}
	return b
}
*/

func jn(a ...string) string           { return strings.Join(a, "") }
func log(a ...interface{})            { fmt.Fprintln(os.Stderr, a...) }
func logf(f string, a ...interface{}) { fmt.Fprintf(os.Stderr, f, a...) }
func page(wasm []c) []c {
	var b bytes.Buffer
	b.WriteString(head)
	b.WriteString(base64.StdEncoding.EncodeToString(wasm))
	b.WriteString(tail)
	return b.Bytes()
}

func (m module) cout() []c {
	var b bytes.Buffer
	b.WriteString(chead)
	for _, f := range m {
		sig := ""
		for i, a := range f.args {
			if i > 0 {
				sig += ","
			}
			sig += ctyp[a] + " " + "x" + string('0'+byte(i))
		}
		fmt.Fprintf(&b, "%s %s(%s){", ctyp[f.t], f.name, sig)
		if sq, o := f.ast.(seq); o {
			for i, e := range sq {
				if i < len(sq)-1 {
					b.WriteString("return ")
				}
				b.WriteString(cstring(e))
				b.WriteString(";}")
			}
		} else {
			b.WriteString("return ")
			b.WriteString(cstring(f.ast))
			b.WriteString(";}\n")
		}
	}
	return b.Bytes()
}
func (m module) gout() []c {
	var b bytes.Buffer
	b.WriteString(ghead)
	for _, f := range m {
		sig := ""
		for i, a := range f.args {
			if i > 0 {
				sig += ","
			}
			sig += "x" + string('0'+byte(i)) + " " + gtyp[a]
		}
		fmt.Fprintf(&b, "func %s(%s) %s {", f.name, sig, gtyp[f.t])
		if sq, o := f.ast.(seq); o {
			for i, e := range sq {
				if i < len(sq)-1 {
					b.WriteString("return ")
				}
				b.WriteString(gstring(e))
				b.WriteString("}")
			}
		} else {
			b.WriteString("return ")
			b.WriteString(gstring(f.ast))
			b.WriteString("}\n")
		}
	}
	return b.Bytes()
}

const head = `<html>
<head><meta charset="utf-8"><title>w</title></head><body><script>
var us = function(s){var r=new Uint8Array(new ArrayBuffer(s.length));for(var i=0;i<s.length;i++)r[i]=s.charCodeAt(i);return r};
var s = "`
const tail = `"
var u = us(atob(s));
(async() => {
var r=await WebAssembly.instantiate(u)
window.k=r.instance.exports
})()
// browse to file://../index.html
</script>
<pre>
run wasm from js console, e.g:
 k.add(1,2)
</pre>
</body></html>
`

var chead = `#include<stdint.h>
`
var ghead = ``
