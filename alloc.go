package k

import . "github.com/ktye/wg/module"

func init() {
	Memory(1)
}
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
	r = K(uint64(t) << uint64(57))
	x := alloc(n * sz(r))
	SetI32(x, 1)
	SetI32(x+4, n)
	return r | K(x+8)
}
func tp(x K) T         { return T(x >> 57) }
func nn(x K) int32     { return I32(int32(x) - 4) }
func sz(x K) int32     { return 1 << I32clz(uint32(x)) }
func nocount(x K) bool { return x < 64 || (((1 & x >> 56) != 0) && (T(x>>57) != ft)) }
func rx(x K) K {
	if nocount(x) {
		return x
	}
	p := int32(x)
	SetI32(p-8, 1+I32(p-8))
	return x
}
func dx(x K) {
	if nocount(x) {
		return
	}
	a := int32(x - 8)
	rc := I32(a)
	SetI32(a, rc-1)
	if rc == 0 {
		trap(Unref)
	}
	if rc == 1 {
		xt := tp(x)
		if xt <= tt {
			p := int32(x)
			for n := nn(x); n != 0; n-- {
				dx(K(I64(p)))
				p += 8
			}
		}
		free(a, bucket(sz(x)*nn(x)))
	}
}
