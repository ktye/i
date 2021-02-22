package j

import (
	"fmt"
	"strconv"
	"strings"
)

type jj struct{}

func (o jj) J(x uint32) uint32 { return J(x) }
func (o jj) M() []uint32       { return M }

type jer interface {
	J(x uint32) uint32
	M() []uint32
}

func XX(x uint32) string {
	if x&7 != 0 {
		return fmt.Sprintf("nolist(%d)", x)
	}
	n := nn(x)
	xp := 8 + x
	s := ""
	for i := uint32(0); i < n; i++ {
		xi := I(xp)
		if xi == 0 {
			s += "<NULL> "
		} else if xi&7 == 0 {
			s += fmt.Sprintf("L(%d %d/%d) ", xi, I(xi), I(xi+4))
		} else if xi&1 != 0 {
			s += fmt.Sprintf("I(%d) ", xi>>1)
		} else if xi&2 != 0 {
			s += fmt.Sprintf("<%d> ", xi>>2)
		} else {
			s += fmt.Sprintf("E(%d) ", xi)
		}
		xp += 4
	}
	return s
}
func X(x uint32) string { return SX(M, x) }
func SX(m []uint32, x uint32) string {
	if x == 0 {
		panic("XX0")
	} else if x&7 == 0 {
		n := ln(m, x)
		v := make([]string, n)
		for i := uint32(0); i < n; i++ {
			v[i] = SX(m, m[2+i+x>>2])
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
func refcount(x uint32) string {
	rc := func() int {
		if x&7 == 0 {
			return int(I(x))
		}
		return -1
	}
	return fmt.Sprintf("rc(%d)", rc())
}
func cls() {
	s := I(4)
	n := nn(s)
	for i := uint32(0); i < n; i++ {
		dx(I(s + 8 + 4*i))
	}
	sI(4+s, 0)
}
func Leak(j jer) {
	B := M
	m := j.M()
	M = make([]uint32, len(m))
	copy(M, m)
	defer func() { M = B }()

	blank := func(x, n uint32) {
		sI(x+4, n)
		p := x + 8
		for i := uint32(0); i < n; i++ {
			sI(p, 0)
			p += 4
		}
		dx(x)
	}
	if n := nn(M[1]); n != 0 {
		panic("stack is not clear: #" + strconv.Itoa(int(n)))
	}
	blank(M[1], sz)
	parse := M[2]
	root := I(parse + 8)
	dx(root) // I(8) contains only 1 refcounted list
	blank(parse, 5)
	dx(M[3])

	mark()
	p := uint32(32)
	for p < uint32(len(M)) {
		if M[p] != 0 {
			panic(fmt.Errorf("non-free block: %d(%x) rc=%d #=%d", 4*p, 4*p, M[p], M[1+p]))
		}
		n := uint32(1 << bk(M[1+p]))
		p += n >> 2
	}
}
func mark() {
	free := func(x, t uint32) {
		for {
			if x == 0 {
				return
			}
			n := uint32(1<<(t-2) - 2)
			if bk(n) != t {
				panic("mark")
			}
			next := I(x)
			sI(x, 0)
			sI(x+4, n)
			x = next
		}
	}
	for i := uint32(4); i < 32; i++ {
		free(M[i], i)
	}
}
func Dump(j jer, n uint32) uint32 { // type: cifzsld -> 2468ace
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
