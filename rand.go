package main

import (
	. "github.com/ktye/wg/module"
)

func rnd() (r int32) {
	r = rand
	r ^= (r << 13)
	r ^= (r >> 17)
	r ^= (r << 5)
	rand = r
	return r
}
func roll(x K) (r K) { // ?x (atom) ?n(uniform 0..1) ?-n(normal) ?z(binormal)
	xt := tp(x)
	xp := int32(x)
	if xt == it {
		if xp > 0 {
			return kx(240, x) // .rf uniform
		} else {
			r = kx(256, Ki((1+-xp)/2))
			SetI32(int32(r)-12, -xp)
			return K(int32(r)) | K(Ft)<<59 // normal
		}
	}
	if xt == zt {
		dx(x)
		return kx(256, Ki(int32(F64floor(F64(xp))))) //.rz binormal
	}
	return trap(Type)
}
func deal(x, y K) (r K) { // x?y (x atom) n?n(with replacement) -n?n(without) n?L (-#L)?L shuffle
	if tp(x) != it {
		trap(Type)
	}
	xp := int32(x)
	yt := tp(y)
	if yt > 16 {
		yn := nn(y)
		if xp < 0 {
			xp = -xp
			if xp >= yn {
				return atv(y, shuffle(seq(xp), yn)) //(-#L)?L
			} else {
				return atv(y, deal(x, Ki(yn))) // -n?L
			}
		} else {
			return atv(y, randI(nn(y), xp)) // n?L
		}
	}
	if yt != it {
		trap(Type)
	}
	yp := int32(y)
	if xp > 0 {
		return randI(yp, xp) // n?m
	}
	// todo n<<m
	return ntake(-xp, shuffle(seq(yp), -xp)) //-n?m (no duplicates)
}
func randi(n int32) int32 {
	v := uint32(rnd())
	prod := uint64(v) * uint64(n)
	low := uint32(prod)
	if low < uint32(n) {
		thresh := uint32(-n) % uint32(n)
		for low < thresh {
			v = uint32(rnd())
			prod = uint64(v) * uint64(n)
			low = uint32(prod)
		}
	}
	return int32(prod >> 32)
}
func randI(i, n int32) (r K) {
	r = mk(It, n)
	rp := int32(r)
	e := rp + 4*n
	if i == 0 {
		for rp < e {
			SetI32(rp, rnd())
			rp += 4
		}
	} else {
		for rp < e {
			SetI32(rp, randi(i))
			rp += 4
		}
	}
	return r
}
func shuffle(r K, m int32) K { // I, inplace
	rp := int32(r)
	n := nn(r)
	m = mini(n-1, m)
	for i := int32(0); i < m; i++ {
		ii := i + randi(n-i)
		j := rp + 4*(ii-i)
		t := I32(rp)
		SetI32(rp, I32(j))
		SetI32(j, t)
		rp += 4
	}
	return r
}
