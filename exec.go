package main

import (
	. "github.com/ktye/wg/module"
)

type f1 = func(K) K
type f2 = func(K, K) K
type f3 = func(K, K, K) K
type f4 = func(K, K, K, K) K

func quoted(x K) bool { return int32(x) >= 448 && tp(x) == 0 }
func quote(x K) K     { return x + 448 }
func unquote(x K) K   { return x - 448 }

func exec(x K) K {
	xn := nn(x)
	if xn == 0 {
		dx(x)
		return 0
	}
	var a K // accumulator
	p := int32(x)
	e := p + 8*xn
	for p < e {
		u := K(I64(p))
		//fmt.Printf("exec p=%d tp=%d int32=%d case(%d) %s\n", p, tp(u), int32(u), int32(u)>>6, sK(u))
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
				a = Cal(a, l2(pop(), pop()))
			case 3: // 192..255  tetradic
				a = Func[marksrc(u)].(f4)(a, pop(), pop(), pop())
			case 4: // 256       drop
				dx(a)
				a = pop()
			case 5: // 320       jump
				p += int32(a)
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
	dx(pop())
	dx(x)
	return a
}

func marksrc(x K) int32 {
	if p := 0xffffff & int32(x>>32); p != 0 {
		srcp = p
	}
	return int32(x)
}
func push(x K) {
	SetI64(sp, int64(x))
	sp += 8
	if sp == 512 {
		trap(Stack)
	}
}
func pop() (r K) {
	sp -= 8
	if sp < 256 {
		trap(Stack)
	}
	return K(I64(sp))
}
func lst(n K) (r K) {
	rn := int32(n)
	r = mk(Lt, rn)
	rp := int32(r)
	for i := int32(0); i < rn; i++ {
		SetI64(rp, int64(pop()))
		rp += 8
	}
	return uf(r)
}
func nul(x K) K { push(x); return 0 }
func lup(x K) K {
	vp := I32(8) + int32(x)
	r := x0(vp)
	return r
}
func Asn(x, y K) K {
	if tp(x) != st {
		trap(Type)
	}
	vp := I32(8) + int32(x)
	dx(K(I64(vp)))
	SetI64(vp, int64(rx(y)))
	return y
}
func Amd(x, i, v, y K) (r K) {
	//fmt.Printf("amend[%s;%s;%s;%s]\n", sK(x), sK(i), sK(v), sK(y))
	xt := tp(x)
	if xt < 16 {
		trap(Type)
	}
	if xt > Lt {
		if xt == Tt && tp(i)&15 == it {
			return trap(Nyi) // table-assign-rows
		}
		r, x = spl2(x)
		return key(r, Amd(x, Fnd(rx(r), i), v, y), xt)
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
		return sti(x, int32(i), y)
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
func Dmd(x, i, v, y K) (r K) {
	// fmt.Printf("dmend[%s;%s;%s;%s]\n", sK(x), sK(i), sK(v), sK(y))
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
			trap(Rank)
		}
		i = Fst(i)
		if tp(f) != It || tp(x) != Lt { // todo table-matrix-assign
			trap(Nyi) // Dt
		}
		r = use(x)
		for j := int32(0); j < n; j++ {
			rj := int32(r) + 8*I32(int32(f)+4*j)
			SetI64(rj, int64(Amd(K(I64(rj)), rx(i), rx(v), ati(rx(y), j))))
		}
		dx(f)
		dx(i)
		dx(v)
		dx(y)
		return r
	}

	return Amd(x, f, 1, Dmd(Atx(rx(x), f), i, v, y))
}
