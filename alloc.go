package main

import (
	. "github.com/ktye/wg/module"
)

func minit(a, b int32) {
	p := int32(1 << a)
	for i := a; i < b; i++ {
		SetI32(4*i, p)
		SetI32(p, 0)
		p *= 2
	}
	SetI32(128, b)
}
func alloc(n, s int32) int32 {
	size := n * s
	t := bucket(size)
	if int64(n)*int64(s) > 2147483647 /*|| t > 31*/ {
		trap() //grow (oom)
	}
	i := 4 * t
	m := 4 * I32(128)
	for I32(i) == 0 {
		if i >= m {
			m = 4 * grow(i)
		} else {
			i += 4
		}
	}
	a := I32(i)
	SetI32(i, I32(a))
	for j := i - 4; j >= 4*t; j -= 4 {
		u := a + int32(1)<<(j>>2)
		SetI32(u, I32(j))
		SetI32(j, u)
	}
	if a&31 != 0 {
		trap() //memory corruption
	}
	return a
}
func grow(p int32) int32 {
	m := I32(128)                       // old total memory (log2)
	n := 1 + (p >> 2)                   // required total mem (log2)
	g := (1 << (n - 16)) - Memorysize() // grow by 64k blocks

	if g > 0 {
		if Memorygrow(g) < 0 {
			trap() //grow
		}
	}
	minit(m, n)
	return n
}
func mfree(x, bs int32) {
	if x&31 != 0 {
		trap() //memory corruption
	}
	t := 4 * bs
	SetI32(x, I32(t))
	SetI32(t, x)
}
func bucket(size int32) int32 {
	r := int32(32) - I32clz(15+size)
	if r < 5 {
		r = 5
	}
	return r
}
func mk(t T, n int32) K {
	if t < 17 {
		trap() //type
	}
	x := alloc(n, sz(t))
	SetI32(x+12, 1) //rc
	SetI32(x+4, n)
	return ti(t, x+16)
}
func tp(x K) T     { return T(uint64(x) >> 59) }
func nn(x K) int32 { return I32(int32(x) - 12) }
func ep(x K) int32 { return int32(x) + sz(tp(x))*nn(x) }
func sz(t T) int32 {
	if t < 16 {
		return 8
	} else if t < 19 {
		return 1
	} else if t < 21 {
		return 4
	} else if t == Zt {
		return 16
	}
	return 8
}
func rx(x K) K {
	if tp(x) < 5 {
		return x
	}
	p := int32(x) - 4
	SetI32(p, 1+I32(p))
	return x
}
func dx(x K) {
	t := tp(x)
	if t < 5 {
		return
	}
	p := int32(x) - 16
	rc := I32(p + 12)
	SetI32(p+12, rc-1)
	if rc == 0 {
		trap() //unref
	}
	if rc == 1 {
		n := nn(x)
		if t&15 > 6 {
			if t == 14 || t == 24 || t == 25 {
				n = 2 // nat | D | T
			} else if t == 12 || t == 13 {
				n = 3 // prj | lam
			}
			p := int32(x)
			e := p + 8*n
			for p < e {
				dx(K(I64(p)))
				p += 8
			}
		}
		mfree(p, bucket(sz(t)*n))
	}
}
func rl(x K) { // ref list elements
	e := ep(x)
	p := int32(x)
	for e > p {
		e -= 8
		rx(K(I64(e)))
	}
}
func lfree(x K) { // free list non-recursive
	mfree(int32(x)-16, bucket(8*nn(x)))
}
