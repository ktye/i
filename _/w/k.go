// +build ignore

package main

import (
	"fmt"
	"math/bits"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

type c = byte
type s = string
type i = uint32
type j = uint64
type f = float64

var MC []c // MC, MI, MJ, MF share array (see msl)
var MI []i
var MJ []j
var MF []f
var WT []i

type vt1 func(i) i
type vt2 func(i, i) i
type slice struct {
	p uintptr
	l int
	c int
}

const m0 = 16

func main() {
	fn1 := map[string]vt1{"ini": ini, "mki": mki, "til": til, "rev": rev}
	fn2 := map[string]vt2{"mk": mk, "dump": dump}
	MC = make([]c, 1<<m0)
	msl()
	stack := make([]i, 0)
	args := []string{strconv.Itoa(m0), "ini"}
	for _, a := range append(args, os.Args[1:]...) {
		if n, e := strconv.ParseUint(a, 10, 32); e == nil {
			stack = append(stack, i(n))
		} else {
			if f1, o := fn1[a]; o {
				r := f1(stack[len(stack)-1])
				fmt.Printf("%s %d: x%x\n", a, stack[len(stack)-1], r)
				stack[len(stack)-1] = r
				continue
			}
			if f2, o := fn2[a]; o {
				x, y := stack[len(stack)-2], stack[len(stack)-1]
				r := f2(x, y)
				fmt.Printf("%s %d %d: x%x\n", a, x, y, r)
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
	if len(stack) != 2 {
		panic("stack")
	}
	fmt.Println(kst(stack[len(stack)-1]))
	dx(stack[len(stack)-1])
	leak()
}
func ini(x i) i {
	sJ(0, 1130366807310592)
	sI(128, x)
	p := i(256)
	for i := i(8); i < x; i++ {
		sI(4*i, p)
		sI(p, i)
		p *= 2
	}
	return x
}
func msl() { // update slice headers after set/inc MC
	cp := *(*slice)(unsafe.Pointer(&MC))
	ip := *(*slice)(unsafe.Pointer(&MI))
	jp := *(*slice)(unsafe.Pointer(&MJ))
	fp := *(*slice)(unsafe.Pointer(&MF))
	ip.l, ip.c, ip.p = cp.l/2, cp.c/2, cp.p
	jp.l, jp.c, jp.p = ip.l/2, ip.c/2, ip.p
	fp.l, fp.c, fp.p = jp.l, jp.c, jp.p
	MI = *(*[]i)(unsafe.Pointer(&ip))
	MJ = *(*[]j)(unsafe.Pointer(&jp))
	MF = *(*[]f)(unsafe.Pointer(&fp))
	// todo Z
}
func bk(t, n i) i { return i(32 - bits.LeadingZeros32(7+n*i(C(t)))) }
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
	sI(i, I(a))
	for j := i - 4; j >= 4*t; j -= 4 {
		u := a + 1<<(j>>2)
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
func dx(x i) {
	if x > 255 {
		xr := I(x + 4)
		sI(x+4, xr-1)
		if xr == 1 {
			fr(x)
		}
	}
}
func dxr(x, r i) i          { dx(x); return r }
func mki(i i) (r i)         { r = mk(2, 1); sI(r+8, i); return r }
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
func rev(x i) (r i) {
	t, n, xp := v1(x)
	if t == 0 {
		return x
	} else if t < 1 || t > 5 {
		panic("nyi")
	}
	r = mk(t, n)
	w := i(C(t))
	rp := 8 + r + w*(n-1)
	for k := i(0); k < n; k++ {
		copy(MC[rp:rp+w], MC[xp:xp+w])
		xp, rp = xp+w, rp-w
	}
	return dxr(x, r)
}
func trap() { panic("trap") }
func dump(a, n i) i {
	p := a >> 2
	fmt.Printf("%.8x ", a)
	for i := i(0); i < n; i++ {
		x := MI[p+i]
		fmt.Printf(" %.8x", x)
		if i > 0 && (i+1)%8 == 0 {
			fmt.Printf("\n%.8x ", a+4*i+4)
		} else if i > 0 && (i+1)%4 == 0 {
			fmt.Printf(" ")
		}
	}
	fmt.Println()
	return 0
}
func C(a i) c     { return MC[a] } // global get, e.g. I i
func I(a i) i     { return MI[a>>2] }
func J(a i) j     { return MJ[a>>3] }
func F(a i) f     { return MF[a>>3] }
func sC(a i, v c) { MC[a] = v } // global set, e.g. i::v
func sI(a i, v i) { MI[a>>2] = v }
func sJ(a i, v j) { MJ[a>>3] = v }
func sF(a i, v f) { MF[a>>3] = v }
func atoi(s string) i {
	if x, e := strconv.Atoi(s); e == nil {
		return i(x)
	}
	panic("atoi")
}
func mark() { // mark bucket type within free blocks
	for t := i(4); t < 32; t++ {
		for p := MI[t] >> 2; p != 0; p = MI[p] >> 2 {
			MI[2+p] = t
		}
	}
}
func leak() error {
	mark()
	p := i(64)
	for p < i(len(MI)) {
		if MI[p+1] != 0 {
			return fmt.Errorf("non-free block at %d(%x)", p<<2, p<<2)
		}
		t := MI[p+2]
		if t < 4 || t > 31 {
			return fmt.Errorf("illegal bucket type %d at %d(%x)", t, p<<2, p<<2)
		}
		dp := i(1) << t
		p += dp >> 2
	}
	return nil
}
func kst(x i) s {
	a := MI[x>>2]
	t, n := a>>29, a&536870911
	if t != 2 {
		panic("nyi: ~I")
	}
	r := make([]s, n)
	for i := i(0); i < n; i++ {
		r[i] = strconv.Itoa(int(int32(MI[2+i+x>>2])))
	}
	return strings.Join(r, " ")
}
