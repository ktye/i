package k

import (
	"fmt"
	"strings"
)

func (k *K) Run(src []byte) string {
	defer func() {
		if k.Trap == false {
			if r := recover(); r != nil {
				fmt.Println(r)
			}
		}
	}()
	e := k.parse(src)
	r := k.exec(k.fold(e, src), src)
	defer dx(r)
	if len(e) > 0 {
		if _, o := e[len(e)-1].t.(copula); o {
			return ""
		}
	}
	return sstr(r)
}

type parser struct {
	k   *K
	src []byte
	tok []token
	pos int
}
type verb string
type adverb string

type word struct {
	t    []token
	verb bool
}
type expr []word
type list []expr

type arglist []string
type makelist int
type drp struct{}

func (e expr) flat() (r []token) {
	for _, w := range e {
		for _, t := range w.t {
			r = append(r, t)
		}
	}
	return r
}
func (w word) asverb() (verb, int) {
	if len(w.t) == 1 {
		t := w.t[len(w.t)-1]
		if v, o := t.t.(verb); o {
			return v, t.p
		}
	}
	return "", 0
}
func (w word) asadverb() (adverb, int) {
	if len(w.t) > 1 {
		t := w.t[len(w.t)-1]
		if a, o := t.t.(adverb); o {
			//fmt.Printf("asadverb w.t=%#v\n", w.t)
			return a, t.p
		}
	}
	return "", 0
}
func (w word) end() (e int) {
	for _, t := range w.t {
		if t.e > e {
			e = t.e
		}
	}
	return e
}
func (l list) toks() (r []token) {
	n := len(l)
	for i := len(l) - 1; i >= 0; i-- {
		r = append(r, l[i].flat()...)
		if l[i] == nil {
			r = append(r, token{t: nil})
		}
	}
	return append(r, token{t: makelist(n)})
}

func (k *K) parse(b []byte) (r []token) {
	p := parser{k: k, src: b}

	//fmt.Println("parse>")
	k.ctx.push(b, &p.pos)
	defer k.ctx.restore("parse")

	p.tok = k.tok(b)

	r = p.e(p.t()).flat()
	for p.pos < len(p.tok) {
		if p.pos == len(p.tok) {
			break
		}
		n := p.next()
		switch n.t.(type) {
		case semicolon:
		default:
			panic("parse")
		}
		r = append(r, token{t: drp{}, p: n.p, e: n.e})
		r = append(r, p.e(p.t()).flat()...)
	}
	return r
}

func (p *parser) e(x word) (r expr) {
	//fmt.Printf("e> x = %#v\n", x)
	//defer func() { fmt.Printf("e< %#v\n", r) }()
	if x.t == nil {
		return nil
	}
	y := p.t()
	if y.verb && !x.verb {
		r = p.e(p.t())
		if v, pos := y.asverb(); y.t != nil && strings.HasSuffix(string(v), ":") {
			return p.copula(x, strings.TrimSuffix(string(v), ":"), r, pos)
		}
		prj := r == nil
		r = append(r, x)
		if prj {
			return expr{p.proj(r, y)}
		}
		return append(r, p.k.dyadic(y))
	}
	r = p.e(y)
	if t := p.train(r, x); t != nil {
		return t
	}
	if len(x.t) == 1 {
		if a, o := x.t[0].t.(adverb); o && a == "/" {
			x.t[0].t = verb(a) //regex
			x.verb = true
		}
	}
	//var m word
	var jux bool
	var juxp int
	if r != nil {
		x = p.k.monadic(x)
		if x.verb == false {
			jux, juxp = true, x.end()
		}
	}
	r = append(r, x)
	if jux {
		r = append(r, word{t: []token{token{t: f2(p.k.atx), p: juxp, e: 1 + juxp}}})
	}
	return r
}
func (p *parser) proj(r expr, y word) word {
	l := append([]token{token{t: nil}}, r.flat()...)
	l = append(l, token{t: makelist(2)})
	l = append(l, y.t...)
	l = append(l, token{t: p.k.Func["."]}, token{t: Call2{}})
	return word{t: l, verb: true}
}
func (p *parser) train(r expr, x word) expr {
	if len(r) != 1 || r[0].verb == false || x.verb == false {
		return nil
	}

	return expr{word{t: append(append(r.flat(), x.t...), token{t: Link{}}), verb: true}}
}

func (p *parser) t() (w word) {
	//defer func() { fmt.Printf("t< %#v\n", w) }()
	r := p.next()
	switch v := r.t.(type) {
	case terminator:
		p.pos--
		return word{t: nil}
	case semicolon:
		p.pos--
		return word{t: nil}
	case open:
		var al arglist
		a := r.p
		if v == '{' {
			n := p.next()
			if c, o := n.t.(open); o && c == '[' {
				l, _ := p.list()
				al = p.arglist(l)
			} else {
				p.pos--
			}
		}
		l, b := p.list()
		if len(l) == 1 {
			w.t = l[0].flat()
		} else {
			if n := p.next(); n.t != nil {
				p.pos--
				if v, o := n.t.(verb); o && strings.HasSuffix(string(v), ":") {
					al = p.arglist(l)
					return word{t: []token{token{t: al, p: n.p, e: n.e}}}
				}
			}
			w.t = l.toks()
		}
		if v == '{' {
			w = p.lambda(l, a, b, al)
		}
	default:
		if r.t != nil {
			_, v := r.t.(verb)
			w = word{t: []token{r}, verb: v}
		}
	}
aa:
	for {
		n := p.next()
		if n.t == nil {
			p.pos++
			break
		}
		switch v := n.t.(type) {
		case adverb:
			w.t = append(w.t, n)
			w.verb = true
		case open:
			if v == '[' {
				l, _ := p.list()
				if len(l) == 1 {
					w.t = append(l[0].flat(), w.t...)
					w.t = append(w.t, token{t: f2(p.k.atx), p: n.p})
				} else {
					w.t = append(append(l.toks(), w.t...), token{t: f2(p.k.call), p: n.p})
				}
			} else {
				break aa
			}
		default:
			break aa
		}
	}
	p.pos--
	return w
}

func (p *parser) span(t token) string { return string(p.src[t.p:t.e]) }
func (p *parser) list() (list, int) {
	l := list{}
	n := p.next()
	if _, o := n.t.(terminator); o {
		return l, n.p
	} else {
		p.pos--
	}
	for {
		l = append(l, p.e(p.t()))
		n := p.next()
		switch n.t.(type) {
		case terminator:
			return l, n.p
		case semicolon:
		default:
			panic("list")
		}
	}
}
func (p *parser) next() (r token) {
	if p.pos >= len(p.tok) {
		return token{}
	} else {
		p.pos++
		return p.tok[p.pos-1]
	}
}
func (p *parser) copula(x word, f string, y expr, pos int) expr {
	var c copula
	if f == ":" {
		c.g = true
	} else if len(f) > 0 {
		v, o := p.k.Func[f]
		if o == false || v.f2 == nil {
			panic("assign")
		}
		c.f = v.f2
	}
	if len(x.t) != 1 { // index-expr symbol @(or .)
		if len(x.t) < 3 {
			panic("assign")
		}
		c.i = true
		_, c.d = x.t[len(x.t)-3].t.(makelist)
		y = append(y, word{t: x.t[:len(x.t)-2]})
		x.t = []token{x.t[len(x.t)-2]}
	}
	switch v := x.t[0].t.(type) {
	case varname:
		c.s = string(v)
	case arglist:
		c.v = []string(v)
	default:
		panic("assign")
	}
	return append(y, word{t: []token{token{t: c, p: pos}}})
}
func (p *parser) arglist(l list) []string {
	v := make([]string, len(l))
	for i, e := range l {
		if len(e) != 1 {
			panic("arglist")
		}
		w := e[0]
		if len(w.t) != 1 {
			panic("arglist")
		}
		t := w.t[0]
		if n, o := t.t.(varname); o {
			v[i] = string(n)
		} else {
			panic("arglist")
		}
	}
	return v
}
func (p *parser) lambda(l list, a, b int, args []string) word {
	var f λ
	ary := map[string]int{"x": 1, "y": 2, "z": 3}
	if args != nil {
		ary = map[string]int{}
	}
	f.src = p.src[a : 1+b]
	loc := make(map[string]bool)
	var locs []string
	translate := func(t []token) []token {
		for i := range t {
			t[i].p = maxi(0, t[i].p-a)
			t[i].e = maxi(0, t[i].e-a)
		}
		return t
	}
	for i, e := range l {
		toks := e.flat()
		f.code = translate(append(f.code, toks...))
		if i < len(l)-1 {
			f.code = append(f.code, token{t: drp{}, p: a})
		}
		for _, t := range toks {
			if s, o := t.t.(varname); o {
				if a := ary[string(s)]; a > f.ary {
					f.ary = a
				}
			}
			if c, o := t.t.(copula); o && c.s != "" && c.g == false {
				loc[c.s] = true
				locs = append(locs, c.s)
			}
		}
	}
	for i := 0; i < f.ary; i++ {
		s := string('x' + byte(i))
		f.loc = append(f.loc, s)
		loc[s] = false
	}
	if args != nil {
		f.ary = len(args)
		f.loc = args
	}
	for _, s := range locs {
		if _, o := loc[s]; o {
			f.loc = append(f.loc, s)
		}
	}
	f.code = p.k.fold(f.code, f.src)
	f.save = make([]T, len(f.loc))
	f.refcount.init()
	return word{t: []token{token{t: f, p: a}}}
}
func (k *K) monadic(x word) word {
	v, p := x.asverb()
	if v != "" {
		e := p + len(v)
		x.t = append(x.t, token{t: Call1{}, p: p, e: e})
		return x
	}
	a, p := x.asadverb()
	if a != "" {
		e := p + len(v)
		x.t = append(x.t, token{t: Call1{}, p: p, e: e})
		return x
	}
	return x
}
func (k *K) dyadic(x word) word {
	v, p := x.asverb()
	if v != "" {
		e := p + len(v)
		x.t = append(x.t, token{t: Call2{}, p: p, e: e})
		return x
	}
	a, p := x.asadverb()
	if a != "" {
		e := p + len(v)
		x.t = append(x.t, token{t: Call2{}, p: p, e: e})
		return x
	}
	panic("dyadic?")
	return x
}
