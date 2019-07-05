package main

import (
	"strconv"
)

func prs(x k) (r k) { // `p"…"
	t, n := typ(x)
	if t != C || n == atom {
		panic("type")
	}
	if n == 0 {
		dec(x)
		return mk(N, atom)
	}
	p := p{p: 8 + x<<2, e: n + 8 + x<<2, lp: 7 + x<<2, ln: 1}
	r = mk(L, 1)
	m.k[2+r] = mku(0) // ;→`
	for p.p < p.e {   // ex;ex;…
		y := p.ex(p.noun())
		if x == 0 {
			break
		}
		r = lcat(r, y)
		if !p.t(sSem) {
			break
		}
		p.p++
	}
	if p.p < p.e {
		p.xx() // unprocessed input
	}
	_, n = typ(r)
	if n == 2 {
		inc(m.k[3+r])
		dec(r)
		r = m.k[3+r] // r[1]
	}
	dec(x)
	return r
}

type p struct {
	p  k // current position, m.c[p.p:...]
	m  k // pos after matched token (token: m.c[p.p:p.m])
	e  k // pos after last byte available
	ln k // current line number
	lp k // m.c index of last newline
}

func (p *p) t(f func([]c) int) bool { // test for next token
	p.w()
	if p.p == p.e {
		return false
	}
	if n := f(m.c[p.p:p.e]); n > 0 {
		p.m = p.p + k(n)
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
		case 0:
			return
		case ' ', '\t', '\r':
		case '\n':
			if p.p != p.lp+1 {
				p.lp = p.p - 1
				p.ln++
			}
			p.p--
			return
		case '/':
			if p.p == p.lp+2 || m.c[p.p-2] == ' ' || m.c[p.p-2] == '\t' {
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
	panic("parse: " + string(m.c[p.lp+1:p.p+1]) + " <-")
}

// Parsers
func (p *p) ex(x k) (r k) {
	t, _ := typ(x)
	switch t {
	case N:
		return x
	case N + 1, N + 2, N + 3, N + 4:
		if p.t(sAdv) {
			return p.adv(0, x)
		}
		v := p.noun()
		if p.t(sAdv) {
			return p.adv(x, v) // e.g. f g/x
		}
		y := p.ex(v)
		if yt := m.k[y] >> 28; yt == N { // verb only
			dec(y)
			return x
		}
		r = mk(L, 2) // monadic application
		m.k[2+r] = monad(inc(x))
		m.k[3+r] = y
		dec(x)
		return r
	default:
		if p.t(sAdv) { // e.g. 2/3
			return p.adv(0, x)
		}
		a := p.noun()
		if p.t(sAdv) {
			return p.adv(x, a)
		}
		at, _ := typ(a)
		if at == N {
			dec(a)
			return x
		} else if at > N {
			y := p.ex(p.noun())
			if m.k[y]>>28 == N { // projection, e.g. 2+
				dec(y)
				r = mk(L, 2)
				m.k[2+r] = a
				m.k[3+r] = x
				return r
			}
			r = mk(L, 3) // dyadic application (infix)
			m.k[2+r] = a
			m.k[3+r] = x
			m.k[4+r] = y
			return r
		} else if p.t(sAdv) {
			return p.adv(x, a)
		}
		return lcat(lcat(mk(L, 0), x), p.ex(a)) // (x;…) e.g. `p".."
	}
}
func (p *p) noun() (r k) {
	switch {
	case p.t(sHex):
		r = p.a(pHex)
		return p.idxr(r)
	case p.t(sIov):
		return p.idxr(p.a(pIov))
	case p.t(sNum):
		r = p.a(pNum)
		for p.t(sNum) {
			y := p.a(pNum)
			rt, yt, _, _ := typs(r, y)
			if rt < yt {
				r = to(r, yt)
			} else if yt < rt {
				y = to(y, rt)
			}
			r = cat(r, y)
		}
		return p.idxr(r)
	case p.t(sStr):
		r = p.a(pStr)
		return p.idxr(r)
	case p.t(sSym):
		r = p.a(pSym)
		for p.p != p.e && m.c[p.p] == '`' { // `a`b`c without whitespace
			if p.t(sSym) {
				r = cat(r, p.a(pSym))
			}
		}
		return p.idxr(enl(r))
	case p.t(sOcb):
		p.p = p.m
		st := p.m - 1
		r = mk(N+1, 0)                   // lambda is indicated with length 0 but uses 2 fields:
		m.k[3+r] = p.lst(mk(C, 0), sCcb) // #2: parse tree
		dst, n := mk(C, p.p-st), p.p-st
		m.k[2+r] = dst // #1: string representation
		dst = 8 + dst<<2
		copy(m.c[dst:dst+n], m.c[st:p.p])
		args := argn(m.k[3+r], 0)
		if args == 0 {
			panic("valence(lambda)")
		}
		m.k[r] = (N+args)<<28 | 0
		return p.idxr(r)
	case p.t(sOpa):
		p.p = p.m
		r = p.lst(mk(C, 0), sCpa)
		if m.k[r]&atom == 1 {
			// TODO drv
			return p.idxr(fst(r))
		}
		// TODO isverb
		return p.idxr(cat(enlist(mk(N, atom)), r))
	case p.t(sVrb):
		return p.idxr(p.a(pVrb))
	case p.t(sBin):
		return p.idxr(p.a(pBin))
	case p.t(sNam):
		return p.idxr(p.a(pNam))
	}
	return mk(N, atom)
}
func (p *p) idxr(x k) (r k) { // […]
	if p.t(sObr) {
		p.p = p.m
		r = mk(L, 1)
		m.k[2+r] = x
		return p.lst(r, sCbr)
	}
	return x
}
func (p *p) lst(l k, term func([]c) int) (r k) { // append to l
	r = l
	if p.t(term) {
		p.p = p.m
		return r
	}
	for {
		r = lcat(r, p.ex(p.noun()))
		if !p.t(sSem) {
			break
		}
		p.p = p.m
	}
	if p.t(term) == false {
		panic("parse: unclosed list")
	}
	p.p = p.m
	return r
}
func (p *p) adv(x, v k) (r k) { // x(left arg) v(verb)
	vt, vn := typ(v)
	a := p.a(pAdv)
	if p.t(sAdv) {
		b := p.a(pAdv)
		if vt == N+2 && vn == atom && m.k[2+v] < 83 { // force monad ambivalent primitive (reflite p131)
			m.k[v] = (N+1)<<28 | atom
			m.k[2+v] -= 50
		}
		v, vn = lcat(lcat(mk(L, 0), a), v), 0
		a = b
	}
	y := p.ex(p.noun())
	if x == 0 {
		if yt, _ := typ(y); yt == N {
			dec(y)
			return lcat(lcat(mk(L, 0), a), v) // (a;v)
		}
		return lcat(lcat(mk(L, 0), lcat(lcat(mk(L, 0), a), v)), y) // ((a;v);y)
	}
	r = mk(L, 3) // ((a;v);x;y)
	m.k[2+r] = lcat(lcat(mk(L, 0), a), v)
	m.k[3+r] = x
	m.k[4+r] = y
	return r
}
func monad(x k) (r k) { // force monad
	t, _ := typ(x)
	if t == N+1 {
		return x
	} else if t == N+2 {
		r = mk(N+1, atom)
		m.k[2+r] = m.k[2+x] - 50
		if m.k[2+x] > 99 {
			panic("parse monad")
		}
		dec(x)
		return r
	}
	panic("type")
}
func pHex(b []byte) (r k) { // 0x1234 `c|`C
	if n := k(len(b)); n == 3 { // allow short form 0x1
		r = mkc(xtoc(b[2]))
	} else if n%2 != 0 {
		panic("parse hex")
	} else {
		n = (n - 2) / 2
		r, b = mk(C, n), b[2:]
		rc := 8 + r<<2
		for i := k(0); i < n; i++ {
			m.c[rc+i] = (xtoc(b[2*i]) << 4) | xtoc(b[2*i+1])
		}
		if n == 1 {
			m.k[r] = C<<28 | atom
		}
	}
	return r
}
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
func pNum(b []byte) (r k) { // 0|1f|-2.3e+4|1i2: `i|`f|`z
	for i, c := range b {
		if c == 'i' {
			r = to(pNum(b[:i]), Z)
			y := to(pNum(b[i+1:]), F)
			m.f[3+r>>1] = m.f[1+y>>1]
			dec(y)
			return r
		}
	}
	f := 0
	if len(b) > 1 {
		if c := b[len(b)-1]; c == 'f' || c == '.' {
			b = b[:len(b)-1]
			f = 1
		} else if b[0] == '.' {
			b = b[1:]
			f = 2
		}
	}
	if x, o := atoi(b); o {
		if f == 0 {
			r = mki(k(i(x)))
		} else {
			r = mk(F, atom)
			if f == 1 {
				m.f[1+r>>1] = float64(x)
			} else {
				m.f[1+r>>1] = 0.1 * float64(x)
			}
		}
		return r
	}
	if x, err := strconv.ParseFloat(string(b), 64); err == nil { // TODO remove strconv
		r = mk(F, atom)
		m.f[1+r>>1] = x
		return r
	}
	panic("parse number")
}
func pStr(b []byte) (r k) { // "a"|"a\nbc": `c|`C
	r = pQot(b)
	if _, n := typ(r); n == 1 {
		m.k[r] = C<<28 | atom
	}
	return r
}
func pNam(b []byte) (r k) { // name: `n
	return mku(btou(b))
}
func pSym(b []byte) (r k) { // `name|`"name": `n
	if len(b) == 1 {
		r = mku(0)
	} else if len(b) > 1 || b[1] != '"' {
		r = mku(btou(b[1:]))
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
	return srk(r, C, k(len(b)-2), k(p-(8+r<<2)))
}
func pVrb(b []byte) (r k) {
	for i := k(0); i < 34; i++ { // :+-*%&|<>=!~,^#_$?@.0123456789'/\
		if b[0] == m.c[i+136] {
			if len(b) == 1 || i > 49 {
				r = mk(N+2, atom)
				m.k[2+r] = k(i) + 50
			} else {
				r = mk(N+1, atom)
				m.k[2+r] = k(i)
			}
			return r
		}
	}
	panic("pVrb")
}
func pIov(b []byte) (r k) {
	r = mk(N+2, atom)
	m.k[2+r] = 70 + k(b[0]-'0') // ioverb parses always as `2
	return r
}
func pAdv(b []byte) (r k) {
	f := k(80)
	if len(b) > 1 {
		f -= 50
	}
	if b[0] == '/' {
		f++
	} else if b[0] == '\\' {
		f += 2
	}
	r = mk(N+1, atom)
	m.k[2+r] = f
	return r
}
func pBin(b []byte) (r k) { // builtin
	x := mku(btou(b))
	x = atx(inc(m.k[3]), x)
	if xt, xn := typ(x); xt != C && xn != atom {
		panic("parse builtin")
	}
	if f := m.c[8+x<<2]; f < 50 {
		r = mk(N+1, atom)
		m.k[2+r] = k(f)
	} else {
		r = mk(N+2, atom)
		m.k[2+r] = k(f)
	}
	dec(x)
	return r
}

// Scanners return the length of the matched input or 0
func sHex(b []byte) (r int) {
	if !(len(b) > 1 && b[0] == '0' && b[1] == 'x') {
		return 0
	}
	for i, c := range b[2:] {
		if crHx(c) == false {
			return 2 + i
		}
	}
	return len(b)
}
func sNum(b []byte) (r int) {
	n := sFlt(b)
	if n > 0 && len(b) > n && b[n] == 'i' {
		n += 1 + sFlt(b[n+1:])
	}
	return n
}
func sFlt(b []byte) (r int) { // -0.12e-12|1f
	if len(b) > 1 && b[0] == '-' {
		r++
	}
	for i := r; i < len(b); i++ {
		if c := b[i]; cr09(c) {
			r++
		} else {
			if c == '.' {
				d := sDec(b[i+1:])
				if i == 0 && d == 0 {
					return 0 // 1. or .1 is allowed but not .
				}
				r += 1 + sDec(b[i+1:])
			}
			break
		}
	}
	if len(b) > r && b[r] == 'e' {
		r += 1 + sExp(b[r+1:])
	}
	if len(b) > r && b[r] == 'f' {
		r++
	}
	if r == 1 && b[0] == '-' {
		return 0
	}
	return r
}
func sDec(b []byte) (r int) {
	for _, c := range b {
		if !cr09(c) {
			break
		}
		r++
	}
	return r
}
func sExp(b []byte) (r int) {
	if len(b) > 0 && (b[0] == '+' || b[0] == '-') {
		r++
	}
	for i := r; i < len(b); i++ {
		if c := b[i]; !cr09(c) {
			break
		}
		r++
	}
	return r
}
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
	for i, c := range b[1:] {
		if !q && c == '\\' {
			q = true
		} else {
			if q == false && c == '"' {
				return i + 2
			}
			q = false
		}
	}
	return 0
}
func sSym(b []byte) (r int) { // `alp012|`"any\"thing"|`a.b.c
	if b[0] != '`' {
		return 0
	}
	if len(b) > 2 && b[1] == '"' {
		return 1 + sStr(b[1:])
	}
	for i, c := range b[1:] {
		if !(cr09(c) || craZ(c) || c == '.') {
			return 1 + i
		}
	}
	return len(b)
}
func sSem(b []byte) int {
	if b[0] == ';' || b[0] == '\n' {
		return 1
	}
	return 0
}
func sIov(b []byte) int { // ioverb 0: .. 4:
	if len(b) > 1 && b[1] == ':' && b[0] >= '0' && b[0] <= '4' {
		return 2
	}
	return 0
}
func sVrb(b []byte) int {
	if cOps(b[0]) {
		if len(b) > 1 && b[1] == ':' {
			return 2
		}
		return 1
	}
	return 0
}
func sAdv(b []byte) int {
	if cAdv(b[0]) {
		if len(b) > 1 && b[1] == ':' {
			return 2
		}
		return 1
	}
	return 0
}
func sBin(b []byte) int { // builtin
	if b[0] < 'a' || b[0] > 'z' {
		return 0
	}
	x := til(inc(m.k[3]))
	_, n := typ(x)
	xp, max := 8+x<<2, k(len(b))
	for i := k(0); i < n; i++ {
		u := sym(xp + 8*i)
		for j := k(0); j < 8; j++ {
			if c := c(u >> (8 * (7 - j))); c == 0 {
				if j == max || !(cr09(b[j]) || craZ(b[j])) {
					return int(j)
				}
			} else if j == max {
				break
			} else if c != b[j] {
				break
			}
		}
	}
	dec(x)
	return 0
}
func sObr(b []byte) int { return ib(b[0] == '[') }
func sOpa(b []byte) int { return ib(b[0] == '(') }
func sOcb(b []byte) int { return ib(b[0] == '{') }
func sCbr(b []byte) int { return ib(b[0] == ']') }
func sCpa(b []byte) int { return ib(b[0] == ')') }
func sCcb(b []byte) int { return ib(b[0] == '}') }
func cr09(c c) bool     { return c >= '0' && c <= '9' }
func craZ(c c) bool     { return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') }
func crHx(c c) bool     { return cr09(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F') }
func cOps(c c) bool {
	for _, b := range m.c[136 : 136+20] { // :+-*%&|<>=!~,^#_$?@.
		if c == b {
			return true
		}
	}
	return cAdv(c)
}
func cAdv(c c) bool {
	if c == '\'' || c == '/' || c == '\\' {
		return true
	}
	return false
}
func ib(b bool) (r int) {
	if b {
		r = 1
	}
	return r
}
func argn(x, a k) k { // count args of lambda parse tree
	t, n := typ(x)
	switch t {
	case S:
		if n == atom {
			n = 1
		}
		ux, uy, uz := uint64('x')<<56, uint64('y')<<56, uint64('z')<<56
		for i := k(0); i < n; i++ {
			if u := sym((2 + 2*i + x) << 2); u == ux && a < 1 {
				a = 1
			} else if u == uy && a < 2 {
				a = 2
			} else if u == uz && a < 3 {
				a = 3
			}
		}
	case L:
		for i := k(0); i < n; i++ {
			a = argn(m.k[2+i+x], a)
		}
	}
	return a
}
func atoi(b []c) (int, bool) {
	n, s := 0, 1
	for i, c := range b {
		if i == 0 && c == '-' {
			s = -1
		} else if cr09(c) {
			n = 10*n + int(c-'0')
		} else {
			return 0, false
		}
	}
	return s * n, true
}
