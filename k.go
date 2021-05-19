package k

import . "github.com/ktye/wg/module"

func kinit() {
	minit(10, 16)
}

type K uint64
type T uint32

// type: t=x>>57
// size: 1<<clz(x)
// atom: 1&x>>56
// flat: t>tt
// func: 0=t
// pntr: p=uint32(x)
// rc:   I32(p-8)
// len:  I32(p-4)
const ( //                   bytes  atom      vector
	bt T = 65 // bool    1      10000011  10000010
	ct T = 69 // char    1      10001011  10001010
	it T = 17 // int     4      00100011  00100010
	ft T = 15 // float   8      00011111  00011110
	st T = 25 // symbol  4      00110011  00110010
	zt T = 12 // complex (8)    00011001  00011000
	lt T = 8  // list    8                00010000
	dt T = 10 // dict    (8)              00010100
	tt T = 14 // table   (8)              00010100
)

// func t=0
// basic x < 64 (triadic/tetradic)
// xn=2: composition string funclist
// xn=3: derived     string func symb
// xn=4: projection  string func arglist emptylist
// xn=5: lambda      string code locals arity save

// ptr: int32(x)
// p-8       p-4     p
// [refcount][length][data]

func Kb(x bool) K  { return K(ib(x)) | (K(131) << 56) }
func Kc(x int32) K { return K(x) | (K(139) << 56) }
func Ki(x int32) K { return K(x) | (K(35) << 56) }
func Kf(x float64) (r K) {
	r = mk(ft, 1)
	SetF64(int32(r), x)
	return r
}
func Kz(r, i float64) (z K) {
	x := Kf(r)
	var y K
	if i != 0 {
		y = Kf(i)
	}
	z = mk(zt, 1)
	SetI64(int32(z), int64(x))
	SetI64(8+int32(z), int64(y))
	return z
}
func l2(x, y K) (r K) {
	r = mk(lt, 2)
	SetI64(int32(r), int64(x))
	SetI64(8+int32(r), int64(y))
	return r
}
func mkd(x, y K) (r K) {
	nx, ny := cnt(x), cnt(y)
	if nx+ny == int32(-2) {
		x, y = enl(x), enl(y)
	}
	r = mk(dt, 2)
	SetI64(int32(r), int64(x))
	SetI64(8+int32(r), int64(y))
	return r
}
func cnt(x K) int32 { // #x but -1 for atoms
	if x < 64 || (1&(x>>56) != 0) {
		return -1
	}
	t := tp(x)
	if t == dt {
		return cnt(K(I64(int32(x))))
	}
	if t == tt {
		return cnt(first(rx(K(I64(8 + int32(x))))))
	}
	return I32(int32(x) - 4)
}
func ib(x bool) int32 {
	if x {
		return 1
	}
	return 0
}
