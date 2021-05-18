package k

import (
	"fmt"
	"math"
	"strconv"
)

type ftok func(b []byte, l byte) (T, int)
type token struct {
	t T
	p int
	e int
}
type semicolon byte
type terminator byte
type open byte
type varname string

var tokenizers = []ftok{tHex, tBool, tChr, tNums, tSyms, tVerb, tAdv, tVar, tOpen, tTerm}

func (k *K) tok(b []byte) (r []token) {
	p := 0
	l := byte(0)
	for {
		for len(b) > 0 {
			if n := ws(b, l); n > 0 {
				p += n
				l = b[n-1]
				b = b[n:]
			} else {
				break
			}
		}
		if len(b) == 0 {
			return r
		}
		for i, t := range tokenizers {
			v, n := t(b, l)
			if n > 0 {
				if s, o := v.(varname); o {
					if _, o := k.Func[string(s)]; o {
						v = verb(s)
					}
				}
				r = append(r, token{v, p, p + n})
				l = b[n-1]
				b = b[n:]
				p += n
				break
			}
			if i == len(tokenizers)-1 {
				panic("unknown token")
			}
		}
	}
}

func tHex(b []byte, l byte) (T, int) { // 0xab12ef
	if len(b) < 2 || b[0] != '0' || b[1] != 'x' {
		return nil, 0
	}
	hx := func(c byte) (r byte) {
		r = 16
		if is(c, NM) {
			r = c - '0'
		} else if is(c, AZ) {
			r = 10 + c - 'A'
		} else if is(c, az) {
			r = 10 + c - 'a'
		}
		if r > 15 {
			panic("token hex")
		}
		return r
	}
	b = b[2:]
	n := 0
	for _, c := range b {
		if is(c, NM+az+AZ) == false {
			break
		}
		n++
	}
	if n%2 != 0 {
		panic("token: hex is odd")
	}
	v := make([]byte, n/2)
	k := 0
	for i := 0; i < n; i += 2 {
		v[k] = hx(b[1+i]) + 16*hx(b[i])
		k++
	}
	cv := C{v: v}
	cv.init()
	return cv, 2 + n
}

func tBool(b []byte, l byte) (r T, n int) { // 0b 1b 0110b
	if b[0] != '0' && b[0] != '1' {
		return nil, 0
	}
	for _, c := range b {
		if c == '0' || c == '1' {
			n++
		}
	}
	if len(b) > n && b[n] == 'b' {
		return bools(b[:n]), 1 + n
	}
	return nil, 0
}
func bools(b []byte) T {
	if len(b) == 1 {
		return b[0] == '1'
	}
	v := mk(false, len(b)).(B)
	for i := range v.v {
		v.v[i] = b[i] == '1'
	}
	return v
}

func tChr(b []byte, l byte) (T, int) {
	if b[0] != '"' {
		return nil, 0
	}
	q := false
	for i, c := range b[1:] {
		if q == false && c == '\\' {
			q = true
			continue
		}
		if q == false && c == '"' {
			s, e := strconv.Unquote(string(b[:2+i]))
			if e != nil {
				panic(e)
			}
			r := []byte(s)
			if len(r) == 1 {
				return byte(s[0]), 2 + i
			}
			return KC(r), 2 + i
		}
		q = false
	}
	panic("parse\"")
}

func tNums(b []byte, l byte) (r T, n int) {
	if is(l, az|AZ|NM) || l == ')' || l == ']' || l == '}' || l == '"' {
		if b[0] == '-' || b[0] == '.' {
			return nil, 0
		}
	}
	if b[0] == '+' {
		return nil, 0
	}
	r, n = tNum(b)
	if n == 0 {
		return r, n
	}
	rt := numtype(r)
	b = b[n:]
	for len(b) > 1 && b[0] == ' ' {
		t, nn := tNum(b[1:])
		if nn == 0 {
			return r, n
		}
		tt := numtype(t)
		if tt > rt {
			r, rt = uptype(r, tt), tt
		}
		if rt > tt {
			t, tt = uptype(t, rt), rt
		}
		r = cat(r, t)
		n += 1 + nn
		b = b[1+nn:]
	}
	return r, n
}
func tNum(b []byte) (T, int) {
	if len(b) > 2 {
		if s := string(b[:3]); s == "-0w" {
			return math.Inf(-1), 3
		} else if s == "0na" {
			return complex(math.NaN(), math.NaN()), 3
		}
	}
	if len(b) > 1 {
		if s := string(b[:2]); s == "0n" {
			return math.NaN(), 2
		} else if s == "0w" {
			return math.Inf(1), 2
		}
	}
	n := sFloat(b)
	if n == 0 {
		return nil, 0
	}
	s := string(b[:n])
	i, e := strconv.Atoi(s)
	if e == nil {
		if z, n := tComplex(float64(i), b, n); n > 0 {
			return z, n
		}
		return i, n
	}
	f, e := strconv.ParseFloat(s, 64)
	if e == nil {
		if z, n := tComplex(f, b, n); n > 0 {
			return z, n
		}
		return f, n
	}
	panic("parse number")
}
func tComplex(f float64, b []byte, n int) (complex128, int) {
	if len(b) > n && b[n] == 'a' {
		if len(b) == 1+n || is(b[1+n], NM) == false {
			return complex(f, 0), 1 + n
		}
		b = b[1+n:]
		if m := sFloat(b); m == 0 {
			return complex(f, 0), 1 + n
		} else {
			if g, e := strconv.ParseFloat(string(b[:m]), 64); e == nil {
				return rrot(f, g), 1 + n + m
			} else {
				panic("parse complex")
			}
		}
	}
	return 0, 0
}

func tVerb(b []byte, l byte) (T, int) {
	if is(b[0], VB) {
		if len(b) > 1 && b[1] == ':' {
			return verb(b[:2]), 2
		}
		return verb(b[0]), 1
	}
	return nil, 0
}

func tAdv(b []byte, l byte) (r T, n int) {
	if is(b[0], AD) == false {
		return nil, 0
	}

	if len(b) > 1 && b[1] == ':' {
		r, n = adverb(b[:2]), 2
	} else {
		r, n = adverb(b[:1]), 1
	}
	if l == ' ' || l == 10 || l == 0 {
		fmt.Println("adverb-verb")
		r = verb(r.(adverb)) // \out \:? '? ':?
	}
	return r, n
}

func tOpen(b []byte, l byte) (T, int) {
	c := b[0]
	switch c {
	case '(', '[', '{':
		return open(c), 1
	}
	return nil, 0
}

func tTerm(b []byte, l byte) (T, int) {
	c := b[0]
	switch c {
	case 10, ';':
		return semicolon(';'), 1
	case ')', ']', '}':
		return terminator(c), 1
	}
	return nil, 0
}

func tVar(b []byte, l byte) (T, int) { // identifier
	if is(b[0], az+AZ) == false {
		return nil, 0
	}
	for i := range b {
		if is(b[i], az+AZ+NM) == false {
			return varname(b[:i]), i
		}
	}
	return varname(b), len(b)
}

func tSyms(b []byte, l byte) (T, int) { // symbols
	s, n := tSym(b)
	b = b[n:]
	if n == 0 {
		return nil, 0
	}
	r := []string{s}
	for len(b) > 0 {
		s, m := tSym(b)
		if m == 0 {
			break
		}
		r = append(r, s)
		n += m
		b = b[m:]
	}
	if len(r) == 1 {
		return r[0], n
	}
	return KS(r), n
}
func tSym(b []byte) (string, int) { // symbol
	if b[0] != '`' {
		return "", 0
	}
	b = b[1:]
	if len(b) == 0 {
		return "", 1
	}
	if s, n := tVar(b, 0); n > 0 {
		return string(s.(varname)), 1 + n
	} else if b[0] == '"' {
		s, n := tChr(b, 0)
		if n > 0 {
			return string(s.([]byte)), 1 + n
		}
	}
	return "", 1
}

func ws(b []byte, l byte) (n int) {
	for _, c := range b {
		if c == 10 || (c > 32 && c < 127) {
			break
		}
		l = c
		n++
	}
	b = b[n:]
	if len(b) > 0 && b[0] == '/' && (l == 0 || l == 32 || l == 10) {
		return n + tCom(b)
	}
	return n
}
func tCom(b []byte) (n int) {
	for i := range b {
		if i == 10 {
			return 1 + i
		}
	}
	return len(b)
}

func sFloat(s []byte) (i int) {
	lower := func(c byte) byte { return c | ('x' - 'X') }
	var mantissa uint64
	if i >= len(s) {
		return
	}
	switch {
	case s[i] == '+':
		i++
	case s[i] == '-':
		i++
	}
	base := uint64(10)
	maxMantDigits := 19 // 10^19 fits in uint64
	expChar := byte('e')
	if i+2 < len(s) && s[i] == '0' && lower(s[i+1]) == 'x' {
		base = 16
		maxMantDigits = 16 // 16^16 fits in uint64
		i += 2
		expChar = 'p'
	}
	sawdot := false
	sawdigits := false
	nd := 0
	ndMant := 0
	dp := 0
loop:
	for ; i < len(s); i++ {
		switch c := s[i]; true {
		case c == '.':
			if sawdot {
				break loop
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
				mantissa *= base
				mantissa += uint64(c - '0')
				ndMant++
			}
			continue
		}
		break
	}
	if !sawdigits {
		return 0
	}
	if !sawdot {
		dp = nd
	}
	if base == 16 {
		dp *= 4
		ndMant *= 4
	}
	if i < len(s) && lower(s[i]) == expChar {
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
	return
}

const (
	az = 1 << iota //  1 a-z
	AZ             //  2 A-Z
	NM             //  4 numbers     0123456789
	VB             //  8 verbs       :+-*%!&|<>=~,^#_$?@.
	AD             // 16 adverbs     '/\
	TE             // 32 terminators ;)]} (space)
	NW             // 64 nonwhite    33..126
)

var c_ [256]byte

func is(x, m byte) bool { return (m & c_[x]) != 0 }
func init() {
	m := func(s string, b byte) {
		for i := range s {
			c_[s[i]] |= b
		}
	}
	m("abcdefghijklmnopqrstuvwxyz", az)
	m("ABCDEFGHIJKLMNOPQRSTUVWXYZ", AZ)
	m("0123456789", NM)
	m(":+-*%!&|<>=~,^#_$?@.", VB)
	m("'/\\", AD)
	m(";)]} ", TE)
	for i := 33; i < 127; i++ {
		c_[i] |= NW
	}
}
