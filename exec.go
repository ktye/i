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
		// fmt.Printf("exec p=%d tp=%d int32=%d case(%d) %s\n", p, tp(u), int32(u), int32(u)>>6, sK(u))
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
				p += int32(a) * ib(int32(u) == 0)
				a = pop()
				dx(u)
			default: //448..     quoted verb
				push(a)
				a = rx(u - 448)
			}
		}
		p += 8
		continue
	}
	return a
}

func marksrc(x K) int32 {
	srcp = 0xffffff & int32(x>>32)
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
	return r
}
func nul(x K) K { push(x); return 0 }
func lup(x K) K {
	if tp(x) != st {
		trap(Type)
	}
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
	if tp(v) != 0 || int32(v) != 1 {
		y = cal(v, l2(Atx(rx(x), rx(i)), y))
	}
	if tp(x) == It && tp(i) == it && tp(y) == it {
		r = ucat(x, mk(It, 0))
		SetI32(int32(r)+4*int32(i), int32(y))
		return r
	}
	trap(Nyi)
	return x
}
func Dmd(x, i, v, y K) K {
	//fmt.Printf("dmend[%s;%s;%s;%s]\n", sK(x), sK(i), sK(v), sK(y))
	trap(Nyi)
	return y
}
