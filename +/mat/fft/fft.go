package main

import (
	"math"
	"math/bits"
)

const N = 1 << 15

func main() {
	println(N)
	x := make([]complex128, N)
	f := Prepare(N)
	for i := 0; i < 1000; i++ {
		R(x)
		f.Complex(x)
	}
	P(x)
}

func R(x []complex128) {
	for i := range x {
		x[i] = complex(rnd(), rnd())
	}
}
func P(x []complex128) {
	for i := 0; i < 3; i++ {
		print(real(x[i]))
		print(" ")
		print(imag(x[i]))
		print(" ")
	}
	println()
}

var rand_ int32 = 1592653589

func rnd() float64 {
	r := rand_
	r ^= (r << 13)
	r ^= (r >> 17)
	r ^= (r << 5)
	rand_ = r
	return 0.5 + float64(r)/4294967295.0
}


type NFFT struct {
	n, h, l uint16
	p       []uint16
	e       [][]complex128
	i       [][]uint16
}

func Prepare(n uint16) NFFT { //n: power of two
	l := uint16(bits.TrailingZeros16(n))
	e := make([][]complex128, l)
	i := make([][]uint16, l)
	r := rots(n)
	p := perm(n)
	s := n
	h := n >> 1
	t := uint16(1)
	for k := range e {
		E := make([]complex128, h)
		I := make([]uint16, h)
		s >>= 1
		c := 0
		for b := uint16(0); b < s; b++ {
			o := 2 * b * t
			for j := uint16(0); j < t; j++ {
				I[c] = j + o
				E[c] = r[s*j]
				c++
			}
		}
		e[k] = E
		i[k] = I
		t <<= 1
	}
	return NFFT{n: n, h: h, l: l, p: p, e: e, i: i}
}
func (f NFFT) Complex(x []complex128) {
	brswap(x, f.p)
	s := uint16(1)
	for i, el := range f.e {
		l := f.i[i]
		for k := uint16(0); k < f.h; k++ {
			ii := l[k]
			jj := ii + s
			xi := x[ii]
			xj := x[jj]
			ek := el[k]
			x[ii] += xj * ek
			x[jj] = xi - xj*ek
		}
		s <<= 1
	}
}
func brswap(x []complex128, p []uint16) {
	for i := range p {
		if k := p[i]; i < int(k) {
			x[i], x[k] = x[k], x[i]
		}
	}
}
func perm(n uint16) []uint16 {
	r := make([]uint16, n)
	k := uint16(1)
	for n > 1 {
		n >>= 1
		for i := uint16(0); i < k; i++ {
			r[i] <<= 1
			r[i+k] = 1 + r[i]
		}
		k <<= 1
	}
	return r
}
func rots(N uint16) []complex128 {
	E := make([]complex128, N)
	for n := uint16(0); n < N; n++ {
		phi := -2.0 * math.Pi * float64(n) / float64(N)
		s, c := math.Sincos(phi)
		E[n] = complex(c, s)
	}
	return E
}
