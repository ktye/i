package main

func prs(x k) (r k) { // `p"…"
	t, n := typ(x)
	if t != C || n == atom {
		panic("type")
	}
	dec(x) // but keep using it
	if n == 0 {
		return mk(N, atom)
	}
	p := p{p: 8 + x<<2, e: n + 8 + x<<2, ln: 1}
	r = mk(L, 0)
	cat(r, mys(mk(S, 0), uint64(';')<<56))
	for p.p < p.e { // ex;ex;…
		x = p.ex(p.noun())
		if x == 0 {
			break
		}
		r = cat(r, x)
		if !p.t(sSem) {
			break
		}
	}
	if p.p < p.e {
		p.xx() // unprocessed input
	}
	_, n = typ(r)
	if n == 1 {
		inc(m.k[2+r])
		dec(r)
		return m.k[2+r] // r[0]
	}
	return r
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
			p.m = p.p + k(n)
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
	for {
		switch p.get() {
		case ' ', '\t', '\r':
		case '\n':
			p.lp = p.p + 1
			p.ln++
			p.p--
			return
		case '/':
			if p.p == p.ln+1 || m.c[p.p]-1 == ' ' || m.c[p.p]-1 == '\t' {
				for {
					if c := p.get(); c == 0 {
						return
					} else if c == '\n' {
						p.p--
						break
					}
				}
			} else {
				p.p--
				return
			}
		default:
			p.p--
			return
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
		r = p.a(pSym)
		for p.p != p.e && m.c[p.p] == '`' { // `a`b`c without whitespace
			if p.t(sSym) {
				r = cat(r, p.a(pSym))
			}
		}
		if m.k[r]&atom == atom { // `a → ,`a
			m.k[r] = S<<28 + 1
		} else { // `a`b → ,`a`b
			r = enl(r)
		}
		return p.idxr(r)
	case p.t(sNam):
		return p.idxr(p.a(pNam))
	}
	panic("nyi")
}
func (p *p) idxr(x k) (r k) { return x } // TODO

func pNam(b []byte) (r k) { // name: `n
	r = mk(S, atom)
	mys(8+r<<2, btou(b))
	return r
}
func pSym(b []byte) (r k) { // `name|`"name": `n
	if len(b) == 1 || b[1] != '"' {
		r = mk(S, atom)
		mys(8+r<<2, btou(b[1:]))
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
	r = mk(C, k(len(b)-2))
	p := 8 + r<<2
	q := false
	for _, c := range b[1 : len(b)-1] {
		if c == '\\' && !q {
			q = true
		} else {
			if q {
				q = false
				switch c {
				case 'r':
					c = '\r'
				case 'n':
					c = '\n'
				case 't':
					c = '\t'
				}
			}
			m.c[p] = c
			p++
		}
	}
	srk(r, C, k(len(b)-2), k(p-(8+r<<2)))
	return r
}

// Scanners return the length of the matched input or 0
func sNam(b []byte) (r int) {
	for i, c := range b {
		if cr09(c) || craZ(c) { // TODO: dot?
			if i == 0 && cr09(c) {
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
	for i, c := range b[1:] {
		if !(cr09(c) || craZ(c) || c == '.') {
			return i
		}
	}
	return len(b)
}
func sSem(b []byte) (r int) {
	if b[0] == ';' || b[0] == '\n' {
		r = 1
	}
	return r
}
func cr09(c c) bool { return c >= '0' && c <= '9' }
func craZ(c c) bool { return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') }
func cOps(c c) bool {
	for _, b := range []byte("+-%*|&<>=~,^#_$?@.") { // TODO: store in ktree for kwac compat
		if c == b {
			return true
		}
	}
	return false
}
