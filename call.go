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
	//case 2: // df
	//case 3: // pf
	//case 4: // lf
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
func compose(x, y K, yt T) (r K) {
	if yt == ct {
		r = cat1(K(int32(y))|K(Lt)<<59, x)
	} else {
		r = l2(y, x)
	}
	return K(int32(r)) | K(cf)<<59
}
