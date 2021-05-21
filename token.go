package k

import (
	. "github.com/ktye/wg/module"
)

type ftok = func() K

func tok(x K) (r K) {
	var y K
	src = int32(x)
	pp = int32(x)
	pe = pp + nn(x)
	r = mk(Lt, 0)
	for {
		// todo ws
		if pp == pe {
			break
		}
		p := pp
		for i := int32(0); i < 2; i++ { // tnum, tvrb
			y = Func[i].(ftok)()
			if y != 0 {
				y |= K(p << 32)
				r = lcat(r, y)
				break
			}
			if i == 1 { // todo ntoks-1
				trap(Parse)
			}
		}
	}
	return r
}

func tnum() K {
	p := pp
	r := pi()
	if r == 0 && p == pp {
		return 0
	}
	return Ki(int32(r))
}
func pi() (r int64) {
	if I8(pp) == '-' {
		pp++
		p := pp
		r = -pu()
		if r == 0 {
			pp = p - 1
		}
	} else {
		r = pu()
	}
	return r
}
func pu() (r int64) {
	for pp < pe {
		c := I8(pp)
		i := c - '0'
		if i < 0 || i > 9 {
			break
		}
		r = 10*r + int64(i)
		pp++
	}
	return r
}

func tvrb() (r K) {
	i := fndc(cvb, I8(pp))
	if i < 0 {
		return 0
	}
	pp++
	return K(1 + i)
}
