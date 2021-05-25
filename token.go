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
	p := pp - 1 //force srcp > 0
	r = mk(Lt, 0)
	for {
		// todo ws
		if pp == pe {
			break
		}
		for i := int32(192); i < 196; i++ { // tnum, tvrb, tpct, tvar
			y = Func[i].(ftok)()
			if y != 0 {
				if i == 193 {
					y |= K(int64(pp-p) << 32)
				}
				//fmt.Println("mark", (int64(pp-p) << 32), pp-p, y)
				r = cat1(r, y)
				break
			}
			if i == 195 { // todo last-1
				trap(Parse)
			}
		}
	}
	return r
}
func tnms() (r K) {
	r = tnum()
	for pp < pe-1 && I8(pp) == ' ' {
		pp++
		x := tnum()
		if x == 0 {
			break
		}
		r = ncat(r, x)
	}
	return r
}
func tnum() K {
	p := pp
	c := I8(p)
	if p > int32(src) {
		if c == '-' || c == '.' {
			if is(I8(p-1), 64) {
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
	o := int32(1)
	if pp < pe && I8(pp) == 58 { // :
		pp++
		if is(c, 8) {
			o = 2 // ':
		} else {
			o = 97 // +:
		}
	}
	return K(o + index(c, 228, 253))
}
func tpct() (r K) {
	c := I8(pp)
	if is(c, 48) { // ([{}]);
		pp++
		return K(c)
	}
	return 0
}
func tvar() (r K) {
	c := I8(pp)
	if !is(c, 2) {
		return 0
	}
	pp++
	r = Ku(uint64(c))
	for pp < pe {
		c = I8(pp)
		if !is(c, 6) {
			break
		}
		r = cat1(r, K(c)|K(ct)<<59)
		pp++
	}
	return sc(r)
}
func is(x, m int32) bool { return m&I8(100+x) != 0 }
