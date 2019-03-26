package i

type rn = rune
type rv = []rn
type sf func(rv) i

type p struct {
	b rv
	d int
}

/*        k4 (-5!"expr")    k7 (`p@"expr")    ngn (`p@"expr")     i
1       → e                                   ,1                  1                    / atom → itself
(1)     → 1                                   ,1                  1
,1      → e                                   ,(,:;1))            l{",",1}             / list → function
1;2     → (";";1;2)                           (1;2)               l{";",1,2}           / special func ";": LR sequence
(1;2)   → (,;1;2)                             ,(X:;1;2)           l{nil,l{nil,1,2}}    / l[0] nil → list l[1:]
(1;;3)  → (,;1;::;3)                          ,(X:;1;::;3)        l{nil,1,nil,3}
+1      → e                                   ,(+:;1)             l{"+",nil,1}         / monad: first arg is nil
1+      → (+;1)                               ,(+;1)              l{"+",1}             / currying: seconds arg is missing
1+2     → (+;1;2)                             ,(+;1;2)            l{"+",1,2}           / dyad
+[1;2]  → (+;1;2)                             ,(+;1;2)            l{"+",1,2}
1+a     → (+;1;`a)                            ,(+;1;`a)           l{"+",1,"a"}         / symbol evaluates (lookup)
1+`a    → (+;1;,`a)                           ,(+;1;,`a)          l{"+",1,l{"`","a"}}}
1+"a"   → (+;1;,"a")                          ,(+;1;,"a")         l{"+",1,'a'}         / rune character, multiple: rv
1+(1;2) → (+;1;(,;1;2))                       ,(+;1;(X:;1;2))     l{"+",1,l{nil,1,2}}  / nested function
1-2+3   → (-;1;(+2;3))                        ,(-1;(+;2;3))       l{"-",l{"+",2,3}}
+/1 2 3 → e                                   ,((/;+);1 2 3)      l{l{"/","+"}},fv{1,2,3}          / derived verb
1{x+y}2 → e                                   e                   l{"λ",l{l{'+',"x","y"}},"x","y"} / lambda function
{x}[3]  → ({x};3)                             ,({x};3)            l{"λ",l{"x"}}
a[3]:4  → (:;(`a;3);4)                        (:;(`a;3);4;::)     l{":",l{"a",3},4}    / assignment
a[3]+:4 → (+:;(`a;3;4);5)                     (+:;(`a;3;4);5;);;) l{"+:"},l{"a",3,4},5 / modified assignment
*/

func (p *p) a() bool     { return len(p.w().b) > 0 } // buffer available (not empty)
func (p *p) t(f sf) bool { return f(p.w().b) > 0 }   // test and keep token
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

func prs(v interface{}) interface{} { // string contains no comments
	// p := p{b: pbeg(rv(v.(s)))}
	p := p{b: rv(v.(s))}
	if len(p.b) == 0 {
		return l{}
	}
	/*
		r := p.lst(nil, false)
		if len(p.b) > 0 {
			return e("parse:" + string(p.b))
		}
		return r
	*/
	return e("nyi")
}

/*
func (p *p) mul() l { // expr;expr... → (; expr expr ...)
	var m l
	for p.o() {
		if p.m(ws
		       if !p.m(sSem) {
			       break
		       }
		       }
		       for p.o() {
			       e := p.expr(p.noun())
}
*/
/*
func (p *p) lst(tf sf, c bool) (r l) {
	for p.o() {
		if tf != nil && tf(p.b) > 0 {
			break
		}
		for p.m(sSem) {
			if !c {
				r = append(r, nil)
			}
		}
		ex := p.expr(p.noun())
		// find sticky // TODO
		if ex != nil {
			r = append(r, ex)
		} else if !c {
			r = append(r, nil)
		}
		if p.m(sSem) == false {
			break
		}
	}
	if tf != nil {
		p.p(tf)
	}
	return
}
*/
func (p *p) expr(node v) v {
	return e("nyi")
	/*
		if node == nil {
			return nil
		}
		if p.m(sAdv) {
			return p.adv(nil, node)
		}
		if w, o := node.(vrb); o && w.curry == nil {
			pa := p.m(sOpa)
			x := p.noun()
		}
	*/
}
func (p *p) noun() v { return e("nyi") }

/*
	b, n := rv(v.(s)), 0
	for _, r := range b { // trim left
		if !any(r, wsp) {
			break
		}
		n++
	}
	b = pbeg(b[n:])
	if len(b) == 0 {
		return l{}
	}
	for i := len(c) - 1; i >= 0; i-- { // trim right
		if !any(c[i], wsp) {
			break
		}
		c = c[:i]
	}
	if len(c) == 0 {
		return l{}
	}
	var r l
	c, r = pLst(c, nil, false)
	if len(c) != 0 {
		return e("parse")
	}
	return r
*/
/* TODO: do we need this?
func pbeg(b rv) (c rv) { // preserve strings, disambiguate +-, replace \n
	push := func(n int) { c = append(c, b[:n]...); b = b[n:] }
	for {
		if len(b) == 0 {
			return
		}
		if n := sStr(b); n > 0 {
			push(n)
		} else if n := sNum(b); n > 0 {
			push(n)
		} else if any(b[0], "\r\n+-") {
			switch b[0] {
			case '\n':
				c = append(c, ';')
			case '\r':
			default:
				c = append(c, b[0], ' ')
			}
			b = b[1:]
		} else {
			push(1)
		}
	}
}
*/

const dig = "0123456789"
const con = "π"
const sym = `+\-*%!&|<>=~,^#_$?@.`
const uni = `⍉×÷⍳⍸⌊⌽⌈⍋⌸≡∧⍴≢↑⌊↓⍕∪⍎⍣¯ℜℑ√⍟`
const uav = "⍨¨⌿⍀"
const wsp = " \t\r"

// Tokenizers return the rune count of the matched input or 0; input len > 1.
func sNum(s rv) i { // number f | fjf, allow leading +
	n := 0
l:
	for i, r := range s {
		switch {
		case i == 0 && any(r, "-+"):
		case any(r, dig):
		case i > 0 && any(r, "eEaij"): // 1j0 0i1 not 1i
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
func sNam(s rv) i { // name [a-Z][a-Z0-9]*
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
		case a(r) || (i > 0 && any(r, "0123456789")):
		default:
			return i
		}
		n++
	}
	return n
}
func sSym(s rv) i { // symbol `name
	if s[0] != '`' {
		return 0
	}
	if len(s) == 1 {
		return 1
	}
	return 1 + sNam(s[1:])
}
func sStr(s rv) i { // string "str\esc"
	if len(s) < 2 || s[0] != '"' {
		return 0
	}
	h := false
	for i, r := range s {
		switch {
		case i == 0:
		case r == '\\':
			h = !h
		case r == '"' && !h:
			return i + 1
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
func sCnd(s rv) i { return pref(s, `$[`) } // $[
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
