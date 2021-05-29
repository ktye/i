package main

import (
	. "github.com/ktye/wg/module"
)

type f1 = func(K) K
type f2 = func(K, K) K
type f3 = func(K, K, K) K
type f4 = func(K, K, K, K) K

func exec(x K) K {
	//fmt.Println("exec", sK(x))
	var a K
	p := int32(x)
	e := p + 8*nn(x)
	for p < e {
		u := K(I64(p))
		//fmt.Println("exec", tp(u), int32(u), sK(u), u > 2)
		p += 8
		if u > 7 {
			push(a)
			a = u
		} else {
			switch int32(u) {
			case 0:
				a = Lup(a)
			case 1:
				a = Func[marksrc(a)].(f1)(pop())
			case 2:
				a = Func[64+marksrc(a)].(f2)(pop(), pop())
			case 3:
				a = Func[128+marksrc(a)].(f3)(pop(), pop(), pop())
			case 4:
				a = Func[192+marksrc(a)].(f4)(pop(), pop(), pop(), pop())
			case 5: // drop
				dx(a)
				a = pop()
			case 6: // jump
				p += int32(a)

				a = pop()
			case 7: // jump if not
				u = pop()
				p += int32(a) * ib(int32(u) == 0)
				a = pop()
				dx(u)
			default:
				panic(Nyi)
			}
		}
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

func Lup(x K) K {
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
