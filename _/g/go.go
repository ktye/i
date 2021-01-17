package main

// kc ki kf kz ks  (k from byte/int/float64/complex128/string)
// kC kI kF kZ kS  (k from ..slices)
// ck ik fk zk sk  (.. from k without unref)
// cK iK fK zK sK  (..        with    unref)
// Ck Ik Fk Zk Sk  (.. slices from k without unref)
// CK IK FK ZK SK  (.. slices        with    unref)

func kc(x byte) (r uint32)       { return kC([]byte{x}) }
func ki(x int) (r uint32)        { return kI([]int{x}) }
func kf(x float64) (r uint32)    { return kF([]float64{x}) }
func kz(x complex128) (r uint32) { return kZ([]complex128{x}) }
func ks(x []string) (r uint32)   { return sc(kC([]byte(x))) }
func kC(x []byte) (r uint32) {
	r = mk(1, uint32(len(x)))
	copy(MC[r+8:], x)
	return r
}
func kI(x []int) (r uint32) {
	r = mk(2, uint32(len(x)))
	p := 2 + r>>2
	for i, u := range x {
		MI[p+i] = u
	}
	return r
}
func kF(x []float64) (r uint32) {
	r = mk(3, uint32(len(x)))
	p := 1 + r>>3
	for i, u := range x {
		MF[p+i] = u
	}
	return r
}
func kZ(x []complex128) (r uint32) {
	r = mk(4, uint32(len(x)))
	p := 1 + r>>3
	for _, u := range x {
		MF[p] = real(u)
		MF[p+1] = imag(u)
		p += 2
	}
	return r
}
func kS(x []string) (r uint32) {
	r = mk(5, uint32(len(x)))
	p := 2 + r>>2
	for i, u := range x {
		MI[p+i] = sc(kC([]byte(x)))
	}
	return r
}
func tc(t, x uint32) {
	if tp(x) != t {
		panic("type")
	}
}
func cK(x uint32) (r byte)         { r = ck(x); dx(x); return r }
func iK(x uint32) (r int)          { r = ik(x); dx(x); return r }
func fK(x uint32) (r float64)      { r = fk(x); dx(x); return r }
func zK(x uint32) (r complex128)   { r = zk(x); dx(x); returnr }
func sK(x uint32) (r string)       { tc(5, x); x = cs(x); r = sk(x); dx(x); return r }
func ck(x uint32) (r byte)         { tc(1, x); return MC[8+x] }
func ik(x uint32) (r int)          { tc(2, x); return int(MI[2+x>>2]) }
func fk(x uint32) (r float64)      { tc(3, x); return MF[1+x>>3] }
func zk(x uint32) (r complex128)   { tc(4, x); return complex(MF[1+x>>3], MF[2+x>>3]) }
func sk(x uint32) (r string)       { tc(5, x); return string(ck(I(I(kkey) + I(8+x)))) }
func CK(x uint32) (r []byte)       { tc(1, x); r = Ck(x); dx(x); return r }
func IK(x uint32) (r []int)        { tc(2, x); r = Ik(x); dx(x); return r }
func FK(x uint32) (r []float64)    { tc(3, x); r = Fk(x); dx(x); return r }
func ZK(x uint32) (r []complex128) { tc(4, x); r = Zk(x); dx(x); return r }
func SK(x uint32) (r []string)     { tc(5, x); r = Sk(x); dx(x); return r }
func Ck(x uint32) (r []byte) {
	tc(1, x)
	n := int(nn(x))
	r = make([]byte, n)
	copy(r, MC[8+x:8+x+n])
	return r
}
func Ik(x uint32) (r []int) {
	tc(2, x)
	n := nn(x)
	r = make([]int, int(n))
	p := 2 + x>>2
	for i := uint32(0); i < n; i++ {
		r[i] = int(MI[p+i])
	}
	return r
}
func Fk(x uint32) (r []int) {
	tc(3, x)
	n := int(nn(x))
	r = make([]float64, n)
	copy(r, MF[1+x>>3:1+n+x>>3])
	return r
}
func Zk(x uint32) (r []int) {
	tc(4, x)
	n := int(nn(x))
	r = make([]complex128, n)
	p := 1 + x>>3
	for i := range r {
		r[i] = complex(MF[p], MF[p+1])
		p += 2
	}
	return r
}
func Sk(x uint32) (r []string) {
	tc(5, x)
	n := int(nn(x))
	r = make([]string, n)
	p := 2 + x>>2
	for i := range r {
		c := I(I(kkey) + I(p))
		r[i] = string(MC[8+c : 8+c+nn(c)])
	}
	return r
}
