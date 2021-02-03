package j

import (
	"fmt"
	"strconv"
	"strings"
)

func Dx(x uint32) uint32 { return dx(x) }

func X(x uint32) uint32 {
	if x&7 == 0 {
		n := nn(x)
		fmt.Printf("%x[%d]:", x, n)
		for i := uint32(0); i < n; i++ {
			fmt.Printf(" %x", M[2+i+x>>2])
		}
		fmt.Println()
	} else {
		fmt.Printf("%x: %d\n", x>>1)
	}
	return x
}

func XX(x uint32) string {
	//fmt.Printf("xx %d(%x) %b\n", x, x, x)
	if x == 0 {
		panic("XX0")
	} else if x&7 == 0 {
		n := nn(x)
		v := make([]string, n)
		for i := uint32(0); i < n; i++ {
			v[i] = XX(M[2+i+x>>2])
		}
		return "[" + strings.Join(v, " ") + "]"
	} else if x&1 != 0 {
		return strconv.Itoa(int(int32(x) >> 1))
	} else if x&2 != 0 {
		return "<" + strconv.Itoa(int(x)) + ">"
	} else if x&4 != 0 {
		return string([]byte{33 + byte(x>>3)})
	}
	panic("XX")
}

func Leak() {
	dx(M[1])
	dx(M[2])
	dx(M[3])
	//dump(200)
	mark()
	//dump(200)
	p := uint32(32)
	for p < uint32(len(M)) {
		if M[p] != 0 {
			panic(fmt.Errorf("non-free block: %d(%x)", 4*p, 4*p))
		}
		n := uint32(1 << bk(M[1+p]))
		p += n >> 2
	}
	fmt.Println("no leak")
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
