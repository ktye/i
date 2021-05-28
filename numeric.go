package main

import (
	. "github.com/ktye/wg/module"
)

func Neg(x K) (r K) {
	if tp(x) == It {
		nx := nn(x)
		r = mk(It, nx)
		rp := int32(r)
		xp := int32(x)
		for i := int32(0); i < nx; i++ {
			SetI32(rp, -I32(xp))
			xp += 4
			rp += 4
		}
		dx(x)
		return r
	}
	return Ki(-int32(x))
}

func Sqr(x K) K { trap(Nyi); return x }

func nd(x, y K, f int32) (r K) {
	x, y = uptypes(x, y)
	t, n, av := conform(x, y)
	switch t - 2 {
	case 0: // ct
		r = ndc(x, y, f, n, av)
	case 1: // it
		r = ndi(x, y, f, n, av)
	case 2: // ft
		r = ndf(x, y, 1+f, n, av)
	case 3: // zt
		r = ndz(x, y, 2+f, n, av)
	default:
		trap(Type)
	}
	return r
}
func nc(x, y K, f int32) (r K) {
	x, y = uptypes(x, y)
	t, n, av := conform(x, y)
	switch t - 2 {
	case 0: // ct
		r = ndc(x, y, f, n, av)
		if av == 3 {
			r = Kb(int32(r))
		} else {
			r = K(int32(r)) | K(Bt)<<59
		}
	case 1: // it
		r = nci(x, y, f, n, av)
	case 2: // ft
		r = ncf(x, y, 1+f, n, av)
	case 3: // zt
		r = ncz(x, y, 2+f, n, av)
	default:
		trap(Type)
	}
	return r
}
func ndc(x, y K, f, n, av int32) (r K) {
	xp := int32(x)
	yp := int32(y)
	if av == 3 {
		return Kc(Func[f].(f2i)(xp, yp))
	}
	switch av {
	case 0:
		r = use2(x, y)
		rp := int32(r)
		for i := int32(0); i < n; i++ {
			SetI8(rp, Func[f].(f2i)(I8(xp), I8(yp)))
			xp++
			yp++
			rp++
		}
	case 1:
		r = use1(y)
		rp := int32(r)
		for i := int32(0); i < n; i++ {
			SetI8(rp, Func[f].(f2i)(xp, I8(yp)))
			yp++
			rp++
		}
	case 2:
		r = use1(x)
		rp := int32(r)
		for i := int32(0); i < n; i++ {
			SetI8(rp, Func[f].(f2i)(I8(xp), yp))
			xp++
			rp++
		}
	}
	dx(x)
	dx(y)
	return r
}
func ndi(x, y K, f, n, av int32) (r K) {
	xp := int32(x)
	yp := int32(y)
	if av == 3 {
		return Ki(Func[f].(f2i)(xp, yp))
	}
	switch av {
	case 0:
		r = use2(x, y)
		rp := int32(r)
		for i := int32(0); i < n; i++ {
			SetI32(rp, Func[f].(f2i)(I32(xp), I32(yp)))
			xp += 4
			yp += 4
			rp += 4
		}
	case 1:
		r = use1(y)
		rp := int32(r)
		for i := int32(0); i < n; i++ {
			SetI32(rp, Func[f].(f2i)(xp, I32(yp)))
			yp += 4
			rp += 4
		}
	case 2:
		r = use1(x)
		rp := int32(r)
		for i := int32(0); i < n; i++ {
			SetI32(rp, Func[f].(f2i)(I32(xp), yp))
			xp += 4
			rp += 4
		}
	}
	dx(x)
	dx(y)
	return r
}
func nci(x, y K, f, n, av int32) (r K) {
	xp := int32(x)
	yp := int32(y)
	if av == 3 {
		return Kb(Func[f].(f2i)(xp, yp))
	}
	r = mk(Bt, n)
	rp := int32(r)
	switch av {
	case 0:
		for i := int32(0); i < n; i++ {
			SetI8(rp, Func[f].(f2i)(I32(xp), I32(yp)))
			xp += 4
			yp += 4
			rp++
		}
	case 1:
		for i := int32(0); i < n; i++ {
			SetI8(rp, Func[f].(f2i)(xp, I32(yp)))
			yp += 4
			rp++
		}
	case 2:
		for i := int32(0); i < n; i++ {
			SetI8(rp, Func[f].(f2i)(I32(xp), yp))
			xp += 4
			rp++
		}
	}
	dx(x)
	dx(y)
	return r
}
func ndf(x, y K, f, n, av int32) (r K) { trap(Nyi); return x }
func ndz(x, y K, f, n, av int32) (r K) { trap(Nyi); return x }
func ncf(x, y K, f, n, av int32) (r K) { trap(Nyi); return x }
func ncz(x, y K, f, n, av int32) (r K) { trap(Nyi); return x }

type f2b = func(uint64, uint64) uint64
type f2i = func(int32, int32) int32
type f2f = func(float64, float64) float64
type f2c = func(float64, float64) int32

func use2(x, y K) K {
	if I64(int32(y)-8) == 1 {
		return rx(y)
	}
	return use1(x)
}
func use1(x K) K {
	if I64(int32(x)-8) == 1 {
		return rx(x)
	}
	return mk(tp(x), nn(x))
}

// 0:vector-vector, 1:atom-vector, 2:vector-atom, 3: atom-atom
func conform(x, y K) (r T, n int32, av int32) {
	xt, yt := tp(x), tp(y)
	if xt < 16 {
		if yt < 16 {
			return xt, 1, 3
		} else {
			return xt, nn(y), 1
		}
	}
	xn := nn(x)
	if yt < 16 {
		return yt, xn, 2
	}
	if nn(y) != xn {
		trap(Length)
	}
	return xt - 16, xn, 0
}

func Add(x, y K) K              { return nd(x, y, 220) }
func addi(x, y int32) int32     { return x + y }
func addf(x, y float64) float64 { return x + y }

func Sub(x, y K) K              { return nd(x, y, 223) }
func subi(x, y int32) int32     { return x - y }
func subf(x, y float64) float64 { return x - y }

func Mul(x, y K) K              { return nd(x, y, 226) }
func muli(x, y int32) int32     { return x * y }
func mulf(x, y float64) float64 { return x * y }

func Div(x, y K) K              { return nd(x, y, 229) }
func divi(x, y int32) int32     { return x / y }
func divf(x, y float64) float64 { return x / y }

func Min(x, y K) K { return nd(x, y, 235) }
func mini(x, y int32) int32 {
	if x < y {
		return x
	}
	return y
}
func minf(x, y float64) float64 { return F64min(x, y) }

func Max(x, y K) K { return nd(x, y, 238) }
func maxi(x, y int32) int32 {
	if x > y {
		return x
	}
	return y
}
func maxf(x, y float64) float64 { return F64max(x, y) }

func Les(x, y K) K           { return nc(x, y, 241) }
func lti(x, y int32) int32   { return ib(x < y) }
func ltf(x, y float64) int32 { return ib(x < y) }

func Mor(x, y K) K           { return nc(x, y, 244) }
func gti(x, y int32) int32   { return ib(x > y) }
func gtf(x, y float64) int32 { return ib(x > y) }

func Eql(x, y K) K           { return nc(x, y, 247) }
func eqi(x, y int32) int32   { return ib(x == y) }
func eqf(x, y float64) int32 { return ib(isnan(x) && isnan(y) || x == y) }
func eqz()                   { trap(Nyi) }
func isnan(x float64) bool   { return x != x }

func uptypes(x, y K) (K, K) {
	xt, yt := tp(x)&15, tp(y)&15
	rt := T(maxi(int32(xt), int32(yt)))
	if rt == bt {
		rt = it
	}
	if xt < rt {
		x = uptype(x, rt)
	}
	if yt < rt {
		y = uptype(y, rt)
	}
	return x, y
}
func uptype(x K, dst T) (r K) {
	xt := tp(x)
	xp := int32(x)
	if xt < 16 {
		if dst == ct {
			return Kc(xp)
		} else if dst == it {
			return Ki(xp)
		} else if dst == ft {
			return Kf(float64(xp))
		} else if dst == zt {
			if xt == ft {
				return Kz(x, 0)
			} else {
				return Kz(Kf(float64(xp)), 0)
			}
		} else {
			trap(Type)
		}
	}
	if xt < It && dst == ft {
		x, xt = uptype(x, it), It
	}
	if xt < Ft && dst == zt {
		x, xt = uptype(x, ft), Ft
	}
	xn := nn(x)
	r = mk(dst+16, xn)
	rp := int32(r)
	if dst == ct {
		Memorycopy(rp, xp, xn)
	} else if dst == it {
		for i := int32(0); i < xn; i++ {
			SetI32(rp, I8(xp))
			xp++
			rp += 4
		}
	} else if dst == ft {
		for i := int32(0); i < xn; i++ {
			SetF64(rp, float64(I32(xp)))
			xp += 4
			rp += 8
		}
	} else if dst == zt {
		return KZ(rx(x), 0)
	} else {
		trap(Type)
	}
	dx(x)
	return r
}
