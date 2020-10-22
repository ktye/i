package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type c = byte
type s = string
type T c
type fn struct { // name:I:IIF{body} or name.I:IIF{..} (unexported)
	name s
	ex   bool
	data bool   // data section (not a function)
	src  [2]int // line, col
	t    T      // return type
	args int
	locl []T
	lmap map[s]int // local index: args+locals
	sign int       // function signature index
	ast  expr
	bytes.Buffer
}
type sig struct {
	t T
	a []T
}
type module []fn

const (
	V = T(0x00)
	C = T(0x01) // i8
	I = T(0x7f) // i32
	J = T(0x7e) // i64
	F = T(0x7c) // f64
)

var typs = map[c]T{'V': V, 'C': C, 'I': I, 'J': J, 'F': F}
var tnum = map[T]int{C: 0, I: 0, J: 1, F: 2} // same op for C and I
var styp = map[T]s{V: "V", C: "C", I: "I", J: "J", F: "F"}
var alin = map[T]c{C: 0, I: 2, J: 3, F: 3}

func main() {
	var stdin io.Reader = os.Stdin
	var cout, gout bool
	flag.BoolVar(&cout, "c", false, "c output")
	flag.BoolVar(&gout, "go", false, "go output")
	flag.Parse()
	m, tab, data := run(stdin)
	if cout {
		os.Stdout.Write(m.cout(tab, data))
	} else if gout {
		os.Stdout.Write(m.gout(tab, data))
	} else {
		os.Stdout.Write(m.wasm(tab, data))
	}
}
func (t T) String() s { return styp[t] }
func run(r io.Reader) (module, []segment, []dataseg) {
	sFnam, sRety, sArgs, sBody, sData, sCmnt := 0, 1, 2, 3, 4, 5
	rd := bufio.NewReader(r)
	state := sFnam
	line, char := 1, 0
	err := func(s string) { panic(sf("%d:%d: %s", line, char, s)) }
	var m module
	var f fn
	var datas []dataseg
	var data []c
	decode := func(data []c) []c {
		r, e := hex.DecodeString(string(data))
		if e != nil {
			err("data section is no valid hex")
		}
		return r
	}
	for {
		b, e := rd.ReadByte()
		if e == io.EOF || (state == sFnam && b == '\\') {
			mod, tab := m.compile()
			return mod, tab, datas
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
			} else if craZ(b) || cr09(b) {
				f.name += s(b)
			} else if b == '.' {
				state = sRety
			} else if b == ':' {
				state = sRety
				f.ex = true
			} else if b == '!' {
				state = sData
				f.data = true
			} else {
				fmt.Printf("%s\n", string(b))
				err("parse function name")
			}
		case sRety:
			if b == ':' {
				state = sArgs
				continue
			} else if f.t != 0 {
				err("parse return type")
			}
			if b == '{' && f.t == 0 {
				state = sBody // macro
				continue
			}
			f.t = typs[b]
			if f.t == 0 && b != 'V' {
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
		case sData:
			if b == '}' {
				s := dataseg{}
				if off, e := strconv.Atoi(f.name); e != nil {
					err("data section name must be an integer, not: " + f.name)
				} else {
					s.off = off
					s.bytes = decode(data)
				}
				datas = append(datas, s)
				data = nil
				f = fn{}
				state = sFnam
			} else if b != '{' {
				data = append(data, b)
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
func (s sig) String() s {
	r := s.t.String() + ":"
	for i := range s.a {
		r += s.a[i].String()
	}
	return r
}
func (m module) compile() (r module, tab []segment) {
	mac := make(map[s][]c)
	fns := make(map[s]int)
	var fsg []sig
	for _, f := range m {
		if cr09(f.name[0]) {
			continue
		}
		_, x := mac[f.name]
		_, y := fns[f.name]
		if x || y {
			panic(f.name + " already defined")
		}
		if f.args == 0 {
			b := f.Bytes()
			if strings.HasSuffix(f.name, "T") { // type-generic-macro fnT→fnC|fnI|... W:..
				for c, v := range typs {
					name := f.name[:len(f.name)-1] + string(c) // fnC fnI ..
					widt := 1 << alin[v]
					body := string(b[:len(b)-1]) // strip '}'
					body = strings.Replace(body, "T", string(c), -1)
					body = strings.Replace(body, "W", strconv.Itoa(widt), -1) // W→1 4 8 ..
					mac[name] = []byte(body)
				}
			} else {
				mac[f.name] = b[:len(b)-1] // strip '}'
			}
		} else {
			r = append(r, f)
			n := len(r) - 1
			fns[f.name] = n
			sg := make([]T, f.args)
			copy(sg, f.locl)
			fsg = append(fsg, sig{t: f.t, a: sg})
		}
	}
	sgm := make(map[string]int)
	for _, sig := range fsg {
		s := sig.String()
		if _, o := sgm[s]; !o {
			sgm[s] = len(sgm)
		}
	}
	for i, f := range r {
		f.ast = f.parse(mac, fns, fsg, sgm)
		r[i] = f
	}
	for _, f := range m {
		if !cr09(f.name[0]) {
			continue
		}
		tab = append(tab, f.parseTab(fns))
	}
	return r, tab
}

type parser struct {
	mac map[s][]c
	fns map[s]int
	fsg []sig
	sgm map[s]int
	*fn
	p   int
	exp map[int]int
	b   []byte
	tok []byte
}

func (f *fn) parse(mac map[s][]c, fns map[s]int, fsg []sig, sgm map[s]int) expr { // parse function body
	f.lmap = make(map[string]int)
	for i := 0; i < f.args; i++ {
		s := s('x' + c(i))
		if i > 2 {
			s = sf("x%d", i)
		}
		f.lmap[s] = i
	}
	//fmt.Printf("parse %s\n", f.name)
	p := parser{mac: mac, fns: fns, fsg: fsg, sgm: sgm, fn: f, b: strip(f.Bytes()), exp: make(map[int]int)}
	e := p.seq('}')
	if e == nil {
		return nil
	}
	e = p.locals(e, 0)
	if x, s := p.validate(e); x != nil {
		println(f.name)
		return p.xerr(x, s)
	}
	if t := e.rt(); t != p.fn.t {
		if !(t == 0 && p.fn.t == 255) {
			return p.err(sf("%s: return type is %s not %s", f.name, t, p.fn.t))
		}
	}
	return e
}
func (p parser) pos() (r pos) {
	r = pos(p.p)
	for x, y := range p.exp {
		if x < p.p {
			r -= pos(y)
		}
	}
	if r < 0 {
		r = 0
	}
	return r
}
func (f *fn) parseTab(fns map[s]int) (tab segment) { // function table: 8:{f;g;h}
	var e error
	tab.off, e = strconv.Atoi(f.name)
	if e != nil {
		panic("illegal function name: " + f.name)
	}
	v := bytes.Split(bytes.TrimSuffix(f.Bytes(), []c{'}'}), []c{';'})
	tab.names = make([]string, len(v))
	for i := range v {
		tab.names[i] = strings.TrimSpace(s(v[i]))
	}
	return tab
}
func strip(b []c) []c { // strip comments
	lines := bytes.Split(b, []c{'\n'})
	for i, l := range lines {
		for k, c := range l {
			if c == '/' && (k == 0 || l[k-1] == ' ') {
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
				//pos--
			}
			return p.err("\n" + l + "\n" + strings.Repeat(" ", pos) + "^" + e)
		}
		pos -= len(l) + 1
	}
	if len(s) > 0 {
		return p.err("\n" + s + "\n^" + e)
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
	seq.pos = pos(p.p)
	for {
		e := p.ex(p.noun())
		if e != nil {
			seq.argv = append(seq.argv, e)
		} else {
			p.w()
			if len(p.b) == 0 {
				p.xerr(seq, "missing "+s(term))
			}
			if p.b[0] == term {
				p.b = p.b[1:]
				break
			} else if p.b[0] != ';' {
				p.xerr(seq, "expected ;")
			} else {
				p.b = p.b[1:]
			}
		}
	}
	if seq.argv == nil {
		return nil // empty?
	} else if len(seq.argv) == 1 {
		return seq.argv[0]
	}
	return seq
}
func (p *parser) ex(x expr) expr {
	if x == nil {
		return x
	}
	h := p.pos()
	v := p.noun()
	if op, o := x.(opx); o && s(op) == "-" { // fix neg. numbers
		if c, o := v.(con); o {
			c.i = -c.i
			c.f = -c.f
			x = c
			v = p.noun()
		}
	}
	if p.verb(x) {
		if y := p.ex(v); y == nil {
			if t, o := x.(opx); o && s(t) == "!" {
				return p.pTrp()
			}
			return x // verb ?
		} else {
			return p.monadic(x, y, h)
		}
	} else {
		if v == nil {
			return x // noun
		} else if p.verb(v) {
			if y := p.ex(p.noun()); y == nil {
				return p.xerr(p.pos(), sf("verb-verb (missing noun) x=%#v v=%#v", x, v))
			} else {
				return p.dyadic(v, x, y, h)
			}
		} else if t, o := x.(typ); o {
			y := p.ex(v)
			if y == nil {
				y = v
			}
			return lod{t: t.t, argv: argv{y}, pos: pos(h)} // I x
		} else if d, o := x.(dot); o {
			if s, o := v.(seq); o {
				d.argv = s.argv
			} else {
				d.argv = argv{v}
			}
			return d
		} else {
			return p.xerr(pos(h), sf("noun-noun (missing verb) %#v %#v", x, v))
		}
	}
}
func (p *parser) monadic(f, x expr, h pos) expr {
	switch v := f.(type) {
	case opx:
		if s(v) == "?" { // ?x
			return brif{argv: argv{x}, pos: h}
		}
		return v1{s: s(v), argv: argv{x}, pos: h}
	case fun:
		if s, o := x.(seq); o {
			return cal{fun: v, argv: s.argv}
		}
		return cal{fun: v, argv: argv{x}}
	case asn:
		return ret{argv: argv{x}, pos: h}
	default:
		fmt.Printf("h=%d\nf=%#v\nx=%#v\n", h, f, x)
		panic("nyi")
	}
}
func (p *parser) dyadic(f, x, y expr, h pos) expr {
	switch v := f.(type) {
	case asn:
		if v.opx == ":" { // memory
			return sto{argv: argv{x, y}, t: y.rt(), pos: h}
		} else { // local
			if v.opx != "" { // modified
				if v.opx == "?" {
					return p.xerr(v, "?: illegal modified assignment (missing space?)")
				}
				y = p.dyadic(opx(v.opx), x, y, h)
			}
			a := las{pos: h}
			xv, o := x.(loc)
			if o == false {
				fmt.Printf("%s\n", string(p.b))
				fmt.Printf("f=%#v\nx=%#v\ny=%#v\n", f, x, y)
				return p.xerr(a, fmt.Sprintf("assignment expects a symbol on the left: not %T", x))
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
		if a, o := x.(con); o && a.t == I && a.i == 1 { // 1/
			return whl{pos: h, argv: argv{x, y}}
		}
		return nlp{pos: h, argv: argv{x, y}}
	case opx:
		if _, o := v2Tab[s(v)]; o {
			return v2{s: s(v), argv: argv{x, y}, pos: h}
		}
		if _, o := cTab[s(v)]; o {
			return cmp{s: s(v), argv: argv{x, y}, pos: h}
		}
		if s(v) == "." {
			d := dot{pos: h, idx: y}
			if t, o := x.(typ); o == false {
				return p.xerr(d, "dot call requires type on the left")
			} else {
				d.t = t.t
			}
			return d
		}
		if s(v) == "?" || s(v) == "?/" || s(v) == "?'" {
			if xt, o := x.(typ); o {
				sn := 0
				if s(v) == "?'" {
					sn = 1
				}
				return cvt{t: xt.t, argv: argv{y}, pos: h, sign: sn}
			}
			if s(v) == "?" {
				return iff{argv: argv{x, y}, pos: h}
			} else {
				return whl{argv: argv{x, y}, pos: h}
			}
		}
		return p.err("unknown operator(" + s(v) + ")")
	case fun:
		return cal{fun: v, argv: argv{x, y}}
	default:
		panic("nyi")
	}
}
func (p *parser) verb(v expr) bool {
	switch v.(type) {
	case opx, nlp, asn, fun: // todo: others
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
	case p.t(sC('(')):
		return p.seq(')')
	case p.t(sTyp):
		return p.pTyp(p.tok)
	case p.t(sSym):
		if mc, o := p.mac[s(p.tok)]; o { // macro-expansion
			p.b = append(mc, p.b...)
			p.exp[p.p] = len(mc) - len(s(p.tok))
			return p.noun()
		}
		return p.pSym(p.tok)
	case p.t(sC('/')):
		return nlp{}
	case p.t(sCon):
		return p.pCon(p.tok)
	case p.t(sOp):
		e := p.pOp(p.tok)
		if len(p.b) > 0 && p.b[0] == ':' { // w/o space
			p.t(sC(':'))
			t := T(0)
			if len(p.b) > 0 && p.b[0] == '\'' { // ::' (i32.store_8) //todo remove
				p.t(sC('\''))
				t = C
			}
			return asn{e.(opx), t}
		} else if s(e.(opx)) == ":" {
			return asn{opx(""), 0}
		}
		return e
	case p.t(sCnd):
		s := p.seq(']').(seq)
		return cnd(s)
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
			fmt.Fprintf(os.Stderr, "%#v\n", l)
			return p.xerr(e, "cannot assign zero type")
		}
		if x.t == 0 {
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
		if l.t == 0 {
			l.t = p.locl[l.i]
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
			l.n = p.nloc(s(byte('i'+lv))+"n", I) // create limit in jn ..
		}
		l.c = p.nloc(s(byte('i'+lv)), I) // set/create loop counter
		l.argv[1] = p.locals(l.argv[1], lv+1)
		return l
	case v2:
		l.argv[0] = p.locals(l.argv[0], lv)
		l.argv[1] = p.locals(l.argv[1], lv)
		t := l.rt()
		for i := 0; i < 2; i++ {
			if x, o := l.argv[i].(loc); o {
				if x.t == 0 {
					x.t = t // uninitialized local
					l.argv[i] = x
				}
			}
		}
		return l
	case sto:
		l.argv[0] = p.locals(l.argv[0], lv)
		l.argv[1] = p.locals(l.argv[1], lv)
		if l.t == 0 {
			l.t = l.argv[1].rt()
		}
		return l
	case dot:
		tv := make([]T, len(l.argv))
		for i, a := range l.argv {
			l.argv[i] = p.locals(a, lv)
			tv[i] = l.argv[i].rt()
		}
		l.idx = p.locals(l.idx, lv)
		s := sig{t: l.t, a: tv}.String()
		if idx, o := p.sgm[s]; o {
			l.sig = idx
		} else {
			return p.xerr(e, sf("unknown function signature for indirect dot call %s", s))
		}
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
		if tt := p.fn.locl[n]; t != tt {
			fmt.Printf("%v %v %#v %#v\n", tt, t, p.fn.locl, p.fn.lmap)
			p.err(s + " exists with different type")
		}
	} else {
		n = len(p.fn.lmap)
		p.fn.lmap[s] = n
		p.fn.locl = append(p.fn.locl, t)
	}
	return n
}
func sTyp(b []c) int { // C I J F
	if _, o := typs[b[0]]; o == false {
		return 0
	}
	if len(b) > 0 && (craZ(b[1]) || cr09(b[1])) {
		return 0
	}
	return 1
}
func (p *parser) pTyp(b []c) expr { return typ{t: typs[b[0]]} }
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
func (p *parser) pSym(b []c) expr {
	if n, o := p.fns[s(b)]; o {
		return fun{s: s(b), n: n, sig: p.fsg[n], pos: p.pos()}
	}
	return loc{pos: pos(p.p), s: s(b), i: -1}
}
func (p *parser) pTrp() expr {
	n, o := p.fns["trap"]
	if !o {
		return opx("!") // no trap function defined
	}
	return seq{pos: p.pos(), argv: argv{cal{fun: fun{s: "trap", n: n, sig: p.fsg[n], pos: p.pos()}, argv: argv{con{t: I, i: int64(p.src[0])}, con{t: I, i: int64(p.src[1]) + int64(p.pos())}}}, opx("!")}}
}
func (p *parser) pDot(b []c) expr {
	return loc{pos: pos(p.p), s: s(b), i: -1}
}
func sCon(b []c) int { // 123 123i 123j 123f .123 123. -..
	dot := false
	if !cr09(b[0]) {
		return 0
	}
	for i, c := range b {
		if cr09(c) {
			continue
		} else if dot == false && (c == 'i' || c == 'j' || c == 'f') {
			return i + 1
		} else if dot == false && c == '.' {
			dot = true
		} else {
			return i
		}
	}
	return len(b)
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
	if c := b[len(b)-1]; c == 'i' || c == 'j' || c == 'f' {
		b = b[:len(b)-1]
		if c == 'j' {
			r.t = J
		} else if c == 'f' {
			r.t = F
		}
	}
	if u, e := strconv.ParseUint(s(b), 10, 64); e != nil {
		return p.err(e.Error())
	} else {
		r.i = int64(u)
		if r.t == F {
			r.f = math.Float64frombits(uint64(r.i))
		}
	}
	return r
}
func sOp(b []c) int {
	if b[0] == '?' && len(b) > 1 && (b[1] == '/' || b[1] == '\'') { // ?/(while)  ?'(cvt signed)
		return 2
	} else if b[0] == ':' || b[0] == '?' || b[0] == '.' {
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
func sCnd(b []c) int {
	if len(b) > 1 && b[0] == '$' && b[1] == '[' {
		return 2
	}
	return 0
}
func sC(x c) func(b []c) int { return func(b []c) int { return boolvar(b[0] == x) } }

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
type seq struct { // a;b;..
	pos
	argv
}
type cnd struct { // $[a;b;..]
	pos
	argv
}
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
type typ struct { // type C I J F
	pos
	t T
}
type con struct { // numeric constant
	pos
	t T
	i int64
	f float64
}
type cvt struct { // J? convert
	pos
	argv
	t    T
	sign int
}
type fun struct { // f
	pos
	s     s
	n     int
	sig   sig
	indir bool
}
type cal struct { // f x
	fun
	argv
}
type dot struct { // call indirect
	pos
	argv      // function arguments
	t    T    // function return type
	idx  expr // function index
	sig  int  // function signature index
}
type loc struct { // local get
	pos
	t T
	s s
	i int
}
type las struct { // local set
	pos
	argv
}
type sto struct { // x::y
	pos
	argv
	t T
}
type lod struct { // I x  (I'x signed)
	pos
	argv
	t T
}
type ret struct { // :x (return)
	pos
	argv
}
type iff struct { // x?y
	pos
	argv
}
type nlp struct { // x/y loop
	pos
	argv
	n int // index in locl for loop limit
	c int // index in locl for loop counter
}
type whl struct { // x?/y while  1/ while(1)
	pos
	argv
}
type brif struct { // ?x
	pos
	argv
}
type opx s        // operator
type asn struct { // assignments :(local) ::(memory) +:(modified local)
	opx
	t T // C(::')
}
type nop struct{}
type pos int // src position
type indicator interface {
	indicate() int
}

func (p pos) indicate() int { return int(p) }

func getop(tab map[s]code, op s, t T) (r c) {
	ops, ok := tab[op]
	if !ok {
		panic("unknown operator: " + op)
	}
	switch t {
	case C, I:
		r = ops.I
	case J:
		r = ops.J
	case F:
		r = ops.F
	default:
		panic("type(" + op + ")")
	}
	if r == 0 {
		panic("type(" + op + ")")
	}
	return r
}
func cop(tab map[s]code, op s, t T) (o, sn s) {
	ops, ok := tab[op]
	if !ok {
		panic("type")
	}
	o = ops.c
	if strings.Index(o, ";") != -1 {
		v := strings.Split(o, ";")
		o = v[tnum[t]]
	}
	if o[0] == 'S' {
		if t == I {
			sn = "(SI)"
		} else if t == J {
			sn = "(SJ)"
		}
		o = o[1:]
	}
	return o, sn
}
func gop(tab map[s]code, op s, t T) (o, sn s) {
	ops, ok := tab[op]
	if !ok {
		panic("type")
	}
	o = ops.g
	if strings.Index(o, ";") != -1 {
		v := strings.Split(o, ";")
		o = v[tnum[t]]
	}
	if o[0] == 'S' {
		if t == I {
			sn = "SI"
		} else if t == J {
			sn = "SJ"
		}
		o = o[1:]
	}
	return o, sn
}

func (a argv) args() []expr { return a }
func (a argv) x() expr      { return a[0] }
func (a argv) y() expr      { return a[1] }
func (s seq) rt() T         { return s.argv[len(s.argv)-1].rt() }
func (s seq) valid() s { // all but the last expressions in a sequence must have no return type
	for i, e := range s.argv {
		if t := e.rt(); i < len(s.argv)-1 && t != 0 {
			return sf("statement %d/%d has nonzero type %s: %#v", i+1, len(s.argv), t, e)
		}
	}
	return ""
}
func (s seq) bytes() (r []c) {
	for _, e := range s.argv {
		r = append(r, e.bytes()...)
	}
	return r
}
func (s seq) cstr() (r s) {
	t := s.rt()
	if t != 0 {
		r = "("
		for i, a := range s.argv {
			r += cstring(a)
			if i < len(s.argv)-1 {
				r += ","
			}
		}
		return r + ")"
	}
	for _, a := range s.argv {
		r += cstring(a)
	}
	return r
}
func (s seq) gstr() (r s) {
	t := s.rt()
	if t != 0 {
		r = "func()" + styp[t] + "{"
	}
	for i, a := range s.argv {
		if i == len(s.argv)-1 && t != 0 {
			r += "return "
		}
		r += gstring(a)
		if i < len(s.argv)-1 {
			r += "\n"
		}
	}
	if t != 0 {
		r += "}()"
	}
	return r
}
func (v cnd) rt() T { return v.argv[len(v.argv)-1].rt() }
func (v cnd) valid() s {
	n := len(v.argv)
	if n < 3 || n%2 == 0 { // only odd are allowed (with else statement)
		return sf("conditional $[..] has wrong number of cases: %d", n)
	}
	rt := v.rt()
	for i := 0; i < n-1; i += 2 {
		if t := v.argv[i].rt(); t != I {
			return sf("[%d]conditional must be I (%s)", i, t)
		}
	}
	for i := 1; i < n; i += 2 {
		if t := v.argv[i].rt(); t != rt {
			fmt.Println(t, "!=", rt)
			return sf("conditional has mixed types")
		}
	}
	return ""
}
func (v cnd) bytes() (r []c) {
	//t := v.rt()
	a := v.argv
	//if t == 0 && len(a) == 3 {
	//	return catb(a[0].bytes(), []c{0x04, 0x40}, a[1].bytes(), []c{0x05}, a[2].bytes(), []c{0x0b}) // if(a0){a1}else{a2}
	//}
	for i := 0; i < len(a)-1; i += 2 {
		//r = catb(r, a[i].bytes(), []c{0x04, c(t)}, a[i+1].bytes(), []c{0x05})
		r = catb(r, a[i].bytes(), []c{0x04, 0x40}, a[i+1].bytes(), []c{0x05})
	}
	return catb(r, a[len(a)-1].bytes(), bytes.Repeat([]c{0x0b}, len(a)/2))
}
func (v cnd) cstr() (r s) {
	s := "if "
	a := v.argv
	for i := 0; i < len(a)-1; i += 2 {
		if i > 0 {
			s = " else if "
		}
		r += s + "(" + cstring(a[i]) + "){" + cstring(a[i+1]) + "}"
	}
	return r + " else {" + cstring(a[len(a)-1]) + "}"
}
func (v cnd) gstr() (r s) {
	s := "if "
	a := v.argv
	for i := 0; i < len(a)-1; i += 2 {
		if i > 0 {
			s = " else if "
		}
		r += s + gbool(a[i]) + "{" + gstring(a[i+1]) + "}"
	}
	return r + " else {" + gstring(a[len(a)-1]) + "}"
}
func (v v2) rt() T {
	t := v.x().rt()
	if t == 0 { // e.g. uninitialized local (r+:x)
		return v.y().rt()
	}
	return t
}
func (v v2) valid() s {
	if tx, ty := v.x().rt(), v.y().rt(); tx == 0 {
		return sf("left argument has zero type")
	} else if ty == 0 {
		fmt.Fprintf(os.Stderr, "%#v\n", v)
		return sf("right argument has zero type")
	} else if tx != ty {
		fmt.Fprintf(os.Stderr, "%#v\n", v)
		return sf("types mismatch %s %s", tx, ty)
	}
	return ""
}
func (v v2) bytes() []c {
	return append(append(v.x().bytes(), v.y().bytes()...), getop(v2Tab, v.s, v.rt()))
}
func (v v2) cstr() s  { return c2str(v2Tab, v.s, v.rt(), v.x(), v.y()) }
func (v v2) gstr() s  { return g2str(v2Tab, v.s, v.rt(), v.x(), v.y()) }
func (v v1) rt() T    { return v.x().rt() }
func (v v1) valid() s { return ifex(v.x().rt() == 0, "argument has zero type") }
func (v v1) bytes() []c {
	if t := v.rt(); v.s == "-" && t == I {
		return catb([]c{0x41, 0x00}, v.x().bytes(), []c{0x6b}) // 0-x
	} else if v.s == "-" && t == J {
		return catb([]c{0x42, 0x00}, v.x().bytes(), []c{0x7d}) // 0-x
	}
	return append(v.x().bytes(), getop(v1Tab, v.s, v.rt()))
}
func (v v1) cstr() s { o, u := cop(v1Tab, v.s, v.rt()); return jn(o, u, "(", cstring(v.x()), ")") }
func (v v1) gstr() s {
	o, u := gop(v1Tab, v.s, v.rt())
	c := ")"
	if u != "" {
		u = "(" + u
		c += ")"
	}
	return jn(o, u, "(", gstring(v.x()), c)
}
func (v cmp) rt() T    { return I }
func (v cmp) valid() s { return v2(v).valid() }
func (v cmp) bytes() []c {
	return append(append(v.x().bytes(), v.y().bytes()...), getop(cTab, v.s, v.x().rt()))
}
func (v cmp) cstr() s  { return c2str(cTab, v.s, v.x().rt(), v.x(), v.y()) }
func (v cmp) gstr() s  { return "i32b(" + v.bgstr() + ")" }
func (v cmp) bgstr() s { return g2str(cTab, v.s, v.x().rt(), v.x(), v.y()) }
func (v con) rt() T    { return v.t }
func (v con) valid() s { return ifex(v.t == 0, "constant has zero type") }
func (v con) bytes() (r []c) {
	if v.t == I {
		i := v.i
		if i > 2147483647 { // e.g. nai must be positive for gout
			var x int32
			x = int32(i)
			i = int64(x)
		}
		r = append([]c{0x41}, leb(i)...)
	} else {
		r = append([]c{0x41}, leb(v.i)...)
	}
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
		if math.IsNaN(v.f) {
			return "NAN"
		} else if math.IsInf(v.f, 1) {
			return "INFINITY"
		} else {
			s := sf("%v", v.f)
			if strings.Index(s, ".") == -1 && strings.Index(s, "e") == -1 {
				s += ".0"
			}
			return s
		}
	} else if v.t == I && v.i < 0 {
		return sf("%d", uint32(v.i))
	}
	return sf("%d", v.i)
}
func (v con) gstr() s  { return v.cstr() }
func (v cvt) rt() T    { return v.t }
func (v cvt) valid() s { return ifex(v.t == 0, "convert: illegal target type") }
func (v cvt) bytes() []c {
	x := v.x()
	if xt := x.rt(); xt == v.t || xt == C && v.t == I || xt == I && v.t == C {
		return x.bytes()
	}
	tab := map[T]s{
		I: "\x00\x00\xa7\xa7\xab\xaa",
		J: "\xad\xac\x00\x00\xb1\xb0",
		F: "\xb8\xb7\xba\xb9\x00\x00",
	}
	return append(v.x().bytes(), c(tab[v.t][2*tnum[v.x().rt()]+v.sign]))
}
func (v cvt) cstr() s {
	sn := ""
	if v.sign == 1 && v.t == F { // F?'
		sn = "(SI)"
	}
	return jn("(", styp[v.t], ")", sn, cstring(v.x()))
}
func (v cvt) gstr() s {
	sn, cs := "", ""
	if v.sign == 1 && v.t == F { // F?'
		sn, cs = "SI(", ")"
	}
	return jn(styp[v.t], "(", sn, gstring(v.x()), cs, ")")
}                        // todo signed
func (v typ) rt() T      { return 0 }
func (v typ) valid() s   { return "freestanding type" }
func (v typ) bytes() []c { return nil }
func (v fun) rt() T      { return 0 }
func (v fun) valid() s   { return "unapplied func " + v.s }
func (v fun) bytes() []c { return nil }
func (v cal) rt() T      { return v.sig.t }
func (v cal) valid() s {
	if len(v.sig.a) != len(v.argv) {
		return sf("func %s has wrong argn: %d", v.s, len(v.argv))
	}
	for i, a := range v.argv {
		if a.rt() != v.sig.a[i] {
			return sf("func %s arg %d has wrong type", v.s, i+1)
		}
	}
	return ""
}
func (v cal) bytes() (r []c) {
	for _, a := range v.argv {
		r = append(r, a.bytes()...)
	}
	return append(append(r, 0x10), leb(int64(v.n))...)
}
func (v cal) cstr() s {
	av, s := make([]s, len(v.argv)), ""
	for i, a := range v.argv {
		av[i] = cstring(a)
	}
	if v.sig.t == 0 {
		s = ";"
	}
	return jn(v.s, "(", strings.Join(av, ","), ")", s)
}
func (v cal) gstr() s {
	av := make([]s, len(v.argv))
	for i, a := range v.argv {
		av[i] = gstring(a)
	}
	return jn(v.s, "(", strings.Join(av, ","), ")")
}
func (d dot) rt() T    { return d.t }
func (d dot) valid() s { return "" }
func (d dot) bytes() (r []c) {
	for _, a := range d.argv {
		r = append(r, a.bytes()...)
	}
	r = append(r, d.idx.bytes()...)
	r = append(r, 0x11)
	r = append(r, leb(int64(d.sig))...)
	return append(r, 0x00)
}
func (d dot) cstr() s {
	av := make([]s, len(d.argv))
	ptr := fmt.Sprintf("(%s(*)(", d.t)
	for i, a := range d.argv {
		av[i] = cstring(a)
		if i > 0 {
			ptr += ","
		}
		ptr += a.rt().String()
	}
	ptr += "))"
	s := ""
	if d.t == 0 {
		s = ";"
	}
	return jn("(", ptr, "MT[", cstring(d.idx), "])(", strings.Join(av, ","), ")", s)
}
func (d dot) gstr() s {
	av := make([]s, len(d.argv))
	sig := ""
	for i, a := range d.argv {
		av[i] = gstring(a)
		if i > 0 {
			sig += ","
		}
		sig += a.rt().String()
	}
	f := jn("func(", sig, ")", gtyp(d.t))
	return jn("MT[", gstring(d.idx), "].(", f, ")(", strings.Join(av, ","), ")")
}
func (v loc) rt() T      { return v.t }
func (v loc) valid() s   { return ifex(v.t == 0, "local has zero type") }
func (v loc) bytes() []c { return append([]c{0x20}, leb(int64(v.i))...) }
func (v loc) cstr() s    { return locstr(v) }
func (v loc) gstr() s    { return locstr(v) }
func (v las) rt() T      { return 0 }
func (v las) valid() s {
	tx, ty := v.x().rt(), v.y().rt()
	return ifex(tx == 0 || tx != ty, sf("assignment with mismatched types %s %s", tx, ty))
}
func (v las) bytes() []c {
	return append(v.y().bytes(), append([]c{0x21}, leb(int64(v.x().(loc).i))...)...)
}
func (v las) cstr() (r s) { return jn(locstr(v.x()), "=", cstring(v.y()), ";") }
func (v las) gstr() s {
	t := styp[v.x().rt()]
	ys := gstring(v.y())
	if len(ys) > 0 && ys[0] != '(' {
		ys = "(" + ys + ")"
	}
	return jn(locstr(v.x()), "=", t, ys, ";")
}
func (v lod) rt() T {
	if v.t == C {
		return I
	}
	return v.t
}
func (v lod) valid() s { return ifex(v.x().rt() != I, "lod type must be I") }
func (v lod) bytes() (r []c) {
	op := map[T]c{C: 0x2d, I: 0x28, J: 0x29, F: 0x2b}[v.t]
	al := alin[v.t]
	return append(v.x().bytes(), []c{op, al, 0}...)
}
func cgadr(xs string, t T) (r string) { // e.g. "MF[x>>3]"
	r = jn("M", styp[t], "[", xs)
	if t != C {
		r = jn(r, ">>", string('0'+alin[t]))
	}
	return r + "]"
}
func (v lod) cstr() s { return cgadr(cstring(v.x()), v.t) }
func (v lod) gstr() (r s) {
	c := ""
	if v.t == C {
		r, c = "I(", ")"
	}
	return r + cgadr(gstring(v.x()), v.t) + c
}
func (v sto) rt() T { return 0 }
func (v sto) valid() s {
	if v.x().rt() != I {
		return "store addr has wrong type"
	}
	if v.t == 0 {
		return "store has no type"
	}
	return ""
}
func (v sto) bytes() (r []c) {
	y := v.y()
	op := map[T]c{C: 0x3a, I: 0x36, J: 0x37, F: 0x39}[v.t]
	al := alin[v.t]
	return catb(v.x().bytes(), y.bytes(), []c{op, al, 0})
}
func (v sto) cstr() s { // e.g. sI(x, y), not MI[x>>2]=..., because MI could have changed at realloc.
	return jn("s", v.t.String(), "(", cstring(v.x()), ",", cstring(v.y()), ");")
	// return jn(cgadr(cstring(v.x()), v.t), "=(", styp[v.t], ")", cstring(v.y()), ";")
}
func (v sto) gstr() s {
	return jn(cgadr(gstring(v.x()), v.t), "=", styp[v.t], "(", gstring(v.y()), ");")
}
func (v ret) rt() T      { return 0 /*v.x().rt()*/ }
func (v ret) valid() s   { return ifex(v.x().rt() == 0, "return zero type") }
func (v ret) bytes() []c { return append(v.x().bytes(), 0x0f) }
func (v ret) cstr() s    { return jn("R ", cstring(v.x()), ";") }
func (v ret) gstr() s    { return jn("return ", gstring(v.x()), ";") }
func (v iff) rt() T      { return 0 }
func (v iff) valid() s {
	if t := v.x().rt(); t != I && t != C {
		return sf("conditional has wrong type %s", t)
	}
	if t := v.y().rt(); t != 0 {
		return "if statement must not return a value"
	}
	return ""
}
func (v iff) bytes() (r []c) { return catb(v.x().bytes(), []c{0x04, 0x40}, v.y().bytes(), []c{0x0b}) }
func (v iff) cstr() s        { return jn("if(", cstring(v.x()), "){", cstring(v.y()), "}") }
func (v iff) gstr() s        { return jn("if ", gbool(v.x()), "{", gstring(v.y()), "}") }
func (v nlp) rt() T          { return 0 }
func (v nlp) valid() s {
	if xt, yt := v.x().rt(), v.y().rt(); xt != I {
		return sf("loop range is not I: %s", xt)
	} else if yt != 0 {
		return sf("loop body has nonzero type %s", yt)
	}
	return ""
}
func (v nlp) bytes() (r []c) {
	r = v.x().bytes()
	if isexpr(v.x()) {
		r = append(append(r, 0x22), leb(int64(v.n))...) // tee.n for general expressions
	}
	i, n := s(leb(int64(v.c))), s(leb(int64(v.n)))
	//                    if           0   →i   loop
	r = catb(r, []c(sf("\x04\x40\x41\x00\x21%s\x03\x40", i)))
	//                                        i       1   +  tee→i    n   <  continue
	return catb(r, v.y().bytes(), []c(sf("\x20%s\x41\x01\x6a\x22%s\x20%s\x49\x0d\x00\x0b\x0b", i, i, n)))
}
func (v nlp) cstr() (r s) {
	if isexpr(v.x()) {
		r = sf("x%d=%s;", v.n, cstring(v.x()))
	}
	return r + sf("for(x%d=0;x%d<x%d;x%d++){%s}", v.c, v.c, v.n, v.c, cstring(v.y()))
}
func (v nlp) gstr() (r s) {
	if isexpr(v.x()) {
		r = sf("x%d=%s;", v.n, gstring(v.x()))
	}
	return r + sf("for x%d=0;x%d<x%d;x%d++{%s}", v.c, v.c, v.n, v.c, gstring(v.y()))
}
func (v whl) rt() T { return 0 }
func (v whl) valid() s {
	if t := v.x().rt(); t != I {
		return sf("while conditional has wrong type %d", t)
	} else if v.y().rt() != 0 {
		return sf("while body must have no type")
	}
	return ""
}
func (v whl) bytes() (r []c) {
	cnd := sf("%s\x45\x0d\x01", s(v.x().bytes()))
	if _, o := v.x().(con); o {
		cnd = "" // 1/
	}
	//             block   loop     ? y  continue
	return []c(sf("\x02\x40\x03\x40%s%s\x0c\x00\x0b\x0b", cnd, s(v.y().bytes())))
}
func (v whl) cstr() s { return jn("while(", cstring(v.x()), "){", cstring(v.y()), "}") }
func (v whl) gstr() s {
	x := gstring(v.x())
	if x == "1" {
		x = ""
	} else {
		x = gbool(v.x())
	}
	return jn("for ", x, "{", gstring(v.y()), "}")
}
func (v brif) rt() T      { return 0 }
func (v brif) valid() s   { return ifex(v.x().rt() != I, "brif has wrong conditional type") }
func (v brif) bytes() []c { return append(v.x().bytes(), 0x0d, 0x01) } // break outer block
func (v brif) cstr() s    { return jn("if(", cstring(v.x()), ")break;") }
func (v brif) gstr() s    { return jn("if ", gstring(v.x()), "{break}") }
func (v opx) rt() T       { return 0 }
func (v opx) valid() s    { return ifex(s(v) != "!", "nonapplied operator") }
func (v opx) bytes() []c  { return []c{0x00} }
func (v opx) cstr() s     { return "panic();" }
func (v opx) gstr() s     { return `panic("trap")` }
func (v pos) rt() T       { return 0 }
func (v pos) valid() s    { return "position(dummy expr)" }
func (v pos) bytes() []c  { return nil }
func (v nop) rt() T       { return 0 }
func (v nop) valid() s    { return "" }
func (v nop) bytes() []c  { return nil }
func (v nop) cstr() s     { return "()" }
func (v nop) gstr() s     { return "()" }

func locstr(v expr) s { return sf("x%d", v.(loc).i) }
func isexpr(x expr) bool { // general expr that needs an explicit assignment
	switch x.(type) {
	case las:
	case loc:
	default:
		return true
	}
	return false
}
func ifex(c bool, s s) s {
	if c {
		return s
	}
	return ""
}
func gbool(x expr) (r s) {
	if c, o := x.(cmp); o {
		return c.bgstr()
	} else {
		if v, o := x.(v1); o && v.s == "~" {
			if w, o := v.x().(cmp); o && w.s == "~" { // ~x~y
				return gstring(w.x()) + "!=" + gstring(w.y())
			}
			return "0==" + gstring(v.x())
		}
		return "0!=" + gstring(x)
	}
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
		return jn("(", u, embrace(cstring, x), o, u, embrace(cstring, y), ")")
	}
}
func g2str(tab map[s]code, op s, t T, x, y expr) (r s) {
	o, u := gop(tab, op, t)
	uc := ""
	if u != "" {
		u, uc = u+"(", ")"
	}
	if len(o) > 2 {
		return jn(u, o, "(", gstring(x), ",", gstring(y), ")")
	} else {
		return jn("(", u, gstring(x), uc, o, u, gstring(y), uc, ")")
	}
}
func cstring(x expr) s { xs := x.(cstringer); return xs.cstr() }
func gstring(x expr) s { xs := x.(gstringer); return xs.gstr() }
func embrace(sf func(expr) s, x expr) s {
	_, icon := x.(con)
	_, iloc := x.(loc)
	if icon || iloc {
		return sf(x)
	} else {
		return jn("(", sf(x), ")")
	}
}

var v1Tab = map[s]code{
	"-": code{0, 0, 0x9a, "-", "-"},                                            // neg (-I -J is replaced)
	"+": code{0, 0, 0x99, "fabs", "math.Abs"},                                  // abs (+I +J is not allowed)
	"~": code{0x45, 0x50, 0, "!", "n32"},                                       // eqz
	"_": code{1, 1, 0x9c, ";;floor", "math.Floor"},                             // floor (ceil, trunc, nearest?)
	"*": code{0x67, 0x79, 0, "__builtin_clz;__builtin_clzll;", "clz32;clz64;"}, // clz
	"|": code{0x68, 0x79, 0, "__builtin_ctz;__builtin_ctzll;", ""},             // ctz
	"%": code{0, 0, 0x9f, "sqrt", "math.Sqrt"},                                 // sqr
}
var v2Tab = map[s]code{
	`+`:  code{0x6a, 0x7c, 0xa0, "+", "+"},   // add
	`-`:  code{0x6b, 0x7d, 0xa1, "-", "-"},   // sub
	`*`:  code{0x6c, 0x7e, 0xa2, "*", "*"},   // mul
	`%`:  code{0x6e, 0x80, 0xa3, "/", "/"},   // div/div_u
	`%'`: code{0x6d, 0x7f, 0xa3, "S/", "S/"}, // div_s
	`\`:  code{0x70, 0x82, 0, "%", "%"},      // rem_u
	`\'`: code{0x6f, 0x81, 0, "S%", "S%"},    // rem_s
	`&`:  code{0x71, 0x83, 0, "&", "&"},      // and
	`|`:  code{0x72, 0x84, 0, "|", "|"},      // or
	`^`:  code{0x73, 0x85, 0, "^", "^"},      // xor
	`<<`: code{0x74, 0x86, 0, "<<", "<<"},    // shl
	`>>`: code{0x76, 0x88, 0, ">>", ">>"},    // shr_u
	//	`>>'`: code{0x75, 0x87, 0, "S>>", "S>>"},    // shl_s
	//	`<|'`: code{0x77, 0x89, 0, "", ""},          // rotl
	//	`>|'`: code{0x78, 0x8a, 0, "", ""},          // rotr
	//	`&'`:  code{0, 0, 0xa4, "fmin", "math.Max"}, // min
	//	`|'`:  code{0, 0, 0xa5, "fmax", "math.Min"}, // max
}
var cTab = map[s]code{
	"<":   code{0x49, 0x54, 0x63, "<", "<"},     // lt/lt_u
	"<'":  code{0x48, 0x53, 0x63, "S<", "S<"},   // lt_s
	">":   code{0x4b, 0x56, 0x64, ">", ">"},     // gt/gt_u
	">'":  code{0x4a, 0x55, 0x64, "S>", "S>"},   // gt_s
	"<=":  code{0x4d, 0x58, 0x65, "<=", "<="},   // le/le_u
	"<='": code{0x4c, 0x57, 0x65, "S<=", "S<="}, // le_s
	">=":  code{0x4f, 0x5a, 0x66, ">=", ">="},   // ge/ge_u
	">='": code{0x4e, 0x59, 0x66, "S>=", "S>="}, // ge/ge_s
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
func (m module) wasm(tab []segment, data []dataseg) []c {
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
	sec.cat(leb(int64(len(sigv))))
	for _, s := range sigv {
		sec.cat([]c(s))
	}
	sec.out(o)
	// import section(2)
	imports := m.imports()
	if len(imports) > 0 {
		sec = NewSection(2)
		sec.cat(leb(int64(len(imports))))
		for _, f := range imports {
			mod := "ext"
			sec.cat(leb(int64(len(mod))))
			sec.cat([]c(mod))
			sec.cat(leb(int64(len(f.name))))
			sec.cat([]c(f.name))
			sec.cat1(0) // kind
			sec.cat(leb(int64(sigs[s(f.sig())])))
		}
		sec.out(o)
	}
	// function section(3: function signature indexes)
	sec = NewSection(3)
	sec.cat(leb(int64(len(m) - len(imports))))
	for _, f := range m {
		if f.ast != nil {
			sec.cat(leb(int64(sigs[s(f.sig())])))
		}
	}
	sec.out(o)
	// function table section(4)
	if len(tab) > 0 {
		sec = NewSection(4)
		sec.cat1(1)                           // one table
		sec.cat1(0x70)                        // table type
		sec.cat1(0)                           // flags
		sec.cat(leb(int64(segmentsize(tab)))) // size
		sec.out(o)
	}
	// linear memory section(5)
	sec = NewSection(5)
	sec.cat([]c{1, 0, 1}) // 1 initial memory segment, unshared, size 1 block
	sec.out(o)
	// no global section(6)
	// export section(7)
	sec = NewSection(7)
	idx, exp := m.exports()
	sec.cat(leb(int64(1 + len(exp)))) // number of exports (funcs + memory)
	sec.cat(leb(int64(len("mem"))))
	sec.cat([]c("mem"))
	sec.cat1(2) // export kind memory
	sec.cat1(0) // memory index
	for i, f := range exp {
		sec.cat(leb(int64(len(f.name))))
		sec.cat([]c(f.name))
		sec.cat1(0) // function-export
		sec.cat(leb(int64(idx[i])))
	}
	sec.out(o)
	// no start section(8)
	// element section(9)
	if len(tab) > 0 {
		names := make(map[string]int)
		for i, f := range m {
			names[f.name] = i
		}
		sec = NewSection(9)
		sec.cat(leb(int64(len(tab))))
		for _, t := range tab {
			sec.cat1(0) // table index
			sec.cat1(0x41)
			sec.cat(leb(int64(t.off)))
			sec.cat1(0x0b)
			sec.cat(leb(int64(len(t.names))))
			for _, name := range t.names {
				sec.cat(leb(int64(names[name])))
			}
		}
		sec.out(o)
	}
	// code section(10)
	sec = NewSection(10)
	sec.cat(leb(int64(len(m) - len(imports)))) // number of functions
	for _, f := range m {
		if f.ast == nil {
			continue // import
		}
		b := f.code()
		sec.cat(leb(int64(len(b))))
		sec.cat(b)
	}
	sec.out(o)
	// data section(11)
	if len(data) > 0 {
		sec = NewSection(11)
		sec.cat(leb(int64(len(data))))
		for _, d := range data {
			sec.cat1(0)    // memory index
			sec.cat1(0x41) // const.i32 (off is an expr)
			sec.cat(leb(int64(d.off)))
			sec.cat1(0x0b)
			sec.cat(leb(int64(len(d.bytes))))
			sec.cat(d.bytes)
		}
		sec.out(o)
	}
	return o.Bytes()
}
func (m module) imports() (r []fn) {
	for i, f := range m {
		if m[i].ast == nil {
			r = append(r, f)
		}
	}
	return r
}
func (m module) exports() (idx []int, fns []fn) {
	for i, f := range m {
		if f.ex {
			idx = append(idx, i)
			fns = append(fns, f)
		}
	}
	return idx, fns
}
func segmentsize(t []segment) (n int) {
	for _, s := range t {
		if v := s.off + len(s.names); v > n {
			n = v
		}
	}
	return n
}

type section struct {
	t c
	b []c
}
type segment struct {
	off   int
	names []string
}
type dataseg struct {
	off   int
	bytes []c
}

func NewSection(t c) section { return section{t: t} }
func (s *section) cat(b []c) { s.b = append(s.b, b...) }
func (s *section) cat1(b c)  { s.b = append(s.b, b) }
func (s *section) out(w *bytes.Buffer) {
	w.WriteByte(s.t)
	w.Write(leb(int64(len(s.b))))
	w.Write(s.b)
}

func (f fn) sig() (r []c) {
	r = append(r, 0x60)
	r = append(r, leb(int64(f.args))...)
	for i := 0; i < f.args; i++ {
		r = append(r, c(f.locl[i]))
	}
	if f.t == 0 {
		return append(r, 0)
	}
	r = append(r, 1)
	r = append(r, c(f.t))
	return r
}
func (f fn) code() (r []c) {
	if f.ast == nil {
		println("nil ast")
	}
	r = append(r, f.locs()...)
	r = append(r, f.ast.bytes()...)
	return append(r, 0x0b)
}
func (f fn) locs() (r []c) {
	var u []T
	var n []int
	l := f.locl[f.args:]
	for i, t := range l {
		if i > 0 && t == l[i-1] {
			n[len(n)-1]++
		} else {
			u, n = append(u, t), append(n, 1)
		}
	}
	r = leb(int64(len(u)))
	for i, t := range u {
		r = append(r, leb(int64(n[i]))...)
		r = append(r, c(t))
	}
	return r
}

//func leb(v int64) []c { return lebx(v) }
/*
func lebu(v int) []c { // encode unsigned leb128
	if v < 0 {
		panic("lebu")
	}
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
*/
func leb(v int64) []c { // encode signed leb128
	var b []c
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

func (m module) cout(tab []segment, data []dataseg) []c {
	var b bytes.Buffer
	if len(tab) > 0 {
		fmt.Fprintf(&b, "V *MT[%d];\n", segmentsize(tab))
	}
	for _, f := range m { // predeclare functions
		if f.ast == nil {
			continue // import
		}
		st := styp[f.t]
		sig := ""
		for i := 0; i < f.args; i++ {
			if i > 0 {
				sig += ","
			}
			sig += styp[f.locl[i]]
		}
		fmt.Fprintf(&b, "%s %s(%s);", st, f.name, sig)
	}
	fmt.Fprintf(&b, "\n")
	for _, f := range m {
		if f.ast == nil {
			continue // import
		}
		sig, loc := "", ""
		for i := 0; i < f.args; i++ {
			if i > 0 {
				sig += ","
			}
			sig += styp[f.locl[i]] + " " + "x" + s('0'+byte(i))
		}
		for i := f.args; i < len(f.locl); i++ { // declare locals
			t := f.locl[i]
			if t == F {
				loc += sf("%s x%d=.0;", styp[t], i)
			} else {
				loc += sf("%s x%d=0;", styp[t], i)
			}
		}
		st := styp[f.t]
		fmt.Fprintf(&b, "%s %s(%s){%s", st, f.name, sig, loc)
		if sq, o := f.ast.(seq); o {
			for i, e := range sq.argv {
				nl := ""
				if i == len(sq.argv)-1 {
					if f.t != 0 {
						b.WriteString("R ")
					}
					nl = ";}\n"
				}
				b.WriteString(cstring(e))
				b.WriteString(nl)
			}
		} else {
			if f.t != 0 {
				b.WriteString("R ")
			}
			b.WriteString(cstring(f.ast))
			b.WriteString(";}\n")
		}
	}
	if len(tab) > 0 || len(data) > 0 {
		fmt.Fprintf(&b, "V mt_init(){")
		for _, t := range tab {
			for k, s := range t.names {
				fmt.Fprintf(&b, "MT[%d]=%s;", t.off+k, s)
			}
		}
		for i, d := range data {
			v := make([]string, len(d.bytes))
			for j, c := range d.bytes {
				v[j] = strconv.Itoa(int(c))
			}
			fmt.Fprintf(&b, "C d%d[%d]={%s};", i, len(d.bytes), strings.Join(v, ","))
			fmt.Fprintf(&b, "C *mc=MC+%d;for(I i=0;i<%d;i++)mc[i]=d%d[i];\n", d.off, len(d.bytes), i)
		}
		fmt.Fprintf(&b, "}\n")
	}
	return b.Bytes()
}
func (m module) gout(tab []segment, data []dataseg) []c {
	// todo: segments / function pointers
	var b bytes.Buffer
	for _, f := range m {
		if f.ast == nil {
			continue // import
		}
		sig := ""
		for i := 0; i < f.args; i++ {
			if i > 0 {
				sig += ","
			}
			sig += "x" + s('0'+byte(i)) + " " + styp[f.locl[i]]
		}
		r := gtyp(f.t)
		if r != "" {
			r = "(r " + r + ")"
		}
		fmt.Fprintf(&b, "func %s(%s) %s {\n", f.name, sig, r)
		fmt.Fprintln(&b, `//defer func(){fmt.Printf("`+f.name+`: r=%x\n", r)}()`)
		locs := make(map[T][]s)
		lvar, drop := make([]s, len(f.locl)-f.args), make([]s, len(f.locl)-f.args)
		for i := f.args; i < len(f.locl); i++ { // declare locals
			t := f.locl[i]
			s := "x" + strconv.Itoa(i)
			locs[t] = append(locs[t], s)
			lvar[i-f.args] = s
			drop[i-f.args] = "_"
		}
		var v []string
		for i := range locs {
			v = append(v, jn("var ", strings.Join(locs[i], ", "), " ", styp[i], "\n"))
		}
		sort.Strings(v) // reproducible order
		fmt.Fprintf(&b, strings.Join(v, "\n"))
		if len(lvar) > 0 { // _,_,_=x1,x2,x3 (prevent "declared and not used")
			b.WriteString(jn(strings.Join(drop, ","), "=", strings.Join(lvar, ","), "\n"))
		}
		if sq, o := f.ast.(seq); o {
			for i, e := range sq.argv {
				nl := "\n"
				if i == len(sq.argv)-1 {
					if f.t != 0 {
						b.WriteString("return ")
					}
					nl = "}\n"
				}
				b.WriteString(gstring(e))
				b.WriteString(nl)
			}
		} else {
			if f.t != 0 {
				b.WriteString("return ")
			}
			b.WriteString(gstring(f.ast))
			b.WriteString(";}\n")
		}
	}
	if len(tab) > 0 || len(data) > 0 {
		mx := 0
		for _, t := range tab {
			if n := t.off + len(t.names); n > mx {
				mx = n
			}
		}
		fmt.Fprintf(&b, "var MT [%d]interface{}\n", mx)
		fmt.Fprintf(&b, "func mt_init(){\n")
		for _, t := range tab {
			for k, s := range t.names {
				fmt.Fprintf(&b, "MT[%d]=%s;", t.off+k, s)
			}
		}
		for _, d := range data {
			fmt.Fprintf(&b, "copy(MC[%d:], %q)\n", d.off, string(d.bytes))
		}
		fmt.Fprintf(&b, "}\n")
	}
	return b.Bytes()
}
func gtyp(t T) s {
	if t == 0 {
		return ""
	} else {
		return styp[t]
	}
}
