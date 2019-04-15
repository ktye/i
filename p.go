package i

import (
	"math"
	"math/cmplx"
	"strconv"
)

func prs(x v) v { // s: rv
	p := p{b: beg(rv(x.(s)))}
	r := l{}
	for p.a() { // ex;ex...
		ex := p.ex(p.noun())
		if ex == nil {
			break
		}
		r = append(r, ex)
		if p.t(sSem) {
			p.p(sSem)
		}
	}
	if p.a() {
		e("prs:" + string(p.b))
	}
	if len(r) == 0 {
		return nil
	} else if len(r) == 1 {
		return r[0]
	}
	return append(l{";"}, r...)
}

type rv = []rune
type sf func(rv) int

type p struct {
	b rv
	d int
}

func (p *p) a() bool     { return len(p.w().b) > 0 }               // buffer available (not empty)
func (p *p) t(f sf) bool { return len(p.w().b) > 0 && f(p.b) > 0 } // test and keep token
func (p *p) p(f sf) s { // (must)parse and remove token return capture
	n := f(p.w().b)
	if n < 0 {
		e("parse")
		return ""
	}
	s := s(p.b[:n])
	p.b = p.b[n:]
	return s
}
func (p *p) m(f sf) bool { // match and remove token
	p.w()
	if len(p.b) == 0 {
		return false
	}
	n := f(p.b)
	p.b = p.b[n:]
	return n > 0
}
func (p *p) w() *p { // remove wsp
	i := 0
	for _, r := range p.b {
		if !any(r, wsp) {
			break
		}
		i++
	}
	p.b = p.b[i:]
	return p
}
func (p *p) isVrb(x v) bool {
	if iss(x) && sVrb(rv(x.(s))) > 0 {
		return true // primitive
	}
	lv, o := x.(l)
	if o && len(lv) > 1 && iss(lv[0]) {
		s0 := rv(lv[0].(s))
		if sVrb(s0) > 0 && len(lv) == 3 && lv[2] == nil { // curry: e.g. (2+)
			return true
			// If λs are treated as verbs, they would be allowed infix.
			// But this would breaks {lambda}@arg and {lambda}.(args)
			// } else if len(lv) == 2 && (lv[0].(s) == "λ" || sAdv(s0) > 0) {
		} else if len(lv) == 2 && sAdv(s0) > 0 {
			return true
		}
	}
	return false
}
func (p *p) ex(a v) v {
	atNoun := func() bool {
		return p.t(sNum) || p.t(sNam) || p.t(sStr) || p.t(sOpa) || p.t(sOcb)
	}
	if a == nil {
		return nil
	}
	if p.t(sAdv) {
		return p.adv(nil, a)
	}
	if p.isVrb(a) {
		// TODO at (
		x := p.noun()
		// TODO x is verb
		if r := p.ex(x); r != nil {
			a = l{a, r} // monad
		}
	}
	if p.t(sVrb) || p.t(sIov) {
		x := p.noun()
		// TODO force monad
		if p.t(sAdv) {
			return p.adv(a, x)
		}
		if r := p.ex(p.noun()); r != nil {

			return l{x, a, r}
		}
		return l{x, a, nil} // curry
	}
	if atNoun() {
		x := p.noun()
		if p.isVrb(x) {
			y := p.ex(p.noun())
			if y == nil {
				return l{x, a, nil} // curry
			}
			return l{x, a, y}
		}
		if p.t(sAdv) {
			return p.adv(a, x)
		}
		return l{a, p.ex(x)}
	}
	return a
}
func (p *p) noun() v {
	switch {
	// TODO colon
	// TODO ioverb
	case p.t(sNum):
		r := zv{}
		for p.t(sNum) {
			r = append(r, p.num(p.p(sNum)))
		}
		if len(r) == 1 {
			return p.idxr(r[0])
		}
		return p.idxr(r)
	case p.t(sStr):
		r := sv{}
		for p.t(sStr) {
			r = append(r, p.str(p.p(sStr)))
		}
		if len(r) == 1 {
			return l{"`", r[0]}
		}
		return p.idxr(r)
	case p.t(sObr):
		var key sv
		var val = l{nil}
		p.p(sObr)
		for {
			key = append(key, p.p(sNam))
			p.p(sCol)
			val = append(val, p.ex(p.noun()))
			if !p.t(sSem) {
				break
			}
			p.p(sSem)
		}
		p.p(sCbr)
		return p.idxr(l{"!", key, val})
	case p.t(sOcb):
		p.p(sOcb)
		r := p.lst(sCcb)
		λ := l{"λ", r}
		if p.t(sObr) {
			x := p.idxr(nil)
			x.(l)[0] = λ
			return x
		}
		return λ
	case p.t(sOpa):
		p.p(sOpa)
		r := p.lst(sCpa)
		if len(r) == 1 {
			if p.isVrb(r[0]) { // e.g. (2+)
				return p.idxr(p.drv(r[0]))
			}
			return p.idxr(r[0])
		} else if p.isVrb(r) {
			return r // curry
		}
		return p.idxr(append(l{nil}, r...))
	case p.t(sVrb):
		h := p.p(sVrb)
		if p.t(sCol) { // modified assignment
			h += p.p(sCol)
		}
		// TODO [ | dict
		return p.idxr(h)
	case p.t(sNam):
		ref := p.p(sNam)
		if p.t(sCol) {
			p.p(sCol)
			col := ":"
			if p.t(sCol) {
				p.p(sCol)
				col += ":" // :: global assignment
			}
			return l{col, ref, p.ex(p.noun())}
		}
		return p.idxr(ref)
		// TODO compound assign []
	}
	return nil
}
func (p *p) drv(w v) v {
	for p.t(sAdv) {
		a := p.p(sAdv)
		w = l{a, w}
	}
	return w
}
func (p *p) adv(left, w v) v {
	a := p.p(sAdv)
	for p.t(sAdv) {
		b := p.p(sAdv)
		if iss(w) && sVrb(rv(w.(s))) > 0 {
			// reflite p131: force monad for derived primitives
			w = l{l{a, w}, nil}
		} else {
			w = l{a, w}
		}
		a = b
	}
	// TODO [] callright
	r := p.ex(p.noun())
	if left == nil {
		return l{l{a, w}, r}
	}
	return l{l{a, w}, left, r}
}
func (p *p) num(s s) z {
	pf := func(s string) float64 {
		neg := 1.0
		if s[0] == '+' {
			s = s[1:]
		} else if s[0] == '-' {
			s = s[1:]
			neg = -1.0
		}
		var f float64
		switch rv(s)[0] {
		case 'π':
			f = math.Pi
		case '∞':
			f = math.Inf(1)
		case '𝜀':
			f = 1.0E-14
		case 'ø':
			f = math.NaN()
		default:
			var o error
			f, o = strconv.ParseFloat(s, 64)
			if o != nil {
				e("num")
			}
		}
		return neg * f
	}
	for i, r := range s {
		if r == 'a' {
			rp := pf(s[:i])
			switch j := s[i+1:]; j {
			case "0":
				return complex(rp, 0)
			case "90":
				return complex(0, rp)
			case "180":
				return complex(-rp, 0)
			case "270":
				return complex(0, -rp)
			default:
				return cmplx.Rect(rp, math.Pi*pf(j)/180.0)
			}
		} else if r == 'i' {
			return complex(pf(s[:i]), pf(s[i+1:]))
		}
	}
	return complex(pf(s), 0)
}
func (p *p) str(s s) s { // "string" | `name
	if s[0] == '`' {
		return s[1:]
	}
	s, o := strconv.Unquote(s)
	if o != nil {
		e("str")
	}
	return s
}
func (p *p) lst(term sf) l {
	r := l{}
	if p.t(term) {
		p.p(term)
		return r
	}
	for {
		r = append(r, p.ex(p.noun()))
		if !p.t(sSem) {
			break
		}
		p.p(sSem)
	}
	p.p(term)
	return r
}
func (p *p) idxr(x v) v {
	if p.t(sDct) {
		return x
	}
	for p.t(sObr) {
		p.p(sObr)
		r := p.lst(sCbr)
		if len(r) == 0 {
			r = l{nilad(true)}
		}
		return append(l{x}, r...)
	}
	return x
}

// scanner
const dig = "0123456789"
const sym = `+\-*%!&|<>=~,^#_$?@.`
const uni = `⍉×÷⍳∈⌊⌽⌈⍋⍒≡∧⍴≢↑⌊↓⍕∪⍎⍣ℜℑ‖°𝜑√⍟∇`
const uav = "⍨¨⌿⍀"
const wsp = " \t\r"

// Scanners return the rune count of the matched input or 0; input len > 1.
func sNum(s rv) int { // number f | fjf, allow leading +
	if s[0] == 'ø' {
		return 1
	} else if any(s[0], "π∞𝜀") {
		if len(s) > 2 && s[1] == 'i' {
			if n := sNum(s[2:]); n > 0 {
				return n + 2
			}
		}
		return 1
	} else if len(s) > 1 && any(s[0], "-+") && any(s[1], "π∞𝜀") {
		if n := sNum(s[1:]); n > 0 {
			return n + 1
		}
	}
	n := 0
l:
	for i, r := range s {
		switch {
		case i == 0 && any(r, "-+"):
		case any(r, dig):
		case i > 0 && any(r, "eEai"): // 1i0, 0i1 not 1i
			if len(s) < i+2 || i == 1 && any(s[0], "-+") {
				return 0
			}
			if n := sNum(s[i+1:]); n > 0 {
				return i + n + 1
			}
			return i
		case r == '.' && i > 0 && i < len(s)-1 && any(s[i-1], dig) && any(s[i+1], dig):
		default:
			break l
		}
		n++
	}
	if n == 1 && any(s[0], "+-") {
		return 0
	}
	return n
}
func sNam(s rv) int { // name [a-Z][a-Z0-9]*
	a := func(r rune) bool {
		if alpha(r) {
			return true
		}
		return false
	}
	n := 0
	for i, r := range s {
		switch {
		case i == 0 && !a(r):
			return 0
		case a(r) || (i > 0 && any(r, "0123456789")):
		default:
			return i
		}
		n++
	}
	return n
}
func sStr(s rv) int { // string `name | "str\esc"
	if s[0] == '`' {
		if len(s) == 1 {
			return 1
		}
		return 1 + sNam(s[1:])
	}
	if len(s) < 2 || s[0] != '"' {
		return 0
	}
	q := false
	for i, r := range s {
		switch {
		case i == 0:
		case r == '\\':
			q = !q
		case r == '"' && !q:
			return i + 1
		}
		if q && r != '\\' {
			q = false
		}
	}
	return 0
}
func sVrb(s rv) int { // verb single rune ascii or unicode, possibly with : suffix
	if any(s[0], sym) || any(s[0], uni) {
		if len(s) > 1 && s[1] == ':' {
			return 2
		}
		return 1
	}
	return 0
}
func sAsn(s rv) int { // assignment verb:
	if n := sVrb(s); n != 0 && len(s) > n && s[n] == ':' {
		return n + 1
	}
	return 0
}
func sIov(s rv) int { // io verb [0-9]:
	if len(s) < 2 {
		return 0
	}
	if any(s[0], dig) && s[1] == ':' {
		return 2
	}
	return 0
}
func sAdv(s rv) int { // adverb ascii or unicode
	for i, r := range s {
		if i == 0 && any(r, uav) {
			return 1
		}
	}
	if any(s[0], `'/\`) {
		if len(s) > 1 && s[1] == ':' {
			return 2
		}
		return 1
	}
	return 0
}
func sSem(s rv) int { // ;|\n
	if any(s[0], ";\n") {
		return 1
	}
	return 0
}
func sWsp(s rv) int { // whitespace
	for i, r := range s {
		if !any(r, wsp) {
			return i
		}
	}
	return len(s)
}
func sCol(s rv) int { return pref(s, ":") }  // :
func sViw(s rv) int { return pref(s, "::") } // ::
func sDct(s rv) int { // dict [name:
	if len(s) < 2 || s[0] != '[' {
		return 0
	}
	if n := sNam(s[1:]); n > 0 && len(s) > n+1 && s[n+1] == ':' {
		return n + 2
	}
	return 0
}
func sObr(s rv) int { return pref(s, "[") } // [
func sOpa(s rv) int { return pref(s, "(") } // (
func sOcb(s rv) int { return pref(s, "{") } // ]
func sCbr(s rv) int { return pref(s, "]") } // ]
func sCpa(s rv) int { return pref(s, ")") } // )
func sCcb(s rv) int { return pref(s, "}") } // }

func any(r rune, s s) bool {
	for _, x := range s {
		if r == x {
			return true
		}
	}
	return false
}
func pref(r rv, p string) int {
	s := string(r)
	if len(s) < len(p) {
		return 0
	}
	if s[:len(p)] == p {
		return len([]rune(p))
	}
	return 0
}
func alpha(r rune) bool {
	if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
		return true
	}
	return false
}
func beg(x rv) rv { // preserve strings, strip comments, disambiguate +-
	var y rv
	p := func(n int) {
		y = append(y, x[:n]...)
		x = x[n:]
	}
	r, c := ' ', false
L:
	for {
		switch {
		case len(x) == 0:
			break L
		case !c && sStr(x) > 0: // skip strings
			p(sStr(x))
			r = '.'
			continue
		case !c && sNum(x) > 0:
			if any(x[0], "+-") {
				y = append(y, x[0])
				x = x[1:]
				if !(sVrb(rv{r}) > 0 || any(r, " ([{;")) {
					y = append(y, ' ')
				}
			}
			p(sNum(x))
			r = '0'
			continue
		case x[0] == '/' && r == ' ': // if at newline or after blank: comment
			c = true
			fallthrough
		default:
			switch {
			case c && x[0] == '\n':
				c, r = false, ' '
			case c:
				x = x[1:]
				continue
			case any(x[0], " \t"):
				r = ' '
			default:
				r = x[0]
			}
			p(1)
		}
	}
	return y
}

type nilad bool           // mark argument of a niladic call f[] → ("f"; nilad)
func (n nilad) String() s { return "nilad" }
