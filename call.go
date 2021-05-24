package k

import (
	"fmt"

	. "github.com/ktye/wg/module"
)

func Cal(x, y K) K {
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
			rl(y)
			dx(y)
			yp := int32(y)
			return cal2(x, K(I64(yp)), K(8+I64(yp)))
		}
	}
	panic(Nyi)
	return x
}

func cal1(f, x K) (r K) {
	t := tp(f)
	if t == 0 {
		return Func[int32(t)].(f1)(x)
	} else if t == df {
		fp := int32(f)
		d := K(I64(fp))
		a := int32(I64(fp + 8))
		r = Func[85+int32(a)].(f2)(d, x)
		dx(f)
		return r
	}
	trap(Nyi)
	return x
}
func cal2(f, x, y K) K {
	trap(Nyi)
	return x
}
