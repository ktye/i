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
type s = string
type T c
type fn struct { // name:I:IIF{body}
	name s
	src  [2]int // line, col
	t    T      // return type
	args int
	locl []T
	lmap map[s]int // local index: args+locals
	sign int       // function signature index
	ast  expr
	bytes.Buffer
}
type module []fn

const (
	I = T(0x7f) // i32
	J = T(0x7e) // i64
	F = T(0x7c) // f64
)

var typs = map[c]T{'I': I, 'J': J, 'F': F}
var tnum = map[T]int{I: 0, J: 1, F: 2}
var styp = map[T]s{I: "I", J: "J", F: "F"}
var P = I // I(wasm32) J(wasm64)

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
func (t T) String() s {
	if c := map[T]c{I: 'I', J: 'J', F: 'F'}[t]; c == 0 {
		return "0"
	} else {
		return s(c)
	}
}
func run(r io.Reader) module {
	sFnam, sRety, sArgs, sBody, sCmnt := 0, 1, 2, 3, 4
	rd := bufio.NewReader(r)
	state := sFnam
	line, char := 1, 0
	err := func(s string) { panic(sf("%d:%d: %s", line, char, s)) }
	var m module
	var f fn
	for {
		b, e := rd.ReadByte()
		if e == io.EOF {
			return m.compile()
		} else if e != nil {
			panic(e)
		}
		char++
		if b == '\n' {
			line++
			char = 1
		}
		switch state {
		case sFnam:
			if len(f.name) == 0 && b == ' ' || b == '\t' || b == '\n' {
				continue
			} else if len(f.name) == 0 && b == '/' {
				state = sCmnt
			} else if craZ(b) || (len(f.name) > 0 && cr09(b)) {
				f.name += s(b)
			} else if b == ':' {
				state = sRety
			} else {
				err("parse function name")
			}
		case sRety:
			if f.t == 0 {
				if b == '{' {
					state = sBody // macro
				}
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
			if t := typs[b]; t == 0 && f.locl == nil {
				err("parse args")
			} else if t != 0 {
				f.locl = append(f.locl, t)
				f.args++
			} else if b == ' ' || b == '\t' {
				continue
			} else if b == '{' {
				state = sBody
				f.src = [2]int{line, char}
			} else {
				err("parse args")
			}
		case sBody:
			f.WriteByte(b)
			if b == '}' {
				state = sFnam
				m = append(m, f)
				f = fn{}
			}
		case sCmnt:
			if b == '\n' {
				state = sFnam
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

func (m module) compile() (r module) {
	mac := make(map[s][]c)
	fns := make(map[s]int)
	for _, f := range m {
		_, x := mac[f.name]
		_, y := fns[f.name]
		if x || y {
			panic(f.name + " already defined")
		}
		if f.t == 0 {
			mac[f.name] = f.Bytes()
		} else {
			r = append(r, f)
			fns[f.name] = len(r) - 1
		}
	}
	for i, f := range r {
		r[i].ast = f.parse(mac, fns)
	}
	return r
}

type parser struct {
	mac map[s][]c
	fns map[s]int
	*fn
	p   int
	b   []byte
	tok []byte
}

func (f *fn) parse(mac map[s][]c, fns map[s]int) expr { // parse function body
	f.lmap = make(map[string]int)
	for i := 0; i < f.args; i++ {
		s := s('x' + c(i))
		if i > 2 {
			s = sf("x%d", i)
		}
		f.lmap[s] = i
	}
	p := parser{mac: mac, fns: fns, fn: f, b: strip(f.Bytes())}
	e := p.seq('}')
	e = p.locals(e, 0)
	if x, s := p.validate(e); x != nil {
		return p.xerr(x, s)
	}
	if t := e.rt(); t != p.fn.t {
		return p.err(sf("return type is %s not %s", t, p.fn.t))
	}
	return e
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
func (p *parser) err(s s) expr {
	panic(s)
	return nil
}
func (p *parser) xerr(x expr, s s) expr {
	if i, o := x.(indicator); o {
		return p.indicate(i.indicate(), s)
	} else {
		return p.err(s)
	}
}
func (p *parser) indicate(pos int, e s) expr {
	s := s(p.fn.Bytes())
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
				p.err("missing " + s(term))
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
	if len(seq) > 1 { // suppress assignment expressions
		for i, e := range seq[:len(seq)-1] {
			if v, o := e.(las); o {
				v.tee = 0
				seq[i] = v
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
				return p.err(sf("verb-verb (missing noun) %T %T", x, v))
			} else {
				return p.dyadic(v, x, y, pos(h))
			}
		} else {
			return p.err(sf("noun-noun (missing verb) %T %T", x, v))
		}
	}
}
func (p *parser) monadic(f, x expr, h pos) expr {
	switch v := f.(type) {
	case opx:
		return v1{s: s(v), argv: argv{x}, pos: h}
	default:
		panic("nyi")
	}
}
func (p *parser) dyadic(f, x, y expr, h pos) expr {
	switch v := f.(type) {
	case asn:
		if v.opx == "::" { // memory
			panic("nyi memory asn")
		} else { // local
			if v.opx != ":" { // modified
				y = p.dyadic(opx(v.opx[:len(v.opx)-1]), x, y, h)
			}

			a := las{tee: 1}
			xv, o := x.(loc)
			if o == false {
				return p.xerr(x, "assignment expects a symbol on the left")
			}
			if n, o := p.fn.lmap[xv.s]; o == false {
				xv.i = len(p.fn.lmap)
				p.fn.lmap[xv.s] = xv.i
				p.fn.locl = append(p.fn.locl, 0) // type is still unknown
			} else {
				xv.i = n
			}
			a.argv = []expr{xv, y}
			return a
		}
	case nlp:
		if a, o := y.(las); o { // loop over single assignment
			a.tee = 0
			y = a
		}
		return nlp{pos: h, argv: argv{x, y}}
	case opx:
		if _, o := v2Tab[s(v)]; o {
			return v2{s: s(v), argv: argv{x, y}, pos: h}
		}
		if _, o := cTab[s(v)]; o {
			return cmp{s: s(v), argv: argv{x, y}, pos: h}
		}
		return p.err("unknown operator")
	default:
		panic("nyi")
	}
}
func (p *parser) verb(v expr) bool {
	switch v.(type) {
	case opx, nlp, asn: // todo: others
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
	case p.t(sSym):
		return p.pSym(p.tok)
	case p.t(sC('/')):
		return nlp{}
	case p.t(sCon):
		return p.pCon(p.tok)
	case p.t(sOp):
		e := p.pOp(p.tok)
		if len(p.b) > 0 && p.b[0] == ':' { // w/o space
			p.t(sC(':'))
			return asn{e.(opx)}
		} else if s(e.(opx)) == ":" {
			return asn{e.(opx)}
		}
		return p.pOp(p.tok)
	case p.t(sC('(')):
		return p.seq(')')
	default:
		return nil
	}
}
func (p *parser) locals(e expr, lv int) expr {
	switch l := e.(type) {
	case las:
		l.argv[1] = p.locals(l.argv[1], lv)
		x := l.argv[0].(loc)
		if x.t == 0 {
			x.t = p.locl[x.i]
			l.argv[0] = x
		}
		yt := l.argv[1].rt()
		if yt == 0 {
			return p.xerr(e, "cannot assign zero type")
		} else if x.t == 0 {
			p.locl[x.i] = yt
			x.t = yt
			l.argv[0] = x
		} else if x.t != yt {
			return p.xerr(e, sf("local reassignment of type %s with %s", x.t, yt))
		}
		return l
	case loc:
		if n, o := p.fn.lmap[l.s]; o {
			l.i = n
			l.t = p.locl[n]
		} else {
			return p.xerr(l, "undeclared("+l.s+")")
		}
		return l
	case nlp:
		l.argv[0] = p.locals(l.argv[0], lv+1)
		switch x := l.argv[0].(type) {
		case las:
			l.n = x.argv[0].(loc).i
		case loc:
			l.n = x.i
		default:
			l.n = p.nloc(s('i'+lv)+"n", I) // create limit in jn ..
		}
		l.c = p.nloc(s('i'+lv), I) // set/create loop counter
		l.argv[1] = p.locals(l.argv[1], lv+1)
		return l
	default:
		if av, o := e.(argvec); o {
			v := av.args()
			for i, a := range v {
				v[i] = p.locals(a, lv)
			}
		}
		return e
	}
}
func (p *parser) validate(e expr) (expr, s) {
	if av, o := e.(argvec); o {
		for _, e := range av.args() {
			if r, s := p.validate(e); r != nil {
				return r, s
			}
		}
	}
	s := e.valid()
	if s != "" {
		return e, s
	}
	return nil, ""
}
func (p *parser) nloc(s s, t T) int { // local index by name, may create new
	n, o := p.fn.lmap[s]
	if o {
		if p.fn.locl[n] != t {
			p.err(s + " exists with different type")
		}
	} else {
		n = len(p.fn.lmap)
		p.fn.lmap[s] = n
		p.fn.locl = append(p.fn.locl, t)
	}
	return n
}
func sSym(b []c) int { // [aZ][a9]*
	c := b[0]
	if craZ(c) == false {
		return 0
	}
	for i, c := range b {
		if craZ(c) == false && cr09(c) == false {
			return i
		}
	}
	return len(b)
}
func (p *parser) pSym(b []c) expr { return loc{pos: pos(p.p), s: s(b), i: -1} }
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
		if f, e := strconv.ParseFloat(s(b), 64); e != nil {
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
	if i, e := strconv.ParseInt(s(b), 10, 64); e != nil {
		return p.err(e.Error())
	} else {
		r.i = i
	}
	return r
}
func sOp(b []c) int {
	if b[0] == ':' {
		return 1
	}
	for _, n := range []int{3, 2, 1} { // longest match first
		if len(b) >= n && allops[s(b[:n])] {
			return n
		}
	}
	return 0
}
func (p *parser) pOp(b []c) expr { return opx(s(b)) }
func sC(x c) func(b []c) int     { return func(b []c) int { return boolvar(b[0] == x) } }

// intermediate representation for function bodies (typed expression tree)
type expr interface {
	rt() T    // result type, maybe 0
	valid() s // ok("") or err
	bytes() []c
}
type argvec interface {
	args() []expr
}
type argv []expr
type cstringer interface {
	cstr() s
}
type gstringer interface {
	gstr() s
}
type seq argv    // a;b;..
type v2 struct { // x+y unitype
	pos
	argv
	s s // +-*%
}
type v1 struct { // -y
	pos
	argv
	s s
	p int
}
type cmp struct { // x<y..
	pos
	argv
	s s
}
type con struct { // numeric constant
	pos
	t T
	i int64
	f float64
}
type loc struct { // local get
	pos
	t T
	s s
	i int
}
type las struct { // local set
	// loc
	argv
	tee c // 01
}
type nlp struct { // x/y loop
	pos
	argv
	n int // index in locl for loop limit
	c int // index in locl for loop counter
}
type opx s             // operator
type asn struct{ opx } // assignments :(local) ::(memory) +:(modified local)
type pos int           // src position
type indicator interface {
	indicate() int
}

func (p pos) indicate() int { return int(p) }

func getop(tab map[s]code, op s, t T) (r c) {
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
func cop(tab map[s]code, op s, t T) (o, u s) {
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
			u = "(u" + styp[t] + ")"
		}
		o = o[1:]
	}
	return o, u
}
func gop(tab map[s]code, op s, t T) (o, u s) {
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
			u = "u" + styp[t]
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
func (s seq) valid() s { // all but the last expressions in a sequence must have no return type
	for i, e := range s {
		if t := e.rt(); i < len(s)-1 && t != 0 {
			return sf("statement %d/%d has nonzero type %s", i+1, len(s), t)
		} else if i == len(s)-1 && t == 0 {
			return sf("last statement of %d has zero type", i+1)
		}
	}
	return ""
}
func (s seq) bytes() (r []c) {
	for _, e := range s {
		r = append(r, e.bytes()...)
	}
	return r
}
func (v v2) rt() T { return v.x().rt() }
func (v v2) valid() s {
	if tx, ty := v.x().rt(), v.y().rt(); tx == 0 {
		return sf("left argument has zero type")
	} else if ty == 0 {
		return sf("right argument has zero type")
	} else if tx != ty {
		return sf("types mismatch %s %s", tx, ty)
	}
	return ""
}
func (v v2) bytes() []c {
	return append(append(v.x().bytes(), v.y().bytes()...), getop(v2Tab, v.s, v.rt()))
}
func (v v2) cstr() s    { return c2str(v2Tab, v.s, v.rt(), v.x(), v.y()) }
func (v v2) gstr() s    { return g2str(v2Tab, v.s, v.rt(), v.x(), v.y()) }
func (v v1) rt() T      { return v.x().rt() }
func (v v1) valid() s   { return ifex(v.x().rt() == 0, "argument has zero type") }
func (v v1) bytes() []c { return append(v.x().bytes(), getop(v1Tab, v.s, v.rt())) }
func (v v1) cstr() s    { o, u := cop(v1Tab, v.s, v.rt()); return jn(o, "(", u, cstring(v.x()), ")") }
func (v v1) gstr() s {
	o, u := gop(v1Tab, v.s, v.rt())
	return jn(o, u, "((", gstring(v.x()), "))")
}
func (v cmp) rt() T    { return I }
func (v cmp) valid() s { return v2(v).valid() }
func (v cmp) bytes() []c {
	return append(append(v.x().bytes(), v.y().bytes()...), getop(cTab, v.s, v.rt()))
}
func (v cmp) cstr() s  { return c2str(cTab, v.s, v.rt(), v.x(), v.y()) }
func (v cmp) gstr() s  { return g2str(cTab, v.s, v.rt(), v.x(), v.y()) }
func (v con) rt() T    { return v.t }
func (v con) valid() s { return ifex(v.t == 0, "constant has zero type") }
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
func (v con) cstr() s {
	if v.t == F {
		s := sf("%v", v.f)
		if strings.Index(s, ".") == -1 {
			s += ".0"
		}
		return s
	}
	return sf("%d", v.i)
}
func (v con) gstr() s    { return v.cstr() }
func (v loc) rt() T      { return v.t }
func (v loc) valid() s   { return ifex(v.t == 0, "local has zero type") }
func (v loc) bytes() []c { return append([]c{0x20}, lebu(v.i)...) }
func (v loc) cstr() s    { return locstr(v) }
func (v loc) gstr() s    { return locstr(v) }
func (v las) rt() T      { return T(v.tee) * v.y().rt() }
func (v las) valid() s {
	tx, ty := v.x().rt(), v.y().rt()
	return ifex(tx == 0 || tx != ty, sf("assignment with mismatched types %s %s", tx, ty))
}
func (v las) bytes() []c {
	return append(v.y().bytes(), append([]c{0x21 + v.tee}, lebu(v.x().(loc).i)...)...)
}
func (v las) cstr() s { return jn("(", locstr(v.x()), "=", s(v.y().bytes()), ")", s(59-27*v.tee)) }
func (v las) gstr() s {
	if v.tee > 0 {
		return jn("as", styp[v.rt()], "(&", locstr(v.x()), ",", s(v.y().bytes()), ")")
	}
	return jn(locstr(v.x()), "=", s(v.y().bytes()), ";")
}
func (v nlp) rt() T { return 0 }
func (v nlp) valid() s {
	if xt, yt := v.x().rt(), v.y().rt(); xt != I {
		return sf("loop range is not I: %s", xt)
	} else if yt != 0 {
		return sf("loop body has nonzero type %s", yt)
	}
	return ""
}
func (v nlp) bytes() (r []c) {
	i, n := s(lebu(v.c)), s(lebu(v.n))
	//         x                           0  !=  if           0   →i   loop
	r = catb(v.x().bytes(), []c(sf("\x41\x00\x47\x04\x40\x41\x00\x21%s\x03\x40", i)))
	//                                        i       1  +  tee→i    n   <  continue
	return catb(r, v.y().bytes(), []c(sf("\x20%s\x41\x01\x6a\x21%s\x20%s\x49\x0d\x00\x0b", i, i, n)))
} // TODO
func (v nlp) cstr() s {
	return sf("for(x%d=0;x%d<(%s);x%d++){%s}", v.c, v.c, cstring(v.x()), v.c, cstring(v.y()))
}
func (v nlp) gstr() s {
	return sf("for x%d=0;x%d<(%s);x%d++{%s}", v.c, v.c, gstring(v.x()), v.c, gstring(v.y()))
}
func (v opx) rt() T      { return 0 }
func (v opx) valid() s   { return "nonapplied operator" }
func (v opx) bytes() []c { return nil }
func locstr(v expr) s    { return sf("x%d", v.(loc).i) }
func ifex(c bool, s s) s {
	if c {
		return s
	}
	return ""
}

type code struct {
	I, J, F c
	c, g    s
}

func c2str(tab map[s]code, op s, t T, x, y expr) s {
	o, u := cop(tab, op, t)
	if len(o) > 2 {
		return jn(u, o, "(", cstring(x), ",", cstring(y), ")")
	} else {
		return jn("((", u, cstring(x), ")", o, "(", u, cstring(y), "))")
	}
}
func g2str(tab map[s]code, op s, t T, x, y expr) s {
	o, u := cop(tab, op, t)
	u += "("
	if len(o) > 2 {
		return jn(u, o, "(", cstring(x), ",", cstring(y), "))")
	} else {
		return jn("((", u, cstring(x), "))", o, "(", u, cstring(y), ")))")
	}
}
func cstring(x expr) s { xs := x.(cstringer); return xs.cstr() }
func gstring(x expr) s { xs := x.(gstringer); return xs.gstr() }

var v1Tab = map[s]code{
	"-": code{0, 0, 0x9a, "-", "-"},                                                                            // neg (no neg for ints)
	"+": code{0, 0, 0x99, "fabs", "math.Abs"},                                                                  // abs
	"_": code{1, 1, 0x9c, ";;floor", "math.Floor"},                                                             // floor (ceil, trunc, nearest?)
	"*": code{0x67, 0x79, 0, "__builtin_clz;__builtin_clzll;", "Ubits.LeadingZeros32;Ubits.LeadingZeros64;"},   // clz
	"|": code{0x68, 0x79, 0, "__builtin_ctz;__builtin_ctzll;", "Ubits.TrailingZeros32;Ubits;TrailingZeros64;"}, // ctz
	"%": code{0, 0, 0x9f, "sqrt", "math.Sqrt"},                                                                 // sqr
}
var v2Tab = map[s]code{
	`+`:   code{0x6a, 0x7c, 0xa0, "+", "+"},     // add
	`-`:   code{0x6b, 0x7d, 0xa1, "-", "-"},     // sub
	`*`:   code{0x6c, 0x7e, 0xa2, "*", "*"},     // mul
	`%`:   code{0x6e, 0x80, 0xa3, "U/", "U/"},   // div/div_u
	`%'`:  code{0x6d, 0x7f, 0xa3, "/", "/"},     // div_s
	`\`:   code{0x70, 0x82, 0, "U%", "%U"},      // rem_u
	`\'`:  code{0x6f, 0x81, 0, "%", "%"},        // rem_s
	`&`:   code{0x71, 0x83, 0, "&", "&"},        // and
	`|`:   code{0x72, 0x84, 0, "|", "|"},        // or
	`^`:   code{0x73, 0x85, 0, "^", "^"},        // xor
	`<<`:  code{0x74, 0x86, 0, "<<", "<<"},      // shl
	`>>`:  code{0x76, 0x88, 0, "U>>", "U>>"},    // shr_u
	`>>'`: code{0x75, 0x87, 0, ">>", ">>"},      // shl_s
	`<|'`: code{0x77, 0x89, 0, "", ""},          // rotl
	`>|'`: code{0x78, 0x8a, 0, "", ""},          // rotr
	`&'`:  code{0, 0, 0xa4, "fmin", "math.Max"}, // min
	`|'`:  code{0, 0, 0xa5, "fmax", "math.Min"}, // max
}
var cTab = map[s]code{
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
var allops map[s]bool

func init() {
	allops = make(map[s]bool)
	for _, t := range []map[s]code{v1Tab, v2Tab, cTab} {
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
	sigs, sigv := make(map[s]int), make([]s, 0)
	for i, f := range m {
		s := s(f.sig())
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
		sec.cat(lebu(sigs[s(f.sig())]))
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
	r = append(r, lebu(f.args)...)
	for i := 0; i < f.args; i++ {
		r = append(r, c(f.locl[i]))
	}
	r = append(r, 1)
	r = append(r, c(f.t))
	return r
}
func (f fn) code() (r []c) {
	r = append(r, f.locs()...)
	r = append(r, f.ast.bytes()...)
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

func catb(x ...[]c) (r []c) {
	for _, b := range x {
		r = append(r, b...)
	}
	return r
}
func jn(a ...s) s                { return strings.Join(a, "") }
func log(a ...interface{})       { fmt.Fprintln(os.Stderr, a...) }
func logf(f s, a ...interface{}) { fmt.Fprintf(os.Stderr, f, a...) }
func sf(f s, a ...interface{}) s { return fmt.Sprintf(f, a...) }
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
		for i := 0; i < f.args; i++ {
			if i > 0 {
				sig += ","
			}
			sig += styp[f.locl[i]] + " " + "x" + s('0'+byte(i))
		}
		fmt.Fprintf(&b, "%s %s(%s){", styp[f.t], f.name, sig)
		if sq, o := f.ast.(seq); o {
			for i, e := range sq {
				if i < len(sq)-1 {
					b.WriteString("R ")
				}
				b.WriteString(cstring(e))
				b.WriteString(";}")
			}
		} else {
			b.WriteString("R ")
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
		for i := 0; i < f.args; i++ {
			if i > 0 {
				sig += ","
			}
			sig += "x" + s('0'+byte(i)) + " " + styp[f.locl[i]]
		}
		fmt.Fprintf(&b, "func %s(%s) %s {", f.name, sig, styp[f.t])
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

var chead = ``
var ghead = `
type I=int32;type J=int64;type F=float64
func asI(x *I,y I)I{*x=y;return y};func lsJ(x *J,y J)J{*x=y;return y};func asF(x *F,y F)F{*x=y;return y};
`
