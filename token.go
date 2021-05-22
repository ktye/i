package k

import (
	. "github.com/ktye/wg/module"
)

type ftok = func() K

func tok(x K) (r K) {
	var y K
	src = x
	pp = int32(x)
	pe = pp + nn(x)
	r = mk(Lt, 0)
	for {
		// todo ws
		if pp == pe {
			break
		}
		p := pp
		for i := int32(120); i < 122; i++ { // tnum, tvrb
			y = Func[i].(ftok)()
			if y != 0 {
				y |= K(p << 32)
				r = lcat(r, y)
				break
			}
			if i == 121 { // todo last-1
				trap(Parse)
			}
		}
	}
	return r
}

func tnum() K {
	p := pp
	c := I8(p)
	if p > int32(src) {
		if c == '-' || c == '.' {
			if is(I8(p-1), 32) {
				return 0 // e.g. x-1 is (x - 1) not (x -1)
			}
		}
	}
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
		if is(c, 4) == false {
			break
		}
		r = 10*r + int64(c-'0')
		pp++
	}
	return r
}

func tvrb() (r K) {
	c := I8(pp)
	if !is(c, 1) {
		return 0
	}
	pp++
	return K(1 + index(c, 228, 250))
}
func is(x, m int32) bool { return m&I8(100+x) != 0 }
