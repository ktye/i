package main

import (
	. "github.com/ktye/wg/module"
)

func Cal(x, y K) (r K) {
	xt, yt := tp(x), tp(y)
	if xt < 16 {
		if xt == 0 || xt > tt {
			if yt != Lt {
				trap(Nyi)
			}
			r = project(x, y)
			if r != 0 {
				return r
			}
			return cal(x, y)
		}
	}
	panic(Nyi)
	return x
}
func cal(f, x K) (r K) {
	t := tp(f)
	fp := int32(f)
	xn := nn(x)
	xp := int32(x)
	var z K
	if t < pf {
		switch xn - 1 {
		case 0:
			x = Fst(x)
		case 1:
			x, r = spl2(x)
		case 2:
			x, r, z = spl3(x)
		default:
		}
	}
	if t != 0 {
		t -= 9
	}
	switch t {
	case 0: // basic
		switch xn - 1 {
		case 0:
			r = Func[int32(f)].(f1)(x)
		case 1:
			r = Func[fp+64].(f2)(x, r)
		case 2:
			r = Func[fp+192].(f4)(x, r, 1, z)
		case 3:
			z := x3(xp)
			x, r, f = spl3(x)
			r = Func[fp+192].(f4)(x, r, z, f)
		default:
			trap(Rank)
		}
	case 1: // cf
		switch xn - 1 {
		case 0:
			r = calltrain(f, x, 0)
		case 1:
			r = calltrain(f, x, r)
		default:
			trap(Rank)
		}
	case 2: // df
		d := K(I64(fp))
		a := 85 + int32(I64(fp+8))
		switch xn - 1 {
		case 0:
			r = Func[a].(f2)(d, x)
		case 1:
			r = Func[64+a].(f3)(d, x, r)
		default:
			trap(Rank)
		}
	case 3: // pf
		r = callprj(f, x)
	case 4: // lf
		r = lambda(f, x)
	default:
		trap(Type)
	}
	return r
}
func calltrain(f, x, y K) (r K) {
	n := nn(f)
	fp := int32(f)
	if y == 0 {
		r = cal(x0(fp), l1(x))
	} else {
		r = cal(x0(fp), l2(x, y))
	}
	for i := int32(1); i < n; i++ {
		fp += 8
		r = cal(x0(fp), l1(r))
	}
	dx(f)
	return r
}
func callprj(f, x K) K {
	n := nn(x)
	if nn(f) != n {
		trap(Rank)
	}
	var l, i K
	p, fp := nxl(int32(f))
	l, fp = nxl(fp)
	i, fp = nxl(fp)
	x = stv(rx(l), rx(i), x)
	x = Cal(p, x)
	dx(f)
	return x
}
func lambda(f K, x K) (r K) {
	fn := nn(f)
	fp := int32(f)
	c := K(I64(fp))
	lo := K(I64(fp + 8))
	sa := K(I64(fp + 16))
	sp := int32(sa)
	nl := nn(sa)
	vp := I32(8)
	lp := int32(lo)
	xp := int32(x)
	for i := int32(0); i < nl; i++ {
		p := vp + I32(lp)
		SetI64(sp, I64(p))
		if i < fn {
			SetI64(p, I64(xp))
			xp += 8
		} else {
			SetI64(p, 0)
		}
		sp += 8
		lp += 4
	}
	spp, spe := pp, pe
	r = exec(c)
	vp = I32(8)
	sp = int32(sa)
	for i := int32(0); i < nl; i++ {
		p := vp + I32(lp)
		if i < fn {
			dx(K(I64(p)))
		}
		SetI64(p, I64(sp))
		lp += 4
		sp += 8
	}
	pp, pe = spp, spe
	return r
}
func compose(x, y K) (r K) {
	if tp(y) == ct {
		r = cat1(K(int32(y))|K(Lt)<<59, x)
	} else {
		r = l2(y, x)
	}
	return K(int32(r)) | K(cf)<<59
}
func project(f, x K) (r K) {
	xn := nn(x)
	xp := int32(x)
	a := mk(It, 0)
	for i := int32(0); i < xn; i++ {
		if I64(xp) == 0 {
			a = cat1(a, Ki(i))
		}
		xp += 8
	}
	ar := arity(f)
	for i := xn; i < ar; i++ {
		a = cat1(a, Ki(i))
		x = cat1(x, 0)
	}
	an := nn(a)
	if an == 0 {
		dx(a)
		return 0
	}
	r = l3(f, x, a)
	SetI32(int32(r)-4, an)
	return K(int32(r)) | K(pf)<<59
}
func arity(f K) int32 {
	t := tp(f)
	if t > df {
		return nn(f)
	}
	return 2
}
