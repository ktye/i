package main

import "math/rand"

// shuffle x  (char, int, float, complex, symbol)
func rand1(x uint32) uint32 {
	xt, xn, _ := v1(x)
	if xt < 1 || xt > 5 {
		panic("shuffle: type")
	}
	r := use(x)
	p1 := int(r + 8)
	p2, p3, p4 := p1>>2, p1>>3, 1+p1>>3
	f1 := func(i, j int) { MC[p1+i], MC[p1+j] = MC[p1+j], MC[p1+i] }
	f4 := func(i, j int) { MI[p2+i], MI[p2+j] = MI[p2+j], MI[p2+i] }
	f8 := func(i, j int) { MF[p3+i], MF[p3+j] = MF[p3+j], MF[p3+i] }
	f16 := func(i, j int) {
		MF[p3+2*i], MF[p3+2*j] = MF[p3+2*j], MF[p3+2*i]
		MF[p4+2*i], MF[p4+2*j] = MF[p4+2*j], MF[p4+2*i]
	}
	switch xt {
	case 1:
		rand.Shuffle(int(xn), f1)
	case 2, 5:
		rand.Shuffle(int(xn), f4)
	case 3:
		rand.Shuffle(int(xn), f8)
	case 4:
		rand.Shuffle(int(xn), f16)
	default:
		panic("type")
	}
	return r
}

// 2 'r n   n random integers (int32)       randi n
//          n random integers 0..m-1        m/randi n
// 3 'r n   n uniform floats range 0..1     randf n
//-3 'r n   n normal distributed floats     randn n
// 4 'r n   n binormal complex numbers      randz n
func rand2(x, y uint32) uint32 {
	if tp(x) != 2 || nn(x) != 1 || tp(y) != 2 || nn(y) != 1 {
		panic("rand2: type (x and y must be int atoms)")
	}
	a := MI[2+x>>2]
	dx(x)
	b := iK(y)
	switch a {
	case 2:
		if b < 0 {
			panic("rand2: y<0")
		}
		r := mk(2, uint32(b))
		for i := uint32(0); i < uint32(b); i++ {
			MI[2+i+r>>2] = rand.Uint32()
		}
		return r
	case 4294967293: // -3
		r := make([]float64, b)
		for i := range r {
			r[i] = rand.NormFloat64()
		}
		return kF(r)
	case 3:
		r := make([]float64, b)
		for i := range r {
			r[i] = rand.Float64()
		}
		return kF(r)
	case 4:
		r := make([]complex128, b)
		for i := range r {
			r[i] = complex(rand.NormFloat64(), rand.NormFloat64())
		}
		return kZ(r)
	default:
		panic("rand2: x type")
	}
}
