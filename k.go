package main

import (
	"math"
	"unsafe"
)

const ref = `
00 : idn dex    20 0 rdl wrl    40 exi exit  120 ... in       60 prm   140
01 + flp add    21 1 nil nil    41 sqr sqrt  121 ... within   61       141
02 - neg sub    22 2 nil nil    42 sin       122 bin          62       142
03 * fst mul    23 3 nil nil    43 cos       123 ... like     63       143
04 % inv div    24 4 nil nil    44 dev       124 del          64       144
05 & wer min    25 5 nil nil    45 log       125 lgn log      65       145
06 | rev max    26 6 nil nil    46 exp       126 pow exp      66       146
07 < asc les    27 7 nil nil    47 rnd rand  127 rol rand     67       147
08 > dst mor    28 8 lun nil    48 abs       128              68       148
09 = grp eql    29 9 nil nil    49 plo plot  129 plt plot     69       149
                                                                          
10 ! til key    30 ' qtc qot    50 rel real  130 mkz cmplx    70       150
11 ~ not mch    31 / slc sla    51 ima imag  131 fns find     71       151
12 , enl cat    32 \ bsc bsl    52 phi phase 132 rot          72       152
13 ^ srt ept    33 ' ech ecd    53 cnj conj  133              73       153
14 # cnt tak    34 / ovr ovi    54 cnd cond  134              74       154
15 _ flr drp    35 \ scn sci    55 zxp expi  135 rxp expi     75       155
16 $ str cst    36 ' ecp epi    56 dia diag  136              76       156
17 ? unq fnd    37 / jon ecr    57 avg       137 mvg avg      77       157
18 @ tip atx    38 \ spl ecl    58 med       138 pct med      78       158
19 . val cal    39 /     trp    59 vri var   139 cov var      79       159
`

type c = byte
type k = uint32
type i = int32
type f = float64
type z = complex128
type s = string

const (
	C, I, F, Z, S, L, A, N                                    k = 1, 2, 3, 4, 5, 6, 7, 8
	atom, srcp, kkey, kval, stab, asci, dyad                  k = 0x0fffffff, 0x2f, 0x30, 0x31, 0x32, 0x33, 80
	NaI                                                       i = -2147483648
	yb64, yhex, ycsv, ypng, ysel, yudt, ydel, yby, yfrm, ywer k = 257, 258, 259, 260, 261, 262, 263, 264, 265, 266
)

type (
	f1 func(k, k)
	f2 func(k, k, k)
	fc func(k, k) bool
)
type slice struct {
	p uintptr
	l int
	c int
}

//                C  I  F   Z  S  L  A  0+
var lns = [9]k{0, 1, 4, 8, 16, 4, 4, 0, 0}
var m struct { // linear memory (slices share underlying array)
	c []c
	k []k
	f []f
	z []z
}
var null, nans k                    // mk(N,atom), mk(S,nil)
var nan = [7]k{0, 0, 0, 0, 0, 0, 0} // nan[t] (missing) mk(t,atom)..
var unan, inan, uinf = uint64(0x7FF8000000000001), k(0x80000000), uint64(0x7FF0000000000000)
var fnan, finf f = *(*f)(unsafe.Pointer(&unan)), *(*f)(unsafe.Pointer(&uinf))
var cpx = []f1{nil, cpC, cpI, cpF, cpZ, cpI, cpL}      // copy
var eqx = []fc{nil, eqC, eqI, eqF, eqZ, eqS, nil}      // equal
var ltx = []fc{nil, ltC, ltI, ltF, ltZ, ltS}           // less than
var gtx = []fc{nil, gtC, gtI, gtF, gtZ, gtS, nil}      // greater than (gtL causes init loop)
var stx = []func(k, k) k{nil, nil, stI, stF, stZ, stS} // tostring (assumes 56 bytes space at dst)
var tox = []f1{nil, func(r, x k) { m.k[r] = k(i(m.c[x])) }, func(r, x k) { m.f[r] = f(m.c[x]) }, func(r, x k) { m.z[r] = complex(f(m.c[x]), 0) }, func(r, x k) { m.c[r] = c(m.k[x]) }, nil, func(r, x k) {
	m.f[r] = f(i(m.k[x]))
	if i(m.k[x]) == NaI {
		m.f[r] = fnan
	}
}, func(r, x k) {
	m.z[r] = complex(f(i(m.k[x])), 0)
	if i(m.k[x]) == NaI {
		m.z[r] = complex(fnan, fnan)
	}
}, func(r, x k) { m.c[r] = c(m.f[x]) }, func(r, x k) {
	m.k[r] = k(i(f(m.f[x])))
	if math.IsNaN(m.f[x]) {
		m.k[r] = inan
	}
}, nil, func(r, x k) { m.z[r] = complex(m.f[x], 0) }, func(r, x k) { m.c[r] = c(m.f[x<<1]) }, func(r, x k) {
	m.k[r] = k(i(m.f[x<<1]))
	if math.IsNaN(m.f[x<<1]) || math.IsNaN(m.f[1+x<<1]) {
		m.k[r] = inan
	}
}, func(r, x k) { m.f[r] = m.f[x<<1] }}
var table [160]interface{} // function table :+-*%&|<>=!~,^#_$?@.0123456789'/\

func ini(mem []f) { // start function
	table = [160]interface{}{
		//   1                   5                        10                       15
		idn, flp, neg, fst, inv, wer, rev, asc, dsc, grp, til, not, enl, srt, cnt, flr, str, unq, tip, val, //  00- 19
		rdl, nil, nil, nil, nil, nil, nil, nil, lun, deb, qtc, slc, bsc, ech, ovr, scn, ecp, jon, spl, nil, //  20- 39
		nil, sqr, sin, cos, dev, log, exp, rnd, abs, nil, rel, ima, phi, cnj, cnd, zxp, dia, avg, med, vri, //  40- 59
		prm, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, //  60- 79
		dex, add, sub, mul, div, min, max, les, mor, eql, key, mch, cat, ept, tak, drp, cst, fnd, atx, cal, //  80- 99
		wrl, nil, nil, nil, nil, nil, nil, nil, nil, nil, qot, sla, bsl, ecd, ovi, sci, epi, ecr, ecl, nil, // 100-119
		nil, nil, bin, nil, del, lgn, pow, rol, nil, nil, mkz, fns, rot, nil, nil, rxp, nil, mvg, pct, cov, // 120-139
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, // 140-159
	}
	if len(mem) != 1<<13 {
		panic("ini")
	}
	m.f = mem
	msl()
	m.k[2] = 16
	p := k(64)
	m.k[8] = p
	m.k[p] = 8
	for i := 9; i < 16; i++ {
		p *= 2
		m.k[i] = p
		m.k[p] = k(i)
	}
	m.k[0] = (I << 28) | 31
	m.k[1] = 0x70881342
	copy(m.c[136:169], []c(`:+-*%&|<>=!~,^#_$?@.0123456789'/\`))
	copy(m.c[169:177], []c{0, 'c', 'i', 'f', 'z', 'n', '.', 'a'})

	m.k[stab] = spl(mkc(','), mkb([]byte(",b64,hex,csv,png,select,update,delete,by,from,where"))) // symbol table
	nans = mk(S, atom)
	m.k[2+nans] = 0
	null = mk(N, atom)
	nan[C], nan[I], nan[F], nan[Z], nan[S] = mk(C, atom), mk(I, atom), mk(F, atom), mk(Z, atom), nans
	m.c[ptr(nan[C], C)] = 32
	m.k[ptr(nan[I], I)] = k(inan)
	m.f[ptr(nan[F], F)] = fnan
	m.f[1+ptr(nan[Z], F)], m.f[2+ptr(nan[Z], F)] = fnan, fnan
	nan[L] = nan[C] // shared
	n := mk(L, 5)
	for i := k(0); i < 5; i++ {
		m.k[2+n+i] = nan[i+1]
	}
	m.k[asci] = to(jota(256), C)
	m.k[kkey] = mk(S, 0) // k-tree keys
	m.k[kval] = mk(L, 0) // k-tree values
	m.k[3] = mk(A, atom)
	m.k[2+m.k[3]] = mk(S, 0)
	m.k[3+m.k[3]] = mk(C, 0)
	gtx[L] = gtL
	eqx[L] = eqL
	builtins(40, "exit,sqrt,sin,cos,dev,,,,abs,,real,imag,phase,conj,cond,nyi15,diag,,,,prm")                           // monads
	builtins(c(40+dyad), "in,within,bin,like,del,log,exp,rand,,plot,cmplx,find,rot,nyi13,nyi14,expi,nyi16,avg,med,var") // dyads

	dec(asn(mks(".a"), m.k[asci], inc(null)))
	dec(asn(mks(".0"), null, inc(null)))
	dec(asn(mks(".n"), n, inc(null)))
	//dec(asn(inc(nans), key(mk(S, 0), mk(L, 0)), inc(null))) // ktree `
	dec(asn(mks(".f"), mk(C, 0), inc(null)))        // file name
	dec(asn(mks(".c"), mk(C, 0), inc(null)))        // current src
	mkk(".flp", `{(,/x[;!n])@(n*!#x)+/:!n:|/#:'x}`) // transpose
	// mkk(".odo", `{x\:!*/x}`)                             // odometer (replaced by native, \: is very slow)
	mkk(".rot", `{$[x~0;y;0~#y;y];x:(#y)\x;$[0<x;(x_y),x#y;(x#y),x_y]}`)
	mkk(".csv", "{a:`A~@x;(h;d):$[a;(!+x;.+x);(0#`;x)];r:`z=@:'d;h:h@&1+r;d:d@&1+r;d:@[d;&r;abs];d:@[d;1+&r;(180%1p)*phase];$[a;((,\",\"/:$h),\",\"/:'+$:'d);\",\"/:'+$:'d]}")
	mkk(".vsc", "{(t;s):(x[0];x[1]);y:+s\\:'y;$[0=#t;y;?[,/'(`$'t)$'(#t)#y;(&\" \"=t),'1;()]]}")
	mkk(".stats", "{s:_2 exp !#x;l:(!#x;x;y;s*x;s*y);t:+`b`used`free`uB`fB!l,'+/+l;` 1:(\"\n\"/:`m@t),\"\n\";}")
}
func cpC(dst, src k)  { m.c[dst] = m.c[src] }
func cpI(dst, src k)  { m.k[dst] = m.k[src] }
func cpF(dst, src k)  { m.f[dst] = m.f[src] }
func cpZ(dst, src k)  { m.z[dst] = m.z[src] }
func cpL(dst, src k)  { m.k[dst] = inc(m.k[src]) }
func swI(dst, src k)  { m.k[dst], m.k[src] = m.k[src], m.k[dst] }
func eqC(x, y k) bool { return m.c[x] == m.c[y] }
func eqI(x, y k) bool { return i(m.k[x]) == i(m.k[y]) }
func eqF(x, y k) bool { return m.k[x<<1] == m.k[y<<1] && m.k[1+x<<1] == m.k[1+y<<1] } // return m.f[x] == m.f[y] || (m.f[x] != m.f[x] && m.f[y] != m.f[y]) }
func eqZ(x, y k) bool { return eqF(x<<1, y<<1) && eqF(1+x<<1, 1+y<<1) }
func eqS(x, y k) bool { return m.k[x] == m.k[y] }
func eqL(x, y k) bool { return match(m.k[x], m.k[y]) }
func ltC(x, y k) bool { return m.c[x] < m.c[y] }
func gtC(x, y k) bool { return m.c[x] > m.c[y] }
func ltI(x, y k) bool { return i(m.k[x]) < i(m.k[y]) }
func gtI(x, y k) bool { return i(m.k[x]) > i(m.k[y]) }
func ltF(x, y k) bool { return m.f[x] < m.f[y] }
func gtF(x, y k) bool { return m.f[x] > m.f[y] }
func ltZ(x, y k) bool { // real then imag
	if ltF(x<<1, y<<1) {
		return true
	} else if eqF(x<<1, y<<1) {
		return ltF(1+x<<1, 1+y<<1)
	}
	return false
}
func gtZ(x, y k) bool {
	if gtF(x<<1, y<<1) {
		return true
	} else if eqF(x<<1, y<<1) {
		return gtF(1+x<<1, 1+y<<1)
	}
	return false
}
func gtS(x, y k) bool { return m.k[x] != m.k[y] && ltS(x, y) == false }
func ltS(x, y k) bool {
	xp, xn := sc(x)
	yp, yn := sc(y)
	mn := xn
	if yn < mn {
		mn = yn
	}
	for i := k(0); i < mn; i++ {
		if m.c[xp+i] < m.c[yp+i] {
			return true
		} else if m.c[xp+i] > m.c[yp+i] {
			return false
		}
	}
	return xn < yn
}
func sym(s s) (r k) {
	r = mks(s)
	return dex(r, m.k[2+r])
}
func sc(x k) (r, n k) {
	p := m.k[x]
	if p == 0 {
		return 0, 0
	} else if p < 256 {
		return 8 + p + m.k[asci]<<2, 1
	}
	p -= 256
	if mx := m.k[m.k[stab]] & atom; p > mx { // TODO rm
		panic("sc")
	}
	r = m.k[m.k[stab]+2+p]
	return 8 + r<<2, m.k[r] & atom
}
func gtL(x, y k) bool {
	x, y = m.k[x], m.k[y]
	xt, yt, xn, yn := typs(x, y)
	if xt == yt && xt < L {
		return gtx[xt](ptr(x, xt), ptr(y, yt))
	} else if xt != yt {
		return xt > yt
	} else if xt == L {
		mn := atm1(xn)
		if n := atm1(yn); n < xn {
			mn = n
		}
		if mn == 0 {
			return true
		}
		for i := k(0); i < mn; i++ {
			if xi, yi := 2+x+i, 2+y+i; !match(m.k[xi], m.k[yi]) {
				return gtL(xi, yi)
			}
		}
		return atm1(xn) > atm1(yn)
	} else {
		panic("type")
	}
}
func stI(dst, x k) k {
	if m.k[x] == 0x80000000 {
		m.c[dst] = '0'
		m.c[dst+1] = 'N'
		return 2
	}
	v, s, b, n := i(m.k[x]), k(0), m.c[dst:], k(0)
	if v < 0 {
		v, s, b[0] = -v, 1, '-'
		b = b[1:]
	}
	for i := 0; ; i++ {
		b[i] = c(v%10 + '0')
		n++
		if v < 10 {
			break
		}
		v /= 10
	}
	for i := k(0); i < n/2; i++ {
		b[i], b[n-i-1] = b[n-i-1], b[i]
	}
	return s + n
}
func stF(dst, x k) k { return ftoa(dst, m.f[x]) }
func stZ(dst, x k) k {
	re, im := m.f[x<<1], m.f[1+x<<1]
	if im == 0 {
		return ftoa(dst, re)
	}
	n := ftoa(dst, math.Hypot(re, im))
	re = 180.0 * math.Atan2(im, re) / math.Pi
	if re < 0 {
		re += 360.0
	}
	m.c[dst+n] = 'a'
	return 1 + n + ftoa(dst+1+n, re)
}
func stS(dst, x k) k {
	if p := m.k[x]; p == 0 {
		return 0
	} else if p < 256 {
		m.c[dst] = c(p)
		return 1
	}
	c, n := sc(x)
	copy(m.c[dst:], m.c[c:c+n])
	return n
}
func ptr(x, t k) k { // convert k address to type dependend index of data section
	switch t {
	case C:
		return (2 + x) << 2
	case F:
		return (2 + x) >> 1
	case Z:
		return (4 + x) >> 2
	default:
		return 2 + x
	}
}
func mv(dst, src k) {
	t, n := typ(src)
	ln := k(1 << bk(t, n))
	rc := m.k[1+dst]
	dst, src = dst<<2, src<<2
	copy(m.c[dst:dst+ln], m.c[src:src+ln]) // copy bucket
	dst >>= 2
	m.k[dst] = t<<28 | n // restore header
	m.k[1+dst] = rc
}
func builtins(code c, s string) {
	v := spl(mkc(','), mkb([]byte(s)))
	for i := c(0); i < c(m.k[v]&atom); i++ {
		if ln := m.k[m.k[2+v+k(i)]] & atom; ln > 0 {
			m.k[2+m.k[3]] = cat(m.k[2+m.k[3]], c2s(inc(m.k[2+v+k(i)])))
			m.k[3+m.k[3]] = cat(m.k[3+m.k[3]], enl(mkc(code+i)))
		}
	}
	dec(v)
}
func mkk(s s, k s) { dec(asn(mks(s), mkb([]c(k)), inc(null))) } // store k implementation
func kx(s, x k) (r k) { // exec monadic k implementation
	f := lup(inc(s))
	if m.k[f]>>28 == C { // cache on first call
		f = asn(inc(s), evl(prs(f)), inc(null))
	}
	dec(s)
	if m.k[f]>>28 != N+1 {
		panic("type")
	}
	return cal(f, enl(x))
}
func kxy(s, x, y k) (r k) { // execute dyadic k implementation
	f := lup(inc(s))
	if m.k[f]>>28 == C {
		f = asn(inc(s), evl(prs(f)), inc(null))
	}
	dec(s)
	if m.k[f]>>28 != N+2 {
		panic("type")
	}
	return cal(f, l2(x, y))
}
func msl() { // update slice header after increasing m.f
	f := *(*slice)(unsafe.Pointer(&m.f))
	i := *(*slice)(unsafe.Pointer(&m.k))
	i.l, i.c, i.p = f.l*2, f.c*2, f.p
	m.k = *(*[]k)(unsafe.Pointer(&i))
	b := *(*slice)(unsafe.Pointer(&m.c))
	b.l, b.c, b.p = f.l*8, f.c*8, f.p
	m.c = *(*[]c)(unsafe.Pointer(&b))
	zz := *(*slice)(unsafe.Pointer(&m.z))
	zz.l, zz.c, zz.p = f.l/2, f.l/2, f.p
	m.z = *(*[]z)(unsafe.Pointer(&zz))
}
func swap(ws []f) {
	m.f, ws = ws, m.f
	msl()
}
func save(ws []f) { copy(ws, m.f) }
func bk(t, n k) k {
	sz := k(0)
	if t < N {
		sz = lns[t]
	}
	if n != atom {
		sz *= n
	}
	if sz > 1<<31 {
		panic("size")
	}
	return buk(sz + 8) // complex values have an additional 8 byte padding after the header (does not change bucket type)
}
func mk(t, n k) k { // make type t of len n (-1:atom)
	bt := bk(t, n)
	fb, a := k(0), k(0)
	for i := bt; i < 32; i++ { // find next free bucket >= bt
		for k(i) >= m.k[2] {
			s := m.k[2]
			m.k[s] = k(len(m.c)) >> 2
			grw(1 << k(i)) // TODO: run a gc cycle (merge blocks) before growing?
			msl()
			if 1<<(1+s) != len(m.c) {
				panic("grow")
			}
			m.k[2] = s + 1
			m.k[1<<(s-2)] = s // bucket type of new upper half
		}
		if m.k[i] != 0 {
			fb, a = i, m.k[i]
			break
		}
	}
	m.k[fb] = m.k[1+a] // occupy
	if p := m.k[fb]; p > 0 && p < 40 {
		panic("illegal free pointer")
	}
	for i := fb - 1; i >= bt; i-- { // split large buckets
		u := a + 1<<(i-2) // free upper half
		m.k[1+u] = m.k[i]
		m.k[i] = u
		m.k[u] = i
	}
	m.k[a] = n | t<<28 // ok for atoms
	m.k[a+1] = 1       // refcount
	return a
}

func mki(i k) (r k) { r = mk(I, atom); m.k[2+r] = i; return r }
func mkc(c c) (r k) { r = mk(C, atom); m.c[8+r<<2] = c; return r }
func mkb(b []c) (r k) {
	r, n := mk(C, k(len(b))), k(len(b))
	rp := 8 + r<<2
	copy(m.c[rp:rp+n], b)
	return r
}
func mks(s s) (r k)                  { return c2s(mkb([]c(s))) }
func typ(x k) (k, k)                 { return m.k[x] >> 28, m.k[x] & atom } // type and length at addr
func typs(x, y k) (xt, yt, xn, yn k) { xt, xn = typ(x); yt, yn = typ(y); return }
func mtyp(x, n k) (t, cols k, eq bool) { // matrix type
	eq = true
	for i := k(0); i < n; i++ {
		tt, nn := typ(m.k[2+i+x])
		if tt > Z {
			panic("type")
		}
		if i == 0 {
			t, cols = tt, nn
		} else if nn != cols {
			panic("size")
		} else if tt != t {
			eq = false
		}
	}
	return t, cols, eq
}
func inc(x k) k {
	t, n := typ(x)
	switch {
	case t == L:
		if n == atom {
			panic("type")
		}
		for i := k(0); i < n; i++ {
			inc(m.k[2+x+i])
		}
	case t == A || (t == N && n == 2):
		inc(m.k[2+x])
		inc(m.k[3+x])
	case t > N:
		if n == atom && m.k[2+x] > 255 { // derived
			inc(m.k[3+x])
		} else if n == 0 { // lambda
			inc(m.k[2+x])
			inc(m.k[3+x])
		} else if n != atom { // projection, composition
			if n == 1 || n == 3 { // lambda-projection, composition
				inc(m.k[2+x])
			}
			inc(m.k[3+x])
		}
	}
	m.k[1+x]++
	return x
}
func dex(x, r k) k     { dec(x); return r }
func decr(x, y, r k) k { dec(x); dec(y); return r }
func dec(x k) {
	if m.k[x]>>28 == 0 || m.k[1+x] == 0 {
		panic("unref " + hxk(x))
	}
	t, n := typ(x)
	switch {
	case t == L:
		if n == atom {
			panic("type")
		}
		for i := k(0); i < n; i++ {
			dec(m.k[2+x+i])
		}
	case t == A || (t == N && n == 2):
		dec(m.k[2+x])
		dec(m.k[3+x])
	case t > N:
		if n == atom && m.k[2+x] > 255 { // derived
			dec(m.k[3+x])
		} else if n == 0 { // lambda
			dec(m.k[2+x])
			dec(m.k[3+x])
		} else if n != atom { // n: 1, 2 (projection), 3(composition)
			if n == 1 || n == 3 { // lambda-projection, or composition
				dec(m.k[2+x])
			}
			dec(m.k[3+x])
		}
	}
	m.k[1+x]--
	if m.k[1+x] == 0 {
		free(x)
	}
}
func free(x k) {
	t, n := typ(x)
	bt := bk(t, n)
	m.k[x] = bt
	m.k[x+1] = m.k[bt]
	m.k[bt] = x
}
func srk(x, t, n, nn k) (r k) { // shrink bucket
	if m.k[x]>>28 != t || t == L {
		panic("type")
	}
	if bk(t, nn) < bk(t, n) { // alloc not split: prevent small object accumulation
		r = mk(t, nn)
		ln := k(0)
		if t < N {
			ln = nn * lns[t]
		}
		if t == Z {
			ln += 8
		}
		rc, xc := 8+r<<2, 8+x<<2
		copy(m.c[rc:rc+ln], m.c[xc:xc+ln])
		return dex(x, r)
	}
	m.k[x] = t<<28 | nn
	return x
}
func to(x, rt k) (r k) { // numeric conversions for types CIFZ
	if rt == 0 || rt >= L {
		return x
	}
	t, n := typ(x)
	if rt == t {
		return x
	} else if t == L {
		r = mk(L, n)
		for i := k(0); i < n; i++ {
			m.k[2+i+r] = to(inc(m.k[2+x+i]), rt)
		}
		return dex(x, uf(r))
	}
	var g func(k, k)
	if t == S && rt == I { // for symbol conversion to bool
		g = func(r, x k) {
			if m.k[x] == 0 {
				m.k[r] = 0
			} else {
				m.k[r] = 1
			}
		}
	} else {
		g = tox[4*(t-1)+rt-1]
	}
	r = mk(rt, n)
	n = atm1(n)
	xp, rp := ptr(x, t), ptr(r, rt)
	for i := k(0); i < k(n); i++ {
		g(rp+i, xp+i)
	}
	return dex(x, r)
}
func nm(x, rt k, fx []f1) (r k) { // numeric monad
	t, n := typ(x)
	min := C
	if fx[C] == nil {
		min = I
	}
	if fx[I] == nil {
		min = F
	} // TODO: Z only for ff == nil ?
	if min > t { // uptype x
		x, t = to(x, min), min
	}
	if t == Z && fx[Z] == nil { // e.g. real functions
		x, t = to(x, F), F
	}
	if m.k[1+x] == 1 && t < L {
		r = inc(x)
	} else {
		r = mk(t, n)
	}
	n = atm1(n)
	switch t {
	case L:
		for j := k(0); j < k(n); j++ {
			m.k[r+2+j] = nm(inc(m.k[j+2+x]), rt, fx)
		}
	case A:
		if r != x {
			m.k[2+r] = inc(m.k[2+x])
		}
		m.k[3+r] = nm(m.k[3+x], rt, fx)
	case C, I, F, Z:
		rp, xp, f := ptr(r, t), ptr(x, t), fx[t]
		for i := k(0); i < k(n); i++ {
			f(rp+i, xp+i)
		}
	default:
		panic("type")
	}
	if rt != 0 && t > rt && t < L { // only down-type
		r = to(r, rt)
	}
	return dex(x, r)
}
func ntyps(xt, yt k, fx []f2, fc []fc) (t k) {
	if yt > xt {
		xt = yt
	}
	if xt == C && fc == nil && fx[C] == nil {
		xt = I
	}
	if (xt == I && fc == nil && fx[I] == nil) || (xt == Z && fc == nil && fx[Z] == nil) {
		xt = F
	}
	return xt
}
func nd(x, y, rt k, fx []f2, fc []fc) (r k) { // numeric dyad
	xt, yt, xn, yn := typs(x, y)
	if xt == A && yt == A {
		if match(m.k[2+x], m.k[2+y]) == false {
			panic("nyi") // (`a`b!..)+`b`c!..
		}
		r = mk(A, xn)
		m.k[2+r] = inc(m.k[2+x])
		m.k[3+r] = nd(inc(m.k[3+x]), inc(m.k[3+y]), rt, fx, fc)
		return decr(x, y, r)
	} else if xt == A {
		r = mk(A, xn)
		m.k[2+r] = inc(m.k[2+x])
		m.k[3+r] = nd(inc(m.k[3+x]), y, rt, fx, fc)
		return dex(x, r)
	} else if yt == A {
		r = mk(A, yn)
		m.k[2+r] = inc(m.k[2+y])
		m.k[3+r] = nd(x, inc(m.k[3+y]), rt, fx, fc)
		return dex(y, r)
	}
	n, sc := xn, k(0)
	if xn == atom {
		n, sc = yn, 1
	} else if yn == atom {
		n, sc = xn, 2
	} else if xn != yn {
		panic("size")
	}
	if xt == L || yt == L {
		r = mk(L, n)
		a := mk(I, atom)
		for i := k(0); i < n; i++ {
			m.k[2+a] = i
			m.k[2+r+i] = nd(atx(inc(x), inc(a)), atx(inc(y), inc(a)), rt, fx, fc)
		}
		dec(a)
		return decr(x, y, uf(r))
	} else if xt > S || yt > S {
		panic("type")
	}
	t := ntyps(xt, yt, fx, fc)
	if xt != t {
		x, xt = to(x, t), t
	}
	if yt != t {
		y, yt = to(y, t), t
	}
	if fc == nil {
		if m.k[1+x] == 1 && xn != atom {
			r = inc(x)
		} else if m.k[1+y] == 1 && yn != atom {
			r = inc(y)
		} else {
			r = mk(t, n)
		}
	} else if t < L {
		r = mk(I, n)
	} else {
		r = mk(L, n)
	}
	n = atm1(n)
	if fc == nil {
		if xt == S {
			panic("type")
		}
		ns(ptr(r, t), ptr(x, t), ptr(y, t), t, n, n, sc, fx[t])
	} else {
		g := fc[t]
		f := func(rp, xp, yp k) {
			if x := g(xp, yp); x {
				m.k[rp] = 1
			} else {
				m.k[rp] = 0
			}
		}
		ns(ptr(r, I), ptr(x, t), ptr(y, t), t, n, n, sc, f)
	}
	if rt != 0 && t > rt {
		r = to(r, rt)
	}
	return decr(x, y, r)
}
func ns(rp, xp, yp, t, xn, yn, c k, f f2) {
	switch c {
	case 0: // v f v
		for i := k(0); i < xn; i++ {
			f(rp+i, xp+i, yp+i)
		}
	case 1: // a f v
		for i := k(0); i < yn; i++ {
			f(rp+i, xp, yp+i)
		}
	case 2: // v f a
		for i := k(0); i < xn; i++ {
			f(rp+i, xp+i, yp)
		}
	default:
		panic("assert")
	}
}
func idx(x, t k) i { // int from a numeric scalar (trunc, ignore imag)
	switch t {
	case C:
		return i(m.c[8+x<<2])
	case I:
		return i(m.k[2+x])
	case F:
		return i(m.f[1+x>>1])
	case Z:
		return i(m.f[2+x>>1])
	}
	panic("type")
}
func explode(x k) (r k) { // explode an array (or atom) to a list of atoms
	t, n := typ(x)
	if t == L {
		return x
	} else if t > L {
		panic("type")
	}
	n = atm1(n)
	cp, xp := cpx[t], ptr(x, t)
	r = mk(L, n)
	for i := k(0); i < n; i++ {
		rk := mk(t, atom)
		rp := ptr(rk, t)
		cp(rp, xp+i)
		m.k[2+r+i] = rk
	}
	return dex(x, r)
}
func uf(x k) (r k) { // unify lists if possible
	xt, xn := typ(x)
	if xt != L || xn == 0 {
		return x
	}
	if m.k[m.k[2+x]]>>28 == A { // unify list of dicts to table
		rk := m.k[2+m.k[2+x]]
		for i := k(1); i < xn; i++ {
			if ti, ni := typ(m.k[2+i+x]); ti != A || ni != atom {
				return x
			}
			if match(m.k[2+m.k[2+i+x]], rk) == false {
				return x
			}
		}
		r = mk(A, xn)
		m.k[2+r] = inc(rk)
		l := mk(L, xn)
		for i := k(0); i < xn; i++ {
			m.k[2+i+l] = inc(m.k[3+m.k[2+i+x]])
		}
		m.k[3+r] = flp(l)
		return dex(x, r)
	}
	ut := k(0)
	for j := k(0); j < xn; j++ {
		t, n := typ(m.k[2+x+j])
		switch {
		case t >= L || n != atom:
			return x
		case j == 0:
			ut = t
		case t != ut:
			return x
		}
	}
	r = mk(ut, xn)
	cp, rp := cpx[ut], ptr(r, ut)
	for i := k(0); i < xn; i++ {
		cp(rp+i, ptr(m.k[2+x+i], ut))
	}
	return dex(x, r)
}
func idn(x k) (r k) { return x } // :x
func flp(x k) (r k) { // +x
	t, n := typ(x)
	if t > A {
		panic("type")
	} else if t == A {
		if n == atom {
			ln := k(0)
			v := m.k[3+x]
			for i := k(0); i < m.k[v]&atom; i++ {
				vk := atx(inc(v), mki(i))
				vn := m.k[vk] & atom
				if i == 0 {
					ln = vn
					if vn == atom {
						panic("class")
					}
				} else if vn != ln {
					panic("size") // rows have different lengths
				}
				dec(vk)
			}
			r = mk(A, ln)
			m.k[2+r] = inc(m.k[2+x])
			m.k[3+r] = inc(m.k[3+x])
			return dex(x, r)
		} else if m.k[1+x] == 1 {
			m.k[x] = A<<28 | atom
			return x
		} else {
			r = mk(A, atom)
			m.k[2+r] = inc(m.k[2+x])
			m.k[3+r] = inc(m.k[3+x])
			return dex(x, r)
		}
	} else if t < L {
		return x
	}
	return kx(mks(".flp"), x)
}
func neg(x k) k { // -x
	return nm(x, 0, []f1{nil, func(r, x k) { m.c[r] = -m.c[x] }, func(r, x k) { m.k[r] = k(-i(m.k[x])) }, func(r, x k) { m.f[r] = -m.f[x] }, func(r, x k) { m.f[r<<1] = -m.f[x<<1]; m.f[1+r<<1] = -m.f[1+x<<1] }})
}
func fst(x k) (r k) { // *x
	t, n := typ(x)
	if t == A {
		inc(m.k[3+x])
		r = fst(m.k[3+x])
		return dex(x, r)
	}
	if n == atom {
		return x
	} else if n == 0 {
		if t < L {
			return dex(x, inc(nan[t]))
		} else if t == L {
			return dex(x, mk(C, 0))
		} else {
			return dex(x, inc(null))
		}
	}
	if t == L {
		r = m.k[2+x]
		inc(r)
		return dex(x, r)
	}
	r = mk(t, atom)
	switch t {
	case C:
		m.c[8+r<<2] = m.c[8+x<<2]
	case I:
		m.k[2+r] = m.k[2+x]
	case F, S:
		m.f[1+r>>1] = m.f[1+x>>1]
	case Z:
		m.z[1+r>>2] = m.z[1+x>>2]
	default:
		panic("nyi")
	}
	return dex(x, r)
}
func inv(x k) k { // %x
	return nm(x, 0, []f1{nil, nil, nil, func(r, x k) { m.f[r] = 1 / m.f[x] }, func(r, x k) { m.z[r] = 1 / m.z[x] }})
}
func wer(x k) (r k) { // &x
	t, n := typ(x)
	if t == A {
		return dex(x, atx(inc(m.k[2+x]), wer(inc(m.k[3+x]))))
	} else if t != I {
		panic("type")
	}
	n = atm1(n)
	nn := k(0)
	for j := k(0); j < n; j++ {
		if p := i(m.k[2+x+j]); p < 0 {
			panic("domain")
		} else {
			nn += k(p)
		}
	}
	r = mk(I, nn)
	jj := k(0)
	for j := k(0); j < n; j++ {
		for p := k(0); p < m.k[2+x+j]; p++ {
			m.k[2+r+jj] = j
			jj++
		}
	}
	return dex(x, r)
}
func rev(x k) (r k) { // |x
	t, n := typ(x)
	if n == atom || n < 2 {
		if t == A {
			r = mk(t, n)
			m.k[r+2] = rev(inc(m.k[x+2]))
			m.k[r+3] = rev(inc(m.k[x+3]))
			return dex(x, r)
		}
		return x
	}
	r = mk(t, n)
	if t < A {
		cp := cpx[t]
		if t == L {
			cp = cpI
		}
		xp, rp := ptr(x, t), ptr(r, t)
		for i := k(0); i < n; i++ {
			cp(rp+n-1-i, xp+i)
		}
	} else {
		panic("nyi")
	}
	if t == L {
		for i := k(0); i < n; i++ {
			inc(m.k[2+i+r])
		}
	}
	return dex(x, r)
}
func asc(x k) (r k) { // <x
	t, n := typ(x)
	if n == atom || t > L {
		panic("type")
	} else if t == A {
		return arc(x, n, asc)
	}
	r = til(mki(n))
	w := mk(I, n)
	mv(w, r)
	msrt(2+w, 2+r, k(0), n, ptr(x, t), gtx[t])
	return decr(x, w, r)
}
func msrt(x, r, a, b, p k, gt fc) { // merge sort
	if b-a < 2 {
		return
	}
	c := (a + b) / 2
	msrt(r, x, a, c, p, gt)
	msrt(r, x, c, b, p, gt)
	mrge(x, r, a, b, c, p, gt)
}
func mrge(x, r, a, b, c, p k, gt fc) {
	i, j := a, c
	for k := a; k < b; k++ {
		if i >= c || (j < b && gt(p+m.k[x+i], p+m.k[x+j])) {
			m.k[r+k] = m.k[x+j]
			j++
		} else {
			m.k[r+k] = m.k[x+i]
			i++
		}
	}
}
func isrt(x, r, n k, lt fc) { // insertion sort (can be removed)
	v := r
	for i := k(1); i < n; i++ {
		for j := k(i); j > 0 && lt(x+m.k[v+j], x+m.k[v+j-1]); j-- {
			swI(r+j, r+(j-1))
		}
	}
}
func dsc(x k) (r k) { return rev(asc(x)) } // >x
func grp(x k) (r k) { // =x {k!&:'(k:^?x)~/:\:x}
	t, n := typ(x)
	if n == atom {
		return eye(x)
	} else if t == A && n != atom { // =t
		g := grp(flp(inc(m.k[3+x])))
		r = key(flp(key(inc(m.k[2+x]), flp(inc(m.k[2+g])))), inc(m.k[3+g]))
		return decr(g, x, r)
	} else if t > L {
		panic("type")
	} else if n == 0 {
		return dex(x, key(inc(x), take(0, 0, inc(x))))
	}
	kk := srt(unq(inc(x)))
	kp, kn, xp := ptr(kk, t), m.k[kk]&atom, ptr(x, t)
	vv := tak(mki(kn), enl(mk(I, 0))) // (#^?x)#,!0
	if t == L {
		for i := k(0); i < n; i++ {
			for j := k(0); j < kn; j++ {
				if match(m.k[kp+j], m.k[xp+i]) {
					m.k[2+vv+j] = ucat(m.k[2+vv+j], mki(i), I, m.k[m.k[2+vv+j]]&atom, atom)
					break
				}
			}
		}
	} else {
		for i := k(0); i < n; i++ {
			ii := ibin(kp, kn, xp+i, gtx[t])
			m.k[2+vv+ii] = ucat(m.k[2+vv+ii], mki(i), I, m.k[m.k[2+vv+ii]]&atom, atom)
		}
	}
	return dex(x, key(kk, vv))
}
func til(x k) (r k) { // !x
	t, n := typ(x)
	if n != atom {
		return odo(x) // return kx(mks(".odo"), x)
	} else if t == A {
		r = inc(m.k[2+x])
		return dex(x, r)
	} else if t > Z {
		panic("type")
	}
	if nn := idx(x, t); nn < 0 {
		return eye(neg(x))
	} else if nn == 0 {
		return dex(x, mk(t, 0))
	} else {
		return dex(x, to(jota(k(nn)), t))
	}
}
func jota(n k) (r k) { // !n
	r = mk(I, n)
	for j := k(0); j < n; j++ {
		m.k[2+r+j] = j
	}
	return r
}
func odo(x k) (r k) { // !x
	n := m.k[x] & atom
	q := k(1) // */x
	for i := k(0); i < n; i++ {
		ki := m.k[2+i+x]
		if ki == 0 || int32(ki) < 0 {
			panic("domain")
		}
		q *= ki
	}
	r = mk(L, n)
	for i := k(0); i < n; i++ {
		m.k[2+r+i] = mk(I, q)
	}
	rep := k(1)
	for i := k(0); i < n; i++ {
		ri := m.k[1+r+n-i]
		p, kk := k(0), k(0)
		for j := k(0); j < q; j++ {
			m.k[2+ri+j] = p
			kk++
			if kk == rep {
				kk = 0
				p++
				if p == m.k[1+x+n-i] {
					p = 0
				}
			}
		}
		rep *= m.k[1+x+n-i]
	}
	return dex(x, r)
}
func eye(x k) (r k) { // !-n =n (ifz)
	t, n := typ(x)
	if n != atom {
		panic("type")
	}
	ln := idx(x, t)
	if ln < 0 {
		panic("value")
	}
	return dex(x, dia(take(k(ln), 0, to(mki(1), t))))
}
func dia(x k) (r k) { // diag x
	t, n := typ(x)
	if n == atom || t > L {
		panic("type")
	} else if t == L { // v:diag A
		r = mk(L, n)
		for i := k(0); i < n; i++ {
			m.k[2+i+r] = atx(inc(m.k[2+i+x]), mki(i))
		}
		return dex(x, uf(r))
	}
	r = mk(L, n) // A:diag v
	z, cp := to(mki(0), t), cpx[t]
	xp, zp := ptr(x, t), ptr(z, t)
	for i := k(0); i < n; i++ {
		rr := mk(t, n)
		rp := ptr(rr, t)
		for j := k(0); j < n; j++ {
			cp(rp+j, zp)
		}
		cp(rp+i, xp+i)
		m.k[2+i+r] = rr
	}
	return decr(x, z, r)
}
func not(x k) (r k) { // ~x
	t, n := typ(x)
	if n == 0 {
		return dex(x, mk(I, 0))
	} else if t < S {
		return to(eql(mki(0), x), I)
	} else if t == S {
		return eql(inc(nans), x)
	} else if t == L {
		return lrc(x, n, not)
	} else if t == A {
		return arc(x, n, not)
	} else if t == N {
		if n == 2 { // ~ :e
			return dex(x, mki(0))
		}
		return dex(x, mki(1))
	} else if t > N {
		return dex(x, mki(0))
	}
	panic("type")
}
func enl(x k) (r k) { // ,x (collaps uniform)
	t, n := typ(x)
	if t < L && n == atom {
		r = mk(t, 1)
		cp := cpx[t]
		src, dst := ptr(x, t), ptr(r, t)
		cp(dst, src)
		return dex(x, r)
	} else if t == A && n == atom {
		r = mk(A, 1)
		m.k[2+r] = inc(m.k[2+x])
		f := mk(N+1, atom)
		m.k[2+f] = 12 // ,:
		m.k[3+r] = ech(f, inc(m.k[3+x]))
		return dex(x, r)
	}
	r = mk(L, 1)
	m.k[2+r] = x
	return r
}
func enlist(x k) (r k) { // dont unify
	r = mk(L, 1)
	m.k[2+r] = x
	return r
}
func srt(x k) (r k) { return atx(x, asc(inc(x))) } // ^x  TODO: replace with a sort implementation
func cnt(x k) (r k) { // #x
	t, n := typ(x)
	if t == A {
		_, n = typ(m.k[x+2])
	}
	n = atm1(n)
	return dex(x, mki(k(i(n))))
}
func flr(x k) (r k) { // _x
	if t, n := typ(x); t == C { // _c (tolower)
		r = mk(C, n)
		xp, rp := ptr(x, C), ptr(r, C)
		for i := k(0); i < atm1(n); i++ {
			if c := m.c[xp+i]; c >= 'A' && c <= 'Z' {
				m.c[rp+i] = c + 32
			} else {
				m.c[rp+i] = c
			}
		}
		return dex(x, r)
	}
	return nm(x, I, []f1{nil, func(r, x k) { m.c[r] = m.c[x] }, func(r, x k) { m.k[r] = m.k[x] }, func(r, x k) {
		if isnan(m.f[x]) { // go issue 35034 on wasm
			m.f[r] = m.f[x]
			return
		}
		y := f(i(m.f[x]))
		if m.f[x] < y {
			y -= 1.0
		}
		m.f[r] = y
	}, nil}) // TODO: k7 does not convert c to i
}
func str(x k) (r k) { // $x
	t, n := typ(x)
	if t == C {
		return x
	} else if t == S && n == atom {
		p := m.k[2+x]
		if p == 0 {
			return dex(x, mk(C, 0))
		} else if p < 256 {
			r = mk(C, 1)
			m.c[8+r<<2] = c(p)
			return dex(x, r)
		}
		p -= 256
		return dex(x, inc(m.k[2+p+m.k[stab]]))
	}
	if t < L {
		st, xp := stx[t], ptr(x, t)
		if n == atom {
			r = mk(C, 56)
			r = srk(r, C, 56, st(8+r<<2, xp))
		} else {
			r = mk(L, n)
			for i := k(0); i < n; i++ {
				y := mk(C, 56)
				m.k[2+r+i] = srk(y, C, 56, st(8+y<<2, xp+i))
			}
		}
	} else {
		switch {
		case t == L:
			return lrc(x, n, str)
		case t == A:
			return arc(x, n, str)
		case t == N:
			if n == 2 {
				r = inc(m.k[3+x]) // :expr
			} else {
				r = mk(C, 0)
			}
		case t > N:
			f := m.k[2+x]
			if n == 0 || n == 1 { // 0(lambda), 1(lambda projection)
				r = inc(m.k[2+x]) // `C
			} else if n == 3 { // composition
				r = cat(str(inc(m.k[2+x])), str(inc(m.k[3+x])))
			} else if (f >= 39 && f < dyad) || (f >= 39+dyad && f < 2*dyad) { // built-ins
				r = str(atx(inc(m.k[2+m.k[3]]), fst(wer(eql(mki(f), inc(m.k[3+m.k[3]]))))))
			} else if f < 20 || (f >= 30 && f <= 33) { // monad +: /:
				r = mkb([]c{m.c[136+m.k[2+x]], ':'})
			} else if f >= 20 && f < 30 { // monadic ioverb
				r = mkb([]c{'0' + c(f-20), ':', ':'})
			} else if f >= 20+dyad && f < 30+dyad { // dyadic ioverb 3:
				r = mkb([]c{'0' + c(f-20-dyad), ':'})
			} else if f >= dyad && f < 33+dyad { // dyad * /
				r = mkc(m.c[136+m.k[2+x]-dyad])
				m.k[r] = C<<28 | 1
			} else if f > 256 && t == N+1 { // derived verb (see func drv)
				if op := f >> 8; op >= 33 && op <= 38 {
					if op < 36 { // 33-35 ' / \
						op += dyad - 3
					} else { // 36-38 ': /: \:
						op -= 6
					}
					opv := mk(N+1, atom)
					m.k[2+opv] = op
					r = cat(str(inc(m.k[3+x])), str(opv))
				} else {
					panic("type") // unknown verb
				}
			} else {
				panic("assert")
			}
			if n == 1 || n == 2 { // projection
				a := m.k[3+x]
				if n == 2 && f < 2*dyad && match(m.k[3+a], null) {
					r = cat(kst(inc(m.k[2+a])), r) // short form: 2+
				} else {
					a = kst(inc(a))   // arg list
					m.c[8+a<<2] = '[' // convert () to []
					m.c[7+(m.k[a]&atom)+a<<2] = ']'
					r = cat(r, a)
				}
			}
		default:
			panic("nyi")
		}
	}
	return dex(x, r)
}
func unq(x k) (r k) { // ?x
	t, n := typ(x)
	if n == atom {
		panic("nyi") // overloads, random numbers?
	} else if t == A { // what does ?d do?
		panic("type")
	} else if n < 2 {
		return x
	}
	r = mk(t, n)
	eq, cp, src, dst, nn := eqx[t], cpx[t], ptr(x, t), ptr(r, t), k(0)
	for i := k(0); i < n; i++ { // quadratic, should be improved
		u := true
		srci := src + i
		for j := k(0); j < nn; j++ {
			if eq(srci, dst+j) {
				u = false
				break
			}
		}
		if u {
			cp(dst+nn, srci)
			nn++
		}
	}
	if t != L {
		return dex(x, srk(r, t, n, nn))
	}
	for i := nn; i < n; i++ {
		m.k[dst+i] = inc(null)
	}
	return dex(x, take(nn, 0, r))
}
func tip(x k) (r k) { // @x
	t, n := typ(x)
	if match(x, null) {
		return dex(x, c2s(mkb(nil)))
	} else if t >= N {
		return dex(x, c2s(mkc(byte('0'+t-N))))
	} else if t == A && n != atom {
		return dex(x, c2s(mkc('A')))
	}
	return dex(x, c2s(mkc(m.c[169+t])))
}
func val(x k) (r k) { // . x
	switch m.k[x] >> 28 {
	case L, S:
		return evl(x)
	case C:
		return evl(prs(x))
	case A:
		return dex(x, inc(m.k[3+x]))
	case N:
		if m.k[x]&atom == 2 { // :expr
			return dex(x, inc(m.k[2+x]))
		}
	}
	panic("type")
}
func evl(x k) (r k) {
	t, n := typ(x)
	if t != L {
		if t == S && n == 1 {
			return fst(x)
		} else if t == S && n == atom {
			return lup(x)
		}
		return x
	} else if n == 1 {
		return fst(x)
	}
	sp := k(0)
	if m.k[1+x] > 0xFFFF {
		sp = m.k[1+x] >> 16
		m.k[srcp] = sp
		m.k[1+x] &= 0xFFFF
	}
	R, P := func() k {
		m.k[srcp] = sp
		return 0
	}, func(s string) {
		m.k[srcp] = sp
		panic(s)
	}
	if n == 0 {
		return x
	}
	v := m.k[2+x]
	vt, vn := typ(v)
	if vt == S {
		if n == 2 && vn == 1 && m.k[2+v] == 0 { // (,`;..) <- `@y
			return R() + dex(x, ser(evl(inc(m.k[3+x]))))
		}
		if n == 1 { // ,`a`b → `a`b
			return dex(x, inc(v))
		}
		if m.k[2+v] == 0 { // (`;…) → ex;ex…
			for i := k(1); i < n; i++ {
				if i > 1 {
					dec(r)
				}
				r = evl(inc(m.k[2+i+x]))
			}
			return R() + dex(x, r)
		}
	}
	if match(v, null) { // (;…) → list
		r = mk(L, n-1)
		if n > 1 {
			for i := int(n - 2); i >= 0; i-- {
				m.k[2+r+k(i)] = evl(inc(m.k[3+x+k(i)]))
			}
		}
		return R() + dex(x, uf(r))
	} else {
		inc(v)
		iev := false
		if vt == S && vn == atom {
			v = evl(v)
			vt, vn = typ(v)
			iev = true
		}
		if n == 1 && vt > N { // e.g. (-)
			return dex(x, v)
		}
		af := m.k[2+v]
		if vt == N+1 && vn == atom && n == 3 { // : or :: or *: (modified assignmnt)
			if n != 3 {
				P("args")
			}
			f := k(0)
			if af != 0 { // not ::, e.g. *:
				f = mk(N+2, atom)
				m.k[2+f] = m.k[2+v] + dyad
			} else {
				f = inc(null)
			}
			name, val := inc(m.k[3+x]), evl(inc(m.k[4+x]))
			if nt, nn := typ(name); nt == L && nn > 1 {
				if match(m.k[2+name], null) { // (;`a;`b) vector assignment
					name = drp(mki(1), name)
				} else if nn > 1 { // (`a;i) amd | (`a;i;j..) dmd
					idx := drop(1, inc(name)) // inc(m.k[3+name])
					name = fst(name)
					if m.k[idx]>>28 != L {
						P("assert")
					}
					idx = evl(cat(inc(null), idx))
					name, idx = dxn(name, idx)
					if nn := m.k[idx] & atom; nn == 1 {
						dec(amd(name, fst(idx), f, val))
					} else {
						dec(dmd(name, idx, f, val))
					}
					return R() + decr(v, x, inc(null))
				}
			}
			return decr(v, x, asn(name, val, f))
		} else if n > 3 && vt > N && vn == atom && m.k[2+v] == 16+dyad { // $[...] delays evaluation
			return R() + dex(v, swc(drop(1, x)))
		}
		r = mk(L, n-1)
		for i := int(n - 2); i >= 0; i-- {
			m.k[2+r+k(i)] = evl(inc(m.k[3+x+k(i)]))
		}
		dec(x)
		if iev == false {
			v = evl(v)
		}
		vt, vn := typ(v)
		if n > 3 && vt > N && vn == atom {
			g := amd
			switch code := m.k[2+v]; code { // triadics..
			case dyad + 14: // #
				g = sel
			case dyad + 15: // _
				g = udt
			case dyad + 17: // ?
				g = ins
			case dyad + 18: // @
				g = amd
			case dyad + 19: // .
				g = dmd
			default:
				P("args")
			}
			if n == 4 {
				x, a, y := inc(m.k[2+r]), inc(m.k[3+r]), inc(m.k[4+r])
				dec(v) // dec early, allow inplace
				dec(r)
				if m.k[y]>>28 >= N {
					return R() + g(x, a, y, inc(null))
				} else {
					return R() + g(x, a, inc(null), y)
				}
			} else if n == 5 {
				x, a, f, y := inc(m.k[2+r]), inc(m.k[3+r]), inc(m.k[4+r]), inc(m.k[5+r])
				dec(v)
				dec(r)
				return R() + g(x, a, f, y)
			} else {
				P("args")
			}
		} else if n == 3 && vt == N+2 && m.k[2+v] == 19+dyad && m.k[m.k[3+r]]>>28 > N { // composition
			dec(v)
			v = mk(m.k[m.k[3+r]]>>28, 3)
			m.k[2+v] = inc(m.k[2+r])
			m.k[3+v] = inc(m.k[3+r])
			return R() + dex(r, v)
		} else if vt > N && !(vn == atom && vt == N+1 && n-1 == 2 && m.k[2+v] > 255) { // allow dyadic derived
			if n-1 > vt-N {
				P("args") // too many arguments
			}
			for i := n - 1; i < vt-N; i++ { // fill args, e.g. 2+
				r = lcat(r, inc(null))
			}
			if vt > N+1 { // no projection for monads, allow N argument
				for i := k(0); i < m.k[r]&atom; i++ {
					if match(m.k[2+i+r], null) {
						return R() + prj(v, r)
					}
				}
			}
		} else if vt < N && n == 2 { // @
			return atx(v, fst(r))
		} else if vt < N && n > 2 {
			for i := k(0); i < n-2; i++ { // last index may be a vector
				if match(m.k[2+i+r], null) || m.k[m.k[2+i+r]]&atom != atom {
					return R() + atm(v, r)
				}
			}
		}
		return R() + cal(v, r)
	}
	return x
}
func ano(p, e k) (r k) { // annotate source line with error position
	r = cat(lup(mks(".f")), mkc(':')) // TODO: .ano(k)
	re := cat(cat(inc(r), inc(e)), mkc('\n'))
	if p == 0 {
		return decr(r, e, re)
	}
	s := lupo(mks(".c"))
	if s == 0 {
		return decr(r, e, re)
	}
	t, n := typ(s)
	if t != C || p >= n {
		return decr(r, e, re)
	}
	dec(re)
	sp, a, b, l := ptr(s, C), k(0), n, k(1)
	for i := k(0); i < p; i++ {
		if m.c[sp+i] == '\n' {
			a = i + 1
			l++
		}
	}
	for i := p; i < n; i++ {
		if m.c[sp+i] == '\n' {
			b = i
			break
		}
	}
	r = cat(r, str(mki(l))) // file:line:char:error\nsource\n   ^
	r = cat(r, mkc(':'))
	r = cat(r, str(mki(2+p-a)))
	r = cat(r, mkc(':'))
	r = cat(r, e)
	r = cat(r, mkc('\n'))
	r = cat(r, take(b-a, 0, drop(i(a), s)))
	r = cat(r, mkc('\n'))
	if p-a != 0 {
		r = cat(r, take(1+p-a, 0, mkc(' ')))
	}
	r = cat(r, mkc('^'))
	r = cat(r, mkc('\n'))
	return r
}
func swc(x k) (r k) { // $[...]
	n := m.k[x] & atom
	for i := k(0); i < n-1; i += 2 {
		if !is0(evl(inc(m.k[2+i+x]))) {
			return dex(x, evl(inc(m.k[3+i+x])))
		}
	}
	if n%2 == 0 {
		return dex(x, inc(null))
	}
	return dex(x, evl(inc(m.k[1+n+x])))
}
func is0(x k) bool { // for swc
	n := m.k[x] & atom
	if n == atom {
		x = not(x)
		dec(x)
		return m.k[2+x] == k(1)
	}
	dec(x)
	if n == 0 {
		return true
	}
	return false // e.g. $[,0;1;2] is 2
}
func prj(x, y k) (r k) { // convert x to a projection
	t := m.k[x] >> 28
	r = mk(t, 2)
	ln := k(2)
	if f := m.k[2+x]; f < 256 {
		m.k[2+r] = f // #1: function code if < 256
		dec(x)
	} else {
		m.k[2+r] = x // #1: pointer to lambda function if code >= 256
		ln = 1
	}
	m.k[3+r] = y // #2: argument list with holes
	n := k(0)
	for i := k(0); i < m.k[y]&atom; i++ {
		if match(m.k[2+y+i], null) {
			n++
		}
	}
	m.k[r] = k(N+n)<<28 | ln
	return r
}
func kst(x k) (r k) { // `k@x
	t, n := typ(x)
	atm := n == atom
	if atm {
		n = 1
	}
	if n == 0 && t < A {
		r = mk(C, 0)
		rc, rn := 8+r<<2, k(0)
		switch t { // these could also be in the k-tree
		case C:
			rn = putb(rc, rn, []c(`""`))
		case I:
			rn = putb(rc, rn, []c("!0"))
		case F:
			rn = putb(rc, rn, []c("!0f"))
		case Z:
			rn = putb(rc, rn, []c("!0a"))
		case S:
			rn = putb(rc, rn, []c("0#`"))
		case L:
			rn = putb(rc, rn, []c("()"))
		case A:
			rn = putb(rc, rn, []c("+()!()"))
		default:
			panic("nyi")
		}
		m.k[r] = C<<28 | rn
		return dex(x, r)
	}
	switch t {
	case C: // ,"a" "a" "ab" "a\nb" ,0x01 0x010203
		r = mk(C, 2+2*n) // for both "a\nb" or 0x01234 or ,"\n"(short enough)
		rc, rn, xc := 8+r<<2, k(0), 8+x<<2
		if n == 1 && !atm {
			rn = putc(rc, rn, ',')
		}
		rn = putc(rc, rn, '"')
		hex := false
		for i := k(0); i < n; i++ {
			c := m.c[xc+i]
			if c < 32 || c > 126 || c == '"' || c == '\\' {
				if c, o := qt(c); o {
					rn = putc(rc, rn, '\\')
					rn = putc(rc, rn, c)
				} else {
					hex = true
					break
				}
			} else {
				rn = putc(rc, rn, c)
			}
		}
		rn = putc(rc, rn, '"')
		if hex {
			rn = 0
			if n == 1 && !atm {
				rn = 1
			}
			rn = putc(rc, rn, '0')
			rn = putc(rc, rn, 'x')
			for i := k(0); i < n; i++ {
				c1, c2 := hxb(m.c[xc+i])
				rn = putc(rc, rn, c1)
				rn = putc(rc, rn, c2)
			}
		}
		r = srk(r, C, 2+2*n, rn)
	case I, F, Z:
		r = mk(C, 0)
		if n == 1 && !atm {
			m.c[8+r<<2] = ','
			m.k[r] = C<<28 | 1
		}
		rr := mk(C, 56)
		st, xp, rrc := stx[t], ptr(x, t), 8+rr<<2
		sp := mkb([]c{' '})
		for i := k(0); i < n; i++ {
			rn := st(rrc, xp+i)
			m.k[rr] = C<<28 | rn
			r = cat(r, inc(rr))
			if i < n-1 {
				r = cat(r, inc(sp))
			}
		}
		dec(sp)
		m.k[rr] = C<<28 | 56
		dec(rr)
		if t == F || t == Z {
			_, n = typ(r)
			rc, dot := 8+r<<2, false
			for i := k(0); i < n; i++ {
				if c := m.c[rc+i]; t == F && (c == '.' || c == 'n' || c == 'e' || c == 'w') {
					dot = true
					break
				} else if t == Z && (c == 'i' || c == 'a') {
					dot = true
					break
				}
			}
			if !dot {
				if t == F {
					r = cat(r, mkb([]c{'f'}))
				} else {
					r = cat(r, mkb([]c{'a'}))
				}
			}
		}
	case S:
		if atm || n == 1 {
			_, n := sc(2 + x)
			rr := mk(C, n)
			sn, rrc, rn, q := stS(8+rr<<2, ptr(x, S)), 8+rr<<2, k(1), false
			for i := k(0); i < sn; i++ {
				c := m.c[rrc+i]
				if !(cr0Z(c) || c == '.') {
					q = true
				}
				if _, o := qt(c); o {
					rn++
				}
				rn++
			}
			if q {
				rn += 2
			}
			if !atm {
				rn++
			}
			r = mk(C, rn)
			rc, rn := 8+r<<2, k(0)
			if !atm {
				rn = putc(rc, rn, ',')
			}
			rn = putc(rc, rn, '`')
			if q {
				rn = putc(rc, rn, '"')
			}
			for i := k(0); i < sn; i++ {
				c, o := qt(m.c[rrc+i])
				if o {
					rn = putc(rc, rn, '\\')
				}
				rn = putc(rc, rn, c)
			}
			if q {
				rn = putc(rc, rn, '"')
			}
			dec(rr)
		} else {
			r = mk(C, 0)
			ix := mk(I, atom)
			for i := k(0); i < n; i++ {
				m.k[2+ix] = i
				r = cat(r, kst(atx(inc(x), inc(ix))))
			}
			dec(ix)
		}
	case L:
		r = mk(C, 1)
		rc := 8 + r<<2
		m.c[rc] = '('
		if n == 1 {
			m.c[rc] = ','
		}
		y := mkb([]c{';'})
		for i := k(0); i < n; i++ {
			r = cat(r, kst(inc(m.k[2+i+x])))
			if i < n-1 {
				r = cat(r, inc(y))
			}
		}
		if n != 1 {
			m.c[8+y<<2] = ')'
			r = cat(r, y)
		} else {
			dec(y)
		}
	case A:
		r = mk(C, 0)
		if !atm {
			r = cat(r, mkc('+'))
		}
		kv, vv := inc(m.k[2+x]), inc(m.k[3+x])
		kt, nk := typ(kv)
		if nk == 1 {
			kv, vv = fst(kv), fst(vv)
			kt, nk = typ(kv)
		}
		rr, encl := kst(kv), false
		if (kt <= L && nk == 1) || (kt == A) || (kt > A) || (nk == 0 && kt != C && kt != L) {
			encl = true
		}
		y := mk(C, 1)
		if encl {
			m.c[8+y<<2] = '('
			r = cat(r, inc(y))
			r = cat(r, rr)
			m.c[8+y<<2] = ')'
			r = cat(r, inc(y))
		} else {
			r = cat(r, rr)
		}
		m.c[8+y<<2] = '!'
		r = cat(r, y)
		r = cat(r, kst(vv))
	default:
		if t >= N {
			r = str(inc(x))
		} else {
			panic("nyi")
		}
	}
	return dex(x, r)
}
func mat(x k) (r k) { // `m@x (matrix display; should be implemented in k)
	t, n := typ(x)
	if t == L {
		r = mk(L, n)
		isc := true
		for i := k(0); i < n; i++ {
			xi := inc(m.k[2+x+i])
			xxt, xxn := typ(xi)
			if xxt != C || xxn == atom {
				isc = false
			}
			switch {
			case xxt < L && xxn == atom:
				m.k[2+r+i] = enl(str(xi))
			case xxt < L:
				m.k[2+r+i] = str(xi)
			default:
				m.k[2+r+i] = enl(kst(xi))
			}
		}
		if isc {
			return dex(r, x)
		}
	} else if t == A && n == atom {
		r = str(inc(m.k[2+x]))
		v := m.k[3+x]
		if t, n := typ(r); t != L {
			panic("type")
		} else {
			mx := cmc(r, n)
			for i := k(0); i < n; i++ {
				m.k[2+r+i] = cat(cat(pad(mx, m.k[2+r+i]), mkc(':')), kst(atx(inc(v), mki(i))))
			}
		}
		return dex(x, r)
	} else if t == A {
		nc := m.k[m.k[2+x]] & atom
		h := str(inc(m.k[2+x]))
		a := mk(L, 0)
		for i := k(0); i < nc; i++ {
			r = str(atx(inc(m.k[3+x]), mki(i)))
			r = cat(enl(inc(m.k[2+h+i])), r)
			mx := cmc(r, m.k[r]&atom)
			for j := k(0); j < m.k[r]&atom; j++ {
				m.k[2+r+j] = pad(mx, m.k[2+r+j])
			}
			a = lcat(a, r)
		}
		dec(h)
		a = flp(a)
		b := mk(L, nc)
		for i := k(0); i < nc; i++ {
			m.k[2+b+i] = take(m.k[m.k[2+i+m.k[2+a]]]&atom, 0, mkc('-'))
		}
		a = insert(a, b, 1)
		sep := mkc(' ')
		for i := k(0); i < m.k[a]&atom; i++ {
			m.k[2+a+i] = jon(inc(sep), m.k[2+a+i])
		}
		return decr(x, sep, a)
	} else {
		return enl(kst(x))
	}
	r = flp(r) // " "/:+n$+r (join flip pad flip)
	for i := k(0); i < m.k[r]&atom; i++ {
		ri := m.k[2+r+i]
		mx, nc := cmc(ri, m.k[ri]&atom), m.k[ri]&atom
		for j := k(0); j < nc; j++ {
			m.k[2+j+ri] = pad(mx, m.k[2+j+ri])
		}
	}
	r = flp(r)
	sep := mkc(' ')
	for i := k(0); i < m.k[r]&atom; i++ {
		m.k[2+r+i] = jon(inc(sep), m.k[2+r+i])
	}
	return decr(x, sep, r)
}
func ser(x k) (r k) { // `@ (k7 compat)
	t, n := typ(x)
	switch t {
	case C, I, F, Z:
		ln, o := lns[t]*atm1(n), k(0)
		if t == I {
			t = 7
		} else if t == F {
			t = 0x0e
		} else if t == Z {
			t, o = 9, 8 // not k7
		}
		if n == atom {
			r = mk(C, 1+ln)
			rp, xp := ptr(r, C), ptr(x, C)
			m.c[rp] = c(t)
			copy(m.c[1+rp:1+ln+rp], m.c[xp+o:xp+o+ln])
			return dex(x, r)
		}
		r = mk(C, 8+ln) // 8 byte header: type, length
		m.k[2+r] = t << 24
		m.k[3+r] = n
		rp, xp := ptr(2+r, C), ptr(x, C)
		copy(m.c[rp:rp+ln], m.c[xp+o:xp+o+ln])
		return dex(x, r)
	case S:
		if n == atom {
			r = mk(C, 1)
			m.c[8+r<<2] = 0x0f
		} else {
			r = mk(C, 8)
			m.k[2+r] = 0x0f << 24
			m.k[3+r] = n
		}
		for i := k(0); i < atm1(n); i++ {
			c, n := sc(2 + i + x)
			s := mk(C, n+1)
			sp := 8 + s<<2
			copy(m.c[sp:sp+n], m.c[c:c+n])
			m.c[sp+n] = 0
			r = cat(r, s)
		}
		return dex(x, r)
	case L:
		r = mk(C, 8)
		m.k[2+r] = 0
		m.k[3+r] = n
		if n == 0 {
			return dex(x, cat(r, ser(mk(C, 0)))) // prototype
		}
		for i := k(0); i < n; i++ {
			r = cat(r, ser(inc(m.k[2+x+i])))
		}
		return dex(x, r)
	case A:
		p := k(0)
		if n != atom {
			r = mk(C, 16)
			m.k[2+r] = 0x14 << 24
			m.k[3+r] = 1
			p = 2
		} else {
			r = mk(C, 8)
		}
		m.k[2+p+r] = 0x15 << 24
		m.k[3+p+r] = 2
		r = cat(r, ser(inc(m.k[2+x])))
		r = cat(r, ser(inc(m.k[3+x])))
		return dex(x, r)
	default:
		panic("nyi")
	}
}
func res(x k) (r k) { // `?x
	r, x = resn(x)
	return dex(x, r)
}
func resn(x k) (r, xe k) {
	t, n := typ(x)
	if t != C || n < 2 || (m.c[8+x<<2] == 0 && n < 8) {
		panic("type")
	}
	rt, rn, o := k(m.c[8+x<<2]), atom, k(0)
	if rt == 0 {
		rt, rn = m.k[2+x]>>24, m.k[3+x]
		x = drop(8, x)
	} else {
		x = drop(1, x)
	}
	switch rt {
	case 0:
		if rn == atom {
			panic("length")
		}
		r = mk(L, rn)
		for i := k(0); i < rn; i++ {
			m.k[2+r+i], x = resn(x)
		}
		return r, x
	case 1:
		rt = C
	case 7:
		rt = I
	case 0x0e:
		rt = F
	case 9:
		rt, o = Z, 8
	case 0x10: // k7 progression arrays
		r, o = mk(I, rn), m.k[2+x]
		for i := k(0); i < atm1(rn); i++ {
			m.k[2+r+i] = i + o
		}
		return r, drop(4, x)
	case 0x0f:
		r = mk(S, rn)
		p, xn, nn := 8+x<<2, atm1(m.k[x]&atom), k(0)
		for i := k(0); i < atm1(rn); i++ {
			ni := k(0)
			for j := p; j < p+xn; j++ {
				if m.c[j] == 0 {
					break
				}
				ni++
			}
			ds := c2s(mkb(m.c[p : p+ni]))
			m.k[2+r+i] = dex(ds, m.k[2+ds])
			p, xn, nn = p+ni+1, xn-ni-1, nn+1+ni
		}
		return r, drop(i(nn), x)
	case 0x14:
		r, x = resn(x)
		return flp(r), x
	case 0x15:
		r = mk(A, atom)
		m.k[2+r], x = resn(x)
		m.k[3+r], x = resn(x)
		return r, x
	default:
		panic("nyi")
	}
	r = mk(rt, rn) // C,I,F,Z
	p, rp, ln := 8+x<<2, 8+r<<2, atm1(rn)*lns[rt]
	copy(m.c[rp+o:rp+o+ln], m.c[p:p+ln])
	return r, drop(i(ln), x)
}
func sqr(x k) (r k) { // sqrt x
	return nm(x, 0, []f1{nil, nil, nil, func(r, x k) { m.f[r] = math.Sqrt(m.f[x]) }, nil})
}
func sin(x k) (r k) { // sin x
	return nm(x, 0, []f1{nil, nil, nil, func(r, x k) { m.f[r] = math.Sin(m.f[x]) }, nil})
}
func cos(x k) (r k) { // cos x
	return nm(x, 0, []f1{nil, nil, nil, func(r, x k) { m.f[r] = math.Cos(m.f[x]) }, nil})
}
func abs(x k) (r k) { // abs x
	return nm(x, F, []f1{nil, nil, func(r, x k) {
		m.k[r] = m.k[x]
		if i(m.k[r]) < 0 {
			m.k[r] = k(-i(m.k[r]))
		}
	}, func(r, x k) { m.f[r] = math.Abs(m.f[x]) }, func(r, x k) { m.z[r] = complex(math.Hypot(real(m.z[x]), imag(m.z[x])), 0) }})
}
func log(x k) (r k) { // log x
	return nm(x, 0, []f1{nil, nil, nil, func(r, x k) { m.f[r] = math.Log(m.f[x]) }, nil})
}
func exp(x k) (r k) { // exp x
	return nm(x, 0, []f1{nil, nil, nil, func(r, x k) { m.f[r] = math.Exp(m.f[x]) }, nil})
}
func rel(x k) (r k) { return zfn(0, x) } // real x
func ima(x k) (r k) { return zfn(1, x) } // imag x
func phi(x k) (r k) { return zfn(2, x) } // phase x
func zfn(c, x k) (r k) {
	t, n := typ(x)
	if t == L {
		r = mk(L, n)
		for i := k(0); i < n; i++ {
			m.k[2+r+i] = uf(zfn(c, inc(m.k[2+i+x])))
		}
		return dex(x, r)
	} else if t != Z {
		panic("type")
	}
	r = mk(F, n)
	rp, xp := ptr(r, F), ptr(x, Z)<<1
	if c == 1 {
		xp++
		c = 0
	}
	for i := k(0); i < atm1(n); i++ {
		switch c {
		case 0:
			m.f[rp+i] = m.f[xp+2*i]
		case 2:
			m.f[rp+i] = math.Atan2(m.f[xp+2*i+1], m.f[xp+2*i])
		}
	}
	return dex(x, r)
}
func cnj(x k) (r k) {
	t, n := typ(x)
	if t == L {
		return lrc(x, n, cnj)
	} else if t != Z {
		panic("type")
	}
	if m.k[x+1] != 1 {
		r = mk(Z, n)
		mv(r, x)
		dec(x)
	} else {
		r = x
	}
	rp := 2 + ptr(r, F)
	for i := k(0); i < atm1(n); i++ {
		m.f[rp] = -m.f[rp]
		rp += 2
	}
	return r
}
func rxp(x, y k) (r k) { // x expi y
	xt, yt, xn, yn := typs(x, y)
	if yt > F {
		panic("type")
	} else if yt < F {
		y = to(y, F)
	}
	r = mk(Z, yn)
	rp, yp := ptr(r, Z)<<1, ptr(y, F)
	for i := k(0); i < atm1(yn); i++ {
		s, c := math.Sincos(m.f[yp+i])
		m.f[rp+2*i], m.f[rp+2*i+1] = c, s
	}
	if x == 0 {
		return dex(y, r)
	}
	if xt > F {
		panic("type")
	} else if xt < F {
		x = to(x, F)
	}
	dx := k(1)
	if xn == atom {
		dx = 0
	}
	if xn != atom && yn == atom {
		r = take(xn, 0, r)
	}
	rp = ptr(r, Z) << 1
	xp := ptr(x, F)
	for i := k(0); i < atm1(m.k[r]&atom); i++ {
		m.f[rp+2*i] *= m.f[xp]
		m.f[rp+2*i+1] *= m.f[xp]
		xp += dx
	}
	return decr(x, y, r)
}
func zxp(x k) (r k) { // expi i
	t, n := typ(x)
	if t == L {
		return lrc(x, n, zxp)
	} else if t == A {
		return arc(x, n, zxp)
	} else if t > F {
		panic("type")
	}
	return rxp(0, x)
}

func putc(rc, rn k, c c) k { // assumes enough space
	m.c[rc+rn] = c
	return rn + 1
}
func putb(rc, rn k, b []c) k {
	rc += rn
	copy(m.c[rc:rc+k(len(b))], b)
	return rn + k(len(b))
}
func qt(c c) (c, bool) { // quote
	switch c {
	case '"', '\\':
		return c, true
	case '\n':
		return 'n', true
	case '\t':
		return 't', true
	case '\r':
		return 'r', true
	default:
		return c, false
	}
}
func op2(code k) ([]f2, []fc) {
	switch code - dyad {
	case 1: // +
		return []f2{nil, adC, adI, adF, adZ, nil}, nil
	case 2: // -
		return []f2{nil, sbC, sbI, sbF, sbZ, nil}, nil
	case 3: // *
		return []f2{nil, muC, muI, muF, muZ, nil}, nil
	case 4: // %
		return []f2{nil, nil, nil, diF, diZ, nil}, nil
	case 5: // &
		return []f2{nil, miC, miI, miF, miZ, miS}, nil
	case 6: // |
		return []f2{nil, maC, maI, maF, maZ, maS}, nil
	case 7: // <
		return nil, ltx
	case 8: // >
		return nil, gtx
	case 9: // =
		return nil, eqx
	default:
		return nil, nil
	}
}
func add(x, y k) (r k) { f, g := op2(1 + dyad); return nd(x, y, 0, f, g) } // x+y
func sub(x, y k) (r k) { f, g := op2(2 + dyad); return nd(x, y, 0, f, g) } // x-y
func mul(x, y k) (r k) { f, g := op2(3 + dyad); return nd(x, y, 0, f, g) } // x*y
func div(x, y k) (r k) { f, g := op2(4 + dyad); return nd(x, y, 0, f, g) } // x%y
func min(x, y k) (r k) { f, g := op2(5 + dyad); return nd(x, y, 0, f, g) } // x&y
func max(x, y k) (r k) { f, g := op2(6 + dyad); return nd(x, y, 0, f, g) } // x|y
func les(x, y k) (r k) { return nd(x, y, I, nil, ltx) }                    // x<y
func mor(x, y k) (r k) { return nd(x, y, I, nil, gtx) }                    // x>y
func eql(x, y k) (r k) { return nd(x, y, I, nil, eqx) }                    // x=y
func lgn(x, y k) (r k) { // x log y
	return nd(x, y, 0, []f2{nil, nil, nil, func(r, x, y k) { m.f[r] = math.Log(m.f[y]) / math.Log(m.f[x]) }, nil, nil}, nil)
}
func pow(x, y k) (r k) { // x exp y
	return nd(x, y, 0, []f2{nil, nil, nil, func(r, x, y k) { m.f[r] = math.Pow(m.f[x], m.f[y]) }, nil, nil}, nil)
}
func adC(r, x, y k) { m.c[r] = m.c[x] + m.c[y] }
func adI(r, x, y k) { m.k[r] = m.k[x] + m.k[y] }
func adF(r, x, y k) { m.f[r] = m.f[x] + m.f[y] }
func adZ(r, x, y k) { m.z[r] = m.z[x] + m.z[y] }
func sbC(r, x, y k) { m.c[r] = m.c[x] - m.c[y] }
func sbI(r, x, y k) { m.k[r] = m.k[x] - m.k[y] }
func sbF(r, x, y k) { m.f[r] = m.f[x] - m.f[y] }
func sbZ(r, x, y k) { m.z[r] = m.z[x] - m.z[y] }
func muC(r, x, y k) { m.c[r] = m.c[x] * m.c[y] }
func muI(r, x, y k) { m.k[r] = m.k[x] * m.k[y] }
func muF(r, x, y k) { m.f[r] = m.f[x] * m.f[y] }
func muZ(r, x, y k) { m.z[r] = m.z[x] * m.z[y] }
func diC(r, x, y k) { m.c[r] = m.c[y] / m.c[x] }
func diI(r, x, y k) { m.k[r] = k(i(m.k[y]) / i(m.k[x])) }
func diF(r, x, y k) { m.f[r] = m.f[x] / m.f[y] }
func diZ(r, x, y k) { m.z[r] = m.z[x] / m.z[y] }
func baC(r, x, y k) { m.c[r] = m.c[x] * (m.c[y] / m.c[x]) }
func baI(r, x, y k) { m.k[r] = k(i(m.k[x]) * (i(m.k[y]) / i(m.k[x]))) }
func mdC(r, x, y k) { m.c[r] = m.c[x] % m.c[y] }
func mdI(r, x, y k) { m.k[r] = k(imod(i(m.k[x]), i(m.k[y]))) }
func mdF(r, x, y k) { m.f[r] = math.Mod(m.f[x], m.f[y]) }
func miC(r, x, y k) { m.c[r] = m.c[ter(m.c[x] < m.c[y], x, y)] }
func miI(r, x, y k) { m.k[r] = m.k[ter(i(m.k[x]) < i(m.k[y]), x, y)] }
func miF(r, x, y k) { m.f[r] = m.f[ter(m.f[x] < m.f[y], x, y)] }
func miZ(r, x, y k) { m.z[r] = m.z[ter(ltZ(x, y), x, y)] }
func miS(r, x, y k) { m.f[r] = m.f[ter(ltS(x, y), x, y)] }
func maC(r, x, y k) { m.c[r] = m.c[ter(m.c[x] > m.c[y], x, y)] }
func maI(r, x, y k) { m.k[r] = m.k[ter(i(m.k[x]) > i(m.k[y]), x, y)] }
func maF(r, x, y k) { m.f[r] = m.f[ter(m.f[x] > m.f[y], x, y)] }
func maZ(r, x, y k) { m.z[r] = m.z[ter(gtZ(x, y), x, y)] }
func maS(r, x, y k) { m.f[r] = m.f[ter(gtS(x, y), x, y)] }
func ter(b bool, x, y k) k {
	if b {
		return x
	} else {
		return y
	}
}
func key(x, y k) (r k) { // x!y
	_, yt, xn, yn := typs(x, y)
	if xn == atom {
		x, xn = enl(x), 1
		if yt == A && yn == atom {
			y, yn = enlist(y), 1 // dont uf to table
		} else {
			y, yn = enl(y), 1
		}
	}
	if yn == atom {
		y, yn = ext(y, yt, xn), xn
	}
	if xn == 1 && yn != 1 {
		y, yn = enl(y), 1
	}
	if xn != yn {
		panic("length")
	}
	r = mk(A, atom)
	m.k[2+r] = x
	m.k[3+r] = y
	return r
}
func ext(x, t, n k) (r k) { // scalar extension
	if t >= L {
		r = mk(L, n)
		for i := k(0); i < n; i++ {
			m.k[2+i+r] = inc(x)
		}
		return dex(x, r)
	}
	r = mk(t, n)
	xp, rp, cp := ptr(x, t), ptr(r, t), cpx[t]
	for i := k(0); i < n; i++ {
		cp(rp+i, xp)
	}
	return dex(x, r)
}
func mch(x, y k) (r k) { // x~y
	r = mki(0)
	if match(x, y) {
		m.k[2+r] = 1
	}
	return decr(x, y, r)
}
func match(x, y k) (rv bool) { // recursive match
	if x == y {
		return true
	}
	t, n := typ(x)
	tt, nn := typ(y)
	if tt != t || nn != n {
		return false
	}
	n = atm1(n)
	switch t {
	case A:
		if match(m.k[2+x], m.k[2+y]) == false || match(m.k[3+x], m.k[3+y]) == false {
			return false
		}
		return true
	case N:
		return true
	default:
		if t > N {
			return match(m.k[x+2], m.k[y+2]) && match(m.k[x+3], m.k[y+3])
		}
		eq := eqx[t]
		if eq == nil {
			panic("type")
		}
		x, y = ptr(x, t), ptr(y, t)
		for j := k(0); j < n; j++ {
			if eq(x+j, y+j) == false {
				return false
			}
		}
		return true
	}
	return false
}
func cat(x, y k) (r k) { // x,y
	xt, yt, xn, yn := typs(x, y)
	if xn == atom && yn == 0 {
		return dex(y, enl(x))
	} else if yn == 0 {
		return dex(y, x)
	}
	if xt > A {
		x, xt, xn = enlist(x), L, 1
	}
	if yt > A {
		y, yt, yn = enlist(y), L, 1
	}
	switch {
	case xt < L && yt == xt:
		return ucat(x, y, xt, xn, yn)
	case xt == A:
		if yt != A {
			panic("type")
		}
		if xn == atom && yn == atom { // d,d
			// TODO (`b!`c!`d!1),`b!`c!`f!2
			return dex(y, amdv(x, inc(m.k[2+y]), inc(null), inc(m.k[3+y])))
		}
		if match(m.k[2+x], m.k[2+y]) == false {
			panic("nyi") // downtype to dict
		}
		nk, rn := m.k[m.k[2+x]]&atom, atm1(xn)+atm1(yn)
		r = mk(A, rn)
		m.k[2+r] = inc(m.k[2+x])
		m.k[3+r] = mk(L, nk)
		ik := mki(0)
		for i := k(0); i < nk; i++ {
			a := atx(inc(m.k[3+x]), inc(ik))
			if xn == atom {
				a = enl(a)
			}
			b := atx(inc(m.k[3+y]), inc(ik))
			if yn == atom {
				b = enl(b)
			}
			m.k[2+i+m.k[3+r]] = cat(a, b)
			m.k[2+ik]++
		}
		dec(ik)
		return decr(x, y, r)
	default:
		if xt != L {
			x = explode(x)
			xt, xn = typ(x)
		}
		if yt != L {
			y = explode(y)
			yt, yn = typ(y)
		}
	}
	if m.k[x+1] == 1 && bk(L, xn+yn) == bk(L, xn) {
		r = x
		m.k[r] = L<<28 | (xn + yn)
	} else {
		r = mk(L, xn+yn)
		for j := k(0); j < xn; j++ {
			m.k[2+r+j] = inc(m.k[2+x+j])
		}
		dec(x)
	}

	for j := k(0); j < yn; j++ {
		m.k[2+r+xn+j] = inc(m.k[2+y+j])
	}
	return dex(y, uf(r))
}
func ucat(x, y, t, xn, yn k) (r k) { // x, y same type < L
	xn, yn = atm1(xn), atm1(yn)
	cp, xp := cpx[t], ptr(x, t)
	if m.k[x+1] > 1 || bk(t, xn+yn) != bk(t, xn) {
		r = mk(t, xn+yn)
		rp := ptr(r, t)
		for i := k(0); i < xn; i++ {
			cp(rp+i, xp+i)
		}
	} else {
		r = x
		m.k[r] = t<<28 | (xn + yn)
	}
	rp, yp := xn+ptr(r, t), ptr(y, t)
	for i := k(0); i < yn; i++ {
		cp(rp+i, yp+i)
	}
	if r != x {
		dec(x)
	}
	return dex(y, r)
}
func lcat(x, y k) (r k) { // append anything to a list; no unify
	t, nl := typ(x)
	if t != L {
		panic("assert lcat")
	}
	if m.k[x+1] == 1 && bk(L, nl) == bk(L, nl+1) {
		m.k[2+x+nl] = y
		m.k[x] = L<<28 | (nl + 1)
		return x
	}
	r = mk(L, nl+1)
	for i := k(0); i < nl; i++ {
		m.k[2+i+r] = inc(m.k[2+i+x])
	}
	m.k[2+nl+r] = y
	return dex(x, r)
}
func lrc(x, n k, f func(k) k) (r k) { // list rec
	r = mk(L, n)
	for i := k(0); i < n; i++ {
		m.k[2+i+r] = f(inc(m.k[2+x+i]))
	}
	return dex(x, uf(r))
}
func lrc2(x, y, yn k, f func(k, k) k) (r k) {
	r = mk(L, yn)
	for i := k(0); i < yn; i++ {
		m.k[2+i+r] = f(inc(x), inc(m.k[2+y+i]))
	}
	return decr(x, y, uf(r))
}
func arc(x, n k, f func(k) k) (r k) { // dict rec
	r = mk(A, n)
	m.k[2+r] = inc(m.k[2+x])
	m.k[3+r] = f(inc(m.k[3+x]))
	return dex(x, r)
}
func arc2(x, y, yn k, f func(k, k) k) (r k) {
	r = mk(A, yn)
	m.k[2+r] = inc(m.k[2+y])
	m.k[3+r] = f(inc(x), inc(m.k[3+y]))
	return decr(x, y, r)
}
func ept(x, y k) (r k) { // x^y
	t, yt, n, yn := typs(x, y)
	if t == I && n == atom {
		x = til(x)
		t, n = typ(x)
	}
	if t != yt || t > L || n == atom {
		panic("type")
	} else if yn == atom {
		y, yn = enl(y), 1
	}
	eq, b, xp, yp := eqx[t], mk(I, n), ptr(x, t), ptr(y, t)
	all := true
	for i := k(0); i < n; i++ { // TODO: quadratic
		m.k[2+i+b] = 1
		for j := k(0); j < yn; j++ {
			if eq(xp+i, yp+j) {
				m.k[2+i+b], all = 0, false
				break
			}
		}
	}
	if all {
		return decr(b, y, x)
	}
	return dex(y, atx(x, wer(b)))
}
func tak(x, y k) (r k) { // x#y
	xt, yt, xn, yn := typs(x, y)
	if xt > N {
		return fil(x, y, false)
	}
	if yt == A {
		if xt == I {
			r = mk(A, yn)
			m.k[2+r] = tak(inc(x), inc(m.k[2+y]))
			m.k[3+r] = tak(inc(x), inc(m.k[3+y]))
			return decr(x, y, r)
		} else if xt == L {
			for i := k(0); i < xn; i++ { // (:e1;:e2)#t
				y = tak(inc(m.k[2+x+i]), y)
			}
			return dex(x, y)
		} else if xt == N && xn == 2 { // (:e)#t
			return atx(y, wer(atx(inc(y), x)))
		}
		return key(x, atx(y, inc(x)))
	}
	if xt != I {
		panic("type")
	}
	if xn == atom {
		xn = 1
	} else if xn == 0 {
		dec(x)
		return fst(y)
	}
	yn = atm1(yn)
	n, o := m.k[2+x+xn-1], k(0) // n:-1#x
	if i(n) < 0 {
		if yn != 0 {
			o = k(i(yn) + ((i(yn) + i(n)) % i(yn)))
		}
		n = k(-i(n))
	}
	if xn == 1 {
		return dex(x, take(n, o, y))
	}
	r, o = rsh(2+x, xn-1, n, o, y, yn)
	return decr(x, y, r)
}
func rsh(xp, xn, n, o, y, yn k) (r, oo k) { // reshape (with offset): (x,n)#y
	a := m.k[xp]
	if i(a) < 0 {
		panic("domain")
	}
	r = mk(L, a)
	for i := k(0); i < a; i++ {
		if xn > 1 {
			m.k[2+i+r], o = rsh(xp+1, xn-1, n, o, y, yn)
		} else {
			m.k[2+i+r] = take(n, o, inc(y))
			if yn != 0 {
				o = (o + n) % yn
			}
		}
	}
	return r, o
}
func take(n, o, y k) (r k) { // integer index and offset
	t, yn := typ(y)
	cp, yp := cpx[t], ptr(y, t)
	if yn == 0 {
		if t == L {
			return dex(y, take(n, o, enlist(mk(C, 0))))
		}
		return dex(y, take(n, o, inc(nan[t])))
	} else if yn == atom {
		yn = 1
	}
	r = mk(t, n)
	rp := ptr(r, t)
	for i := k(0); i < n; i++ {
		cp(rp+i, yp+((i+o)%yn))
	}
	return dex(y, uf(r))
}
func rot(x, y k) (r k) { return kxy(mks(".rot"), x, y) } // x rot y (rotate)
func drp(x, y k) (r k) { // x_y
	xt, t, xn, yn := typs(x, y)
	if xt > N {
		return fil(x, y, true)
	}
	if xt == N && xn == 2 && t == A && yn != atom { // (:e)_t  delete from t where e
		return dex(x, atx(y, wer(not(env(m.k[2+y], m.k[2+x], inc(m.k[3+y]), m.k[m.k[2+y]]&atom))))) // t@&~t@(:e)
	} else if t == A { // x_d  x_t  delete x from t
		u := ept(inc(m.k[y+2]), x)
		if r == m.k[y+2] {
			return y
		}
		r = mk(A, yn)
		m.k[2+r] = inc(u)
		m.k[3+r] = atx(y, u)
		return r
	} else if xt != I {
		panic("type")
	} else if yn == atom {
		if m.k[2+x] == 0 { // 0_atom
			return dex(x, enl(y))
		} else { // i_atom
			return decr(x, y, mk(t, 0))
		}
	} else if xn != atom {
		return cut(x, y)
	}
	n := m.k[2+x]
	return dex(x, uf(drop(i(n), y)))
}
func drop(x i, y k) (r k) { // integer index; does not unify
	t, yn := typ(y)
	if yn == atom {
		panic("type")
	} else if yn == 0 {
		return y
	} else if x >= i(yn) {
		return dex(y, mk(t, 0))
	}
	n, neg, o := k(x), false, k(x)
	if x < 0 {
		n, neg, o = k(-x), true, 0
	}
	yp, cp := ptr(y, t), cpx[t]
	if m.k[1+y] == 1 && t != L {
		if neg {
			return srk(y, t, yn, yn-n)
		}
		for i := k(0); i < yn-n; i++ {
			cp(yp+i, yp+o+i)
		}
		return uf(srk(y, t, yn, yn-n)) // uf? TODO rm
	}
	r = mk(t, yn-n)
	rp := ptr(r, t)
	for i := k(0); i < yn-n; i++ {
		cp(rp+i, yp+o+i)
	}
	return dex(y, r)
}
func fil(x, y k, drop bool) (r k) { // f#x f_x filter
	xt, yt, _, yn := typs(x, y)
	if xt != N+1 {
		panic("type")
	}
	v := y
	if yt == A {
		v = m.k[3+y]
		yn = m.k[v] & atom
	}
	yn = atm1(yn)
	idx := mk(I, yn)
	z := mki(0)
	for i := k(0); i < yn; i++ {
		r = cal(inc(x), enl(atx(inc(v), mki(i))))
		if m.k[r]>>28 != I {
			panic("type")
		}
		if match(r, z) == drop {
			m.k[2+idx+i] = m.k[2+r]
			if drop {
				m.k[2+idx+i] = 1
			}
		} else {
			m.k[2+idx+i] = 0
		}
		dec(r)
	}
	if yt == A {
		return decr(x, z, tak(atx(inc(m.k[2+y]), wer(idx)), y))
	}
	return decr(x, z, atx(y, wer(idx)))
}
func cut(x, y k) (r k) { // x_y
	xt, yt, xn, yn := typs(x, y)
	if xt != I || yn == atom {
		panic("type")
	}
	for i := k(0); i < xn; i++ {
		if a := m.k[2+x+i]; int32(a) < 0 || (i > 0 && m.k[1+x+i] > a) || a > yn {
			panic("domain")
		}
	}
	r = mk(L, xn)
	cp, yp := cpx[yt], ptr(y, yt)
	for i := k(0); i < xn; i++ {
		nn := yn
		if i < xn-1 {
			nn = m.k[3+i+x]
		}
		ln := nn - m.k[2+i+x]
		a := mk(yt, ln)
		yp, ap := yp+m.k[2+i+x], ptr(a, yt)
		for j := k(0); j < ln; j++ {
			cp(ap+j, yp+j)
		}
		yp += ln
		m.k[2+i+r] = uf(a)
	}
	return decr(x, y, r)
}
func cst(x, y k) (r k) { // x$y
	xt, yt, xn, yn := typs(x, y)
	if xt == I && xn == atom && yt < L && yn != atom { // [-]n$y (pad)
		return dex(x, pad(m.k[2+x], y))
	} else if xt == I && yt == C {
		return tak(x, y) // i#"c"
	} else if xt != S || xn != atom {
		panic("type")
	}
	if yt == L && yn > 0 {
		r = mk(L, yn)
		for i := k(0); i < yn; i++ {
			m.k[2+r+i] = cst(inc(x), inc(m.k[2+y+i]))
		}
		return decr(x, y, uf(r))
	} else if yt == A {
		r = mk(yt, yn)
		m.k[2+r] = inc(m.k[2+y])
		m.k[3+r] = cst(x, inc(m.k[3+y]))
		return dex(y, r)
	}
	ss, sn := sc(2 + x)
	s := m.c[ss]
	if yt == C { // strconv
		if sn == 0 || s == 'n' { // `$x `n|x
			return dex(x, c2s(y))
		} else if s == 'c' || s == ' ' { // `c$x `" "$x(for `csv?)
			return dex(x, y)
		} else if s == '*' { // `"*"$ for `csv?
			return dex(x, enlist(y))
		}
		num, o := aton(m.c[8+y<<2 : atm1(yn)+8+y<<2])
		if !o {
			num = dex(num, inc(nan[I]))
		}
		switch s {
		case 'i': // `i$x
			r = to(num, I)
		case 'f':
			r = to(num, F)
		case 'z':
			r = to(num, Z)
		default:
			panic("value")
		}
		return decr(x, y, r)
	}
	t, o := k(0), k(169)
	for i := o; i < o+15; i++ {
		if s == m.c[i] {
			t = i - o
		}
	}
	if yn == 0 && t > 0 && t < A {
		return decr(x, y, mk(t, 0))
	}
	if t < 1 || t >= L || yt >= L {
		panic("type")
	}
	return dex(x, to(y, t)) // TODO other conversions?
}
func c2s(x k) (r k) { // `$c
	r = mk(S, atom)
	if n := m.k[x] & atom; n == 0 {
		m.k[2+r] = 0
		return dex(x, r)
	} else if n == atom || n == 1 {
		m.k[2+r] = k(m.c[8+x<<2])
		return dex(x, r)
	}
	s := m.k[stab]
	sp, n := 2+s, m.k[s]&atom
	for i := k(0); i < n; i++ {
		if match(m.k[sp+i], x) {
			m.k[2+r] = i + 256
			return dex(x, r)
		}
	}
	m.k[2+r] = n + 256
	m.k[stab] = lcat(s, x)
	return r
}
func pad(n, y k) (r k) { // n$y
	t, yn := typ(y)
	if t > S || yn == atom { // k7 allows only C
		panic("type")
	}
	yp := ptr(y, t)
	mi, ma := yp, yp+yn
	if yn == n {
		return y
	} else if i(n) < 0 {
		n = k(-i(n))
		yp += yn - n
	}
	r = mk(t, n)
	rp, np, cp := ptr(r, t), ptr(nan[t], t), cpx[t]
	for i := k(0); i < n; i++ {
		if yp+i < mi || yp+i >= ma {
			cp(rp+i, np)
		} else {
			cp(rp+i, yp+i)
		}
	}
	return dex(y, r)
}
func fnd(x, y k) (r k) { // x?y
	t, yt, xn, yn := typs(x, y)
	if t == S && yt != S {
		switch m.k[2+x] {
		case 0: // `?
			return dex(x, res(y))
		case yb64: // `b64?
			return dex(x, b46(y))
		case yhex: // `hex?
			return dex(x, xeh(y))
		case ycsv: // `csv?
			return dex(x, vsc(mk(L, 0), y))
		default:
			panic("type")
		}
	} else if t == A {
		return dex(x, atx(inc(m.k[2+x]), fnd(inc(m.k[3+x]), y)))
	}
	if xn == atom || t != yt {
		panic("type")
	}
	r = mk(I, yn)
	yn = atm1(yn)
	eq, xp, yp := eqx[t], ptr(x, t), ptr(y, t)
	for j := k(0); j < yn; j++ {
		n := xn // TODO: or 0N?
		for i := k(0); i < xn; i++ {
			if eq(xp+i, yp+j) {
				n = i
				break
			}
		}
		m.k[2+j+r] = n
	}
	return decr(x, y, r)
}
func fnk(sv, sp k) k { // find key, nk on undefined
	t, n := typ(sv)
	if t != S || n == atom {
		panic("type")
	}
	for i := k(0); i < n; i++ {
		if m.k[2+sv+i] == sp {
			return i
		}
	}
	return n
}
func fns(x, y k) (r k) { // x find y
	xt, yt, xn, yn := typs(x, y)
	if xt != yt || xt != C || xn == atom || yn == atom {
		panic("type")
	}
	xp, yp := 8+x<<2, 8+y<<2
	r = mk(L, 0)
	for i := k(0); i < xn; {
		if m.c[xp+i] == m.c[yp] {
			n := k(1)
			for j := k(1); j < yn; j++ {
				if m.c[xp+i+j] == m.c[yp+j] {
					n++
				} else {
					break
				}
			}
			if n == yn {
				ri := mk(I, 2)
				m.k[2+ri] = i
				m.k[3+ri] = yn
				r = lcat(r, ri)
				i += yn
			} else {
				i++
			}
		} else {
			i++
		}
	}
	return decr(x, y, r)
}
func atx(x, y k) (r k) { // x@y
	xt, yt, xn, yn := typs(x, y)
	if xn == atom && xt == S {
		switch m.k[2+x] {
		case k('p'):
			return dex(x, prs(y))
		case k('k'):
			return dex(x, kst(y))
		case k('m'):
			return dex(x, mat(y))
		case ycsv:
			return dex(x, csv(y))
		case yb64:
			return dex(x, b64(y))
		case yhex:
			return dex(x, hex(y))
		default:
			panic("class")
		}
	} else if xt > N {
		return cal(x, enlist(y))
	} else if xn == atom && xt < A { // class error in k7
		if yn == atom {
			return dex(y, x)
		}
		return dex(y, take(yn, 0, x)) // (#y)#x
	}
	switch {
	case xt < L && yt == I:
		cp, xp := cpx[xt], ptr(x, xt)
		r = mk(xt, yn)
		yn = atm1(yn)
		rp, yp, np := ptr(r, xt), 2+y, ptr(nan[xt], xt)
		for i := k(0); i < yn; i++ {
			if ix := m.k[yp+i]; ix < 0 || ix >= xn {
				cp(rp+i, np)
			} else {
				cp(rp+i, xp+ix)
			}
		}
		return decr(x, y, r)
	case xt == L && yt == I:
		if yn == atom {
			if xi := m.k[2+y]; int32(xi) < 0 || xi >= xn {
				r = lnan(inc(x))
			} else {
				r = inc(m.k[2+x+m.k[2+y]])
			}
		} else {
			r = mk(L, yn)
			for i := k(0); i < yn; i++ {
				if xi := m.k[2+y+i]; int32(xi) < 0 || xi >= xn {
					m.k[2+r+i] = lnan(inc(x))
				} else {
					m.k[2+r+i] = inc(m.k[2+x+xi])
				}
			}
			r = uf(r)
		}
		return decr(x, y, r)
	case xt == A:
		keys := m.k[2+x]
		if y == keys { // x[!x]
			r = inc(m.k[3+x])
			return decr(x, y, r)
		} else if yt == A && yn == atom { // d@d (expr)
			nk := m.k[m.k[2+y]] & atom
			if nk == 0 {
				return dex(y, x)
			}
			r = mk(A, atom)
			yk, yv := m.k[2+y], m.k[3+y]
			m.k[2+r] = inc(yk)
			m.k[3+r] = mk(L, nk)
			nv := k(0)
			for i := k(0); i < nk; i++ {
				ri := atx(inc(x), atx(inc(yv), mki(i)))
				m.k[2+i+m.k[3+r]] = ri
				if n := m.k[ri] & atom; i == 0 {
					nv = n
				} else if n != nv {
					nv = atom
				}
			}
			if xn != atom { // preserve table if possible
				m.k[r] = A<<28 | nv
			}
			if m.k[r]&atom == atom {
				m.k[3+r] = uf(m.k[3+r])
			}
			return decr(x, y, r)
		} else if yt == N && yn == 2 { // (d|t)@:expr
			return decr(x, y, env(m.k[2+x], m.k[2+y], inc(m.k[3+x]), m.k[m.k[2+x]]&atom))
		} else if xn != atom { // t[I;S]
			idx := k(0)
			if yt == I {
				idx = y
			} else if yt == S { // t[S]
				return atx(flp(x), y)
			} else if yt == L && yn == 2 {
				idx = inc(m.k[2+y])
				x = flp(tak(inc(m.k[3+y]), flp(x)))
				dec(y)
			} else {
				panic("table-index")
			}
			r = mk(A, m.k[idx]&atom)
			v := mk(L, m.k[m.k[3+x]]&atom)
			m.k[2+r] = inc(m.k[2+x])
			for i := k(0); i < atm1(m.k[v]&atom); i++ {
				m.k[2+i+v] = atx(inc(m.k[2+i+m.k[3+x]]), inc(idx))
			}
			m.k[3+r] = uf(v)
			return decr(x, idx, r)
		}
		kt, nk := typ(keys)
		vt, _ := typ(m.k[3+x])
		if kt != yt {
			panic("type")
		} else if nk == 0 {
			if yn == atom {
				return decr(x, y, fst(take(1, 0, inc(m.k[3+x]))))
			} else {
				return decr(x, y, take(yn, 0, inc(m.k[3+x])))
			}
		}
		r = mk(vt, atm1(yn))
		cp, eq, kp, vp, rp, yp := cpx[vt], eqx[kt], ptr(keys, kt), ptr(m.k[3+x], vt), ptr(r, vt), ptr(y, yt)
		for i := k(0); i < atm1(yn); i++ {
			for j := k(0); j < nk; j++ {
				if eq(kp+j, yp+i) {
					cp(rp+i, vp+j)
					break
				}
				if j == nk-1 {
					if vt < L {
						cp(rp+i, ptr(nan[vt], vt))
					} else if vt == L {
						m.k[2+i+r] = lnan(inc(m.k[3+x]))
					} else {
						panic("index")
					}
				}
			}
		}
		if yn == atom {
			r = fst(r)
		}
		return decr(x, y, uf(r))
	case yt == F || yt == Z:
		return dot(x, y)
	case yt == L:
		if t := m.k[m.k[2+y]] >> 28; t == F || t == Z {
			return dot(x, y)
		}
		r = mk(L, yn)
		for i := k(0); i < yn; i++ {
			m.k[2+r+i] = atx(inc(x), inc(m.k[2+y+i]))
		}
		return decr(x, y, r)
	// case xt == L:
	//	missing element for a list is nan[type of first element]
	case match(y, null):
		return dex(y, x)
	default:
		panic("atx")
	}
}
func atm(x, y k) (r k) { // x@y (matrix indexing)
	xt, yt, xn, yn := typs(x, y)
	if xt == A && yt == L && yn == 2 { // d[`n;I] t[I;`n]
		a := mk(L, 2)
		if xn != atom {
			y = rev(y)
		}
		m.k[2+a] = fnd(inc(m.k[2+x]), inc(m.k[2+y]))
		m.k[3+a] = inc(m.k[3+y])
		y = dex(y, a)
		x = dex(x, inc(m.k[3+x]))
	} else if xt != L || xn == atom || yt != L || yn != 2 {
		panic("type")
	}
	a, b := inc(m.k[2+y]), inc(m.k[3+y])
	dec(y)
	//at, bt, an, _ := typs(a, b)
	a0, b0, an := match(a, null), match(b, null), m.k[a]&atom
	if a0 && b0 { // x[;]→x? or force a rectangular matrix?
		r = jota(cmc(x, xn))
		return decr(a, b, atm(x, l2(inc(r), r)))
	} else if b0 { // x[a;]
		return dex(b, atx(x, a))
	} else if a0 { // x[;b]
		dec(a)
		r, an = mk(L, xn), xn
		for i := k(0); i < xn; i++ {
			m.k[2+r+i] = inc(m.k[2+x+i])
		}
		dec(x)
	} else {
		r = atx(x, a)
	}
	if m.k[b]>>28 != I {
		panic("type")
	}
	for i := k(0); i < an; i++ {
		m.k[2+r+i] = atx(m.k[2+i+r], inc(b))
	}
	return dex(b, uf(r))
}
func cmc(x, n k) (r k) { // count matrix cols
	r = 0
	for i := k(0); i < n; i++ {
		if c := m.k[m.k[2+i+x]] & atom; c != atom && c > r {
			r = c
		}
	}
	return r
}
func csv(x k) (r k) { return kx(mks(".csv"), x) } // `csv@x
func vsc(x, y k) (r k) { // `csv?y  x 0: y, ("ii";"|")0:("2|3";"3|4";"4|5")
	xt, yt, xn, yn := typs(x, y)
	if yt != L || yn == 0 {
		return dex(x, y)
	}
	p := k(0)
	if xn == 0 { // detect format
		z, tb := spl(mkc(','), inc(m.k[2+y])), k(1)
		zn := m.k[z] & atom
		p = mk(C, zn)
		for i := k(0); i < zn; i++ {
			fp, fn := 8+m.k[2+i+z]<<2, m.k[m.k[2+i+z]]&atom
			if c := m.c[fp]; fn > 0 && (c == '-' || cr09(c)) {
				tb, m.c[8+i+p<<2] = 0, 'f'
				if _, ok := atoi(m.c[fp : fp+fn]); ok {
					m.c[8+i+p<<2] = 'i'
				}
			} else { // k7 identifies single chars as c not n
				m.c[8+i+p<<2] = 'n'
			}
		}
		if tb == 1 {
			return dex(p, flp(key(cst(inc(nans), z), vsc(x, drop(1, y)))))
		}
		x, xt, xn = decr(x, z, l2(p, mkc(','))), L, 2
	} else if xt == L && xn == 2 && m.k[m.k[3+x]]>>28 == C {
		p = m.k[2+x]
	} else {
		panic("type")
	}
	cc := wer(eql(mkc('z'), inc(p)))
	p = atx(p, wer(add(mki(1), eql(mkc('z'), inc(p)))))          // p@&1+"z"="ifcz"
	p = amdv(p, wer(eql(mkc('z'), inc(p))), inc(null), mkc('f')) // @[p;&"z"=p;"f"]
	m.k[2+x] = p
	r = kxy(mks(".vsc"), x, y)
	nc := m.k[cc] & atom
	if nc == 0 {
		return dex(cc, r)
	}
	nr := m.k[r]&atom - nc
	z, o := mk(L, nr), k(0)
	for i := k(0); i < nr; i++ {
		if o < nc && m.k[2+cc+o] == i {
			a := ptr(m.k[3+r+i+o], F)
			for j := k(0); j < m.k[m.k[3+r+i+o]]&atom; j++ {
				m.f[a+j] *= math.Pi / 180.0
			}
			m.k[2+z+i] = rxp(inc(m.k[2+r+i+o]), inc(m.k[3+r+i+o]))
			o++
		} else {
			m.k[2+z+i] = inc(m.k[2+r+i+o])
		}
	}
	return decr(r, cc, z)
}
func hex(x k) (r k) { // `hex@x
	t, n := typ(x)
	if t != C {
		panic("type")
	}
	n = atm1(n)
	r = mk(C, 2*n)
	rp, xp := ptr(r, C), ptr(x, C)
	for i := k(0); i < n; i++ {
		m.c[rp+2*i], m.c[rp+2*i+1] = hxb(m.c[xp+i])
	}
	return dex(x, r)
}
func xeh(x k) (r k) { // `hex?x
	t, n := typ(x)
	xp := ptr(x, C)
	if t != C || n%2 != 0 {
		panic("type")
	} else if m.c[xp] == '0' && m.c[xp+1] == 'x' {
		xp, n = xp+2, n-2
	}
	r = mk(C, n/2)
	rp := ptr(r, C)
	for i := k(0); i < n/2; i++ {
		h := m.c[xp+2*i]
		l := m.c[xp+2*i+1]
		if !crHx(h) || !(crHx(l)) {
			panic("value")
		}
		m.c[rp+i] = (xtoc(h) << 4) | xtoc(l)
	}
	return dex(x, r)
}
func b64(x k) (r k) { // `b64@x
	if t, n := typ(x); t != C || n == atom {
		panic("type")
	}
	panic("nyi") // {p:(3\x)#"=";x,:(3\#x)#0x00;..}
}
func b46(x k) (r k) { // `64?x
	panic("nyi")
}
func cal(x, y k) (r k) { // x.y
	xt, _, xn, yn := typs(x, y)
	if xt <= A { // TODO dict
		if yn != atom {
			if yn == 0 {
				return dex(y, x)
			} else if yn == 1 {
				return atx(x, fst(y))
			}
			return cal(cal(x, fst(inc(y))), drop(1, y)) // at depth
		}
		return atx(x, y)
	}
	y = explode(y)
	if xn == 1 || xn == 2 { // convert projected to full call
		l := m.k[x+3] // arg list with holes
		n := m.k[l] & atom
		if n != xt-N+yn {
			panic("valence")
		}
		a, l, yi := mk(L, n), m.k[x+3], k(0) // a: full arg vector
		for i := k(0); i < n; i++ {
			if v := m.k[2+l+i]; match(v, null) {
				m.k[2+a+i] = inc(m.k[2+y+yi])
				yi++
			} else {
				m.k[2+a+i] = inc(v)
			}
		}
		dec(y)
		r = m.k[x+2]
		if f := m.k[x+2]; xn == 1 { // lambda projection
			r, xn = inc(f), 0
		} else {
			r = mk(N+1, atom)
			m.k[2+r] = f
			if f > 255 {
				panic("assert proj")
			}
			if f >= dyad {
				m.k[r] = (N+2)<<28 | atom
			}
		}
		dec(x)
		x, y, xt, yn = r, a, N+n, n
	}
	if xn == 0 {
		return lambda(x, y)
	} else if xn == 3 {
		u, v := inc(m.k[2+x]), inc(m.k[3+x])
		dec(x)
		return cal(u, enl(cal(v, y)))
	}
	code := m.k[2+x]
	if code > 255 { // derived
		code >>= 8
		switch yn {
		case 1:
			f := table[code].(func(k, k) k)
			a, b := inc(m.k[3+x]), inc(m.k[2+y])
			dec(y)
			return dex(x, f(a, b))
		case 2:
			f := table[code+dyad].(func(k, k, k) k)
			a, b, c := inc(m.k[3+x]), inc(m.k[2+y]), inc(m.k[3+y])
			dec(y)
			return dex(x, f(a, b, c))
		default:
			panic("valence")
		}
	}
	switch xt {
	case N + 1:
		if yn != 1 {
			panic("valence")
		}
		f, a := table[code].(func(k) k), inc(m.k[2+y])
		dec(y)
		return dex(x, f(a))
	case N + 2:
		if yn != 2 {
			panic("valence")
		}
		f, a, b := table[code].(func(k, k) k), inc(m.k[2+y]), inc(m.k[3+y])
		dec(y)
		return dex(x, f(a, b))
	default:
		panic("nyi")
	}
}
func cal2(f, x, y k) (r k) { return cal(f, l2(x, y)) }
func lambda(x, y k) (r k) { // call lambda
	v := (m.k[x] >> 28) - N
	if v < 1 || v > 9 {
		panic("valence")
	}
	if yt, yn := typ(y); yt != L || yn != v {
		panic("args")
	}
	loc, l := m.k[2+m.k[3+x]], m.k[3+m.k[3+x]] // arguments and parse tree
	lt, nl := typ(l)
	if nl == 0 {
		return decr(x, y, inc(null))
	} else if lt != L {
		panic("type")
	}
	return dex(x, env(loc, l, y, v))
}
func env(e, l, y, v k) (r k) { // call-with-env
	n := m.k[e] & atom
	vv := mk(L, n) // save old values
	for i := k(0); i < n; i++ {
		name := atx(inc(e), mki(i))
		m.k[2+i+vv] = lupo(inc(name))
		if i < v {
			dec(asn(name, inc(m.k[2+y+i]), inc(null)))
		} else {
			dec(asn(name, mk(I, 0), inc(null)))
		}
	}
	dec(y)
	r = evl(inc(l))
	for i := k(0); i < n; i++ { // restore old values
		name, w := atx(inc(e), mki(i)), m.k[2+i+vv]
		if w == 0 {
			m.k[2+i+vv] = inc(null)
			dec(del(name))
		} else {
			dec(asn(name, inc(w), inc(null)))
		}
	}
	return dex(vv, r)
}
func ltr(x k) (r k) { // `p@{..} (lambda tree, k7: .{...})
	t, n := typ(x)
	if t < N+1 || n != 0 {
		panic("nolambda")
	}
	r = mk(L, 4)
	args, tree, v := m.k[2+m.k[3+x]], m.k[3+m.k[3+x]], t-N
	m.k[2+r] = inc(tree)             // parse tree
	m.k[3+r] = inc(m.k[2+x])         // string
	m.k[4+r] = take(v, 0, inc(args)) // args
	m.k[5+r] = drop(i(v), inc(args)) // locals
	return dex(x, r)
}
func qot(x k) (r k) { return drv(0, x) } // '
func sla(x k) (r k) { return drv(1, x) } // /
func bsl(x k) (r k) { return drv(2, x) } // \
func qtc(x k) (r k) { return drv(3, x) } // ':
func slc(x k) (r k) { return drv(4, x) } // /:
func bsc(x k) (r k) { return drv(5, x) } // \:
func drv(op k, x k) (r k) { // derived function
	r = mk(N+1, atom)
	m.k[2+r] = (33 + op) << 8 // op: 0(') 1(/) 2(\) 3(':) 4(/:) 5(\:)
	m.k[3+r] = x
	return r
}
func ech(f, x k) (r k) { // f'x
	if t := m.k[f] >> 28; t < N { // x'y (bar: x*x/y)
		if t > I {
			panic("type")
		}
		if m.k[x]>>28 > I {
			x = to(x, I)
		}
		return nd(f, x, 0, []f2{nil, baC, baI, nil, nil, nil}, nil)
	}
	t, n := typ(x)
	if t == A {
		r = mk(A, atom)
		m.k[2+r] = inc(m.k[2+x])
		m.k[3+r] = ech(f, inc(m.k[3+x]))
		return dex(x, r)
	}
	if n == atom {
		return atx(f, x)
	}

	// simple version:
	//  r = mk(L, n)
	//  for i := k(0); i < n; i++ {
	//  	m.k[2+r+i] = atx(inc(f), atx(inc(x), mki(i)))
	//  }
	//  return decr(f, x, uf(r))
	// this is optimized for unitype result vectors and in-place:
	uni, r, rp, rt, ri, cp := true, k(0), k(0), k(0), k(0), f1(nil)
	for i := k(0); i < n; i++ {
		ri = atx(inc(f), atx(inc(x), mki(i)))
		if tt, nn := typ(ri); i == 0 {
			if tt < L && nn == atom { // uni
				if m.k[1+x] == 1 && tt == t { // inplace
					r, rt, m.k[1+x] = x, tt, 2
				} else {
					r, rt = mk(tt, n), tt
				}
				rp, cp = ptr(r, tt), cpx[rt]
			} else {
				r, uni = mk(L, 0), false
			}
		} else if uni && tt != rt || nn != atom { // upgrade
			r, uni = explode(take(i, 0, r)), false
		}
		if uni {
			cp(rp+i, ptr(ri, rt))
			dec(ri)
		} else {
			r = lcat(r, ri)
		}
	}
	return decr(f, x, r)
}
func ecd(f, x, y k) (r k) { // x f'y (each pair)
	xt, yt, xn, yn := typs(x, y)
	if xn == atom && yn == atom {
		return cal2(f, x, y)
	}
	if xn == atom {
		x, xn = ext(x, xt, yn), yn
	} else if yn == atom {
		y, yn = ext(y, yt, xn), xn
	}
	if xn != yn {
		panic("length")
	}
	r = mk(L, xn)
	for i := k(0); i < xn; i++ {
		m.k[2+i+r] = cal2(inc(f), atx(inc(x), mki(i)), atx(inc(y), mki(i)))
	}
	dec(f)
	return decr(x, y, uf(r))
}
func ecp(f, x k) (r k) { // f':x (each prior)
	ft, t, fn, n := typs(f, x)
	if ft == I && fn == atom { // n':y (window)
		return dex(f, win(m.k[2+f], n, x))
	}
	if t == A {
		return arc2(f, x, n, ecp)
	} else if t > L || n == atom {
		panic("class")
	}
	x0 := k(0)
	if code := m.k[2+f]; m.k[f]&atom == atom && code < 256 {
		switch code - dyad {
		//case 12: // ,
		//return decr(x, f, mk(t, 0))
		case 1, 2, 6: // +-|
			x0 = to(mki(0), t)
		case 3, 5: // *&
			x0 = to(mki(1), t)
		}
	}
	if x0 == 0 {
		x0 = inc(nan[t])
	}
	if x, t, op := scop1(f, x); op != nil {
		r = mk(t, n)
		rp, xp := ptr(r, t), ptr(x, t)
		op(rp, xp, ptr(x0, t))
		for i := k(1); i < n; i++ {
			op(rp+i, xp+i, xp+i-1)
		}
		return decr(x0, x, r)
	}
	r = mk(L, n)
	m.k[2+r] = cal2(inc(f), atx(inc(x), mki(0)), x0)
	for i := k(1); i < n; i++ {
		m.k[2+r+i] = cal2(inc(f), atx(inc(x), mki(i)), atx(inc(x), mki(i-1)))
	}
	return decr(f, x, uf(r))
}
func win(n, ny, y k) (r k) { // n':y (window)
	if ny == atom || i(n) < 0 || n > ny {
		panic("length")
	} else if n == 0 { // 0':y → 3':0,y,0
		return win(3, ny+2, cat(cat(mki(0), y), mki(0)))
	}
	a := jota(n)
	r = mk(L, 1+ny-n)
	m.k[2+r] = atx(inc(y), inc(a))
	for i := k(0); i < ny-n; i++ {
		for j := k(0); j < n; j++ {
			m.k[2+a+j]++
		}
		m.k[3+i+r] = atx(inc(y), inc(a))
	}
	return decr(a, y, r)
}
func epi(f, x, y k) (r k) { // x f':y (each prior initial)
	n := m.k[y] & atom
	if m.k[x]&atom == atom && m.k[f]&atom == atom {
		if x, y, t, op := scop2(f, x, y); op != nil {
			r = mk(t, n)
			rp, xp, yp := ptr(r, t), ptr(x, t), ptr(y, t)
			op(rp, yp, xp)
			for i := k(1); i < n; i++ {
				op(rp+i, yp+i, yp+i-1)
			}
			return decr(x, y, r)
		}
	}
	r = mk(L, n)
	m.k[2+r] = cal2(inc(f), atx(inc(y), mki(0)), x)
	for i := k(1); i < n; i++ {
		m.k[2+r+i] = cal2(inc(f), atx(inc(y), mki(i)), atx(inc(y), mki(i-1)))
	}
	return decr(f, y, uf(r))
}
func ecr(f, x, y k) (r k) { // x f/: y
	xt, yt, _, yn := typs(x, y)
	if xt > L || yt > L { // TODO A
		panic("type")
	} else if yn == atom {
		return cal2(f, x, y)
	}
	r = mk(L, yn)
	for i := k(0); i < yn; i++ {
		m.k[2+r+i] = cal2(inc(f), inc(x), atx(inc(y), mki(i)))
	}
	dec(f)
	return decr(x, y, uf(r))
}
func ecl(f, x, y k) (r k) { // x f\: y
	xt, yt, xn, _ := typs(x, y)
	if xt > L || yt > L { // TODO A
		panic("type")
	} else if xn == atom {
		return cal2(f, x, y)
	}
	r = mk(L, xn)
	for i := k(0); i < xn; i++ {
		m.k[2+r+i] = cal2(inc(f), atx(inc(x), mki(i)), inc(y))
	}
	dec(f)
	return decr(x, y, uf(r))
}
func ovr(f, x k) (r k) { // f/x
	if t, n := typ(f); t < N { // x/y (idiv, dot)
		if t == L || n != atom {
			return dot(f, x)
		} else if t > I {
			panic("type")
		}
		if m.k[x]>>28 > I {
			x = to(x, I)
		}
		return nd(f, x, 0, []f2{nil, diC, diI, nil, nil, nil}, nil)
	} else if t == N+1 { // fixed
		x0 := inc(x)
		for {
			r = x
			x = atx(inc(f), inc(r))
			if match(x, r) || match(x, x0) {
				break
			}
			dec(r)
		}
		dec(x0)
		return decr(f, r, x)
	}
	if t, n := typ(x); n == 0 { // verb depended values for empty y
		if code := m.k[2+f]; m.k[f]&atom == atom && code < 256 {
			switch code - dyad {
			case 12: // ,
				return decr(x, f, mk(t, 0))
			case 1, 2, 6: // +-|
				if t > C {
					return decr(x, f, to(mki(0), t))
				}
			case 3, 5: // *&
				if t > C {
					return decr(x, f, to(mki(1), t))
				}
			}
		}
		return decr(x, f, inc(nan[t]))
	} else if n == atom && t != A {
		return dex(f, x)
	}
	if x, t, op := scop1(f, x); op != nil {
		n := m.k[x] & atom
		if n == atom {
			return x
		}
		if m.k[f]&atom == atom && m.k[2+f] == 1+dyad { // +/
			return sum(x)
		}
		r := mk(t, atom)
		rp, xp, cp := ptr(r, t), ptr(x, t), cpx[t]
		cp(rp, xp)
		for i := k(1); i < n; i++ {
			op(rp, rp, xp+i)
		}
		return dex(x, r)
	}
	return ovsc(f, x, false)
}
func scn(f, x k) (r k) { // f\x
	if xt := m.k[f] >> 28; xt < N { // x\y (mod, solve, qr, inv)
		if xt == L {
			return slv(f, x)
		} else if xt > F {
			panic("type")
		}
		return nd(x, f, 0, []f2{nil, mdC, mdI, mdF, nil, nil}, nil)
	} else if xt == N+1 { // scan fixed
		l := mk(L, 0)
		x0 := inc(x)
		for {
			r = x
			l = lcat(l, inc(r))
			x = atx(inc(f), inc(r))
			if match(x, r) || match(x, x0) {
				dec(x)
				dec(x0)
				break
			}
			dec(r)
		}
		return decr(f, r, uf(l))
	} else if x, t, op := scop1(f, x); op != nil {
		n := m.k[x] & atom
		if n == atom {
			return x
		}
		r := mk(t, n)
		rp, xp, cp := ptr(r, t), ptr(x, t), cpx[t]
		cp(rp, xp)
		for i := k(1); i < n; i++ {
			op(rp+i, rp+i-1, xp+i)
		}
		return dex(x, r)
	}
	return ovsc(f, x, true)
}
func ovsc(f, x k, scan bool) (r k) {
	t, n := typ(x)
	if t == A {
		if scan {
			r = mk(A, atom)
			m.k[2+r] = inc(m.k[2+x])
			m.k[3+r] = ovsc(f, inc(m.k[3+x]), scan)
			return dex(x, r)
		} else {
			r = inc(m.k[3+x])
			dec(x)
			x = r
			t, n = typ(x)
		}
	}
	if t > L {
		panic("type")
	}
	if n == atom {
		x = enl(x)
	}
	if n == 0 {
		return dex(f, x)
	} else if n == 1 {
		return dex(f, fst(x))
	}
	if scan {
		r = mk(L, n)
		m.k[2+r] = atx(inc(x), mki(0))
		p := m.k[2+r]
		for i := k(1); i < n; i++ {
			p = cal2(inc(f), inc(p), atx(inc(x), mki(i)))
			m.k[2+i+r] = p
		}
		return decr(x, f, uf(r))
	} else {
		r = atx(inc(x), mki(0))
		for i := k(1); i < n; i++ {
			r = cal2(inc(f), r, atx(inc(x), mki(i)))
		}
		return decr(x, f, r)
	}
}
func ovi(f, x, y k) (r k) { // x f/y
	xt, _, xn, yn := typs(x, y)
	if xt > N {
		return whl(f, x, y)
	} else if m.k[f]>>28 == N+1 && xt == I { // for
		n := m.k[2+x]
		r = y
		for i := k(0); i < n; i++ {
			r = atx(inc(f), r)
		}
		return decr(f, x, r)
	} else if yn == atom {
		panic("class")
	}
	if x, y, t, op := scop2(f, x, y); op != nil {
		r := mk(t, xn)
		xn, yn = atm1(xn), atm1(yn)
		rp, xp, yp, cp := ptr(r, t), ptr(x, t), ptr(y, t), cpx[t]
		for i := k(0); i < xn; i++ {
			cp(rp, xp+i)
			for j := k(0); j < yn; j++ {
				op(rp, rp, yp+j)
			}
			rp++
		}
		return decr(x, y, r)
	}
	r = x
	for i := k(0); i < yn; i++ {
		r = cal2(inc(f), r, atx(inc(y), mki(i)))
	}
	return decr(y, f, r)
}
func sci(f, x, y k) (r k) { // x f\y
	xt, _, xn, yn := typs(x, y)
	if xt > N {
		return whls(f, x, y)
	} else if m.k[f]>>28 == N+1 && xt == I { // scan-for
		n := m.k[2+x]
		r = y
		l := lcat(mk(L, 0), inc(r))
		for i := k(0); i < n; i++ {
			r = atx(inc(f), r)
			l = lcat(l, inc(r))
		}
		dec(r)
		return decr(f, x, uf(l))
	} else if yn == atom {
		panic("class")
	}
	if xn == atom {
		if x, y, t, op := scop2(f, x, y); op != nil {
			r = mk(t, yn)
			xp, yp, rp := ptr(x, t), ptr(y, t), ptr(r, t)
			op(rp, xp, yp)
			for i := k(1); i < yn; i++ {
				op(rp+i, rp+i-1, yp+i)
			}
			return decr(x, y, r)
		}
	}
	r = mk(L, yn)
	for i := k(0); i < yn; i++ {
		x = cal2(inc(f), x, atx(inc(y), mki(i)))
		m.k[2+i+r] = inc(x)
	}
	dec(f)
	return decr(x, y, uf(r))
}
func scop1(f, x k) (k, k, f2) {
	if m.k[f]&atom != atom || m.k[2+f] > 255 {
		return x, 0, nil
	}
	t := m.k[x] >> 28
	if t >= L {
		return x, 0, nil
	}
	opx, _ := op2(m.k[2+f])
	if opx == nil {
		return x, 0, nil
	}
	dec(f)
	op := opx[t]
	if t == C && opx[C] == nil {
		t = I
	}
	if (t == I && opx[I] == nil) || (t == Z && opx[Z] == nil) {
		t = F
	}
	x = to(x, t)
	return to(x, t), t, op
}
func scop2(f, x, y k) (k, k, k, f2) {
	if m.k[f]&atom != atom || m.k[2+f] > 255 {
		return x, y, 0, nil
	}
	xt, yt := m.k[x]>>28, m.k[y]>>28
	if xt >= L || yt >= L {
		return x, y, 0, nil
	}
	opx, _ := op2(m.k[2+f])
	if opx == nil {
		return x, y, 0, nil
	}
	dec(f)
	t := ntyps(xt, yt, opx, nil)
	op := opx[t]
	if xt < t {
		x, xt = to(x, t), t
	}
	if yt < t {
		y, yt = to(y, t), t
	}
	return x, y, t, op
}
func whl(f, x, y k) (r k) { // g f/y
	r = y
	for {
		b := atx(inc(x), inc(r))
		br := m.k[2+b]
		if bt, bn := typ(b); bt != I || bn != atom {
			panic("type")
		}
		dec(b)
		if br == 0 {
			break
		}
		r = atx(inc(f), r)
	}
	return decr(f, x, r)
}
func whls(f, x, y k) (r k) { // g f\y
	r = y
	l, b, br := lcat(mk(L, 0), inc(r)), k(0), k(0)
	for {
		b = atx(inc(x), inc(r))
		br = m.k[2+b]
		if bt, bn := typ(b); bt != I || bn != atom {
			panic("type")
		}
		dec(b)
		if br == 0 {
			break
		}
		r = atx(inc(f), r)
		l = lcat(l, inc(r))
	}
	dec(r)
	return decr(f, x, uf(l))
}
func rdl(x k) (r k) { // 0:x
	rd := table[21].(func(k) k)
	return spl(inc(nans), rd(x))
}
func wrl(x, y k) (r k) { // x 0:y
	if m.k[x]>>28 == L {
		return vsc(x, y)
	}
	w := table[21+dyad].(func(k, k) k)
	return w(x, jon(mkc('\n'), y))
}
func lun(x k) (r k) { // 8:x or .. \x (display)
	wr := table[21+dyad].(func(k, k) k)
	dec(wr(inc(nans), cat(kst(inc(x)), mkc('\n'))))
	return x
}
func deb(x k) (r k) { // 9:x (println)
	r = kst(inc(x))
	t, n := typ(r)
	if t != C {
		return dex(r, x)
	}
	p := ptr(r, C)
	println(string(m.c[p : p+n]))
	return dex(r, x)
}
func lod(x k) (r k) {
	dec(asn(mks(".f"), tak(min(mki(8), cnt(inc(x))), inc(x)), inc(null)))
	evp(red(x))
	return inc(null)
}
func cmd(x k) (r k) {
	xp := 8 + x<<2
	switch m.c[xp] {
	case 'a':
		return dex(x, adler())
	case 'b':
		return dex(x, stats())
	case 'v':
		return dex(x, lsv())
	case 'c':
		return dex(x, clv())
	case 'h':
		return dex(x, hlp())
	case 'l':
		return lod(trm(x))
	case 'k':
		return key(inc(m.k[kkey]), inc(m.k[kval])) // dump ktree
	case 's':
		if r = lupo(mks(".stk")); r != 0 {
			w := table[21+dyad].(func(k, k) k)
			dec(w(inc(nans), cat(r, mkc('\n'))))
		}
		return dex(x, inc(null))
	case 'm': // \m x (matrix display)
		w := table[21+dyad].(func(k, k) k)
		dec(w(inc(nans), cat(jon(mkc('\n'), mat(val(trm(x)))), mkc('\n'))))
		return inc(null)
	case '\\':
		exi := table[40].(func(k) k)
		if m.k[x]&atom > 1 {
			return dex(x, exi(mki(1)))
		}
		return dex(x, exi(mki(0)))
	default:
		panic("undefined")
	}
}
func trm(x k) (r k) { // trim 1st char and additional spaces
	xp, n, p := ptr(x, C), m.k[x]&atom, i(1)
	for i := k(1); i < n; i++ {
		if m.c[xp+i] == ' ' {
			p++
		} else {
			break
		}
	}
	return drop(p, x)
}
func evp(x k) { // parse-eval-print
	t, n := typ(x)
	if t != C {
		panic("type")
	}
	if n > 1 && m.c[8+x<<2] == '\\' {
		out(cmd(drop(1, x)))
		return
	}
	r, asn := par(x, 8+x<<2), false
	a, s := r, inc(nans)
	if t, n := typ(a); t == L && n > 1 && match(m.k[2+a], s) {
		dec(s)
		a = m.k[1+n+a] // last of multiple statements
	}
	if t, n := typ(a); t == L && n > 1 {
		if f := m.k[2+a]; m.k[f]>>28 == N+2 && m.k[2+f] == dyad { // (::;`x;v)
			asn = true
		} else if m.k[f]>>28 == N+1 && m.k[f]&atom == atom && n == 3 { // (*:;`x;v) modified assignment
			asn = true
		}
	}
	r = evl(r)
	if asn {
		dec(r)
		return
	}
	out(r)
}
func out(x k) {
	if match(x, null) {
		return
	}
	w := table[21+dyad].(func(k, k) k)
	dec(w(inc(nans), cat(kst(x), mkc('\n'))))
}
func spl(x, y k) (r k) { // x\:y (split)
	xt, yt, xn, yn := typs(x, y)
	if xt == I && yt == I { // encode ⊤
		return enc(x, y)
	}
	if yt != C || yn == atom {
		panic("type")
	}
	if yn == 0 {
		return decr(x, y, mk(L, 0)) // k7 returns () instead of ,""
	}
	yp := ptr(y, C)
	if xt == S && match(x, nan[S]) {
		if yn > 0 && m.c[yp+yn-1] == '\n' { // `\:y ignores trailing newline
			y, yn = drop(-1, y), yn-1
			yp = ptr(y, C)
		}
		dec(x)
		x, xt, xn = mkc('\n'), C, atom
	}
	if xt != C || xn != atom {
		panic("type")
	}
	k0 := k(0)
	idx := cat(cat(mki(k0-1), wer(eql(inc(y), x))), mki(yn))
	_, n := typ(idx)
	r = mk(L, n-1)
	for i := k(0); i < n-1; i++ {
		a, b := 1+m.k[2+i+idx], m.k[3+i+idx]
		m.k[2+i+r] = mkb(m.c[yp+a : yp+b])
	}
	return decr(idx, y, r)
}
func jon(x, y k) (r k) { // x/:y (join)
	xt, yt, _, yn := typs(x, y)
	if xt == I && yt == I { // decode ⊥
		return dcd(x, y)
	} else if xt == N+1 {
		f := table[39].(func(k, k) k)
		return f(x, y)
	}
	if yt != L {
		panic("type")
	}
	if xt == S {
		dec(x)
		y, yn = lcat(y, mk(C, 0)), yn+1
		x, xt = mkc('\n'), C
	} else if xt != C {
		panic("type")
	}
	if yn == 0 {
		return decr(x, y, mk(C, 0))
	}
	yn = atm1(yn)
	r = cat(mk(C, 0), inc(m.k[2+y]))
	for i := k(1); i < yn; i++ {
		if e := m.k[2+i+y]; m.k[e]>>28 != C {
			panic("type")
		} else {
			r = cat(cat(r, inc(x)), inc(e))
		}
	}
	return decr(x, y, r)
}
func imod(x, y i) (r i) { // k: y\x go: x%y differs for x<0
	if y < 0 {
		panic("domain")
	}
	r = x % y
	if r < 0 {
		r += y
	}
	return r
}
func enc(x, y k) (r k) { // x\:y (encode y in base x)
	if ny := m.k[y] & atom; ny != atom {
		r = mk(L, ny)
		mx := k(0)
		for i := k(0); i < ny; i++ {
			m.k[2+i+r] = enc(inc(x), mki(m.k[2+i+y]))
			if n := m.k[m.k[2+i+r]] & atom; n > mx {
				mx = n
			}
		}
		for i := k(0); i < ny; i++ {
			if n := m.k[m.k[2+i+r]] & atom; n == 0 {
				dec(m.k[2+i+r])
				m.k[2+i+r] = take(mx, 0, mki(0))
			} else if n < mx {
				m.k[2+i+r] = amd(take(mx, 0, mki(0)), tak(mki(n-mx), jota(mx)), inc(null), m.k[2+i+r])
			}
		}
		return decr(x, y, flp(r))
	}
	x = rev(x)
	xn, a, xp, n, xx := m.k[x]&atom, k(0), 2+x, i(m.k[2+y]), i(0)
	if n == 0 && xn == atom {
		return decr(x, y, mk(I, 0))
	} else if n == 0 {
		return decr(x, y, take(xn, 0, mki(0)))
	} else if n < 0 {
		panic("domain")
	} else if xn != atom {
		a = 1
	} else if xn == 0 {
		panic("type")
	}
	r = mk(I, 0)
	for j := k(0); ; j++ {
		if aj := a * j; aj >= xn {
			xx = i(m.k[xp+xn-1])
		} else {
			xx = i(m.k[xp+aj])
		}
		m := imod(n, xx)
		r = cat(r, mki(k(m)))
		n /= xx
		if (xn == atom && n <= 0) || (xn != atom && j >= xn-1) {
			break
		}
	}
	return decr(x, y, rev(r))
}
func dcd(x, y k) (r k) { // x/:y (decode y given in base x) {{z+y*x}/[0;x;y]}
	xn, yn, a := m.k[x]&atom, m.k[y]&atom, k(1)
	if yn == atom {
		panic("class")
	} else if xn == atom {
		a = 0
	} else if xn != yn {
		panic("length")
	}
	s, xp, yp := k(0), 2+x, 2+y
	for i := k(0); i < yn; i++ {
		s = s*m.k[xp+i*a] + m.k[yp+i]
	}
	return decr(x, y, mki(s))
}
func bin(x, y k) (r k) { // x bin y
	t, yt, xn, yn := typs(x, y)
	if t == L {
		y = explode(y)
	} else if yt != t || t > S || xn == atom {
		panic("type")
	}
	r = mk(I, yn)
	yn = atm1(yn)
	gt, xp, yp := gtx[t], ptr(x, t), ptr(y, t)
	for i := k(0); i < yn; i++ {
		m.k[2+i+r] = ibin(xp, xn, yp+i, gt)
	}
	return decr(x, y, r)
}
func ibin(xp, n, yp k, gt func(x, y k) bool) (r k) {
	i, j, h := k(0), n-1, k(0)
	for int32(i) <= int32(j) {
		h = (i + j) >> 1
		if gt(xp+h, yp) {
			j = h - 1
		} else {
			i = h + 1
		}
	}
	return i - 1
}
func sel(t, c, b, a k) (r k) { // #[t;c;b;a] select
	if tt, n := typ(t); tt != A || n == atom {
		panic("type")
	}
	if cn := m.k[c] & atom; cn != 0 {
		t = tak(inc(c), t) // t where c
	}
	dec(c)
	g := false
	if match(b, null) || (m.k[b]>>28 == A && m.k[m.k[2+b]]&atom == 0) {
		dec(b)
	} else {
		d := k(0)
		t, d = gby(t, b)
		dec(d)
		g = true
	}
	if match(a, null) == false {
		if g {
			l, n := m.k[3+t], m.k[m.k[3+t]]&atom
			for i := k(0); i < n; i++ {
				m.k[2+i+l] = atx(m.k[2+i+l], inc(a))
			}
			m.k[3+t] = uf(m.k[3+t])
		} else {
			t = atx(t, inc(a))
			if tt, nn := typ(t); tt == A && nn == atom {
				t = enl(t)
			}
		}
	}
	return dex(a, t)
}
func gby(t, b k) (k, k) { // select by b from t
	var a k
	if bt := m.k[b] >> 28; bt == S { // #[t;();`b;(0#`)!()]  (`b_t)@=t`b
		a = inc(b)
	} else if bt == A { // select a:b from t  (`a_t)@=t@`a!b
		a = inc(m.k[2+b])
	} else {
		panic("type")
	}
	g := grp(atx(inc(t), inc(b)))
	return decr(b, t, atx(drp(a, inc(t)), inc(g))), g
}
func udt(t, c, b, a k) (r k) { // _[t;c;b;a] update
	if t, n := typ(t); t == S {
		panic("nyi") // inplace
	} else if t != A || n == atom {
		panic("type")
	}
	if !match(b, null) { // update-by
		q := inc(t)
		if cn := m.k[c] & atom; cn != 0 {
			q = tak(inc(c), q) // t where c
		}
		dec(c)
		g, q := gby(q, b)

		l, n := m.k[3+g], m.k[m.k[3+g]]&atom
		for i := k(0); i < n; i++ {
			m.k[2+i+l] = atx(m.k[2+i+l], inc(a))
		}
		dec(a)

		if m.k[q]&atom == atom { // _[t;c;`b;a]
			g, q = val(g), val(q)
		} else { // _[t;c;`b!(:e);a]
			g, q = dex(g, val(inc(m.k[3+g]))), dex(q, val(inc(m.k[3+q])))
		}
		for i := k(0); i < n; i++ {
			ri := inc(m.k[2+i+q])
			ki := inc(m.k[2+m.k[2+g+i]])
			vi := inc(m.k[3+m.k[2+g+i]])
			t = dmdv(t, l2(ri, ki), inc(null), vi)
		}
		return decr(g, q, t)
	}
	dec(b)
	if m.k[c]&atom != 0 {
		c = wer(atx(inc(t), c))
	} else {
		c = dex(c, jota(m.k[t]&atom))
	}

	v := atx(inc(t), a)
	vt, _ := typ(v)
	if vt != A {
		panic("type")
	}

	rn := m.k[t] & atom
	r = mk(A, m.k[t]&atom)
	m.k[2+r] = inc(m.k[2+t])
	m.k[3+r] = mk(L, m.k[m.k[3+t]]&atom)
	for i := k(0); i < m.k[m.k[3+t]]&atom; i++ {
		m.k[2+i+m.k[3+r]] = inc(m.k[2+i+m.k[3+t]])
	}

	kv := m.k[2+v]
	nk := m.k[kv] & atom
	cols, u := m.k[3+r], k(0)
	for i := k(0); i < nk; i++ {
		j := fnk(m.k[2+r], m.k[2+i+kv])
		u = atx(inc(m.k[3+v]), mki(i))
		if j == nk {
			m.k[2+r] = cat(m.k[2+r], atx(inc(kv), mki(i)))
			m.k[3+r] = lcat(m.k[3+r], take(rn, 0, mk(m.k[u]>>28, 0))) // nulls
			nk++
		}
		m.k[2+j+cols] = amdv(m.k[2+j+cols], inc(c), inc(null), u)
	}
	dec(v)
	m.k[3+r] = cols
	return decr(t, c, r)
}
func ins(x, y, f, z k) (r k) { // ?[x;y;f;z] splice
	if !match(f, null) {
		panic("nyi")
	}
	dec(f)
	if t, n := typ(y); t == I && n == atom {
		y = enl(cat(y, mki(0)))
	} else if t == I && n == 2 {
		y = enl(y)
	} else if t != L {
		panic("type")
	}
	xt, zt, xn, zn := typs(x, z)
	if xt == L && zt < L {
		z = explode(z)
		zt, zn = typ(z)
	}
	if xt != zt || xt > L || xn == atom {
		panic("type")
	}
	yn, zn := m.k[y]&atom, atm1(zn)
	rn := xn + yn*zn
	for i := k(0); i < yn; i++ {
		yi := m.k[2+y+i]
		t, n := typ(yi)
		if t != I || n != 2 {
			panic("type")
		}
		if m.k[2+yi] > xn || m.k[3+yi] > xn+1 {
			panic("length")
		}
		rn -= m.k[3+yi]
	}
	r = mk(xt, rn)
	cp, xp, zp, rp, x0 := cpx[xt], ptr(x, xt), ptr(z, zt), ptr(r, xt), k(0)
	for i := k(0); i < yn; i++ {
		yi := m.k[2+y+i]
		y0, dy := m.k[2+yi]-x0, m.k[3+yi]
		for j := k(0); j < y0; j++ {
			cp(rp+j, xp+j)
		}
		rp, xp, x0 = rp+y0, xp+y0+dy, y0+dy
		for j := k(0); j < zn; j++ {
			cp(rp+j, zp+j)
		}
		rp, xn = rp+zn, xn-y0-dy
	}
	for j := k(0); j < xn; j++ {
		cp(rp+j, xp+j)
	}
	dec(y)
	return decr(x, z, r)
}
func insert(x, y, idx k) (r k) { // insert y into x at k
	t, yt, n, yn := typs(x, y)
	if t > L || n == atom || (t != L && (t != yt || yn != atom)) {
		panic("type")
	}
	cp, xp, yp := cpx[t], ptr(x, t), ptr(y, t)
	if m.k[x+1] == 1 && bk(t, n) == bk(t, n+1) {
		if t == L {
			cp = cpI
		}
		m.k[x] = t<<28 | (n + 1)
		for i := n; i > idx; i-- {
			cp(xp+i, xp+i-1)
		}
		if t == L {
			m.k[2+x+idx] = y
		} else {
			cp(xp+idx, yp)
			dec(y)
		}
		return x
	}
	r = mk(t, n+1)
	rp := ptr(r, t)
	for i := k(0); i < n+1; i++ {
		if i == idx {
			if t == L {
				m.k[2+idx+r] = y
			} else {
				cp(rp+i, yp)
				dec(y)
			}
		} else {
			cp(rp+i, xp)
			xp++
		}
	}
	return dex(x, r)
}
func unsert(x, idx k) (r k) { // delete index from x
	t, n := typ(x)
	if t == atom || t > L || idx >= n {
		panic("type")
	}
	cp, xp := cpx[t], ptr(x, t)
	if m.k[x+1] == 1 && bk(t, n-1) == bk(t, n) {
		if t == L {
			cp = cpI
			dec(m.k[2+x+idx])
		}
		for i := k(idx); i < n-1; i++ {
			cp(xp+i, xp+i+1)
		}
		m.k[x] = t<<28 | (n - 1)
		return x
	}
	r = mk(t, n-1)
	rp := ptr(r, t)
	for i := k(0); i < idx; i++ {
		cp(rp+i, xp+i)
	}
	for i := k(idx) + 1; i < n; i++ {
		cp(rp+i-1, xp+i)
	}
	return dex(x, r)
}
func asn(x, y, f k) (r k) { // `x:y
	_, yt, xn, yn := typs(x, y)
	if xn != atom {
		if yn == atom {
			y, yn = ext(y, yt, xn), xn
		}
		if yn != xn {
			panic("length")
		}
		for i := k(0); i < xn; i++ {
			dec(asn(atx(inc(x), mki(i)), atx(inc(y), mki(i)), inc(f)))
		}
		dec(f)
		return decr(x, y, inc(null))
	}
	if match(f, null) == false {
		y = cal2(f, lup(inc(x)), y)
	} else {
		dec(f)
	}
	keys, vals := m.k[kkey], m.k[kval]
	if ix, exists := varn(ptr(x, S)); exists {
		dec(m.k[2+vals+ix])
		m.k[2+vals+ix] = inc(y)
		return dex(x, y)
	} else {
		c, n := sc(2 + x)
		dot := false
		for i := k(0); i < n; i++ {
			if m.c[c+i] == '.' {
				dot = true
				break
			}
		}
		if !dot {
			m.k[kkey] = insert(keys, x, ix)
			m.k[kval] = insert(vals, inc(y), ix)
			return y
		}
		x = spl(mkc('.'), str(x))
		s, p, v := c2s(fst(inc(x))), mk(S, 0), k(0)
		for i := k(1); i < m.k[x]&atom; i++ {
			p = cat(p, c2s(inc(m.k[2+i+x])))
		}
		ix, exists = varn(2 + s)
		if exists {
			v = m.k[2+vals+ix]
		} else {
			v = key(mk(S, 0), mk(L, 0))
		}
		v = das(v, p, inc(y))
		if !exists {
			m.k[kkey] = insert(keys, s, ix)
			m.k[kval] = insert(vals, v, ix)
		} else {
			dec(s)
			m.k[2+vals+ix] = v
		}
		return dex(x, y)
	}
}
func das(d, p, y k) (r k) { // dot assign y to dict at path
	dt, dn := typ(d)
	if dt != A || dn != atom {
		panic("value")
	}
	kv, vv := m.k[2+d], m.k[3+d]
	kt, vt, kn, _ := typs(kv, vv)
	if kt != S || vt != L {
		panic("value")
	}
	h := fst(inc(p))
	j := fnd(inc(kv), inc(h))
	j, p = dex(j, m.k[2+j]), drop(1, p)
	v, rn := k(0), kn
	if j == kn {
		v, rn = key(mk(S, 0), mk(L, 0)), rn+1
	} else {
		v = inc(m.k[2+j+vv])
	}
	rk := unq(cat(inc(kv), h)) // TODO rc1 use d
	rv := mk(L, rn)
	for i := k(0); i < rn; i++ {
		if i == j {
			if m.k[p]&atom == 0 {
				m.k[2+rv+i] = decr(v, p, y)
			} else {
				m.k[2+rv+i] = das(v, p, y)
			}
		} else {
			m.k[2+rv+i] = inc(m.k[2+vv+i])
		}
	}
	r = mk(A, atom)
	m.k[2+r], m.k[3+r] = rk, rv
	return dex(d, r)
}
func dxn(name, idx k) (k, k) { // `a.b[`c] → `a `b`c
	c, n := sc(2 + name)
	for i := k(0); i < n; i++ {
		if m.c[c+i] == '.' {
			break
		} else if i == n-1 {
			return name, idx
		}
	}
	x := spl(mkc('.'), str(name))
	name = c2s(fst(inc(x)))
	x = drop(1, x)
	r := mk(S, m.k[x]&atom)
	for i := k(0); i < m.k[x]&atom; i++ {
		s := c2s(inc(m.k[2+i+x]))
		m.k[2+r+i] = dex(s, m.k[2+s])
	}
	return dex(x, name), cat(r, idx)
}
func mut(x, y k) { // modify-inplace
	if ix, exists := varn(ptr(x, S)); !exists {
		panic(undef(x))
	} else {
		m.k[2+ix+m.k[kval]] = y
		dec(x)
	}
}
func amd(x, a, f, y k) (r k) { // @[x;i;f;y]
	t, n := typ(x)
	if t != S {
		return amdv(x, a, f, y)
	} else if n != atom {
		panic("type")
	}
	v := lup(inc(x))
	if m.k[1+v] == 2 { // in-place
		dec(v)
		mut(inc(x), amdv(v, a, f, y))
		return x
	}
	dec(asn(inc(x), amdv(v, a, f, y), inc(null)))
	return x
}
func dmd(x, a, f, y k) (r k) { // .[x;i;f;y]
	t, n := typ(x)
	if t != S {
		return dmdv(x, a, f, y)
	} else if n != atom {
		panic("type")
	}
	v := lup(inc(x))
	if m.k[1+v] == 2 {
		dec(v)
		mut(inc(x), dmdv(v, a, f, y))
		return x
	}
	dec(asn(inc(x), dmdv(v, a, f, y), inc(null)))
	return x
}
func amdv(x, a, f, y k) (r k) { // amd on value(x)
	xt, at, xn, an := typs(x, a)
	if xt == A {
		if m.k[m.k[2+x]]&atom == 0 {
			return decr(x, f, key(a, y))
		}
		if xn != atom && at == I { // t[I]:d t[I]:t
			yt, _ := typ(y)
			if yt != A || !match(m.k[2+x], m.k[2+y]) {
				panic("type")
			}
			r = mk(A, xn)
			nk := m.k[m.k[2+x]] & atom
			m.k[2+r] = inc(m.k[2+x])
			m.k[3+r] = mk(L, nk)
			for i := k(0); i < nk; i++ {
				m.k[2+i+m.k[3+r]] = amdv(inc(m.k[2+i+m.k[3+x]]), inc(a), inc(f), atx(inc(m.k[3+y]), mki(i)))
			}
			return decr(x, decr(a, f, y), r)
		}
		r = mk(A, atom)
		m.k[2+r] = inc(m.k[2+x])
		m.k[3+r] = inc(m.k[3+x])
		al := inc(a)
		if m.k[al]&atom == atom {
			al = enl(a)
		}
		idx, n := fnd(inc(m.k[2+x]), inc(a)), m.k[m.k[2+x]]&atom // i:(!x)?`a
		u := unq(atx(al, wer(eql(mki(n), idx))))                 // ?(!x)@&(#!x)=i
		if m.k[u]&atom > 0 {
			e := take(m.k[u]&atom, 0, mk(m.k[m.k[3+r]]>>28, 0))
			m.k[2+r] = cat(m.k[2+r], u)
			m.k[3+r] = cat(m.k[3+r], e)
		} else {
			dec(u)
		}
		m.k[3+r] = amdv(m.k[3+r], fnd(inc(m.k[2+r]), a), f, y)
		return dex(x, r)
	}
	if an == 0 {
		dec(f)
		return decr(a, y, x)
	} else if at != I {
		panic("type")
	}
	if m.k[f]>>28 != N {
		if match(y, null) {
			dec(y)
			y = cal(f, enl(atx(inc(x), inc(a))))
		} else {
			y = cal2(f, atx(inc(x), inc(a)), y)
		}
	} else {
		dec(f)
	}
	yt, yn := typ(y)
	if an == atom && yn != atom { // replace, e.g. x[1]:2 3
		if xt < L {
			x, xt = explode(x), L
			xn = m.k[x] & atom
		} else if xt != L {
			panic("type")
		}
		j := m.k[2+a]
		if int32(j) < 0 || j >= xn {
			panic("index")
		}
		if m.k[x+1] == 1 {
			dec(m.k[2+j+x])
			m.k[2+j+x] = y
			return dex(a, x)
		}
		r = mk(L, xn)
		for i := k(0); i < xn; i++ {
			if i == j {
				m.k[2+i+r] = y
			} else {
				m.k[2+i+r] = inc(m.k[2+i+x])
			}
		}
		return decr(x, a, r)
	}
	if an != atom && yn == atom {
		y, yn = ext(y, yt, an), an
	}
	if an != yn {
		panic("length")
	}
	if m.k[x+1] == 1 { // in-place
		if xt == L {
			if yn == atom {
				y, yn = enl(y), 1
			}
			for i := k(0); i < yn; i++ {
				j := m.k[2+i+a]
				if int32(j) < 0 || j >= xn {
					panic("index")
				}
				dec(m.k[2+j+x])
				m.k[2+j+x] = atx(inc(y), mki(i))
			}
			x = uf(x)
		} else if xt >= L || xt != yt {
			panic("type")
		} else {
			xp, yp, cp := ptr(x, xt), ptr(y, yt), cpx[xt]
			for i := k(0); i < atm1(yn); i++ {
				j := m.k[2+i+a]
				if int32(j) < 0 || j >= xn {
					panic("index")
				}
				cp(xp+j, yp+i)
			}
		}
		return decr(a, y, x)
	}
	r = mk(xt, xn)
	if xt == L {
		for i := k(0); i < xn; i++ {
			m.k[2+i+r] = inc(m.k[2+i+x])
		}
		if yn == atom {
			y, yn = enl(y), 1
		}
		for i := k(0); i < yn; i++ {
			j := m.k[2+i+a]
			if int32(j) < 0 || j >= xn {
				panic("index")
			}
			dec(m.k[2+j+r])
			m.k[2+j+r] = atx(inc(y), mki(i))
		}
		if xn != 1 || m.k[m.k[2+r]]>>28 != A || m.k[m.k[2+r]]&atom != atom {
			r = uf(r) // don't unify (d)
		}
	} else if xt >= L || xt != yt {
		panic("type")
	} else {
		mv(r, x)
		rp, yp, cp := ptr(r, xt), ptr(y, yt), cpx[xt]
		yn = atm1(yn)
		for i := k(0); i < yn; i++ {
			j := m.k[2+i+a]
			if int32(j) < 0 || j >= xn {
				panic("index")
			}
			cp(rp+j, yp+i)
		}
	}
	dec(a)
	return decr(x, y, r)
}
func dmdv(x, a, f, y k) (r k) { // dmd on value(x)
	t, xt, n, xn := typs(a, x)
	if n == atom {
		return amdv(x, a, f, y)
	} else if n == 1 {
		return amdv(x, fst(a), f, y)
	} else if n == 0 {
		panic("domain")
	} else if n == 2 && t == L && xt == A && xn != atom { // t[I;S]..
		idx := l2(fnd(inc(m.k[2+x]), inc(m.k[3+a])), inc(m.k[2+a]))
		if match(m.k[2+a], null) {
			dec(m.k[3+idx])
			m.k[3+idx] = jota(xn)
		}
		r = mk(A, xn)
		m.k[2+r] = inc(m.k[2+x])
		m.k[3+r] = dmdv(inc(m.k[3+x]), idx, f, y)
		return decr(x, a, r)
	} else if n == 2 && t == L && (m.k[m.k[2+a]]&atom != atom || match(m.k[2+a], null)) { // matrix assign
		a0, a1, yi := inc(m.k[2+a]), inc(m.k[3+a]), k(0)
		dec(a)
		if match(a0, null) {
			a0 = dex(a0, jota(m.k[x]&atom))
		}
		for i := k(0); i < m.k[a0]&atom; i++ {
			ai := atx(inc(a0), mki(i))
			yi = inc(null)
			if m.k[y]>>28 != N {
				yi = dex(null, atx(inc(y), mki(i)))
			}
			x = amdv(inc(x), inc(ai), inc(null), amdv(atx(x, ai), inc(a1), inc(f), yi))
		}
		return decr(f, decr(a0, a1, y), x)
	}
	a0 := fst(inc(a))
	return amdv(inc(x), inc(a0), inc(null), dmdv(atx(x, a0), drop(1, a), f, y))
}
func lup(x k) (r k) { // lookup
	if r = lupo(inc(x)); r == 0 {
		panic(undef(x))
	}
	return dex(x, r)
}
func lupo(x k) (r k) { // lup, 0 on undefined
	ix, o := varn(ptr(x, S))
	if !o { // try a.b.c
		x = spl(mkc('.'), str(x))
		n := m.k[x] & atom
		if n < 2 {
			return dex(x, 0)
		}
		r = lupo(c2s(inc(m.k[2+x])))
		if r == 0 || m.k[r]>>28 != A {
			return dex(x, 0)
		}
		for i := k(1); i < n; i++ {
			kv, vv, nk, s := m.k[2+r], m.k[3+r], m.k[m.k[2+r]]&atom, k(0)
			if m.k[kv]>>28 != S || m.k[vv]>>28 != L {
				panic("value")
			}
			s = c2s(inc(m.k[2+x+i]))
			s = dex(s, m.k[2+s])
			for j := k(0); j < nk; j++ {
				if m.k[2+j+kv] == s {
					r = dex(r, inc(m.k[2+j+vv]))
					break
				} else if j == nk-1 {
					return decr(x, r, 0)
				}
			}
		}
		return dex(x, r)
	}
	vals := m.k[kval]
	r = inc(m.k[2+vals+ix])
	return dex(x, r)
}
func undef(x k) (r s) {
	c, n := sc(2 + x)
	dec(x)
	return "undefined:" + s(m.c[c:c+n])
}
func varn(xp k) (idx k, exists bool) {
	keys := m.k[kkey]
	kp := ptr(keys, S)
	kn := m.k[keys] & atom
	ix := ibin(kp, kn, xp, gtS)
	if ix < kn && eqS(kp+ix, xp) {
		return ix, true
	}
	return ix + 1, false
}
func vars(dummy k) (r k) { dec(dummy); return inc(m.k[kkey]) }
func del(x k) (r k) { // delete variable
	t, n := typ(x)
	if t != S {
		panic("type")
	}
	n = atm1(n)
	xp := ptr(x, S)
	for i := k(0); i < n; i++ {
		if idx, o := varn(xp + i); o {
			m.k[kkey] = unsert(m.k[kkey], idx)
			m.k[kval] = unsert(m.k[kval], idx)
		}
	}
	return dex(x, inc(null))
}
func clear() { // clear variables
	dec(m.k[kkey])
	dec(m.k[kval]) // m.k[kkey] and m.k[kval] should not be free blocks
	m.k[kkey] = mk(S, 0)
	m.k[kval] = mk(L, 0)
}
func lsv() (r k)     { return inc(m.k[kkey]) }                                  // \v (list variables)
func clv() (r k)     { clear(); return inc(null) }                              // \c (clear variables)
func hlp() (r k)     { return cat(mkb(m.c[136:168]), kst(inc(m.k[2+m.k[3]]))) } // \h
func hxb(x c) (c, c) { h := "0123456789abcdef"; return h[x>>4], h[x&0x0F] }
func hxk(x k) s {
	b := []c{'0', 'x', '0', '0', '0', '0', '0', '0', '0', '0'}
	for j := k(0); j < 4; j++ {
		n := 8 * (3 - j)
		b[2+2*j], b[3+2*j] = hxb(c((x & (0xFF << n)) >> n))
	}
	return s(b)
}
func prm(x k) (r k) { // prm x (permutations)
	t, n := typ(x)
	if t < L && n == atom {
		if idx(x, t) < 0 {
			panic("type")
		}
		x = til(x)
		n = m.k[x] & atom
	} else if n == atom || t > Z {
		panic("type")
	}
	z := mk(t, n)
	mv(z, x)
	return decr(x, z, heap(z, n, mk(L, 0)))
}
func heap(x, n, r k) k {
	t, xn := typ(x)
	if n == 1 {
		ri := mk(t, xn)
		mv(ri, x)
		return lcat(r, ri)
	}
	y := mk(t, atom)
	xp, yp, cp := ptr(x, t), ptr(y, t), cpx[t]
	for i := k(0); i < n; i++ {
		r = heap(x, n-1, r)
		xpi := xp + i
		if n%2 == 1 {
			xpi = xp
		}
		cp(yp, xpi)
		cp(xpi, xp+n-1)
		cp(xp+n-1, yp)
	}
	dec(y)
	return r
}
func rol(x, y k) (r k) { // x rand y (roll, deal)
	xt, yt, xn, yn := typs(x, y)
	if xt != I || xn != atom {
		panic("type")
	}
	n := i(m.k[2+x])
	if yn == atom {
		if yt == I && m.k[2+x] > 0 && m.k[2+y] > 0 {
			rn, n := m.k[2+x], 1+m.k[2+y]
			r = mk(I, rn)
			for i := k(0); i < rn; i++ {
				m.k[2+r+i] = rndn(n)
			}
			return decr(x, y, r)
		} else if yt == I { // n rand m → n rand !m
			yn = m.k[2+y]
			y = til(y)
		} else if yt == F { // n rand yf uniform [0,y] #n
			if n < 0 {
				panic("type")
			}
			r = rnd(x)
			rp := ptr(r, F)
			c := m.f[ptr(y, F)]
			for i := k(0); i < m.k[r]&atom; i++ {
				m.f[rp+i] *= c
			}
			return dex(y, r)
		} else if yt == C { // n rand "A"
			dec(y)
			y, yt, yn = mk(C, 26), C, 26
			for i := k(0); i < yn; i++ {
				m.c[8+i+y<<2] = c('A' + i)
			}
		}
	}
	if n < 0 { // -n rand y (draw -n, no repetitions)
		n = -n
		if k(n) > yn {
			panic("type")
		}
		r = mk(I, k(n))
		for i := k(0); i < k(n); i++ {
			m.k[2+r+i] = 888
		}
		for i := k(0); i < k(n); i++ { // fisher-yates (inside-out)
			j := rndn(i + 1)
			if i != j {
				m.k[2+r+i] = m.k[2+r+j]
			}
			m.k[2+r+j] = i
		}
		return dex(x, atx(y, r))
	} else { // n rand y (draw n with repetitions)
		r = mk(I, k(n))
		for i := k(0); i < k(n); i++ {
			m.k[2+r+i] = rndn(yn)
		}
		return dex(x, atx(y, r))
	}
	return decr(x, y, r)
}
func rnd(x k) (r k) { // rand x
	t, xn := typ(x)
	if xn != atom { // draw a random element from x
		return atx(x, mki(rndn(xn)))
	} else if t == Z { // rand Ni: binormal #N
		n := idx(x, Z)
		if n < 0 {
			panic("type")
		}
		r = mk(Z, k(n))
		rc := 2 + r>>1
		for i := k(0); i < 2*k(n); i += 2 {
			m.f[rc+i], m.f[rc+1+i] = normal()
		}
	} else if t == I {
		n := i(m.k[2+x])
		if n < 0 { // rand -N: normal #N
			n = -n
			r = mk(F, k(n))
			rp := 1 + r>>1
			m.f[rp], _ = normal()
			rp++
			for i := k(0); i < k(n-1); i += 2 {
				m.f[rp+i], m.f[rp+i+1] = normal()
			}
		} else { // rand N: uniform [0, 1] #N
			r = mk(F, k(n))
			rp := 1 + r>>1
			for i := k(0); i < k(n); i++ {
				m.f[rp+i] = f(rng()) / f(0xFFFFFFFF)
			}
		}
	} else {
		panic("type")
	}
	return dex(x, r)
}
func rng() (r k) { // xor-shift
	r = m.k[1]
	r ^= (r << 13)
	r ^= (r >> 17)
	r ^= (r << 5)
	m.k[1] = r
	return r
}
func rndn(n k) (r k) { // random [0,n) math/rand/rand.go:int32n (lemire)
	r = rng()
	p := uint64(r) * uint64(n)
	if l := k(p); l < k(n) {
		t := k(-i(n)) % n
		for l < t {
			r = rng()
			p = uint64(r) * uint64(n)
			l = k(p)
		}
	}
	return k(p >> 32)
}
func normal() (f, f) { // marsaglia polar
	var u, v, s f
	for s == 0 || s >= 1 {
		u = (f(rng())/f(0xFFFFFFFF))*2.0 - 1.0
		v = (f(rng())/f(0xFFFFFFFF))*2.0 - 1.0
		s = u*u + v*v
	}
	s = math.Sqrt(-2.0 * math.Log(s) / s)
	return u * s, v * s
}
func norm(xp, n k, sqrt bool) (r f) { // vector norm L2
	s := 0.0
	for i := xp; i < xp+n; i++ {
		if x := m.f[i]; x != 0 {
			if x = math.Abs(x); math.IsNaN(x) || math.IsInf(x, 1) {
				return x
			} else if s < x {
				t := s / x
				r, s = 1+r*t*t, x
			} else {
				t := x / s
				r += t * t
			}
		}
	}
	if sqrt {
		return s * math.Sqrt(r)
	}
	return s * s * r
}
func dot(x, y k) (r k) { // x@y x/y (matrix multiplication)
	//    v/  v →  a       n/n   → (atom)
	// (,v)/  v → ,a     1 n/n   → 1
	// (,v)/+,v →,,a     1 n/n 1 → 1 1
	//    v/  m →  v       n/n m → m
	//    m/  v →  v     m n/n   → m
	//    m/+,v →+,v     m n/n 1 → m 1
	//(+,v)/ ,v →  m     n 1/1 n → n n
	//    m/  m →  m     m n/n r → m r
	xt, yt, xn, yn := typs(x, y)
	if yt == L {
		y = flp(y)
		yn = m.k[y] & atom
	}
	switch {
	case xt > L || yt > L:
		panic("type")
	case xt < L && yt < L: // v/v → a
		r = mk(xt, atom)
		rp := ptr(r, xt)
		vdot(rp, x, y, xt)
	case xt == L && yt < L: // (,v)/v | m/v
		r = mk(yt, xn)
		rp := ptr(r, yt)
		for i := k(0); i < xn; i++ {
			vdot(rp+i, m.k[2+x+i], y, yt)
		}
	case xt < L && yt == L: // v/m
		r = mk(xt, yn)
		rp := ptr(r, xt)
		for i := k(0); i < yn; i++ {
			vdot(rp+i, x, m.k[2+y+i], xt)
		}
	case xt == L && yt == L:
		r = mk(L, xn)
		t := m.k[m.k[2+x]] >> 28
		for i := k(0); i < xn; i++ {
			rr := mk(t, yn)
			rp := ptr(rr, t)
			for j := k(0); j < yn; j++ {
				vdot(rp+j, m.k[2+i+x], m.k[2+j+y], t)
			}
			m.k[2+r+i] = rr
		}
	}
	return decr(x, y, r)
}
func vdot(rp, x, y, t k) { // x+/y for vectors
	xt, yt, xn, yn := typs(x, y)
	if xt != t || yt != t || xt > Z {
		panic("type")
	} else if xn != yn || xn == atom {
		panic("size")
	}
	xp, yp := ptr(x, t), ptr(y, t)
	switch t {
	case I:
		s := k(0)
		for i := k(0); i < xn; i++ {
			s += m.k[xp+i] * m.k[yp+i]
		}
		m.k[rp] = s
	case F:
		s := 0.0
		for i := k(0); i < xn; i++ {
			s += m.f[xp+i] * m.f[yp+i]
		}
		m.f[rp] = s
	case Z:
		s := 0i
		for i := k(0); i < xn; i++ {
			s += m.z[xp+i] * m.z[yp+i]
		}
		m.z[rp] = s
	}
}
func qrd(x k) (r k) { // x\0 (qr decomposition)
	lt, rows := typ(x)
	if lt != L {
		panic("type")
	} else if rows == 0 || rows == atom {
		panic("empty")
	}
	t, cols := typ(m.k[2+x])
	if cols == atom || (t != F && t != Z) {
		panic("type")
	}
	h, d := mk(t, cols*rows), mk(t, cols) // h: qr compact storage of Q and R without diag, d: diag R
	hp, hpr, dp, cp := ptr(h, t), k(0), ptr(d, t), cpx[t]
	for i := k(0); i < rows; i++ {
		xp := ptr(m.k[2+i+x], t)
		for j := k(0); j < cols; j++ {
			cp(hp+rows*j+i, xp+j) // h: transpose of x
		}
	}
	var s, a f
	for j := k(0); j < cols; j++ {
		hpj := hp + j*rows // H size: cols x rows!
		if t == F {
			s = norm(hpj+j, rows-j, true) // H[j][j:]
			if m.f[hpj+j] > 0 {
				m.f[dp+j] = -s
			} else {
				m.f[dp+j] = s
			}
			a = 1.0 / math.Sqrt(s*(s+math.Abs(m.f[hpj+j])))
			m.f[hpj+j] -= m.f[dp+j]
			for k := j; k < rows; k++ {
				m.f[hpj+k] *= a
			}
			for i := j + 1; i < cols; i++ {
				hpi, ss := hp+i*rows, 0.0
				for k := j; k < rows; k++ {
					ss += m.f[hpj+k] * m.f[hpi+k]
				}
				for k := j; k < rows; k++ {
					m.f[hpi+k] -= ss * m.f[hpj+k]
				}
			}
		} else {
			hpr = hpj << 1
			re, im := m.f[hpr+2*j], m.f[hpr+2*j+1]
			si, co := math.Sincos(math.Atan2(im, re))
			s = norm(hpr+2*j, 2*(rows-j), true)
			m.z[dp+j] = complex(-s*co, -s*si)
			a = 1.0 / math.Sqrt(s*(s+math.Hypot(re, im)))
			m.z[hpj+j] -= m.z[dp+j]
			for k := j; k < rows; k++ {
				m.f[hpr+2*k] *= a
				m.f[hpr+2*k+1] *= a
			}
			for i := j + 1; i < cols; i++ {
				hpi, ss := hp+i*rows, 0i
				for k := j; k < rows; k++ {
					ss += conj(m.z[hpj+k]) * m.z[hpi+k]
				}
				for k := j; k < rows; k++ {
					m.z[hpi+k] -= ss * m.z[hpj+k]
				}
			}
		}
		if s == 0 {
			panic("singular")
		}
	}
	r = mk(L, 4)
	m.k[2+r], m.k[3+r], m.k[4+r], m.k[5+r] = mki(rows), mki(cols), h, d
	return dex(x, r)
}
func cnd(x k) (r k) { // cond x (using max row sum of R from QR)
	// matlab: for A (rows x cols), cond A should be equal to: [q, r] = qr(A); cond(r(1:rows,1:rows), inf)
	if m.k[x]&atom != 4 || m.k[m.k[2+x]]&atom != atom {
		x = qrd(x)
	} // else assume input is qr
	rows, cols, h, d := m.k[2+m.k[2+x]], m.k[2+m.k[3+x]], m.k[4+x], m.k[5+x]
	if rows*cols != m.k[h]&atom {
		panic("type")
	}
	t, ln := m.k[h]>>28, (cols*(cols+1))/2
	r = mk(t, ln)
	dp, hp, rp, cp, n := ptr(d, t), ptr(h, t), ptr(r, t), cpx[t], k(0)
	for i := k(0); i < cols; i++ { // store triangular matrix r
		cp(rp+n, dp+i)
		n++
		for k := i + 1; k < cols; k++ {
			cp(rp+n, hp+cols*k+i)
			n++
		}
	}
	n1 := trn(r, t, cols) // cond: norm(inv R) * n1:norm(R)
	z, o := mk(t, atom), mk(t, atom)
	if t == F {
		m.f[1+z>>1], m.f[1+o>>1] = 0, 1
	} else {
		m.z[1+z>>2], m.z[1+o>>2] = 0, 1
	}
	e := mk(t, cols)
	ep, zp, op := ptr(e, t), ptr(z, t), ptr(o, t)
	for i := k(0); i < cols; i++ { // solve R x = I
		for j := k(0); j < cols; j++ {
			cp(ep+j, zp)
		}
		cp(ep+i, op)
		rsv(hp, dp, ep, rows, cols, t)
		n = i
		for j := k(0); j < i+1; j++ {
			cp(rp+n, ep+j)
			n += cols - j - 1
		}
	}
	dec(e)
	dec(z)
	dec(o)
	n2 := trn(r, t, cols) // cond (inv R)
	dec(r)
	r = mk(F, atom)
	m.f[1+r>>1] = n1 * n2
	return dex(x, r)
}
func trn(r, t, n k) f { // inf-norm of triangular matrix
	mx, rp := 0.0, ptr(r, F)
	for i := k(0); i < n; i++ {
		s := 0.0
		if t == F {
			for j := k(0); j < n-i; j++ {
				s += math.Abs(m.f[rp+j])
			}
			rp += n - i
		} else if t == Z {
			for j := k(0); j < 2*(n-i); j++ {
				s += math.Abs(m.f[1+rp+j])
			}
			rp += 2 * (n - i)
		}
		if s > mx {
			mx = s
		}
	}
	return mx
}
func slv(x, y k) (r k) { // x\y (solve)
	xt, yt, xn, yn := typs(x, y)
	if xt != L {
		panic("type")
	} else if yt == I && yn == atom {
		if n := m.k[2+y]; n == 0 { // qr: x\0
			return dex(y, qrd(x))
		} else if n == 1 { // inv: x\1
			n := m.k[m.k[2+x]] & atom
			y, yn = eye(n), n
		}
	}
	x0 := m.k[2+x]
	if m.k[x0]&atom != atom {
		x = qrd(x)
		x0, xn = m.k[2+x], 4
	} else if xn != 4 {
		panic("type") // x is no qr
	}
	rows := m.k[2+x0]
	cols := m.k[2+m.k[3+x]]
	h, d := m.k[4+x], m.k[5+x]
	if yt == L {
		y = flp(y)
		yn = m.k[y] & atom
		r = mk(L, yn)
		for i := k(0); i < atm1(yn); i++ {
			m.k[2+r+i] = qrs(rows, cols, h, d, inc(m.k[2+y+i]))
		}
		return decr(x, y, r)
	}
	r = qrs(rows, cols, h, d, y)
	return dex(x, r)
}
func qrs(rows, cols, h, d, y k) (r k) {
	t, ht, n, _ := typs(y, h)
	if ht < t || t > Z {
		panic("type")
	} else if t < ht {
		y, t = to(y, ht), ht
	}
	if rows != n {
		panic("size")
	}
	if m.k[1+y] != 1 {
		r = mk(t, n)
		mv(r, y)
		dec(y)
	} else {
		r = y
	}
	rp, hp, dp := ptr(r, t), ptr(h, t), ptr(d, t)
	for i := k(0); i < cols; i++ { // multiply: r = Q^T y
		if t == F {
			var s f
			hpi := ptr(h, F) + i*rows
			for k := i; k < rows; k++ {
				s += m.f[hpi+k] * m.f[rp+k]
			}
			for k := i; k < rows; k++ {
				m.f[rp+k] -= m.f[hpi+k] * s
			}
		} else {
			var s z
			hpi := ptr(h, Z) + i*rows
			for k := i; k < rows; k++ {
				s += conj(m.z[hpi+k]) * m.z[rp+k]
			}
			for k := i; k < rows; k++ {
				m.z[rp+k] -= m.z[hpi+k] * s
			}
		}
	}
	rsv(hp, dp, rp, rows, cols, t)
	return srk(r, t, n, cols)
}
func rsv(hp, dp, rp, rows, cols, t k) { // solve R x = y
	for i := cols - 1; ; i-- { // back-substitution
		if t == F {
			for j := i + 1; j < cols; j++ {
				m.f[rp+i] -= m.f[hp+j*rows+i] * m.f[rp+j]
			}
			m.f[rp+i] /= m.f[dp+i]
		} else {
			for j := i + 1; j < cols; j++ {
				m.z[rp+i] -= m.z[hp+j*rows+i] * m.z[rp+j]
			}
			m.z[rp+i] /= m.z[dp+i]
		}
		if i == 0 {
			break
		}
	}
}
func conj(x z) z { return complex(real(x), -imag(x)) }
func mkz(x, y k) (r k) { // x cmplx y
	xt, yt, xn, yn := typs(x, y)
	n, dx, dy := atom, k(1), k(1)
	if xn == atom {
		dx, n = 0, yn
	} else if yn == atom {
		dy, n = 0, xn
	} else {
		n = xn
	}
	if dx+dy == 2 && xn != yn {
		panic("size")
	}
	if xt != F {
		x = to(x, F)
	}
	if yt != F {
		y = to(y, F)
	}
	r = mk(Z, n)
	xp, yp, rp := ptr(x, F), ptr(y, F), ptr(r, Z)
	for i := k(0); i < atm1(n); i++ {
		m.z[rp+i] = complex(m.f[xp], m.f[yp])
		xp += dx
		yp += dy
	}
	return decr(x, y, r)
}
func sum(x k) (r k) {
	t, n := typ(x)
	if n == atom && t != A {
		return x
	}
	switch t {
	case C:
		s := c(0)
		for i := ptr(x, C); i < ptr(x, C)+atm1(n); i++ {
			s += m.c[i]
		}
		r = mkc(s)
	case I:
		s := k(0)
		for i := ptr(x, I); i < ptr(x, I)+atm1(n); i++ {
			s += m.k[i]
		}
		r = mki(s)
	case F:
		r = mk(F, atom)
		m.f[ptr(r, F)] = fsum(ptr(x, F), n)
	case Z:
		r = mk(Z, atom)
		m.z[ptr(r, Z)] = zsum(ptr(x, Z), n)
	default:
		panic("type")
	}
	return dex(x, r)
}
func fsum(xp, n k) (r f) { // pairwise
	if n < 128 {
		for i := k(0); i < n; i++ {
			r += m.f[xp+i]
		}
		return r
	}
	nn := n >> 1
	return fsum(xp, nn) + fsum(xp+nn, n-nn)
}
func zsum(xp, n k) (r z) { // pairwise
	if n < 128 {
		for i := k(0); i < n; i++ {
			r += m.z[xp+i]
		}
		return r
	}
	nn := n >> 1
	return zsum(xp, nn) + zsum(xp+nn, n-nn)
}
func dev(x k) (r k) { // dev x
	t, n := typ(x)
	if t == L {
		return lrc(x, n, dev)
	} else if t == A {
		return arc(x, n, dev)
	} else if n == atom || n < 1 {
		return dex(x, to(mki(0), F))
	}
	if t == Z { // dev z: standard deviations in the principal axes
		c := vri(x)
		cp := ptr(c, F)
		vx, vy, vxy := m.f[cp], m.f[cp+1], m.f[cp+2]
		tr, det := vx+vy, vx*vy-vxy*vxy
		e1, e2 := tr/2+math.Sqrt(tr*tr/4-det), tr/2-math.Sqrt(tr*tr/4-det)
		v1, v2 := complex(e1-vy, vxy), complex(e2-vy, vxy)
		s1, s2 := math.Sqrt(e1), math.Sqrt(e2)
		r = mk(Z, 2)
		rp := ptr(r, Z)
		m.z[rp], m.z[rp+1] = v1*complex(s1/math.Hypot(real(v1), imag(v1)), 0), v2*complex(s2/math.Hypot(real(v2), imag(v2)), 0)
		return dex(c, r)
	}
	r = vri(x)
	m.f[1+r>>1] = math.Sqrt(m.f[1+r>>1])
	return r
}
func vri(x k) (r k) { // var x
	t, n := typ(x)
	if t == L {
		return lrc(x, n, vri)
	} else if t == A {
		return arc(x, n, vri)
	} else if n == atom || n < 1 {
		return dex(x, to(mki(0), F))
	}
	if t < F {
		x = to(x, F)
	} else if t == Z {
		return cov(rel(inc(x)), ima(x))
	} else if t > Z {
		panic("type")
	}
	r = mk(F, atom)
	s2, _ := varf(x, n)
	m.f[1+r>>1] = s2
	return r
}
func cov(x, y k) (r k) { // x var y
	xt, yt, xn, yn := typs(x, y)
	if xt != yt || xn != yn || xn == atom || xn < 2 || xt > F {
		panic("type")
	} else if xt < F {
		x, xt = to(x, F), F
		y, yt = to(y, F), F
	}
	vx, ax := varf(inc(x), xn)
	vy, ay := varf(inc(y), yn)
	if m.k[x+1] != 1 {
		r = mk(F, xn)
		mv(r, x)
		dec(x)
		x = r
	}
	xp, yp := ptr(x, F), ptr(y, F)
	for i := k(0); i < xn; i++ {
		m.f[xp+i] -= ax
		m.f[xp+i] *= m.f[yp+i] - ay
	}
	x = sum(x)
	m.f[ptr(x, F)] /= f(xn - 1)
	r = mk(F, 3)
	m.f[1+r>>1] = vx
	m.f[2+r>>1] = vy
	m.f[3+r>>1] = m.f[1+x>>1]
	return decr(x, y, r)
}
func varf(x, n k) (f, f) { // var, avg
	a := avg(inc(x))
	if m.k[x+1] != 1 {
		r := mk(F, n)
		mv(r, x)
		dec(x)
		x = r
	}
	xp, af, t := ptr(x, F), m.f[ptr(a, F)], 0.0
	for i := k(0); i < n; i++ {
		t = m.f[xp+i] - af
		m.f[xp+i] = t * t
	}
	dec(a)
	a = sum(x)
	s2 := m.f[ptr(a, F)] / f(n-1)
	dec(a)
	return s2, af
}
func avg(x k) (r k) { // avg x
	t, n := typ(x)
	if t == L {
		return lrc(x, n, avg)
	} else if t == A {
		return arc(x, n, avg)
	} else if n == atom {
		return x
	}
	x = sum(x)
	nf := f(n)
	switch t {
	case C:
		r = mk(F, atom)
		m.f[1+r>>1] = f(m.c[ptr(x, C)]) / nf
	case I:
		r = mk(F, atom)
		m.f[1+r>>1] = f(m.k[2+x]) / nf
	case F:
		r = mk(F, atom)
		m.f[1+r>>1] = f(m.f[1+x>>1]) / nf
	case Z:
		r = mk(Z, atom)
		rp, xp := ptr(r, Z)<<1, ptr(x, Z)<<1
		m.f[rp] = f(m.f[xp]) / nf
		m.f[rp+1] = f(m.f[xp+1]) / nf
	default:
		panic("type")
	}
	return dex(x, r)
}
func mvg(x, y k) (r k) { // x avg y
	xt, yt, xn, yn := typs(x, y)
	if yt == L {
		return lrc2(x, y, yn, mvg)
	} else if yt == A {
		return arc2(x, y, yn, mvg)
	} else if yn < 2 || yn == atom {
		return dex(x, y)
	}
	if xn != atom {
		panic("type")
	}
	if yt < F {
		y, yt = to(y, F), F
	} else if yt > Z {
		panic("type")
	}
	if m.k[1+y] == 1 {
		r = y
	} else {
		r = mk(yt, yn)
		mv(r, y)
		dec(y)
	}
	switch xt {
	case I:
		n := idx(x, I)
		if n < 0 {
			panic("value")
		} else if n == 0 { // 0 avg y (cummulative moving average)
			rp := ptr(r, yt)
			if yt == F {
				s := 0.0
				for i := k(0); i < yn; i++ {
					s += (m.f[rp+i] - s) / f(i+1)
					m.f[rp+i] = s
				}
			} else {
				s := 0i
				for i := k(0); i < yn; i++ {
					s += (m.z[rp+i] - s) / complex(f(i+1), 0)
					m.z[rp+i] = s
				}
			}
			return dex(x, r)
		}
		// n avg y: moving average window size n
		b, p := mk(yt, k(n)), k(0)
		bp, rp := ptr(b, yt), ptr(r, yt)
		if yt == F {
			s := 0.0
			for i := k(0); i < yn; i++ {
				s += m.f[rp+i]
				if i < k(n) {
					m.f[bp+p] = m.f[rp+i]
					m.f[rp+i] = s / f(i+1)
				} else {
					s -= m.f[bp+p]
					m.f[bp+p] = m.f[rp+i]
					m.f[rp+i] = s / f(n)
				}
				p++
				if p == k(n) {
					p = 0
				}
			}
		} else {
			s := 0i
			for i := k(0); i < yn; i++ {
				s += m.z[rp+i]
				if i < k(n) {
					m.z[bp+p] = m.z[rp+i]
					m.z[rp+i] = s / complex(f(i+1), 0)
				} else {
					s -= m.z[bp+p]
					m.z[bp+p] = m.z[rp+i]
					m.z[rp+i] = s / complex(f(n), 0)
				}
				p++
				if p == k(n) {
					p = 0
				}
			}
		}
		return decr(x, b, r)
	case F: // exponential moving average
		if yt == F {
			a := m.f[ptr(x, F)]
			b, rp := 1-a, ptr(r, F)
			t := m.f[rp]
			for i := k(0); i < yn; i++ {
				m.f[rp+i], t = a*m.f[rp+i]+b*t, m.f[rp+i]
			}
		} else {
			a := complex(m.f[ptr(x, F)], 0)
			b, rp := 1-a, ptr(r, Z)
			t := m.z[rp]
			for i := k(0); i < yn; i++ {
				m.z[rp+i], t = a*m.z[rp+i]+b*t, m.z[rp+i]
			}
		}
		return dex(x, r)
	}
	panic("type")
}
func med(x k) (r k) { // med x
	t, n := typ(x)
	if t == L {
		return lrc(x, n, med)
	} else if t == A {
		return arc(x, n, med)
	} else if n == atom {
		return x
	}
	x = srt(x)
	return atx(x, mki(n/2))
}
func pct(x, y k) (r k) { // x med y (0.95 med y, -0.95f med y, 0 med y)
	xt, yt, xn, yn := typs(x, y)
	if yt == L {
		return lrc2(x, y, yn, pct)
	} else if yt == A {
		return arc2(x, y, yn, pct)
	} else if yn < 2 || yn == atom {
		return dex(x, y)
	}
	if xn != atom {
		panic("type")
	}
	switch xt {
	case F:
		p := m.f[1+x>>1]
		if p < 0 { // -p med y (normal distribution)
			if yt == Z {
				vx, _ := varf(rel(inc(y)), yn)
				vy, _ := varf(ima(y), yn)
				b := 3.2
				p95 := 1.97 * math.Pow(math.Pow(math.Sqrt(vx), b)+math.Pow(math.Sqrt(vy), b), 1/b)
				fac := math.Sqrt(-2.0*math.Log(1.0+p)) / math.Sqrt(-2.0*math.Log(0.05)) // p<0
				r = mk(F, atom)
				m.f[1+r>>1] = fac * p95
				return dex(x, r)
			}
			s2, av := varf(y, yn)
			r := mk(F, atom)
			m.f[1+r>>1] = av + math.Sqrt(s2)*math.Sqrt2*math.Erfinv(2.0*(-p)-1.0)
			return dex(x, r)
		} else if p < 1 { // p med y (percentile)
			if yt == Z {
				y = abs(y)
			}
			nf := p * f(yn-1)
			n := k(nf)
			y = srt(y)
			if n == yn-1 {
				return dex(x, to(atx(y, mki(yn-1)), F))
			}
			d := mk(I, 2)
			m.k[2+d] = n
			m.k[3+d] = n + 1
			y = to(atx(y, d), F)
			r = mk(F, atom)
			m.f[1+r>>1] = (nf-f(n))*m.f[1+y>>1] + (f(n+1)-nf)*m.f[2+y>>1]
			return decr(x, y, r)
		} else {
			panic("value")
		}
	case I:
		if xi := m.k[2+x]; xi == 0 { // 0 med y (cummulative running median)
			if yt == Z || yt > S {
				panic("type")
			}
			if m.k[1+y] != 1 {
				r = mk(F, yn)
				mv(r, y)
				dec(y)
				y = r
			}
			dec(x)
			x = mk(yt, 0)
			xp, yp, gt, cp := ptr(x, yt), ptr(y, yt), gtx[yt], cpx[yt]
			for i := k(0); i < yn; i++ {
				ix := 1 + ibin(xp, i, yp+i, gt)
				x = insert(x, atx(inc(y), mki(i)), ix)
				xp = ptr(x, yt)
				cp(yp+i, xp+(i+1)/2)
			}
			return dex(x, y)

		} else { // n med y (running median, window size n)
			panic("nyi")
			// e.g. www.stat.cmu.edu/~ryantibs/median/binmedian.c
		}
	default:
		panic("type")
	}
}
func lnan(x k) (r k) { // nan value of a list
	t, n := typ(x)
	if t == L {
		if n == 0 {
			return dex(x, mk(C, 0))
		} else {
			x = fst(x)
			t, n = typ(x)
		}
	}
	return lrna(x)
	// TODO: for functions e.g. "(-;1)@3" k7 returns mk(N+2, 0), which conflicts with lambda.
}
func lrna(x k) (r k) {
	t, n := typ(x)
	if t < L {
		r = dex(x, inc(nan[t]))
		if n != atom {
			r = take(n, 0, r)
		}
		return r
	} else if t == L {
		return lrc(x, n, lrna)
	} else if t == A {
		return arc(x, n, lrna)
	} else {
		return dex(x, mk(C, 0))
	}
}

func isnan(x f) bool { return x != x }
func atm1(n k) k {
	if n == atom {
		return 1
	}
	return n
}
func buk(x uint32) (n k) { // from https://golang.org/src/math/bits/bits.go (Len32)
	x--
	if x >= 1<<16 {
		x >>= 16
		n = 16
	}
	if x >= 1<<8 {
		x >>= 8
		n += 8
	}
	n += k(l8t[x])
	if n < 4 {
		return 4
	}
	return n
}

// E:E;e|e e:nve|te| t:n|v v:tA|V n:t[E]|(E)|{E}|N
func prs(x k) (r k) { return par(x, 0) } // `p"…"
func par(x, sto k) (r k) {
	t, n := typ(x)
	if t != C || n == atom {
		if t > N && n == 0 {
			return ltr(x)
		}
		panic("type")
	}
	if n == 0 {
		dec(x)
		return inc(null)
	}
	if sto != 0 {
		dec(asn(mks(".c"), inc(x), inc(null)))
	}
	p := p{p: 8 + x<<2, e: n + 8 + x<<2, lp: 7 + x<<2, sto: sto}
	r = mk(L, 1)
	m.k[2+r] = inc(nans) // ;→`
	for p.p <= p.e {     // ex;ex;…
		r = lcat(r, p.ex(p.noun()))
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
	p   k // current position, m.c[p.p:...]
	m   k // pos after matched token (token: m.c[p.p:p.m])
	e   k // pos after last byte available
	lp  k // m.c index of last newline
	sto k // ~0: store src pointer in nodes (start index in m.c)
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
		case ' ':
			if p.p != p.e && m.c[p.p] == '\\' { // 1+ \x
				m.c[p.p] = 0x08 // lunettes
			}
		case '\t', '\r':
		case '\n':
			if p.p != p.lp+1 {
				p.lp = p.p - 1
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
func (p *p) nNum() bool { // minus part of a number
	if m.c[p.p] != '-' || p.p == p.lp+1 {
		return true
	}
	if c := m.c[p.p-1]; cr0Z(c) || c == ')' || c == ']' {
		p.m = p.p
		return false // verb: exceptions (kref p28)
	}
	return true
}
func (p *p) xx() {
	panic("parse: " + string(m.c[p.lp+1:p.p+1]) + " <-")
}
func (p *p) ex(x k) (r k) { // e:nve|te| t:n|v v:tA|V n:t[E]|(E)|{E}|N
	if match(x, null) {
		return x
	}
	r = p.verb(x)
	if r != 0 && p.t(sObr) {
		x = p.idxr(r)
		r = 0
	}
	if r == 0 { // n
		ps := p.p
		r = p.noun()
		if match(r, null) {
			return dex(r, x) // n
		}
		if m.k[x]>>28 == N && m.k[2+x] == 0 {
			r = p.ex(r)
			m.k[2+x] = r
			m.k[3+x] = cat(mkc(' '), mkb(m.c[ps-1:p.p]))
			return x
		} else if v := p.verb(r); v == 0 {
			return p.store(ps, compose(l2(x, p.ex(r)))) // te
		} else {
			if y := p.ex(p.noun()); match(y, null) {
				return p.store(ps, l2(v, dex(y, x))) // e.g. 2+
			} else if m.k[v]>>28 == N+2 && m.k[2+v] == dyad+18 { // @
				return dex(v, p.store(ps, l2(x, y))) // x@y
			} else if cmpvrb(y) {
				if tv, nv := typ(v); tv == N+2 && nv == atom && m.k[2+v] == dyad {
					// TODO: also for global assign?
					return p.store(ps, l3(iasn(v), x, y)) // n:y (not a composition)
				}
				return p.store(ps, compose(l2(l2(v, x), y))) // 2+ *
			} else {
				return p.store(ps, l3(iasn(v), x, y)) // nve
			}
		}
	} else {
		ps := p.p
		x = p.ex(p.noun())
		if match(x, null) {
			return dex(x, r) // v
		} else {
			return p.store(ps, compose(l2(monad(r), x))) // ve
		}
	}
}
func iasn(v k) k { // mark (local) infix assignment
	if t, n := typ(v); t == N+2 && n == atom && m.k[v+2] == dyad {
		m.k[v] = (N+1)<<28 | atom
		m.k[v+2] = 0
		m.k[v+3] = 1 // marker
	}
	return v
}
func (p *p) store(ps k, x k) k { // store source position in refcount's high bits
	if m.k[x]>>28 != L {
		panic("type")
	}
	if p.sto != 0 {
		m.k[1+x] |= (ps - 1 - p.sto) << 16
	}
	return x
}
func (p *p) verb(x k) (r k) { // v:tA|V
	if p.t(sAdv) {
		for {
			r = mk(L, 2)
			m.k[2+r] = p.a(pAdv)
			m.k[3+r] = x
			x = r
			if p.t(sAdv) == false {
				break
			}
		}
		return r
	}
	if t, n := typ(x); t > N && n == atom {
		return x // V
	}
	return 0
}
func (p *p) noun() (r k) {
	switch {
	case p.t(sHex):
		r = p.a(pHex)
		return p.idxr(r)
	case p.t(sIov):
		return p.idxr(p.a(pIov))
	case p.t(sNum) && p.nNum():
		r = p.a(pNum)
		for p.p != p.e && m.c[p.p] == ' ' {
			if n := sNum(m.c[p.p+1 : p.e]); n == 0 {
				break
			} else {
				p.p++
				p.m = p.p + k(n)
			}
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
		args := uf(p.idxa(mk(L, 0)))
		ln := m.k[args] & atom
		r = mk(N+1, 0) // lambda is indicated with length 0 but uses 2 fields:
		tree := p.lst(mk(L, 0), sCcb)
		if m.k[tree]&atom == 1 {
			tree = fst(tree)
		} else {
			tree = cat(inc(nans), tree)
		}
		dst, n := mk(C, p.p-st), p.p-st
		m.k[2+r] = dst  // #1: string representation
		m.k[3+r] = tree // #2: parse tree
		dst = 8 + dst<<2
		copy(m.c[dst:dst+n], m.c[st:p.p])
		if ln == 0 {
			dec(args)
			ln = argn(m.k[3+r], 0)
			if ln == 0 || ln > 3 {
				panic("valence(lambda)")
			}
			args = mk(S, ln)
			ln = argn(m.k[3+r], 0)
			for i := k(0); i < ln; i++ {
				m.k[2+args+i] = 'x' + i // `x..`z
			}
		}
		if N+ln > 15 { // type overflow (4bits)
			panic("args")
		}
		m.k[3+r] = l2(unq(locl(m.k[3+r], args)), m.k[3+r])
		m.k[r] = (N+ln)<<28 | 0
		return p.idxr(r)
	case p.t(sOpa):
		p.p = p.m
		r = p.lst(mk(L, 0), sCpa)
		if n := m.k[r] & atom; n == 0 {
			return p.idxr(r)
		} else if n == 1 {
			if m.k[m.k[2+r]]>>28 > N { // verb
				return p.idxr(r)
			}
			// TODO drv
			return p.idxr(fst(r))
		}
		return p.idxr(cat(enlist(inc(null)), r))
	case p.t(sBin):
		return p.idxr(p.a(pBin))
	case p.t(sQlq):
		sto := p.sto
		p.sto = 0
		r = p.sql(p.a(pSql))
		p.sto = sto
		return r
	case p.t(sQlv):
		return inc(null)
	case p.t(sNam):
		n := p.a(pNam)
		if m.k[n]>>28 == L {
			return p.idxa(n)
		}
		return p.idxr(n)
	case p.t(sVrb):
		r = p.a(pVrb)
		if m.k[r]>>28 == N+2 && m.k[2+r] == dyad {
			if c := m.c[p.p-2]; m.c[p.p] != '[' && (p.p == p.lp+2 || c == ';' || c == ' ' || c == '[' || c == '(') {
				r = dex(r, mk(N, 2)) // :expr
				m.k[2+r] = 0
				m.k[3+r] = 0
				return r
			}
		}
		return p.idxr(r)
	}
	return inc(null)
}
func (p *p) idxr(x k) (r k) { // […]
	for p.t(sObr) {
		p.p = p.m
		r = mk(L, 1)
		m.k[2+r] = x
		x = p.lst(r, sCbr)
	}
	return x
}
func (p *p) idxa(x k) (r k) { // a.b[..]
	if t := m.k[x] >> 28; t != L {
		panic("assert idxa")
	}
	if p.t(sObr) {
		p.p = p.m
		x = p.lst(x, sCbr)
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
		if t := m.k[r] >> 28; t != L {
			panic("assert lst not L")
		}
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
func (p *p) sql(x k) (r k) { // select|update|delete [ex] [by expr] from t [where ex]
	var t, tc, c, cc, b, bc, a, ac k // t(table), c(where), b(group-by), a(aggregate)
	if !p.t(sQlv) {                  // select
		a, ac = p.mustex()
		a = sqle(a, ac)
	} else if m.k[2+x] == yudt {
		panic("parse")
	}
	if !p.t(sQlv) {
		panic("parse")
	}
	s := p.a(pSql)
	if m.k[2+s] == yby { // by
		b, bc = p.mustex()
		b = sqle(b, bc)
		if !p.t(sQlv) {
			panic("parse")
		}
		s = dex(s, p.a(pSql))
	}
	if m.k[2+s] == yfrm { // from
		t, tc = p.mustex()
		decr(s, tc, 0)
	} else {
		panic("parse")
	}
	if p.t(sQlv) { // where
		s = p.a(pSql)
		if m.k[2+s] != ywer {
			panic("parse")
		}
		c, cc = p.mustex()
		c = dex(s, val(sqle(c, cc)))
		if m.k[c]&atom == 1 {
			c = fst(c)
		} else {
			c = enlist(c)
		}
	} else {
		c = mk(L, 0)
	}
	f := mk(N+2, atom)
	m.k[2+f] = 14 + dyad // #
	if m.k[2+x] == yudt || m.k[2+x] == ydel {
		m.k[2+f]++ // _
	}
	r = mk(L, 3)
	m.k[2+r] = f
	switch m.k[2+x] {
	case ysel, yudt:
		if a == 0 && b == 0 {
			m.k[3+r] = c // select from t where c | select from t
			m.k[4+r] = t
			break
		}
		m.k[3+r] = t
		m.k[4+r] = c
		if b == 0 { // select a from t
			r = lcat(r, a)
			break
		} else {
			r = lcat(r, b)
			if a == 0 {
				r = lcat(r, key(mk(S, 0), mk(L, 0)))
			} else {
				r = lcat(r, a)
			}
		}
	case ydel:
		if m.k[c]&atom == 0 {
			c = dex(c, 0)
		}
		if b != 0 || (c != 0 && a != 0) || (a == 0 && b == 0 && c == 0) {
			panic("parse")
		}
		m.k[4+r] = t
		if a != 0 {
			m.k[3+r] = til(a) // delete a from t
		} else {
			m.k[3+r] = c // delete from t where c
		}
	default:
		panic("parse")
	}
	return dex(x, r)
}
func scat(x k) (r k) { // a → ,`a | a,b,.. → ,`a`b`..
	t, n := typ(x)
	if t == S {
		return x
	} else if t != L || n != 3 || m.k[3+x] != 12+dyad || m.k[4+x]>>28 != S {
		panic("parse")
	}
	return dex(x, cat(inc(m.k[2+x]), scat(inc(m.k[3+x]))))
}
func sqle(x, xc k) (r k) { // parse sql subphrase
	n := m.k[xc] & atom
	d, b, p, c, a := key(mk(S, 0), mk(L, 0)), k(0), 8+xc<<2, c(0), k(0)
	for i := k(0); i < n; i++ { // split at , not within (,) and parse each subphrase as ex
		c = m.c[p+i]
		if c == '"' {
			i += k(sStr(m.c[p+i : p+n-i]))
			continue
		} else if c == '(' {
			b++
		} else if c == ')' {
			if b == 0 {
				panic("parse")
			}
			b--
		}
		if c == ',' && b == 0 {
			d = cat(d, sqlp(mkb(m.c[p+a:p+i])))
			a = i + 1
		} else if i == n-1 {
			if m.c[p+n-1] == ' ' {
				n--
			}
			d = cat(d, sqlp(mkb(m.c[p+a:p+n])))
		}
	}
	return decr(x, xc, d)
}
func sqlp(x k) (r k) {
	n := m.k[x] & atom
	if n == 0 {
		return dex(x, key(mk(S, 0), mk(L, 0)))
	}
	p, s := p{p: 8 + x<<2, e: n + 8 + x<<2, lp: 7 + x<<2, sto: 0}, k(0)
	if p.t(sNam) {
		s = p.a(pNam)
		if !p.t(sVrb) {
			if p.p == p.e {
				return key(inc(s), toex(s, x))
			}
			panic("parse")
		}
		v := p.a(pVrb)
		if t, n := typ(v); t == N+2 && n == atom && m.k[2+v] == dyad {
			x = decr(x, v, mkb(m.c[p.p:p.e]))
			return key(s, toex(prs(inc(x)), x))
		}
		decr(s, v, 0)
	}
	e := toex(prs(inc(x)), x)
	s = mk(S, atom)
	m.k[2+x] = k('x')
	s = lsym(inc(m.k[2+e]), s)
	return key(s, e)
}
func lsym(x, s k) (r k) { // last symbol in tree
	t, n := typ(x)
	if t == S {
		return dex(s, x)
	} else if t != L {
		return dex(x, s)
	}
	for i := k(0); i < n; i++ {
		s = lsym(inc(m.k[2+i+x]), s)
	}
	return dex(x, s)
}
func toex(x, c k) (r k) {
	if m.c[8+c<<2] == ' ' {
		c = drop(1, c)
	}
	r = l2(x, cat(mkc(' '), cat(mkc(':'), c)))
	m.k[r] = N<<28 | 2
	return r
}
func (p *p) mustex() (r, c k) {
	s := p.p
	r = p.ex(p.noun())
	if match(r, null) {
		panic("parse")
	}
	return r, mkb(m.c[s:p.p])
}
func monad(x k) (r k) { // force monad
	t, _ := typ(x)
	if t == N+2 {
		r = mk(N+1, atom)
		m.k[2+r] = m.k[2+x] - dyad
		if m.k[2+x] >= 2*dyad {
			panic("parse monad")
		}
		return dex(x, r)
	}
	return x
}
func compose(x k) (r k) { // composition
	t, n := typ(x)
	if t != L {
		panic("assert")
	} else if n != 2 {
		return x
	}
	if cmpvrb(m.k[2+x]) && cmpvrb(m.k[3+x]) {
		r = mk(N+2, atom)
		m.k[2+r] = 19 + dyad // cal
		return cat(r, x)
	}
	return x
}
func cmpvrb(x k) bool { // is allowed in composition
	t, n := typ(x)
	if n == atom && (t == N+1 || t == N+2) {
		return true
	} else if t != L {
		return false
	}
	u, v := m.k[2+x], m.k[3+x]
	if n == 2 && m.k[u]>>28 == N+2 {
		if code := m.k[2+m.k[2+x]]; code == 0 || code == 80 { // (:;..)
			panic("assignment in composition")
			return false // assignment, not composition
		}
		return true // 1+
	} else if n == 3 && m.k[u]>>28 == N+2 && m.k[2+u] == 19+dyad && cmpvrb(v) {
		return true // (.;v;w)
	} else if n == 2 && t == L {
		code := m.k[2+m.k[2+x]]
		if code > dyad {
			code -= dyad
		}
		if m.k[m.k[2+x]]&atom == atom && code >= 30 && code <= 32 {
			return true
		}
	}
	return false
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
func pNum(b []byte) (r k) {
	r, o := aton(b)
	if !o {
		if len(b) > 0 && b[len(b)-1] == 'e' {
			return dex(r, inc(nan[F]))
		}
		panic("number")
	}
	return r
}
func pStr(b []byte) (r k) { // "a"|"a\nbc": `c|`C
	r = pQot(b)
	if _, n := typ(r); n == 1 {
		m.k[r] = C<<28 | atom
	}
	return r
}
func pDot(b []byte) (r k) { return c2s(mkb(b[1:])) } // .name: `n
func pNam(b []byte) (r k) { return c2s(mkb(b)) }
func pSym(b []byte) (r k) { // `name|`"name": `n
	if len(b) == 1 {
		r = inc(nans)
	} else if len(b) > 1 && b[1] != '"' {
		r = c2s(mkb(b[1:]))
	} else {
		r = c2s(pQot(b[1:]))
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
			if len(b) == 1 {
				r = mk(N+2, atom)
				if i >= 30 {
					m.k[r] = atom | (N+1)<<28
				}
				m.k[2+r] = k(i) + dyad
			} else {
				r = mk(N+1, atom)
				m.k[2+r] = k(i)
			}
			m.k[3+r] = 0 // clear infix assignment
			return r
		}
	}
	panic("pVrb")
}
func pIov(b []byte) (r k) {
	r = mk(N+2, atom)
	m.k[2+r] = dyad + 20 + k(b[0]-'0') // ioverb parses always as `2
	if b[0] == 0x08 {
		m.k[2+r] = dyad + 28
	}
	return r
}
func pAdv(b []byte) (r k) {
	f := k(dyad + 30)
	if len(b) > 1 {
		f -= dyad
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
	x := atx(inc(m.k[3]), c2s(mkb(b)))
	if xt, xn := typ(x); xt != C && xn != atom {
		panic("parse builtin")
	}
	if f := m.c[8+x<<2]; f < c(dyad) {
		r = mk(N+1, atom)
		m.k[2+r] = k(f)
	} else {
		r = mk(N+2, atom)
		m.k[2+r] = k(f)
	}
	dec(x)
	return r
}
func pSql(b []byte) (r k) { // ksql
	s := [6]k{ysel, yudt, ydel, yby, yfrm, ywer}
	q := c2s(mkb(b))
	for i := k(0); i < 6; i++ {
		if s[i] == m.k[2+q] {
			return q
		}
	}
	panic("pSql")
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
	n := 0
	if len(b) > 1 && b[0] == '0' && b[1] == 'N' {
		n = 2
	} else {
		n = sFlt(b)
	}
	if n > 0 && len(b) > n && (b[n] == 'i' || b[n] == 'a' || b[n] == 'p') {
		n += 1 + sFlt(b[n+1:])
	}
	return n
}
func sFlt(b []byte) (r int) { // -0.12e-12|1f
	if len(b) > 1 && b[0] == '0' && (b[1] == 'n' || b[1] == 'w') {
		return 2
	} else if len(b) > 2 && b[2] == 'w' && b[1] == '0' && b[0] == '-' {
		return 3
	}
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
	if len(b) > r && r > 0 && b[r] == 'e' {
		r += 1 + sExp(b[r+1:])
	}
	if r > 0 && len(b) > r && b[r] == 'f' {
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
	o := false
	for i, c := range b {
		if cr09(c) || craZ(c) {
			o = true
			if i == 0 && cr09(c) {
				return 0
			}
		} else if c != '.' {
			if o {
				return i
			}
			return 0
		}
	}
	if o {
		return len(b)
	}
	return 0
}
func sDot(b []byte) (r int) { // .name
	if b[0] != '.' || len(b) < 2 {
		return 0
	}
	r = 1 + sNam(b[1:])
	if r > 1 {
		return r
	}
	return 0
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
		if !(cr0Z(c) || c == '.') {
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
func sIov(b []byte) int { // ioverb 0: .. 9:
	if len(b) > 1 && b[1] == ':' && b[0] >= '0' && b[0] <= '9' {
		return 2
	} else if b[0] == 0x08 { // lunettes
		return 1
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
	n := sNam(b)
	if n == 0 {
		return 0
	}
	names := m.k[2+m.k[3]]
	a, max := fnd(inc(names), c2s(mkb(b[:n]))), cnt(inc(names))
	if match(a, max) {
		n = 0
	}
	dec(a)
	dec(max)
	return n
}

func sQlq(b []byte) int { // ksql select|update|delete
	if r := sQls(b); r == 7 { // "select " (including space)
		return r - 1
	}
	return 0
}
func sQlv(b []byte) int { // ksql from|by|where
	if r := sQls(b); r > 0 && r < 7 {
		return r - 1
	}
	return 0
}
func sQls(b []byte) int {
	s := [6]k{ysel, yudt, ydel, yby, yfrm, ywer}
	for i := k(0); i < 6; i++ {
		c := m.k[2+s[i]+m.k[stab]-256]
		p, n := 8+c<<2, m.k[c]&atom
		if k(len(b)) > n {
			for j := k(0); j < n; j++ {
				if b[j] != m.c[p+j] {
					break
				}
				if j == n-1 && b[n] == ' ' {
					return int(n + 1)
				}
			}
		}
	}
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
func cr0Z(c c) bool     { return cr09(c) || craZ(c) }
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
func argn(x, ln k) k { // count args of lambda parse tree
	t, n := typ(x)
	switch t {
	case S:
		if u := m.k[2+x]; n == atom && u >= 'x' && u <= 'z' && ln < 1+u-'x' {
			ln = 1 + u - 'x'
		}
	case L:
		for i := k(0); i < n; i++ {
			ln = argn(m.k[2+i+x], ln)
		}
	}
	return ln
}
func locl(x, l k) k { // local list of lambda parse tree
	t, n := typ(x)
	if t == L {
		for i := k(0); i < n; i++ {
			l = locl(m.k[2+i+x], l)
		}
		if n == 3 {
			v := m.k[3+x]
			if f := m.k[2+x]; m.k[f]>>28 == N+1 && m.k[2+f] == 0 && m.k[3+f] == 1 { // infix local assignment
				inc(v)
				if m.k[v]>>28 == L { // (a;b):..
					v = drop(1, v)
				}
				l = unq(cat(l, v))
			}
		}
	}
	return l
}
func l2(x, y k) (r k) {
	r = mk(L, 2)
	m.k[2+r], m.k[3+r] = x, y
	return r
}
func l3(x, y, z k) (r k) {
	r = mk(L, 3)
	m.k[2+r], m.k[3+r], m.k[4+r] = x, y, z
	return r
}
func aton(b []byte) (r k, o bool) { // 0|1f|2p|-2.3e+4|1i2|1a90: `i|`f|`z
	if len(b) == 0 {
		return r, false
	}
	if len(b) > 1 && b[len(b)-1] == 'p' { // 2p→2*π
		r, o = aton(b[:len(b)-1])
		if !o {
			return r, false
		}
		if m.k[r]>>28 == I {
			m.f[1+r>>1] = f(i(m.k[2+r])) * math.Pi
			m.k[r] = F<<28 | atom
		} else {
			m.f[1+r>>1] *= math.Pi
		}
		return r, true
	}
	for i, c := range b {
		if c == 'i' || c == 'a' {
			r, o = aton(b[:i])
			if !o {
				return r, false
			}
			r = to(r, Z)
			if i == len(b)-1 {
				return r, true
			}
			y, o := aton(b[i+1:])
			if !o {
				return y, false
			}
			y = to(y, F)
			if c == 'i' {
				m.f[3+r>>1] = m.f[1+y>>1]
			} else {
				var s, c f
				switch a := m.f[1+y>>1]; a { // avoid rounding errors
				case 0:
					s, c = 0, 1
				case 90:
					s, c = 1, 0
				case 180:
					s, c = 0, -1
				case 270:
					s, c = -1, 0
				default:
					s, c = math.Sincos(math.Pi * a / 180.0)
				}
				m.f[2+r>>1], m.f[3+r>>1] = m.f[2+r>>1]*c, m.f[2+r>>1]*s
			}
			dec(y)
			return r, true
		}
	}
	if len(b) == 2 && b[0] == '0' { // 0N 0n 0w
		if b[1] == 'N' {
			return inc(nan[I]), true
		} else if b[1] == 'n' {
			return inc(nan[F]), true
		} else if b[1] == 'w' {
			r = mk(F, atom)
			m.f[1+r>>1] = finf
			return r, true
		}
	} else if len(b) == 3 && b[0] == '-' && b[1] == '0' && b[2] == 'w' {
		r = mk(F, atom)
		m.f[1+r>>1] = -finf
		return r, true
	}
	f := 0.0
	if len(b) > 1 {
		if c := b[len(b)-1]; c == 'f' || c == '.' {
			b = b[:len(b)-1]
			f = 1.0
		} else if b[0] == '.' {
			b = b[1:]
			if len(b) > 21 {
				return inc(null), false
			}
			f = 1.0 / e10[len(b)]
		}
	}
	if x, o := atoi(b); o {
		if f == 0 {
			r = mki(k(i(x)))
		} else {
			r = mk(F, atom)
			if x == 0 && b[0] == '-' { // -0f
				f = -f
			}
			m.f[1+r>>1] = f * float64(x)
		}
		return r, true
	}
	if x, o := atof(b); o {
		r = mk(F, atom)
		m.f[1+r>>1] = x
		return r, true
	}
	return inc(null), false
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
func atof(b []c) (f, bool) {
	man, exp, neg, o := flt(b)
	if !o {
		return 0, false
	}
	v := f(man)
	if exp < 0 {
		for exp < 0 {
			if exp < -21 {
				exp += 21
				v /= e10[21]
			} else {
				v /= e10[-exp]
				exp = 0
			}
		}
	} else if exp > 0 {
		for exp > 0 {
			if exp > 21 {
				exp -= 21
				v *= e10[21]
			} else {
				v *= e10[exp]
				exp = 0
			}
		}
	}
	if neg {
		return -v, true
	}
	return v, true
}
func ftoa(dst k, v f) k {
	switch {
	case v == 0:
		if 0 != (1<<63)&*(*uint64)(unsafe.Pointer(&v)) {
			m.c[dst] = '-'
			m.c[dst+1] = '0'
			return 2
		}
		m.c[dst] = '0'
		return 1
	case v != v:
		m.c[dst] = '0'
		m.c[dst+1] = 'n'
		return 2
	case v+v == v && v > 0:
		m.c[dst] = '0'
		m.c[dst+1] = 'w'
		return 2
	case v+v == v && v < 0:
		m.c[dst] = '-'
		m.c[dst+1] = '0'
		m.c[dst+2] = 'w'
		return 3
	}
	var b []c
	e, sn, n, d, t := 0, k(0), k(7), k(1), k(1)
	if v < 0 {
		b = m.c[dst:]
		b[0], v, sn = '-', -v, 1
	} else {
		b = m.c[dst-1:]
	}
	for v >= 10 {
		e++
		v /= 10
	}
	for v < 1 {
		e--
		v *= 10
	}
	h := 5.0
	for i := k(0); i < n; i++ {
		h /= 10
	}
	v += h
	if v >= 10 {
		e++
		v /= 10
	}
	for i := k(0); i < n; i++ {
		s := int(v)
		b[i+2] = byte(s + '0')
		v -= float64(s)
		v *= 10
	}
	for b[1+n] == '0' && n > 1 {
		n--
	}
	b[1] = b[2]
	if n == 1 {
		d = 0
	}
	if e == 0 { // 1
		b[2] = '.'
		return n + d + sn
	} else if e < 0 && e > -5 { // 0.01234
		for i := int(n); i >= 0; i-- {
			b[2+i-e] = b[2+i]
		}
		for i := 1; i < -e; i++ {
			b[2+i] = '0'
		}
		b[1] = '0'
		b[2] = '.'
		return n + sn + 1 + k(-e)
	} else if e > 0 && e < 7 { // 123.456
		if n <= k(e) {
			n = k(e) + 1
		}
		b[2] = '.'
		for i := 0; i < e; i++ {
			b[2+i], b[3+i] = b[3+i], b[2+i]
		}
		if n == k(e)+1 {
			d = 0
		}
		return n + sn + d
	} else { // 1.234
		b[2] = '.'
	}
	t = 1 + d
	b[n+t] = 'e'
	if e < 0 {
		t++
		e = -e
		b[n+t] = '-'
	}
	uu := false
	if u := c(e / 100); u > 0 {
		t++
		uu = true
		b[n+t] = u + '0'
	}
	if u := c(e/10) % 10; uu || u > 0 {
		t++
		uu = true
		b[n+t] = u + '0'
	}
	t++
	b[n+t] = c(e%10) + '0'
	return n + t + sn
}
func low(c c) c { return c | ('x' - 'X') }
func flt(s []c) (man uint64, exp int, neg, ok bool) {
	i := 0
	if i >= len(s) {
		return
	}
	switch {
	case s[i] == '+':
		i++
	case s[i] == '-':
		neg = true
		i++
	}
	maxMantDigits := 19 // 10^19 fits in uint64
	expChar := byte('e')
	sawdot := false
	sawdigits := false
	nd := 0
	ndMant := 0
	dp := 0
	for ; i < len(s); i++ {
		switch c := s[i]; true {
		case c == '.':
			if sawdot {
				return
			}
			sawdot = true
			dp = nd
			continue
		case '0' <= c && c <= '9':
			sawdigits = true
			if c == '0' && nd == 0 { // ignore leading zeros
				dp--
				continue
			}
			nd++
			if ndMant < maxMantDigits {
				man *= 10
				man += uint64(c - '0')
				ndMant++
			}
			continue
		}
		break
	}
	if !sawdigits {
		return
	}
	if !sawdot {
		dp = nd
	}
	if i < len(s) && low(s[i]) == expChar {
		i++
		if i >= len(s) {
			return
		}
		esign := 1
		if s[i] == '+' {
			i++
		} else if s[i] == '-' {
			i++
			esign = -1
		}
		if i >= len(s) || s[i] < '0' || s[i] > '9' {
			return
		}
		e := 0
		for ; i < len(s) && ('0' <= s[i] && s[i] <= '9'); i++ {
			if e < 10000 {
				e = e*10 + int(s[i]) - '0'
			}
		}
		dp += e * esign
	}
	if i != len(s) {
		return
	}
	if man != 0 {
		exp = dp - ndMant
	}
	ok = true
	return
}
func adler() (r k) {
	a, b := k(1), k(0)
	for i := k(0); i < k(len(m.c)); i++ {
		a = (a + k(m.c[i])) % 65521
		b = (b + a) % 65521
	}
	return mki(a | b<<16)
}
func stats() (r k) { // \b (memory stats used/free buckets)
	u, f := mk(I, 32), mk(I, 32)
	for i := k(0); i < 32; i++ {
		m.k[2+i+u], m.k[2+i+f] = 0, 0
	}
	x, o := k(0), k(0)
	for x < 1<<(m.k[2]-2) {
		if xt := m.k[x] >> 28; xt == 0 {
			p := m.k[x]
			if p < 4 || p > 31 {
				panic("free block type")
			}
			m.k[2+f+p]++
			o = 1 << (p - 2)
		} else {
			t, n := typ(x)
			p := bk(t, n)
			if p < 4 || p > 31 {
				panic("used block type")
			}
			m.k[2+u+p]++
			n = atm1(n)
			o = 1 << (p - 2)
		}
		x += o
	}
	return kxy(mks(".stats"), u, f)
}

var e10 = [22]f{1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9, 1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16, 1e17, 1e18, 1e19, 1e20, 1e21}
var l8t = [256]c{
	0x00, 0x01, 0x02, 0x02, 0x03, 0x03, 0x03, 0x03, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06,
	0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07,
	0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
	0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
}
