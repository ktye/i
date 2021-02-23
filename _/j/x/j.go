package x

import (
	"fmt"
	"math/bits"
	"strconv"
	"strings"
)

type J interface {
	J(x uint32) uint32
	M() []uint32
}

func X(m []uint32, x uint32) string {
	if x == 0 {
		panic("XX0")
	} else if x&7 == 0 {
		n := ln(m, x)
		v := make([]string, n)
		for i := uint32(0); i < n; i++ {
			v[i] = X(m, m[2+i+x>>2])
		}
		return "[" + strings.Join(v, " ") + "]"
	} else if x&1 != 0 {
		return strconv.Itoa(int(int32(x) >> 1))
	} else if x&2 != 0 {
		return sy(x)
	} else if x&4 != 0 {
		return string([]byte{33 + byte(x>>3)})
	}
	panic("XX")
}

func ln(m []uint32, x uint32) uint32 { return m[1+(x>>2)] }
func sy(x uint32) string {
	var b []byte
	x >>= 2
	for x > 0 {
		b = append(b, '`'+byte(x%32))
		x >>= 5
	}
	return string(reverse(b))
}
func reverse(b []byte) []byte {
	n := len(b)
	r := make([]byte, n)
	if n == 0 {
		return r
	}
	n--
	for i := 0; i < len(b); i++ {
		r[i] = b[n-i]
	}
	return r
}

func Leak(j J) {
	m := j.M()
	M := make([]uint32, len(m))
	copy(M, m)

	blank := func(x, n uint32) {
		sI(M, x+4, n)
		p := x + 8
		for i := uint32(0); i < n; i++ {
			sI(M, p, 0)
			p += 4
		}
		dx(M, x)
	}
	if n := nn(M, M[1]); n != 0 {
		panic("stack is not clear: #" + strconv.Itoa(int(n)))
	}
	blank(M[1], M[0])
	parse := M[2]
	root := I(M, parse+8)
	dx(M, root) // I(8) contains only 1 refcounted list
	blank(parse, 5)
	dx(M, M[3])

	mark(M)
	p := uint32(32)
	for p < uint32(len(M)) {
		if M[p] != 0 {
			panic(fmt.Errorf("non-free block: %d(%x) rc=%d #=%d", 4*p, 4*p, M[p], M[1+p]))
		}
		n := uint32(1 << bk(M[1+p]))
		p += n >> 2
	}
}
func nn(m []uint32, x uint32) uint32 { return I(m, 4+x) }
func sI(m []uint32, x, y uint32)     { m[x>>2] = y }
func I(m []uint32, x uint32) uint32  { return m[x>>2] }
func dx(m []uint32, x uint32) {
	if x != 0 && x&7 == 0 {
		if I(m, x) == 0 {
			panic("dx on free")
		}
		r := I(m, x) - 1
		sI(m, x, r)
		if r == 0 {
			n := I(m, x+4)
			p := x + 8
			for i := uint32(0); i < n; i++ {
				dx(m, I(m, p))
				p += 4
			}
			fr(m, x)
		}
	}
}
func bk(n uint32) (r uint32) {
	r = uint32(32 - bits.LeadingZeros32(7+4*n))
	if r < 4 {
		return 4
	}
	return r
}
func fr(m []uint32, x uint32) {
	p := 4 * bk(I(m, 4+x))
	sI(m, x, I(m, p))
	sI(m, p, x)
}
func mark(m []uint32) {
	free := func(x, t uint32) {
		for {
			if x == 0 {
				return
			}
			n := uint32(1<<(t-2) - 2)
			if bk(n) != t {
				panic("mark")
			}
			next := I(m, x)
			sI(m, x, 0)
			sI(m, x+4, n)
			x = next
		}
	}
	for i := uint32(4); i < 32; i++ {
		free(m[i], i)
	}
}
func Dump(j J, n uint32) uint32 { // type: cifzsld -> 2468ace
	m := j.M()
	fmt.Printf("%.8x ", 0)
	for i := uint32(0); i < n; i++ {
		x := m[i]
		fmt.Printf(" %.8x", x)
		if i > 0 && (i+1)%8 == 0 {
			fmt.Printf("\n%.8x ", i+1)
		} else if i > 0 && (i+1)%4 == 0 {
			fmt.Printf(" ")
		}
	}
	fmt.Println()
	return 0
}
