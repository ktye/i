package k

import (
	. "github.com/ktye/wg/module"
)

type f1 = func(K) K
type f2 = func(K, K) K

func exec(x K) K {
	var a K
	pp = int32(x)
	pe = pp + 8*nn(x)
	for pp < pe {
		u := K(I64(pp))
		// fmt.Println("exec", tp(u), int32(u), sK(u))
		pp += 8
		if u > 2 {
			push(a)
			a = u
		} else {
			switch int32(u) {
			case 0:
				a = Func[marksrc(a)].(f1)(pop())
			case 1:
				a = Func[64+marksrc(a)].(f2)(pop(), pop())
			default:
				panic(Nyi)
			}
		}
	}
	return a
}
func marksrc(x K) int32 {
	srcp = 0xffffff & int32(x>>32)
	// fmt.Println("call func", int32(x))
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
