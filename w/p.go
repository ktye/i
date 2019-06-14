package main

func prs(x k) (r k) { // `p"…"
	t, n := typ(x)
	if t != C || n == atom {
		panic("type")
	}
	dec(x) // but keep using it
	p := p{p: 8 + x<<2, e: n + 8 + x<<2, ln: 1}
	r = mk(L, 0)
	cat(r, mys(mk(S, 0), uint64(';')<<56))
	for p.a() { // ex;ex;…
		x = p.ex(p.noun())
		if x == 0 {
			break
		}
		cat(r, x)
		if !p.s(p.t(sSem)) {
			break
		}
	}
	if p.a() {
		p.xx() // unprocessed input
	}
	_, n = typ(r)
	if n == 1 {
		dec(r)
		return mk(C, 0) // nil, empty, `0, 0?
	} else if n == 2 {
		x = m.k[2+r+1]
		inc(x)
		dec(r)
		return x
	}
	return x
}

type p struct {
	p  k // current position, m.c[p.p:...]
	m  k // pos after matched token (token: m.c[p.p:p.m])
	e  k // pos after last byte available
	ln k // current line number
	lp k // m.c index of start of current line
}

func (p *p) t(f func([]c) int) bool { // test for next token
	p.w()
	if p.p == p.e {
		return false
	} else {
		if n := f(m.c[p.p:p.e]); n > 0 {
			p.m = p.p + n
		}
	}
	return p.m > p.p
}
func (p *p) a(f func([]c) k) (r k) { // accept, parse and advance
	n := p.m - p.p
	p.p = p.m
	return f(m.c[p.p-n : p.p])
}
func (p *p) w() { // remove whitespace and count lines
	// TODO remove comments
	for {
		switch p.get() {
		case ' ', '\t', '\r':
		case '\n':
			p.lp = p.p + 1
			p.ln++
			p.p--
			return p
		default:
			p.p--
			return p
		}
	}
}
func (p *p) get() c {
	if p.p == p.e {
		return 0
	}
	p.p++
	return m.c[p.p-1]
}
func (p *p) xx() {
	panic("parse: " + string(m.c[p.lp:p.p+1]) + " <-")
}

// Parsers
func (p *p) ex(x k) (r k) {
	switch {
	//case p.t(sNam): // TODO ... atNoun
	}
	return x
}
func (p *p) noun() (r k) {
	switch {
	case p.t(sSym):
		return p.a(pSym)
	case p.t(sNam):
		return p.a(pNam) // TODO idxr
	}
	panic("nyi")
}

func pNam(b []byte) (r k) { return btos(b) } // name: `n
func pSym(b []byte) (r k) { // `name|`"name": `n
	if len(b) == 1 || b[1] != '"' {
		r = mk(S, atom)
		mys(rc, btou(b[1:]))
	} else {
		r = pQot(b[1:])
		_, n := typ(r)
		m.k[r] = C<<28 | atom
		rc := 8 + r<<2
		mys(rc, btou(m.c[rc:rc+n]))
	}
	return r
}
func pQot(b []byte) (r k) { // "a\nb": `C
	btos // TODO
}

// Scanners return the length of the matched input or 0
func sNam(b []byte) (r int) {
	for i, c := range b {
		if cr09(c) || craZ(c) { // TODO: dot?
			if i == 0 && c09 {
				return 0
			}
		} else {
			return i
		}
	}
	return len(b)
}
func sStr(b []byte) (r int) {
	if len(b) < 2 || b[0] != '"' {
		return 0
	}
	q := false
	for i := 1; i < len(b); i++ {
		switch b[i] {
		case '\\':
			q = !q
		case '"':
			if !q {
				return i
			}
		}
	}
	return 0
}
func sSym(b []byte) (r int) { // `alp012|`"any\"thing"|`a.b.c
	if b[0] != '`' {
		return 0
	}
	if len(b) > 2 || b[1] == '"' {
		return 1 + sStr(b[1:])
	}
	for i := 1; i < len(b); i++ {
		if !(cr09(c) || craZ(c) || c == '.') {
			return i
		}
	}
	return len(b)
}
func sSem(b []byte) (r k) { return bol(b[0] == ';' || b[0] == '\n') }
func cr09(c c) bool       { return c >= '0' && c <= '9' }
func craZ(c c) bool       { return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') }
func cOps(c c) bool {
	for i, b := range []byte("+-%*|&<>=~,^#_$?@.") { // TODO: store in ktree for kwac compat
		if c == b {
			return true
		}
	}
	return false
}
