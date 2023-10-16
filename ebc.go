package main

/*
import (
	. "github.com/ktye/wg/module"
)

// +266 bytes
func lz(x K) K { // long from zipped (lz4) one block: 11_-8_x
	r := mk(Ct, 0)
	p := int32(x)
	for {
		t := I8(p) //token
		p++
		l := t >> 4 //literal length
		if l == 15 {
			for {
				l += I8(p)
				p++
				if I8(p-1) != 255 {
					break
				}
			}
		}
		r = ucat(r, mk(Ct, l))
		Memorycopy(ep(r)-l, p, l) //literal
		p += l
		if p >= nn(x) {
			dx(x)
			return r
		}

		o := I8(p) | I8(p+1)<<8 //offset
		p += 2
		l = 4 + t&15 //match length
		if l == 19 {
			for {
				l += I8(p)
				p++
				if I8(p-1) != 255 {
					break
				}
			}
		}
		for l != 0 {
			l--
			r = cat1(r, Kc(I8(ep(r)-o)))
			continue
		}
	}
	return 0
}

func Des(x K) K { //deserialize
	pp := int32(x)
	r = des()
	dx(x)
	return r
}
func des() K { //deserialize
	t := I8(pp)
	n := I32(pp+1)
	pp += 1 + 4*I32B(t > 6) //all but atoms followed by len
	r := K(0)
	if t == 0 {
		r = K(I8(p))
	} else if t&15 < 7 {
		if t < 16 {
			t += 16
			r = mk(t, 1)
			n = sz(t)
			Memorycopy(int32(r), n)
			r = Fst(r)
		} else {
			n *= sz(t)
			Memorycopy(int32(r), n)
			r = mk(t, n)
		}
		pp += n
	} else {
		//10(comp) 11(derv) 12(proj) 13(lamb)
		//23(list) 24(dict) 25(table)
		r = mk(L, 0)
		for n != 0 {
			n--
			r = cat1(r, ser())
		}
		r = int32(r) | K(t)<<59
	}
	return r
}
*/
