// +build ignore

package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/bits"
	"os"
	"strconv"
)

type c = byte
type i = uint32
type j = uint64
type f = float64

var M []c

func main() {
	M = make([]c, 64*1024)
	ini(16)
	for _, a := range os.Args[1:] {
		r := mk(2, atoi(a))
		fmt.Printf("r=%x\n", r)
	}
	dump(0, 4100)
}
func ini(x i) {
	sJ(0, 1130366807310592)
	p := i(512)
	for i := i(9); i < x; i++ {
		sI(4*i, p)
		sI(p, i)
		p *= 2
		fmt.Printf("i=%d p=%d [%d]=%d [%d]=%d\n", i, p, 4*i, p, p, i)
	}
	sI(128, x)
}
func bk(t, n i) i { return i(32 - bits.LeadingZeros32(7+n*C(t))) }
func mk(x, y i) i {
	t := bk(x, y)
	i := 4 * t
	for I(i) == 0 {
		i += 4
	}
	if i == 32 {
		panic("oom")
	}
	a := I(i)
	sI(i, I(4+a))
	for j := i - 4; j >= 4*t; j -= 4 {
		u := a + 1<<(j-2)
		sI(u, I(j))
		sI(j, u)
	}
	sI(a, y|x<<29)
	sI(a+4, 1)
	return a
}

func dump(a, n i) {
	fmt.Printf("%.8x  ", 0)
	for i, b := range M[a : a+n] {
		hi, lo := hxb(b)
		fmt.Printf("%c%c", hi, lo)
		if i > 0 && (i+1)%32 == 0 {
			fmt.Printf("\n%.8x  ", i+1)
		} else if i > 0 && (i+1)%16 == 0 {
			fmt.Printf("  ")
		} else if i > 0 && (i+1)%4 == 0 {
			fmt.Printf(" ")
		}
	}
	fmt.Println()
}
func hxb(x c) (c, c) { h := "0123456789abcdef"; return h[x>>4], h[x&0x0F] }
func C(a i) i        { return i(M[a]) }
func I(a i) i        { return binary.LittleEndian.Uint32(M[a : a+4]) }
func J(a i) j        { return binary.LittleEndian.Uint64(M[a : a+8]) }
func F(a i) f        { return math.Float64frombits(J(a)) }
func sI(a i, v i)    { binary.LittleEndian.PutUint32(M[a:a+4], v) }
func sJ(a i, v j)    { binary.LittleEndian.PutUint64(M[a:a+8], v) }
func sF(a i, v f)    { binary.LittleEndian.PutUint64(M[a:a+8], math.Float64bits(v)) }
func atoi(s string) i {
	if x, e := strconv.Atoi(s); e == nil {
		return i(x)
	}
	panic("atoi")
}
