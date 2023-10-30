package main

import (
	. "github.com/ktye/wg/module"
)

func Cal(x, y K) K {
	xt := tp(x)
	y = explode(y)
	if isfunc(xt) != 0 {
		return cal(x, y)
	}
	return atdepth(x, y)
}
func isfunc(t T) int32 { return I32B(t == 0 || (t < 16 && t > tt)) }

func cal(f, x K) K {
	r := K(0)
	z := K(0)
	y := K(0)
	t := tp(f)
	fp := int32(f)
	xn := nn(x)
	if t < df {
		switch xn - 1 {
		case 0:
			x = Fst(x)
		case 1:
			r = x1(x)
			x = r0(x)
		default:
			r = x1(x)
			z = x2(x)
			if xn == 4 {
				y = x3(x)
			}
			x = r0(x)
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
			r = Func[fp+192].(f4)(x, r, z, y)
		default:
			trap() //rank
			r = 0
		}
		r = r
	case 1: // cf
		switch xn - 1 {
		case 0:
			r = calltrain(f, x, 0)
		case 1:
			r = calltrain(f, x, r)
		default:
			trap() //rank
			r = 0
		}
		r = r
	case 2: // df
		d := x0(f)
		a := 85 + int32(I64(fp+8))
		r = Func[a].(f2)(d, x)
	case 3: // pf
		r = callprj(f, x)
	case 4: // lf
		r = lambda(f, x)
	case 5: // xf
		r = native(f, x)
	default:
		trap() //type
		r = 0
	}
	dx(f)
	return r
}

func calltrain(f, x, y K) K {
	r := K(0)
	n := nn(f)
	if y == 0 {
		r = cal(x0(f), l1(x))
	} else {
		r = cal(x0(f), l2(x, y))
	}
	for n > 1 {
		n--
		r = cal(x0(f+8), l1(r))
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
		trap() //rank
	}
	return Cal(x0(f), stv(x1(f), x2(f), x))
}
func native(f K, x K) K {
	fn := nn(f)
	xn := nn(x)
	if xn < fn {
		rx(f)
		return prj(f, x)
	}
	if xn != fn {
		trap() //rank
	}
	return K(Native(int64(x0(f)), int64(x))) // +/api: KR
}
func lambda(f K, x K) K {
	r := lac(f, x)
	if r != 0 {
		return r
	}
	pro(f, x)
	r = exec(x0(f))
	epi()
	return r
}
func lac(f K, x K) K { //test if lambda call has correct number of args
	fn := nn(f)
	xn := nn(x)
	if xn < fn {
		rx(f)
		return prj(f, x)
	}
	if xn != fn {
		trap() //rank
	}
	return 0
}
func pro(f K, x K) { //save(push) locals, assign args
	lo := K(I64(int32(f) + 8))
	n := nn(lo)
	a := nn(f)
	z := mk(Zt, n) //use a complex vector to store symbols+values w/o refcounting
	zp := int32(z)
	xp := ep(x)
	vp := I32(8)
	for n > 0 {
		n -= 1
		p := I32(int32(lo) + 4*n)
		SetI32(zp, p)
		p += vp
		SetI64(zp+8, I64(p))
		if n < a { //args
			xp -= 8
			SetI64(p, I64(xp))
		} else { //locals
			SetI64(p, 0)
		}
		zp += 16
	}
	rl(x)
	dx(x)
	push(z)
}
func epi() { //restore saved arguments
	z := pop()
	if tp(z) != Zt {
		trap() //type/stack
	}
	zp := int32(z)
	e := ep(z)
	for zp < e {
		p := I32(8) + I32(zp)
		dx(K(I64(p)))
		SetI64(p, I64(zp+8))
		zp += 16
	}
	dx(z)
}

func com(x, y K) K { return ti(cf, int32(l2(y, x))) } // compose
func prj(f, x K) K { // project
	r := K(0)
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
		r = x1(f)
		y := x2(f)
		f = r0(f)
		x = stv(r, rx(y), x)
		a = Drp(a, y)
	}
	r = l3(f, x, a)
	SetI32(int32(r)-12, an)
	return ti(pf, int32(r))
}
func arity(f K) int32 {
	if tp(f) > df {
		return nn(f)
	}
	return 2
}
