package j

import (
	"fmt"
	"strconv"
	"strings"
)

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
			panic("Y0")
		} else if xi&7 == 0 {
			s += fmt.Sprintf("L(%d %d/%d) ", xi, I(xi), I(xi+4))
		} else if xi&1 != 0 {
			s += fmt.Sprintf("I(%d) ", xi>>1)
		} else {
			s += fmt.Sprintf("E(%d) ", xi)
		}
		xp += 4
	}
	return s
}
func X(x uint32) string {
	//fmt.Printf("xx %d(%x) %b\n", x, x, x)
	if x == 0 {
		panic("XX0")
	} else if x&7 == 0 {
		n := nn(x)
		v := make([]string, n)
		for i := uint32(0); i < n; i++ {
			v[i] = X(M[2+i+x>>2])
		}
		return "[" + strings.Join(v, " ") + "]"
	} else if x&1 != 0 {
		return strconv.Itoa(int(int32(x) >> 1))
	} else if x&2 != 0 {
		// return "<" + strconv.Itoa(int(x)) + ">"
		return sy(x)
	} else if x&4 != 0 {
		return string([]byte{33 + byte(x>>3)})
	}
	panic("XX")
}
func annotate(est, rst, p uint32) {
	n := nn(est)
	ep := 8 + est
	r := make([]uint32, n)
	for i := uint32(0); i < n-1; i++ {
		r[i] = I(12+rst+4*i) >> 1
	}
	r[n-1] = p + 4

	e := "(est) "
	o := "(rst) "
	for i := uint32(0); i < n; i++ {
		de, do := anno(I(ep), r[i])
		e, o = e+de, o+do
		ep += 4
	}
	fmt.Printf("%s\n%s\n", e, o)
}
func anno(e, p uint32) (string, string) {
	n := nn(e)
	ep := e + 8
	s, o := "[", " "
	if ep == p {
		o = "^"
	}
	for i := uint32(0); i < n; i++ {
		ds := X(I(ep)) + " "
		if ep+4 == p {
			o += "^"
		} else {
			o += " "
		}
		o += strings.Repeat(" ", len(ds)-1)
		s += ds
		ep += 4
	}
	return s + "] ", o + "  "
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

func Leak() { LeakExec(P, 0, 0) }
func LeakExec(a, b, c uint32) {
	B := make([]uint32, len(M))
	copy(B, M)
	defer func() { copy(M, B) }()
	dx(a)
	dx(b)
	dx(c)
	dx(M[1])
	dx(M[2])
	//dx(P)
	//dump(200)
	mark()
	//dump(200)
	p := uint32(32)
	for p < uint32(len(M)) {
		if M[p] != 0 {
			// fmt.Println(X(4 * p))
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
func dump(n uint32) uint32 { // type: cifzsld -> 2468ace
	fmt.Printf("%.8x ", 0)
	for i := uint32(0); i < n; i++ {
		x := M[i]
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
