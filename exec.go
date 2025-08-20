package main

import (
	. "github.com/ktye/wg/module"
)

type f1 = func(K) K
type f2 = func(K, K) K
type f3 = func(K, K, K) K
type f4 = func(K, K, K, K) K

func quoted(x K) int32 { return I32B(int32(x) >= 448 && tp(x) == 0) }
func quote(x K) K      { return x + 448 }
func unquote(x K) K    { return x - 448 }

func exec(x K) K {
	var b, c K
	srcp = 0
	a := K(0) // accumulator
	xn := nn(x)
	if xn == 0 {
		dx(x)
		return 0
	}
	p := int32(x)
	e := p + 8*xn
	for p < e {
		u := K(I64(p))
		if tp(u) != 0 {
			push(a)
			a = rx(u)
		} else {
			switch int32(u) >> 6 {
			case 0: //   0..63   monadic
				a = Func[marksrc(u)].(f1)(a)
			case 1: //  64..127  dyadic
				a = Func[marksrc(u)].(f2)(a, pop())
			case 2: // 128       dyadic indirect
				marksrc(a)
				b = pop()
				a = Cal(a, l2(b, pop()))
			case 3: // 192..255  tetradic
				b = pop()
				c = pop()
				a = Func[marksrc(u)].(f4)(a, b, c, pop())
			case 4: // 256       drop
				dx(a)
				a = pop()
			case 5: // 320       jump
				p = p + int32(a)
				a = pop()
			case 6: // 384       jump if not
				u = pop()
				p += int32(a) * I32B(int32(u) == 0)
				dx(u)
				a = pop()
			default: //448..     quoted verb
				push(a)
				a = rx(u - 448)
			}
		}
		p += 8
		continue
	}
	pop() //0
	dx(x)
	return a
}
func marksrc(x K) int32 {
	if p := h48(x); p != 0 {
		srcp = p
	}
	return int32(x)
}
func push(x K) {
	SetI64(sp, int64(x))
	sp += 8
	if sp == 4096 { //512 {
		trap() //stack overflow
	}
}
func pop() K {
	sp -= 8
	if sp < 2048 {
		trap() //stack underflow
	}
	return K(I64(sp))
}
func lst(n K) K {
	r := mk(Lt, int32(n))
	rp := int32(r)
	e := ep(r)
	for rp < e {
		SetI64(rp, int64(pop()))
		rp += 8
	}
	return uf(r)
}
func nul(x K) K { push(x); return 0 }
func lup(x K) K {
	vp := I32(8) + int32(x)
	return x0(K(vp))
}
func Asn(x, y K) K {
	if tp(x) != st {
		trap() //type
	}
	vp := I32(8) + int32(x)
	dx(K(I64(vp)))
	SetI64(vp, int64(rx(y)))
	return y
}
func Amd(x, i, v, y K) K {
	xt := tp(x)
	if xt == st {
		a := lup(x)
		if rc := I32(int32(a)-4); rc == 2 { //enable reuse for @[`x;i;+;y]
			dx(a)
			p := int32(a)
			a = rx(Amd(a, i, v, y))
			if int32(a) != p {
				SetI64(I32(8)+int32(x), int64((a)))
			}
			return a
		}
		return Asn(x, Amd(a, i, v, y))
		//return Asn(x, Amd(lup(x), i, v, y))
	}
	if xt < 16 {
		trap() //type
	}
	if tp(i) == Lt { // @[;;v;]/[x;y;i]
		n := nn(i)
		for j := int32(0); j < n; j++ {
			x = Amd(x, ati(rx(i), j), rx(v), ati(rx(y), j))
		}
		dx(i)
		dxy(v, y)
		return x
	}
	if xt > Lt {
		r := x0(x)
		x = r1(x)
		if xt == Tt && tp(i)&15 == it { // table-assign-rows
			if tp(y) > Lt {
				y = Val(y)
			}
			return key(r, Dmd(x, l2(0, i), v, y), xt)
		}
		r = Unq(Cat(r, rx(i)))
		return key(r, Amd(ntake(nn(r), x), Fnd(rx(r), i), v, y), xt)
	}
	if i == 0 {
		if v == 1 {
			if tp(y) < 16 {
				y = ntake(nn(x), y)
			}
			dx(x)
			return y
		}
		return Cal(v, l2(x, y))
	}
	if tp(v) != 0 || v != 1 {
		y = cal(v, l2(Atx(rx(x), rx(i)), y))
	}
	ti, yt := tp(i), tp(y)
	if xt&15 != yt&15 {
		x, xt = explode(x), Lt
	}
	if ti == it {
		if xt != yt+16 {
			x = explode(x)
		}
		return sti(use(x), int32(i), y)
	}
	if yt < 16 {
		y = ntake(nn(i), y)
		yt = tp(y)
	}
	if xt == Lt {
		y = explode(y)
	}
	return stv(x, i, y)
}
func Dmd(x, i, v, y K) K {
	if tp(x) == st {
		return Asn(x, Dmd(lup(x), i, v, y))
	}
	i = explode(i)
	f := Fst(rx(i))
	if nn(i) == 1 {
		dx(i)
		return Amd(x, f, v, y)
	}
	if f == 0 {
		f = seq(nn(x))
	}
	i = ndrop(1, i)
	if tp(f) > 16 { // matrix-assign
		n := nn(f)
		if nn(i) != 1 {
			trap() //rank
		}
		i = Fst(i)
		if tp(f) == It && tp(x) == Tt {
			t := rx(x0(x))
			return key(t, Dmd(r1(x), l2(Fnd(t, i), f), v, y), Tt)
		}
		if tp(f) != It || tp(x) != Lt {
			trap() // nyi Dt
		}
		x = use(x)
		for j := int32(0); j < n; j++ {
			rj := int32(x) + 8*I32(int32(f)+4*j)
			SetI64(rj, int64(Amd(K(I64(rj)), rx(i), rx(v), ati(rx(y), j))))
		}
		dxy(f, i)
		dxy(v, y)
		return x
	}
	return Amd(x, f, 1, Dmd(Atx(rx(x), f), i, v, y))
}
