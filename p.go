package i

import (
	"math"
	"math/cmplx"
	"strconv"
)

/*         k4 (-5!"expr")    k7 (`p@"expr")    ngn (`p@"expr")     i
	 →                                     ::                  nil
1        → e                 e                 ,1                  1                    / atom → itself
(1)      → 1                 1                 ,1                  1
,1       → e                 (,:;1)            ,(,:;1))            l{",",1}             / list → function
1;2      → (";";1;2)         (/;1;2)           (1;2)               l{";",1,2}           / special func ";": LR sequence
(1;2)    → (,;1;2)           (\;1;2)           ,(X:;1;2)           l{nil,l{nil,1,2}}    / l[0] nil → list l[1:]
(1;;3)   → (,;1;::;3)        (\;1;::;2)        ,(X:;1;::;3)        l{nil,1,nil,3}
+1       → e                 (+:;1)            ,(+:;1)             l{"+",nil,1}         / monad: first arg is nil
1+       → (+;1)             (+;1)             ,(+;1)              l{"+",1}             / currying: seconds arg is missing
1+2      → (+;1;2)           (+;1;2)           ,(+;1;2)            l{"+",1,2}           / dyad
+[1;2]   → (+;1;2)           ?                 ,(+;1;2)            l{"+",1,2}
1+a      → (+;1;`a)          (+;1;`a)          ,(+;1;`a)           l{"+",1,"a"}         / symbol evaluates (lookup)
1+`a     → (+;1;,`a)         (+;1;,`a)         ,(+;1;,`a)          l{"+",1,l{"`","a"}}}
1+"a"    → (+;1;,"a")        (+;1;,"a")        ,(+;1;,"a")         l{"+",1,'a'}         / rune character, multiple: rv
1+(1;2)  → (+;1;(,;1;2))     (+;1;(\;1;2))     ,(+;1;(X:;1;2))     l{"+",1,l{nil,1,2}}  / nested function
1-2+3    → (-;1;(+2;3))      (-;1;(+;2;3))     ,(-;1;(+;2;3))      l{"-",l{"+",2,3}}
+/1 2 3  → e                 ((/;+);1 2 3)     ,((/;+);1 2 3)      l{l{"/","+"}},fv{1,2,3}          / derived verb
1{x+y}2  → e                 e                 e                   l{"λ",l{l{'+',"x","y"}},"x","y"} / lambda function
{x}[3]   → ({x};3)           ({x};3)           ,({x};3)            l{"λ",l{"x"}}
a[3]:4   → (:;(`a;3);4)      (::;(`a;3);4)     (:;(`a;3);4;::)     l{":",l{"a",3},4}    / assignment
a[3;4]+:5→ (+:;(`a;3;4);5)   (+:;(`a;3;4);5)   (+:;(`a;3;4);5;);;) l{"+:"},l{"a",3,4},5 / modified assignment
*/

func prs(s v) v { // s: rv (no comments)
	p := p{b: s.(rv)}
	r := l{}
	for p.a() { // ex;ex...
		ex := p.ex(p.noun())
		if ex == nil {
			break
		}
		r = append(r, ex)
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

type rn = rune
type rv = []rn
type sf func(rv) i

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
func (p *p) ex(a v) v {
	// TODO
	if a == nil {
		return nil
	}
	// TODO
	return a
}
func (p *p) noun() v {
	switch {
	// TODO colon
	// TODO ioverb
	case p.t(sNum):
		r := zv{}
		c := false
		for p.t(sNum) {
			x, ic := p.num(p.p(sNum))
			r, c = append(r, x), c || ic
		}
		switch {
		case !c && len(r) == 1:
			return real(r[0])
		case !c:
			fv := make(fv, len(r))
			for i := range r {
				fv[i] = real(r[i])
			}
			return fv
		case c && len(r) == 1:
			return r[0]
		}
		return r
	case p.t(sSym):
		r := sv{}
		for p.t(sSym) {
			r = append(r, p.sym(p.p(sSym)))
		}
		if len(r) == 1 {
			return l{"`", r[0]}
		}
		return r
	case p.t(sStr):
		r := p.str(p.p(sStr))
		if len(r) == 1 {
			return r[0]
		}
		return r
	case p.t(sObr):
		var key, val l
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
		return l{"!", key, val}
	// TODO {}
	// TODO ()
	// TODO verb
	case p.t(sNam):
		ref := p.p(sNam)
		println("ref", ref)
		if p.t(sCol) {
			p.p(sCol)
			return l{":", ref, p.ex(p.noun())}
		}
		return ref
		// TODO compound assign []
	}
	return nil
}
func (p *p) num(s s) (z, bool) {
	pf := func(s string) f {
		f, o := strconv.ParseFloat(s, 64)
		if o != nil {
			e("num")
		}
		return f
	}
	for i, r := range s {
		if r == 'a' {
			rp := pf(s[:i])
			switch j := s[i+1:]; j {
			case "0":
				return complex(rp, 0), true
			case "90":
				return complex(0, rp), true
			case "180":
				return complex(-rp, 0), true
			case "270":
				return complex(0, -rp), true
			default:
				return cmplx.Rect(rp, math.Pi*pf(j)/180.0), true
			}
		} else if r == 'i' {
			return complex(pf(s[:i]), pf(s[i+1:])), true
		}
	}
	return complex(pf(s), 0), false
}
func (p *p) sym(s s) s { // `a | `"a"
	if len(s) < 4 || s[1] != '"' {
		return s[1:]
	}
	return s[2 : len(s)-1]
}
func (p *p) str(s s) rv {
	s, o := strconv.Unquote(s)
	if o != nil {
		e("str")
	}
	return rv(s)
}

// scanner
const dig = "0123456789"
const con = "π"
const sym = `+\-*%!&|<>=~,^#_$?@.`
const uni = `⍉×÷⍳⍸⌊⌽⌈⍋⌸≡∧⍴≢↑⌊↓⍕∪⍎⍣¯ℜℑ√⍟`
const uav = "⍨¨⌿⍀"
const wsp = " \t\r"

// Scanners return the rune count of the matched input or 0; input len > 1.
func sNum(s rv) i { // number f | fjf, allow leading +
	n := 0
l:
	for i, r := range s {
		switch {
		case i == 0 && any(r, "-+"):
		case any(r, dig):
		case i > 0 && any(r, "eEai"): // 1i0, 0i1 not 1i
			if len(s) < i+1 || i == 1 && any(s[0], "-+") {
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
func sNam(s rv) i { // name [a-Z][_a-Z0-9]*
	a := func(r rn) bool {
		if alpha(r) {
			return true
		}
		return false
	}
	n := 0
	for i, r := range s {
		switch {
		case i == 0 && any(r, con):
			return 1
		case i == 0 && !a(r):
			return 0
		case a(r) || (i > 0 && any(r, "_0123456789")):
		default:
			return i
		}
		n++
	}
	return n
}
func sSym(s rv) i { // symbol `name|`string
	if s[0] != '`' {
		return 0
	}
	if len(s) == 1 {
		return 1
	} else if n := sStr(s[1:]); n > 0 {
		return 1 + n
	}
	return 1 + sNam(s[1:])
}
func sStr(s rv) i { // string "str\esc"
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
func sVrb(s rv) i { // verb single rune ascii or unicode
	for _, r := range s {
		if any(r, sym) || any(r, uni) {
			return 1
		}
		return 0
	}
	return 0
}
func sAsn(s rv) i { // assignment verb:
	if n := sVrb(s); n != 0 && len(s) > n && s[n] == ':' {
		return n + 1
	}
	return 0
}
func sIov(s rv) i { // io verb [0-9]:
	if len(s) < 2 {
		return 0
	}
	if any(s[0], dig) && s[1] == ':' {
		return 2
	}
	return 0
}
func sAdv(s rv) i { // adverb ascii or unicode
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
func sSem(s rv) i { // ;|\n
	if any(s[0], ";\n") {
		return 1
	}
	return 0
}
func sWsp(s rv) i { // whitespace
	for i, r := range s {
		if !any(r, wsp) {
			return i
		}
	}
	return len(s)
}
func sCol(s rv) i { return pref(s, ":") }  // :
func sViw(s rv) i { return pref(s, "::") } // ::
func sDct(s rv) i { // dict [name:
	if len(s) < 2 || s[0] != '[' {
		return 0
	}
	if n := sNam(s[1:]); n > 0 && len(s) > n+1 && s[n+1] == ':' {
		return n + 2
	}
	return 0
}
func sObr(s rv) i { return pref(s, "[") } // [
func sOpa(s rv) i { return pref(s, "(") } // (
func sOcb(s rv) i { return pref(s, "{") } // ]
func sCbr(s rv) i { return pref(s, "]") } // ]
func sCpa(s rv) i { return pref(s, ")") } // )
func sCcb(s rv) i { return pref(s, "}") } // }

func any(r rn, s s) bool {
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
func alpha(r rn) bool {
	if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
		return true
	}
	return false
}
