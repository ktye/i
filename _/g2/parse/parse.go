package main

import (
	"bufio"
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strconv"
	"time"
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

// 1 → 1
// 0.0 → 0
// 0. → 0
// .0 → 0
func pFloat(b []byte) (A, []byte) {
	if n := sFloat(b); n > 0 {
		f, _ := strconv.ParseFloat(string(b[:n]), 64)
		return f, b[n:]
	}
	return nil, b
}

// 1a30 → (0.8660254037844387+0.49999999999999994i)
func pComplex(b []byte) (A, []byte) {
	f, b := pFloat(b)
	if f == nil {
		return nil, b
	}
	r := f.(float64)
	a := 0.0
	if len(b) > 0 && b[0] == 'a' {
		b = b[1:]
		if len(b) > 0 {
			f, b = pFloat(b)
			if f != nil {
				a = f.(float64)
			}
		}
	}
	return cmplx.Rect(r, math.Pi*a/180.0), b
}
func pTime(b []byte) (A, []byte) {
	if len(b) > 1 && b[0] == '0' && b[1] == 'T' {
		return time.Time{}, b[2:]
	}
	if len(b) < 19 {
		return nil, b
	}
	t := 10
	if b[19] == '.' {
		for i := 20; i < len(b); i++ {
			if !is(b[i], NM) {
				break
			}
			t = i
		}
	}
	d, e := time.Parse("2006.01.02T15:14:05", string(b[:t]))
	if e != nil {
		panic(e)
	}
	return d, b[:t]
}
func pDuration(b []byte) (A, []byte) {
	var d time.Duration
	var r A
	for len(b) > 0 {
		if dt, o, t := pDurationPart(b); o {
			d += dt
			r = d
			b = t
		} else {
			break
		}
	}
	return r, b
}
func pDurationPart(b []byte) (time.Duration, bool, []byte) {
	if n := sFloat(b); n > 0 && len(b) > n { // ns us ms s m h
		c := b[n]
		if c == 'n' || c == 'u' || c == 'm' {
			if len(b) > n+1 && b[1+n] == 's' {
				n++
			}
		}
		if d, e := time.ParseDuration(string(b[:1+n])); e == nil {
			return d, true, b[1+n:]
		}
	}
	return 0, false, b
}
func pNum(b []byte) (A, []byte) { // atom 12 -3 1.2e-12
	n := sFloat(b)
	if n == 0 {
		return nil, b
	}
	if len(b) > n {
		switch b[n] {
		case 'a':
			return pComplex(b)
		case '.', 'T':
			return pTime(b)
		case 'n', 'u', 'm', 's', 'h':
			return pDuration(b)
		}
	}
	return pFloat(b)
}
func pVrb(b []byte) (A, []byte) {
	if is(b[0], VB) {
		return Verb(b[0]), b[1:]
	}
	if is(b[0], AD) {
		if len(b) > 1 && b[1] == ':' {
			return Verb(b[:2]), b[2:]
		}
		return Verb(b[0]), b[1:]
	}
	return nil, b
}
func token(b []byte) (r []A) {
	for len(b) > 0 {
		for i, p := range Parsers {
			b = ws(b)
			if len(b) == 0 {
				break
			}
			if x, t := p(b); x != nil {
				r = append(r, x)
				b = t
				break
			}
			if i == len(Parsers)-1 {
				panic(fmt.Errorf("parse:%s", string(b)))
			}
		}
	}
	return r
}
func ws(b []byte) []byte {
	for i, u := range b {
		if u > 32 {
			return b[i:]
		}
	}
	return nil
}

func main() {
	scn := bufio.NewScanner(os.Stdin)
	for scn.Scan() {
		fmt.Println(token([]byte(scn.Text())))
	}
}

var Parsers []func([]byte) (A, []byte)
var c_ [256]byte

func is(x, m byte) bool { return (m & c_[x]) != 0 }
func init() {
	Parsers = []func([]byte) (A, []byte){pNum, pVrb}
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
