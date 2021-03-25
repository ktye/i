package main

import (
	"errors"
	"fmt"
	"math"
	"math/cmplx"
	"strconv"
)

const (
	az = 1 << iota //  1 a-z
	AZ             //  2 A-Z
	NM             //  4 numbers     0123456789
	VB             //  8 verbs       :+-*%!&|<>=~,^#_$?@.
	AD             // 16 adverbs     '/\
	TE             // 32 terminators ;)]} (space)
	NW             // 64 nonwhite    33..126
)

// "a" → 97
// "\n" → 10
// "abc" → [97 98 99]
// "a\nb" → [97 10 98]
// 0x → []
// 0xa0 → [16]
// 0x01fF → [1 255]
func pChr(b []byte, l byte) (A, []byte, error) {
	if b[0] == '"' {
		r, b, err := pQuote(b[1:])
		if len(r) == 1 {
			return r[0], b, err
		}
		return r, b, err
	}
	if len(b) > 1 && b[0] == '0' && b[1] == 'x' {
		return pHex(b[2:])
	}
	return nil, b, nil
}
func pQuote(b []byte) ([]byte, []byte, error) {
	r := make([]byte, 0)
	q := false
	for i, c := range b {
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
			} else if c == '"' {
				return r, b[1+i:], nil
			}
			r = append(r, c)
		}
	}
	return nil, b, errors.New(`unmatched "`)
}
func pHex(b []byte) (r []byte, t []byte, e error) {
	hc := func(c byte) byte {
		if is(NM, c) {
			return c - '0'
		}
		if is(az, c) {
			return c - 'a'
		}
		if is(AZ, c) {
			return c - 'A'
		}
		return 16
	}
	r = make([]byte, 0)
	for len(b) > 1 {
		c := hc(b[0])
		if c == 16 {
			break
		}
		x := 16 * c
		c = hc(b[1])
		if c == 16 {
			break
		}
		r = append(r, x+c)
		b = b[2:]
	}
	return r, b, nil
}

// 1 → 1
// 0.0 → 0
// 0. → 0
// .0 → 0
func pFloat(b []byte) (A, []byte, error) {
	if n := sFloat(b); n > 0 {
		f, e := strconv.ParseFloat(string(b[:n]), 64)
		return f, b[n:], e
	}
	return nil, b, nil
}

// 1a30 → (0.8660254037844387+0.49999999999999994i)
func pComplex(b []byte) (A, []byte, error) {
	f, b, e := pFloat(b)
	if f == nil || e != nil {
		return f, b, e
	}
	r := f.(float64)
	a := 0.0
	if len(b) > 0 && b[0] == 'a' {
		b = b[1:]
		if len(b) > 0 {
			f, b, e = pFloat(b)
			if f != nil && e == nil {
				a = f.(float64)
			}
		}
	}
	return cmplx.Rect(r, math.Pi*a/180.0), b, e
}
func pNum(b []byte, l byte) (A, []byte, error) { // atom 12 -3 1.2e-12
	n := sFloat(b)
	if n == 0 {
		return nil, b, nil
	}
	if b[0] == '-' && is(l, az|AZ|NM) || l == ')' || l == ']' || l == '"' {
		return nil, b, nil
	}
	if len(b) > n {
		switch b[n] {
		case 'a':
			return pComplex(b)
		}
	}
	return pFloat(b)
}
func pVrb(b []byte, l byte) (A, []byte, error) {
	if is(b[0], VB) {
		return Verb(b[0]), b[1:], nil
	}
	if is(b[0], AD) {
		if len(b) > 1 && b[1] == ':' {
			return Verb(b[:2]), b[2:], nil
		}
		return Verb(b[0]), b[1:], nil
	}
	return nil, b, nil
}
func token(b []byte) (r List, sp []int, e error) {
	var parsers = []func([]byte, byte) (A, []byte, error){pChr, pNum, pVrb}
	o := b
	q := 0
	n := len(b)
	l := byte(0)
	pos := func(b []byte) {
		q += n - len(b)
		n = len(b)
		if q > 0 {
			l = o[q]
		}
	}
	for len(b) > 0 {
		for i, p := range parsers {
			b = ws(b)
			pos(b)
			if len(b) == 0 {
				break
			}
			if x, t, e := p(b, l); e != nil {
				return nil, nil, SrcError{Src: o, Pos: q, Err: e}
			} else if x != nil {
				r = append(r, x)
				b = t
				pos(b)
				sp = append(sp, q)
				break
			}
			if i == len(parsers)-1 {
				panic(fmt.Errorf("token:%s", string(b)))
			}
		}
	}
	return r, sp, nil
}
func ws(b []byte) []byte {
	for i, u := range b {
		if u > 32 {
			return b[i:]
		}
	}
	return nil
}

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
