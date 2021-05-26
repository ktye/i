package k

import (
	"fmt"

	. "github.com/ktye/wg/module"
)

func Cal(x, y K) (r K) {
	fmt.Println("Cal", sK(x), sK(y))
	xt, yt := tp(x), tp(y)
	if xt < 16 {
		if xt == 0 || xt > tt {
			if yt != Lt {
				trap(Nyi)
			}
			if nn(y) != 2 {
				trap(Rank)
			}
			r = project(x, y)
			if r != 0 {
				return r
			}
			yp := int32(y)
			r = cal2(x, rx(K(I64(yp))), rx(K(I64(8+yp))))
			dx(y)
			return r
		}
	}
	panic(Nyi)
	return x
}

func cal1(f, x K) (r K) {
	t, xt := tp(f), tp(x)
	if xt == 0 || (xt > tt && xt < 16) {
		return compose(f, x, xt)
	}
	if t != 0 {
		t -= 9
	}
	switch t {
	case 0:
		r = Func[int32(f)].(f1)(x)
	case 1: // cf
		r = calltrain(f, x, 0)
	case 2: // df
		fp := int32(f)
		d := K(I64(fp))
		a := int32(I64(fp + 8))
		r = Func[85+int32(a)].(f2)(d, x)
		dx(f)
	case 3: // pf
		r = callprj(f, x, 0)
	case 4: // lf
		r = lambda(f, l1(x))
	default:
		trap(Nyi)
	}
	return r
}
func cal2(f, x, y K) (r K) {
	t := tp(f)
	if t != 0 {
		t -= 9
	}
	switch t {
	case 0:
		r = Func[int32(f)+64].(f2)(x, y)
	case 1: // cf
		r = calltrain(f, x, y)
	case 2: // df
		trap(Nyi)
	case 3: // pf
		trap(Nyi)
	case 4: // lf
		r = lambda(f, l2(x, y))
	default:
		trap(Type)
	}
	return r
}
func calltrain(f, x, y K) (r K) {
	n := nn(f)
	fp := int32(f)
	if y == 0 {
		r = cal1(rx(K(I64(fp))), x)
	} else {
		r = cal2(rx(K(I64(fp))), x, y)
	}
	for i := int32(1); i < n; i++ {
		fp += 8
		r = cal1(rx(K(I64(fp))), r)
	}
	dx(f)
	return r
}
func callprj(f, x, y K) (r K) {
	ar := nn(f)
	fp := int32(f)
	if y == 0 && ar > 1 {
		return project(f, l1(x))
	} else if ar > 2 {
		return project(f, l2(x, y))
	}
	p := rx(K(I64(fp)))
	l := rx(K(I64(fp + 8)))
	i := rx(K(I64(fp + 16)))
	if y == 0 {
		dx(i)
		r = sti(l, I32(int32(i)), x)
	} else {
		r = stv(l, i, l2(x, y))
	}
	r = Cal(p, r)
	dx(f)
	return r
}
func lambda(f K, x K) (r K) {
	fn := nn(f)
	if nn(x) > fn {
		return project(f, x)
	}
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
func compose(x, y K, yt T) (r K) {
	if yt == ct {
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
