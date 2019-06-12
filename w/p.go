package main

func prs(x k) (r k) { // `p"…"
	t, n := typ(x)
	if t != C || n == atom {
		panic("type")
	}
	dec(x) // but keep using it
	p := p{p: 8 + x<<2, n: n}
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
	p  k // m.c index to current position (after consumed part)
	c  k // start of matched token (end is p)
	n  k // bytes available
	ln k // current line number
	lp k // m.c index of start of current line
}

func (p *p) a() bool { return p.w().n > 0 } // available
func (p *p) t(f func([]byte) k) k { // test but keep token
	if p.w().n == 0 {
		return 0
	} else {
		return f(p.b())
	}
}
func (p *p) s(n k) bool { p.c = p.p; p.p += n; p.n -= n; return n > 0 } // shift
func (p *p) m(f func([]byte) k) []c { // must match token, advance and return length
	if n := p.t(f); n > 0 {
		p.s(n)
		return p.b()
	}
	p.xx()
	return nil
}
func (p *p) any(f []func([]c) k) bool {
	if p.w().n == 0 {
		return false
	}
	for i := range f {
		if f[i](p.b()) > 0 {
			return true
		}
	}
	return false
}
func (p *p) b() []byte { // remaining bytes
	return m.c[p.p : p.p+p.n]
}
func (p *p) w() *p { // advance to next token, replace and count newlines
	// TODO remove comments
	p.bak()
	for {
		switch p.get() {
		case ' ', '\t', '\r':
		case '\n':
			p.lp = p.p
			p.ln++
			p.bak()
			return p
		default:
			p.bak()
			return p
		}
	}
}
func (p *p) get() c {
	if p.n == 0 {
		return 0
	}
	p.p++
	p.n--
	return m.c[p.p-1]
}
func (p *p) bak() { p.p++; p.n-- }
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
	case p.t(sNam) > 0:
		return btos(p.m(sNam)) // TODO idxr
		// TODO colon
	}
	panic("nyi")
}

// Scanners return the length of the matched input or 0
func sNam(b []byte) (r k) {
	for i, c := range b {
		if (c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			if i == 0 && c >= '0' && c <= '9' {
				return 0
			}
		} else {
			return k(i)
		}
	}
	return k(len(b))
}
func sSem(b []byte) (r k) { return bol(b[0] == ';' || b[0] == '\n') }
