// +build ignore

package main

import (
	"fmt"
	"math/bits"
	"os"
	"strconv"
	"unsafe"
)

type c = byte
type i = uint32
type j = uint64
type f = float64

var MC []c // MC, MI, MJ, MF share array (see msl)
var MI []i
var MJ []j
var MF []f

type vt1 func(i) i
type vt2 func(i, i) i
type slice struct {
	p uintptr
	l int
	c int
}

func main() {
	fn1 := map[string]vt1{"ini": ini, "mki": mki, "til": til}
	fn2 := map[string]vt2{"mk": mk, "dump": dump}
	MC = make([]c, 64*1024)
	msl()
	stack := make([]i, 0)
	args := []string{"16", "ini"}
	for _, a := range append(args, os.Args[1:]...) {
		if n, e := strconv.ParseUint(a, 10, 32); e == nil {
			stack = append(stack, i(n))
		} else {
			if f1, o := fn1[a]; o {
				r := f1(stack[len(stack)-1])
				fmt.Printf("%s %d: %d\n", a, stack[len(stack)-1], r)
				stack[len(stack)-1] = r
				continue
			}
			if f2, o := fn2[a]; o {
				x, y := stack[len(stack)-2], stack[len(stack)-1]
				r := f2(x, y)
				fmt.Printf("%s %d %d: %d\n", a, x, y, r)
				if a == "dump" {
					stack = stack[:len(stack)-2]
				} else {
					stack = stack[:len(stack)-1]
					stack[len(stack)-1] = r
				}
			} else {
				panic("unknown func: " + a)
			}
		}
	}
}
func ini(x i) i {
	sJ(0, 1130366807310592)
	p := i(64)
	for i := i(8); i < x; i++ {
		MI
		sI(4*i, p)
		sI(p, i)
		p *= 2
	}
	sI(128, x)
	return x
}
func msl() { // update slice header after increasing MC
	c := *(*slice)(unsafe.Pointer(&MC))
	i := *(*slice)(unsafe.Pointer(&MI))
	j := *(*slice)(unsafe.Pointer(&MJ))
	f := *(*slice)(unsafe.Pointer(&MF))
	i.l, i.c, i.p = c.l/2, c.c/2, c.p
	j.l, j.c, j.p = i.l/2, i.c/2, i.p
	f.l, f.c, f.p = j.l, j.c, j.p
	MI = *(*[]I)(unsafe.Pointer(&i))
	MJ = *(*[]J)(unsafe.Pointer(&j))
	MF = *(*[]F)(unsafe.Pointer(&f))
	// todo Z
}
func bk(t, n i) i { return i(32 - bits.LeadingZeros32(7+n*C(t))) }
func mk(x, y i) i {
	t := bk(x, y)
	i := 4 * t
	for I(i) == 0 {
		i += 4
	}
	if i == 128 {
		panic("oom")
	}
	a := I(i)
	sI(i, I(4+a))
	for j := i - 4; j >= 4*t; j -= 4 {
		u := a + 1<<(j/4)
		sI(u, I(j))
		sI(j, u)
	}
	sI(a, y|x<<29)
	sI(a+4, 1)
	return a
}
func fr(x i) {
	xt, xn, _ := v1(x)
	t := 4 * bk(xt, xn)
	sI(x, I(t))
	sI(t, x)
}
func decr(x i) {
	if x > 255 {
		println("decr", x)
		xr := I(x + 4)
		sI(x+4, xr-1)
		if xr == 1 {
			fr(x)
		}
	}
}
func dxr(x, r i) i { decr(x); return r }
func mki(i i) (r i) {
	r = mk(2, 1)
	sI(r+8, i)
	return r
}
func v1(x i) (xt, xn, xp i) { u := I(x); return u >> 29, u & 536870911, 8 + x }
func til(x i) (r i) {
	xt, _, xp := v1(x)
	if xt != 2 {
		trap()
	}
	n := I(xp)
	r = mk(xt, n)
	rp := 8 + r
	for i := i(0); i < n; i++ {
		sI(rp, i)
		rp += 4
	}
	return dxr(x, r)
}
func trap() { panic("trap") }
func dump(a, n i) i {
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
	return 0
}
func hxb(x c) (c, c) { h := "0123456789abcdef"; return h[x>>4], h[x&0x0F] }
func C(o, i i) i     { return i(MC[o+i]) } // global get, e.g. I o,i
func I(o, i i) i     { return MI[o>>2+i] }
func J(o, i i) j     { return MJ[o>>3+i] }
func F(o, i i) f     { return MF[o>>3+i] }
func sC(o, i i, v c) { MC[o+i] = v } // global set, e.g. (o,i)::v
func sI(o, i i, v i) { MI[o<<2+i] = v }
func sJ(o, i i, v j) { MJ[o<<3+i] = v }
func sF(o, i i, v f) { MF[o<<3+i] = v }
func atoi(s string) i {
	if x, e := strconv.Atoi(s); e == nil {
		return i(x)
	}
	panic("atoi")
}
