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
)

type c = byte
type T c
type fn struct { // name:I:IIF::body..
	name string
	src  [2]int // line, col
	t    T      // return type
	args []T
	locl []T
	sign int
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
var P = I // I(wasm32) J(wasm64)

type module []fn

func main() {
	var html bool
	flag.BoolVar(&html, "html", false, "html output")
	flag.Parse()
	b := run(os.Stdin)
	if html {
		b = page(b)
	}
	os.Stdout.Write(b)
}
func run(r io.Reader) []c {
	rd := bufio.NewReader(r)
	state := sNewl
	line, char, hi := 1, 0, true
	err := func(s string) { fatal(fmt.Errorf("%d:%d: %s", line, char, s)) }
	var m module
	var f fn
	var p c
	for {
		b, e := rd.ReadByte()
		if e == io.EOF {
			if f.name != "" {
				m = append(m, f)
			}
			return m.emit()
		}
		char++
		if b == '\n' {
			line++
			char = 1
		}
		fatal(e)
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
					panic("parse return type")
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
			if b == '}' {
				state = sCmnt
				f.parse()
			} else {
				f.WriteByte(b)
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
func fatal(e error) {
	if e != nil {
		panic(e)
	}
}

type parser struct {
	name       string
	line, char int
	b          []byte
	tok        []byte
}

func (f fn) parse() expr { // parse function body
	p := parser{name: f.name, line: f.src[0], char: f.src[1], b: strip(f.Bytes())}
	e := p.seq('}')
	// todo build locals, validate, bytes
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
func (p *parser) err(s string) { panic(fmt.Errorf("%d:%d(%s) %s", p.line, p.char, p.name, s)) }
func (p *parser) w() {
	for len(p.b) > 0 {
		if c := p.b[0]; c == ' ' || c == '\t' || c == '\n' {
			p.char++
			if c == '\n' {
				p.line++
				p.char = 0
			}
			p.b = p.b[1:]
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
		p.char += n
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
func (p *parser) ex(v expr) expr {
	panic("nyi")
	return nil
}
func (p *parser) noun() expr {
	p.w()
	if len(p.b) == 0 {
		return nil
	}
	switch {
	case p.t(sArg):
		return pArg(p.tok)
	case p.t(sSym):
		return pSym(p.tok)
	case p.t(sCon):
		return pCon(p.tok)
	case p.t(sOp):
		return pOp(p.tok)
	case p.t(sC('(')):
		return p.seq(')')
	default:
		return nil
	}
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
func pArg(b []c) expr {
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
func pSym(b []c) expr { return loc{s: string(b)} }
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
			return neg + i
		}
	}
	return neg + len(b)
}
func pCon(b []c) expr {
	var r con
	if bytes.IndexByte(b, '.') != -1 {
		if f, err := strconv.ParseFloat(string(b), 64); err != nil {
			panic(err)
		} else {
			r.t = F
			r.f = f
		}
	}
	r.t = I
	if c := b[len(b)-1]; c == 'i' || c == 'j' {
		b = b[:len(b)-1]
		if c == 'j' {
			r.t = J
		}
	}
	if i, err := strconv.ParseInt(string(b), 10, 64); err != nil {
		panic(err)
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
func pOp(b []c) expr         { return opx(string(b)) }
func sC(x c) func(b []c) int { return func(b []c) int { return boolvar(b[0] == x) } }

// intermediate representation for function bodies (typed expression tree)
type expr interface {
	rt() T // result type, maybe 0
	args() []expr
	valid() bool
	bytes() []c
}
type seq []expr  // a;b;..
type v2 struct { // x+y unitype
	s    string // +-*%
	l, r expr
}
type v1 struct { // -y
	s string
	a expr
}
type cmp struct { // x<y..
	s    string
	l, r expr
}
type arg struct { // x y ..
	t T
	n int
}
type loc struct { // abc
	arg
	s string
}
type con struct { // numeric constant
	t T
	i int64
	f float64
}
type opx string // operator

func getop(tab map[string][3]c, op string, t T) (r c) {
	ops, ok := tab[op]
	if !ok {
		panic("type")
	}
	i, ok := tnum[t]
	r = ops[i]
	if !ok || r == 0 {
		panic("type")
	}
	return r
}

func (s seq) rt() T        { return s[len(s)-1].rt() }
func (s seq) args() []expr { return s }
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
func (v v2) rt() T         { return v.r.rt() }
func (v v2) args() []expr  { return []expr{v.l, v.r} }
func (v v2) valid() bool   { return v.r.rt() == v.l.rt() && v.r.rt() != 0 }
func (v v2) bytes() []c    { return append(append(v.l.bytes(), v.r.bytes()...), getop(v2Tab, v.s, v.rt())) }
func (v v1) rt() T         { return v.a.rt() }
func (v v1) args() []expr  { return []expr{v.a} }
func (v v1) valid() bool   { return v.a.rt() != 0 }
func (v v1) bytes() []c    { return append(v.a.bytes(), getop(v1Tab, v.s, v.rt())) }
func (v cmp) rt() T        { return I }
func (v cmp) args() []expr { return []expr{v.l, v.r} }
func (v cmp) valid() bool  { return v.r.rt() == v.l.rt() && v.r.rt() != 0 }
func (v cmp) bytes() []c   { return append(append(v.l.bytes(), v.r.bytes()...), getop(cTab, v.s, v.rt())) }
func (v arg) rt() T        { return v.t }
func (v arg) args() []expr { return nil }
func (v arg) valid() bool  { return v.n >= 0 && v.t != 0 }
func (v arg) bytes() []c   { return append([]c{0x20}, lebu(int(v.n))...) }
func (v con) rt() T        { return v.t }
func (v con) valid() bool  { return v.t != 0 }
func (v con) args() []expr { return nil }
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
func (v opx) rt() T        { return 0 }
func (v opx) args() []expr { return nil }
func (v opx) valid() bool  { return false }
func (v opx) bytes() []c   { return nil }

var v1Tab = map[string][3]c{
	"-": [3]c{0, 0, 0x9a},    // neg (no neg for ints)
	"+": [3]c{0, 0, 0x99},    // abs
	"_": [3]c{1, 1, 0x9c},    // floor (ceil, trunc, nearest?)
	"*": [3]c{0x67, 0x79, 0}, // clz
	"|": [3]c{0x68, 0x79, 0}, // ctz
	"%": [3]c{0, 0, 0x9f},    // sqr
	// TODO min/max
}
var v2Tab = map[string][3]c{
	`+`:   [3]c{0x6a, 0x7c, 0xa0}, // add
	`-`:   [3]c{0x6b, 0x7d, 0xa1}, // sub
	`*`:   [3]c{0x6c, 0x7e, 0x94}, // mul
	`%`:   [3]c{0x6e, 0x80, 0xa3}, // div/div_u
	`%'`:  [3]c{0x6d, 0x7f, 0xa3}, // div_s
	`\`:   [3]c{0x70, 0x82, 0},    // rem_u
	`\'`:  [3]c{0x6f, 0x81, 0},    // rem_s
	`&`:   [3]c{0x71, 0x83, 0},    // and
	`|`:   [3]c{0x72, 0x84, 0},    // or
	`^`:   [3]c{0x73, 0x85, 0},    // xor
	`<<`:  [3]c{0x74, 0x86, 0},    // shl
	`>>`:  [3]c{0x76, 0x88, 0},    // shr_u
	`>>'`: [3]c{0x75, 0x87, 0},    // shl_s
	`<|'`: [3]c{0x77, 0x89, 0},    // rotl
	`>|'`: [3]c{0x78, 0x8a, 0},    // rotr
}
var cTab = map[string][3]c{
	"<":   [3]c{0x49, 0x54, 0x63}, // lt/lt_u
	"<'":  [3]c{0x48, 0x53, 0x63}, // lt_s
	">":   [3]c{0x4b, 0x56, 0x64}, // gt/gt_u
	">'":  [3]c{0x4a, 0x55, 0x64}, // gt_s
	"<=":  [3]c{0x4d, 0x58, 0x65}, // le/le_u
	"<='": [3]c{0x4c, 0x57, 0x65}, // le_s
	">=":  [3]c{0x4f, 0x5a, 0x66}, // ge/ge_u
	">='": [3]c{0x4e, 0x59, 0x66}, // ge/ge_u
	"~":   [3]c{0x46, 0x51, 0x61}, // eq
	"!":   [3]c{0x47, 0x52, 0x62}, // ne
}
var allops map[string]bool

func init() {
	allops = make(map[string]bool)
	for _, t := range []map[string][3]c{v1Tab, v2Tab, cTab} {
		for s := range t {
			allops[s] = true
		}
	}
}

// emit byte code
func (m module) emit() []c {
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

func log(a ...interface{})            { fmt.Fprintln(os.Stderr, a...) }
func logf(f string, a ...interface{}) { fmt.Fprintf(os.Stderr, f, a...) }
func page(wasm []c) []c {
	var b bytes.Buffer
	b.Write([]c(head))
	b.WriteString(base64.StdEncoding.EncodeToString(wasm))
	b.Write([]c(tail))
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
