package x

import (
	"fmt"
	"math/bits"
	"strconv"
	"strings"
)

const SZ = 126 // stack size (#elements)

type J interface {
	J(x uint32) uint32
	M() []uint32
}

func X(m []uint32, x uint32) string {
	if x == 0 {
		panic("XX0")
	} else if x&7 == 0 {
		n := nn(m, x)
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
func cpy(x []uint32) (r []uint32) { r = make([]uint32, len(x)); copy(r, x); return r }
func Leak(j J) {
	m := cpy(j.M())

	drp := func(x uint32) {
		n := nn(m, x)
		for i := uint32(0); i < n; i++ {
			dx(m, m[2+i+x>>2])
		}
		m[1+x>>2] = 0
	}
	cls := func(x uint32) {
		if n := nn(m, x); n != 0 {
			panic(fmt.Errorf("stack %d is not empty: #%d", x, n))
		}
		for i := uint32(0); i < SZ; i++ {
			sI(m, x+8+4*i, 0)
		}
		sI(m, x+4, SZ)
		dx(m, x)
	}
	drp(m[1])
	cls(m[1])
	cls(m[2])

	clear := func(n uint32, y uint32) {
		x := 8 + m[3]
		if I(m, x+8*n) != y {
			panic("clear wrong symbol")
		}
		sI(m, x+4+8*n, 0)
	}
	clear(0, N)
	clear(1, Y)
	clear(2, C)
	dx(m, m[3])

	mark(m)
	p := uint32(32)
	for p < uint32(len(m)) {
		if m[p] != 0 {
			Dump(m, 200)
			panic(fmt.Errorf("non-free block: %d(%x) rc=%d #=%d", p, p, m[p], m[1+p]))
		}
		n := uint32(1 << bk(m[1+p]))
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
func Dump(m []uint32, n uint32) uint32 { // type: cifzsld -> 2468ace
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
func lu(m []uint32, y uint32) uint32 {
	p := fn(m, I(m, 12), y)
	if p == 0 {
		panic("undefined: " + X(m, y))
	}
	return I(m, p)
}
func fn(m []uint32, x, y uint32) uint32 {
	n := nn(m, x) / 2
	p := x + 8
	for i := uint32(0); i < n; i++ {
		if I(m, p) == y {
			return p + 4
		}
		p += 8
	}
	return 0
}

const N, Y, C uint32 = 110, 114, 118
