package main

import (
	"bufio"
	"fmt"
	"math/big"
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

func pUint(b []byte) (int, []byte) {
	a := 0
	i := 0
	for i < len(b) {
		if is(b[i], NM) == false {
			if i == 0 {
				return -1, b
			} else {
				break
			}
		}
		a *= 10
		a += int(b[i] - '0')
		i++
	}
	return a, b[i:]
}
func pInt(b []byte) (A, []byte) { // atom -3 234
	neg := 1
	if b[0] == '-' && len(b) > 1 {
		neg = -1
		b = b[1:]
	}
	if u, r := pUint(b); u >= 0 {
		return Int(neg * u), r
	}
	return nil, b
}
func pFloat(b []byte) (A, []byte) {
	if n, ok := sFloat(b); ok {
		f, _ := strconv.ParseFloat(string(b[:n]), 64)
		return Float(f), b[n:]
	}
	return nil, b
}
func pBigInt(b []byte) (A, []byte) {
	for i := range b {
		if b[i] == 'i' {
			var u big.Int
			u.SetString(string(b[:i]), 10)
			return &Big{Int: &u}, b[1+i:]
		}
	}
	return nil, b
}
func pBigFloat(b []byte) (A, []byte) { return nil, b }
func pComplex(b []byte) (A, []byte)  { return nil, b }
func pDuration(b []byte) (A, []byte) {
	var d Duration
	var r A
	for len(b) > 0 {
		if a, t := pDurationPart(b); a != nil {
			d += a.(Duration)
			r = d
			b = t
		} else {
			break
		}
	}
	return r, b
}
func pDurationPart(b []byte) (A, []byte) {
	if n, ok := sFloat(b); ok && len(b) > n { // ns us ms s m h
		c := b[n]
		if c == 'n' || c == 'u' || c == 'm' {
			if len(b) > n+1 && b[1+n] == 's' {
				n++
			}
		}
		if d, e := time.ParseDuration(string(b[:1+n])); e == nil {
			return Duration(d), b[1+n:]
		}
	}
	return nil, b
}
func pNum(b []byte) (A, []byte) { // atom 12 -3 1.2e-12
	x, r := pInt(b)
	if x == nil {
		return nil, b
	}
	if len(r) == 0 {
		return x, r
	}
	switch r[0] {
	case '.', 'e', 'E':
		return pFloat(b)
	case 'f':
		return pBigFloat(b)
	case 'i':
		return pBigInt(b)
	case 'a':
		return pComplex(b)
	case 'n', 'u', 'm', 's', 'h':
		return pDuration(b)
	}
	return x, r
}
func pVrb(b []byte) (A, []byte) {
	if is(b[0], VB) {
		return V(b[0]), b[1:]
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

func sFloat(s []byte) (i int, ok bool) {
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
		case base == 16 && 'a' <= lower(c) && lower(c) <= 'f':
			sawdigits = true
			nd++
			if ndMant < maxMantDigits {
				mantissa *= 16
				mantissa += uint64(lower(c) - 'a' + 10)
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
	ok = true
	return
}
