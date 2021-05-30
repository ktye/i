package main

import (
	. "github.com/ktye/wg/module"
)

func minit(a, b int32) {
	p := int32(1 << a)
	for i := a; i < b; i++ {
		SetI32(4*i, p)
		p *= 2
	}
	SetI32(128, b)
}
func alloc(size int32) int32 {
	t := bucket(size)
	i := 4 * t
	m := 4 * I32(128)
	for I32(i) == 0 {
		if i >= m {
			trap(Grow)
		}
		i += 4
	}
	a := I32(i)
	SetI32(i, I32(a))
	for j := i - 4; j >= 4*t; j -= 4 {
		u := a + 1<<(j>>2)
		SetI32(u, I32(j))
		SetI32(j, u)
	}
	return a
}
func free(x, bs int32) {
	t := 4 * bs
	SetI32(x, I32(t))
	SetI32(t, x)
}
func bucket(size int32) (r int32) {
	r = 32 - I32clz(uint32(7+size))
	if r < 4 {
		r = 4
	}
	return r
}

func mk(t T, n int32) (r K) {
	if t < 17 {
		trap(Value)
	}
	r = K(uint64(t) << uint64(59))
	x := alloc(n * sz(t))
	SetI32(x, 1)
	SetI32(x+4, n)
	return r | K(x+8)
}
func tp(x K) T     { return T(x >> 59) }
func nn(x K) int32 { return I32(int32(x) - 4) }
func sz(t T) int32 {
	if t < 19 {
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
	p := int32(x) - 8
	SetI32(p, 1+I32(p))
	return x
}
func dx(x K) {
	t := tp(x)
	if t < 5 {
		return
	}
	p := int32(x) - 8
	rc := I32(p)
	SetI32(p, rc-1)
	if rc == 0 {
		trap(Unref)
	}
	if rc == 1 {
		n := nn(x)
		if t&15 > 6 {
			if t == 22 || t == 24 || t == 25 {
				n = 2
			} else if t == 12 {
				n = 3 // prj
			} else if t == 13 {
				n = 4 // lam
			}
			p := int32(x)
			for i := int32(0); i < n; i++ {
				dx(K(I64(p)))
				p += 8
			}
		}
		free(p, bucket(sz(t)*n))
	}
}
func rl(x K) { // ref list elements
	xp := int32(x)
	xn := nn(x)
	for i := int32(0); i < xn; i++ {
		x0(xp)
		xp += 8
	}
}
