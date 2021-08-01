package main

import (
	. "github.com/ktye/wg/module"
)

func Cal(x, y K) (r K) {
	xt := tp(x)
	yt := tp(y)
	if tp(x) == lf && yt > Lt {
		return callenv(x, y)
	}
	y = explode(y)
	if isfunc(xt) != 0 {
		return cal(x, y)
	}
	return atdepth(x, y)
}
func isfunc(t T) int32 { return I32B(t == 0 || (t < 16 && t > tt)) }

func cal(f, x K) (r K) {
	t := tp(f)
	fp := int32(f)
	xn := nn(x)
	xp := int32(x)
	var z K
	if t < df {
		switch xn - 1 {
		case 0:
			x = Fst(x)
		case 1:
			x, r = spl2(x)
		case 2:
			x, r, z = spl3(x)
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
			r = Func[fp+192].(f4)(x, r, f, z)
		default:
			r = trap(Rank)
		}
		r = r
	case 1: // cf
		switch xn - 1 {
		case 0:
			r = calltrain(f, x, 0)
		case 1:
			r = calltrain(f, x, r)
		default:
			r = trap(Rank)
		}
		r = r
	case 2: // df
		//d := K(I64(fp))
		d := x0(fp)
		a := 85 + int32(I64(fp+8))
		r = Func[a].(f2)(d, x)
	case 3: // pf
		r = callprj(f, x)
	case 4: // lf
		r = lambda(f, x)
	default:
		r = trap(Type)
	}
	dx(f)
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
	return r
}
func callprj(f, x K) K {
	n := nn(x)
	fn := nn(f)
	if fn != n {
		if n < fn {
			rx(f)
			return prj(f, x)
		}
		trap(Rank)
	}
	p, l, i := spl3(rx(f))
	x = stv(l, i, x)
	x = Cal(p, x)
	return x
}
func callenv(f, x K) (r K) { // {..}.d  {..}.t
	r, x = spl2(x)
	if tp(r) != St {
		trap(Type)
	}
	n := nn(r)
	fp := int32(f)
	r = slam(l3(x0(fp), Unq(ucat(r, x1(fp))), x2(fp)), n)
	dx(f)
	return cal(r, explode(x))
}
func lambda(f K, x K) (r K) {
	fn := nn(f)
	xn := nn(x)
	if xn < fn {
		rx(f)
		return prj(f, x)
	}
	if xn != fn {
		trap(Rank)
	}
	fp := int32(f)
	c := K(I64(fp))
	lo := K(I64(fp + 8))
	nl := nn(lo)
	sa := mk(It, 2*nl) //K(I64(fp + 16))
	sp := int32(sa)
	vp := I32(8)
	lp := int32(lo)
	xp := int32(x)
	rl(x)
	dx(x)
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
	r = exec(rx(c))
	vp = I32(8)
	sp = int32(sa)
	lp = int32(lo)
	for i := int32(0); i < nl; i++ {
		p := vp + I32(lp)
		dx(K(I64(p)))
		SetI64(p, I64(sp))
		SetI64(sp, 0)
		lp += 4
		sp += 8
	}
	dx(sa)
	pp, pe = spp, spe
	return r
}
func com(x, y K) K { return K(int32(l2(y, x))) | K(cf)<<59 } // compose
func prj(f, x K) (r K) { // project
	if isfunc(tp(f)) == 0 {
		return atdepth(f, x)
	}
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
	if tp(f) == pf { // collapse
		var y K
		f, r, y = spl3(f)
		x = stv(r, rx(y), x)
		a = Drp(a, y)
	}
	r = l3(f, x, a)
	SetI32(int32(r)-12, an)
	return K(int32(r)) | K(pf)<<59
}
func arity(f K) int32 {
	t := tp(f)
	if t > df {
		return nn(f)
	}
	return 2
}
