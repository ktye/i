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
		pp += 8
		//fmt.Println("exec", sK(u))
		if tp(u) != 0 {
			push(u)
			continue
		}
		t := int32(u)

		if t < 80 {
			push(Func[t].(f1)(pop()))
		} else if t < 120 {
			a = pop()
			push(Func[t].(f2)(a, pop()))
		} else {
			push(u)
		}
	}
	return pop()
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
